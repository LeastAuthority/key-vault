package backend

import (
	"context"
	vault "github.com/bloxapp/KeyVault"
	enc "github.com/bloxapp/KeyVault/encryptors"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCreateAndListWallet(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "wallets/?",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			//logical.ListOperation:   b.listWallets,
			logical.UpdateOperation: b.createWallet,
		},
		HelpSynopsis: "List all the Ethereum 2.0 wallets maintained by the plugin backend and create new wallets.",
		HelpDescription: `

	LIST - list all wallets
    POST - create a new wallet

    `,
		Fields: map[string]*framework.FieldSchema{
			"walletName": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Wallet name",
				Default:     "",
			},
		},
	}
}

func (b *backend) createWallet(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("walletName").(string)
	options := vault.WalletOptions{}
	options.SetEncryptor(enc.NewPlainTextEncryptor())
	options.SetWalletPassword("")
	options.SetWalletName(walletName)
	options.SetStore(store.NewHashicorpVaultStore(req.Storage, ctx))
	vlt, err := vault.NewKeyVault(&options)
	if err != nil {
		return nil,err
	}

	//err = vlt.Wallet.Unlock([]byte(""))
	//if err != nil {
	//	return nil,err
	//}

	//account, err := vlt.Wallet.CreateAccount("account1", []byte(""))
	//vlt.Wallet.CreateAccount()

	return &logical.Response{
		Data: map[string]interface{}{
			"id": vlt.Wallet.ID().String(),
			"name": vlt.Wallet.Name(),
		},
	}, nil
}

func (b *backend) listWallets(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	options := vault.WalletOptions{}
	options.SetStore(store.NewHashicorpVaultStore(req.Storage, ctx))
	vlt, err := vault. NewKeyVault(&options)
	if err != nil {
		return nil,err
	}
	//vlt.Wallet.CreateAccount()

	return &logical.Response{
		Data: map[string]interface{}{
			"id": vlt.Wallet.ID().String(),
			"name": vlt.Wallet.Name(),
		},
	}, nil
}