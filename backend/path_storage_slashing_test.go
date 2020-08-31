package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/bloxapp/eth-key-manager/core"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/backend/store"
)

func TestSlashingStorage(t *testing.T) {
	b, _ := getBackend(t)
	inMemStore, accountID, err := baseInmemStorage()
	require.NoError(t, err)

	account, err := inMemStore.OpenAccount(accountID)
	require.NoError(t, err)
	publicKey := hex.EncodeToString(account.ValidatorPublicKey().Marshal())

	t.Run("successfully setup slashing history", func(t *testing.T) {
		slashingHistory, err := json.Marshal(struct {
			Attestations []*core.BeaconAttestation `json:"attestations"`
			Proposals    []*core.BeaconBlockHeader `json:"proposals"`
		}{
			Attestations: []*core.BeaconAttestation{
				{
					Slot:            123123,
					BeaconBlockRoot: []byte{1, 2, 3},
					CommitteeIndex:  1,
					Source: &core.Checkpoint{
						Root:  []byte{1, 2, 3},
						Epoch: 123,
					},
					Target: &core.Checkpoint{
						Root:  []byte{1, 2, 3},
						Epoch: 123,
					},
				},
			},
			Proposals: []*core.BeaconBlockHeader{
				{
					Slot:          123123,
					ProposerIndex: 1,
					BodyRoot:      []byte{1, 2, 3},
					ParentRoot:    []byte{1, 2, 3},
					StateRoot:     []byte{1, 2, 3},
				},
			},
		})
		require.NoError(t, err)

		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		req.Data = map[string]interface{}{
			"public_key": publicKey,
			"data":       hex.EncodeToString(slashingHistory),
		}
		_, err = store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.True(t, res.Data["status"].(bool))
	})

	t.Run("rejects setup slashing history with wrong JSON format", func(t *testing.T) {
		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		req.Data = map[string]interface{}{
			"public_key": publicKey,
			"data":       hex.EncodeToString([]byte("slashinghistory")),
		}
		_, err := store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.EqualValues(t, 400, res.Data["http_status_code"])
	})

	t.Run("rejects setup slashing history with wrong HEX format", func(t *testing.T) {
		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		req.Data = map[string]interface{}{
			"public_key": publicKey,
			"data":       "slashinghistory",
		}
		_, err := store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.EqualValues(t, 400, res.Data["http_status_code"])
	})

	t.Run("rejects setup slashing history for unknown public key", func(t *testing.T) {
		slashingHistory, err := json.Marshal(struct {
			Attestations []*core.BeaconAttestation `json:"attestations"`
			Proposals    []*core.BeaconBlockHeader `json:"proposals"`
		}{
			Attestations: []*core.BeaconAttestation{},
			Proposals:    []*core.BeaconBlockHeader{},
		})
		require.NoError(t, err)

		ctx := context.Background()
		req := logical.TestRequest(t, logical.CreateOperation, "storage/slashing")
		req.Data = map[string]interface{}{
			"public_key": "ab0cb36c4ce5ddabdc38a1d6868c871328539ebde5fea89686b2cd6332bf4cc5f9c48a501d1d6d87bf916d0e1b01ead963e1b6ce52075e26dc65bad535ecfad0",
			"data":       hex.EncodeToString(slashingHistory),
		}
		_, err = store.FromInMemoryStore(ctx, inMemStore, req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(ctx, req)
		require.NoError(t, err)
		require.EqualValues(t, 404, res.Data["http_status_code"])
	})
}
