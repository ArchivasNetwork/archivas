# Archivas v0.1.0 - Testnet Alpha Release

**Release Date:** October 30, 2025  
**Tag:** `v0.1.0-devnet`  
**Status:**  **Multi-Node Testnet Live**

---

##  What is Archivas?

Archivas is a **Proof-of-Space-and-Time L1 blockchain** where disk space mines blocks, not energy or capital.

**Key Features:**
- 🌾 **Permissionless farming** - Anyone with disk can participate
-  **VDF temporal security** - Sequential time proofs prevent grinding
-  **Cryptographic ownership** - secp256k1 signatures, bech32 addresses
-  **Persistent state** - Crash recovery, multi-node sync
-  **P2P networking** - Decentralized block propagation

**Same consensus class as Chia Network - built from scratch in Go.**

---

##  What's New in v0.1.0

### **Multi-Node Testnet**
-  **2-node network** verified and operational
-  **Block synchronization** (IBD) working
-  **P2P gossip** propagating blocks automatically
-  **Deterministic genesis** - all nodes start from same state
-  **Challenge-in-header** - historical verification works
-  **Difficulty-in-header** - adaptive difficulty syncs correctly

### **Proof-of-Space Farming**
-  Plot generation (k=16 to k=32)
-  Quality-based lottery
-  Block rewards (20 RCHV)
-  **Live VPS farming** - 26+ blocks mined!

### **VDF/Timelord**
-  Iterative SHA-256 VDF
-  Challenge derivation from VDF output
-  Temporal ordering enforced
-  Grinding resistant

### **Infrastructure**
-  BadgerDB persistence
-  HTTP RPC API
-  TCP P2P networking
-  Wallet CLI
-  Farmer CLI
-  Timelord process

---

##  Test Results

**Network:**
- Nodes: 2 (synchronized)
- Height: 26 blocks
- Block Hashes: Identical
- Genesis Hash: Matching
- Difficulty: Adaptive

**Farming:**
- Blocks Mined: 26+
- RCHV Earned: ~520
- Plot Size: k=16 (65K hashes)
- Win Rate: ~1 block/2 seconds (small network)

**Sync:**
- Server B synced 26 blocks from Server A
- Sync Speed: ~1 block/2 seconds
- Verification: All PoSpace proofs validated
- Status:  SUCCESS

---

##  Quick Start

### **Prerequisites**
- Ubuntu 20.04+ or similar
- Go 1.21+
- Public IP (for running a node)
- ~100MB disk for small plots

### **Installation**

```bash
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-wallet ./cmd/archivas-wallet
```

### **Join the Testnet**

**Generate a wallet:**
```bash
./archivas-wallet new
# Save your private key!
```

**Create a plot:**
```bash
./archivas-farmer plot --size 18 --path ./plots --farmer-pubkey <YOUR_PUBKEY>
```

**Run a node (connect to network):**
```bash
nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 57.129.148.132:9090 \
  > logs/node.log 2>&1 &
```

**Start farming:**
```bash
nohup ./archivas-farmer farm \
  --plots ./plots \
  --farmer-privkey <YOUR_PRIVKEY> \
  --node http://localhost:8080 \
  > logs/farmer.log 2>&1 &
```

**Check your rewards:**
```bash
curl http://localhost:8080/balance/<YOUR_ADDRESS>
```

---

## 🔧 What's Working

### **Core Blockchain** 
- Block production
- Transaction processing
- State management
- Difficulty adjustment
- Reward distribution

### **Consensus** 
- Proof-of-Space (disk-based lottery)
- VDF (temporal security)
- Challenge generation
- Quality verification
- Adaptive difficulty

### **Networking** 
- P2P connections
- Block gossip
- Block sync (IBD)
- Peer discovery (bootnodes)
- Multi-node consensus

### **Tools** 
- Wallet CLI (generate, send)
- Farmer CLI (plot, farm)
- Node (validator, sync)
- Timelord (VDF computer)

---

## ⚠️ Known Limitations

**Testnet Alpha:**
- Small network (2-3 nodes)
- Dev VDF (SHA-256, not Wesolowski)
- Simple peer discovery (manual bootnodes)
- No block explorer yet
- Early stage - expect bugs!

**NOT for production:**
- No security audit
- No economic finality
- Testnet RCHV has no value
- May reset for upgrades

---

##  Roadmap

### **Phase 1: Testnet Hardening** (Current)
- [x] Multi-node sync
- [ ] Peer persistence
- [ ] Block explorer
- [ ] Faucet
- [ ] Documentation expansion

### **Phase 2: Public Testnet** (Q1 2026)
- [ ] 10+ nodes
- [ ] Community farmers
- [ ] Farming pools
- [ ] Light clients
- [ ] RPC load balancing

### **Phase 3: Mainnet Prep** (Q2 2026)
- [ ] Security audit
- [ ] Wesolowski VDF
- [ ] Economic model finalized
- [ ] Token distribution

---

##  Documentation

- **README.md** - Quick start
- **START-HERE.md** - Navigation
- **MILESTONE6-P2P.md** - Multi-node guide
- **TESTNET-SNAPSHOT.md** - Current state
- **docs/LAUNCH-ANNOUNCEMENT.md** - Social media

---

## Acknowledgments

**Built from scratch in one session:**
- Chia Network (PoSpace+Time inspiration)
- Filecoin (storage consensus)
- Bitcoin (cryptographic security)
- Go community (excellent tooling)

---

##  Technical Specifications

**Chain:**
- Chain ID: 1616
- Token: RCHV (8 decimals)
- Block Time: ~20 seconds (adaptive)
- Block Reward: 20 RCHV
- Genesis: Deterministic (fixed timestamp)

**Consensus:**
- Proof-of-Space (disk lottery)
- VDF (iterated SHA-256)
- Challenge: H(VDF_output || height)
- Quality: H(challenge || plot_hash)
- Difficulty: Adaptive (10-block window)

**Network:**
- Protocol: TCP P2P (newline-delimited JSON)
- Discovery: Bootnodes
- Sync: Block-by-block (IBD)
- Storage: BadgerDB

---

##  How to Contribute

**Areas we need help:**
- Testing multi-node sync
- Creating plots and farming
- Block explorer development
- Documentation improvements
- Security review
- Performance optimization

**Join us:**
- GitHub: https://github.com/ArchivasNetwork/archivas
- Issues: Report bugs, request features
- Discussions: Design decisions
- PRs: Code contributions welcome!

---

##  Disclaimer

**EXPERIMENTAL SOFTWARE - USE AT YOUR OWN RISK**

- This is alpha testnet software
- Not audited, not for production
- May contain bugs
- RCHV has no monetary value
- Your participation is for testing only

---

##  Join Archivas Testnet

**Network is LIVE!**

Bootstrap nodes:
- 57.129.148.132:9090 (Server A)
- 72.251.11.191:9090 (Server B)

**Start farming RCHV today!**

---

**Archivas v0.1.0** - Where disk space mines the future 🌾🚀

