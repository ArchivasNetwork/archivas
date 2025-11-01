package consensus

import (
	"fmt"
	"time"

	"github.com/iljanemesis/archivas/pospace"
)

// Consensus implements the Archivas consensus protocol
type Consensus struct {
	DifficultyTarget uint64
	TargetBlockTime  time.Duration
}

const (
	// InitialDifficulty in new domain (1e9-1e12)
	// Start at 15M for ~20s blocks (adjusted from testing)
	InitialDifficulty = 15_000_000
)

// NewConsensus creates a new consensus instance
func NewConsensus() *Consensus {
	return &Consensus{
		DifficultyTarget: InitialDifficulty,
		TargetBlockTime:  20 * time.Second,
	}
}

// VerifyProofOfSpace verifies a Proof-of-Space for a block
func (c *Consensus) VerifyProofOfSpace(proof *pospace.Proof, challenge [32]byte) error {
	// Verify the proof is valid
	if !pospace.VerifyProof(proof, challenge, c.DifficultyTarget) {
		return fmt.Errorf("invalid proof of space: quality %d does not meet difficulty %d", proof.Quality, c.DifficultyTarget)
	}

	return nil
}

// UpdateDifficulty updates the difficulty based on actual vs target block time
// Uses inverted proportional control with damping
func (c *Consensus) UpdateDifficulty(actualBlockTime time.Duration) {
	target := c.TargetBlockTime.Seconds()
	actual := actualBlockTime.Seconds()
	
	if actual == 0 || actual < 0.1 {
		return
	}
	
	// Inverted: scale = target / observed
	// If blocks too fast (2s vs 20s target), scale = 10 (harder!)
	scale := target / actual
	
	// Damping factor (alpha = 0.5)
	alpha := 0.5
	adjustedScale := 1.0 + alpha*(scale-1.0)
	
	// Apply to difficulty
	newDiff := uint64(float64(c.DifficultyTarget) * adjustedScale)
	
	// Clamp to domain [1e8, 1e12]
	minDiff := uint64(100_000_000)   // 100M
	maxDiff := uint64(1_000_000_000_000) // 1T (QMAX)
	
	if newDiff < minDiff {
		newDiff = minDiff
	}
	if newDiff > maxDiff {
		newDiff = maxDiff
	}
	
	c.DifficultyTarget = newDiff
}
