package backend

import (
	"context"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
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

	req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
	storage := req.Storage
	res, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
	req.Storage = storage
	res, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ListOperation, "wallets/wallet1/accounts/")
	req.Storage = storage
	res, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	//data = map[string]interface{}{
	//	"walletName":  "wallet1",
	//	"accountName": "account1",
	//}
	//req.Data = data
	if res.Error() != nil {
		t.Error(res.Error())
	}
}
