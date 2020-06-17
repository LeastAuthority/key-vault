package backend

import (
	"context"
	"encoding/hex"

	vault "github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/slashing_protection"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/bloxapp/KeyVault/validator_signer"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func signsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name") + "/accounts/" + framework.GenericNameRegex("account_name") + "/sign-attestation",
			HelpSynopsis:    "Sign",
			HelpDescription: `Sign attestation`,
			Fields: map[string]*framework.FieldSchema{
				"wallet_name":  &framework.FieldSchema{Type: framework.TypeString},
				"account_name": &framework.FieldSchema{Type: framework.TypeString},
				"domain": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Domain",
					Default:     "",
				},
				"slot": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data Slot",
					Default:     1,
				},
				"committeeIndex": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data CommitteeIndex",
					Default:     5,
				},
				"beaconBlockRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data BeaconBlockRoot",
					Default:     "test",
				},
				"sourceEpoch": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data Source Epoch",
					Default:     6,
				},
				"sourceRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data Source Root",
					Default:     "test2",
				},
				"targetEpoch": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data Target Epoch",
					Default:     45,
				},
				"targetRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data Target Root",
					Default:     "test3",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletsAccountSignAttestation,
			},
		},
	}
}

func (b *backend) pathWalletsAccountSignAttestation(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
		return nil, err
	}
	//wallet, err := portfolio.CreateWallet(walletName)
	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, err
	}

	//_, err = wallet.CreateValidatorAccount(accountName)
	//if err != nil {
	//	return nil, err
	//}

	protector := slashing_protection.NewNormalProtection(storage)
	signer := validator_signer.NewSimpleSigner(wallet, protector)
	res, err := signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
		Id:     &pb.SignBeaconAttestationRequest_Account{Account: accountName},
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
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": res.GetSignature(),
		},
	}, nil
}

func ignoreError(val interface{}, err error) interface{} {
	return val
}
