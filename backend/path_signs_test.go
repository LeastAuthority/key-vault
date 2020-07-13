package backend

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestSignAttestation(t *testing.T) {
	b, _ := getBackend(t)
	req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
	storage := req.Storage
	_, err := b.HandleRequest(context.Background(), req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
	req.Storage = storage
	_, err = b.HandleRequest(context.Background(), req)
	require.NoError(t, err)

	t.Run("Successfully Sign Attestation", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1/sign-attestation")
		req.Storage = storage
		data := map[string]interface{}{
			"domain": "",
		}
		req.Data = data
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
	})

	t.Run("Sign Attestation in non existing portfolio", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1/sign-attestation")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to open key vault: key vault not found")
	})

	t.Run("Sign Attestation in unknown wallet", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/unknown_wallet/accounts/account1/sign-attestation")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to retrieve wallet by name: no wallet found")
	})

	t.Run("Sign Attestation of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/unknown_account/sign-attestation")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign data: account not found")
	})
}

func TestSignProposal(t *testing.T) {
	b, _ := getBackend(t)
	req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
	storage := req.Storage
	_, err := b.HandleRequest(context.Background(), req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
	req.Storage = storage
	_, err = b.HandleRequest(context.Background(), req)
	require.NoError(t, err)

	t.Run("Successfully Sign Proposal", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1/sign-proposal")
		req.Storage = storage
		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data)
	})

	t.Run("Sign Proposal in non existing portfolio", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1/sign-proposal")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to open key vault: key vault not found")
	})

	t.Run("Sign Proposal in unknown wallet", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/unknown_wallet/accounts/account1/sign-proposal")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to retrieve wallet by name: no wallet found")
	})

	t.Run("Sign Proposal of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/unknown_account/sign-proposal")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign data: account not found")
	})
}

func TestSignAggregation(t *testing.T) {
	b, _ := getBackend(t)
	req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
	storage := req.Storage
	_, err := b.HandleRequest(context.Background(), req)
	require.NoError(t, err)

	req = logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
	req.Storage = storage
	res, err := b.HandleRequest(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res.Data)

	t.Run("Successfully Sign Aggregation", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1/sign-aggregation")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
	})

	t.Run("Sign Aggregation in non existing portfolio", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1/sign-aggregation")
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to open key vault: key vault not found")
	})

	t.Run("Sign Aggregation in unknown wallet", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/unknown_wallet/accounts/account1/sign-aggregation")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to retrieve wallet by name: no wallet found")
	})

	t.Run("Sign Aggregation of unknown account", func(t *testing.T) {
		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/unknown_account/sign-aggregation")
		req.Storage = storage
		_, err := b.HandleRequest(context.Background(), req)
		require.EqualError(t, err, "failed to sign data: account not found")
	})
}
