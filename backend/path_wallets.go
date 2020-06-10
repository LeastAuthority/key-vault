package backend

import (
	"context"
	vault "github.com/bloxapp/KeyVault"
	enc "github.com/bloxapp/KeyVault/encryptors"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func walletsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "wallets/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathWalletsList,
			},
			HelpSynopsis: "List all the Ethereum 2.0 wallets at a path",
			HelpDescription: `
			All the Ethereum 2.0 wallets will be listed.
			`,
		},
		&framework.Path{
			Pattern:      "wallets/" + framework.GenericNameRegex("name"),
			HelpSynopsis: "Create an Ethereum 2.0 wallet",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				//logical.ReadOperation:   b.pathAccountsRead,
				logical.CreateOperation: b.pathWalletCreate,
				//logical.UpdateOperation: b.pathAccountUpdate,
				//logical.DeleteOperation: b.pathAccountsDelete,
			},
		},
		&framework.Path{
			Pattern:      "wallets/" + framework.GenericNameRegex("name") + "/accounts/" + framework.GenericNameRegex("account"),
			HelpSynopsis: "Create an Ethereum 2.0 account",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
				"account": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				//logical.ReadOperation:   b.pathAccountsRead,
				logical.CreateOperation: b.pathWalletsAccountCreate,
				//logical.UpdateOperation: b.pathAccountUpdate,
				//logical.DeleteOperation: b.pathAccountsDelete,
			},
		},
		&framework.Path{
			Pattern:      "wallets/" + framework.GenericNameRegex("name") + "/accounts/",
			HelpSynopsis: "List wallet accounts",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.pathWalletAccountsList,
			},
		},
	}
	//return &framework.Path{
	//	Pattern: "wallets/?",
	//	Callbacks: map[logical.Operation]framework.OperationFunc{
	//		logical.ListOperation:   b.listWallets,
	//		logical.CreateOperation: b.createWallet,
	//	},
	//	HelpSynopsis: "List all the Ethereum 2.0 wallets maintained by the plugin backend and create new wallets.",
	//	HelpDescription: `
	//
	//LIST - list all wallets
    //POST - create a new wallet
	//
    //`,
	//	Fields: map[string]*framework.FieldSchema{
	//		"walletName": &framework.FieldSchema{
	//			Type:        framework.TypeString,
	//			Description: "Wallet name",
	//			Default:     "",
	//		},
	//	},
	//	ExistenceCheck: b.pathExistenceCheck,
	//}
}

func (b *backend) pathWalletCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("name").(string)
	options := vault.WalletOptions{}
	options.SetEncryptor(enc.NewPlainTextEncryptor())
	options.SetWalletPassword("")
	options.SetWalletName(walletName)
	options.SetStore(store.NewHashicorpVaultStore(req.Storage, ctx))
	vlt, err := vault.NewKeyVault(&options)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"id":   vlt.Wallet.ID().String(),
			"name": vlt.Wallet.Name(),
		},
	}, nil
}

func (b *backend) pathWalletsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	storeInstance := store.NewHashicorpVaultStore(req.Storage, ctx)
	wallets := map[string]bool{}
	for w := range storeInstance.RetrieveWallets() {
		wallets[string(w)] = true
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"wallets": wallets,
		},
	}, nil
}

func (b *backend) pathWalletsAccountCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("name").(string)
	accountName := data.Get("account").(string)
	options := vault.WalletOptions{}
	options.SetEncryptor(enc.NewPlainTextEncryptor())
	options.SetWalletName(walletName)
	options.SetStore(store.NewHashicorpVaultStore(req.Storage, ctx))
	vlt, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil,err
	}

	err = vlt.Wallet.Unlock([]byte(""))
	if err != nil {
		return nil,err
	}

	account, err := vlt.Wallet.CreateAccount(accountName, []byte(""))

	return &logical.Response{
		Data: map[string]interface{}{
			"walletName": vlt.Wallet.Name(),
			"accountName": account.Name(),
			"path": account.Path(),
		},
	}, nil
}

func (b *backend) pathWalletAccountsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("name").(string)
	storeInstance := store.NewHashicorpVaultStore(req.Storage, ctx)

	options := vault.WalletOptions{}
	options.SetEncryptor(enc.NewPlainTextEncryptor())
	options.SetWalletName(walletName)
	options.SetStore(storeInstance)
	vlt, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil,err
	}
	err = vlt.Wallet.Unlock([]byte(""))
	if err != nil {
		return nil,err
	}

	accounts := map[string]bool{}
	for w := range storeInstance.RetrieveAccounts(vlt.Wallet.ID()) {
		accounts[string(w)] = true
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"accounts": accounts,
		},
	}, nil
}