# Introduction to Archivas

## What is Archivas?

Archivas is a **Proof-of-Space-and-Time Layer 1 blockchain** that secures its network with disk space and sequential time proofs instead of energy-intensive computation or capital lockup.

### Why Archivas?

**Traditional blockchains have problems:**
- **Proof-of-Work:** Wastes enormous energy (~150 TWh/year for Bitcoin)
- **Proof-of-Stake:** "Rich get richer" with high capital barriers

**Archivas uses Proof-of-Space-and-Time:**
- 🌾 **Permissionless** - Anyone with disk space can farm
- ⚡ **Energy Efficient** - No wasteful computation, just disk storage
- 🔐 **Secure** - VDF prevents grinding, PoSpace prevents centralization
- ⏰ **Fair** - Disk space determines odds, not wealth or specialized hardware

### How It Works

**Proof-of-Space:**
- Farmers create "plots" - large files filled with precomputed hashes
- When a block is needed, farmers search their plots for the best "proof"
- The farmer with the winning proof produces the block and earns RCHV

**Verifiable Delay Functions (VDF):**
- Timelords continuously compute a sequential function that takes real time
- VDF output is used to generate challenges
- Cannot be parallelized or skipped - provides temporal security

**Combined Security:**
- Blocks require BOTH a winning Proof-of-Space AND valid VDF proof
- Same consensus class as Chia Network
- Grinding-resistant, temporally ordered, fair lottery

## Key Features

### Consensus
- **Proof-of-Space:** Disk-based block lottery
- **VDF:** Sequential time proofs (iterated SHA-256 for devnet)
- **Adaptive Difficulty:** Self-regulating ~20 second block times
- **Block Rewards:** 20 RCHV per block

### Technology
- **Language:** Go 1.21+
- **Cryptography:** secp256k1 (ECDSA signatures)
- **Addresses:** Bech32 format (arcv prefix)
- **Storage:** BadgerDB (embedded key-value store)
- **Networking:** TCP P2P with JSON protocol

### Token
- **Symbol:** RCHV
- **Decimals:** 8
- **Block Reward:** 20.00000000 RCHV
- **Target Block Time:** ~20 seconds
- **Supply:** Determined by block rewards and halving schedule (TBD)

## Current Status

**Testnet v0.1.0:** Multi-node operational

- **Nodes:** 2+ (synchronized)
- **Height:** 78+ blocks
- **Status:** 🟢 LIVE
- **Network:** Public testnet (join anytime!)

## Same Class as Chia Network

Archivas implements the same fundamental consensus model as Chia Network:
- Proof-of-Space for permissionless participation
- VDF for temporal security and grinding resistance
- Adaptive difficulty for consistent block times

Built from scratch in Go with a modular, extensible architecture.

## Use Cases

**Current (Testnet):**
- Learn about Proof-of-Space consensus
- Test disk-based farming
- Experiment with VDF
- Build blockchain applications

**Future (Mainnet):**
- Permissionless data storage layer
- Energy-efficient value transfer
- Storage-backed smart contracts
- Decentralized archival network

## Community

- **GitHub:** https://github.com/ArchivasNetwork/archivas
- **Discussions:** https://github.com/ArchivasNetwork/archivas/discussions
- **Issues:** https://github.com/ArchivasNetwork/archivas/issues

**Join us in building the future of storage-based consensus!** 🌾

---

**Next:** [Getting Started →](Getting-Started.md)

