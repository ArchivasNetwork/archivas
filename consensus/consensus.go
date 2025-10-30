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

// NewConsensus creates a new consensus instance
func NewConsensus() *Consensus {
	return &Consensus{
		DifficultyTarget: 1125899906842624, // 2^50 initial difficulty
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
func (c *Consensus) UpdateDifficulty(actualBlockTime time.Duration) {
	// Simple adjustment for now (smoothing in difficulty.go)
	target := c.TargetBlockTime.Seconds()
	actual := actualBlockTime.Seconds()
	
	if actual == 0 {
		return
	}
	
	ratio := target / actual
	newDiff := uint64(float64(c.DifficultyTarget) * ratio)
	
	// Limit changes to 2x per adjustment
	if newDiff > c.DifficultyTarget*2 {
		newDiff = c.DifficultyTarget * 2
	}
	if newDiff < c.DifficultyTarget/2 {
		newDiff = c.DifficultyTarget / 2
	}
	
	// Minimum difficulty
	if newDiff < 1000000 {
		newDiff = 1000000
	}
	
	c.DifficultyTarget = newDiff
}

