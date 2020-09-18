package tests

import (
	"encoding/hex"
	"testing"

	"github.com/bloxapp/eth2-key-manager/core"

	"github.com/bloxapp/eth2-key-manager/slashing_protection"
	"github.com/bloxapp/eth2-key-manager/stores/in_memory"
	"github.com/bloxapp/eth2-key-manager/validator_signer"
	"github.com/stretchr/testify/require"
	v1 "github.com/wealdtech/eth2-signer-api/pb/v1"

	"github.com/bloxapp/key-vault/e2e"
	"github.com/bloxapp/key-vault/e2e/shared"
)

// AttestationSigning tests sign attestation endpoint.
type AttestationSigning struct {
}

// Name returns the name of the test.
func (test *AttestationSigning) Name() string {
	return "Test attestation signing"
}

// Run run the test.
func (test *AttestationSigning) Run(t *testing.T) {
	setup := e2e.Setup(t)

	// setup vault with db
	setup.UpdateConfig(t)
	storage := setup.UpdateStorage(t)
	account := shared.RetrieveAccount(t, storage)
	require.NotNil(t, account)
	pubKeyBytes := account.ValidatorPublicKey().Marshal()

	// Get wallet
	wallet, err := storage.OpenWallet()
	require.NoError(t, err)

	dataToSign := map[string]interface{}{
		"public_key":      hex.EncodeToString(pubKeyBytes),
		"domain":          "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
		"slot":            284115,
		"committeeIndex":  2,
		"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
		"sourceEpoch":     8877,
		"sourceRoot":      "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
		"targetEpoch":     8878,
		"targetRoot":      "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
	}

	// Sign data
	protector := slashing_protection.NewNormalProtection(in_memory.NewInMemStore(core.MainNetwork))
	var signer validator_signer.ValidatorSigner = validator_signer.NewSimpleSigner(wallet, protector)

	res, err := signer.SignBeaconAttestation(test.dataToAttestationRequest(t, pubKeyBytes, dataToSign))
	require.NoError(t, err)

	// Send sign attestation request
	sig, err := setup.SignAttestation(dataToSign)
	require.NoError(t, err)

	require.Equal(t, res.GetSignature(), sig)
}

func (test *AttestationSigning) dataToAttestationRequest(t *testing.T, pubKey []byte, data map[string]interface{}) *v1.SignBeaconAttestationRequest {
	// Decode domain
	domainBytes, err := hex.DecodeString(data["domain"].(string))
	require.NoError(t, err)

	// Decode block root
	beaconBlockRoot, err := hex.DecodeString(data["beaconBlockRoot"].(string))
	require.NoError(t, err)

	// Decode source root
	sourceRootBytes, err := hex.DecodeString(data["sourceRoot"].(string))
	require.NoError(t, err)

	// Decode target root
	targetRootBytes, err := hex.DecodeString(data["targetRoot"].(string))
	require.NoError(t, err)

	return &v1.SignBeaconAttestationRequest{
		Id:     &v1.SignBeaconAttestationRequest_PublicKey{PublicKey: pubKey},
		Domain: domainBytes,
		Data: &v1.AttestationData{
			Slot:            uint64(data["slot"].(int)),
			CommitteeIndex:  uint64(data["committeeIndex"].(int)),
			BeaconBlockRoot: beaconBlockRoot,
			Source: &v1.Checkpoint{
				Epoch: uint64(data["sourceEpoch"].(int)),
				Root:  sourceRootBytes,
			},
			Target: &v1.Checkpoint{
				Epoch: uint64(data["targetEpoch"].(int)),
				Root:  targetRootBytes,
			},
		},
	}
}
