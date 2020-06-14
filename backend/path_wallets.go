package backend

import (
	"context"
	"encoding/hex"
	vault "github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/core"
	"github.com/bloxapp/KeyVault/slashing_protection"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/bloxapp/KeyVault/validator_signer"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
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
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name"),
			HelpSynopsis:    "Create an Ethereum 2.0 wallet",
			HelpDescription: ``,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletCreate,
			},
		},
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
			Pattern:      "wallets/" + framework.GenericNameRegex("wallet_name") + "/accounts/" + framework.GenericNameRegex("account_name") + "/sign",
			HelpSynopsis: "Sign",
			HelpDescription: ` Sign attestation`,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name": &framework.FieldSchema{Type: framework.TypeString},
				"account_name": &framework.FieldSchema{Type: framework.TypeString},
				"domain": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Domain",
					Default:     "",
				},
				"slot": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data Slot",
					Default:     0,
				},
				"committeeIndex": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data CommitteeIndex",
					Default:     0,
				},
				"beaconBlockRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data BeaconBlockRoot",
					Default:     "",
				},
				"sourceEpoch": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data Source Epoch",
					Default:     0,
				},
				"sourceRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data Source Root",
					Default:     "",
				},
				"targetEpoch": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data Target Epoch",
					Default:     0,
				},
				"targetRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data Target Root",
					Default:     "",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletsAccountSign,
			},
		},
	}
}

func (b *backend) pathWalletCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)
	portfolio, err := vault.NewKeyVault(&options)
	if err != nil {
		return nil, err
	}
	wallet, err := portfolio.CreateWallet(walletName)
	if err != nil {
		return nil, err
	}

	//portfolio, err = vault.OpenKeyVault(&options)
	//if err != nil {
	//	return nil, err
	//}
	//
	//wallet, err = portfolio.CreateWallet(walletName)
	//if err != nil {
	//	return nil, err
	//}

	return &logical.Response{
		Data: map[string]interface{}{
			"wallet": wallet,
		},
	}, nil
}

func (b *backend) pathWalletsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)
	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, err
	}
	wallets := make([]core.Wallet, 0)
	for w := range portfolio.Wallets() {
		wallets = append(wallets, w)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"wallets": wallets,
		},
	}, nil
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
			"account":  account,
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

func (b *backend) pathWalletsAccountSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	accountName := data.Get("account_name").(string)
	domain := data.Get("domain").(string)
	slot := data.Get("slot").(int)
	committeeIndex := data.Get("committeeIndex").(int)
	beaconBlockRoot := data.Get("beaconBlockRoot").(string)
	sourceEpoch := data.Get("sourceEpoch").(int)
	sourceRoot := data.Get("sourceRoot").(string)
	targetEpoch := data.Get("targetEpoch").(int)
	targetRoot := data.Get("targetRoot").(string)
	storage := store.NewHashicorpVaultStore(req.Storage, ctx)
	options := vault.PortfolioOptions{}
	options.SetStorage(storage)
	//portfolio, err := vault.NewKeyVault(&options)
	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil,err
	}
	//wallet, err := portfolio.CreateWallet(walletName)
	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, err
	}
	//account, err := wallet.CreateValidatorAccount(accountName)
	account, err := wallet.AccountByName(accountName)
	if err != nil {
		return nil,err
	}

	protector := slashing_protection.NewNormalProtection(storage)
	signer := validator_signer.NewSimpleSigner(wallet, protector)
	res, err := signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
		Id:     &pb.SignBeaconAttestationRequest_Account{Account: account.Name()},
		Domain: ignoreError(hex.DecodeString(domain)).([]byte),
		Data: &pb.AttestationData{
			Slot:            uint64(slot),
			CommitteeIndex:  uint64(committeeIndex),
			BeaconBlockRoot: ignoreError(hex.DecodeString(beaconBlockRoot)).([]byte),
			Source: &pb.Checkpoint{
				Epoch: uint64(sourceEpoch),
				Root:  ignoreError(hex.DecodeString(sourceRoot)).([]byte),
			},
			Target: &pb.Checkpoint{
				Epoch: uint64(targetEpoch),
				Root:  ignoreError(hex.DecodeString(targetRoot)).([]byte),
			},
		},
	})
	if err != nil {
		return nil,err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": res.Signature,
		},
	}, nil
}

func ignoreError(val interface{}, err error)interface{} {
	return val
}
