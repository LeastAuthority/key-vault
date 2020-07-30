package backend

import (
	"context"
	"encoding/hex"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"testing"
)

func ignoreError(val interface{}, err error) interface{} {
	return val
}

func basicAttestationData() map[string]interface{} {
	return map[string]interface{}{
		"domain": "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
		"slot": 284115,
		"committeeIndex": 2,
		"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
		"sourceEpoch": 8877,
		"sourceRoot": "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
		"targetEpoch": 8878,
		"targetRoot": "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
	}
}

func TestAttestationSlashing(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Successfully Sign Attestation", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicAttestationData()
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
		require.Equal(t,
			ignoreError(hex.DecodeString("b3234e48fa4d7b9df6f743aad1fa1c54889b3a1cff0649441731a129359c7ad568a2fce3181ed2b767a369684974f67a1960ec139595aa5347883698ab0af2236310cf4f1d59483abe2cefcfc3a79b453a7ffea4d2268aad314fdac5b468984f")).([]byte),
			res.Data["signature"],
			)
	})

	t.Run("Sign duplicated Attestation (exactly same), should sign", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first attestation
		req.Data = basicAttestationData()
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)

		// duplicated attestation
		req.Data = basicAttestationData()
		res, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
		require.Equal(t,
			ignoreError(hex.DecodeString("b3234e48fa4d7b9df6f743aad1fa1c54889b3a1cff0649441731a129359c7ad568a2fce3181ed2b767a369684974f67a1960ec139595aa5347883698ab0af2236310cf4f1d59483abe2cefcfc3a79b453a7ffea4d2268aad314fdac5b468984f")).([]byte),
			res.Data["signature"],
		)
	})

	t.Run("Sign double Attestation (different block root), should return error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first attestation
		req.Data = basicAttestationData()
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)

		// slashable attestation
		data := basicAttestationData()
		data["beaconBlockRoot"] = "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0f"
		req.Data = data
		res, err = b.HandleRequest(context.Background(), req)
		require.Error(t, err)
		require.EqualError(t, err, "failed to sign attestation: slashable attestation (DoubleVote), not signing")
		require.Nil(t, res)
	})

	t.Run("Sign double Attestation (different source root), should return error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first attestation
		req.Data = basicAttestationData()
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)

		// slashable attestation
		data := basicAttestationData()
		data["sourceRoot"] = "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33e"
		req.Data = data
		res, err = b.HandleRequest(context.Background(), req)
		require.Error(t, err)
		require.EqualError(t, err, "failed to sign attestation: slashable attestation (DoubleVote), not signing")
		require.Nil(t, res)
	})

	t.Run("Sign double Attestation (different target root), should return error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first attestation
		req.Data = basicAttestationData()
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)

		// slashable attestation
		data := basicAttestationData()
		data["targetRoot"] = "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb1"
		req.Data = data
		res, err = b.HandleRequest(context.Background(), req)
		require.Error(t, err)
		require.EqualError(t, err, "failed to sign attestation: slashable attestation (DoubleVote), not signing")
		require.Nil(t, res)
	})

	t.Run("Sign Attestation (different domain), should sign", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first attestation
		req.Data = basicAttestationData()
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)

		// slashable attestation
		data := basicAttestationData()
		data["domain"] = "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dad"
		req.Data = data
		res, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
		require.Equal(t,
			ignoreError(hex.DecodeString("a720ff1e58bbafa59a7e6adc0e0bc1991fede153982a63bdf00e9a898165c1f39408da5170a765bacff1ff0bfe33425e016abf5e5b9a7e8c1323fb26be0a185a78b9a6d798756fd056ba1072d3b449ac8a695e10c5c979c2cfca37394b6bc0af")).([]byte),
			res.Data["signature"],
		)
	})

	t.Run("Sign surrounding Attestation, should error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first attestation
		req.Data = basicAttestationData()
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// add another attestation building on the base
		// 8877 <- 8878 <- 8879
		req.Data = map[string]interface{}{
			"domain": "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot": 284116,
			"committeeIndex": 2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch": 8878,
			"sourceRoot": "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch": 8879,
			"targetRoot": "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)


		// surround previous vote
		// 8877 <- 8878 <- 8879
		// 	<- 8880
		// slashable
		req.Data = map[string]interface{}{
			"domain": "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot": 284117,
			"committeeIndex": 2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch": 8877,
			"sourceRoot": "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch": 8880,
			"targetRoot": "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		res, err := b.HandleRequest(context.Background(), req)
		require.Error(t, err)
		require.EqualError(t, err, "failed to sign attestation: slashable attestation (SurroundingVote), not signing")
		require.Nil(t, res)
	})

	t.Run("Sign surrounded Attestation, should error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		// first attestation
		req.Data = basicAttestationData()
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// add another attestation building on the base
		// 8877 <- 8878 <- 8879 <----------------------9000
		req.Data = map[string]interface{}{
			"domain": "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot": 284116,
			"committeeIndex": 2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch": 8878,
			"sourceRoot": "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch": 9000,
			"targetRoot": "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)


		// surround previous vote
		// 8877 <- 8878 <- 8879 <----------------------9000
		// 								8900 <- 8901
		// slashable
		req.Data = map[string]interface{}{
			"domain": "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot": 284117,
			"committeeIndex": 2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch": 8900,
			"sourceRoot": "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch": 8901,
			"targetRoot": "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		res, err := b.HandleRequest(context.Background(), req)
		require.Error(t, err)
		require.EqualError(t, err, "failed to sign attestation: slashable attestation (SurroundedVote), not signing")
		require.Nil(t, res)
	})
}
