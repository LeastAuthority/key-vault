package main

import (
	"fmt"
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e"
	"log"
)

func main() {
	err := e2e.SetupE2EEnv()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	fmt.Printf("\n\nfinished\n\n")
}
