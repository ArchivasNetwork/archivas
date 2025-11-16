package types

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/pospace"
)

// Block represents a blockchain block with Proof-of-Space and EVM execution
// Phase 2: Extended with EVM fields for Betanet
type Block struct {
	// PoST Consensus fields
	Height        uint64
	TimestampUnix int64
	PrevHash      [32]byte
	Difficulty    uint64         // Difficulty target when mined
	Challenge     [32]byte       // The challenge used to win this block
	Proof         *pospace.Proof // Proof-of-Space
	FarmerAddr    address.EVMAddress // Address to receive block reward

	// Cumulative work for fork resolution
	CumulativeWork uint64 // Total work from genesis to this block

	// EVM Execution fields (Phase 2)
	StateRoot    [32]byte // Root hash of the world state after executing this block
	ReceiptsRoot [32]byte // Root hash of transaction receipts
	GasUsed      uint64   // Total gas consumed by transactions in this block
	GasLimit     uint64   // Maximum gas allowed in this block

	// Transactions (both legacy and EVM)
	Txs []Transaction // Transactions in this block

	// Optional: Bloom filter for log searching (can be computed from receipts)
	// LogsBloom [256]byte
}

// BlockHeader represents just the header portion of a block
type BlockHeader struct {
	Height         uint64
	TimestampUnix  int64
	PrevHash       [32]byte
	Difficulty     uint64
	Challenge      [32]byte
	FarmerAddr     address.EVMAddress
	CumulativeWork uint64
	StateRoot      [32]byte
	ReceiptsRoot   [32]byte
	GasUsed        uint64
	GasLimit       uint64
	TxRoot         [32]byte // Merkle root of transactions
}

// Hash computes the block hash
// Includes all PoST fields + EVM state roots for integrity
func (b *Block) Hash() [32]byte {
	header := b.Header()
	return header.Hash()
}

// Header extracts the header from a full block
func (b *Block) Header() *BlockHeader {
	return &BlockHeader{
		Height:         b.Height,
		TimestampUnix:  b.TimestampUnix,
		PrevHash:       b.PrevHash,
		Difficulty:     b.Difficulty,
		Challenge:      b.Challenge,
		FarmerAddr:     b.FarmerAddr,
		CumulativeWork: b.CumulativeWork,
		StateRoot:      b.StateRoot,
		ReceiptsRoot:   b.ReceiptsRoot,
		GasUsed:        b.GasUsed,
		GasLimit:       b.GasLimit,
		TxRoot:         b.computeTxRoot(),
	}
}

// Hash computes the header hash
func (h *BlockHeader) Hash() [32]byte {
	// Serialize header for hashing
	// Format: Height|Timestamp|PrevHash|Difficulty|Challenge|FarmerAddr|
	//         CumulativeWork|StateRoot|ReceiptsRoot|GasUsed|GasLimit|TxRoot
	
	data := make([]byte, 0, 256)
	
	// Height (8 bytes)
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, h.Height)
	data = append(data, heightBytes...)
	
	// Timestamp (8 bytes)
	tsBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tsBytes, uint64(h.TimestampUnix))
	data = append(data, tsBytes...)
	
	// PrevHash (32 bytes)
	data = append(data, h.PrevHash[:]...)
	
	// Difficulty (8 bytes)
	diffBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(diffBytes, h.Difficulty)
	data = append(data, diffBytes...)
	
	// Challenge (32 bytes)
	data = append(data, h.Challenge[:]...)
	
	// FarmerAddr (20 bytes)
	data = append(data, h.FarmerAddr.Bytes()...)
	
	// CumulativeWork (8 bytes)
	workBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(workBytes, h.CumulativeWork)
	data = append(data, workBytes...)
	
	// StateRoot (32 bytes)
	data = append(data, h.StateRoot[:]...)
	
	// ReceiptsRoot (32 bytes)
	data = append(data, h.ReceiptsRoot[:]...)
	
	// GasUsed (8 bytes)
	gasUsedBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(gasUsedBytes, h.GasUsed)
	data = append(data, gasUsedBytes...)
	
	// GasLimit (8 bytes)
	gasLimitBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(gasLimitBytes, h.GasLimit)
	data = append(data, gasLimitBytes...)
	
	// TxRoot (32 bytes)
	data = append(data, h.TxRoot[:]...)
	
	// Compute SHA256 hash
	return Hash256(data)
}

// computeTxRoot computes the merkle root of transactions
// Simplified version - in production, use proper Merkle tree
func (b *Block) computeTxRoot() [32]byte {
	if len(b.Txs) == 0 {
		return [32]byte{} // Empty root
	}
	
	// Simple approach: hash concatenated tx hashes
	// TODO: Implement proper Merkle tree
	data := make([]byte, 0)
	for _, tx := range b.Txs {
		txHash := tx.Hash()
		data = append(data, txHash[:]...)
	}
	
	return Hash256(data)
}

// String returns a human-readable representation
func (b *Block) String() string {
	return fmt.Sprintf(
		"Block{Height: %d, Time: %d, PrevHash: %s, Farmer: %s, Txs: %d, Gas: %d/%d, StateRoot: %s}",
		b.Height,
		b.TimestampUnix,
		hex.EncodeToString(b.PrevHash[:8]),
		b.FarmerAddr.Hex()[:10]+"...",
		len(b.Txs),
		b.GasUsed,
		b.GasLimit,
		hex.EncodeToString(b.StateRoot[:8]),
	)
}

// IsEVMEnabled returns true if this block supports EVM execution
// Determined by protocol version from genesis
func (b *Block) IsEVMEnabled() bool {
	// For Betanet (protocol v2), EVM is always enabled
	// For devnet-legacy (protocol v1), EVM is disabled
	// This will be set based on the network profile at runtime
	return b.GasLimit > 0
}

// EmptyStateRoot returns the hash of an empty state trie
// This is the Keccak256 hash of an empty trie
func EmptyStateRoot() [32]byte {
	// This is the well-known empty trie hash in Ethereum
	// 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421
	root := [32]byte{}
	bytes, _ := hex.DecodeString("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	copy(root[:], bytes)
	return root
}

// EmptyReceiptsRoot returns the hash of an empty receipts trie
func EmptyReceiptsRoot() [32]byte {
	return EmptyStateRoot() // Same as empty state root
}

