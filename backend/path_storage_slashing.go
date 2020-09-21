package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"

	vault "github.com/bloxapp/eth2-key-manager"
	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/bloxapp/eth2-key-manager/wallet_hd"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"

	"github.com/bloxapp/key-vault/backend/store"
	"github.com/bloxapp/key-vault/utils/errorex"
)

// Endpoints patterns
const (
	// SlashingStoragePattern is the path pattern for slashing storage endpoint
	SlashingStoragePattern = "storage/slashing"
)

// SlashingHistory contains slashing history data.
type SlashingHistory struct {
	Attestations []*core.BeaconAttestation `json:"attestations"`
	Proposals    []*core.BeaconBlockHeader `json:"proposals"`
}

func storageSlashingPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         SlashingStoragePattern,
			HelpSynopsis:    "Manage slashing storage",
			HelpDescription: `Manage KeyVault slashing storage`,
			ExistenceCheck:  b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSlashingStorageBatchUpdate,
				logical.ReadOperation:   b.pathSlashingStorageBatchRead,
			},
		},
	}
}

func (b *backend) pathSlashingStorageBatchUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Load config
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, err
	}

	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage, config.Network)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)

	// Open wallet
	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet")
	}

	// Load accounts slashing history
	for publicKey, data := range req.Data {
		account, err := wallet.AccountByPublicKey(publicKey)
		if err != nil {
			if err == wallet_hd.ErrAccountNotFound {
				return b.notFoundResponse()
			}

			return nil, errors.Wrap(err, "failed to retrieve account")
		}

		// Store slashing data
		if err := storeAccountSlashingHistory(storage, account, data.(string)); err != nil {
			return b.prepareErrorResponse(err)
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}

func (b *backend) pathSlashingStorageBatchRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Load config
	config, err := b.configured(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage, config.Network)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)

	// Open wallet
	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet")
	}

	// Load accounts slashing history
	responseData := make(map[string]interface{})
	for _, account := range wallet.Accounts() {
		// Load slashing history
		slashingHistory, err := loadAccountSlashingHistory(storage, account)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load slashing history")
		}

		responseData[hex.EncodeToString(account.ValidatorPublicKey().Marshal())] = slashingHistory
	}

	return &logical.Response{
		Data: responseData,
	}, nil
}

func loadAccountSlashingHistory(storage *store.HashicorpVaultStore, account core.ValidatorAccount) (string, error) {
	var slashingHistory SlashingHistory
	var err error

	// Fetch attestations
	if slashingHistory.Attestations, err = storage.ListAllAttestations(account.ValidatorPublicKey()); err != nil {
		return "", errors.Wrap(err, "failed to list attestations data")
	}

	// Fetch proposals
	if slashingHistory.Proposals, err = storage.ListAllProposals(account.ValidatorPublicKey()); err != nil {
		return "", errors.Wrap(err, "failed to list proposals data")
	}

	slashingHistoryEncoded, err := json.Marshal(slashingHistory)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal slashing history")
	}

	return hex.EncodeToString(slashingHistoryEncoded), nil
}

func storeAccountSlashingHistory(storage *store.HashicorpVaultStore, account core.ValidatorAccount, slashingData string) error {
	// HEX decode slashing history
	slashingHistoryBytes, err := hex.DecodeString(slashingData)
	if err != nil {
		return errorex.NewErrBadRequest(err.Error())
	}

	// JSON unmarshal slashing history
	var slashingHistory SlashingHistory
	if err := json.Unmarshal(slashingHistoryBytes, &slashingHistory); err != nil {
		return errorex.NewErrBadRequest(err.Error())
	}

	// Store attestation history
	for _, attestation := range slashingHistory.Attestations {
		// Save attestation
		if err := storage.SaveAttestation(account.ValidatorPublicKey(), attestation); err != nil {
			return errors.Wrapf(err, "failed to save attestation for slot %d", attestation.Slot)
		}
	}

	// Store proposal history
	for _, proposal := range slashingHistory.Proposals {
		// Save proposals
		if err := storage.SaveProposal(account.ValidatorPublicKey(), proposal); err != nil {
			return errors.Wrapf(err, "failed to save proposal for slot %d", proposal.Slot)
		}
	}

	return nil
}
