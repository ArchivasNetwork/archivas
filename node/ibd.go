package node

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// IBDState tracks Initial Block Download progress
type IBDState struct {
	InProgress    bool   `json:"inProgress"`
	StartHeight   uint64 `json:"startHeight"`
	CurrentHeight uint64 `json:"currentHeight"`
	TargetHeight  uint64 `json:"targetHeight"`
	PeerURL       string `json:"peerURL"`
	StartedAt     int64  `json:"startedAt"`
}

// IBDConfig holds IBD parameters
type IBDConfig struct {
	BatchSize         int           // Blocks per request (default 512)
	IBDThreshold      uint64        // Gap to trigger IBD (default 200)
	CatchUpThreshold  uint64        // Gap to exit IBD (default 50)
	RetryDelay        time.Duration // Initial retry delay (default 5s)
	MaxRetries        int           // Max retries per batch (default 5)
	StateFile         string        // Path to ibd_state.json
	ProgressInterval  time.Duration // Log progress every N seconds (default 5s)
}

// DefaultIBDConfig returns sensible defaults
func DefaultIBDConfig(dataDir string) *IBDConfig {
	return &IBDConfig{
		BatchSize:         512,
		IBDThreshold:      200,
		CatchUpThreshold:  50,
		RetryDelay:        5 * time.Second,
		MaxRetries:        5,
		StateFile:         filepath.Join(dataDir, "ibd_state.json"),
		ProgressInterval:  5 * time.Second,
	}
}

// NodeIBDInterface defines what IBD needs from the node
type NodeIBDInterface interface {
	GetCurrentHeight() uint64
	ApplyBlock(blockData json.RawMessage) error
}

// IBDManager handles Initial Block Download
type IBDManager struct {
	config *IBDConfig
	node   NodeIBDInterface
	state  *IBDState
}

// NewIBDManager creates a new IBD manager
func NewIBDManager(config *IBDConfig, node NodeIBDInterface) *IBDManager {
	return &IBDManager{
		config: config,
		node:   node,
		state:  &IBDState{},
	}
}

// LoadState loads IBD state from disk (for resume support)
func (m *IBDManager) LoadState() error {
	data, err := os.ReadFile(m.config.StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No state file = fresh start
		}
		return err
	}

	return json.Unmarshal(data, m.state)
}

// SaveState persists IBD state to disk
func (m *IBDManager) SaveState() error {
	data, err := json.Marshal(m.state)
	if err != nil {
		return err
	}

	return os.WriteFile(m.config.StateFile, data, 0644)
}

// ClearState removes IBD state file (called when sync complete)
func (m *IBDManager) ClearState() error {
	return os.Remove(m.config.StateFile)
}

// ShouldRunIBD checks if IBD is needed
func (m *IBDManager) ShouldRunIBD(localHeight, remoteHeight uint64) bool {
	gap := remoteHeight - localHeight
	return gap >= m.config.IBDThreshold
}

