package backend

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"

	"github.com/bloxapp/key-vault/backend/store"
	"github.com/bloxapp/key-vault/utils/errorex"
)

// Factory returns the backend factory
func Factory(version string) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b := newBackend(version)
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// newBackend returns the backend
func newBackend(version string) *backend {
	b := &backend{
		Version: version,
	}
	b.Backend = &framework.Backend{
		Help: "",
		Paths: prepareStoreMiddleware(framework.PathAppend(
			versionPaths(b),
			storagePaths(b),
			storageSlashingPaths(b),
			accountsPaths(b),
			signsPaths(b),
		)),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"wallet/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	return b
}

// backend implements the Backend for this plugin
type backend struct {
	*framework.Backend
	Version string
}

func (b *backend) pathExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		b.Logger().Error("Path existence check failed", err)
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}

func (b *backend) notFoundResponse() (*logical.Response, error) {
	return logical.RespondWithStatusCode(&logical.Response{
		Data: map[string]interface{}{
			"message":     "account not found",
			"status_code": http.StatusNotFound,
		},
	}, nil, http.StatusNotFound)
}

func (b *backend) prepareErrorResponse(originError error) (*logical.Response, error) {
	switch err := errors.Cause(originError).(type) {
	case *errorex.ErrBadRequest:
		return err.ToLogicalResponse()
	case nil:
		return nil, nil
	default:
		return logical.ErrorResponse(originError.Error()), nil
	}
}

func prepareStoreMiddleware(paths []*framework.Path) []*framework.Path {
	wrapperFunc := func(callback framework.OperationFunc) framework.OperationFunc {
		return framework.OperationFunc(func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
			network := data.Get("network").(string)
			switch network {
			case "test":
				req.Storage = store.NewLogical(req.Storage, "test")
				break
			case "launchtest":
				req.Storage = store.NewLogical(req.Storage, "launchtest")
				break
			default:
				return nil, errors.New("unsupported network")
			}
			return callback(ctx, req, data)
		})
	}

	for i := range paths {
		paths[i].Pattern = framework.GenericNameRegex("network") + "/" + paths[i].Pattern
		if paths[i].Fields == nil {
			paths[i].Fields = make(map[string]*framework.FieldSchema)
		}
		paths[i].Fields["network"] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: "Blockchain network",
			Default:     "",
		}
		for j, callback := range paths[i].Callbacks {
			paths[i].Callbacks[j] = wrapperFunc(callback)
		}
	}
	return paths
}
