package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/bloxapp/eth-key-manager/core"
	"github.com/bloxapp/eth-key-manager/stores/in_memory"
	"github.com/bloxapp/eth-key-manager/wallet_hd"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/backend/store"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func baseInmemStorage() (*in_memory.InMemStore, error) {
	inMemStore := in_memory.NewInMemStore()

	// wallet
	wallet := wallet_hd.NewHDWallet(&core.WalletContext{Storage: inMemStore})
	err := inMemStore.SaveWallet(wallet)
	if err != nil {
		return nil, err
	}

	// account
	acc, err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), "test_account")
	if err != nil {
		return nil, err
	}

	if err := inMemStore.SaveAccount(acc); err != nil {
		return nil, err
	}

	return inMemStore, nil
}

func baseHashicorpStorage(logicalStorage logical.Storage, ctx context.Context) (*store.HashicorpVaultStore, error) {
	inMem, err := baseInmemStorage()
	if err != nil {
		return nil, err
	}
	return store.FromInMemoryStore(ctx, inMem, logicalStorage)
}

func TestStorage(t *testing.T) {
	b, _ := getBackend(t)
	inMemStore, err := baseInmemStorage()
	require.NoError(t, err)
	var logicalStorage logical.Storage

	// marshal and to string
	byts, err := json.Marshal(inMemStore)
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
		wallet, err := inMemStore.OpenWallet()
		require.NoError(t, err)
		acc, err := wallet.AccountByPublicKey("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279")
		require.NoError(t, err)

		vault := store.NewHashicorpVaultStore(context.Background(), logicalStorage)
		wallet2, err := vault.OpenWallet()
		require.NoError(t, err)
		require.Equal(t, wallet.ID().String(), wallet2.ID().String())

		acc2, err := wallet2.AccountByPublicKey("ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279")
		require.NoError(t, err)
		require.Equal(t, acc.ID().String(), acc2.ID().String())
	})
}
