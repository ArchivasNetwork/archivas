package consensus

import (
	"time"
)

// SmoothedDifficulty calculates difficulty using 10-block moving average
func SmoothedDifficulty(recentBlocks []BlockInfo, targetBlockTime time.Duration, initialDifficulty uint64) uint64 {
	if len(recentBlocks) < 2 {
		return initialDifficulty
	}
	
	// Use last 10 blocks for smoothing
	window := 10
	if len(recentBlocks) < window {
		window = len(recentBlocks)
	}
	
	blocks := recentBlocks[len(recentBlocks)-window:]
	
	// Calculate actual average time between blocks
	totalTime := int64(0)
	for i := 1; i < len(blocks); i++ {
		timeDiff := blocks[i].Timestamp - blocks[i-1].Timestamp
		totalTime += timeDiff
	}
	
	avgTime := totalTime / int64(len(blocks)-1)
	targetSec := int64(targetBlockTime.Seconds())
	
	if avgTime == 0 {
		return initialDifficulty
	}
	
	// Get last difficulty
	lastDiff := blocks[len(blocks)-1].Difficulty
	
	// Adjust: if blocks too fast, increase difficulty; if too slow, decrease
	// ratio = targetTime / actualTime
	// newDiff = oldDiff * ratio
	
	ratio := float64(targetSec) / float64(avgTime)
	newDiff := uint64(float64(lastDiff) * ratio)
	
	// Limit adjustment to 2x per window to prevent wild swings
	maxIncrease := lastDiff * 2
	maxDecrease := lastDiff / 2
	
	if newDiff > maxIncrease {
		newDiff = maxIncrease
	}
	if newDiff < maxDecrease {
		newDiff = maxDecrease
	}
	
	// Ensure minimum difficulty
	minDiff := uint64(1000000)
	if newDiff < minDiff {
		newDiff = minDiff
	}
	
	return newDiff
}

// BlockInfo contains minimal block data for difficulty calculation
type BlockInfo struct {
	Height     uint64
	Timestamp  int64
	Difficulty uint64
}
