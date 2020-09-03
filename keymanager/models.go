package keymanager

// SignAttestationRequest is the request body of vault sign attestation endpoint.
type SignAttestationRequest struct {
	PubKey          string `json:"public_key"`
	Domain          string `json:"domain"`
	Slot            uint64 `json:"slot"`
	CommitteeIndex  uint64 `json:"committeeIndex"`
	BeaconBlockRoot string `json:"beaconBlockRoot"`
	SourceEpoch     uint64 `json:"sourceEpoch"`
	SourceRoot      string `json:"sourceRoot"`
	TargetEpoch     uint64 `json:"targetEpoch"`
	TargetRoot      string `json:"targetRoot"`
}

// SignProposalRequest is the request body of vault sign proposal endpoint.
type SignProposalRequest struct {
	PubKey        string `json:"public_key"`
	Domain        string `json:"domain"`
	Slot          uint64 `json:"slot"`
	ProposerIndex uint64 `json:"proposerIndex"`
	ParentRoot    string `json:"parentRoot"`
	StateRoot     string `json:"stateRoot"`
	BodyRoot      string `json:"bodyRoot"`
}

// SignAggregationRequest is the request body of vault sign aggregation endpoint.
type SignAggregationRequest struct {
	PubKey     string `json:"public_key"`
	Domain     string `json:"domain"`
	DataToSign string `json:"dataToSign"`
}

// SignResponse is the vault sign response model.
type SignResponse struct {
	Data SignatureModel `json:"data"`
}

// SignatureModel represents vault signature model.
type SignatureModel struct {
	Signature string `json:"signature"`
}
