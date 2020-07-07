package backend

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory returns the backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend()
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns the backend
func Backend() (*backend, error) {
	var b backend
	b.Backend = &framework.Backend{
		Help: "",
		Paths: framework.PathAppend(
			portfoliosPaths(&b),
			walletsPaths(&b),
			accountsPaths(&b),
			signsPaths(&b),
		),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"wallets/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}
	return &b, nil
}

// backend implements the Backend for this plugin
type backend struct {
	*framework.Backend
}

func (b *backend) pathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		b.Logger().Error("Path existence check failed", err)
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}
