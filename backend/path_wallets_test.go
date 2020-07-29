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
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/")
		storage = req.Storage
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data["wallet"])
	})

	t.Run("Create Wallet when one already exists", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "KeyVault wallet already exists")
	})
}