package main

import (
	"log"
	"os"

	vault "github.com/bloxapp/eth2-key-manager"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"

	"github.com/bloxapp/key-vault/backend"
)

// Version contains the current version of app binary.
// Basically, this is the commit hash
var Version = "latest"

func init() {
	// This is needed for signing methods
	vault.InitCrypto()
}

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:]) // Ignore command, strictly parse flags

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: backend.Factory(Version),
		TLSProviderFunc:    tlsProviderFunc,
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
