package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e/shared"
)

// BaseSetup implements mechanism, to setup base env for e2e tests.
type BaseSetup struct {
	WorkingDir string
	RootKey    string
	baseURL    string
}

// SignAttestation tests the sign attestation endpoint.
func (setup *BaseSetup) SignAttestation(data map[string]interface{}) ([]byte, error) {
	// body
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// build req
	targetURL := fmt.Sprintf("%s/v1/ethereum/accounts/sign-attestation", setup.baseURL)
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil
	}
	req.Header.Set("Authorization", "Bearer "+setup.RootKey)

	// Do request
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// Read response body
	respBodyByts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse to json
	retObj := make(map[string]interface{})
	err = json.Unmarshal(respBodyByts, &retObj)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", retObj["errors"].([]interface{})[0])
	}

	sigStr := retObj["data"].(map[string]interface{})["signature"].(string)
	ret, err := hex.DecodeString(sigStr)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UpdateStorage updates the storage.
func (setup *BaseSetup) UpdateStorage(t *testing.T) error {
	// get store
	store, err := shared.BaseInmemStorage(t)
	require.NoError(t, err)

	// encode store
	byts, err := json.Marshal(store)
	require.NoError(t, err)

	encodedStore := hex.EncodeToString(byts)

	// body
	body, err := json.Marshal(map[string]string{
		"data": encodedStore,
	})
	require.NoError(t, err)

	// build req
	targetURL := fmt.Sprintf("%s/v1/ethereum/storage", setup.baseURL)
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewBuffer(body))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+setup.RootKey)

	// Do request
	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	require.NoError(t, err)

	// Read response body
	respBodyByts, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	respBody := string(respBodyByts)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: could not update vault: respose: %s", respBody)
	}

	fmt.Printf("e2e: setup hashicorp vault db\n")

	return nil
}
