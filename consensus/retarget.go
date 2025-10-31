package consensus

import (
	"time"
)

// RetargetConfig holds difficulty retargeting parameters
type RetargetConfig struct {
	BlockTimeTarget     time.Duration // Target block time (e.g., 30s)
	RetargetInterval    int           // Adjust every N blocks
	Alpha               float64       // EMA smoothing factor (0..1)
	MaxIncrease         float64       // Max increase per interval (e.g., 1.5 = 50% up)
	MaxDecrease         float64       // Max decrease per interval (e.g., 0.67 = 33% down)
	MinDifficulty       uint64        // Floor
}

// DefaultRetargetConfig returns sensible defaults
func DefaultRetargetConfig() RetargetConfig {
	return RetargetConfig{
		BlockTimeTarget:     30 * time.Second,
		RetargetInterval:    10, // Adjust every 10 blocks (faster)
		Alpha:               0.45, // More responsive (was 0.25)
		MaxIncrease:         1.25, // Slower increase (25%)
		MaxDecrease:         0.30, // AGGRESSIVE decrease (70% drop allowed!)
		MinDifficulty:       10_000_000, // Floor at 10M (not 1M)
	}
}

// RetargetDifficulty calculates new difficulty using EMA
func RetargetDifficulty(currentDiff uint64, observedBlockTime time.Duration, config RetargetConfig) uint64 {
	if observedBlockTime == 0 {
		return currentDiff
	}

	// Calculate scale factor
	target := config.BlockTimeTarget.Seconds()
	observed := observedBlockTime.Seconds()
	scaleFactor := observed / target

	// Calculate ideal new difficulty
	idealDiff := float64(currentDiff) * scaleFactor

	// Apply EMA smoothing
	// newDiff = (1-α)*currentDiff + α*idealDiff
	alpha := config.Alpha
	newDiff := (1.0-alpha)*float64(currentDiff) + alpha*idealDiff

	// Apply bounds
	maxUp := float64(currentDiff) * config.MaxIncrease
	maxDown := float64(currentDiff) * config.MaxDecrease

	if newDiff > maxUp {
		newDiff = maxUp
	}
	if newDiff < maxDown {
		newDiff = maxDown
	}

	// Apply floor
	result := uint64(newDiff)
	if result < config.MinDifficulty {
		result = config.MinDifficulty
	}

	return result
}

// CalculateObservedBlockTime computes average block time from recent blocks
func CalculateObservedBlockTime(timestamps []int64) time.Duration {
	if len(timestamps) < 2 {
		return 0
	}

	totalTime := int64(0)
	for i := 1; i < len(timestamps); i++ {
		diff := timestamps[i] - timestamps[i-1]
		totalTime += diff
	}

	avgSeconds := totalTime / int64(len(timestamps)-1)
	return time.Duration(avgSeconds) * time.Second
}

