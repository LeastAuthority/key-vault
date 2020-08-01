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
	types "github.com/wealdtech/go-eth2-types/v2"
	"testing"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func baseInmemStorage() (*in_memory.InMemStore, error) {
	types.InitBLS()

	store := in_memory.NewInMemStore()

	// wallet
	wallet := wallet_hd.NewHDWallet(&core.WalletContext{Storage:store})
	err := store.SaveWallet(wallet)
	if err != nil {
		return nil, err
	}

	// account
	acc,err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), "test_account")
	if err != nil {
		return nil, err
	}
	err = store.SaveAccount(acc)
	if err != nil {
		return nil, err
	}

	return store, nil
}


func baseHashicorpStorage(logicalStorage logical.Storage, ctx context.Context) (*hashicorp.HashicorpVaultStore, error) {
	inMem, err := baseInmemStorage()
	if err != nil {
		return nil, err
	}
	return hashicorp.FromInMemoryStore(inMem, logicalStorage, ctx)
}

func TestPushUpdate(t *testing.T) {
	require.NoError(t, types.InitBLS())

	b, _ := getBackend(t)
	store, err := baseInmemStorage()
	require.NoError(t, err)
	var logicalStorage logical.Storage

	// marshal and to string
	byts, err := json.Marshal(store)
	require.NoError(t, err)
	data := hex.EncodeToString(byts)

	// test
	t.Run("import from in-memory to hashicorp vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "storage")
		logicalStorage = req.Storage
		req.Data = map[string]interface{}{
			"data": data,
		}
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.True(t, res.Data["status"].(bool))
	})

	t.Run("verify wallet and account", func(t *testing.T) {
		// get wallet and account
		wallet, err := store.OpenWallet()
		require.NoError(t, err)
		acc, err := wallet.AccountByName("test_account")
		require.NoError(t, err)

		vault := hashicorp.NewHashicorpVaultStore(logicalStorage, context.Background())
		wallet2,err := vault.OpenWallet()
		require.NoError(t, err)
		require.Equal(t, wallet.ID().String(), wallet2.ID().String())

		acc2, err := wallet2.AccountByName("test_account")
		require.NoError(t, err)
		require.Equal(t, acc.ID().String(), acc2.ID().String())
	})
}
