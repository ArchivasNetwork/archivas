# Archivas Current Status

**Last Updated:** November 3, 2025  
**Version:** v1.2.0  
**Network:** archivas-devnet-v4  

---

## 🟢 Production Testnet Live

**Public Endpoint:** `https://seed.archivas.ai`  
**Explorer:** https://archivas-explorer-production.up.railway.app  
**Grafana:** http://57.129.148.132:3001  

---

## 📊 Network Statistics

| Metric | Value |
|--------|-------|
| **Block Height** | 64,000+ |
| **Total Blocks Mined** | ~64,000 |
| **Active Farmers** | 2 servers (7 k28 plots) |
| **Total Plot Space** | ~55 GB (1.6B hashes) |
| **Block Time** | ~20-30 seconds |
| **Difficulty** | 1,000,000 (normalized QMAX) |
| **Block Reward** | 20 RCHV |
| **Total Supply** | ~1,280,000 RCHV |
| **Uptime** | 25+ days |

---

## 🏗️ Infrastructure

### Seed Node (Server A - 57.129.148.132)

**Services:**
- Archivas Node (v1.2.0)
- Farmer (1 k28 plot)
- Nginx reverse proxy (TLS termination)
- Prometheus (metrics, port 9091)
- Grafana (dashboard, port 3001)

**Endpoints:**
- Public RPC: `https://seed.archivas.ai`
- Node RPC: `127.0.0.1:8080` (localhost only)
- P2P: `:9090` (public)
- Prometheus: `http://localhost:9091` (localhost only)
- Grafana: `http://57.129.148.132:3001`

**Security:**
- TLS 1.2+ with Let's Encrypt (auto-renewing)
- `/metrics` blocked from public access (404)
- Rate limiting on `/submit` (10 req/min)
- CORS configured for explorer

### Farmer Node (Server C - 57.129.148.134)

**Services:**
- Farmer (6 k28 plots)

**Configuration:**
- Connects to: `https://seed.archivas.ai`
- Plots: 6 × k28 (~330 GB)
- Uptime: 500+ hours

---

## 🌾 Farming Statistics

### Server A Farmer
- **Address:** `arcv1t3huuyd08er3yfnmk9c935rmx3wdh5j6m2uc9d`
- **Balance:** ~1,062,000 RCHV
- **Plots:** 1 × k28 (268M hashes)
- **Blocks Won:** ~53,100 (~83% of network)
- **Nonce:** 8 (8 outgoing transactions)

### Server C Farmer
- **Address:** `arcv1xjgsguj9e4assk4a24pfkm6cl92jwrgfxsru7c`
- **Balance:** ~197,000 RCHV
- **Plots:** 6 × k28 (1.6B hashes)
- **Blocks Won:** ~9,850 (~15% of network)
- **Nonce:** 1 (1 outgoing transaction)

**Note:** Server A has localhost advantage, resulting in higher win rate than expected.

---

## 💰 Wallet & Transactions

### Active Wallets

| Address | Balance | Purpose |
|---------|---------|---------|
| `arcv1t3huuyd08...` | 1,062,000 RCHV | Server A farmer |
| `arcv1xjgsguj9...` | 197,000 RCHV | Server C farmer |
| `arcv1fq46fw68...` | 925 RCHV | Test wallet |
| `arcv17ten7wh5...` | 160 RCHV | Test wallet |
| `arcv1r8hjkh5p...` | 100 RCHV | External transfer |
| `arcv1n4saqzuj...` | 0 RCHV | New wallet 4 |
| `arcv173h2uxmj...` | 0 RCHV | New wallet 5 |

### Transaction History
- **Total Transfers:** 1,160+ RCHV successfully transferred
- **Transaction Types:** Coinbase (block rewards) + Transfers
- **Nonce System:** Working correctly (sequential)
- **Fee Market:** Fixed fee (0.001 RCHV)

---

## 🔧 Technical Stack

### Core
- **Language:** Go 1.24.9
- **Consensus:** Proof-of-Space-and-Time
- **Storage:** BadgerDB (60,000+ blocks, ~2 GB)
- **Networking:** Custom P2P protocol
- **Cryptography:** Ed25519 (wallets), Blake2b (hashing)

### Wallet & Keys
- **Key Derivation:** BIP39 (24-word mnemonics) + SLIP-0010 (Ed25519)
- **Address Format:** Bech32 (`arcv` prefix, 20-byte Blake2b hash)
- **Transaction Format:** RFC 8785 canonical JSON + Ed25519 signatures
- **Domain Separation:** `Archivas-TxV1`

