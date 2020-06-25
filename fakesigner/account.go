package fakesigner

import (
	"encoding/hex"
	"encoding/json"

	"github.com/bloxapp/KeyVault/core"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/shared/keystore"
)

// ReplacePublicKey replaces account's public key
func ReplacePublicKey(acc core.Account) (map[string]interface{}, error) {
	derivableKey, err := keystore.DecryptKey([]byte(key), "changeit")
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt key")
	}

	data, err := json.Marshal(acc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal account")
	}

	account := make(map[string]interface{})
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal account")
	}

	key, ok := account["key"].(map[string]interface{})
	if !ok {
		key = make(map[string]interface{})
	}

	key["pubkey"] = hex.EncodeToString(derivableKey.PublicKey.Marshal())

	account["key"] = key
	return account, nil
}
