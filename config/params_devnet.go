package config

const (
	ChainName              = "Archivas Devnet"
	ChainID                = 1616
	DenomSymbol            = "RCHV"
	DenomDecimals          = 8
	TargetBlockTimeSeconds = 20
	InitialBlockReward     = 20_00000000 // 20.00000000 RCHV
)

// GenesisAlloc moved to genesis.go
// Legacy allocation for backward compatibility (will be removed)
var LegacyGenesisAlloc = map[string]int64{
	"arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7": 1_000_000_000_00000000, // 1B RCHV
}

