package p2p

import (
	"encoding/json"
)

// MessageType represents different P2P message types
type MessageType uint8

const (
	MsgTypePing        MessageType = 1
	MsgTypePong        MessageType = 2
	MsgTypeNewBlock    MessageType = 3
	MsgTypeGetBlock    MessageType = 4
	MsgTypeBlockData   MessageType = 5
	MsgTypeGetStatus   MessageType = 6
	MsgTypeStatus      MessageType = 7
	MsgTypeGossipPeers MessageType = 8
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
