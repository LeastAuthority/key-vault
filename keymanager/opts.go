package keymanager

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Config contains cccconfiguration of the remote HTTP keymanager
type Config struct {
	Location    string `json:"location"`
	AccessToken string `json:"access_token"`
	PubKey      string `json:"public_key"`
}

var remoteOptsHelp = `The remote key manager connects to a walletd instance.  The options are:
  - location This is the address of remote HTTP wallet. E.g. http://localhost:8200
  - access_token This is an access token of a remote vault wallet.
  - public_key This is a public key that belongs to the wallet.

An sample remote HTTP keymanager options file (with annotations; these should be removed if
using this as a template) is:

  {
	"location":    	"http://host.example.com:8200", // Connect to remote HTTP wallet at http://host.example.com on port 8200
    "access_token": "x.somesupersecrettoken",   	// Access token of the wallet service above
    "public_key":   "hexencodedpublickey"  			// Public key of an account that belongs to the wallet.
  }`

// UnmarshalConfigFile attempts to JSON unmarshal a keymanager
// configuration file into the *Config{} struct.
func UnmarshalConfigFile(r io.ReadCloser) (*Config, error) {
	enc, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read config")
	}
	defer r.Close()

	cfg := &Config{}
	if err := json.Unmarshal(enc, cfg); err != nil {
		return nil, errors.Wrap(err, "could not JSON unmarshal")
	}
	return cfg, nil
}
