package store_test

import (
	"context"
	"testing"

	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/bloxapp/vault-plugin-secrets-eth2.0/backend/store"
)

func getSlashingStorage() core.SlashingStore {
	return store.NewHashicorpVaultStore(&logical.InmemStorage{}, context.Background())
}

func TestSavingProposal(t *testing.T) {
	stores.TestingSaveProposal(getSlashingStorage(), t)
}

func TestSavingAttestation(t *testing.T) {
	stores.TestingSaveAttestation(getSlashingStorage(), t)
}

func TestSavingLatestAttestation(t *testing.T) {
	stores.TestingSaveLatestAttestation(getSlashingStorage(), t)
}

func TestRetrieveEmptyLatestAttestation(t *testing.T) {
	stores.TestingRetrieveEmptyLatestAttestation(getSlashingStorage(), t)
}

func TestListingAttestation(t *testing.T) {
	stores.TestingListingAttestation(getSlashingStorage(), t)
}
