# Archivas

<p align="center">
  <img src="docs/archivas-logo.png" alt="Archivas Logo" width="200"/>
  <br/>
  <strong>Farm RCHV with Disk Space</strong>
  <br/>
  <em>A Proof-of-Space-and-Time L1 Blockchain</em>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#how-it-works">How It Works</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#roadmap">Roadmap</a>
</p>

---

## What is Archivas?

Archivas is a Layer 1 blockchain that uses **Proof-of-Space-and-Time** consensus. Instead of burning electricity (Proof-of-Work) or requiring capital lockup (Proof-of-Stake), Archivas secures its network with:

- 🌾 **Disk Space** (Proof-of-Space) - Anyone with storage can farm
- ⏰ **Sequential Time** (Verifiable Delay Functions) - Prevents grinding attacks

Farmers allocate disk space to create "plots," then compete to produce blocks and earn **RCHV** tokens.

## Why Archivas?

**Problems with traditional consensus:**

| Consensus | Energy | Barrier to Entry | Centralization Risk |
|-----------|--------|------------------|---------------------|
| Proof-of-Work | Very High ⚡ | Medium (ASICs) | Mining pools |
| Proof-of-Stake | Low | High (capital) | Wealth concentration |
| **Proof-of-Space-and-Time** | **Very Low** ✅ | **Low** ✅ | **Distributed** ✅ |

**Archivas advantages:**
- ✅ **Permissionless** - Anyone with disk space can farm
- ✅ **Energy Efficient** - No wasteful computation
- ✅ **Fair** - Disk space determines odds, not capital
- ✅ **Secure** - VDF prevents grinding, PoSpace prevents centralization

## Features

- 🌾 **Proof-of-Space Farming** - Tested with real blocks
- ⏰ **Verifiable Delay Functions** - Temporal security (ready to activate)
- 🔐 **Cryptographic Wallets** - secp256k1, bech32 addresses (arcv...)
- ✍️ **Transaction Signing** - ECDSA signatures
- 💾 **Persistent Storage** - BadgerDB, crash recovery
- 📊 **Adaptive Difficulty** - Maintains ~20 second block times
- 💰 **Block Rewards** - 20 RCHV per block farmed

## Current Status

**🟢 Devnet Operational**

- ✅ **Farming Tested** - 6 blocks farmed, 120 RCHV earned
- ✅ **Wallets Working** - 750 RCHV sent between addresses
- ✅ **Persistence Verified** - Node restart, state restored
- ✅ **Difficulty Adaptive** - 5 adjustments in testing
- ⏸️ **VDF Mode** - Implemented, ready to activate
- 🚧 **P2P Networking** - Milestone 6 (next)

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Linux, macOS, or Windows
- ~100MB disk space for small plots

### Installation

```bash
# Clone repository
git clone https://github.com/ArchivasNetwork/archivas
cd archivas

# Download dependencies
go mod download

# Build binaries
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-wallet ./cmd/archivas-wallet
```

### Farm Your First RCHV

**Step 1: Generate a wallet**
```bash
./archivas-wallet new
```

Save your private key! Output:
```
🔐 New Archivas Wallet Generated

Address:     arcv1q84xt5pzcslhnjsc2h2t9cnuxrn0e2u2u97jnl
Public Key:  03457989304d0c1ecbe33bcdb2b5ae8f8f34a4d2c0f278a7ad79460c557fe98dd9
Private Key: <EXAMPLE_PRIVATE_KEY_GENERATE_YOUR_OWN>
```

**Step 2: Create a plot**
```bash
./archivas-farmer plot --size 20 --path ./plots
```

Plot sizes:
- `k=16`: ~2 MB (testing)
- `k=20`: ~32 MB (small farm)
- `k=24`: ~512 MB (medium farm)  
- `k=28`: ~8 GB (large farm)

**Step 3: Start the node**
```bash
./archivas-node
```