// RunIBD performs Initial Block Download from a peer
func (m *IBDManager) RunIBD(peerURL string) error {
	// Fetch remote tip
	remoteTip, err := m.fetchRemoteTip(peerURL)
	if err != nil {
		return fmt.Errorf("failed to fetch remote tip: %w", err)
	}

	localHeight := m.node.GetCurrentHeight()

	// Check if IBD needed
	if !m.ShouldRunIBD(localHeight, remoteTip) {
		log.Printf("[IBD] Node is close to tip (local=%d, remote=%d, gap=%d), skipping IBD",
			localHeight, remoteTip, remoteTip-localHeight)
		return nil
	}

	// Initialize state
	m.state = &IBDState{
		InProgress:    true,
		StartHeight:   localHeight,
		CurrentHeight: localHeight,
		TargetHeight:  remoteTip,
		PeerURL:       peerURL,
		StartedAt:     time.Now().Unix(),
	}

	if err := m.SaveState(); err != nil {
		log.Printf("[IBD] Warning: failed to save state: %v", err)
	}

	log.Printf("[IBD] Starting sync: local=%d remote=%d (%.1f%% behind, %d blocks to download)",
		localHeight, remoteTip, float64(localHeight)*100/float64(remoteTip), remoteTip-localHeight)

	// Download loop
	lastLogTime := time.Now()
	consecutiveFailures := 0

	for {
		currentHeight := m.node.GetCurrentHeight()
		m.state.CurrentHeight = currentHeight

		// Check if caught up
		if remoteTip-currentHeight <= m.config.CatchUpThreshold {
			log.Printf("[IBD] Complete! Synced to height %d (gap: %d blocks)",
				currentHeight, remoteTip-currentHeight)
			m.ClearState()
			return nil
		}

		// Download next batch
		blocks, newRemoteTip, err := m.fetchBlockBatch(peerURL, currentHeight+1, m.config.BatchSize)
		if err != nil {
			consecutiveFailures++
			if consecutiveFailures >= m.config.MaxRetries {
				return fmt.Errorf("IBD failed after %d retries: %w", m.config.MaxRetries, err)
			}

			retryDelay := m.config.RetryDelay * time.Duration(consecutiveFailures)
			log.Printf("[IBD] Fetch error (attempt %d/%d): %v, retrying in %v...",
				consecutiveFailures, m.config.MaxRetries, err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Reset failure counter on success
		consecutiveFailures = 0

		// Update remote tip (chain may be growing)
		if newRemoteTip > remoteTip {
			remoteTip = newRemoteTip
			m.state.TargetHeight = remoteTip
		}

		// Apply blocks
		for _, blockData := range blocks {
			if err := m.node.ApplyBlock(blockData); err != nil {
				return fmt.Errorf("failed to apply block: %w", err)
			}
		}

		// Progress logging
		if time.Since(lastLogTime) > m.config.ProgressInterval {
			currentHeight = m.node.GetCurrentHeight()
			pct := float64(currentHeight) * 100 / float64(remoteTip)
			remaining := remoteTip - currentHeight
			blocksDownloaded := currentHeight - m.state.StartHeight
			elapsed := time.Since(time.Unix(m.state.StartedAt, 0))
			rate := float64(blocksDownloaded) / elapsed.Seconds()
			eta := time.Duration(float64(remaining)/rate) * time.Second

			log.Printf("[IBD] Progress: %d/%d (%.1f%% complete, %d blocks remaining, %.1f blocks/sec, ETA %v)",
				currentHeight, remoteTip, pct, remaining, rate, eta.Round(time.Second))

			lastLogTime = time.Now()

			// Save state periodically
			if err := m.SaveState(); err != nil {
				log.Printf("[IBD] Warning: failed to save state: %v", err)
			}
		}

		// If no blocks received, might be caught up
		if len(blocks) == 0 {
			break
		}
	}

	m.ClearState()
	return nil
}

// fetchRemoteTip gets current chain tip from peer
func (m *IBDManager) fetchRemoteTip(peerURL string) (uint64, error) {
	url := fmt.Sprintf("%s/chainTip", peerURL)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Height string `json:"height"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	height, err := strconv.ParseUint(result.Height, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid height: %s", result.Height)
	}

	return height, nil
}

// fetchBlockBatch fetches a batch of blocks from /blocks/range
func (m *IBDManager) fetchBlockBatch(peerURL string, fromHeight uint64, limit int) ([]json.RawMessage, uint64, error) {
	url := fmt.Sprintf("%s/blocks/range?from=%d&limit=%d", peerURL, fromHeight, limit)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	var result struct {
		From   uint64            `json:"from"`
		To     uint64            `json:"to"`
		Blocks []json.RawMessage `json:"blocks"`
		Tip    uint64            `json:"tip"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	if result.Tip == 0 {
		return nil, 0, fmt.Errorf("peer returned tip=0")
	}

	return result.Blocks, result.Tip, nil
}

// RunIBDWithRetry attempts IBD against multiple peers
func (m *IBDManager) RunIBDWithRetry(peerURLs []string) error {
	for i, peerURL := range peerURLs {
		log.Printf("[IBD] Attempting sync from peer %d/%d: %s", i+1, len(peerURLs), peerURL)

		if err := m.RunIBD(peerURL); err != nil {
			log.Printf("[IBD] Failed with peer %s: %v", peerURL, err)
			continue
		}

		log.Printf("[IBD] Successfully synced via %s", peerURL)
		return nil
	}

	return fmt.Errorf("IBD failed with all %d peers", len(peerURLs))
}

