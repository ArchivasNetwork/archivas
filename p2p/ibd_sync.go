package p2p

import (
	"log"

	"github.com/ArchivasNetwork/archivas/metrics"
)

// StartIBD initiates Initial Block Download from connected peers
// v1.1.1: Batched sync for catching up to network tip
func (n *Network) StartIBD(fromHeight uint64) {
	n.RLock()
	if len(n.peers) == 0 {
		n.RUnlock()
		log.Printf("[sync] no peers available for IBD")
		return
	}
	
	// Pick first available peer (could be improved with peer selection logic)
	var syncPeer *Peer
	for _, p := range n.peers {
		syncPeer = p
		break
	}
	n.RUnlock()

	if syncPeer == nil {
		log.Printf("[sync] no sync peer available")
		return
	}

	log.Printf("[sync] starting IBD from height %d via peer %s", fromHeight, syncPeer.Address)

	// Send initial request
	req := RequestBlocksMessage{
		FromHeight: fromHeight,
		MaxBlocks:  uint32(n.ibdBatchSize),
	}

	metrics.IncIBDRequestedBatches()
	n.SendMessage(syncPeer, MsgTypeRequestBlocks, req)
}

