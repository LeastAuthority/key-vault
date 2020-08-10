package backend

import (
	"context"
	"encoding/hex"

	vault "github.com/bloxapp/KeyVault"
	"github.com/bloxapp/KeyVault/slashing_protection"
	"github.com/bloxapp/KeyVault/validator_signer"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/pkg/errors"
	v1 "github.com/wealdtech/eth2-signer-api/pb/v1"

	"github.com/bloxapp/vault-plugin-secrets-eth2.0/backend/store"
)

func signsPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern:         "accounts/sign-attestation",
			HelpSynopsis:    "Sign attestation",
			HelpDescription: `Sign attestation`,
			Fields: map[string]*framework.FieldSchema{
				"public_key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Public key of the account",
					Default:     "",
				},
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
				"useFakeSigner": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "True if the fake signer should be used",
					Default:     false,
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSignAttestation,
			},
		},
		&framework.Path{
			Pattern:         "accounts/sign-proposal",
			HelpSynopsis:    "Sign proposal",
			HelpDescription: `Sign proposal`,
			Fields: map[string]*framework.FieldSchema{
				"public_key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Public key of the account",
					Default:     "",
				},
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
				"useFakeSigner": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "True if the fake signer should be used",
					Default:     false,
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSignProposal,
			},
		},
		&framework.Path{
			Pattern:         "accounts/sign-aggregation",
			HelpSynopsis:    "Sign aggregation",
			HelpDescription: `Sign aggregation`,
			Fields: map[string]*framework.FieldSchema{
				"public_key": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Public key of the account",
					Default:     "",
				},
				"domain": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Domain",
					Default:     "",
				},
				"dataToSign": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Data to sign",
					Default:     "",
				},
				"useFakeSigner": &framework.FieldSchema{
					Type:        framework.TypeBool,
					Description: "True if the fake signer should be used",
					Default:     false,
				},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.pathSignAggregation,
			},
		},
	}
}

func (b *backend) pathSignAttestation(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)

	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet")
	}

	account, err := wallet.AccountByPublicKey(data.Get("public_key").(string))
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve account")
	}

	// try to lock signature lock, if it fails return error
	lock := DBLock{storage: req.Storage, id: account.ID()}
	err = lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.UnLock()

	publicKey := data.Get("public_key").(string)
	domain := data.Get("domain").(string)
	slot := data.Get("slot").(int)
	committeeIndex := data.Get("committeeIndex").(int)
	beaconBlockRoot := data.Get("beaconBlockRoot").(string)
	sourceEpoch := data.Get("sourceEpoch").(int)
	sourceRoot := data.Get("sourceRoot").(string)
	targetEpoch := data.Get("targetEpoch").(int)
	targetRoot := data.Get("targetRoot").(string)

	// Decode public key
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode public key")
	}

	// Decode domain
	domainBytes, err := hex.DecodeString(domain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode domain")
	}

	// Decode beacon block root
	beaconBlockRootBytes, err := hex.DecodeString(beaconBlockRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode beacon block root")
	}

	// Decode source root
	sourceRootBytes, err := hex.DecodeString(sourceRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode source root")
	}

	// Decode target root
	targetRootBytes, err := hex.DecodeString(targetRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode target root")
	}

	protector := slashing_protection.NewNormalProtection(storage)
	var signer validator_signer.ValidatorSigner = validator_signer.NewSimpleSigner(wallet, protector)

	res, err := signer.SignBeaconAttestation(&v1.SignBeaconAttestationRequest{
		Id:     &v1.SignBeaconAttestationRequest_PublicKey{PublicKey: publicKeyBytes},
		Domain: domainBytes,
		Data: &v1.AttestationData{
			Slot:            uint64(slot),
			CommitteeIndex:  uint64(committeeIndex),
			BeaconBlockRoot: beaconBlockRootBytes,
			Source: &v1.Checkpoint{
				Epoch: uint64(sourceEpoch),
				Root:  sourceRootBytes,
			},
			Target: &v1.Checkpoint{
				Epoch: uint64(targetEpoch),
				Root:  targetRootBytes,
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign attestation")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hex.EncodeToString(res.GetSignature()),
		},
	}, nil
}

func (b *backend) pathSignProposal(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)

	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet by name")
	}

	account, err := wallet.AccountByPublicKey(data.Get("public_key").(string))
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve account")
	}

	// try to lock signature lock, if it fails return error
	lock := DBLock{storage: req.Storage, id: account.ID()}
	err = lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.UnLock()

	publicKey := data.Get("public_key").(string)
	domain := data.Get("domain").(string)
	slot := data.Get("slot").(int)
	proposerIndex := data.Get("proposerIndex").(int)
	parentRoot := data.Get("parentRoot").(string)
	stateRoot := data.Get("stateRoot").(string)
	bodyRoot := data.Get("bodyRoot").(string)

	// Decode public key
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode public key")
	}

	// Decode domain
	domainBytes, err := hex.DecodeString(domain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode domain")
	}

	// Decode parent root
	parentRootBytes, err := hex.DecodeString(parentRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode parent root")
	}

	// Decode state root
	stateRootBytes, err := hex.DecodeString(stateRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode state root")
	}

	// Decode body root
	bodyRootBytes, err := hex.DecodeString(bodyRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode body root")
	}

	proposalRequest := &v1.SignBeaconProposalRequest{
		Id:     &v1.SignBeaconProposalRequest_PublicKey{PublicKey: publicKeyBytes},
		Domain: domainBytes,
		Data: &v1.BeaconBlockHeader{
			Slot:          uint64(slot),
			ProposerIndex: uint64(proposerIndex),
			ParentRoot:    parentRootBytes,
			StateRoot:     stateRootBytes,
			BodyRoot:      bodyRootBytes,
		},
	}

	protector := slashing_protection.NewNormalProtection(storage)
	var signer validator_signer.ValidatorSigner = validator_signer.NewSimpleSigner(wallet, protector)

	res, err := signer.SignBeaconProposal(proposalRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign data")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hex.EncodeToString(res.GetSignature()),
		},
	}, nil
}

func (b *backend) pathSignAggregation(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// bring up KeyVault and wallet
	storage := store.NewHashicorpVaultStore(ctx, req.Storage)
	options := vault.KeyVaultOptions{}
	options.SetStorage(storage)

	kv, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open key vault")
	}

	wallet, err := kv.Wallet()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve wallet by name")
	}

	// try to lock signature lock, if it fails return error
	lock := DBLock{storage: req.Storage, id: wallet.ID()}
	err = lock.Lock()
	if err != nil {
		return nil, err
	}
	defer lock.UnLock()

	publicKey := data.Get("public_key").(string)
	domain := data.Get("domain").(string)
	dataToSign := data.Get("dataToSign").(string)

	// Decode public key
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode public key")
	}

	// Decode domain
	domainBytes, err := hex.DecodeString(domain)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode domain")
	}

	// Decode data to sign
	dataToSignBytes, err := hex.DecodeString(dataToSign)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HEX decode data to sign")
	}

	proposalRequest := &v1.SignRequest{
		Id:     &v1.SignRequest_PublicKey{PublicKey: publicKeyBytes},
		Domain: domainBytes,
		Data:   dataToSignBytes,
	}

	protector := slashing_protection.NewNormalProtection(storage)
	var signer validator_signer.ValidatorSigner = validator_signer.NewSimpleSigner(wallet, protector)

	res, err := signer.Sign(proposalRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign data")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hex.EncodeToString(res.GetSignature()),
		},
	}, nil
}
