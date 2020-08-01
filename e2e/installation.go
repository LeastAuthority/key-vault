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

const (
	pluginContainerName         = "vault-plugin-secrets-eth20_vault_1"
	logSignalingPluginInstalled = "core: successfully reloaded plugin: plugin=ethsign path=ethereum/"
	dataDirSuffix				= "/data/"
	rootTokenSuffix             = dataDirSuffix + "keys/vault.root.token"
)

func (setup *E2EBaseSetup) Cleanup() error {
	// Cleanup data dir
	dataDir := fmt.Sprintf("%s%s", setup.WorkingDir, dataDirSuffix)
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
		return nil
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
			//fmt.Println(newLine)

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

	ret := string(content)
	ret = strings.TrimSuffix(ret, "\n")
	return ret, nil
}

func SetupE2EEnv() (*E2EBaseSetup,error) {
	ret := &E2EBaseSetup{}

	workingDir, err := os.Getwd()
	workingDir = strings.ReplaceAll(workingDir, "/e2e/tests", "") // since tests run from 2e2/tests.. remove that from working dir
	if err != nil {
		return nil, err
	}
	fmt.Printf("e2e: working dir: %s\n", workingDir)
	ret.WorkingDir = workingDir

	// step 1 - Cleanup
	err = ret.Cleanup()
	if err != nil {
		return nil,err
	}
	fmt.Printf("e2e: Cleanup done\n")

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
	fmt.Printf("e2e: Plugin installed and running\n")

	// step 4 - get root access token
	token, err := rootAccessToken(workingDir)
	if err != nil {
		return nil,err
	}
	fmt.Printf("e2e: root token: %s\n", token)

	return &E2EBaseSetup{
		RootKey: token,
		baseUrl: "http://localhost:8200",
	}, nil
}
