package backend

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func setupStorageWithWalletAndAccounts(storage logical.Storage) error {
	_,err := baseHashicorpStorage(storage, context.Background())
	return err
}

func TestSignAttestation(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Successfully Sign Attestation", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		data := map[string]interface{}{
			"domain": "",
		}
		req.Data = data
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
	})

	t.Run("Sign Attestation in non existing key vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-attestation")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to open key vault: wallet not found")
	})

	t.Run("Sign Attestation of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/unknown_account/sign-attestation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		_, err = b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign attestation: account not found")
	})
}

func TestSignProposal(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Successfully Sign Proposal", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-proposal")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
	})

	t.Run("Sign Proposal in non existing key vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-proposal")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to open key vault: wallet not found")
	})

	t.Run("Sign Proposal of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/unknown_account/sign-proposal")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		_, err = b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign data: account not found")
	})
}

func TestSignAggregation(t *testing.T) {
	b, _ := getBackend(t)

	t.Run("Successfully Sign Aggregation", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-aggregation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		_, err = b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
	})

	t.Run("Sign Aggregation in non existing key vault", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/test_account/sign-aggregation")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to open key vault: wallet not found")
	})

	t.Run("Sign Aggregation of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallet/accounts/unknown_account/sign-aggregation")

		// setup storage
		err := setupStorageWithWalletAndAccounts(req.Storage)
		require.NoError(t, err)

		_, err = b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign data: account not found")
	})
}
