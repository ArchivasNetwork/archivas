package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

const (
	HeartbeatTTL = 120 * time.Second // Nodes must heartbeat every 2 minutes
)

type NodeRegistration struct {
	Address   string `json:"address"`   // Archivas address (arcv1...)
	P2PAddr   string `json:"p2pAddr"`   // host:port for P2P
	RPCAddr   string `json:"rpcAddr"`   // host:port for RPC
	NetworkID string `json:"networkId"` // network identifier
	PubKey    string `json:"pubkey"`    // hex public key
	Nonce     uint64 `json:"nonce"`     // anti-replay
	Signature string `json:"signature"` // signature of registration

	// Updated by heartbeats
	TipHeight uint64    `json:"tipHeight,omitempty"`
	PeerCount int       `json:"peerCount,omitempty"`
	LastSeen  time.Time `json:"lastSeen"`
}

type Registry struct {
	mu        sync.RWMutex
	nodes     map[string]*NodeRegistration // key = p2pAddr
	networkID string
}

func main() {
	port := flag.String("port", ":8088", "HTTP listen port")
	networkID := flag.String("network-id", "archivas-devnet-v3", "Accepted network ID")
	flag.Parse()

	registry := &Registry{
		nodes:     make(map[string]*NodeRegistration),
		networkID: *networkID,
	}

	// Cleanup stale nodes every minute
	go registry.cleanupLoop()

	http.HandleFunc("/register", registry.handleRegister)
	http.HandleFunc("/heartbeat", registry.handleHeartbeat)
	http.HandleFunc("/peers", registry.handlePeers)
	http.HandleFunc("/nodes", registry.handleNodes)
	http.HandleFunc("/health", registry.handleHealth)
	http.HandleFunc("/", registry.handleUI)

	log.Printf("📋 Archivas Node Registry")
	log.Printf("   Network: %s", *networkID)
	log.Printf("   Port: %s", *port)
	log.Println()
	log.Printf("Starting registry server...")

	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatal(err)
	}
}

func (r *Registry) handleRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reg NodeRegistration
	if err := json.NewDecoder(req.Body).Decode(&reg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate network ID
	if reg.NetworkID != r.networkID {
		http.Error(w, fmt.Sprintf("Network mismatch: got %s, want %s", reg.NetworkID, r.networkID), http.StatusBadRequest)
		return
	}

	// Verify signature
	if !r.verifyRegistration(&reg) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Store registration
	r.mu.Lock()
	reg.LastSeen = time.Now()
	r.nodes[reg.P2PAddr] = &reg
	r.mu.Unlock()

	log.Printf("✅ Registered node: %s (p2p=%s, rpc=%s)", reg.Address, reg.P2PAddr, reg.RPCAddr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

func (r *Registry) handleHeartbeat(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var hb struct {
		P2PAddr   string `json:"p2pAddr"`
		TipHeight uint64 `json:"tipHeight"`
		PeerCount int    `json:"peerCount"`
		Signature string `json:"signature"`
	}

	if err := json.NewDecoder(req.Body).Decode(&hb); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update existing node
	r.mu.Lock()
	node, exists := r.nodes[hb.P2PAddr]
	if exists {
		node.TipHeight = hb.TipHeight
		node.PeerCount = hb.PeerCount
		node.LastSeen = time.Now()
	}
	r.mu.Unlock()

	if !exists {
		http.Error(w, "Node not registered", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (r *Registry) handlePeers(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peers := make([]string, 0, len(r.nodes))
	now := time.Now()

	for p2pAddr, node := range r.nodes {
		if now.Sub(node.LastSeen) < HeartbeatTTL {
			peers = append(peers, p2pAddr)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"peers": peers,
		"count": len(peers),
	})
}

func (r *Registry) handleNodes(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	nodes := make([]*NodeRegistration, 0, len(r.nodes))
	now := time.Now()

	for _, node := range r.nodes {
		if now.Sub(node.LastSeen) < HeartbeatTTL {
			nodes = append(nodes, node)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"nodes": nodes,
		"count": len(nodes),
	})
}

func (r *Registry) handleHealth(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	activeNodes := 0
	now := time.Now()
	for _, node := range r.nodes {
		if now.Sub(node.LastSeen) < HeartbeatTTL {
			activeNodes++
		}
	}
	r.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":          true,
		"activeNodes": activeNodes,
		"totalNodes":  len(r.nodes),
	})
}

func (r *Registry) verifyRegistration(reg *NodeRegistration) bool {
	// Parse public key
	pubKeyBytes, err := hex.DecodeString(reg.PubKey)
	if err != nil || len(pubKeyBytes) != 33 {
		return false
	}

	pubKey, err := secp256k1.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false
	}

	// Compute message hash: H(address|p2pAddr|rpcAddr|networkId|nonce)
	msg := fmt.Sprintf("%s|%s|%s|%s|%d", reg.Address, reg.P2PAddr, reg.RPCAddr, reg.NetworkID, reg.Nonce)
	hash := sha256.Sum256([]byte(msg))

	// Parse signature
	sigBytes, err := hex.DecodeString(reg.Signature)
	if err != nil {
		return false
	}

	sig, err := ecdsa.ParseDERSignature(sigBytes)
	if err != nil {
		return false
	}

	// Verify
	return sig.Verify(hash[:], pubKey)
}

func (r *Registry) cleanupLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		removed := 0

		for p2pAddr, node := range r.nodes {
			if now.Sub(node.LastSeen) > HeartbeatTTL {
				delete(r.nodes, p2pAddr)
				removed++
				log.Printf("🧹 Removed stale node: %s (last seen %v ago)", p2pAddr, now.Sub(node.LastSeen))
			}
		}

		r.mu.Unlock()

		if removed > 0 {
			log.Printf("Cleanup: removed %d stale nodes", removed)
		}
	}
}

