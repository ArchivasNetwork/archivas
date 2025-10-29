package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Peer represents a connected peer
type Peer struct {
	Address    string
	Conn       net.Conn
	LastSeen   time.Time
	Height     uint64
	Reader     *bufio.Reader
	Writer     *bufio.Writer
	writeMutex sync.Mutex
}

// Network handles peer-to-peer networking
type Network struct {
	sync.RWMutex
	peers       map[string]*Peer
	listener    net.Listener
	nodeHandler NodeHandler
	listenAddr  string
}

// NodeHandler interface for node callbacks
type NodeHandler interface {
	OnNewBlock(height uint64, hash [32]byte, fromPeer string)
	OnBlockRequest(height uint64) (interface{}, error)
	GetStatus() (height uint64, difficulty uint64, tipHash [32]byte)
}

// NewNetwork creates a new P2P network
func NewNetwork(listenAddr string, handler NodeHandler) *Network {
	return &Network{
		peers:       make(map[string]*Peer),
		nodeHandler: handler,
		listenAddr:  listenAddr,
	}
}

// Start starts the P2P network listener
func (n *Network) Start() error {
	listener, err := net.Listen("tcp", n.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start P2P listener: %w", err)
	}

	n.listener = listener
	log.Printf("[p2p] listening on %s", n.listenAddr)

	// Accept incoming connections
	go n.acceptLoop()

	return nil
}

// Stop stops the P2P network
func (n *Network) Stop() error {
	if n.listener != nil {
		return n.listener.Close()
	}
	return nil
}

// ConnectPeer connects to a remote peer
func (n *Network) ConnectPeer(address string) error {
	log.Printf("[p2p] connecting to peer %s", address)

	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to peer: %w", err)
	}

	peer := &Peer{
		Address:  address,
		Conn:     conn,
		LastSeen: time.Now(),
		Reader:   bufio.NewReader(conn),
		Writer:   bufio.NewWriter(conn),
	}

	n.Lock()
	n.peers[address] = peer
	n.Unlock()

	log.Printf("[p2p] connected to peer %s", address)

	// Start handling messages from this peer
	go n.handlePeer(peer)

	// Send initial status request
	go func() {
		time.Sleep(1 * time.Second)
		n.SendMessage(peer, MsgTypeGetStatus, GetStatusMessage{})
	}()

	return nil
}

// acceptLoop accepts incoming connections
func (n *Network) acceptLoop() {
	for {
		conn, err := n.listener.Accept()
		if err != nil {
			log.Printf("[p2p] accept error: %v", err)
			return
		}

		peer := &Peer{
			Address:  conn.RemoteAddr().String(),
			Conn:     conn,
			LastSeen: time.Now(),
			Reader:   bufio.NewReader(conn),
			Writer:   bufio.NewWriter(conn),
		}

		n.Lock()
		n.peers[peer.Address] = peer
		n.Unlock()

		log.Printf("[p2p] accepted connection from %s", peer.Address)

		go n.handlePeer(peer)
	}
}

// handlePeer handles messages from a peer
func (n *Network) handlePeer(peer *Peer) {
	defer func() {
		peer.Conn.Close()
		n.Lock()
		delete(n.peers, peer.Address)
		n.Unlock()
		log.Printf("[p2p] peer %s disconnected", peer.Address)
	}()

	for {
		// Read message (newline-delimited JSON)
		line, err := peer.Reader.ReadBytes('\n')
		if err != nil {
			log.Printf("[p2p] read error from %s: %v", peer.Address, err)
			return
		}

		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			log.Printf("[p2p] invalid message from %s: %v", peer.Address, err)
			continue
		}

		peer.LastSeen = time.Now()

		// Handle message
		n.handleMessage(peer, &msg)
	}
}

