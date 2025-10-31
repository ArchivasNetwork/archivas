package registry

import (
	"encoding/json"
	"time"
)

// ValidatorInfo represents a registered validator
type ValidatorInfo struct {
	PubKey       []byte `json:"pubKey"`       // Node signing key
	Address      string `json:"address"`      // arcv1... address
	Moniker      string `json:"moniker"`      // Display name (max 50 chars)
	Endpoint     string `json:"endpoint"`     // P2P endpoint
	Version      string `json:"version"`      // e.g., "v0.7.0-alpha"
	Commit       string `json:"commit"`       // Git short SHA
	Role         string `json:"role"`         // "validator", "full", "archivist"
	Capacity     uint64 `json:"capacity"`     // Self-declared capacity points
	RegisteredAt int64  `json:"registeredAt"` // Unix timestamp
	LastSeen     int64  `json:"lastSeen"`     // Last heartbeat timestamp
	Signature    []byte `json:"signature"`    // Signature over canonical JSON
}

// HealthSample represents a health check sample
type HealthSample struct {
	Timestamp      int64   `json:"timestamp"`      // Unix timestamp
	Height         uint64  `json:"height"`         // Current block height
	BlocksPerHour  float64 `json:"blocksPerHour"`  // Block production rate
	PeerCount      int     `json:"peerCount"`      // Connected peers
	CPUPercent     float64 `json:"cpuPercent"`     // CPU usage (-1 if unavailable)
	MemRSSMB       float64 `json:"memRSSMB"`       // Memory RSS in MB (-1 if unavailable)
	LatencyMS      float64 `json:"latencyMS"`      // Median latency to peers
}

// Proposal represents a governance proposal
type Proposal struct {
	ID          string          `json:"id"`          // "P-000001"
	Kind        string          `json:"kind"`        // "PARAM_CHANGE", "TEXT"
	Title       string          `json:"title"`       // Max 200 chars
	Description string          `json:"description"` // Max 2000 chars
	Payload     json.RawMessage `json:"payload"`     // Param changes for PARAM_CHANGE
	CreatorPub  []byte          `json:"creatorPub"`  // Proposer's public key
	CreatedAt   int64           `json:"createdAt"`   // Unix timestamp
	Status      string          `json:"status"`      // "OPEN", "CLOSED", "REJECTED", "PASSED"
	StartsAt    int64           `json:"startsAt"`    // Voting start (0 = immediate)
	EndsAt      int64           `json:"endsAt"`      // Voting end
	MinYesPct   float64         `json:"minYesPct"`   // e.g., 0.5 (50%)
	QuorumPct   float64         `json:"quorumPct"`   // e.g., 0.4 (40%)
	Signature   []byte          `json:"signature"`   // Creator's signature
}

// Vote represents a governance vote
type Vote struct {
	ProposalID string  `json:"proposalId"`
	VoterPub   []byte  `json:"voterPub"`   // Voter's public key
	Choice     string  `json:"choice"`     // "YES", "NO", "ABSTAIN"
	Weight     float64 `json:"weight"`     // Computed voting weight
	Timestamp  int64   `json:"timestamp"`  // Unix timestamp
	Signature  []byte  `json:"signature"`  // Voter's signature
}

// TallyResult represents vote tally results
type TallyResult struct {
	ProposalID    string  `json:"proposalId"`
	TotalWeight   float64 `json:"totalWeight"`   // Total eligible weight
	YesWeight     float64 `json:"yesWeight"`     // Yes votes weight
	NoWeight      float64 `json:"noWeight"`      // No votes weight
	AbstainWeight float64 `json:"abstainWeight"` // Abstain weight
	Quorum        float64 `json:"quorum"`        // Actual quorum (voted/total)
	Passed        bool    `json:"passed"`        // Did it pass?
	ComputedAt    int64   `json:"computedAt"`    // When tally was computed
}

// ParamSet represents governable parameters
type ParamSet struct {
	MinCompatibleVersion   string `json:"minCompatibleVersion"`
	BlockPropagateTargetMS int    `json:"blockPropagateTargetMS"`
	MaxPeers               int    `json:"maxPeers"`
	GossipTxRatePerSec     int    `json:"gossipTxRatePerSec"`
	GossipBlockRatePerSec  int    `json:"gossipBlockRatePerSec"`
	ReorgDepthMax          int    `json:"reorgDepthMax"`
	APIRateLimitPerMin     int    `json:"apiRateLimitPerMin"`
}

// DefaultParams returns default governable parameters
func DefaultParams() ParamSet {
	return ParamSet{
		MinCompatibleVersion:   "v0.6.0-alpha",
		BlockPropagateTargetMS: 2000,
		MaxPeers:               20,
		GossipTxRatePerSec:     100,
		GossipBlockRatePerSec:  10,
		ReorgDepthMax:          100,
		APIRateLimitPerMin:     60,
	}
}

// UptimeScore calculates uptime score from heartbeat samples
func UptimeScore(samples []HealthSample, windowDuration time.Duration) float64 {
	if len(samples) == 0 {
		return 0.0
	}

	// Expected: 1 sample per minute for 24 hours = 1440 samples
	expectedSamples := windowDuration.Minutes()
	actualSamples := float64(len(samples))

	score := actualSamples / expectedSamples
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// CalculateVotingWeight computes voting weight based on uptime and capacity
func CalculateVotingWeight(uptimeScore float64, capacity uint64) float64 {
	// Eligibility threshold
	if uptimeScore < 0.60 {
		return 0.0
	}

	// Base weight from uptime: 0.5 to 1.0
	baseWeight := 0.5 + 0.5*uptimeScore

	// Capacity factor (soft differentiation)
	// capacityFactor = max(1.0, log2(1 + Capacity/100))
	if capacity == 0 {
		capacity = 100
	}

	// Simplified: just use capacity as small multiplier
	capacityFactor := 1.0
	if capacity > 100 {
		capacityFactor = 1.0 + float64(capacity-100)/1000.0
	}
	if capacityFactor > 2.0 {
		capacityFactor = 2.0 // Cap at 2x
	}

	return baseWeight * capacityFactor
}