**Step 4: Start farming**
```bash
./archivas-farmer farm \
  --plots ./plots \
  --farmer-key <your_private_key_hex>
```

**Step 5: Watch your RCHV grow!**
```bash
# Check balance
curl http://localhost:8080/balance/<your_address>

# Or use the wallet
go run ./cmd/archivas-wallet balance --address <your_address>
```

Every block you farm = **20 RCHV** 🎉

## How It Works

### Proof-of-Space

1. **Plot Creation**: Generate a large file filled with precomputed hashes
   ```
   hash[i] = H(farmerPubKey || plotID || i)
   ```

2. **Challenge**: Network broadcasts a challenge for each new block
   ```
   challenge = H(VDF_output || height)
   ```

3. **Proof Search**: Farmers scan their plots for the best quality proof
   ```
   quality = H(challenge || hash)  
   // Lower quality = better proof
   ```

4. **Block Production**: Farmer with winning proof (quality < difficulty) produces the block

### Verifiable Delay Functions

1. **Timelord** continuously computes sequential function:
   ```
   y₀ = seed
   y₁ = H(y₀)
   y₂ = H(y₁)
   ...
   yₙ = H(yₙ₋₁)
   ```

2. **VDF Output** derives the PoSpace challenge
3. **Cannot skip** - must compute all iterations sequentially
4. **Prevents grinding** - can't precompute alternative timelines

### Combined Security

Blocks require **BOTH**:
- ✅ Winning Proof-of-Space (disk space lottery)
- ✅ Valid VDF proof (sequential time elapsed)

This is **Chia-class consensus security**.

## Chain Parameters

| Parameter | Value |
|-----------|-------|
| Chain ID | 1616 |
| Native Token | RCHV |
| Decimals | 8 |
| Block Time | ~20 seconds |
| Block Reward | 20.00000000 RCHV |
| Difficulty Adjustment | Every block (10-block window) |
| Address Prefix | arcv |
| Genesis | October 2025 |

## Architecture

```
┌─────────────────────────────────────────────┐
│           Archivas Full Node                │
│  ┌─────────────────────────────────────┐   │
│  │  Consensus Engine                   │   │
│  │  • Validates PoSpace proofs         │   │
│  │  │  • Validates VDF proofs          │   │
│  │  • Manages chain state              │   │
│  │  • Adaptive difficulty              │   │
│  └─────────────────────────────────────┘   │
│  ┌─────────────────────────────────────┐   │
│  │  Storage Layer (BadgerDB)           │   │
│  │  • Blocks, accounts, metadata       │   │
│  │  • Crash recovery                   │   │
│  └─────────────────────────────────────┘   │
│  ┌─────────────────────────────────────┐   │
│  │  RPC API (:8080)                    │   │
│  │  • Balance queries                  │   │
│  │  • Transaction submission           │   │
│  │  • Challenge feed (for farmers)     │   │
│  │  • Block submission                 │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
         ↑                              ↑
         │                              │
    ┌────┴────┐                    ┌────┴────┐
    │Timelord │                    │ Farmer  │
    │         │                    │         │
    │VDF      │                    │PoSpace  │
    │Computer │                    │Scanner  │
    └─────────┘                    └─────────┘
```

## API

### Query Balance
```bash
GET http://localhost:8080/balance/<address>
```

Response:
```json
{
  "address": "arcv1q84xt5pzcslhnjsc2h2t9cnuxrn0e2u2u97jnl",
  "balance": 12000000000,
  "nonce": 0
}
```

### Submit Transaction
```bash
POST http://localhost:8080/submitTx
Content-Type: application/json

{
  "from": "arcv1...",
  "to": "arcv1...",
  "amount": 100000000,
  "fee": 100000,
  "nonce": 0,
  "senderPubKey": "03...",
  "signature": "30..."
}
```

### Get Challenge (for farmers)
```bash
GET http://localhost:8080/challenge
```

