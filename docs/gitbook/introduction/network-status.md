# Network Status

**Last Updated:** November 3, 2025  
**Network:** archivas-devnet-v4  
**Status:** 🟢 OPERATIONAL  

---

## Live Statistics

### Blockchain

| Metric | Value | Status |
|--------|-------|--------|
| **Block Height** | 64,000+ | 🟢 Growing |
| **Uptime** | 25+ days | 🟢 Stable |
| **Total Supply** | ~1,280,000 RCHV | 📈 Increasing |
| **Block Time** | ~25 seconds | ✅ Consistent |
| **Difficulty** | 1,000,000 | 🔄 Adaptive |

### Network Health

| Metric | Value | Status |
|--------|-------|--------|
| **Active Farmers** | 2 servers | 🟢 Mining |
| **Total Plots** | 7 × k28 (~55 GB) | 🌾 Farming |
| **Connected Peers** | 2-3 | 🌐 Synced |
| **Mempool** | 0-1 pending | 💤 Clear |

### Infrastructure

| Service | URL | Status |
|---------|-----|--------|
| **Public RPC** | https://seed.archivas.ai | 🟢 200 OK |
| **Block Explorer** | https://archivas-explorer.up.railway.app | 🟢 Live |
| **Grafana** | http://57.129.148.132:3001 | 🟢 Metrics |
| **Prometheus** | localhost:9091 | 🟢 Scraping |

---

## Real-Time Data

### Current Block

```bash
curl https://seed.archivas.ai/chainTip
```

Response:
```json
{
  "height": "64103",
  "hash": "0x3f4a...",
  "difficulty": "1000000"
}
```

### Recent Blocks

```bash
curl https://seed.archivas.ai/blocks/recent?limit=5
```

Shows latest 5 blocks with farmer addresses and transaction counts.

---

## Farming Statistics

### Server A (Seed Node)
- **Address:** `arcv1t3huuyd08er3yfnmk9c935rmx3wdh5j6m2uc9d`
- **Balance:** ~1,062,000 RCHV
- **Plots:** 1 × k28 (268M hashes, ~8 GB)
- **Blocks Won:** ~53,100 (~83%)
- **Location:** 57.129.148.132
- **Uptime:** 608+ hours

### Server C (Farmer Node)
- **Address:** `arcv1xjgsguj9e4assk4a24pfkm6cl92jwrgfxsru7c`
- **Balance:** ~197,000 RCHV
- **Plots:** 6 × k28 (1.6B hashes, ~48 GB)
- **Blocks Won:** ~9,850 (~15%)
- **Location:** 57.129.148.134
- **Uptime:** 486+ hours

**Note:** Server A has localhost advantage, resulting in higher win rate.

---

## Transaction Activity

### Confirmed Transfers

**Total Transferred:** 1,160+ RCHV

**Recent Transactions:**
1. 300 RCHV → Test Wallet (nonce 0)
2. 100 RCHV → Test Wallet 3 (nonce 1)
3. 100 RCHV → External (nonce 2)
4. 300 RCHV → Test Wallet (nonce 3)
5. 50 RCHV → Test Wallet 3 (nonce 4)
6. 10 RCHV → Server C → Server A (nonce 0)
7. 10 RCHV → Test Wallet 3 (nonce 5)
8. 300 RCHV → Test Wallet (nonce 6)
9. 25 RCHV → Test Wallet (nonce 7)

**All transactions confirmed and balances verified!**

---

## Performance Metrics

### Block Production

- **Average block time:** 25 seconds
- **Fastest block:** <1 second  
- **Longest gap:** ~2 minutes
- **Consistency:** 99% within 30 seconds

### API Response Times

- `/chainTip`: ~30-50ms
- `/account/<addr>`: ~50-100ms
- `/blocks/recent`: ~100-200ms
- `/block/<height>`: ~80-150ms
- `/tx/recent`: ~150-250ms

### Network Latency

- **Server A ↔ Server C:** <5ms (same provider)
- **Public API (worldwide):** 50-300ms

---

## Monitoring Dashboards

### Grafana Dashboard

**URL:** http://57.129.148.132:3001

**Panels:**
- Tip Height (real-time)
- Connected Peers
- Mining Difficulty
- Block Production Rate
- RPC Request Volume
- System Resources

### Prometheus Metrics

**Available internally** (not exposed publicly for security):
- `archivas_tip_height`
- `archivas_difficulty`
- `archivas_peer_count`
- `archivas_blocks_total`
- `archivas_submit_accepted_total`
- `archivas_rpc_requests_total`

