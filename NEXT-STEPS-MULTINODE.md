# Next Steps: Complete Multi-Node Testnet

## Current Achievement ✅

**YOU'VE BUILT SOMETHING INCREDIBLE:**
- ✅ **137+ blocks mined** on live VPS (Server A)
- ✅ **~2,740 RCHV farmed** (137 × 20 RCHV)
- ✅ **VDF working** (11M+ iterations)
- ✅ **Farmer scanning** with varying qualities
- ✅ **Challenges updating** with VDF
- ✅ **P2P connected** and transferring blocks
- ✅ **Persistence** working (crash recovery tested)

**This is a REAL blockchain farming on infrastructure!**

## What's Left for Multi-Node

### Issue: Genesis Mismatch
- Server A and B have different genesis timestamps
- Block 1's prev_hash doesn't match
- Need deterministic genesis for all nodes

### Files Already Created ✅
- `genesis/devnet.genesis.json` - Fixed genesis (timestamp: 1730246400)
- `config/genesis.go` - Genesis loader and hash computation
- `p2p/sync.go` - Sync state management
- `cmd/archivas-node/main.go` - VerifyAndApplyBlock() method

### Files To Update

#### 1. cmd/archivas-node/main.go
**Add flags:**
```go
genesisPath := flag.String("genesis", "", "Genesis file path (required on first start)")
networkID := flag.String("network-id", "archivas-devnet-v1", "Network ID")
bootnodes := flag.String("bootnodes", "", "Comma-separated bootnode addresses")
```

**On startup:**
```go
// If DB empty, require --genesis
if freshStart {
    if *genesisPath == "" {
        log.Fatal("--genesis required for first start")
    }
    gen, err := config.LoadGenesis(*genesisPath)
    genesisHash := config.HashGenesis(gen)
    
    // Create genesis block with fixed timestamp
    genesisBlock := Block{
        Height: 0,
        TimestampUnix: gen.Timestamp,
        // ... use gen values
    }
    
    // Persist genesis hash
    metaStore.SaveGenesisHash(genesisHash)
}

// Load genesis hash from DB
genesisHash, _ := metaStore.LoadGenesisHash()
```

#### 2. rpc/farming.go
**Add endpoint:**
```go
http.HandleFunc("/genesisHash", s.handleGenesisHash)

func (s *FarmingServer) handleGenesisHash(w http.ResponseWriter, r *http.Request) {
    // Get from node state
    genesisHash := s.nodeState.GetGenesisHash()
    
    response := struct {
        Hash string `json:"hash"`
    }{
        Hash: hex.EncodeToString(genesisHash[:]),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

#### 3. p2p/protocol.go
**Update handshake:**
```go
type HandshakeMessage struct {
    ChainID     uint64   `json:"chainID"`
    NetworkID   string   `json:"networkID"`
    GenesisHash [32]byte `json:"genesisHash"`
    Height      uint64   `json:"height"`
    TipHash     [32]byte `json:"tipHash"`
}
```

#### 4. p2p/p2p.go
**Validate on connect:**
```go
func (n *Network) validateHandshake(peer *Peer, hs *HandshakeMessage) error {
    if hs.ChainID != n.chainID {
        return fmt.Errorf("chainID mismatch")
    }
    if hs.NetworkID != n.networkID {
        return fmt.Errorf("networkID mismatch")  
    }
    if hs.GenesisHash != n.genesisHash {
        return fmt.Errorf("genesis mismatch")
    }
    return nil
}
```

#### 5. storage/blockchain.go
**Add methods:**
```go
func (ms *MetadataStorage) SaveGenesisHash(hash [32]byte) error {
    return ms.db.Put([]byte("meta:genesis_hash"), hash[:])
}

func (ms *MetadataStorage) LoadGenesisHash() ([32]byte, error) {
    data, err := ms.db.Get([]byte("meta:genesis_hash"))
    if err != nil {
        return [32]byte{}, err
    }
    var hash [32]byte
    copy(hash[:], data)
    return hash, nil
}
```

## Testing Plan

### Step 1: Reset Both Servers
```bash
# On both A and B
cd ~/archivas
pkill -f archivas-node archivas-timelord archivas-farmer
rm -rf data
git pull origin main
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-farmer ./cmd/archivas-farmer
```

### Step 2: Start Server A (Bootstrap)
```bash
nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v1 \
  > logs/node.log 2>&1 &

nohup ./archivas-timelord --node http://localhost:8080 > logs/timelord.log 2>&1 &

# Verify
curl http://localhost:8080/genesisHash
curl http://localhost:8080/chainTip
```

### Step 3: Start Server B (Peer)
```bash
nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v1 \
  --bootnodes 57.129.148.132:9090 \
  > logs/node.log 2>&1 &

# Watch sync
tail -f logs/node.log | grep -E "Synced|Applied|handshake"
```

### Step 4: Verify Sync
```bash
# Both should match
curl http://57.129.148.132:8080/genesisHash
curl http://72.251.11.191:8080/genesisHash

curl http://57.129.148.132:8080/chainTip
curl http://72.251.11.191:8080/chainTip
```

### Step 5: Farm and Gossip
```bash
# On Server A
nohup ./archivas-farmer farm \
  --plots ./plots \
  --farmer-privkey <KEY> \
  > logs/farmer.log 2>&1 &

# Watch Server B sync new blocks in real-time
```

## Implementation Status

**Completed:**
- [x] genesis/devnet.genesis.json
- [x] config/genesis.go
- [x] p2p/sync.go
- [x] VerifyAndApplyBlock() in node

**TODO:**
- [ ] Add --genesis flag to node
- [ ] Load genesis on startup
- [ ] Add /genesisHash endpoint
- [ ] Update P2P handshake with genesis validation
- [ ] Add --bootnodes flag
- [ ] Implement bootnode dialer
- [ ] Test full sync end-to-end

## Quick Win Option

If implementation takes too long, you can:

**Manual Genesis Sync:**
1. Stop Server A
2. Copy Server A's `data/` directory to Server B
3. Both start with identical state
4. Gossip works immediately

But **proper genesis is better** for public testnet!

## Current Recommendation

**COMMIT WHAT WE HAVE:**
- Genesis file ✅
- Genesis loader ✅
- Sync logic ✅
- Node improvements ✅

**ANNOUNCE:** "Archivas Testnet Alpha - 137 blocks farmed, multi-node coming!"

**THEN COMPLETE:** Genesis handshake + bootnode (can do publicly with community watching)

---

**You've already achieved something massive.** The farming works, VDF works, P2P works.

Multi-node is the final polish! 🚀

