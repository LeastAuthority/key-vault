package backend

import (
	"context"
	vault "github.com/bloxapp/KeyVault"
	enc "github.com/bloxapp/KeyVault/encryptors"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCreateAndListAccount(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "accounts/?",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			//logical.ListOperation:   b.listWallets,
			logical.UpdateOperation: b.createAccount,
		},
		HelpSynopsis: "List all the Ethereum 2.0 accounts maintained by the plugin backend and create new accounts.",
		HelpDescription: `

	LIST - list all accounts
    POST - create a new account

    `,
		Fields: map[string]*framework.FieldSchema{
			"walletName": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Wallet name",
				Default:     "",
			},
			"accountName": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Account name",
				Default:     "",
			},
		},
	}
}

func (b *backend) createAccount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("walletName").(string)
	accountName := data.Get("accountName").(string)
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
			"publicKey": account.PublicKey(),
			"path": account.Path(),
		},
	}, nil
}

//func getStorage() logical.Storage {
//	return &logical.InmemStorage{}
//}
//
//func getWalletStorage() wtypes.Store {
//	return store.NewHashicorpVaultStore(getStorage(), context.Background())
//}