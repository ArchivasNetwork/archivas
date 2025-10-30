package health

import (
	"sync"
	"time"
)

// ChainHealth tracks blockchain health metrics
type ChainHealth struct {
	mu sync.RWMutex

	StartTime        time.Time
	LastBlockTime    time.Time
	BlockTimes       []time.Duration // Last 100 block times
	TotalBlocks      uint64
	AverageBlockTime time.Duration
}

// NewChainHealth creates a new health tracker
func NewChainHealth() *ChainHealth {
	return &ChainHealth{
		StartTime:  time.Now(),
		BlockTimes: make([]time.Duration, 0, 100),
	}
}

// RecordBlock records a new block and updates metrics
func (h *ChainHealth) RecordBlock() {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()

	if !h.LastBlockTime.IsZero() {
		blockTime := now.Sub(h.LastBlockTime)

		// Keep last 100 block times
		h.BlockTimes = append(h.BlockTimes, blockTime)
		if len(h.BlockTimes) > 100 {
			h.BlockTimes = h.BlockTimes[1:]
		}

		// Calculate average
		if len(h.BlockTimes) > 0 {
			var sum time.Duration
			for _, bt := range h.BlockTimes {
				sum += bt
			}
			h.AverageBlockTime = sum / time.Duration(len(h.BlockTimes))
		}
	}

	h.LastBlockTime = now
	h.TotalBlocks++
}

// GetStats returns current health statistics
func (h *ChainHealth) GetStats() Stats {
	h.mu.RLock()
	defer h.mu.RUnlock()

	uptime := time.Since(h.StartTime)

	var blocksPerHour float64
	if uptime.Hours() > 0 {
		blocksPerHour = float64(h.TotalBlocks) / uptime.Hours()
	}

	return Stats{
		Uptime:           uptime,
		UptimePercent:    100.0, // Simplified - would need restart tracking
		TotalBlocks:      h.TotalBlocks,
		AverageBlockTime: h.AverageBlockTime,
		BlocksPerHour:    blocksPerHour,
		LastBlockTime:    h.LastBlockTime,
	}
}

// Stats contains health statistics
type Stats struct {
	Uptime           time.Duration
	UptimePercent    float64
	TotalBlocks      uint64
	AverageBlockTime time.Duration
	BlocksPerHour    float64
	LastBlockTime    time.Time
}
