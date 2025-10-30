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
	syncState   *SyncState
	peerStore   PeerStore
}

// NodeHandler interface for node callbacks
type NodeHandler interface {
	OnNewBlock(height uint64, hash [32]byte, fromPeer string)
	OnBlockRequest(height uint64) (interface{}, error)
	GetStatus() (height uint64, difficulty uint64, tipHash [32]byte)
	// Block importing
	LocalHeight() uint64
	HasBlock(height uint64) bool
	VerifyAndApplyBlock(blockJSON json.RawMessage) error
}

// NewNetwork creates a new P2P network
func NewNetwork(listenAddr string, handler NodeHandler) *Network {
	return &Network{
		peers:       make(map[string]*Peer),
		nodeHandler: handler,
		listenAddr:  listenAddr,
		syncState:   NewSyncState(),
		peerStore:   nil, // Will be set via SetPeerStore
	}
}

// SetPeerStore sets the peer store for persistence
func (n *Network) SetPeerStore(store PeerStore) {
	n.Lock()
	defer n.Unlock()
	n.peerStore = store
	
	// Auto-dial stored peers
	if store != nil {
		go func() {
			time.Sleep(2 * time.Second)
			peers, _ := store.List()
			for _, addr := range peers {
				if addr != "" {
					n.ConnectPeer(addr)
				}
			}
		}()
		
		// Start peer gossip loop
		go n.peerGossipLoop()
	}
}

// peerGossipLoop periodically shares known peers with connected peers
func (n *Network) peerGossipLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		n.gossipPeers()
	}
}

// gossipPeers shares our known peers with connected peers
func (n *Network) gossipPeers() {
	// Get list of known peers from store
	if n.peerStore == nil {
		return
	}
	
	knownPeers, err := n.peerStore.List()
	if err != nil || len(knownPeers) == 0 {
		return
	}
	
	// Create gossip message
	msg := GossipPeersMessage{
		Peers: knownPeers,
	}
	
	// Send to all connected peers
	n.RLock()
	peers := make([]*Peer, 0, len(n.peers))
	for _, peer := range n.peers {
		peers = append(peers, peer)
	}
	n.RUnlock()
	
	for _, peer := range peers {
		if err := n.SendMessage(peer, MsgTypeGossipPeers, msg); err != nil {
			log.Printf("[p2p] failed to send GOSSIP_PEERS to %s: %v", peer.Address, err)
		}
	}
	
	log.Printf("[p2p] gossiped %d known peers to %d connected peers", len(knownPeers), len(peers))
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

	// CRITICAL: Register peer BEFORE starting handler
	n.Lock()
	n.peers[address] = peer
	peerCount := len(n.peers)
	
	// Persist to peer store
	if n.peerStore != nil {
		n.peerStore.Add(address)
	}
	n.Unlock()

	log.Printf("[p2p] connected to peer %s (total peers: %d, persisted)", address, peerCount)

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

		// CRITICAL: Register peer BEFORE starting handler
		n.Lock()
		n.peers[peer.Address] = peer
		peerCount := len(n.peers)
		n.Unlock()

		log.Printf("[p2p] accepted connection from %s (total peers: %d)", peer.Address, peerCount)

		go n.handlePeer(peer)
	}
}

