# Archivas Roadmap

## Vision

**Build a permissionless data layer secured by commodity storage where disk space provides both consensus and utility.**

---

## Phase 1: Testnet Alpha (Current) ✅

**Goal:** Prove multi-node PoSpace+Time consensus works

**Status:** COMPLETE

**Achievements:**
- [x] Core blockchain implementation
- [x] Proof-of-Space farming
- [x] VDF/Timelord
- [x] Multi-node sync
- [x] Deterministic genesis
- [x] P2P networking
- [x] Live VPS deployment

**Metrics:**
- Nodes: 2+
- Height: 78+ blocks
- RCHV: ~1,560 farmed
- Status: 🟢 LIVE

---

## Phase 2: Testnet Beta (Q4 2025)

**Goal:** Grow network, improve UX, build tooling

### Infrastructure
- [ ] **10+ nodes** distributed globally
- [ ] **Persistent peer store** (auto-reconnect)
- [ ] **DNS seeding** (seed.archivas.network)
- [ ] **Load balancer** for RPC endpoints
- [ ] **Monitoring dashboard** (Grafana)

### Tooling
- [ ] **Block explorer** (web UI)
- [ ] **Faucet** for test RCHV
- [ ] **Wallet GUI** (Electron or web)
- [ ] **Pool protocol** for group farming
- [ ] **Light client** (mobile-friendly)

### Documentation
- [ ] **GitBook site** (docs.archivas.network)
- [ ] **API documentation** (OpenAPI/Swagger)
- [ ] **Video tutorials**
- [ ] **Farming calculator**
- [ ] **Translated docs** (ES, ZH, etc.)

### Community
- [ ] **Discord server**
- [ ] **Telegram group**
- [ ] **Twitter/X presence**
- [ ] **Blog/Medium** (technical updates)
- [ ] **Bounty program**

**Target:** 50+ community farmers, 10+ nodes

---

## Phase 3: Testnet RC (Q1 2026)

**Goal:** Production readiness, security hardening

### Security
- [ ] **External audit** (Trail of Bits, etc.)
- [ ] **Bug bounty program** (up to $50K)
- [ ] **Formal verification** of consensus
- [ ] **Penetration testing**
- [ ] **Economic analysis**

### Performance
- [ ] **State pruning** (archive vs. full nodes)
- [ ] **Block compression**
- [ ] **Parallel transaction execution**
- [ ] **Optimized plot scanning**
- [ ] **VDF acceleration** (GPU/ASIC)

### VDF Upgrade
- [ ] **Wesolowski VDF** implementation
- [ ] **Succinct proofs** (O(log n) verification)
- [ ] **Multi-timelord** competition
- [ ] **VDF fallback** mechanisms

### Network
- [ ] **NAT traversal** (STUN/TURN)
- [ ] **Encrypted connections** (TLS)
- [ ] **Peer reputation** system
- [ ] **DDoS protection**
- [ ] **Rate limiting**

**Target:** Production-grade testnet, 100+ nodes

---

## Phase 4: Mainnet Prep (Q2 2026)

**Goal:** Economic model, token distribution, final testing

### Economics
- [ ] **Halving schedule** finalized
- [ ] **Total supply** determined
- [ ] **Fee market** design
- [ ] **Foundation allocation**
- [ ] **Inflation model**

### Distribution
- [ ] **Airdrop** to testnet participants
- [ ] **Farming rewards** allocation
- [ ] **Team/advisors** vesting
- [ ] **Treasury** multisig
- [ ] **DEX liquidity** planning

### Final Testing
- [ ] **Stress test** (1000+ nodes simulated)
- [ ] **Network partition** recovery
- [ ] **Long-range attack** testing
- [ ] **Economic attack** simulations
- [ ] **Disaster recovery** procedures

### Compliance
- [ ] **Legal review**
- [ ] **Token classification**
- [ ] **Regulatory compliance** (where applicable)
- [ ] **Terms of service**
- [ ] **Privacy policy**

**Target:** Mainnet-ready, audited, compliant

---

## Phase 5: Mainnet Launch (Q3 2026)

**Goal:** Public launch of Archivas L1

### Launch
- [ ] **Mainnet genesis** (public ceremony)
- [ ] **Network launch** (coordinated start)
- [ ] **Token distribution** (airdrops, etc.)
- [ ] **Exchange listings** (DEX first, CEX later)
- [ ] **Marketing campaign**

