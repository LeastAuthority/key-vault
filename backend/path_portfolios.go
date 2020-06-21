package backend

import (
	"context"
	"encoding/hex"

	vault "github.com/bloxapp/KeyVault"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func portfoliosPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "export",
			HelpSynopsis:    "Export seed",
			HelpDescription: `Export seed`,
			ExistenceCheck:  b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletExport,
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
	_, err := vault.OpenKeyVault(&options)
	//_, err := vault.NewKeyVault(&options)
	if err != nil {
		return nil, err
	}

	seed, _ := storage.SecurelyFetchPortfolioSeed()

	return &logical.Response{
		Data: map[string]interface{}{
			"seed": hex.EncodeToString(seed),
		},
	}, nil
}

func (b *backend) pathWalletImport(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	seed := data.Get("seed").(string)
	seedDecoded, _ := hex.DecodeString(seed)
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)
	options.SetSeed(seedDecoded)
	//_, err := vault.NewKeyVault(&options)
	_, err := vault.ImportKeyVault(&options)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"success": true,
		},
	}, nil
}
