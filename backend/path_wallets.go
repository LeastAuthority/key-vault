package backend

import (
	"context"

	vault "github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
)

func walletsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "wallets/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathWalletsList,
			},
			HelpSynopsis: "List all the Ethereum 2.0 wallets at a path",
			HelpDescription: `
			All the Ethereum 2.0 wallets will be listed.
			`,
		},
		&framework.Path{
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name"),
			HelpSynopsis:    "Create an Ethereum 2.0 wallet",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletCreate,
			},
		},
	}
}

func (b *backend) pathWalletCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)
	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		switch e := err.(type) {
		// TODO move to portfolio create path
		case *vault.NotExistError:
			portfolio, err = vault.NewKeyVault(&options)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create key vault")
			}
		default:
			return nil, errors.Wrap(e, "failed to open key vault")
		}
	}

	wallet, err := portfolio.CreateWallet(walletName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new wallet")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"wallet": wallet,
		},
	}, nil
}

func (b *backend) pathWalletsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)

	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	var wallets []core.Wallet
	for w := range portfolio.Wallets() {
		wallets = append(wallets, w)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"wallets": wallets,
		},
	}, nil
}
