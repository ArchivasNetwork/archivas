package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	nodeURL := "http://localhost:8080"
	if len(os.Args) > 1 {
		nodeURL = os.Args[1]
	}

	log.Println("[timelord] Archivas Timelord starting...")
	log.Printf("[timelord] Node: %s\n", nodeURL)

	var currentSeed []byte
	var currentIterations uint64
	var currentOutput []byte

	ticker := time.NewTicker(TickInterval)
	defer ticker.Stop()

	for range ticker.C {
		// Get current chain tip
		tip, err := getChainTip(nodeURL)
		if err != nil {
			log.Printf("[timelord] ⚠️  Error getting chain tip: %v", err)
			continue
		}

		// Compute new seed
		newSeed := computeSeed(tip.BlockHash, tip.Height)

		// Check if we need to reset (new tip)
		if !bytes.Equal(newSeed, currentSeed) {
			log.Printf("[timelord] 🔄 New chain tip detected (height=%d)", tip.Height)
			log.Printf("[timelord] 🌱 New seed: %x", newSeed[:8])
			currentSeed = newSeed
			currentIterations = 0
			currentOutput = make([]byte, len(currentSeed))
			copy(currentOutput, currentSeed)
		}

		// Advance VDF
		newIterations := currentIterations + StepSize
		final, _ := vdf.ComputeSequential(currentSeed, newIterations, CheckpointStep)

		currentIterations = newIterations
		currentOutput = final

		log.Printf("[timelord] iter=%d output=%x", currentIterations, currentOutput[:8])

		// Send update to node
		if err := sendVDFUpdate(nodeURL, currentSeed, currentIterations, currentOutput); err != nil {
			log.Printf("[timelord] ⚠️  Error sending VDF update: %v", err)
		}
	}
}

func getChainTip(nodeURL string) (*ChainTipResponse, error) {
	resp, err := http.Get(nodeURL + "/chainTip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tip ChainTipResponse
	if err := json.NewDecoder(resp.Body).Decode(&tip); err != nil {
		return nil, err
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
		return err
	}

	resp, err := http.Post(nodeURL+"/vdf/update", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
