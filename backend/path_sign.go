package backend

import (
	"context"
	vault "github.com/bloxapp/KeyVault"
	store "github.com/bloxapp/KeyVault/stores/hashicorp"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	keystore "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

func pathSign(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "sign/?",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.signTx,
		},
		//Operations: map[logical.Operation]framework.OperationHandler{
		//	logical.CreateOperation: b.signTx
		//},
		HelpSynopsis: "Sign a provided transaction object.",
		HelpDescription: `

    Sign a transaction.

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
			"domain": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Domain",
				Default:     "",
			},
			"slot": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Data Slot",
				Default:     "",
			},
			"committeeIndex": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Data CommitteeIndex",
				Default:     "",
			},
			"beaconBlockRoot": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Data BeaconBlockRoot",
				Default:     "",
			},
			"sourceEpoch": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Data Source Epoch",
				Default:     "",
			},
			"sourceRoot": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Data Source Root",
				Default:     "",
			},
			"targetEpoch": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Data Target Epoch",
				Default:     "",
			},
			"targetRoot": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Data Target Root",
				Default:     "",
			},
		},
		ExistenceCheck: b.pathExistenceCheck,
	}
}

func (b *backend) signTx(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	walletName := data.Get("walletName").(string)
	//accountName := data.Get("accountName").(string)
	//domain := data.Get("domain").(string)
	//slot := data.Get("slot").(uint64)
	//committeeIndex := data.Get("committeeIndex").(uint64)
	//beaconBlockRoot := data.Get("beaconBlockRoot").(string)
	//sourceEpoch := data.Get("sourceEpoch").(uint64)
	//sourceRoot := data.Get("sourceRoot").(string)
	//targetEpoch := data.Get("targetEpoch").(uint64)
	//targetRoot := data.Get("targetRoot").(string)
	options := vault.WalletOptions{}
	options.SetEncryptor(keystore.New())
	options.SetWalletName(walletName)
	options.SetStore(store.NewHashicorpVaultStore(req.Storage, ctx))
	options.EnableSimpleSigner(true)
	vlt, err := vault.OpenKeyVault(&options)
	if err != nil {
		return nil,err
	}

	err = vlt.Wallet.Unlock([]byte(""))
	if err != nil {
		return nil,err
	}
	//
	//res, err := vlt.Signer.SignBeaconAttestation(&pb.SignBeaconAttestationRequest{
	//	Id:     &pb.SignBeaconAttestationRequest_Account{Account: accountName},
	//	Domain: ignoreError(hex.DecodeString(domain)).([]byte),
	//	Data: &pb.AttestationData{
	//		Slot:            slot,
	//		CommitteeIndex:  committeeIndex,
	//		BeaconBlockRoot: ignoreError(hex.DecodeString(beaconBlockRoot)).([]byte),
	//		Source: &pb.Checkpoint{
	//			Epoch: sourceEpoch,
	//			Root:  ignoreError(hex.DecodeString(sourceRoot)).([]byte),
	//		},
	//		Target: &pb.Checkpoint{
	//			Epoch: targetEpoch,
	//			Root:  ignoreError(hex.DecodeString(targetRoot)).([]byte),
	//		},
	//	},
	//})
	//if err != nil {
	//	return nil,err
	//}

	return &logical.Response{
		Data: map[string]interface{}{
			"data": vlt,
		},
	}, nil
}

func ignoreError(val interface{}, err error)interface{} {
	return val
}