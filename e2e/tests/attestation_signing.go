package tests

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/bloxapp/vault-plugin-secrets-eth2.0/e2e"
)

func ignoreError(val interface{}, err error) interface{} {
	return val
}

type AttestationSigning struct {

}

func (t *AttestationSigning)Name() string {
	return "Test attestation signing"
}

func (t *AttestationSigning)Run() error {
	setup, err := e2e.SetupE2EEnv()
	if err != nil {
		return err
	}

	err = setup.PushUpdatedDb()
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

	expecetd := ignoreError(hex.DecodeString("b3234e48fa4d7b9df6f743aad1fa1c54889b3a1cff0649441731a129359c7ad568a2fce3181ed2b767a369684974f67a1960ec139595aa5347883698ab0af2236310cf4f1d59483abe2cefcfc3a79b453a7ffea4d2268aad314fdac5b468984f")).([]byte)
	if bytes.Compare(sig, expecetd) != 0 {
		return fmt.Errorf("e2e: attestation signature not valid")
	}
	return nil
}

