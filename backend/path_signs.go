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
	"github.com/pkg/errors"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

func signsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name") + "/accounts/" + framework.GenericNameRegex("account_name") + "/sign-attestation",
			HelpSynopsis:    "Sign attestation",
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
				logical.CreateOperation: b.pathWalletsAccountSignAttestation,
			},
		},
		&framework.Path{
			Pattern:         "wallets/" + framework.GenericNameRegex("wallet_name") + "/accounts/" + framework.GenericNameRegex("account_name") + "/sign-proposal",
			HelpSynopsis:    "Sign proposal",
			HelpDescription: `Sign proposal`,
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
					Default:     0,
				},
				"proposerIndex": &framework.FieldSchema{
					Type:        framework.TypeInt,
					Description: "Data ProposerIndex",
					Default:     0,
				},
				"parentRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data ParentRoot",
					Default:     "",
				},
				"stateRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data StateRoot",
					Default:     "",
				},
				"bodyRoot": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data BodyRoot",
					Default:     "",
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathWalletsAccountSignProposal,
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

	portfolio, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := portfolio.WalletByName(walletName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet by name")
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
		return nil, errors.Wrap(err, "failed to sign data")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": res.GetSignature(),
		},
	}, nil
}

func (b *backend) pathWalletsAccountSignProposal(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("wallet_name").(string)
	accountName := data.Get("account_name").(string)
	domain := data.Get("domain").(string)
	slot := data.Get("slot").(int)
	proposerIndex := data.Get("proposerIndex").(int)
	parentRoot := data.Get("parentRoot").(string)
	stateRoot := data.Get("stateRoot").(string)
	bodyRoot := data.Get("bodyRoot").(string)

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

	proposalRequest := &pb.SignBeaconProposalRequest{
		Id:     &pb.SignBeaconProposalRequest_Account{Account: accountName},
		Domain: ignoreError(hex.DecodeString(domain)).([]byte),
		Data: &pb.BeaconBlockHeader{
			Slot:          uint64(slot),
			ProposerIndex: uint64(proposerIndex),
			ParentRoot:    ignoreError(hex.DecodeString(parentRoot)).([]byte),
			StateRoot:     ignoreError(hex.DecodeString(stateRoot)).([]byte),
			BodyRoot:      ignoreError(hex.DecodeString(bodyRoot)).([]byte),
		},
	}

	protector := slashing_protection.NewNormalProtection(storage)
	signer := validator_signer.NewSimpleSigner(wallet, protector)
	res, err := signer.SignBeaconProposal(proposalRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign data")
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
