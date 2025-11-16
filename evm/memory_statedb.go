package evm

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/types"
)

// MemoryStateDB is a simple in-memory implementation of StateDB
// Phase 2: This is a functional but simplified implementation
// Phase 3: Will be replaced with proper Merkle Patricia Trie
type MemoryStateDB struct {
	mu sync.RWMutex
	
	accounts map[address.EVMAddress]*Account
	logs     map[[32]byte][]*Log // logs by tx hash
	
	// Snapshot support
	snapshots []map[address.EVMAddress]*Account
	
	// Current root hash
	root [32]byte
}

// Account represents an EVM account state
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	CodeHash [32]byte
	Code     []byte
	Storage  map[[32]byte][32]byte
}

// NewMemoryStateDB creates a new in-memory state database
func NewMemoryStateDB() *MemoryStateDB {
	return &MemoryStateDB{
		accounts:  make(map[address.EVMAddress]*Account),
		logs:      make(map[[32]byte][]*Log),
		snapshots: make([]map[address.EVMAddress]*Account, 0),
		root:      types.EmptyStateRoot(),
	}
}

// getAccount returns the account for an address, creating it if needed
func (db *MemoryStateDB) getAccount(addr address.EVMAddress) *Account {
	if acc, exists := db.accounts[addr]; exists {
		return acc
	}
	
	// Create new account
	acc := &Account{
		Nonce:   0,
		Balance: big.NewInt(0),
		Storage: make(map[[32]byte][32]byte),
	}
	db.accounts[addr] = acc
	return acc
}

// Exist checks if an account exists
func (db *MemoryStateDB) Exist(addr address.EVMAddress) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	_, exists := db.accounts[addr]
	return exists
}

// Empty checks if an account is empty (no nonce, balance, or code)
func (db *MemoryStateDB) Empty(addr address.EVMAddress) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	acc, exists := db.accounts[addr]
	if !exists {
		return true
	}
	
	return acc.Nonce == 0 && 
		acc.Balance.Sign() == 0 && 
		len(acc.Code) == 0
}

// GetBalance returns the balance of an account
func (db *MemoryStateDB) GetBalance(addr address.EVMAddress) *big.Int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	acc := db.getAccount(addr)
	return new(big.Int).Set(acc.Balance)
}

// AddBalance adds amount to an account's balance
func (db *MemoryStateDB) AddBalance(addr address.EVMAddress, amount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	acc := db.getAccount(addr)
	acc.Balance = new(big.Int).Add(acc.Balance, amount)
}

// SubBalance subtracts amount from an account's balance
func (db *MemoryStateDB) SubBalance(addr address.EVMAddress, amount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	acc := db.getAccount(addr)
	acc.Balance = new(big.Int).Sub(acc.Balance, amount)
}

// SetBalance sets an account's balance
func (db *MemoryStateDB) SetBalance(addr address.EVMAddress, amount *big.Int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	acc := db.getAccount(addr)
	acc.Balance = new(big.Int).Set(amount)
}

// GetNonce returns the nonce of an account
func (db *MemoryStateDB) GetNonce(addr address.EVMAddress) uint64 {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	acc := db.getAccount(addr)
	return acc.Nonce
}

// SetNonce sets the nonce of an account
func (db *MemoryStateDB) SetNonce(addr address.EVMAddress, nonce uint64) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	acc := db.getAccount(addr)
	acc.Nonce = nonce
}

// GetCode returns the code of a contract
func (db *MemoryStateDB) GetCode(addr address.EVMAddress) []byte {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	acc := db.getAccount(addr)
	if len(acc.Code) == 0 {
		return nil
	}
	
	code := make([]byte, len(acc.Code))
	copy(code, acc.Code)
	return code
}

// GetCodeSize returns the size of a contract's code
func (db *MemoryStateDB) GetCodeSize(addr address.EVMAddress) int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	acc := db.getAccount(addr)
	return len(acc.Code)
}

// GetCodeHash returns the hash of a contract's code
func (db *MemoryStateDB) GetCodeHash(addr address.EVMAddress) [32]byte {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	acc := db.getAccount(addr)
	return acc.CodeHash
}

// SetCode sets the code of a contract
func (db *MemoryStateDB) SetCode(addr address.EVMAddress, code []byte) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	acc := db.getAccount(addr)
	acc.Code = make([]byte, len(code))
	copy(acc.Code, code)
	acc.CodeHash = types.Hash256(code)
}

