package evm

import (
	"math/big"

	"github.com/ArchivasNetwork/archivas/network"
)

// NewChainConfigFromGenesis creates a ChainConfig from genesis
func NewChainConfigFromGenesis(genesis *network.GenesisFile) *ChainConfig {
	config := &ChainConfig{
		ChainID: big.NewInt(int64(genesis.NetworkID)),
	}

	if genesis.EVMConfig != nil {
		config.HomesteadBlock = big.NewInt(int64(genesis.EVMConfig.HomesteadBlock))
		config.EIP150Block = big.NewInt(int64(genesis.EVMConfig.EIP150Block))
		config.EIP155Block = big.NewInt(int64(genesis.EVMConfig.EIP155Block))
		config.EIP158Block = big.NewInt(int64(genesis.EVMConfig.EIP158Block))
		config.ByzantiumBlock = big.NewInt(int64(genesis.EVMConfig.ByzantiumBlock))
		config.ConstantinopleBlock = big.NewInt(int64(genesis.EVMConfig.ConstantinopleBlock))
		config.PetersburgBlock = big.NewInt(int64(genesis.EVMConfig.PetersburgBlock))
		config.IstanbulBlock = big.NewInt(int64(genesis.EVMConfig.IstanbulBlock))
		config.BerlinBlock = big.NewInt(int64(genesis.EVMConfig.BerlinBlock))
		config.LondonBlock = big.NewInt(int64(genesis.EVMConfig.LondonBlock))
	}

	return config
}

// DefaultBetanetConfig returns the default Betanet EVM configuration
func DefaultBetanetConfig() *ChainConfig {
	return &ChainConfig{
		ChainID:             big.NewInt(102), // Betanet network ID
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
	}
}

