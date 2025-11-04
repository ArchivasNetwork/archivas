# Archivas Public Testnet Launch

**Date:** November 2025  
**Version:** v1.2.0  
**Network:** archivas-devnet-v4  
**Status:** 🟢 LIVE - 25+ days uptime, 64,000+ blocks  

---

## Twitter/X Thread

### Tweet 1 (The Hook)
```
🌾 Archivas Public Testnet is LIVE

A production-ready Proof-of-Space-and-Time blockchain.

✅ 64,000+ blocks mined over 25 days
✅ Public RPC: https://seed.archivas.ai
✅ Live explorer with real transactions
✅ TypeScript SDK ready

Farm RCHV with disk space. No energy waste, no staking.

🧵👇
```

### Tweet 2 (The Problem)
```
Traditional blockchains have critical flaws:

❌ PoW: Burns energy worth $10B+/year (Bitcoin)
❌ PoS: "Rich get richer" - high capital barriers

We need consensus that's:
✅ Permissionless (anyone can join)
✅ Energy efficient
✅ Fair (not capital-based)
```

### Tweet 3 (The Solution)
```
Archivas uses Proof-of-Space-and-Time:

🌾 SPACE: Farmers allocate disk storage, create "plots"
⏰ TIME: Timelords compute VDFs (sequential time proofs)

Blocks require BOTH:
• Winning PoSpace proof (disk lottery)
• Valid VDF output (prevents grinding)

This is Chia-class security.
```

### Tweet 4 (How It Works)
```
How Archivas mining works:

1️⃣ Create plots (precomputed hashes on disk)
2️⃣ Run farmer software (scans plots for proofs)
3️⃣ Find winning proof → submit block
4️⃣ Earn 20 RCHV block reward

More disk space = more lottery tickets.

But you can't pre-compute because VDF resets each block.
```

### Tweet 5 (Technical Achievement)
```
What's running in production:

✅ 64,000+ blocks over 25 days (99.9% uptime)
✅ 2 active farmers across 2 servers
✅ 1,260,000+ RCHV minted
✅ Public HTTPS API with TLS
✅ TypeScript SDK + Next.js explorer
✅ Ed25519 wallets with BIP39 mnemonics
✅ 1,160+ RCHV transferred in 8 transactions

Real chain. Real farming. Real transfers.
```

### Tweet 6 (Live Network Stats)
```
Current network stats:

📊 Height: 64,000+ blocks
💰 Supply: ~1,280,000 RCHV
⚡ Block time: ~25 seconds (stable)
🌾 Farmers: 2 servers, 7 k28 plots (~55 GB)
⏱️ Uptime: 25+ days continuous
🔄 Difficulty: Adaptive retargeting working

Public endpoint: https://seed.archivas.ai
Live explorer: archivas-explorer.up.railway.app
```

### Tweet 7 (Developer Tools)
```
Developer ecosystem ready:

🟢 Public RPC API (v1.2.0) - HTTPS, CORS, rate limited
🟢 TypeScript SDK - BIP39, Ed25519, transaction signing
🟢 Block Explorer - Next.js, real-time updates
🟢 CLI Tools - Wallet management, transfers
🟢 Prometheus + Grafana - Full observability

Build on Archivas:
📖 API docs
🔧 SDK: @archivas/sdk (ready for npm)
🌐 Explorer: https://archivas-explorer.up.railway.app
```

### Tweet 8 (Getting Started - Developers)
```
Build on Archivas:

```typescript
// TypeScript SDK
import { Derivation, createRpcClient } from '@archivas/sdk';

const mnemonic = Derivation.mnemonicGenerate();
const kp = await Derivation.fromMnemonic(mnemonic);
const address = Derivation.toAddress(kp.publicKey);

const rpc = createRpcClient({ 
  baseUrl: 'https://seed.archivas.ai' 
});
const balance = await rpc.getBalance(address);
```

Docs: github.com/ArchivasNetwork/archivas
SDK: github.com/ArchivasNetwork/archivas-sdk
```

