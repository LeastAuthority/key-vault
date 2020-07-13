package backend

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
)

func getBackend(t *testing.T) (logical.Backend, logical.Storage) {
	config := &logical.BackendConfig{
		Logger:      logging.NewVaultLogger(log.Trace),
		System:      &logical.StaticSystemView{},
		StorageView: &logical.InmemStorage{},
		BackendUUID: "test",
	}

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("unable to create backend: %v", err)
	}

	// Wait for the upgrade to finish
	time.Sleep(time.Second)

	return b, config.StorageView
}

func TestWalletCreate(t *testing.T) {
	b, _ := getBackend(t)
	var storage logical.Storage

	t.Run("Successfully Create Wallet", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
		storage = req.Storage
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data["wallet"])
	})

	t.Run("Create Wallet with empty name", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/ ")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "unsupported path")
	})

	t.Run("Create Wallet with existing name", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to create new wallet: wallet \"wallet1\" already exists")
	})

	// TODO
	//t.Run("Create Wallet with too long name (more than 128 characters)", func(t *testing.T) {
	//
	//})

	//data = map[string]interface{}{
	//	"walletName":  "wallet1",
	//	"accountName": "account1",
	//}
	//req.Data = data
}

func TestWalletsList(t *testing.T) {
	b, _ := getBackend(t)
	var storage logical.Storage

	t.Run("Successfully List Wallets", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
		storage = req.Storage
		_, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		req = logical.TestRequest(t, logical.ListOperation, "wallets/")
		req.Storage = storage
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data["wallets"])
	})

	t.Run("List Wallets in non existing portfolio", func(t *testing.T) {
		req := logical.TestRequest(t, logical.ListOperation, "wallets/")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to open key vault: key vault not found")
	})
}