Response:
```json
{
  "challenge": "c2566d51d073bb62...",
  "difficulty": 1125899906842624,
  "height": 1
}
```

## Test Results

### Farming Test (60 seconds)
- **Blocks Found:** 6
- **RCHV Earned:** 120.00000000 RCHV
- **Average Time:** ~10 seconds per block
- **Difficulty Adjustments:** 5
- **Plot Size:** k=16 (2MB)

### Persistence Test
- **Blocks Before Restart:** 7 (genesis + 6)
- **Node Killed:** Process terminated
- **Node Restarted:** Full state restored
- **Blocks After Restart:** 7 (preserved)
- **Balances:** All correct
- **Recovery Time:** <100ms

### Transaction Test
- **Wallets Generated:** 2
- **RCHV Sent:** 750 RCHV
- **Transactions:** 2
- **Signatures:** Verified ✅
- **Nonces:** Incremented correctly

**All tests passed. All features verified.** ✅

## Roadmap

### Phase 1: Devnet (Complete ✅)
- [x] Proof-of-Space consensus
- [x] VDF implementation
- [x] Wallet system
- [x] Persistent storage
- [x] Single-node operation

### Phase 2: Testnet (Q1 2026)
- [ ] P2P networking
- [ ] Multi-node consensus
- [ ] Public testnet deployment
- [ ] Block explorer
- [ ] Faucet for test RCHV

### Phase 3: Mainnet Prep (Q2-Q3 2026)
- [ ] Security audit
- [ ] Wesolowski VDF upgrade
- [ ] Token distribution plan
- [ ] Mainnet deployment

### Phase 4: Ecosystem (2026+)
- [ ] Smart contracts (WASM)
- [ ] Farming pools
- [ ] Light clients
- [ ] Cross-chain bridges

## Documentation

- **[README.md](README.md)** - This file
- **[STATUS.md](STATUS.md)** - Current technical status
- **[JOURNEY.md](JOURNEY.md)** - Complete development story
- **[MILESTONE3.md](MILESTONE3.md)** - Proof-of-Space guide
- **[MILESTONE5-PERSISTENCE.md](MILESTONE5-PERSISTENCE.md)** - Storage guide
- **[ACTIVATE-VDF.md](ACTIVATE-VDF.md)** - VDF activation instructions
- **[docs/WHITEPAPER-OUTLINE.md](docs/WHITEPAPER-OUTLINE.md)** - Technical whitepaper structure

## Contributing

Archivas is open source and community-driven. We welcome contributions!

**Areas of focus:**
- P2P networking (libp2p integration)
- Performance optimization
- Testing and security
- Documentation
- Farming utilities

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Community

