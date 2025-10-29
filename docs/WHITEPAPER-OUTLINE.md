# Archivas Whitepaper - Draft Outline

## Title
**Archivas: A Proof-of-Space-and-Time Layer 1 Blockchain**

## Abstract (300 words)

Archivas is a Layer 1 blockchain that achieves consensus through Proof-of-Space-and-Time, eliminating the energy waste of Proof-of-Work and the capital centralization of Proof-of-Stake. The protocol combines disk-based storage proofs with verifiable sequential computation to create a permissionless, energy-efficient, and fair consensus mechanism.

Key contributions:
- Modular PoSpace+VDF architecture in Go
- Adaptive difficulty targeting consistent block times
- Persistent state with crash recovery
- Production-tested farming implementation

We demonstrate that disk space combined with sequential time proofs can secure a blockchain without energy waste or capital barriers, making consensus participation accessible to anyone with commodity storage.

---

## 1. Introduction

### 1.1 The Blockchain Consensus Problem
- Security vs. decentralization vs. scalability
- Limitations of Proof-of-Work (energy waste)
- Limitations of Proof-of-Stake (capital requirements)
- Need for alternative consensus mechanisms

### 1.2 Storage-Based Consensus
- Proof-of-Space overview
- Chia Network's pioneering work
- Filecoin's storage proofs
- Gap: Need modular, auditable implementation

### 1.3 Archivas Vision
- Permissionless farming with commodity hardware
- Energy-efficient consensus
- Fair participation (disk space, not capital)
- Foundation for storage-utility layer

---

## 2. Background & Related Work

### 2.1 Proof-of-Work
- Bitcoin's consensus mechanism
- Energy consumption (~150 TWh/year)
- ASIC centralization
- Security model and assumptions

### 2.2 Proof-of-Stake
- Ethereum's transition
- Capital requirements
- "Nothing-at-stake" problem
- Validator economics

### 2.3 Proof-of-Space
- Original PoSpace papers (Dziembowski et al.)
- Chia Network implementation
- Spacemint protocol
- Burstcoin/Signum

### 2.4 Verifiable Delay Functions
- Wesolowski VDF (RSA groups)
- Pietrzak VDF (class groups)
- Applications in consensus
- Timelord architecture

---

## 3. Archivas Protocol

### 3.1 Overview
- Three-component architecture (node, timelord, farmer)
- Block structure
- Transaction model
- State management

### 3.2 Proof-of-Space

**3.2.1 Plot Generation**
- Plot file format
- Hash table construction: `H(farmerPubKey || plotID || index)`
- Storage requirements (k=16 to k=32)
- Plot reusability

**3.2.2 Challenge-Response**
- Challenge generation from VDF output
- Quality calculation: `Q = H(challenge || plotHash)`
- Proof structure
- Verification algorithm

**3.2.3 Difficulty Adjustment**
- Target block time: 20 seconds
- Moving average window: 10 blocks
- Adjustment bounds: 0.5x to 2.x
- Difficulty range: 2^40 to 2^60

### 3.3 Verifiable Delay Functions

**3.3.1 VDF Algorithm (Devnet)**
- Iterated SHA-256: `y_n = H(y_{n-1})`
- Sequential property
- Verification through recomputation
- Checkpoint system

**3.3.2 VDF Algorithm (Mainnet)**
- Wesolowski construction
- RSA group operations
- Proof generation and verification
- Performance characteristics

**3.3.3 Timelord Process**
- Continuous VDF computation
- Seed derivation: `H(blockHash || height)`
- Reset behavior on new blocks
- Network synchronization

### 3.4 Combined PoSpace+Time

**3.4.1 Block Production**
1. Timelord computes VDF from chain tip
2. Node generates challenge from VDF output
3. Farmers search plots for best proof
4. Farmer submits block (PoSpace + VDF)
5. Node validates both proofs
6. Block accepted, farmer rewarded

**3.4.2 Security Properties**
- Grinding resistance
- Temporal ordering
- No precomputation attacks
- Fairness guarantees

**3.4.3 Attack Scenarios**
- Long-range attacks (prevented by VDF)
- Nothing-at-stake (prevented by PoSpace)
- Selfish mining analysis
- 51% attack requirements

---

## 4. Implementation

### 4.1 Architecture

**4.1.1 Package Structure**
- config: Chain parameters
- ledger: State transitions, transactions
- wallet: Cryptography, signing
- pospace: Plot generation, proof verification
- vdf: Sequential computation
- consensus: Difficulty, challenge generation
- storage: Persistent state (BadgerDB)
- rpc: HTTP API
- p2p: Network layer (future)

**4.1.2 Data Structures**
- Block format
- Transaction structure
- Account state
- Plot file format
- VDF state

### 4.2 Cryptography

**4.2.1 Addresses**
- secp256k1 keypairs
- Bech32 encoding (arcv prefix)
- Address derivation: `H(pubKey)`