// handlePeer handles messages from a peer
func (n *Network) handlePeer(peer *Peer) {
	defer func() {
		peer.Conn.Close()
		n.Lock()
		delete(n.peers, peer.Address)
		peerCount := len(n.peers)
		n.Unlock()
		
		// Don't remove from store on disconnect (will retry later)
		log.Printf("[p2p] peer %s disconnected (remaining peers: %d)", peer.Address, peerCount)
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

	// CRITICAL: Use snapshot to avoid holding lock during sends
	n.RLock()
	peers := make([]*Peer, 0, len(n.peers))
	for _, peer := range n.peers {
		peers = append(peers, peer)
	}
	peerCount := len(n.peers)
	n.RUnlock()

	// Log peer count BEFORE sending
	log.Printf("[p2p] broadcasting block %d to %d peers (map has %d entries)", height, len(peers), peerCount)

	for _, peer := range peers {
		if err := n.SendMessage(peer, MsgTypeNewBlock, msg); err != nil {
			log.Printf("[p2p] failed to send NEW_BLOCK to %s: %v", peer.Address, err)
		} else {
			log.Printf("[p2p] sent NEW_BLOCK height=%d to %s", height, peer.Address)
		}
	}

	log.Printf("[p2p] broadcast complete: block %d sent to %d/%d peers", height, len(peers), peerCount)
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
	case MsgTypeGossipPeers:
		n.handleGossipPeers(peer, msg.Payload)
	default:
		log.Printf("[p2p] unknown message type %d from %s", msg.Type, peer.Address)
	}
}

// handleGossipPeers handles received peer gossip
func (n *Network) handleGossipPeers(peer *Peer, payload json.RawMessage) {
	var gossip GossipPeersMessage
	if err := json.Unmarshal(payload, &gossip); err != nil {
		log.Printf("[p2p] invalid GOSSIP_PEERS from %s", peer.Address)
		return
	}
	
	log.Printf("[p2p] received %d peers from %s", len(gossip.Peers), peer.Address)
	
	// Connect to new peers (max 20 total, avoid self)
	n.RLock()
	currentCount := len(n.peers)
	n.RUnlock()
	
	if currentCount >= 20 {
		return // Already at max
	}
	
	for _, addr := range gossip.Peers {
		// Skip if already connected
		n.RLock()
		_, exists := n.peers[addr]
		n.RUnlock()
		
		if exists || addr == n.listenAddr || addr == "" {
			continue
		}
		
		// Try to connect
		go func(a string) {
			time.Sleep(time.Duration(len(a)%5) * time.Second) // Stagger
			n.ConnectPeer(a)
		}(addr)
		
		// Don't overwhelm - connect slowly
		if currentCount >= 20 {
			break
		}
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
	
	// Try to apply the block
	if err := n.nodeHandler.VerifyAndApplyBlock(blockData.BlockJSON); err != nil {
		log.Printf("[p2p] failed to apply block %d: %v", blockData.Height, err)
		
		// If block is out of order, queue it and request missing blocks
		localHeight := n.nodeHandler.LocalHeight()
		if blockData.Height > localHeight+1 {
			log.Printf("[p2p] block %d is ahead of us (height=%d), requesting gap blocks", blockData.Height, localHeight)
			
			// Request missing blocks
			for h := localHeight + 1; h < blockData.Height; h++ {
				if !n.nodeHandler.HasBlock(h) && n.syncState.WantBlock(h) {
					req := GetBlockMessage{Height: h}
					n.SendMessage(peer, MsgTypeGetBlock, req)
				}
			}
			
			// Queue this block for later
			n.syncState.QueueBlock(blockData.Height, blockData.BlockJSON)
		}
		return
	}
	
	n.syncState.GotBlock(blockData.Height)
	
	// Try to apply any queued blocks that are now sequential
	for {
		nextHeight := n.nodeHandler.LocalHeight() + 1
		if queuedData, ok := n.syncState.GetQueuedBlock(nextHeight); ok {
			if err := n.nodeHandler.VerifyAndApplyBlock(queuedData); err != nil {
				log.Printf("[p2p] failed to apply queued block %d: %v", nextHeight, err)
				break
			}
			log.Printf("[p2p] ✅ Applied queued block %d", nextHeight)
		} else {
			break
		}
	}
	
	// If peer is still ahead, request next block
	localHeight := n.nodeHandler.LocalHeight()
	if peer.Height > localHeight {
		nextHeight := localHeight + 1
		if n.syncState.WantBlock(nextHeight) {
			req := GetBlockMessage{Height: nextHeight}
			n.SendMessage(peer, MsgTypeGetBlock, req)
		}
	}
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
