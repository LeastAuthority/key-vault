package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func TestPushUpdate(t *testing.T) {
	b, _ := getBackend(t)
	store := in_memory.NewInMemStore()
	var logicalStorage logical.Storage

	// wallet
	wallet := wallet_hd.NewHDWallet(&core.WalletContext{Storage:store})
	err := store.SaveWallet(wallet)
	require.NoError(t, err)

	// account
	acc,err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), "acc")
	require.NoError(t, err)
	err = store.SaveAccount(acc)
	require.NoError(t, err)

	// marshal and to string
	byts, err := json.Marshal(store)
	require.NoError(t, err)
	data := hex.EncodeToString(byts)

	// test
	t.Run("import from in-memory to hashicorp vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.UpdateOperation, "admin/pushUpdate")
		logicalStorage = req.Storage
		req.Data = map[string]interface{}{
			"data": data,
		}
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.True(t, res.Data["status"].(bool))
	})

	t.Run("verify wallet and account", func(t *testing.T) {
		vault := hashicorp.NewHashicorpVaultStore(logicalStorage, context.Background())
		wallet2,err := vault.OpenWallet()
		require.NoError(t, err)
		require.Equal(t, wallet.ID().String(), wallet2.ID().String())

		acc2, err := wallet2.AccountByName("acc")
		require.NoError(t, err)
		require.Equal(t, acc.ID().String(), acc2.ID().String())
	})
}
