package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e/shared"
	"io/ioutil"
	"net/http"
)

type E2EBaseSetup struct {
	WorkingDir string
	RootKey string
	baseUrl string
}

func (setup *E2EBaseSetup) SignAttestation(account string, data map[string]interface{}) ([]byte,error) {
	// body
	body, err := json.Marshal(data)
	if err != nil {
		return nil,err
	}

	// build req
	targetURL := fmt.Sprintf("%s/v1/ethereum/wallet/accounts/%s/sign-attestation", setup.baseUrl, account)
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
	if err != nil {
		return nil,nil
	}
	req.Header.Set("Authorization", "Bearer " + setup.RootKey)

	// Do request
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil,err
	}
	// Read response body
	respBodyByts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	respBody := string(respBodyByts)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,fmt.Errorf("error: could not sign attesttation: respose: %s\n", respBody)
	}

	// parse to json
	retObj := make(map[string]interface{})
	err = json.Unmarshal(respBodyByts, &retObj)
	if err != nil {
		return nil, err
	}

	sigStr := retObj["data"].(map[string]interface{})["signature"].(string)
	ret, err := hex.DecodeString(sigStr)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (setup *E2EBaseSetup) PushUpdatedDb() error {
	// get store
	store,err := shared.BaseInmemStorage()
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

	fmt.Printf("e2e: setup hashicorp vault db")

	return nil
}
