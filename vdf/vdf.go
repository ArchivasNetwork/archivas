package vdf

import (
	"crypto/sha256"
)

// StepHash returns H(input) using SHA-256
func StepHash(input []byte) []byte {
	h := sha256.Sum256(input)
	out := make([]byte, 32)
	copy(out, h[:])
	return out
}

// ComputeSequential runs the VDF for `iterations` steps starting from `seed`.
// Returns the final output and a slice of checkpoints (every `checkpointStep` steps).
func ComputeSequential(seed []byte, iterations uint64, checkpointStep uint64) (final []byte, checkpoints [][]byte) {
	cur := make([]byte, len(seed))
	copy(cur, seed)

	checkpoints = make([][]byte, 0)

	for i := uint64(1); i <= iterations; i++ {
		cur = StepHash(cur)
		if checkpointStep > 0 && (i%checkpointStep == 0) {
			snap := make([]byte, len(cur))
			copy(snap, cur)
			checkpoints = append(checkpoints, snap)
		}
	}
	final = cur
	return
}

// VerifySequential recomputes from seed for `iterations`
// and checks that the claimed final output matches.
func VerifySequential(seed []byte, iterations uint64, claimedFinal []byte) bool {
	cur := make([]byte, len(seed))
	copy(cur, seed)
	for i := uint64(1); i <= iterations; i++ {
		cur = StepHash(cur)
	}
	// compare
	if len(cur) != len(claimedFinal) {
		return false
	}
	for i := range cur {
		if cur[i] != claimedFinal[i] {
			return false
		}
	}
	return true
}
