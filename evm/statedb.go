package evm

import (
	"math/big"

	"github.com/ArchivasNetwork/archivas/address"
)

// StateDB is the interface for EVM state access
// This abstracts the underlying state trie implementation
type StateDB interface {
	// Account queries
	Exist(addr address.EVMAddress) bool
	Empty(addr address.EVMAddress) bool
	
	// Balance
	GetBalance(addr address.EVMAddress) *big.Int
	AddBalance(addr address.EVMAddress, amount *big.Int)
	SubBalance(addr address.EVMAddress, amount *big.Int)
	SetBalance(addr address.EVMAddress, amount *big.Int)
	
	// Nonce
	GetNonce(addr address.EVMAddress) uint64
	SetNonce(addr address.EVMAddress, nonce uint64)
	
	// Code
	GetCode(addr address.EVMAddress) []byte
	GetCodeSize(addr address.EVMAddress) int
	GetCodeHash(addr address.EVMAddress) [32]byte
	SetCode(addr address.EVMAddress, code []byte)
	
	// Storage
	GetState(addr address.EVMAddress, key [32]byte) [32]byte
	SetState(addr address.EVMAddress, key [32]byte, value [32]byte)
	
	// State management
	SetRoot(root [32]byte) error
	GetRoot() [32]byte
	Commit() ([32]byte, error)
	
	// Snapshots (for revert support)
	Snapshot() int
	RevertToSnapshot(id int)
	
	// Logs
	AddLog(log *Log)
	GetLogs(txHash [32]byte) []*Log
}

// Log represents an EVM event log
type Log struct {
	Address address.EVMAddress
	Topics  [][32]byte
	Data    []byte
	TxHash  [32]byte
	TxIndex uint32
	BlockHeight uint64
	Index   uint32
}

