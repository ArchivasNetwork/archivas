package consensus

import (
	"fmt"

	"github.com/iljanemesis/archivas/pospace"
)

// Consensus implements the Archivas consensus protocol
type Consensus struct {
	DifficultyTarget uint64
}

// NewConsensus creates a new consensus instance
func NewConsensus() *Consensus {
	return &Consensus{
		DifficultyTarget: InitialDifficulty,
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

// UpdateDifficulty updates the difficulty based on recent block times
func (c *Consensus) UpdateDifficulty(recentBlockTimes []int64) {
	c.DifficultyTarget = CalculateDifficulty(recentBlockTimes, c.DifficultyTarget)
}

