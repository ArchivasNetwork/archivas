package types

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ArchivasNetwork/archivas/address"
)

// Transaction is the interface for all transaction types
type Transaction interface {
	Hash() [32]byte
	Type() TxType
	From() address.EVMAddress
	Nonce() uint64
	Gas() uint64
	GasPrice() *big.Int
}

// TxType represents the type of transaction
type TxType uint8

const (
	TxTypeLegacy    TxType = 0 // Legacy Ethereum transaction (EIP-155)
	TxTypeEIP2930   TxType = 1 // EIP-2930 access list transaction
	TxTypeEIP1559   TxType = 2 // EIP-1559 dynamic fee transaction
	TxTypeEVMCall   TxType = 3 // EVM contract call or transfer (legacy Archivas)
	TxTypeEVMDeploy TxType = 4 // EVM contract deployment (legacy Archivas)
)

// EVMTransaction represents an EVM-compatible transaction
type EVMTransaction struct {
	TypeFlag TxType // Transaction type
	NonceVal uint64 // Account nonce
	GasPriceVal *big.Int // Gas price in wei
	GasLimitVal uint64 // Gas limit
	FromAddr address.EVMAddress // Sender address
	ToAddr   *address.EVMAddress // Recipient (nil for contract creation)
	ValueVal *big.Int // Value in wei
	DataVal  []byte   // Contract bytecode or call data
	
	// Signature (ECDSA)
	V *big.Int
	R *big.Int
	S *big.Int
}

// Hash computes the transaction hash
func (tx *EVMTransaction) Hash() [32]byte {
	data := make([]byte, 0, 256)
	
	// Type
	data = append(data, byte(tx.TypeFlag))
	
	// Nonce
	nonceBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		nonceBytes[i] = byte(tx.NonceVal >> (8 * (7 - i)))
	}
	data = append(data, nonceBytes...)
	
	// GasPrice
	if tx.GasPriceVal != nil {
		data = append(data, tx.GasPriceVal.Bytes()...)
	}
	
	// GasLimit
	gasLimitBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		gasLimitBytes[i] = byte(tx.GasLimitVal >> (8 * (7 - i)))
	}
	data = append(data, gasLimitBytes...)
	
	// From
	data = append(data, tx.FromAddr.Bytes()...)
	
	// To (if not nil)
	if tx.ToAddr != nil {
		data = append(data, tx.ToAddr.Bytes()...)
	}
	
	// Value
	if tx.ValueVal != nil {
		data = append(data, tx.ValueVal.Bytes()...)
	}
	
	// Data
	data = append(data, tx.DataVal...)
	
	return Hash256(data)
}

// Type returns the transaction type
func (tx *EVMTransaction) Type() TxType {
	return tx.TypeFlag
}

// From returns the sender address
func (tx *EVMTransaction) From() address.EVMAddress {
	return tx.FromAddr
}

// Nonce returns the account nonce
func (tx *EVMTransaction) Nonce() uint64 {
	return tx.NonceVal
}

// Gas returns the gas limit
func (tx *EVMTransaction) Gas() uint64 {
	return tx.GasLimitVal
}

// GasPrice returns the gas price
func (tx *EVMTransaction) GasPrice() *big.Int {
	return tx.GasPriceVal
}

// To returns the recipient address (nil for contract creation)
func (tx *EVMTransaction) To() *address.EVMAddress {
	return tx.ToAddr
}

// Value returns the transaction value
func (tx *EVMTransaction) Value() *big.Int {
	return tx.ValueVal
}

// Data returns the transaction data
func (tx *EVMTransaction) Data() []byte {
	return tx.DataVal
}

// IsContractCreation returns true if this is a contract deployment
func (tx *EVMTransaction) IsContractCreation() bool {
	return tx.ToAddr == nil
}

// String returns a human-readable representation
func (tx *EVMTransaction) String() string {
	toStr := "CREATE"
	if tx.ToAddr != nil {
		toStr = tx.ToAddr.Hex()[:10] + "..."
	}
	
	return fmt.Sprintf(
		"EVMTx{Type: %d, From: %s, To: %s, Value: %s, Gas: %d, Nonce: %d}",
		tx.TypeFlag,
		tx.FromAddr.Hex()[:10]+"...",
		toStr,
		tx.ValueVal.String(),
		tx.GasLimitVal,
		tx.NonceVal,
	)
}

// Receipt represents a transaction receipt
type Receipt struct {
	TxHash          [32]byte           // Transaction hash
	BlockHeight     uint64             // Block height
	TxIndex         uint32             // Transaction index in block
	From            address.EVMAddress // Sender
	To              *address.EVMAddress // Recipient (nil for contract creation)
	ContractAddress *address.EVMAddress // Created contract address (nil if not a deployment)
	GasUsed         uint64             // Gas consumed
	Status          uint8              // 1 = success, 0 = failure
	Logs            []Log              // Event logs
	
	// Cumulative gas used in the block up to and including this transaction
	CumulativeGasUsed uint64
}

// Log represents an EVM event log
type Log struct {
	Address address.EVMAddress // Contract that emitted the log
	Topics  [][32]byte         // Indexed topics
	Data    []byte             // Non-indexed data
}

// Hash computes the receipt hash
func (r *Receipt) Hash() [32]byte {
	data := make([]byte, 0, 256)
	
	data = append(data, r.TxHash[:]...)
	
	heightBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		heightBytes[i] = byte(r.BlockHeight >> (8 * (7 - i)))
	}
	data = append(data, heightBytes...)
	
	data = append(data, byte(r.Status))
	
	gasUsedBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		gasUsedBytes[i] = byte(r.GasUsed >> (8 * (7 - i)))
	}
	data = append(data, gasUsedBytes...)
	
	return Hash256(data)
}

// String returns a human-readable representation
func (r *Receipt) String() string {
	status := "✅ Success"
	if r.Status == 0 {
		status = "❌ Failed"
	}
	
	return fmt.Sprintf(
		"Receipt{TxHash: %s, Block: %d, Status: %s, GasUsed: %d, Logs: %d}",
		hex.EncodeToString(r.TxHash[:8]),
		r.BlockHeight,
		status,
		r.GasUsed,
		len(r.Logs),
	)
}

