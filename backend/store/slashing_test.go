package store_test

import (
	"context"
	"testing"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/stores"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/bloxapp/key-vault/backend/store"
)

func getSlashingStorage() core.SlashingStore {
	return store.NewHashicorpVaultStore(context.Background(), &logical.InmemStorage{})
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
