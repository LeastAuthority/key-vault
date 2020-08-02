package tests

import (
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e"
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e/shared"
	"github.com/stretchr/testify/require"
	"testing"
)


type AttestationSigning struct {

}

func (test *AttestationSigning)Name() string {
	return "Test attestation signing"
}

func (test *AttestationSigning)Run(t *testing.T) {
	setup, err := e2e.SetupE2EEnv()
	require.NoError(t, err)

	// setup vault with db
	err = setup.PushUpdatedDb()
	require.NoError(t, err)

	// sign
	sig,err := setup.SignAttestation(
		"test_account",
		map[string]interface{}{
			"domain": "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot": 284115,
			"committeeIndex": 2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch": 8877,
			"sourceRoot": "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch": 8878,
			"targetRoot": "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		},
	)
	require.NoError(t, err)

	expecetd := shared.HexToBytes("a53b6728fc2cc52abb0059da9b2e7cb01f33cd95fd6c9db7f2b821fa58a58d5ef2bc5dda058d570a7f240bf24b335eee066b2ab8dbf5a989157dd51b647733665f7c1be0d1c285b02efdbb37cd4e0ace0529b8e02c944386e3b110c32b019c63")
	require.Equal(t, expecetd, sig)

	// cleanup
	require.NoError(t, setup.Cleanup())
}

