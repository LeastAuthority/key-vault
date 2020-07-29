package backend

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

//func TestAccountCreate(t *testing.T) {
//	b, _ := getBackend(t)
//	req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
//	storage := req.Storage
//	_, err := b.HandleRequest(context.Background(), req)
//	require.NoError(t, err)
//
//	t.Run("Successfully Create Account", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
//		req.Storage = storage
//		res, err := b.HandleRequest(context.Background(), req)
//		require.NoError(t, err)
//		require.NotNil(t, res.Data["account"])
//	})
//
//	t.Run("Successfully Create Account using private key", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account2")
//		req.Storage = storage
//		data := map[string]interface{}{
//			"key": "08f880baca147c33bea67304a54da5252a1c2b972993d2dabc8e6625240116be",
//		}
//		req.Data = data
//		res, err := b.HandleRequest(context.Background(), req)
//		require.NoError(t, err)
//		require.NotNil(t, res.Data["account"])
//	})
//
//	t.Run("Fail on create account using random key", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account2")
//		req.Storage = storage
//		data := map[string]interface{}{
//			"key": "2131868",
//		}
//		req.Data = data
//		res, err := b.HandleRequest(context.Background(), req)
//		require.Nil(t, res)
//		require.EqualError(t, err, "failed to HEX decode key: encoding/hex: odd length hex string")
//	})
//
//	t.Run("Create Account with empty name", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/ ")
//		req.Storage = storage
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "unsupported path")
//	})
//
//	t.Run("Create Account with existing name", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
//		req.Storage = storage
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to create a new validator account: account \"account1\" already exists")
//	})
//
//	t.Run("Create Account in non existing portfolio", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to open key vault: key vault not found")
//	})
//
//	t.Run("Create Account under unknown wallet", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.CreateOperation, "wallets/unknown_wallet/accounts/account1")
//		req.Storage = storage
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to retrieve wallet by name: no wallet found")
//	})
//
//	// TODO
//	//t.Run("Create Account with too long name (more than 128 characters)", func(t *testing.T) {
//	//
//	//})
//}

//func TestAccountRead(t *testing.T) {
//	b, _ := getBackend(t)
//	req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
//	storage := req.Storage
//	_, err := b.HandleRequest(context.Background(), req)
//	require.NoError(t, err)
//
//	req = logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
//	req.Storage = storage
//	_, err = b.HandleRequest(context.Background(), req)
//	require.NoError(t, err)
//
//	t.Run("Successfully Read Account", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/wallet1/accounts/account1")
//		req.Storage = storage
//		res, err := b.HandleRequest(context.Background(), req)
//		require.NoError(t, err)
//		require.NotNil(t, res.Data["account"])
//	})
//
//	t.Run("Read Account in non existing portfolio", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/wallet1/accounts/account1")
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to open key vault: key vault not found")
//	})
//
//	t.Run("Read Account of unknown wallet", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/unknown_wallet/accounts/account1")
//		req.Storage = storage
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to retrieve wallet by name: no wallet found")
//	})
//
//	t.Run("Read unknown account", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/wallet1/accounts/unknown_account")
//		req.Storage = storage
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to read a validator account: account not found")
//	})
//}

func TestAccountsList(t *testing.T) {
	b, _ := getBackend(t)


	t.Run("Successfully List Accounts", func(t *testing.T) {
		req := logical.TestRequest(t, logical.ListOperation, "wallet/accounts/")

		// setup logical storage
		_,err := baseHashicorpStorage(req.Storage, context.Background())
		require.NoError(t, err)


		res, err := b.HandleRequest(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res.Data["accounts"])
		require.Len(t, res.Data["accounts"], 1)
		require.Equal(t, res.Data["accounts"].([]map[string]string)[0]["name"], "acc")
	})

	//
	//req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
	//storage := req.Storage
	//_, err := b.HandleRequest(context.Background(), req)
	//require.NoError(t, err)
	//
	//req = logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
	//req.Storage = storage
	//_, err = b.HandleRequest(context.Background(), req)
	//require.NoError(t, err)
	//
	//t.Run("Successfully List Accounts", func(t *testing.T) {
	//	req := logical.TestRequest(t, logical.ListOperation, "wallets/wallet1/accounts/")
	//	req.Storage = storage
	//	res, err := b.HandleRequest(context.Background(), req)
	//	require.NoError(t, err)
	//	require.NotNil(t, res.Data["accounts"])
	//})
	//
	//t.Run("List Accounts in non existing portfolio", func(t *testing.T) {
	//	req := logical.TestRequest(t, logical.ListOperation, "wallets/wallet1/accounts/")
	//	_, err := b.HandleRequest(context.Background(), req)
	//	require.EqualError(t, err, "failed to open key vault: key vault not found")
	//})
	//
	//t.Run("List Accounts under unknown wallet", func(t *testing.T) {
	//	req := logical.TestRequest(t, logical.ListOperation, "wallets/unknown_wallet/accounts/")
	//	req.Storage = storage
	//	_, err := b.HandleRequest(context.Background(), req)
	//	require.EqualError(t, err, "failed to retrieve wallet by name: no wallet found")
	//})
}

//func TestAccountDepositData(t *testing.T) {
//	b, _ := getBackend(t)
//	req := logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1")
//	storage := req.Storage
//	_, err := b.HandleRequest(context.Background(), req)
//	require.NoError(t, err)
//
//	req = logical.TestRequest(t, logical.CreateOperation, "wallets/wallet1/accounts/account1")
//	req.Storage = storage
//	_, err = b.HandleRequest(context.Background(), req)
//	require.NoError(t, err)
//
//	t.Run("Successfully Get Account deposit data", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/wallet1/accounts/account1/deposit-data/")
//		req.Storage = storage
//		res, err := b.HandleRequest(context.Background(), req)
//		require.NoError(t, err)
//		require.NotNil(t, res.Data)
//	})
//
//	t.Run("Get Account deposit data in non existing portfolio", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/wallet1/accounts/account1/deposit-data/")
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to open key vault: key vault not found")
//	})
//
//	t.Run("Get Account deposit data of unknown wallet", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/unknown_wallet/accounts/account1/deposit-data/")
//		req.Storage = storage
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to retrieve wallet by name: no wallet found")
//	})
//
//	t.Run("Get Account deposit data of unknown account", func(t *testing.T) {
//		req := logical.TestRequest(t, logical.ReadOperation, "wallets/wallet1/accounts/unknown_account/deposit-data/")
//		req.Storage = storage
//		_, err := b.HandleRequest(context.Background(), req)
//		require.EqualError(t, err, "failed to retrieve account by name: account not found")
//	})
//}
