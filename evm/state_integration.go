package evm

import (
	"math/big"
	
	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/ledger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/holiman/uint256"
)

// CreateGethStateDB creates a real go-ethereum StateDB from Archivas WorldState
// This allows full EVM contract execution with proper state management
func CreateGethStateDB(ws *ledger.WorldState) (*state.StateDB, error) {
	// Create an in-memory database for the state trie
	memdb := rawdb.NewMemoryDatabase()
	
	// Create trie database
	trieDB := triedb.NewDatabase(memdb, &triedb.Config{
		Preimages: true,
	})
	
	// Create StateDB with empty root (will be populated)
	statedb, err := state.New(common.Hash{}, state.NewDatabase(trieDB, nil))
	if err != nil {
		return nil, err
	}
	
	// Populate StateDB with current WorldState
	for arcvAddr, account := range ws.Accounts {
		// Parse ARCV address to get EVM address
		evmAddr, err := address.ParseAddress(arcvAddr, "arcv")
		if err != nil {
			continue // Skip invalid addresses
		}
		
		commonAddr := common.Address(evmAddr)
		
		// Set balance (convert RCHV to Wei: 8 decimals -> 18 decimals)
		balanceRCHV := big.NewInt(account.Balance)
		balanceWei := new(big.Int).Mul(balanceRCHV, big.NewInt(10_000_000_000))
		statedb.SetBalance(commonAddr, uint256.MustFromBig(balanceWei), tracing.BalanceIncreaseGenesisBalance)
		
		// Set nonce
		statedb.SetNonce(commonAddr, account.Nonce, tracing.NonceChangeUnspecified)
		
		// Set code if any (for contracts)
		// Note: WorldState doesn't store code yet, but when it does, set it here
		// if len(account.Code) > 0 {
		// 	 statedb.SetCode(commonAddr, account.Code)
		// }
	}
	
	return statedb, nil
}

// SyncStateDBToWorldState copies changes from StateDB back to WorldState after execution
// This ensures our native state stays in sync with EVM state
func SyncStateDBToWorldState(statedb *state.StateDB, ws *ledger.WorldState) error {
	// Get all modified addresses from StateDB
	// Note: This is a simplified approach - in production we'd track dirty addresses
	
	// For now, we sync back the accounts we know about
	for arcvAddr := range ws.Accounts {
		evmAddr, err := address.ParseAddress(arcvAddr, "arcv")
		if err != nil {
			continue
		}
		
		commonAddr := common.Address(evmAddr)
		
		// Get balance from StateDB (in Wei)
		balanceWei := statedb.GetBalance(commonAddr)
		// Convert Wei to RCHV (18 decimals -> 8 decimals)
		balanceRCHV := new(big.Int).Div(balanceWei.ToBig(), big.NewInt(10_000_000_000))
		
		// Get nonce
		nonce := statedb.GetNonce(commonAddr)
		
		// Update WorldState
		if account, exists := ws.Accounts[arcvAddr]; exists {
			account.Balance = balanceRCHV.Int64()
			account.Nonce = nonce
		} else {
			// New account created during execution
			ws.Accounts[arcvAddr] = &ledger.AccountState{
				Balance: balanceRCHV.Int64(),
				Nonce:   nonce,
			}
		}
	}
	
	return nil
}
