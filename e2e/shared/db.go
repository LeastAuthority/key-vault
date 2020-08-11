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

// AccountName is the test account name.
const AccountName = "test_account"

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
	acc, err := wallet.CreateValidatorAccount(seed, AccountName)
	if err != nil {
		return nil, err
	}
	err = store.SaveAccount(acc)
	if err != nil {
		return nil, err
	}

	return store, nil
}

// RetrieveAccount retrieves test account fro the storage.
func RetrieveAccount(t *testing.T, store core.Storage) core.ValidatorAccount {
	accounts, err := store.ListAccounts()
	require.NoError(t, err)

	for _, acc := range accounts {
		if acc.Name() == AccountName {
			return acc
		}
	}
	return nil
}