// SendMessage sends a message to a peer
func (n *Network) SendMessage(peer *Peer, msgType MessageType, payload interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := Message{
		Type:    msgType,
		Payload: payloadJSON,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	peer.writeMutex.Lock()
	defer peer.writeMutex.Unlock()

	// Write newline-delimited JSON
	if _, err := peer.Writer.Write(append(data, '\n')); err != nil {
		return err
	}

	return peer.Writer.Flush()
}

// BroadcastNewBlock announces a new block to all peers
func (n *Network) BroadcastNewBlock(height uint64, hash [32]byte) {
	msg := NewBlockMessage{
		Height: height,
		Hash:   hash,
	}

	n.RLock()
	peers := make([]*Peer, 0, len(n.peers))
	for _, peer := range n.peers {
		peers = append(peers, peer)
	}
	n.RUnlock()

	for _, peer := range peers {
		if err := n.SendMessage(peer, MsgTypeNewBlock, msg); err != nil {
			log.Printf("[p2p] failed to send NEW_BLOCK to %s: %v", peer.Address, err)
		}
	}

	log.Printf("[p2p] broadcasted block %d to %d peers", height, len(peers))
}

// handleMessage processes incoming messages
func (n *Network) handleMessage(peer *Peer, msg *Message) {
	switch msg.Type {
	case MsgTypePing:
		n.handlePing(peer, msg.Payload)
	case MsgTypePong:
		n.handlePong(peer, msg.Payload)
	case MsgTypeNewBlock:
		n.handleNewBlock(peer, msg.Payload)
	case MsgTypeGetBlock:
		n.handleGetBlock(peer, msg.Payload)
	case MsgTypeBlockData:
		n.handleBlockData(peer, msg.Payload)
	case MsgTypeGetStatus:
		n.handleGetStatus(peer, msg.Payload)
	case MsgTypeStatus:
		n.handleStatus(peer, msg.Payload)
	default:
		log.Printf("[p2p] unknown message type %d from %s", msg.Type, peer.Address)
	}
}

func (n *Network) handlePing(peer *Peer, payload json.RawMessage) {
	var ping PingMessage
	json.Unmarshal(payload, &ping)

	pong := PongMessage{Timestamp: time.Now().Unix()}
	n.SendMessage(peer, MsgTypePong, pong)
}

func (n *Network) handlePong(peer *Peer, payload json.RawMessage) {
	// Peer is alive
	peer.LastSeen = time.Now()
}

func (n *Network) handleNewBlock(peer *Peer, payload json.RawMessage) {
	var newBlock NewBlockMessage
	if err := json.Unmarshal(payload, &newBlock); err != nil {
		log.Printf("[p2p] invalid NEW_BLOCK from %s", peer.Address)
		return
	}

	log.Printf("[p2p] received NEW_BLOCK height=%d from %s", newBlock.Height, peer.Address)

	// Notify node
	if n.nodeHandler != nil {
		n.nodeHandler.OnNewBlock(newBlock.Height, newBlock.Hash, peer.Address)
	}
}

func (n *Network) handleGetBlock(peer *Peer, payload json.RawMessage) {
	var getBlock GetBlockMessage
	if err := json.Unmarshal(payload, &getBlock); err != nil {
		log.Printf("[p2p] invalid GET_BLOCK from %s", peer.Address)
		return
	}

	log.Printf("[p2p] peer %s requested block %d", peer.Address, getBlock.Height)

	// Get block from node
	if n.nodeHandler != nil {
		blockData, err := n.nodeHandler.OnBlockRequest(getBlock.Height)
		if err != nil {
			log.Printf("[p2p] failed to get block %d: %v", getBlock.Height, err)
			return
		}

		blockJSON, _ := json.Marshal(blockData)
		response := BlockDataMessage{
			Height:    getBlock.Height,
			BlockJSON: blockJSON,
		}

		n.SendMessage(peer, MsgTypeBlockData, response)
	}
}

func (n *Network) handleBlockData(peer *Peer, payload json.RawMessage) {
	var blockData BlockDataMessage
	if err := json.Unmarshal(payload, &blockData); err != nil {
		log.Printf("[p2p] invalid BLOCK_DATA from %s", peer.Address)
		return
	}

	log.Printf("[p2p] received BLOCK_DATA height=%d from %s", blockData.Height, peer.Address)

	// TODO: Store this block (will be handled by sync logic)
}

func (n *Network) handleGetStatus(peer *Peer, payload json.RawMessage) {
	if n.nodeHandler != nil {
		height, difficulty, tipHash := n.nodeHandler.GetStatus()
		status := StatusMessage{
			Height:     height,
			Difficulty: difficulty,
			TipHash:    tipHash,
		}
		n.SendMessage(peer, MsgTypeStatus, status)
	}
}

func (n *Network) handleStatus(peer *Peer, payload json.RawMessage) {
	var status StatusMessage
	if err := json.Unmarshal(payload, &status); err != nil {
		log.Printf("[p2p] invalid STATUS from %s", peer.Address)
		return
	}

	log.Printf("[p2p] peer %s status: height=%d difficulty=%d", peer.Address, status.Height, status.Difficulty)
	peer.Height = status.Height
}

// GetPeerCount returns number of connected peers
func (n *Network) GetPeerCount() int {
	n.RLock()
	defer n.RUnlock()
	return len(n.peers)
}

// RequestBlock requests a specific block from any peer
func (n *Network) RequestBlock(height uint64) error {
	n.RLock()
	defer n.RUnlock()

	if len(n.peers) == 0 {
		return fmt.Errorf("no peers available")
	}

	// Request from first available peer
	for _, peer := range n.peers {
		req := GetBlockMessage{Height: height}
		return n.SendMessage(peer, MsgTypeGetBlock, req)
	}

	return fmt.Errorf("no peers available")
}

// SyncFromPeers syncs chain from peers
func (n *Network) SyncFromPeers(currentHeight uint64) (highestPeer uint64, err error) {
	n.RLock()
	defer n.RUnlock()

	highest := currentHeight

	// Find highest peer
	for _, peer := range n.peers {
		if peer.Height > highest {
			highest = peer.Height
		}
	}

	if highest == currentHeight {
		log.Println("[p2p] already at tip, no sync needed")
		return currentHeight, nil
	}

	log.Printf("[p2p] peers are ahead, need to sync from %d to %d", currentHeight, highest)
	return highest, nil
}