// GetState returns a storage value
func (db *MemoryStateDB) GetState(addr address.EVMAddress, key [32]byte) [32]byte {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	acc := db.getAccount(addr)
	if value, exists := acc.Storage[key]; exists {
		return value
	}
	return [32]byte{}
}

// SetState sets a storage value
func (db *MemoryStateDB) SetState(addr address.EVMAddress, key [32]byte, value [32]byte) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	acc := db.getAccount(addr)
	acc.Storage[key] = value
}

// SetRoot sets the state root
func (db *MemoryStateDB) SetRoot(root [32]byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	// In a real implementation, this would load state from the trie
	// For Phase 2, we just store the root hash
	db.root = root
	return nil
}

// GetRoot returns the current state root
func (db *MemoryStateDB) GetRoot() [32]byte {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	return db.root
}

// Commit commits state changes and returns the new state root
func (db *MemoryStateDB) Commit() ([32]byte, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	// Compute state root by hashing all accounts
	// In a real implementation, this would use a Merkle Patricia Trie
	data := make([]byte, 0)
	
	// Sort addresses for deterministic hashing
	addresses := make([]address.EVMAddress, 0, len(db.accounts))
	for addr := range db.accounts {
		addresses = append(addresses, addr)
	}
	
	// Simple approach: hash concatenated account data
	for _, addr := range addresses {
		acc := db.accounts[addr]
		
		// Add address
		data = append(data, addr.Bytes()...)
		
		// Add nonce
		nonceBytes := make([]byte, 8)
		for i := 0; i < 8; i++ {
			nonceBytes[i] = byte(acc.Nonce >> (8 * (7 - i)))
		}
		data = append(data, nonceBytes...)
		
		// Add balance
		data = append(data, acc.Balance.Bytes()...)
		
		// Add code hash
		data = append(data, acc.CodeHash[:]...)
	}
	
	// If no accounts, return empty root
	if len(data) == 0 {
		db.root = types.EmptyStateRoot()
		return db.root, nil
	}
	
	db.root = types.Hash256(data)
	return db.root, nil
}

// Snapshot creates a snapshot of the current state
func (db *MemoryStateDB) Snapshot() int {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	// Deep copy accounts
	snapshot := make(map[address.EVMAddress]*Account)
	for addr, acc := range db.accounts {
		accCopy := &Account{
			Nonce:    acc.Nonce,
			Balance:  new(big.Int).Set(acc.Balance),
			CodeHash: acc.CodeHash,
			Code:     make([]byte, len(acc.Code)),
			Storage:  make(map[[32]byte][32]byte),
		}
		copy(accCopy.Code, acc.Code)
		for k, v := range acc.Storage {
			accCopy.Storage[k] = v
		}
		snapshot[addr] = accCopy
	}
	
	db.snapshots = append(db.snapshots, snapshot)
	return len(db.snapshots) - 1
}

// RevertToSnapshot reverts state to a snapshot
func (db *MemoryStateDB) RevertToSnapshot(id int) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if id < 0 || id >= len(db.snapshots) {
		return
	}
	
	// Restore snapshot
	db.accounts = db.snapshots[id]
	
	// Remove snapshots after this one
	db.snapshots = db.snapshots[:id]
}

// AddLog adds a log entry
func (db *MemoryStateDB) AddLog(log *Log) {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	logs := db.logs[log.TxHash]
	logs = append(logs, log)
	db.logs[log.TxHash] = logs
}

// GetLogs returns logs for a transaction
func (db *MemoryStateDB) GetLogs(txHash [32]byte) []*Log {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	return db.logs[txHash]
}

// String returns a human-readable representation
func (db *MemoryStateDB) String() string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	return fmt.Sprintf(
		"MemoryStateDB{Accounts: %d, Root: %s}",
		len(db.accounts),
		hex.EncodeToString(db.root[:8]),
	)
}

// LoadGenesisState loads genesis allocations into the state
func (db *MemoryStateDB) LoadGenesisState(allocations map[address.EVMAddress]*big.Int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	for addr, balance := range allocations {
		acc := db.getAccount(addr)
		acc.Balance = new(big.Int).Set(balance)
	}
	
	// Commit genesis state
	_, err := db.Commit()
	return err
}