func (r *Registry) handleUI(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	activeNodes := make([]*NodeRegistration, 0)
	now := time.Now()
	for _, node := range r.nodes {
		if now.Sub(node.LastSeen) < HeartbeatTTL {
			activeNodes = append(activeNodes, node)
		}
	}
	r.mu.RUnlock()

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Archivas Node Registry</title>
    <meta http-equiv="refresh" content="30">
    <style>
        body { font-family: 'Courier New', monospace; max-width: 1200px; margin: 0 auto; padding: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; text-align: center; }
        .stats { display: flex; justify-content: space-around; margin: 30px 0; }
        .stat { background: #2a2a2a; padding: 20px; border: 1px solid #00ff00; text-align: center; flex: 1; margin: 0 10px; }
        .stat-value { font-size: 32px; font-weight: bold; color: #00ff00; }
        .stat-label { font-size: 12px; color: #888; margin-top: 5px; }
        .node { background: #2a2a2a; border: 1px solid #00ff00; padding: 15px; margin: 10px 0; }
        .node-header { font-size: 18px; font-weight: bold; margin-bottom: 10px; }
        .node-detail { margin: 5px 0; font-size: 14px; }
        .label { color: #888; }
        .online { color: #00ff00; }
        .refresh { text-align: center; font-size: 12px; color: #666; margin-top: 30px; }
        a { color: #00ff00; }
    </style>
</head>
<body>
    <h1>🌾 Archivas Node Registry</h1>
    
    <div class="stats">
        <div class="stat">
            <div class="stat-value">` + fmt.Sprintf("%d", len(activeNodes)) + `</div>
            <div class="stat-label">ACTIVE NODES</div>
        </div>
        <div class="stat">
            <div class="stat-value">` + r.networkID + `</div>
            <div class="stat-label">NETWORK ID</div>
        </div>
        <div class="stat">
            <div class="stat-value">` + fmt.Sprintf("%d", len(r.nodes)) + `</div>
            <div class="stat-label">TOTAL REGISTERED</div>
        </div>
    </div>

    <h2>📡 Active Nodes</h2>
`

	if len(activeNodes) == 0 {
		html += `<p style="text-align: center; color: #666;">No active nodes (waiting for heartbeats...)</p>`
	} else {
		for _, node := range activeNodes {
			timeSince := now.Sub(node.LastSeen).Round(time.Second)
			html += `
    <div class="node">
        <div class="node-header">` + node.Address + ` <span class="online">●</span></div>
        <div class="node-detail"><span class="label">P2P:</span> ` + node.P2PAddr + `</div>
        <div class="node-detail"><span class="label">RPC:</span> ` + node.RPCAddr + `</div>
`
			if node.TipHeight > 0 {
				html += `        <div class="node-detail"><span class="label">Height:</span> ` + fmt.Sprintf("%d", node.TipHeight) + `</div>`
			}
			if node.PeerCount > 0 {
				html += `        <div class="node-detail"><span class="label">Peers:</span> ` + fmt.Sprintf("%d", node.PeerCount) + `</div>`
			}
			html += `        <div class="node-detail"><span class="label">Last Seen:</span> ` + fmt.Sprintf("%v ago", timeSince) + `</div>
    </div>
`
		}
	}

	html += `
    <div class="refresh">
        <p>Auto-refresh every 30 seconds | <a href="/nodes">JSON API</a> | <a href="/peers">Peer List</a></p>
        <p style="margin-top: 20px; font-size: 11px;">
            <a href="https://github.com/ArchivasNetwork/archivas">GitHub</a> | 
            <a href="http://57.129.148.132:8082">Block Explorer</a> |
            <a href="http://57.129.148.132:3000">Grafana</a>
        </p>
    </div>
</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}
