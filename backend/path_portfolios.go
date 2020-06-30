package backend

import (
	"context"
	"encoding/hex"

	vault "github.com/bloxapp/KeyVault"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
)

func portfoliosPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "export",
			HelpSynopsis:    "Export seed",
			HelpDescription: `Export seed`,
			ExistenceCheck:  b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathWalletExport,
			},
		},
		&framework.Path{
			Pattern:         "import",
			HelpSynopsis:    "Import seed",
			HelpDescription: `Import seed`,
			Fields: map[string]*framework.FieldSchema{
				"seed": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletImport,
			},
		},
	}
}

func (b *backend) pathWalletExport(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)

	if _, err := vault.OpenKeyVault(&options); err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	seed, err := storage.SecurelyFetchPortfolioSeed()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch portfolio seed")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"seed": hex.EncodeToString(seed),
		},
	}, nil
}

func (b *backend) pathWalletImport(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	seed := data.Get("seed").(string)
	seedDecoded, err := hex.DecodeString(seed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hex decode the given seed")
	}

	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)
	options.SetSeed(seedDecoded)

	if _, err := vault.ImportKeyVault(&options); err != nil {
		return nil, errors.Wrap(err, "failed to import key vault")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"success": true,
		},
	}, nil
}
