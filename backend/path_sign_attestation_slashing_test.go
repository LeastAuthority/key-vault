package backend

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func ignoreError(val interface{}, err error) interface{} {
	return val
}

func basicAttestationData() map[string]interface{} {
	return map[string]interface{}{
		"public_key":      "ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279",
		"domain":          "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
		"slot":            284115,
		"committeeIndex":  2,
		"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
		"sourceEpoch":     8877,
		"sourceRoot":      "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
		"targetEpoch":     8878,
		"targetRoot":      "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
	}
}

func TestAttestationSlashing(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Successfully Sign Attestation", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		req.Data = basicAttestationData()
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
		require.Equal(t,
			"a53b6728fc2cc52abb0059da9b2e7cb01f33cd95fd6c9db7f2b821fa58a58d5ef2bc5dda058d570a7f240bf24b335eee066b2ab8dbf5a989157dd51b647733665f7c1be0d1c285b02efdbb37cd4e0ace0529b8e02c944386e3b110c32b019c63",
			res.Data["signature"],
		)
	})

	t.Run("Sign duplicated Attestation (exactly same), should sign", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

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
			"a53b6728fc2cc52abb0059da9b2e7cb01f33cd95fd6c9db7f2b821fa58a58d5ef2bc5dda058d570a7f240bf24b335eee066b2ab8dbf5a989157dd51b647733665f7c1be0d1c285b02efdbb37cd4e0ace0529b8e02c944386e3b110c32b019c63",
			res.Data["signature"],
		)
	})

	t.Run("Sign double Attestation (different block root), should return error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

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
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

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
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

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
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

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
			"918b8b0128d263a8a7f42b8c19e57c92844cd826ef6c8c7eb78b789591a934b0a89dd615ad381808cac4ae464066fb85067b806bf359321c4600b1390e037de31ad54151fbc63dc18c1e6bcec7f56e18100cc1430091e0eb28a42c4787ce5f86",
			res.Data["signature"],
		)
	})

	t.Run("Sign surrounding Attestation, should error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

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
			"public_key":      "ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279",
			"domain":          "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot":            284116,
			"committeeIndex":  2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch":     8878,
			"sourceRoot":      "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch":     8879,
			"targetRoot":      "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// surround previous vote
		// 8877 <- 8878 <- 8879
		// 	<- 8880
		// slashable
		req.Data = map[string]interface{}{
			"public_key":      "ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279",
			"domain":          "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot":            284117,
			"committeeIndex":  2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch":     8877,
			"sourceRoot":      "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch":     8880,
			"targetRoot":      "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		res, err := b.HandleRequest(context.Background(), req)
		require.Error(t, err)
		require.EqualError(t, err, "failed to sign attestation: slashable attestation (SurroundingVote), not signing")
		require.Nil(t, res)
	})

	t.Run("Sign surrounded Attestation, should error", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "accounts/sign-attestation")
		setupBaseStorage(t, req)

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
			"public_key":      "ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279",
			"domain":          "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot":            284116,
			"committeeIndex":  2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch":     8878,
			"sourceRoot":      "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch":     9000,
			"targetRoot":      "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)

		// surround previous vote
		// 8877 <- 8878 <- 8879 <----------------------9000
		// 								8900 <- 8901
		// slashable
		req.Data = map[string]interface{}{
			"public_key":      "ab321d63b7b991107a5667bf4fe853a266c2baea87d33a41c7e39a5641bfd3b5434b76f1229d452acb45ba86284e3279",
			"domain":          "01000000f071c66c6561d0b939feb15f513a019d99a84bd85635221e3ad42dac",
			"slot":            284117,
			"committeeIndex":  2,
			"beaconBlockRoot": "7b5679277ca45ea74e1deebc9d3e8c0e7d6c570b3cfaf6884be144a81dac9a0e",
			"sourceEpoch":     8900,
			"sourceRoot":      "7402fdc1ce16d449d637c34a172b349a12b2bae8d6d77e401006594d8057c33d",
			"targetEpoch":     8901,
			"targetRoot":      "17959acc370274756fa5e9fdd7e7adf17204f49cc8457e49438c42c4883cbfb0",
		}
		res, err := b.HandleRequest(context.Background(), req)
		require.Error(t, err)
		require.EqualError(t, err, "failed to sign attestation: slashable attestation (SurroundedVote), not signing")
		require.Nil(t, res)
	})
}
