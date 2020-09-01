package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	vault "github.com/bloxapp/eth-key-manager"
	"github.com/bloxapp/eth-key-manager/core"
	"github.com/bloxapp/eth-key-manager/stores/in_memory"
	"github.com/bloxapp/eth-key-manager/wallet_hd"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"

	"github.com/bloxapp/key-vault/backend/store"
)

// Endpoints patterns
const (
	// StoragePattern is the path pattern for storage endpoint
	StoragePattern = "storage"

	// SlashingStoragePattern is the path pattern for slashing storage endpoint
	SlashingStoragePattern = "storage/%s/slashing"
)

// SlashingHistory contains slashing history data.
type SlashingHistory struct {
	Attestations []*core.BeaconAttestation `json:"attestations"`
	Proposals    []*core.BeaconBlockHeader `json:"proposals"`
}

func storagePaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         StoragePattern,
			HelpSynopsis:    "Update storage",
			HelpDescription: `Manage KeyVault storage`,
			Fields: map[string]*framework.FieldSchema{
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "storage to update",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathStorageUpdate,
			},
		},
		&framework.Path{
			Pattern:         fmt.Sprintf(SlashingStoragePattern, framework.GenericNameRegex("public_key")),
			HelpSynopsis:    "Update slashing storage",
			HelpDescription: `Manage KeyVault slashing storage`,
			Fields: map[string]*framework.FieldSchema{
				"public_key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Public key of the account",
					Default:     "",
				},
				"data": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "slashing storage to update",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSlashingStorageUpdate,
				logical.ReadOperation:   b.pathSlashingStorageRead,
			},
		},
	}
}

func (b *backend) pathStorageUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	storage := data.Get("data").(string)
	storageBytes, err := hex.DecodeString(storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode storage")
	}

	var inMemStore *in_memory.InMemStore
	err = json.Unmarshal(storageBytes, &inMemStore)
	if err != nil {
		return nil, errors.Wrap(err, "failed to JSON un-marshal storage")
	}

	_, err = store.FromInMemoryStore(ctx, inMemStore, req.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update storage")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}

func (b *backend) pathSlashingStorageUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)

	// Get request payload
	publicKey := data.Get("public_key").(string)
	slashingData := data.Get("data").(string)

	// Open wallet
	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet")
	}

	account, err := wallet.AccountByPublicKey(publicKey)
	if err != nil {
		if err == wallet_hd.ErrAccountNotFound {
			return b.notFoundResponse()
		}

		return nil, errors.Wrap(err, "failed to retrieve account")
	}

	// HEX decode slashing history
	slashingHistoryBytes, err := hex.DecodeString(slashingData)
	if err != nil {
		return b.badRequestResponse(errors.Wrap(err, "failed to HEX decode slashing storage").Error())
	}

	// JSON unmarshal slashing history
	var slashingHistory SlashingHistory
	if err := json.Unmarshal(slashingHistoryBytes, &slashingHistory); err != nil {
		return b.badRequestResponse(errors.Wrap(err, "failed to unmarshal slashing history").Error())
	}

	// Store attestation history
	for _, attestation := range slashingHistory.Attestations {
		// Save attestation
		if err := storage.SaveAttestation(account.ValidatorPublicKey(), attestation); err != nil {
			return nil, errors.Wrapf(err, "failed to save attestation for slot %d", attestation.Slot)
		}
	}

	// Store proposal history
	for _, proposal := range slashingHistory.Proposals {
		// Save proposals
		if err := storage.SaveProposal(account.ValidatorPublicKey(), proposal); err != nil {
			return nil, errors.Wrapf(err, "failed to save proposal for slot %d", proposal.Slot)
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}

func (b *backend) pathSlashingStorageRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)

	// Get request payload
	publicKey := data.Get("public_key").(string)

	// Open wallet
	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet")
	}

	account, err := wallet.AccountByPublicKey(publicKey)
	if err != nil {
		if err == wallet_hd.ErrAccountNotFound {
			return b.notFoundResponse()
		}

		return nil, errors.Wrap(err, "failed to retrieve account")
	}

	var slashingHistory SlashingHistory

	// Fetch attestations
	if slashingHistory.Attestations, err = storage.ListAllAttestations(account.ValidatorPublicKey()); err != nil {
		return nil, errors.Wrap(err, "failed to list attestations data")
	}

	// Fetch proposals
	if slashingHistory.Proposals, err = storage.ListAllProposals(account.ValidatorPublicKey()); err != nil {
		return nil, errors.Wrap(err, "failed to list proposals data")
	}

	slashingHistoryEncoded, err := json.Marshal(slashingHistory)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal slashing history")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"data": hex.EncodeToString(slashingHistoryEncoded),
		},
	}, nil
}
