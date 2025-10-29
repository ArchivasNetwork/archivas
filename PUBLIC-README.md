# Archivas

**A Proof-of-Space-and-Time L1 blockchain. Farm RCHV with disk space.**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![Status](https://img.shields.io/badge/status-devnet-yellow.svg)](https://github.com/iljanemesis/archivas)

## What is Archivas?

Archivas is a Layer 1 blockchain that uses **Proof-of-Space-and-Time** consensus. Instead of burning electricity (Proof-of-Work) or requiring capital lockup (Proof-of-Stake), Archivas secures its network with:

- **Disk space** (Proof-of-Space) - Anyone with storage can farm
- **Sequential time** (Verifiable Delay Functions) - Prevents grinding attacks

Farmers allocate disk space to create "plots," then compete to produce blocks and earn **RCHV**, the native token.

## Key Features

- 🌾 **Permissionless Farming** - No minimum stake, just disk space
- ⚡ **Energy Efficient** - No wasteful computation
- 🔐 **Cryptographic Security** - secp256k1 signatures, bech32 addresses
- ⏰ **Temporal Ordering** - VDF prevents precomputation attacks
- 📊 **Adaptive Difficulty** - Self-regulating ~20 second block times
- 💾 **Persistent State** - Survives restarts, runs 24/7

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Linux, macOS, or Windows

### Installation

```bash
git clone https://github.com/iljanemesis/archivas
cd archivas
go mod download
```

### Build

```bash
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-wallet ./cmd/archivas-wallet
```

### Farm RCHV

**1. Generate a farmer wallet:**
```bash
./archivas-wallet new
```

**2. Create a plot:**
```bash
./archivas-farmer plot --size 20 --path ./plots
```

Plot sizes:
- k=16: ~2 MB (testing)
- k=20: ~32 MB (small farm)
- k=24: ~512 MB (medium farm)
- k=28: ~8 GB (large farm)

**3. Start the node:**
```bash
./archivas-node
```

**4. Farm blocks:**
```bash
./archivas-farmer farm \
  --plots ./plots \
  --farmer-key <your_private_key_hex>
```

**5. Check your rewards:**
```bash
curl http://localhost:8080/balance/<your_address>
```

Every block you farm = **20 RCHV**!

## How It Works

### Proof-of-Space

Farmers create "plots" - large files filled with precomputed cryptographic hashes. When a new block is needed, the network broadcasts a "challenge." Farmers scan their plots to find the best "proof" (lowest quality hash). The farmer with the winning proof produces the block and earns the reward.

**Think of it like:** Your hard drive is a lottery ticket. More disk space = more tickets = higher chance to win.

### Verifiable Delay Functions (VDF)

A timelord process continuously computes a sequential function that takes real time. This VDF output is used to generate challenges, ensuring:

- **No grinding** - Can't precompute future blocks
- **Temporal security** - Blocks have provable time ordering
- **Fair lottery** - Challenge unpredictable until VDF computed

### Combined Security

Archivas requires BOTH:
- A winning Proof-of-Space (disk space)
- A valid VDF proof (sequential time)

This is the same security model as Chia Network.

## Chain Parameters (Devnet)

| Parameter | Value |
|-----------|-------|
| Chain ID | 1616 |
| Native Token | RCHV |
| Decimals | 8 |
| Block Time | ~20 seconds |
| Block Reward | 20.00000000 RCHV |
| Address Prefix | arcv |
| Consensus | Proof-of-Space-and-Time |

## Architecture

```
┌─────────────────┐
│ archivas-node   │ ← Validates PoSpace+VDF, manages state
└────────┬────────┘
         │
    ┌────┴─────────────────┐
    │                      │
┌───▼──────────┐  ┌───────▼────────┐
│ timelord     │  │    farmer      │
│ (VDF)        │  │   (PoSpace)    │
└──────────────┘  └────────────────┘
```

**Node** - Validates blocks, manages chain state, exposes RPC  
**Timelord** - Computes VDF, publishes time proofs  
**Farmer** - Creates plots, finds PoSpace proofs, submits blocks  

## API

### Query Balance
```bash
GET /balance/<address>
```

```json
{
  "address": "arcv1q84xt5pzcslhnjsc2h2t9cnuxrn0e2u2u97jnl",
  "balance": 12000000000,
  "nonce": 0
}
```

### Submit Transaction
```bash
POST /submitTx
```

```json
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

### Get Current Challenge (for farmers)
```bash
GET /challenge
```

```json
{
  "challenge": "c2566d51d073bb62...",
  "difficulty": 1125899906842624,
  "height": 1
}
```

## Development Status

### ✅ Implemented
- [x] Proof-of-Space consensus
- [x] Verifiable Delay Functions
- [x] Cryptographic wallets (secp256k1)
- [x] Transaction signing
- [x] Adaptive difficulty
- [x] Block rewards
- [x] Persistent storage
- [x] RPC API

### 🚧 In Progress
- [ ] VDF mode activation (Milestone 4 ready, needs testing)
- [ ] P2P networking (Milestone 6)
- [ ] Multi-node testnet

### 📋 Planned
- [ ] Database compaction
- [ ] State snapshots
- [ ] Block explorer
- [ ] Wesolowski/Pietrzak VDF (mainnet)

## Project Structure

```
archivas/
├── cmd/              Command-line binaries
│   ├── archivas-node/
│   ├── archivas-farmer/
│   ├── archivas-timelord/
│   └── archivas-wallet/
├── consensus/        Difficulty, challenges, validation
├── pospace/          Proof-of-Space implementation
├── vdf/              Verifiable Delay Functions
├── ledger/           State management, transactions
├── wallet/           Cryptography, signing
├── mempool/          Transaction pool
├── rpc/              HTTP API
├── storage/          Persistent storage (BadgerDB)
└── config/           Chain parameters
```

## Consensus Algorithm

### Block Production

1. **Timelord** computes VDF from current chain tip
2. **Node** generates challenge from VDF output
3. **Farmer** searches plots for best proof
4. **Farmer** submits block with PoSpace proof + VDF info
5. **Node** validates both proofs, accepts block
6. **Farmer** receives 20 RCHV block reward
7. **Repeat**

### Security

- **Grinding Resistance** - VDF takes real sequential time
- **Fairness** - Disk space determines lottery odds
- **Temporal Ordering** - Blocks have provable time sequence
- **Cryptographic Ownership** - Only private keys can spend

## Roadmap

### Phase 1: Devnet (Current)
- ✅ Single-node operation
- ✅ PoSpace farming working
- ✅ VDF implementation ready
- ✅ Persistent storage

### Phase 2: Testnet (Next)
- [ ] P2P networking
- [ ] Multi-node consensus
- [ ] Public testnet launch
- [ ] Faucet for test RCHV

### Phase 3: Mainnet (Future)
- [ ] Security audit
- [ ] Wesolowski VDF
- [ ] Token distribution
- [ ] Public launch

## Community

- **GitHub:** [github.com/iljanemesis/archivas](https://github.com/iljanemesis/archivas)
- **Docs:** See `/docs` folder or `JOURNEY.md`
- **Issues:** [github.com/iljanemesis/archivas/issues](https://github.com/iljanemesis/archivas/issues)

## Contributing

Archivas is in active development. Contributions welcome!

Areas of focus:
- P2P networking
- Performance optimization
- Testing and security
- Documentation

## License

Apache 2.0 (TBD - choose your license)

## Acknowledgments

Inspired by:
- Chia Network (PoSpace+Time consensus)
- Filecoin (storage-based consensus)
- Bitcoin (UTXO model, cryptographic security)

Built with:
- Go (systems language)
- BadgerDB (key-value store)
- secp256k1 (cryptography)

---

**Archivas: Farming the future of decentralized storage.** 🌾

