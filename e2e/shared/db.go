package shared

import (
	"encoding/hex"
	"testing"

	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/bloxapp/KeyVault/wallet_hd"
	"github.com/stretchr/testify/require"
	types "github.com/wealdtech/go-eth2-types/v2"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

// BaseInmemStorage creates the in-memory storage and creates the base account.
func BaseInmemStorage(t *testing.T) (*in_memory.InMemStore, error) {
	types.InitBLS()
	store := in_memory.NewInMemStore()

	entropy, err := core.GenerateNewEntropy()
	require.NoError(t, err)

	seed, err := core.SeedFromEntropy(entropy, "test_password")
	require.NoError(t, err)

	// wallet
	wallet := wallet_hd.NewHDWallet(&core.WalletContext{Storage: store})
	if err := store.SaveWallet(wallet); err != nil {
		return nil, err
	}

	// account
	// acc, err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), "test_account")
	acc, err := wallet.CreateValidatorAccount(seed, "test_account")
	if err != nil {
		return nil, err
	}
	err = store.SaveAccount(acc)
	if err != nil {
		return nil, err
	}

	return store, nil
}
