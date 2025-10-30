package consensus

import (
	"fmt"
)

// ReorgDetector handles chain reorganization detection and execution
type ReorgDetector struct {
	// Configuration
	MaxReorgDepth int // Maximum depth we'll reorg (safety limit)
}

// NewReorgDetector creates a new reorganization detector
func NewReorgDetector() *ReorgDetector {
	return &ReorgDetector{
		MaxReorgDepth: 100, // Don't reorg more than 100 blocks deep
	}
}

// DetectReorg checks if a competing chain should replace the current chain
// Returns: needsReorg, forkHeight, error
func (r *ReorgDetector) DetectReorg(currentTipWork, newChainWork uint64, commonHeight uint64, currentHeight uint64) (bool, uint64, error) {
	// Calculate reorg depth
	reorgDepth := currentHeight - commonHeight

	// Safety check: don't reorg too deep
	if reorgDepth > uint64(r.MaxReorgDepth) {
		return false, 0, fmt.Errorf("reorg too deep: %d blocks (max: %d)", reorgDepth, r.MaxReorgDepth)
	}

	// Compare cumulative work
	if newChainWork > currentTipWork {
		// New chain has more work - reorg needed
		return true, commonHeight, nil
	}

	// Current chain has more work - no reorg
	return false, 0, nil
}

// ReorgInfo contains information about a chain reorganization
type ReorgInfo struct {
	ForkHeight    uint64 // Height where chains diverged
	OldTipHeight  uint64 // Old chain tip
	NewTipHeight  uint64 // New chain tip
	BlocksRemoved uint64 // Blocks rolled back
	BlocksAdded   uint64 // Blocks from new chain
}

// CalculateReorgInfo computes reorg statistics
func CalculateReorgInfo(forkHeight, oldTip, newTip uint64) ReorgInfo {
	return ReorgInfo{
		ForkHeight:    forkHeight,
		OldTipHeight:  oldTip,
		NewTipHeight:  newTip,
		BlocksRemoved: oldTip - forkHeight,
		BlocksAdded:   newTip - forkHeight,
	}
}
