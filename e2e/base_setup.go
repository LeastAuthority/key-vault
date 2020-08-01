package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/bloxapp/KeyVault/wallet_hd"
	types "github.com/wealdtech/go-eth2-types/v2"
	"io/ioutil"
	"net/http"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func baseInmemStorage() (*in_memory.InMemStore, error) {
	types.InitBLS()
	store := in_memory.NewInMemStore()

	// wallet
	wallet := wallet_hd.NewHDWallet(&core.WalletContext{Storage:store})
	err := store.SaveWallet(wallet)
	if err != nil {
		return nil, err
	}

	// account
	acc,err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), "test_account")
	if err != nil {
		return nil, err
	}
	err = store.SaveAccount(acc)
	if err != nil {
		return nil, err
	}

	return store, nil
}

type E2EBaseSetup struct {
	RootKey string
	baseUrl string
}

func (setup *E2EBaseSetup) PushUpdatedDb() error {
	// get store
	store,err := baseInmemStorage()
	if err != nil {
		return err
	}

	// encode store
	byts, err := json.Marshal(store)
	if err != nil {
		return nil
	}
	encodedStore := hex.EncodeToString(byts)

	// body
	body, err := json.Marshal(map[string]string{
		"data": encodedStore,
	})

	// build req
	targetURL := fmt.Sprintf("%s/v1/ethereum/admin/pushUpdate", setup.baseUrl)
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer " + setup.RootKey)

	// Do request
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	// Read response body
	respBodyByts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBody := string(respBodyByts)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: could not update vault: respose: %s", respBody)
	}

	return nil
}
