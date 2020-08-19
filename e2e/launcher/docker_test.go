package launcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLaunch(t *testing.T) {
	imageName := "vault-plugin-secrets-eth20_vault:latest"
	if envImageName := os.Getenv("VAULT_PLUGIN_IMAGE"); len(envImageName) > 0 {
		imageName = envImageName
	}

	launcher, err := New(logrus.New(), imageName)
	require.NoError(t, err)

	config, err := launcher.Launch(context.Background(), "")
	require.NoError(t, err)
	require.NotEmpty(t, config.ID)
	require.NotEmpty(t, config.URL)
	require.NotEmpty(t, config.Token)

	t.Cleanup(func() {
		err = launcher.Stop(context.Background(), config.ID)
		require.NoError(t, err)
	})

	// Prepare request
	req, err := http.NewRequest("LIST", fmt.Sprintf("%s/v1/ethereum/accounts", config.URL), nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+config.Token)

	// Send request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
	require.Equal(t, `{"errors":["1 error occurred:\n\t* failed to open key vault: wallet not found\n\n"]}`, strings.TrimSpace(string(body)))
}
