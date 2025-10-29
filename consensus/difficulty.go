package consensus

import (
	"github.com/iljanemesis/archivas/config"
)

const (
	// DifficultyAdjustmentWindow is the number of blocks to consider for difficulty adjustment
	DifficultyAdjustmentWindow = 10

	// InitialDifficulty is the starting difficulty (lower = harder)
	// This is a uint64 threshold - quality must be LESS than this to win
	InitialDifficulty = uint64(1 << 50) // ~1.1e15, adjust based on plot size for ~20s blocks
)

// CalculateDifficulty computes the difficulty target based on recent block times
func CalculateDifficulty(recentBlockTimes []int64, currentDifficulty uint64) uint64 {
	if len(recentBlockTimes) < 2 {
		return InitialDifficulty
	}

	// Calculate average time between blocks
	var totalTime int64
	for i := 1; i < len(recentBlockTimes); i++ {
		totalTime += recentBlockTimes[i] - recentBlockTimes[i-1]
	}
	avgBlockTime := totalTime / int64(len(recentBlockTimes)-1)

	targetBlockTime := int64(config.TargetBlockTimeSeconds)

	// Adjust difficulty
	// If blocks are too fast, make difficulty harder (lower threshold)
	// If blocks are too slow, make difficulty easier (higher threshold)
	
	ratio := float64(avgBlockTime) / float64(targetBlockTime)

	// Clamp adjustment to prevent wild swings
	if ratio < 0.5 {
		ratio = 0.5
	} else if ratio > 2.0 {
		ratio = 2.0
	}

	newDifficulty := uint64(float64(currentDifficulty) * ratio)

	// Ensure we don't go to extremes
	minDiff := uint64(1 << 40)  // Minimum (hardest)
	maxDiff := uint64(1 << 60)  // Maximum (easiest)

	if newDifficulty < minDiff {
		newDifficulty = minDiff
	} else if newDifficulty > maxDiff {
		newDifficulty = maxDiff
	}

	return newDifficulty
}

// GetDifficultyTarget returns the current difficulty target
func GetDifficultyTarget() uint64 {
	// For devnet, start with a moderate difficulty
	// In production, this would be calculated from chain state
	return InitialDifficulty
}

