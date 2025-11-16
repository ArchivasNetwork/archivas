package network

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ArchivasNetwork/archivas/address"
)

// GenesisFile represents the complete genesis configuration
type GenesisFile struct {
	ChainName             string                 `json:"chain_name"`
	ChainID               string                 `json:"chain_id"`
	NetworkID             uint64                 `json:"network_id"`
	ProtocolVersion       int                    `json:"protocol_version"`
	GenesisTime           string                 `json:"genesis_time"`
	InitialDifficulty     uint64                 `json:"initial_difficulty"`
	TargetBlockTimeSeconds int                   `json:"target_block_time_seconds"`
	DifficultyAdjustWindow int                   `json:"difficulty_adjustment_window"`
	MaxBlockSize          uint64                 `json:"max_block_size"`
	EVMConfig             *EVMConfig             `json:"evm_config,omitempty"`
	InitialState          InitialState           `json:"initial_state"`
	ConsensusParams       ConsensusParams        `json:"consensus_params"`
	Allocations           []Allocation           `json:"allocations"`
}

// EVMConfig contains EVM hard fork block numbers
type EVMConfig struct {
	ChainID             uint64 `json:"chain_id"`
	HomesteadBlock      uint64 `json:"homestead_block"`
	EIP150Block         uint64 `json:"eip150_block"`
	EIP155Block         uint64 `json:"eip155_block"`
	EIP158Block         uint64 `json:"eip158_block"`
	ByzantiumBlock      uint64 `json:"byzantium_block"`
	ConstantinopleBlock uint64 `json:"constantinople_block"`
	PetersburgBlock     uint64 `json:"petersburg_block"`
	IstanbulBlock       uint64 `json:"istanbul_block"`
	BerlinBlock         uint64 `json:"berlin_block"`
	LondonBlock         uint64 `json:"london_block"`
}

// InitialState represents the EVM state at genesis
type InitialState struct {
	StateRoot    string             `json:"state_root"`
	ReceiptsRoot string             `json:"receipts_root"`
	Accounts     []GenesisAccount   `json:"accounts"`
}

// GenesisAccount represents an account in the genesis state
type GenesisAccount struct {
	Address string `json:"address"` // Can be 0x or arcv format
	Balance string `json:"balance"` // Wei or smallest unit
	Code    string `json:"code,omitempty"`    // Contract bytecode (hex)
	Storage map[string]string `json:"storage,omitempty"` // Storage slots
	Nonce   uint64 `json:"nonce,omitempty"`
}

// ConsensusParams contains PoST consensus parameters
type ConsensusParams struct {
	PoST PoSTParams `json:"post"`
}

// PoSTParams contains Proof of Space and Time parameters
type PoSTParams struct {
	KSize                  int    `json:"k_size"`
	MinPlotSize            uint64 `json:"min_plot_size"`
	ChallengeDelay         int    `json:"challenge_delay"`
	SignagePointInterval   int    `json:"signage_point_interval"`
}

// Allocation represents an initial token allocation
type Allocation struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

// LoadGenesis loads and validates a genesis file
func LoadGenesis(path string) (*GenesisFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read genesis file: %w", err)
	}

	var genesis GenesisFile
	if err := json.Unmarshal(data, &genesis); err != nil {
		return nil, fmt.Errorf("failed to parse genesis JSON: %w", err)
	}

	// Validate
	if err := ValidateGenesis(&genesis); err != nil {
		return nil, fmt.Errorf("invalid genesis: %w", err)
	}

	return &genesis, nil
}

// ValidateGenesis performs basic validation on a genesis file
func ValidateGenesis(g *GenesisFile) error {
	if g.ChainName == "" {
		return fmt.Errorf("chain_name cannot be empty")
	}
	if g.ChainID == "" {
		return fmt.Errorf("chain_id cannot be empty")
	}
	if g.NetworkID == 0 {
		return fmt.Errorf("network_id cannot be zero")
	}
	if g.ProtocolVersion <= 0 {
		return fmt.Errorf("protocol_version must be positive")
	}
	
	// Validate genesis time
	if _, err := time.Parse(time.RFC3339, g.GenesisTime); err != nil {
		return fmt.Errorf("invalid genesis_time format: %w", err)
	}

	// Validate all account addresses
	for i, acc := range g.InitialState.Accounts {
		if _, err := address.ParseAddress(acc.Address, "arcv"); err != nil {
			return fmt.Errorf("invalid address in account %d: %w", i, err)
		}
	}

	// Validate allocations
	for i, alloc := range g.Allocations {
		if _, err := address.ParseAddress(alloc.Address, "arcv"); err != nil {
			return fmt.Errorf("invalid address in allocation %d: %w", i, err)
		}
	}

	return nil
}

// ComputeGenesisHash computes a deterministic hash of the genesis file
// This is used for P2P identity verification
func ComputeGenesisHash(g *GenesisFile) (string, error) {
	// Serialize to canonical JSON (sorted keys)
	data, err := json.Marshal(g)
	if err != nil {
		return "", fmt.Errorf("failed to marshal genesis: %w", err)
	}

	hash := sha256.Sum256(data)
	return "0x" + hex.EncodeToString(hash[:]), nil
}

// MatchesProfile checks if a genesis file matches a network profile
func MatchesProfile(genesis *GenesisFile, profile *NetworkProfile) error {
	if genesis.ChainID != profile.ChainID {
		return fmt.Errorf("chain ID mismatch: genesis=%s, profile=%s", genesis.ChainID, profile.ChainID)
	}
	if genesis.NetworkID != profile.NetworkID {
		return fmt.Errorf("network ID mismatch: genesis=%d, profile=%d", genesis.NetworkID, profile.NetworkID)
	}
	if genesis.ProtocolVersion != profile.ProtocolVersion {
		return fmt.Errorf("protocol version mismatch: genesis=%d, profile=%d", genesis.ProtocolVersion, profile.ProtocolVersion)
	}
	return nil
}

