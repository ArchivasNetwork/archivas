package network

import (
	"fmt"
	"time"
)

// NetworkProfile defines the configuration for a specific Archivas network
type NetworkProfile struct {
	Name             string        // Human-readable name (e.g., "betanet", "devnet-legacy")
	ChainID          string        // Chain identifier (e.g., "archivas-betanet-1")
	NetworkID        uint64        // Numeric network ID for P2P
	ProtocolVersion  int           // Protocol version for compatibility checks
	GenesisPath      string        // Path to genesis file
	DefaultSeeds     []string      // Default seed/bootnode addresses
	DefaultRPCPort   int           // Default RPC port
	DefaultP2PPort   int           // Default P2P port
	TargetBlockTime  time.Duration // Target time between blocks
	InitialDifficulty uint64       // Starting PoST difficulty
	Bech32Prefix     string        // Address prefix for Bech32 encoding (e.g., "arcv")
}

// NetworkProfiles is the global registry of available networks
var NetworkProfiles = map[string]*NetworkProfile{
	"betanet": {
		Name:              "betanet",
		ChainID:           "archivas-betanet-1",
		NetworkID:         1644,
		ProtocolVersion:   2,
		GenesisPath:       "configs/genesis-betanet.json",
		DefaultSeeds: []string{
			"seed1.betanet.archivas.ai:9090",
			"seed2.betanet.archivas.ai:9090",
		},
		DefaultRPCPort:    8545, // Standard Ethereum RPC port
		DefaultP2PPort:    9090,
		TargetBlockTime:   20 * time.Second,
		InitialDifficulty: 15000000,
		Bech32Prefix:      "arcv",
	},
	"devnet-legacy": {
		Name:              "devnet-legacy",
		ChainID:           "archivas-devnet-1",
		NetworkID:         1,
		ProtocolVersion:   1,
		GenesisPath:       "genesis/devnet.genesis.json",
		DefaultSeeds: []string{
			"seed.archivas.ai:9090",
			"seed2.archivas.ai:30303",
		},
		DefaultRPCPort:    8080,
		DefaultP2PPort:    9090,
		TargetBlockTime:   20 * time.Second,
		InitialDifficulty: 15000000,
		Bech32Prefix:      "arcv", // Keep same prefix for compatibility
	},
}

// GetProfile returns the network profile for the given name
func GetProfile(name string) (*NetworkProfile, error) {
	profile, exists := NetworkProfiles[name]
	if !exists {
		return nil, fmt.Errorf("unknown network: %s (available: betanet, devnet-legacy)", name)
	}
	return profile, nil
}

// DefaultNetwork returns the default network name
func DefaultNetwork() string {
	return "betanet"
}

// ValidateProfile performs basic validation on a network profile
func ValidateProfile(profile *NetworkProfile) error {
	if profile.Name == "" {
		return fmt.Errorf("network name cannot be empty")
	}
	if profile.ChainID == "" {
		return fmt.Errorf("chain ID cannot be empty")
	}
	if profile.NetworkID == 0 {
		return fmt.Errorf("network ID cannot be zero")
	}
	if profile.ProtocolVersion <= 0 {
		return fmt.Errorf("protocol version must be positive")
	}
	if profile.GenesisPath == "" {
		return fmt.Errorf("genesis path cannot be empty")
	}
	if profile.Bech32Prefix == "" {
		return fmt.Errorf("bech32 prefix cannot be empty")
	}
	return nil
}

