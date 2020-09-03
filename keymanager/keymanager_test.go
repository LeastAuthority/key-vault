package keymanager_test

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	basekeymanager "github.com/prysmaticlabs/prysm/validator/keymanager/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/keymanager"
	"github.com/bloxapp/key-vault/utils/endpoint"
)

const (
	defaultAccountPublicKey = "965586b5d05c851873f26cb736ed42de96591674772576e7b43848cd7a5c2827a5c5228034fdd55be0e9dc0f0cbc91d7"
	defaultAccessToken      = "supersecureaccesstoken"
)

func TestSignGeneric(t *testing.T) {
	accountPubKey, err := hex.DecodeString(defaultAccountPublicKey)
	require.NoError(t, err)

	domain := make([]byte, 32)
	rand.Read(domain)

	data := make([]byte, 32)
	rand.Read(data)

	beaconState, privKeys := testutil.DeterministicGenesisState(t, 100)
	block, err := testutil.GenerateFullBlock(beaconState, privKeys, nil, 0)
	require.NoError(t, err)

	expectedSignature, err := bls.SignatureFromBytes(block.GetSignature())
	require.NoError(t, err)

	actualSignature := hex.EncodeToString(block.GetSignature())

	var protect sync.Mutex
	var currentMethod http.HandlerFunc
	s := newTestRemoteWallet(t, func(writer http.ResponseWriter, request *http.Request) {
		currentMethod(writer, request)
	})
	defer s.Close()

	wallet, _, err := keymanager.NewKeyManager(logrus.NewEntry(logrus.New()), &keymanager.Config{
		Location:    s.URL,
		AccessToken: defaultAccessToken,
		PubKey:      defaultAccountPublicKey,
	})
	require.NoError(t, err)

	runTest := func(t *testing.T, statusCode int, signature []byte, f func(wallet *keymanager.KeyManager)) {
		protect.Lock()
		currentMethod = func(writer http.ResponseWriter, request *http.Request) {
			require.Equal(t, http.MethodPost, request.Method)
			require.Equal(t, endpoint.Build("accounts/sign-aggregation"), request.URL.Path)

			var req keymanager.SignAggregationRequest
			require.NoError(t, json.NewDecoder(request.Body).Decode(&req))

			require.Equal(t, defaultAccountPublicKey, req.PubKey)
			require.Equal(t, hex.EncodeToString(domain), req.Domain)
			require.Equal(t, hex.EncodeToString(data), req.DataToSign)

			if statusCode == http.StatusOK {
				respBody := &logical.Response{
					Data: map[string]interface{}{
						"signature": string(signature),
					},
				}
				require.NoError(t, json.NewEncoder(writer).Encode(respBody))
			} else {
				writer.WriteHeader(statusCode)
			}
		}
		protect.Unlock()

		f(wallet)
	}

	t.Run("successfully signed data", func(t *testing.T) {
		runTest(t, http.StatusOK, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignGeneric(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(data), bytesutil.ToBytes32(domain))
			require.NoError(t, err)
			require.NotNil(t, actualSignature)
			require.Equal(t, expectedSignature, actualSignature)
		})
	})

	t.Run("rejects with undefined account", func(t *testing.T) {
		undefinedAccount := make([]byte, 48)
		rand.Read(undefinedAccount)

		runTest(t, http.StatusOK, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignGeneric(bytesutil.ToBytes48(undefinedAccount), bytesutil.ToBytes32(data), bytesutil.ToBytes32(domain))
			require.Error(t, err, basekeymanager.ErrNoSuchKey.Error())
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with denied", func(t *testing.T) {
		runTest(t, http.StatusUnauthorized, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignGeneric(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(data), bytesutil.ToBytes32(domain))
			require.Error(t, err, basekeymanager.ErrDenied.Error())
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with failed", func(t *testing.T) {
		runTest(t, http.StatusInternalServerError, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignGeneric(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(data), bytesutil.ToBytes32(domain))
			require.Error(t, err, basekeymanager.ErrCannotSign.Error())
			require.Nil(t, actualSignature)
		})
	})
}

