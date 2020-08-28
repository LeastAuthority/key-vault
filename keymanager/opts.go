package keymanager

type remoteOpts struct {
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
