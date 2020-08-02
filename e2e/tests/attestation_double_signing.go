package tests

import (
	"fmt"
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e"
	"github.com/stretchr/testify/require"
	"testing"
)


type AttestationDoubleSigning struct {

}

func (test *AttestationDoubleSigning)Name() string {
	return "Test double attestation signing, different block root"
}

func (test *AttestationDoubleSigning)Run(t *testing.T) {
	setup, err := e2e.SetupE2EEnv()
	require.NoError(t, err)

	// setup vault with db
	err = setup.PushUpdatedDb()
	require.NoError(t, err)

	// first sig
	_,err = setup.SignAttestation(
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

	// second sig, different block root
	_,err = setup.SignAttestation(
		"test_account",
		map[string]interface{}{
			"domain": "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot": 284115,
			"committeeIndex": 2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0f",
			"sourceEpoch": 8877,
			"sourceRoot": "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch": 8878,
			"targetRoot": "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		},
	)
	expectedErr := fmt.Sprintf("1 error occurred:\n\t* failed to sign attestation: slashable attestation (DoubleVote), not signing\n\n")
	require.Error(t, err)
	require.EqualError(t, err, expectedErr)

	// cleanup
	require.NoError(t, setup.Cleanup())
}

