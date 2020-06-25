package fakesigner

import (
	"github.com/bloxapp/KeyVault/validator_signer"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/shared/keystore"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
)

const key = `{
    	"publickey": "b6be1db4b56ed5e51da58fa1427fb41ea2886170d6563f6f7a49a2632e048a72f3c8476dac93ffe97762e50dea6ae9ae",
    	"crypto": {
    		"cipher": "aes-128-ctr",
    		"ciphertext": "b53442e082cdfbed1b78f73f69218f25cd3546c515ea58d3d7219dc9dec214e6",
    		"cipherparams": {
    			"iv": "bf683d060daedc8ed25ca501658af025"
    		},
    		"kdf": "scrypt",
    		"kdfparams": {
    			"dklen": 32,
    			"n": 262144,
    			"p": 1,
    			"r": 8,
    			"salt": "03b0a06f5a8192248c2698e8400f12ed6ac241b33dff85853fd58631dcf365a1"
    		},
    		"mac": "8eb8e3078669076299b83e81fca09ce2da68bce3e4f27756f99619d975cb740c"
    	},
    	"id": "6026441e-0e10-4b3e-a663-ae0b449ae09a"
    }`

// fakeSigner implements validator_signer.ValidatorSigner.
// Uses hardcoded keys to sign data.
type fakeSigner struct {
	key *keystore.Key
}

// New is the constructor of fakeSigner.
func New() (validator_signer.ValidatorSigner, error) {
	key, err := keystore.DecryptKey([]byte(key), "changeit")
	if err != nil {
		return nil, err
	}

	return &fakeSigner{
		key: key,
	}, nil
}

// ListAccounts implements validator_signer.ValidatorSigner interface.
func (fs *fakeSigner) ListAccounts(req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	panic("not implemented")
}

// SignBeaconProposal implements validator_signer.ValidatorSigner interface.
func (fs *fakeSigner) SignBeaconProposal(req *pb.SignBeaconProposalRequest) (*pb.SignResponse, error) {
	preparedData, err := validator_signer.PrepareProposalReqForSigning(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare data for signing")
	}

	sig := fs.key.SecretKey.Sign(preparedData[:])

	return &pb.SignResponse{
		State:     pb.ResponseState_SUCCEEDED,
		Signature: sig.Marshal(),
	}, nil
}

// SignBeaconAttestation implements validator_signer.ValidatorSigner interface.
func (fs *fakeSigner) SignBeaconAttestation(req *pb.SignBeaconAttestationRequest) (*pb.SignResponse, error) {
	preparedData, err := validator_signer.PrepareAttestationReqForSigning(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare data for signing")
	}

	sig := fs.key.SecretKey.Sign(preparedData[:])

	return &pb.SignResponse{
		State:     pb.ResponseState_SUCCEEDED,
		Signature: sig.Marshal(),
	}, nil
}

// Sign implements validator_signer.ValidatorSigner interface.
func (fs *fakeSigner) Sign(req *pb.SignRequest) (*pb.SignResponse, error) {
	preparedData, err := validator_signer.PrepareReqForSigning(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare data for signing")
	}

	sig := fs.key.SecretKey.Sign(preparedData[:])

	return &pb.SignResponse{
		State:     pb.ResponseState_SUCCEEDED,
		Signature: sig.Marshal(),
	}, nil
}
