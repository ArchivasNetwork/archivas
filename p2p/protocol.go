package p2p

import (
	"encoding/json"
)

// MessageType represents different P2P message types
type MessageType uint8

const (
	MsgTypeHandshake   MessageType = 0  // v1.1.1: First message, validates compatibility
	MsgTypePing        MessageType = 1
	MsgTypePong        MessageType = 2
	MsgTypeNewBlock    MessageType = 3
	MsgTypeGetBlock    MessageType = 4
	MsgTypeBlockData   MessageType = 5
	MsgTypeGetStatus   MessageType = 6
	MsgTypeStatus      MessageType = 7
	MsgTypeGossipPeers MessageType = 8
	// v0.5.0: Enhanced gossip
	MsgTypeInv         MessageType = 9  // Inventory announcement
	MsgTypeReq         MessageType = 10 // Request data
	MsgTypeRes         MessageType = 11 // Response with data
	MsgTypeTxBroadcast MessageType = 12 // Transaction broadcast
	// v1.1.1: Efficient IBD
	MsgTypeRequestBlocks MessageType = 13 // Request block range
	MsgTypeBlocksBatch   MessageType = 14 // Batch of blocks
)

// Message represents a P2P protocol message
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// PingMessage is sent to check if peer is alive
type PingMessage struct {
	Timestamp int64 `json:"timestamp"`
}

// PongMessage is the response to Ping
type PongMessage struct {
	Timestamp int64 `json:"timestamp"`
}

// NewBlockMessage announces a new block
type NewBlockMessage struct {
	Height uint64   `json:"height"`
	Hash   [32]byte `json:"hash"`
}

// GetBlockMessage requests a specific block
type GetBlockMessage struct {
	Height uint64 `json:"height"`
}

// BlockDataMessage contains full block data
type BlockDataMessage struct {
	Height    uint64          `json:"height"`
	BlockJSON json.RawMessage `json:"blockData"`
}

// GetStatusMessage requests peer's current status
type GetStatusMessage struct{}

// StatusMessage contains peer's chain status
type StatusMessage struct {
	Height     uint64   `json:"height"`
	Difficulty uint64   `json:"difficulty"`
	TipHash    [32]byte `json:"tipHash"`
}

// GossipPeersMessage is sent to share known peer addresses with network validation
type GossipPeersMessage struct {
	Addrs  []string `json:"addrs"`  // "ip:port" addresses
	SeenAt int64    `json:"seenAt"` // unix timestamp
	NetID  string   `json:"netId"`  // network identifier for validation
}

// InvMessage announces available data (blocks or txs)
type InvMessage struct {
	Type  string   `json:"type"` // "block" or "tx"
	Hashes []string `json:"hashes"` // hex-encoded hashes
}

// ReqMessage requests specific data
type ReqMessage struct {
	Type  string   `json:"type"` // "block" or "tx"
	Hashes []string `json:"hashes"` // hex-encoded hashes
}

// ResMessage contains requested data
type ResMessage struct {
	Type string          `json:"type"` // "block" or "tx"
	Data json.RawMessage `json:"data"` // actual data
}

// TxBroadcastMessage broadcasts a transaction
type TxBroadcastMessage struct {
	TxHash string          `json:"txHash"`
	TxData json.RawMessage `json:"txData"`
}

// RequestBlocksMessage requests a range of blocks for IBD
// v1.1.1: Batched sync for efficient initial block download
type RequestBlocksMessage struct {
	FromHeight uint64 `json:"fromHeight"` // Starting height (inclusive)
	MaxBlocks  uint32 `json:"maxBlocks"`  // Batch size hint (capped at 512)
}

// BlocksBatchMessage contains a batch of blocks
// v1.1.1: Response to RequestBlocks for efficient IBD
type BlocksBatchMessage struct {
	FromHeight uint64            `json:"fromHeight"` // First block height in this batch
	Count      uint32            `json:"count"`      // Number of blocks in this batch
	Blocks     []json.RawMessage `json:"blocks"`     // Block data (JSON serialized)
	TipHeight  uint64            `json:"tipHeight"`  // Sender's current tip (for progress tracking)
	EOF        bool              `json:"eof"`        // true if this is the last batch (caught up)
}

// HandshakeMessage validates peer compatibility before allowing connection
// v1.1.1: Prevents incompatible nodes from connecting
type HandshakeMessage struct {
	GenesisHash        [32]byte `json:"genesisHash"`        // Must match exactly
	NetworkID          string   `json:"networkID"`          // Must match exactly
	ProtocolVersion    string   `json:"protocolVersion"`    // Must be compatible
	DifficultyParamsID string   `json:"difficultyParamsID"` // Must match exactly
	NodeVersion        string   `json:"nodeVersion"`        // For logging/debugging
}
