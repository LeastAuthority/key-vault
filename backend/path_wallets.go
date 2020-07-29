package backend

import (
	"context"

	vault "github.com/bloxapp/KeyVault"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
)

func walletsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "wallet/",
			HelpSynopsis:    "Create an Ethereum 2.0 wallet",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{

			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletCreate,
			},
		},
	}
}

func (b *backend) pathWalletCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := &vault.KeyVaultOptions{}
	options.SetStorage(storage)

	// check a wallet doesn't exist
	existingKv, err := vault.OpenKeyVault(options)
	if err == nil || existingKv != nil {
		return nil, errors.New("KeyVault wallet already exists")
	}


	// create new KeyVault (which creates new wallet)
	kv, err := vault.NewKeyVault(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wallet")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"wallet": wallet,
		},
	}, nil
}
