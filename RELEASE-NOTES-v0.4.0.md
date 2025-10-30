# Archivas v0.4.0 - Long-Run Testnet Release

**Release Date:** October 30, 2025  
**Network:** archivas-devnet-v3  
**Status:** Community-Ready Testnet  

---

## 🎯 Release Highlights

### Community Onboarding
- ✅ **Built-in Faucet** - Get free test RCHV instantly (`GET /faucet?address=X`)
- ✅ **Registry Web UI** - Live dashboard showing all active nodes
- ✅ **One-Command Setup** - JOIN-TESTNET.md guide for easy onboarding
- ✅ **Automatic Discovery** - Peer gossip makes joining seamless

### Developer Experience
- ✅ **CORS Support** - External apps can query the blockchain
- ✅ **New Endpoints** - `/recentBlocks`, `/block/<height>`, `/health`
- ✅ **Detailed Health Stats** - Uptime, avg block time, blocks/hour
- ✅ **Complete API Docs** - Full endpoint reference

### Consensus & Stability
- ✅ **Difficulty Smoothing** - 10-block moving average for stable block times
- ✅ **Health Tracking** - Comprehensive chain health metrics
- ✅ **Enhanced Logging** - Better visibility for operators

---

## 📊 Testnet Status (At Release)

**Blockchain:**
- Height: 1878+ blocks
- RCHV Farmed: ~37,560
- Block Time: ~20 seconds (smoothed)
- Uptime: 24/7 capable

**Network:**
- Nodes: 2+ (ready for community growth)
- Peer Discovery: Automatic (gossip)
- Synchronization: Real-time (< 2s)

**Observability:**
- Grafana: http://57.129.148.132:3000
- Prometheus: http://57.129.148.132:9091
- Registry: http://57.129.148.132:8088
- Explorer: http://57.129.148.132:8082

---

## 🆕 New Features

### 1. Built-in Faucet

**Endpoint:** `GET /faucet?address=<your_address>`

**Example:**
```bash
curl "http://57.129.148.132:8080/faucet?address=arcv1..."
```

**Response:**
```json
{
  "success": true,
  "message": "20 RCHV sent! Should arrive in next block (~20 seconds)",
  "amount": "20.00000000 RCHV"
}
```

**Features:**
- 20 RCHV per drip
- Rate limit: 1 drip/hour per IP
- Automatic tx creation and signing
- Goes into next block

### 2. Registry Web UI

**URL:** http://57.129.148.132:8088

**Shows:**
- Active nodes count
- Network ID
- Live peer list
- Node details (P2P, RPC, height, peers)
- Auto-refresh every 30s

### 3. Developer API Endpoints

**GET /recentBlocks?count=N**
- Returns last N blocks (max 100)
- Includes: height, hash, timestamp, difficulty, farmer, txCount

**GET /block/<height>**
- Returns full block details
- Includes all transactions
- Challenge and proof data

**GET /health**
- Detailed chain health
- Uptime, avg block time, blocks/hour
- Last block timestamp

**CORS Enabled:**
- All endpoints support cross-origin requests
- External dashboards/explorers can integrate

### 4. Difficulty Smoothing

**Algorithm:**
- 10-block moving average
- Limits: 2x max increase, 0.5x min decrease
- Prevents difficulty spikes
- Smoother ~20s block times

---

## 🔧 API Changes

### New Endpoints
- `GET /faucet?address=<addr>` - Request test RCHV
- `GET /recentBlocks?count=<n>` - Last N blocks
- `GET /block/<height>` - Block by height
- `GET /health` - Detailed health stats
- `GET /` (registry) - Web UI dashboard

### Enhanced Endpoints
- `/healthz` - Now includes more metrics
- `/peers` - Shows connected + known peers
- `/metrics` - More Prometheus metrics

---

## 📚 New Documentation

- **docs/JOIN-TESTNET.md** - Complete joining guide
- **docs/FAUCET-API.md** - Faucet usage and API
- **docs/REGISTRY.md** - Registry documentation (updated)
- **docs/OBSERVABILITY.md** - Monitoring guide (updated)

---

## 🚀 Upgrading from v0.3.0

### Quick Upgrade

```bash
cd ~/archivas
git pull origin main
git checkout v0.4.0-alpha

# Rebuild
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-registry ./cmd/archivas-registry

# Restart (no data loss!)
pkill -f archivas-node archivas-registry
sleep 3

# Start node (same command as before)
./archivas-node --rpc :8080 --p2p :9090 --db ./data ...

# Start registry with new UI
./archivas-registry --port :8088 --network-id archivas-devnet-v3
```

**No database migration needed!** Just rebuild and restart.

---

## 🌍 Join the Testnet

**For New Users:**

See **docs/JOIN-TESTNET.md** for complete guide.

**Quick Start:**
```bash
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
go build -o archivas-node ./cmd/archivas-node

./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 57.129.148.132:9090,72.251.11.191:9090 \
  --enable-gossip
```

**Get Test RCHV:**
```bash
curl "http://57.129.148.132:8080/faucet?address=YOUR_ADDRESS"
```

---

## 🐛 Bug Fixes

- Fixed: Peer gossip network ID validation
- Fixed: Prometheus port conflict (9091 instead of 9090)
- Fixed: Chain fork resolution during sync
- Fixed: Mempool transaction handling

---

## ⚡ Performance Improvements

- Difficulty smoothing reduces variance
- Better peer connection management
- Optimized block verification
- Improved logging efficiency

---

## 🔒 Security

- Network ID validation in gossip
- Faucet rate limiting (anti-spam)
- CORS properly configured
- Registry signature verification

---

## 📖 Resources

**Links:**
- Repository: https://github.com/ArchivasNetwork/archivas
- Explorer: http://57.129.148.132:8082
- Registry: http://57.129.148.132:8088
- Grafana: http://57.129.148.132:3000
- Faucet: http://57.129.148.132:8080/faucet?address=YOUR_ADDR

**Documentation:**
- docs/JOIN-TESTNET.md
- docs/OBSERVABILITY.md
- docs/REGISTRY.md
- OPERATIONS.md

**Support:**
- GitHub Issues: https://github.com/ArchivasNetwork/archivas/issues
- Discussions: https://github.com/ArchivasNetwork/archivas/discussions

---

## 🎊 What's Next

**v0.5.0 Roadmap:**
- Wesolowski VDF (production-grade)
- Smart contracts (WASM)
- Light clients
- Network analytics
- Mainnet preparation

**Community Goals:**
- Reach 10,000 blocks
- Onboard 10+ community nodes
- Sustain 24/7 uptime
- Build ecosystem tools

---

## 🏆 Acknowledgments

Built from scratch in one epic 20-hour development session.

**Thank you to:**
- Early testers and community members
- Everyone running nodes
- Contributors and supporters

---

**Archivas v0.4.0 - The Community Testnet is Here!** 🌾

**Start farming, start building, start participating!** 🚀

---

**Released:** October 30, 2025  
**Tag:** v0.4.0-alpha  
**Network:** archivas-devnet-v3  
**Status:** 🟢 LIVE  