**4.2.2 Transaction Signing**
- ECDSA signatures
- Transaction hash construction
- Signature verification
- Replay protection (nonces)

### 4.3 Consensus

**4.3.1 Block Validation**
```
ValidateBlock(block):
  1. Verify VDF(seed, iterations) = output
  2. Verify challenge = H(VDFOutput || height)
  3. Verify PoSpace proof quality < difficulty
  4. Verify all transactions
  5. Apply state transitions
  6. Mint block reward to farmer
```

**4.3.2 Fork Choice**
- Longest chain rule
- Chain weight by cumulative difficulty
- Reorganization handling

### 4.4 Storage

**4.4.1 Database Schema**
- Block storage: `blk:<height> → block_data`
- Account storage: `acc:<address> → {balance, nonce}`
- Metadata: tip_height, difficulty, vdf_state

**4.4.2 Crash Recovery**
- Persistence after each block
- Atomic updates
- Recovery procedure on restart

---

## 5. Economics

### 5.1 Token (RCHV)

**5.1.1 Supply**
- Block reward: 20 RCHV
- Halving schedule: Every 1,051,200 blocks (~243 days)
- Terminal supply: ~42M RCHV (TBD - Bitcoin-like)
- Genesis allocation: 0% premine (or 10% foundation)

**5.1.2 Transaction Fees**
- Fee market (base fee + priority)
- Fee burning vs. farmer rewards (TBD)

### 5.2 Farming Economics

**5.2.1 Expected Returns**
- ROI calculation based on plot size
- Electricity costs (minimal - disk only)
- Hardware requirements
- Comparison to PoW mining

**5.2.2 Network Decentralization**
- Gini coefficient target
- Plot distribution analysis
- Barriers to entry
- Pool protocols (future)

---

## 6. Security Analysis

### 6.1 Threat Model
- Malicious farmers
- Timelord attacks
- Network-level attacks
- State bloat

### 6.2 Formal Security
- Safety proofs
- Liveness guarantees
- Byzantine fault tolerance
- Economic security bounds

### 6.3 Comparison

| Chain | Energy | Barrier | Decentralization | Security |
|-------|--------|---------|------------------|----------|
| Bitcoin (PoW) | Very High | Medium (ASICs) | Medium | Very High |
| Ethereum (PoS) | Low | High (32 ETH) | Medium | High |
| Chia (PoSpace+Time) | Very Low | Low (disk) | High | High |
| **Archivas** | **Very Low** | **Low** | **High** | **High** |

---

## 7. Performance

### 7.1 Benchmarks
- Plot generation speed
- Block validation time
- VDF computation rate
- Database I/O performance
- Sync speed

### 7.2 Scalability
- Transaction throughput (current)
- State size growth
- Pruning strategies
- Future improvements (sharding, rollups)

---

## 8. Future Work

### 8.1 Short-term
- P2P networking
- Multi-node testnet
- Block explorer
- Farming pools

### 8.2 Medium-term
- Wesolowski/Pietrzak VDF
- State snapshots
- Light clients
- Cross-chain bridges

### 8.3 Long-term
- Smart contracts (WASM)
- Useful storage (archival proofs)
- ZK-rollups
- Modular design (execution layer separation)

---

## 9. Conclusion

Archivas demonstrates that Proof-of-Space-and-Time consensus is not only viable but practical. By combining disk space lottery with sequential time proofs, we achieve:

- ✅ Permissionless participation
- ✅ Energy efficiency
- ✅ Temporal security
- ✅ Fair distribution

The devnet is operational, farming is verified, and persistence is tested. Archivas is ready to evolve from single-node devnet to multi-node testnet to public mainnet.

We invite developers, farmers, and researchers to join us in building the future of storage-based consensus.

---

## Appendices

### A. Protocol Parameters
- Complete parameter listing
- Rationale for each choice
- Devnet vs. testnet vs. mainnet differences

### B. API Reference
- RPC endpoints
- Request/response formats
- Error codes

### C. Farming Guide
- Hardware recommendations
- Plot size selection
- Expected earnings calculator
- Troubleshooting

### D. Developer Guide
- Package documentation
- Contributing guidelines
- Code structure
- Testing procedures

---

## References

1. Nakamoto, S. (2008). Bitcoin: A Peer-to-Peer Electronic Cash System
2. Buterin, V. et al. (2014). Ethereum: A Next-Generation Smart Contract Platform
3. Dziembowski, S. et al. (2015). Proofs of Space
4. Cohen, B. (2019). Chia Network Greenpaper
5. Wesolowski, B. (2018). Efficient Verifiable Delay Functions
6. Protocol Labs (2017). Filecoin: A Decentralized Storage Network

---

**Version:** 0.1 (Draft)  
**Date:** October 2025  
**Authors:** Archivas Labs  
**License:** Apache 2.0


