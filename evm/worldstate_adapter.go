package evm

import (
	"log"
	"math/big"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/ledger"
)

// WorldStateAdapter bridges WorldState (arcv addresses, int64) with EVM StateDB (0x addresses, big.Int)
// This allows existing RCHV balances to be accessible via ETH RPC calls
type WorldStateAdapter struct {
	*MemoryStateDB
	worldState *ledger.WorldState
}

// NewWorldStateAdapter creates a new adapter that syncs WorldState with EVM StateDB
func NewWorldStateAdapter(ws *ledger.WorldState) *WorldStateAdapter {
	adapter := &WorldStateAdapter{
		MemoryStateDB: NewMemoryStateDB(),
		worldState:    ws,
	}
	
	// Initial sync: load all existing balances from WorldState into EVM StateDB
	adapter.SyncFromWorldState()
	
	return adapter
}

// SyncFromWorldState syncs all balances from WorldState to EVM StateDB
// This ensures ETH RPC calls see the correct RCHV balances
func (a *WorldStateAdapter) SyncFromWorldState() {
	if a.worldState == nil {
		return
	}
	
	synced := 0
	for arcvAddr, acct := range a.worldState.Accounts {
		// Convert arcv address to EVM address
		evmAddr, err := address.ParseAddress(arcvAddr, "arcv")
		if err != nil {
			log.Printf("[evm-adapter] Failed to parse address %s: %v", arcvAddr, err)
			continue
		}
		
		// Convert int64 balance to big.Int with Wei conversion
		// Archivas: 1 RCHV = 10^8 base units (like Bitcoin satoshis)
		// Ethereum: 1 ETH = 10^18 Wei
		// Multiply by 10^10 to convert Archivas base units to Wei
		balance := big.NewInt(acct.Balance)
		multiplier := big.NewInt(10000000000) // 10^10
		balance.Mul(balance, multiplier)
		
		// Set balance in EVM StateDB
		a.MemoryStateDB.SetBalance(evmAddr, balance)
		a.MemoryStateDB.SetNonce(evmAddr, acct.Nonce)
		
		synced++
	}
	
	log.Printf("[evm-adapter] Synced %d accounts from WorldState to EVM StateDB", synced)
}

// GetBalance overrides to always return fresh data from WorldState
func (a *WorldStateAdapter) GetBalance(addr address.EVMAddress) *big.Int {
	// Try to find matching arcv address in WorldState
	// First, convert EVM address to arcv format
	arcvAddr, err := address.EncodeARCVAddress(addr, "arcv")
	if err != nil {
		// If conversion fails, fall back to MemoryStateDB
		return a.MemoryStateDB.GetBalance(addr)
	}
	
	// Check if this address exists in WorldState
	if acct, exists := a.worldState.Accounts[arcvAddr]; exists {
		// Archivas uses 8 decimals (10^8 base units per RCHV)
		// Ethereum uses 18 decimals (10^18 Wei per ETH)
		// Convert: multiply by 10^10 to convert Archivas base units to Wei
		balance := big.NewInt(acct.Balance)
		multiplier := big.NewInt(10000000000) // 10^10
		balance.Mul(balance, multiplier)
		return balance
	}
	
	// Fall back to MemoryStateDB for addresses not in WorldState
	return a.MemoryStateDB.GetBalance(addr)
}

// GetNonce overrides to return fresh data from WorldState
func (a *WorldStateAdapter) GetNonce(addr address.EVMAddress) uint64 {
	// Convert to arcv and check WorldState
	arcvAddr, err := address.EncodeARCVAddress(addr, "arcv")
	if err != nil {
		return a.MemoryStateDB.GetNonce(addr)
	}
	
	if acct, exists := a.worldState.Accounts[arcvAddr]; exists {
		return acct.Nonce
	}
	
	return a.MemoryStateDB.GetNonce(addr)
}

// Exist overrides to check both WorldState and EVM StateDB
func (a *WorldStateAdapter) Exist(addr address.EVMAddress) bool {
	// Check WorldState first
	arcvAddr, err := address.EncodeARCVAddress(addr, "arcv")
	if err == nil {
		if _, exists := a.worldState.Accounts[arcvAddr]; exists {
			return true
		}
	}
	
	// Fall back to EVM StateDB (for contract addresses, etc.)
	return a.MemoryStateDB.Exist(addr)
}

