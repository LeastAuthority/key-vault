package tests

import (
	"testing"

	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e"
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e/shared"
	"github.com/stretchr/testify/require"
)

type AttestationSigning struct {
}

func (test *AttestationSigning) Name() string {
	return "Test attestation signing"
}

func (test *AttestationSigning) Run(t *testing.T) {
	setup := e2e.SetupE2EEnv(t)

	// setup vault with db
	err := setup.UpdateStorage()
	require.NoError(t, err)

	// sign
	sig, err := setup.SignAttestation(
		map[string]interface{}{
			"public_key":      "ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279",
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

	expecetd := shared.HexToBytes("a53b6728fc2cc52abb0059da9b2e7cb01f33cd95fd6c9db7f2b821fa58a58d5ef2bc5dda058d570a7f240bf24b335eee066b2ab8dbf5a989157dd51b647733665f7c1be0d1c285b02efdbb37cd4e0ace0529b8e02c944386e3b110c32b019c63")
	require.Equal(t, expecetd, sig)

	// cleanup
	setup.Cleanup(t)
}
