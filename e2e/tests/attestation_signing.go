package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e"
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e/shared"
)

// AttestationSigning tests sign attestation endpoint.
type AttestationSigning struct {
}

// Name returns the name of the test.
func (test *AttestationSigning) Name() string {
	return "Test attestation signing"
}

// Run run the test.
func (test *AttestationSigning) Run(t *testing.T) {
	setup := e2e.SetupE2EEnv(t)

	// setup vault with db
	storage := setup.UpdateStorage(t)
	account := shared.RetrieveAccount(t, storage)
	pubKey := hex.EncodeToString(account.ValidatorPublicKey().Marshal())
	fmt.Println("pubKey", pubKey)

	// sign
	sig, err := setup.SignAttestation(
		map[string]interface{}{
			"public_key":      pubKey,
			"domain":          "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot":            284115,
			"committeeIndex":  2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch":     8877,
			"sourceRoot":      "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch":     8878,
			"targetRoot":      "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		},
	)
	require.NoError(t, err)

	actualHex := hex.EncodeToString(sig)
	expected := "89807a4941144deb709e4f74827475d422914511c4936885c52099c2484f6001b5aed37c4a5d64723aaf620a7a96aa630f59a9add3117ab4c6c9ebf8bd66415acaf9199702b6e6a54a3aba90debb0615e20c9df3b08df60b31b378f543040ee8"
	require.Equal(t, expected, actualHex)

	// cleanup
	setup.Cleanup(t)
}
