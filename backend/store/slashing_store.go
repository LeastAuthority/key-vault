package store

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bloxapp/eth-key-manager/core"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// Paths
const (
	WalletAttestationsBase      = "attestations/%s/"
	WalletAttestationPath       = WalletAttestationsBase + "%d"     // account/attestation
	WalletLatestAttestationPath = WalletAttestationsBase + "latest" // account/latest
	WalletProposalsBase         = "proposals/%s/"                   // account/proposal
	WalletProposalsPath         = WalletProposalsBase + "%d"        // account/proposal
)

// SaveAttestation implements Storage interface.
func (store *HashicorpVaultStore) SaveAttestation(key e2types.PublicKey, req *core.BeaconAttestation) error {
	path := fmt.Sprintf(WalletAttestationPath, store.identfierFromKey(key), req.Target.Epoch)
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal attestation request")
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// RetrieveAttestation implements Storage imterface.
func (store *HashicorpVaultStore) RetrieveAttestation(key e2types.PublicKey, epoch uint64) (*core.BeaconAttestation, error) {
	path := fmt.Sprintf(WalletAttestationPath, store.identfierFromKey(key), epoch)
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record from storage with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	var ret *core.BeaconAttestation
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal beacon attestation object")
	}

	return ret, nil
}

// ListAttestations both epochStart and epochEnd reflect saved attestations by their target epoch
func (store *HashicorpVaultStore) ListAttestations(key e2types.PublicKey, epochStart uint64, epochEnd uint64) ([]*core.BeaconAttestation, error) {
	ret := make([]*core.BeaconAttestation, 0)

	for epoch := epochStart; epoch <= epochEnd; epoch++ {
		att, err := store.RetrieveAttestation(key, epoch)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve attestation with epoch %d", epoch)
		}

		if att != nil {
			ret = append(ret, att)
		}
	}

	return ret, nil
}

// ListAllAttestations returns all attestation data from the DB
func (store *HashicorpVaultStore) ListAllAttestations(key e2types.PublicKey) ([]*core.BeaconAttestation, error) {
	path := fmt.Sprintf(WalletAttestationsBase, store.identfierFromKey(key))
	entries, err := store.storage.List(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list records from storage with path '%s'", path)
	}

	// Return nothing if there is no record
	if len(entries) == 0 {
		return nil, nil
	}

	attestations := make([]*core.BeaconAttestation, 0)
	for _, entry := range entries {
		epoch, err := strconv.Atoi(entry)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid epoch number %s", entry)
		}

		ret, err := store.RetrieveAttestation(key, uint64(epoch))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve beacon attestation for epoch %d", epoch)
		}

		attestations = append(attestations, ret)
	}

	return attestations, nil
}

// SaveProposal implements Storage imterface.
func (store *HashicorpVaultStore) SaveProposal(key e2types.PublicKey, req *core.BeaconBlockHeader) error {
	path := fmt.Sprintf(WalletProposalsPath, store.identfierFromKey(key), req.Slot)
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal proposal request")
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// RetrieveProposal implements Storage interface.
func (store *HashicorpVaultStore) RetrieveProposal(key e2types.PublicKey, slot uint64) (*core.BeaconBlockHeader, error) {
	path := fmt.Sprintf(WalletProposalsPath, store.identfierFromKey(key), slot)
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	var ret *core.BeaconBlockHeader
	if err = json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal beacon block header object")
	}

	return ret, nil
}

// ListAllProposals returns all proposal data from the DB
func (store *HashicorpVaultStore) ListAllProposals(key e2types.PublicKey) ([]*core.BeaconBlockHeader, error) {
	path := fmt.Sprintf(WalletProposalsBase, store.identfierFromKey(key))
	entries, err := store.storage.List(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list records from storage with path '%s'", path)
	}

	// Return nothing if there is no record
	if len(entries) == 0 {
		return nil, nil
	}

	proposals := make([]*core.BeaconBlockHeader, 0)
	for _, entry := range entries {
		epoch, err := strconv.Atoi(entry)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid epoch number %s", entry)
		}

		ret, err := store.RetrieveProposal(key, uint64(epoch))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve beacon proposal for epoch %d", epoch)
		}

		proposals = append(proposals, ret)
	}

	return proposals, nil
}

// SaveLatestAttestation implements Storage interface.
func (store *HashicorpVaultStore) SaveLatestAttestation(key e2types.PublicKey, req *core.BeaconAttestation) error {
	path := fmt.Sprintf(WalletLatestAttestationPath, store.identfierFromKey(key))
	data, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "failed to marshal beacon attestation object")
	}

	entry := &logical.StorageEntry{
		Key:      path,
		Value:    data,
		SealWrap: false,
	}
	return store.storage.Put(store.ctx, entry)
}

// RetrieveLatestAttestation implements Storage imterface.
func (store *HashicorpVaultStore) RetrieveLatestAttestation(key e2types.PublicKey) (*core.BeaconAttestation, error) {
	path := fmt.Sprintf(WalletLatestAttestationPath, store.identfierFromKey(key))
	entry, err := store.storage.Get(store.ctx, path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get record with path '%s'", path)
	}

	// Return nothing if there is no record
	if entry == nil {
		return nil, nil
	}

	var ret *core.BeaconAttestation
	if err := json.Unmarshal(entry.Value, &ret); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal beacon attestation object")
	}

	return ret, nil
}

func (store *HashicorpVaultStore) identfierFromKey(key e2types.PublicKey) string {
	return hex.EncodeToString(key.Marshal())
}
