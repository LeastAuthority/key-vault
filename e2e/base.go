package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func cleanup() error {
	const containerName = "vault-plugin-secrets-eth20_vault_1"
	// check if running
	listCmd := exec.Command("docker", "container", "ps", "-ls")
	byts,err := listCmd.CombinedOutput()
	if err != nil {
		return err
	}
	containerList := string(byts)
	if !strings.Contains(containerList, containerName) {
		return nil
	}

	// kill process
	killCmd := exec.Command("docker", "kill",containerName)
	err = killCmd.Run()
	if err != nil {
		return err
	}

	// remove container and volumes
	removeCmd := exec.Command("docker", "rm",containerName, "--volumes")
	err = removeCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func SetupE2EEnv() error {
	// step 1 - cleanup
	err := cleanup()
	if err != nil {
		return err
	}
	fmt.Printf("cleanup done\n")

	// step 2 - run docker compose
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Printf("working dir: %s\n", workingDir)

	cmd := exec.Command("docker-compose", "up","vault")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
