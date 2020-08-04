package e2e

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	pluginContainerName         = "vault-plugin-secrets-eth20_vault_1"
	logSignalingPluginInstalled = "core: successfully reloaded plugin: plugin=ethsign path=ethereum/"
	dataDirSuffix               = "/data/"
	rootTokenSuffix             = dataDirSuffix + "keys/vault.root.token"
)

// Cleanup cleans up the environment
func (setup *BaseSetup) Cleanup(t *testing.T) {
	// Cleanup data dir
	dataDir := fmt.Sprintf("%s%s", setup.WorkingDir, dataDirSuffix)
	_, err := os.Stat(dataDir)
	if !os.IsNotExist(err) {
		err := os.RemoveAll(dataDir)
		require.NoError(t, err)
		fmt.Printf("Cleanup: deleted data dir\n")
	}

	// check if running
	listCmd := exec.Command("docker", "container", "ps", "-ls")
	byts, err := listCmd.CombinedOutput()
	require.NoError(t, err)

	containerList := string(byts)
	if !strings.Contains(containerList, pluginContainerName) {
		return
	}

	// kill process
	killCmd := exec.Command("docker", "kill", pluginContainerName)
	_ = killCmd.Run()
	//if err != nil {
	//	return err
	//}
	fmt.Printf("Cleanup: killing previous container\n")

	// remove container and volumes
	removeCmd := exec.Command("docker", "rm", pluginContainerName, "--volumes")
	err = removeCmd.Run()
	require.NoError(t, err)
	fmt.Printf("Cleanup: deleting previous container\n")
}

func pluginRunning(closer io.ReadCloser) <-chan bool {
	ret := make(chan bool)

	scanner := bufio.NewScanner(closer)

	go func() {
		for scanner.Scan() {
			newLine := scanner.Text()
			fmt.Println(newLine)

			if strings.Contains(newLine, logSignalingPluginInstalled) {
				ret <- true
				close(ret)
				return
			}
		}
	}()

	return ret
}

func rootAccessToken(t *testing.T, workingDir string) (string, error) {
	var data bytes.Buffer
	cmd := exec.Command("docker-compose", "exec", "-T", "vault", "cat", rootTokenSuffix)
	cmd.Stdout = &data
	require.NoError(t, cmd.Run())

	token, err := ioutil.ReadAll(&data)
	require.NoError(t, err)

	return strings.TrimSuffix(string(token), "\n"), nil
}

var buildOnce sync.Once

// SetupE2EEnv sets up environment for e2e tests
func SetupE2EEnv(t *testing.T) *BaseSetup {
	return &BaseSetup{
		RootKey: "sometoken",
		baseURL: "http://localhost:8200",
	}

	ret := &BaseSetup{}

	workingDir, err := os.Getwd()
	workingDir = strings.ReplaceAll(workingDir, "/e2e/tests", "") // since tests run from 2e2/tests.. remove that from working dir
	require.NoError(t, err)
	fmt.Printf("e2e: working dir: %s\n", workingDir)
	ret.WorkingDir = workingDir

	// step 1 - Cleanup
	ret.Cleanup(t)
	fmt.Printf("e2e: Cleanup done\n")

	// step 2 - build (once per run)
	buildOnce.Do(func() {
		build := exec.Command("docker-compose", "build", "vault")
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr
		require.NoError(t, build.Run())
		fmt.Printf("e2e: Built Vault docker\n")
	})
	require.NoError(t, err)

	// step 3 - run docker compose
	cmd := exec.Command("docker-compose", "up", "vault")
	pipe, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)

	// step 4 - wait for plugin to be active
	<-pluginRunning(pipe)
	fmt.Printf("e2e: Plugin installed and running\n")

	// step 5 - get root access token
	token, err := rootAccessToken(t, workingDir)
	require.NoError(t, err)
	fmt.Printf("e2e: root token: %s\n", token)

	return &BaseSetup{
		RootKey: token,
		baseURL: "http://localhost:8200",
	}
}