func TestSignProposal(t *testing.T) {
	accountPubKey, err := hex.DecodeString(defaultAccountPublicKey)
	require.NoError(t, err)

	domain := make([]byte, 32)
	rand.Read(domain)
	data := &ethpb.BeaconBlockHeader{
		Slot:          10,
		ProposerIndex: 10,
		StateRoot:     []byte{1, 2, 3},
		ParentRoot:    []byte{1, 2, 3},
		BodyRoot:      []byte{1, 2, 3},
	}

	beaconState, privKeys := testutil.DeterministicGenesisState(t, 100)
	block, err := testutil.GenerateFullBlock(beaconState, privKeys, nil, 0)
	require.NoError(t, err)

	expectedSignature, err := bls.SignatureFromBytes(block.GetSignature())
	require.NoError(t, err)

	actualSignature := hex.EncodeToString(block.GetSignature())

	var protect sync.Mutex
	var currentMethod http.HandlerFunc
	s := newTestRemoteWallet(t, func(writer http.ResponseWriter, request *http.Request) {
		currentMethod(writer, request)
	})
	defer s.Close()

	wallet, _, err := keymanager.NewKeyManager(logrus.NewEntry(logrus.New()), &keymanager.Config{
		Location:    s.URL,
		AccessToken: defaultAccessToken,
		PubKey:      defaultAccountPublicKey,
	})
	require.NoError(t, err)

	runTest := func(t *testing.T, statusCode int, signature []byte, f func(wallet *keymanager.KeyManager)) {
		protect.Lock()
		currentMethod = func(writer http.ResponseWriter, request *http.Request) {
			require.Equal(t, http.MethodPost, request.Method)
			require.Equal(t, endpoint.Build("accounts/sign-proposal"), request.URL.Path)

			var req keymanager.SignProposalRequest
			require.NoError(t, json.NewDecoder(request.Body).Decode(&req))

			require.Equal(t, defaultAccountPublicKey, req.PubKey)
			require.Equal(t, hex.EncodeToString(domain), req.Domain)
			require.Equal(t, int(data.GetSlot()), int(req.Slot))
			require.Equal(t, int(data.GetProposerIndex()), int(req.ProposerIndex))
			require.Equal(t, hex.EncodeToString(data.GetParentRoot()), req.ParentRoot)
			require.Equal(t, hex.EncodeToString(data.GetStateRoot()), req.StateRoot)
			require.Equal(t, hex.EncodeToString(data.GetBodyRoot()), req.BodyRoot)

			if statusCode == http.StatusOK {
				respBody := &logical.Response{
					Data: map[string]interface{}{
						"signature": string(signature),
					},
				}
				require.NoError(t, json.NewEncoder(writer).Encode(respBody))
			} else {
				writer.WriteHeader(statusCode)
			}
		}
		protect.Unlock()

		f(wallet)
	}

	t.Run("successfully signed data", func(t *testing.T) {
		runTest(t, http.StatusOK, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignProposal(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(domain), data)
			require.NoError(t, err)
			require.NotNil(t, actualSignature)
			require.Equal(t, expectedSignature, actualSignature)
		})
	})

	t.Run("rejects with undefined account", func(t *testing.T) {
		undefinedAccount := make([]byte, 48)
		rand.Read(undefinedAccount)

		runTest(t, http.StatusOK, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignProposal(bytesutil.ToBytes48(undefinedAccount), bytesutil.ToBytes32(domain), data)
			require.Error(t, err, basekeymanager.ErrNoSuchKey.Error())
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with denied", func(t *testing.T) {
		runTest(t, http.StatusUnauthorized, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignProposal(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(domain), data)
			require.Error(t, err, basekeymanager.ErrDenied.Error())
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with failed", func(t *testing.T) {
		runTest(t, http.StatusInternalServerError, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignProposal(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(domain), data)
			require.Error(t, err, basekeymanager.ErrCannotSign.Error())
			require.Nil(t, actualSignature)
		})
	})
}

