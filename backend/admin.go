package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	types "github.com/wealdtech/go-eth2-types/v2"
)

func adminPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "admin/pushUpdate",
			HelpSynopsis:    "Push and replace KeyVault data",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"data": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathPushUpdate,
			},
		},
	}
}

func (b *backend) pathPushUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := types.InitBLS() // very important
	if err != nil {
		return nil,err
	}

	keyVaultData := data.Get("data").(string)
	byts, err := hex.DecodeString(keyVaultData)
	if err != nil {
		return nil,err
	}

	var inMemStore *in_memory.InMemStore
	err = json.Unmarshal(byts, &inMemStore)
	if err != nil {
		return nil,err
	}

	_, err = hashicorp.FromInMemoryStore(inMemStore, req.Storage, ctx)
	if err != nil {
		return nil,err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"status": true,
		},
	}, nil
}
