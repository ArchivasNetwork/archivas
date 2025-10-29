package p2p

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// BlockImporter interface for importing blocks from peers
type BlockImporter interface {
	LocalHeight() uint64
	LocalTipHash() [32]byte
	HasBlock(height uint64) bool
	GetBlock(height uint64) (interface{}, error)
	VerifyAndApplyBlock(blockJSON json.RawMessage) error
}

// SyncState tracks ongoing sync operations
type SyncState struct {
	sync.Mutex
	importing     bool
	wantedHeights map[uint64]bool
	queuedBlocks  map[uint64]json.RawMessage
}

// NewSyncState creates a new sync state
func NewSyncState() *SyncState {
	return &SyncState{
		wantedHeights: make(map[uint64]bool),
		queuedBlocks:  make(map[uint64]json.RawMessage),
	}
}

// WantBlock marks a block height as wanted
func (s *SyncState) WantBlock(height uint64) bool {
	s.Lock()
	defer s.Unlock()
	
	if s.wantedHeights[height] {
		return false // already wanting
	}
	s.wantedHeights[height] = true
	return true
}

// GotBlock marks a block as received
func (s *SyncState) GotBlock(height uint64) {
	s.Lock()
	defer s.Unlock()
	delete(s.wantedHeights, height)
}

// QueueBlock stores an out-of-order block
func (s *SyncState) QueueBlock(height uint64, blockData json.RawMessage) {
	s.Lock()
	defer s.Unlock()
	s.queuedBlocks[height] = blockData
}

// GetQueuedBlock retrieves a queued block
func (s *SyncState) GetQueuedBlock(height uint64) (json.RawMessage, bool) {
	s.Lock()
	defer s.Unlock()
	data, ok := s.queuedBlocks[height]
	if ok {
		delete(s.queuedBlocks, height)
	}
	return data, ok
}

// ProcessBlockData handles received block data
func ProcessBlockData(importer BlockImporter, syncState *SyncState, height uint64, blockData json.RawMessage) error {
	localHeight := importer.LocalHeight()
	
	log.Printf("[sync] Received block %d (local height: %d)", height, localHeight)
	
	// If this block is next in sequence, try to apply it
	if height == localHeight+1 {
		if err := importer.VerifyAndApplyBlock(blockData); err != nil {
			return fmt.Errorf("failed to apply block %d: %w", height, err)
		}
		log.Printf("[sync] ✅ Applied block %d", height)
		syncState.GotBlock(height)
		
		// Try to apply any queued blocks that are now sequential
		for {
			nextHeight := importer.LocalHeight() + 1
			if queuedData, ok := syncState.GetQueuedBlock(nextHeight); ok {
				log.Printf("[sync] Applying queued block %d", nextHeight)
				if err := importer.VerifyAndApplyBlock(queuedData); err != nil {
					log.Printf("[sync] ⚠️  Failed to apply queued block %d: %v", nextHeight, err)
					break
				}
				log.Printf("[sync] ✅ Applied queued block %d", nextHeight)
			} else {
				break
			}
		}
	} else if height > localHeight+1 {
		// Out of order - queue it and request missing blocks
		log.Printf("[sync] Queueing out-of-order block %d (gap: %d blocks)", height, height-localHeight-1)
		syncState.QueueBlock(height, blockData)
		
		// Request missing blocks in gap
		for h := localHeight + 1; h < height; h++ {
			if syncState.WantBlock(h) {
				log.Printf("[sync] Requesting missing block %d", h)
				// Caller should send GET_BLOCK for height h
			}
		}
	} else {
		// Already have this block or it's behind us
		log.Printf("[sync] Skipping block %d (already have it or behind)", height)
		syncState.GotBlock(height)
	}
	
	return nil
}