func TestSignAttestation(t *testing.T) {
	accountPubKey, err := hex.DecodeString(defaultAccountPublicKey)
	require.NoError(t, err)

	domain := make([]byte, 32)
	rand.Read(domain)
	data := &ethpb.AttestationData{
		Slot:            10,
		CommitteeIndex:  10,
		BeaconBlockRoot: []byte{1, 2, 3},
		Source: &ethpb.Checkpoint{
			Epoch: 101010,
			Root:  []byte{2, 3, 4},
		},
		Target: &ethpb.Checkpoint{
			Epoch: 202020,
			Root:  []byte{5, 6, 7},
		},
	}

	beaconState, privKeys := testutil.DeterministicGenesisState(t, 100)
	block, err := testutil.GenerateFullBlock(beaconState, privKeys, nil, 0)
	require.NoError(t, err)

	expectedSignature, err := bls.SignatureFromBytes(block.GetSignature())
	require.NoError(t, err)

	actualSignature := hex.EncodeToString(block.GetSignature())

	var protect sync.Mutex
	var currentMethod http.HandlerFunc
	s := newTestRemoteWallet(t, func(writer http.ResponseWriter, request *http.Request) {
		currentMethod(writer, request)
	})
	defer s.Close()

	wallet, _, err := keymanager.NewKeyManager(logrus.NewEntry(logrus.New()), &keymanager.Config{
		Location:    s.URL,
		AccessToken: defaultAccessToken,
		PubKey:      defaultAccountPublicKey,
	})
	require.NoError(t, err)

	runTest := func(t *testing.T, statusCode int, signature []byte, f func(wallet *keymanager.KeyManager)) {
		protect.Lock()
		currentMethod = func(writer http.ResponseWriter, request *http.Request) {
			require.Equal(t, http.MethodPost, request.Method)
			require.Equal(t, endpoint.Build("accounts/sign-attestation"), request.URL.Path)

			var req keymanager.SignAttestationRequest
			require.NoError(t, json.NewDecoder(request.Body).Decode(&req))

			require.Equal(t, defaultAccountPublicKey, req.PubKey)
			require.Equal(t, hex.EncodeToString(domain), req.Domain)
			require.Equal(t, int(data.GetSlot()), int(req.Slot))
			require.Equal(t, int(data.GetCommitteeIndex()), int(req.CommitteeIndex))
			require.Equal(t, hex.EncodeToString(data.GetBeaconBlockRoot()), req.BeaconBlockRoot)
			require.Equal(t, int(data.GetSource().GetEpoch()), int(req.SourceEpoch))
			require.Equal(t, hex.EncodeToString(data.GetSource().GetRoot()), req.SourceRoot)
			require.Equal(t, int(data.GetTarget().GetEpoch()), int(req.TargetEpoch))
			require.Equal(t, hex.EncodeToString(data.GetTarget().GetRoot()), req.TargetRoot)

			if statusCode == http.StatusOK {
				respBody := &logical.Response{
					Data: map[string]interface{}{
						"signature": string(signature),
					},
				}
				require.NoError(t, json.NewEncoder(writer).Encode(respBody))
			} else {
				writer.WriteHeader(statusCode)
			}
		}
		protect.Unlock()

		f(wallet)
	}

	t.Run("successfully signed data", func(t *testing.T) {
		runTest(t, http.StatusOK, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(domain), data)
			require.NoError(t, err)
			require.NotNil(t, actualSignature)
			require.Equal(t, expectedSignature, actualSignature)
		})
	})

	t.Run("rejects with undefined account", func(t *testing.T) {
		undefinedAccount := make([]byte, 48)
		rand.Read(undefinedAccount)

		runTest(t, http.StatusOK, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(bytesutil.ToBytes48(undefinedAccount), bytesutil.ToBytes32(domain), data)
			require.Error(t, err, basekeymanager.ErrNoSuchKey.Error())
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with denied", func(t *testing.T) {
		runTest(t, http.StatusUnauthorized, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(domain), data)
			require.Error(t, err, basekeymanager.ErrDenied.Error())
			require.Nil(t, actualSignature)
		})
	})

	t.Run("rejects with failed", func(t *testing.T) {
		runTest(t, http.StatusInternalServerError, []byte(actualSignature), func(wallet *keymanager.KeyManager) {
			actualSignature, err := wallet.SignAttestation(bytesutil.ToBytes48(accountPubKey), bytesutil.ToBytes32(domain), data)
			require.Error(t, err, basekeymanager.ErrCannotSign.Error())
			require.Nil(t, actualSignature)
		})
	})
}

func newTestRemoteWallet(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler(writer, request)
	}))

	return s
}
