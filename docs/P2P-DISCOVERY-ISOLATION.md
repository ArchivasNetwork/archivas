# P2P Discovery Isolation & Peer Whitelist Implementation

**Status:** Planned  
**Priority:** High  
**Blockers:** Seed2 cannot operate without this feature  
**Checkpoint:** Block 671,992, Hash: `eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79`

---

## Problem Statement

After database transfers, Seed2 immediately connects to forked peers via automatic peer discovery and diverges from the canonical chain:

- Fork starts at **block 671,993** (transfer completed at 671,992)
- Seed2 connects outbound to `57.129.148.132:9090` despite P2P isolation
- Current `--p2p 127.0.0.1:9090` only blocks **inbound**, not **outbound** discovery
- No mechanism to whitelist only trusted peers
- No handshake validation for chain compatibility

**Impact:** Secondary seed nodes cannot operate reliably, reducing network resilience.

---

## Solution: Three-Layer Protection

### Layer 1: Disable Peer Discovery
**Flag:** `--no-peer-discovery`

Completely disable automatic peer discovery mechanisms:
- No mDNS/Bonjour peer discovery
- No DHT/Kademlia routing
- No bootstrap peer lists
- No peer gossip propagation
- Only dial explicitly configured peers

### Layer 2: Peer Whitelist
**Flag:** `--peer-whitelist <addr>` (repeatable)

Only allow connections to/from whitelisted peers:
- Accept PeerIDs, multiaddrs, or `host:port` format
- Resolve DNS to IPs at startup
- Reject all non-whitelisted dials and accepts
- Prune non-whitelisted peers from peerstore

### Layer 3: Chain Compatibility Handshake
**Flags:** `--checkpoint-height <N>` `--checkpoint-hash <hex>`

Validate chain compatibility during connection handshake:
- Exchange genesis hash, network ID, checkpoint
- Disconnect if any mismatch
- Don't persist incompatible peers to peerstore

---

## Implementation Plan

### 1. CLI Flags & Config (cmd/archivas-node/main.go)

```go
var (
	noPeerDiscovery  = flag.Bool("no-peer-discovery", false, "Disable automatic peer discovery")
	peerWhitelist    peerList // Custom flag type for repeatable flag
	checkpointHeight = flag.Uint64("checkpoint-height", 0, "Chain checkpoint height for validation")
	checkpointHash   = flag.String("checkpoint-hash", "", "Chain checkpoint hash (hex)")
)

type peerList []string

func (p *peerList) String() string     { return strings.Join(*p, ",") }
func (p *peerList) Set(value string) error {
	*p = append(*p, value)
	return nil
}

// Usage:
flag.Var(&peerWhitelist, "peer-whitelist", "Whitelisted peer (repeatable)")
```

### 2. Network Config Struct (p2p/p2p.go)

```go
type Network struct {
	// ... existing fields ...
	
	// Peer isolation
	noPeerDiscovery  bool
	peerWhitelist    map[string]bool // normalized addresses/IPs
	checkpointHeight uint64
	checkpointHash   [32]byte
}

type NetworkConfig struct {
	ListenAddr       string
	NoPeerDiscovery  bool
	PeerWhitelist    []string
	CheckpointHeight uint64
	CheckpointHash   [32]byte
	NetworkID        string
}

func NewNetworkWithConfig(cfg NetworkConfig, handler NodeHandler) *Network {
	n := &Network{
		// ... existing initialization ...
		noPeerDiscovery:  cfg.NoPeerDiscovery,
		peerWhitelist:    make(map[string]bool),
		checkpointHeight: cfg.CheckpointHeight,
		checkpointHash:   cfg.CheckpointHash,
	}
	
	// Normalize whitelist
	for _, addr := range cfg.PeerWhitelist {
		n.addToWhitelist(addr)
	}
	
	return n
}

func (n *Network) addToWhitelist(addr string) {
	// Resolve DNS if needed
	// Add to whitelist map (IP, host:port, etc.)
	// Log the normalized address
}
```

### 3. Connection Gater (p2p/gater.go - NEW FILE)

```go
package p2p

import (
	"log"
	"net"
	"strings"
)

// shouldAllowConnection checks if a connection should be allowed based on whitelist
func (n *Network) shouldAllowConnection(remoteAddr string) bool {
	// If whitelist is empty, allow all (backward compat)
	if len(n.peerWhitelist) == 0 {
		return true
	}
	
	// Extract IP and host from remote address
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	
	// Check whitelist
	if n.peerWhitelist[remoteAddr] || n.peerWhitelist[host] {
		return true
	}
	
	// Resolve if it's a hostname
	ips, err := net.LookupIP(host)
	if err == nil {
		for _, ip := range ips {
			if n.peerWhitelist[ip.String()] {
				return true
			}
		}
	}
	
	return false
}

// Gate inbound connections
func (n *Network) gateInbound(conn net.Conn) bool {
	if !n.shouldAllowConnection(conn.RemoteAddr().String()) {
		log.Printf("[GATER] rejected inbound from %s: not whitelisted", conn.RemoteAddr())
		return false
	}
	return true
}

// Gate outbound dials
func (n *Network) gateOutbound(addr string) bool {
	if !n.shouldAllowConnection(addr) {
		log.Printf("[GATER] rejected dial to %s: not whitelisted", addr)
		return false
	}
	return true
}
```