### Tweet 9 (Proven Technology)
```
Why Archivas is production-ready:

✅ 25+ days continuous uptime
✅ 64,000+ blocks without consensus issues
✅ Multi-server coordination working
✅ State sync, IBD, reorg handling tested
✅ Transactions confirmed and verified
✅ Security: TLS, rate limiting, CORS
✅ Monitoring: Prometheus + Grafana

Not a prototype. Not a demo.
This is a working PoST blockchain.
```

### Tweet 10 (The Call to Action)
```
Join Archivas Testnet:

🌍 Public RPC: https://seed.archivas.ai
🔍 Explorer: archivas-explorer.up.railway.app
📦 SDK: @archivas/sdk (TypeScript)
⚙️ Grafana: Real-time metrics
📖 Docs: Full API reference

Developers: Build wallets, explorers, dApps
Farmers: Run nodes, earn RCHV
Researchers: Study PoST consensus

🌾 github.com/ArchivasNetwork/archivas

#Archivas #RCHV #ProofOfSpace #Web3
```

---

## Reddit/Discord/Forum Post

### Title
**Archivas (RCHV): New Proof-of-Space-and-Time L1 - Devnet Live, Farming Verified**

### Body

Hey everyone,

I'm excited to share **Archivas**, a new Layer 1 blockchain I've been building.

**What is Archivas?**

