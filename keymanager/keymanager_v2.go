package keymanager

import (
	"context"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	v2keymanager "github.com/prysmaticlabs/prysm/validator/keymanager/v2"
)

// To make sure VaultRemoteHTTPWallet implements v2keymanager.IKeymanager interface
var _ v2keymanager.IKeymanager = &VaultRemoteHTTPWallet{}

// Sign implements KeyManager-v2 interface.
func (km *VaultRemoteHTTPWallet) Sign(ctx context.Context, req *validatorpb.SignRequest) (bls.Signature, error) {
	if bytesutil.ToBytes48(req.GetPublicKey()) != km.pubKey {
		return nil, ErrNoSuchKey
	}

	domain := bytesutil.ToBytes32(req.GetSignatureDomain())
	switch data := req.GetObject().(type) {
	case *validatorpb.SignRequest_Block:
		return km.SignProposal(km.pubKey, domain, &ethpb.BeaconBlockHeader{
			Slot:          data.Block.GetSlot(),
			ProposerIndex: data.Block.GetProposerIndex(),
			StateRoot:     data.Block.GetStateRoot(),
			ParentRoot:    data.Block.GetParentRoot(),
			BodyRoot:      req.GetSigningRoot(),
		})
	case *validatorpb.SignRequest_AttestationData:
		return km.SignAttestation(km.pubKey, domain, data.AttestationData)
	case *validatorpb.SignRequest_Slot:
		return km.SignGeneric(km.pubKey, bytesutil.ToBytes32(req.GetSigningRoot()), domain)
	default:
		return nil, ErrUnsupportedSigning
	}
}

// FetchValidatingPublicKeys implements KeyManager-v2 interface.
func (km *VaultRemoteHTTPWallet) FetchValidatingPublicKeys(_ context.Context) ([][48]byte, error) {
	return [][48]byte{km.pubKey}, nil
}
