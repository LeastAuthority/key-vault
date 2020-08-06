package store_test

import (
	"context"
	"testing"

	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/bloxapp/vault-plugin-secrets-eth2.0/backend/store"
)

func getStorage() logical.Storage {
	return &logical.InmemStorage{}
}

func getWalletStorage() core.Storage {
	return store.NewHashicorpVaultStore(getStorage(), context.Background())
}

func TestOpeningAccounts(t *testing.T) {
	stores.TestingOpenAccounts(getWalletStorage(), t)
}

func TestNonExistingWallet(t *testing.T) {
	stores.TestingNonExistingWallet(getWalletStorage(), t)
}

func TestWalletStorage(t *testing.T) {
	stores.TestingWalletStorage(getWalletStorage(), t)
}
