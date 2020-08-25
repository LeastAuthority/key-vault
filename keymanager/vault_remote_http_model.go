package keymanager

// VaultSignAttestationRequest is the request body of vault sign attestation endpoint.
type VaultSignAttestationRequest struct {
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

// VaultSignProposalRequest is the request body of vault sign proposal endpoint.
type VaultSignProposalRequest struct {
	PubKey        string `json:"public_key"`
	Domain        string `json:"domain"`
	Slot          uint64 `json:"slot"`
	ProposerIndex uint64 `json:"proposerIndex"`
	ParentRoot    string `json:"parentRoot"`
	StateRoot     string `json:"stateRoot"`
	BodyRoot      string `json:"bodyRoot"`
}

// VaultSignAggregationRequest is the request body of vault sign aggregation endpoint.
type VaultSignAggregationRequest struct {
	PubKey     string `json:"public_key"`
	Domain     string `json:"domain"`
	DataToSign string `json:"dataToSign"`
}

// VaultSignResponse is the vault sign response model.
type VaultSignResponse struct {
	Data VaultSignatureModel `json:"data"`
}

// VaultSignatureModel represents vault signature model.
type VaultSignatureModel struct {
	Signature string `json:"signature"`
}