### APIs
- **RPC Version:** v1.2.0
- **Protocol:** HTTP/2, TLS 1.2+
- **Format:** JSON (all numeric fields as strings)
- **CORS:** Enabled for public access
- **Rate Limiting:** 10 req/min on `/submit`

---

## 🛠️ Developer Tools

### TypeScript SDK
- **Repository:** https://github.com/ArchivasNetwork/archivas-sdk
- **Status:** Complete, ready for npm publish
- **Features:** BIP39, SLIP-0010, Ed25519, Bech32, Transaction signing

### Block Explorer
- **Repository:** https://github.com/ArchivasNetwork/archivas-explorer
- **Live:** https://archivas-explorer-production.up.railway.app
- **Stack:** Next.js 16 + TypeScript + Tailwind CSS
- **Features:** Real-time updates, account lookup, block browsing

### CLI Tools
- `archivas-node` - Full node with RPC and P2P
- `archivas-farmer` - Farming client
- `archivas-cli` - Wallet and transaction tools (Ed25519)
- `archivas-wallet` - Legacy wallet (secp256k1)

---

## 🔐 Security

### Infrastructure
- ✅ TLS with auto-renewing Let's Encrypt certificate
- ✅ Rate limiting (10 req/min on `/submit`)
- ✅ Internal metrics not exposed publicly
- ✅ Node RPC bound to localhost only (not 0.0.0.0)
- ✅ Firewall rules (deny 8080, allow 80/443)

### Cryptography
- ✅ Ed25519 signatures (64 bytes)
- ✅ Blake2b-256 hashing
- ✅ Bech32 address encoding
- ✅ RFC 8785 canonical JSON (deterministic)
- ✅ Domain separation for transaction hashing

---

## 📈 Monitoring

### Prometheus Metrics
- `archivas_tip_height` - Current blockchain height
- `archivas_difficulty` - Mining difficulty target
- `archivas_peer_count` - Connected P2P peers
- `archivas_blocks_total` - Total blocks processed
- `archivas_submit_accepted_total` - Valid proof submissions
- `archivas_rpc_requests_total` - RPC endpoint usage

### Grafana Dashboards
- **URL:** http://57.129.148.132:3001
- **Dashboard:** Archivas Network Overview
- **Panels:** Tip Height, Peer Count, Difficulty, Block Rate, RPC Requests

---

## 🚀 Recent Milestones

### v1.2.0 (Current)
- ✅ Explorer listing endpoints (`/blocks/recent`, `/tx/recent`)
- ✅ Transaction type field (`coinbase` vs `transfer`)
- ✅ CORS duplication fixed
- ✅ Farmer field terminology

### v1.1.1-infra
- ✅ Seed node infrastructure (seed.archivas.ai)
- ✅ Nginx reverse proxy with TLS
- ✅ Let's Encrypt certificate (auto-renewing)
- ✅ HTTP/2 enabled
- ✅ Security hardening

### v1.1.0
- ✅ Wallet primitives frozen
- ✅ Ed25519 keypairs with BIP39/SLIP-0010
- ✅ Transaction v1 schema
- ✅ Public wallet API

---

## 🎯 Roadmap

### Phase 1: Testnet Stability (COMPLETE)
- ✅ Multi-node P2P networking
- ✅ Block synchronization
- ✅ Public RPC endpoint
- ✅ Basic monitoring

### Phase 2: Developer Tools (COMPLETE)
- ✅ TypeScript SDK
- ✅ Block explorer
- ✅ Wallet CLI tools
- ✅ API documentation

### Phase 3: Public Launch (IN PROGRESS)
- ✅ Public seed node (seed.archivas.ai)
- ✅ Explorer deployment
- ⏳ SDK published to npm
- ⏳ Faucet for testnet RCHV
- ⏳ Community documentation

### Phase 4: Production Hardening (PLANNED)
- Security audit
- Performance optimization
- State pruning
- Snapshot sync
- Checkpoint system

### Phase 5: Mainnet (2026)
- Economic model finalization
- Token distribution plan
- Mainnet deployment
- Exchange listings

---

## 🔗 Links

- **Core Repository:** https://github.com/ArchivasNetwork/archivas
- **TypeScript SDK:** https://github.com/ArchivasNetwork/archivas-sdk
- **Block Explorer:** https://github.com/ArchivasNetwork/archivas-explorer
- **Public RPC:** https://seed.archivas.ai
- **Live Explorer:** https://archivas-explorer-production.up.railway.app

---

## 📞 Community

- **GitHub Issues:** https://github.com/ArchivasNetwork/archivas/issues
- **Documentation:** https://github.com/ArchivasNetwork/archivas/tree/main/docs

---

**Archivas is a production-ready Proof-of-Space-and-Time blockchain with a live, operational testnet.** 🌾

