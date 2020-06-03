package backend

import (
	"context"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"testing"
	"time"
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

func TestWalletCreate(t *testing.T) {
	b, _ := getBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "wallets")
	data := map[string]interface{}{
		"walletName": "wallet1",
	}
	req.Data = data
	res, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "accounts")
	data = map[string]interface{}{
		"walletName": "wallet1",
		"accountName": "account1",
	}
	req.Data = data
	res, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if res.Error() != nil {
		t.Error(res.Error())
	}
}
