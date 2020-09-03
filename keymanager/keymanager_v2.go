package keymanager

import (
	"context"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	v2keymanager "github.com/prysmaticlabs/prysm/validator/keymanager/v2"
)

// To make sure KeyManagerV2 implements v2keymanager.IKeymanager interface
var _ v2keymanager.IKeymanager = &KeyManagerV2{}

// KeyManagerV2 implements prysm keymanager-v2 interface.
type KeyManagerV2 struct {
	km *KeyManager
}

// NewKeyManagerV2 is the constructor of NewKeyManagerV2.
func NewKeyManagerV2(km *KeyManager) *KeyManagerV2 {
	return &KeyManagerV2{
		km: km,
	}
}

// Sign implements KeyManager-v2 interface.
func (km *KeyManagerV2) Sign(ctx context.Context, req *validatorpb.SignRequest) (bls.Signature, error) {
	if bytesutil.ToBytes48(req.GetPublicKey()) != km.km.pubKey {
		return nil, ErrNoSuchKey
	}

	domain := bytesutil.ToBytes32(req.GetSignatureDomain())
	switch data := req.GetObject().(type) {
	case *validatorpb.SignRequest_Block:
		return km.km.SignProposal(km.km.pubKey, domain, &ethpb.BeaconBlockHeader{
			Slot:          data.Block.GetSlot(),
			ProposerIndex: data.Block.GetProposerIndex(),
			StateRoot:     data.Block.GetStateRoot(),
			ParentRoot:    data.Block.GetParentRoot(),
			BodyRoot:      req.GetSigningRoot(),
		})
	case *validatorpb.SignRequest_AttestationData:
		return km.km.SignAttestation(km.km.pubKey, domain, data.AttestationData)
	case *validatorpb.SignRequest_AggregateAttestationAndProof:
		return km.km.SignGeneric(km.km.pubKey, bytesutil.ToBytes32(req.GetSigningRoot()), domain)
	case *validatorpb.SignRequest_Slot:
		return km.km.SignGeneric(km.km.pubKey, bytesutil.ToBytes32(req.GetSigningRoot()), domain)
	default:
		return nil, ErrUnsupportedSigning
	}
}

// FetchValidatingPublicKeys implements KeyManager-v2 interface.
func (km *KeyManagerV2) FetchValidatingPublicKeys(_ context.Context) ([][48]byte, error) {
	return [][48]byte{km.km.pubKey}, nil
}
