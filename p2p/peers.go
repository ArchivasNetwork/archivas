package p2p

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/iljanemesis/archivas/storage"
)

// PeerInfo represents a known peer
type PeerInfo struct {
	Address  string    `json:"address"`
	LastSeen time.Time `json:"lastSeen"`
	Score    int       `json:"score"` // Reputation score
}

// PeerStore manages known peers with persistence
type PeerStore struct {
	sync.RWMutex
	peers map[string]*PeerInfo
	db    *storage.DB
}

// NewPeerStore creates a new peer store
func NewPeerStore(db *storage.DB) *PeerStore {
	ps := &PeerStore{
		peers: make(map[string]*PeerInfo),
		db:    db,
	}
	
	// Load peers from database
	ps.loadFromDB()
	
	return ps
}

// Add adds or updates a peer
func (ps *PeerStore) Add(address string) {
	ps.Lock()
	defer ps.Unlock()
	
	if peer, exists := ps.peers[address]; exists {
		peer.LastSeen = time.Now()
		peer.Score++
	} else {
		ps.peers[address] = &PeerInfo{
			Address:  address,
			LastSeen: time.Now(),
			Score:    1,
		}
	}
	
	// Persist to DB
	ps.saveToDB(address)
}

// Remove removes a peer
func (ps *PeerStore) Remove(address string) {
	ps.Lock()
	defer ps.Unlock()
	
	delete(ps.peers, address)
	
	// Remove from DB
	key := []byte("peer:" + address)
	ps.db.Delete(key)
}

// All returns all known peers
func (ps *PeerStore) All() []string {
	ps.RLock()
	defer ps.RUnlock()
	
	addrs := make([]string, 0, len(ps.peers))
	for addr := range ps.peers {
		addrs = append(addrs, addr)
	}
	return addrs
}

// Recent returns peers seen in last N minutes
func (ps *PeerStore) Recent(minutes int) []string {
	ps.RLock()
	defer ps.RUnlock()
	
	cutoff := time.Now().Add(-time.Duration(minutes) * time.Minute)
	recent := make([]string, 0)
	
	for addr, info := range ps.peers {
		if info.LastSeen.After(cutoff) {
			recent = append(recent, addr)
		}
	}
	
	return recent
}

// Count returns number of known peers
func (ps *PeerStore) Count() int {
	ps.RLock()
	defer ps.RUnlock()
	return len(ps.peers)
}

// saveToDB persists a peer to database
func (ps *PeerStore) saveToDB(address string) {
	peer, ok := ps.peers[address]
	if !ok {
		return
	}
	
	key := []byte("peer:" + address)
	data, _ := json.Marshal(peer)
	ps.db.Put(key, data)
}

// loadFromDB loads all peers from database
func (ps *PeerStore) loadFromDB() {
	// TODO: Implement database scan for peer: prefix
	// For now, peers are ephemeral (will be added on connect)
}

// GossipMessage contains peer addresses to share
type GossipPeersMessage struct {
	Peers []string `json:"peers"`
}

// GetGossipPeers returns peers to share with others
func (ps *PeerStore) GetGossipPeers(max int) []string {
	recent := ps.Recent(60) // Peers seen in last hour
	
	if len(recent) <= max {
		return recent
	}
	
	// Return up to max peers
	return recent[:max]
}

// MergePeers merges received peers into our store
func (ps *PeerStore) MergePeers(addresses []string) []string {
	ps.Lock()
	defer ps.Unlock()
	
	newPeers := make([]string, 0)
	
	for _, addr := range addresses {
		if _, exists := ps.peers[addr]; !exists {
			ps.peers[addr] = &PeerInfo{
				Address:  addr,
				LastSeen: time.Now(),
				Score:    0,
			}
			newPeers = append(newPeers, addr)
			ps.saveToDB(addr)
		}
	}
	
	return newPeers
}

