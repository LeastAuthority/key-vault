package shared

import (
	"encoding/hex"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	"github.com/bloxapp/KeyVault/wallet_hd"
	types "github.com/wealdtech/go-eth2-types/v2"
)

func _byteArray(input string) []byte {
	res, _ := hex.DecodeString(input)
	return res
}

func BaseInmemStorage() (*in_memory.InMemStore, error) {
	types.InitBLS()
	store := in_memory.NewInMemStore()

	// wallet
	wallet := wallet_hd.NewHDWallet(&core.WalletContext{Storage:store})
	err := store.SaveWallet(wallet)
	if err != nil {
		return nil, err
	}

	// account
	acc,err := wallet.CreateValidatorAccount(_byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff"), "test_account")
	if err != nil {
		return nil, err
	}
	err = store.SaveAccount(acc)
	if err != nil {
		return nil, err
	}

	return store, nil
}
