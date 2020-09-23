package e2e

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/bloxapp/eth2-key-manager/core"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/key-vault/e2e/launcher"
	"github.com/bloxapp/key-vault/e2e/shared"
)

var (
	dockerLauncher *launcher.Docker
	// basePath       = os.Getenv("GOPATH") + "/src/github.com/bloxapp/key-vault"
	basePath = os.Getenv("PWD") + "/../../"
)

// ServiceError represents service error type.
type ServiceError struct {
	Data map[string]interface{}
}

// NewServiceError is the constructor of ServiceError.
func NewServiceError(data map[string]interface{}) *ServiceError {
	return &ServiceError{
		Data: data,
	}
}

// Error implements error interface.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("%#v", e.Data)
}

// ErrorValue returns error value from the data
func (e *ServiceError) ErrorValue() string {
	return e.Data["errors"].([]interface{})[0].(string)
}

// DataValue returns "field" value from data
func (e *ServiceError) DataValue(field string) interface{} {
	return e.Data["data"].(map[string]interface{})[field]
}

func init() {
	var err error
	imageName := "key-vault:" + uuid.New()
	if dockerLauncher, err = launcher.New(imageName, basePath); err != nil {
		log.Fatal(err)
	}
}

// BaseSetup implements mechanism, to setup base env for e2e tests.
type BaseSetup struct {
	WorkingDir string
	RootKey    string
	baseURL    string
}

// Setup sets up environment for e2e tests
func Setup(t *testing.T) *BaseSetup {
	conf, err := dockerLauncher.Launch(context.Background(), uuid.New())
	require.NoError(t, err)
	t.Cleanup(func() {
		err := dockerLauncher.Stop(context.Background(), conf.ID)
		require.NoError(t, err)
	})

	return &BaseSetup{
		RootKey: conf.Token,
		baseURL: conf.URL,
	}
}

// SignAttestation tests the sign attestation endpoint.
func (setup *BaseSetup) SignAttestation(data map[string]interface{}) ([]byte, error) {
	// body
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// build req
	targetURL := fmt.Sprintf("%s/v1/ethereum/test/accounts/sign-attestation", setup.baseURL)
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
		return nil, NewServiceError(retObj)
	}

	sigStr := retObj["data"].(map[string]interface{})["signature"].(string)
	ret, err := hex.DecodeString(sigStr)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UpdateConfig updates the config.
func (setup *BaseSetup) UpdateConfig(t *testing.T) {
	// body
	body, err := json.Marshal(map[string]string{
		"network":      "test",
		"genesis_time": "2020-08-04 13:00:08 UTC",
	})
	require.NoError(t, err)

	// build req
	targetURL := fmt.Sprintf("%s/v1/ethereum/test/config", setup.baseURL)
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

	require.Equal(t, http.StatusOK, resp.StatusCode, respBody)

	fmt.Printf("e2e: setup hashicorp vault db\n")
}

// UpdateStorage updates the storage.
func (setup *BaseSetup) UpdateStorage(t *testing.T) core.Storage {
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
	targetURL := fmt.Sprintf("%s/v1/ethereum/test/storage", setup.baseURL)
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

	require.Equal(t, http.StatusOK, resp.StatusCode, respBody)

	fmt.Printf("e2e: setup hashicorp vault db\n")

	return store
}