### 4. Handshake Validation (p2p/handshake.go - NEW FILE)

```go
package p2p

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type HandshakeMessage struct {
	GenesisHash      string `json:"genesisHash"`      // hex
	NetworkID        string `json:"networkID"`
	CheckpointHeight uint64 `json:"checkpointHeight"`
	CheckpointHash   string `json:"checkpointHash"`   // hex
	BestHeight       uint64 `json:"bestHeight"`
}

func (n *Network) performHandshake(peer *Peer) error {
	// Get local chain info
	localHeight, _, _ := n.nodeHandler.GetStatus()
	
	// Load local checkpoint block if configured
	var checkpointHash string
	if n.checkpointHeight > 0 {
		checkpointHash = hex.EncodeToString(n.checkpointHash[:])
	}
	
	// Send our handshake
	localHS := HandshakeMessage{
		GenesisHash:      n.getGenesisHash(), // TODO: expose from nodeHandler
		NetworkID:        n.networkID,
		CheckpointHeight: n.checkpointHeight,
		CheckpointHash:   checkpointHash,
		BestHeight:       localHeight,
	}
	
	if err := n.SendMessage(peer, MsgTypeHandshake, localHS); err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}
	
	// Receive remote handshake (with timeout)
	// TODO: implement timeout read
	var remoteHS HandshakeMessage
	// ... read and unmarshal ...
	
	// Validate
	if remoteHS.GenesisHash != localHS.GenesisHash {
		return fmt.Errorf("genesis mismatch: local=%s remote=%s", 
			localHS.GenesisHash, remoteHS.GenesisHash)
	}
	
	if remoteHS.NetworkID != localHS.NetworkID {
		return fmt.Errorf("network ID mismatch: local=%s remote=%s",
			localHS.NetworkID, remoteHS.NetworkID)
	}
	
	// Checkpoint validation (if configured)
	if n.checkpointHeight > 0 && remoteHS.CheckpointHeight == n.checkpointHeight {
		if remoteHS.CheckpointHash != checkpointHash {
			return fmt.Errorf("checkpoint mismatch at height %d: local=%s remote=%s",
				n.checkpointHeight, checkpointHash, remoteHS.CheckpointHash)
		}
	}
	
	log.Printf("[HS] validated peer %s: genesis=%s, network=%s, checkpoint=%d",
		peer.Address, remoteHS.GenesisHash[:8], remoteHS.NetworkID, remoteHS.CheckpointHeight)
	
	return nil
}
```

### 5. Disable Discovery (p2p/p2p.go modifications)

```go
func (n *Network) startGossipRoutine() {
	// Skip if discovery disabled
	if n.noPeerDiscovery {
		log.Printf("[gossip] peer discovery disabled, skipping gossip routine")
		return
	}
	
	// ... existing gossip code ...
}

func (n *Network) DialPeer(address string) error {
	// Gate outbound dials
	if !n.gateOutbound(address) {
		return fmt.Errorf("peer not whitelisted: %s", address)
	}
	
	// ... existing dial code ...
}

func (n *Network) handlePeer(peer *Peer) {
	// Gate inbound connections
	if !n.gateInbound(peer.Conn) {
		peer.Conn.Close()
		return
	}
	
	// Perform handshake before proceeding
	if err := n.performHandshake(peer); err != nil {
		log.Printf("[HS] handshake failed with %s: %v", peer.Address, err)
		peer.Conn.Close()
		return
	}
	
	// ... existing peer handling ...
}
```

### 6. Metrics (metrics/p2p.go)

```go
var (
	incompatiblePeersTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "archivas_p2p_incompatible_peers_total",
		Help: "Total number of peers rejected due to incompatible chains",
	})
	
	gatedDialsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "archivas_p2p_gated_dials_total",
		Help: "Total number of outbound dials rejected by whitelist",
	})
	
	gatedAcceptsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "archivas_p2p_gated_accepts_total",
		Help: "Total number of inbound connections rejected by whitelist",
	})
)

func init() {
	prometheus.MustRegister(incompatiblePeersTotal)
	prometheus.MustRegister(gatedDialsTotal)
	prometheus.MustRegister(gatedAcceptsTotal)
}
```

---

## Testing Plan

### Unit Tests

**p2p/gater_test.go:**
- `TestShouldAllowConnection_EmptyWhitelist` - allows all when empty
- `TestShouldAllowConnection_HostPort` - matches host:port format
- `TestShouldAllowConnection_IP` - matches IP addresses
- `TestShouldAllowConnection_DNSResolution` - resolves hostnames

