package e2e

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type E2EBaseSetup struct {
	RootKey string
	baseUrl string
}

const (
	pluginContainerName         = "vault-plugin-secrets-eth20_vault_1"
	logSignalingPluginInstalled = "core: successfully reloaded plugin: plugin=ethsign path=ethereum/"
	dataDirSuffix				= "/data/"
	rootTokenSuffix             = dataDirSuffix + "keys/vault.root.token"
)

func Cleanup(workingDir string) error {
	// Cleanup data dir
	dataDir := fmt.Sprintf("%s%s", workingDir, dataDirSuffix)
	_, err := os.Stat(dataDir)
	if !os.IsNotExist(err) {
		err := os.RemoveAll(dataDir)
		if err != nil {
			return err
		}
		fmt.Printf("Cleanup: deleted data dir\n")
	}


	// check if running
	listCmd := exec.Command("docker", "container", "ps", "-ls")
	byts,err := listCmd.CombinedOutput()
	if err != nil {
		return err
	}
	containerList := string(byts)
	if !strings.Contains(containerList, pluginContainerName) {
		fmt.Printf("Cleanup: NO previous container\n")
		return nil
	}

	// kill process
	killCmd := exec.Command("docker", "kill", pluginContainerName)
	err = killCmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Cleanup: killing previous container\n")

	// remove container and volumes
	removeCmd := exec.Command("docker", "rm", pluginContainerName, "--volumes")
	err = removeCmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Cleanup: deleting previous container\n")

	return nil
}

func pluginRunning(closer io.ReadCloser) <- chan bool {
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

func rootAccessToken(workingDir string) (string, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf("%s%s", workingDir, rootTokenSuffix))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func SetupE2EEnv() (*E2EBaseSetup,error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Printf("working dir: %s\n", workingDir)

	// step 1 - Cleanup
	err = Cleanup(workingDir)
	if err != nil {
		return nil,err
	}
	fmt.Printf("Cleanup done\n")

	// step 2 - run docker compose
	cmd := exec.Command("docker-compose", "up","vault")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil,err
	}

	err = cmd.Start()
	if err != nil {
		return nil,err
	}

	// step 3 - wait for plugin to be active
	<- pluginRunning(pipe)
	fmt.Printf("Plugin installed and running\n")

	// step 4 - get root access token
	token, err := rootAccessToken(workingDir)
	if err != nil {
		return nil,err
	}
	fmt.Printf("root token: %s\n", token)

	return &E2EBaseSetup{
		RootKey: token,
		baseUrl: "https://localhost:8200",
	}, nil
}
