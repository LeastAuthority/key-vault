package backend

import (
	"context"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"sort"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func getBackend(t *testing.T) (logical.Backend, logical.Storage) {
	config := &logical.BackendConfig{
		Logger:      logging.NewVaultLogger(log.Trace),
		System:      &logical.StaticSystemView{},
		StorageView: &logical.InmemStorage{},
		BackendUUID: "test",
	}

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("unable to create backend: %v", err)
	}

	// Wait for the upgrade to finish
	time.Sleep(time.Second)

	return b, config.StorageView
}

func TestAccountsList(t *testing.T) {
	b, _ := getBackend(t)


	t.Run("Successfully List Accounts", func(t *testing.T) {
		req := logical.TestRequest(t, logical.ListOperation, "wallet/accounts/")

		// setup logical storage
		_,err := baseHashicorpStorage(req.Storage, context.Background())
		require.NoError(t, err)


		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data["accounts"])
		require.Len(t, res.Data["accounts"], 1)
		require.Equal(t, res.Data["accounts"].([]map[string]string)[0]["name"], "test_account")

		// make sure only the following fields are present to prevent accidental secret sharing
		keys := make([]string, 0)
		for k := range res.Data["accounts"].([]map[string]string)[0] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		require.Equal(t, keys, []string{"id","name","validationPubKey","withdrawalPubKey"})
	})
}