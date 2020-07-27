package backend

import (
	"context"
	"encoding/hex"

	vault "github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/eth1_deposit"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
)

func accountsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name") + "/accounts/",
			HelpSynopsis:    "List wallet accounts",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathWalletAccountsList,
			},
		},
		&framework.Path{
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name") + "/accounts/" + framework.GenericNameRegex("account_name"),
			HelpSynopsis:    "Create/Read an Ethereum 2.0 account",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name":  &framework.FieldSchema{Type: framework.TypeString},
				"account_name": &framework.FieldSchema{Type: framework.TypeString},
				"key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Private Key",
					Default:     "",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletsAccountCreate,
				logical.ReadOperation:   b.pathWalletsAccountRead,
			},
		},
		&framework.Path{
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name") + "/accounts/" + framework.GenericNameRegex("account_name") + "/deposit-data/",
			HelpSynopsis:    "Get Deposit Data",
			HelpDescription: `Get ETH1 Deposit Data`,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name":  &framework.FieldSchema{Type: framework.TypeString},
				"account_name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathWalletsAccountDepositData,
			},
		},
	}
}

func (b *backend) pathWalletsAccountCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	accountName := data.Get("account_name").(string)
	key := data.Get("key").(string)
	var keyBytes []byte
	var err error

	if len(key) != 0 {
		// Decode key
		keyBytes, err = hex.DecodeString(key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to HEX decode key")
		}
	}

	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)

	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet by name")
	}

	account, err := wallet.CreateValidatorAccount(accountName, keyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a new validator account")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"account": account,
		},
	}, nil
}

func (b *backend) pathWalletsAccountRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	accountName := data.Get("account_name").(string)

	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)

	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet by name")
	}

	account, err := wallet.AccountByName(accountName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read a validator account")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"account": account,
		},
	}, nil
}

func (b *backend) pathWalletAccountsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)

	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet by name")
	}

	var accounts []core.Account
	for a := range wallet.Accounts() {
		accounts = append(accounts, a)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"accounts": accounts,
		},
	}, nil
}

func (b *backend) pathWalletsAccountDepositData(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	accountName := data.Get("account_name").(string)
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)

	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet by name")
	}

	account, err := wallet.AccountByName(accountName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve account by name")
	}

	withdrawal, err := wallet.GetWithdrawalAccount()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get withdrawal account")
	}

	depositData, root, err := eth1_deposit.DepositData(account, withdrawal, eth1_deposit.MaxEffectiveBalanceInGwei)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get deposit data")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"amount":                depositData.GetAmount(),
			"publicKey":             hex.EncodeToString(depositData.GetPublicKey()),
			"signature":             hex.EncodeToString(depositData.GetSignature()),
			"withdrawalCredentials": hex.EncodeToString(depositData.GetWithdrawalCredentials()),
			"depositDataRoot":       hex.EncodeToString(root[:]),
		},
	}, nil
}
