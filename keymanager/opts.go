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
