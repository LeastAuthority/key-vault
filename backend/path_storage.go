package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault"

	"github.com/pkg/errors"

	"github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func storagePaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "storage",
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
	}
}

func (b *backend) pathStorageUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	KeyVault.InitCrypto()

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

	_, err = hashicorp.FromInMemoryStore(inMemStore, req.Storage, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update storage")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}
