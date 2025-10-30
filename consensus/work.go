package consensus

// CalculateWork computes the work contribution of a block based on its difficulty
// Higher difficulty = more work
func CalculateWork(difficulty uint64) uint64 {
	// Simple work calculation: work is proportional to difficulty
	// In production, this would be: work = 2^256 / difficulty
	// For now, we use difficulty directly as work
	return difficulty
}

// CompareChains returns true if chain A has more cumulative work than chain B
func CompareChains(workA, workB uint64) bool {
	return workA > workB
}

// ComputeCumulativeWork calculates total work from a chain of blocks
func ComputeCumulativeWork(blocks []BlockWork) uint64 {
	var total uint64
	for _, block := range blocks {
		total += CalculateWork(block.Difficulty)
	}
	return total
}

// BlockWork contains minimal block data for work calculation
type BlockWork struct {
	Height     uint64
	Difficulty uint64
}

