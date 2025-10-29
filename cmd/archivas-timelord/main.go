package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/iljanemesis/archivas/vdf"
)

const (
	StepSize       = 500 // iterations per tick
	CheckpointStep = 100 // checkpoint every N iterations
	TickInterval   = 1000 * time.Millisecond
)

type ChainTipResponse struct {
	BlockHash [32]byte `json:"blockHash"`
	Height    uint64   `json:"height"`
}

type VDFUpdateRequest struct {
	Seed       []byte `json:"seed"`
	Iterations uint64 `json:"iterations"`
	Output     []byte `json:"output"`
}

func main() {
	// Parse CLI flags
	nodeURL := flag.String("node", "http://localhost:8080", "Node RPC URL")
	stepSize := flag.Uint64("step", 500, "VDF iterations per tick")
	flag.Parse()

	log.Println("[timelord] Archivas Timelord starting...")
	log.Printf("[timelord] Using node RPC base: %s", *nodeURL)
	log.Printf("[timelord] Step size: %d iterations/tick", *stepSize)
	log.Printf("[timelord] Will poll: %s/chainTip", *nodeURL)
	log.Printf("[timelord] Will POST to: %s/vdf/update", *nodeURL)

	var currentSeed []byte
	var currentIterations uint64
	var currentOutput []byte

	ticker := time.NewTicker(TickInterval)
	defer ticker.Stop()

	for range ticker.C {
		// Get current chain tip
		tip, err := getChainTip(*nodeURL)
		if err != nil {
			log.Printf("[timelord] ⚠️  Error getting chain tip: %v", err)
			continue
		}

		// Compute new seed
		newSeed := computeSeed(tip.BlockHash, tip.Height)

		// Check if we need to reset (new tip)
		if !bytes.Equal(newSeed, currentSeed) {
			log.Printf("[timelord] 🔄 New tip detected: height=%d hash=%x", tip.Height, tip.BlockHash[:8])
			log.Printf("[timelord] 🌱 Resetting VDF seed: %x", newSeed[:8])
			currentSeed = newSeed
			currentIterations = 0
			currentOutput = make([]byte, len(currentSeed))
			copy(currentOutput, currentSeed)
		}

		// Advance VDF
		newIterations := currentIterations + *stepSize
		final, _ := vdf.ComputeSequential(currentSeed, newIterations, CheckpointStep)

		currentIterations = newIterations
		currentOutput = final

		log.Printf("[timelord] seed=%x iter=%d output=%x", currentSeed[:8], currentIterations, currentOutput[:8])

		// Send update to node
		if err := sendVDFUpdate(*nodeURL, currentSeed, currentIterations, currentOutput); err != nil {
			log.Printf("[timelord] ⚠️  Error sending VDF update: %v", err)
		}
	}
}

func getChainTip(nodeURL string) (*ChainTipResponse, error) {
	url := fmt.Sprintf("%s/chainTip", nodeURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d from %s: %s", resp.StatusCode, url, string(body))
	}

	var tip ChainTipResponse
	if err := json.NewDecoder(resp.Body).Decode(&tip); err != nil {
		return nil, fmt.Errorf("failed to decode response from %s: %w", url, err)
	}

	return &tip, nil
}

func computeSeed(blockHash [32]byte, height uint64) []byte {
	h := sha256.New()
	h.Write(blockHash[:])
	binary.Write(h, binary.BigEndian, height)
	sum := h.Sum(nil)
	return sum
}

func sendVDFUpdate(nodeURL string, seed []byte, iterations uint64, output []byte) error {
	update := VDFUpdateRequest{
		Seed:       seed,
		Iterations: iterations,
		Output:     output,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal VDF update: %w", err)
	}

	url := fmt.Sprintf("%s/vdf/update", nodeURL)
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to POST to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d from %s: %s", resp.StatusCode, url, string(body))
	}

	return nil
}