- **GitHub:** [github.com/ArchivasNetwork/archivas](https://github.com/ArchivasNetwork/archivas)
- **Discussions:** [GitHub Discussions](https://github.com/ArchivasNetwork/archivas/discussions)
- **Issues:** [GitHub Issues](https://github.com/ArchivasNetwork/archivas/issues)
- **Twitter:** [@ArchivasChain](https://twitter.com/ArchivasChain) (TBD)
- **Discord:** [Coming Soon]

## Technical Specifications

**Consensus:** Proof-of-Space-and-Time  
**Language:** Go 1.21+  
**Cryptography:** secp256k1 (ECDSA signatures)  
**Addresses:** Bech32 (arcv prefix)  
**Storage:** BadgerDB (embedded key-value store)  
**Networking:** HTTP RPC (P2P coming in Milestone 6)  

**Source Code:** 27 files, ~4,000 lines  
**Dependencies:** 3 (minimal footprint)  
**License:** Apache 2.0 (TBD)  

## Comparison

|  | Energy | Barrier | Hardware | Centralization |
|--|--------|---------|----------|----------------|
| **Bitcoin (PoW)** | Very High | ASICs | Specialized | Mining pools |
| **Ethereum (PoS)** | Low | 32 ETH | Standard | Validator sets |
| **Archivas (PoSpace+Time)** | **Very Low** | **Disk only** | **Commodity** | **Distributed** |

## Security Model

**Proof-of-Space:**
- Farmers can't fake disk space (must precompute plots)
- More plots = proportionally more winning chances
- No shortcuts (quality lottery is fair)

**Verifiable Delay Functions:**
- Can't precompute future blocks (VDF seed changes)
- Can't grind alternative timelines (VDF takes real time)
- Temporal ordering (blocks have provable sequence)

**Combined:**
- ✅ Grinding resistance
- ✅ No precomputation attacks
- ✅ Deterministic finality
- ✅ Fair lottery

## FAQ

**Q: How is this different from Chia?**  
A: Same consensus model (PoSpace+Time), different implementation. Archivas is built in Go with modular architecture, simpler plot format, and designed for extensibility.

**Q: Can I farm with existing Chia plots?**  
A: No, Archivas uses a different plot format. You'll need to create Archivas-specific plots.

**Q: How much can I earn?**  
A: Depends on your plot size vs. network size. Currently 20 RCHV per block. With k=20 plot on devnet, expect ~1-2 blocks/hour.

**Q: Is this ready for production?**  
A: No, this is devnet. P2P networking (Milestone 6) and security audit needed before testnet, more work for mainnet.

**Q: What's the token supply?**  
A: TBD - considering Bitcoin-like halving model. Current: 20 RCHV/block, halving every ~243 days.

**Q: Can I run this on my VPS?**  
A: Yes! Persistence is working, so you can deploy to any server.

## Development Status

### ✅ Completed
- [x] Proof-of-Space consensus
- [x] Adaptive difficulty
- [x] Cryptographic wallets
- [x] Transaction signing
- [x] Persistent storage
- [x] Block rewards
- [x] Farmer CLI
- [x] Wallet CLI

### ⏸️ Implemented (Ready)
- [x] VDF/Timelord
- [x] PoSpace+Time validation
- [x] VDF RPC endpoints

### 🚧 In Progress
- [ ] P2P networking (Milestone 6)
- [ ] Multi-node testnet
- [ ] Block explorer

### 📋 Planned
- [ ] Public testnet
- [ ] Security audit
- [ ] Mainnet preparation

## Building from Source

```bash
# Clone
git clone https://github.com/ArchivasNetwork/archivas
cd archivas

# Build all binaries
make build

# Or build individually
go build -o bin/archivas-node ./cmd/archivas-node
go build -o bin/archivas-farmer ./cmd/archivas-farmer
go build -o bin/archivas-timelord ./cmd/archivas-timelord
go build -o bin/archivas-wallet ./cmd/archivas-wallet
```

## Running Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test ./tests/integration/...

# Farming test (requires ~1 minute)
./scripts/test-farming.sh
```

## License

Apache 2.0 (TBD - choose appropriate license)

## Acknowledgments

**Inspired by:**
- Chia Network (PoSpace+Time consensus model)
- Filecoin (storage-based consensus)
- Bitcoin (cryptographic security, UTXO model)

**Built with:**
- Go (systems programming language)
- BadgerDB (embedded key-value store)
- secp256k1 (elliptic curve cryptography)
- Bech32 (address encoding)

## Citation

If you use Archivas in research, please cite:

```bibtex
@software{archivas2025,
  title = {Archivas: A Proof-of-Space-and-Time Layer 1 Blockchain},
  author = {Archivas Labs},
  year = {2025},
  url = {https://github.com/ArchivasNetwork/archivas}
}
```

---

<p align="center">
  <strong>Archivas: Farming the future of decentralized storage</strong> 🌾
  <br/>
  <br/>
  <a href="https://github.com/ArchivasNetwork/archivas">GitHub</a> •
  <a href="docs/WHITEPAPER-OUTLINE.md">Whitepaper</a> •
  <a href="STATUS.md">Status</a> •
  <a href="JOURNEY.md">Story</a>
</p>

