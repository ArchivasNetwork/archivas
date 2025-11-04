# Current Status

**Last Updated:** November 3, 2025  
**Version:** v1.2.0  
**Network:** archivas-devnet-v4  
**Status:** 🟢 LIVE

---

## Network Overview

Archivas is a **live, operational** Proof-of-Space-and-Time blockchain with:
- Public RPC endpoint with HTTPS
- Active farmers on 2 servers
- Block explorer with real-time updates
- TypeScript SDK for developers
- Full monitoring and observability

This is **not a demo** - it's a production-grade testnet that has been running continuously for 25+ days.

---

## Current Statistics

### Blockchain
- **Height:** 64,000+ blocks
- **Uptime:** 25+ days continuous
- **Total Supply:** ~1,280,000 RCHV
- **Block Time:** ~25 seconds average
- **Difficulty:** 1,000,000 (QMAX normalized)

### Network
- **Active Farmers:** 2
- **Total Plots:** 7 × k28 (~55 GB)
- **Plot Hashes:** 1.6 billion
- **Peers Connected:** 2-3

### Transactions
- **Total Transfers:** 1,160+ RCHV
- **Transaction Count:** 8 confirmed
- **Mempool:** Working correctly
- **Fee:** 0.001 RCHV standard

---

## Infrastructure

### Seed Node (seed.archivas.ai)
- **Location:** 57.129.148.132
- **Services:** Node, Farmer, Nginx, Prometheus, Grafana
- **Uptime:** 99.9%
- **TLS:** Let's Encrypt (auto-renewing)
- **Security:** Rate limiting, CORS, localhost binding

### Farmer Nodes
1. **Server A:** 1 × k28 plot, 1,062,000 RCHV earned
2. **Server C:** 6 × k28 plots, 197,000 RCHV earned

---

## Public Services

| Service | URL | Status |
|---------|-----|--------|
| Public RPC | https://seed.archivas.ai | 🟢 Live |
| Block Explorer | https://archivas-explorer-production.up.railway.app | 🟢 Live |
| Grafana Dashboard | http://57.129.148.132:3001 | 🟢 Live |
| TypeScript SDK | https://github.com/ArchivasNetwork/archivas-sdk | ✅ Ready |

---

## Recent Milestones

### ✅ Completed

**v1.2.0 (November 2025)**
- Explorer listing endpoints (`/blocks/recent`, `/tx/recent`)
- Transaction type field (coinbase vs transfer)
- CORS fixes for browser compatibility
- Farmer terminology (was "miner")

**v1.1.1-infra (November 2025)**
- Seed node at seed.archivas.ai
- Nginx reverse proxy with TLS
- Let's Encrypt certificate
- HTTP/2 enabled
- Security hardening

**v1.1.0 (November 2025)**
- Wallet primitives frozen
- Ed25519 with BIP39/SLIP-0010
- Transaction v1 schema
- Public wallet API

**Earlier (October 2025)**
- Initial block download (IBD)
- P2P networking
- Block synchronization
- Difficulty retargeting
- PoSpace verification fixes

---

## What's Working

### ✅ Core Functionality
- Block production and validation
- Proof-of-Space farming
- VDF challenge generation
- Difficulty adjustment
- State persistence (BadgerDB)
- Crash recovery

### ✅ Networking
- P2P block propagation
- Peer discovery and gossip
- Initial block download (IBD)
- Multi-server coordination
- Handshake with genesis validation

### ✅ Wallets & Transactions
- BIP39 mnemonic generation
- Ed25519 key derivation (SLIP-0010)
- Bech32 address encoding
- Transaction signing and verification
- Nonce sequencing
- Balance updates

### ✅ APIs & Tools
- Public RPC with 12+ endpoints
- TypeScript SDK (wallet, RPC client)
- Block explorer (Next.js)
- CLI tools (keygen, sign, broadcast)
- Prometheus metrics
- Grafana dashboards

### ✅ Security & Operations
- TLS with auto-renewal
- CORS for public access
- Rate limiting
- Localhost-only metrics
- Firewall configuration
- Health checks

---

## Known Limitations

### Not Yet Implemented
- ⏳ Historical transaction indexing (only recent)
- ⏳ Advanced VDF (using iterated SHA256)
- ⏳ State pruning
- ⏳ Snapshot sync
- ⏳ Mempool transaction broadcasting
- ⏳ Fee market (currently fixed)

### Testnet Constraints
- ⚠️ **NOT production-ready** - testnet only
- ⚠️ **No security audit** - use at own risk
- ⚠️ **May reset** - data not guaranteed
- ⚠️ **Geographic centralization** - 2 servers in same region
- ⚠️ **Limited farmers** - small network

---

## Performance Metrics

### Block Production
- **Average block time:** 25 seconds
- **Fastest block:** <1 second
- **Longest gap:** ~2 minutes
- **99th percentile:** <60 seconds

### Farming
- **Server A (1 plot):** ~83% win rate (localhost advantage)
- **Server C (6 plots):** ~17% win rate
- **Expected distribution:** 14% / 86% (based on plot count)

### API Response Times
- `/chainTip`: <50ms
- `/account/<addr>`: <100ms
- `/blocks/recent`: <200ms
- `/block/<height>`: <150ms

---

## Next Steps

### Immediate (Nov-Dec 2025)
- Publish TypeScript SDK to npm
- Add faucet for testnet RCHV
- Improve transaction history indexing
- Add more geographic diversity (3rd farmer)

### Short-term (Q1 2026)
- Implement Wesolowski VDF
- Add state pruning
- Snapshot sync for fast bootstrapping
- Fee market dynamics
- Mempool improvements

### Medium-term (Q2 2026)
- Security audit
- Performance optimization
- Economic model finalization
- Community governance
- Public mainnet preparation

---

## Resources

**Documentation:**
- [API Reference](api-reference.md)
- [SDK Guide](sdk-guide.md)
- [Farming Guide](setup-farmer.md)

**Live Services:**
- [Public RPC](https://seed.archivas.ai)
- [Block Explorer](https://archivas-explorer-production.up.railway.app)
- [GitHub](https://github.com/ArchivasNetwork/archivas)

**Monitoring:**
- [Grafana Dashboard](http://57.129.148.132:3001)
- [Current Status](../CURRENT-STATUS.md)

---

**Archivas testnet is operational and ready for developers and farmers to join!** 🌾

