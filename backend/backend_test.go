package backend

import (
	"testing"

	vault "github.com/bloxapp/eth2-key-manager"
)

func TestMain(t *testing.M) {
	vault.InitCrypto()
	t.Run()
}