---

## Security Status

### Infrastructure
- ✅ TLS 1.2+ with auto-renewing certificate
- ✅ HTTP/2 enabled
- ✅ Rate limiting (10 req/min on /submit)
- ✅ CORS configured for browsers
- ✅ Internal metrics blocked from public
- ✅ Firewall rules active

### Consensus
- ✅ 64,000+ blocks without issues
- ✅ No reorganizations
- ✅ No double-spends
- ✅ All proofs verified
- ✅ Difficulty adjusting correctly

### Cryptography
- ✅ Ed25519 signatures (audited algorithm)
- ✅ Blake2b hashing (secure, fast)
- ✅ BIP39 mnemonics (industry standard)
- ✅ RFC 8785 canonical JSON (deterministic)

---

## Known Issues & Limitations

### Testnet Constraints
- ⚠️ **Not audited:** Use at own risk
- ⚠️ **May reset:** Data not guaranteed permanent
- ⚠️ **Geographic centralization:** 2 servers in same region
- ⚠️ **Small network:** Only 2 farmers

### Missing Features
- ⏳ Advanced VDF (Wesolowski/Pietrzak)
- ⏳ State pruning
- ⏳ Snapshot sync
- ⏳ Dynamic fee market
- ⏳ Historical transaction indexing

### Under Development
- 🔄 Explorer improvements (full tx history)
- 🔄 SDK published to npm
- 🔄 Faucet for testnet RCHV
- 🔄 More geographic distribution

---

## Comparison: Archivas vs Chia

| Feature | Chia Network | Archivas |
|---------|-------------|----------|
| **Launch** | May 2021 | Nov 2025 |
| **Language** | Python | Go |
| **Plot Format** | Custom binary | Custom binary |
| **Wallet** | BLS12-381 | Ed25519 |
| **Public API** | Limited | HTTPS + CORS |
| **TypeScript SDK** | Community | Official |
| **Block Explorer** | Multiple 3rd party | Official (Next.js) |
| **Monitoring** | Basic | Prometheus + Grafana |
| **Block Time** | 18-19 sec | 20-30 sec |
| **Block Reward** | 2 XCH | 20 RCHV |
| **Uptime (current)** | 3.5+ years | 25+ days |
| **Maturity** | Production | Testnet |

---

## Why Choose Archivas?

### For Developers
- ✅ Modern API (HTTPS, JSON, CORS)
- ✅ TypeScript SDK with full types
- ✅ Clear documentation
- ✅ Fast iteration (small codebase)
- ✅ Active development

### For Farmers
- ✅ Simple setup (one binary)
- ✅ Low overhead (Go performance)
- ✅ Same economics as Chia (space-based)
- ✅ Can run alongside Chia

### For Users
- ✅ Fast transactions (~25 seconds)
- ✅ Low fees (~0.001 RCHV)
- ✅ Standard wallets (BIP39)
- ✅ Block explorer

---

## Goals

### Technical Goals
1. Prove PoST works in Go
2. Build modular, extensible architecture
3. Create developer-friendly API
4. Achieve production reliability

### Community Goals
1. Educate about PoST consensus
2. Build tools ecosystem (SDK, explorer, wallets)
3. Grow farmer community
4. Foster open development

### Long-term Vision
1. Mainnet launch with security audit
2. Storage-backed smart contracts
3. Decentralized archival network
4. Integration with existing storage protocols

---

## Success Metrics

**What we've proven:**
- ✅ PoST consensus works reliably
- ✅ Multi-server coordination is stable
- ✅ Public API can handle production load
- ✅ Transactions process correctly
- ✅ State persistence works
- ✅ Difficulty adjustment stabilizes block times

**Next to prove:**
- Geographic distribution (3+ regions)
- Scale to 10+ farmers
- Handle high transaction volume
- Advanced VDF (Wesolowski)
- Economic sustainability

---

## Get Involved

**Try it yourself:**
1. Query the API: `curl https://seed.archivas.ai/chainTip`
2. Use the SDK: `npm install @archivas/sdk`
3. Browse blocks: https://archivas-explorer.up.railway.app
4. View metrics: http://57.129.148.132:3001

**Contribute:**
- GitHub: https://github.com/ArchivasNetwork/archivas
- Issues: Report bugs or request features
- PRs: Contribute code improvements

---

**Next:** Learn how to [Get Started](../getting-started/quick-start.md) with Archivas!

