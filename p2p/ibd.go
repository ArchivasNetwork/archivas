package p2p

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ArchivasNetwork/archivas/metrics"
)

// handleRequestBlocks serves a batch of blocks for IBD
// v1.1.1: Batched sync for efficient initial block download
func (n *Network) handleRequestBlocks(peer *Peer, payload json.RawMessage) {
	var req RequestBlocksMessage
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("[p2p] invalid REQUEST_BLOCKS from %s: %v", peer.Address, err)
		return
	}

	log.Printf("[p2p] peer %s requested blocks from=%d max=%d", peer.Address, req.FromHeight, req.MaxBlocks)

	// v1.1.1: Backpressure - limit concurrent IBD streams
	n.Lock()
	if n.ibdInflight >= n.ibdMaxConcurrent {
		n.Unlock()
		log.Printf("[p2p] IBD backpressure: rejecting request from %s (inflight=%d)", peer.Address, n.ibdInflight)
		// Send empty batch with EOF=false to signal busy
		emptyBatch := BlocksBatchMessage{
			FromHeight: req.FromHeight,
			Count:      0,
			Blocks:     []json.RawMessage{},
			TipHeight:  0,
			EOF:        false,
		}
		n.SendMessage(peer, MsgTypeBlocksBatch, emptyBatch)
		return
	}
	n.ibdInflight++
	metrics.UpdateIBDInflight(n.ibdInflight)
	n.Unlock()

	defer func() {
		n.Lock()
		n.ibdInflight--
		metrics.UpdateIBDInflight(n.ibdInflight)
		n.Unlock()
	}()

	// Cap batch size
	maxBlocks := req.MaxBlocks
	if maxBlocks == 0 || maxBlocks > 512 {
		maxBlocks = 512
	}

	// Validate fromHeight
	fromHeight := req.FromHeight
	if fromHeight < 1 {
		fromHeight = 1 // Genesis already known via handshake
	}

	// Get blocks from node handler
	if n.nodeHandler == nil {
		log.Printf("[p2p] no node handler, cannot serve blocks")
		return
	}

	blocks, tipHeight, eof, err := n.nodeHandler.OnBlocksRangeRequest(fromHeight, maxBlocks)
	if err != nil {
		log.Printf("[p2p] failed to get blocks from=%d: %v", fromHeight, err)
		return
	}

	// Rate limit to avoid starving RPC
	if len(blocks) > 0 {
		time.Sleep(5 * time.Millisecond) // 5ms per batch = ~100 batches/sec max
	}

	// Build response
	batch := BlocksBatchMessage{
		FromHeight: fromHeight,
		Count:      uint32(len(blocks)),
		Blocks:     blocks,
		TipHeight:  tipHeight,
		EOF:        eof,
	}

	log.Printf("[p2p] serving batch to %s: from=%d count=%d tip=%d eof=%v", 
		peer.Address, fromHeight, len(blocks), tipHeight, eof)

	n.SendMessage(peer, MsgTypeBlocksBatch, batch)
}

// handleBlocksBatch processes a batch of blocks received during IBD
// v1.1.1: Client-side IBD handler
func (n *Network) handleBlocksBatch(peer *Peer, payload json.RawMessage) {
	var batch BlocksBatchMessage
	if err := json.Unmarshal(payload, &batch); err != nil {
		log.Printf("[p2p] invalid BLOCKS_BATCH from %s: %v", peer.Address, err)
		return
	}

	log.Printf("[p2p] received batch from=%d count=%d tip=%d eof=%v from %s",
		batch.FromHeight, batch.Count, batch.TipHeight, batch.EOF, peer.Address)

	// Track metrics
	metrics.IncIBDReceivedBatches()

	if n.nodeHandler == nil {
		return
	}

	// Apply blocks sequentially
	localHeight := n.nodeHandler.LocalHeight()
	applied := 0

	for i, blockJSON := range batch.Blocks {
		expectedHeight := batch.FromHeight + uint64(i)
		
		// Verify sequential
		if expectedHeight != localHeight+1 {
			log.Printf("[p2p] block %d out of sequence (expected %d), skipping batch", 
				expectedHeight, localHeight+1)
			break
		}

		// Apply block
		if err := n.nodeHandler.VerifyAndApplyBlock(blockJSON); err != nil {
			log.Printf("[p2p] failed to apply block %d: %v", expectedHeight, err)
			break
		}

		applied++
		localHeight = n.nodeHandler.LocalHeight()
		
		// Log progress every 100 blocks
		if applied%100 == 0 {
			log.Printf("[sync] progress: applied %d blocks, now at height %d (tip: %d)", 
				applied, localHeight, batch.TipHeight)
		}
	}

	if applied > 0 {
		log.Printf("[sync] applied %d blocks from batch, now at height %d", applied, localHeight)
		metrics.IncIBDBlocksApplied(applied)
	}

	// If this was the last batch, we're caught up
	if batch.EOF && batch.Count == 0 {
		log.Printf("[sync] caught up to tip at height %d", localHeight)
		return
	}

	// If we're still behind, request next batch
	if localHeight < batch.TipHeight {
		nextReq := RequestBlocksMessage{
			FromHeight: localHeight + 1,
			MaxBlocks:  uint32(n.ibdBatchSize),
		}
		metrics.IncIBDRequestedBatches()
		n.SendMessage(peer, MsgTypeRequestBlocks, nextReq)
		log.Printf("[sync] requesting next batch from=%d", nextReq.FromHeight)
	}
}

