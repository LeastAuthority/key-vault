package backend

import (
	"context"

	vault "github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/eth1_deposit"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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
			HelpSynopsis:    "Create an Ethereum 2.0 account",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name":  &framework.FieldSchema{Type: framework.TypeString},
				"account_name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletsAccountCreate,
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
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)
	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, err
	}
	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, err
	}
	account, err := wallet.CreateValidatorAccount(accountName)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, err
	}
	accounts := make([]core.Account, 0)
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
		return nil, err
	}
	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, err
	}
	account, err := wallet.AccountByName(accountName)
	if err != nil {
		return nil, err
	}
	withdrawal, err := wallet.GetWithdrawalAccount()

	depositData, _, err := eth1_deposit.DepositData(account, withdrawal, eth1_deposit.MaxEffectiveBalanceInGwei)

	return &logical.Response{
		Data: map[string]interface{}{
			"depositData": depositData,
		},
	}, nil
}
