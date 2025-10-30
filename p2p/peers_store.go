package p2p

import (
	"encoding/json"
	"os"
	"sync"
)

// PeerStore interface for peer persistence
type PeerStore interface {
	Add(addr string) error
	Remove(addr string) error
	List() ([]string, error)
}

// FilePeerStore implements PeerStore using a JSON file
type FilePeerStore struct {
	mu   sync.Mutex
	path string
	data *peerStoreData
}

type peerStoreData struct {
	Peers []string `json:"peers"`
}

// NewFilePeerStore creates a new file-based peer store
func NewFilePeerStore(path string) (PeerStore, error) {
	store := &FilePeerStore{
		path: path,
		data: &peerStoreData{
			Peers: make([]string, 0),
		},
	}
	
	// Load existing peers if file exists
	if err := store.load(); err != nil {
		// File doesn't exist yet, that's OK
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	
	return store, nil
}

// Add adds a peer address
func (s *FilePeerStore) Add(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Check if already exists
	for _, p := range s.data.Peers {
		if p == addr {
			return nil // Already have it
		}
	}
	
	s.data.Peers = append(s.data.Peers, addr)
	return s.save()
}

// Remove removes a peer address
func (s *FilePeerStore) Remove(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	filtered := make([]string, 0)
	for _, p := range s.data.Peers {
		if p != addr {
			filtered = append(filtered, p)
		}
	}
	
	s.data.Peers = filtered
	return s.save()
}

// List returns all known peer addresses
func (s *FilePeerStore) List() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Return copy
	peers := make([]string, len(s.data.Peers))
	copy(peers, s.data.Peers)
	return peers, nil
}

// load reads peers from disk
func (s *FilePeerStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, s.data)
}

// save writes peers to disk (atomic)
func (s *FilePeerStore) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	
	// Atomic write: write to temp, then rename
	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}
	
	return os.Rename(tmpPath, s.path)
}