Archivas is a Proof-of-Space-and-Time blockchain (similar to Chia Network's consensus model). Instead of wasting energy (PoW) or requiring capital lockup (PoS), Archivas secures its network with:

- **Disk space** (Proof-of-Space) - Anyone with storage can farm
- **Sequential time** (Verifiable Delay Functions) - Prevents grinding attacks

The native token is **RCHV**. Farmers earn 20 RCHV per block by allocating disk space and finding winning proofs.

**Current Status:**

- 🟢 **Devnet operational** - Running and tested
- 🟢 **PoSpace farming works** - Verified with real blocks (120 RCHV earned in testing)
- 🟢 **Persistent storage** - Node survives restarts
- 🟢 **Adaptive difficulty** - Auto-adjusts for ~20s block times
- 🟡 **VDF/Timelord** - Implemented, ready to activate
- 🔴 **P2P networking** - Next milestone

**Technical Details:**

- Written in Go (~4,000 lines)
- Consensus: Proof-of-Space-and-Time
- Addresses: Bech32 (arcv prefix)
- Cryptography: secp256k1 ECDSA
- Storage: BadgerDB
- Block time: ~20 seconds
- Block reward: 20 RCHV

**How to Farm:**

```bash
# Clone and build
git clone https://github.com/ArchivasNetwork/archivas
cd archivas
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer

# Create a plot
./archivas-farmer plot --size 20 --path ./plots

# Start farming
./archivas-node
./archivas-farmer farm --plots ./plots --farmer-key <your_key>
```

**What's Tested:**
- ✅ Wallet generation and transaction signing
- ✅ Plot creation (k=16 in 156ms, k=20 in ~2s)
- ✅ Block farming (6 blocks in 60s test)
- ✅ Block rewards (120 RCHV verified on-chain)
- ✅ Node persistence (restart recovery tested)

**Roadmap:**

- **Now:** Single-node devnet with PoSpace
- **Q1 2026:** Multi-node testnet with P2P
- **Q2 2026:** Public testnet, block explorer
- **Q3 2026:** Security audit, mainnet prep

**Why I Built This:**

I wanted to prove you can build production-grade consensus from scratch. Archivas demonstrates that Proof-of-Space-and-Time is viable, efficient, and fair.

**Repository:** https://github.com/ArchivasNetwork/archivas

Looking for:
- Early farmers to test the network
- Go developers for P2P networking
- Feedback on consensus design
- Node operators for testnet

Questions welcome!

---

## HackerNews Post

### Title
**Archivas: A Proof-of-Space-and-Time blockchain built from scratch in Go**

### Body

I built Archivas over the past week - a Layer 1 blockchain using Proof-of-Space-and-Time consensus (similar to Chia Network's model).

**The idea:** Secure a blockchain with disk space and sequential time proofs instead of energy waste (PoW) or capital requirements (PoS).

**What works:**
- Farmers create "plots" (precomputed hash tables on disk)
- Farmers scan plots for winning proofs when blocks are needed
- Timelords compute VDFs (Verifiable Delay Functions) to prevent grinding
- Blocks require BOTH a valid PoSpace proof AND a valid VDF
- Farmers earn 20 RCHV per block mined

**Current status:**
- ~4,000 lines of Go
- Complete implementation (wallet, farmer, timelord, node)
- Tested: farmed 6 blocks, earned 120 RCHV
- Persistent storage with BadgerDB
- Adaptive difficulty working
- Single-node devnet operational

**Test results:**
- Plot generation: 156ms for 65K hashes (k=16)
- Farming: Found 6 winning proofs in 60 seconds
- Rewards: 120 RCHV earned and verified on-chain
- Persistence: Node restarted, full state recovered

**Tech stack:**
- Language: Go
- Consensus: PoSpace (disk space) + VDF (sequential time)
- Crypto: secp256k1, bech32 addresses
- Storage: BadgerDB
- Networking: HTTP RPC (P2P next)

**Next steps:**
- Milestone 6: libp2p for multi-node consensus
- Milestone 7: Public testnet
- Milestone 8: Security audit, mainnet

**Why this matters:**
- Proof-of-Space is energy-efficient and permissionless
- VDFs add temporal security (can't grind or precompute)
- Demonstrates this consensus model is viable at the L1 layer

**Code:** https://github.com/ArchivasNetwork/archivas

Open to feedback on the consensus design, implementation, or roadmap. Happy to answer technical questions!

---

## Medium/Substack Article

### Title
**Building Archivas: A Proof-of-Space-and-Time Blockchain from Scratch**

### Subtitle
*How we went from zero to farming blocks in one development sprint*

### Outline

**Introduction**
- The blockchain trilemma
- Problems with PoW and PoS
- Why Proof-of-Space-and-Time?

**Part 1: The Vision**
- Chia Network showed PoSpace works
- Filecoin showed storage has value  
- Gap: Need modular PoSpace+Time in Go
- Goal: Permissionless storage-based consensus

**Part 2: The Implementation**
- Milestone 1: Foundation (Go monorepo)
- Milestone 2: Cryptographic Ownership (wallets, signatures)
- Milestone 3: Proof-of-Space (farming tested!)
- Milestone 4: VDF/Timelord (temporal security)
- Milestone 5: Persistence (production-ready)

**Part 3: How It Works**
- Plot generation (precomputed hashes)
- Challenge-response (quality lottery)
- VDF computation (sequential time)
- Block validation (PoSpace + VDF)
- Adaptive difficulty
- Block rewards

**Part 4: Test Results**
- 6 blocks farmed in 60 seconds
- 120 RCHV earned
- Persistence verified (restart test)
- Difficulty adapted 5 times
- All consensus rules enforced

**Part 5: What's Next**
- P2P networking (multi-node testnet)
- Public testnet launch
- Block explorer
- Security audit
- Mainnet preparation

**Conclusion**
- Archivas proves PoSpace+Time is viable
- Production architecture from day 1
- Open source, community-driven
- Join us in building the future of storage consensus

---

## GitHub Issues to Create

### Development Roadmap

**Milestone 6: P2P Networking**
- [ ] Implement libp2p peer discovery
- [ ] Add block propagation protocol
- [ ] Implement sync from peers
- [ ] Add CLI flags for peer connections
- [ ] Test 3-node network

**Milestone 7: Public Testnet**
- [ ] Deploy to 3 public VPS nodes
- [ ] Activate VDF mode
- [ ] Create faucet for test RCHV
- [ ] Build simple block explorer
- [ ] Write testnet participation guide

**Milestone 8: Production Hardening**
- [ ] Upgrade to Wesolowski/Pietrzak VDF
- [ ] Add state pruning
- [ ] Implement checkpoints
- [ ] Security audit (external)
- [ ] Performance optimization
- [ ] Mainnet deployment guide

---

Let me know which you'd like first and I'll generate it!


