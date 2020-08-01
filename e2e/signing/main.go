package main

import (
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e"
	"log"
)

func main() {
	setup, err := e2e.SetupE2EEnv()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	err = setup.PushUpdatedDb()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	
}
