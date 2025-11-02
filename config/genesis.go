package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// GenesisAlloc represents a genesis allocation
type GenesisAlloc struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"` // base units
}

// GenesisDoc represents the genesis document
type GenesisDoc struct {
	ChainName          string         `json:"chainName"`
	ChainID            uint64         `json:"chainID"`
	Denom              string         `json:"denom"`
	Decimals           uint8          `json:"decimals"`
	Timestamp          int64          `json:"timestamp"`          // fixed unix timestamp
	Seed               string         `json:"seed"`               // network seed
	NetworkID          string         `json:"networkID"`          // v1.1.1: network identifier
	ProtocolVersion    string         `json:"protocolVersion"`    // v1.1.1: protocol version
	DifficultyParamsID string         `json:"difficultyParamsID"` // v1.1.1: difficulty params identifier
	InitialDifficulty  uint64         `json:"initialDifficulty"`  // v1.1.1: starting difficulty
	Allocations        []GenesisAlloc `json:"allocations"`
}

// LoadGenesis loads genesis from a JSON file
func LoadGenesis(path string) (*GenesisDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis file: %w", err)
	}

	var gen GenesisDoc
	if err := json.Unmarshal(data, &gen); err != nil {
		return nil, fmt.Errorf("failed to parse genesis JSON: %w", err)
	}

	return &gen, nil
}

// HashGenesis computes deterministic hash of genesis
func HashGenesis(gen *GenesisDoc) [32]byte {
	// Use canonical JSON encoding for deterministic hash
	// Sort allocations by address for consistency
	sortedAllocs := make([]GenesisAlloc, len(gen.Allocations))
	copy(sortedAllocs, gen.Allocations)
	sort.Slice(sortedAllocs, func(i, j int) bool {
		return sortedAllocs[i].Address < sortedAllocs[j].Address
	})

	canonical := struct {
		ChainName          string         `json:"chainName"`
		ChainID            uint64         `json:"chainID"`
		Denom              string         `json:"denom"`
		Decimals           uint8          `json:"decimals"`
		Timestamp          int64          `json:"timestamp"`
		Seed               string         `json:"seed"`
		NetworkID          string         `json:"networkID"`
		ProtocolVersion    string         `json:"protocolVersion"`
		DifficultyParamsID string         `json:"difficultyParamsID"`
		InitialDifficulty  uint64         `json:"initialDifficulty"`
		Allocations        []GenesisAlloc `json:"allocations"`
	}{
		ChainName:          gen.ChainName,
		ChainID:            gen.ChainID,
		Denom:              gen.Denom,
		Decimals:           gen.Decimals,
		Timestamp:          gen.Timestamp,
		Seed:               gen.Seed,
		NetworkID:          gen.NetworkID,
		ProtocolVersion:    gen.ProtocolVersion,
		DifficultyParamsID: gen.DifficultyParamsID,
		InitialDifficulty:  gen.InitialDifficulty,
		Allocations:        sortedAllocs,
	}

	data, _ := json.Marshal(canonical)
	return sha256.Sum256(data)
}

// ValidateGenesis checks if genesis params match expected values
// v1.1.1: Prevents starting with incompatible genesis
func ValidateGenesis(gen *GenesisDoc, expectedDifficultyParamsID, expectedProtocolVersion string) error {
	if gen.DifficultyParamsID != expectedDifficultyParamsID {
		return fmt.Errorf("genesis difficulty params mismatch: got %s, expected %s", 
			gen.DifficultyParamsID, expectedDifficultyParamsID)
	}
	if gen.ProtocolVersion != expectedProtocolVersion {
		return fmt.Errorf("genesis protocol version mismatch: got %s, expected %s",
			gen.ProtocolVersion, expectedProtocolVersion)
	}
	return nil
}

// GenesisAllocToMap converts allocations to map for world state
func GenesisAllocToMap(allocs []GenesisAlloc) map[string]int64 {
	result := make(map[string]int64)
	for _, alloc := range allocs {
		result[alloc.Address] = int64(alloc.Amount)
	}
	return result
}