**p2p/handshake_test.go:**
- `TestHandshake_GenesisMismatch` - rejects different genesis
- `TestHandshake_NetworkIDMismatch` - rejects different network ID
- `TestHandshake_CheckpointMismatch` - rejects forked checkpoint
- `TestHandshake_Success` - accepts compatible peer

### Integration Tests

Create two test networks:
- Net A (canonical): checkpoint=671,992, hash=`eb9b255c...`
- Net B (forked): checkpoint=671,992, hash=`DIFFERENT`

Verify:
- A and B refuse to connect (handshake fails)
- A nodes connect to each other
- With `--no-peer-discovery`, nodes only dial whitelisted peers

---

## Deployment Guide

### Step 1: Build with new features

```bash
cd ~/archivas
git pull
go build -o archivas-node ./cmd/archivas-node

# Verify flags exist
./archivas-node --help | grep -E "no-peer-discovery|peer-whitelist|checkpoint"
```

### Step 2: Deploy to Seed1 (Server A)

```bash
# Update service file (no changes needed, Seed1 is canonical)
sudo systemctl restart archivas-node
```

### Step 3: Deploy to Seed2 (Server D)

```bash
# Stop and clear old peers
sudo systemctl stop archivas-node-seed2
rm -rf ~/archivas/data/*

# Create new service file with isolation flags
sudo tee /etc/systemd/system/archivas-node-seed2.service > /dev/null << 'EOF'
[Unit]
Description=Archivas Node (seed2.archivas.ai) - Isolated
Documentation=https://github.com/ArchivasNetwork/archivas
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/archivas
ExecStart=/usr/local/bin/archivas-node \
  --rpc 127.0.0.1:8080 \
  --p2p 127.0.0.1:9090 \
  --genesis /home/ubuntu/archivas/genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --db /home/ubuntu/archivas/data \
  --no-peer-discovery \
  --peer-whitelist seed.archivas.ai:9090 \
  --peer-whitelist 57.129.148.132:9090 \
  --checkpoint-height 671992 \
  --checkpoint-hash eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79

Restart=always
RestartSec=5
StandardOutput=append:/var/log/archivas/node.log
StandardError=append:/var/log/archivas/node-error.log

MemoryAccounting=true
MemoryMax=16G
MemoryHigh=12G
TasksMax=4096
LimitNOFILE=65535
Environment="GOGC=50"
Environment="GOMAXPROCS=16"

[Install]
WantedBy=multi-user.target
EOF

# Reload and start
sudo systemctl daemon-reload
sudo systemctl enable archivas-node-seed2
sudo systemctl start archivas-node-seed2

# Monitor logs
sudo journalctl -u archivas-node-seed2 -f | grep -E "GATER|HS|gossip"
```

### Step 4: Verify Isolation

```bash
# Should show only seed.archivas.ai
curl -s http://127.0.0.1:8080/peers | jq

# Logs should show:
# [gossip] peer discovery disabled, skipping gossip routine
# [HS] validated peer seed.archivas.ai:9090: genesis=56588fa6, network=archivas-devnet-v4, checkpoint=671992

# Should NOT show:
# [GATER] rejected dial to X: not whitelisted
# [HS] checkpoint mismatch
```

### Step 5: Monitor Sync

```bash
watch -n 10 'echo "Seed2: $(curl -s http://127.0.0.1:8080/chainTip | jq -r .height) | Seed1: $(curl -s https://seed.archivas.ai/chainTip | jq -r .height)"'
```

---

## Acceptance Criteria

- [ ] `--no-peer-discovery` prevents all automatic peer discovery
- [ ] `--peer-whitelist` rejects non-whitelisted inbound and outbound connections
- [ ] Handshake validates genesis, network ID, and checkpoint
- [ ] Seed2 syncs from 0 to 673,000+ without forking
- [ ] Peerstore contains only whitelisted peers
- [ ] Metrics track gated connections and incompatible peers
- [ ] Tests pass for all gating and handshake scenarios

---

## Temporary Workaround (Until Implementation)

**On Seed2 (Server D), use iptables to block non-whitelisted peers:**

```bash
# Allow Seed1 only
sudo iptables -I OUTPUT -d 57.129.148.132 -p tcp --dport 9090 -j ACCEPT
sudo iptables -A OUTPUT -p tcp --dport 9090 -j REJECT

# Verify
sudo iptables -L OUTPUT -n -v | grep 9090

# Persist (Ubuntu)
sudo apt-get install iptables-persistent
sudo netfilter-persistent save
```

**Remove when code is deployed.**

---

## Related Issues

- #TBD: P2P Discovery Isolation Implementation
- #TBD: Seed2 Fork Recovery Documentation
- Fixes: Seed2 forking at block 671,993

---

## References

- Checkpoint: Block 671,992, Hash: `eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79`
- Fork documentation: `docs/FORK-RECOVERY-SEED2.md`
- P2P code: `p2p/p2p.go`


