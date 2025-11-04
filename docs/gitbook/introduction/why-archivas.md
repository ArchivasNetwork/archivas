# Why Archivas?

Archivas is not the first Proof-of-Space-and-Time blockchain. So why build another one?

---

## What Makes Archivas Different

### 1. Built in Go (Not Python)

**Chia Network:** Python-based, complex codebase  
**Archivas:** Clean Go implementation

**Advantages:**
- ✅ **Performance:** Go's concurrency and speed
- ✅ **Simplicity:** ~10K lines vs Chia's 100K+
- ✅ **Type Safety:** Compile-time checks
- ✅ **Single Binary:** Easy deployment
- ✅ **Cross-platform:** Works everywhere Go does

---

### 2. Modular & Extensible

Archivas is designed to be hackable and extensible.

**Core Modules:**
```
archivas/
├── consensus/     # Pluggable difficulty algorithms
├── pospace/       # PoSpace verification
├── vdf/           # VDF (can swap implementations)
├── storage/       # Storage engine (BadgerDB)
├── p2p/           # Network protocol
├── rpc/           # Public API
├── wallet/        # Cryptographic primitives
└── pkg/           # Shared utilities
```

Each module has clear interfaces and can be swapped or upgraded independently.

---

### 3. Developer-First

**Public API from Day 1:**
- ✅ HTTPS endpoint with TLS
- ✅ CORS enabled for web apps
- ✅ TypeScript SDK
- ✅ Block explorer
- ✅ Complete documentation
- ✅ Rate limiting and security

**Chia Network:** Requires running a full node to interact  
**Archivas:** Public API at https://seed.archivas.ai

---

### 4. Modern Cryptography

**Wallets:**
- ✅ **Ed25519** (not secp256k1)
- ✅ **BIP39** mnemonics (24 words)
- ✅ **SLIP-0010** key derivation
- ✅ **Bech32** addresses (`arcv` prefix)
- ✅ **Blake2b** hashing
- ✅ **RFC 8785** canonical JSON

**Benefits:**
- Industry-standard wallet compatibility
- Hardware wallet support (future)
- Smaller signatures (64 bytes)
- Faster verification

---

### 5. Production Infrastructure

**From the start:**
- ✅ Public HTTPS API
- ✅ Prometheus metrics
- ✅ Grafana dashboards
- ✅ Nginx reverse proxy
- ✅ Let's Encrypt TLS
- ✅ Rate limiting
- ✅ Health checks

**Not an afterthought - built for production from day 1.**

---

### 6. Proven Reliability

**25+ days continuous operation:**
- ✅ 64,000+ blocks without consensus failures
- ✅ Multi-server coordination working
- ✅ 1,160+ RCHV transferred successfully
- ✅ State sync and crash recovery tested
- ✅ 99.9% uptime

**This is not a prototype or demo - it's a working blockchain.**

---

## What Archivas Keeps from Chia

### Core Consensus
- ✅ Proof-of-Space lottery mechanism
- ✅ VDF for temporal ordering
- ✅ Challenge-response protocol
- ✅ Difficulty adjustment

### Security Model
- ✅ Grinding-resistant
- ✅ Nothing-at-stake protection
- ✅ Sybil-resistant
- ✅ 51% attack requires 51% of storage

---

## What Archivas Improves

### Simplicity
- **Chia:** Complex Python codebase, many dependencies
- **Archivas:** Clean Go code, minimal dependencies

### Accessibility
- **Chia:** Steep learning curve, requires deep blockchain knowledge
- **Archivas:** Clear documentation, simple API, TypeScript SDK

### Developer Experience
- **Chia:** Must run full node, complex RPC
- **Archivas:** Public HTTPS API, modern tooling, CORS support

### Deployment
- **Chia:** Multiple processes (farmer, timelord, full node, wallet)
- **Archivas:** Single `archivas-node` binary

---

## Use Cases

### Current (Testnet)

**For Developers:**
- Build wallets and payment apps
- Integrate blockchain storage
- Experiment with PoST consensus
- Learn cryptography (Ed25519, BIP39)

**For Farmers:**
- Test farming profitability
- Optimize plot sizes
- Learn about PoSpace
- Earn testnet RCHV

**For Researchers:**
- Study PoST consensus in practice
- Analyze difficulty algorithms
- Test VDF implementations
- Measure decentralization

### Future (Mainnet)

**Storage Layer:**
- Decentralized file storage
- Archival records
- IPFS/Arweave alternative
- Verifiable data availability

**Value Transfer:**
- Fast, low-fee payments
- Cross-border transfers
- Programmable money
- DeFi primitives

**Smart Contracts:**
- Storage-backed contracts
- Proof-of-Replication markets
- Data availability guarantees
- Decentralized CDN

---

## Comparison Matrix

| Feature | Bitcoin (PoW) | Ethereum (PoS) | Chia (PoST) | Archivas (PoST) |
|---------|---------------|----------------|-------------|-----------------|
| **Energy Use** | Very High | Low | Very Low | Very Low |
| **Capital Requirement** | Medium-High | Very High | Low | Low |
| **Hardware** | ASIC/GPU | Consumer | HDD/SSD | HDD/SSD |
| **Permissionless** | ✅ | ✅ | ✅ | ✅ |
| **Decentralized** | ⚠️ Mining pools | ⚠️ Staking services | ✅ | ✅ |
| **Language** | C++ | Go/Rust | Python | **Go** |
| **Public API** | ✅ | ✅ | Limited | **✅ HTTPS** |
| **TypeScript SDK** | 3rd party | 3rd party | Limited | **✅ Official** |
| **Block Time** | 10 min | 12 sec | 18-19 sec | **20-30 sec** |
| **Developer Friendly** | ⚠️ | ✅ | ⚠️ | **✅** |

---

## Vision

**Short-term:** Prove that PoST can work as a modular, developer-friendly L1.

**Medium-term:** Build a community of developers and farmers around energy-efficient consensus.

**Long-term:** Archivas becomes the go-to platform for storage-backed applications and value transfer.

---

## Why Now?

**The timing is right:**
- ✅ Chia proved PoST works (2021-2025)
- ✅ Storage is cheap and abundant ($15/TB)
- ✅ Developers want eco-friendly blockchains
- ✅ Users demand low fees and fast confirmations
- ✅ Go ecosystem is mature for blockchain development

**Archivas builds on these foundations to create something better.**

---

## For Chia Farmers

If you're already farming Chia, Archivas is easy to add:

**What's familiar:**
- Same consensus model (PoSpace + VDF)
- Same plot concept (k-size, hash tables)
- Same terminology (farmer, timelord, proof, quality)

**What's different:**
- **Lighter:** Go binary vs Python processes
- **Simpler:** One config file vs multiple
- **Modern API:** HTTPS + TypeScript vs local RPC

You can farm **both** Chia and Archivas with the same hardware!

---

## Join the Movement

Archivas is **live and operational** - not vaporware, not a whitepaper.

✅ 25+ days uptime  
✅ 64,000+ blocks  
✅ Public API  
✅ Block explorer  
✅ TypeScript SDK  

**Start building today:** https://seed.archivas.ai

---

**Next:** Check the [Network Status](network-status.md) to see live statistics.