### Infrastructure
- [ ] **Public RPC** endpoints
- [ ] **Block explorers** (multiple)
- [ ] **Faucets** (for onboarding)
- [ ] **Mobile wallets**
- [ ] **Hardware wallet** support

### Ecosystem
- [ ] **Developer grants** program
- [ ] **Ecosystem fund** (10M RCHV example)
- [ ] **Accelerator** for dApps
- [ ] **Partnerships** (storage providers, etc.)
- [ ] **Community governance**

**Target:** 1000+ nodes, global distribution

---

## Phase 6: Ecosystem Growth (2026+)

**Goal:** Build utility and adoption

### Smart Contracts
- [ ] **WASM runtime**
- [ ] **Contract deployment** API
- [ ] **Gas model**
- [ ] **Contract verification**
- [ ] **SDK** (TypeScript, Rust)

### Storage Utility
- [ ] **Proof-of-Archival** (useful storage)
- [ ] **IPFS integration**
- [ ] **Filecoin bridge**
- [ ] **Data marketplace**
- [ ] **Decentralized CDN**

### DeFi
- [ ] **DEX** (native RCHV pairs)
- [ ] **Lending** protocol
- [ ] **Stablecoins**
- [ ] **Derivatives**
- [ ] **Bridges** (to ETH, BTC, etc.)

### Layer 2
- [ ] **Optimistic rollups**
- [ ] **ZK rollups**
- [ ] **State channels**
- [ ] **Sidechains**
- [ ] **Plasma**

**Target:** Thriving ecosystem, 10K+ users

---

## Technical Milestones

### Consensus Evolution

**Phase 1:** PoSpace only (devnet v1)  
**Phase 2:** PoSpace + SHA-256 VDF (devnet v2-v3)  
**Phase 3:** PoSpace + Wesolowski VDF (testnet RC)  
**Phase 4:** PoSpace + VDF + Finality gadget (mainnet)

### Networking Evolution

**Phase 1:** Manual peers (current)  
**Phase 2:** Bootnode discovery (current)  
**Phase 3:** DHT peer discovery  
**Phase 4:** libp2p integration  
**Phase 5:** Kademlia routing

### Storage Evolution

**Phase 1:** In-memory + BadgerDB (current)  
**Phase 2:** State pruning  
**Phase 3:** Snapshots + fast sync  
**Phase 4:** Distributed storage  
**Phase 5:** Archival nodes + light clients

---

## Research Directions

### Future Research

**VDF:**
- Class group VDFs
- Verifiable computation
- Hardware acceleration
- Proof aggregation

**PoSpace:**
- Beyond Hellman plots
- Compressed proofs
- Useful storage integration
- Space-time tradeoffs

**Consensus:**
- Byzantine fault tolerance proofs
- Finality gadgets
- Cross-chain security
- Economic security bounds

**Scalability:**
- Sharding designs
- Data availability sampling
- Fraud proofs
- Validity proofs (ZK)

---

## Community Roadmap

**User input drives priorities!**

**Vote on features:**
- GitHub Discussions
- Community calls
- Governance proposals (future)

**Current asks:**
1. Block explorer (high priority)
2. Mobile wallet (high demand)
3. Farming pools (community request)
4. DEX integration (future)

---

## Long-term Vision

### Year 1 (2026)
- Mainnet launch
- 1K+ nodes
- 10K+ users
- Basic ecosystem

### Year 2 (2027)
- Smart contracts
- DeFi applications
- 100K+ users
- Exchange listings

### Year 3 (2028)
- Storage utility live
- L2 scaling solutions
- 1M+ users
- Decentralized governance

### Year 5 (2030)
- Archivas as infrastructure layer
- Hundreds of dApps
- Millions of users
- Sustainable ecosystem

**Mission:** Make storage-based consensus the standard for permissionless networks.

---

## How to Influence

**Developers:**
- Build on Archivas
- Contribute code
- Propose features
- Review PRs

**Farmers:**
- Run nodes
- Provide feedback
- Test features
- Report issues

**Community:**
- Spread awareness
- Create content
- Answer questions
- Organize events

**Everyone:**
- Join discussions
- Vote on proposals
- Participate in governance
- Build the future! 🌾

---

**Questions?** [Join the discussion!](https://github.com/ArchivasNetwork/archivas/discussions)  
**Want to help?** [See Developer Docs](Developer-Docs.md)  
**Back:** [← Developer Docs](Developer-Docs.md)

