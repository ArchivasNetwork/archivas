package rpc

// Shared RPC types used across farming, VDF, and wallet endpoints

// ChainTipResponse represents the current chain tip (for timelord)
type ChainTipResponse struct {
	BlockHash  [32]byte `json:"blockHash"`
	Height     uint64   `json:"height"`
	Difficulty uint64   `json:"difficulty"`
}

// VDFUpdateRequest represents a VDF update from timelord
type VDFUpdateRequest struct {
	Seed       []byte `json:"seed"`
	Iterations uint64 `json:"iterations"`
	Output     []byte `json:"output"`
}

// Note: BalanceResponse and SubmitTxResponse are defined in rpc.go

// ChallengeResponse represents the challenge info for farmers
type ChallengeResponse struct {
	Challenge  [32]byte `json:"challenge"`
	Difficulty uint64   `json:"difficulty"`
	Height     uint64   `json:"height"`
	VDF        *VDFInfo `json:"vdf,omitempty"` // Optional VDF info
}

// VDFInfo represents VDF state in challenge response
type VDFInfo struct {
	Seed       []byte `json:"seed"`
	Iterations uint64 `json:"iterations"`
	Output     []byte `json:"output"`
}

// BlockSubmission represents a block submission from a farmer (used by farming.go)
// Note: Proof will need type assertion to *pospace.Proof when used

