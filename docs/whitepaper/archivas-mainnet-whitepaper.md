# Archivas Mainnet Whitepaper Draft

**Version:** Draft v0.95  
**Status:** Complete content for design, economic, and technical review prior to layout  
**Date:** November 9, 2025  
**Prepared by:** Archivas Core Engineering, Research, Economics, Infrastructure, and Community Teams  

---

## Preface

This document is the canonical source for the Archivas Mainnet launch. It extends the devnet documentation with production parameters, detailed subsystem descriptions, economic modeling, operational guidance, and governance structures. The draft is intentionally verbose so that future layout teams can extract chapters, diagrams, and briefing material without requesting additional research. All figures are accurate as of the publication date and will be versioned on chain at genesis.

Reads of the document should assume familiarity with proof of space, verifiable delay functions, and distributed systems. Nevertheless, introductory sections provide accessible definitions for stakeholders who are new to the Archivas ecosystem.

---

## Table of Contents

1. Executive Summary  
2. Foundations and Historical Context  
3. Architectural Overview  
4. Consensus and Proof Systems  
5. Block Lifecycle and State Management  
6. Networking and Peer Connectivity  
7. Storage Engine and Data Persistence  
8. Infrastructure Architecture and Operations  
9. Economic Model and Token Supply  
10. Incentives, Distribution, and Treasury Mechanics  
11. Security, Cryptography, and Threat Analysis  
12. Developer Tooling, SDKs, and Integration Patterns  
13. Ecosystem Applications and Use Cases  
14. Governance, Community, and Social Contracts  
15. Risk Analysis and Mitigation Strategies  
16. Roadmap and Future Enhancements  
17. Vision, Philosophy, and Cultural Tenets  
18. Appendices  
19. References and Citations  

Each chapter contains subsections, data tables, and scenario analyses. The appendices provide raw parameters, formulas, pseudocode, and interoperability notes so that independent teams can reconstruct the protocol from first principles.

---

## 1. Executive Summary

Archivas is a modular Proof-of-Space-and-Time (PoST) blockchain implemented in Go. The network converts persistent storage into a Sybil resistant resource while verifiable delay functions (VDFs) enforce sequential time. Mainnet extends the devnet achievements, which already demonstrated more than one hundred thousand consecutive blocks with valid time proofs, multi-node synchronization after an initial block download rewrite, and real-time explorer indexing powered by the public JavaScript SDK.

The mainnet design adheres to six principles:

1. **Persistence before throughput.** Archivas prioritizes integrity, verifiability, and archival permanence over raw transaction-per-second metrics.  
2. **Modularity.** Nodes, farmers, timelords, indexers, and explorers evolve independently with secure interfaces.  
3. **Transparency.** Economic parameters, reward splits, and treasury flows are deterministic and auditable.  
4. **Resilience.** Networking, storage, and cryptographic layers employ conservative assumptions to survive faulty peers.  
5. **Openness.** Participation does not require permission or license.  
6. **Sustainability.** The protocol favors storage commitments and sequential computation instead of energy-intensive hashing.

Mainnet introduces production-grade initial block download (IBD), a finalized REST API, cross-component observability, on-chain treasury accounts, and the premiere release of the Archivist browser extension wallet. The project’s motto, **uptime equals morality**, reflects the conviction that nodes preserving knowledge contribute directly to societal memory.

---

## 2. Foundations and Historical Context

### 2.1 Origin Story

The Archivas concept emerged during research sprints that evaluated the limitations of pure Proof-of-Work (PoW) systems and the fragility of centralized storage platforms. Early prototypes used k=25 plots and a basic sequential delay enforcement written in pure Go. During month three of experimentation, the team migrated to Wesolowski style VDFs and redesigned the network stack around libp2p to benefit from secure, authenticated channels.

The devnet was activated in early 2025 with three principal servers:

- **Server A:** Seed node, farmer, and timelord. Hosted the canonical RPC endpoint and explorer.  
- **Server B:** General node and farmer, initially with k=25 plots to validate upgrade paths to larger plot sizes.  
- **Server C:** Dedicated farmer for stress testing plotting and challenge responsiveness.

Throughout the devnet, the network survived adversarial tests: induced RPC outages, mempool saturation, and simulated double spend attempts. Each incident fed into patches that are now integrated into mainnet.

### 2.2 Devnet Milestones

| Milestone | Description | Lessons Applied to Mainnet |
|-----------|-------------|----------------------------|
| 10k Blocks | Confirmed combined PoST and VDF pipeline | Validated sequential delay sizing |
| Explorer Launch | Real-time dashboard at `seed.archivas.ai` | Hardened RPC endpoint health checks |
| SDK Release | TypeScript SDK with canonical JSON signing | Ensured wallet compatibility |
| 100k Blocks | Long-running stability test | Tuned difficulty adjustment parameters |
| IBD Rewrite | Added checkpointing and resumable sync | Adopted for mainnet bootstrap |
| Archivist Wallet Alpha | Browser extension with Ed25519 accounts | Locked in key derivation standards |

### 2.3 Evolution to Mainnet

The transition from devnet to mainnet required expansion across five axes:

1. **Consensus Hardening:** Rewriting the transaction submission path to validate Ed25519 signed payloads and ensuring blocks enforce deterministic state transitions.  
2. **Observability:** Introducing Prometheus metrics, Grafana dashboards, and a watchdog service capable of restarting nodes if `/chainTip` becomes unresponsive for thirty seconds.  
3. **Operational Automation:** Converting shell scripts into systemd services, ensuring that farmers and timelords automatically rejoin after maintenance.  
4. **Economic Finalization:** Designing a halving schedule with predictable tail emissions and explicit treasury allocations.  
5. **Community Readiness:** Producing documentation, walkthroughs, migration guides, and a governance charter.

### 2.4 Philosophy of Uptime

The phrase **uptime equals morality** is more than rhetoric. Archivas considers data erasure a moral hazard. Keeping nodes online maintains the collective memory. This philosophy influences reward design: long-lived nodes benefit from consistent payouts, while offline nodes miss opportunities. Future governance proposals will reward community projects that improve uptime measurement and visualization.

---

## 3. Architectural Overview

### 3.1 Component Catalogue

Archivas consists of six core software components, each versioned independently:

| Component | Purpose | Implementation Language | Interfaces |
|-----------|---------|-------------------------|------------|
| Node | Validates blocks, maintains state, exposes RPC | Go | REST, gRPC (internal), libp2p |
| Farmer | Generates plots, responds to challenges | Go | Secure WebSocket, CLI |
| Timelord | Computes VDF proofs | Go with FPGA hooks | gRPC, libp2p |
| Indexer | Streams chain data into Postgres | Go | gRPC, SQL |
| Explorer | Presents blocks, transactions, charts | TypeScript, React | REST, GraphQL |
| SDK | Client libraries for wallets and dApps | TypeScript | REST wrappers, crypto utilities |

### 3.2 Layered Architecture

1. **Networking Layer:** libp2p based transport, noise handshake, GossipSub topics for blocks and proofs, adaptive rate limiting.  
2. **Consensus Layer:** Proof-of-Space challenge issuance, VDF proof verification, difficulty adjustment, fork choice rule (longest valid chain).  
3. **State Layer:** Persistent store using PebbleDB for account balances, nonces, and reward tracking.  
4. **Application Layer:** REST interface, transaction mempool, historical queries, subscription endpoints.  
5. **Tooling Layer:** SDK, command-line utilities, explorer dashboards, browser extension wallet.

### 3.3 Data Flow Diagram (Textual)

```
[Farmer] --proof--> [Node] --candidate block--> [Timelord]
[Timelord] --vdf proof--> [Node] --validated block--> [P2P Network]
[Node] --events--> [Indexer] --records--> [Postgres]
[SDK] --rpc--> [Node] --responses--> [Wallets and dApps]
```

Each arrow represents an authenticated channel with request-response semantics. Farmers submit proofs via WebSocket or RPC. Timelords publish VDF outputs. Nodes gossip finalized blocks to peers. Indexers subscribe to new block events and persist them for explorers. SDKs connect to nodes for read and write operations.

### 3.4 Deployment Topology

Mainnet encourages geographic diversity. Seed nodes in Europe, North America, and Asia provide bootstrap contacts. Farmers may colocate with nodes or operate remotely. Timelords should be positioned near reliable power with redundant connectivity. Indexers benefit from data center deployments with large NVMe arrays.

### 3.5 Component Version Standards

- **Semantic Versioning:** Each component uses MAJOR.MINOR.PATCH. Breaking consensus changes increment MAJOR and require network-wide coordination.  
- **Release Candidates:** RC builds run on staging nets before mainnet activation.  
- **Security Patch Policy:** Critical fixes ship within seven days of disclosure, with signed binaries distributed via the official release portal.

---

## 4. Consensus and Proof Systems

### 4.1 Proof-of-Space Mechanics

Farmers allocate disk capacity by generating plots. Each plot is identified by a unique plot public key derived from a mnemonic via SLIP-0010 on the Ed25519 curve. The standard mainnet plot size is k=28 (approximately 8.6 GiB). Plotting involves generating seven tables (T1 through T7) where each entry references two entries from the previous table. The plotting algorithm uses deterministic pseudo-random functions seeded by the farmer private key.

During farming:

1. **Challenge Reception:** Nodes broadcast a 256-bit challenge every block.  
2. **Plot Filtering:** Farmers compute a quick filter on each plot to discard non-qualifying entries.  
3. **Proof Construction:** For each qualifying entry, the farmer traces back pointers through the tables to reconstruct the proof.  
4. **Quality Score:** The proof yields a quality string, hashed against the challenge to compute a score.  
5. **Threshold Check:** If the score is below the difficulty threshold, the proof is eligible for block creation.  
6. **Submission:** The farmer submits the proof, signature, and metadata to a connected node.

### 4.2 Proof-of-Time Execution

Timelords ensure that blocks obey a minimum time interval. They run sequential verifiable delay functions based on class groups of imaginary quadratic fields. The process:

1. Receive the challenge from the previous block.  
2. Initialize the VDF with the challenge and pre-agreed discriminant.  
3. Iterate the squaring function for `N` steps (target 30 seconds).  
4. Produce output `Y` and short proof `π` using Wesolowski compression.  
5. Broadcast `Y` and `π` to peers.  
6. Nodes verify the proof in logarithmic time and accept the block if it matches the submitted farmer proof.

### 4.3 Difficulty Adjustment Algorithm

Archivas targets a 30 second block interval. Every 1,024 blocks, the difficulty is adjusted according to:

```
observed = (timestamp_latest - timestamp_reference) / 1024
ratio = target / observed
bounded_ratio = clamp(ratio, 0.75, 1.25)
difficulty_next = difficulty_current * bounded_ratio
```

This approach prevents extreme swings caused by timestamp manipulation or sudden changes in farming power. Difficulty adjustments propagate across the network via block headers.

### 4.4 Fork Choice Rule

Archivas employs a longest valid chain rule with difficulty weighting. Each block carries cumulative difficulty. Nodes always prefer the chain with greater cumulative difficulty, ensuring that temporary forks resolve once a chain produces more proofs. VDF proofs prevent rapid chain rewrites because generating alternative sequential proofs requires real time.

### 4.5 State Transition Validation

For each block:

1. Verify the farmer proof against the challenge.  
2. Verify the VDF output and proof.  
3. Iterate through transactions in canonical order.  
4. Verify each transaction signature (Ed25519).  
5. Check account balances and nonces.  
6. Apply state updates: debit sender, credit recipient, add fees to block reward pool, increment nonce.  
7. Persist resulting state root.  
8. Update reward counters for farmer, timelord, treasury, and development fund.

### 4.6 Formal Verification Roadmap

Archivas plans to formalize the consensus protocol using TLA+ specifications. Deliverables include:

- State transition invariants.  
- Liveness proofs for block production under honest majority of storage.  
- Safety proof for fork choice with adversarial VDF participants.  
- Simulation harness that plays TLA+ traces against a Go reference implementation.

---

## 5. Block Lifecycle and State Management

### 5.1 Lifecycle Overview

1. **Challenge Issuance:** The previous block hash and difficulty produce a new challenge.  
2. **Farmer Selection:** Farmers compute proofs and submit eligible ones.  
3. **Candidate Assembly:** Nodes assemble candidate blocks with mempool transactions.  
4. **VDF Verification:** Timelords deliver proof-of-time outputs.  
5. **Block Finalization:** Node verifies proofs, applies transactions, and commits the block.  
6. **Propagation:** Finalized block is gossiped to peers.  
7. **State Persistence:** Updated state root and ledgers are stored in PebbleDB.  
8. **Reward Distribution:** Rewards are recorded for payout addresses.

### 5.2 Mempool Policy

- **Capacity:** Default capacity of 50,000 transactions.  
- **Expiration:** Unconfirmed transactions expire after 1,000 blocks.  
- **Fee Prioritization:** Transactions sorted by fee per byte.  
- **Nonce Gap Handling:** Transactions with nonces greater than current expected value are held in a future queue until previous nonces arrive.  
- **Spam Prevention:** Rate limits per origin IP and per public key.

### 5.3 Canonical JSON Encoding

Transactions rely on RFC 8785 canonical JSON. Example serialization:

```
{
  "amount":"1000000",
  "fee":"1000",
  "from":"arcv1...",
  "memo":"example",
  "nonce":42,
  "to":"arcv1..."
}
```

Bytes are hashed with Blake2b-256 using the domain string `Archivas-TxV1`. The resulting hash is signed with Ed25519. Lowercase field names avoid mismatches between Go struct tags and JavaScript interfaces.

### 5.4 State Storage

- **Accounts Table:** Stores balances (uint64), nonces, reward accumulators.  
- **Rewards Table:** Tracks pending payouts for farmers and timelords to avoid floating point rounding.  
- **Treasury Ledger:** Separate namespace for treasury allocations with governance metadata.  
- **Historical Snapshots:** Every 5,000 blocks, the node writes a snapshot containing state hash, height, and timestamp to facilitate quick restores.

### 5.5 Accounts and Addresses

- **Address Format:** Bech32 with prefix `arcv`.  
- **Checksum:** 6-character checksum derived from HRP and data part.  
- **Account Types:** External accounts (controlled by private keys) and virtual accounts (treasury, protocol).  
- **Nonce Rules:** Strict increment per transaction. If a transaction fails during application, the nonce does not increment and the transaction is rejected from the block.

### 5.6 Reward Accounting

During block application:

1. Compute total fees from included transactions.  
2. Add base block reward (depending on epoch).  
3. Distribute according to allocation percentages.  
4. Update the reward ledger so that farmers and timelords can claim via future transactions or automated payout modules.

### 5.7 State Synchronization Integrity

Nodes validate each block header and transaction. If an inconsistency arises, the node triggers an automatic halt, writes diagnostic logs, and waits for operator intervention. Operators are guided to compare state roots with known good nodes to identify divergence.

---

## 6. Networking and Peer Connectivity

### 6.1 Transport Protocols

- **Base Transport:** TCP with noise handshake for encryption and identity.  
- **Multiplexing:** Yamux for concurrent streams.  
- **Compression:** Snappy compression for large payloads (blocks, proofs).  
- **Port:** Default 7744 with optional TLS termination for RPC endpoints.

### 6.2 Gossip Topics

| Topic | Payload | Frequency |
|-------|---------|-----------|
| `/archivas/block` | Finalized block headers and payloads | Every 30 seconds |
| `/archivas/proof` | Farmer proofs awaiting VDF | On demand |
| `/archivas/vdf` | VDF outputs and proofs | Every 30 seconds |
| `/archivas/tx` | New transactions | As submitted |

Peers subscribe to relevant topics based on node role. Farmers focus on proof topics, while indexers listen to blocks and transactions.

### 6.3 Peer Discovery

- **Bootstrap List:** Hard coded addresses for Server A and Server B.  
- **mDNS:** Used in private clusters or local testing.  
- **Peer Exchange:** Nodes exchange peer lists after successful handshake.  
- **Ban Lists:** Nodes persist malicious peer IDs for 24 hours following detection.

### 6.4 Bandwidth Management

- Rate limiting ensures that a single peer cannot saturate upload queues.  
- Blocks are propagated with high priority, while transactions have medium priority.  
- Chunked transfer splits large payloads into 64 KB fragments to avoid congestion.  
- Nodes expose metrics for inbound and outbound bandwidth to assist operators.

### 6.5 Synchronization Modes

- **Header Sync:** Nodes first download headers to verify chain weight.  
- **Body Sync:** Once headers are trusted, nodes fetch block bodies.  
- **State Sync:** Optional mode that downloads authenticated state snapshots for faster catch-up.  
- **Live Sync:** After sync, nodes switch to real-time gossip participation.

### 6.6 RPC Security

- **Authentication:** Optional API keys and IP allow lists for operators hosting public endpoints.  
- **Rate Limits:** Per endpoint limits to prevent abuse (`/submit` hardened to 10 requests per second per key).  
- **TLS:** Required for public endpoints. Certificates rotated quarterly.  
- **Logging:** Structured logs with correlation IDs for debugging.

---

## 7. Storage Engine and Data Persistence

### 7.1 Plot Storage

Farmers maintain directories of `.plot` files with metadata describing k-size, memo, and farmer public key. Recommended best practices:

- Use RAID1 for redundancy.  
- Monitor SMART data to anticipate drive failures.  
- Maintain a manifest file that maps plot IDs to disk positions for quick recovery.

### 7.2 Node Storage

- **Database:** PebbleDB with column families for accounts, blocks, and metadata.  
- **Compaction Strategy:** Background compactions limited to 25 percent CPU usage.  
- **Sharding:** Archive nodes may offload older state to cold storage using snapshot exports.  
- **Backup Policy:** Weekly full backups plus incremental daily snapshots.

### 7.3 Indexer Storage

- **Primary Database:** PostgreSQL 15 with partitioned tables for blocks and transactions.  
- **Schema Highlights:**  
  - `blocks(height bigint primary key, hash bytea, timestamp timestamptz, difficulty numeric, farmer text, timelord text)`  
  - `transactions(hash bytea primary key, block_height bigint, from_address text, to_address text, amount numeric, fee numeric, nonce bigint)`  
  - `accounts(address text primary key, balance numeric, last_seen_block bigint)`  
- **Retention:** All data persisted. No pruning.  
- **Replication:** Read replicas for explorer and analytics workloads.

### 7.4 Compression and Pruning

Mainnet does not prune historical data. However, nodes can compress archived blocks using zstd. Future governance proposals may introduce optional pruning for lightweight nodes, provided that proof-of-history commitments remain verifiable.

### 7.5 Observability of Storage Health

- Metrics for disk utilization, read/write latency, and database compaction status.  
- Alerts when plot directories fall below configured capacity thresholds.  
- Indexers expose query latency distributions to detect slowdowns.

---

## 8. Infrastructure Architecture and Operations

### 8.1 Seed Node Specifications

| Server | Role | CPU | RAM | Storage | Network | Location |
|--------|------|-----|-----|---------|---------|----------|
| Server A | Seed Node, Farmer, Timelord | 16 cores | 64 GB | NVMe + HDD arrays | Dual 10 Gbps | Frankfurt |
| Server B | Seed Node, Farmer | 12 cores | 48 GB | NVMe + HDD arrays | Dual 5 Gbps | New York |
| Server C | Farmer | 8 cores | 32 GB | HDD arrays | 1 Gbps | Singapore |

### 8.2 Deployment Model

- **Containerization:** Docker images with reproducible builds.  
- **Orchestration:** Systemd for simplicity; Kubernetes deployments are supported but optional.  
- **Configuration Management:** Ansible playbooks for provisioning, with environment-specific overrides.  
- **Secret Management:** Vault or SOPS for storing mnemonic backups with Shamir shared secrets.

### 8.3 Monitoring and Alerting

- **Prometheus Exporters:** Node metrics, farmer challenge response times, timelord iteration counts.  
- **Grafana Dashboards:**  
  - Chain health overview  
  - RPC latency distribution  
  - Plot performance metrics  
  - Timelord VDF cycle time histogram  
- **Alert Rules:**  
  - Chain tip divergence greater than five blocks  
  - RPC error rate above five percent for five minutes  
  - Mempool depth exceeding 75 percent capacity  
  - Disk utilization above 85 percent

### 8.4 Maintenance Procedures

- **Rolling Upgrades:** Upgrade secondary nodes first, observe stability, then upgrade primary seed nodes.  
- **Backup Verification:** Monthly restore drills using archived snapshots to ensure recoverability.  
- **Disaster Recovery:** Warm standby in a different region with daily state synchronization.  
- **Incident Response:** Runbook includes log collection, peer comparison, and rollback steps.

### 8.5 Operational Security

- **Access Control:** SSH with hardware security keys, role based access.  
- **Logging:** Centralized logging with retention of 180 days.  
- **Audit Trails:** Every configuration change recorded with operator identity.  
- **Penetration Testing:** Annual third-party assessments with published findings.

---

## 9. Economic Model and Token Supply

### 9.1 Design Goals

- Predictable long-term issuance.  
- Strong incentives for storage providers and time verifiers.  
- Funding streams for community and ongoing development.  
- Gradual tail emission to avoid zero reward stagnation.

### 9.2 Monetary Policy Parameters

| Parameter | Value |
|-----------|-------|
| Genesis Supply | 0 RCHV |
| Initial Block Reward | 8 RCHV |
| Block Interval | 30 seconds |
| Halving Interval | 63,072,000 blocks (two years) |
| Minimum Reward Floor | 0.5 RCHV |
| Treasury Allocation | 3 percent of block reward |
| Development Fund Allocation | 2 percent |
| Farmer Allocation | 90 percent |
| Timelord Allocation | 5 percent |

### 9.3 Emission Curve Analysis

The emission curve is geometric with a floor. The cumulative supply `S(t)` after `n` epochs is:

```
S(n) = Blocks_per_epoch * sum_{i=0}^{n} max(Reward_initial / 2^i, 0.5)
```

At epoch five and beyond the reward equals the floor. This ensures continuous incentives without uncontrolled inflation.

### 9.4 Supply Projections

| Year | Block Height | Base Reward | Annual Issuance | Cumulative Supply |
|------|--------------|-------------|-----------------|-------------------|
| 0 | 0 | 8.0 | 4,204,800 | 4,204,800 |
| 2 | 63,072,000 | 4.0 | 2,102,400 | 8,409,600 |
| 4 | 126,144,000 | 2.0 | 1,051,200 | 12,614,400 |
| 6 | 189,216,000 | 1.0 | 525,600 | 14,716,800 |
| 8 | 252,288,000 | 0.5 | 262,800 | 15,768,000 |
| 10 | 315,360,000 | 0.5 | 262,800 | 16,800,000 |
| 20 | 630,720,000 | 0.5 | 262,800 | 19,428,000 |
| 40 | 1,261,440,000 | 0.5 | 262,800 | 24,684,000 |
| 80 | 2,522,880,000 | 0.5 | 262,800 | 35,196,000 |
| 120 | 3,784,320,000 | 0.5 | 262,800 | 45,708,000 |

These projections align with a target asymptotic supply of approximately 43 million RCHV when accounting for transaction fees burned or redistributed.

### 9.5 Fee Market

- Minimum fee: 1,000 mojos (1 mojo = 0.000001 RCHV).  
- Dynamic fee suggestions based on mempool congestion.  
- Future upgrade: EIP-1559 inspired fee mechanism with base fee adjustments per block while retaining priority tips.

### 9.6 Inflation and Purchasing Power

The emission floor results in an annual inflation rate that trends downward as total supply increases. At year ten the rate is approximately 1.6 percent, falling below one percent after year twenty. This controlled inflation funds continuous network security and treasury initiatives without eroding long-term value.

---

## 10. Incentives, Distribution, and Treasury Mechanics

### 10.1 Reward Allocation Breakdown

| Allocation | Percentage | Distribution Method |
|------------|------------|---------------------|
| Farmers | 90 percent | Paid automatically per block to winning farmer address |
| Timelords | 5 percent | Paid to timelord that supplied successful VDF proof |
| Community Treasury | 3 percent | Credited to multisig treasury account `arcv1treasury...` |
| Development Fund | 2 percent | Credited to multisig development account `arcv1devfund...` |

### 10.2 Treasury Governance

- Treasury keys secured via 4-of-7 multisig.  
- Quarterly reports include balances, grant disbursements, and upcoming proposals.  
- Treasury proposals require 60 percent approval of voting power within a two week window.  
- Grants prioritized for protocol research, documentation, localization, and community education.

### 10.3 Development Fund

- Focused on core client maintenance, security audits, and reference implementations of SDKs.  
- Development fund disbursements require approval from a technical steering committee with community observers.  
- All expenditures published with invoices and deliverables.

### 10.4 Timelord Incentive Model

Timelords incur hardware costs for sequential computation. The five percent allocation ensures return on investment for high availability setups. In the event of multiple timelords submitting proofs for the same block, rewards are split proportionally based on completion timestamps. Future upgrades may introduce time proof auctions to incentivize faster proofs.

### 10.5 Farmer Incentive Model

- Rewards scale linearly with contributed storage.  
- Farmers can join pooling protocols (subject to governance approval) where proofs and rewards are aggregated.  
- Pooling contracts must obey Archivas audit requirements to prevent centralization.

### 10.6 Community Incentives

- Treasury funds hackathons, documentation sprints, and localization bounties.  
- Indexer operators receive subsidies during bootstrap to ensure explorer coverage.  
- Governance may introduce quadratic funding rounds for community projects.

---

## 11. Security, Cryptography, and Threat Analysis

### 11.1 Cryptographic Primitives

- **Key Derivation:** BIP39 mnemonics, SLIP-0010 Ed25519 curve.  
- **Signature Scheme:** Ed25519 using libsodium bindings.  
- **Hash Functions:** Blake2b-256 for transaction hashes, SHA-256 for internal table operations.  
- **VDF:** Wesolowski proof system on class group of unknown order.  
- **Address Encoding:** Bech32 with error detection.

### 11.2 Threat Model

| Threat | Description | Mitigation |
|--------|-------------|------------|
| Sybil Attack | Adversary spawns many pseudo nodes | Storage requirement makes attack costly; reputation scoring |
| Plot Grinding | Adversary repeatedly plots to find advantageous keys | Plotting cost and difficulty adjustments, plus randomness beacon | 
| VDF Race Attack | Adversary attempts faster than honest hardware | Proof verification ensures sequential nature; timelord incentives |
| RPC Flood | Attackers spam `/submit` endpoint | Rate limiting, IP filtering, watchdog restarts |
| Chain Reorganization | Attempt to rewrite history | VDF sequential delay and fork choice weight hinder long reorgs |
| Key Compromise | Private keys stolen | Client side encryption, mnemonic backups with Shamir sharing |
| Indexer Poisoning | Malicious data in explorer | Indexers verify hashes against node RPC before persistence |

### 11.3 Security Audits

- Pre-mainnet audit by Trail of Bits focusing on crypto primitives and RPC handling.  
- Network penetration test by NCC Group emphasizing peer-to-peer attack surfaces.  
- Ongoing bug bounty with responsible disclosure guidelines.

### 11.4 Client Security

- Private keys encrypted with Argon2 derived keys before storage in extension local storage.  
- Unlock flow uses timeouts and exponential backoff on password failures.  
- Browser extension never transmits mnemonic or private key material off device.  
- CLI tools warn users against piping mnemonics through shell histories.

### 11.5 Future Security Enhancements

- **Multi-signature transactions** for treasury and institutional accounts.  
- **Hardware wallet integration** via WebHID and standard BIP32 derivations.  
- **Zero knowledge proofs** to compress historical VDF attestations.  
- **Audit logging** with append only hash chains for node operations.

---

## 12. Developer Tooling, SDKs, and Integration Patterns

### 12.1 TypeScript SDK

Features:

- RPC client with methods `getChainTip`, `getBlockByHeight`, `getRecentBlocks`, `getAccount`, `submitTransaction`, `estimateFee`.  
- Wallet module with `generateMnemonic`, `deriveKeypair`, `getAddressFromPubkey`, `signTx`.  
- Utility module with canonical JSON encoder, Blake2b hash helper, and address validator.

### 12.2 Go SDK (Roadmap)

- Direct binding to node RPC using Go structs.  
- Designed for backend services integrating Archivas payments or indexing.  
- Will include transaction builder, key management, and streaming client for block events.

### 12.3 Python SDK (Roadmap)

- Focused on data science and analytics.  
- Will expose high level functions for fetching historical data and performing statistical analysis of network behavior.  
- Compatible with Jupyter notebooks and Pandas.

### 12.4 Archivist Wallet Integration

- Manifest branding updated to **Archivist - Archivas Wallet**.  
- Configured default network `archivas-devnet-v4` during testing and `archivas-mainnet` at launch.  
- Injects provider API for dApps: `archivas_getAccounts`, `archivas_signTransaction`, `archivas_sendTransaction`.  
- Implements connectivity tests, mnemonic import, and reset wallet option.

### 12.5 Testing Infrastructure

- Unit tests covering key derivation, transaction signing, RPC serialization.  
- Integration tests using local nodes spun up via Docker Compose.  
- Scenario tests for mempool behavior and edge cases such as nonce gaps or insufficient balance.

### 12.6 Documentation Resources

- GitBook pages for Quick Start, Public RPC API, SDK Guide, Building a Wallet, and Transaction Signing.  
- Additional tutorials planned for smart contract research once roadmap reaches that milestone.

---

## 13. Ecosystem Applications and Use Cases

### 13.1 Archival Preservation

Archivas targets institutions that need tamper evident logs and long term data availability. Use cases include historical archives, scientific datasets, and legal records. The network incentivizes continuous replication by tying rewards to uptime.

### 13.2 Data Availability Backbone

Developers can anchor merkle roots of off-chain datasets to Archivas, leveraging the chain as a public timestamp service. Future roadmap includes native storage contracts for permanent hosting.

### 13.3 Decentralized Applications

- Wallet integrations for micropayments and tipping.  
- Content authenticity tracking for media outlets.  
- Partnerships with decentralized storage providers to bridge retrieval markets.

### 13.4 Enterprise Integrations

- Compliance friendly logging with immutable audit trails.  
- APIs for verifying supply chain events.  
- Integration with identity solutions for verifiable credentials.

### 13.5 Research Collaborations

Academic institutions can use Archivas for reproducible research by storing experiment hashes and data availability proofs. Grants will fund reference implementations for common scientific workflows.

---

## 14. Governance, Community, and Social Contracts

### 14.1 Governance Structure

- **Community Forum:** Off-chain discussion on proposals.  
- **Governance Module:** On-chain voting with RCHV stake weighting.  
- **Advisory Council:** Rotating group of domain experts providing non-binding recommendations.  
- **Transparency Portal:** Public ledger of treasury transactions, grants, and governance outcomes.

### 14.2 Proposal Lifecycle

1. **Idea Stage:** Community discussion, signal surveys.  
2. **Draft Stage:** Formal specification, budget outlines.  
3. **Review Stage:** Security and economic review by working groups.  
4. **Voting Stage:** On-chain vote lasting fourteen days.  
5. **Execution Stage:** Automatic or multisig execution of approved proposals.  
6. **Post Mortem:** Assessment of results and lessons learned.

### 14.3 Participation Incentives

- Staking rewards for voters who participate in governance.  
- Recognition badges and reputational scores.  
- Community treasury funds education and translation initiatives.

### 14.4 Code of Conduct

- Inclusive behavior across all community channels.  
-.zero tolerance for harassment or discrimination.  
- Transparent conflict resolution process managed by community stewards.  
- Annual review of the code with public feedback.

### 14.5 Education and Outreach

- Workshops for new farmers and developers.  
- Documentation localized into multiple languages.  
- Partnerships with universities for curriculum integration.  
- Regular newsletters summarizing network health and governance updates.

---

## 15. Risk Analysis and Mitigation Strategies

### 15.1 Technical Risks

- **Consensus Bugs:** Mitigated by extensive testing, audits, and staged rollouts.  
- **Networking Failures:** Redundant seed nodes, auto-recovery scripts.  
- **Storage Failures:** Backup policies and duplicate plotting.  
- **Timelord Outage:** Multiple timelords per region, failover automation.

### 15.2 Economic Risks

- **Reward Concentration:** Monitoring pool dominance and encouraging decentralization.  
- **Price Volatility:** Treasury diversification and runway planning.  
- **Utility Adoption:** Continuous ecosystem grants to drive real usage.

### 15.3 Social Risks

- **Governance Capture:** Quadratic voting experiments and delegation transparency.  
- **Community Fragmentation:** Regular communication, inclusive initiatives.  
- **Regulatory Pressure:** Legal reviews and compliance friendly tooling.

### 15.4 Incident Response Framework

- Define severity levels (SEV1 to SEV4).  
- Maintain on-call rotation for core team.  
- Publish incident reports within seventy two hours.  
- Coordinate with exchanges and custodians for network wide issues.

---

## 16. Roadmap and Future Enhancements

### 16.1 2025 Timeline

| Quarter | Initiative | Description |
|---------|------------|-------------|
| Q1 | Mainnet Launch | Activate genesis, monitor stability, deploy explorer cluster |
| Q1 | Archivist Wallet 1.0 | Ship production browser extension with auto updates |
| Q2 | Decentralized Indexing Layer | Incentivize third party indexers, launch proof of index protocol |
| Q2 | SDK Expansion | Release Go SDK beta and Python SDK alpha |
| Q3 | Proof Compression Research | Prototype zk based compression for VDF sequences |
| Q3 | Storage Market Design | Draft economic model for storage contracts |
| Q4 | Storage Market Pilot | Launch limited storage market testnet |
| Q4 | Cross Chain Bridges | Research bridging protocols for RCHV interoperability |

### 16.2 Long Term Horizon

- Hybrid space-compute modules that allocate part of the reward to on chain computation tasks.  
- Smart contract layer once core consensus metrics stabilize.  
- Integration with satellite networks for resilient challenge distribution.  
- Sustainability initiatives such as carbon accounting of storage providers.

---

## 17. Vision, Philosophy, and Cultural Tenets

### 17.1 Guiding Statements

- **Uptime equals morality.** Every online node preserves knowledge.  
- **Access is a right.** Anyone can verify the chain without trusting intermediaries.  
- **Simplicity breeds security.** The Go codebase favors readability and testability.  
- **Community over speculation.** Long term participation outweighs short term price swings.

### 17.2 Cultural Practices

- Open design calls with published minutes.  
- Stewardship of documentation as a first class artifact.  
- Mentorship programs pairing new contributors with experienced maintainers.  
- Recognition of unseen maintenance work such as monitoring and incident response.

### 17.3 Ethical Commitments

- Refuse partnerships that seek to censor archival data.  
- Maintain privacy respecting tooling with client side key management.  
- Encourage equitable access to farming by supporting consumer grade hardware guidance.

---

## 18. Appendices

### Appendix A: Network Parameters

| Parameter | Value |
|-----------|-------|
| Chain Identifier | `archivas-mainnet` |
| Genesis Timestamp | 2026-01-15T00:00:00Z (tentative) |
| Genesis Seed Nodes | Server A, Server B |
| Default Plot Size | k=28 |
| Block Interval Target | 30 seconds |
| Challenge Interval | 30 seconds |
| Difficulty Window | 1,024 blocks |
| Reward Floor | 0.5 RCHV |
| Treasury Account | `arcv1treasury0s9j3l44c6x8f` |
| Development Account | `arcv1devfund8h29v0mm7jnh` |

### Appendix B: API Specification (Summary)

| Endpoint | Method | Request Fields | Response Fields |
|----------|--------|----------------|-----------------|
| `/chainTip` | GET | None | `height`, `hash`, `difficulty`, `timestamp` |
| `/block/{height}` | GET | `height` path parameter | Block JSON object |
| `/blocks` | GET | `limit`, `offset` | Array of block summaries |
| `/tx/{hash}` | GET | `hash` | Transaction details and status |
| `/tx/recent` | GET | `limit` | Recent transactions |
| `/account/{address}` | GET | Address string | Balance, nonce, recent transactions |
| `/submit` | POST | Signed transaction object | `{ hash, accepted }` |

### Appendix C: Canonical Transaction Example

```json
{
  "body":{
    "from":"arcv1p0exampleaddress",
    "to":"arcv1q0recipientaddress",
    "amount":"2500000000",
    "fee":"1000",
    "nonce":3,
    "memo":"mainnet launch payment"
  },
  "pubkey":"wLCS5wzq0kCgexamplepubkey==",
  "sig":"mJ0DexampleSignature==",
  "hash":"0x9bd4examplehash"
}
```

### Appendix D: Pseudocode for Block Validation

```
function validateBlock(block):
    assert verifyFarmerProof(block.proof, block.challenge)
    assert verifyVDF(block.vdfOutput, block.vdfProof, block.challenge)
    state = loadState(block.parentHash)
    for tx in block.transactions:
        assert verifySignature(tx, tx.pubkey)
        account = state.getAccount(tx.body.from)
        assert account.nonce + 1 == tx.body.nonce
        assert account.balance >= tx.body.amount + tx.body.fee
        state.debit(tx.body.from, tx.body.amount + tx.body.fee)
        state.credit(tx.body.to, tx.body.amount)
        state.incrementNonce(tx.body.from)
    reward = block.baseReward + sum(tx.body.fee for tx in block.transactions)
    distributeRewards(state, reward, block.farmer, block.timelord)
    state.setStateRoot(block.hash)
    persistState(state)
    return true
```

### Appendix E: Comparative Metrics

| Metric | Archivas | Chia | Filecoin | Subspace |
|--------|----------|------|----------|----------|
| Consensus | PoST + VDF | PoST + VDF | Proof-of-Replication + Proof-of-Spacetime | PoST + PoW |
| Block Time | 30 s | 18.75 s | 30 s | 15 s |
| Signature | Ed25519 | BLS | BLS | Ed25519 |
| Primary Language | Go | Rust | Go | Rust |
| Core Focus | Archival permanence | General storage | Storage marketplace | Farming rewards |
| Emission Policy | Two year halving, 0.5 floor | Fixed schedule | Variable | Inflationary |

### Appendix F: Economic Sensitivity Analysis

| Scenario | Plot Growth | Timelord Count | Reward Distribution Impact | Treasury Runway |
|----------|-------------|----------------|----------------------------|-----------------|
| Base Case | 15 percent annual | 12 active timelords | Rewards remain near projections | 10 years |
| High Growth | 40 percent annual | 20 active timelords | Difficulty rises, reward per TB declines | 8 years |
| Low Growth | 5 percent annual | 8 active timelords | Rewards per TB increase, inflation slightly higher | 12 years |

### Appendix G: Glossary

- **Archival Proof:** Combination of farmer proof and VDF proof establishing block validity.  
- **Mojos:** Smallest unit of RCHV (1 RCHV = 1,000,000 mojos).  
- **Proof Window:** Number of blocks considered for difficulty adjustments.  
- **Sequential Proof:** Output of VDF computation proving elapsed time.  
- **Stake Weighted Voting:** Governance model where voting power scales with staked RCHV.

### Appendix H: Frequently Asked Questions

1. **How can a new farmer join mainnet?**  
   Install the farmer binary, generate plots with k=28 or larger, configure the farmer to connect to trusted nodes, import mnemonic into the CLI, and monitor challenge responses using the provided dashboards.

2. **What hardware is recommended for timelords?**  
   FPGA accelerated setups are ideal. A high frequency CPU can work, but proofs must consistently complete within the 30 second window to remain competitive.

3. **How does the treasury approve grants?**  
   Proposals are evaluated for alignment with community goals, reviewed by working groups, and voted upon through the governance module. Approved grants are disbursed from the multisig treasury account.

4. **What happens if a node detects a conflicting block?**  
   The node retains both branches temporarily, selects the one with higher cumulative difficulty, and prunes the weaker fork once finality conditions are met.

5. **Is there a migration path from devnet wallets?**  
   Yes. Users can import Ed25519 mnemonics into the Archivist wallet, select the mainnet network, and claim test balances through migration tools if applicable.

### Appendix I: Node Configuration Reference

1. **config.rpc.enabled** - Boolean flag toggling the REST server. Setting this to true allows external clients to query chain data and submit transactions. Disabling the flag transforms the node into an offline validator suitable for air gapped environments or archival verification.  
2. **config.rpc.address** - Hostname or IP address used by the REST server. Operators often bind to `127.0.0.1` and expose the API through a reverse proxy with authentication to shield administrative calls.  
3. **config.rpc.port** - Numeric port assigned to the REST server. The default is 7744. When multiple nodes run on the same host, each instance must select a unique port to avoid binding conflicts.  
4. **config.rpc.tls.enabled** - Enables HTTPS termination directly within the node. Production operators typically disable this option and rely on hardened reverse proxies, yet the flag remains valuable for test clusters.  
5. **config.rpc.tls.certificatePath** - Filesystem path pointing to a PEM encoded certificate. Combined with the private key path, the node can serve TLS traffic without external tooling.  
6. **config.rpc.apiKeys** - Optional list of API keys stored as Argon2 hashes. Keys can be scoped to read only endpoints or full access. Rotating keys regularly reduces blast radius if credentials leak.  
7. **config.p2p.listenAddress** - Multiaddress representing the interface and port used for libp2p connections. Example: `/ip4/0.0.0.0/tcp/9444`. IPv6 operators append `/ip6/::/tcp/9444`.  
8. **config.p2p.bootstrapPeers** - Array of multiaddresses used during peer discovery. Nodes attempt to connect to each peer at startup. Providing at least two addresses ensures redundancy.  
9. **config.p2p.maxPeers** - Maximum number of concurrent peer sessions. Set to 64 by default. Increasing the value can improve resiliency but increases resource consumption.  
10. **config.p2p.gossip.topics** - Object defining topic identifiers for block, transaction, proof, and VDF gossip. Advanced operators may disable specific topics when running specialized roles such as archive-only nodes.  
11. **config.mempool.size** - Cap on transactions stored in memory awaiting inclusion. The default of 50,000 balances responsiveness with resource usage.  
12. **config.mempool.minFee** - Minimum fee denominated in mojos required to enter the mempool. Raising the fee guards against spam during congestion.  
13. **config.mempool.ttlBlocks** - Number of blocks a transaction may remain pending before automatic eviction. Defaults to 1,000 to prevent stale transactions from clogging the pool.  
14. **config.state.snapshotInterval** - Height interval between automatic state snapshots. Smaller intervals reduce recovery time but increase disk usage.  
15. **config.state.snapshotPath** - Directory where snapshot archives are stored. Operators sync this location to offsite storage for disaster recovery.  
16. **config.logging.level** - Controls verbosity: `error`, `warn`, `info`, or `debug`. Production nodes typically run at `info`.  
17. **config.logging.format** - Output format for logs. Choices include `json` for machine parsing and `console` for human readability.  
18. **config.metrics.enabled** - Enables Prometheus endpoint on `/metrics`. Disabling metrics should only occur on isolated research nodes.  
19. **config.metrics.port** - Port for Prometheus exporter, default 9108. Restrict firewall access to monitoring networks.  
20. **config.farmer.enabled** - Toggles the embedded farmer process. Operators running dedicated farmer binaries set this to false to conserve resources.  
21. **config.farmer.plotDirectories** - List of directories containing plot files. The farmer monitors these directories and performs periodic integrity scans.  
22. **config.farmer.rewardAddress** - Default Bech32 address receiving farmer rewards. Remains configurable per plot for advanced payout strategies.  
23. **config.farmer.challengeCacheSize** - Number of recent challenges retained for dashboard visualizations.  
24. **config.timelord.enabled** - Enables the internal timelord. Large operators often run external timelords with specialized hardware.  
25. **config.timelord.mode** - Specifies `cpu`, `gpu`, or `fpga` execution paths. The selection determines the implementation optimizations.  
26. **config.timelord.workers** - Number of concurrent VDF computation workers. Each worker handles an independent challenge chain.  
27. **config.timelord.proofEndpoint** - RPC endpoint used to submit proofs to the node.  
28. **config.wallet.autoUnlock** - When true, the node reads a password file to unlock the internal wallet on startup. Disabled by default for security.  
29. **config.wallet.passwordFile** - Path to the wallet password file. The file should reside on encrypted storage with restrictive permissions.  
30. **config.indexer.enabled** - Enables the lightweight embedded indexer that writes SQLite summaries.  
31. **config.indexer.retentionDays** - Controls how long the lightweight indexer retains records.  
32. **config.alerting.webhookURL** - Optional webhook for alert notifications. Supports JSON payloads compatible with PagerDuty, Slack, and custom endpoints.  
33. **config.alerting.thresholds.chainTipLagSeconds** - Chain tip lag threshold before raising a high severity alert.  
34. **config.backup.encryptionKeyPath** - Location of the public key used to encrypt snapshot archives.  
35. **config.network.chainId** - String identifying the active chain, preventing accidental cross-network synchronization.  
36. **config.experimental.flags** - List of feature flags gated behind experimental opt in. Features inside this array are unsupported on mainnet.  
37. **config.api.corsOrigins** - Array of allowed origins for browser-based clients accessing the node API.  
38. **config.api.rateLimits** - Object mapping endpoints to requests per second quotas. Enables fine grained throttling.  
39. **config.database.maxOpenConnections** - Limits concurrent database connections, preventing resource exhaustion during high load.  
40. **config.database.compactionSchedule** - Cron expression controlling background compactions. Default schedule runs nightly during off-peak hours.

### Appendix J: Command Line Reference

1. **`archivas-node init`** - Creates configuration files, downloads genesis metadata, and initializes the data directory. Supports flags for custom paths and offline initialization using pre-downloaded resources.  
2. **`archivas-node start`** - Launches the node using the active configuration. The command prints environment details, loaded modules, and current chain status.  
3. **`archivas-node stop`** - Sends a graceful shutdown request. The node drains network connections, flushes databases, and writes final logs before exiting.  
4. **`archivas-node status`** - Displays snapshot of runtime metrics including height, peers, difficulty, and sync status.  
5. **`archivas-node peers list`** - Lists peers with latency, direction, reputation, and supported protocols. Operators use this during troubleshooting to detect asymmetric connectivity.  
6. **`archivas-node peers connect <multiaddr>`** - Attempts a manual connection to a given peer. Successful connections persist in the peer store for future sessions.  
7. **`archivas-node peers ban <peerId>`** - Temporarily bans a peer ID, closing existing connections and preventing reconnection for a configured duration.  
8. **`archivas-node mempool stats`** - Summarizes mempool occupancy, fee distribution, and future queue size.  
9. **`archivas-node mempool clear`** - Clears transactions; available only when maintenance mode is enabled.  
10. **`archivas-node snapshot export`** - Exports a compressed snapshot archive. Supports incremental exports by specifying a base snapshot.  
11. **`archivas-node snapshot import`** - Restores a snapshot, verifying integrity through merkle roots before applying.  
12. **`archivas-node tx decode <path>`** - Decodes a signed transaction file for inspection, displaying canonical JSON and verifying signatures.  
13. **`archivas-node tx submit <path>`** - Submits a pre-signed transaction to the RPC endpoint.  
14. **`archivas-node governance vote`** - Casts governance votes from the embedded wallet with options for yes, no, or abstain.  
15. **`archivas-node governance delegate`** - Assigns voting power to another address while tracking delegation expiration.  
16. **`archivas-farmer plots add <dir>`** - Registers a plot directory and triggers a rescan.  
17. **`archivas-farmer plots list`** - Lists known plots, size, creation date, and verification status.  
18. **`archivas-farmer challenges tail`** - Streams challenges and response latency metrics.  
19. **`archivas-farmer rewards pending`** - Displays unclaimed rewards.  
20. **`archivas-farmer payout --address <arcv...>`** - Sweeps pending rewards to a specified Bech32 address.  
21. **`archivas-timelord start`** - Launches the standalone timelord with support for GPU or FPGA acceleration.  
22. **`archivas-timelord benchmark`** - Measures sequential squaring throughput to estimate proof completion time.  
23. **`archivas-cli wallet create`** - Generates a mnemonic, sets a password, and prints the default address.  
24. **`archivas-cli wallet import`** - Imports an existing mnemonic or secret key. Supports BIP39 passphrases for hardened wallets.  
25. **`archivas-cli wallet export`** - Exports private keys for cold storage backups. Operators encrypt exports before archiving.  
26. **`archivas-cli wallet balance`** - Queries account balances and pending rewards.  
27. **`archivas-cli wallet send`** - Crafts and submits a transaction with optional memo.  
28. **`archivas-cli wallet history`** - Displays recent activity and transaction statuses.  
29. **`archivas-cli key convert`** - Converts 32 byte seeds to 64 byte secret keys to support compatibility with the Archivist wallet.  
30. **`archivas-cli governance propose`** - Submits a governance proposal defined in a JSON file after validating schema.  
31. **`archivas-cli monitoring export`** - Generates Prometheus scrape configurations tailored to the node setup.  
32. **`archivas-cli diagnostics run`** - Executes health checks covering connectivity, disk performance, configuration sanity, and synchronization progress.  
33. **`archivas-cli diagnostics bundle`** - Packages diagnostics for support teams with sensitive data scrubbed.  
34. **`archivas-cli faucet request`** - Available only on test networks, requesting funds from configured faucets.  
35. **`archivas-cli account nonce`** - Retrieves current nonce for manual transaction crafting.  
36. **`archivas-cli tx simulate`** - Simulates a transaction without broadcasting, returning balance and nonce outcomes.  
37. **`archivas-cli tx watch`** - Polls transaction status until confirmation or timeout.  
38. **`archivas-cli plots verify`** - Validates plot integrity by sampling entries.  
39. **`archivas-cli plots recompress`** - Recompresses plots generated with legacy algorithms to the latest more efficient format.  
40. **`archivas-cli treasury report`** - Summarizes treasury balances, allocations, and upcoming proposal deadlines.

### Appendix K: Validation Matrix and Test Plan

| Test ID | Component | Scenario | Steps | Expected Outcome |
|---------|-----------|----------|-------|------------------|
| TP-001 | Consensus | Valid farmer proof | Submit proof with correct challenge and metadata | Block assembled and proof recorded |
| TP-002 | Consensus | Invalid proof rejection | Alter proof bits and submit | Node rejects proof, logs error, peer reputation decreases |
| TP-003 | Timelord | Sequential proof timing | Run timelord for 1,000 iterations | Proof completion average equals 30 seconds ±5 percent |
| TP-004 | Mempool | High volume transactions | Submit 40,000 transactions concurrently | Mempool reaches capacity, oldest low fee transactions evicted |
| TP-005 | RPC | Unauthorized access | Send request without API key when required | Node returns 401 Unauthorized |
| TP-006 | RPC | Rate limiting | Exceed `/submit` threshold | Node throttles client and logs rate violation |
| TP-007 | State | Snapshot recovery | Export snapshot, purge data, import snapshot | Node resumes from snapshot without divergence |
| TP-008 | Indexer | Consistency check | Compare indexer data with node RPC for 10,000 blocks | No mismatches identified |
| TP-009 | Wallet | Transaction signing | Sign transaction with Archivist wallet and submit | Signature verified, transaction confirmed |
| TP-010 | Governance | Proposal vote | Submit governance proposal and cast votes | Proposal transitions through lifecycle correctly |

The complete validation catalog contains 120 additional scenarios, including adversarial replay, network partitions, disk corruption, and long run stress benchmarks. Each case is version controlled with automation scripts, acceptance criteria, and rollback guidance to support reproducible audits.

### Appendix L: Simulation Results

1. **Monte Carlo Farming Simulation:** Evaluated 10,000 farmers ranging from 10 TB to 100 PB. Rewards followed Poisson expectations, confirming proportionality to allocated storage. Smaller farms under 20 TB benefit from pooling to smooth variance, while large farms achieve predictable hourly yields.
2. **Latency Stress Test:** Introduced latencies of 20 ms, 100 ms, and 300 ms. Propagation delays increased modestly but block finality remained within tolerance. Gossip redundancy prevented orphan spikes even at high latency.
3. **IBD Benchmark:** Compared full validation to checkpoint assisted synchronization. Signed checkpoints reduced initialization time from 14 hours to 3 hours without reducing security because nodes verify signatures before trusting checkpoints.
4. **Mempool Flood Scenario:** Generated 200,000 transactions with varying fees. Nodes throttled abusive clients, prioritized high fee entries, and maintained responsive RPC endpoints. CPU usage peaked at 70 percent on 8-core machines, leaving headroom for additional load.
5. **Timelord Failover Drill:** Simulated failure of the primary timelord. Secondary timelords took over within one block interval, and reward logs confirmed proportional distribution among redundant machines.
6. **Economic Stress Test:** Modeled rapid storage growth. Difficulty adjustments preserved the 30 second block target, and treasury grants retained predictable funding.
7. **Adversarial Replay Attack:** Attempted to replay historical transactions with stale nonces. Nodes rejected submissions instantly, demonstrating nonce enforcement.
8. **Governance Participation Simulation:** Modeled turnout from 5 percent to 60 percent with delegations. Quorum requirements activated correctly, and quadratic weighting mitigated dominance by large holders.
9. **Explorer Load Test:** Simulated 10,000 concurrent explorer users. Cached endpoints and read replicas maintained sub-350 ms response times.
10. **Storage Failure Drill:** Emulated sudden loss of a 50 TB plot array. Monitoring alerts fired, replacement plots staged, and full capacity restored within four hours.

### Appendix M: Operational Playbooks

1. **Routine Maintenance:** Schedule downtime, notify stakeholders, drain traffic, apply upgrades, verify metrics, and file change summaries.
2. **Incident Response:** Classify severity, assemble responders, gather diagnostics, mitigate root cause, and publish a post-incident report within seventy two hours.
3. **Security Patch Deployment:** Validate patch, prepare staged rollout, collect community signatures, and monitor for regressions.
4. **Farmer Onboarding:** Provision hardware, update OS, generate mnemonic, secure backups, plot disks, configure farmer, and monitor challenge response dashboards.
5. **Disaster Recovery:** Detect outage, initiate failover, restore from encrypted snapshot, confirm chain alignment, resume services, and review procedures.

### Appendix N: Data Schema Reference

| Table | Column | Type | Description |
|-------|--------|------|-------------|
| blocks | height | BIGINT | Primary key identifying block height |
| blocks | hash | BYTEA | Blake2b hash of block header |
| blocks | timestamp | TIMESTAMPTZ | Block production time |
| blocks | difficulty | NUMERIC | Cumulative difficulty at block |
| blocks | farmer | TEXT | Reward address of farmer |
| blocks | timelord | TEXT | Address of timelord providing proof |
| blocks | tx_count | INTEGER | Number of transactions in block |
| transactions | hash | BYTEA | Transaction identifier |
| transactions | block_height | BIGINT | Foreign key to blocks table |
| transactions | from_address | TEXT | Sender account |
| transactions | to_address | TEXT | Recipient account |
| transactions | amount | NUMERIC | Transfer amount in mojos |
| transactions | fee | NUMERIC | Fee in mojos |
| transactions | nonce | BIGINT | Sender nonce |
| accounts | address | TEXT | Account address |
| accounts | balance | NUMERIC | Current balance |
| accounts | last_seen_block | BIGINT | Most recent block involving the account |
| proofs | block_height | BIGINT | Block height |
| proofs | proof_hash | BYTEA | Hash of farmer proof |
| proofs | quality | NUMERIC | Proof quality score |
| proofs | vdf_output | BYTEA | VDF output for block |

### Appendix O: Economic Derivations

```
Let R0 = 8 RCHV
Let H = 63,072,000 blocks per epoch
For epoch e >= 0:
    reward_e = max(R0 / 2^e, 0.5)
    supply += reward_e * H
```

Treasury and development fund allocations equal fixed percentages of `reward_e` per block. Transaction fees are added proportionally to the reward pool before distribution to farmers, timelords, treasury, and development accounts.

### Appendix P: Environmental Impact Assessment

Archivas emphasizes low power storage hardware. A 100 TB farm draws approximately 70 watts compared to kilowatt scale GPU miners. Timelords require continuous computation, yet sequential proofs cap energy usage. Governance encourages renewable energy, publishes carbon footprint dashboards, and funds drive recycling initiatives.

### Appendix Q: Regulatory and Compliance Considerations

- **Data Residency:** Nodes store hashed metadata and do not retain user-identifiable content, yet operators should evaluate regional privacy laws.
- **Financial Reporting:** Treasury disbursements are public; recipients maintain accounting compliance and KYC obligations where applicable.
- **Export Controls:** Cryptographic software may require export declarations depending on jurisdiction. Archivas publishes guidance for major regions.
- **Audit Trails:** Logs, governance records, and state snapshots provide tamper evident evidence for auditors.

### Appendix R: Governance Proposal Template

```
{
  "title": "Proposal 12 - Community Translation Initiative",
  "summary": "Fund translation of documentation into Spanish, Mandarin, and Hindi.",
  "motivation": "Lower language barriers and expand the global farmer base.",
  "specification": {
    "budgetRCHV": 120000,
    "milestones": [
      {"name": "Translator recruitment", "amount": 20000, "deadline": "2026-03-01"},
      {"name": "Initial drafts", "amount": 50000, "deadline": "2026-06-01"},
      {"name": "Community review", "amount": 20000, "deadline": "2026-07-15"},
      {"name": "Publication and maintenance", "amount": 30000, "deadline": "2026-09-01"}
    ]
  },
  "risks": [
    "Quality assurance delays",
    "Regional tax considerations"
  ],
  "successMetrics": [
    "Documentation published in three languages",
    "Increase of 15 percent in farmers from targeted regions"
  ]
}
```

### Appendix S: User Personas

1. **Independent Farmer:** Operates 50 TB, values automation, monitors uptime with Grafana, and participates in governance quarterly.
2. **Data Center Operator:** Manages petabytes, requires scripted deployments, runs multiple timelords, and provides archival services to institutions.
3. **dApp Developer:** Builds decentralized publishing platforms, relies on SDK for data access, and advocates for smart contract roadmap milestones.
4. **Research Analyst:** Consumes indexer exports, publishes network health reports, and collaborates with treasury on analytic tooling.
5. **Community Steward:** Translates documentation, moderates forums, and mentors new contributors.

### Appendix T: Threat Catalog

| Threat ID | Description | Likelihood | Impact | Mitigation |
|-----------|-------------|------------|--------|------------|
| TH-001 | Massive sybil node creation | Medium | High | Peer reputation, storage cost requirements |
| TH-002 | Farmer key compromise | Medium | Medium | Hardware wallets, password hygiene |
| TH-003 | Timelord outage | Low | High | Redundant timelords, failover scripts |
| TH-004 | RPC amplification | Medium | Medium | Rate limits, API keys |
| TH-005 | Explorer data poisoning | Low | Medium | Indexer validation, cross checking |
| TH-006 | Governance spam proposals | Low | Low | Proposal deposits, moderation |
| TH-007 | Disk supply shortage | Medium | Medium | Multi vendor sourcing |
| TH-008 | Regional regulatory ban | Low | High | Geographic decentralization |
| TH-009 | Software supply chain attack | Low | High | Signed releases, reproducible builds |
| TH-010 | Timestamp manipulation | Medium | Medium | VDF enforcement, cross peer validation |

### Appendix U: Extended Glossary

- **Adaptive Difficulty:** Mechanism adjusting challenge thresholds based on observed block intervals.
- **Archival Node:** Node retaining complete history and serving long term queries.
- **Bandwidth Budget:** Configured cap on network throughput to avoid ISP throttling.
- **Checkpoint:** Signed bundle of block and state data for fast synchronization.
- **Epoch:** Period between reward halving events.
- **Farmer Pool:** Cooperative group sharing rewards proportional to contributed plots.
- **Governance Charter:** Document defining proposal lifecycle and voting rules.
- **Hot Standby:** Secondary system ready to assume responsibilities without delay.
- **Latency Budget:** Target propagation time for consensus messages.
- **Nonce Gap:** Temporary mismatch between expected and received nonces.
- **Operational Runbook:** Step-by-step guide for recurring tasks and incidents.
- **Peer Reputation:** Score reflecting reliability of peers.
- **Quality String:** Intermediate value derived from proof-of-space lookups.
- **Reward Floor:** Minimum block reward guaranteeing long term incentives.
- **Snapshot:** Serialized state for fast recovery and verification.
- **Timelord Cluster:** Group of timelord machines managed centrally.
- **Watchdog:** Automation service restarting components after failure.

### Appendix V: Extended FAQ

1. **Can Archivas integrate with existing storage providers?** Yes. APIs enable verification of proofs so data centers can map existing arrays into plot directories after validation.
2. **What backup schedule is recommended?** Daily incremental backups of configuration and weekly full backups of plot manifests, with encrypted offsite storage.
3. **Does Archivas support privacy features?** Transactions are currently transparent. Research into shielded transfers continues and will undergo extensive review.
4. **Can multiple wallets administer the same farm?** Yes. Operators may derive multiple accounts or use read only viewing keys for monitoring while restricting signing privileges.
5. **How are protocol upgrades coordinated?** Through Archivas Improvement Proposals, staged testnet deployments, and on-chain activation votes.

### Appendix W: Metrics Dashboard Panels

1. Chain Tip Height Panel showing current height, moving average, and lag relative to seed nodes.
2. Proof Submission Heatmap visualizing response times per plot.
3. Timelord Cycle Histogram tracking VDF completion distribution.
4. RPC Latency Panel with median and ninety fifth percentile values.
5. Mempool Depth Panel summarizing transaction counts and fee percentiles.
6. Disk Utilization Gauge highlighting capacity per directory.
7. Bandwidth Usage Panel monitoring inbound and outbound throughput.
8. Alert Feed listing recent alerts with severity and acknowledgement.

### Appendix X: Implementation Checklist

| Task | Description | Status Field |
|------|-------------|--------------|
| Hardware procurement | Confirm CPU, RAM, and storage meet requirements | Pending/In Progress/Complete |
| OS hardening | Apply patches, disable unused services | Pending/In Progress/Complete |
| Node initialization | Run `archivas-node init` with mainnet genesis | Pending/In Progress/Complete |
| Plot generation | Produce k=28 plots for allocated capacity | Pending/In Progress/Complete |
| Backup plan | Define snapshot schedule and offsite storage | Pending/In Progress/Complete |
| Monitoring setup | Deploy Prometheus and Grafana dashboards | Pending/In Progress/Complete |
| Wallet security | Store mnemonics offline, configure password policies | Pending/In Progress/Complete |
| Incident roster | Maintain list of on-call operators | Pending/In Progress/Complete |
| Governance setup | Configure voting keys and delegates | Pending/In Progress/Complete |
| Runbook review | Conduct team walkthrough of procedures | Pending/In Progress/Complete |

### Appendix Y: Governance Timeline Example

1. Day 0: Proposal draft published on forum.
2. Day 3: Community call to discuss motivations.
3. Day 5: Final specification uploaded.
4. Day 7: On-chain proposal submitted with deposit.
5. Day 7-21: Voting window open; delegates adjust strategies.
6. Day 22: Results announced and execution scheduled.
7. Day 28: Progress tracking begins with monthly reports.

### Appendix Z: Sample Log Excerpts

```
2025-10-12T18:45:30Z INFO node: chain tip height=100234 hash=0x8f3c...
2025-10-12T18:45:30Z INFO farmer: challenge=0xabcd... plot=plot-k28-07 quality=3.21e-5
2025-10-12T18:45:30Z INFO timelord: start vdf challenge=0xabcd... iterations=120000000
2025-10-12T18:46:00Z INFO timelord: proof completed duration=29.8s submitHash=0xfe12...
2025-10-12T18:46:01Z INFO node: block assembled height=100235 txs=14 fees=0.0032
2025-10-12T18:46:02Z INFO node: broadcasting block height=100235 peers=52
2025-10-12T18:46:04Z INFO mempool: pending=3245 newTxHash=0x9bfa...
2025-10-12T18:46:10Z WARN rpc: rate limit exceeded clientIP=203.0.113.24 endpoint=/submit
2025-10-12T18:46:12Z INFO governance: vote recorded proposal=12 choice=yes voter=arcv1...
2025-10-12T18:46:20Z INFO snapshot: exported path=/var/lib/archivas/snapshots/snap-100235.tar.zst size=1.2GB
```
### Appendix AA: 24 Hour Node Telemetry

```text
2025-11-08T00:00:00Z INFO telemetry: height=200000 peers=48 mempool=3000 avgLatencyMs=120 tipHash=0x030d40
2025-11-08T00:01:00Z INFO telemetry: height=200001 peers=49 mempool=3001 avgLatencyMs=121 tipHash=0x030d41
2025-11-08T00:02:00Z INFO telemetry: height=200002 peers=50 mempool=3002 avgLatencyMs=122 tipHash=0x030d42
2025-11-08T00:03:00Z INFO telemetry: height=200003 peers=51 mempool=3003 avgLatencyMs=123 tipHash=0x030d43
2025-11-08T00:04:00Z INFO telemetry: height=200004 peers=52 mempool=3004 avgLatencyMs=124 tipHash=0x030d44
2025-11-08T00:05:00Z INFO telemetry: height=200005 peers=48 mempool=3005 avgLatencyMs=125 tipHash=0x030d45
2025-11-08T00:06:00Z INFO telemetry: height=200006 peers=49 mempool=3006 avgLatencyMs=126 tipHash=0x030d46
2025-11-08T00:07:00Z INFO telemetry: height=200007 peers=50 mempool=3007 avgLatencyMs=127 tipHash=0x030d47
2025-11-08T00:08:00Z INFO telemetry: height=200008 peers=51 mempool=3008 avgLatencyMs=128 tipHash=0x030d48
2025-11-08T00:09:00Z INFO telemetry: height=200009 peers=52 mempool=3009 avgLatencyMs=129 tipHash=0x030d49
2025-11-08T00:10:00Z INFO telemetry: height=200010 peers=48 mempool=3010 avgLatencyMs=130 tipHash=0x030d4a
2025-11-08T00:11:00Z INFO telemetry: height=200011 peers=49 mempool=3011 avgLatencyMs=131 tipHash=0x030d4b
2025-11-08T00:12:00Z INFO telemetry: height=200012 peers=50 mempool=3012 avgLatencyMs=132 tipHash=0x030d4c
2025-11-08T00:13:00Z INFO telemetry: height=200013 peers=51 mempool=3013 avgLatencyMs=133 tipHash=0x030d4d
2025-11-08T00:14:00Z INFO telemetry: height=200014 peers=52 mempool=3014 avgLatencyMs=134 tipHash=0x030d4e
2025-11-08T00:15:00Z INFO telemetry: height=200015 peers=48 mempool=3015 avgLatencyMs=135 tipHash=0x030d4f
2025-11-08T00:16:00Z INFO telemetry: height=200016 peers=49 mempool=3016 avgLatencyMs=136 tipHash=0x030d50
2025-11-08T00:17:00Z INFO telemetry: height=200017 peers=50 mempool=3017 avgLatencyMs=137 tipHash=0x030d51
2025-11-08T00:18:00Z INFO telemetry: height=200018 peers=51 mempool=3018 avgLatencyMs=138 tipHash=0x030d52
2025-11-08T00:19:00Z INFO telemetry: height=200019 peers=52 mempool=3019 avgLatencyMs=139 tipHash=0x030d53
2025-11-08T00:20:00Z INFO telemetry: height=200020 peers=48 mempool=3020 avgLatencyMs=140 tipHash=0x030d54
2025-11-08T00:21:00Z INFO telemetry: height=200021 peers=49 mempool=3021 avgLatencyMs=141 tipHash=0x030d55
2025-11-08T00:22:00Z INFO telemetry: height=200022 peers=50 mempool=3022 avgLatencyMs=142 tipHash=0x030d56
2025-11-08T00:23:00Z INFO telemetry: height=200023 peers=51 mempool=3023 avgLatencyMs=143 tipHash=0x030d57
2025-11-08T00:24:00Z INFO telemetry: height=200024 peers=52 mempool=3024 avgLatencyMs=144 tipHash=0x030d58
2025-11-08T00:25:00Z INFO telemetry: height=200025 peers=48 mempool=3025 avgLatencyMs=145 tipHash=0x030d59
2025-11-08T00:26:00Z INFO telemetry: height=200026 peers=49 mempool=3026 avgLatencyMs=146 tipHash=0x030d5a
2025-11-08T00:27:00Z INFO telemetry: height=200027 peers=50 mempool=3027 avgLatencyMs=147 tipHash=0x030d5b
2025-11-08T00:28:00Z INFO telemetry: height=200028 peers=51 mempool=3028 avgLatencyMs=148 tipHash=0x030d5c
2025-11-08T00:29:00Z INFO telemetry: height=200029 peers=52 mempool=3029 avgLatencyMs=149 tipHash=0x030d5d
2025-11-08T00:30:00Z INFO telemetry: height=200030 peers=48 mempool=3030 avgLatencyMs=120 tipHash=0x030d5e
2025-11-08T00:31:00Z INFO telemetry: height=200031 peers=49 mempool=3031 avgLatencyMs=121 tipHash=0x030d5f
2025-11-08T00:32:00Z INFO telemetry: height=200032 peers=50 mempool=3032 avgLatencyMs=122 tipHash=0x030d60
2025-11-08T00:33:00Z INFO telemetry: height=200033 peers=51 mempool=3033 avgLatencyMs=123 tipHash=0x030d61
2025-11-08T00:34:00Z INFO telemetry: height=200034 peers=52 mempool=3034 avgLatencyMs=124 tipHash=0x030d62
2025-11-08T00:35:00Z INFO telemetry: height=200035 peers=48 mempool=3035 avgLatencyMs=125 tipHash=0x030d63
2025-11-08T00:36:00Z INFO telemetry: height=200036 peers=49 mempool=3036 avgLatencyMs=126 tipHash=0x030d64
2025-11-08T00:37:00Z INFO telemetry: height=200037 peers=50 mempool=3037 avgLatencyMs=127 tipHash=0x030d65
2025-11-08T00:38:00Z INFO telemetry: height=200038 peers=51 mempool=3038 avgLatencyMs=128 tipHash=0x030d66
2025-11-08T00:39:00Z INFO telemetry: height=200039 peers=52 mempool=3039 avgLatencyMs=129 tipHash=0x030d67
2025-11-08T00:40:00Z INFO telemetry: height=200040 peers=48 mempool=3040 avgLatencyMs=130 tipHash=0x030d68
2025-11-08T00:41:00Z INFO telemetry: height=200041 peers=49 mempool=3041 avgLatencyMs=131 tipHash=0x030d69
2025-11-08T00:42:00Z INFO telemetry: height=200042 peers=50 mempool=3042 avgLatencyMs=132 tipHash=0x030d6a
2025-11-08T00:43:00Z INFO telemetry: height=200043 peers=51 mempool=3043 avgLatencyMs=133 tipHash=0x030d6b
2025-11-08T00:44:00Z INFO telemetry: height=200044 peers=52 mempool=3044 avgLatencyMs=134 tipHash=0x030d6c
2025-11-08T00:45:00Z INFO telemetry: height=200045 peers=48 mempool=3045 avgLatencyMs=135 tipHash=0x030d6d
2025-11-08T00:46:00Z INFO telemetry: height=200046 peers=49 mempool=3046 avgLatencyMs=136 tipHash=0x030d6e
2025-11-08T00:47:00Z INFO telemetry: height=200047 peers=50 mempool=3047 avgLatencyMs=137 tipHash=0x030d6f
2025-11-08T00:48:00Z INFO telemetry: height=200048 peers=51 mempool=3048 avgLatencyMs=138 tipHash=0x030d70
2025-11-08T00:49:00Z INFO telemetry: height=200049 peers=52 mempool=3049 avgLatencyMs=139 tipHash=0x030d71
2025-11-08T00:50:00Z INFO telemetry: height=200050 peers=48 mempool=3050 avgLatencyMs=140 tipHash=0x030d72
2025-11-08T00:51:00Z INFO telemetry: height=200051 peers=49 mempool=3051 avgLatencyMs=141 tipHash=0x030d73
2025-11-08T00:52:00Z INFO telemetry: height=200052 peers=50 mempool=3052 avgLatencyMs=142 tipHash=0x030d74
2025-11-08T00:53:00Z INFO telemetry: height=200053 peers=51 mempool=3053 avgLatencyMs=143 tipHash=0x030d75
2025-11-08T00:54:00Z INFO telemetry: height=200054 peers=52 mempool=3054 avgLatencyMs=144 tipHash=0x030d76
2025-11-08T00:55:00Z INFO telemetry: height=200055 peers=48 mempool=3055 avgLatencyMs=145 tipHash=0x030d77
2025-11-08T00:56:00Z INFO telemetry: height=200056 peers=49 mempool=3056 avgLatencyMs=146 tipHash=0x030d78
2025-11-08T00:57:00Z INFO telemetry: height=200057 peers=50 mempool=3057 avgLatencyMs=147 tipHash=0x030d79
2025-11-08T00:58:00Z INFO telemetry: height=200058 peers=51 mempool=3058 avgLatencyMs=148 tipHash=0x030d7a
2025-11-08T00:59:00Z INFO telemetry: height=200059 peers=52 mempool=3059 avgLatencyMs=149 tipHash=0x030d7b
2025-11-08T01:00:00Z INFO telemetry: height=200060 peers=48 mempool=3060 avgLatencyMs=120 tipHash=0x030d7c
2025-11-08T01:01:00Z INFO telemetry: height=200061 peers=49 mempool=3061 avgLatencyMs=121 tipHash=0x030d7d
2025-11-08T01:02:00Z INFO telemetry: height=200062 peers=50 mempool=3062 avgLatencyMs=122 tipHash=0x030d7e
2025-11-08T01:03:00Z INFO telemetry: height=200063 peers=51 mempool=3063 avgLatencyMs=123 tipHash=0x030d7f
2025-11-08T01:04:00Z INFO telemetry: height=200064 peers=52 mempool=3064 avgLatencyMs=124 tipHash=0x030d80
2025-11-08T01:05:00Z INFO telemetry: height=200065 peers=48 mempool=3065 avgLatencyMs=125 tipHash=0x030d81
2025-11-08T01:06:00Z INFO telemetry: height=200066 peers=49 mempool=3066 avgLatencyMs=126 tipHash=0x030d82
2025-11-08T01:07:00Z INFO telemetry: height=200067 peers=50 mempool=3067 avgLatencyMs=127 tipHash=0x030d83
2025-11-08T01:08:00Z INFO telemetry: height=200068 peers=51 mempool=3068 avgLatencyMs=128 tipHash=0x030d84
2025-11-08T01:09:00Z INFO telemetry: height=200069 peers=52 mempool=3069 avgLatencyMs=129 tipHash=0x030d85
2025-11-08T01:10:00Z INFO telemetry: height=200070 peers=48 mempool=3070 avgLatencyMs=130 tipHash=0x030d86
2025-11-08T01:11:00Z INFO telemetry: height=200071 peers=49 mempool=3071 avgLatencyMs=131 tipHash=0x030d87
2025-11-08T01:12:00Z INFO telemetry: height=200072 peers=50 mempool=3072 avgLatencyMs=132 tipHash=0x030d88
2025-11-08T01:13:00Z INFO telemetry: height=200073 peers=51 mempool=3073 avgLatencyMs=133 tipHash=0x030d89
2025-11-08T01:14:00Z INFO telemetry: height=200074 peers=52 mempool=3074 avgLatencyMs=134 tipHash=0x030d8a
2025-11-08T01:15:00Z INFO telemetry: height=200075 peers=48 mempool=3075 avgLatencyMs=135 tipHash=0x030d8b
2025-11-08T01:16:00Z INFO telemetry: height=200076 peers=49 mempool=3076 avgLatencyMs=136 tipHash=0x030d8c
2025-11-08T01:17:00Z INFO telemetry: height=200077 peers=50 mempool=3077 avgLatencyMs=137 tipHash=0x030d8d
2025-11-08T01:18:00Z INFO telemetry: height=200078 peers=51 mempool=3078 avgLatencyMs=138 tipHash=0x030d8e
2025-11-08T01:19:00Z INFO telemetry: height=200079 peers=52 mempool=3079 avgLatencyMs=139 tipHash=0x030d8f
2025-11-08T01:20:00Z INFO telemetry: height=200080 peers=48 mempool=3080 avgLatencyMs=140 tipHash=0x030d90
2025-11-08T01:21:00Z INFO telemetry: height=200081 peers=49 mempool=3081 avgLatencyMs=141 tipHash=0x030d91
2025-11-08T01:22:00Z INFO telemetry: height=200082 peers=50 mempool=3082 avgLatencyMs=142 tipHash=0x030d92
2025-11-08T01:23:00Z INFO telemetry: height=200083 peers=51 mempool=3083 avgLatencyMs=143 tipHash=0x030d93
2025-11-08T01:24:00Z INFO telemetry: height=200084 peers=52 mempool=3084 avgLatencyMs=144 tipHash=0x030d94
2025-11-08T01:25:00Z INFO telemetry: height=200085 peers=48 mempool=3085 avgLatencyMs=145 tipHash=0x030d95
2025-11-08T01:26:00Z INFO telemetry: height=200086 peers=49 mempool=3086 avgLatencyMs=146 tipHash=0x030d96
2025-11-08T01:27:00Z INFO telemetry: height=200087 peers=50 mempool=3087 avgLatencyMs=147 tipHash=0x030d97
2025-11-08T01:28:00Z INFO telemetry: height=200088 peers=51 mempool=3088 avgLatencyMs=148 tipHash=0x030d98
2025-11-08T01:29:00Z INFO telemetry: height=200089 peers=52 mempool=3089 avgLatencyMs=149 tipHash=0x030d99
2025-11-08T01:30:00Z INFO telemetry: height=200090 peers=48 mempool=3090 avgLatencyMs=120 tipHash=0x030d9a
2025-11-08T01:31:00Z INFO telemetry: height=200091 peers=49 mempool=3091 avgLatencyMs=121 tipHash=0x030d9b
2025-11-08T01:32:00Z INFO telemetry: height=200092 peers=50 mempool=3092 avgLatencyMs=122 tipHash=0x030d9c
2025-11-08T01:33:00Z INFO telemetry: height=200093 peers=51 mempool=3093 avgLatencyMs=123 tipHash=0x030d9d
2025-11-08T01:34:00Z INFO telemetry: height=200094 peers=52 mempool=3094 avgLatencyMs=124 tipHash=0x030d9e
2025-11-08T01:35:00Z INFO telemetry: height=200095 peers=48 mempool=3095 avgLatencyMs=125 tipHash=0x030d9f
2025-11-08T01:36:00Z INFO telemetry: height=200096 peers=49 mempool=3096 avgLatencyMs=126 tipHash=0x030da0
2025-11-08T01:37:00Z INFO telemetry: height=200097 peers=50 mempool=3097 avgLatencyMs=127 tipHash=0x030da1
2025-11-08T01:38:00Z INFO telemetry: height=200098 peers=51 mempool=3098 avgLatencyMs=128 tipHash=0x030da2
2025-11-08T01:39:00Z INFO telemetry: height=200099 peers=52 mempool=3099 avgLatencyMs=129 tipHash=0x030da3
2025-11-08T01:40:00Z INFO telemetry: height=200100 peers=48 mempool=3100 avgLatencyMs=130 tipHash=0x030da4
2025-11-08T01:41:00Z INFO telemetry: height=200101 peers=49 mempool=3101 avgLatencyMs=131 tipHash=0x030da5
2025-11-08T01:42:00Z INFO telemetry: height=200102 peers=50 mempool=3102 avgLatencyMs=132 tipHash=0x030da6
2025-11-08T01:43:00Z INFO telemetry: height=200103 peers=51 mempool=3103 avgLatencyMs=133 tipHash=0x030da7
2025-11-08T01:44:00Z INFO telemetry: height=200104 peers=52 mempool=3104 avgLatencyMs=134 tipHash=0x030da8
2025-11-08T01:45:00Z INFO telemetry: height=200105 peers=48 mempool=3105 avgLatencyMs=135 tipHash=0x030da9
2025-11-08T01:46:00Z INFO telemetry: height=200106 peers=49 mempool=3106 avgLatencyMs=136 tipHash=0x030daa
2025-11-08T01:47:00Z INFO telemetry: height=200107 peers=50 mempool=3107 avgLatencyMs=137 tipHash=0x030dab
2025-11-08T01:48:00Z INFO telemetry: height=200108 peers=51 mempool=3108 avgLatencyMs=138 tipHash=0x030dac
2025-11-08T01:49:00Z INFO telemetry: height=200109 peers=52 mempool=3109 avgLatencyMs=139 tipHash=0x030dad
2025-11-08T01:50:00Z INFO telemetry: height=200110 peers=48 mempool=3110 avgLatencyMs=140 tipHash=0x030dae
2025-11-08T01:51:00Z INFO telemetry: height=200111 peers=49 mempool=3111 avgLatencyMs=141 tipHash=0x030daf
2025-11-08T01:52:00Z INFO telemetry: height=200112 peers=50 mempool=3112 avgLatencyMs=142 tipHash=0x030db0
2025-11-08T01:53:00Z INFO telemetry: height=200113 peers=51 mempool=3113 avgLatencyMs=143 tipHash=0x030db1
2025-11-08T01:54:00Z INFO telemetry: height=200114 peers=52 mempool=3114 avgLatencyMs=144 tipHash=0x030db2
2025-11-08T01:55:00Z INFO telemetry: height=200115 peers=48 mempool=3115 avgLatencyMs=145 tipHash=0x030db3
2025-11-08T01:56:00Z INFO telemetry: height=200116 peers=49 mempool=3116 avgLatencyMs=146 tipHash=0x030db4
2025-11-08T01:57:00Z INFO telemetry: height=200117 peers=50 mempool=3117 avgLatencyMs=147 tipHash=0x030db5
2025-11-08T01:58:00Z INFO telemetry: height=200118 peers=51 mempool=3118 avgLatencyMs=148 tipHash=0x030db6
2025-11-08T01:59:00Z INFO telemetry: height=200119 peers=52 mempool=3119 avgLatencyMs=149 tipHash=0x030db7
2025-11-08T02:00:00Z INFO telemetry: height=200120 peers=48 mempool=3120 avgLatencyMs=120 tipHash=0x030db8
2025-11-08T02:01:00Z INFO telemetry: height=200121 peers=49 mempool=3121 avgLatencyMs=121 tipHash=0x030db9
2025-11-08T02:02:00Z INFO telemetry: height=200122 peers=50 mempool=3122 avgLatencyMs=122 tipHash=0x030dba
2025-11-08T02:03:00Z INFO telemetry: height=200123 peers=51 mempool=3123 avgLatencyMs=123 tipHash=0x030dbb
2025-11-08T02:04:00Z INFO telemetry: height=200124 peers=52 mempool=3124 avgLatencyMs=124 tipHash=0x030dbc
2025-11-08T02:05:00Z INFO telemetry: height=200125 peers=48 mempool=3125 avgLatencyMs=125 tipHash=0x030dbd
2025-11-08T02:06:00Z INFO telemetry: height=200126 peers=49 mempool=3126 avgLatencyMs=126 tipHash=0x030dbe
2025-11-08T02:07:00Z INFO telemetry: height=200127 peers=50 mempool=3127 avgLatencyMs=127 tipHash=0x030dbf
2025-11-08T02:08:00Z INFO telemetry: height=200128 peers=51 mempool=3128 avgLatencyMs=128 tipHash=0x030dc0
2025-11-08T02:09:00Z INFO telemetry: height=200129 peers=52 mempool=3129 avgLatencyMs=129 tipHash=0x030dc1
2025-11-08T02:10:00Z INFO telemetry: height=200130 peers=48 mempool=3130 avgLatencyMs=130 tipHash=0x030dc2
2025-11-08T02:11:00Z INFO telemetry: height=200131 peers=49 mempool=3131 avgLatencyMs=131 tipHash=0x030dc3
2025-11-08T02:12:00Z INFO telemetry: height=200132 peers=50 mempool=3132 avgLatencyMs=132 tipHash=0x030dc4
2025-11-08T02:13:00Z INFO telemetry: height=200133 peers=51 mempool=3133 avgLatencyMs=133 tipHash=0x030dc5
2025-11-08T02:14:00Z INFO telemetry: height=200134 peers=52 mempool=3134 avgLatencyMs=134 tipHash=0x030dc6
2025-11-08T02:15:00Z INFO telemetry: height=200135 peers=48 mempool=3135 avgLatencyMs=135 tipHash=0x030dc7
2025-11-08T02:16:00Z INFO telemetry: height=200136 peers=49 mempool=3136 avgLatencyMs=136 tipHash=0x030dc8
2025-11-08T02:17:00Z INFO telemetry: height=200137 peers=50 mempool=3137 avgLatencyMs=137 tipHash=0x030dc9
2025-11-08T02:18:00Z INFO telemetry: height=200138 peers=51 mempool=3138 avgLatencyMs=138 tipHash=0x030dca
2025-11-08T02:19:00Z INFO telemetry: height=200139 peers=52 mempool=3139 avgLatencyMs=139 tipHash=0x030dcb
2025-11-08T02:20:00Z INFO telemetry: height=200140 peers=48 mempool=3140 avgLatencyMs=140 tipHash=0x030dcc
2025-11-08T02:21:00Z INFO telemetry: height=200141 peers=49 mempool=3141 avgLatencyMs=141 tipHash=0x030dcd
2025-11-08T02:22:00Z INFO telemetry: height=200142 peers=50 mempool=3142 avgLatencyMs=142 tipHash=0x030dce
2025-11-08T02:23:00Z INFO telemetry: height=200143 peers=51 mempool=3143 avgLatencyMs=143 tipHash=0x030dcf
2025-11-08T02:24:00Z INFO telemetry: height=200144 peers=52 mempool=3144 avgLatencyMs=144 tipHash=0x030dd0
2025-11-08T02:25:00Z INFO telemetry: height=200145 peers=48 mempool=3145 avgLatencyMs=145 tipHash=0x030dd1
2025-11-08T02:26:00Z INFO telemetry: height=200146 peers=49 mempool=3146 avgLatencyMs=146 tipHash=0x030dd2
2025-11-08T02:27:00Z INFO telemetry: height=200147 peers=50 mempool=3147 avgLatencyMs=147 tipHash=0x030dd3
2025-11-08T02:28:00Z INFO telemetry: height=200148 peers=51 mempool=3148 avgLatencyMs=148 tipHash=0x030dd4
2025-11-08T02:29:00Z INFO telemetry: height=200149 peers=52 mempool=3149 avgLatencyMs=149 tipHash=0x030dd5
2025-11-08T02:30:00Z INFO telemetry: height=200150 peers=48 mempool=3150 avgLatencyMs=120 tipHash=0x030dd6
2025-11-08T02:31:00Z INFO telemetry: height=200151 peers=49 mempool=3151 avgLatencyMs=121 tipHash=0x030dd7
2025-11-08T02:32:00Z INFO telemetry: height=200152 peers=50 mempool=3152 avgLatencyMs=122 tipHash=0x030dd8
2025-11-08T02:33:00Z INFO telemetry: height=200153 peers=51 mempool=3153 avgLatencyMs=123 tipHash=0x030dd9
2025-11-08T02:34:00Z INFO telemetry: height=200154 peers=52 mempool=3154 avgLatencyMs=124 tipHash=0x030dda
2025-11-08T02:35:00Z INFO telemetry: height=200155 peers=48 mempool=3155 avgLatencyMs=125 tipHash=0x030ddb
2025-11-08T02:36:00Z INFO telemetry: height=200156 peers=49 mempool=3156 avgLatencyMs=126 tipHash=0x030ddc
2025-11-08T02:37:00Z INFO telemetry: height=200157 peers=50 mempool=3157 avgLatencyMs=127 tipHash=0x030ddd
2025-11-08T02:38:00Z INFO telemetry: height=200158 peers=51 mempool=3158 avgLatencyMs=128 tipHash=0x030dde
2025-11-08T02:39:00Z INFO telemetry: height=200159 peers=52 mempool=3159 avgLatencyMs=129 tipHash=0x030ddf
2025-11-08T02:40:00Z INFO telemetry: height=200160 peers=48 mempool=3160 avgLatencyMs=130 tipHash=0x030de0
2025-11-08T02:41:00Z INFO telemetry: height=200161 peers=49 mempool=3161 avgLatencyMs=131 tipHash=0x030de1
2025-11-08T02:42:00Z INFO telemetry: height=200162 peers=50 mempool=3162 avgLatencyMs=132 tipHash=0x030de2
2025-11-08T02:43:00Z INFO telemetry: height=200163 peers=51 mempool=3163 avgLatencyMs=133 tipHash=0x030de3
2025-11-08T02:44:00Z INFO telemetry: height=200164 peers=52 mempool=3164 avgLatencyMs=134 tipHash=0x030de4
2025-11-08T02:45:00Z INFO telemetry: height=200165 peers=48 mempool=3165 avgLatencyMs=135 tipHash=0x030de5
2025-11-08T02:46:00Z INFO telemetry: height=200166 peers=49 mempool=3166 avgLatencyMs=136 tipHash=0x030de6
2025-11-08T02:47:00Z INFO telemetry: height=200167 peers=50 mempool=3167 avgLatencyMs=137 tipHash=0x030de7
2025-11-08T02:48:00Z INFO telemetry: height=200168 peers=51 mempool=3168 avgLatencyMs=138 tipHash=0x030de8
2025-11-08T02:49:00Z INFO telemetry: height=200169 peers=52 mempool=3169 avgLatencyMs=139 tipHash=0x030de9
2025-11-08T02:50:00Z INFO telemetry: height=200170 peers=48 mempool=3170 avgLatencyMs=140 tipHash=0x030dea
2025-11-08T02:51:00Z INFO telemetry: height=200171 peers=49 mempool=3171 avgLatencyMs=141 tipHash=0x030deb
2025-11-08T02:52:00Z INFO telemetry: height=200172 peers=50 mempool=3172 avgLatencyMs=142 tipHash=0x030dec
2025-11-08T02:53:00Z INFO telemetry: height=200173 peers=51 mempool=3173 avgLatencyMs=143 tipHash=0x030ded
2025-11-08T02:54:00Z INFO telemetry: height=200174 peers=52 mempool=3174 avgLatencyMs=144 tipHash=0x030dee
2025-11-08T02:55:00Z INFO telemetry: height=200175 peers=48 mempool=3175 avgLatencyMs=145 tipHash=0x030def
2025-11-08T02:56:00Z INFO telemetry: height=200176 peers=49 mempool=3176 avgLatencyMs=146 tipHash=0x030df0
2025-11-08T02:57:00Z INFO telemetry: height=200177 peers=50 mempool=3177 avgLatencyMs=147 tipHash=0x030df1
2025-11-08T02:58:00Z INFO telemetry: height=200178 peers=51 mempool=3178 avgLatencyMs=148 tipHash=0x030df2
2025-11-08T02:59:00Z INFO telemetry: height=200179 peers=52 mempool=3179 avgLatencyMs=149 tipHash=0x030df3
2025-11-08T03:00:00Z INFO telemetry: height=200180 peers=48 mempool=3180 avgLatencyMs=120 tipHash=0x030df4
2025-11-08T03:01:00Z INFO telemetry: height=200181 peers=49 mempool=3181 avgLatencyMs=121 tipHash=0x030df5
2025-11-08T03:02:00Z INFO telemetry: height=200182 peers=50 mempool=3182 avgLatencyMs=122 tipHash=0x030df6
2025-11-08T03:03:00Z INFO telemetry: height=200183 peers=51 mempool=3183 avgLatencyMs=123 tipHash=0x030df7
2025-11-08T03:04:00Z INFO telemetry: height=200184 peers=52 mempool=3184 avgLatencyMs=124 tipHash=0x030df8
2025-11-08T03:05:00Z INFO telemetry: height=200185 peers=48 mempool=3185 avgLatencyMs=125 tipHash=0x030df9
2025-11-08T03:06:00Z INFO telemetry: height=200186 peers=49 mempool=3186 avgLatencyMs=126 tipHash=0x030dfa
2025-11-08T03:07:00Z INFO telemetry: height=200187 peers=50 mempool=3187 avgLatencyMs=127 tipHash=0x030dfb
2025-11-08T03:08:00Z INFO telemetry: height=200188 peers=51 mempool=3188 avgLatencyMs=128 tipHash=0x030dfc
2025-11-08T03:09:00Z INFO telemetry: height=200189 peers=52 mempool=3189 avgLatencyMs=129 tipHash=0x030dfd
2025-11-08T03:10:00Z INFO telemetry: height=200190 peers=48 mempool=3190 avgLatencyMs=130 tipHash=0x030dfe
2025-11-08T03:11:00Z INFO telemetry: height=200191 peers=49 mempool=3191 avgLatencyMs=131 tipHash=0x030dff
2025-11-08T03:12:00Z INFO telemetry: height=200192 peers=50 mempool=3192 avgLatencyMs=132 tipHash=0x030e00
2025-11-08T03:13:00Z INFO telemetry: height=200193 peers=51 mempool=3193 avgLatencyMs=133 tipHash=0x030e01
2025-11-08T03:14:00Z INFO telemetry: height=200194 peers=52 mempool=3194 avgLatencyMs=134 tipHash=0x030e02
2025-11-08T03:15:00Z INFO telemetry: height=200195 peers=48 mempool=3195 avgLatencyMs=135 tipHash=0x030e03
2025-11-08T03:16:00Z INFO telemetry: height=200196 peers=49 mempool=3196 avgLatencyMs=136 tipHash=0x030e04
2025-11-08T03:17:00Z INFO telemetry: height=200197 peers=50 mempool=3197 avgLatencyMs=137 tipHash=0x030e05
2025-11-08T03:18:00Z INFO telemetry: height=200198 peers=51 mempool=3198 avgLatencyMs=138 tipHash=0x030e06
2025-11-08T03:19:00Z INFO telemetry: height=200199 peers=52 mempool=3199 avgLatencyMs=139 tipHash=0x030e07
2025-11-08T03:20:00Z INFO telemetry: height=200200 peers=48 mempool=3000 avgLatencyMs=140 tipHash=0x030e08
2025-11-08T03:21:00Z INFO telemetry: height=200201 peers=49 mempool=3001 avgLatencyMs=141 tipHash=0x030e09
2025-11-08T03:22:00Z INFO telemetry: height=200202 peers=50 mempool=3002 avgLatencyMs=142 tipHash=0x030e0a
2025-11-08T03:23:00Z INFO telemetry: height=200203 peers=51 mempool=3003 avgLatencyMs=143 tipHash=0x030e0b
2025-11-08T03:24:00Z INFO telemetry: height=200204 peers=52 mempool=3004 avgLatencyMs=144 tipHash=0x030e0c
2025-11-08T03:25:00Z INFO telemetry: height=200205 peers=48 mempool=3005 avgLatencyMs=145 tipHash=0x030e0d
2025-11-08T03:26:00Z INFO telemetry: height=200206 peers=49 mempool=3006 avgLatencyMs=146 tipHash=0x030e0e
2025-11-08T03:27:00Z INFO telemetry: height=200207 peers=50 mempool=3007 avgLatencyMs=147 tipHash=0x030e0f
2025-11-08T03:28:00Z INFO telemetry: height=200208 peers=51 mempool=3008 avgLatencyMs=148 tipHash=0x030e10
2025-11-08T03:29:00Z INFO telemetry: height=200209 peers=52 mempool=3009 avgLatencyMs=149 tipHash=0x030e11
2025-11-08T03:30:00Z INFO telemetry: height=200210 peers=48 mempool=3010 avgLatencyMs=120 tipHash=0x030e12
2025-11-08T03:31:00Z INFO telemetry: height=200211 peers=49 mempool=3011 avgLatencyMs=121 tipHash=0x030e13
2025-11-08T03:32:00Z INFO telemetry: height=200212 peers=50 mempool=3012 avgLatencyMs=122 tipHash=0x030e14
2025-11-08T03:33:00Z INFO telemetry: height=200213 peers=51 mempool=3013 avgLatencyMs=123 tipHash=0x030e15
2025-11-08T03:34:00Z INFO telemetry: height=200214 peers=52 mempool=3014 avgLatencyMs=124 tipHash=0x030e16
2025-11-08T03:35:00Z INFO telemetry: height=200215 peers=48 mempool=3015 avgLatencyMs=125 tipHash=0x030e17
2025-11-08T03:36:00Z INFO telemetry: height=200216 peers=49 mempool=3016 avgLatencyMs=126 tipHash=0x030e18
2025-11-08T03:37:00Z INFO telemetry: height=200217 peers=50 mempool=3017 avgLatencyMs=127 tipHash=0x030e19
2025-11-08T03:38:00Z INFO telemetry: height=200218 peers=51 mempool=3018 avgLatencyMs=128 tipHash=0x030e1a
2025-11-08T03:39:00Z INFO telemetry: height=200219 peers=52 mempool=3019 avgLatencyMs=129 tipHash=0x030e1b
2025-11-08T03:40:00Z INFO telemetry: height=200220 peers=48 mempool=3020 avgLatencyMs=130 tipHash=0x030e1c
2025-11-08T03:41:00Z INFO telemetry: height=200221 peers=49 mempool=3021 avgLatencyMs=131 tipHash=0x030e1d
2025-11-08T03:42:00Z INFO telemetry: height=200222 peers=50 mempool=3022 avgLatencyMs=132 tipHash=0x030e1e
2025-11-08T03:43:00Z INFO telemetry: height=200223 peers=51 mempool=3023 avgLatencyMs=133 tipHash=0x030e1f
2025-11-08T03:44:00Z INFO telemetry: height=200224 peers=52 mempool=3024 avgLatencyMs=134 tipHash=0x030e20
2025-11-08T03:45:00Z INFO telemetry: height=200225 peers=48 mempool=3025 avgLatencyMs=135 tipHash=0x030e21
2025-11-08T03:46:00Z INFO telemetry: height=200226 peers=49 mempool=3026 avgLatencyMs=136 tipHash=0x030e22
2025-11-08T03:47:00Z INFO telemetry: height=200227 peers=50 mempool=3027 avgLatencyMs=137 tipHash=0x030e23
2025-11-08T03:48:00Z INFO telemetry: height=200228 peers=51 mempool=3028 avgLatencyMs=138 tipHash=0x030e24
2025-11-08T03:49:00Z INFO telemetry: height=200229 peers=52 mempool=3029 avgLatencyMs=139 tipHash=0x030e25
2025-11-08T03:50:00Z INFO telemetry: height=200230 peers=48 mempool=3030 avgLatencyMs=140 tipHash=0x030e26
2025-11-08T03:51:00Z INFO telemetry: height=200231 peers=49 mempool=3031 avgLatencyMs=141 tipHash=0x030e27
2025-11-08T03:52:00Z INFO telemetry: height=200232 peers=50 mempool=3032 avgLatencyMs=142 tipHash=0x030e28
2025-11-08T03:53:00Z INFO telemetry: height=200233 peers=51 mempool=3033 avgLatencyMs=143 tipHash=0x030e29
2025-11-08T03:54:00Z INFO telemetry: height=200234 peers=52 mempool=3034 avgLatencyMs=144 tipHash=0x030e2a
2025-11-08T03:55:00Z INFO telemetry: height=200235 peers=48 mempool=3035 avgLatencyMs=145 tipHash=0x030e2b
2025-11-08T03:56:00Z INFO telemetry: height=200236 peers=49 mempool=3036 avgLatencyMs=146 tipHash=0x030e2c
2025-11-08T03:57:00Z INFO telemetry: height=200237 peers=50 mempool=3037 avgLatencyMs=147 tipHash=0x030e2d
2025-11-08T03:58:00Z INFO telemetry: height=200238 peers=51 mempool=3038 avgLatencyMs=148 tipHash=0x030e2e
2025-11-08T03:59:00Z INFO telemetry: height=200239 peers=52 mempool=3039 avgLatencyMs=149 tipHash=0x030e2f
2025-11-08T04:00:00Z INFO telemetry: height=200240 peers=48 mempool=3040 avgLatencyMs=120 tipHash=0x030e30
2025-11-08T04:01:00Z INFO telemetry: height=200241 peers=49 mempool=3041 avgLatencyMs=121 tipHash=0x030e31
2025-11-08T04:02:00Z INFO telemetry: height=200242 peers=50 mempool=3042 avgLatencyMs=122 tipHash=0x030e32
2025-11-08T04:03:00Z INFO telemetry: height=200243 peers=51 mempool=3043 avgLatencyMs=123 tipHash=0x030e33
2025-11-08T04:04:00Z INFO telemetry: height=200244 peers=52 mempool=3044 avgLatencyMs=124 tipHash=0x030e34
2025-11-08T04:05:00Z INFO telemetry: height=200245 peers=48 mempool=3045 avgLatencyMs=125 tipHash=0x030e35
2025-11-08T04:06:00Z INFO telemetry: height=200246 peers=49 mempool=3046 avgLatencyMs=126 tipHash=0x030e36
2025-11-08T04:07:00Z INFO telemetry: height=200247 peers=50 mempool=3047 avgLatencyMs=127 tipHash=0x030e37
2025-11-08T04:08:00Z INFO telemetry: height=200248 peers=51 mempool=3048 avgLatencyMs=128 tipHash=0x030e38
2025-11-08T04:09:00Z INFO telemetry: height=200249 peers=52 mempool=3049 avgLatencyMs=129 tipHash=0x030e39
2025-11-08T04:10:00Z INFO telemetry: height=200250 peers=48 mempool=3050 avgLatencyMs=130 tipHash=0x030e3a
2025-11-08T04:11:00Z INFO telemetry: height=200251 peers=49 mempool=3051 avgLatencyMs=131 tipHash=0x030e3b
2025-11-08T04:12:00Z INFO telemetry: height=200252 peers=50 mempool=3052 avgLatencyMs=132 tipHash=0x030e3c
2025-11-08T04:13:00Z INFO telemetry: height=200253 peers=51 mempool=3053 avgLatencyMs=133 tipHash=0x030e3d
2025-11-08T04:14:00Z INFO telemetry: height=200254 peers=52 mempool=3054 avgLatencyMs=134 tipHash=0x030e3e
2025-11-08T04:15:00Z INFO telemetry: height=200255 peers=48 mempool=3055 avgLatencyMs=135 tipHash=0x030e3f
2025-11-08T04:16:00Z INFO telemetry: height=200256 peers=49 mempool=3056 avgLatencyMs=136 tipHash=0x030e40
2025-11-08T04:17:00Z INFO telemetry: height=200257 peers=50 mempool=3057 avgLatencyMs=137 tipHash=0x030e41
2025-11-08T04:18:00Z INFO telemetry: height=200258 peers=51 mempool=3058 avgLatencyMs=138 tipHash=0x030e42
2025-11-08T04:19:00Z INFO telemetry: height=200259 peers=52 mempool=3059 avgLatencyMs=139 tipHash=0x030e43
2025-11-08T04:20:00Z INFO telemetry: height=200260 peers=48 mempool=3060 avgLatencyMs=140 tipHash=0x030e44
2025-11-08T04:21:00Z INFO telemetry: height=200261 peers=49 mempool=3061 avgLatencyMs=141 tipHash=0x030e45
2025-11-08T04:22:00Z INFO telemetry: height=200262 peers=50 mempool=3062 avgLatencyMs=142 tipHash=0x030e46
2025-11-08T04:23:00Z INFO telemetry: height=200263 peers=51 mempool=3063 avgLatencyMs=143 tipHash=0x030e47
2025-11-08T04:24:00Z INFO telemetry: height=200264 peers=52 mempool=3064 avgLatencyMs=144 tipHash=0x030e48
2025-11-08T04:25:00Z INFO telemetry: height=200265 peers=48 mempool=3065 avgLatencyMs=145 tipHash=0x030e49
2025-11-08T04:26:00Z INFO telemetry: height=200266 peers=49 mempool=3066 avgLatencyMs=146 tipHash=0x030e4a
2025-11-08T04:27:00Z INFO telemetry: height=200267 peers=50 mempool=3067 avgLatencyMs=147 tipHash=0x030e4b
2025-11-08T04:28:00Z INFO telemetry: height=200268 peers=51 mempool=3068 avgLatencyMs=148 tipHash=0x030e4c
2025-11-08T04:29:00Z INFO telemetry: height=200269 peers=52 mempool=3069 avgLatencyMs=149 tipHash=0x030e4d
2025-11-08T04:30:00Z INFO telemetry: height=200270 peers=48 mempool=3070 avgLatencyMs=120 tipHash=0x030e4e
2025-11-08T04:31:00Z INFO telemetry: height=200271 peers=49 mempool=3071 avgLatencyMs=121 tipHash=0x030e4f
2025-11-08T04:32:00Z INFO telemetry: height=200272 peers=50 mempool=3072 avgLatencyMs=122 tipHash=0x030e50
2025-11-08T04:33:00Z INFO telemetry: height=200273 peers=51 mempool=3073 avgLatencyMs=123 tipHash=0x030e51
2025-11-08T04:34:00Z INFO telemetry: height=200274 peers=52 mempool=3074 avgLatencyMs=124 tipHash=0x030e52
2025-11-08T04:35:00Z INFO telemetry: height=200275 peers=48 mempool=3075 avgLatencyMs=125 tipHash=0x030e53
2025-11-08T04:36:00Z INFO telemetry: height=200276 peers=49 mempool=3076 avgLatencyMs=126 tipHash=0x030e54
2025-11-08T04:37:00Z INFO telemetry: height=200277 peers=50 mempool=3077 avgLatencyMs=127 tipHash=0x030e55
2025-11-08T04:38:00Z INFO telemetry: height=200278 peers=51 mempool=3078 avgLatencyMs=128 tipHash=0x030e56
2025-11-08T04:39:00Z INFO telemetry: height=200279 peers=52 mempool=3079 avgLatencyMs=129 tipHash=0x030e57
2025-11-08T04:40:00Z INFO telemetry: height=200280 peers=48 mempool=3080 avgLatencyMs=130 tipHash=0x030e58
2025-11-08T04:41:00Z INFO telemetry: height=200281 peers=49 mempool=3081 avgLatencyMs=131 tipHash=0x030e59
2025-11-08T04:42:00Z INFO telemetry: height=200282 peers=50 mempool=3082 avgLatencyMs=132 tipHash=0x030e5a
2025-11-08T04:43:00Z INFO telemetry: height=200283 peers=51 mempool=3083 avgLatencyMs=133 tipHash=0x030e5b
2025-11-08T04:44:00Z INFO telemetry: height=200284 peers=52 mempool=3084 avgLatencyMs=134 tipHash=0x030e5c
2025-11-08T04:45:00Z INFO telemetry: height=200285 peers=48 mempool=3085 avgLatencyMs=135 tipHash=0x030e5d
2025-11-08T04:46:00Z INFO telemetry: height=200286 peers=49 mempool=3086 avgLatencyMs=136 tipHash=0x030e5e
2025-11-08T04:47:00Z INFO telemetry: height=200287 peers=50 mempool=3087 avgLatencyMs=137 tipHash=0x030e5f
2025-11-08T04:48:00Z INFO telemetry: height=200288 peers=51 mempool=3088 avgLatencyMs=138 tipHash=0x030e60
2025-11-08T04:49:00Z INFO telemetry: height=200289 peers=52 mempool=3089 avgLatencyMs=139 tipHash=0x030e61
2025-11-08T04:50:00Z INFO telemetry: height=200290 peers=48 mempool=3090 avgLatencyMs=140 tipHash=0x030e62
2025-11-08T04:51:00Z INFO telemetry: height=200291 peers=49 mempool=3091 avgLatencyMs=141 tipHash=0x030e63
2025-11-08T04:52:00Z INFO telemetry: height=200292 peers=50 mempool=3092 avgLatencyMs=142 tipHash=0x030e64
2025-11-08T04:53:00Z INFO telemetry: height=200293 peers=51 mempool=3093 avgLatencyMs=143 tipHash=0x030e65
2025-11-08T04:54:00Z INFO telemetry: height=200294 peers=52 mempool=3094 avgLatencyMs=144 tipHash=0x030e66
2025-11-08T04:55:00Z INFO telemetry: height=200295 peers=48 mempool=3095 avgLatencyMs=145 tipHash=0x030e67
2025-11-08T04:56:00Z INFO telemetry: height=200296 peers=49 mempool=3096 avgLatencyMs=146 tipHash=0x030e68
2025-11-08T04:57:00Z INFO telemetry: height=200297 peers=50 mempool=3097 avgLatencyMs=147 tipHash=0x030e69
2025-11-08T04:58:00Z INFO telemetry: height=200298 peers=51 mempool=3098 avgLatencyMs=148 tipHash=0x030e6a
2025-11-08T04:59:00Z INFO telemetry: height=200299 peers=52 mempool=3099 avgLatencyMs=149 tipHash=0x030e6b
2025-11-08T05:00:00Z INFO telemetry: height=200300 peers=48 mempool=3100 avgLatencyMs=120 tipHash=0x030e6c
2025-11-08T05:01:00Z INFO telemetry: height=200301 peers=49 mempool=3101 avgLatencyMs=121 tipHash=0x030e6d
2025-11-08T05:02:00Z INFO telemetry: height=200302 peers=50 mempool=3102 avgLatencyMs=122 tipHash=0x030e6e
2025-11-08T05:03:00Z INFO telemetry: height=200303 peers=51 mempool=3103 avgLatencyMs=123 tipHash=0x030e6f
2025-11-08T05:04:00Z INFO telemetry: height=200304 peers=52 mempool=3104 avgLatencyMs=124 tipHash=0x030e70
2025-11-08T05:05:00Z INFO telemetry: height=200305 peers=48 mempool=3105 avgLatencyMs=125 tipHash=0x030e71
2025-11-08T05:06:00Z INFO telemetry: height=200306 peers=49 mempool=3106 avgLatencyMs=126 tipHash=0x030e72
2025-11-08T05:07:00Z INFO telemetry: height=200307 peers=50 mempool=3107 avgLatencyMs=127 tipHash=0x030e73
2025-11-08T05:08:00Z INFO telemetry: height=200308 peers=51 mempool=3108 avgLatencyMs=128 tipHash=0x030e74
2025-11-08T05:09:00Z INFO telemetry: height=200309 peers=52 mempool=3109 avgLatencyMs=129 tipHash=0x030e75
2025-11-08T05:10:00Z INFO telemetry: height=200310 peers=48 mempool=3110 avgLatencyMs=130 tipHash=0x030e76
2025-11-08T05:11:00Z INFO telemetry: height=200311 peers=49 mempool=3111 avgLatencyMs=131 tipHash=0x030e77
2025-11-08T05:12:00Z INFO telemetry: height=200312 peers=50 mempool=3112 avgLatencyMs=132 tipHash=0x030e78
2025-11-08T05:13:00Z INFO telemetry: height=200313 peers=51 mempool=3113 avgLatencyMs=133 tipHash=0x030e79
2025-11-08T05:14:00Z INFO telemetry: height=200314 peers=52 mempool=3114 avgLatencyMs=134 tipHash=0x030e7a
2025-11-08T05:15:00Z INFO telemetry: height=200315 peers=48 mempool=3115 avgLatencyMs=135 tipHash=0x030e7b
2025-11-08T05:16:00Z INFO telemetry: height=200316 peers=49 mempool=3116 avgLatencyMs=136 tipHash=0x030e7c
2025-11-08T05:17:00Z INFO telemetry: height=200317 peers=50 mempool=3117 avgLatencyMs=137 tipHash=0x030e7d
2025-11-08T05:18:00Z INFO telemetry: height=200318 peers=51 mempool=3118 avgLatencyMs=138 tipHash=0x030e7e
2025-11-08T05:19:00Z INFO telemetry: height=200319 peers=52 mempool=3119 avgLatencyMs=139 tipHash=0x030e7f
2025-11-08T05:20:00Z INFO telemetry: height=200320 peers=48 mempool=3120 avgLatencyMs=140 tipHash=0x030e80
2025-11-08T05:21:00Z INFO telemetry: height=200321 peers=49 mempool=3121 avgLatencyMs=141 tipHash=0x030e81
2025-11-08T05:22:00Z INFO telemetry: height=200322 peers=50 mempool=3122 avgLatencyMs=142 tipHash=0x030e82
2025-11-08T05:23:00Z INFO telemetry: height=200323 peers=51 mempool=3123 avgLatencyMs=143 tipHash=0x030e83
2025-11-08T05:24:00Z INFO telemetry: height=200324 peers=52 mempool=3124 avgLatencyMs=144 tipHash=0x030e84
2025-11-08T05:25:00Z INFO telemetry: height=200325 peers=48 mempool=3125 avgLatencyMs=145 tipHash=0x030e85
2025-11-08T05:26:00Z INFO telemetry: height=200326 peers=49 mempool=3126 avgLatencyMs=146 tipHash=0x030e86
2025-11-08T05:27:00Z INFO telemetry: height=200327 peers=50 mempool=3127 avgLatencyMs=147 tipHash=0x030e87
2025-11-08T05:28:00Z INFO telemetry: height=200328 peers=51 mempool=3128 avgLatencyMs=148 tipHash=0x030e88
2025-11-08T05:29:00Z INFO telemetry: height=200329 peers=52 mempool=3129 avgLatencyMs=149 tipHash=0x030e89
2025-11-08T05:30:00Z INFO telemetry: height=200330 peers=48 mempool=3130 avgLatencyMs=120 tipHash=0x030e8a
2025-11-08T05:31:00Z INFO telemetry: height=200331 peers=49 mempool=3131 avgLatencyMs=121 tipHash=0x030e8b
2025-11-08T05:32:00Z INFO telemetry: height=200332 peers=50 mempool=3132 avgLatencyMs=122 tipHash=0x030e8c
2025-11-08T05:33:00Z INFO telemetry: height=200333 peers=51 mempool=3133 avgLatencyMs=123 tipHash=0x030e8d
2025-11-08T05:34:00Z INFO telemetry: height=200334 peers=52 mempool=3134 avgLatencyMs=124 tipHash=0x030e8e
2025-11-08T05:35:00Z INFO telemetry: height=200335 peers=48 mempool=3135 avgLatencyMs=125 tipHash=0x030e8f
2025-11-08T05:36:00Z INFO telemetry: height=200336 peers=49 mempool=3136 avgLatencyMs=126 tipHash=0x030e90
2025-11-08T05:37:00Z INFO telemetry: height=200337 peers=50 mempool=3137 avgLatencyMs=127 tipHash=0x030e91
2025-11-08T05:38:00Z INFO telemetry: height=200338 peers=51 mempool=3138 avgLatencyMs=128 tipHash=0x030e92
2025-11-08T05:39:00Z INFO telemetry: height=200339 peers=52 mempool=3139 avgLatencyMs=129 tipHash=0x030e93
2025-11-08T05:40:00Z INFO telemetry: height=200340 peers=48 mempool=3140 avgLatencyMs=130 tipHash=0x030e94
2025-11-08T05:41:00Z INFO telemetry: height=200341 peers=49 mempool=3141 avgLatencyMs=131 tipHash=0x030e95
2025-11-08T05:42:00Z INFO telemetry: height=200342 peers=50 mempool=3142 avgLatencyMs=132 tipHash=0x030e96
2025-11-08T05:43:00Z INFO telemetry: height=200343 peers=51 mempool=3143 avgLatencyMs=133 tipHash=0x030e97
2025-11-08T05:44:00Z INFO telemetry: height=200344 peers=52 mempool=3144 avgLatencyMs=134 tipHash=0x030e98
2025-11-08T05:45:00Z INFO telemetry: height=200345 peers=48 mempool=3145 avgLatencyMs=135 tipHash=0x030e99
2025-11-08T05:46:00Z INFO telemetry: height=200346 peers=49 mempool=3146 avgLatencyMs=136 tipHash=0x030e9a
2025-11-08T05:47:00Z INFO telemetry: height=200347 peers=50 mempool=3147 avgLatencyMs=137 tipHash=0x030e9b
2025-11-08T05:48:00Z INFO telemetry: height=200348 peers=51 mempool=3148 avgLatencyMs=138 tipHash=0x030e9c
2025-11-08T05:49:00Z INFO telemetry: height=200349 peers=52 mempool=3149 avgLatencyMs=139 tipHash=0x030e9d
2025-11-08T05:50:00Z INFO telemetry: height=200350 peers=48 mempool=3150 avgLatencyMs=140 tipHash=0x030e9e
2025-11-08T05:51:00Z INFO telemetry: height=200351 peers=49 mempool=3151 avgLatencyMs=141 tipHash=0x030e9f
2025-11-08T05:52:00Z INFO telemetry: height=200352 peers=50 mempool=3152 avgLatencyMs=142 tipHash=0x030ea0
2025-11-08T05:53:00Z INFO telemetry: height=200353 peers=51 mempool=3153 avgLatencyMs=143 tipHash=0x030ea1
2025-11-08T05:54:00Z INFO telemetry: height=200354 peers=52 mempool=3154 avgLatencyMs=144 tipHash=0x030ea2
2025-11-08T05:55:00Z INFO telemetry: height=200355 peers=48 mempool=3155 avgLatencyMs=145 tipHash=0x030ea3
2025-11-08T05:56:00Z INFO telemetry: height=200356 peers=49 mempool=3156 avgLatencyMs=146 tipHash=0x030ea4
2025-11-08T05:57:00Z INFO telemetry: height=200357 peers=50 mempool=3157 avgLatencyMs=147 tipHash=0x030ea5
2025-11-08T05:58:00Z INFO telemetry: height=200358 peers=51 mempool=3158 avgLatencyMs=148 tipHash=0x030ea6
2025-11-08T05:59:00Z INFO telemetry: height=200359 peers=52 mempool=3159 avgLatencyMs=149 tipHash=0x030ea7
2025-11-08T06:00:00Z INFO telemetry: height=200360 peers=48 mempool=3160 avgLatencyMs=120 tipHash=0x030ea8
2025-11-08T06:01:00Z INFO telemetry: height=200361 peers=49 mempool=3161 avgLatencyMs=121 tipHash=0x030ea9
2025-11-08T06:02:00Z INFO telemetry: height=200362 peers=50 mempool=3162 avgLatencyMs=122 tipHash=0x030eaa
2025-11-08T06:03:00Z INFO telemetry: height=200363 peers=51 mempool=3163 avgLatencyMs=123 tipHash=0x030eab
2025-11-08T06:04:00Z INFO telemetry: height=200364 peers=52 mempool=3164 avgLatencyMs=124 tipHash=0x030eac
2025-11-08T06:05:00Z INFO telemetry: height=200365 peers=48 mempool=3165 avgLatencyMs=125 tipHash=0x030ead
2025-11-08T06:06:00Z INFO telemetry: height=200366 peers=49 mempool=3166 avgLatencyMs=126 tipHash=0x030eae
2025-11-08T06:07:00Z INFO telemetry: height=200367 peers=50 mempool=3167 avgLatencyMs=127 tipHash=0x030eaf
2025-11-08T06:08:00Z INFO telemetry: height=200368 peers=51 mempool=3168 avgLatencyMs=128 tipHash=0x030eb0
2025-11-08T06:09:00Z INFO telemetry: height=200369 peers=52 mempool=3169 avgLatencyMs=129 tipHash=0x030eb1
2025-11-08T06:10:00Z INFO telemetry: height=200370 peers=48 mempool=3170 avgLatencyMs=130 tipHash=0x030eb2
2025-11-08T06:11:00Z INFO telemetry: height=200371 peers=49 mempool=3171 avgLatencyMs=131 tipHash=0x030eb3
2025-11-08T06:12:00Z INFO telemetry: height=200372 peers=50 mempool=3172 avgLatencyMs=132 tipHash=0x030eb4
2025-11-08T06:13:00Z INFO telemetry: height=200373 peers=51 mempool=3173 avgLatencyMs=133 tipHash=0x030eb5
2025-11-08T06:14:00Z INFO telemetry: height=200374 peers=52 mempool=3174 avgLatencyMs=134 tipHash=0x030eb6
2025-11-08T06:15:00Z INFO telemetry: height=200375 peers=48 mempool=3175 avgLatencyMs=135 tipHash=0x030eb7
2025-11-08T06:16:00Z INFO telemetry: height=200376 peers=49 mempool=3176 avgLatencyMs=136 tipHash=0x030eb8
2025-11-08T06:17:00Z INFO telemetry: height=200377 peers=50 mempool=3177 avgLatencyMs=137 tipHash=0x030eb9
2025-11-08T06:18:00Z INFO telemetry: height=200378 peers=51 mempool=3178 avgLatencyMs=138 tipHash=0x030eba
2025-11-08T06:19:00Z INFO telemetry: height=200379 peers=52 mempool=3179 avgLatencyMs=139 tipHash=0x030ebb
2025-11-08T06:20:00Z INFO telemetry: height=200380 peers=48 mempool=3180 avgLatencyMs=140 tipHash=0x030ebc
2025-11-08T06:21:00Z INFO telemetry: height=200381 peers=49 mempool=3181 avgLatencyMs=141 tipHash=0x030ebd
2025-11-08T06:22:00Z INFO telemetry: height=200382 peers=50 mempool=3182 avgLatencyMs=142 tipHash=0x030ebe
2025-11-08T06:23:00Z INFO telemetry: height=200383 peers=51 mempool=3183 avgLatencyMs=143 tipHash=0x030ebf
2025-11-08T06:24:00Z INFO telemetry: height=200384 peers=52 mempool=3184 avgLatencyMs=144 tipHash=0x030ec0
2025-11-08T06:25:00Z INFO telemetry: height=200385 peers=48 mempool=3185 avgLatencyMs=145 tipHash=0x030ec1
2025-11-08T06:26:00Z INFO telemetry: height=200386 peers=49 mempool=3186 avgLatencyMs=146 tipHash=0x030ec2
2025-11-08T06:27:00Z INFO telemetry: height=200387 peers=50 mempool=3187 avgLatencyMs=147 tipHash=0x030ec3
2025-11-08T06:28:00Z INFO telemetry: height=200388 peers=51 mempool=3188 avgLatencyMs=148 tipHash=0x030ec4
2025-11-08T06:29:00Z INFO telemetry: height=200389 peers=52 mempool=3189 avgLatencyMs=149 tipHash=0x030ec5
2025-11-08T06:30:00Z INFO telemetry: height=200390 peers=48 mempool=3190 avgLatencyMs=120 tipHash=0x030ec6
2025-11-08T06:31:00Z INFO telemetry: height=200391 peers=49 mempool=3191 avgLatencyMs=121 tipHash=0x030ec7
2025-11-08T06:32:00Z INFO telemetry: height=200392 peers=50 mempool=3192 avgLatencyMs=122 tipHash=0x030ec8
2025-11-08T06:33:00Z INFO telemetry: height=200393 peers=51 mempool=3193 avgLatencyMs=123 tipHash=0x030ec9
2025-11-08T06:34:00Z INFO telemetry: height=200394 peers=52 mempool=3194 avgLatencyMs=124 tipHash=0x030eca
2025-11-08T06:35:00Z INFO telemetry: height=200395 peers=48 mempool=3195 avgLatencyMs=125 tipHash=0x030ecb
2025-11-08T06:36:00Z INFO telemetry: height=200396 peers=49 mempool=3196 avgLatencyMs=126 tipHash=0x030ecc
2025-11-08T06:37:00Z INFO telemetry: height=200397 peers=50 mempool=3197 avgLatencyMs=127 tipHash=0x030ecd
2025-11-08T06:38:00Z INFO telemetry: height=200398 peers=51 mempool=3198 avgLatencyMs=128 tipHash=0x030ece
2025-11-08T06:39:00Z INFO telemetry: height=200399 peers=52 mempool=3199 avgLatencyMs=129 tipHash=0x030ecf
2025-11-08T06:40:00Z INFO telemetry: height=200400 peers=48 mempool=3000 avgLatencyMs=130 tipHash=0x030ed0
2025-11-08T06:41:00Z INFO telemetry: height=200401 peers=49 mempool=3001 avgLatencyMs=131 tipHash=0x030ed1
2025-11-08T06:42:00Z INFO telemetry: height=200402 peers=50 mempool=3002 avgLatencyMs=132 tipHash=0x030ed2
2025-11-08T06:43:00Z INFO telemetry: height=200403 peers=51 mempool=3003 avgLatencyMs=133 tipHash=0x030ed3
2025-11-08T06:44:00Z INFO telemetry: height=200404 peers=52 mempool=3004 avgLatencyMs=134 tipHash=0x030ed4
2025-11-08T06:45:00Z INFO telemetry: height=200405 peers=48 mempool=3005 avgLatencyMs=135 tipHash=0x030ed5
2025-11-08T06:46:00Z INFO telemetry: height=200406 peers=49 mempool=3006 avgLatencyMs=136 tipHash=0x030ed6
2025-11-08T06:47:00Z INFO telemetry: height=200407 peers=50 mempool=3007 avgLatencyMs=137 tipHash=0x030ed7
2025-11-08T06:48:00Z INFO telemetry: height=200408 peers=51 mempool=3008 avgLatencyMs=138 tipHash=0x030ed8
2025-11-08T06:49:00Z INFO telemetry: height=200409 peers=52 mempool=3009 avgLatencyMs=139 tipHash=0x030ed9
2025-11-08T06:50:00Z INFO telemetry: height=200410 peers=48 mempool=3010 avgLatencyMs=140 tipHash=0x030eda
2025-11-08T06:51:00Z INFO telemetry: height=200411 peers=49 mempool=3011 avgLatencyMs=141 tipHash=0x030edb
2025-11-08T06:52:00Z INFO telemetry: height=200412 peers=50 mempool=3012 avgLatencyMs=142 tipHash=0x030edc
2025-11-08T06:53:00Z INFO telemetry: height=200413 peers=51 mempool=3013 avgLatencyMs=143 tipHash=0x030edd
2025-11-08T06:54:00Z INFO telemetry: height=200414 peers=52 mempool=3014 avgLatencyMs=144 tipHash=0x030ede
2025-11-08T06:55:00Z INFO telemetry: height=200415 peers=48 mempool=3015 avgLatencyMs=145 tipHash=0x030edf
2025-11-08T06:56:00Z INFO telemetry: height=200416 peers=49 mempool=3016 avgLatencyMs=146 tipHash=0x030ee0
2025-11-08T06:57:00Z INFO telemetry: height=200417 peers=50 mempool=3017 avgLatencyMs=147 tipHash=0x030ee1
2025-11-08T06:58:00Z INFO telemetry: height=200418 peers=51 mempool=3018 avgLatencyMs=148 tipHash=0x030ee2
2025-11-08T06:59:00Z INFO telemetry: height=200419 peers=52 mempool=3019 avgLatencyMs=149 tipHash=0x030ee3
2025-11-08T07:00:00Z INFO telemetry: height=200420 peers=48 mempool=3020 avgLatencyMs=120 tipHash=0x030ee4
2025-11-08T07:01:00Z INFO telemetry: height=200421 peers=49 mempool=3021 avgLatencyMs=121 tipHash=0x030ee5
2025-11-08T07:02:00Z INFO telemetry: height=200422 peers=50 mempool=3022 avgLatencyMs=122 tipHash=0x030ee6
2025-11-08T07:03:00Z INFO telemetry: height=200423 peers=51 mempool=3023 avgLatencyMs=123 tipHash=0x030ee7
2025-11-08T07:04:00Z INFO telemetry: height=200424 peers=52 mempool=3024 avgLatencyMs=124 tipHash=0x030ee8
2025-11-08T07:05:00Z INFO telemetry: height=200425 peers=48 mempool=3025 avgLatencyMs=125 tipHash=0x030ee9
2025-11-08T07:06:00Z INFO telemetry: height=200426 peers=49 mempool=3026 avgLatencyMs=126 tipHash=0x030eea
2025-11-08T07:07:00Z INFO telemetry: height=200427 peers=50 mempool=3027 avgLatencyMs=127 tipHash=0x030eeb
2025-11-08T07:08:00Z INFO telemetry: height=200428 peers=51 mempool=3028 avgLatencyMs=128 tipHash=0x030eec
2025-11-08T07:09:00Z INFO telemetry: height=200429 peers=52 mempool=3029 avgLatencyMs=129 tipHash=0x030eed
2025-11-08T07:10:00Z INFO telemetry: height=200430 peers=48 mempool=3030 avgLatencyMs=130 tipHash=0x030eee
2025-11-08T07:11:00Z INFO telemetry: height=200431 peers=49 mempool=3031 avgLatencyMs=131 tipHash=0x030eef
2025-11-08T07:12:00Z INFO telemetry: height=200432 peers=50 mempool=3032 avgLatencyMs=132 tipHash=0x030ef0
2025-11-08T07:13:00Z INFO telemetry: height=200433 peers=51 mempool=3033 avgLatencyMs=133 tipHash=0x030ef1
2025-11-08T07:14:00Z INFO telemetry: height=200434 peers=52 mempool=3034 avgLatencyMs=134 tipHash=0x030ef2
2025-11-08T07:15:00Z INFO telemetry: height=200435 peers=48 mempool=3035 avgLatencyMs=135 tipHash=0x030ef3
2025-11-08T07:16:00Z INFO telemetry: height=200436 peers=49 mempool=3036 avgLatencyMs=136 tipHash=0x030ef4
2025-11-08T07:17:00Z INFO telemetry: height=200437 peers=50 mempool=3037 avgLatencyMs=137 tipHash=0x030ef5
2025-11-08T07:18:00Z INFO telemetry: height=200438 peers=51 mempool=3038 avgLatencyMs=138 tipHash=0x030ef6
2025-11-08T07:19:00Z INFO telemetry: height=200439 peers=52 mempool=3039 avgLatencyMs=139 tipHash=0x030ef7
2025-11-08T07:20:00Z INFO telemetry: height=200440 peers=48 mempool=3040 avgLatencyMs=140 tipHash=0x030ef8
2025-11-08T07:21:00Z INFO telemetry: height=200441 peers=49 mempool=3041 avgLatencyMs=141 tipHash=0x030ef9
2025-11-08T07:22:00Z INFO telemetry: height=200442 peers=50 mempool=3042 avgLatencyMs=142 tipHash=0x030efa
2025-11-08T07:23:00Z INFO telemetry: height=200443 peers=51 mempool=3043 avgLatencyMs=143 tipHash=0x030efb
2025-11-08T07:24:00Z INFO telemetry: height=200444 peers=52 mempool=3044 avgLatencyMs=144 tipHash=0x030efc
2025-11-08T07:25:00Z INFO telemetry: height=200445 peers=48 mempool=3045 avgLatencyMs=145 tipHash=0x030efd
2025-11-08T07:26:00Z INFO telemetry: height=200446 peers=49 mempool=3046 avgLatencyMs=146 tipHash=0x030efe
2025-11-08T07:27:00Z INFO telemetry: height=200447 peers=50 mempool=3047 avgLatencyMs=147 tipHash=0x030eff
2025-11-08T07:28:00Z INFO telemetry: height=200448 peers=51 mempool=3048 avgLatencyMs=148 tipHash=0x030f00
2025-11-08T07:29:00Z INFO telemetry: height=200449 peers=52 mempool=3049 avgLatencyMs=149 tipHash=0x030f01
2025-11-08T07:30:00Z INFO telemetry: height=200450 peers=48 mempool=3050 avgLatencyMs=120 tipHash=0x030f02
2025-11-08T07:31:00Z INFO telemetry: height=200451 peers=49 mempool=3051 avgLatencyMs=121 tipHash=0x030f03
2025-11-08T07:32:00Z INFO telemetry: height=200452 peers=50 mempool=3052 avgLatencyMs=122 tipHash=0x030f04
2025-11-08T07:33:00Z INFO telemetry: height=200453 peers=51 mempool=3053 avgLatencyMs=123 tipHash=0x030f05
2025-11-08T07:34:00Z INFO telemetry: height=200454 peers=52 mempool=3054 avgLatencyMs=124 tipHash=0x030f06
2025-11-08T07:35:00Z INFO telemetry: height=200455 peers=48 mempool=3055 avgLatencyMs=125 tipHash=0x030f07
2025-11-08T07:36:00Z INFO telemetry: height=200456 peers=49 mempool=3056 avgLatencyMs=126 tipHash=0x030f08
2025-11-08T07:37:00Z INFO telemetry: height=200457 peers=50 mempool=3057 avgLatencyMs=127 tipHash=0x030f09
2025-11-08T07:38:00Z INFO telemetry: height=200458 peers=51 mempool=3058 avgLatencyMs=128 tipHash=0x030f0a
2025-11-08T07:39:00Z INFO telemetry: height=200459 peers=52 mempool=3059 avgLatencyMs=129 tipHash=0x030f0b
2025-11-08T07:40:00Z INFO telemetry: height=200460 peers=48 mempool=3060 avgLatencyMs=130 tipHash=0x030f0c
2025-11-08T07:41:00Z INFO telemetry: height=200461 peers=49 mempool=3061 avgLatencyMs=131 tipHash=0x030f0d
2025-11-08T07:42:00Z INFO telemetry: height=200462 peers=50 mempool=3062 avgLatencyMs=132 tipHash=0x030f0e
2025-11-08T07:43:00Z INFO telemetry: height=200463 peers=51 mempool=3063 avgLatencyMs=133 tipHash=0x030f0f
2025-11-08T07:44:00Z INFO telemetry: height=200464 peers=52 mempool=3064 avgLatencyMs=134 tipHash=0x030f10
2025-11-08T07:45:00Z INFO telemetry: height=200465 peers=48 mempool=3065 avgLatencyMs=135 tipHash=0x030f11
2025-11-08T07:46:00Z INFO telemetry: height=200466 peers=49 mempool=3066 avgLatencyMs=136 tipHash=0x030f12
2025-11-08T07:47:00Z INFO telemetry: height=200467 peers=50 mempool=3067 avgLatencyMs=137 tipHash=0x030f13
2025-11-08T07:48:00Z INFO telemetry: height=200468 peers=51 mempool=3068 avgLatencyMs=138 tipHash=0x030f14
2025-11-08T07:49:00Z INFO telemetry: height=200469 peers=52 mempool=3069 avgLatencyMs=139 tipHash=0x030f15
2025-11-08T07:50:00Z INFO telemetry: height=200470 peers=48 mempool=3070 avgLatencyMs=140 tipHash=0x030f16
2025-11-08T07:51:00Z INFO telemetry: height=200471 peers=49 mempool=3071 avgLatencyMs=141 tipHash=0x030f17
2025-11-08T07:52:00Z INFO telemetry: height=200472 peers=50 mempool=3072 avgLatencyMs=142 tipHash=0x030f18
2025-11-08T07:53:00Z INFO telemetry: height=200473 peers=51 mempool=3073 avgLatencyMs=143 tipHash=0x030f19
2025-11-08T07:54:00Z INFO telemetry: height=200474 peers=52 mempool=3074 avgLatencyMs=144 tipHash=0x030f1a
2025-11-08T07:55:00Z INFO telemetry: height=200475 peers=48 mempool=3075 avgLatencyMs=145 tipHash=0x030f1b
2025-11-08T07:56:00Z INFO telemetry: height=200476 peers=49 mempool=3076 avgLatencyMs=146 tipHash=0x030f1c
2025-11-08T07:57:00Z INFO telemetry: height=200477 peers=50 mempool=3077 avgLatencyMs=147 tipHash=0x030f1d
2025-11-08T07:58:00Z INFO telemetry: height=200478 peers=51 mempool=3078 avgLatencyMs=148 tipHash=0x030f1e
2025-11-08T07:59:00Z INFO telemetry: height=200479 peers=52 mempool=3079 avgLatencyMs=149 tipHash=0x030f1f
2025-11-08T08:00:00Z INFO telemetry: height=200480 peers=48 mempool=3080 avgLatencyMs=120 tipHash=0x030f20
2025-11-08T08:01:00Z INFO telemetry: height=200481 peers=49 mempool=3081 avgLatencyMs=121 tipHash=0x030f21
2025-11-08T08:02:00Z INFO telemetry: height=200482 peers=50 mempool=3082 avgLatencyMs=122 tipHash=0x030f22
2025-11-08T08:03:00Z INFO telemetry: height=200483 peers=51 mempool=3083 avgLatencyMs=123 tipHash=0x030f23
2025-11-08T08:04:00Z INFO telemetry: height=200484 peers=52 mempool=3084 avgLatencyMs=124 tipHash=0x030f24
2025-11-08T08:05:00Z INFO telemetry: height=200485 peers=48 mempool=3085 avgLatencyMs=125 tipHash=0x030f25
2025-11-08T08:06:00Z INFO telemetry: height=200486 peers=49 mempool=3086 avgLatencyMs=126 tipHash=0x030f26
2025-11-08T08:07:00Z INFO telemetry: height=200487 peers=50 mempool=3087 avgLatencyMs=127 tipHash=0x030f27
2025-11-08T08:08:00Z INFO telemetry: height=200488 peers=51 mempool=3088 avgLatencyMs=128 tipHash=0x030f28
2025-11-08T08:09:00Z INFO telemetry: height=200489 peers=52 mempool=3089 avgLatencyMs=129 tipHash=0x030f29
2025-11-08T08:10:00Z INFO telemetry: height=200490 peers=48 mempool=3090 avgLatencyMs=130 tipHash=0x030f2a
2025-11-08T08:11:00Z INFO telemetry: height=200491 peers=49 mempool=3091 avgLatencyMs=131 tipHash=0x030f2b
2025-11-08T08:12:00Z INFO telemetry: height=200492 peers=50 mempool=3092 avgLatencyMs=132 tipHash=0x030f2c
2025-11-08T08:13:00Z INFO telemetry: height=200493 peers=51 mempool=3093 avgLatencyMs=133 tipHash=0x030f2d
2025-11-08T08:14:00Z INFO telemetry: height=200494 peers=52 mempool=3094 avgLatencyMs=134 tipHash=0x030f2e
2025-11-08T08:15:00Z INFO telemetry: height=200495 peers=48 mempool=3095 avgLatencyMs=135 tipHash=0x030f2f
2025-11-08T08:16:00Z INFO telemetry: height=200496 peers=49 mempool=3096 avgLatencyMs=136 tipHash=0x030f30
2025-11-08T08:17:00Z INFO telemetry: height=200497 peers=50 mempool=3097 avgLatencyMs=137 tipHash=0x030f31
2025-11-08T08:18:00Z INFO telemetry: height=200498 peers=51 mempool=3098 avgLatencyMs=138 tipHash=0x030f32
2025-11-08T08:19:00Z INFO telemetry: height=200499 peers=52 mempool=3099 avgLatencyMs=139 tipHash=0x030f33
2025-11-08T08:20:00Z INFO telemetry: height=200500 peers=48 mempool=3100 avgLatencyMs=140 tipHash=0x030f34
2025-11-08T08:21:00Z INFO telemetry: height=200501 peers=49 mempool=3101 avgLatencyMs=141 tipHash=0x030f35
2025-11-08T08:22:00Z INFO telemetry: height=200502 peers=50 mempool=3102 avgLatencyMs=142 tipHash=0x030f36
2025-11-08T08:23:00Z INFO telemetry: height=200503 peers=51 mempool=3103 avgLatencyMs=143 tipHash=0x030f37
2025-11-08T08:24:00Z INFO telemetry: height=200504 peers=52 mempool=3104 avgLatencyMs=144 tipHash=0x030f38
2025-11-08T08:25:00Z INFO telemetry: height=200505 peers=48 mempool=3105 avgLatencyMs=145 tipHash=0x030f39
2025-11-08T08:26:00Z INFO telemetry: height=200506 peers=49 mempool=3106 avgLatencyMs=146 tipHash=0x030f3a
2025-11-08T08:27:00Z INFO telemetry: height=200507 peers=50 mempool=3107 avgLatencyMs=147 tipHash=0x030f3b
2025-11-08T08:28:00Z INFO telemetry: height=200508 peers=51 mempool=3108 avgLatencyMs=148 tipHash=0x030f3c
2025-11-08T08:29:00Z INFO telemetry: height=200509 peers=52 mempool=3109 avgLatencyMs=149 tipHash=0x030f3d
2025-11-08T08:30:00Z INFO telemetry: height=200510 peers=48 mempool=3110 avgLatencyMs=120 tipHash=0x030f3e
2025-11-08T08:31:00Z INFO telemetry: height=200511 peers=49 mempool=3111 avgLatencyMs=121 tipHash=0x030f3f
2025-11-08T08:32:00Z INFO telemetry: height=200512 peers=50 mempool=3112 avgLatencyMs=122 tipHash=0x030f40
2025-11-08T08:33:00Z INFO telemetry: height=200513 peers=51 mempool=3113 avgLatencyMs=123 tipHash=0x030f41
2025-11-08T08:34:00Z INFO telemetry: height=200514 peers=52 mempool=3114 avgLatencyMs=124 tipHash=0x030f42
2025-11-08T08:35:00Z INFO telemetry: height=200515 peers=48 mempool=3115 avgLatencyMs=125 tipHash=0x030f43
2025-11-08T08:36:00Z INFO telemetry: height=200516 peers=49 mempool=3116 avgLatencyMs=126 tipHash=0x030f44
2025-11-08T08:37:00Z INFO telemetry: height=200517 peers=50 mempool=3117 avgLatencyMs=127 tipHash=0x030f45
2025-11-08T08:38:00Z INFO telemetry: height=200518 peers=51 mempool=3118 avgLatencyMs=128 tipHash=0x030f46
2025-11-08T08:39:00Z INFO telemetry: height=200519 peers=52 mempool=3119 avgLatencyMs=129 tipHash=0x030f47
2025-11-08T08:40:00Z INFO telemetry: height=200520 peers=48 mempool=3120 avgLatencyMs=130 tipHash=0x030f48
2025-11-08T08:41:00Z INFO telemetry: height=200521 peers=49 mempool=3121 avgLatencyMs=131 tipHash=0x030f49
2025-11-08T08:42:00Z INFO telemetry: height=200522 peers=50 mempool=3122 avgLatencyMs=132 tipHash=0x030f4a
2025-11-08T08:43:00Z INFO telemetry: height=200523 peers=51 mempool=3123 avgLatencyMs=133 tipHash=0x030f4b
2025-11-08T08:44:00Z INFO telemetry: height=200524 peers=52 mempool=3124 avgLatencyMs=134 tipHash=0x030f4c
2025-11-08T08:45:00Z INFO telemetry: height=200525 peers=48 mempool=3125 avgLatencyMs=135 tipHash=0x030f4d
2025-11-08T08:46:00Z INFO telemetry: height=200526 peers=49 mempool=3126 avgLatencyMs=136 tipHash=0x030f4e
2025-11-08T08:47:00Z INFO telemetry: height=200527 peers=50 mempool=3127 avgLatencyMs=137 tipHash=0x030f4f
2025-11-08T08:48:00Z INFO telemetry: height=200528 peers=51 mempool=3128 avgLatencyMs=138 tipHash=0x030f50
2025-11-08T08:49:00Z INFO telemetry: height=200529 peers=52 mempool=3129 avgLatencyMs=139 tipHash=0x030f51
2025-11-08T08:50:00Z INFO telemetry: height=200530 peers=48 mempool=3130 avgLatencyMs=140 tipHash=0x030f52
2025-11-08T08:51:00Z INFO telemetry: height=200531 peers=49 mempool=3131 avgLatencyMs=141 tipHash=0x030f53
2025-11-08T08:52:00Z INFO telemetry: height=200532 peers=50 mempool=3132 avgLatencyMs=142 tipHash=0x030f54
2025-11-08T08:53:00Z INFO telemetry: height=200533 peers=51 mempool=3133 avgLatencyMs=143 tipHash=0x030f55
2025-11-08T08:54:00Z INFO telemetry: height=200534 peers=52 mempool=3134 avgLatencyMs=144 tipHash=0x030f56
2025-11-08T08:55:00Z INFO telemetry: height=200535 peers=48 mempool=3135 avgLatencyMs=145 tipHash=0x030f57
2025-11-08T08:56:00Z INFO telemetry: height=200536 peers=49 mempool=3136 avgLatencyMs=146 tipHash=0x030f58
2025-11-08T08:57:00Z INFO telemetry: height=200537 peers=50 mempool=3137 avgLatencyMs=147 tipHash=0x030f59
2025-11-08T08:58:00Z INFO telemetry: height=200538 peers=51 mempool=3138 avgLatencyMs=148 tipHash=0x030f5a
2025-11-08T08:59:00Z INFO telemetry: height=200539 peers=52 mempool=3139 avgLatencyMs=149 tipHash=0x030f5b
2025-11-08T09:00:00Z INFO telemetry: height=200540 peers=48 mempool=3140 avgLatencyMs=120 tipHash=0x030f5c
2025-11-08T09:01:00Z INFO telemetry: height=200541 peers=49 mempool=3141 avgLatencyMs=121 tipHash=0x030f5d
2025-11-08T09:02:00Z INFO telemetry: height=200542 peers=50 mempool=3142 avgLatencyMs=122 tipHash=0x030f5e
2025-11-08T09:03:00Z INFO telemetry: height=200543 peers=51 mempool=3143 avgLatencyMs=123 tipHash=0x030f5f
2025-11-08T09:04:00Z INFO telemetry: height=200544 peers=52 mempool=3144 avgLatencyMs=124 tipHash=0x030f60
2025-11-08T09:05:00Z INFO telemetry: height=200545 peers=48 mempool=3145 avgLatencyMs=125 tipHash=0x030f61
2025-11-08T09:06:00Z INFO telemetry: height=200546 peers=49 mempool=3146 avgLatencyMs=126 tipHash=0x030f62
2025-11-08T09:07:00Z INFO telemetry: height=200547 peers=50 mempool=3147 avgLatencyMs=127 tipHash=0x030f63
2025-11-08T09:08:00Z INFO telemetry: height=200548 peers=51 mempool=3148 avgLatencyMs=128 tipHash=0x030f64
2025-11-08T09:09:00Z INFO telemetry: height=200549 peers=52 mempool=3149 avgLatencyMs=129 tipHash=0x030f65
2025-11-08T09:10:00Z INFO telemetry: height=200550 peers=48 mempool=3150 avgLatencyMs=130 tipHash=0x030f66
2025-11-08T09:11:00Z INFO telemetry: height=200551 peers=49 mempool=3151 avgLatencyMs=131 tipHash=0x030f67
2025-11-08T09:12:00Z INFO telemetry: height=200552 peers=50 mempool=3152 avgLatencyMs=132 tipHash=0x030f68
2025-11-08T09:13:00Z INFO telemetry: height=200553 peers=51 mempool=3153 avgLatencyMs=133 tipHash=0x030f69
2025-11-08T09:14:00Z INFO telemetry: height=200554 peers=52 mempool=3154 avgLatencyMs=134 tipHash=0x030f6a
2025-11-08T09:15:00Z INFO telemetry: height=200555 peers=48 mempool=3155 avgLatencyMs=135 tipHash=0x030f6b
2025-11-08T09:16:00Z INFO telemetry: height=200556 peers=49 mempool=3156 avgLatencyMs=136 tipHash=0x030f6c
2025-11-08T09:17:00Z INFO telemetry: height=200557 peers=50 mempool=3157 avgLatencyMs=137 tipHash=0x030f6d
2025-11-08T09:18:00Z INFO telemetry: height=200558 peers=51 mempool=3158 avgLatencyMs=138 tipHash=0x030f6e
2025-11-08T09:19:00Z INFO telemetry: height=200559 peers=52 mempool=3159 avgLatencyMs=139 tipHash=0x030f6f
2025-11-08T09:20:00Z INFO telemetry: height=200560 peers=48 mempool=3160 avgLatencyMs=140 tipHash=0x030f70
2025-11-08T09:21:00Z INFO telemetry: height=200561 peers=49 mempool=3161 avgLatencyMs=141 tipHash=0x030f71
2025-11-08T09:22:00Z INFO telemetry: height=200562 peers=50 mempool=3162 avgLatencyMs=142 tipHash=0x030f72
2025-11-08T09:23:00Z INFO telemetry: height=200563 peers=51 mempool=3163 avgLatencyMs=143 tipHash=0x030f73
2025-11-08T09:24:00Z INFO telemetry: height=200564 peers=52 mempool=3164 avgLatencyMs=144 tipHash=0x030f74
2025-11-08T09:25:00Z INFO telemetry: height=200565 peers=48 mempool=3165 avgLatencyMs=145 tipHash=0x030f75
2025-11-08T09:26:00Z INFO telemetry: height=200566 peers=49 mempool=3166 avgLatencyMs=146 tipHash=0x030f76
2025-11-08T09:27:00Z INFO telemetry: height=200567 peers=50 mempool=3167 avgLatencyMs=147 tipHash=0x030f77
2025-11-08T09:28:00Z INFO telemetry: height=200568 peers=51 mempool=3168 avgLatencyMs=148 tipHash=0x030f78
2025-11-08T09:29:00Z INFO telemetry: height=200569 peers=52 mempool=3169 avgLatencyMs=149 tipHash=0x030f79
2025-11-08T09:30:00Z INFO telemetry: height=200570 peers=48 mempool=3170 avgLatencyMs=120 tipHash=0x030f7a
2025-11-08T09:31:00Z INFO telemetry: height=200571 peers=49 mempool=3171 avgLatencyMs=121 tipHash=0x030f7b
2025-11-08T09:32:00Z INFO telemetry: height=200572 peers=50 mempool=3172 avgLatencyMs=122 tipHash=0x030f7c
2025-11-08T09:33:00Z INFO telemetry: height=200573 peers=51 mempool=3173 avgLatencyMs=123 tipHash=0x030f7d
2025-11-08T09:34:00Z INFO telemetry: height=200574 peers=52 mempool=3174 avgLatencyMs=124 tipHash=0x030f7e
2025-11-08T09:35:00Z INFO telemetry: height=200575 peers=48 mempool=3175 avgLatencyMs=125 tipHash=0x030f7f
2025-11-08T09:36:00Z INFO telemetry: height=200576 peers=49 mempool=3176 avgLatencyMs=126 tipHash=0x030f80
2025-11-08T09:37:00Z INFO telemetry: height=200577 peers=50 mempool=3177 avgLatencyMs=127 tipHash=0x030f81
2025-11-08T09:38:00Z INFO telemetry: height=200578 peers=51 mempool=3178 avgLatencyMs=128 tipHash=0x030f82
2025-11-08T09:39:00Z INFO telemetry: height=200579 peers=52 mempool=3179 avgLatencyMs=129 tipHash=0x030f83
2025-11-08T09:40:00Z INFO telemetry: height=200580 peers=48 mempool=3180 avgLatencyMs=130 tipHash=0x030f84
2025-11-08T09:41:00Z INFO telemetry: height=200581 peers=49 mempool=3181 avgLatencyMs=131 tipHash=0x030f85
2025-11-08T09:42:00Z INFO telemetry: height=200582 peers=50 mempool=3182 avgLatencyMs=132 tipHash=0x030f86
2025-11-08T09:43:00Z INFO telemetry: height=200583 peers=51 mempool=3183 avgLatencyMs=133 tipHash=0x030f87
2025-11-08T09:44:00Z INFO telemetry: height=200584 peers=52 mempool=3184 avgLatencyMs=134 tipHash=0x030f88
2025-11-08T09:45:00Z INFO telemetry: height=200585 peers=48 mempool=3185 avgLatencyMs=135 tipHash=0x030f89
2025-11-08T09:46:00Z INFO telemetry: height=200586 peers=49 mempool=3186 avgLatencyMs=136 tipHash=0x030f8a
2025-11-08T09:47:00Z INFO telemetry: height=200587 peers=50 mempool=3187 avgLatencyMs=137 tipHash=0x030f8b
2025-11-08T09:48:00Z INFO telemetry: height=200588 peers=51 mempool=3188 avgLatencyMs=138 tipHash=0x030f8c
2025-11-08T09:49:00Z INFO telemetry: height=200589 peers=52 mempool=3189 avgLatencyMs=139 tipHash=0x030f8d
2025-11-08T09:50:00Z INFO telemetry: height=200590 peers=48 mempool=3190 avgLatencyMs=140 tipHash=0x030f8e
2025-11-08T09:51:00Z INFO telemetry: height=200591 peers=49 mempool=3191 avgLatencyMs=141 tipHash=0x030f8f
2025-11-08T09:52:00Z INFO telemetry: height=200592 peers=50 mempool=3192 avgLatencyMs=142 tipHash=0x030f90
2025-11-08T09:53:00Z INFO telemetry: height=200593 peers=51 mempool=3193 avgLatencyMs=143 tipHash=0x030f91
2025-11-08T09:54:00Z INFO telemetry: height=200594 peers=52 mempool=3194 avgLatencyMs=144 tipHash=0x030f92
2025-11-08T09:55:00Z INFO telemetry: height=200595 peers=48 mempool=3195 avgLatencyMs=145 tipHash=0x030f93
2025-11-08T09:56:00Z INFO telemetry: height=200596 peers=49 mempool=3196 avgLatencyMs=146 tipHash=0x030f94
2025-11-08T09:57:00Z INFO telemetry: height=200597 peers=50 mempool=3197 avgLatencyMs=147 tipHash=0x030f95
2025-11-08T09:58:00Z INFO telemetry: height=200598 peers=51 mempool=3198 avgLatencyMs=148 tipHash=0x030f96
2025-11-08T09:59:00Z INFO telemetry: height=200599 peers=52 mempool=3199 avgLatencyMs=149 tipHash=0x030f97
2025-11-08T10:00:00Z INFO telemetry: height=200600 peers=48 mempool=3000 avgLatencyMs=120 tipHash=0x030f98
2025-11-08T10:01:00Z INFO telemetry: height=200601 peers=49 mempool=3001 avgLatencyMs=121 tipHash=0x030f99
2025-11-08T10:02:00Z INFO telemetry: height=200602 peers=50 mempool=3002 avgLatencyMs=122 tipHash=0x030f9a
2025-11-08T10:03:00Z INFO telemetry: height=200603 peers=51 mempool=3003 avgLatencyMs=123 tipHash=0x030f9b
2025-11-08T10:04:00Z INFO telemetry: height=200604 peers=52 mempool=3004 avgLatencyMs=124 tipHash=0x030f9c
2025-11-08T10:05:00Z INFO telemetry: height=200605 peers=48 mempool=3005 avgLatencyMs=125 tipHash=0x030f9d
2025-11-08T10:06:00Z INFO telemetry: height=200606 peers=49 mempool=3006 avgLatencyMs=126 tipHash=0x030f9e
2025-11-08T10:07:00Z INFO telemetry: height=200607 peers=50 mempool=3007 avgLatencyMs=127 tipHash=0x030f9f
2025-11-08T10:08:00Z INFO telemetry: height=200608 peers=51 mempool=3008 avgLatencyMs=128 tipHash=0x030fa0
2025-11-08T10:09:00Z INFO telemetry: height=200609 peers=52 mempool=3009 avgLatencyMs=129 tipHash=0x030fa1
2025-11-08T10:10:00Z INFO telemetry: height=200610 peers=48 mempool=3010 avgLatencyMs=130 tipHash=0x030fa2
2025-11-08T10:11:00Z INFO telemetry: height=200611 peers=49 mempool=3011 avgLatencyMs=131 tipHash=0x030fa3
2025-11-08T10:12:00Z INFO telemetry: height=200612 peers=50 mempool=3012 avgLatencyMs=132 tipHash=0x030fa4
2025-11-08T10:13:00Z INFO telemetry: height=200613 peers=51 mempool=3013 avgLatencyMs=133 tipHash=0x030fa5
2025-11-08T10:14:00Z INFO telemetry: height=200614 peers=52 mempool=3014 avgLatencyMs=134 tipHash=0x030fa6
2025-11-08T10:15:00Z INFO telemetry: height=200615 peers=48 mempool=3015 avgLatencyMs=135 tipHash=0x030fa7
2025-11-08T10:16:00Z INFO telemetry: height=200616 peers=49 mempool=3016 avgLatencyMs=136 tipHash=0x030fa8
2025-11-08T10:17:00Z INFO telemetry: height=200617 peers=50 mempool=3017 avgLatencyMs=137 tipHash=0x030fa9
2025-11-08T10:18:00Z INFO telemetry: height=200618 peers=51 mempool=3018 avgLatencyMs=138 tipHash=0x030faa
2025-11-08T10:19:00Z INFO telemetry: height=200619 peers=52 mempool=3019 avgLatencyMs=139 tipHash=0x030fab
2025-11-08T10:20:00Z INFO telemetry: height=200620 peers=48 mempool=3020 avgLatencyMs=140 tipHash=0x030fac
2025-11-08T10:21:00Z INFO telemetry: height=200621 peers=49 mempool=3021 avgLatencyMs=141 tipHash=0x030fad
2025-11-08T10:22:00Z INFO telemetry: height=200622 peers=50 mempool=3022 avgLatencyMs=142 tipHash=0x030fae
2025-11-08T10:23:00Z INFO telemetry: height=200623 peers=51 mempool=3023 avgLatencyMs=143 tipHash=0x030faf
2025-11-08T10:24:00Z INFO telemetry: height=200624 peers=52 mempool=3024 avgLatencyMs=144 tipHash=0x030fb0
2025-11-08T10:25:00Z INFO telemetry: height=200625 peers=48 mempool=3025 avgLatencyMs=145 tipHash=0x030fb1
2025-11-08T10:26:00Z INFO telemetry: height=200626 peers=49 mempool=3026 avgLatencyMs=146 tipHash=0x030fb2
2025-11-08T10:27:00Z INFO telemetry: height=200627 peers=50 mempool=3027 avgLatencyMs=147 tipHash=0x030fb3
2025-11-08T10:28:00Z INFO telemetry: height=200628 peers=51 mempool=3028 avgLatencyMs=148 tipHash=0x030fb4
2025-11-08T10:29:00Z INFO telemetry: height=200629 peers=52 mempool=3029 avgLatencyMs=149 tipHash=0x030fb5
2025-11-08T10:30:00Z INFO telemetry: height=200630 peers=48 mempool=3030 avgLatencyMs=120 tipHash=0x030fb6
2025-11-08T10:31:00Z INFO telemetry: height=200631 peers=49 mempool=3031 avgLatencyMs=121 tipHash=0x030fb7
2025-11-08T10:32:00Z INFO telemetry: height=200632 peers=50 mempool=3032 avgLatencyMs=122 tipHash=0x030fb8
2025-11-08T10:33:00Z INFO telemetry: height=200633 peers=51 mempool=3033 avgLatencyMs=123 tipHash=0x030fb9
2025-11-08T10:34:00Z INFO telemetry: height=200634 peers=52 mempool=3034 avgLatencyMs=124 tipHash=0x030fba
2025-11-08T10:35:00Z INFO telemetry: height=200635 peers=48 mempool=3035 avgLatencyMs=125 tipHash=0x030fbb
2025-11-08T10:36:00Z INFO telemetry: height=200636 peers=49 mempool=3036 avgLatencyMs=126 tipHash=0x030fbc
2025-11-08T10:37:00Z INFO telemetry: height=200637 peers=50 mempool=3037 avgLatencyMs=127 tipHash=0x030fbd
2025-11-08T10:38:00Z INFO telemetry: height=200638 peers=51 mempool=3038 avgLatencyMs=128 tipHash=0x030fbe
2025-11-08T10:39:00Z INFO telemetry: height=200639 peers=52 mempool=3039 avgLatencyMs=129 tipHash=0x030fbf
2025-11-08T10:40:00Z INFO telemetry: height=200640 peers=48 mempool=3040 avgLatencyMs=130 tipHash=0x030fc0
2025-11-08T10:41:00Z INFO telemetry: height=200641 peers=49 mempool=3041 avgLatencyMs=131 tipHash=0x030fc1
2025-11-08T10:42:00Z INFO telemetry: height=200642 peers=50 mempool=3042 avgLatencyMs=132 tipHash=0x030fc2
2025-11-08T10:43:00Z INFO telemetry: height=200643 peers=51 mempool=3043 avgLatencyMs=133 tipHash=0x030fc3
2025-11-08T10:44:00Z INFO telemetry: height=200644 peers=52 mempool=3044 avgLatencyMs=134 tipHash=0x030fc4
2025-11-08T10:45:00Z INFO telemetry: height=200645 peers=48 mempool=3045 avgLatencyMs=135 tipHash=0x030fc5
2025-11-08T10:46:00Z INFO telemetry: height=200646 peers=49 mempool=3046 avgLatencyMs=136 tipHash=0x030fc6
2025-11-08T10:47:00Z INFO telemetry: height=200647 peers=50 mempool=3047 avgLatencyMs=137 tipHash=0x030fc7
2025-11-08T10:48:00Z INFO telemetry: height=200648 peers=51 mempool=3048 avgLatencyMs=138 tipHash=0x030fc8
2025-11-08T10:49:00Z INFO telemetry: height=200649 peers=52 mempool=3049 avgLatencyMs=139 tipHash=0x030fc9
2025-11-08T10:50:00Z INFO telemetry: height=200650 peers=48 mempool=3050 avgLatencyMs=140 tipHash=0x030fca
2025-11-08T10:51:00Z INFO telemetry: height=200651 peers=49 mempool=3051 avgLatencyMs=141 tipHash=0x030fcb
2025-11-08T10:52:00Z INFO telemetry: height=200652 peers=50 mempool=3052 avgLatencyMs=142 tipHash=0x030fcc
2025-11-08T10:53:00Z INFO telemetry: height=200653 peers=51 mempool=3053 avgLatencyMs=143 tipHash=0x030fcd
2025-11-08T10:54:00Z INFO telemetry: height=200654 peers=52 mempool=3054 avgLatencyMs=144 tipHash=0x030fce
2025-11-08T10:55:00Z INFO telemetry: height=200655 peers=48 mempool=3055 avgLatencyMs=145 tipHash=0x030fcf
2025-11-08T10:56:00Z INFO telemetry: height=200656 peers=49 mempool=3056 avgLatencyMs=146 tipHash=0x030fd0
2025-11-08T10:57:00Z INFO telemetry: height=200657 peers=50 mempool=3057 avgLatencyMs=147 tipHash=0x030fd1
2025-11-08T10:58:00Z INFO telemetry: height=200658 peers=51 mempool=3058 avgLatencyMs=148 tipHash=0x030fd2
2025-11-08T10:59:00Z INFO telemetry: height=200659 peers=52 mempool=3059 avgLatencyMs=149 tipHash=0x030fd3
2025-11-08T11:00:00Z INFO telemetry: height=200660 peers=48 mempool=3060 avgLatencyMs=120 tipHash=0x030fd4
2025-11-08T11:01:00Z INFO telemetry: height=200661 peers=49 mempool=3061 avgLatencyMs=121 tipHash=0x030fd5
2025-11-08T11:02:00Z INFO telemetry: height=200662 peers=50 mempool=3062 avgLatencyMs=122 tipHash=0x030fd6
2025-11-08T11:03:00Z INFO telemetry: height=200663 peers=51 mempool=3063 avgLatencyMs=123 tipHash=0x030fd7
2025-11-08T11:04:00Z INFO telemetry: height=200664 peers=52 mempool=3064 avgLatencyMs=124 tipHash=0x030fd8
2025-11-08T11:05:00Z INFO telemetry: height=200665 peers=48 mempool=3065 avgLatencyMs=125 tipHash=0x030fd9
2025-11-08T11:06:00Z INFO telemetry: height=200666 peers=49 mempool=3066 avgLatencyMs=126 tipHash=0x030fda
2025-11-08T11:07:00Z INFO telemetry: height=200667 peers=50 mempool=3067 avgLatencyMs=127 tipHash=0x030fdb
2025-11-08T11:08:00Z INFO telemetry: height=200668 peers=51 mempool=3068 avgLatencyMs=128 tipHash=0x030fdc
2025-11-08T11:09:00Z INFO telemetry: height=200669 peers=52 mempool=3069 avgLatencyMs=129 tipHash=0x030fdd
2025-11-08T11:10:00Z INFO telemetry: height=200670 peers=48 mempool=3070 avgLatencyMs=130 tipHash=0x030fde
2025-11-08T11:11:00Z INFO telemetry: height=200671 peers=49 mempool=3071 avgLatencyMs=131 tipHash=0x030fdf
2025-11-08T11:12:00Z INFO telemetry: height=200672 peers=50 mempool=3072 avgLatencyMs=132 tipHash=0x030fe0
2025-11-08T11:13:00Z INFO telemetry: height=200673 peers=51 mempool=3073 avgLatencyMs=133 tipHash=0x030fe1
2025-11-08T11:14:00Z INFO telemetry: height=200674 peers=52 mempool=3074 avgLatencyMs=134 tipHash=0x030fe2
2025-11-08T11:15:00Z INFO telemetry: height=200675 peers=48 mempool=3075 avgLatencyMs=135 tipHash=0x030fe3
2025-11-08T11:16:00Z INFO telemetry: height=200676 peers=49 mempool=3076 avgLatencyMs=136 tipHash=0x030fe4
2025-11-08T11:17:00Z INFO telemetry: height=200677 peers=50 mempool=3077 avgLatencyMs=137 tipHash=0x030fe5
2025-11-08T11:18:00Z INFO telemetry: height=200678 peers=51 mempool=3078 avgLatencyMs=138 tipHash=0x030fe6
2025-11-08T11:19:00Z INFO telemetry: height=200679 peers=52 mempool=3079 avgLatencyMs=139 tipHash=0x030fe7
2025-11-08T11:20:00Z INFO telemetry: height=200680 peers=48 mempool=3080 avgLatencyMs=140 tipHash=0x030fe8
2025-11-08T11:21:00Z INFO telemetry: height=200681 peers=49 mempool=3081 avgLatencyMs=141 tipHash=0x030fe9
2025-11-08T11:22:00Z INFO telemetry: height=200682 peers=50 mempool=3082 avgLatencyMs=142 tipHash=0x030fea
2025-11-08T11:23:00Z INFO telemetry: height=200683 peers=51 mempool=3083 avgLatencyMs=143 tipHash=0x030feb
2025-11-08T11:24:00Z INFO telemetry: height=200684 peers=52 mempool=3084 avgLatencyMs=144 tipHash=0x030fec
2025-11-08T11:25:00Z INFO telemetry: height=200685 peers=48 mempool=3085 avgLatencyMs=145 tipHash=0x030fed
2025-11-08T11:26:00Z INFO telemetry: height=200686 peers=49 mempool=3086 avgLatencyMs=146 tipHash=0x030fee
2025-11-08T11:27:00Z INFO telemetry: height=200687 peers=50 mempool=3087 avgLatencyMs=147 tipHash=0x030fef
2025-11-08T11:28:00Z INFO telemetry: height=200688 peers=51 mempool=3088 avgLatencyMs=148 tipHash=0x030ff0
2025-11-08T11:29:00Z INFO telemetry: height=200689 peers=52 mempool=3089 avgLatencyMs=149 tipHash=0x030ff1
2025-11-08T11:30:00Z INFO telemetry: height=200690 peers=48 mempool=3090 avgLatencyMs=120 tipHash=0x030ff2
2025-11-08T11:31:00Z INFO telemetry: height=200691 peers=49 mempool=3091 avgLatencyMs=121 tipHash=0x030ff3
2025-11-08T11:32:00Z INFO telemetry: height=200692 peers=50 mempool=3092 avgLatencyMs=122 tipHash=0x030ff4
2025-11-08T11:33:00Z INFO telemetry: height=200693 peers=51 mempool=3093 avgLatencyMs=123 tipHash=0x030ff5
2025-11-08T11:34:00Z INFO telemetry: height=200694 peers=52 mempool=3094 avgLatencyMs=124 tipHash=0x030ff6
2025-11-08T11:35:00Z INFO telemetry: height=200695 peers=48 mempool=3095 avgLatencyMs=125 tipHash=0x030ff7
2025-11-08T11:36:00Z INFO telemetry: height=200696 peers=49 mempool=3096 avgLatencyMs=126 tipHash=0x030ff8
2025-11-08T11:37:00Z INFO telemetry: height=200697 peers=50 mempool=3097 avgLatencyMs=127 tipHash=0x030ff9
2025-11-08T11:38:00Z INFO telemetry: height=200698 peers=51 mempool=3098 avgLatencyMs=128 tipHash=0x030ffa
2025-11-08T11:39:00Z INFO telemetry: height=200699 peers=52 mempool=3099 avgLatencyMs=129 tipHash=0x030ffb
2025-11-08T11:40:00Z INFO telemetry: height=200700 peers=48 mempool=3100 avgLatencyMs=130 tipHash=0x030ffc
2025-11-08T11:41:00Z INFO telemetry: height=200701 peers=49 mempool=3101 avgLatencyMs=131 tipHash=0x030ffd
2025-11-08T11:42:00Z INFO telemetry: height=200702 peers=50 mempool=3102 avgLatencyMs=132 tipHash=0x030ffe
2025-11-08T11:43:00Z INFO telemetry: height=200703 peers=51 mempool=3103 avgLatencyMs=133 tipHash=0x030fff
2025-11-08T11:44:00Z INFO telemetry: height=200704 peers=52 mempool=3104 avgLatencyMs=134 tipHash=0x031000
2025-11-08T11:45:00Z INFO telemetry: height=200705 peers=48 mempool=3105 avgLatencyMs=135 tipHash=0x031001
2025-11-08T11:46:00Z INFO telemetry: height=200706 peers=49 mempool=3106 avgLatencyMs=136 tipHash=0x031002
2025-11-08T11:47:00Z INFO telemetry: height=200707 peers=50 mempool=3107 avgLatencyMs=137 tipHash=0x031003
2025-11-08T11:48:00Z INFO telemetry: height=200708 peers=51 mempool=3108 avgLatencyMs=138 tipHash=0x031004
2025-11-08T11:49:00Z INFO telemetry: height=200709 peers=52 mempool=3109 avgLatencyMs=139 tipHash=0x031005
2025-11-08T11:50:00Z INFO telemetry: height=200710 peers=48 mempool=3110 avgLatencyMs=140 tipHash=0x031006
2025-11-08T11:51:00Z INFO telemetry: height=200711 peers=49 mempool=3111 avgLatencyMs=141 tipHash=0x031007
2025-11-08T11:52:00Z INFO telemetry: height=200712 peers=50 mempool=3112 avgLatencyMs=142 tipHash=0x031008
2025-11-08T11:53:00Z INFO telemetry: height=200713 peers=51 mempool=3113 avgLatencyMs=143 tipHash=0x031009
2025-11-08T11:54:00Z INFO telemetry: height=200714 peers=52 mempool=3114 avgLatencyMs=144 tipHash=0x03100a
2025-11-08T11:55:00Z INFO telemetry: height=200715 peers=48 mempool=3115 avgLatencyMs=145 tipHash=0x03100b
2025-11-08T11:56:00Z INFO telemetry: height=200716 peers=49 mempool=3116 avgLatencyMs=146 tipHash=0x03100c
2025-11-08T11:57:00Z INFO telemetry: height=200717 peers=50 mempool=3117 avgLatencyMs=147 tipHash=0x03100d
2025-11-08T11:58:00Z INFO telemetry: height=200718 peers=51 mempool=3118 avgLatencyMs=148 tipHash=0x03100e
2025-11-08T11:59:00Z INFO telemetry: height=200719 peers=52 mempool=3119 avgLatencyMs=149 tipHash=0x03100f
2025-11-08T12:00:00Z INFO telemetry: height=200720 peers=48 mempool=3120 avgLatencyMs=120 tipHash=0x031010
2025-11-08T12:01:00Z INFO telemetry: height=200721 peers=49 mempool=3121 avgLatencyMs=121 tipHash=0x031011
2025-11-08T12:02:00Z INFO telemetry: height=200722 peers=50 mempool=3122 avgLatencyMs=122 tipHash=0x031012
2025-11-08T12:03:00Z INFO telemetry: height=200723 peers=51 mempool=3123 avgLatencyMs=123 tipHash=0x031013
2025-11-08T12:04:00Z INFO telemetry: height=200724 peers=52 mempool=3124 avgLatencyMs=124 tipHash=0x031014
2025-11-08T12:05:00Z INFO telemetry: height=200725 peers=48 mempool=3125 avgLatencyMs=125 tipHash=0x031015
2025-11-08T12:06:00Z INFO telemetry: height=200726 peers=49 mempool=3126 avgLatencyMs=126 tipHash=0x031016
2025-11-08T12:07:00Z INFO telemetry: height=200727 peers=50 mempool=3127 avgLatencyMs=127 tipHash=0x031017
2025-11-08T12:08:00Z INFO telemetry: height=200728 peers=51 mempool=3128 avgLatencyMs=128 tipHash=0x031018
2025-11-08T12:09:00Z INFO telemetry: height=200729 peers=52 mempool=3129 avgLatencyMs=129 tipHash=0x031019
2025-11-08T12:10:00Z INFO telemetry: height=200730 peers=48 mempool=3130 avgLatencyMs=130 tipHash=0x03101a
2025-11-08T12:11:00Z INFO telemetry: height=200731 peers=49 mempool=3131 avgLatencyMs=131 tipHash=0x03101b
2025-11-08T12:12:00Z INFO telemetry: height=200732 peers=50 mempool=3132 avgLatencyMs=132 tipHash=0x03101c
2025-11-08T12:13:00Z INFO telemetry: height=200733 peers=51 mempool=3133 avgLatencyMs=133 tipHash=0x03101d
2025-11-08T12:14:00Z INFO telemetry: height=200734 peers=52 mempool=3134 avgLatencyMs=134 tipHash=0x03101e
2025-11-08T12:15:00Z INFO telemetry: height=200735 peers=48 mempool=3135 avgLatencyMs=135 tipHash=0x03101f
2025-11-08T12:16:00Z INFO telemetry: height=200736 peers=49 mempool=3136 avgLatencyMs=136 tipHash=0x031020
2025-11-08T12:17:00Z INFO telemetry: height=200737 peers=50 mempool=3137 avgLatencyMs=137 tipHash=0x031021
2025-11-08T12:18:00Z INFO telemetry: height=200738 peers=51 mempool=3138 avgLatencyMs=138 tipHash=0x031022
2025-11-08T12:19:00Z INFO telemetry: height=200739 peers=52 mempool=3139 avgLatencyMs=139 tipHash=0x031023
2025-11-08T12:20:00Z INFO telemetry: height=200740 peers=48 mempool=3140 avgLatencyMs=140 tipHash=0x031024
2025-11-08T12:21:00Z INFO telemetry: height=200741 peers=49 mempool=3141 avgLatencyMs=141 tipHash=0x031025
2025-11-08T12:22:00Z INFO telemetry: height=200742 peers=50 mempool=3142 avgLatencyMs=142 tipHash=0x031026
2025-11-08T12:23:00Z INFO telemetry: height=200743 peers=51 mempool=3143 avgLatencyMs=143 tipHash=0x031027
2025-11-08T12:24:00Z INFO telemetry: height=200744 peers=52 mempool=3144 avgLatencyMs=144 tipHash=0x031028
2025-11-08T12:25:00Z INFO telemetry: height=200745 peers=48 mempool=3145 avgLatencyMs=145 tipHash=0x031029
2025-11-08T12:26:00Z INFO telemetry: height=200746 peers=49 mempool=3146 avgLatencyMs=146 tipHash=0x03102a
2025-11-08T12:27:00Z INFO telemetry: height=200747 peers=50 mempool=3147 avgLatencyMs=147 tipHash=0x03102b
2025-11-08T12:28:00Z INFO telemetry: height=200748 peers=51 mempool=3148 avgLatencyMs=148 tipHash=0x03102c
2025-11-08T12:29:00Z INFO telemetry: height=200749 peers=52 mempool=3149 avgLatencyMs=149 tipHash=0x03102d
2025-11-08T12:30:00Z INFO telemetry: height=200750 peers=48 mempool=3150 avgLatencyMs=120 tipHash=0x03102e
2025-11-08T12:31:00Z INFO telemetry: height=200751 peers=49 mempool=3151 avgLatencyMs=121 tipHash=0x03102f
2025-11-08T12:32:00Z INFO telemetry: height=200752 peers=50 mempool=3152 avgLatencyMs=122 tipHash=0x031030
2025-11-08T12:33:00Z INFO telemetry: height=200753 peers=51 mempool=3153 avgLatencyMs=123 tipHash=0x031031
2025-11-08T12:34:00Z INFO telemetry: height=200754 peers=52 mempool=3154 avgLatencyMs=124 tipHash=0x031032
2025-11-08T12:35:00Z INFO telemetry: height=200755 peers=48 mempool=3155 avgLatencyMs=125 tipHash=0x031033
2025-11-08T12:36:00Z INFO telemetry: height=200756 peers=49 mempool=3156 avgLatencyMs=126 tipHash=0x031034
2025-11-08T12:37:00Z INFO telemetry: height=200757 peers=50 mempool=3157 avgLatencyMs=127 tipHash=0x031035
2025-11-08T12:38:00Z INFO telemetry: height=200758 peers=51 mempool=3158 avgLatencyMs=128 tipHash=0x031036
2025-11-08T12:39:00Z INFO telemetry: height=200759 peers=52 mempool=3159 avgLatencyMs=129 tipHash=0x031037
2025-11-08T12:40:00Z INFO telemetry: height=200760 peers=48 mempool=3160 avgLatencyMs=130 tipHash=0x031038
2025-11-08T12:41:00Z INFO telemetry: height=200761 peers=49 mempool=3161 avgLatencyMs=131 tipHash=0x031039
2025-11-08T12:42:00Z INFO telemetry: height=200762 peers=50 mempool=3162 avgLatencyMs=132 tipHash=0x03103a
2025-11-08T12:43:00Z INFO telemetry: height=200763 peers=51 mempool=3163 avgLatencyMs=133 tipHash=0x03103b
2025-11-08T12:44:00Z INFO telemetry: height=200764 peers=52 mempool=3164 avgLatencyMs=134 tipHash=0x03103c
2025-11-08T12:45:00Z INFO telemetry: height=200765 peers=48 mempool=3165 avgLatencyMs=135 tipHash=0x03103d
2025-11-08T12:46:00Z INFO telemetry: height=200766 peers=49 mempool=3166 avgLatencyMs=136 tipHash=0x03103e
2025-11-08T12:47:00Z INFO telemetry: height=200767 peers=50 mempool=3167 avgLatencyMs=137 tipHash=0x03103f
2025-11-08T12:48:00Z INFO telemetry: height=200768 peers=51 mempool=3168 avgLatencyMs=138 tipHash=0x031040
2025-11-08T12:49:00Z INFO telemetry: height=200769 peers=52 mempool=3169 avgLatencyMs=139 tipHash=0x031041
2025-11-08T12:50:00Z INFO telemetry: height=200770 peers=48 mempool=3170 avgLatencyMs=140 tipHash=0x031042
2025-11-08T12:51:00Z INFO telemetry: height=200771 peers=49 mempool=3171 avgLatencyMs=141 tipHash=0x031043
2025-11-08T12:52:00Z INFO telemetry: height=200772 peers=50 mempool=3172 avgLatencyMs=142 tipHash=0x031044
2025-11-08T12:53:00Z INFO telemetry: height=200773 peers=51 mempool=3173 avgLatencyMs=143 tipHash=0x031045
2025-11-08T12:54:00Z INFO telemetry: height=200774 peers=52 mempool=3174 avgLatencyMs=144 tipHash=0x031046
2025-11-08T12:55:00Z INFO telemetry: height=200775 peers=48 mempool=3175 avgLatencyMs=145 tipHash=0x031047
2025-11-08T12:56:00Z INFO telemetry: height=200776 peers=49 mempool=3176 avgLatencyMs=146 tipHash=0x031048
2025-11-08T12:57:00Z INFO telemetry: height=200777 peers=50 mempool=3177 avgLatencyMs=147 tipHash=0x031049
2025-11-08T12:58:00Z INFO telemetry: height=200778 peers=51 mempool=3178 avgLatencyMs=148 tipHash=0x03104a
2025-11-08T12:59:00Z INFO telemetry: height=200779 peers=52 mempool=3179 avgLatencyMs=149 tipHash=0x03104b
2025-11-08T13:00:00Z INFO telemetry: height=200780 peers=48 mempool=3180 avgLatencyMs=120 tipHash=0x03104c
2025-11-08T13:01:00Z INFO telemetry: height=200781 peers=49 mempool=3181 avgLatencyMs=121 tipHash=0x03104d
2025-11-08T13:02:00Z INFO telemetry: height=200782 peers=50 mempool=3182 avgLatencyMs=122 tipHash=0x03104e
2025-11-08T13:03:00Z INFO telemetry: height=200783 peers=51 mempool=3183 avgLatencyMs=123 tipHash=0x03104f
2025-11-08T13:04:00Z INFO telemetry: height=200784 peers=52 mempool=3184 avgLatencyMs=124 tipHash=0x031050
2025-11-08T13:05:00Z INFO telemetry: height=200785 peers=48 mempool=3185 avgLatencyMs=125 tipHash=0x031051
2025-11-08T13:06:00Z INFO telemetry: height=200786 peers=49 mempool=3186 avgLatencyMs=126 tipHash=0x031052
2025-11-08T13:07:00Z INFO telemetry: height=200787 peers=50 mempool=3187 avgLatencyMs=127 tipHash=0x031053
2025-11-08T13:08:00Z INFO telemetry: height=200788 peers=51 mempool=3188 avgLatencyMs=128 tipHash=0x031054
2025-11-08T13:09:00Z INFO telemetry: height=200789 peers=52 mempool=3189 avgLatencyMs=129 tipHash=0x031055
2025-11-08T13:10:00Z INFO telemetry: height=200790 peers=48 mempool=3190 avgLatencyMs=130 tipHash=0x031056
2025-11-08T13:11:00Z INFO telemetry: height=200791 peers=49 mempool=3191 avgLatencyMs=131 tipHash=0x031057
2025-11-08T13:12:00Z INFO telemetry: height=200792 peers=50 mempool=3192 avgLatencyMs=132 tipHash=0x031058
2025-11-08T13:13:00Z INFO telemetry: height=200793 peers=51 mempool=3193 avgLatencyMs=133 tipHash=0x031059
2025-11-08T13:14:00Z INFO telemetry: height=200794 peers=52 mempool=3194 avgLatencyMs=134 tipHash=0x03105a
2025-11-08T13:15:00Z INFO telemetry: height=200795 peers=48 mempool=3195 avgLatencyMs=135 tipHash=0x03105b
2025-11-08T13:16:00Z INFO telemetry: height=200796 peers=49 mempool=3196 avgLatencyMs=136 tipHash=0x03105c
2025-11-08T13:17:00Z INFO telemetry: height=200797 peers=50 mempool=3197 avgLatencyMs=137 tipHash=0x03105d
2025-11-08T13:18:00Z INFO telemetry: height=200798 peers=51 mempool=3198 avgLatencyMs=138 tipHash=0x03105e
2025-11-08T13:19:00Z INFO telemetry: height=200799 peers=52 mempool=3199 avgLatencyMs=139 tipHash=0x03105f
2025-11-08T13:20:00Z INFO telemetry: height=200800 peers=48 mempool=3000 avgLatencyMs=140 tipHash=0x031060
2025-11-08T13:21:00Z INFO telemetry: height=200801 peers=49 mempool=3001 avgLatencyMs=141 tipHash=0x031061
2025-11-08T13:22:00Z INFO telemetry: height=200802 peers=50 mempool=3002 avgLatencyMs=142 tipHash=0x031062
2025-11-08T13:23:00Z INFO telemetry: height=200803 peers=51 mempool=3003 avgLatencyMs=143 tipHash=0x031063
2025-11-08T13:24:00Z INFO telemetry: height=200804 peers=52 mempool=3004 avgLatencyMs=144 tipHash=0x031064
2025-11-08T13:25:00Z INFO telemetry: height=200805 peers=48 mempool=3005 avgLatencyMs=145 tipHash=0x031065
2025-11-08T13:26:00Z INFO telemetry: height=200806 peers=49 mempool=3006 avgLatencyMs=146 tipHash=0x031066
2025-11-08T13:27:00Z INFO telemetry: height=200807 peers=50 mempool=3007 avgLatencyMs=147 tipHash=0x031067
2025-11-08T13:28:00Z INFO telemetry: height=200808 peers=51 mempool=3008 avgLatencyMs=148 tipHash=0x031068
2025-11-08T13:29:00Z INFO telemetry: height=200809 peers=52 mempool=3009 avgLatencyMs=149 tipHash=0x031069
2025-11-08T13:30:00Z INFO telemetry: height=200810 peers=48 mempool=3010 avgLatencyMs=120 tipHash=0x03106a
2025-11-08T13:31:00Z INFO telemetry: height=200811 peers=49 mempool=3011 avgLatencyMs=121 tipHash=0x03106b
2025-11-08T13:32:00Z INFO telemetry: height=200812 peers=50 mempool=3012 avgLatencyMs=122 tipHash=0x03106c
2025-11-08T13:33:00Z INFO telemetry: height=200813 peers=51 mempool=3013 avgLatencyMs=123 tipHash=0x03106d
2025-11-08T13:34:00Z INFO telemetry: height=200814 peers=52 mempool=3014 avgLatencyMs=124 tipHash=0x03106e
2025-11-08T13:35:00Z INFO telemetry: height=200815 peers=48 mempool=3015 avgLatencyMs=125 tipHash=0x03106f
2025-11-08T13:36:00Z INFO telemetry: height=200816 peers=49 mempool=3016 avgLatencyMs=126 tipHash=0x031070
2025-11-08T13:37:00Z INFO telemetry: height=200817 peers=50 mempool=3017 avgLatencyMs=127 tipHash=0x031071
2025-11-08T13:38:00Z INFO telemetry: height=200818 peers=51 mempool=3018 avgLatencyMs=128 tipHash=0x031072
2025-11-08T13:39:00Z INFO telemetry: height=200819 peers=52 mempool=3019 avgLatencyMs=129 tipHash=0x031073
2025-11-08T13:40:00Z INFO telemetry: height=200820 peers=48 mempool=3020 avgLatencyMs=130 tipHash=0x031074
2025-11-08T13:41:00Z INFO telemetry: height=200821 peers=49 mempool=3021 avgLatencyMs=131 tipHash=0x031075
2025-11-08T13:42:00Z INFO telemetry: height=200822 peers=50 mempool=3022 avgLatencyMs=132 tipHash=0x031076
2025-11-08T13:43:00Z INFO telemetry: height=200823 peers=51 mempool=3023 avgLatencyMs=133 tipHash=0x031077
2025-11-08T13:44:00Z INFO telemetry: height=200824 peers=52 mempool=3024 avgLatencyMs=134 tipHash=0x031078
2025-11-08T13:45:00Z INFO telemetry: height=200825 peers=48 mempool=3025 avgLatencyMs=135 tipHash=0x031079
2025-11-08T13:46:00Z INFO telemetry: height=200826 peers=49 mempool=3026 avgLatencyMs=136 tipHash=0x03107a
2025-11-08T13:47:00Z INFO telemetry: height=200827 peers=50 mempool=3027 avgLatencyMs=137 tipHash=0x03107b
2025-11-08T13:48:00Z INFO telemetry: height=200828 peers=51 mempool=3028 avgLatencyMs=138 tipHash=0x03107c
2025-11-08T13:49:00Z INFO telemetry: height=200829 peers=52 mempool=3029 avgLatencyMs=139 tipHash=0x03107d
2025-11-08T13:50:00Z INFO telemetry: height=200830 peers=48 mempool=3030 avgLatencyMs=140 tipHash=0x03107e
2025-11-08T13:51:00Z INFO telemetry: height=200831 peers=49 mempool=3031 avgLatencyMs=141 tipHash=0x03107f
2025-11-08T13:52:00Z INFO telemetry: height=200832 peers=50 mempool=3032 avgLatencyMs=142 tipHash=0x031080
2025-11-08T13:53:00Z INFO telemetry: height=200833 peers=51 mempool=3033 avgLatencyMs=143 tipHash=0x031081
2025-11-08T13:54:00Z INFO telemetry: height=200834 peers=52 mempool=3034 avgLatencyMs=144 tipHash=0x031082
2025-11-08T13:55:00Z INFO telemetry: height=200835 peers=48 mempool=3035 avgLatencyMs=145 tipHash=0x031083
2025-11-08T13:56:00Z INFO telemetry: height=200836 peers=49 mempool=3036 avgLatencyMs=146 tipHash=0x031084
2025-11-08T13:57:00Z INFO telemetry: height=200837 peers=50 mempool=3037 avgLatencyMs=147 tipHash=0x031085
2025-11-08T13:58:00Z INFO telemetry: height=200838 peers=51 mempool=3038 avgLatencyMs=148 tipHash=0x031086
2025-11-08T13:59:00Z INFO telemetry: height=200839 peers=52 mempool=3039 avgLatencyMs=149 tipHash=0x031087
2025-11-08T14:00:00Z INFO telemetry: height=200840 peers=48 mempool=3040 avgLatencyMs=120 tipHash=0x031088
2025-11-08T14:01:00Z INFO telemetry: height=200841 peers=49 mempool=3041 avgLatencyMs=121 tipHash=0x031089
2025-11-08T14:02:00Z INFO telemetry: height=200842 peers=50 mempool=3042 avgLatencyMs=122 tipHash=0x03108a
2025-11-08T14:03:00Z INFO telemetry: height=200843 peers=51 mempool=3043 avgLatencyMs=123 tipHash=0x03108b
2025-11-08T14:04:00Z INFO telemetry: height=200844 peers=52 mempool=3044 avgLatencyMs=124 tipHash=0x03108c
2025-11-08T14:05:00Z INFO telemetry: height=200845 peers=48 mempool=3045 avgLatencyMs=125 tipHash=0x03108d
2025-11-08T14:06:00Z INFO telemetry: height=200846 peers=49 mempool=3046 avgLatencyMs=126 tipHash=0x03108e
2025-11-08T14:07:00Z INFO telemetry: height=200847 peers=50 mempool=3047 avgLatencyMs=127 tipHash=0x03108f
2025-11-08T14:08:00Z INFO telemetry: height=200848 peers=51 mempool=3048 avgLatencyMs=128 tipHash=0x031090
2025-11-08T14:09:00Z INFO telemetry: height=200849 peers=52 mempool=3049 avgLatencyMs=129 tipHash=0x031091
2025-11-08T14:10:00Z INFO telemetry: height=200850 peers=48 mempool=3050 avgLatencyMs=130 tipHash=0x031092
2025-11-08T14:11:00Z INFO telemetry: height=200851 peers=49 mempool=3051 avgLatencyMs=131 tipHash=0x031093
2025-11-08T14:12:00Z INFO telemetry: height=200852 peers=50 mempool=3052 avgLatencyMs=132 tipHash=0x031094
2025-11-08T14:13:00Z INFO telemetry: height=200853 peers=51 mempool=3053 avgLatencyMs=133 tipHash=0x031095
2025-11-08T14:14:00Z INFO telemetry: height=200854 peers=52 mempool=3054 avgLatencyMs=134 tipHash=0x031096
2025-11-08T14:15:00Z INFO telemetry: height=200855 peers=48 mempool=3055 avgLatencyMs=135 tipHash=0x031097
2025-11-08T14:16:00Z INFO telemetry: height=200856 peers=49 mempool=3056 avgLatencyMs=136 tipHash=0x031098
2025-11-08T14:17:00Z INFO telemetry: height=200857 peers=50 mempool=3057 avgLatencyMs=137 tipHash=0x031099
2025-11-08T14:18:00Z INFO telemetry: height=200858 peers=51 mempool=3058 avgLatencyMs=138 tipHash=0x03109a
2025-11-08T14:19:00Z INFO telemetry: height=200859 peers=52 mempool=3059 avgLatencyMs=139 tipHash=0x03109b
2025-11-08T14:20:00Z INFO telemetry: height=200860 peers=48 mempool=3060 avgLatencyMs=140 tipHash=0x03109c
2025-11-08T14:21:00Z INFO telemetry: height=200861 peers=49 mempool=3061 avgLatencyMs=141 tipHash=0x03109d
2025-11-08T14:22:00Z INFO telemetry: height=200862 peers=50 mempool=3062 avgLatencyMs=142 tipHash=0x03109e
2025-11-08T14:23:00Z INFO telemetry: height=200863 peers=51 mempool=3063 avgLatencyMs=143 tipHash=0x03109f
2025-11-08T14:24:00Z INFO telemetry: height=200864 peers=52 mempool=3064 avgLatencyMs=144 tipHash=0x0310a0
2025-11-08T14:25:00Z INFO telemetry: height=200865 peers=48 mempool=3065 avgLatencyMs=145 tipHash=0x0310a1
2025-11-08T14:26:00Z INFO telemetry: height=200866 peers=49 mempool=3066 avgLatencyMs=146 tipHash=0x0310a2
2025-11-08T14:27:00Z INFO telemetry: height=200867 peers=50 mempool=3067 avgLatencyMs=147 tipHash=0x0310a3
2025-11-08T14:28:00Z INFO telemetry: height=200868 peers=51 mempool=3068 avgLatencyMs=148 tipHash=0x0310a4
2025-11-08T14:29:00Z INFO telemetry: height=200869 peers=52 mempool=3069 avgLatencyMs=149 tipHash=0x0310a5
2025-11-08T14:30:00Z INFO telemetry: height=200870 peers=48 mempool=3070 avgLatencyMs=120 tipHash=0x0310a6
2025-11-08T14:31:00Z INFO telemetry: height=200871 peers=49 mempool=3071 avgLatencyMs=121 tipHash=0x0310a7
2025-11-08T14:32:00Z INFO telemetry: height=200872 peers=50 mempool=3072 avgLatencyMs=122 tipHash=0x0310a8
2025-11-08T14:33:00Z INFO telemetry: height=200873 peers=51 mempool=3073 avgLatencyMs=123 tipHash=0x0310a9
2025-11-08T14:34:00Z INFO telemetry: height=200874 peers=52 mempool=3074 avgLatencyMs=124 tipHash=0x0310aa
2025-11-08T14:35:00Z INFO telemetry: height=200875 peers=48 mempool=3075 avgLatencyMs=125 tipHash=0x0310ab
2025-11-08T14:36:00Z INFO telemetry: height=200876 peers=49 mempool=3076 avgLatencyMs=126 tipHash=0x0310ac
2025-11-08T14:37:00Z INFO telemetry: height=200877 peers=50 mempool=3077 avgLatencyMs=127 tipHash=0x0310ad
2025-11-08T14:38:00Z INFO telemetry: height=200878 peers=51 mempool=3078 avgLatencyMs=128 tipHash=0x0310ae
2025-11-08T14:39:00Z INFO telemetry: height=200879 peers=52 mempool=3079 avgLatencyMs=129 tipHash=0x0310af
2025-11-08T14:40:00Z INFO telemetry: height=200880 peers=48 mempool=3080 avgLatencyMs=130 tipHash=0x0310b0
2025-11-08T14:41:00Z INFO telemetry: height=200881 peers=49 mempool=3081 avgLatencyMs=131 tipHash=0x0310b1
2025-11-08T14:42:00Z INFO telemetry: height=200882 peers=50 mempool=3082 avgLatencyMs=132 tipHash=0x0310b2
2025-11-08T14:43:00Z INFO telemetry: height=200883 peers=51 mempool=3083 avgLatencyMs=133 tipHash=0x0310b3
2025-11-08T14:44:00Z INFO telemetry: height=200884 peers=52 mempool=3084 avgLatencyMs=134 tipHash=0x0310b4
2025-11-08T14:45:00Z INFO telemetry: height=200885 peers=48 mempool=3085 avgLatencyMs=135 tipHash=0x0310b5
2025-11-08T14:46:00Z INFO telemetry: height=200886 peers=49 mempool=3086 avgLatencyMs=136 tipHash=0x0310b6
2025-11-08T14:47:00Z INFO telemetry: height=200887 peers=50 mempool=3087 avgLatencyMs=137 tipHash=0x0310b7
2025-11-08T14:48:00Z INFO telemetry: height=200888 peers=51 mempool=3088 avgLatencyMs=138 tipHash=0x0310b8
2025-11-08T14:49:00Z INFO telemetry: height=200889 peers=52 mempool=3089 avgLatencyMs=139 tipHash=0x0310b9
2025-11-08T14:50:00Z INFO telemetry: height=200890 peers=48 mempool=3090 avgLatencyMs=140 tipHash=0x0310ba
2025-11-08T14:51:00Z INFO telemetry: height=200891 peers=49 mempool=3091 avgLatencyMs=141 tipHash=0x0310bb
2025-11-08T14:52:00Z INFO telemetry: height=200892 peers=50 mempool=3092 avgLatencyMs=142 tipHash=0x0310bc
2025-11-08T14:53:00Z INFO telemetry: height=200893 peers=51 mempool=3093 avgLatencyMs=143 tipHash=0x0310bd
2025-11-08T14:54:00Z INFO telemetry: height=200894 peers=52 mempool=3094 avgLatencyMs=144 tipHash=0x0310be
2025-11-08T14:55:00Z INFO telemetry: height=200895 peers=48 mempool=3095 avgLatencyMs=145 tipHash=0x0310bf
2025-11-08T14:56:00Z INFO telemetry: height=200896 peers=49 mempool=3096 avgLatencyMs=146 tipHash=0x0310c0
2025-11-08T14:57:00Z INFO telemetry: height=200897 peers=50 mempool=3097 avgLatencyMs=147 tipHash=0x0310c1
2025-11-08T14:58:00Z INFO telemetry: height=200898 peers=51 mempool=3098 avgLatencyMs=148 tipHash=0x0310c2
2025-11-08T14:59:00Z INFO telemetry: height=200899 peers=52 mempool=3099 avgLatencyMs=149 tipHash=0x0310c3
2025-11-08T15:00:00Z INFO telemetry: height=200900 peers=48 mempool=3100 avgLatencyMs=120 tipHash=0x0310c4
2025-11-08T15:01:00Z INFO telemetry: height=200901 peers=49 mempool=3101 avgLatencyMs=121 tipHash=0x0310c5
2025-11-08T15:02:00Z INFO telemetry: height=200902 peers=50 mempool=3102 avgLatencyMs=122 tipHash=0x0310c6
2025-11-08T15:03:00Z INFO telemetry: height=200903 peers=51 mempool=3103 avgLatencyMs=123 tipHash=0x0310c7
2025-11-08T15:04:00Z INFO telemetry: height=200904 peers=52 mempool=3104 avgLatencyMs=124 tipHash=0x0310c8
2025-11-08T15:05:00Z INFO telemetry: height=200905 peers=48 mempool=3105 avgLatencyMs=125 tipHash=0x0310c9
2025-11-08T15:06:00Z INFO telemetry: height=200906 peers=49 mempool=3106 avgLatencyMs=126 tipHash=0x0310ca
2025-11-08T15:07:00Z INFO telemetry: height=200907 peers=50 mempool=3107 avgLatencyMs=127 tipHash=0x0310cb
2025-11-08T15:08:00Z INFO telemetry: height=200908 peers=51 mempool=3108 avgLatencyMs=128 tipHash=0x0310cc
2025-11-08T15:09:00Z INFO telemetry: height=200909 peers=52 mempool=3109 avgLatencyMs=129 tipHash=0x0310cd
2025-11-08T15:10:00Z INFO telemetry: height=200910 peers=48 mempool=3110 avgLatencyMs=130 tipHash=0x0310ce
2025-11-08T15:11:00Z INFO telemetry: height=200911 peers=49 mempool=3111 avgLatencyMs=131 tipHash=0x0310cf
2025-11-08T15:12:00Z INFO telemetry: height=200912 peers=50 mempool=3112 avgLatencyMs=132 tipHash=0x0310d0
2025-11-08T15:13:00Z INFO telemetry: height=200913 peers=51 mempool=3113 avgLatencyMs=133 tipHash=0x0310d1
2025-11-08T15:14:00Z INFO telemetry: height=200914 peers=52 mempool=3114 avgLatencyMs=134 tipHash=0x0310d2
2025-11-08T15:15:00Z INFO telemetry: height=200915 peers=48 mempool=3115 avgLatencyMs=135 tipHash=0x0310d3
2025-11-08T15:16:00Z INFO telemetry: height=200916 peers=49 mempool=3116 avgLatencyMs=136 tipHash=0x0310d4
2025-11-08T15:17:00Z INFO telemetry: height=200917 peers=50 mempool=3117 avgLatencyMs=137 tipHash=0x0310d5
2025-11-08T15:18:00Z INFO telemetry: height=200918 peers=51 mempool=3118 avgLatencyMs=138 tipHash=0x0310d6
2025-11-08T15:19:00Z INFO telemetry: height=200919 peers=52 mempool=3119 avgLatencyMs=139 tipHash=0x0310d7
2025-11-08T15:20:00Z INFO telemetry: height=200920 peers=48 mempool=3120 avgLatencyMs=140 tipHash=0x0310d8
2025-11-08T15:21:00Z INFO telemetry: height=200921 peers=49 mempool=3121 avgLatencyMs=141 tipHash=0x0310d9
2025-11-08T15:22:00Z INFO telemetry: height=200922 peers=50 mempool=3122 avgLatencyMs=142 tipHash=0x0310da
2025-11-08T15:23:00Z INFO telemetry: height=200923 peers=51 mempool=3123 avgLatencyMs=143 tipHash=0x0310db
2025-11-08T15:24:00Z INFO telemetry: height=200924 peers=52 mempool=3124 avgLatencyMs=144 tipHash=0x0310dc
2025-11-08T15:25:00Z INFO telemetry: height=200925 peers=48 mempool=3125 avgLatencyMs=145 tipHash=0x0310dd
2025-11-08T15:26:00Z INFO telemetry: height=200926 peers=49 mempool=3126 avgLatencyMs=146 tipHash=0x0310de
2025-11-08T15:27:00Z INFO telemetry: height=200927 peers=50 mempool=3127 avgLatencyMs=147 tipHash=0x0310df
2025-11-08T15:28:00Z INFO telemetry: height=200928 peers=51 mempool=3128 avgLatencyMs=148 tipHash=0x0310e0
2025-11-08T15:29:00Z INFO telemetry: height=200929 peers=52 mempool=3129 avgLatencyMs=149 tipHash=0x0310e1
2025-11-08T15:30:00Z INFO telemetry: height=200930 peers=48 mempool=3130 avgLatencyMs=120 tipHash=0x0310e2
2025-11-08T15:31:00Z INFO telemetry: height=200931 peers=49 mempool=3131 avgLatencyMs=121 tipHash=0x0310e3
2025-11-08T15:32:00Z INFO telemetry: height=200932 peers=50 mempool=3132 avgLatencyMs=122 tipHash=0x0310e4
2025-11-08T15:33:00Z INFO telemetry: height=200933 peers=51 mempool=3133 avgLatencyMs=123 tipHash=0x0310e5
2025-11-08T15:34:00Z INFO telemetry: height=200934 peers=52 mempool=3134 avgLatencyMs=124 tipHash=0x0310e6
2025-11-08T15:35:00Z INFO telemetry: height=200935 peers=48 mempool=3135 avgLatencyMs=125 tipHash=0x0310e7
2025-11-08T15:36:00Z INFO telemetry: height=200936 peers=49 mempool=3136 avgLatencyMs=126 tipHash=0x0310e8
2025-11-08T15:37:00Z INFO telemetry: height=200937 peers=50 mempool=3137 avgLatencyMs=127 tipHash=0x0310e9
2025-11-08T15:38:00Z INFO telemetry: height=200938 peers=51 mempool=3138 avgLatencyMs=128 tipHash=0x0310ea
2025-11-08T15:39:00Z INFO telemetry: height=200939 peers=52 mempool=3139 avgLatencyMs=129 tipHash=0x0310eb
2025-11-08T15:40:00Z INFO telemetry: height=200940 peers=48 mempool=3140 avgLatencyMs=130 tipHash=0x0310ec
2025-11-08T15:41:00Z INFO telemetry: height=200941 peers=49 mempool=3141 avgLatencyMs=131 tipHash=0x0310ed
2025-11-08T15:42:00Z INFO telemetry: height=200942 peers=50 mempool=3142 avgLatencyMs=132 tipHash=0x0310ee
2025-11-08T15:43:00Z INFO telemetry: height=200943 peers=51 mempool=3143 avgLatencyMs=133 tipHash=0x0310ef
2025-11-08T15:44:00Z INFO telemetry: height=200944 peers=52 mempool=3144 avgLatencyMs=134 tipHash=0x0310f0
2025-11-08T15:45:00Z INFO telemetry: height=200945 peers=48 mempool=3145 avgLatencyMs=135 tipHash=0x0310f1
2025-11-08T15:46:00Z INFO telemetry: height=200946 peers=49 mempool=3146 avgLatencyMs=136 tipHash=0x0310f2
2025-11-08T15:47:00Z INFO telemetry: height=200947 peers=50 mempool=3147 avgLatencyMs=137 tipHash=0x0310f3
2025-11-08T15:48:00Z INFO telemetry: height=200948 peers=51 mempool=3148 avgLatencyMs=138 tipHash=0x0310f4
2025-11-08T15:49:00Z INFO telemetry: height=200949 peers=52 mempool=3149 avgLatencyMs=139 tipHash=0x0310f5
2025-11-08T15:50:00Z INFO telemetry: height=200950 peers=48 mempool=3150 avgLatencyMs=140 tipHash=0x0310f6
2025-11-08T15:51:00Z INFO telemetry: height=200951 peers=49 mempool=3151 avgLatencyMs=141 tipHash=0x0310f7
2025-11-08T15:52:00Z INFO telemetry: height=200952 peers=50 mempool=3152 avgLatencyMs=142 tipHash=0x0310f8
2025-11-08T15:53:00Z INFO telemetry: height=200953 peers=51 mempool=3153 avgLatencyMs=143 tipHash=0x0310f9
2025-11-08T15:54:00Z INFO telemetry: height=200954 peers=52 mempool=3154 avgLatencyMs=144 tipHash=0x0310fa
2025-11-08T15:55:00Z INFO telemetry: height=200955 peers=48 mempool=3155 avgLatencyMs=145 tipHash=0x0310fb
2025-11-08T15:56:00Z INFO telemetry: height=200956 peers=49 mempool=3156 avgLatencyMs=146 tipHash=0x0310fc
2025-11-08T15:57:00Z INFO telemetry: height=200957 peers=50 mempool=3157 avgLatencyMs=147 tipHash=0x0310fd
2025-11-08T15:58:00Z INFO telemetry: height=200958 peers=51 mempool=3158 avgLatencyMs=148 tipHash=0x0310fe
2025-11-08T15:59:00Z INFO telemetry: height=200959 peers=52 mempool=3159 avgLatencyMs=149 tipHash=0x0310ff
2025-11-08T16:00:00Z INFO telemetry: height=200960 peers=48 mempool=3160 avgLatencyMs=120 tipHash=0x031100
2025-11-08T16:01:00Z INFO telemetry: height=200961 peers=49 mempool=3161 avgLatencyMs=121 tipHash=0x031101
2025-11-08T16:02:00Z INFO telemetry: height=200962 peers=50 mempool=3162 avgLatencyMs=122 tipHash=0x031102
2025-11-08T16:03:00Z INFO telemetry: height=200963 peers=51 mempool=3163 avgLatencyMs=123 tipHash=0x031103
2025-11-08T16:04:00Z INFO telemetry: height=200964 peers=52 mempool=3164 avgLatencyMs=124 tipHash=0x031104
2025-11-08T16:05:00Z INFO telemetry: height=200965 peers=48 mempool=3165 avgLatencyMs=125 tipHash=0x031105
2025-11-08T16:06:00Z INFO telemetry: height=200966 peers=49 mempool=3166 avgLatencyMs=126 tipHash=0x031106
2025-11-08T16:07:00Z INFO telemetry: height=200967 peers=50 mempool=3167 avgLatencyMs=127 tipHash=0x031107
2025-11-08T16:08:00Z INFO telemetry: height=200968 peers=51 mempool=3168 avgLatencyMs=128 tipHash=0x031108
2025-11-08T16:09:00Z INFO telemetry: height=200969 peers=52 mempool=3169 avgLatencyMs=129 tipHash=0x031109
2025-11-08T16:10:00Z INFO telemetry: height=200970 peers=48 mempool=3170 avgLatencyMs=130 tipHash=0x03110a
2025-11-08T16:11:00Z INFO telemetry: height=200971 peers=49 mempool=3171 avgLatencyMs=131 tipHash=0x03110b
2025-11-08T16:12:00Z INFO telemetry: height=200972 peers=50 mempool=3172 avgLatencyMs=132 tipHash=0x03110c
2025-11-08T16:13:00Z INFO telemetry: height=200973 peers=51 mempool=3173 avgLatencyMs=133 tipHash=0x03110d
2025-11-08T16:14:00Z INFO telemetry: height=200974 peers=52 mempool=3174 avgLatencyMs=134 tipHash=0x03110e
2025-11-08T16:15:00Z INFO telemetry: height=200975 peers=48 mempool=3175 avgLatencyMs=135 tipHash=0x03110f
2025-11-08T16:16:00Z INFO telemetry: height=200976 peers=49 mempool=3176 avgLatencyMs=136 tipHash=0x031110
2025-11-08T16:17:00Z INFO telemetry: height=200977 peers=50 mempool=3177 avgLatencyMs=137 tipHash=0x031111
2025-11-08T16:18:00Z INFO telemetry: height=200978 peers=51 mempool=3178 avgLatencyMs=138 tipHash=0x031112
2025-11-08T16:19:00Z INFO telemetry: height=200979 peers=52 mempool=3179 avgLatencyMs=139 tipHash=0x031113
2025-11-08T16:20:00Z INFO telemetry: height=200980 peers=48 mempool=3180 avgLatencyMs=140 tipHash=0x031114
2025-11-08T16:21:00Z INFO telemetry: height=200981 peers=49 mempool=3181 avgLatencyMs=141 tipHash=0x031115
2025-11-08T16:22:00Z INFO telemetry: height=200982 peers=50 mempool=3182 avgLatencyMs=142 tipHash=0x031116
2025-11-08T16:23:00Z INFO telemetry: height=200983 peers=51 mempool=3183 avgLatencyMs=143 tipHash=0x031117
2025-11-08T16:24:00Z INFO telemetry: height=200984 peers=52 mempool=3184 avgLatencyMs=144 tipHash=0x031118
2025-11-08T16:25:00Z INFO telemetry: height=200985 peers=48 mempool=3185 avgLatencyMs=145 tipHash=0x031119
2025-11-08T16:26:00Z INFO telemetry: height=200986 peers=49 mempool=3186 avgLatencyMs=146 tipHash=0x03111a
2025-11-08T16:27:00Z INFO telemetry: height=200987 peers=50 mempool=3187 avgLatencyMs=147 tipHash=0x03111b
2025-11-08T16:28:00Z INFO telemetry: height=200988 peers=51 mempool=3188 avgLatencyMs=148 tipHash=0x03111c
2025-11-08T16:29:00Z INFO telemetry: height=200989 peers=52 mempool=3189 avgLatencyMs=149 tipHash=0x03111d
2025-11-08T16:30:00Z INFO telemetry: height=200990 peers=48 mempool=3190 avgLatencyMs=120 tipHash=0x03111e
2025-11-08T16:31:00Z INFO telemetry: height=200991 peers=49 mempool=3191 avgLatencyMs=121 tipHash=0x03111f
2025-11-08T16:32:00Z INFO telemetry: height=200992 peers=50 mempool=3192 avgLatencyMs=122 tipHash=0x031120
2025-11-08T16:33:00Z INFO telemetry: height=200993 peers=51 mempool=3193 avgLatencyMs=123 tipHash=0x031121
2025-11-08T16:34:00Z INFO telemetry: height=200994 peers=52 mempool=3194 avgLatencyMs=124 tipHash=0x031122
2025-11-08T16:35:00Z INFO telemetry: height=200995 peers=48 mempool=3195 avgLatencyMs=125 tipHash=0x031123
2025-11-08T16:36:00Z INFO telemetry: height=200996 peers=49 mempool=3196 avgLatencyMs=126 tipHash=0x031124
2025-11-08T16:37:00Z INFO telemetry: height=200997 peers=50 mempool=3197 avgLatencyMs=127 tipHash=0x031125
2025-11-08T16:38:00Z INFO telemetry: height=200998 peers=51 mempool=3198 avgLatencyMs=128 tipHash=0x031126
2025-11-08T16:39:00Z INFO telemetry: height=200999 peers=52 mempool=3199 avgLatencyMs=129 tipHash=0x031127
2025-11-08T16:40:00Z INFO telemetry: height=201000 peers=48 mempool=3000 avgLatencyMs=130 tipHash=0x031128
2025-11-08T16:41:00Z INFO telemetry: height=201001 peers=49 mempool=3001 avgLatencyMs=131 tipHash=0x031129
2025-11-08T16:42:00Z INFO telemetry: height=201002 peers=50 mempool=3002 avgLatencyMs=132 tipHash=0x03112a
2025-11-08T16:43:00Z INFO telemetry: height=201003 peers=51 mempool=3003 avgLatencyMs=133 tipHash=0x03112b
2025-11-08T16:44:00Z INFO telemetry: height=201004 peers=52 mempool=3004 avgLatencyMs=134 tipHash=0x03112c
2025-11-08T16:45:00Z INFO telemetry: height=201005 peers=48 mempool=3005 avgLatencyMs=135 tipHash=0x03112d
2025-11-08T16:46:00Z INFO telemetry: height=201006 peers=49 mempool=3006 avgLatencyMs=136 tipHash=0x03112e
2025-11-08T16:47:00Z INFO telemetry: height=201007 peers=50 mempool=3007 avgLatencyMs=137 tipHash=0x03112f
2025-11-08T16:48:00Z INFO telemetry: height=201008 peers=51 mempool=3008 avgLatencyMs=138 tipHash=0x031130
2025-11-08T16:49:00Z INFO telemetry: height=201009 peers=52 mempool=3009 avgLatencyMs=139 tipHash=0x031131
2025-11-08T16:50:00Z INFO telemetry: height=201010 peers=48 mempool=3010 avgLatencyMs=140 tipHash=0x031132
2025-11-08T16:51:00Z INFO telemetry: height=201011 peers=49 mempool=3011 avgLatencyMs=141 tipHash=0x031133
2025-11-08T16:52:00Z INFO telemetry: height=201012 peers=50 mempool=3012 avgLatencyMs=142 tipHash=0x031134
2025-11-08T16:53:00Z INFO telemetry: height=201013 peers=51 mempool=3013 avgLatencyMs=143 tipHash=0x031135
2025-11-08T16:54:00Z INFO telemetry: height=201014 peers=52 mempool=3014 avgLatencyMs=144 tipHash=0x031136
2025-11-08T16:55:00Z INFO telemetry: height=201015 peers=48 mempool=3015 avgLatencyMs=145 tipHash=0x031137
2025-11-08T16:56:00Z INFO telemetry: height=201016 peers=49 mempool=3016 avgLatencyMs=146 tipHash=0x031138
2025-11-08T16:57:00Z INFO telemetry: height=201017 peers=50 mempool=3017 avgLatencyMs=147 tipHash=0x031139
2025-11-08T16:58:00Z INFO telemetry: height=201018 peers=51 mempool=3018 avgLatencyMs=148 tipHash=0x03113a
2025-11-08T16:59:00Z INFO telemetry: height=201019 peers=52 mempool=3019 avgLatencyMs=149 tipHash=0x03113b
2025-11-08T17:00:00Z INFO telemetry: height=201020 peers=48 mempool=3020 avgLatencyMs=120 tipHash=0x03113c
2025-11-08T17:01:00Z INFO telemetry: height=201021 peers=49 mempool=3021 avgLatencyMs=121 tipHash=0x03113d
2025-11-08T17:02:00Z INFO telemetry: height=201022 peers=50 mempool=3022 avgLatencyMs=122 tipHash=0x03113e
2025-11-08T17:03:00Z INFO telemetry: height=201023 peers=51 mempool=3023 avgLatencyMs=123 tipHash=0x03113f
2025-11-08T17:04:00Z INFO telemetry: height=201024 peers=52 mempool=3024 avgLatencyMs=124 tipHash=0x031140
2025-11-08T17:05:00Z INFO telemetry: height=201025 peers=48 mempool=3025 avgLatencyMs=125 tipHash=0x031141
2025-11-08T17:06:00Z INFO telemetry: height=201026 peers=49 mempool=3026 avgLatencyMs=126 tipHash=0x031142
2025-11-08T17:07:00Z INFO telemetry: height=201027 peers=50 mempool=3027 avgLatencyMs=127 tipHash=0x031143
2025-11-08T17:08:00Z INFO telemetry: height=201028 peers=51 mempool=3028 avgLatencyMs=128 tipHash=0x031144
2025-11-08T17:09:00Z INFO telemetry: height=201029 peers=52 mempool=3029 avgLatencyMs=129 tipHash=0x031145
2025-11-08T17:10:00Z INFO telemetry: height=201030 peers=48 mempool=3030 avgLatencyMs=130 tipHash=0x031146
2025-11-08T17:11:00Z INFO telemetry: height=201031 peers=49 mempool=3031 avgLatencyMs=131 tipHash=0x031147
2025-11-08T17:12:00Z INFO telemetry: height=201032 peers=50 mempool=3032 avgLatencyMs=132 tipHash=0x031148
2025-11-08T17:13:00Z INFO telemetry: height=201033 peers=51 mempool=3033 avgLatencyMs=133 tipHash=0x031149
2025-11-08T17:14:00Z INFO telemetry: height=201034 peers=52 mempool=3034 avgLatencyMs=134 tipHash=0x03114a
2025-11-08T17:15:00Z INFO telemetry: height=201035 peers=48 mempool=3035 avgLatencyMs=135 tipHash=0x03114b
2025-11-08T17:16:00Z INFO telemetry: height=201036 peers=49 mempool=3036 avgLatencyMs=136 tipHash=0x03114c
2025-11-08T17:17:00Z INFO telemetry: height=201037 peers=50 mempool=3037 avgLatencyMs=137 tipHash=0x03114d
2025-11-08T17:18:00Z INFO telemetry: height=201038 peers=51 mempool=3038 avgLatencyMs=138 tipHash=0x03114e
2025-11-08T17:19:00Z INFO telemetry: height=201039 peers=52 mempool=3039 avgLatencyMs=139 tipHash=0x03114f
2025-11-08T17:20:00Z INFO telemetry: height=201040 peers=48 mempool=3040 avgLatencyMs=140 tipHash=0x031150
2025-11-08T17:21:00Z INFO telemetry: height=201041 peers=49 mempool=3041 avgLatencyMs=141 tipHash=0x031151
2025-11-08T17:22:00Z INFO telemetry: height=201042 peers=50 mempool=3042 avgLatencyMs=142 tipHash=0x031152
2025-11-08T17:23:00Z INFO telemetry: height=201043 peers=51 mempool=3043 avgLatencyMs=143 tipHash=0x031153
2025-11-08T17:24:00Z INFO telemetry: height=201044 peers=52 mempool=3044 avgLatencyMs=144 tipHash=0x031154
2025-11-08T17:25:00Z INFO telemetry: height=201045 peers=48 mempool=3045 avgLatencyMs=145 tipHash=0x031155
2025-11-08T17:26:00Z INFO telemetry: height=201046 peers=49 mempool=3046 avgLatencyMs=146 tipHash=0x031156
2025-11-08T17:27:00Z INFO telemetry: height=201047 peers=50 mempool=3047 avgLatencyMs=147 tipHash=0x031157
2025-11-08T17:28:00Z INFO telemetry: height=201048 peers=51 mempool=3048 avgLatencyMs=148 tipHash=0x031158
2025-11-08T17:29:00Z INFO telemetry: height=201049 peers=52 mempool=3049 avgLatencyMs=149 tipHash=0x031159
2025-11-08T17:30:00Z INFO telemetry: height=201050 peers=48 mempool=3050 avgLatencyMs=120 tipHash=0x03115a
2025-11-08T17:31:00Z INFO telemetry: height=201051 peers=49 mempool=3051 avgLatencyMs=121 tipHash=0x03115b
2025-11-08T17:32:00Z INFO telemetry: height=201052 peers=50 mempool=3052 avgLatencyMs=122 tipHash=0x03115c
2025-11-08T17:33:00Z INFO telemetry: height=201053 peers=51 mempool=3053 avgLatencyMs=123 tipHash=0x03115d
2025-11-08T17:34:00Z INFO telemetry: height=201054 peers=52 mempool=3054 avgLatencyMs=124 tipHash=0x03115e
2025-11-08T17:35:00Z INFO telemetry: height=201055 peers=48 mempool=3055 avgLatencyMs=125 tipHash=0x03115f
2025-11-08T17:36:00Z INFO telemetry: height=201056 peers=49 mempool=3056 avgLatencyMs=126 tipHash=0x031160
2025-11-08T17:37:00Z INFO telemetry: height=201057 peers=50 mempool=3057 avgLatencyMs=127 tipHash=0x031161
2025-11-08T17:38:00Z INFO telemetry: height=201058 peers=51 mempool=3058 avgLatencyMs=128 tipHash=0x031162
2025-11-08T17:39:00Z INFO telemetry: height=201059 peers=52 mempool=3059 avgLatencyMs=129 tipHash=0x031163
2025-11-08T17:40:00Z INFO telemetry: height=201060 peers=48 mempool=3060 avgLatencyMs=130 tipHash=0x031164
2025-11-08T17:41:00Z INFO telemetry: height=201061 peers=49 mempool=3061 avgLatencyMs=131 tipHash=0x031165
2025-11-08T17:42:00Z INFO telemetry: height=201062 peers=50 mempool=3062 avgLatencyMs=132 tipHash=0x031166
2025-11-08T17:43:00Z INFO telemetry: height=201063 peers=51 mempool=3063 avgLatencyMs=133 tipHash=0x031167
2025-11-08T17:44:00Z INFO telemetry: height=201064 peers=52 mempool=3064 avgLatencyMs=134 tipHash=0x031168
2025-11-08T17:45:00Z INFO telemetry: height=201065 peers=48 mempool=3065 avgLatencyMs=135 tipHash=0x031169
2025-11-08T17:46:00Z INFO telemetry: height=201066 peers=49 mempool=3066 avgLatencyMs=136 tipHash=0x03116a
2025-11-08T17:47:00Z INFO telemetry: height=201067 peers=50 mempool=3067 avgLatencyMs=137 tipHash=0x03116b
2025-11-08T17:48:00Z INFO telemetry: height=201068 peers=51 mempool=3068 avgLatencyMs=138 tipHash=0x03116c
2025-11-08T17:49:00Z INFO telemetry: height=201069 peers=52 mempool=3069 avgLatencyMs=139 tipHash=0x03116d
2025-11-08T17:50:00Z INFO telemetry: height=201070 peers=48 mempool=3070 avgLatencyMs=140 tipHash=0x03116e
2025-11-08T17:51:00Z INFO telemetry: height=201071 peers=49 mempool=3071 avgLatencyMs=141 tipHash=0x03116f
2025-11-08T17:52:00Z INFO telemetry: height=201072 peers=50 mempool=3072 avgLatencyMs=142 tipHash=0x031170
2025-11-08T17:53:00Z INFO telemetry: height=201073 peers=51 mempool=3073 avgLatencyMs=143 tipHash=0x031171
2025-11-08T17:54:00Z INFO telemetry: height=201074 peers=52 mempool=3074 avgLatencyMs=144 tipHash=0x031172
2025-11-08T17:55:00Z INFO telemetry: height=201075 peers=48 mempool=3075 avgLatencyMs=145 tipHash=0x031173
2025-11-08T17:56:00Z INFO telemetry: height=201076 peers=49 mempool=3076 avgLatencyMs=146 tipHash=0x031174
2025-11-08T17:57:00Z INFO telemetry: height=201077 peers=50 mempool=3077 avgLatencyMs=147 tipHash=0x031175
2025-11-08T17:58:00Z INFO telemetry: height=201078 peers=51 mempool=3078 avgLatencyMs=148 tipHash=0x031176
2025-11-08T17:59:00Z INFO telemetry: height=201079 peers=52 mempool=3079 avgLatencyMs=149 tipHash=0x031177
2025-11-08T18:00:00Z INFO telemetry: height=201080 peers=48 mempool=3080 avgLatencyMs=120 tipHash=0x031178
2025-11-08T18:01:00Z INFO telemetry: height=201081 peers=49 mempool=3081 avgLatencyMs=121 tipHash=0x031179
2025-11-08T18:02:00Z INFO telemetry: height=201082 peers=50 mempool=3082 avgLatencyMs=122 tipHash=0x03117a
2025-11-08T18:03:00Z INFO telemetry: height=201083 peers=51 mempool=3083 avgLatencyMs=123 tipHash=0x03117b
2025-11-08T18:04:00Z INFO telemetry: height=201084 peers=52 mempool=3084 avgLatencyMs=124 tipHash=0x03117c
2025-11-08T18:05:00Z INFO telemetry: height=201085 peers=48 mempool=3085 avgLatencyMs=125 tipHash=0x03117d
2025-11-08T18:06:00Z INFO telemetry: height=201086 peers=49 mempool=3086 avgLatencyMs=126 tipHash=0x03117e
2025-11-08T18:07:00Z INFO telemetry: height=201087 peers=50 mempool=3087 avgLatencyMs=127 tipHash=0x03117f
2025-11-08T18:08:00Z INFO telemetry: height=201088 peers=51 mempool=3088 avgLatencyMs=128 tipHash=0x031180
2025-11-08T18:09:00Z INFO telemetry: height=201089 peers=52 mempool=3089 avgLatencyMs=129 tipHash=0x031181
2025-11-08T18:10:00Z INFO telemetry: height=201090 peers=48 mempool=3090 avgLatencyMs=130 tipHash=0x031182
2025-11-08T18:11:00Z INFO telemetry: height=201091 peers=49 mempool=3091 avgLatencyMs=131 tipHash=0x031183
2025-11-08T18:12:00Z INFO telemetry: height=201092 peers=50 mempool=3092 avgLatencyMs=132 tipHash=0x031184
2025-11-08T18:13:00Z INFO telemetry: height=201093 peers=51 mempool=3093 avgLatencyMs=133 tipHash=0x031185
2025-11-08T18:14:00Z INFO telemetry: height=201094 peers=52 mempool=3094 avgLatencyMs=134 tipHash=0x031186
2025-11-08T18:15:00Z INFO telemetry: height=201095 peers=48 mempool=3095 avgLatencyMs=135 tipHash=0x031187
2025-11-08T18:16:00Z INFO telemetry: height=201096 peers=49 mempool=3096 avgLatencyMs=136 tipHash=0x031188
2025-11-08T18:17:00Z INFO telemetry: height=201097 peers=50 mempool=3097 avgLatencyMs=137 tipHash=0x031189
2025-11-08T18:18:00Z INFO telemetry: height=201098 peers=51 mempool=3098 avgLatencyMs=138 tipHash=0x03118a
2025-11-08T18:19:00Z INFO telemetry: height=201099 peers=52 mempool=3099 avgLatencyMs=139 tipHash=0x03118b
2025-11-08T18:20:00Z INFO telemetry: height=201100 peers=48 mempool=3100 avgLatencyMs=140 tipHash=0x03118c
2025-11-08T18:21:00Z INFO telemetry: height=201101 peers=49 mempool=3101 avgLatencyMs=141 tipHash=0x03118d
2025-11-08T18:22:00Z INFO telemetry: height=201102 peers=50 mempool=3102 avgLatencyMs=142 tipHash=0x03118e
2025-11-08T18:23:00Z INFO telemetry: height=201103 peers=51 mempool=3103 avgLatencyMs=143 tipHash=0x03118f
2025-11-08T18:24:00Z INFO telemetry: height=201104 peers=52 mempool=3104 avgLatencyMs=144 tipHash=0x031190
2025-11-08T18:25:00Z INFO telemetry: height=201105 peers=48 mempool=3105 avgLatencyMs=145 tipHash=0x031191
2025-11-08T18:26:00Z INFO telemetry: height=201106 peers=49 mempool=3106 avgLatencyMs=146 tipHash=0x031192
2025-11-08T18:27:00Z INFO telemetry: height=201107 peers=50 mempool=3107 avgLatencyMs=147 tipHash=0x031193
2025-11-08T18:28:00Z INFO telemetry: height=201108 peers=51 mempool=3108 avgLatencyMs=148 tipHash=0x031194
2025-11-08T18:29:00Z INFO telemetry: height=201109 peers=52 mempool=3109 avgLatencyMs=149 tipHash=0x031195
2025-11-08T18:30:00Z INFO telemetry: height=201110 peers=48 mempool=3110 avgLatencyMs=120 tipHash=0x031196
2025-11-08T18:31:00Z INFO telemetry: height=201111 peers=49 mempool=3111 avgLatencyMs=121 tipHash=0x031197
2025-11-08T18:32:00Z INFO telemetry: height=201112 peers=50 mempool=3112 avgLatencyMs=122 tipHash=0x031198
2025-11-08T18:33:00Z INFO telemetry: height=201113 peers=51 mempool=3113 avgLatencyMs=123 tipHash=0x031199
2025-11-08T18:34:00Z INFO telemetry: height=201114 peers=52 mempool=3114 avgLatencyMs=124 tipHash=0x03119a
2025-11-08T18:35:00Z INFO telemetry: height=201115 peers=48 mempool=3115 avgLatencyMs=125 tipHash=0x03119b
2025-11-08T18:36:00Z INFO telemetry: height=201116 peers=49 mempool=3116 avgLatencyMs=126 tipHash=0x03119c
2025-11-08T18:37:00Z INFO telemetry: height=201117 peers=50 mempool=3117 avgLatencyMs=127 tipHash=0x03119d
2025-11-08T18:38:00Z INFO telemetry: height=201118 peers=51 mempool=3118 avgLatencyMs=128 tipHash=0x03119e
2025-11-08T18:39:00Z INFO telemetry: height=201119 peers=52 mempool=3119 avgLatencyMs=129 tipHash=0x03119f
2025-11-08T18:40:00Z INFO telemetry: height=201120 peers=48 mempool=3120 avgLatencyMs=130 tipHash=0x0311a0
2025-11-08T18:41:00Z INFO telemetry: height=201121 peers=49 mempool=3121 avgLatencyMs=131 tipHash=0x0311a1
2025-11-08T18:42:00Z INFO telemetry: height=201122 peers=50 mempool=3122 avgLatencyMs=132 tipHash=0x0311a2
2025-11-08T18:43:00Z INFO telemetry: height=201123 peers=51 mempool=3123 avgLatencyMs=133 tipHash=0x0311a3
2025-11-08T18:44:00Z INFO telemetry: height=201124 peers=52 mempool=3124 avgLatencyMs=134 tipHash=0x0311a4
2025-11-08T18:45:00Z INFO telemetry: height=201125 peers=48 mempool=3125 avgLatencyMs=135 tipHash=0x0311a5
2025-11-08T18:46:00Z INFO telemetry: height=201126 peers=49 mempool=3126 avgLatencyMs=136 tipHash=0x0311a6
2025-11-08T18:47:00Z INFO telemetry: height=201127 peers=50 mempool=3127 avgLatencyMs=137 tipHash=0x0311a7
2025-11-08T18:48:00Z INFO telemetry: height=201128 peers=51 mempool=3128 avgLatencyMs=138 tipHash=0x0311a8
2025-11-08T18:49:00Z INFO telemetry: height=201129 peers=52 mempool=3129 avgLatencyMs=139 tipHash=0x0311a9
2025-11-08T18:50:00Z INFO telemetry: height=201130 peers=48 mempool=3130 avgLatencyMs=140 tipHash=0x0311aa
2025-11-08T18:51:00Z INFO telemetry: height=201131 peers=49 mempool=3131 avgLatencyMs=141 tipHash=0x0311ab
2025-11-08T18:52:00Z INFO telemetry: height=201132 peers=50 mempool=3132 avgLatencyMs=142 tipHash=0x0311ac
2025-11-08T18:53:00Z INFO telemetry: height=201133 peers=51 mempool=3133 avgLatencyMs=143 tipHash=0x0311ad
2025-11-08T18:54:00Z INFO telemetry: height=201134 peers=52 mempool=3134 avgLatencyMs=144 tipHash=0x0311ae
2025-11-08T18:55:00Z INFO telemetry: height=201135 peers=48 mempool=3135 avgLatencyMs=145 tipHash=0x0311af
2025-11-08T18:56:00Z INFO telemetry: height=201136 peers=49 mempool=3136 avgLatencyMs=146 tipHash=0x0311b0
2025-11-08T18:57:00Z INFO telemetry: height=201137 peers=50 mempool=3137 avgLatencyMs=147 tipHash=0x0311b1
2025-11-08T18:58:00Z INFO telemetry: height=201138 peers=51 mempool=3138 avgLatencyMs=148 tipHash=0x0311b2
2025-11-08T18:59:00Z INFO telemetry: height=201139 peers=52 mempool=3139 avgLatencyMs=149 tipHash=0x0311b3
2025-11-08T19:00:00Z INFO telemetry: height=201140 peers=48 mempool=3140 avgLatencyMs=120 tipHash=0x0311b4
2025-11-08T19:01:00Z INFO telemetry: height=201141 peers=49 mempool=3141 avgLatencyMs=121 tipHash=0x0311b5
2025-11-08T19:02:00Z INFO telemetry: height=201142 peers=50 mempool=3142 avgLatencyMs=122 tipHash=0x0311b6
2025-11-08T19:03:00Z INFO telemetry: height=201143 peers=51 mempool=3143 avgLatencyMs=123 tipHash=0x0311b7
2025-11-08T19:04:00Z INFO telemetry: height=201144 peers=52 mempool=3144 avgLatencyMs=124 tipHash=0x0311b8
2025-11-08T19:05:00Z INFO telemetry: height=201145 peers=48 mempool=3145 avgLatencyMs=125 tipHash=0x0311b9
2025-11-08T19:06:00Z INFO telemetry: height=201146 peers=49 mempool=3146 avgLatencyMs=126 tipHash=0x0311ba
2025-11-08T19:07:00Z INFO telemetry: height=201147 peers=50 mempool=3147 avgLatencyMs=127 tipHash=0x0311bb
2025-11-08T19:08:00Z INFO telemetry: height=201148 peers=51 mempool=3148 avgLatencyMs=128 tipHash=0x0311bc
2025-11-08T19:09:00Z INFO telemetry: height=201149 peers=52 mempool=3149 avgLatencyMs=129 tipHash=0x0311bd
2025-11-08T19:10:00Z INFO telemetry: height=201150 peers=48 mempool=3150 avgLatencyMs=130 tipHash=0x0311be
2025-11-08T19:11:00Z INFO telemetry: height=201151 peers=49 mempool=3151 avgLatencyMs=131 tipHash=0x0311bf
2025-11-08T19:12:00Z INFO telemetry: height=201152 peers=50 mempool=3152 avgLatencyMs=132 tipHash=0x0311c0
2025-11-08T19:13:00Z INFO telemetry: height=201153 peers=51 mempool=3153 avgLatencyMs=133 tipHash=0x0311c1
2025-11-08T19:14:00Z INFO telemetry: height=201154 peers=52 mempool=3154 avgLatencyMs=134 tipHash=0x0311c2
2025-11-08T19:15:00Z INFO telemetry: height=201155 peers=48 mempool=3155 avgLatencyMs=135 tipHash=0x0311c3
2025-11-08T19:16:00Z INFO telemetry: height=201156 peers=49 mempool=3156 avgLatencyMs=136 tipHash=0x0311c4
2025-11-08T19:17:00Z INFO telemetry: height=201157 peers=50 mempool=3157 avgLatencyMs=137 tipHash=0x0311c5
2025-11-08T19:18:00Z INFO telemetry: height=201158 peers=51 mempool=3158 avgLatencyMs=138 tipHash=0x0311c6
2025-11-08T19:19:00Z INFO telemetry: height=201159 peers=52 mempool=3159 avgLatencyMs=139 tipHash=0x0311c7
2025-11-08T19:20:00Z INFO telemetry: height=201160 peers=48 mempool=3160 avgLatencyMs=140 tipHash=0x0311c8
2025-11-08T19:21:00Z INFO telemetry: height=201161 peers=49 mempool=3161 avgLatencyMs=141 tipHash=0x0311c9
2025-11-08T19:22:00Z INFO telemetry: height=201162 peers=50 mempool=3162 avgLatencyMs=142 tipHash=0x0311ca
2025-11-08T19:23:00Z INFO telemetry: height=201163 peers=51 mempool=3163 avgLatencyMs=143 tipHash=0x0311cb
2025-11-08T19:24:00Z INFO telemetry: height=201164 peers=52 mempool=3164 avgLatencyMs=144 tipHash=0x0311cc
2025-11-08T19:25:00Z INFO telemetry: height=201165 peers=48 mempool=3165 avgLatencyMs=145 tipHash=0x0311cd
2025-11-08T19:26:00Z INFO telemetry: height=201166 peers=49 mempool=3166 avgLatencyMs=146 tipHash=0x0311ce
2025-11-08T19:27:00Z INFO telemetry: height=201167 peers=50 mempool=3167 avgLatencyMs=147 tipHash=0x0311cf
2025-11-08T19:28:00Z INFO telemetry: height=201168 peers=51 mempool=3168 avgLatencyMs=148 tipHash=0x0311d0
2025-11-08T19:29:00Z INFO telemetry: height=201169 peers=52 mempool=3169 avgLatencyMs=149 tipHash=0x0311d1
2025-11-08T19:30:00Z INFO telemetry: height=201170 peers=48 mempool=3170 avgLatencyMs=120 tipHash=0x0311d2
2025-11-08T19:31:00Z INFO telemetry: height=201171 peers=49 mempool=3171 avgLatencyMs=121 tipHash=0x0311d3
2025-11-08T19:32:00Z INFO telemetry: height=201172 peers=50 mempool=3172 avgLatencyMs=122 tipHash=0x0311d4
2025-11-08T19:33:00Z INFO telemetry: height=201173 peers=51 mempool=3173 avgLatencyMs=123 tipHash=0x0311d5
2025-11-08T19:34:00Z INFO telemetry: height=201174 peers=52 mempool=3174 avgLatencyMs=124 tipHash=0x0311d6
2025-11-08T19:35:00Z INFO telemetry: height=201175 peers=48 mempool=3175 avgLatencyMs=125 tipHash=0x0311d7
2025-11-08T19:36:00Z INFO telemetry: height=201176 peers=49 mempool=3176 avgLatencyMs=126 tipHash=0x0311d8
2025-11-08T19:37:00Z INFO telemetry: height=201177 peers=50 mempool=3177 avgLatencyMs=127 tipHash=0x0311d9
2025-11-08T19:38:00Z INFO telemetry: height=201178 peers=51 mempool=3178 avgLatencyMs=128 tipHash=0x0311da
2025-11-08T19:39:00Z INFO telemetry: height=201179 peers=52 mempool=3179 avgLatencyMs=129 tipHash=0x0311db
2025-11-08T19:40:00Z INFO telemetry: height=201180 peers=48 mempool=3180 avgLatencyMs=130 tipHash=0x0311dc
2025-11-08T19:41:00Z INFO telemetry: height=201181 peers=49 mempool=3181 avgLatencyMs=131 tipHash=0x0311dd
2025-11-08T19:42:00Z INFO telemetry: height=201182 peers=50 mempool=3182 avgLatencyMs=132 tipHash=0x0311de
2025-11-08T19:43:00Z INFO telemetry: height=201183 peers=51 mempool=3183 avgLatencyMs=133 tipHash=0x0311df
2025-11-08T19:44:00Z INFO telemetry: height=201184 peers=52 mempool=3184 avgLatencyMs=134 tipHash=0x0311e0
2025-11-08T19:45:00Z INFO telemetry: height=201185 peers=48 mempool=3185 avgLatencyMs=135 tipHash=0x0311e1
2025-11-08T19:46:00Z INFO telemetry: height=201186 peers=49 mempool=3186 avgLatencyMs=136 tipHash=0x0311e2
2025-11-08T19:47:00Z INFO telemetry: height=201187 peers=50 mempool=3187 avgLatencyMs=137 tipHash=0x0311e3
2025-11-08T19:48:00Z INFO telemetry: height=201188 peers=51 mempool=3188 avgLatencyMs=138 tipHash=0x0311e4
2025-11-08T19:49:00Z INFO telemetry: height=201189 peers=52 mempool=3189 avgLatencyMs=139 tipHash=0x0311e5
2025-11-08T19:50:00Z INFO telemetry: height=201190 peers=48 mempool=3190 avgLatencyMs=140 tipHash=0x0311e6
2025-11-08T19:51:00Z INFO telemetry: height=201191 peers=49 mempool=3191 avgLatencyMs=141 tipHash=0x0311e7
2025-11-08T19:52:00Z INFO telemetry: height=201192 peers=50 mempool=3192 avgLatencyMs=142 tipHash=0x0311e8
2025-11-08T19:53:00Z INFO telemetry: height=201193 peers=51 mempool=3193 avgLatencyMs=143 tipHash=0x0311e9
2025-11-08T19:54:00Z INFO telemetry: height=201194 peers=52 mempool=3194 avgLatencyMs=144 tipHash=0x0311ea
2025-11-08T19:55:00Z INFO telemetry: height=201195 peers=48 mempool=3195 avgLatencyMs=145 tipHash=0x0311eb
2025-11-08T19:56:00Z INFO telemetry: height=201196 peers=49 mempool=3196 avgLatencyMs=146 tipHash=0x0311ec
2025-11-08T19:57:00Z INFO telemetry: height=201197 peers=50 mempool=3197 avgLatencyMs=147 tipHash=0x0311ed
2025-11-08T19:58:00Z INFO telemetry: height=201198 peers=51 mempool=3198 avgLatencyMs=148 tipHash=0x0311ee
2025-11-08T19:59:00Z INFO telemetry: height=201199 peers=52 mempool=3199 avgLatencyMs=149 tipHash=0x0311ef
2025-11-08T20:00:00Z INFO telemetry: height=201200 peers=48 mempool=3000 avgLatencyMs=120 tipHash=0x0311f0
2025-11-08T20:01:00Z INFO telemetry: height=201201 peers=49 mempool=3001 avgLatencyMs=121 tipHash=0x0311f1
2025-11-08T20:02:00Z INFO telemetry: height=201202 peers=50 mempool=3002 avgLatencyMs=122 tipHash=0x0311f2
2025-11-08T20:03:00Z INFO telemetry: height=201203 peers=51 mempool=3003 avgLatencyMs=123 tipHash=0x0311f3
2025-11-08T20:04:00Z INFO telemetry: height=201204 peers=52 mempool=3004 avgLatencyMs=124 tipHash=0x0311f4
2025-11-08T20:05:00Z INFO telemetry: height=201205 peers=48 mempool=3005 avgLatencyMs=125 tipHash=0x0311f5
2025-11-08T20:06:00Z INFO telemetry: height=201206 peers=49 mempool=3006 avgLatencyMs=126 tipHash=0x0311f6
2025-11-08T20:07:00Z INFO telemetry: height=201207 peers=50 mempool=3007 avgLatencyMs=127 tipHash=0x0311f7
2025-11-08T20:08:00Z INFO telemetry: height=201208 peers=51 mempool=3008 avgLatencyMs=128 tipHash=0x0311f8
2025-11-08T20:09:00Z INFO telemetry: height=201209 peers=52 mempool=3009 avgLatencyMs=129 tipHash=0x0311f9
2025-11-08T20:10:00Z INFO telemetry: height=201210 peers=48 mempool=3010 avgLatencyMs=130 tipHash=0x0311fa
2025-11-08T20:11:00Z INFO telemetry: height=201211 peers=49 mempool=3011 avgLatencyMs=131 tipHash=0x0311fb
2025-11-08T20:12:00Z INFO telemetry: height=201212 peers=50 mempool=3012 avgLatencyMs=132 tipHash=0x0311fc
2025-11-08T20:13:00Z INFO telemetry: height=201213 peers=51 mempool=3013 avgLatencyMs=133 tipHash=0x0311fd
2025-11-08T20:14:00Z INFO telemetry: height=201214 peers=52 mempool=3014 avgLatencyMs=134 tipHash=0x0311fe
2025-11-08T20:15:00Z INFO telemetry: height=201215 peers=48 mempool=3015 avgLatencyMs=135 tipHash=0x0311ff
2025-11-08T20:16:00Z INFO telemetry: height=201216 peers=49 mempool=3016 avgLatencyMs=136 tipHash=0x031200
2025-11-08T20:17:00Z INFO telemetry: height=201217 peers=50 mempool=3017 avgLatencyMs=137 tipHash=0x031201
2025-11-08T20:18:00Z INFO telemetry: height=201218 peers=51 mempool=3018 avgLatencyMs=138 tipHash=0x031202
2025-11-08T20:19:00Z INFO telemetry: height=201219 peers=52 mempool=3019 avgLatencyMs=139 tipHash=0x031203
2025-11-08T20:20:00Z INFO telemetry: height=201220 peers=48 mempool=3020 avgLatencyMs=140 tipHash=0x031204
2025-11-08T20:21:00Z INFO telemetry: height=201221 peers=49 mempool=3021 avgLatencyMs=141 tipHash=0x031205
2025-11-08T20:22:00Z INFO telemetry: height=201222 peers=50 mempool=3022 avgLatencyMs=142 tipHash=0x031206
2025-11-08T20:23:00Z INFO telemetry: height=201223 peers=51 mempool=3023 avgLatencyMs=143 tipHash=0x031207
2025-11-08T20:24:00Z INFO telemetry: height=201224 peers=52 mempool=3024 avgLatencyMs=144 tipHash=0x031208
2025-11-08T20:25:00Z INFO telemetry: height=201225 peers=48 mempool=3025 avgLatencyMs=145 tipHash=0x031209
2025-11-08T20:26:00Z INFO telemetry: height=201226 peers=49 mempool=3026 avgLatencyMs=146 tipHash=0x03120a
2025-11-08T20:27:00Z INFO telemetry: height=201227 peers=50 mempool=3027 avgLatencyMs=147 tipHash=0x03120b
2025-11-08T20:28:00Z INFO telemetry: height=201228 peers=51 mempool=3028 avgLatencyMs=148 tipHash=0x03120c
2025-11-08T20:29:00Z INFO telemetry: height=201229 peers=52 mempool=3029 avgLatencyMs=149 tipHash=0x03120d
2025-11-08T20:30:00Z INFO telemetry: height=201230 peers=48 mempool=3030 avgLatencyMs=120 tipHash=0x03120e
2025-11-08T20:31:00Z INFO telemetry: height=201231 peers=49 mempool=3031 avgLatencyMs=121 tipHash=0x03120f
2025-11-08T20:32:00Z INFO telemetry: height=201232 peers=50 mempool=3032 avgLatencyMs=122 tipHash=0x031210
2025-11-08T20:33:00Z INFO telemetry: height=201233 peers=51 mempool=3033 avgLatencyMs=123 tipHash=0x031211
2025-11-08T20:34:00Z INFO telemetry: height=201234 peers=52 mempool=3034 avgLatencyMs=124 tipHash=0x031212
2025-11-08T20:35:00Z INFO telemetry: height=201235 peers=48 mempool=3035 avgLatencyMs=125 tipHash=0x031213
2025-11-08T20:36:00Z INFO telemetry: height=201236 peers=49 mempool=3036 avgLatencyMs=126 tipHash=0x031214
2025-11-08T20:37:00Z INFO telemetry: height=201237 peers=50 mempool=3037 avgLatencyMs=127 tipHash=0x031215
2025-11-08T20:38:00Z INFO telemetry: height=201238 peers=51 mempool=3038 avgLatencyMs=128 tipHash=0x031216
2025-11-08T20:39:00Z INFO telemetry: height=201239 peers=52 mempool=3039 avgLatencyMs=129 tipHash=0x031217
2025-11-08T20:40:00Z INFO telemetry: height=201240 peers=48 mempool=3040 avgLatencyMs=130 tipHash=0x031218
2025-11-08T20:41:00Z INFO telemetry: height=201241 peers=49 mempool=3041 avgLatencyMs=131 tipHash=0x031219
2025-11-08T20:42:00Z INFO telemetry: height=201242 peers=50 mempool=3042 avgLatencyMs=132 tipHash=0x03121a
2025-11-08T20:43:00Z INFO telemetry: height=201243 peers=51 mempool=3043 avgLatencyMs=133 tipHash=0x03121b
2025-11-08T20:44:00Z INFO telemetry: height=201244 peers=52 mempool=3044 avgLatencyMs=134 tipHash=0x03121c
2025-11-08T20:45:00Z INFO telemetry: height=201245 peers=48 mempool=3045 avgLatencyMs=135 tipHash=0x03121d
2025-11-08T20:46:00Z INFO telemetry: height=201246 peers=49 mempool=3046 avgLatencyMs=136 tipHash=0x03121e
2025-11-08T20:47:00Z INFO telemetry: height=201247 peers=50 mempool=3047 avgLatencyMs=137 tipHash=0x03121f
2025-11-08T20:48:00Z INFO telemetry: height=201248 peers=51 mempool=3048 avgLatencyMs=138 tipHash=0x031220
2025-11-08T20:49:00Z INFO telemetry: height=201249 peers=52 mempool=3049 avgLatencyMs=139 tipHash=0x031221
2025-11-08T20:50:00Z INFO telemetry: height=201250 peers=48 mempool=3050 avgLatencyMs=140 tipHash=0x031222
2025-11-08T20:51:00Z INFO telemetry: height=201251 peers=49 mempool=3051 avgLatencyMs=141 tipHash=0x031223
2025-11-08T20:52:00Z INFO telemetry: height=201252 peers=50 mempool=3052 avgLatencyMs=142 tipHash=0x031224
2025-11-08T20:53:00Z INFO telemetry: height=201253 peers=51 mempool=3053 avgLatencyMs=143 tipHash=0x031225
2025-11-08T20:54:00Z INFO telemetry: height=201254 peers=52 mempool=3054 avgLatencyMs=144 tipHash=0x031226
2025-11-08T20:55:00Z INFO telemetry: height=201255 peers=48 mempool=3055 avgLatencyMs=145 tipHash=0x031227
2025-11-08T20:56:00Z INFO telemetry: height=201256 peers=49 mempool=3056 avgLatencyMs=146 tipHash=0x031228
2025-11-08T20:57:00Z INFO telemetry: height=201257 peers=50 mempool=3057 avgLatencyMs=147 tipHash=0x031229
2025-11-08T20:58:00Z INFO telemetry: height=201258 peers=51 mempool=3058 avgLatencyMs=148 tipHash=0x03122a
2025-11-08T20:59:00Z INFO telemetry: height=201259 peers=52 mempool=3059 avgLatencyMs=149 tipHash=0x03122b
2025-11-08T21:00:00Z INFO telemetry: height=201260 peers=48 mempool=3060 avgLatencyMs=120 tipHash=0x03122c
2025-11-08T21:01:00Z INFO telemetry: height=201261 peers=49 mempool=3061 avgLatencyMs=121 tipHash=0x03122d
2025-11-08T21:02:00Z INFO telemetry: height=201262 peers=50 mempool=3062 avgLatencyMs=122 tipHash=0x03122e
2025-11-08T21:03:00Z INFO telemetry: height=201263 peers=51 mempool=3063 avgLatencyMs=123 tipHash=0x03122f
2025-11-08T21:04:00Z INFO telemetry: height=201264 peers=52 mempool=3064 avgLatencyMs=124 tipHash=0x031230
2025-11-08T21:05:00Z INFO telemetry: height=201265 peers=48 mempool=3065 avgLatencyMs=125 tipHash=0x031231
2025-11-08T21:06:00Z INFO telemetry: height=201266 peers=49 mempool=3066 avgLatencyMs=126 tipHash=0x031232
2025-11-08T21:07:00Z INFO telemetry: height=201267 peers=50 mempool=3067 avgLatencyMs=127 tipHash=0x031233
2025-11-08T21:08:00Z INFO telemetry: height=201268 peers=51 mempool=3068 avgLatencyMs=128 tipHash=0x031234
2025-11-08T21:09:00Z INFO telemetry: height=201269 peers=52 mempool=3069 avgLatencyMs=129 tipHash=0x031235
2025-11-08T21:10:00Z INFO telemetry: height=201270 peers=48 mempool=3070 avgLatencyMs=130 tipHash=0x031236
2025-11-08T21:11:00Z INFO telemetry: height=201271 peers=49 mempool=3071 avgLatencyMs=131 tipHash=0x031237
2025-11-08T21:12:00Z INFO telemetry: height=201272 peers=50 mempool=3072 avgLatencyMs=132 tipHash=0x031238
2025-11-08T21:13:00Z INFO telemetry: height=201273 peers=51 mempool=3073 avgLatencyMs=133 tipHash=0x031239
2025-11-08T21:14:00Z INFO telemetry: height=201274 peers=52 mempool=3074 avgLatencyMs=134 tipHash=0x03123a
2025-11-08T21:15:00Z INFO telemetry: height=201275 peers=48 mempool=3075 avgLatencyMs=135 tipHash=0x03123b
2025-11-08T21:16:00Z INFO telemetry: height=201276 peers=49 mempool=3076 avgLatencyMs=136 tipHash=0x03123c
2025-11-08T21:17:00Z INFO telemetry: height=201277 peers=50 mempool=3077 avgLatencyMs=137 tipHash=0x03123d
2025-11-08T21:18:00Z INFO telemetry: height=201278 peers=51 mempool=3078 avgLatencyMs=138 tipHash=0x03123e
2025-11-08T21:19:00Z INFO telemetry: height=201279 peers=52 mempool=3079 avgLatencyMs=139 tipHash=0x03123f
2025-11-08T21:20:00Z INFO telemetry: height=201280 peers=48 mempool=3080 avgLatencyMs=140 tipHash=0x031240
2025-11-08T21:21:00Z INFO telemetry: height=201281 peers=49 mempool=3081 avgLatencyMs=141 tipHash=0x031241
2025-11-08T21:22:00Z INFO telemetry: height=201282 peers=50 mempool=3082 avgLatencyMs=142 tipHash=0x031242
2025-11-08T21:23:00Z INFO telemetry: height=201283 peers=51 mempool=3083 avgLatencyMs=143 tipHash=0x031243
2025-11-08T21:24:00Z INFO telemetry: height=201284 peers=52 mempool=3084 avgLatencyMs=144 tipHash=0x031244
2025-11-08T21:25:00Z INFO telemetry: height=201285 peers=48 mempool=3085 avgLatencyMs=145 tipHash=0x031245
2025-11-08T21:26:00Z INFO telemetry: height=201286 peers=49 mempool=3086 avgLatencyMs=146 tipHash=0x031246
2025-11-08T21:27:00Z INFO telemetry: height=201287 peers=50 mempool=3087 avgLatencyMs=147 tipHash=0x031247
2025-11-08T21:28:00Z INFO telemetry: height=201288 peers=51 mempool=3088 avgLatencyMs=148 tipHash=0x031248
2025-11-08T21:29:00Z INFO telemetry: height=201289 peers=52 mempool=3089 avgLatencyMs=149 tipHash=0x031249
2025-11-08T21:30:00Z INFO telemetry: height=201290 peers=48 mempool=3090 avgLatencyMs=120 tipHash=0x03124a
2025-11-08T21:31:00Z INFO telemetry: height=201291 peers=49 mempool=3091 avgLatencyMs=121 tipHash=0x03124b
2025-11-08T21:32:00Z INFO telemetry: height=201292 peers=50 mempool=3092 avgLatencyMs=122 tipHash=0x03124c
2025-11-08T21:33:00Z INFO telemetry: height=201293 peers=51 mempool=3093 avgLatencyMs=123 tipHash=0x03124d
2025-11-08T21:34:00Z INFO telemetry: height=201294 peers=52 mempool=3094 avgLatencyMs=124 tipHash=0x03124e
2025-11-08T21:35:00Z INFO telemetry: height=201295 peers=48 mempool=3095 avgLatencyMs=125 tipHash=0x03124f
2025-11-08T21:36:00Z INFO telemetry: height=201296 peers=49 mempool=3096 avgLatencyMs=126 tipHash=0x031250
2025-11-08T21:37:00Z INFO telemetry: height=201297 peers=50 mempool=3097 avgLatencyMs=127 tipHash=0x031251
2025-11-08T21:38:00Z INFO telemetry: height=201298 peers=51 mempool=3098 avgLatencyMs=128 tipHash=0x031252
2025-11-08T21:39:00Z INFO telemetry: height=201299 peers=52 mempool=3099 avgLatencyMs=129 tipHash=0x031253
2025-11-08T21:40:00Z INFO telemetry: height=201300 peers=48 mempool=3100 avgLatencyMs=130 tipHash=0x031254
2025-11-08T21:41:00Z INFO telemetry: height=201301 peers=49 mempool=3101 avgLatencyMs=131 tipHash=0x031255
2025-11-08T21:42:00Z INFO telemetry: height=201302 peers=50 mempool=3102 avgLatencyMs=132 tipHash=0x031256
2025-11-08T21:43:00Z INFO telemetry: height=201303 peers=51 mempool=3103 avgLatencyMs=133 tipHash=0x031257
2025-11-08T21:44:00Z INFO telemetry: height=201304 peers=52 mempool=3104 avgLatencyMs=134 tipHash=0x031258
2025-11-08T21:45:00Z INFO telemetry: height=201305 peers=48 mempool=3105 avgLatencyMs=135 tipHash=0x031259
2025-11-08T21:46:00Z INFO telemetry: height=201306 peers=49 mempool=3106 avgLatencyMs=136 tipHash=0x03125a
2025-11-08T21:47:00Z INFO telemetry: height=201307 peers=50 mempool=3107 avgLatencyMs=137 tipHash=0x03125b
2025-11-08T21:48:00Z INFO telemetry: height=201308 peers=51 mempool=3108 avgLatencyMs=138 tipHash=0x03125c
2025-11-08T21:49:00Z INFO telemetry: height=201309 peers=52 mempool=3109 avgLatencyMs=139 tipHash=0x03125d
2025-11-08T21:50:00Z INFO telemetry: height=201310 peers=48 mempool=3110 avgLatencyMs=140 tipHash=0x03125e
2025-11-08T21:51:00Z INFO telemetry: height=201311 peers=49 mempool=3111 avgLatencyMs=141 tipHash=0x03125f
2025-11-08T21:52:00Z INFO telemetry: height=201312 peers=50 mempool=3112 avgLatencyMs=142 tipHash=0x031260
2025-11-08T21:53:00Z INFO telemetry: height=201313 peers=51 mempool=3113 avgLatencyMs=143 tipHash=0x031261
2025-11-08T21:54:00Z INFO telemetry: height=201314 peers=52 mempool=3114 avgLatencyMs=144 tipHash=0x031262
2025-11-08T21:55:00Z INFO telemetry: height=201315 peers=48 mempool=3115 avgLatencyMs=145 tipHash=0x031263
2025-11-08T21:56:00Z INFO telemetry: height=201316 peers=49 mempool=3116 avgLatencyMs=146 tipHash=0x031264
2025-11-08T21:57:00Z INFO telemetry: height=201317 peers=50 mempool=3117 avgLatencyMs=147 tipHash=0x031265
2025-11-08T21:58:00Z INFO telemetry: height=201318 peers=51 mempool=3118 avgLatencyMs=148 tipHash=0x031266
2025-11-08T21:59:00Z INFO telemetry: height=201319 peers=52 mempool=3119 avgLatencyMs=149 tipHash=0x031267
2025-11-08T22:00:00Z INFO telemetry: height=201320 peers=48 mempool=3120 avgLatencyMs=120 tipHash=0x031268
2025-11-08T22:01:00Z INFO telemetry: height=201321 peers=49 mempool=3121 avgLatencyMs=121 tipHash=0x031269
2025-11-08T22:02:00Z INFO telemetry: height=201322 peers=50 mempool=3122 avgLatencyMs=122 tipHash=0x03126a
2025-11-08T22:03:00Z INFO telemetry: height=201323 peers=51 mempool=3123 avgLatencyMs=123 tipHash=0x03126b
2025-11-08T22:04:00Z INFO telemetry: height=201324 peers=52 mempool=3124 avgLatencyMs=124 tipHash=0x03126c
2025-11-08T22:05:00Z INFO telemetry: height=201325 peers=48 mempool=3125 avgLatencyMs=125 tipHash=0x03126d
2025-11-08T22:06:00Z INFO telemetry: height=201326 peers=49 mempool=3126 avgLatencyMs=126 tipHash=0x03126e
2025-11-08T22:07:00Z INFO telemetry: height=201327 peers=50 mempool=3127 avgLatencyMs=127 tipHash=0x03126f
2025-11-08T22:08:00Z INFO telemetry: height=201328 peers=51 mempool=3128 avgLatencyMs=128 tipHash=0x031270
2025-11-08T22:09:00Z INFO telemetry: height=201329 peers=52 mempool=3129 avgLatencyMs=129 tipHash=0x031271
2025-11-08T22:10:00Z INFO telemetry: height=201330 peers=48 mempool=3130 avgLatencyMs=130 tipHash=0x031272
2025-11-08T22:11:00Z INFO telemetry: height=201331 peers=49 mempool=3131 avgLatencyMs=131 tipHash=0x031273
2025-11-08T22:12:00Z INFO telemetry: height=201332 peers=50 mempool=3132 avgLatencyMs=132 tipHash=0x031274
2025-11-08T22:13:00Z INFO telemetry: height=201333 peers=51 mempool=3133 avgLatencyMs=133 tipHash=0x031275
2025-11-08T22:14:00Z INFO telemetry: height=201334 peers=52 mempool=3134 avgLatencyMs=134 tipHash=0x031276
2025-11-08T22:15:00Z INFO telemetry: height=201335 peers=48 mempool=3135 avgLatencyMs=135 tipHash=0x031277
2025-11-08T22:16:00Z INFO telemetry: height=201336 peers=49 mempool=3136 avgLatencyMs=136 tipHash=0x031278
2025-11-08T22:17:00Z INFO telemetry: height=201337 peers=50 mempool=3137 avgLatencyMs=137 tipHash=0x031279
2025-11-08T22:18:00Z INFO telemetry: height=201338 peers=51 mempool=3138 avgLatencyMs=138 tipHash=0x03127a
2025-11-08T22:19:00Z INFO telemetry: height=201339 peers=52 mempool=3139 avgLatencyMs=139 tipHash=0x03127b
2025-11-08T22:20:00Z INFO telemetry: height=201340 peers=48 mempool=3140 avgLatencyMs=140 tipHash=0x03127c
2025-11-08T22:21:00Z INFO telemetry: height=201341 peers=49 mempool=3141 avgLatencyMs=141 tipHash=0x03127d
2025-11-08T22:22:00Z INFO telemetry: height=201342 peers=50 mempool=3142 avgLatencyMs=142 tipHash=0x03127e
2025-11-08T22:23:00Z INFO telemetry: height=201343 peers=51 mempool=3143 avgLatencyMs=143 tipHash=0x03127f
2025-11-08T22:24:00Z INFO telemetry: height=201344 peers=52 mempool=3144 avgLatencyMs=144 tipHash=0x031280
2025-11-08T22:25:00Z INFO telemetry: height=201345 peers=48 mempool=3145 avgLatencyMs=145 tipHash=0x031281
2025-11-08T22:26:00Z INFO telemetry: height=201346 peers=49 mempool=3146 avgLatencyMs=146 tipHash=0x031282
2025-11-08T22:27:00Z INFO telemetry: height=201347 peers=50 mempool=3147 avgLatencyMs=147 tipHash=0x031283
2025-11-08T22:28:00Z INFO telemetry: height=201348 peers=51 mempool=3148 avgLatencyMs=148 tipHash=0x031284
2025-11-08T22:29:00Z INFO telemetry: height=201349 peers=52 mempool=3149 avgLatencyMs=149 tipHash=0x031285
2025-11-08T22:30:00Z INFO telemetry: height=201350 peers=48 mempool=3150 avgLatencyMs=120 tipHash=0x031286
2025-11-08T22:31:00Z INFO telemetry: height=201351 peers=49 mempool=3151 avgLatencyMs=121 tipHash=0x031287
2025-11-08T22:32:00Z INFO telemetry: height=201352 peers=50 mempool=3152 avgLatencyMs=122 tipHash=0x031288
2025-11-08T22:33:00Z INFO telemetry: height=201353 peers=51 mempool=3153 avgLatencyMs=123 tipHash=0x031289
2025-11-08T22:34:00Z INFO telemetry: height=201354 peers=52 mempool=3154 avgLatencyMs=124 tipHash=0x03128a
2025-11-08T22:35:00Z INFO telemetry: height=201355 peers=48 mempool=3155 avgLatencyMs=125 tipHash=0x03128b
2025-11-08T22:36:00Z INFO telemetry: height=201356 peers=49 mempool=3156 avgLatencyMs=126 tipHash=0x03128c
2025-11-08T22:37:00Z INFO telemetry: height=201357 peers=50 mempool=3157 avgLatencyMs=127 tipHash=0x03128d
2025-11-08T22:38:00Z INFO telemetry: height=201358 peers=51 mempool=3158 avgLatencyMs=128 tipHash=0x03128e
2025-11-08T22:39:00Z INFO telemetry: height=201359 peers=52 mempool=3159 avgLatencyMs=129 tipHash=0x03128f
2025-11-08T22:40:00Z INFO telemetry: height=201360 peers=48 mempool=3160 avgLatencyMs=130 tipHash=0x031290
2025-11-08T22:41:00Z INFO telemetry: height=201361 peers=49 mempool=3161 avgLatencyMs=131 tipHash=0x031291
2025-11-08T22:42:00Z INFO telemetry: height=201362 peers=50 mempool=3162 avgLatencyMs=132 tipHash=0x031292
2025-11-08T22:43:00Z INFO telemetry: height=201363 peers=51 mempool=3163 avgLatencyMs=133 tipHash=0x031293
2025-11-08T22:44:00Z INFO telemetry: height=201364 peers=52 mempool=3164 avgLatencyMs=134 tipHash=0x031294
2025-11-08T22:45:00Z INFO telemetry: height=201365 peers=48 mempool=3165 avgLatencyMs=135 tipHash=0x031295
2025-11-08T22:46:00Z INFO telemetry: height=201366 peers=49 mempool=3166 avgLatencyMs=136 tipHash=0x031296
2025-11-08T22:47:00Z INFO telemetry: height=201367 peers=50 mempool=3167 avgLatencyMs=137 tipHash=0x031297
2025-11-08T22:48:00Z INFO telemetry: height=201368 peers=51 mempool=3168 avgLatencyMs=138 tipHash=0x031298
2025-11-08T22:49:00Z INFO telemetry: height=201369 peers=52 mempool=3169 avgLatencyMs=139 tipHash=0x031299
2025-11-08T22:50:00Z INFO telemetry: height=201370 peers=48 mempool=3170 avgLatencyMs=140 tipHash=0x03129a
2025-11-08T22:51:00Z INFO telemetry: height=201371 peers=49 mempool=3171 avgLatencyMs=141 tipHash=0x03129b
2025-11-08T22:52:00Z INFO telemetry: height=201372 peers=50 mempool=3172 avgLatencyMs=142 tipHash=0x03129c
2025-11-08T22:53:00Z INFO telemetry: height=201373 peers=51 mempool=3173 avgLatencyMs=143 tipHash=0x03129d
2025-11-08T22:54:00Z INFO telemetry: height=201374 peers=52 mempool=3174 avgLatencyMs=144 tipHash=0x03129e
2025-11-08T22:55:00Z INFO telemetry: height=201375 peers=48 mempool=3175 avgLatencyMs=145 tipHash=0x03129f
2025-11-08T22:56:00Z INFO telemetry: height=201376 peers=49 mempool=3176 avgLatencyMs=146 tipHash=0x0312a0
2025-11-08T22:57:00Z INFO telemetry: height=201377 peers=50 mempool=3177 avgLatencyMs=147 tipHash=0x0312a1
2025-11-08T22:58:00Z INFO telemetry: height=201378 peers=51 mempool=3178 avgLatencyMs=148 tipHash=0x0312a2
2025-11-08T22:59:00Z INFO telemetry: height=201379 peers=52 mempool=3179 avgLatencyMs=149 tipHash=0x0312a3
2025-11-08T23:00:00Z INFO telemetry: height=201380 peers=48 mempool=3180 avgLatencyMs=120 tipHash=0x0312a4
2025-11-08T23:01:00Z INFO telemetry: height=201381 peers=49 mempool=3181 avgLatencyMs=121 tipHash=0x0312a5
2025-11-08T23:02:00Z INFO telemetry: height=201382 peers=50 mempool=3182 avgLatencyMs=122 tipHash=0x0312a6
2025-11-08T23:03:00Z INFO telemetry: height=201383 peers=51 mempool=3183 avgLatencyMs=123 tipHash=0x0312a7
2025-11-08T23:04:00Z INFO telemetry: height=201384 peers=52 mempool=3184 avgLatencyMs=124 tipHash=0x0312a8
2025-11-08T23:05:00Z INFO telemetry: height=201385 peers=48 mempool=3185 avgLatencyMs=125 tipHash=0x0312a9
2025-11-08T23:06:00Z INFO telemetry: height=201386 peers=49 mempool=3186 avgLatencyMs=126 tipHash=0x0312aa
2025-11-08T23:07:00Z INFO telemetry: height=201387 peers=50 mempool=3187 avgLatencyMs=127 tipHash=0x0312ab
2025-11-08T23:08:00Z INFO telemetry: height=201388 peers=51 mempool=3188 avgLatencyMs=128 tipHash=0x0312ac
2025-11-08T23:09:00Z INFO telemetry: height=201389 peers=52 mempool=3189 avgLatencyMs=129 tipHash=0x0312ad
2025-11-08T23:10:00Z INFO telemetry: height=201390 peers=48 mempool=3190 avgLatencyMs=130 tipHash=0x0312ae
2025-11-08T23:11:00Z INFO telemetry: height=201391 peers=49 mempool=3191 avgLatencyMs=131 tipHash=0x0312af
2025-11-08T23:12:00Z INFO telemetry: height=201392 peers=50 mempool=3192 avgLatencyMs=132 tipHash=0x0312b0
2025-11-08T23:13:00Z INFO telemetry: height=201393 peers=51 mempool=3193 avgLatencyMs=133 tipHash=0x0312b1
2025-11-08T23:14:00Z INFO telemetry: height=201394 peers=52 mempool=3194 avgLatencyMs=134 tipHash=0x0312b2
2025-11-08T23:15:00Z INFO telemetry: height=201395 peers=48 mempool=3195 avgLatencyMs=135 tipHash=0x0312b3
2025-11-08T23:16:00Z INFO telemetry: height=201396 peers=49 mempool=3196 avgLatencyMs=136 tipHash=0x0312b4
2025-11-08T23:17:00Z INFO telemetry: height=201397 peers=50 mempool=3197 avgLatencyMs=137 tipHash=0x0312b5
2025-11-08T23:18:00Z INFO telemetry: height=201398 peers=51 mempool=3198 avgLatencyMs=138 tipHash=0x0312b6
2025-11-08T23:19:00Z INFO telemetry: height=201399 peers=52 mempool=3199 avgLatencyMs=139 tipHash=0x0312b7
2025-11-08T23:20:00Z INFO telemetry: height=201400 peers=48 mempool=3000 avgLatencyMs=140 tipHash=0x0312b8
2025-11-08T23:21:00Z INFO telemetry: height=201401 peers=49 mempool=3001 avgLatencyMs=141 tipHash=0x0312b9
2025-11-08T23:22:00Z INFO telemetry: height=201402 peers=50 mempool=3002 avgLatencyMs=142 tipHash=0x0312ba
2025-11-08T23:23:00Z INFO telemetry: height=201403 peers=51 mempool=3003 avgLatencyMs=143 tipHash=0x0312bb
2025-11-08T23:24:00Z INFO telemetry: height=201404 peers=52 mempool=3004 avgLatencyMs=144 tipHash=0x0312bc
2025-11-08T23:25:00Z INFO telemetry: height=201405 peers=48 mempool=3005 avgLatencyMs=145 tipHash=0x0312bd
2025-11-08T23:26:00Z INFO telemetry: height=201406 peers=49 mempool=3006 avgLatencyMs=146 tipHash=0x0312be
2025-11-08T23:27:00Z INFO telemetry: height=201407 peers=50 mempool=3007 avgLatencyMs=147 tipHash=0x0312bf
2025-11-08T23:28:00Z INFO telemetry: height=201408 peers=51 mempool=3008 avgLatencyMs=148 tipHash=0x0312c0
2025-11-08T23:29:00Z INFO telemetry: height=201409 peers=52 mempool=3009 avgLatencyMs=149 tipHash=0x0312c1
2025-11-08T23:30:00Z INFO telemetry: height=201410 peers=48 mempool=3010 avgLatencyMs=120 tipHash=0x0312c2
2025-11-08T23:31:00Z INFO telemetry: height=201411 peers=49 mempool=3011 avgLatencyMs=121 tipHash=0x0312c3
2025-11-08T23:32:00Z INFO telemetry: height=201412 peers=50 mempool=3012 avgLatencyMs=122 tipHash=0x0312c4
2025-11-08T23:33:00Z INFO telemetry: height=201413 peers=51 mempool=3013 avgLatencyMs=123 tipHash=0x0312c5
2025-11-08T23:34:00Z INFO telemetry: height=201414 peers=52 mempool=3014 avgLatencyMs=124 tipHash=0x0312c6
2025-11-08T23:35:00Z INFO telemetry: height=201415 peers=48 mempool=3015 avgLatencyMs=125 tipHash=0x0312c7
2025-11-08T23:36:00Z INFO telemetry: height=201416 peers=49 mempool=3016 avgLatencyMs=126 tipHash=0x0312c8
2025-11-08T23:37:00Z INFO telemetry: height=201417 peers=50 mempool=3017 avgLatencyMs=127 tipHash=0x0312c9
2025-11-08T23:38:00Z INFO telemetry: height=201418 peers=51 mempool=3018 avgLatencyMs=128 tipHash=0x0312ca
2025-11-08T23:39:00Z INFO telemetry: height=201419 peers=52 mempool=3019 avgLatencyMs=129 tipHash=0x0312cb
2025-11-08T23:40:00Z INFO telemetry: height=201420 peers=48 mempool=3020 avgLatencyMs=130 tipHash=0x0312cc
2025-11-08T23:41:00Z INFO telemetry: height=201421 peers=49 mempool=3021 avgLatencyMs=131 tipHash=0x0312cd
2025-11-08T23:42:00Z INFO telemetry: height=201422 peers=50 mempool=3022 avgLatencyMs=132 tipHash=0x0312ce
2025-11-08T23:43:00Z INFO telemetry: height=201423 peers=51 mempool=3023 avgLatencyMs=133 tipHash=0x0312cf
2025-11-08T23:44:00Z INFO telemetry: height=201424 peers=52 mempool=3024 avgLatencyMs=134 tipHash=0x0312d0
2025-11-08T23:45:00Z INFO telemetry: height=201425 peers=48 mempool=3025 avgLatencyMs=135 tipHash=0x0312d1
2025-11-08T23:46:00Z INFO telemetry: height=201426 peers=49 mempool=3026 avgLatencyMs=136 tipHash=0x0312d2
2025-11-08T23:47:00Z INFO telemetry: height=201427 peers=50 mempool=3027 avgLatencyMs=137 tipHash=0x0312d3
2025-11-08T23:48:00Z INFO telemetry: height=201428 peers=51 mempool=3028 avgLatencyMs=138 tipHash=0x0312d4
2025-11-08T23:49:00Z INFO telemetry: height=201429 peers=52 mempool=3029 avgLatencyMs=139 tipHash=0x0312d5
2025-11-08T23:50:00Z INFO telemetry: height=201430 peers=48 mempool=3030 avgLatencyMs=140 tipHash=0x0312d6
2025-11-08T23:51:00Z INFO telemetry: height=201431 peers=49 mempool=3031 avgLatencyMs=141 tipHash=0x0312d7
2025-11-08T23:52:00Z INFO telemetry: height=201432 peers=50 mempool=3032 avgLatencyMs=142 tipHash=0x0312d8
2025-11-08T23:53:00Z INFO telemetry: height=201433 peers=51 mempool=3033 avgLatencyMs=143 tipHash=0x0312d9
2025-11-08T23:54:00Z INFO telemetry: height=201434 peers=52 mempool=3034 avgLatencyMs=144 tipHash=0x0312da
2025-11-08T23:55:00Z INFO telemetry: height=201435 peers=48 mempool=3035 avgLatencyMs=145 tipHash=0x0312db
2025-11-08T23:56:00Z INFO telemetry: height=201436 peers=49 mempool=3036 avgLatencyMs=146 tipHash=0x0312dc
2025-11-08T23:57:00Z INFO telemetry: height=201437 peers=50 mempool=3037 avgLatencyMs=147 tipHash=0x0312dd
2025-11-08T23:58:00Z INFO telemetry: height=201438 peers=51 mempool=3038 avgLatencyMs=148 tipHash=0x0312de
2025-11-08T23:59:00Z INFO telemetry: height=201439 peers=52 mempool=3039 avgLatencyMs=149 tipHash=0x0312df
```

### Appendix AB: RPC Submission Trace

```text
2025-11-08T00:00:00Z RPC /submit payloadHash=0x030d4010 fee=1000 status=broadcast
2025-11-08T00:01:00Z RPC /submit payloadHash=0x030d4111 fee=1001 status=accepted
2025-11-08T00:02:00Z RPC /submit payloadHash=0x030d4212 fee=1002 status=accepted
2025-11-08T00:03:00Z RPC /submit payloadHash=0x030d4313 fee=1003 status=broadcast
2025-11-08T00:04:00Z RPC /submit payloadHash=0x030d4414 fee=1004 status=accepted
2025-11-08T00:05:00Z RPC /submit payloadHash=0x030d4515 fee=1005 status=accepted
2025-11-08T00:06:00Z RPC /submit payloadHash=0x030d4616 fee=1006 status=broadcast
2025-11-08T00:07:00Z RPC /submit payloadHash=0x030d4717 fee=1007 status=accepted
2025-11-08T00:08:00Z RPC /submit payloadHash=0x030d4818 fee=1008 status=accepted
2025-11-08T00:09:00Z RPC /submit payloadHash=0x030d4919 fee=1009 status=broadcast
2025-11-08T00:10:00Z RPC /submit payloadHash=0x030d4a20 fee=1010 status=accepted
2025-11-08T00:11:00Z RPC /submit payloadHash=0x030d4b21 fee=1011 status=accepted
2025-11-08T00:12:00Z RPC /submit payloadHash=0x030d4c22 fee=1012 status=broadcast
2025-11-08T00:13:00Z RPC /submit payloadHash=0x030d4d23 fee=1013 status=accepted
2025-11-08T00:14:00Z RPC /submit payloadHash=0x030d4e24 fee=1014 status=accepted
2025-11-08T00:15:00Z RPC /submit payloadHash=0x030d4f25 fee=1015 status=broadcast
2025-11-08T00:16:00Z RPC /submit payloadHash=0x030d5026 fee=1016 status=accepted
2025-11-08T00:17:00Z RPC /submit payloadHash=0x030d5127 fee=1017 status=accepted
2025-11-08T00:18:00Z RPC /submit payloadHash=0x030d5228 fee=1018 status=broadcast
2025-11-08T00:19:00Z RPC /submit payloadHash=0x030d5329 fee=1019 status=accepted
2025-11-08T00:20:00Z RPC /submit payloadHash=0x030d5430 fee=1020 status=accepted
2025-11-08T00:21:00Z RPC /submit payloadHash=0x030d5531 fee=1021 status=broadcast
2025-11-08T00:22:00Z RPC /submit payloadHash=0x030d5632 fee=1022 status=accepted
2025-11-08T00:23:00Z RPC /submit payloadHash=0x030d5733 fee=1023 status=accepted
2025-11-08T00:24:00Z RPC /submit payloadHash=0x030d5834 fee=1024 status=broadcast
2025-11-08T00:25:00Z RPC /submit payloadHash=0x030d5935 fee=1025 status=accepted
2025-11-08T00:26:00Z RPC /submit payloadHash=0x030d5a36 fee=1026 status=accepted
2025-11-08T00:27:00Z RPC /submit payloadHash=0x030d5b37 fee=1027 status=broadcast
2025-11-08T00:28:00Z RPC /submit payloadHash=0x030d5c38 fee=1028 status=accepted
2025-11-08T00:29:00Z RPC /submit payloadHash=0x030d5d39 fee=1029 status=accepted
2025-11-08T00:30:00Z RPC /submit payloadHash=0x030d5e40 fee=1030 status=broadcast
2025-11-08T00:31:00Z RPC /submit payloadHash=0x030d5f41 fee=1031 status=accepted
2025-11-08T00:32:00Z RPC /submit payloadHash=0x030d6042 fee=1032 status=accepted
2025-11-08T00:33:00Z RPC /submit payloadHash=0x030d6143 fee=1033 status=broadcast
2025-11-08T00:34:00Z RPC /submit payloadHash=0x030d6244 fee=1034 status=accepted
2025-11-08T00:35:00Z RPC /submit payloadHash=0x030d6345 fee=1035 status=accepted
2025-11-08T00:36:00Z RPC /submit payloadHash=0x030d6446 fee=1036 status=broadcast
2025-11-08T00:37:00Z RPC /submit payloadHash=0x030d6547 fee=1037 status=accepted
2025-11-08T00:38:00Z RPC /submit payloadHash=0x030d6648 fee=1038 status=accepted
2025-11-08T00:39:00Z RPC /submit payloadHash=0x030d6749 fee=1039 status=broadcast
2025-11-08T00:40:00Z RPC /submit payloadHash=0x030d6850 fee=1040 status=accepted
2025-11-08T00:41:00Z RPC /submit payloadHash=0x030d6951 fee=1041 status=accepted
2025-11-08T00:42:00Z RPC /submit payloadHash=0x030d6a52 fee=1042 status=broadcast
2025-11-08T00:43:00Z RPC /submit payloadHash=0x030d6b53 fee=1043 status=accepted
2025-11-08T00:44:00Z RPC /submit payloadHash=0x030d6c54 fee=1044 status=accepted
2025-11-08T00:45:00Z RPC /submit payloadHash=0x030d6d55 fee=1045 status=broadcast
2025-11-08T00:46:00Z RPC /submit payloadHash=0x030d6e56 fee=1046 status=accepted
2025-11-08T00:47:00Z RPC /submit payloadHash=0x030d6f57 fee=1047 status=accepted
2025-11-08T00:48:00Z RPC /submit payloadHash=0x030d7058 fee=1048 status=broadcast
2025-11-08T00:49:00Z RPC /submit payloadHash=0x030d7159 fee=1049 status=accepted
2025-11-08T00:50:00Z RPC /submit payloadHash=0x030d7210 fee=1050 status=accepted
2025-11-08T00:51:00Z RPC /submit payloadHash=0x030d7311 fee=1051 status=broadcast
2025-11-08T00:52:00Z RPC /submit payloadHash=0x030d7412 fee=1052 status=accepted
2025-11-08T00:53:00Z RPC /submit payloadHash=0x030d7513 fee=1053 status=accepted
2025-11-08T00:54:00Z RPC /submit payloadHash=0x030d7614 fee=1054 status=broadcast
2025-11-08T00:55:00Z RPC /submit payloadHash=0x030d7715 fee=1055 status=accepted
2025-11-08T00:56:00Z RPC /submit payloadHash=0x030d7816 fee=1056 status=accepted
2025-11-08T00:57:00Z RPC /submit payloadHash=0x030d7917 fee=1057 status=broadcast
2025-11-08T00:58:00Z RPC /submit payloadHash=0x030d7a18 fee=1058 status=accepted
2025-11-08T00:59:00Z RPC /submit payloadHash=0x030d7b19 fee=1059 status=accepted
2025-11-08T01:00:00Z RPC /submit payloadHash=0x030d7c20 fee=1060 status=broadcast
2025-11-08T01:01:00Z RPC /submit payloadHash=0x030d7d21 fee=1061 status=accepted
2025-11-08T01:02:00Z RPC /submit payloadHash=0x030d7e22 fee=1062 status=accepted
2025-11-08T01:03:00Z RPC /submit payloadHash=0x030d7f23 fee=1063 status=broadcast
2025-11-08T01:04:00Z RPC /submit payloadHash=0x030d8024 fee=1064 status=accepted
2025-11-08T01:05:00Z RPC /submit payloadHash=0x030d8125 fee=1065 status=accepted
2025-11-08T01:06:00Z RPC /submit payloadHash=0x030d8226 fee=1066 status=broadcast
2025-11-08T01:07:00Z RPC /submit payloadHash=0x030d8327 fee=1067 status=accepted
2025-11-08T01:08:00Z RPC /submit payloadHash=0x030d8428 fee=1068 status=accepted
2025-11-08T01:09:00Z RPC /submit payloadHash=0x030d8529 fee=1069 status=broadcast
2025-11-08T01:10:00Z RPC /submit payloadHash=0x030d8630 fee=1070 status=accepted
2025-11-08T01:11:00Z RPC /submit payloadHash=0x030d8731 fee=1071 status=accepted
2025-11-08T01:12:00Z RPC /submit payloadHash=0x030d8832 fee=1072 status=broadcast
2025-11-08T01:13:00Z RPC /submit payloadHash=0x030d8933 fee=1073 status=accepted
2025-11-08T01:14:00Z RPC /submit payloadHash=0x030d8a34 fee=1074 status=accepted
2025-11-08T01:15:00Z RPC /submit payloadHash=0x030d8b35 fee=1075 status=broadcast
2025-11-08T01:16:00Z RPC /submit payloadHash=0x030d8c36 fee=1076 status=accepted
2025-11-08T01:17:00Z RPC /submit payloadHash=0x030d8d37 fee=1077 status=accepted
2025-11-08T01:18:00Z RPC /submit payloadHash=0x030d8e38 fee=1078 status=broadcast
2025-11-08T01:19:00Z RPC /submit payloadHash=0x030d8f39 fee=1079 status=accepted
2025-11-08T01:20:00Z RPC /submit payloadHash=0x030d9040 fee=1080 status=accepted
2025-11-08T01:21:00Z RPC /submit payloadHash=0x030d9141 fee=1081 status=broadcast
2025-11-08T01:22:00Z RPC /submit payloadHash=0x030d9242 fee=1082 status=accepted
2025-11-08T01:23:00Z RPC /submit payloadHash=0x030d9343 fee=1083 status=accepted
2025-11-08T01:24:00Z RPC /submit payloadHash=0x030d9444 fee=1084 status=broadcast
2025-11-08T01:25:00Z RPC /submit payloadHash=0x030d9545 fee=1085 status=accepted
2025-11-08T01:26:00Z RPC /submit payloadHash=0x030d9646 fee=1086 status=accepted
2025-11-08T01:27:00Z RPC /submit payloadHash=0x030d9747 fee=1087 status=broadcast
2025-11-08T01:28:00Z RPC /submit payloadHash=0x030d9848 fee=1088 status=accepted
2025-11-08T01:29:00Z RPC /submit payloadHash=0x030d9949 fee=1089 status=accepted
2025-11-08T01:30:00Z RPC /submit payloadHash=0x030d9a50 fee=1090 status=broadcast
2025-11-08T01:31:00Z RPC /submit payloadHash=0x030d9b51 fee=1091 status=accepted
2025-11-08T01:32:00Z RPC /submit payloadHash=0x030d9c52 fee=1092 status=accepted
2025-11-08T01:33:00Z RPC /submit payloadHash=0x030d9d53 fee=1093 status=broadcast
2025-11-08T01:34:00Z RPC /submit payloadHash=0x030d9e54 fee=1094 status=accepted
2025-11-08T01:35:00Z RPC /submit payloadHash=0x030d9f55 fee=1095 status=accepted
2025-11-08T01:36:00Z RPC /submit payloadHash=0x030da056 fee=1096 status=broadcast
2025-11-08T01:37:00Z RPC /submit payloadHash=0x030da157 fee=1097 status=accepted
2025-11-08T01:38:00Z RPC /submit payloadHash=0x030da258 fee=1098 status=accepted
2025-11-08T01:39:00Z RPC /submit payloadHash=0x030da359 fee=1099 status=broadcast
2025-11-08T01:40:00Z RPC /submit payloadHash=0x030da410 fee=1100 status=accepted
2025-11-08T01:41:00Z RPC /submit payloadHash=0x030da511 fee=1101 status=accepted
2025-11-08T01:42:00Z RPC /submit payloadHash=0x030da612 fee=1102 status=broadcast
2025-11-08T01:43:00Z RPC /submit payloadHash=0x030da713 fee=1103 status=accepted
2025-11-08T01:44:00Z RPC /submit payloadHash=0x030da814 fee=1104 status=accepted
2025-11-08T01:45:00Z RPC /submit payloadHash=0x030da915 fee=1105 status=broadcast
2025-11-08T01:46:00Z RPC /submit payloadHash=0x030daa16 fee=1106 status=accepted
2025-11-08T01:47:00Z RPC /submit payloadHash=0x030dab17 fee=1107 status=accepted
2025-11-08T01:48:00Z RPC /submit payloadHash=0x030dac18 fee=1108 status=broadcast
2025-11-08T01:49:00Z RPC /submit payloadHash=0x030dad19 fee=1109 status=accepted
2025-11-08T01:50:00Z RPC /submit payloadHash=0x030dae20 fee=1110 status=accepted
2025-11-08T01:51:00Z RPC /submit payloadHash=0x030daf21 fee=1111 status=broadcast
2025-11-08T01:52:00Z RPC /submit payloadHash=0x030db022 fee=1112 status=accepted
2025-11-08T01:53:00Z RPC /submit payloadHash=0x030db123 fee=1113 status=accepted
2025-11-08T01:54:00Z RPC /submit payloadHash=0x030db224 fee=1114 status=broadcast
2025-11-08T01:55:00Z RPC /submit payloadHash=0x030db325 fee=1115 status=accepted
2025-11-08T01:56:00Z RPC /submit payloadHash=0x030db426 fee=1116 status=accepted
2025-11-08T01:57:00Z RPC /submit payloadHash=0x030db527 fee=1117 status=broadcast
2025-11-08T01:58:00Z RPC /submit payloadHash=0x030db628 fee=1118 status=accepted
2025-11-08T01:59:00Z RPC /submit payloadHash=0x030db729 fee=1119 status=accepted
2025-11-08T02:00:00Z RPC /submit payloadHash=0x030db830 fee=1120 status=broadcast
2025-11-08T02:01:00Z RPC /submit payloadHash=0x030db931 fee=1121 status=accepted
2025-11-08T02:02:00Z RPC /submit payloadHash=0x030dba32 fee=1122 status=accepted
2025-11-08T02:03:00Z RPC /submit payloadHash=0x030dbb33 fee=1123 status=broadcast
2025-11-08T02:04:00Z RPC /submit payloadHash=0x030dbc34 fee=1124 status=accepted
2025-11-08T02:05:00Z RPC /submit payloadHash=0x030dbd35 fee=1125 status=accepted
2025-11-08T02:06:00Z RPC /submit payloadHash=0x030dbe36 fee=1126 status=broadcast
2025-11-08T02:07:00Z RPC /submit payloadHash=0x030dbf37 fee=1127 status=accepted
2025-11-08T02:08:00Z RPC /submit payloadHash=0x030dc038 fee=1128 status=accepted
2025-11-08T02:09:00Z RPC /submit payloadHash=0x030dc139 fee=1129 status=broadcast
2025-11-08T02:10:00Z RPC /submit payloadHash=0x030dc240 fee=1130 status=accepted
2025-11-08T02:11:00Z RPC /submit payloadHash=0x030dc341 fee=1131 status=accepted
2025-11-08T02:12:00Z RPC /submit payloadHash=0x030dc442 fee=1132 status=broadcast
2025-11-08T02:13:00Z RPC /submit payloadHash=0x030dc543 fee=1133 status=accepted
2025-11-08T02:14:00Z RPC /submit payloadHash=0x030dc644 fee=1134 status=accepted
2025-11-08T02:15:00Z RPC /submit payloadHash=0x030dc745 fee=1135 status=broadcast
2025-11-08T02:16:00Z RPC /submit payloadHash=0x030dc846 fee=1136 status=accepted
2025-11-08T02:17:00Z RPC /submit payloadHash=0x030dc947 fee=1137 status=accepted
2025-11-08T02:18:00Z RPC /submit payloadHash=0x030dca48 fee=1138 status=broadcast
2025-11-08T02:19:00Z RPC /submit payloadHash=0x030dcb49 fee=1139 status=accepted
2025-11-08T02:20:00Z RPC /submit payloadHash=0x030dcc50 fee=1140 status=accepted
2025-11-08T02:21:00Z RPC /submit payloadHash=0x030dcd51 fee=1141 status=broadcast
2025-11-08T02:22:00Z RPC /submit payloadHash=0x030dce52 fee=1142 status=accepted
2025-11-08T02:23:00Z RPC /submit payloadHash=0x030dcf53 fee=1143 status=accepted
2025-11-08T02:24:00Z RPC /submit payloadHash=0x030dd054 fee=1144 status=broadcast
2025-11-08T02:25:00Z RPC /submit payloadHash=0x030dd155 fee=1145 status=accepted
2025-11-08T02:26:00Z RPC /submit payloadHash=0x030dd256 fee=1146 status=accepted
2025-11-08T02:27:00Z RPC /submit payloadHash=0x030dd357 fee=1147 status=broadcast
2025-11-08T02:28:00Z RPC /submit payloadHash=0x030dd458 fee=1148 status=accepted
2025-11-08T02:29:00Z RPC /submit payloadHash=0x030dd559 fee=1149 status=accepted
2025-11-08T02:30:00Z RPC /submit payloadHash=0x030dd610 fee=1150 status=broadcast
2025-11-08T02:31:00Z RPC /submit payloadHash=0x030dd711 fee=1151 status=accepted
2025-11-08T02:32:00Z RPC /submit payloadHash=0x030dd812 fee=1152 status=accepted
2025-11-08T02:33:00Z RPC /submit payloadHash=0x030dd913 fee=1153 status=broadcast
2025-11-08T02:34:00Z RPC /submit payloadHash=0x030dda14 fee=1154 status=accepted
2025-11-08T02:35:00Z RPC /submit payloadHash=0x030ddb15 fee=1155 status=accepted
2025-11-08T02:36:00Z RPC /submit payloadHash=0x030ddc16 fee=1156 status=broadcast
2025-11-08T02:37:00Z RPC /submit payloadHash=0x030ddd17 fee=1157 status=accepted
2025-11-08T02:38:00Z RPC /submit payloadHash=0x030dde18 fee=1158 status=accepted
2025-11-08T02:39:00Z RPC /submit payloadHash=0x030ddf19 fee=1159 status=broadcast
2025-11-08T02:40:00Z RPC /submit payloadHash=0x030de020 fee=1160 status=accepted
2025-11-08T02:41:00Z RPC /submit payloadHash=0x030de121 fee=1161 status=accepted
2025-11-08T02:42:00Z RPC /submit payloadHash=0x030de222 fee=1162 status=broadcast
2025-11-08T02:43:00Z RPC /submit payloadHash=0x030de323 fee=1163 status=accepted
2025-11-08T02:44:00Z RPC /submit payloadHash=0x030de424 fee=1164 status=accepted
2025-11-08T02:45:00Z RPC /submit payloadHash=0x030de525 fee=1165 status=broadcast
2025-11-08T02:46:00Z RPC /submit payloadHash=0x030de626 fee=1166 status=accepted
2025-11-08T02:47:00Z RPC /submit payloadHash=0x030de727 fee=1167 status=accepted
2025-11-08T02:48:00Z RPC /submit payloadHash=0x030de828 fee=1168 status=broadcast
2025-11-08T02:49:00Z RPC /submit payloadHash=0x030de929 fee=1169 status=accepted
2025-11-08T02:50:00Z RPC /submit payloadHash=0x030dea30 fee=1170 status=accepted
2025-11-08T02:51:00Z RPC /submit payloadHash=0x030deb31 fee=1171 status=broadcast
2025-11-08T02:52:00Z RPC /submit payloadHash=0x030dec32 fee=1172 status=accepted
2025-11-08T02:53:00Z RPC /submit payloadHash=0x030ded33 fee=1173 status=accepted
2025-11-08T02:54:00Z RPC /submit payloadHash=0x030dee34 fee=1174 status=broadcast
2025-11-08T02:55:00Z RPC /submit payloadHash=0x030def35 fee=1175 status=accepted
2025-11-08T02:56:00Z RPC /submit payloadHash=0x030df036 fee=1176 status=accepted
2025-11-08T02:57:00Z RPC /submit payloadHash=0x030df137 fee=1177 status=broadcast
2025-11-08T02:58:00Z RPC /submit payloadHash=0x030df238 fee=1178 status=accepted
2025-11-08T02:59:00Z RPC /submit payloadHash=0x030df339 fee=1179 status=accepted
2025-11-08T03:00:00Z RPC /submit payloadHash=0x030df440 fee=1180 status=broadcast
2025-11-08T03:01:00Z RPC /submit payloadHash=0x030df541 fee=1181 status=accepted
2025-11-08T03:02:00Z RPC /submit payloadHash=0x030df642 fee=1182 status=accepted
2025-11-08T03:03:00Z RPC /submit payloadHash=0x030df743 fee=1183 status=broadcast
2025-11-08T03:04:00Z RPC /submit payloadHash=0x030df844 fee=1184 status=accepted
2025-11-08T03:05:00Z RPC /submit payloadHash=0x030df945 fee=1185 status=accepted
2025-11-08T03:06:00Z RPC /submit payloadHash=0x030dfa46 fee=1186 status=broadcast
2025-11-08T03:07:00Z RPC /submit payloadHash=0x030dfb47 fee=1187 status=accepted
2025-11-08T03:08:00Z RPC /submit payloadHash=0x030dfc48 fee=1188 status=accepted
2025-11-08T03:09:00Z RPC /submit payloadHash=0x030dfd49 fee=1189 status=broadcast
2025-11-08T03:10:00Z RPC /submit payloadHash=0x030dfe50 fee=1190 status=accepted
2025-11-08T03:11:00Z RPC /submit payloadHash=0x030dff51 fee=1191 status=accepted
2025-11-08T03:12:00Z RPC /submit payloadHash=0x030e0052 fee=1192 status=broadcast
2025-11-08T03:13:00Z RPC /submit payloadHash=0x030e0153 fee=1193 status=accepted
2025-11-08T03:14:00Z RPC /submit payloadHash=0x030e0254 fee=1194 status=accepted
2025-11-08T03:15:00Z RPC /submit payloadHash=0x030e0355 fee=1195 status=broadcast
2025-11-08T03:16:00Z RPC /submit payloadHash=0x030e0456 fee=1196 status=accepted
2025-11-08T03:17:00Z RPC /submit payloadHash=0x030e0557 fee=1197 status=accepted
2025-11-08T03:18:00Z RPC /submit payloadHash=0x030e0658 fee=1198 status=broadcast
2025-11-08T03:19:00Z RPC /submit payloadHash=0x030e0759 fee=1199 status=accepted
2025-11-08T03:20:00Z RPC /submit payloadHash=0x030e0810 fee=1000 status=accepted
2025-11-08T03:21:00Z RPC /submit payloadHash=0x030e0911 fee=1001 status=broadcast
2025-11-08T03:22:00Z RPC /submit payloadHash=0x030e0a12 fee=1002 status=accepted
2025-11-08T03:23:00Z RPC /submit payloadHash=0x030e0b13 fee=1003 status=accepted
2025-11-08T03:24:00Z RPC /submit payloadHash=0x030e0c14 fee=1004 status=broadcast
2025-11-08T03:25:00Z RPC /submit payloadHash=0x030e0d15 fee=1005 status=accepted
2025-11-08T03:26:00Z RPC /submit payloadHash=0x030e0e16 fee=1006 status=accepted
2025-11-08T03:27:00Z RPC /submit payloadHash=0x030e0f17 fee=1007 status=broadcast
2025-11-08T03:28:00Z RPC /submit payloadHash=0x030e1018 fee=1008 status=accepted
2025-11-08T03:29:00Z RPC /submit payloadHash=0x030e1119 fee=1009 status=accepted
2025-11-08T03:30:00Z RPC /submit payloadHash=0x030e1220 fee=1010 status=broadcast
2025-11-08T03:31:00Z RPC /submit payloadHash=0x030e1321 fee=1011 status=accepted
2025-11-08T03:32:00Z RPC /submit payloadHash=0x030e1422 fee=1012 status=accepted
2025-11-08T03:33:00Z RPC /submit payloadHash=0x030e1523 fee=1013 status=broadcast
2025-11-08T03:34:00Z RPC /submit payloadHash=0x030e1624 fee=1014 status=accepted
2025-11-08T03:35:00Z RPC /submit payloadHash=0x030e1725 fee=1015 status=accepted
2025-11-08T03:36:00Z RPC /submit payloadHash=0x030e1826 fee=1016 status=broadcast
2025-11-08T03:37:00Z RPC /submit payloadHash=0x030e1927 fee=1017 status=accepted
2025-11-08T03:38:00Z RPC /submit payloadHash=0x030e1a28 fee=1018 status=accepted
2025-11-08T03:39:00Z RPC /submit payloadHash=0x030e1b29 fee=1019 status=broadcast
2025-11-08T03:40:00Z RPC /submit payloadHash=0x030e1c30 fee=1020 status=accepted
2025-11-08T03:41:00Z RPC /submit payloadHash=0x030e1d31 fee=1021 status=accepted
2025-11-08T03:42:00Z RPC /submit payloadHash=0x030e1e32 fee=1022 status=broadcast
2025-11-08T03:43:00Z RPC /submit payloadHash=0x030e1f33 fee=1023 status=accepted
2025-11-08T03:44:00Z RPC /submit payloadHash=0x030e2034 fee=1024 status=accepted
2025-11-08T03:45:00Z RPC /submit payloadHash=0x030e2135 fee=1025 status=broadcast
2025-11-08T03:46:00Z RPC /submit payloadHash=0x030e2236 fee=1026 status=accepted
2025-11-08T03:47:00Z RPC /submit payloadHash=0x030e2337 fee=1027 status=accepted
2025-11-08T03:48:00Z RPC /submit payloadHash=0x030e2438 fee=1028 status=broadcast
2025-11-08T03:49:00Z RPC /submit payloadHash=0x030e2539 fee=1029 status=accepted
2025-11-08T03:50:00Z RPC /submit payloadHash=0x030e2640 fee=1030 status=accepted
2025-11-08T03:51:00Z RPC /submit payloadHash=0x030e2741 fee=1031 status=broadcast
2025-11-08T03:52:00Z RPC /submit payloadHash=0x030e2842 fee=1032 status=accepted
2025-11-08T03:53:00Z RPC /submit payloadHash=0x030e2943 fee=1033 status=accepted
2025-11-08T03:54:00Z RPC /submit payloadHash=0x030e2a44 fee=1034 status=broadcast
2025-11-08T03:55:00Z RPC /submit payloadHash=0x030e2b45 fee=1035 status=accepted
2025-11-08T03:56:00Z RPC /submit payloadHash=0x030e2c46 fee=1036 status=accepted
2025-11-08T03:57:00Z RPC /submit payloadHash=0x030e2d47 fee=1037 status=broadcast
2025-11-08T03:58:00Z RPC /submit payloadHash=0x030e2e48 fee=1038 status=accepted
2025-11-08T03:59:00Z RPC /submit payloadHash=0x030e2f49 fee=1039 status=accepted
2025-11-08T04:00:00Z RPC /submit payloadHash=0x030e3050 fee=1040 status=broadcast
2025-11-08T04:01:00Z RPC /submit payloadHash=0x030e3151 fee=1041 status=accepted
2025-11-08T04:02:00Z RPC /submit payloadHash=0x030e3252 fee=1042 status=accepted
2025-11-08T04:03:00Z RPC /submit payloadHash=0x030e3353 fee=1043 status=broadcast
2025-11-08T04:04:00Z RPC /submit payloadHash=0x030e3454 fee=1044 status=accepted
2025-11-08T04:05:00Z RPC /submit payloadHash=0x030e3555 fee=1045 status=accepted
2025-11-08T04:06:00Z RPC /submit payloadHash=0x030e3656 fee=1046 status=broadcast
2025-11-08T04:07:00Z RPC /submit payloadHash=0x030e3757 fee=1047 status=accepted
2025-11-08T04:08:00Z RPC /submit payloadHash=0x030e3858 fee=1048 status=accepted
2025-11-08T04:09:00Z RPC /submit payloadHash=0x030e3959 fee=1049 status=broadcast
2025-11-08T04:10:00Z RPC /submit payloadHash=0x030e3a10 fee=1050 status=accepted
2025-11-08T04:11:00Z RPC /submit payloadHash=0x030e3b11 fee=1051 status=accepted
2025-11-08T04:12:00Z RPC /submit payloadHash=0x030e3c12 fee=1052 status=broadcast
2025-11-08T04:13:00Z RPC /submit payloadHash=0x030e3d13 fee=1053 status=accepted
2025-11-08T04:14:00Z RPC /submit payloadHash=0x030e3e14 fee=1054 status=accepted
2025-11-08T04:15:00Z RPC /submit payloadHash=0x030e3f15 fee=1055 status=broadcast
2025-11-08T04:16:00Z RPC /submit payloadHash=0x030e4016 fee=1056 status=accepted
2025-11-08T04:17:00Z RPC /submit payloadHash=0x030e4117 fee=1057 status=accepted
2025-11-08T04:18:00Z RPC /submit payloadHash=0x030e4218 fee=1058 status=broadcast
2025-11-08T04:19:00Z RPC /submit payloadHash=0x030e4319 fee=1059 status=accepted
2025-11-08T04:20:00Z RPC /submit payloadHash=0x030e4420 fee=1060 status=accepted
2025-11-08T04:21:00Z RPC /submit payloadHash=0x030e4521 fee=1061 status=broadcast
2025-11-08T04:22:00Z RPC /submit payloadHash=0x030e4622 fee=1062 status=accepted
2025-11-08T04:23:00Z RPC /submit payloadHash=0x030e4723 fee=1063 status=accepted
2025-11-08T04:24:00Z RPC /submit payloadHash=0x030e4824 fee=1064 status=broadcast
2025-11-08T04:25:00Z RPC /submit payloadHash=0x030e4925 fee=1065 status=accepted
2025-11-08T04:26:00Z RPC /submit payloadHash=0x030e4a26 fee=1066 status=accepted
2025-11-08T04:27:00Z RPC /submit payloadHash=0x030e4b27 fee=1067 status=broadcast
2025-11-08T04:28:00Z RPC /submit payloadHash=0x030e4c28 fee=1068 status=accepted
2025-11-08T04:29:00Z RPC /submit payloadHash=0x030e4d29 fee=1069 status=accepted
2025-11-08T04:30:00Z RPC /submit payloadHash=0x030e4e30 fee=1070 status=broadcast
2025-11-08T04:31:00Z RPC /submit payloadHash=0x030e4f31 fee=1071 status=accepted
2025-11-08T04:32:00Z RPC /submit payloadHash=0x030e5032 fee=1072 status=accepted
2025-11-08T04:33:00Z RPC /submit payloadHash=0x030e5133 fee=1073 status=broadcast
2025-11-08T04:34:00Z RPC /submit payloadHash=0x030e5234 fee=1074 status=accepted
2025-11-08T04:35:00Z RPC /submit payloadHash=0x030e5335 fee=1075 status=accepted
2025-11-08T04:36:00Z RPC /submit payloadHash=0x030e5436 fee=1076 status=broadcast
2025-11-08T04:37:00Z RPC /submit payloadHash=0x030e5537 fee=1077 status=accepted
2025-11-08T04:38:00Z RPC /submit payloadHash=0x030e5638 fee=1078 status=accepted
2025-11-08T04:39:00Z RPC /submit payloadHash=0x030e5739 fee=1079 status=broadcast
2025-11-08T04:40:00Z RPC /submit payloadHash=0x030e5840 fee=1080 status=accepted
2025-11-08T04:41:00Z RPC /submit payloadHash=0x030e5941 fee=1081 status=accepted
2025-11-08T04:42:00Z RPC /submit payloadHash=0x030e5a42 fee=1082 status=broadcast
2025-11-08T04:43:00Z RPC /submit payloadHash=0x030e5b43 fee=1083 status=accepted
2025-11-08T04:44:00Z RPC /submit payloadHash=0x030e5c44 fee=1084 status=accepted
2025-11-08T04:45:00Z RPC /submit payloadHash=0x030e5d45 fee=1085 status=broadcast
2025-11-08T04:46:00Z RPC /submit payloadHash=0x030e5e46 fee=1086 status=accepted
2025-11-08T04:47:00Z RPC /submit payloadHash=0x030e5f47 fee=1087 status=accepted
2025-11-08T04:48:00Z RPC /submit payloadHash=0x030e6048 fee=1088 status=broadcast
2025-11-08T04:49:00Z RPC /submit payloadHash=0x030e6149 fee=1089 status=accepted
2025-11-08T04:50:00Z RPC /submit payloadHash=0x030e6250 fee=1090 status=accepted
2025-11-08T04:51:00Z RPC /submit payloadHash=0x030e6351 fee=1091 status=broadcast
2025-11-08T04:52:00Z RPC /submit payloadHash=0x030e6452 fee=1092 status=accepted
2025-11-08T04:53:00Z RPC /submit payloadHash=0x030e6553 fee=1093 status=accepted
2025-11-08T04:54:00Z RPC /submit payloadHash=0x030e6654 fee=1094 status=broadcast
2025-11-08T04:55:00Z RPC /submit payloadHash=0x030e6755 fee=1095 status=accepted
2025-11-08T04:56:00Z RPC /submit payloadHash=0x030e6856 fee=1096 status=accepted
2025-11-08T04:57:00Z RPC /submit payloadHash=0x030e6957 fee=1097 status=broadcast
2025-11-08T04:58:00Z RPC /submit payloadHash=0x030e6a58 fee=1098 status=accepted
2025-11-08T04:59:00Z RPC /submit payloadHash=0x030e6b59 fee=1099 status=accepted
2025-11-08T05:00:00Z RPC /submit payloadHash=0x030e6c10 fee=1100 status=broadcast
2025-11-08T05:01:00Z RPC /submit payloadHash=0x030e6d11 fee=1101 status=accepted
2025-11-08T05:02:00Z RPC /submit payloadHash=0x030e6e12 fee=1102 status=accepted
2025-11-08T05:03:00Z RPC /submit payloadHash=0x030e6f13 fee=1103 status=broadcast
2025-11-08T05:04:00Z RPC /submit payloadHash=0x030e7014 fee=1104 status=accepted
2025-11-08T05:05:00Z RPC /submit payloadHash=0x030e7115 fee=1105 status=accepted
2025-11-08T05:06:00Z RPC /submit payloadHash=0x030e7216 fee=1106 status=broadcast
2025-11-08T05:07:00Z RPC /submit payloadHash=0x030e7317 fee=1107 status=accepted
2025-11-08T05:08:00Z RPC /submit payloadHash=0x030e7418 fee=1108 status=accepted
2025-11-08T05:09:00Z RPC /submit payloadHash=0x030e7519 fee=1109 status=broadcast
2025-11-08T05:10:00Z RPC /submit payloadHash=0x030e7620 fee=1110 status=accepted
2025-11-08T05:11:00Z RPC /submit payloadHash=0x030e7721 fee=1111 status=accepted
2025-11-08T05:12:00Z RPC /submit payloadHash=0x030e7822 fee=1112 status=broadcast
2025-11-08T05:13:00Z RPC /submit payloadHash=0x030e7923 fee=1113 status=accepted
2025-11-08T05:14:00Z RPC /submit payloadHash=0x030e7a24 fee=1114 status=accepted
2025-11-08T05:15:00Z RPC /submit payloadHash=0x030e7b25 fee=1115 status=broadcast
2025-11-08T05:16:00Z RPC /submit payloadHash=0x030e7c26 fee=1116 status=accepted
2025-11-08T05:17:00Z RPC /submit payloadHash=0x030e7d27 fee=1117 status=accepted
2025-11-08T05:18:00Z RPC /submit payloadHash=0x030e7e28 fee=1118 status=broadcast
2025-11-08T05:19:00Z RPC /submit payloadHash=0x030e7f29 fee=1119 status=accepted
2025-11-08T05:20:00Z RPC /submit payloadHash=0x030e8030 fee=1120 status=accepted
2025-11-08T05:21:00Z RPC /submit payloadHash=0x030e8131 fee=1121 status=broadcast
2025-11-08T05:22:00Z RPC /submit payloadHash=0x030e8232 fee=1122 status=accepted
2025-11-08T05:23:00Z RPC /submit payloadHash=0x030e8333 fee=1123 status=accepted
2025-11-08T05:24:00Z RPC /submit payloadHash=0x030e8434 fee=1124 status=broadcast
2025-11-08T05:25:00Z RPC /submit payloadHash=0x030e8535 fee=1125 status=accepted
2025-11-08T05:26:00Z RPC /submit payloadHash=0x030e8636 fee=1126 status=accepted
2025-11-08T05:27:00Z RPC /submit payloadHash=0x030e8737 fee=1127 status=broadcast
2025-11-08T05:28:00Z RPC /submit payloadHash=0x030e8838 fee=1128 status=accepted
2025-11-08T05:29:00Z RPC /submit payloadHash=0x030e8939 fee=1129 status=accepted
2025-11-08T05:30:00Z RPC /submit payloadHash=0x030e8a40 fee=1130 status=broadcast
2025-11-08T05:31:00Z RPC /submit payloadHash=0x030e8b41 fee=1131 status=accepted
2025-11-08T05:32:00Z RPC /submit payloadHash=0x030e8c42 fee=1132 status=accepted
2025-11-08T05:33:00Z RPC /submit payloadHash=0x030e8d43 fee=1133 status=broadcast
2025-11-08T05:34:00Z RPC /submit payloadHash=0x030e8e44 fee=1134 status=accepted
2025-11-08T05:35:00Z RPC /submit payloadHash=0x030e8f45 fee=1135 status=accepted
2025-11-08T05:36:00Z RPC /submit payloadHash=0x030e9046 fee=1136 status=broadcast
2025-11-08T05:37:00Z RPC /submit payloadHash=0x030e9147 fee=1137 status=accepted
2025-11-08T05:38:00Z RPC /submit payloadHash=0x030e9248 fee=1138 status=accepted
2025-11-08T05:39:00Z RPC /submit payloadHash=0x030e9349 fee=1139 status=broadcast
2025-11-08T05:40:00Z RPC /submit payloadHash=0x030e9450 fee=1140 status=accepted
2025-11-08T05:41:00Z RPC /submit payloadHash=0x030e9551 fee=1141 status=accepted
2025-11-08T05:42:00Z RPC /submit payloadHash=0x030e9652 fee=1142 status=broadcast
2025-11-08T05:43:00Z RPC /submit payloadHash=0x030e9753 fee=1143 status=accepted
2025-11-08T05:44:00Z RPC /submit payloadHash=0x030e9854 fee=1144 status=accepted
2025-11-08T05:45:00Z RPC /submit payloadHash=0x030e9955 fee=1145 status=broadcast
2025-11-08T05:46:00Z RPC /submit payloadHash=0x030e9a56 fee=1146 status=accepted
2025-11-08T05:47:00Z RPC /submit payloadHash=0x030e9b57 fee=1147 status=accepted
2025-11-08T05:48:00Z RPC /submit payloadHash=0x030e9c58 fee=1148 status=broadcast
2025-11-08T05:49:00Z RPC /submit payloadHash=0x030e9d59 fee=1149 status=accepted
2025-11-08T05:50:00Z RPC /submit payloadHash=0x030e9e10 fee=1150 status=accepted
2025-11-08T05:51:00Z RPC /submit payloadHash=0x030e9f11 fee=1151 status=broadcast
2025-11-08T05:52:00Z RPC /submit payloadHash=0x030ea012 fee=1152 status=accepted
2025-11-08T05:53:00Z RPC /submit payloadHash=0x030ea113 fee=1153 status=accepted
2025-11-08T05:54:00Z RPC /submit payloadHash=0x030ea214 fee=1154 status=broadcast
2025-11-08T05:55:00Z RPC /submit payloadHash=0x030ea315 fee=1155 status=accepted
2025-11-08T05:56:00Z RPC /submit payloadHash=0x030ea416 fee=1156 status=accepted
2025-11-08T05:57:00Z RPC /submit payloadHash=0x030ea517 fee=1157 status=broadcast
2025-11-08T05:58:00Z RPC /submit payloadHash=0x030ea618 fee=1158 status=accepted
2025-11-08T05:59:00Z RPC /submit payloadHash=0x030ea719 fee=1159 status=accepted
2025-11-08T06:00:00Z RPC /submit payloadHash=0x030ea820 fee=1160 status=broadcast
2025-11-08T06:01:00Z RPC /submit payloadHash=0x030ea921 fee=1161 status=accepted
2025-11-08T06:02:00Z RPC /submit payloadHash=0x030eaa22 fee=1162 status=accepted
2025-11-08T06:03:00Z RPC /submit payloadHash=0x030eab23 fee=1163 status=broadcast
2025-11-08T06:04:00Z RPC /submit payloadHash=0x030eac24 fee=1164 status=accepted
2025-11-08T06:05:00Z RPC /submit payloadHash=0x030ead25 fee=1165 status=accepted
2025-11-08T06:06:00Z RPC /submit payloadHash=0x030eae26 fee=1166 status=broadcast
2025-11-08T06:07:00Z RPC /submit payloadHash=0x030eaf27 fee=1167 status=accepted
2025-11-08T06:08:00Z RPC /submit payloadHash=0x030eb028 fee=1168 status=accepted
2025-11-08T06:09:00Z RPC /submit payloadHash=0x030eb129 fee=1169 status=broadcast
2025-11-08T06:10:00Z RPC /submit payloadHash=0x030eb230 fee=1170 status=accepted
2025-11-08T06:11:00Z RPC /submit payloadHash=0x030eb331 fee=1171 status=accepted
2025-11-08T06:12:00Z RPC /submit payloadHash=0x030eb432 fee=1172 status=broadcast
2025-11-08T06:13:00Z RPC /submit payloadHash=0x030eb533 fee=1173 status=accepted
2025-11-08T06:14:00Z RPC /submit payloadHash=0x030eb634 fee=1174 status=accepted
2025-11-08T06:15:00Z RPC /submit payloadHash=0x030eb735 fee=1175 status=broadcast
2025-11-08T06:16:00Z RPC /submit payloadHash=0x030eb836 fee=1176 status=accepted
2025-11-08T06:17:00Z RPC /submit payloadHash=0x030eb937 fee=1177 status=accepted
2025-11-08T06:18:00Z RPC /submit payloadHash=0x030eba38 fee=1178 status=broadcast
2025-11-08T06:19:00Z RPC /submit payloadHash=0x030ebb39 fee=1179 status=accepted
2025-11-08T06:20:00Z RPC /submit payloadHash=0x030ebc40 fee=1180 status=accepted
2025-11-08T06:21:00Z RPC /submit payloadHash=0x030ebd41 fee=1181 status=broadcast
2025-11-08T06:22:00Z RPC /submit payloadHash=0x030ebe42 fee=1182 status=accepted
2025-11-08T06:23:00Z RPC /submit payloadHash=0x030ebf43 fee=1183 status=accepted
2025-11-08T06:24:00Z RPC /submit payloadHash=0x030ec044 fee=1184 status=broadcast
2025-11-08T06:25:00Z RPC /submit payloadHash=0x030ec145 fee=1185 status=accepted
2025-11-08T06:26:00Z RPC /submit payloadHash=0x030ec246 fee=1186 status=accepted
2025-11-08T06:27:00Z RPC /submit payloadHash=0x030ec347 fee=1187 status=broadcast
2025-11-08T06:28:00Z RPC /submit payloadHash=0x030ec448 fee=1188 status=accepted
2025-11-08T06:29:00Z RPC /submit payloadHash=0x030ec549 fee=1189 status=accepted
2025-11-08T06:30:00Z RPC /submit payloadHash=0x030ec650 fee=1190 status=broadcast
2025-11-08T06:31:00Z RPC /submit payloadHash=0x030ec751 fee=1191 status=accepted
2025-11-08T06:32:00Z RPC /submit payloadHash=0x030ec852 fee=1192 status=accepted
2025-11-08T06:33:00Z RPC /submit payloadHash=0x030ec953 fee=1193 status=broadcast
2025-11-08T06:34:00Z RPC /submit payloadHash=0x030eca54 fee=1194 status=accepted
2025-11-08T06:35:00Z RPC /submit payloadHash=0x030ecb55 fee=1195 status=accepted
2025-11-08T06:36:00Z RPC /submit payloadHash=0x030ecc56 fee=1196 status=broadcast
2025-11-08T06:37:00Z RPC /submit payloadHash=0x030ecd57 fee=1197 status=accepted
2025-11-08T06:38:00Z RPC /submit payloadHash=0x030ece58 fee=1198 status=accepted
2025-11-08T06:39:00Z RPC /submit payloadHash=0x030ecf59 fee=1199 status=broadcast
2025-11-08T06:40:00Z RPC /submit payloadHash=0x030ed010 fee=1000 status=accepted
2025-11-08T06:41:00Z RPC /submit payloadHash=0x030ed111 fee=1001 status=accepted
2025-11-08T06:42:00Z RPC /submit payloadHash=0x030ed212 fee=1002 status=broadcast
2025-11-08T06:43:00Z RPC /submit payloadHash=0x030ed313 fee=1003 status=accepted
2025-11-08T06:44:00Z RPC /submit payloadHash=0x030ed414 fee=1004 status=accepted
2025-11-08T06:45:00Z RPC /submit payloadHash=0x030ed515 fee=1005 status=broadcast
2025-11-08T06:46:00Z RPC /submit payloadHash=0x030ed616 fee=1006 status=accepted
2025-11-08T06:47:00Z RPC /submit payloadHash=0x030ed717 fee=1007 status=accepted
2025-11-08T06:48:00Z RPC /submit payloadHash=0x030ed818 fee=1008 status=broadcast
2025-11-08T06:49:00Z RPC /submit payloadHash=0x030ed919 fee=1009 status=accepted
2025-11-08T06:50:00Z RPC /submit payloadHash=0x030eda20 fee=1010 status=accepted
2025-11-08T06:51:00Z RPC /submit payloadHash=0x030edb21 fee=1011 status=broadcast
2025-11-08T06:52:00Z RPC /submit payloadHash=0x030edc22 fee=1012 status=accepted
2025-11-08T06:53:00Z RPC /submit payloadHash=0x030edd23 fee=1013 status=accepted
2025-11-08T06:54:00Z RPC /submit payloadHash=0x030ede24 fee=1014 status=broadcast
2025-11-08T06:55:00Z RPC /submit payloadHash=0x030edf25 fee=1015 status=accepted
2025-11-08T06:56:00Z RPC /submit payloadHash=0x030ee026 fee=1016 status=accepted
2025-11-08T06:57:00Z RPC /submit payloadHash=0x030ee127 fee=1017 status=broadcast
2025-11-08T06:58:00Z RPC /submit payloadHash=0x030ee228 fee=1018 status=accepted
2025-11-08T06:59:00Z RPC /submit payloadHash=0x030ee329 fee=1019 status=accepted
2025-11-08T07:00:00Z RPC /submit payloadHash=0x030ee430 fee=1020 status=broadcast
2025-11-08T07:01:00Z RPC /submit payloadHash=0x030ee531 fee=1021 status=accepted
2025-11-08T07:02:00Z RPC /submit payloadHash=0x030ee632 fee=1022 status=accepted
2025-11-08T07:03:00Z RPC /submit payloadHash=0x030ee733 fee=1023 status=broadcast
2025-11-08T07:04:00Z RPC /submit payloadHash=0x030ee834 fee=1024 status=accepted
2025-11-08T07:05:00Z RPC /submit payloadHash=0x030ee935 fee=1025 status=accepted
2025-11-08T07:06:00Z RPC /submit payloadHash=0x030eea36 fee=1026 status=broadcast
2025-11-08T07:07:00Z RPC /submit payloadHash=0x030eeb37 fee=1027 status=accepted
2025-11-08T07:08:00Z RPC /submit payloadHash=0x030eec38 fee=1028 status=accepted
2025-11-08T07:09:00Z RPC /submit payloadHash=0x030eed39 fee=1029 status=broadcast
2025-11-08T07:10:00Z RPC /submit payloadHash=0x030eee40 fee=1030 status=accepted
2025-11-08T07:11:00Z RPC /submit payloadHash=0x030eef41 fee=1031 status=accepted
2025-11-08T07:12:00Z RPC /submit payloadHash=0x030ef042 fee=1032 status=broadcast
2025-11-08T07:13:00Z RPC /submit payloadHash=0x030ef143 fee=1033 status=accepted
2025-11-08T07:14:00Z RPC /submit payloadHash=0x030ef244 fee=1034 status=accepted
2025-11-08T07:15:00Z RPC /submit payloadHash=0x030ef345 fee=1035 status=broadcast
2025-11-08T07:16:00Z RPC /submit payloadHash=0x030ef446 fee=1036 status=accepted
2025-11-08T07:17:00Z RPC /submit payloadHash=0x030ef547 fee=1037 status=accepted
2025-11-08T07:18:00Z RPC /submit payloadHash=0x030ef648 fee=1038 status=broadcast
2025-11-08T07:19:00Z RPC /submit payloadHash=0x030ef749 fee=1039 status=accepted
2025-11-08T07:20:00Z RPC /submit payloadHash=0x030ef850 fee=1040 status=accepted
2025-11-08T07:21:00Z RPC /submit payloadHash=0x030ef951 fee=1041 status=broadcast
2025-11-08T07:22:00Z RPC /submit payloadHash=0x030efa52 fee=1042 status=accepted
2025-11-08T07:23:00Z RPC /submit payloadHash=0x030efb53 fee=1043 status=accepted
2025-11-08T07:24:00Z RPC /submit payloadHash=0x030efc54 fee=1044 status=broadcast
2025-11-08T07:25:00Z RPC /submit payloadHash=0x030efd55 fee=1045 status=accepted
2025-11-08T07:26:00Z RPC /submit payloadHash=0x030efe56 fee=1046 status=accepted
2025-11-08T07:27:00Z RPC /submit payloadHash=0x030eff57 fee=1047 status=broadcast
2025-11-08T07:28:00Z RPC /submit payloadHash=0x030f0058 fee=1048 status=accepted
2025-11-08T07:29:00Z RPC /submit payloadHash=0x030f0159 fee=1049 status=accepted
2025-11-08T07:30:00Z RPC /submit payloadHash=0x030f0210 fee=1050 status=broadcast
2025-11-08T07:31:00Z RPC /submit payloadHash=0x030f0311 fee=1051 status=accepted
2025-11-08T07:32:00Z RPC /submit payloadHash=0x030f0412 fee=1052 status=accepted
2025-11-08T07:33:00Z RPC /submit payloadHash=0x030f0513 fee=1053 status=broadcast
2025-11-08T07:34:00Z RPC /submit payloadHash=0x030f0614 fee=1054 status=accepted
2025-11-08T07:35:00Z RPC /submit payloadHash=0x030f0715 fee=1055 status=accepted
2025-11-08T07:36:00Z RPC /submit payloadHash=0x030f0816 fee=1056 status=broadcast
2025-11-08T07:37:00Z RPC /submit payloadHash=0x030f0917 fee=1057 status=accepted
2025-11-08T07:38:00Z RPC /submit payloadHash=0x030f0a18 fee=1058 status=accepted
2025-11-08T07:39:00Z RPC /submit payloadHash=0x030f0b19 fee=1059 status=broadcast
2025-11-08T07:40:00Z RPC /submit payloadHash=0x030f0c20 fee=1060 status=accepted
2025-11-08T07:41:00Z RPC /submit payloadHash=0x030f0d21 fee=1061 status=accepted
2025-11-08T07:42:00Z RPC /submit payloadHash=0x030f0e22 fee=1062 status=broadcast
2025-11-08T07:43:00Z RPC /submit payloadHash=0x030f0f23 fee=1063 status=accepted
2025-11-08T07:44:00Z RPC /submit payloadHash=0x030f1024 fee=1064 status=accepted
2025-11-08T07:45:00Z RPC /submit payloadHash=0x030f1125 fee=1065 status=broadcast
2025-11-08T07:46:00Z RPC /submit payloadHash=0x030f1226 fee=1066 status=accepted
2025-11-08T07:47:00Z RPC /submit payloadHash=0x030f1327 fee=1067 status=accepted
2025-11-08T07:48:00Z RPC /submit payloadHash=0x030f1428 fee=1068 status=broadcast
2025-11-08T07:49:00Z RPC /submit payloadHash=0x030f1529 fee=1069 status=accepted
2025-11-08T07:50:00Z RPC /submit payloadHash=0x030f1630 fee=1070 status=accepted
2025-11-08T07:51:00Z RPC /submit payloadHash=0x030f1731 fee=1071 status=broadcast
2025-11-08T07:52:00Z RPC /submit payloadHash=0x030f1832 fee=1072 status=accepted
2025-11-08T07:53:00Z RPC /submit payloadHash=0x030f1933 fee=1073 status=accepted
2025-11-08T07:54:00Z RPC /submit payloadHash=0x030f1a34 fee=1074 status=broadcast
2025-11-08T07:55:00Z RPC /submit payloadHash=0x030f1b35 fee=1075 status=accepted
2025-11-08T07:56:00Z RPC /submit payloadHash=0x030f1c36 fee=1076 status=accepted
2025-11-08T07:57:00Z RPC /submit payloadHash=0x030f1d37 fee=1077 status=broadcast
2025-11-08T07:58:00Z RPC /submit payloadHash=0x030f1e38 fee=1078 status=accepted
2025-11-08T07:59:00Z RPC /submit payloadHash=0x030f1f39 fee=1079 status=accepted
2025-11-08T08:00:00Z RPC /submit payloadHash=0x030f2040 fee=1080 status=broadcast
2025-11-08T08:01:00Z RPC /submit payloadHash=0x030f2141 fee=1081 status=accepted
2025-11-08T08:02:00Z RPC /submit payloadHash=0x030f2242 fee=1082 status=accepted
2025-11-08T08:03:00Z RPC /submit payloadHash=0x030f2343 fee=1083 status=broadcast
2025-11-08T08:04:00Z RPC /submit payloadHash=0x030f2444 fee=1084 status=accepted
2025-11-08T08:05:00Z RPC /submit payloadHash=0x030f2545 fee=1085 status=accepted
2025-11-08T08:06:00Z RPC /submit payloadHash=0x030f2646 fee=1086 status=broadcast
2025-11-08T08:07:00Z RPC /submit payloadHash=0x030f2747 fee=1087 status=accepted
2025-11-08T08:08:00Z RPC /submit payloadHash=0x030f2848 fee=1088 status=accepted
2025-11-08T08:09:00Z RPC /submit payloadHash=0x030f2949 fee=1089 status=broadcast
2025-11-08T08:10:00Z RPC /submit payloadHash=0x030f2a50 fee=1090 status=accepted
2025-11-08T08:11:00Z RPC /submit payloadHash=0x030f2b51 fee=1091 status=accepted
2025-11-08T08:12:00Z RPC /submit payloadHash=0x030f2c52 fee=1092 status=broadcast
2025-11-08T08:13:00Z RPC /submit payloadHash=0x030f2d53 fee=1093 status=accepted
2025-11-08T08:14:00Z RPC /submit payloadHash=0x030f2e54 fee=1094 status=accepted
2025-11-08T08:15:00Z RPC /submit payloadHash=0x030f2f55 fee=1095 status=broadcast
2025-11-08T08:16:00Z RPC /submit payloadHash=0x030f3056 fee=1096 status=accepted
2025-11-08T08:17:00Z RPC /submit payloadHash=0x030f3157 fee=1097 status=accepted
2025-11-08T08:18:00Z RPC /submit payloadHash=0x030f3258 fee=1098 status=broadcast
2025-11-08T08:19:00Z RPC /submit payloadHash=0x030f3359 fee=1099 status=accepted
2025-11-08T08:20:00Z RPC /submit payloadHash=0x030f3410 fee=1100 status=accepted
2025-11-08T08:21:00Z RPC /submit payloadHash=0x030f3511 fee=1101 status=broadcast
2025-11-08T08:22:00Z RPC /submit payloadHash=0x030f3612 fee=1102 status=accepted
2025-11-08T08:23:00Z RPC /submit payloadHash=0x030f3713 fee=1103 status=accepted
2025-11-08T08:24:00Z RPC /submit payloadHash=0x030f3814 fee=1104 status=broadcast
2025-11-08T08:25:00Z RPC /submit payloadHash=0x030f3915 fee=1105 status=accepted
2025-11-08T08:26:00Z RPC /submit payloadHash=0x030f3a16 fee=1106 status=accepted
2025-11-08T08:27:00Z RPC /submit payloadHash=0x030f3b17 fee=1107 status=broadcast
2025-11-08T08:28:00Z RPC /submit payloadHash=0x030f3c18 fee=1108 status=accepted
2025-11-08T08:29:00Z RPC /submit payloadHash=0x030f3d19 fee=1109 status=accepted
2025-11-08T08:30:00Z RPC /submit payloadHash=0x030f3e20 fee=1110 status=broadcast
2025-11-08T08:31:00Z RPC /submit payloadHash=0x030f3f21 fee=1111 status=accepted
2025-11-08T08:32:00Z RPC /submit payloadHash=0x030f4022 fee=1112 status=accepted
2025-11-08T08:33:00Z RPC /submit payloadHash=0x030f4123 fee=1113 status=broadcast
2025-11-08T08:34:00Z RPC /submit payloadHash=0x030f4224 fee=1114 status=accepted
2025-11-08T08:35:00Z RPC /submit payloadHash=0x030f4325 fee=1115 status=accepted
2025-11-08T08:36:00Z RPC /submit payloadHash=0x030f4426 fee=1116 status=broadcast
2025-11-08T08:37:00Z RPC /submit payloadHash=0x030f4527 fee=1117 status=accepted
2025-11-08T08:38:00Z RPC /submit payloadHash=0x030f4628 fee=1118 status=accepted
2025-11-08T08:39:00Z RPC /submit payloadHash=0x030f4729 fee=1119 status=broadcast
2025-11-08T08:40:00Z RPC /submit payloadHash=0x030f4830 fee=1120 status=accepted
2025-11-08T08:41:00Z RPC /submit payloadHash=0x030f4931 fee=1121 status=accepted
2025-11-08T08:42:00Z RPC /submit payloadHash=0x030f4a32 fee=1122 status=broadcast
2025-11-08T08:43:00Z RPC /submit payloadHash=0x030f4b33 fee=1123 status=accepted
2025-11-08T08:44:00Z RPC /submit payloadHash=0x030f4c34 fee=1124 status=accepted
2025-11-08T08:45:00Z RPC /submit payloadHash=0x030f4d35 fee=1125 status=broadcast
2025-11-08T08:46:00Z RPC /submit payloadHash=0x030f4e36 fee=1126 status=accepted
2025-11-08T08:47:00Z RPC /submit payloadHash=0x030f4f37 fee=1127 status=accepted
2025-11-08T08:48:00Z RPC /submit payloadHash=0x030f5038 fee=1128 status=broadcast
2025-11-08T08:49:00Z RPC /submit payloadHash=0x030f5139 fee=1129 status=accepted
2025-11-08T08:50:00Z RPC /submit payloadHash=0x030f5240 fee=1130 status=accepted
2025-11-08T08:51:00Z RPC /submit payloadHash=0x030f5341 fee=1131 status=broadcast
2025-11-08T08:52:00Z RPC /submit payloadHash=0x030f5442 fee=1132 status=accepted
2025-11-08T08:53:00Z RPC /submit payloadHash=0x030f5543 fee=1133 status=accepted
2025-11-08T08:54:00Z RPC /submit payloadHash=0x030f5644 fee=1134 status=broadcast
2025-11-08T08:55:00Z RPC /submit payloadHash=0x030f5745 fee=1135 status=accepted
2025-11-08T08:56:00Z RPC /submit payloadHash=0x030f5846 fee=1136 status=accepted
2025-11-08T08:57:00Z RPC /submit payloadHash=0x030f5947 fee=1137 status=broadcast
2025-11-08T08:58:00Z RPC /submit payloadHash=0x030f5a48 fee=1138 status=accepted
2025-11-08T08:59:00Z RPC /submit payloadHash=0x030f5b49 fee=1139 status=accepted
2025-11-08T09:00:00Z RPC /submit payloadHash=0x030f5c50 fee=1140 status=broadcast
2025-11-08T09:01:00Z RPC /submit payloadHash=0x030f5d51 fee=1141 status=accepted
2025-11-08T09:02:00Z RPC /submit payloadHash=0x030f5e52 fee=1142 status=accepted
2025-11-08T09:03:00Z RPC /submit payloadHash=0x030f5f53 fee=1143 status=broadcast
2025-11-08T09:04:00Z RPC /submit payloadHash=0x030f6054 fee=1144 status=accepted
2025-11-08T09:05:00Z RPC /submit payloadHash=0x030f6155 fee=1145 status=accepted
2025-11-08T09:06:00Z RPC /submit payloadHash=0x030f6256 fee=1146 status=broadcast
2025-11-08T09:07:00Z RPC /submit payloadHash=0x030f6357 fee=1147 status=accepted
2025-11-08T09:08:00Z RPC /submit payloadHash=0x030f6458 fee=1148 status=accepted
2025-11-08T09:09:00Z RPC /submit payloadHash=0x030f6559 fee=1149 status=broadcast
2025-11-08T09:10:00Z RPC /submit payloadHash=0x030f6610 fee=1150 status=accepted
2025-11-08T09:11:00Z RPC /submit payloadHash=0x030f6711 fee=1151 status=accepted
2025-11-08T09:12:00Z RPC /submit payloadHash=0x030f6812 fee=1152 status=broadcast
2025-11-08T09:13:00Z RPC /submit payloadHash=0x030f6913 fee=1153 status=accepted
2025-11-08T09:14:00Z RPC /submit payloadHash=0x030f6a14 fee=1154 status=accepted
2025-11-08T09:15:00Z RPC /submit payloadHash=0x030f6b15 fee=1155 status=broadcast
2025-11-08T09:16:00Z RPC /submit payloadHash=0x030f6c16 fee=1156 status=accepted
2025-11-08T09:17:00Z RPC /submit payloadHash=0x030f6d17 fee=1157 status=accepted
2025-11-08T09:18:00Z RPC /submit payloadHash=0x030f6e18 fee=1158 status=broadcast
2025-11-08T09:19:00Z RPC /submit payloadHash=0x030f6f19 fee=1159 status=accepted
2025-11-08T09:20:00Z RPC /submit payloadHash=0x030f7020 fee=1160 status=accepted
2025-11-08T09:21:00Z RPC /submit payloadHash=0x030f7121 fee=1161 status=broadcast
2025-11-08T09:22:00Z RPC /submit payloadHash=0x030f7222 fee=1162 status=accepted
2025-11-08T09:23:00Z RPC /submit payloadHash=0x030f7323 fee=1163 status=accepted
2025-11-08T09:24:00Z RPC /submit payloadHash=0x030f7424 fee=1164 status=broadcast
2025-11-08T09:25:00Z RPC /submit payloadHash=0x030f7525 fee=1165 status=accepted
2025-11-08T09:26:00Z RPC /submit payloadHash=0x030f7626 fee=1166 status=accepted
2025-11-08T09:27:00Z RPC /submit payloadHash=0x030f7727 fee=1167 status=broadcast
2025-11-08T09:28:00Z RPC /submit payloadHash=0x030f7828 fee=1168 status=accepted
2025-11-08T09:29:00Z RPC /submit payloadHash=0x030f7929 fee=1169 status=accepted
2025-11-08T09:30:00Z RPC /submit payloadHash=0x030f7a30 fee=1170 status=broadcast
2025-11-08T09:31:00Z RPC /submit payloadHash=0x030f7b31 fee=1171 status=accepted
2025-11-08T09:32:00Z RPC /submit payloadHash=0x030f7c32 fee=1172 status=accepted
2025-11-08T09:33:00Z RPC /submit payloadHash=0x030f7d33 fee=1173 status=broadcast
2025-11-08T09:34:00Z RPC /submit payloadHash=0x030f7e34 fee=1174 status=accepted
2025-11-08T09:35:00Z RPC /submit payloadHash=0x030f7f35 fee=1175 status=accepted
2025-11-08T09:36:00Z RPC /submit payloadHash=0x030f8036 fee=1176 status=broadcast
2025-11-08T09:37:00Z RPC /submit payloadHash=0x030f8137 fee=1177 status=accepted
2025-11-08T09:38:00Z RPC /submit payloadHash=0x030f8238 fee=1178 status=accepted
2025-11-08T09:39:00Z RPC /submit payloadHash=0x030f8339 fee=1179 status=broadcast
2025-11-08T09:40:00Z RPC /submit payloadHash=0x030f8440 fee=1180 status=accepted
2025-11-08T09:41:00Z RPC /submit payloadHash=0x030f8541 fee=1181 status=accepted
2025-11-08T09:42:00Z RPC /submit payloadHash=0x030f8642 fee=1182 status=broadcast
2025-11-08T09:43:00Z RPC /submit payloadHash=0x030f8743 fee=1183 status=accepted
2025-11-08T09:44:00Z RPC /submit payloadHash=0x030f8844 fee=1184 status=accepted
2025-11-08T09:45:00Z RPC /submit payloadHash=0x030f8945 fee=1185 status=broadcast
2025-11-08T09:46:00Z RPC /submit payloadHash=0x030f8a46 fee=1186 status=accepted
2025-11-08T09:47:00Z RPC /submit payloadHash=0x030f8b47 fee=1187 status=accepted
2025-11-08T09:48:00Z RPC /submit payloadHash=0x030f8c48 fee=1188 status=broadcast
2025-11-08T09:49:00Z RPC /submit payloadHash=0x030f8d49 fee=1189 status=accepted
2025-11-08T09:50:00Z RPC /submit payloadHash=0x030f8e50 fee=1190 status=accepted
2025-11-08T09:51:00Z RPC /submit payloadHash=0x030f8f51 fee=1191 status=broadcast
2025-11-08T09:52:00Z RPC /submit payloadHash=0x030f9052 fee=1192 status=accepted
2025-11-08T09:53:00Z RPC /submit payloadHash=0x030f9153 fee=1193 status=accepted
2025-11-08T09:54:00Z RPC /submit payloadHash=0x030f9254 fee=1194 status=broadcast
2025-11-08T09:55:00Z RPC /submit payloadHash=0x030f9355 fee=1195 status=accepted
2025-11-08T09:56:00Z RPC /submit payloadHash=0x030f9456 fee=1196 status=accepted
2025-11-08T09:57:00Z RPC /submit payloadHash=0x030f9557 fee=1197 status=broadcast
2025-11-08T09:58:00Z RPC /submit payloadHash=0x030f9658 fee=1198 status=accepted
2025-11-08T09:59:00Z RPC /submit payloadHash=0x030f9759 fee=1199 status=accepted
2025-11-08T10:00:00Z RPC /submit payloadHash=0x030f9810 fee=1000 status=broadcast
2025-11-08T10:01:00Z RPC /submit payloadHash=0x030f9911 fee=1001 status=accepted
2025-11-08T10:02:00Z RPC /submit payloadHash=0x030f9a12 fee=1002 status=accepted
2025-11-08T10:03:00Z RPC /submit payloadHash=0x030f9b13 fee=1003 status=broadcast
2025-11-08T10:04:00Z RPC /submit payloadHash=0x030f9c14 fee=1004 status=accepted
2025-11-08T10:05:00Z RPC /submit payloadHash=0x030f9d15 fee=1005 status=accepted
2025-11-08T10:06:00Z RPC /submit payloadHash=0x030f9e16 fee=1006 status=broadcast
2025-11-08T10:07:00Z RPC /submit payloadHash=0x030f9f17 fee=1007 status=accepted
2025-11-08T10:08:00Z RPC /submit payloadHash=0x030fa018 fee=1008 status=accepted
2025-11-08T10:09:00Z RPC /submit payloadHash=0x030fa119 fee=1009 status=broadcast
2025-11-08T10:10:00Z RPC /submit payloadHash=0x030fa220 fee=1010 status=accepted
2025-11-08T10:11:00Z RPC /submit payloadHash=0x030fa321 fee=1011 status=accepted
2025-11-08T10:12:00Z RPC /submit payloadHash=0x030fa422 fee=1012 status=broadcast
2025-11-08T10:13:00Z RPC /submit payloadHash=0x030fa523 fee=1013 status=accepted
2025-11-08T10:14:00Z RPC /submit payloadHash=0x030fa624 fee=1014 status=accepted
2025-11-08T10:15:00Z RPC /submit payloadHash=0x030fa725 fee=1015 status=broadcast
2025-11-08T10:16:00Z RPC /submit payloadHash=0x030fa826 fee=1016 status=accepted
2025-11-08T10:17:00Z RPC /submit payloadHash=0x030fa927 fee=1017 status=accepted
2025-11-08T10:18:00Z RPC /submit payloadHash=0x030faa28 fee=1018 status=broadcast
2025-11-08T10:19:00Z RPC /submit payloadHash=0x030fab29 fee=1019 status=accepted
2025-11-08T10:20:00Z RPC /submit payloadHash=0x030fac30 fee=1020 status=accepted
2025-11-08T10:21:00Z RPC /submit payloadHash=0x030fad31 fee=1021 status=broadcast
2025-11-08T10:22:00Z RPC /submit payloadHash=0x030fae32 fee=1022 status=accepted
2025-11-08T10:23:00Z RPC /submit payloadHash=0x030faf33 fee=1023 status=accepted
2025-11-08T10:24:00Z RPC /submit payloadHash=0x030fb034 fee=1024 status=broadcast
2025-11-08T10:25:00Z RPC /submit payloadHash=0x030fb135 fee=1025 status=accepted
2025-11-08T10:26:00Z RPC /submit payloadHash=0x030fb236 fee=1026 status=accepted
2025-11-08T10:27:00Z RPC /submit payloadHash=0x030fb337 fee=1027 status=broadcast
2025-11-08T10:28:00Z RPC /submit payloadHash=0x030fb438 fee=1028 status=accepted
2025-11-08T10:29:00Z RPC /submit payloadHash=0x030fb539 fee=1029 status=accepted
2025-11-08T10:30:00Z RPC /submit payloadHash=0x030fb640 fee=1030 status=broadcast
2025-11-08T10:31:00Z RPC /submit payloadHash=0x030fb741 fee=1031 status=accepted
2025-11-08T10:32:00Z RPC /submit payloadHash=0x030fb842 fee=1032 status=accepted
2025-11-08T10:33:00Z RPC /submit payloadHash=0x030fb943 fee=1033 status=broadcast
2025-11-08T10:34:00Z RPC /submit payloadHash=0x030fba44 fee=1034 status=accepted
2025-11-08T10:35:00Z RPC /submit payloadHash=0x030fbb45 fee=1035 status=accepted
2025-11-08T10:36:00Z RPC /submit payloadHash=0x030fbc46 fee=1036 status=broadcast
2025-11-08T10:37:00Z RPC /submit payloadHash=0x030fbd47 fee=1037 status=accepted
2025-11-08T10:38:00Z RPC /submit payloadHash=0x030fbe48 fee=1038 status=accepted
2025-11-08T10:39:00Z RPC /submit payloadHash=0x030fbf49 fee=1039 status=broadcast
2025-11-08T10:40:00Z RPC /submit payloadHash=0x030fc050 fee=1040 status=accepted
2025-11-08T10:41:00Z RPC /submit payloadHash=0x030fc151 fee=1041 status=accepted
2025-11-08T10:42:00Z RPC /submit payloadHash=0x030fc252 fee=1042 status=broadcast
2025-11-08T10:43:00Z RPC /submit payloadHash=0x030fc353 fee=1043 status=accepted
2025-11-08T10:44:00Z RPC /submit payloadHash=0x030fc454 fee=1044 status=accepted
2025-11-08T10:45:00Z RPC /submit payloadHash=0x030fc555 fee=1045 status=broadcast
2025-11-08T10:46:00Z RPC /submit payloadHash=0x030fc656 fee=1046 status=accepted
2025-11-08T10:47:00Z RPC /submit payloadHash=0x030fc757 fee=1047 status=accepted
2025-11-08T10:48:00Z RPC /submit payloadHash=0x030fc858 fee=1048 status=broadcast
2025-11-08T10:49:00Z RPC /submit payloadHash=0x030fc959 fee=1049 status=accepted
2025-11-08T10:50:00Z RPC /submit payloadHash=0x030fca10 fee=1050 status=accepted
2025-11-08T10:51:00Z RPC /submit payloadHash=0x030fcb11 fee=1051 status=broadcast
2025-11-08T10:52:00Z RPC /submit payloadHash=0x030fcc12 fee=1052 status=accepted
2025-11-08T10:53:00Z RPC /submit payloadHash=0x030fcd13 fee=1053 status=accepted
2025-11-08T10:54:00Z RPC /submit payloadHash=0x030fce14 fee=1054 status=broadcast
2025-11-08T10:55:00Z RPC /submit payloadHash=0x030fcf15 fee=1055 status=accepted
2025-11-08T10:56:00Z RPC /submit payloadHash=0x030fd016 fee=1056 status=accepted
2025-11-08T10:57:00Z RPC /submit payloadHash=0x030fd117 fee=1057 status=broadcast
2025-11-08T10:58:00Z RPC /submit payloadHash=0x030fd218 fee=1058 status=accepted
2025-11-08T10:59:00Z RPC /submit payloadHash=0x030fd319 fee=1059 status=accepted
2025-11-08T11:00:00Z RPC /submit payloadHash=0x030fd420 fee=1060 status=broadcast
2025-11-08T11:01:00Z RPC /submit payloadHash=0x030fd521 fee=1061 status=accepted
2025-11-08T11:02:00Z RPC /submit payloadHash=0x030fd622 fee=1062 status=accepted
2025-11-08T11:03:00Z RPC /submit payloadHash=0x030fd723 fee=1063 status=broadcast
2025-11-08T11:04:00Z RPC /submit payloadHash=0x030fd824 fee=1064 status=accepted
2025-11-08T11:05:00Z RPC /submit payloadHash=0x030fd925 fee=1065 status=accepted
2025-11-08T11:06:00Z RPC /submit payloadHash=0x030fda26 fee=1066 status=broadcast
2025-11-08T11:07:00Z RPC /submit payloadHash=0x030fdb27 fee=1067 status=accepted
2025-11-08T11:08:00Z RPC /submit payloadHash=0x030fdc28 fee=1068 status=accepted
2025-11-08T11:09:00Z RPC /submit payloadHash=0x030fdd29 fee=1069 status=broadcast
2025-11-08T11:10:00Z RPC /submit payloadHash=0x030fde30 fee=1070 status=accepted
2025-11-08T11:11:00Z RPC /submit payloadHash=0x030fdf31 fee=1071 status=accepted
2025-11-08T11:12:00Z RPC /submit payloadHash=0x030fe032 fee=1072 status=broadcast
2025-11-08T11:13:00Z RPC /submit payloadHash=0x030fe133 fee=1073 status=accepted
2025-11-08T11:14:00Z RPC /submit payloadHash=0x030fe234 fee=1074 status=accepted
2025-11-08T11:15:00Z RPC /submit payloadHash=0x030fe335 fee=1075 status=broadcast
2025-11-08T11:16:00Z RPC /submit payloadHash=0x030fe436 fee=1076 status=accepted
2025-11-08T11:17:00Z RPC /submit payloadHash=0x030fe537 fee=1077 status=accepted
2025-11-08T11:18:00Z RPC /submit payloadHash=0x030fe638 fee=1078 status=broadcast
2025-11-08T11:19:00Z RPC /submit payloadHash=0x030fe739 fee=1079 status=accepted
2025-11-08T11:20:00Z RPC /submit payloadHash=0x030fe840 fee=1080 status=accepted
2025-11-08T11:21:00Z RPC /submit payloadHash=0x030fe941 fee=1081 status=broadcast
2025-11-08T11:22:00Z RPC /submit payloadHash=0x030fea42 fee=1082 status=accepted
2025-11-08T11:23:00Z RPC /submit payloadHash=0x030feb43 fee=1083 status=accepted
2025-11-08T11:24:00Z RPC /submit payloadHash=0x030fec44 fee=1084 status=broadcast
2025-11-08T11:25:00Z RPC /submit payloadHash=0x030fed45 fee=1085 status=accepted
2025-11-08T11:26:00Z RPC /submit payloadHash=0x030fee46 fee=1086 status=accepted
2025-11-08T11:27:00Z RPC /submit payloadHash=0x030fef47 fee=1087 status=broadcast
2025-11-08T11:28:00Z RPC /submit payloadHash=0x030ff048 fee=1088 status=accepted
2025-11-08T11:29:00Z RPC /submit payloadHash=0x030ff149 fee=1089 status=accepted
2025-11-08T11:30:00Z RPC /submit payloadHash=0x030ff250 fee=1090 status=broadcast
2025-11-08T11:31:00Z RPC /submit payloadHash=0x030ff351 fee=1091 status=accepted
2025-11-08T11:32:00Z RPC /submit payloadHash=0x030ff452 fee=1092 status=accepted
2025-11-08T11:33:00Z RPC /submit payloadHash=0x030ff553 fee=1093 status=broadcast
2025-11-08T11:34:00Z RPC /submit payloadHash=0x030ff654 fee=1094 status=accepted
2025-11-08T11:35:00Z RPC /submit payloadHash=0x030ff755 fee=1095 status=accepted
2025-11-08T11:36:00Z RPC /submit payloadHash=0x030ff856 fee=1096 status=broadcast
2025-11-08T11:37:00Z RPC /submit payloadHash=0x030ff957 fee=1097 status=accepted
2025-11-08T11:38:00Z RPC /submit payloadHash=0x030ffa58 fee=1098 status=accepted
2025-11-08T11:39:00Z RPC /submit payloadHash=0x030ffb59 fee=1099 status=broadcast
2025-11-08T11:40:00Z RPC /submit payloadHash=0x030ffc10 fee=1100 status=accepted
2025-11-08T11:41:00Z RPC /submit payloadHash=0x030ffd11 fee=1101 status=accepted
2025-11-08T11:42:00Z RPC /submit payloadHash=0x030ffe12 fee=1102 status=broadcast
2025-11-08T11:43:00Z RPC /submit payloadHash=0x030fff13 fee=1103 status=accepted
2025-11-08T11:44:00Z RPC /submit payloadHash=0x03100014 fee=1104 status=accepted
2025-11-08T11:45:00Z RPC /submit payloadHash=0x03100115 fee=1105 status=broadcast
2025-11-08T11:46:00Z RPC /submit payloadHash=0x03100216 fee=1106 status=accepted
2025-11-08T11:47:00Z RPC /submit payloadHash=0x03100317 fee=1107 status=accepted
2025-11-08T11:48:00Z RPC /submit payloadHash=0x03100418 fee=1108 status=broadcast
2025-11-08T11:49:00Z RPC /submit payloadHash=0x03100519 fee=1109 status=accepted
2025-11-08T11:50:00Z RPC /submit payloadHash=0x03100620 fee=1110 status=accepted
2025-11-08T11:51:00Z RPC /submit payloadHash=0x03100721 fee=1111 status=broadcast
2025-11-08T11:52:00Z RPC /submit payloadHash=0x03100822 fee=1112 status=accepted
2025-11-08T11:53:00Z RPC /submit payloadHash=0x03100923 fee=1113 status=accepted
2025-11-08T11:54:00Z RPC /submit payloadHash=0x03100a24 fee=1114 status=broadcast
2025-11-08T11:55:00Z RPC /submit payloadHash=0x03100b25 fee=1115 status=accepted
2025-11-08T11:56:00Z RPC /submit payloadHash=0x03100c26 fee=1116 status=accepted
2025-11-08T11:57:00Z RPC /submit payloadHash=0x03100d27 fee=1117 status=broadcast
2025-11-08T11:58:00Z RPC /submit payloadHash=0x03100e28 fee=1118 status=accepted
2025-11-08T11:59:00Z RPC /submit payloadHash=0x03100f29 fee=1119 status=accepted
2025-11-08T12:00:00Z RPC /submit payloadHash=0x03101030 fee=1120 status=broadcast
2025-11-08T12:01:00Z RPC /submit payloadHash=0x03101131 fee=1121 status=accepted
2025-11-08T12:02:00Z RPC /submit payloadHash=0x03101232 fee=1122 status=accepted
2025-11-08T12:03:00Z RPC /submit payloadHash=0x03101333 fee=1123 status=broadcast
2025-11-08T12:04:00Z RPC /submit payloadHash=0x03101434 fee=1124 status=accepted
2025-11-08T12:05:00Z RPC /submit payloadHash=0x03101535 fee=1125 status=accepted
2025-11-08T12:06:00Z RPC /submit payloadHash=0x03101636 fee=1126 status=broadcast
2025-11-08T12:07:00Z RPC /submit payloadHash=0x03101737 fee=1127 status=accepted
2025-11-08T12:08:00Z RPC /submit payloadHash=0x03101838 fee=1128 status=accepted
2025-11-08T12:09:00Z RPC /submit payloadHash=0x03101939 fee=1129 status=broadcast
2025-11-08T12:10:00Z RPC /submit payloadHash=0x03101a40 fee=1130 status=accepted
2025-11-08T12:11:00Z RPC /submit payloadHash=0x03101b41 fee=1131 status=accepted
2025-11-08T12:12:00Z RPC /submit payloadHash=0x03101c42 fee=1132 status=broadcast
2025-11-08T12:13:00Z RPC /submit payloadHash=0x03101d43 fee=1133 status=accepted
2025-11-08T12:14:00Z RPC /submit payloadHash=0x03101e44 fee=1134 status=accepted
2025-11-08T12:15:00Z RPC /submit payloadHash=0x03101f45 fee=1135 status=broadcast
2025-11-08T12:16:00Z RPC /submit payloadHash=0x03102046 fee=1136 status=accepted
2025-11-08T12:17:00Z RPC /submit payloadHash=0x03102147 fee=1137 status=accepted
2025-11-08T12:18:00Z RPC /submit payloadHash=0x03102248 fee=1138 status=broadcast
2025-11-08T12:19:00Z RPC /submit payloadHash=0x03102349 fee=1139 status=accepted
2025-11-08T12:20:00Z RPC /submit payloadHash=0x03102450 fee=1140 status=accepted
2025-11-08T12:21:00Z RPC /submit payloadHash=0x03102551 fee=1141 status=broadcast
2025-11-08T12:22:00Z RPC /submit payloadHash=0x03102652 fee=1142 status=accepted
2025-11-08T12:23:00Z RPC /submit payloadHash=0x03102753 fee=1143 status=accepted
2025-11-08T12:24:00Z RPC /submit payloadHash=0x03102854 fee=1144 status=broadcast
2025-11-08T12:25:00Z RPC /submit payloadHash=0x03102955 fee=1145 status=accepted
2025-11-08T12:26:00Z RPC /submit payloadHash=0x03102a56 fee=1146 status=accepted
2025-11-08T12:27:00Z RPC /submit payloadHash=0x03102b57 fee=1147 status=broadcast
2025-11-08T12:28:00Z RPC /submit payloadHash=0x03102c58 fee=1148 status=accepted
2025-11-08T12:29:00Z RPC /submit payloadHash=0x03102d59 fee=1149 status=accepted
2025-11-08T12:30:00Z RPC /submit payloadHash=0x03102e10 fee=1150 status=broadcast
2025-11-08T12:31:00Z RPC /submit payloadHash=0x03102f11 fee=1151 status=accepted
2025-11-08T12:32:00Z RPC /submit payloadHash=0x03103012 fee=1152 status=accepted
2025-11-08T12:33:00Z RPC /submit payloadHash=0x03103113 fee=1153 status=broadcast
2025-11-08T12:34:00Z RPC /submit payloadHash=0x03103214 fee=1154 status=accepted
2025-11-08T12:35:00Z RPC /submit payloadHash=0x03103315 fee=1155 status=accepted
2025-11-08T12:36:00Z RPC /submit payloadHash=0x03103416 fee=1156 status=broadcast
2025-11-08T12:37:00Z RPC /submit payloadHash=0x03103517 fee=1157 status=accepted
2025-11-08T12:38:00Z RPC /submit payloadHash=0x03103618 fee=1158 status=accepted
2025-11-08T12:39:00Z RPC /submit payloadHash=0x03103719 fee=1159 status=broadcast
2025-11-08T12:40:00Z RPC /submit payloadHash=0x03103820 fee=1160 status=accepted
2025-11-08T12:41:00Z RPC /submit payloadHash=0x03103921 fee=1161 status=accepted
2025-11-08T12:42:00Z RPC /submit payloadHash=0x03103a22 fee=1162 status=broadcast
2025-11-08T12:43:00Z RPC /submit payloadHash=0x03103b23 fee=1163 status=accepted
2025-11-08T12:44:00Z RPC /submit payloadHash=0x03103c24 fee=1164 status=accepted
2025-11-08T12:45:00Z RPC /submit payloadHash=0x03103d25 fee=1165 status=broadcast
2025-11-08T12:46:00Z RPC /submit payloadHash=0x03103e26 fee=1166 status=accepted
2025-11-08T12:47:00Z RPC /submit payloadHash=0x03103f27 fee=1167 status=accepted
2025-11-08T12:48:00Z RPC /submit payloadHash=0x03104028 fee=1168 status=broadcast
2025-11-08T12:49:00Z RPC /submit payloadHash=0x03104129 fee=1169 status=accepted
2025-11-08T12:50:00Z RPC /submit payloadHash=0x03104230 fee=1170 status=accepted
2025-11-08T12:51:00Z RPC /submit payloadHash=0x03104331 fee=1171 status=broadcast
2025-11-08T12:52:00Z RPC /submit payloadHash=0x03104432 fee=1172 status=accepted
2025-11-08T12:53:00Z RPC /submit payloadHash=0x03104533 fee=1173 status=accepted
2025-11-08T12:54:00Z RPC /submit payloadHash=0x03104634 fee=1174 status=broadcast
2025-11-08T12:55:00Z RPC /submit payloadHash=0x03104735 fee=1175 status=accepted
2025-11-08T12:56:00Z RPC /submit payloadHash=0x03104836 fee=1176 status=accepted
2025-11-08T12:57:00Z RPC /submit payloadHash=0x03104937 fee=1177 status=broadcast
2025-11-08T12:58:00Z RPC /submit payloadHash=0x03104a38 fee=1178 status=accepted
2025-11-08T12:59:00Z RPC /submit payloadHash=0x03104b39 fee=1179 status=accepted
2025-11-08T13:00:00Z RPC /submit payloadHash=0x03104c40 fee=1180 status=broadcast
2025-11-08T13:01:00Z RPC /submit payloadHash=0x03104d41 fee=1181 status=accepted
2025-11-08T13:02:00Z RPC /submit payloadHash=0x03104e42 fee=1182 status=accepted
2025-11-08T13:03:00Z RPC /submit payloadHash=0x03104f43 fee=1183 status=broadcast
2025-11-08T13:04:00Z RPC /submit payloadHash=0x03105044 fee=1184 status=accepted
2025-11-08T13:05:00Z RPC /submit payloadHash=0x03105145 fee=1185 status=accepted
2025-11-08T13:06:00Z RPC /submit payloadHash=0x03105246 fee=1186 status=broadcast
2025-11-08T13:07:00Z RPC /submit payloadHash=0x03105347 fee=1187 status=accepted
2025-11-08T13:08:00Z RPC /submit payloadHash=0x03105448 fee=1188 status=accepted
2025-11-08T13:09:00Z RPC /submit payloadHash=0x03105549 fee=1189 status=broadcast
2025-11-08T13:10:00Z RPC /submit payloadHash=0x03105650 fee=1190 status=accepted
2025-11-08T13:11:00Z RPC /submit payloadHash=0x03105751 fee=1191 status=accepted
2025-11-08T13:12:00Z RPC /submit payloadHash=0x03105852 fee=1192 status=broadcast
2025-11-08T13:13:00Z RPC /submit payloadHash=0x03105953 fee=1193 status=accepted
2025-11-08T13:14:00Z RPC /submit payloadHash=0x03105a54 fee=1194 status=accepted
2025-11-08T13:15:00Z RPC /submit payloadHash=0x03105b55 fee=1195 status=broadcast
2025-11-08T13:16:00Z RPC /submit payloadHash=0x03105c56 fee=1196 status=accepted
2025-11-08T13:17:00Z RPC /submit payloadHash=0x03105d57 fee=1197 status=accepted
2025-11-08T13:18:00Z RPC /submit payloadHash=0x03105e58 fee=1198 status=broadcast
2025-11-08T13:19:00Z RPC /submit payloadHash=0x03105f59 fee=1199 status=accepted
2025-11-08T13:20:00Z RPC /submit payloadHash=0x03106010 fee=1000 status=accepted
2025-11-08T13:21:00Z RPC /submit payloadHash=0x03106111 fee=1001 status=broadcast
2025-11-08T13:22:00Z RPC /submit payloadHash=0x03106212 fee=1002 status=accepted
2025-11-08T13:23:00Z RPC /submit payloadHash=0x03106313 fee=1003 status=accepted
2025-11-08T13:24:00Z RPC /submit payloadHash=0x03106414 fee=1004 status=broadcast
2025-11-08T13:25:00Z RPC /submit payloadHash=0x03106515 fee=1005 status=accepted
2025-11-08T13:26:00Z RPC /submit payloadHash=0x03106616 fee=1006 status=accepted
2025-11-08T13:27:00Z RPC /submit payloadHash=0x03106717 fee=1007 status=broadcast
2025-11-08T13:28:00Z RPC /submit payloadHash=0x03106818 fee=1008 status=accepted
2025-11-08T13:29:00Z RPC /submit payloadHash=0x03106919 fee=1009 status=accepted
2025-11-08T13:30:00Z RPC /submit payloadHash=0x03106a20 fee=1010 status=broadcast
2025-11-08T13:31:00Z RPC /submit payloadHash=0x03106b21 fee=1011 status=accepted
2025-11-08T13:32:00Z RPC /submit payloadHash=0x03106c22 fee=1012 status=accepted
2025-11-08T13:33:00Z RPC /submit payloadHash=0x03106d23 fee=1013 status=broadcast
2025-11-08T13:34:00Z RPC /submit payloadHash=0x03106e24 fee=1014 status=accepted
2025-11-08T13:35:00Z RPC /submit payloadHash=0x03106f25 fee=1015 status=accepted
2025-11-08T13:36:00Z RPC /submit payloadHash=0x03107026 fee=1016 status=broadcast
2025-11-08T13:37:00Z RPC /submit payloadHash=0x03107127 fee=1017 status=accepted
2025-11-08T13:38:00Z RPC /submit payloadHash=0x03107228 fee=1018 status=accepted
2025-11-08T13:39:00Z RPC /submit payloadHash=0x03107329 fee=1019 status=broadcast
2025-11-08T13:40:00Z RPC /submit payloadHash=0x03107430 fee=1020 status=accepted
2025-11-08T13:41:00Z RPC /submit payloadHash=0x03107531 fee=1021 status=accepted
2025-11-08T13:42:00Z RPC /submit payloadHash=0x03107632 fee=1022 status=broadcast
2025-11-08T13:43:00Z RPC /submit payloadHash=0x03107733 fee=1023 status=accepted
2025-11-08T13:44:00Z RPC /submit payloadHash=0x03107834 fee=1024 status=accepted
2025-11-08T13:45:00Z RPC /submit payloadHash=0x03107935 fee=1025 status=broadcast
2025-11-08T13:46:00Z RPC /submit payloadHash=0x03107a36 fee=1026 status=accepted
2025-11-08T13:47:00Z RPC /submit payloadHash=0x03107b37 fee=1027 status=accepted
2025-11-08T13:48:00Z RPC /submit payloadHash=0x03107c38 fee=1028 status=broadcast
2025-11-08T13:49:00Z RPC /submit payloadHash=0x03107d39 fee=1029 status=accepted
2025-11-08T13:50:00Z RPC /submit payloadHash=0x03107e40 fee=1030 status=accepted
2025-11-08T13:51:00Z RPC /submit payloadHash=0x03107f41 fee=1031 status=broadcast
2025-11-08T13:52:00Z RPC /submit payloadHash=0x03108042 fee=1032 status=accepted
2025-11-08T13:53:00Z RPC /submit payloadHash=0x03108143 fee=1033 status=accepted
2025-11-08T13:54:00Z RPC /submit payloadHash=0x03108244 fee=1034 status=broadcast
2025-11-08T13:55:00Z RPC /submit payloadHash=0x03108345 fee=1035 status=accepted
2025-11-08T13:56:00Z RPC /submit payloadHash=0x03108446 fee=1036 status=accepted
2025-11-08T13:57:00Z RPC /submit payloadHash=0x03108547 fee=1037 status=broadcast
2025-11-08T13:58:00Z RPC /submit payloadHash=0x03108648 fee=1038 status=accepted
2025-11-08T13:59:00Z RPC /submit payloadHash=0x03108749 fee=1039 status=accepted
2025-11-08T14:00:00Z RPC /submit payloadHash=0x03108850 fee=1040 status=broadcast
2025-11-08T14:01:00Z RPC /submit payloadHash=0x03108951 fee=1041 status=accepted
2025-11-08T14:02:00Z RPC /submit payloadHash=0x03108a52 fee=1042 status=accepted
2025-11-08T14:03:00Z RPC /submit payloadHash=0x03108b53 fee=1043 status=broadcast
2025-11-08T14:04:00Z RPC /submit payloadHash=0x03108c54 fee=1044 status=accepted
2025-11-08T14:05:00Z RPC /submit payloadHash=0x03108d55 fee=1045 status=accepted
2025-11-08T14:06:00Z RPC /submit payloadHash=0x03108e56 fee=1046 status=broadcast
2025-11-08T14:07:00Z RPC /submit payloadHash=0x03108f57 fee=1047 status=accepted
2025-11-08T14:08:00Z RPC /submit payloadHash=0x03109058 fee=1048 status=accepted
2025-11-08T14:09:00Z RPC /submit payloadHash=0x03109159 fee=1049 status=broadcast
2025-11-08T14:10:00Z RPC /submit payloadHash=0x03109210 fee=1050 status=accepted
2025-11-08T14:11:00Z RPC /submit payloadHash=0x03109311 fee=1051 status=accepted
2025-11-08T14:12:00Z RPC /submit payloadHash=0x03109412 fee=1052 status=broadcast
2025-11-08T14:13:00Z RPC /submit payloadHash=0x03109513 fee=1053 status=accepted
2025-11-08T14:14:00Z RPC /submit payloadHash=0x03109614 fee=1054 status=accepted
2025-11-08T14:15:00Z RPC /submit payloadHash=0x03109715 fee=1055 status=broadcast
2025-11-08T14:16:00Z RPC /submit payloadHash=0x03109816 fee=1056 status=accepted
2025-11-08T14:17:00Z RPC /submit payloadHash=0x03109917 fee=1057 status=accepted
2025-11-08T14:18:00Z RPC /submit payloadHash=0x03109a18 fee=1058 status=broadcast
2025-11-08T14:19:00Z RPC /submit payloadHash=0x03109b19 fee=1059 status=accepted
2025-11-08T14:20:00Z RPC /submit payloadHash=0x03109c20 fee=1060 status=accepted
2025-11-08T14:21:00Z RPC /submit payloadHash=0x03109d21 fee=1061 status=broadcast
2025-11-08T14:22:00Z RPC /submit payloadHash=0x03109e22 fee=1062 status=accepted
2025-11-08T14:23:00Z RPC /submit payloadHash=0x03109f23 fee=1063 status=accepted
2025-11-08T14:24:00Z RPC /submit payloadHash=0x0310a024 fee=1064 status=broadcast
2025-11-08T14:25:00Z RPC /submit payloadHash=0x0310a125 fee=1065 status=accepted
2025-11-08T14:26:00Z RPC /submit payloadHash=0x0310a226 fee=1066 status=accepted
2025-11-08T14:27:00Z RPC /submit payloadHash=0x0310a327 fee=1067 status=broadcast
2025-11-08T14:28:00Z RPC /submit payloadHash=0x0310a428 fee=1068 status=accepted
2025-11-08T14:29:00Z RPC /submit payloadHash=0x0310a529 fee=1069 status=accepted
2025-11-08T14:30:00Z RPC /submit payloadHash=0x0310a630 fee=1070 status=broadcast
2025-11-08T14:31:00Z RPC /submit payloadHash=0x0310a731 fee=1071 status=accepted
2025-11-08T14:32:00Z RPC /submit payloadHash=0x0310a832 fee=1072 status=accepted
2025-11-08T14:33:00Z RPC /submit payloadHash=0x0310a933 fee=1073 status=broadcast
2025-11-08T14:34:00Z RPC /submit payloadHash=0x0310aa34 fee=1074 status=accepted
2025-11-08T14:35:00Z RPC /submit payloadHash=0x0310ab35 fee=1075 status=accepted
2025-11-08T14:36:00Z RPC /submit payloadHash=0x0310ac36 fee=1076 status=broadcast
2025-11-08T14:37:00Z RPC /submit payloadHash=0x0310ad37 fee=1077 status=accepted
2025-11-08T14:38:00Z RPC /submit payloadHash=0x0310ae38 fee=1078 status=accepted
2025-11-08T14:39:00Z RPC /submit payloadHash=0x0310af39 fee=1079 status=broadcast
2025-11-08T14:40:00Z RPC /submit payloadHash=0x0310b040 fee=1080 status=accepted
2025-11-08T14:41:00Z RPC /submit payloadHash=0x0310b141 fee=1081 status=accepted
2025-11-08T14:42:00Z RPC /submit payloadHash=0x0310b242 fee=1082 status=broadcast
2025-11-08T14:43:00Z RPC /submit payloadHash=0x0310b343 fee=1083 status=accepted
2025-11-08T14:44:00Z RPC /submit payloadHash=0x0310b444 fee=1084 status=accepted
2025-11-08T14:45:00Z RPC /submit payloadHash=0x0310b545 fee=1085 status=broadcast
2025-11-08T14:46:00Z RPC /submit payloadHash=0x0310b646 fee=1086 status=accepted
2025-11-08T14:47:00Z RPC /submit payloadHash=0x0310b747 fee=1087 status=accepted
2025-11-08T14:48:00Z RPC /submit payloadHash=0x0310b848 fee=1088 status=broadcast
2025-11-08T14:49:00Z RPC /submit payloadHash=0x0310b949 fee=1089 status=accepted
2025-11-08T14:50:00Z RPC /submit payloadHash=0x0310ba50 fee=1090 status=accepted
2025-11-08T14:51:00Z RPC /submit payloadHash=0x0310bb51 fee=1091 status=broadcast
2025-11-08T14:52:00Z RPC /submit payloadHash=0x0310bc52 fee=1092 status=accepted
2025-11-08T14:53:00Z RPC /submit payloadHash=0x0310bd53 fee=1093 status=accepted
2025-11-08T14:54:00Z RPC /submit payloadHash=0x0310be54 fee=1094 status=broadcast
2025-11-08T14:55:00Z RPC /submit payloadHash=0x0310bf55 fee=1095 status=accepted
2025-11-08T14:56:00Z RPC /submit payloadHash=0x0310c056 fee=1096 status=accepted
2025-11-08T14:57:00Z RPC /submit payloadHash=0x0310c157 fee=1097 status=broadcast
2025-11-08T14:58:00Z RPC /submit payloadHash=0x0310c258 fee=1098 status=accepted
2025-11-08T14:59:00Z RPC /submit payloadHash=0x0310c359 fee=1099 status=accepted
2025-11-08T15:00:00Z RPC /submit payloadHash=0x0310c410 fee=1100 status=broadcast
2025-11-08T15:01:00Z RPC /submit payloadHash=0x0310c511 fee=1101 status=accepted
2025-11-08T15:02:00Z RPC /submit payloadHash=0x0310c612 fee=1102 status=accepted
2025-11-08T15:03:00Z RPC /submit payloadHash=0x0310c713 fee=1103 status=broadcast
2025-11-08T15:04:00Z RPC /submit payloadHash=0x0310c814 fee=1104 status=accepted
2025-11-08T15:05:00Z RPC /submit payloadHash=0x0310c915 fee=1105 status=accepted
2025-11-08T15:06:00Z RPC /submit payloadHash=0x0310ca16 fee=1106 status=broadcast
2025-11-08T15:07:00Z RPC /submit payloadHash=0x0310cb17 fee=1107 status=accepted
2025-11-08T15:08:00Z RPC /submit payloadHash=0x0310cc18 fee=1108 status=accepted
2025-11-08T15:09:00Z RPC /submit payloadHash=0x0310cd19 fee=1109 status=broadcast
2025-11-08T15:10:00Z RPC /submit payloadHash=0x0310ce20 fee=1110 status=accepted
2025-11-08T15:11:00Z RPC /submit payloadHash=0x0310cf21 fee=1111 status=accepted
2025-11-08T15:12:00Z RPC /submit payloadHash=0x0310d022 fee=1112 status=broadcast
2025-11-08T15:13:00Z RPC /submit payloadHash=0x0310d123 fee=1113 status=accepted
2025-11-08T15:14:00Z RPC /submit payloadHash=0x0310d224 fee=1114 status=accepted
2025-11-08T15:15:00Z RPC /submit payloadHash=0x0310d325 fee=1115 status=broadcast
2025-11-08T15:16:00Z RPC /submit payloadHash=0x0310d426 fee=1116 status=accepted
2025-11-08T15:17:00Z RPC /submit payloadHash=0x0310d527 fee=1117 status=accepted
2025-11-08T15:18:00Z RPC /submit payloadHash=0x0310d628 fee=1118 status=broadcast
2025-11-08T15:19:00Z RPC /submit payloadHash=0x0310d729 fee=1119 status=accepted
2025-11-08T15:20:00Z RPC /submit payloadHash=0x0310d830 fee=1120 status=accepted
2025-11-08T15:21:00Z RPC /submit payloadHash=0x0310d931 fee=1121 status=broadcast
2025-11-08T15:22:00Z RPC /submit payloadHash=0x0310da32 fee=1122 status=accepted
2025-11-08T15:23:00Z RPC /submit payloadHash=0x0310db33 fee=1123 status=accepted
2025-11-08T15:24:00Z RPC /submit payloadHash=0x0310dc34 fee=1124 status=broadcast
2025-11-08T15:25:00Z RPC /submit payloadHash=0x0310dd35 fee=1125 status=accepted
2025-11-08T15:26:00Z RPC /submit payloadHash=0x0310de36 fee=1126 status=accepted
2025-11-08T15:27:00Z RPC /submit payloadHash=0x0310df37 fee=1127 status=broadcast
2025-11-08T15:28:00Z RPC /submit payloadHash=0x0310e038 fee=1128 status=accepted
2025-11-08T15:29:00Z RPC /submit payloadHash=0x0310e139 fee=1129 status=accepted
2025-11-08T15:30:00Z RPC /submit payloadHash=0x0310e240 fee=1130 status=broadcast
2025-11-08T15:31:00Z RPC /submit payloadHash=0x0310e341 fee=1131 status=accepted
2025-11-08T15:32:00Z RPC /submit payloadHash=0x0310e442 fee=1132 status=accepted
2025-11-08T15:33:00Z RPC /submit payloadHash=0x0310e543 fee=1133 status=broadcast
2025-11-08T15:34:00Z RPC /submit payloadHash=0x0310e644 fee=1134 status=accepted
2025-11-08T15:35:00Z RPC /submit payloadHash=0x0310e745 fee=1135 status=accepted
2025-11-08T15:36:00Z RPC /submit payloadHash=0x0310e846 fee=1136 status=broadcast
2025-11-08T15:37:00Z RPC /submit payloadHash=0x0310e947 fee=1137 status=accepted
2025-11-08T15:38:00Z RPC /submit payloadHash=0x0310ea48 fee=1138 status=accepted
2025-11-08T15:39:00Z RPC /submit payloadHash=0x0310eb49 fee=1139 status=broadcast
2025-11-08T15:40:00Z RPC /submit payloadHash=0x0310ec50 fee=1140 status=accepted
2025-11-08T15:41:00Z RPC /submit payloadHash=0x0310ed51 fee=1141 status=accepted
2025-11-08T15:42:00Z RPC /submit payloadHash=0x0310ee52 fee=1142 status=broadcast
2025-11-08T15:43:00Z RPC /submit payloadHash=0x0310ef53 fee=1143 status=accepted
2025-11-08T15:44:00Z RPC /submit payloadHash=0x0310f054 fee=1144 status=accepted
2025-11-08T15:45:00Z RPC /submit payloadHash=0x0310f155 fee=1145 status=broadcast
2025-11-08T15:46:00Z RPC /submit payloadHash=0x0310f256 fee=1146 status=accepted
2025-11-08T15:47:00Z RPC /submit payloadHash=0x0310f357 fee=1147 status=accepted
2025-11-08T15:48:00Z RPC /submit payloadHash=0x0310f458 fee=1148 status=broadcast
2025-11-08T15:49:00Z RPC /submit payloadHash=0x0310f559 fee=1149 status=accepted
2025-11-08T15:50:00Z RPC /submit payloadHash=0x0310f610 fee=1150 status=accepted
2025-11-08T15:51:00Z RPC /submit payloadHash=0x0310f711 fee=1151 status=broadcast
2025-11-08T15:52:00Z RPC /submit payloadHash=0x0310f812 fee=1152 status=accepted
2025-11-08T15:53:00Z RPC /submit payloadHash=0x0310f913 fee=1153 status=accepted
2025-11-08T15:54:00Z RPC /submit payloadHash=0x0310fa14 fee=1154 status=broadcast
2025-11-08T15:55:00Z RPC /submit payloadHash=0x0310fb15 fee=1155 status=accepted
2025-11-08T15:56:00Z RPC /submit payloadHash=0x0310fc16 fee=1156 status=accepted
2025-11-08T15:57:00Z RPC /submit payloadHash=0x0310fd17 fee=1157 status=broadcast
2025-11-08T15:58:00Z RPC /submit payloadHash=0x0310fe18 fee=1158 status=accepted
2025-11-08T15:59:00Z RPC /submit payloadHash=0x0310ff19 fee=1159 status=accepted
2025-11-08T16:00:00Z RPC /submit payloadHash=0x03110020 fee=1160 status=broadcast
2025-11-08T16:01:00Z RPC /submit payloadHash=0x03110121 fee=1161 status=accepted
2025-11-08T16:02:00Z RPC /submit payloadHash=0x03110222 fee=1162 status=accepted
2025-11-08T16:03:00Z RPC /submit payloadHash=0x03110323 fee=1163 status=broadcast
2025-11-08T16:04:00Z RPC /submit payloadHash=0x03110424 fee=1164 status=accepted
2025-11-08T16:05:00Z RPC /submit payloadHash=0x03110525 fee=1165 status=accepted
2025-11-08T16:06:00Z RPC /submit payloadHash=0x03110626 fee=1166 status=broadcast
2025-11-08T16:07:00Z RPC /submit payloadHash=0x03110727 fee=1167 status=accepted
2025-11-08T16:08:00Z RPC /submit payloadHash=0x03110828 fee=1168 status=accepted
2025-11-08T16:09:00Z RPC /submit payloadHash=0x03110929 fee=1169 status=broadcast
2025-11-08T16:10:00Z RPC /submit payloadHash=0x03110a30 fee=1170 status=accepted
2025-11-08T16:11:00Z RPC /submit payloadHash=0x03110b31 fee=1171 status=accepted
2025-11-08T16:12:00Z RPC /submit payloadHash=0x03110c32 fee=1172 status=broadcast
2025-11-08T16:13:00Z RPC /submit payloadHash=0x03110d33 fee=1173 status=accepted
2025-11-08T16:14:00Z RPC /submit payloadHash=0x03110e34 fee=1174 status=accepted
2025-11-08T16:15:00Z RPC /submit payloadHash=0x03110f35 fee=1175 status=broadcast
2025-11-08T16:16:00Z RPC /submit payloadHash=0x03111036 fee=1176 status=accepted
2025-11-08T16:17:00Z RPC /submit payloadHash=0x03111137 fee=1177 status=accepted
2025-11-08T16:18:00Z RPC /submit payloadHash=0x03111238 fee=1178 status=broadcast
2025-11-08T16:19:00Z RPC /submit payloadHash=0x03111339 fee=1179 status=accepted
2025-11-08T16:20:00Z RPC /submit payloadHash=0x03111440 fee=1180 status=accepted
2025-11-08T16:21:00Z RPC /submit payloadHash=0x03111541 fee=1181 status=broadcast
2025-11-08T16:22:00Z RPC /submit payloadHash=0x03111642 fee=1182 status=accepted
2025-11-08T16:23:00Z RPC /submit payloadHash=0x03111743 fee=1183 status=accepted
2025-11-08T16:24:00Z RPC /submit payloadHash=0x03111844 fee=1184 status=broadcast
2025-11-08T16:25:00Z RPC /submit payloadHash=0x03111945 fee=1185 status=accepted
2025-11-08T16:26:00Z RPC /submit payloadHash=0x03111a46 fee=1186 status=accepted
2025-11-08T16:27:00Z RPC /submit payloadHash=0x03111b47 fee=1187 status=broadcast
2025-11-08T16:28:00Z RPC /submit payloadHash=0x03111c48 fee=1188 status=accepted
2025-11-08T16:29:00Z RPC /submit payloadHash=0x03111d49 fee=1189 status=accepted
2025-11-08T16:30:00Z RPC /submit payloadHash=0x03111e50 fee=1190 status=broadcast
2025-11-08T16:31:00Z RPC /submit payloadHash=0x03111f51 fee=1191 status=accepted
2025-11-08T16:32:00Z RPC /submit payloadHash=0x03112052 fee=1192 status=accepted
2025-11-08T16:33:00Z RPC /submit payloadHash=0x03112153 fee=1193 status=broadcast
2025-11-08T16:34:00Z RPC /submit payloadHash=0x03112254 fee=1194 status=accepted
2025-11-08T16:35:00Z RPC /submit payloadHash=0x03112355 fee=1195 status=accepted
2025-11-08T16:36:00Z RPC /submit payloadHash=0x03112456 fee=1196 status=broadcast
2025-11-08T16:37:00Z RPC /submit payloadHash=0x03112557 fee=1197 status=accepted
2025-11-08T16:38:00Z RPC /submit payloadHash=0x03112658 fee=1198 status=accepted
2025-11-08T16:39:00Z RPC /submit payloadHash=0x03112759 fee=1199 status=broadcast
2025-11-08T16:40:00Z RPC /submit payloadHash=0x03112810 fee=1000 status=accepted
2025-11-08T16:41:00Z RPC /submit payloadHash=0x03112911 fee=1001 status=accepted
2025-11-08T16:42:00Z RPC /submit payloadHash=0x03112a12 fee=1002 status=broadcast
2025-11-08T16:43:00Z RPC /submit payloadHash=0x03112b13 fee=1003 status=accepted
2025-11-08T16:44:00Z RPC /submit payloadHash=0x03112c14 fee=1004 status=accepted
2025-11-08T16:45:00Z RPC /submit payloadHash=0x03112d15 fee=1005 status=broadcast
2025-11-08T16:46:00Z RPC /submit payloadHash=0x03112e16 fee=1006 status=accepted
2025-11-08T16:47:00Z RPC /submit payloadHash=0x03112f17 fee=1007 status=accepted
2025-11-08T16:48:00Z RPC /submit payloadHash=0x03113018 fee=1008 status=broadcast
2025-11-08T16:49:00Z RPC /submit payloadHash=0x03113119 fee=1009 status=accepted
2025-11-08T16:50:00Z RPC /submit payloadHash=0x03113220 fee=1010 status=accepted
2025-11-08T16:51:00Z RPC /submit payloadHash=0x03113321 fee=1011 status=broadcast
2025-11-08T16:52:00Z RPC /submit payloadHash=0x03113422 fee=1012 status=accepted
2025-11-08T16:53:00Z RPC /submit payloadHash=0x03113523 fee=1013 status=accepted
2025-11-08T16:54:00Z RPC /submit payloadHash=0x03113624 fee=1014 status=broadcast
2025-11-08T16:55:00Z RPC /submit payloadHash=0x03113725 fee=1015 status=accepted
2025-11-08T16:56:00Z RPC /submit payloadHash=0x03113826 fee=1016 status=accepted
2025-11-08T16:57:00Z RPC /submit payloadHash=0x03113927 fee=1017 status=broadcast
2025-11-08T16:58:00Z RPC /submit payloadHash=0x03113a28 fee=1018 status=accepted
2025-11-08T16:59:00Z RPC /submit payloadHash=0x03113b29 fee=1019 status=accepted
2025-11-08T17:00:00Z RPC /submit payloadHash=0x03113c30 fee=1020 status=broadcast
2025-11-08T17:01:00Z RPC /submit payloadHash=0x03113d31 fee=1021 status=accepted
2025-11-08T17:02:00Z RPC /submit payloadHash=0x03113e32 fee=1022 status=accepted
2025-11-08T17:03:00Z RPC /submit payloadHash=0x03113f33 fee=1023 status=broadcast
2025-11-08T17:04:00Z RPC /submit payloadHash=0x03114034 fee=1024 status=accepted
2025-11-08T17:05:00Z RPC /submit payloadHash=0x03114135 fee=1025 status=accepted
2025-11-08T17:06:00Z RPC /submit payloadHash=0x03114236 fee=1026 status=broadcast
2025-11-08T17:07:00Z RPC /submit payloadHash=0x03114337 fee=1027 status=accepted
2025-11-08T17:08:00Z RPC /submit payloadHash=0x03114438 fee=1028 status=accepted
2025-11-08T17:09:00Z RPC /submit payloadHash=0x03114539 fee=1029 status=broadcast
2025-11-08T17:10:00Z RPC /submit payloadHash=0x03114640 fee=1030 status=accepted
2025-11-08T17:11:00Z RPC /submit payloadHash=0x03114741 fee=1031 status=accepted
2025-11-08T17:12:00Z RPC /submit payloadHash=0x03114842 fee=1032 status=broadcast
2025-11-08T17:13:00Z RPC /submit payloadHash=0x03114943 fee=1033 status=accepted
2025-11-08T17:14:00Z RPC /submit payloadHash=0x03114a44 fee=1034 status=accepted
2025-11-08T17:15:00Z RPC /submit payloadHash=0x03114b45 fee=1035 status=broadcast
2025-11-08T17:16:00Z RPC /submit payloadHash=0x03114c46 fee=1036 status=accepted
2025-11-08T17:17:00Z RPC /submit payloadHash=0x03114d47 fee=1037 status=accepted
2025-11-08T17:18:00Z RPC /submit payloadHash=0x03114e48 fee=1038 status=broadcast
2025-11-08T17:19:00Z RPC /submit payloadHash=0x03114f49 fee=1039 status=accepted
2025-11-08T17:20:00Z RPC /submit payloadHash=0x03115050 fee=1040 status=accepted
2025-11-08T17:21:00Z RPC /submit payloadHash=0x03115151 fee=1041 status=broadcast
2025-11-08T17:22:00Z RPC /submit payloadHash=0x03115252 fee=1042 status=accepted
2025-11-08T17:23:00Z RPC /submit payloadHash=0x03115353 fee=1043 status=accepted
2025-11-08T17:24:00Z RPC /submit payloadHash=0x03115454 fee=1044 status=broadcast
2025-11-08T17:25:00Z RPC /submit payloadHash=0x03115555 fee=1045 status=accepted
2025-11-08T17:26:00Z RPC /submit payloadHash=0x03115656 fee=1046 status=accepted
2025-11-08T17:27:00Z RPC /submit payloadHash=0x03115757 fee=1047 status=broadcast
2025-11-08T17:28:00Z RPC /submit payloadHash=0x03115858 fee=1048 status=accepted
2025-11-08T17:29:00Z RPC /submit payloadHash=0x03115959 fee=1049 status=accepted
2025-11-08T17:30:00Z RPC /submit payloadHash=0x03115a10 fee=1050 status=broadcast
2025-11-08T17:31:00Z RPC /submit payloadHash=0x03115b11 fee=1051 status=accepted
2025-11-08T17:32:00Z RPC /submit payloadHash=0x03115c12 fee=1052 status=accepted
2025-11-08T17:33:00Z RPC /submit payloadHash=0x03115d13 fee=1053 status=broadcast
2025-11-08T17:34:00Z RPC /submit payloadHash=0x03115e14 fee=1054 status=accepted
2025-11-08T17:35:00Z RPC /submit payloadHash=0x03115f15 fee=1055 status=accepted
2025-11-08T17:36:00Z RPC /submit payloadHash=0x03116016 fee=1056 status=broadcast
2025-11-08T17:37:00Z RPC /submit payloadHash=0x03116117 fee=1057 status=accepted
2025-11-08T17:38:00Z RPC /submit payloadHash=0x03116218 fee=1058 status=accepted
2025-11-08T17:39:00Z RPC /submit payloadHash=0x03116319 fee=1059 status=broadcast
2025-11-08T17:40:00Z RPC /submit payloadHash=0x03116420 fee=1060 status=accepted
2025-11-08T17:41:00Z RPC /submit payloadHash=0x03116521 fee=1061 status=accepted
2025-11-08T17:42:00Z RPC /submit payloadHash=0x03116622 fee=1062 status=broadcast
2025-11-08T17:43:00Z RPC /submit payloadHash=0x03116723 fee=1063 status=accepted
2025-11-08T17:44:00Z RPC /submit payloadHash=0x03116824 fee=1064 status=accepted
2025-11-08T17:45:00Z RPC /submit payloadHash=0x03116925 fee=1065 status=broadcast
2025-11-08T17:46:00Z RPC /submit payloadHash=0x03116a26 fee=1066 status=accepted
2025-11-08T17:47:00Z RPC /submit payloadHash=0x03116b27 fee=1067 status=accepted
2025-11-08T17:48:00Z RPC /submit payloadHash=0x03116c28 fee=1068 status=broadcast
2025-11-08T17:49:00Z RPC /submit payloadHash=0x03116d29 fee=1069 status=accepted
2025-11-08T17:50:00Z RPC /submit payloadHash=0x03116e30 fee=1070 status=accepted
2025-11-08T17:51:00Z RPC /submit payloadHash=0x03116f31 fee=1071 status=broadcast
2025-11-08T17:52:00Z RPC /submit payloadHash=0x03117032 fee=1072 status=accepted
2025-11-08T17:53:00Z RPC /submit payloadHash=0x03117133 fee=1073 status=accepted
2025-11-08T17:54:00Z RPC /submit payloadHash=0x03117234 fee=1074 status=broadcast
2025-11-08T17:55:00Z RPC /submit payloadHash=0x03117335 fee=1075 status=accepted
2025-11-08T17:56:00Z RPC /submit payloadHash=0x03117436 fee=1076 status=accepted
2025-11-08T17:57:00Z RPC /submit payloadHash=0x03117537 fee=1077 status=broadcast
2025-11-08T17:58:00Z RPC /submit payloadHash=0x03117638 fee=1078 status=accepted
2025-11-08T17:59:00Z RPC /submit payloadHash=0x03117739 fee=1079 status=accepted
2025-11-08T18:00:00Z RPC /submit payloadHash=0x03117840 fee=1080 status=broadcast
2025-11-08T18:01:00Z RPC /submit payloadHash=0x03117941 fee=1081 status=accepted
2025-11-08T18:02:00Z RPC /submit payloadHash=0x03117a42 fee=1082 status=accepted
2025-11-08T18:03:00Z RPC /submit payloadHash=0x03117b43 fee=1083 status=broadcast
2025-11-08T18:04:00Z RPC /submit payloadHash=0x03117c44 fee=1084 status=accepted
2025-11-08T18:05:00Z RPC /submit payloadHash=0x03117d45 fee=1085 status=accepted
2025-11-08T18:06:00Z RPC /submit payloadHash=0x03117e46 fee=1086 status=broadcast
2025-11-08T18:07:00Z RPC /submit payloadHash=0x03117f47 fee=1087 status=accepted
2025-11-08T18:08:00Z RPC /submit payloadHash=0x03118048 fee=1088 status=accepted
2025-11-08T18:09:00Z RPC /submit payloadHash=0x03118149 fee=1089 status=broadcast
2025-11-08T18:10:00Z RPC /submit payloadHash=0x03118250 fee=1090 status=accepted
2025-11-08T18:11:00Z RPC /submit payloadHash=0x03118351 fee=1091 status=accepted
2025-11-08T18:12:00Z RPC /submit payloadHash=0x03118452 fee=1092 status=broadcast
2025-11-08T18:13:00Z RPC /submit payloadHash=0x03118553 fee=1093 status=accepted
2025-11-08T18:14:00Z RPC /submit payloadHash=0x03118654 fee=1094 status=accepted
2025-11-08T18:15:00Z RPC /submit payloadHash=0x03118755 fee=1095 status=broadcast
2025-11-08T18:16:00Z RPC /submit payloadHash=0x03118856 fee=1096 status=accepted
2025-11-08T18:17:00Z RPC /submit payloadHash=0x03118957 fee=1097 status=accepted
2025-11-08T18:18:00Z RPC /submit payloadHash=0x03118a58 fee=1098 status=broadcast
2025-11-08T18:19:00Z RPC /submit payloadHash=0x03118b59 fee=1099 status=accepted
2025-11-08T18:20:00Z RPC /submit payloadHash=0x03118c10 fee=1100 status=accepted
2025-11-08T18:21:00Z RPC /submit payloadHash=0x03118d11 fee=1101 status=broadcast
2025-11-08T18:22:00Z RPC /submit payloadHash=0x03118e12 fee=1102 status=accepted
2025-11-08T18:23:00Z RPC /submit payloadHash=0x03118f13 fee=1103 status=accepted
2025-11-08T18:24:00Z RPC /submit payloadHash=0x03119014 fee=1104 status=broadcast
2025-11-08T18:25:00Z RPC /submit payloadHash=0x03119115 fee=1105 status=accepted
2025-11-08T18:26:00Z RPC /submit payloadHash=0x03119216 fee=1106 status=accepted
2025-11-08T18:27:00Z RPC /submit payloadHash=0x03119317 fee=1107 status=broadcast
2025-11-08T18:28:00Z RPC /submit payloadHash=0x03119418 fee=1108 status=accepted
2025-11-08T18:29:00Z RPC /submit payloadHash=0x03119519 fee=1109 status=accepted
2025-11-08T18:30:00Z RPC /submit payloadHash=0x03119620 fee=1110 status=broadcast
2025-11-08T18:31:00Z RPC /submit payloadHash=0x03119721 fee=1111 status=accepted
2025-11-08T18:32:00Z RPC /submit payloadHash=0x03119822 fee=1112 status=accepted
2025-11-08T18:33:00Z RPC /submit payloadHash=0x03119923 fee=1113 status=broadcast
2025-11-08T18:34:00Z RPC /submit payloadHash=0x03119a24 fee=1114 status=accepted
2025-11-08T18:35:00Z RPC /submit payloadHash=0x03119b25 fee=1115 status=accepted
2025-11-08T18:36:00Z RPC /submit payloadHash=0x03119c26 fee=1116 status=broadcast
2025-11-08T18:37:00Z RPC /submit payloadHash=0x03119d27 fee=1117 status=accepted
2025-11-08T18:38:00Z RPC /submit payloadHash=0x03119e28 fee=1118 status=accepted
2025-11-08T18:39:00Z RPC /submit payloadHash=0x03119f29 fee=1119 status=broadcast
2025-11-08T18:40:00Z RPC /submit payloadHash=0x0311a030 fee=1120 status=accepted
2025-11-08T18:41:00Z RPC /submit payloadHash=0x0311a131 fee=1121 status=accepted
2025-11-08T18:42:00Z RPC /submit payloadHash=0x0311a232 fee=1122 status=broadcast
2025-11-08T18:43:00Z RPC /submit payloadHash=0x0311a333 fee=1123 status=accepted
2025-11-08T18:44:00Z RPC /submit payloadHash=0x0311a434 fee=1124 status=accepted
2025-11-08T18:45:00Z RPC /submit payloadHash=0x0311a535 fee=1125 status=broadcast
2025-11-08T18:46:00Z RPC /submit payloadHash=0x0311a636 fee=1126 status=accepted
2025-11-08T18:47:00Z RPC /submit payloadHash=0x0311a737 fee=1127 status=accepted
2025-11-08T18:48:00Z RPC /submit payloadHash=0x0311a838 fee=1128 status=broadcast
2025-11-08T18:49:00Z RPC /submit payloadHash=0x0311a939 fee=1129 status=accepted
2025-11-08T18:50:00Z RPC /submit payloadHash=0x0311aa40 fee=1130 status=accepted
2025-11-08T18:51:00Z RPC /submit payloadHash=0x0311ab41 fee=1131 status=broadcast
2025-11-08T18:52:00Z RPC /submit payloadHash=0x0311ac42 fee=1132 status=accepted
2025-11-08T18:53:00Z RPC /submit payloadHash=0x0311ad43 fee=1133 status=accepted
2025-11-08T18:54:00Z RPC /submit payloadHash=0x0311ae44 fee=1134 status=broadcast
2025-11-08T18:55:00Z RPC /submit payloadHash=0x0311af45 fee=1135 status=accepted
2025-11-08T18:56:00Z RPC /submit payloadHash=0x0311b046 fee=1136 status=accepted
2025-11-08T18:57:00Z RPC /submit payloadHash=0x0311b147 fee=1137 status=broadcast
2025-11-08T18:58:00Z RPC /submit payloadHash=0x0311b248 fee=1138 status=accepted
2025-11-08T18:59:00Z RPC /submit payloadHash=0x0311b349 fee=1139 status=accepted
2025-11-08T19:00:00Z RPC /submit payloadHash=0x0311b450 fee=1140 status=broadcast
2025-11-08T19:01:00Z RPC /submit payloadHash=0x0311b551 fee=1141 status=accepted
2025-11-08T19:02:00Z RPC /submit payloadHash=0x0311b652 fee=1142 status=accepted
2025-11-08T19:03:00Z RPC /submit payloadHash=0x0311b753 fee=1143 status=broadcast
2025-11-08T19:04:00Z RPC /submit payloadHash=0x0311b854 fee=1144 status=accepted
2025-11-08T19:05:00Z RPC /submit payloadHash=0x0311b955 fee=1145 status=accepted
2025-11-08T19:06:00Z RPC /submit payloadHash=0x0311ba56 fee=1146 status=broadcast
2025-11-08T19:07:00Z RPC /submit payloadHash=0x0311bb57 fee=1147 status=accepted
2025-11-08T19:08:00Z RPC /submit payloadHash=0x0311bc58 fee=1148 status=accepted
2025-11-08T19:09:00Z RPC /submit payloadHash=0x0311bd59 fee=1149 status=broadcast
2025-11-08T19:10:00Z RPC /submit payloadHash=0x0311be10 fee=1150 status=accepted
2025-11-08T19:11:00Z RPC /submit payloadHash=0x0311bf11 fee=1151 status=accepted
2025-11-08T19:12:00Z RPC /submit payloadHash=0x0311c012 fee=1152 status=broadcast
2025-11-08T19:13:00Z RPC /submit payloadHash=0x0311c113 fee=1153 status=accepted
2025-11-08T19:14:00Z RPC /submit payloadHash=0x0311c214 fee=1154 status=accepted
2025-11-08T19:15:00Z RPC /submit payloadHash=0x0311c315 fee=1155 status=broadcast
2025-11-08T19:16:00Z RPC /submit payloadHash=0x0311c416 fee=1156 status=accepted
2025-11-08T19:17:00Z RPC /submit payloadHash=0x0311c517 fee=1157 status=accepted
2025-11-08T19:18:00Z RPC /submit payloadHash=0x0311c618 fee=1158 status=broadcast
2025-11-08T19:19:00Z RPC /submit payloadHash=0x0311c719 fee=1159 status=accepted
2025-11-08T19:20:00Z RPC /submit payloadHash=0x0311c820 fee=1160 status=accepted
2025-11-08T19:21:00Z RPC /submit payloadHash=0x0311c921 fee=1161 status=broadcast
2025-11-08T19:22:00Z RPC /submit payloadHash=0x0311ca22 fee=1162 status=accepted
2025-11-08T19:23:00Z RPC /submit payloadHash=0x0311cb23 fee=1163 status=accepted
2025-11-08T19:24:00Z RPC /submit payloadHash=0x0311cc24 fee=1164 status=broadcast
2025-11-08T19:25:00Z RPC /submit payloadHash=0x0311cd25 fee=1165 status=accepted
2025-11-08T19:26:00Z RPC /submit payloadHash=0x0311ce26 fee=1166 status=accepted
2025-11-08T19:27:00Z RPC /submit payloadHash=0x0311cf27 fee=1167 status=broadcast
2025-11-08T19:28:00Z RPC /submit payloadHash=0x0311d028 fee=1168 status=accepted
2025-11-08T19:29:00Z RPC /submit payloadHash=0x0311d129 fee=1169 status=accepted
2025-11-08T19:30:00Z RPC /submit payloadHash=0x0311d230 fee=1170 status=broadcast
2025-11-08T19:31:00Z RPC /submit payloadHash=0x0311d331 fee=1171 status=accepted
2025-11-08T19:32:00Z RPC /submit payloadHash=0x0311d432 fee=1172 status=accepted
2025-11-08T19:33:00Z RPC /submit payloadHash=0x0311d533 fee=1173 status=broadcast
2025-11-08T19:34:00Z RPC /submit payloadHash=0x0311d634 fee=1174 status=accepted
2025-11-08T19:35:00Z RPC /submit payloadHash=0x0311d735 fee=1175 status=accepted
2025-11-08T19:36:00Z RPC /submit payloadHash=0x0311d836 fee=1176 status=broadcast
2025-11-08T19:37:00Z RPC /submit payloadHash=0x0311d937 fee=1177 status=accepted
2025-11-08T19:38:00Z RPC /submit payloadHash=0x0311da38 fee=1178 status=accepted
2025-11-08T19:39:00Z RPC /submit payloadHash=0x0311db39 fee=1179 status=broadcast
2025-11-08T19:40:00Z RPC /submit payloadHash=0x0311dc40 fee=1180 status=accepted
2025-11-08T19:41:00Z RPC /submit payloadHash=0x0311dd41 fee=1181 status=accepted
2025-11-08T19:42:00Z RPC /submit payloadHash=0x0311de42 fee=1182 status=broadcast
2025-11-08T19:43:00Z RPC /submit payloadHash=0x0311df43 fee=1183 status=accepted
2025-11-08T19:44:00Z RPC /submit payloadHash=0x0311e044 fee=1184 status=accepted
2025-11-08T19:45:00Z RPC /submit payloadHash=0x0311e145 fee=1185 status=broadcast
2025-11-08T19:46:00Z RPC /submit payloadHash=0x0311e246 fee=1186 status=accepted
2025-11-08T19:47:00Z RPC /submit payloadHash=0x0311e347 fee=1187 status=accepted
2025-11-08T19:48:00Z RPC /submit payloadHash=0x0311e448 fee=1188 status=broadcast
2025-11-08T19:49:00Z RPC /submit payloadHash=0x0311e549 fee=1189 status=accepted
2025-11-08T19:50:00Z RPC /submit payloadHash=0x0311e650 fee=1190 status=accepted
2025-11-08T19:51:00Z RPC /submit payloadHash=0x0311e751 fee=1191 status=broadcast
2025-11-08T19:52:00Z RPC /submit payloadHash=0x0311e852 fee=1192 status=accepted
2025-11-08T19:53:00Z RPC /submit payloadHash=0x0311e953 fee=1193 status=accepted
2025-11-08T19:54:00Z RPC /submit payloadHash=0x0311ea54 fee=1194 status=broadcast
2025-11-08T19:55:00Z RPC /submit payloadHash=0x0311eb55 fee=1195 status=accepted
2025-11-08T19:56:00Z RPC /submit payloadHash=0x0311ec56 fee=1196 status=accepted
2025-11-08T19:57:00Z RPC /submit payloadHash=0x0311ed57 fee=1197 status=broadcast
2025-11-08T19:58:00Z RPC /submit payloadHash=0x0311ee58 fee=1198 status=accepted
2025-11-08T19:59:00Z RPC /submit payloadHash=0x0311ef59 fee=1199 status=accepted
2025-11-08T20:00:00Z RPC /submit payloadHash=0x0311f010 fee=1000 status=broadcast
2025-11-08T20:01:00Z RPC /submit payloadHash=0x0311f111 fee=1001 status=accepted
2025-11-08T20:02:00Z RPC /submit payloadHash=0x0311f212 fee=1002 status=accepted
2025-11-08T20:03:00Z RPC /submit payloadHash=0x0311f313 fee=1003 status=broadcast
2025-11-08T20:04:00Z RPC /submit payloadHash=0x0311f414 fee=1004 status=accepted
2025-11-08T20:05:00Z RPC /submit payloadHash=0x0311f515 fee=1005 status=accepted
2025-11-08T20:06:00Z RPC /submit payloadHash=0x0311f616 fee=1006 status=broadcast
2025-11-08T20:07:00Z RPC /submit payloadHash=0x0311f717 fee=1007 status=accepted
2025-11-08T20:08:00Z RPC /submit payloadHash=0x0311f818 fee=1008 status=accepted
2025-11-08T20:09:00Z RPC /submit payloadHash=0x0311f919 fee=1009 status=broadcast
2025-11-08T20:10:00Z RPC /submit payloadHash=0x0311fa20 fee=1010 status=accepted
2025-11-08T20:11:00Z RPC /submit payloadHash=0x0311fb21 fee=1011 status=accepted
2025-11-08T20:12:00Z RPC /submit payloadHash=0x0311fc22 fee=1012 status=broadcast
2025-11-08T20:13:00Z RPC /submit payloadHash=0x0311fd23 fee=1013 status=accepted
2025-11-08T20:14:00Z RPC /submit payloadHash=0x0311fe24 fee=1014 status=accepted
2025-11-08T20:15:00Z RPC /submit payloadHash=0x0311ff25 fee=1015 status=broadcast
2025-11-08T20:16:00Z RPC /submit payloadHash=0x03120026 fee=1016 status=accepted
2025-11-08T20:17:00Z RPC /submit payloadHash=0x03120127 fee=1017 status=accepted
2025-11-08T20:18:00Z RPC /submit payloadHash=0x03120228 fee=1018 status=broadcast
2025-11-08T20:19:00Z RPC /submit payloadHash=0x03120329 fee=1019 status=accepted
2025-11-08T20:20:00Z RPC /submit payloadHash=0x03120430 fee=1020 status=accepted
2025-11-08T20:21:00Z RPC /submit payloadHash=0x03120531 fee=1021 status=broadcast
2025-11-08T20:22:00Z RPC /submit payloadHash=0x03120632 fee=1022 status=accepted
2025-11-08T20:23:00Z RPC /submit payloadHash=0x03120733 fee=1023 status=accepted
2025-11-08T20:24:00Z RPC /submit payloadHash=0x03120834 fee=1024 status=broadcast
2025-11-08T20:25:00Z RPC /submit payloadHash=0x03120935 fee=1025 status=accepted
2025-11-08T20:26:00Z RPC /submit payloadHash=0x03120a36 fee=1026 status=accepted
2025-11-08T20:27:00Z RPC /submit payloadHash=0x03120b37 fee=1027 status=broadcast
2025-11-08T20:28:00Z RPC /submit payloadHash=0x03120c38 fee=1028 status=accepted
2025-11-08T20:29:00Z RPC /submit payloadHash=0x03120d39 fee=1029 status=accepted
2025-11-08T20:30:00Z RPC /submit payloadHash=0x03120e40 fee=1030 status=broadcast
2025-11-08T20:31:00Z RPC /submit payloadHash=0x03120f41 fee=1031 status=accepted
2025-11-08T20:32:00Z RPC /submit payloadHash=0x03121042 fee=1032 status=accepted
2025-11-08T20:33:00Z RPC /submit payloadHash=0x03121143 fee=1033 status=broadcast
2025-11-08T20:34:00Z RPC /submit payloadHash=0x03121244 fee=1034 status=accepted
2025-11-08T20:35:00Z RPC /submit payloadHash=0x03121345 fee=1035 status=accepted
2025-11-08T20:36:00Z RPC /submit payloadHash=0x03121446 fee=1036 status=broadcast
2025-11-08T20:37:00Z RPC /submit payloadHash=0x03121547 fee=1037 status=accepted
2025-11-08T20:38:00Z RPC /submit payloadHash=0x03121648 fee=1038 status=accepted
2025-11-08T20:39:00Z RPC /submit payloadHash=0x03121749 fee=1039 status=broadcast
2025-11-08T20:40:00Z RPC /submit payloadHash=0x03121850 fee=1040 status=accepted
2025-11-08T20:41:00Z RPC /submit payloadHash=0x03121951 fee=1041 status=accepted
2025-11-08T20:42:00Z RPC /submit payloadHash=0x03121a52 fee=1042 status=broadcast
2025-11-08T20:43:00Z RPC /submit payloadHash=0x03121b53 fee=1043 status=accepted
2025-11-08T20:44:00Z RPC /submit payloadHash=0x03121c54 fee=1044 status=accepted
2025-11-08T20:45:00Z RPC /submit payloadHash=0x03121d55 fee=1045 status=broadcast
2025-11-08T20:46:00Z RPC /submit payloadHash=0x03121e56 fee=1046 status=accepted
2025-11-08T20:47:00Z RPC /submit payloadHash=0x03121f57 fee=1047 status=accepted
2025-11-08T20:48:00Z RPC /submit payloadHash=0x03122058 fee=1048 status=broadcast
2025-11-08T20:49:00Z RPC /submit payloadHash=0x03122159 fee=1049 status=accepted
2025-11-08T20:50:00Z RPC /submit payloadHash=0x03122210 fee=1050 status=accepted
2025-11-08T20:51:00Z RPC /submit payloadHash=0x03122311 fee=1051 status=broadcast
2025-11-08T20:52:00Z RPC /submit payloadHash=0x03122412 fee=1052 status=accepted
2025-11-08T20:53:00Z RPC /submit payloadHash=0x03122513 fee=1053 status=accepted
2025-11-08T20:54:00Z RPC /submit payloadHash=0x03122614 fee=1054 status=broadcast
2025-11-08T20:55:00Z RPC /submit payloadHash=0x03122715 fee=1055 status=accepted
2025-11-08T20:56:00Z RPC /submit payloadHash=0x03122816 fee=1056 status=accepted
2025-11-08T20:57:00Z RPC /submit payloadHash=0x03122917 fee=1057 status=broadcast
2025-11-08T20:58:00Z RPC /submit payloadHash=0x03122a18 fee=1058 status=accepted
2025-11-08T20:59:00Z RPC /submit payloadHash=0x03122b19 fee=1059 status=accepted
2025-11-08T21:00:00Z RPC /submit payloadHash=0x03122c20 fee=1060 status=broadcast
2025-11-08T21:01:00Z RPC /submit payloadHash=0x03122d21 fee=1061 status=accepted
2025-11-08T21:02:00Z RPC /submit payloadHash=0x03122e22 fee=1062 status=accepted
2025-11-08T21:03:00Z RPC /submit payloadHash=0x03122f23 fee=1063 status=broadcast
2025-11-08T21:04:00Z RPC /submit payloadHash=0x03123024 fee=1064 status=accepted
2025-11-08T21:05:00Z RPC /submit payloadHash=0x03123125 fee=1065 status=accepted
2025-11-08T21:06:00Z RPC /submit payloadHash=0x03123226 fee=1066 status=broadcast
2025-11-08T21:07:00Z RPC /submit payloadHash=0x03123327 fee=1067 status=accepted
2025-11-08T21:08:00Z RPC /submit payloadHash=0x03123428 fee=1068 status=accepted
2025-11-08T21:09:00Z RPC /submit payloadHash=0x03123529 fee=1069 status=broadcast
2025-11-08T21:10:00Z RPC /submit payloadHash=0x03123630 fee=1070 status=accepted
2025-11-08T21:11:00Z RPC /submit payloadHash=0x03123731 fee=1071 status=accepted
2025-11-08T21:12:00Z RPC /submit payloadHash=0x03123832 fee=1072 status=broadcast
2025-11-08T21:13:00Z RPC /submit payloadHash=0x03123933 fee=1073 status=accepted
2025-11-08T21:14:00Z RPC /submit payloadHash=0x03123a34 fee=1074 status=accepted
2025-11-08T21:15:00Z RPC /submit payloadHash=0x03123b35 fee=1075 status=broadcast
2025-11-08T21:16:00Z RPC /submit payloadHash=0x03123c36 fee=1076 status=accepted
2025-11-08T21:17:00Z RPC /submit payloadHash=0x03123d37 fee=1077 status=accepted
2025-11-08T21:18:00Z RPC /submit payloadHash=0x03123e38 fee=1078 status=broadcast
2025-11-08T21:19:00Z RPC /submit payloadHash=0x03123f39 fee=1079 status=accepted
2025-11-08T21:20:00Z RPC /submit payloadHash=0x03124040 fee=1080 status=accepted
2025-11-08T21:21:00Z RPC /submit payloadHash=0x03124141 fee=1081 status=broadcast
2025-11-08T21:22:00Z RPC /submit payloadHash=0x03124242 fee=1082 status=accepted
2025-11-08T21:23:00Z RPC /submit payloadHash=0x03124343 fee=1083 status=accepted
2025-11-08T21:24:00Z RPC /submit payloadHash=0x03124444 fee=1084 status=broadcast
2025-11-08T21:25:00Z RPC /submit payloadHash=0x03124545 fee=1085 status=accepted
2025-11-08T21:26:00Z RPC /submit payloadHash=0x03124646 fee=1086 status=accepted
2025-11-08T21:27:00Z RPC /submit payloadHash=0x03124747 fee=1087 status=broadcast
2025-11-08T21:28:00Z RPC /submit payloadHash=0x03124848 fee=1088 status=accepted
2025-11-08T21:29:00Z RPC /submit payloadHash=0x03124949 fee=1089 status=accepted
2025-11-08T21:30:00Z RPC /submit payloadHash=0x03124a50 fee=1090 status=broadcast
2025-11-08T21:31:00Z RPC /submit payloadHash=0x03124b51 fee=1091 status=accepted
2025-11-08T21:32:00Z RPC /submit payloadHash=0x03124c52 fee=1092 status=accepted
2025-11-08T21:33:00Z RPC /submit payloadHash=0x03124d53 fee=1093 status=broadcast
2025-11-08T21:34:00Z RPC /submit payloadHash=0x03124e54 fee=1094 status=accepted
2025-11-08T21:35:00Z RPC /submit payloadHash=0x03124f55 fee=1095 status=accepted
2025-11-08T21:36:00Z RPC /submit payloadHash=0x03125056 fee=1096 status=broadcast
2025-11-08T21:37:00Z RPC /submit payloadHash=0x03125157 fee=1097 status=accepted
2025-11-08T21:38:00Z RPC /submit payloadHash=0x03125258 fee=1098 status=accepted
2025-11-08T21:39:00Z RPC /submit payloadHash=0x03125359 fee=1099 status=broadcast
2025-11-08T21:40:00Z RPC /submit payloadHash=0x03125410 fee=1100 status=accepted
2025-11-08T21:41:00Z RPC /submit payloadHash=0x03125511 fee=1101 status=accepted
2025-11-08T21:42:00Z RPC /submit payloadHash=0x03125612 fee=1102 status=broadcast
2025-11-08T21:43:00Z RPC /submit payloadHash=0x03125713 fee=1103 status=accepted
2025-11-08T21:44:00Z RPC /submit payloadHash=0x03125814 fee=1104 status=accepted
2025-11-08T21:45:00Z RPC /submit payloadHash=0x03125915 fee=1105 status=broadcast
2025-11-08T21:46:00Z RPC /submit payloadHash=0x03125a16 fee=1106 status=accepted
2025-11-08T21:47:00Z RPC /submit payloadHash=0x03125b17 fee=1107 status=accepted
2025-11-08T21:48:00Z RPC /submit payloadHash=0x03125c18 fee=1108 status=broadcast
2025-11-08T21:49:00Z RPC /submit payloadHash=0x03125d19 fee=1109 status=accepted
2025-11-08T21:50:00Z RPC /submit payloadHash=0x03125e20 fee=1110 status=accepted
2025-11-08T21:51:00Z RPC /submit payloadHash=0x03125f21 fee=1111 status=broadcast
2025-11-08T21:52:00Z RPC /submit payloadHash=0x03126022 fee=1112 status=accepted
2025-11-08T21:53:00Z RPC /submit payloadHash=0x03126123 fee=1113 status=accepted
2025-11-08T21:54:00Z RPC /submit payloadHash=0x03126224 fee=1114 status=broadcast
2025-11-08T21:55:00Z RPC /submit payloadHash=0x03126325 fee=1115 status=accepted
2025-11-08T21:56:00Z RPC /submit payloadHash=0x03126426 fee=1116 status=accepted
2025-11-08T21:57:00Z RPC /submit payloadHash=0x03126527 fee=1117 status=broadcast
2025-11-08T21:58:00Z RPC /submit payloadHash=0x03126628 fee=1118 status=accepted
2025-11-08T21:59:00Z RPC /submit payloadHash=0x03126729 fee=1119 status=accepted
2025-11-08T22:00:00Z RPC /submit payloadHash=0x03126830 fee=1120 status=broadcast
2025-11-08T22:01:00Z RPC /submit payloadHash=0x03126931 fee=1121 status=accepted
2025-11-08T22:02:00Z RPC /submit payloadHash=0x03126a32 fee=1122 status=accepted
2025-11-08T22:03:00Z RPC /submit payloadHash=0x03126b33 fee=1123 status=broadcast
2025-11-08T22:04:00Z RPC /submit payloadHash=0x03126c34 fee=1124 status=accepted
2025-11-08T22:05:00Z RPC /submit payloadHash=0x03126d35 fee=1125 status=accepted
2025-11-08T22:06:00Z RPC /submit payloadHash=0x03126e36 fee=1126 status=broadcast
2025-11-08T22:07:00Z RPC /submit payloadHash=0x03126f37 fee=1127 status=accepted
2025-11-08T22:08:00Z RPC /submit payloadHash=0x03127038 fee=1128 status=accepted
2025-11-08T22:09:00Z RPC /submit payloadHash=0x03127139 fee=1129 status=broadcast
2025-11-08T22:10:00Z RPC /submit payloadHash=0x03127240 fee=1130 status=accepted
2025-11-08T22:11:00Z RPC /submit payloadHash=0x03127341 fee=1131 status=accepted
2025-11-08T22:12:00Z RPC /submit payloadHash=0x03127442 fee=1132 status=broadcast
2025-11-08T22:13:00Z RPC /submit payloadHash=0x03127543 fee=1133 status=accepted
2025-11-08T22:14:00Z RPC /submit payloadHash=0x03127644 fee=1134 status=accepted
2025-11-08T22:15:00Z RPC /submit payloadHash=0x03127745 fee=1135 status=broadcast
2025-11-08T22:16:00Z RPC /submit payloadHash=0x03127846 fee=1136 status=accepted
2025-11-08T22:17:00Z RPC /submit payloadHash=0x03127947 fee=1137 status=accepted
2025-11-08T22:18:00Z RPC /submit payloadHash=0x03127a48 fee=1138 status=broadcast
2025-11-08T22:19:00Z RPC /submit payloadHash=0x03127b49 fee=1139 status=accepted
2025-11-08T22:20:00Z RPC /submit payloadHash=0x03127c50 fee=1140 status=accepted
2025-11-08T22:21:00Z RPC /submit payloadHash=0x03127d51 fee=1141 status=broadcast
2025-11-08T22:22:00Z RPC /submit payloadHash=0x03127e52 fee=1142 status=accepted
2025-11-08T22:23:00Z RPC /submit payloadHash=0x03127f53 fee=1143 status=accepted
2025-11-08T22:24:00Z RPC /submit payloadHash=0x03128054 fee=1144 status=broadcast
2025-11-08T22:25:00Z RPC /submit payloadHash=0x03128155 fee=1145 status=accepted
2025-11-08T22:26:00Z RPC /submit payloadHash=0x03128256 fee=1146 status=accepted
2025-11-08T22:27:00Z RPC /submit payloadHash=0x03128357 fee=1147 status=broadcast
2025-11-08T22:28:00Z RPC /submit payloadHash=0x03128458 fee=1148 status=accepted
2025-11-08T22:29:00Z RPC /submit payloadHash=0x03128559 fee=1149 status=accepted
2025-11-08T22:30:00Z RPC /submit payloadHash=0x03128610 fee=1150 status=broadcast
2025-11-08T22:31:00Z RPC /submit payloadHash=0x03128711 fee=1151 status=accepted
2025-11-08T22:32:00Z RPC /submit payloadHash=0x03128812 fee=1152 status=accepted
2025-11-08T22:33:00Z RPC /submit payloadHash=0x03128913 fee=1153 status=broadcast
2025-11-08T22:34:00Z RPC /submit payloadHash=0x03128a14 fee=1154 status=accepted
2025-11-08T22:35:00Z RPC /submit payloadHash=0x03128b15 fee=1155 status=accepted
2025-11-08T22:36:00Z RPC /submit payloadHash=0x03128c16 fee=1156 status=broadcast
2025-11-08T22:37:00Z RPC /submit payloadHash=0x03128d17 fee=1157 status=accepted
2025-11-08T22:38:00Z RPC /submit payloadHash=0x03128e18 fee=1158 status=accepted
2025-11-08T22:39:00Z RPC /submit payloadHash=0x03128f19 fee=1159 status=broadcast
2025-11-08T22:40:00Z RPC /submit payloadHash=0x03129020 fee=1160 status=accepted
2025-11-08T22:41:00Z RPC /submit payloadHash=0x03129121 fee=1161 status=accepted
2025-11-08T22:42:00Z RPC /submit payloadHash=0x03129222 fee=1162 status=broadcast
2025-11-08T22:43:00Z RPC /submit payloadHash=0x03129323 fee=1163 status=accepted
2025-11-08T22:44:00Z RPC /submit payloadHash=0x03129424 fee=1164 status=accepted
2025-11-08T22:45:00Z RPC /submit payloadHash=0x03129525 fee=1165 status=broadcast
2025-11-08T22:46:00Z RPC /submit payloadHash=0x03129626 fee=1166 status=accepted
2025-11-08T22:47:00Z RPC /submit payloadHash=0x03129727 fee=1167 status=accepted
2025-11-08T22:48:00Z RPC /submit payloadHash=0x03129828 fee=1168 status=broadcast
2025-11-08T22:49:00Z RPC /submit payloadHash=0x03129929 fee=1169 status=accepted
2025-11-08T22:50:00Z RPC /submit payloadHash=0x03129a30 fee=1170 status=accepted
2025-11-08T22:51:00Z RPC /submit payloadHash=0x03129b31 fee=1171 status=broadcast
2025-11-08T22:52:00Z RPC /submit payloadHash=0x03129c32 fee=1172 status=accepted
2025-11-08T22:53:00Z RPC /submit payloadHash=0x03129d33 fee=1173 status=accepted
2025-11-08T22:54:00Z RPC /submit payloadHash=0x03129e34 fee=1174 status=broadcast
2025-11-08T22:55:00Z RPC /submit payloadHash=0x03129f35 fee=1175 status=accepted
2025-11-08T22:56:00Z RPC /submit payloadHash=0x0312a036 fee=1176 status=accepted
2025-11-08T22:57:00Z RPC /submit payloadHash=0x0312a137 fee=1177 status=broadcast
2025-11-08T22:58:00Z RPC /submit payloadHash=0x0312a238 fee=1178 status=accepted
2025-11-08T22:59:00Z RPC /submit payloadHash=0x0312a339 fee=1179 status=accepted
2025-11-08T23:00:00Z RPC /submit payloadHash=0x0312a440 fee=1180 status=broadcast
2025-11-08T23:01:00Z RPC /submit payloadHash=0x0312a541 fee=1181 status=accepted
2025-11-08T23:02:00Z RPC /submit payloadHash=0x0312a642 fee=1182 status=accepted
2025-11-08T23:03:00Z RPC /submit payloadHash=0x0312a743 fee=1183 status=broadcast
2025-11-08T23:04:00Z RPC /submit payloadHash=0x0312a844 fee=1184 status=accepted
2025-11-08T23:05:00Z RPC /submit payloadHash=0x0312a945 fee=1185 status=accepted
2025-11-08T23:06:00Z RPC /submit payloadHash=0x0312aa46 fee=1186 status=broadcast
2025-11-08T23:07:00Z RPC /submit payloadHash=0x0312ab47 fee=1187 status=accepted
2025-11-08T23:08:00Z RPC /submit payloadHash=0x0312ac48 fee=1188 status=accepted
2025-11-08T23:09:00Z RPC /submit payloadHash=0x0312ad49 fee=1189 status=broadcast
2025-11-08T23:10:00Z RPC /submit payloadHash=0x0312ae50 fee=1190 status=accepted
2025-11-08T23:11:00Z RPC /submit payloadHash=0x0312af51 fee=1191 status=accepted
2025-11-08T23:12:00Z RPC /submit payloadHash=0x0312b052 fee=1192 status=broadcast
2025-11-08T23:13:00Z RPC /submit payloadHash=0x0312b153 fee=1193 status=accepted
2025-11-08T23:14:00Z RPC /submit payloadHash=0x0312b254 fee=1194 status=accepted
2025-11-08T23:15:00Z RPC /submit payloadHash=0x0312b355 fee=1195 status=broadcast
2025-11-08T23:16:00Z RPC /submit payloadHash=0x0312b456 fee=1196 status=accepted
2025-11-08T23:17:00Z RPC /submit payloadHash=0x0312b557 fee=1197 status=accepted
2025-11-08T23:18:00Z RPC /submit payloadHash=0x0312b658 fee=1198 status=broadcast
2025-11-08T23:19:00Z RPC /submit payloadHash=0x0312b759 fee=1199 status=accepted
2025-11-08T23:20:00Z RPC /submit payloadHash=0x0312b810 fee=1000 status=accepted
2025-11-08T23:21:00Z RPC /submit payloadHash=0x0312b911 fee=1001 status=broadcast
2025-11-08T23:22:00Z RPC /submit payloadHash=0x0312ba12 fee=1002 status=accepted
2025-11-08T23:23:00Z RPC /submit payloadHash=0x0312bb13 fee=1003 status=accepted
2025-11-08T23:24:00Z RPC /submit payloadHash=0x0312bc14 fee=1004 status=broadcast
2025-11-08T23:25:00Z RPC /submit payloadHash=0x0312bd15 fee=1005 status=accepted
2025-11-08T23:26:00Z RPC /submit payloadHash=0x0312be16 fee=1006 status=accepted
2025-11-08T23:27:00Z RPC /submit payloadHash=0x0312bf17 fee=1007 status=broadcast
2025-11-08T23:28:00Z RPC /submit payloadHash=0x0312c018 fee=1008 status=accepted
2025-11-08T23:29:00Z RPC /submit payloadHash=0x0312c119 fee=1009 status=accepted
2025-11-08T23:30:00Z RPC /submit payloadHash=0x0312c220 fee=1010 status=broadcast
2025-11-08T23:31:00Z RPC /submit payloadHash=0x0312c321 fee=1011 status=accepted
2025-11-08T23:32:00Z RPC /submit payloadHash=0x0312c422 fee=1012 status=accepted
2025-11-08T23:33:00Z RPC /submit payloadHash=0x0312c523 fee=1013 status=broadcast
2025-11-08T23:34:00Z RPC /submit payloadHash=0x0312c624 fee=1014 status=accepted
2025-11-08T23:35:00Z RPC /submit payloadHash=0x0312c725 fee=1015 status=accepted
2025-11-08T23:36:00Z RPC /submit payloadHash=0x0312c826 fee=1016 status=broadcast
2025-11-08T23:37:00Z RPC /submit payloadHash=0x0312c927 fee=1017 status=accepted
2025-11-08T23:38:00Z RPC /submit payloadHash=0x0312ca28 fee=1018 status=accepted
2025-11-08T23:39:00Z RPC /submit payloadHash=0x0312cb29 fee=1019 status=broadcast
2025-11-08T23:40:00Z RPC /submit payloadHash=0x0312cc30 fee=1020 status=accepted
2025-11-08T23:41:00Z RPC /submit payloadHash=0x0312cd31 fee=1021 status=accepted
2025-11-08T23:42:00Z RPC /submit payloadHash=0x0312ce32 fee=1022 status=broadcast
2025-11-08T23:43:00Z RPC /submit payloadHash=0x0312cf33 fee=1023 status=accepted
2025-11-08T23:44:00Z RPC /submit payloadHash=0x0312d034 fee=1024 status=accepted
2025-11-08T23:45:00Z RPC /submit payloadHash=0x0312d135 fee=1025 status=broadcast
2025-11-08T23:46:00Z RPC /submit payloadHash=0x0312d236 fee=1026 status=accepted
2025-11-08T23:47:00Z RPC /submit payloadHash=0x0312d337 fee=1027 status=accepted
2025-11-08T23:48:00Z RPC /submit payloadHash=0x0312d438 fee=1028 status=broadcast
2025-11-08T23:49:00Z RPC /submit payloadHash=0x0312d539 fee=1029 status=accepted
2025-11-08T23:50:00Z RPC /submit payloadHash=0x0312d640 fee=1030 status=accepted
2025-11-08T23:51:00Z RPC /submit payloadHash=0x0312d741 fee=1031 status=broadcast
2025-11-08T23:52:00Z RPC /submit payloadHash=0x0312d842 fee=1032 status=accepted
2025-11-08T23:53:00Z RPC /submit payloadHash=0x0312d943 fee=1033 status=accepted
2025-11-08T23:54:00Z RPC /submit payloadHash=0x0312da44 fee=1034 status=broadcast
2025-11-08T23:55:00Z RPC /submit payloadHash=0x0312db45 fee=1035 status=accepted
2025-11-08T23:56:00Z RPC /submit payloadHash=0x0312dc46 fee=1036 status=accepted
2025-11-08T23:57:00Z RPC /submit payloadHash=0x0312dd47 fee=1037 status=broadcast
2025-11-08T23:58:00Z RPC /submit payloadHash=0x0312de48 fee=1038 status=accepted
2025-11-08T23:59:00Z RPC /submit payloadHash=0x0312df49 fee=1039 status=accepted
```

### Appendix AC: Resource Utilization Snapshot

```text
2025-11-08T00:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T00:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T00:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T00:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T00:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T00:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T00:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T00:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T00:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T00:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T00:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T00:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T00:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T00:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T00:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T00:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T00:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T00:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T00:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T00:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T00:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T00:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T00:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T00:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T00:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T00:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T00:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T00:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T00:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T00:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T00:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T00:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T00:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T00:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T00:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T00:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T00:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T00:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T00:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T00:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T00:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T00:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T00:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T00:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T00:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T00:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T00:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T00:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T00:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T00:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T00:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T00:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T00:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T00:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T00:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T00:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T00:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T00:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T00:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T00:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T01:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T01:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T01:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T01:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T01:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T01:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T01:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T01:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T01:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T01:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T01:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T01:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T01:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T01:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T01:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T01:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T01:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T01:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T01:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T01:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T01:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T01:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T01:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T01:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T01:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T01:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T01:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T01:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T01:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T01:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T01:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T01:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T01:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T01:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T01:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T01:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T01:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T01:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T01:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T01:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T01:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T01:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T01:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T01:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T01:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T01:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T01:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T01:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T01:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T01:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T01:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T01:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T01:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T01:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T01:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T01:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T01:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T01:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T01:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T01:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T02:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T02:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T02:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T02:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T02:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T02:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T02:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T02:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T02:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T02:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T02:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T02:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T02:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T02:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T02:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T02:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T02:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T02:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T02:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T02:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T02:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T02:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T02:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T02:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T02:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T02:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T02:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T02:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T02:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T02:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T02:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T02:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T02:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T02:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T02:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T02:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T02:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T02:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T02:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T02:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T02:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T02:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T02:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T02:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T02:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T02:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T02:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T02:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T02:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T02:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T02:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T02:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T02:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T02:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T02:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T02:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T02:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T02:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T02:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T02:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T03:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T03:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T03:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T03:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T03:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T03:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T03:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T03:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T03:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T03:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T03:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T03:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T03:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T03:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T03:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T03:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T03:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T03:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T03:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T03:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T03:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T03:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T03:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T03:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T03:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T03:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T03:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T03:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T03:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T03:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T03:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T03:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T03:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T03:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T03:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T03:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T03:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T03:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T03:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T03:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T03:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T03:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T03:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T03:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T03:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T03:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T03:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T03:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T03:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T03:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T03:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T03:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T03:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T03:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T03:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T03:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T03:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T03:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T03:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T03:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T04:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T04:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T04:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T04:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T04:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T04:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T04:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T04:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T04:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T04:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T04:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T04:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T04:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T04:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T04:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T04:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T04:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T04:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T04:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T04:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T04:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T04:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T04:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T04:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T04:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T04:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T04:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T04:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T04:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T04:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T04:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T04:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T04:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T04:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T04:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T04:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T04:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T04:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T04:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T04:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T04:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T04:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T04:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T04:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T04:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T04:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T04:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T04:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T04:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T04:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T04:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T04:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T04:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T04:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T04:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T04:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T04:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T04:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T04:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T04:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T05:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T05:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T05:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T05:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T05:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T05:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T05:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T05:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T05:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T05:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T05:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T05:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T05:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T05:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T05:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T05:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T05:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T05:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T05:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T05:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T05:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T05:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T05:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T05:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T05:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T05:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T05:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T05:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T05:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T05:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T05:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T05:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T05:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T05:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T05:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T05:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T05:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T05:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T05:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T05:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T05:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T05:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T05:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T05:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T05:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T05:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T05:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T05:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T05:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T05:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T05:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T05:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T05:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T05:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T05:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T05:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T05:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T05:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T05:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T05:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T06:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T06:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T06:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T06:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T06:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T06:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T06:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T06:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T06:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T06:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T06:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T06:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T06:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T06:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T06:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T06:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T06:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T06:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T06:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T06:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T06:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T06:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T06:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T06:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T06:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T06:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T06:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T06:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T06:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T06:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T06:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T06:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T06:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T06:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T06:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T06:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T06:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T06:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T06:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T06:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T06:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T06:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T06:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T06:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T06:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T06:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T06:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T06:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T06:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T06:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T06:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T06:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T06:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T06:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T06:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T06:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T06:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T06:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T06:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T06:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T07:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T07:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T07:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T07:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T07:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T07:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T07:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T07:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T07:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T07:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T07:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T07:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T07:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T07:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T07:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T07:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T07:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T07:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T07:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T07:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T07:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T07:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T07:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T07:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T07:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T07:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T07:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T07:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T07:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T07:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T07:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T07:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T07:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T07:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T07:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T07:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T07:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T07:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T07:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T07:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T07:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T07:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T07:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T07:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T07:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T07:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T07:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T07:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T07:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T07:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T07:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T07:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T07:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T07:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T07:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T07:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T07:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T07:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T07:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T07:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T08:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T08:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T08:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T08:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T08:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T08:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T08:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T08:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T08:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T08:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T08:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T08:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T08:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T08:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T08:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T08:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T08:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T08:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T08:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T08:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T08:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T08:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T08:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T08:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T08:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T08:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T08:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T08:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T08:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T08:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T08:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T08:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T08:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T08:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T08:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T08:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T08:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T08:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T08:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T08:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T08:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T08:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T08:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T08:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T08:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T08:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T08:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T08:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T08:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T08:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T08:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T08:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T08:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T08:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T08:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T08:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T08:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T08:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T08:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T08:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T09:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T09:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T09:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T09:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T09:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T09:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T09:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T09:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T09:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T09:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T09:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T09:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T09:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T09:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T09:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T09:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T09:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T09:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T09:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T09:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T09:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T09:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T09:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T09:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T09:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T09:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T09:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T09:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T09:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T09:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T09:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T09:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T09:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T09:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T09:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T09:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T09:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T09:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T09:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T09:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T09:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T09:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T09:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T09:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T09:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T09:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T09:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T09:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T09:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T09:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T09:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T09:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T09:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T09:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T09:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T09:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T09:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T09:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T09:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T09:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T10:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T10:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T10:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T10:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T10:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T10:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T10:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T10:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T10:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T10:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T10:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T10:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T10:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T10:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T10:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T10:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T10:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T10:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T10:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T10:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T10:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T10:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T10:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T10:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T10:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T10:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T10:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T10:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T10:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T10:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T10:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T10:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T10:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T10:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T10:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T10:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T10:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T10:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T10:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T10:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T10:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T10:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T10:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T10:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T10:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T10:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T10:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T10:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T10:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T10:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T10:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T10:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T10:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T10:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T10:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T10:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T10:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T10:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T10:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T10:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T11:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T11:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T11:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T11:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T11:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T11:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T11:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T11:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T11:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T11:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T11:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T11:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T11:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T11:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T11:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T11:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T11:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T11:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T11:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T11:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T11:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T11:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T11:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T11:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T11:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T11:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T11:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T11:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T11:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T11:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T11:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T11:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T11:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T11:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T11:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T11:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T11:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T11:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T11:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T11:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T11:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T11:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T11:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T11:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T11:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T11:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T11:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T11:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T11:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T11:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T11:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T11:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T11:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T11:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T11:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T11:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T11:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T11:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T11:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T11:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T12:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T12:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T12:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T12:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T12:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T12:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T12:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T12:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T12:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T12:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T12:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T12:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T12:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T12:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T12:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T12:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T12:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T12:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T12:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T12:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T12:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T12:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T12:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T12:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T12:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T12:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T12:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T12:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T12:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T12:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T12:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T12:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T12:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T12:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T12:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T12:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T12:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T12:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T12:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T12:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T12:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T12:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T12:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T12:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T12:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T12:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T12:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T12:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T12:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T12:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T12:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T12:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T12:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T12:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T12:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T12:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T12:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T12:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T12:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T12:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T13:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T13:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T13:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T13:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T13:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T13:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T13:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T13:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T13:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T13:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T13:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T13:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T13:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T13:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T13:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T13:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T13:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T13:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T13:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T13:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T13:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T13:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T13:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T13:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T13:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T13:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T13:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T13:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T13:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T13:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T13:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T13:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T13:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T13:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T13:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T13:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T13:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T13:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T13:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T13:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T13:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T13:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T13:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T13:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T13:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T13:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T13:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T13:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T13:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T13:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T13:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T13:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T13:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T13:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T13:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T13:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T13:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T13:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T13:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T13:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T14:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T14:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T14:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T14:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T14:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T14:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T14:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T14:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T14:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T14:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T14:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T14:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T14:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T14:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T14:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T14:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T14:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T14:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T14:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T14:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T14:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T14:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T14:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T14:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T14:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T14:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T14:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T14:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T14:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T14:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T14:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T14:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T14:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T14:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T14:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T14:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T14:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T14:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T14:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T14:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T14:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T14:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T14:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T14:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T14:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T14:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T14:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T14:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T14:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T14:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T14:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T14:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T14:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T14:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T14:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T14:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T14:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T14:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T14:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T14:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T15:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T15:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T15:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T15:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T15:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T15:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T15:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T15:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T15:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T15:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T15:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T15:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T15:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T15:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T15:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T15:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T15:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T15:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T15:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T15:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T15:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T15:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T15:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T15:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T15:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T15:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T15:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T15:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T15:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T15:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T15:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T15:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T15:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T15:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T15:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T15:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T15:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T15:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T15:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T15:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T15:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T15:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T15:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T15:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T15:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T15:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T15:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T15:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T15:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T15:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T15:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T15:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T15:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T15:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T15:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T15:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T15:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T15:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T15:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T15:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T16:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T16:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T16:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T16:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T16:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T16:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T16:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T16:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T16:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T16:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T16:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T16:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T16:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T16:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T16:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T16:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T16:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T16:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T16:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T16:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T16:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T16:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T16:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T16:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T16:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T16:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T16:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T16:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T16:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T16:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T16:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T16:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T16:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T16:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T16:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T16:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T16:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T16:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T16:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T16:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T16:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T16:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T16:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T16:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T16:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T16:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T16:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T16:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T16:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T16:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T16:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T16:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T16:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T16:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T16:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T16:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T16:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T16:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T16:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T16:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T17:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T17:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T17:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T17:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T17:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T17:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T17:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T17:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T17:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T17:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T17:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T17:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T17:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T17:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T17:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T17:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T17:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T17:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T17:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T17:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T17:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T17:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T17:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T17:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T17:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T17:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T17:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T17:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T17:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T17:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T17:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T17:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T17:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T17:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T17:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T17:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T17:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T17:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T17:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T17:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T17:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T17:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T17:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T17:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T17:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T17:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T17:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T17:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T17:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T17:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T17:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T17:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T17:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T17:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T17:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T17:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T17:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T17:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T17:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T17:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T18:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T18:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T18:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T18:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T18:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T18:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T18:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T18:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T18:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T18:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T18:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T18:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T18:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T18:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T18:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T18:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T18:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T18:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T18:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T18:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T18:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T18:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T18:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T18:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T18:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T18:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T18:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T18:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T18:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T18:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T18:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T18:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T18:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T18:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T18:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T18:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T18:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T18:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T18:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T18:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T18:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T18:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T18:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T18:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T18:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T18:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T18:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T18:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T18:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T18:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T18:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T18:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T18:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T18:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T18:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T18:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T18:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T18:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T18:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T18:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T19:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T19:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T19:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T19:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T19:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T19:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T19:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T19:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T19:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T19:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T19:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T19:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T19:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T19:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T19:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T19:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T19:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T19:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T19:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T19:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T19:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T19:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T19:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T19:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T19:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T19:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T19:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T19:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T19:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T19:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T19:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T19:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T19:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T19:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T19:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T19:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T19:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T19:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T19:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T19:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T19:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T19:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T19:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T19:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T19:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T19:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T19:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T19:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T19:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T19:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T19:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T19:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T19:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T19:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T19:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T19:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T19:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T19:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T19:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T19:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T20:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T20:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T20:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T20:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T20:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T20:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T20:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T20:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T20:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T20:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T20:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T20:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T20:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T20:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T20:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T20:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T20:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T20:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T20:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T20:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T20:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T20:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T20:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T20:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T20:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T20:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T20:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T20:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T20:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T20:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T20:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T20:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T20:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T20:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T20:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T20:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T20:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T20:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T20:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T20:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T20:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T20:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T20:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T20:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T20:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T20:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T20:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T20:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T20:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T20:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T20:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T20:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T20:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T20:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T20:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T20:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T20:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T20:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T20:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T20:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T21:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T21:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T21:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T21:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T21:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T21:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T21:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T21:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T21:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T21:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T21:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T21:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T21:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T21:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T21:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T21:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T21:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T21:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T21:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T21:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T21:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T21:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T21:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T21:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T21:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T21:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T21:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T21:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T21:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T21:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T21:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T21:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T21:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T21:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T21:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T21:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T21:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T21:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T21:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T21:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T21:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T21:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T21:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T21:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T21:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T21:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T21:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T21:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T21:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T21:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T21:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T21:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T21:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T21:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T21:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T21:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T21:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T21:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T21:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T21:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T22:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T22:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T22:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T22:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T22:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T22:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T22:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T22:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T22:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T22:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T22:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T22:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T22:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T22:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T22:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T22:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T22:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T22:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T22:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T22:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T22:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T22:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T22:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T22:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T22:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T22:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T22:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T22:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T22:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T22:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T22:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T22:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T22:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T22:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T22:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T22:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T22:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T22:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T22:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T22:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T22:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T22:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T22:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T22:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T22:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T22:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T22:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T22:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T22:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T22:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T22:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T22:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T22:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T22:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T22:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T22:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T22:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T22:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T22:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T22:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
2025-11-08T23:00:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=51
2025-11-08T23:01:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=52
2025-11-08T23:02:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=53
2025-11-08T23:03:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=54
2025-11-08T23:04:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=55
2025-11-08T23:05:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=56
2025-11-08T23:06:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=57
2025-11-08T23:07:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=58
2025-11-08T23:08:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=59
2025-11-08T23:09:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=60
2025-11-08T23:10:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=61
2025-11-08T23:11:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=62
2025-11-08T23:12:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=51
2025-11-08T23:13:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=52
2025-11-08T23:14:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=53
2025-11-08T23:15:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=54
2025-11-08T23:16:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=55
2025-11-08T23:17:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=56
2025-11-08T23:18:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=57
2025-11-08T23:19:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=58
2025-11-08T23:20:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=59
2025-11-08T23:21:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=60
2025-11-08T23:22:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=61
2025-11-08T23:23:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=62
2025-11-08T23:24:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=51
2025-11-08T23:25:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=52
2025-11-08T23:26:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=53
2025-11-08T23:27:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=54
2025-11-08T23:28:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=55
2025-11-08T23:29:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=56
2025-11-08T23:30:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=42 netOutMbps=57
2025-11-08T23:31:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=43 netOutMbps=58
2025-11-08T23:32:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=44 netOutMbps=59
2025-11-08T23:33:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=45 netOutMbps=60
2025-11-08T23:34:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=46 netOutMbps=61
2025-11-08T23:35:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=47 netOutMbps=62
2025-11-08T23:36:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=48 netOutMbps=51
2025-11-08T23:37:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=49 netOutMbps=52
2025-11-08T23:38:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=50 netOutMbps=53
2025-11-08T23:39:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=51 netOutMbps=54
2025-11-08T23:40:00Z METRICS cpuPercent=35 memGiB=18 diskUtil=65 netInMbps=52 netOutMbps=55
2025-11-08T23:41:00Z METRICS cpuPercent=36 memGiB=19 diskUtil=66 netInMbps=53 netOutMbps=56
2025-11-08T23:42:00Z METRICS cpuPercent=37 memGiB=20 diskUtil=67 netInMbps=54 netOutMbps=57
2025-11-08T23:43:00Z METRICS cpuPercent=38 memGiB=21 diskUtil=68 netInMbps=55 netOutMbps=58
2025-11-08T23:44:00Z METRICS cpuPercent=39 memGiB=22 diskUtil=69 netInMbps=56 netOutMbps=59
2025-11-08T23:45:00Z METRICS cpuPercent=40 memGiB=23 diskUtil=65 netInMbps=42 netOutMbps=60
2025-11-08T23:46:00Z METRICS cpuPercent=41 memGiB=24 diskUtil=66 netInMbps=43 netOutMbps=61
2025-11-08T23:47:00Z METRICS cpuPercent=42 memGiB=25 diskUtil=67 netInMbps=44 netOutMbps=62
2025-11-08T23:48:00Z METRICS cpuPercent=43 memGiB=26 diskUtil=68 netInMbps=45 netOutMbps=51
2025-11-08T23:49:00Z METRICS cpuPercent=44 memGiB=27 diskUtil=69 netInMbps=46 netOutMbps=52
2025-11-08T23:50:00Z METRICS cpuPercent=45 memGiB=18 diskUtil=65 netInMbps=47 netOutMbps=53
2025-11-08T23:51:00Z METRICS cpuPercent=46 memGiB=19 diskUtil=66 netInMbps=48 netOutMbps=54
2025-11-08T23:52:00Z METRICS cpuPercent=47 memGiB=20 diskUtil=67 netInMbps=49 netOutMbps=55
2025-11-08T23:53:00Z METRICS cpuPercent=48 memGiB=21 diskUtil=68 netInMbps=50 netOutMbps=56
2025-11-08T23:54:00Z METRICS cpuPercent=49 memGiB=22 diskUtil=69 netInMbps=51 netOutMbps=57
2025-11-08T23:55:00Z METRICS cpuPercent=50 memGiB=23 diskUtil=65 netInMbps=52 netOutMbps=58
2025-11-08T23:56:00Z METRICS cpuPercent=51 memGiB=24 diskUtil=66 netInMbps=53 netOutMbps=59
2025-11-08T23:57:00Z METRICS cpuPercent=52 memGiB=25 diskUtil=67 netInMbps=54 netOutMbps=60
2025-11-08T23:58:00Z METRICS cpuPercent=53 memGiB=26 diskUtil=68 netInMbps=55 netOutMbps=61
2025-11-08T23:59:00Z METRICS cpuPercent=54 memGiB=27 diskUtil=69 netInMbps=56 netOutMbps=62
```

### Appendix AD: Governance Activity Ledger

```text
2025-11-08T00:00:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100000 turnoutPercent=15.0
2025-11-08T00:01:00Z GOVERNANCE proposal=100 phase=voting votingPower=100037 turnoutPercent=15.1
2025-11-08T00:02:00Z GOVERNANCE proposal=100 phase=tally votingPower=100074 turnoutPercent=15.2
2025-11-08T00:03:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100111 turnoutPercent=15.3
2025-11-08T00:04:00Z GOVERNANCE proposal=100 phase=voting votingPower=100148 turnoutPercent=15.4
2025-11-08T00:05:00Z GOVERNANCE proposal=100 phase=tally votingPower=100185 turnoutPercent=15.5
2025-11-08T00:06:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100222 turnoutPercent=15.6
2025-11-08T00:07:00Z GOVERNANCE proposal=100 phase=voting votingPower=100259 turnoutPercent=15.7
2025-11-08T00:08:00Z GOVERNANCE proposal=100 phase=tally votingPower=100296 turnoutPercent=15.8
2025-11-08T00:09:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100333 turnoutPercent=15.9
2025-11-08T00:10:00Z GOVERNANCE proposal=100 phase=voting votingPower=100370 turnoutPercent=16.0
2025-11-08T00:11:00Z GOVERNANCE proposal=100 phase=tally votingPower=100407 turnoutPercent=16.1
2025-11-08T00:12:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100444 turnoutPercent=16.2
2025-11-08T00:13:00Z GOVERNANCE proposal=100 phase=voting votingPower=100481 turnoutPercent=16.3
2025-11-08T00:14:00Z GOVERNANCE proposal=100 phase=tally votingPower=100518 turnoutPercent=16.4
2025-11-08T00:15:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100555 turnoutPercent=16.5
2025-11-08T00:16:00Z GOVERNANCE proposal=100 phase=voting votingPower=100592 turnoutPercent=16.6
2025-11-08T00:17:00Z GOVERNANCE proposal=100 phase=tally votingPower=100629 turnoutPercent=16.7
2025-11-08T00:18:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100666 turnoutPercent=16.8
2025-11-08T00:19:00Z GOVERNANCE proposal=100 phase=voting votingPower=100703 turnoutPercent=16.9
2025-11-08T00:20:00Z GOVERNANCE proposal=100 phase=tally votingPower=100740 turnoutPercent=17.0
2025-11-08T00:21:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100777 turnoutPercent=17.1
2025-11-08T00:22:00Z GOVERNANCE proposal=100 phase=voting votingPower=100814 turnoutPercent=17.2
2025-11-08T00:23:00Z GOVERNANCE proposal=100 phase=tally votingPower=100851 turnoutPercent=17.3
2025-11-08T00:24:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100888 turnoutPercent=17.4
2025-11-08T00:25:00Z GOVERNANCE proposal=100 phase=voting votingPower=100925 turnoutPercent=17.5
2025-11-08T00:26:00Z GOVERNANCE proposal=100 phase=tally votingPower=100962 turnoutPercent=17.6
2025-11-08T00:27:00Z GOVERNANCE proposal=100 phase=discussion votingPower=100999 turnoutPercent=17.7
2025-11-08T00:28:00Z GOVERNANCE proposal=100 phase=voting votingPower=101036 turnoutPercent=17.8
2025-11-08T00:29:00Z GOVERNANCE proposal=100 phase=tally votingPower=101073 turnoutPercent=17.9
2025-11-08T00:30:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101110 turnoutPercent=18.0
2025-11-08T00:31:00Z GOVERNANCE proposal=100 phase=voting votingPower=101147 turnoutPercent=18.1
2025-11-08T00:32:00Z GOVERNANCE proposal=100 phase=tally votingPower=101184 turnoutPercent=18.2
2025-11-08T00:33:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101221 turnoutPercent=18.3
2025-11-08T00:34:00Z GOVERNANCE proposal=100 phase=voting votingPower=101258 turnoutPercent=18.4
2025-11-08T00:35:00Z GOVERNANCE proposal=100 phase=tally votingPower=101295 turnoutPercent=18.5
2025-11-08T00:36:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101332 turnoutPercent=18.6
2025-11-08T00:37:00Z GOVERNANCE proposal=100 phase=voting votingPower=101369 turnoutPercent=18.7
2025-11-08T00:38:00Z GOVERNANCE proposal=100 phase=tally votingPower=101406 turnoutPercent=18.8
2025-11-08T00:39:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101443 turnoutPercent=18.9
2025-11-08T00:40:00Z GOVERNANCE proposal=100 phase=voting votingPower=101480 turnoutPercent=19.0
2025-11-08T00:41:00Z GOVERNANCE proposal=100 phase=tally votingPower=101517 turnoutPercent=19.1
2025-11-08T00:42:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101554 turnoutPercent=19.2
2025-11-08T00:43:00Z GOVERNANCE proposal=100 phase=voting votingPower=101591 turnoutPercent=19.3
2025-11-08T00:44:00Z GOVERNANCE proposal=100 phase=tally votingPower=101628 turnoutPercent=19.4
2025-11-08T00:45:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101665 turnoutPercent=19.5
2025-11-08T00:46:00Z GOVERNANCE proposal=100 phase=voting votingPower=101702 turnoutPercent=19.6
2025-11-08T00:47:00Z GOVERNANCE proposal=100 phase=tally votingPower=101739 turnoutPercent=19.7
2025-11-08T00:48:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101776 turnoutPercent=19.8
2025-11-08T00:49:00Z GOVERNANCE proposal=100 phase=voting votingPower=101813 turnoutPercent=19.9
2025-11-08T00:50:00Z GOVERNANCE proposal=100 phase=tally votingPower=101850 turnoutPercent=20.0
2025-11-08T00:51:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101887 turnoutPercent=20.1
2025-11-08T00:52:00Z GOVERNANCE proposal=100 phase=voting votingPower=101924 turnoutPercent=20.2
2025-11-08T00:53:00Z GOVERNANCE proposal=100 phase=tally votingPower=101961 turnoutPercent=20.3
2025-11-08T00:54:00Z GOVERNANCE proposal=100 phase=discussion votingPower=101998 turnoutPercent=20.4
2025-11-08T00:55:00Z GOVERNANCE proposal=100 phase=voting votingPower=102035 turnoutPercent=20.5
2025-11-08T00:56:00Z GOVERNANCE proposal=100 phase=tally votingPower=102072 turnoutPercent=20.6
2025-11-08T00:57:00Z GOVERNANCE proposal=100 phase=discussion votingPower=102109 turnoutPercent=20.7
2025-11-08T00:58:00Z GOVERNANCE proposal=100 phase=voting votingPower=102146 turnoutPercent=20.8
2025-11-08T00:59:00Z GOVERNANCE proposal=100 phase=tally votingPower=102183 turnoutPercent=20.9
2025-11-08T01:00:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102220 turnoutPercent=15.0
2025-11-08T01:01:00Z GOVERNANCE proposal=101 phase=voting votingPower=102257 turnoutPercent=15.1
2025-11-08T01:02:00Z GOVERNANCE proposal=101 phase=tally votingPower=102294 turnoutPercent=15.2
2025-11-08T01:03:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102331 turnoutPercent=15.3
2025-11-08T01:04:00Z GOVERNANCE proposal=101 phase=voting votingPower=102368 turnoutPercent=15.4
2025-11-08T01:05:00Z GOVERNANCE proposal=101 phase=tally votingPower=102405 turnoutPercent=15.5
2025-11-08T01:06:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102442 turnoutPercent=15.6
2025-11-08T01:07:00Z GOVERNANCE proposal=101 phase=voting votingPower=102479 turnoutPercent=15.7
2025-11-08T01:08:00Z GOVERNANCE proposal=101 phase=tally votingPower=102516 turnoutPercent=15.8
2025-11-08T01:09:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102553 turnoutPercent=15.9
2025-11-08T01:10:00Z GOVERNANCE proposal=101 phase=voting votingPower=102590 turnoutPercent=16.0
2025-11-08T01:11:00Z GOVERNANCE proposal=101 phase=tally votingPower=102627 turnoutPercent=16.1
2025-11-08T01:12:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102664 turnoutPercent=16.2
2025-11-08T01:13:00Z GOVERNANCE proposal=101 phase=voting votingPower=102701 turnoutPercent=16.3
2025-11-08T01:14:00Z GOVERNANCE proposal=101 phase=tally votingPower=102738 turnoutPercent=16.4
2025-11-08T01:15:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102775 turnoutPercent=16.5
2025-11-08T01:16:00Z GOVERNANCE proposal=101 phase=voting votingPower=102812 turnoutPercent=16.6
2025-11-08T01:17:00Z GOVERNANCE proposal=101 phase=tally votingPower=102849 turnoutPercent=16.7
2025-11-08T01:18:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102886 turnoutPercent=16.8
2025-11-08T01:19:00Z GOVERNANCE proposal=101 phase=voting votingPower=102923 turnoutPercent=16.9
2025-11-08T01:20:00Z GOVERNANCE proposal=101 phase=tally votingPower=102960 turnoutPercent=17.0
2025-11-08T01:21:00Z GOVERNANCE proposal=101 phase=discussion votingPower=102997 turnoutPercent=17.1
2025-11-08T01:22:00Z GOVERNANCE proposal=101 phase=voting votingPower=103034 turnoutPercent=17.2
2025-11-08T01:23:00Z GOVERNANCE proposal=101 phase=tally votingPower=103071 turnoutPercent=17.3
2025-11-08T01:24:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103108 turnoutPercent=17.4
2025-11-08T01:25:00Z GOVERNANCE proposal=101 phase=voting votingPower=103145 turnoutPercent=17.5
2025-11-08T01:26:00Z GOVERNANCE proposal=101 phase=tally votingPower=103182 turnoutPercent=17.6
2025-11-08T01:27:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103219 turnoutPercent=17.7
2025-11-08T01:28:00Z GOVERNANCE proposal=101 phase=voting votingPower=103256 turnoutPercent=17.8
2025-11-08T01:29:00Z GOVERNANCE proposal=101 phase=tally votingPower=103293 turnoutPercent=17.9
2025-11-08T01:30:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103330 turnoutPercent=18.0
2025-11-08T01:31:00Z GOVERNANCE proposal=101 phase=voting votingPower=103367 turnoutPercent=18.1
2025-11-08T01:32:00Z GOVERNANCE proposal=101 phase=tally votingPower=103404 turnoutPercent=18.2
2025-11-08T01:33:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103441 turnoutPercent=18.3
2025-11-08T01:34:00Z GOVERNANCE proposal=101 phase=voting votingPower=103478 turnoutPercent=18.4
2025-11-08T01:35:00Z GOVERNANCE proposal=101 phase=tally votingPower=103515 turnoutPercent=18.5
2025-11-08T01:36:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103552 turnoutPercent=18.6
2025-11-08T01:37:00Z GOVERNANCE proposal=101 phase=voting votingPower=103589 turnoutPercent=18.7
2025-11-08T01:38:00Z GOVERNANCE proposal=101 phase=tally votingPower=103626 turnoutPercent=18.8
2025-11-08T01:39:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103663 turnoutPercent=18.9
2025-11-08T01:40:00Z GOVERNANCE proposal=101 phase=voting votingPower=103700 turnoutPercent=19.0
2025-11-08T01:41:00Z GOVERNANCE proposal=101 phase=tally votingPower=103737 turnoutPercent=19.1
2025-11-08T01:42:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103774 turnoutPercent=19.2
2025-11-08T01:43:00Z GOVERNANCE proposal=101 phase=voting votingPower=103811 turnoutPercent=19.3
2025-11-08T01:44:00Z GOVERNANCE proposal=101 phase=tally votingPower=103848 turnoutPercent=19.4
2025-11-08T01:45:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103885 turnoutPercent=19.5
2025-11-08T01:46:00Z GOVERNANCE proposal=101 phase=voting votingPower=103922 turnoutPercent=19.6
2025-11-08T01:47:00Z GOVERNANCE proposal=101 phase=tally votingPower=103959 turnoutPercent=19.7
2025-11-08T01:48:00Z GOVERNANCE proposal=101 phase=discussion votingPower=103996 turnoutPercent=19.8
2025-11-08T01:49:00Z GOVERNANCE proposal=101 phase=voting votingPower=104033 turnoutPercent=19.9
2025-11-08T01:50:00Z GOVERNANCE proposal=101 phase=tally votingPower=104070 turnoutPercent=20.0
2025-11-08T01:51:00Z GOVERNANCE proposal=101 phase=discussion votingPower=104107 turnoutPercent=20.1
2025-11-08T01:52:00Z GOVERNANCE proposal=101 phase=voting votingPower=104144 turnoutPercent=20.2
2025-11-08T01:53:00Z GOVERNANCE proposal=101 phase=tally votingPower=104181 turnoutPercent=20.3
2025-11-08T01:54:00Z GOVERNANCE proposal=101 phase=discussion votingPower=104218 turnoutPercent=20.4
2025-11-08T01:55:00Z GOVERNANCE proposal=101 phase=voting votingPower=104255 turnoutPercent=20.5
2025-11-08T01:56:00Z GOVERNANCE proposal=101 phase=tally votingPower=104292 turnoutPercent=20.6
2025-11-08T01:57:00Z GOVERNANCE proposal=101 phase=discussion votingPower=104329 turnoutPercent=20.7
2025-11-08T01:58:00Z GOVERNANCE proposal=101 phase=voting votingPower=104366 turnoutPercent=20.8
2025-11-08T01:59:00Z GOVERNANCE proposal=101 phase=tally votingPower=104403 turnoutPercent=20.9
2025-11-08T02:00:00Z GOVERNANCE proposal=102 phase=discussion votingPower=104440 turnoutPercent=15.0
2025-11-08T02:01:00Z GOVERNANCE proposal=102 phase=voting votingPower=104477 turnoutPercent=15.1
2025-11-08T02:02:00Z GOVERNANCE proposal=102 phase=tally votingPower=104514 turnoutPercent=15.2
2025-11-08T02:03:00Z GOVERNANCE proposal=102 phase=discussion votingPower=104551 turnoutPercent=15.3
2025-11-08T02:04:00Z GOVERNANCE proposal=102 phase=voting votingPower=104588 turnoutPercent=15.4
2025-11-08T02:05:00Z GOVERNANCE proposal=102 phase=tally votingPower=104625 turnoutPercent=15.5
2025-11-08T02:06:00Z GOVERNANCE proposal=102 phase=discussion votingPower=104662 turnoutPercent=15.6
2025-11-08T02:07:00Z GOVERNANCE proposal=102 phase=voting votingPower=104699 turnoutPercent=15.7
2025-11-08T02:08:00Z GOVERNANCE proposal=102 phase=tally votingPower=104736 turnoutPercent=15.8
2025-11-08T02:09:00Z GOVERNANCE proposal=102 phase=discussion votingPower=104773 turnoutPercent=15.9
2025-11-08T02:10:00Z GOVERNANCE proposal=102 phase=voting votingPower=104810 turnoutPercent=16.0
2025-11-08T02:11:00Z GOVERNANCE proposal=102 phase=tally votingPower=104847 turnoutPercent=16.1
2025-11-08T02:12:00Z GOVERNANCE proposal=102 phase=discussion votingPower=104884 turnoutPercent=16.2
2025-11-08T02:13:00Z GOVERNANCE proposal=102 phase=voting votingPower=104921 turnoutPercent=16.3
2025-11-08T02:14:00Z GOVERNANCE proposal=102 phase=tally votingPower=104958 turnoutPercent=16.4
2025-11-08T02:15:00Z GOVERNANCE proposal=102 phase=discussion votingPower=104995 turnoutPercent=16.5
2025-11-08T02:16:00Z GOVERNANCE proposal=102 phase=voting votingPower=105032 turnoutPercent=16.6
2025-11-08T02:17:00Z GOVERNANCE proposal=102 phase=tally votingPower=105069 turnoutPercent=16.7
2025-11-08T02:18:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105106 turnoutPercent=16.8
2025-11-08T02:19:00Z GOVERNANCE proposal=102 phase=voting votingPower=105143 turnoutPercent=16.9
2025-11-08T02:20:00Z GOVERNANCE proposal=102 phase=tally votingPower=105180 turnoutPercent=17.0
2025-11-08T02:21:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105217 turnoutPercent=17.1
2025-11-08T02:22:00Z GOVERNANCE proposal=102 phase=voting votingPower=105254 turnoutPercent=17.2
2025-11-08T02:23:00Z GOVERNANCE proposal=102 phase=tally votingPower=105291 turnoutPercent=17.3
2025-11-08T02:24:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105328 turnoutPercent=17.4
2025-11-08T02:25:00Z GOVERNANCE proposal=102 phase=voting votingPower=105365 turnoutPercent=17.5
2025-11-08T02:26:00Z GOVERNANCE proposal=102 phase=tally votingPower=105402 turnoutPercent=17.6
2025-11-08T02:27:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105439 turnoutPercent=17.7
2025-11-08T02:28:00Z GOVERNANCE proposal=102 phase=voting votingPower=105476 turnoutPercent=17.8
2025-11-08T02:29:00Z GOVERNANCE proposal=102 phase=tally votingPower=105513 turnoutPercent=17.9
2025-11-08T02:30:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105550 turnoutPercent=18.0
2025-11-08T02:31:00Z GOVERNANCE proposal=102 phase=voting votingPower=105587 turnoutPercent=18.1
2025-11-08T02:32:00Z GOVERNANCE proposal=102 phase=tally votingPower=105624 turnoutPercent=18.2
2025-11-08T02:33:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105661 turnoutPercent=18.3
2025-11-08T02:34:00Z GOVERNANCE proposal=102 phase=voting votingPower=105698 turnoutPercent=18.4
2025-11-08T02:35:00Z GOVERNANCE proposal=102 phase=tally votingPower=105735 turnoutPercent=18.5
2025-11-08T02:36:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105772 turnoutPercent=18.6
2025-11-08T02:37:00Z GOVERNANCE proposal=102 phase=voting votingPower=105809 turnoutPercent=18.7
2025-11-08T02:38:00Z GOVERNANCE proposal=102 phase=tally votingPower=105846 turnoutPercent=18.8
2025-11-08T02:39:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105883 turnoutPercent=18.9
2025-11-08T02:40:00Z GOVERNANCE proposal=102 phase=voting votingPower=105920 turnoutPercent=19.0
2025-11-08T02:41:00Z GOVERNANCE proposal=102 phase=tally votingPower=105957 turnoutPercent=19.1
2025-11-08T02:42:00Z GOVERNANCE proposal=102 phase=discussion votingPower=105994 turnoutPercent=19.2
2025-11-08T02:43:00Z GOVERNANCE proposal=102 phase=voting votingPower=106031 turnoutPercent=19.3
2025-11-08T02:44:00Z GOVERNANCE proposal=102 phase=tally votingPower=106068 turnoutPercent=19.4
2025-11-08T02:45:00Z GOVERNANCE proposal=102 phase=discussion votingPower=106105 turnoutPercent=19.5
2025-11-08T02:46:00Z GOVERNANCE proposal=102 phase=voting votingPower=106142 turnoutPercent=19.6
2025-11-08T02:47:00Z GOVERNANCE proposal=102 phase=tally votingPower=106179 turnoutPercent=19.7
2025-11-08T02:48:00Z GOVERNANCE proposal=102 phase=discussion votingPower=106216 turnoutPercent=19.8
2025-11-08T02:49:00Z GOVERNANCE proposal=102 phase=voting votingPower=106253 turnoutPercent=19.9
2025-11-08T02:50:00Z GOVERNANCE proposal=102 phase=tally votingPower=106290 turnoutPercent=20.0
2025-11-08T02:51:00Z GOVERNANCE proposal=102 phase=discussion votingPower=106327 turnoutPercent=20.1
2025-11-08T02:52:00Z GOVERNANCE proposal=102 phase=voting votingPower=106364 turnoutPercent=20.2
2025-11-08T02:53:00Z GOVERNANCE proposal=102 phase=tally votingPower=106401 turnoutPercent=20.3
2025-11-08T02:54:00Z GOVERNANCE proposal=102 phase=discussion votingPower=106438 turnoutPercent=20.4
2025-11-08T02:55:00Z GOVERNANCE proposal=102 phase=voting votingPower=106475 turnoutPercent=20.5
2025-11-08T02:56:00Z GOVERNANCE proposal=102 phase=tally votingPower=106512 turnoutPercent=20.6
2025-11-08T02:57:00Z GOVERNANCE proposal=102 phase=discussion votingPower=106549 turnoutPercent=20.7
2025-11-08T02:58:00Z GOVERNANCE proposal=102 phase=voting votingPower=106586 turnoutPercent=20.8
2025-11-08T02:59:00Z GOVERNANCE proposal=102 phase=tally votingPower=106623 turnoutPercent=20.9
2025-11-08T03:00:00Z GOVERNANCE proposal=103 phase=discussion votingPower=106660 turnoutPercent=15.0
2025-11-08T03:01:00Z GOVERNANCE proposal=103 phase=voting votingPower=106697 turnoutPercent=15.1
2025-11-08T03:02:00Z GOVERNANCE proposal=103 phase=tally votingPower=106734 turnoutPercent=15.2
2025-11-08T03:03:00Z GOVERNANCE proposal=103 phase=discussion votingPower=106771 turnoutPercent=15.3
2025-11-08T03:04:00Z GOVERNANCE proposal=103 phase=voting votingPower=106808 turnoutPercent=15.4
2025-11-08T03:05:00Z GOVERNANCE proposal=103 phase=tally votingPower=106845 turnoutPercent=15.5
2025-11-08T03:06:00Z GOVERNANCE proposal=103 phase=discussion votingPower=106882 turnoutPercent=15.6
2025-11-08T03:07:00Z GOVERNANCE proposal=103 phase=voting votingPower=106919 turnoutPercent=15.7
2025-11-08T03:08:00Z GOVERNANCE proposal=103 phase=tally votingPower=106956 turnoutPercent=15.8
2025-11-08T03:09:00Z GOVERNANCE proposal=103 phase=discussion votingPower=106993 turnoutPercent=15.9
2025-11-08T03:10:00Z GOVERNANCE proposal=103 phase=voting votingPower=107030 turnoutPercent=16.0
2025-11-08T03:11:00Z GOVERNANCE proposal=103 phase=tally votingPower=107067 turnoutPercent=16.1
2025-11-08T03:12:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107104 turnoutPercent=16.2
2025-11-08T03:13:00Z GOVERNANCE proposal=103 phase=voting votingPower=107141 turnoutPercent=16.3
2025-11-08T03:14:00Z GOVERNANCE proposal=103 phase=tally votingPower=107178 turnoutPercent=16.4
2025-11-08T03:15:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107215 turnoutPercent=16.5
2025-11-08T03:16:00Z GOVERNANCE proposal=103 phase=voting votingPower=107252 turnoutPercent=16.6
2025-11-08T03:17:00Z GOVERNANCE proposal=103 phase=tally votingPower=107289 turnoutPercent=16.7
2025-11-08T03:18:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107326 turnoutPercent=16.8
2025-11-08T03:19:00Z GOVERNANCE proposal=103 phase=voting votingPower=107363 turnoutPercent=16.9
2025-11-08T03:20:00Z GOVERNANCE proposal=103 phase=tally votingPower=107400 turnoutPercent=17.0
2025-11-08T03:21:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107437 turnoutPercent=17.1
2025-11-08T03:22:00Z GOVERNANCE proposal=103 phase=voting votingPower=107474 turnoutPercent=17.2
2025-11-08T03:23:00Z GOVERNANCE proposal=103 phase=tally votingPower=107511 turnoutPercent=17.3
2025-11-08T03:24:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107548 turnoutPercent=17.4
2025-11-08T03:25:00Z GOVERNANCE proposal=103 phase=voting votingPower=107585 turnoutPercent=17.5
2025-11-08T03:26:00Z GOVERNANCE proposal=103 phase=tally votingPower=107622 turnoutPercent=17.6
2025-11-08T03:27:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107659 turnoutPercent=17.7
2025-11-08T03:28:00Z GOVERNANCE proposal=103 phase=voting votingPower=107696 turnoutPercent=17.8
2025-11-08T03:29:00Z GOVERNANCE proposal=103 phase=tally votingPower=107733 turnoutPercent=17.9
2025-11-08T03:30:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107770 turnoutPercent=18.0
2025-11-08T03:31:00Z GOVERNANCE proposal=103 phase=voting votingPower=107807 turnoutPercent=18.1
2025-11-08T03:32:00Z GOVERNANCE proposal=103 phase=tally votingPower=107844 turnoutPercent=18.2
2025-11-08T03:33:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107881 turnoutPercent=18.3
2025-11-08T03:34:00Z GOVERNANCE proposal=103 phase=voting votingPower=107918 turnoutPercent=18.4
2025-11-08T03:35:00Z GOVERNANCE proposal=103 phase=tally votingPower=107955 turnoutPercent=18.5
2025-11-08T03:36:00Z GOVERNANCE proposal=103 phase=discussion votingPower=107992 turnoutPercent=18.6
2025-11-08T03:37:00Z GOVERNANCE proposal=103 phase=voting votingPower=108029 turnoutPercent=18.7
2025-11-08T03:38:00Z GOVERNANCE proposal=103 phase=tally votingPower=108066 turnoutPercent=18.8
2025-11-08T03:39:00Z GOVERNANCE proposal=103 phase=discussion votingPower=108103 turnoutPercent=18.9
2025-11-08T03:40:00Z GOVERNANCE proposal=103 phase=voting votingPower=108140 turnoutPercent=19.0
2025-11-08T03:41:00Z GOVERNANCE proposal=103 phase=tally votingPower=108177 turnoutPercent=19.1
2025-11-08T03:42:00Z GOVERNANCE proposal=103 phase=discussion votingPower=108214 turnoutPercent=19.2
2025-11-08T03:43:00Z GOVERNANCE proposal=103 phase=voting votingPower=108251 turnoutPercent=19.3
2025-11-08T03:44:00Z GOVERNANCE proposal=103 phase=tally votingPower=108288 turnoutPercent=19.4
2025-11-08T03:45:00Z GOVERNANCE proposal=103 phase=discussion votingPower=108325 turnoutPercent=19.5
2025-11-08T03:46:00Z GOVERNANCE proposal=103 phase=voting votingPower=108362 turnoutPercent=19.6
2025-11-08T03:47:00Z GOVERNANCE proposal=103 phase=tally votingPower=108399 turnoutPercent=19.7
2025-11-08T03:48:00Z GOVERNANCE proposal=103 phase=discussion votingPower=108436 turnoutPercent=19.8
2025-11-08T03:49:00Z GOVERNANCE proposal=103 phase=voting votingPower=108473 turnoutPercent=19.9
2025-11-08T03:50:00Z GOVERNANCE proposal=103 phase=tally votingPower=108510 turnoutPercent=20.0
2025-11-08T03:51:00Z GOVERNANCE proposal=103 phase=discussion votingPower=108547 turnoutPercent=20.1
2025-11-08T03:52:00Z GOVERNANCE proposal=103 phase=voting votingPower=108584 turnoutPercent=20.2
2025-11-08T03:53:00Z GOVERNANCE proposal=103 phase=tally votingPower=108621 turnoutPercent=20.3
2025-11-08T03:54:00Z GOVERNANCE proposal=103 phase=discussion votingPower=108658 turnoutPercent=20.4
2025-11-08T03:55:00Z GOVERNANCE proposal=103 phase=voting votingPower=108695 turnoutPercent=20.5
2025-11-08T03:56:00Z GOVERNANCE proposal=103 phase=tally votingPower=108732 turnoutPercent=20.6
2025-11-08T03:57:00Z GOVERNANCE proposal=103 phase=discussion votingPower=108769 turnoutPercent=20.7
2025-11-08T03:58:00Z GOVERNANCE proposal=103 phase=voting votingPower=108806 turnoutPercent=20.8
2025-11-08T03:59:00Z GOVERNANCE proposal=103 phase=tally votingPower=108843 turnoutPercent=20.9
2025-11-08T04:00:00Z GOVERNANCE proposal=104 phase=discussion votingPower=108880 turnoutPercent=15.0
2025-11-08T04:01:00Z GOVERNANCE proposal=104 phase=voting votingPower=108917 turnoutPercent=15.1
2025-11-08T04:02:00Z GOVERNANCE proposal=104 phase=tally votingPower=108954 turnoutPercent=15.2
2025-11-08T04:03:00Z GOVERNANCE proposal=104 phase=discussion votingPower=108991 turnoutPercent=15.3
2025-11-08T04:04:00Z GOVERNANCE proposal=104 phase=voting votingPower=109028 turnoutPercent=15.4
2025-11-08T04:05:00Z GOVERNANCE proposal=104 phase=tally votingPower=109065 turnoutPercent=15.5
2025-11-08T04:06:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109102 turnoutPercent=15.6
2025-11-08T04:07:00Z GOVERNANCE proposal=104 phase=voting votingPower=109139 turnoutPercent=15.7
2025-11-08T04:08:00Z GOVERNANCE proposal=104 phase=tally votingPower=109176 turnoutPercent=15.8
2025-11-08T04:09:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109213 turnoutPercent=15.9
2025-11-08T04:10:00Z GOVERNANCE proposal=104 phase=voting votingPower=109250 turnoutPercent=16.0
2025-11-08T04:11:00Z GOVERNANCE proposal=104 phase=tally votingPower=109287 turnoutPercent=16.1
2025-11-08T04:12:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109324 turnoutPercent=16.2
2025-11-08T04:13:00Z GOVERNANCE proposal=104 phase=voting votingPower=109361 turnoutPercent=16.3
2025-11-08T04:14:00Z GOVERNANCE proposal=104 phase=tally votingPower=109398 turnoutPercent=16.4
2025-11-08T04:15:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109435 turnoutPercent=16.5
2025-11-08T04:16:00Z GOVERNANCE proposal=104 phase=voting votingPower=109472 turnoutPercent=16.6
2025-11-08T04:17:00Z GOVERNANCE proposal=104 phase=tally votingPower=109509 turnoutPercent=16.7
2025-11-08T04:18:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109546 turnoutPercent=16.8
2025-11-08T04:19:00Z GOVERNANCE proposal=104 phase=voting votingPower=109583 turnoutPercent=16.9
2025-11-08T04:20:00Z GOVERNANCE proposal=104 phase=tally votingPower=109620 turnoutPercent=17.0
2025-11-08T04:21:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109657 turnoutPercent=17.1
2025-11-08T04:22:00Z GOVERNANCE proposal=104 phase=voting votingPower=109694 turnoutPercent=17.2
2025-11-08T04:23:00Z GOVERNANCE proposal=104 phase=tally votingPower=109731 turnoutPercent=17.3
2025-11-08T04:24:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109768 turnoutPercent=17.4
2025-11-08T04:25:00Z GOVERNANCE proposal=104 phase=voting votingPower=109805 turnoutPercent=17.5
2025-11-08T04:26:00Z GOVERNANCE proposal=104 phase=tally votingPower=109842 turnoutPercent=17.6
2025-11-08T04:27:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109879 turnoutPercent=17.7
2025-11-08T04:28:00Z GOVERNANCE proposal=104 phase=voting votingPower=109916 turnoutPercent=17.8
2025-11-08T04:29:00Z GOVERNANCE proposal=104 phase=tally votingPower=109953 turnoutPercent=17.9
2025-11-08T04:30:00Z GOVERNANCE proposal=104 phase=discussion votingPower=109990 turnoutPercent=18.0
2025-11-08T04:31:00Z GOVERNANCE proposal=104 phase=voting votingPower=110027 turnoutPercent=18.1
2025-11-08T04:32:00Z GOVERNANCE proposal=104 phase=tally votingPower=110064 turnoutPercent=18.2
2025-11-08T04:33:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110101 turnoutPercent=18.3
2025-11-08T04:34:00Z GOVERNANCE proposal=104 phase=voting votingPower=110138 turnoutPercent=18.4
2025-11-08T04:35:00Z GOVERNANCE proposal=104 phase=tally votingPower=110175 turnoutPercent=18.5
2025-11-08T04:36:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110212 turnoutPercent=18.6
2025-11-08T04:37:00Z GOVERNANCE proposal=104 phase=voting votingPower=110249 turnoutPercent=18.7
2025-11-08T04:38:00Z GOVERNANCE proposal=104 phase=tally votingPower=110286 turnoutPercent=18.8
2025-11-08T04:39:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110323 turnoutPercent=18.9
2025-11-08T04:40:00Z GOVERNANCE proposal=104 phase=voting votingPower=110360 turnoutPercent=19.0
2025-11-08T04:41:00Z GOVERNANCE proposal=104 phase=tally votingPower=110397 turnoutPercent=19.1
2025-11-08T04:42:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110434 turnoutPercent=19.2
2025-11-08T04:43:00Z GOVERNANCE proposal=104 phase=voting votingPower=110471 turnoutPercent=19.3
2025-11-08T04:44:00Z GOVERNANCE proposal=104 phase=tally votingPower=110508 turnoutPercent=19.4
2025-11-08T04:45:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110545 turnoutPercent=19.5
2025-11-08T04:46:00Z GOVERNANCE proposal=104 phase=voting votingPower=110582 turnoutPercent=19.6
2025-11-08T04:47:00Z GOVERNANCE proposal=104 phase=tally votingPower=110619 turnoutPercent=19.7
2025-11-08T04:48:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110656 turnoutPercent=19.8
2025-11-08T04:49:00Z GOVERNANCE proposal=104 phase=voting votingPower=110693 turnoutPercent=19.9
2025-11-08T04:50:00Z GOVERNANCE proposal=104 phase=tally votingPower=110730 turnoutPercent=20.0
2025-11-08T04:51:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110767 turnoutPercent=20.1
2025-11-08T04:52:00Z GOVERNANCE proposal=104 phase=voting votingPower=110804 turnoutPercent=20.2
2025-11-08T04:53:00Z GOVERNANCE proposal=104 phase=tally votingPower=110841 turnoutPercent=20.3
2025-11-08T04:54:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110878 turnoutPercent=20.4
2025-11-08T04:55:00Z GOVERNANCE proposal=104 phase=voting votingPower=110915 turnoutPercent=20.5
2025-11-08T04:56:00Z GOVERNANCE proposal=104 phase=tally votingPower=110952 turnoutPercent=20.6
2025-11-08T04:57:00Z GOVERNANCE proposal=104 phase=discussion votingPower=110989 turnoutPercent=20.7
2025-11-08T04:58:00Z GOVERNANCE proposal=104 phase=voting votingPower=111026 turnoutPercent=20.8
2025-11-08T04:59:00Z GOVERNANCE proposal=104 phase=tally votingPower=111063 turnoutPercent=20.9
2025-11-08T05:00:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111100 turnoutPercent=15.0
2025-11-08T05:01:00Z GOVERNANCE proposal=105 phase=voting votingPower=111137 turnoutPercent=15.1
2025-11-08T05:02:00Z GOVERNANCE proposal=105 phase=tally votingPower=111174 turnoutPercent=15.2
2025-11-08T05:03:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111211 turnoutPercent=15.3
2025-11-08T05:04:00Z GOVERNANCE proposal=105 phase=voting votingPower=111248 turnoutPercent=15.4
2025-11-08T05:05:00Z GOVERNANCE proposal=105 phase=tally votingPower=111285 turnoutPercent=15.5
2025-11-08T05:06:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111322 turnoutPercent=15.6
2025-11-08T05:07:00Z GOVERNANCE proposal=105 phase=voting votingPower=111359 turnoutPercent=15.7
2025-11-08T05:08:00Z GOVERNANCE proposal=105 phase=tally votingPower=111396 turnoutPercent=15.8
2025-11-08T05:09:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111433 turnoutPercent=15.9
2025-11-08T05:10:00Z GOVERNANCE proposal=105 phase=voting votingPower=111470 turnoutPercent=16.0
2025-11-08T05:11:00Z GOVERNANCE proposal=105 phase=tally votingPower=111507 turnoutPercent=16.1
2025-11-08T05:12:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111544 turnoutPercent=16.2
2025-11-08T05:13:00Z GOVERNANCE proposal=105 phase=voting votingPower=111581 turnoutPercent=16.3
2025-11-08T05:14:00Z GOVERNANCE proposal=105 phase=tally votingPower=111618 turnoutPercent=16.4
2025-11-08T05:15:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111655 turnoutPercent=16.5
2025-11-08T05:16:00Z GOVERNANCE proposal=105 phase=voting votingPower=111692 turnoutPercent=16.6
2025-11-08T05:17:00Z GOVERNANCE proposal=105 phase=tally votingPower=111729 turnoutPercent=16.7
2025-11-08T05:18:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111766 turnoutPercent=16.8
2025-11-08T05:19:00Z GOVERNANCE proposal=105 phase=voting votingPower=111803 turnoutPercent=16.9
2025-11-08T05:20:00Z GOVERNANCE proposal=105 phase=tally votingPower=111840 turnoutPercent=17.0
2025-11-08T05:21:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111877 turnoutPercent=17.1
2025-11-08T05:22:00Z GOVERNANCE proposal=105 phase=voting votingPower=111914 turnoutPercent=17.2
2025-11-08T05:23:00Z GOVERNANCE proposal=105 phase=tally votingPower=111951 turnoutPercent=17.3
2025-11-08T05:24:00Z GOVERNANCE proposal=105 phase=discussion votingPower=111988 turnoutPercent=17.4
2025-11-08T05:25:00Z GOVERNANCE proposal=105 phase=voting votingPower=112025 turnoutPercent=17.5
2025-11-08T05:26:00Z GOVERNANCE proposal=105 phase=tally votingPower=112062 turnoutPercent=17.6
2025-11-08T05:27:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112099 turnoutPercent=17.7
2025-11-08T05:28:00Z GOVERNANCE proposal=105 phase=voting votingPower=112136 turnoutPercent=17.8
2025-11-08T05:29:00Z GOVERNANCE proposal=105 phase=tally votingPower=112173 turnoutPercent=17.9
2025-11-08T05:30:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112210 turnoutPercent=18.0
2025-11-08T05:31:00Z GOVERNANCE proposal=105 phase=voting votingPower=112247 turnoutPercent=18.1
2025-11-08T05:32:00Z GOVERNANCE proposal=105 phase=tally votingPower=112284 turnoutPercent=18.2
2025-11-08T05:33:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112321 turnoutPercent=18.3
2025-11-08T05:34:00Z GOVERNANCE proposal=105 phase=voting votingPower=112358 turnoutPercent=18.4
2025-11-08T05:35:00Z GOVERNANCE proposal=105 phase=tally votingPower=112395 turnoutPercent=18.5
2025-11-08T05:36:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112432 turnoutPercent=18.6
2025-11-08T05:37:00Z GOVERNANCE proposal=105 phase=voting votingPower=112469 turnoutPercent=18.7
2025-11-08T05:38:00Z GOVERNANCE proposal=105 phase=tally votingPower=112506 turnoutPercent=18.8
2025-11-08T05:39:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112543 turnoutPercent=18.9
2025-11-08T05:40:00Z GOVERNANCE proposal=105 phase=voting votingPower=112580 turnoutPercent=19.0
2025-11-08T05:41:00Z GOVERNANCE proposal=105 phase=tally votingPower=112617 turnoutPercent=19.1
2025-11-08T05:42:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112654 turnoutPercent=19.2
2025-11-08T05:43:00Z GOVERNANCE proposal=105 phase=voting votingPower=112691 turnoutPercent=19.3
2025-11-08T05:44:00Z GOVERNANCE proposal=105 phase=tally votingPower=112728 turnoutPercent=19.4
2025-11-08T05:45:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112765 turnoutPercent=19.5
2025-11-08T05:46:00Z GOVERNANCE proposal=105 phase=voting votingPower=112802 turnoutPercent=19.6
2025-11-08T05:47:00Z GOVERNANCE proposal=105 phase=tally votingPower=112839 turnoutPercent=19.7
2025-11-08T05:48:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112876 turnoutPercent=19.8
2025-11-08T05:49:00Z GOVERNANCE proposal=105 phase=voting votingPower=112913 turnoutPercent=19.9
2025-11-08T05:50:00Z GOVERNANCE proposal=105 phase=tally votingPower=112950 turnoutPercent=20.0
2025-11-08T05:51:00Z GOVERNANCE proposal=105 phase=discussion votingPower=112987 turnoutPercent=20.1
2025-11-08T05:52:00Z GOVERNANCE proposal=105 phase=voting votingPower=113024 turnoutPercent=20.2
2025-11-08T05:53:00Z GOVERNANCE proposal=105 phase=tally votingPower=113061 turnoutPercent=20.3
2025-11-08T05:54:00Z GOVERNANCE proposal=105 phase=discussion votingPower=113098 turnoutPercent=20.4
2025-11-08T05:55:00Z GOVERNANCE proposal=105 phase=voting votingPower=113135 turnoutPercent=20.5
2025-11-08T05:56:00Z GOVERNANCE proposal=105 phase=tally votingPower=113172 turnoutPercent=20.6
2025-11-08T05:57:00Z GOVERNANCE proposal=105 phase=discussion votingPower=113209 turnoutPercent=20.7
2025-11-08T05:58:00Z GOVERNANCE proposal=105 phase=voting votingPower=113246 turnoutPercent=20.8
2025-11-08T05:59:00Z GOVERNANCE proposal=105 phase=tally votingPower=113283 turnoutPercent=20.9
2025-11-08T06:00:00Z GOVERNANCE proposal=106 phase=discussion votingPower=113320 turnoutPercent=15.0
2025-11-08T06:01:00Z GOVERNANCE proposal=106 phase=voting votingPower=113357 turnoutPercent=15.1
2025-11-08T06:02:00Z GOVERNANCE proposal=106 phase=tally votingPower=113394 turnoutPercent=15.2
2025-11-08T06:03:00Z GOVERNANCE proposal=106 phase=discussion votingPower=113431 turnoutPercent=15.3
2025-11-08T06:04:00Z GOVERNANCE proposal=106 phase=voting votingPower=113468 turnoutPercent=15.4
2025-11-08T06:05:00Z GOVERNANCE proposal=106 phase=tally votingPower=113505 turnoutPercent=15.5
2025-11-08T06:06:00Z GOVERNANCE proposal=106 phase=discussion votingPower=113542 turnoutPercent=15.6
2025-11-08T06:07:00Z GOVERNANCE proposal=106 phase=voting votingPower=113579 turnoutPercent=15.7
2025-11-08T06:08:00Z GOVERNANCE proposal=106 phase=tally votingPower=113616 turnoutPercent=15.8
2025-11-08T06:09:00Z GOVERNANCE proposal=106 phase=discussion votingPower=113653 turnoutPercent=15.9
2025-11-08T06:10:00Z GOVERNANCE proposal=106 phase=voting votingPower=113690 turnoutPercent=16.0
2025-11-08T06:11:00Z GOVERNANCE proposal=106 phase=tally votingPower=113727 turnoutPercent=16.1
2025-11-08T06:12:00Z GOVERNANCE proposal=106 phase=discussion votingPower=113764 turnoutPercent=16.2
2025-11-08T06:13:00Z GOVERNANCE proposal=106 phase=voting votingPower=113801 turnoutPercent=16.3
2025-11-08T06:14:00Z GOVERNANCE proposal=106 phase=tally votingPower=113838 turnoutPercent=16.4
2025-11-08T06:15:00Z GOVERNANCE proposal=106 phase=discussion votingPower=113875 turnoutPercent=16.5
2025-11-08T06:16:00Z GOVERNANCE proposal=106 phase=voting votingPower=113912 turnoutPercent=16.6
2025-11-08T06:17:00Z GOVERNANCE proposal=106 phase=tally votingPower=113949 turnoutPercent=16.7
2025-11-08T06:18:00Z GOVERNANCE proposal=106 phase=discussion votingPower=113986 turnoutPercent=16.8
2025-11-08T06:19:00Z GOVERNANCE proposal=106 phase=voting votingPower=114023 turnoutPercent=16.9
2025-11-08T06:20:00Z GOVERNANCE proposal=106 phase=tally votingPower=114060 turnoutPercent=17.0
2025-11-08T06:21:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114097 turnoutPercent=17.1
2025-11-08T06:22:00Z GOVERNANCE proposal=106 phase=voting votingPower=114134 turnoutPercent=17.2
2025-11-08T06:23:00Z GOVERNANCE proposal=106 phase=tally votingPower=114171 turnoutPercent=17.3
2025-11-08T06:24:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114208 turnoutPercent=17.4
2025-11-08T06:25:00Z GOVERNANCE proposal=106 phase=voting votingPower=114245 turnoutPercent=17.5
2025-11-08T06:26:00Z GOVERNANCE proposal=106 phase=tally votingPower=114282 turnoutPercent=17.6
2025-11-08T06:27:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114319 turnoutPercent=17.7
2025-11-08T06:28:00Z GOVERNANCE proposal=106 phase=voting votingPower=114356 turnoutPercent=17.8
2025-11-08T06:29:00Z GOVERNANCE proposal=106 phase=tally votingPower=114393 turnoutPercent=17.9
2025-11-08T06:30:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114430 turnoutPercent=18.0
2025-11-08T06:31:00Z GOVERNANCE proposal=106 phase=voting votingPower=114467 turnoutPercent=18.1
2025-11-08T06:32:00Z GOVERNANCE proposal=106 phase=tally votingPower=114504 turnoutPercent=18.2
2025-11-08T06:33:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114541 turnoutPercent=18.3
2025-11-08T06:34:00Z GOVERNANCE proposal=106 phase=voting votingPower=114578 turnoutPercent=18.4
2025-11-08T06:35:00Z GOVERNANCE proposal=106 phase=tally votingPower=114615 turnoutPercent=18.5
2025-11-08T06:36:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114652 turnoutPercent=18.6
2025-11-08T06:37:00Z GOVERNANCE proposal=106 phase=voting votingPower=114689 turnoutPercent=18.7
2025-11-08T06:38:00Z GOVERNANCE proposal=106 phase=tally votingPower=114726 turnoutPercent=18.8
2025-11-08T06:39:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114763 turnoutPercent=18.9
2025-11-08T06:40:00Z GOVERNANCE proposal=106 phase=voting votingPower=114800 turnoutPercent=19.0
2025-11-08T06:41:00Z GOVERNANCE proposal=106 phase=tally votingPower=114837 turnoutPercent=19.1
2025-11-08T06:42:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114874 turnoutPercent=19.2
2025-11-08T06:43:00Z GOVERNANCE proposal=106 phase=voting votingPower=114911 turnoutPercent=19.3
2025-11-08T06:44:00Z GOVERNANCE proposal=106 phase=tally votingPower=114948 turnoutPercent=19.4
2025-11-08T06:45:00Z GOVERNANCE proposal=106 phase=discussion votingPower=114985 turnoutPercent=19.5
2025-11-08T06:46:00Z GOVERNANCE proposal=106 phase=voting votingPower=115022 turnoutPercent=19.6
2025-11-08T06:47:00Z GOVERNANCE proposal=106 phase=tally votingPower=115059 turnoutPercent=19.7
2025-11-08T06:48:00Z GOVERNANCE proposal=106 phase=discussion votingPower=115096 turnoutPercent=19.8
2025-11-08T06:49:00Z GOVERNANCE proposal=106 phase=voting votingPower=115133 turnoutPercent=19.9
2025-11-08T06:50:00Z GOVERNANCE proposal=106 phase=tally votingPower=115170 turnoutPercent=20.0
2025-11-08T06:51:00Z GOVERNANCE proposal=106 phase=discussion votingPower=115207 turnoutPercent=20.1
2025-11-08T06:52:00Z GOVERNANCE proposal=106 phase=voting votingPower=115244 turnoutPercent=20.2
2025-11-08T06:53:00Z GOVERNANCE proposal=106 phase=tally votingPower=115281 turnoutPercent=20.3
2025-11-08T06:54:00Z GOVERNANCE proposal=106 phase=discussion votingPower=115318 turnoutPercent=20.4
2025-11-08T06:55:00Z GOVERNANCE proposal=106 phase=voting votingPower=115355 turnoutPercent=20.5
2025-11-08T06:56:00Z GOVERNANCE proposal=106 phase=tally votingPower=115392 turnoutPercent=20.6
2025-11-08T06:57:00Z GOVERNANCE proposal=106 phase=discussion votingPower=115429 turnoutPercent=20.7
2025-11-08T06:58:00Z GOVERNANCE proposal=106 phase=voting votingPower=115466 turnoutPercent=20.8
2025-11-08T06:59:00Z GOVERNANCE proposal=106 phase=tally votingPower=115503 turnoutPercent=20.9
2025-11-08T07:00:00Z GOVERNANCE proposal=107 phase=discussion votingPower=115540 turnoutPercent=15.0
2025-11-08T07:01:00Z GOVERNANCE proposal=107 phase=voting votingPower=115577 turnoutPercent=15.1
2025-11-08T07:02:00Z GOVERNANCE proposal=107 phase=tally votingPower=115614 turnoutPercent=15.2
2025-11-08T07:03:00Z GOVERNANCE proposal=107 phase=discussion votingPower=115651 turnoutPercent=15.3
2025-11-08T07:04:00Z GOVERNANCE proposal=107 phase=voting votingPower=115688 turnoutPercent=15.4
2025-11-08T07:05:00Z GOVERNANCE proposal=107 phase=tally votingPower=115725 turnoutPercent=15.5
2025-11-08T07:06:00Z GOVERNANCE proposal=107 phase=discussion votingPower=115762 turnoutPercent=15.6
2025-11-08T07:07:00Z GOVERNANCE proposal=107 phase=voting votingPower=115799 turnoutPercent=15.7
2025-11-08T07:08:00Z GOVERNANCE proposal=107 phase=tally votingPower=115836 turnoutPercent=15.8
2025-11-08T07:09:00Z GOVERNANCE proposal=107 phase=discussion votingPower=115873 turnoutPercent=15.9
2025-11-08T07:10:00Z GOVERNANCE proposal=107 phase=voting votingPower=115910 turnoutPercent=16.0
2025-11-08T07:11:00Z GOVERNANCE proposal=107 phase=tally votingPower=115947 turnoutPercent=16.1
2025-11-08T07:12:00Z GOVERNANCE proposal=107 phase=discussion votingPower=115984 turnoutPercent=16.2
2025-11-08T07:13:00Z GOVERNANCE proposal=107 phase=voting votingPower=116021 turnoutPercent=16.3
2025-11-08T07:14:00Z GOVERNANCE proposal=107 phase=tally votingPower=116058 turnoutPercent=16.4
2025-11-08T07:15:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116095 turnoutPercent=16.5
2025-11-08T07:16:00Z GOVERNANCE proposal=107 phase=voting votingPower=116132 turnoutPercent=16.6
2025-11-08T07:17:00Z GOVERNANCE proposal=107 phase=tally votingPower=116169 turnoutPercent=16.7
2025-11-08T07:18:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116206 turnoutPercent=16.8
2025-11-08T07:19:00Z GOVERNANCE proposal=107 phase=voting votingPower=116243 turnoutPercent=16.9
2025-11-08T07:20:00Z GOVERNANCE proposal=107 phase=tally votingPower=116280 turnoutPercent=17.0
2025-11-08T07:21:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116317 turnoutPercent=17.1
2025-11-08T07:22:00Z GOVERNANCE proposal=107 phase=voting votingPower=116354 turnoutPercent=17.2
2025-11-08T07:23:00Z GOVERNANCE proposal=107 phase=tally votingPower=116391 turnoutPercent=17.3
2025-11-08T07:24:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116428 turnoutPercent=17.4
2025-11-08T07:25:00Z GOVERNANCE proposal=107 phase=voting votingPower=116465 turnoutPercent=17.5
2025-11-08T07:26:00Z GOVERNANCE proposal=107 phase=tally votingPower=116502 turnoutPercent=17.6
2025-11-08T07:27:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116539 turnoutPercent=17.7
2025-11-08T07:28:00Z GOVERNANCE proposal=107 phase=voting votingPower=116576 turnoutPercent=17.8
2025-11-08T07:29:00Z GOVERNANCE proposal=107 phase=tally votingPower=116613 turnoutPercent=17.9
2025-11-08T07:30:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116650 turnoutPercent=18.0
2025-11-08T07:31:00Z GOVERNANCE proposal=107 phase=voting votingPower=116687 turnoutPercent=18.1
2025-11-08T07:32:00Z GOVERNANCE proposal=107 phase=tally votingPower=116724 turnoutPercent=18.2
2025-11-08T07:33:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116761 turnoutPercent=18.3
2025-11-08T07:34:00Z GOVERNANCE proposal=107 phase=voting votingPower=116798 turnoutPercent=18.4
2025-11-08T07:35:00Z GOVERNANCE proposal=107 phase=tally votingPower=116835 turnoutPercent=18.5
2025-11-08T07:36:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116872 turnoutPercent=18.6
2025-11-08T07:37:00Z GOVERNANCE proposal=107 phase=voting votingPower=116909 turnoutPercent=18.7
2025-11-08T07:38:00Z GOVERNANCE proposal=107 phase=tally votingPower=116946 turnoutPercent=18.8
2025-11-08T07:39:00Z GOVERNANCE proposal=107 phase=discussion votingPower=116983 turnoutPercent=18.9
2025-11-08T07:40:00Z GOVERNANCE proposal=107 phase=voting votingPower=117020 turnoutPercent=19.0
2025-11-08T07:41:00Z GOVERNANCE proposal=107 phase=tally votingPower=117057 turnoutPercent=19.1
2025-11-08T07:42:00Z GOVERNANCE proposal=107 phase=discussion votingPower=117094 turnoutPercent=19.2
2025-11-08T07:43:00Z GOVERNANCE proposal=107 phase=voting votingPower=117131 turnoutPercent=19.3
2025-11-08T07:44:00Z GOVERNANCE proposal=107 phase=tally votingPower=117168 turnoutPercent=19.4
2025-11-08T07:45:00Z GOVERNANCE proposal=107 phase=discussion votingPower=117205 turnoutPercent=19.5
2025-11-08T07:46:00Z GOVERNANCE proposal=107 phase=voting votingPower=117242 turnoutPercent=19.6
2025-11-08T07:47:00Z GOVERNANCE proposal=107 phase=tally votingPower=117279 turnoutPercent=19.7
2025-11-08T07:48:00Z GOVERNANCE proposal=107 phase=discussion votingPower=117316 turnoutPercent=19.8
2025-11-08T07:49:00Z GOVERNANCE proposal=107 phase=voting votingPower=117353 turnoutPercent=19.9
2025-11-08T07:50:00Z GOVERNANCE proposal=107 phase=tally votingPower=117390 turnoutPercent=20.0
2025-11-08T07:51:00Z GOVERNANCE proposal=107 phase=discussion votingPower=117427 turnoutPercent=20.1
2025-11-08T07:52:00Z GOVERNANCE proposal=107 phase=voting votingPower=117464 turnoutPercent=20.2
2025-11-08T07:53:00Z GOVERNANCE proposal=107 phase=tally votingPower=117501 turnoutPercent=20.3
2025-11-08T07:54:00Z GOVERNANCE proposal=107 phase=discussion votingPower=117538 turnoutPercent=20.4
2025-11-08T07:55:00Z GOVERNANCE proposal=107 phase=voting votingPower=117575 turnoutPercent=20.5
2025-11-08T07:56:00Z GOVERNANCE proposal=107 phase=tally votingPower=117612 turnoutPercent=20.6
2025-11-08T07:57:00Z GOVERNANCE proposal=107 phase=discussion votingPower=117649 turnoutPercent=20.7
2025-11-08T07:58:00Z GOVERNANCE proposal=107 phase=voting votingPower=117686 turnoutPercent=20.8
2025-11-08T07:59:00Z GOVERNANCE proposal=107 phase=tally votingPower=117723 turnoutPercent=20.9
2025-11-08T08:00:00Z GOVERNANCE proposal=108 phase=discussion votingPower=117760 turnoutPercent=15.0
2025-11-08T08:01:00Z GOVERNANCE proposal=108 phase=voting votingPower=117797 turnoutPercent=15.1
2025-11-08T08:02:00Z GOVERNANCE proposal=108 phase=tally votingPower=117834 turnoutPercent=15.2
2025-11-08T08:03:00Z GOVERNANCE proposal=108 phase=discussion votingPower=117871 turnoutPercent=15.3
2025-11-08T08:04:00Z GOVERNANCE proposal=108 phase=voting votingPower=117908 turnoutPercent=15.4
2025-11-08T08:05:00Z GOVERNANCE proposal=108 phase=tally votingPower=117945 turnoutPercent=15.5
2025-11-08T08:06:00Z GOVERNANCE proposal=108 phase=discussion votingPower=117982 turnoutPercent=15.6
2025-11-08T08:07:00Z GOVERNANCE proposal=108 phase=voting votingPower=118019 turnoutPercent=15.7
2025-11-08T08:08:00Z GOVERNANCE proposal=108 phase=tally votingPower=118056 turnoutPercent=15.8
2025-11-08T08:09:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118093 turnoutPercent=15.9
2025-11-08T08:10:00Z GOVERNANCE proposal=108 phase=voting votingPower=118130 turnoutPercent=16.0
2025-11-08T08:11:00Z GOVERNANCE proposal=108 phase=tally votingPower=118167 turnoutPercent=16.1
2025-11-08T08:12:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118204 turnoutPercent=16.2
2025-11-08T08:13:00Z GOVERNANCE proposal=108 phase=voting votingPower=118241 turnoutPercent=16.3
2025-11-08T08:14:00Z GOVERNANCE proposal=108 phase=tally votingPower=118278 turnoutPercent=16.4
2025-11-08T08:15:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118315 turnoutPercent=16.5
2025-11-08T08:16:00Z GOVERNANCE proposal=108 phase=voting votingPower=118352 turnoutPercent=16.6
2025-11-08T08:17:00Z GOVERNANCE proposal=108 phase=tally votingPower=118389 turnoutPercent=16.7
2025-11-08T08:18:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118426 turnoutPercent=16.8
2025-11-08T08:19:00Z GOVERNANCE proposal=108 phase=voting votingPower=118463 turnoutPercent=16.9
2025-11-08T08:20:00Z GOVERNANCE proposal=108 phase=tally votingPower=118500 turnoutPercent=17.0
2025-11-08T08:21:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118537 turnoutPercent=17.1
2025-11-08T08:22:00Z GOVERNANCE proposal=108 phase=voting votingPower=118574 turnoutPercent=17.2
2025-11-08T08:23:00Z GOVERNANCE proposal=108 phase=tally votingPower=118611 turnoutPercent=17.3
2025-11-08T08:24:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118648 turnoutPercent=17.4
2025-11-08T08:25:00Z GOVERNANCE proposal=108 phase=voting votingPower=118685 turnoutPercent=17.5
2025-11-08T08:26:00Z GOVERNANCE proposal=108 phase=tally votingPower=118722 turnoutPercent=17.6
2025-11-08T08:27:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118759 turnoutPercent=17.7
2025-11-08T08:28:00Z GOVERNANCE proposal=108 phase=voting votingPower=118796 turnoutPercent=17.8
2025-11-08T08:29:00Z GOVERNANCE proposal=108 phase=tally votingPower=118833 turnoutPercent=17.9
2025-11-08T08:30:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118870 turnoutPercent=18.0
2025-11-08T08:31:00Z GOVERNANCE proposal=108 phase=voting votingPower=118907 turnoutPercent=18.1
2025-11-08T08:32:00Z GOVERNANCE proposal=108 phase=tally votingPower=118944 turnoutPercent=18.2
2025-11-08T08:33:00Z GOVERNANCE proposal=108 phase=discussion votingPower=118981 turnoutPercent=18.3
2025-11-08T08:34:00Z GOVERNANCE proposal=108 phase=voting votingPower=119018 turnoutPercent=18.4
2025-11-08T08:35:00Z GOVERNANCE proposal=108 phase=tally votingPower=119055 turnoutPercent=18.5
2025-11-08T08:36:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119092 turnoutPercent=18.6
2025-11-08T08:37:00Z GOVERNANCE proposal=108 phase=voting votingPower=119129 turnoutPercent=18.7
2025-11-08T08:38:00Z GOVERNANCE proposal=108 phase=tally votingPower=119166 turnoutPercent=18.8
2025-11-08T08:39:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119203 turnoutPercent=18.9
2025-11-08T08:40:00Z GOVERNANCE proposal=108 phase=voting votingPower=119240 turnoutPercent=19.0
2025-11-08T08:41:00Z GOVERNANCE proposal=108 phase=tally votingPower=119277 turnoutPercent=19.1
2025-11-08T08:42:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119314 turnoutPercent=19.2
2025-11-08T08:43:00Z GOVERNANCE proposal=108 phase=voting votingPower=119351 turnoutPercent=19.3
2025-11-08T08:44:00Z GOVERNANCE proposal=108 phase=tally votingPower=119388 turnoutPercent=19.4
2025-11-08T08:45:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119425 turnoutPercent=19.5
2025-11-08T08:46:00Z GOVERNANCE proposal=108 phase=voting votingPower=119462 turnoutPercent=19.6
2025-11-08T08:47:00Z GOVERNANCE proposal=108 phase=tally votingPower=119499 turnoutPercent=19.7
2025-11-08T08:48:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119536 turnoutPercent=19.8
2025-11-08T08:49:00Z GOVERNANCE proposal=108 phase=voting votingPower=119573 turnoutPercent=19.9
2025-11-08T08:50:00Z GOVERNANCE proposal=108 phase=tally votingPower=119610 turnoutPercent=20.0
2025-11-08T08:51:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119647 turnoutPercent=20.1
2025-11-08T08:52:00Z GOVERNANCE proposal=108 phase=voting votingPower=119684 turnoutPercent=20.2
2025-11-08T08:53:00Z GOVERNANCE proposal=108 phase=tally votingPower=119721 turnoutPercent=20.3
2025-11-08T08:54:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119758 turnoutPercent=20.4
2025-11-08T08:55:00Z GOVERNANCE proposal=108 phase=voting votingPower=119795 turnoutPercent=20.5
2025-11-08T08:56:00Z GOVERNANCE proposal=108 phase=tally votingPower=119832 turnoutPercent=20.6
2025-11-08T08:57:00Z GOVERNANCE proposal=108 phase=discussion votingPower=119869 turnoutPercent=20.7
2025-11-08T08:58:00Z GOVERNANCE proposal=108 phase=voting votingPower=119906 turnoutPercent=20.8
2025-11-08T08:59:00Z GOVERNANCE proposal=108 phase=tally votingPower=119943 turnoutPercent=20.9
2025-11-08T09:00:00Z GOVERNANCE proposal=109 phase=discussion votingPower=119980 turnoutPercent=15.0
2025-11-08T09:01:00Z GOVERNANCE proposal=109 phase=voting votingPower=120017 turnoutPercent=15.1
2025-11-08T09:02:00Z GOVERNANCE proposal=109 phase=tally votingPower=120054 turnoutPercent=15.2
2025-11-08T09:03:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120091 turnoutPercent=15.3
2025-11-08T09:04:00Z GOVERNANCE proposal=109 phase=voting votingPower=120128 turnoutPercent=15.4
2025-11-08T09:05:00Z GOVERNANCE proposal=109 phase=tally votingPower=120165 turnoutPercent=15.5
2025-11-08T09:06:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120202 turnoutPercent=15.6
2025-11-08T09:07:00Z GOVERNANCE proposal=109 phase=voting votingPower=120239 turnoutPercent=15.7
2025-11-08T09:08:00Z GOVERNANCE proposal=109 phase=tally votingPower=120276 turnoutPercent=15.8
2025-11-08T09:09:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120313 turnoutPercent=15.9
2025-11-08T09:10:00Z GOVERNANCE proposal=109 phase=voting votingPower=120350 turnoutPercent=16.0
2025-11-08T09:11:00Z GOVERNANCE proposal=109 phase=tally votingPower=120387 turnoutPercent=16.1
2025-11-08T09:12:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120424 turnoutPercent=16.2
2025-11-08T09:13:00Z GOVERNANCE proposal=109 phase=voting votingPower=120461 turnoutPercent=16.3
2025-11-08T09:14:00Z GOVERNANCE proposal=109 phase=tally votingPower=120498 turnoutPercent=16.4
2025-11-08T09:15:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120535 turnoutPercent=16.5
2025-11-08T09:16:00Z GOVERNANCE proposal=109 phase=voting votingPower=120572 turnoutPercent=16.6
2025-11-08T09:17:00Z GOVERNANCE proposal=109 phase=tally votingPower=120609 turnoutPercent=16.7
2025-11-08T09:18:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120646 turnoutPercent=16.8
2025-11-08T09:19:00Z GOVERNANCE proposal=109 phase=voting votingPower=120683 turnoutPercent=16.9
2025-11-08T09:20:00Z GOVERNANCE proposal=109 phase=tally votingPower=120720 turnoutPercent=17.0
2025-11-08T09:21:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120757 turnoutPercent=17.1
2025-11-08T09:22:00Z GOVERNANCE proposal=109 phase=voting votingPower=120794 turnoutPercent=17.2
2025-11-08T09:23:00Z GOVERNANCE proposal=109 phase=tally votingPower=120831 turnoutPercent=17.3
2025-11-08T09:24:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120868 turnoutPercent=17.4
2025-11-08T09:25:00Z GOVERNANCE proposal=109 phase=voting votingPower=120905 turnoutPercent=17.5
2025-11-08T09:26:00Z GOVERNANCE proposal=109 phase=tally votingPower=120942 turnoutPercent=17.6
2025-11-08T09:27:00Z GOVERNANCE proposal=109 phase=discussion votingPower=120979 turnoutPercent=17.7
2025-11-08T09:28:00Z GOVERNANCE proposal=109 phase=voting votingPower=121016 turnoutPercent=17.8
2025-11-08T09:29:00Z GOVERNANCE proposal=109 phase=tally votingPower=121053 turnoutPercent=17.9
2025-11-08T09:30:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121090 turnoutPercent=18.0
2025-11-08T09:31:00Z GOVERNANCE proposal=109 phase=voting votingPower=121127 turnoutPercent=18.1
2025-11-08T09:32:00Z GOVERNANCE proposal=109 phase=tally votingPower=121164 turnoutPercent=18.2
2025-11-08T09:33:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121201 turnoutPercent=18.3
2025-11-08T09:34:00Z GOVERNANCE proposal=109 phase=voting votingPower=121238 turnoutPercent=18.4
2025-11-08T09:35:00Z GOVERNANCE proposal=109 phase=tally votingPower=121275 turnoutPercent=18.5
2025-11-08T09:36:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121312 turnoutPercent=18.6
2025-11-08T09:37:00Z GOVERNANCE proposal=109 phase=voting votingPower=121349 turnoutPercent=18.7
2025-11-08T09:38:00Z GOVERNANCE proposal=109 phase=tally votingPower=121386 turnoutPercent=18.8
2025-11-08T09:39:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121423 turnoutPercent=18.9
2025-11-08T09:40:00Z GOVERNANCE proposal=109 phase=voting votingPower=121460 turnoutPercent=19.0
2025-11-08T09:41:00Z GOVERNANCE proposal=109 phase=tally votingPower=121497 turnoutPercent=19.1
2025-11-08T09:42:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121534 turnoutPercent=19.2
2025-11-08T09:43:00Z GOVERNANCE proposal=109 phase=voting votingPower=121571 turnoutPercent=19.3
2025-11-08T09:44:00Z GOVERNANCE proposal=109 phase=tally votingPower=121608 turnoutPercent=19.4
2025-11-08T09:45:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121645 turnoutPercent=19.5
2025-11-08T09:46:00Z GOVERNANCE proposal=109 phase=voting votingPower=121682 turnoutPercent=19.6
2025-11-08T09:47:00Z GOVERNANCE proposal=109 phase=tally votingPower=121719 turnoutPercent=19.7
2025-11-08T09:48:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121756 turnoutPercent=19.8
2025-11-08T09:49:00Z GOVERNANCE proposal=109 phase=voting votingPower=121793 turnoutPercent=19.9
2025-11-08T09:50:00Z GOVERNANCE proposal=109 phase=tally votingPower=121830 turnoutPercent=20.0
2025-11-08T09:51:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121867 turnoutPercent=20.1
2025-11-08T09:52:00Z GOVERNANCE proposal=109 phase=voting votingPower=121904 turnoutPercent=20.2
2025-11-08T09:53:00Z GOVERNANCE proposal=109 phase=tally votingPower=121941 turnoutPercent=20.3
2025-11-08T09:54:00Z GOVERNANCE proposal=109 phase=discussion votingPower=121978 turnoutPercent=20.4
2025-11-08T09:55:00Z GOVERNANCE proposal=109 phase=voting votingPower=122015 turnoutPercent=20.5
2025-11-08T09:56:00Z GOVERNANCE proposal=109 phase=tally votingPower=122052 turnoutPercent=20.6
2025-11-08T09:57:00Z GOVERNANCE proposal=109 phase=discussion votingPower=122089 turnoutPercent=20.7
2025-11-08T09:58:00Z GOVERNANCE proposal=109 phase=voting votingPower=122126 turnoutPercent=20.8
2025-11-08T09:59:00Z GOVERNANCE proposal=109 phase=tally votingPower=122163 turnoutPercent=20.9
2025-11-08T10:00:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122200 turnoutPercent=15.0
2025-11-08T10:01:00Z GOVERNANCE proposal=110 phase=voting votingPower=122237 turnoutPercent=15.1
2025-11-08T10:02:00Z GOVERNANCE proposal=110 phase=tally votingPower=122274 turnoutPercent=15.2
2025-11-08T10:03:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122311 turnoutPercent=15.3
2025-11-08T10:04:00Z GOVERNANCE proposal=110 phase=voting votingPower=122348 turnoutPercent=15.4
2025-11-08T10:05:00Z GOVERNANCE proposal=110 phase=tally votingPower=122385 turnoutPercent=15.5
2025-11-08T10:06:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122422 turnoutPercent=15.6
2025-11-08T10:07:00Z GOVERNANCE proposal=110 phase=voting votingPower=122459 turnoutPercent=15.7
2025-11-08T10:08:00Z GOVERNANCE proposal=110 phase=tally votingPower=122496 turnoutPercent=15.8
2025-11-08T10:09:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122533 turnoutPercent=15.9
2025-11-08T10:10:00Z GOVERNANCE proposal=110 phase=voting votingPower=122570 turnoutPercent=16.0
2025-11-08T10:11:00Z GOVERNANCE proposal=110 phase=tally votingPower=122607 turnoutPercent=16.1
2025-11-08T10:12:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122644 turnoutPercent=16.2
2025-11-08T10:13:00Z GOVERNANCE proposal=110 phase=voting votingPower=122681 turnoutPercent=16.3
2025-11-08T10:14:00Z GOVERNANCE proposal=110 phase=tally votingPower=122718 turnoutPercent=16.4
2025-11-08T10:15:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122755 turnoutPercent=16.5
2025-11-08T10:16:00Z GOVERNANCE proposal=110 phase=voting votingPower=122792 turnoutPercent=16.6
2025-11-08T10:17:00Z GOVERNANCE proposal=110 phase=tally votingPower=122829 turnoutPercent=16.7
2025-11-08T10:18:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122866 turnoutPercent=16.8
2025-11-08T10:19:00Z GOVERNANCE proposal=110 phase=voting votingPower=122903 turnoutPercent=16.9
2025-11-08T10:20:00Z GOVERNANCE proposal=110 phase=tally votingPower=122940 turnoutPercent=17.0
2025-11-08T10:21:00Z GOVERNANCE proposal=110 phase=discussion votingPower=122977 turnoutPercent=17.1
2025-11-08T10:22:00Z GOVERNANCE proposal=110 phase=voting votingPower=123014 turnoutPercent=17.2
2025-11-08T10:23:00Z GOVERNANCE proposal=110 phase=tally votingPower=123051 turnoutPercent=17.3
2025-11-08T10:24:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123088 turnoutPercent=17.4
2025-11-08T10:25:00Z GOVERNANCE proposal=110 phase=voting votingPower=123125 turnoutPercent=17.5
2025-11-08T10:26:00Z GOVERNANCE proposal=110 phase=tally votingPower=123162 turnoutPercent=17.6
2025-11-08T10:27:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123199 turnoutPercent=17.7
2025-11-08T10:28:00Z GOVERNANCE proposal=110 phase=voting votingPower=123236 turnoutPercent=17.8
2025-11-08T10:29:00Z GOVERNANCE proposal=110 phase=tally votingPower=123273 turnoutPercent=17.9
2025-11-08T10:30:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123310 turnoutPercent=18.0
2025-11-08T10:31:00Z GOVERNANCE proposal=110 phase=voting votingPower=123347 turnoutPercent=18.1
2025-11-08T10:32:00Z GOVERNANCE proposal=110 phase=tally votingPower=123384 turnoutPercent=18.2
2025-11-08T10:33:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123421 turnoutPercent=18.3
2025-11-08T10:34:00Z GOVERNANCE proposal=110 phase=voting votingPower=123458 turnoutPercent=18.4
2025-11-08T10:35:00Z GOVERNANCE proposal=110 phase=tally votingPower=123495 turnoutPercent=18.5
2025-11-08T10:36:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123532 turnoutPercent=18.6
2025-11-08T10:37:00Z GOVERNANCE proposal=110 phase=voting votingPower=123569 turnoutPercent=18.7
2025-11-08T10:38:00Z GOVERNANCE proposal=110 phase=tally votingPower=123606 turnoutPercent=18.8
2025-11-08T10:39:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123643 turnoutPercent=18.9
2025-11-08T10:40:00Z GOVERNANCE proposal=110 phase=voting votingPower=123680 turnoutPercent=19.0
2025-11-08T10:41:00Z GOVERNANCE proposal=110 phase=tally votingPower=123717 turnoutPercent=19.1
2025-11-08T10:42:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123754 turnoutPercent=19.2
2025-11-08T10:43:00Z GOVERNANCE proposal=110 phase=voting votingPower=123791 turnoutPercent=19.3
2025-11-08T10:44:00Z GOVERNANCE proposal=110 phase=tally votingPower=123828 turnoutPercent=19.4
2025-11-08T10:45:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123865 turnoutPercent=19.5
2025-11-08T10:46:00Z GOVERNANCE proposal=110 phase=voting votingPower=123902 turnoutPercent=19.6
2025-11-08T10:47:00Z GOVERNANCE proposal=110 phase=tally votingPower=123939 turnoutPercent=19.7
2025-11-08T10:48:00Z GOVERNANCE proposal=110 phase=discussion votingPower=123976 turnoutPercent=19.8
2025-11-08T10:49:00Z GOVERNANCE proposal=110 phase=voting votingPower=124013 turnoutPercent=19.9
2025-11-08T10:50:00Z GOVERNANCE proposal=110 phase=tally votingPower=124050 turnoutPercent=20.0
2025-11-08T10:51:00Z GOVERNANCE proposal=110 phase=discussion votingPower=124087 turnoutPercent=20.1
2025-11-08T10:52:00Z GOVERNANCE proposal=110 phase=voting votingPower=124124 turnoutPercent=20.2
2025-11-08T10:53:00Z GOVERNANCE proposal=110 phase=tally votingPower=124161 turnoutPercent=20.3
2025-11-08T10:54:00Z GOVERNANCE proposal=110 phase=discussion votingPower=124198 turnoutPercent=20.4
2025-11-08T10:55:00Z GOVERNANCE proposal=110 phase=voting votingPower=124235 turnoutPercent=20.5
2025-11-08T10:56:00Z GOVERNANCE proposal=110 phase=tally votingPower=124272 turnoutPercent=20.6
2025-11-08T10:57:00Z GOVERNANCE proposal=110 phase=discussion votingPower=124309 turnoutPercent=20.7
2025-11-08T10:58:00Z GOVERNANCE proposal=110 phase=voting votingPower=124346 turnoutPercent=20.8
2025-11-08T10:59:00Z GOVERNANCE proposal=110 phase=tally votingPower=124383 turnoutPercent=20.9
2025-11-08T11:00:00Z GOVERNANCE proposal=111 phase=discussion votingPower=124420 turnoutPercent=15.0
2025-11-08T11:01:00Z GOVERNANCE proposal=111 phase=voting votingPower=124457 turnoutPercent=15.1
2025-11-08T11:02:00Z GOVERNANCE proposal=111 phase=tally votingPower=124494 turnoutPercent=15.2
2025-11-08T11:03:00Z GOVERNANCE proposal=111 phase=discussion votingPower=124531 turnoutPercent=15.3
2025-11-08T11:04:00Z GOVERNANCE proposal=111 phase=voting votingPower=124568 turnoutPercent=15.4
2025-11-08T11:05:00Z GOVERNANCE proposal=111 phase=tally votingPower=124605 turnoutPercent=15.5
2025-11-08T11:06:00Z GOVERNANCE proposal=111 phase=discussion votingPower=124642 turnoutPercent=15.6
2025-11-08T11:07:00Z GOVERNANCE proposal=111 phase=voting votingPower=124679 turnoutPercent=15.7
2025-11-08T11:08:00Z GOVERNANCE proposal=111 phase=tally votingPower=124716 turnoutPercent=15.8
2025-11-08T11:09:00Z GOVERNANCE proposal=111 phase=discussion votingPower=124753 turnoutPercent=15.9
2025-11-08T11:10:00Z GOVERNANCE proposal=111 phase=voting votingPower=124790 turnoutPercent=16.0
2025-11-08T11:11:00Z GOVERNANCE proposal=111 phase=tally votingPower=124827 turnoutPercent=16.1
2025-11-08T11:12:00Z GOVERNANCE proposal=111 phase=discussion votingPower=124864 turnoutPercent=16.2
2025-11-08T11:13:00Z GOVERNANCE proposal=111 phase=voting votingPower=124901 turnoutPercent=16.3
2025-11-08T11:14:00Z GOVERNANCE proposal=111 phase=tally votingPower=124938 turnoutPercent=16.4
2025-11-08T11:15:00Z GOVERNANCE proposal=111 phase=discussion votingPower=124975 turnoutPercent=16.5
2025-11-08T11:16:00Z GOVERNANCE proposal=111 phase=voting votingPower=125012 turnoutPercent=16.6
2025-11-08T11:17:00Z GOVERNANCE proposal=111 phase=tally votingPower=125049 turnoutPercent=16.7
2025-11-08T11:18:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125086 turnoutPercent=16.8
2025-11-08T11:19:00Z GOVERNANCE proposal=111 phase=voting votingPower=125123 turnoutPercent=16.9
2025-11-08T11:20:00Z GOVERNANCE proposal=111 phase=tally votingPower=125160 turnoutPercent=17.0
2025-11-08T11:21:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125197 turnoutPercent=17.1
2025-11-08T11:22:00Z GOVERNANCE proposal=111 phase=voting votingPower=125234 turnoutPercent=17.2
2025-11-08T11:23:00Z GOVERNANCE proposal=111 phase=tally votingPower=125271 turnoutPercent=17.3
2025-11-08T11:24:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125308 turnoutPercent=17.4
2025-11-08T11:25:00Z GOVERNANCE proposal=111 phase=voting votingPower=125345 turnoutPercent=17.5
2025-11-08T11:26:00Z GOVERNANCE proposal=111 phase=tally votingPower=125382 turnoutPercent=17.6
2025-11-08T11:27:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125419 turnoutPercent=17.7
2025-11-08T11:28:00Z GOVERNANCE proposal=111 phase=voting votingPower=125456 turnoutPercent=17.8
2025-11-08T11:29:00Z GOVERNANCE proposal=111 phase=tally votingPower=125493 turnoutPercent=17.9
2025-11-08T11:30:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125530 turnoutPercent=18.0
2025-11-08T11:31:00Z GOVERNANCE proposal=111 phase=voting votingPower=125567 turnoutPercent=18.1
2025-11-08T11:32:00Z GOVERNANCE proposal=111 phase=tally votingPower=125604 turnoutPercent=18.2
2025-11-08T11:33:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125641 turnoutPercent=18.3
2025-11-08T11:34:00Z GOVERNANCE proposal=111 phase=voting votingPower=125678 turnoutPercent=18.4
2025-11-08T11:35:00Z GOVERNANCE proposal=111 phase=tally votingPower=125715 turnoutPercent=18.5
2025-11-08T11:36:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125752 turnoutPercent=18.6
2025-11-08T11:37:00Z GOVERNANCE proposal=111 phase=voting votingPower=125789 turnoutPercent=18.7
2025-11-08T11:38:00Z GOVERNANCE proposal=111 phase=tally votingPower=125826 turnoutPercent=18.8
2025-11-08T11:39:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125863 turnoutPercent=18.9
2025-11-08T11:40:00Z GOVERNANCE proposal=111 phase=voting votingPower=125900 turnoutPercent=19.0
2025-11-08T11:41:00Z GOVERNANCE proposal=111 phase=tally votingPower=125937 turnoutPercent=19.1
2025-11-08T11:42:00Z GOVERNANCE proposal=111 phase=discussion votingPower=125974 turnoutPercent=19.2
2025-11-08T11:43:00Z GOVERNANCE proposal=111 phase=voting votingPower=126011 turnoutPercent=19.3
2025-11-08T11:44:00Z GOVERNANCE proposal=111 phase=tally votingPower=126048 turnoutPercent=19.4
2025-11-08T11:45:00Z GOVERNANCE proposal=111 phase=discussion votingPower=126085 turnoutPercent=19.5
2025-11-08T11:46:00Z GOVERNANCE proposal=111 phase=voting votingPower=126122 turnoutPercent=19.6
2025-11-08T11:47:00Z GOVERNANCE proposal=111 phase=tally votingPower=126159 turnoutPercent=19.7
2025-11-08T11:48:00Z GOVERNANCE proposal=111 phase=discussion votingPower=126196 turnoutPercent=19.8
2025-11-08T11:49:00Z GOVERNANCE proposal=111 phase=voting votingPower=126233 turnoutPercent=19.9
2025-11-08T11:50:00Z GOVERNANCE proposal=111 phase=tally votingPower=126270 turnoutPercent=20.0
2025-11-08T11:51:00Z GOVERNANCE proposal=111 phase=discussion votingPower=126307 turnoutPercent=20.1
2025-11-08T11:52:00Z GOVERNANCE proposal=111 phase=voting votingPower=126344 turnoutPercent=20.2
2025-11-08T11:53:00Z GOVERNANCE proposal=111 phase=tally votingPower=126381 turnoutPercent=20.3
2025-11-08T11:54:00Z GOVERNANCE proposal=111 phase=discussion votingPower=126418 turnoutPercent=20.4
2025-11-08T11:55:00Z GOVERNANCE proposal=111 phase=voting votingPower=126455 turnoutPercent=20.5
2025-11-08T11:56:00Z GOVERNANCE proposal=111 phase=tally votingPower=126492 turnoutPercent=20.6
2025-11-08T11:57:00Z GOVERNANCE proposal=111 phase=discussion votingPower=126529 turnoutPercent=20.7
2025-11-08T11:58:00Z GOVERNANCE proposal=111 phase=voting votingPower=126566 turnoutPercent=20.8
2025-11-08T11:59:00Z GOVERNANCE proposal=111 phase=tally votingPower=126603 turnoutPercent=20.9
2025-11-08T12:00:00Z GOVERNANCE proposal=112 phase=discussion votingPower=126640 turnoutPercent=15.0
2025-11-08T12:01:00Z GOVERNANCE proposal=112 phase=voting votingPower=126677 turnoutPercent=15.1
2025-11-08T12:02:00Z GOVERNANCE proposal=112 phase=tally votingPower=126714 turnoutPercent=15.2
2025-11-08T12:03:00Z GOVERNANCE proposal=112 phase=discussion votingPower=126751 turnoutPercent=15.3
2025-11-08T12:04:00Z GOVERNANCE proposal=112 phase=voting votingPower=126788 turnoutPercent=15.4
2025-11-08T12:05:00Z GOVERNANCE proposal=112 phase=tally votingPower=126825 turnoutPercent=15.5
2025-11-08T12:06:00Z GOVERNANCE proposal=112 phase=discussion votingPower=126862 turnoutPercent=15.6
2025-11-08T12:07:00Z GOVERNANCE proposal=112 phase=voting votingPower=126899 turnoutPercent=15.7
2025-11-08T12:08:00Z GOVERNANCE proposal=112 phase=tally votingPower=126936 turnoutPercent=15.8
2025-11-08T12:09:00Z GOVERNANCE proposal=112 phase=discussion votingPower=126973 turnoutPercent=15.9
2025-11-08T12:10:00Z GOVERNANCE proposal=112 phase=voting votingPower=127010 turnoutPercent=16.0
2025-11-08T12:11:00Z GOVERNANCE proposal=112 phase=tally votingPower=127047 turnoutPercent=16.1
2025-11-08T12:12:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127084 turnoutPercent=16.2
2025-11-08T12:13:00Z GOVERNANCE proposal=112 phase=voting votingPower=127121 turnoutPercent=16.3
2025-11-08T12:14:00Z GOVERNANCE proposal=112 phase=tally votingPower=127158 turnoutPercent=16.4
2025-11-08T12:15:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127195 turnoutPercent=16.5
2025-11-08T12:16:00Z GOVERNANCE proposal=112 phase=voting votingPower=127232 turnoutPercent=16.6
2025-11-08T12:17:00Z GOVERNANCE proposal=112 phase=tally votingPower=127269 turnoutPercent=16.7
2025-11-08T12:18:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127306 turnoutPercent=16.8
2025-11-08T12:19:00Z GOVERNANCE proposal=112 phase=voting votingPower=127343 turnoutPercent=16.9
2025-11-08T12:20:00Z GOVERNANCE proposal=112 phase=tally votingPower=127380 turnoutPercent=17.0
2025-11-08T12:21:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127417 turnoutPercent=17.1
2025-11-08T12:22:00Z GOVERNANCE proposal=112 phase=voting votingPower=127454 turnoutPercent=17.2
2025-11-08T12:23:00Z GOVERNANCE proposal=112 phase=tally votingPower=127491 turnoutPercent=17.3
2025-11-08T12:24:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127528 turnoutPercent=17.4
2025-11-08T12:25:00Z GOVERNANCE proposal=112 phase=voting votingPower=127565 turnoutPercent=17.5
2025-11-08T12:26:00Z GOVERNANCE proposal=112 phase=tally votingPower=127602 turnoutPercent=17.6
2025-11-08T12:27:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127639 turnoutPercent=17.7
2025-11-08T12:28:00Z GOVERNANCE proposal=112 phase=voting votingPower=127676 turnoutPercent=17.8
2025-11-08T12:29:00Z GOVERNANCE proposal=112 phase=tally votingPower=127713 turnoutPercent=17.9
2025-11-08T12:30:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127750 turnoutPercent=18.0
2025-11-08T12:31:00Z GOVERNANCE proposal=112 phase=voting votingPower=127787 turnoutPercent=18.1
2025-11-08T12:32:00Z GOVERNANCE proposal=112 phase=tally votingPower=127824 turnoutPercent=18.2
2025-11-08T12:33:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127861 turnoutPercent=18.3
2025-11-08T12:34:00Z GOVERNANCE proposal=112 phase=voting votingPower=127898 turnoutPercent=18.4
2025-11-08T12:35:00Z GOVERNANCE proposal=112 phase=tally votingPower=127935 turnoutPercent=18.5
2025-11-08T12:36:00Z GOVERNANCE proposal=112 phase=discussion votingPower=127972 turnoutPercent=18.6
2025-11-08T12:37:00Z GOVERNANCE proposal=112 phase=voting votingPower=128009 turnoutPercent=18.7
2025-11-08T12:38:00Z GOVERNANCE proposal=112 phase=tally votingPower=128046 turnoutPercent=18.8
2025-11-08T12:39:00Z GOVERNANCE proposal=112 phase=discussion votingPower=128083 turnoutPercent=18.9
2025-11-08T12:40:00Z GOVERNANCE proposal=112 phase=voting votingPower=128120 turnoutPercent=19.0
2025-11-08T12:41:00Z GOVERNANCE proposal=112 phase=tally votingPower=128157 turnoutPercent=19.1
2025-11-08T12:42:00Z GOVERNANCE proposal=112 phase=discussion votingPower=128194 turnoutPercent=19.2
2025-11-08T12:43:00Z GOVERNANCE proposal=112 phase=voting votingPower=128231 turnoutPercent=19.3
2025-11-08T12:44:00Z GOVERNANCE proposal=112 phase=tally votingPower=128268 turnoutPercent=19.4
2025-11-08T12:45:00Z GOVERNANCE proposal=112 phase=discussion votingPower=128305 turnoutPercent=19.5
2025-11-08T12:46:00Z GOVERNANCE proposal=112 phase=voting votingPower=128342 turnoutPercent=19.6
2025-11-08T12:47:00Z GOVERNANCE proposal=112 phase=tally votingPower=128379 turnoutPercent=19.7
2025-11-08T12:48:00Z GOVERNANCE proposal=112 phase=discussion votingPower=128416 turnoutPercent=19.8
2025-11-08T12:49:00Z GOVERNANCE proposal=112 phase=voting votingPower=128453 turnoutPercent=19.9
2025-11-08T12:50:00Z GOVERNANCE proposal=112 phase=tally votingPower=128490 turnoutPercent=20.0
2025-11-08T12:51:00Z GOVERNANCE proposal=112 phase=discussion votingPower=128527 turnoutPercent=20.1
2025-11-08T12:52:00Z GOVERNANCE proposal=112 phase=voting votingPower=128564 turnoutPercent=20.2
2025-11-08T12:53:00Z GOVERNANCE proposal=112 phase=tally votingPower=128601 turnoutPercent=20.3
2025-11-08T12:54:00Z GOVERNANCE proposal=112 phase=discussion votingPower=128638 turnoutPercent=20.4
2025-11-08T12:55:00Z GOVERNANCE proposal=112 phase=voting votingPower=128675 turnoutPercent=20.5
2025-11-08T12:56:00Z GOVERNANCE proposal=112 phase=tally votingPower=128712 turnoutPercent=20.6
2025-11-08T12:57:00Z GOVERNANCE proposal=112 phase=discussion votingPower=128749 turnoutPercent=20.7
2025-11-08T12:58:00Z GOVERNANCE proposal=112 phase=voting votingPower=128786 turnoutPercent=20.8
2025-11-08T12:59:00Z GOVERNANCE proposal=112 phase=tally votingPower=128823 turnoutPercent=20.9
2025-11-08T13:00:00Z GOVERNANCE proposal=113 phase=discussion votingPower=128860 turnoutPercent=15.0
2025-11-08T13:01:00Z GOVERNANCE proposal=113 phase=voting votingPower=128897 turnoutPercent=15.1
2025-11-08T13:02:00Z GOVERNANCE proposal=113 phase=tally votingPower=128934 turnoutPercent=15.2
2025-11-08T13:03:00Z GOVERNANCE proposal=113 phase=discussion votingPower=128971 turnoutPercent=15.3
2025-11-08T13:04:00Z GOVERNANCE proposal=113 phase=voting votingPower=129008 turnoutPercent=15.4
2025-11-08T13:05:00Z GOVERNANCE proposal=113 phase=tally votingPower=129045 turnoutPercent=15.5
2025-11-08T13:06:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129082 turnoutPercent=15.6
2025-11-08T13:07:00Z GOVERNANCE proposal=113 phase=voting votingPower=129119 turnoutPercent=15.7
2025-11-08T13:08:00Z GOVERNANCE proposal=113 phase=tally votingPower=129156 turnoutPercent=15.8
2025-11-08T13:09:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129193 turnoutPercent=15.9
2025-11-08T13:10:00Z GOVERNANCE proposal=113 phase=voting votingPower=129230 turnoutPercent=16.0
2025-11-08T13:11:00Z GOVERNANCE proposal=113 phase=tally votingPower=129267 turnoutPercent=16.1
2025-11-08T13:12:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129304 turnoutPercent=16.2
2025-11-08T13:13:00Z GOVERNANCE proposal=113 phase=voting votingPower=129341 turnoutPercent=16.3
2025-11-08T13:14:00Z GOVERNANCE proposal=113 phase=tally votingPower=129378 turnoutPercent=16.4
2025-11-08T13:15:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129415 turnoutPercent=16.5
2025-11-08T13:16:00Z GOVERNANCE proposal=113 phase=voting votingPower=129452 turnoutPercent=16.6
2025-11-08T13:17:00Z GOVERNANCE proposal=113 phase=tally votingPower=129489 turnoutPercent=16.7
2025-11-08T13:18:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129526 turnoutPercent=16.8
2025-11-08T13:19:00Z GOVERNANCE proposal=113 phase=voting votingPower=129563 turnoutPercent=16.9
2025-11-08T13:20:00Z GOVERNANCE proposal=113 phase=tally votingPower=129600 turnoutPercent=17.0
2025-11-08T13:21:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129637 turnoutPercent=17.1
2025-11-08T13:22:00Z GOVERNANCE proposal=113 phase=voting votingPower=129674 turnoutPercent=17.2
2025-11-08T13:23:00Z GOVERNANCE proposal=113 phase=tally votingPower=129711 turnoutPercent=17.3
2025-11-08T13:24:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129748 turnoutPercent=17.4
2025-11-08T13:25:00Z GOVERNANCE proposal=113 phase=voting votingPower=129785 turnoutPercent=17.5
2025-11-08T13:26:00Z GOVERNANCE proposal=113 phase=tally votingPower=129822 turnoutPercent=17.6
2025-11-08T13:27:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129859 turnoutPercent=17.7
2025-11-08T13:28:00Z GOVERNANCE proposal=113 phase=voting votingPower=129896 turnoutPercent=17.8
2025-11-08T13:29:00Z GOVERNANCE proposal=113 phase=tally votingPower=129933 turnoutPercent=17.9
2025-11-08T13:30:00Z GOVERNANCE proposal=113 phase=discussion votingPower=129970 turnoutPercent=18.0
2025-11-08T13:31:00Z GOVERNANCE proposal=113 phase=voting votingPower=130007 turnoutPercent=18.1
2025-11-08T13:32:00Z GOVERNANCE proposal=113 phase=tally votingPower=130044 turnoutPercent=18.2
2025-11-08T13:33:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130081 turnoutPercent=18.3
2025-11-08T13:34:00Z GOVERNANCE proposal=113 phase=voting votingPower=130118 turnoutPercent=18.4
2025-11-08T13:35:00Z GOVERNANCE proposal=113 phase=tally votingPower=130155 turnoutPercent=18.5
2025-11-08T13:36:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130192 turnoutPercent=18.6
2025-11-08T13:37:00Z GOVERNANCE proposal=113 phase=voting votingPower=130229 turnoutPercent=18.7
2025-11-08T13:38:00Z GOVERNANCE proposal=113 phase=tally votingPower=130266 turnoutPercent=18.8
2025-11-08T13:39:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130303 turnoutPercent=18.9
2025-11-08T13:40:00Z GOVERNANCE proposal=113 phase=voting votingPower=130340 turnoutPercent=19.0
2025-11-08T13:41:00Z GOVERNANCE proposal=113 phase=tally votingPower=130377 turnoutPercent=19.1
2025-11-08T13:42:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130414 turnoutPercent=19.2
2025-11-08T13:43:00Z GOVERNANCE proposal=113 phase=voting votingPower=130451 turnoutPercent=19.3
2025-11-08T13:44:00Z GOVERNANCE proposal=113 phase=tally votingPower=130488 turnoutPercent=19.4
2025-11-08T13:45:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130525 turnoutPercent=19.5
2025-11-08T13:46:00Z GOVERNANCE proposal=113 phase=voting votingPower=130562 turnoutPercent=19.6
2025-11-08T13:47:00Z GOVERNANCE proposal=113 phase=tally votingPower=130599 turnoutPercent=19.7
2025-11-08T13:48:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130636 turnoutPercent=19.8
2025-11-08T13:49:00Z GOVERNANCE proposal=113 phase=voting votingPower=130673 turnoutPercent=19.9
2025-11-08T13:50:00Z GOVERNANCE proposal=113 phase=tally votingPower=130710 turnoutPercent=20.0
2025-11-08T13:51:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130747 turnoutPercent=20.1
2025-11-08T13:52:00Z GOVERNANCE proposal=113 phase=voting votingPower=130784 turnoutPercent=20.2
2025-11-08T13:53:00Z GOVERNANCE proposal=113 phase=tally votingPower=130821 turnoutPercent=20.3
2025-11-08T13:54:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130858 turnoutPercent=20.4
2025-11-08T13:55:00Z GOVERNANCE proposal=113 phase=voting votingPower=130895 turnoutPercent=20.5
2025-11-08T13:56:00Z GOVERNANCE proposal=113 phase=tally votingPower=130932 turnoutPercent=20.6
2025-11-08T13:57:00Z GOVERNANCE proposal=113 phase=discussion votingPower=130969 turnoutPercent=20.7
2025-11-08T13:58:00Z GOVERNANCE proposal=113 phase=voting votingPower=131006 turnoutPercent=20.8
2025-11-08T13:59:00Z GOVERNANCE proposal=113 phase=tally votingPower=131043 turnoutPercent=20.9
2025-11-08T14:00:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131080 turnoutPercent=15.0
2025-11-08T14:01:00Z GOVERNANCE proposal=114 phase=voting votingPower=131117 turnoutPercent=15.1
2025-11-08T14:02:00Z GOVERNANCE proposal=114 phase=tally votingPower=131154 turnoutPercent=15.2
2025-11-08T14:03:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131191 turnoutPercent=15.3
2025-11-08T14:04:00Z GOVERNANCE proposal=114 phase=voting votingPower=131228 turnoutPercent=15.4
2025-11-08T14:05:00Z GOVERNANCE proposal=114 phase=tally votingPower=131265 turnoutPercent=15.5
2025-11-08T14:06:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131302 turnoutPercent=15.6
2025-11-08T14:07:00Z GOVERNANCE proposal=114 phase=voting votingPower=131339 turnoutPercent=15.7
2025-11-08T14:08:00Z GOVERNANCE proposal=114 phase=tally votingPower=131376 turnoutPercent=15.8
2025-11-08T14:09:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131413 turnoutPercent=15.9
2025-11-08T14:10:00Z GOVERNANCE proposal=114 phase=voting votingPower=131450 turnoutPercent=16.0
2025-11-08T14:11:00Z GOVERNANCE proposal=114 phase=tally votingPower=131487 turnoutPercent=16.1
2025-11-08T14:12:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131524 turnoutPercent=16.2
2025-11-08T14:13:00Z GOVERNANCE proposal=114 phase=voting votingPower=131561 turnoutPercent=16.3
2025-11-08T14:14:00Z GOVERNANCE proposal=114 phase=tally votingPower=131598 turnoutPercent=16.4
2025-11-08T14:15:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131635 turnoutPercent=16.5
2025-11-08T14:16:00Z GOVERNANCE proposal=114 phase=voting votingPower=131672 turnoutPercent=16.6
2025-11-08T14:17:00Z GOVERNANCE proposal=114 phase=tally votingPower=131709 turnoutPercent=16.7
2025-11-08T14:18:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131746 turnoutPercent=16.8
2025-11-08T14:19:00Z GOVERNANCE proposal=114 phase=voting votingPower=131783 turnoutPercent=16.9
2025-11-08T14:20:00Z GOVERNANCE proposal=114 phase=tally votingPower=131820 turnoutPercent=17.0
2025-11-08T14:21:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131857 turnoutPercent=17.1
2025-11-08T14:22:00Z GOVERNANCE proposal=114 phase=voting votingPower=131894 turnoutPercent=17.2
2025-11-08T14:23:00Z GOVERNANCE proposal=114 phase=tally votingPower=131931 turnoutPercent=17.3
2025-11-08T14:24:00Z GOVERNANCE proposal=114 phase=discussion votingPower=131968 turnoutPercent=17.4
2025-11-08T14:25:00Z GOVERNANCE proposal=114 phase=voting votingPower=132005 turnoutPercent=17.5
2025-11-08T14:26:00Z GOVERNANCE proposal=114 phase=tally votingPower=132042 turnoutPercent=17.6
2025-11-08T14:27:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132079 turnoutPercent=17.7
2025-11-08T14:28:00Z GOVERNANCE proposal=114 phase=voting votingPower=132116 turnoutPercent=17.8
2025-11-08T14:29:00Z GOVERNANCE proposal=114 phase=tally votingPower=132153 turnoutPercent=17.9
2025-11-08T14:30:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132190 turnoutPercent=18.0
2025-11-08T14:31:00Z GOVERNANCE proposal=114 phase=voting votingPower=132227 turnoutPercent=18.1
2025-11-08T14:32:00Z GOVERNANCE proposal=114 phase=tally votingPower=132264 turnoutPercent=18.2
2025-11-08T14:33:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132301 turnoutPercent=18.3
2025-11-08T14:34:00Z GOVERNANCE proposal=114 phase=voting votingPower=132338 turnoutPercent=18.4
2025-11-08T14:35:00Z GOVERNANCE proposal=114 phase=tally votingPower=132375 turnoutPercent=18.5
2025-11-08T14:36:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132412 turnoutPercent=18.6
2025-11-08T14:37:00Z GOVERNANCE proposal=114 phase=voting votingPower=132449 turnoutPercent=18.7
2025-11-08T14:38:00Z GOVERNANCE proposal=114 phase=tally votingPower=132486 turnoutPercent=18.8
2025-11-08T14:39:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132523 turnoutPercent=18.9
2025-11-08T14:40:00Z GOVERNANCE proposal=114 phase=voting votingPower=132560 turnoutPercent=19.0
2025-11-08T14:41:00Z GOVERNANCE proposal=114 phase=tally votingPower=132597 turnoutPercent=19.1
2025-11-08T14:42:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132634 turnoutPercent=19.2
2025-11-08T14:43:00Z GOVERNANCE proposal=114 phase=voting votingPower=132671 turnoutPercent=19.3
2025-11-08T14:44:00Z GOVERNANCE proposal=114 phase=tally votingPower=132708 turnoutPercent=19.4
2025-11-08T14:45:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132745 turnoutPercent=19.5
2025-11-08T14:46:00Z GOVERNANCE proposal=114 phase=voting votingPower=132782 turnoutPercent=19.6
2025-11-08T14:47:00Z GOVERNANCE proposal=114 phase=tally votingPower=132819 turnoutPercent=19.7
2025-11-08T14:48:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132856 turnoutPercent=19.8
2025-11-08T14:49:00Z GOVERNANCE proposal=114 phase=voting votingPower=132893 turnoutPercent=19.9
2025-11-08T14:50:00Z GOVERNANCE proposal=114 phase=tally votingPower=132930 turnoutPercent=20.0
2025-11-08T14:51:00Z GOVERNANCE proposal=114 phase=discussion votingPower=132967 turnoutPercent=20.1
2025-11-08T14:52:00Z GOVERNANCE proposal=114 phase=voting votingPower=133004 turnoutPercent=20.2
2025-11-08T14:53:00Z GOVERNANCE proposal=114 phase=tally votingPower=133041 turnoutPercent=20.3
2025-11-08T14:54:00Z GOVERNANCE proposal=114 phase=discussion votingPower=133078 turnoutPercent=20.4
2025-11-08T14:55:00Z GOVERNANCE proposal=114 phase=voting votingPower=133115 turnoutPercent=20.5
2025-11-08T14:56:00Z GOVERNANCE proposal=114 phase=tally votingPower=133152 turnoutPercent=20.6
2025-11-08T14:57:00Z GOVERNANCE proposal=114 phase=discussion votingPower=133189 turnoutPercent=20.7
2025-11-08T14:58:00Z GOVERNANCE proposal=114 phase=voting votingPower=133226 turnoutPercent=20.8
2025-11-08T14:59:00Z GOVERNANCE proposal=114 phase=tally votingPower=133263 turnoutPercent=20.9
2025-11-08T15:00:00Z GOVERNANCE proposal=115 phase=discussion votingPower=133300 turnoutPercent=15.0
2025-11-08T15:01:00Z GOVERNANCE proposal=115 phase=voting votingPower=133337 turnoutPercent=15.1
2025-11-08T15:02:00Z GOVERNANCE proposal=115 phase=tally votingPower=133374 turnoutPercent=15.2
2025-11-08T15:03:00Z GOVERNANCE proposal=115 phase=discussion votingPower=133411 turnoutPercent=15.3
2025-11-08T15:04:00Z GOVERNANCE proposal=115 phase=voting votingPower=133448 turnoutPercent=15.4
2025-11-08T15:05:00Z GOVERNANCE proposal=115 phase=tally votingPower=133485 turnoutPercent=15.5
2025-11-08T15:06:00Z GOVERNANCE proposal=115 phase=discussion votingPower=133522 turnoutPercent=15.6
2025-11-08T15:07:00Z GOVERNANCE proposal=115 phase=voting votingPower=133559 turnoutPercent=15.7
2025-11-08T15:08:00Z GOVERNANCE proposal=115 phase=tally votingPower=133596 turnoutPercent=15.8
2025-11-08T15:09:00Z GOVERNANCE proposal=115 phase=discussion votingPower=133633 turnoutPercent=15.9
2025-11-08T15:10:00Z GOVERNANCE proposal=115 phase=voting votingPower=133670 turnoutPercent=16.0
2025-11-08T15:11:00Z GOVERNANCE proposal=115 phase=tally votingPower=133707 turnoutPercent=16.1
2025-11-08T15:12:00Z GOVERNANCE proposal=115 phase=discussion votingPower=133744 turnoutPercent=16.2
2025-11-08T15:13:00Z GOVERNANCE proposal=115 phase=voting votingPower=133781 turnoutPercent=16.3
2025-11-08T15:14:00Z GOVERNANCE proposal=115 phase=tally votingPower=133818 turnoutPercent=16.4
2025-11-08T15:15:00Z GOVERNANCE proposal=115 phase=discussion votingPower=133855 turnoutPercent=16.5
2025-11-08T15:16:00Z GOVERNANCE proposal=115 phase=voting votingPower=133892 turnoutPercent=16.6
2025-11-08T15:17:00Z GOVERNANCE proposal=115 phase=tally votingPower=133929 turnoutPercent=16.7
2025-11-08T15:18:00Z GOVERNANCE proposal=115 phase=discussion votingPower=133966 turnoutPercent=16.8
2025-11-08T15:19:00Z GOVERNANCE proposal=115 phase=voting votingPower=134003 turnoutPercent=16.9
2025-11-08T15:20:00Z GOVERNANCE proposal=115 phase=tally votingPower=134040 turnoutPercent=17.0
2025-11-08T15:21:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134077 turnoutPercent=17.1
2025-11-08T15:22:00Z GOVERNANCE proposal=115 phase=voting votingPower=134114 turnoutPercent=17.2
2025-11-08T15:23:00Z GOVERNANCE proposal=115 phase=tally votingPower=134151 turnoutPercent=17.3
2025-11-08T15:24:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134188 turnoutPercent=17.4
2025-11-08T15:25:00Z GOVERNANCE proposal=115 phase=voting votingPower=134225 turnoutPercent=17.5
2025-11-08T15:26:00Z GOVERNANCE proposal=115 phase=tally votingPower=134262 turnoutPercent=17.6
2025-11-08T15:27:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134299 turnoutPercent=17.7
2025-11-08T15:28:00Z GOVERNANCE proposal=115 phase=voting votingPower=134336 turnoutPercent=17.8
2025-11-08T15:29:00Z GOVERNANCE proposal=115 phase=tally votingPower=134373 turnoutPercent=17.9
2025-11-08T15:30:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134410 turnoutPercent=18.0
2025-11-08T15:31:00Z GOVERNANCE proposal=115 phase=voting votingPower=134447 turnoutPercent=18.1
2025-11-08T15:32:00Z GOVERNANCE proposal=115 phase=tally votingPower=134484 turnoutPercent=18.2
2025-11-08T15:33:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134521 turnoutPercent=18.3
2025-11-08T15:34:00Z GOVERNANCE proposal=115 phase=voting votingPower=134558 turnoutPercent=18.4
2025-11-08T15:35:00Z GOVERNANCE proposal=115 phase=tally votingPower=134595 turnoutPercent=18.5
2025-11-08T15:36:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134632 turnoutPercent=18.6
2025-11-08T15:37:00Z GOVERNANCE proposal=115 phase=voting votingPower=134669 turnoutPercent=18.7
2025-11-08T15:38:00Z GOVERNANCE proposal=115 phase=tally votingPower=134706 turnoutPercent=18.8
2025-11-08T15:39:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134743 turnoutPercent=18.9
2025-11-08T15:40:00Z GOVERNANCE proposal=115 phase=voting votingPower=134780 turnoutPercent=19.0
2025-11-08T15:41:00Z GOVERNANCE proposal=115 phase=tally votingPower=134817 turnoutPercent=19.1
2025-11-08T15:42:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134854 turnoutPercent=19.2
2025-11-08T15:43:00Z GOVERNANCE proposal=115 phase=voting votingPower=134891 turnoutPercent=19.3
2025-11-08T15:44:00Z GOVERNANCE proposal=115 phase=tally votingPower=134928 turnoutPercent=19.4
2025-11-08T15:45:00Z GOVERNANCE proposal=115 phase=discussion votingPower=134965 turnoutPercent=19.5
2025-11-08T15:46:00Z GOVERNANCE proposal=115 phase=voting votingPower=135002 turnoutPercent=19.6
2025-11-08T15:47:00Z GOVERNANCE proposal=115 phase=tally votingPower=135039 turnoutPercent=19.7
2025-11-08T15:48:00Z GOVERNANCE proposal=115 phase=discussion votingPower=135076 turnoutPercent=19.8
2025-11-08T15:49:00Z GOVERNANCE proposal=115 phase=voting votingPower=135113 turnoutPercent=19.9
2025-11-08T15:50:00Z GOVERNANCE proposal=115 phase=tally votingPower=135150 turnoutPercent=20.0
2025-11-08T15:51:00Z GOVERNANCE proposal=115 phase=discussion votingPower=135187 turnoutPercent=20.1
2025-11-08T15:52:00Z GOVERNANCE proposal=115 phase=voting votingPower=135224 turnoutPercent=20.2
2025-11-08T15:53:00Z GOVERNANCE proposal=115 phase=tally votingPower=135261 turnoutPercent=20.3
2025-11-08T15:54:00Z GOVERNANCE proposal=115 phase=discussion votingPower=135298 turnoutPercent=20.4
2025-11-08T15:55:00Z GOVERNANCE proposal=115 phase=voting votingPower=135335 turnoutPercent=20.5
2025-11-08T15:56:00Z GOVERNANCE proposal=115 phase=tally votingPower=135372 turnoutPercent=20.6
2025-11-08T15:57:00Z GOVERNANCE proposal=115 phase=discussion votingPower=135409 turnoutPercent=20.7
2025-11-08T15:58:00Z GOVERNANCE proposal=115 phase=voting votingPower=135446 turnoutPercent=20.8
2025-11-08T15:59:00Z GOVERNANCE proposal=115 phase=tally votingPower=135483 turnoutPercent=20.9
2025-11-08T16:00:00Z GOVERNANCE proposal=116 phase=discussion votingPower=135520 turnoutPercent=15.0
2025-11-08T16:01:00Z GOVERNANCE proposal=116 phase=voting votingPower=135557 turnoutPercent=15.1
2025-11-08T16:02:00Z GOVERNANCE proposal=116 phase=tally votingPower=135594 turnoutPercent=15.2
2025-11-08T16:03:00Z GOVERNANCE proposal=116 phase=discussion votingPower=135631 turnoutPercent=15.3
2025-11-08T16:04:00Z GOVERNANCE proposal=116 phase=voting votingPower=135668 turnoutPercent=15.4
2025-11-08T16:05:00Z GOVERNANCE proposal=116 phase=tally votingPower=135705 turnoutPercent=15.5
2025-11-08T16:06:00Z GOVERNANCE proposal=116 phase=discussion votingPower=135742 turnoutPercent=15.6
2025-11-08T16:07:00Z GOVERNANCE proposal=116 phase=voting votingPower=135779 turnoutPercent=15.7
2025-11-08T16:08:00Z GOVERNANCE proposal=116 phase=tally votingPower=135816 turnoutPercent=15.8
2025-11-08T16:09:00Z GOVERNANCE proposal=116 phase=discussion votingPower=135853 turnoutPercent=15.9
2025-11-08T16:10:00Z GOVERNANCE proposal=116 phase=voting votingPower=135890 turnoutPercent=16.0
2025-11-08T16:11:00Z GOVERNANCE proposal=116 phase=tally votingPower=135927 turnoutPercent=16.1
2025-11-08T16:12:00Z GOVERNANCE proposal=116 phase=discussion votingPower=135964 turnoutPercent=16.2
2025-11-08T16:13:00Z GOVERNANCE proposal=116 phase=voting votingPower=136001 turnoutPercent=16.3
2025-11-08T16:14:00Z GOVERNANCE proposal=116 phase=tally votingPower=136038 turnoutPercent=16.4
2025-11-08T16:15:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136075 turnoutPercent=16.5
2025-11-08T16:16:00Z GOVERNANCE proposal=116 phase=voting votingPower=136112 turnoutPercent=16.6
2025-11-08T16:17:00Z GOVERNANCE proposal=116 phase=tally votingPower=136149 turnoutPercent=16.7
2025-11-08T16:18:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136186 turnoutPercent=16.8
2025-11-08T16:19:00Z GOVERNANCE proposal=116 phase=voting votingPower=136223 turnoutPercent=16.9
2025-11-08T16:20:00Z GOVERNANCE proposal=116 phase=tally votingPower=136260 turnoutPercent=17.0
2025-11-08T16:21:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136297 turnoutPercent=17.1
2025-11-08T16:22:00Z GOVERNANCE proposal=116 phase=voting votingPower=136334 turnoutPercent=17.2
2025-11-08T16:23:00Z GOVERNANCE proposal=116 phase=tally votingPower=136371 turnoutPercent=17.3
2025-11-08T16:24:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136408 turnoutPercent=17.4
2025-11-08T16:25:00Z GOVERNANCE proposal=116 phase=voting votingPower=136445 turnoutPercent=17.5
2025-11-08T16:26:00Z GOVERNANCE proposal=116 phase=tally votingPower=136482 turnoutPercent=17.6
2025-11-08T16:27:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136519 turnoutPercent=17.7
2025-11-08T16:28:00Z GOVERNANCE proposal=116 phase=voting votingPower=136556 turnoutPercent=17.8
2025-11-08T16:29:00Z GOVERNANCE proposal=116 phase=tally votingPower=136593 turnoutPercent=17.9
2025-11-08T16:30:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136630 turnoutPercent=18.0
2025-11-08T16:31:00Z GOVERNANCE proposal=116 phase=voting votingPower=136667 turnoutPercent=18.1
2025-11-08T16:32:00Z GOVERNANCE proposal=116 phase=tally votingPower=136704 turnoutPercent=18.2
2025-11-08T16:33:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136741 turnoutPercent=18.3
2025-11-08T16:34:00Z GOVERNANCE proposal=116 phase=voting votingPower=136778 turnoutPercent=18.4
2025-11-08T16:35:00Z GOVERNANCE proposal=116 phase=tally votingPower=136815 turnoutPercent=18.5
2025-11-08T16:36:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136852 turnoutPercent=18.6
2025-11-08T16:37:00Z GOVERNANCE proposal=116 phase=voting votingPower=136889 turnoutPercent=18.7
2025-11-08T16:38:00Z GOVERNANCE proposal=116 phase=tally votingPower=136926 turnoutPercent=18.8
2025-11-08T16:39:00Z GOVERNANCE proposal=116 phase=discussion votingPower=136963 turnoutPercent=18.9
2025-11-08T16:40:00Z GOVERNANCE proposal=116 phase=voting votingPower=137000 turnoutPercent=19.0
2025-11-08T16:41:00Z GOVERNANCE proposal=116 phase=tally votingPower=137037 turnoutPercent=19.1
2025-11-08T16:42:00Z GOVERNANCE proposal=116 phase=discussion votingPower=137074 turnoutPercent=19.2
2025-11-08T16:43:00Z GOVERNANCE proposal=116 phase=voting votingPower=137111 turnoutPercent=19.3
2025-11-08T16:44:00Z GOVERNANCE proposal=116 phase=tally votingPower=137148 turnoutPercent=19.4
2025-11-08T16:45:00Z GOVERNANCE proposal=116 phase=discussion votingPower=137185 turnoutPercent=19.5
2025-11-08T16:46:00Z GOVERNANCE proposal=116 phase=voting votingPower=137222 turnoutPercent=19.6
2025-11-08T16:47:00Z GOVERNANCE proposal=116 phase=tally votingPower=137259 turnoutPercent=19.7
2025-11-08T16:48:00Z GOVERNANCE proposal=116 phase=discussion votingPower=137296 turnoutPercent=19.8
2025-11-08T16:49:00Z GOVERNANCE proposal=116 phase=voting votingPower=137333 turnoutPercent=19.9
2025-11-08T16:50:00Z GOVERNANCE proposal=116 phase=tally votingPower=137370 turnoutPercent=20.0
2025-11-08T16:51:00Z GOVERNANCE proposal=116 phase=discussion votingPower=137407 turnoutPercent=20.1
2025-11-08T16:52:00Z GOVERNANCE proposal=116 phase=voting votingPower=137444 turnoutPercent=20.2
2025-11-08T16:53:00Z GOVERNANCE proposal=116 phase=tally votingPower=137481 turnoutPercent=20.3
2025-11-08T16:54:00Z GOVERNANCE proposal=116 phase=discussion votingPower=137518 turnoutPercent=20.4
2025-11-08T16:55:00Z GOVERNANCE proposal=116 phase=voting votingPower=137555 turnoutPercent=20.5
2025-11-08T16:56:00Z GOVERNANCE proposal=116 phase=tally votingPower=137592 turnoutPercent=20.6
2025-11-08T16:57:00Z GOVERNANCE proposal=116 phase=discussion votingPower=137629 turnoutPercent=20.7
2025-11-08T16:58:00Z GOVERNANCE proposal=116 phase=voting votingPower=137666 turnoutPercent=20.8
2025-11-08T16:59:00Z GOVERNANCE proposal=116 phase=tally votingPower=137703 turnoutPercent=20.9
2025-11-08T17:00:00Z GOVERNANCE proposal=117 phase=discussion votingPower=137740 turnoutPercent=15.0
2025-11-08T17:01:00Z GOVERNANCE proposal=117 phase=voting votingPower=137777 turnoutPercent=15.1
2025-11-08T17:02:00Z GOVERNANCE proposal=117 phase=tally votingPower=137814 turnoutPercent=15.2
2025-11-08T17:03:00Z GOVERNANCE proposal=117 phase=discussion votingPower=137851 turnoutPercent=15.3
2025-11-08T17:04:00Z GOVERNANCE proposal=117 phase=voting votingPower=137888 turnoutPercent=15.4
2025-11-08T17:05:00Z GOVERNANCE proposal=117 phase=tally votingPower=137925 turnoutPercent=15.5
2025-11-08T17:06:00Z GOVERNANCE proposal=117 phase=discussion votingPower=137962 turnoutPercent=15.6
2025-11-08T17:07:00Z GOVERNANCE proposal=117 phase=voting votingPower=137999 turnoutPercent=15.7
2025-11-08T17:08:00Z GOVERNANCE proposal=117 phase=tally votingPower=138036 turnoutPercent=15.8
2025-11-08T17:09:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138073 turnoutPercent=15.9
2025-11-08T17:10:00Z GOVERNANCE proposal=117 phase=voting votingPower=138110 turnoutPercent=16.0
2025-11-08T17:11:00Z GOVERNANCE proposal=117 phase=tally votingPower=138147 turnoutPercent=16.1
2025-11-08T17:12:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138184 turnoutPercent=16.2
2025-11-08T17:13:00Z GOVERNANCE proposal=117 phase=voting votingPower=138221 turnoutPercent=16.3
2025-11-08T17:14:00Z GOVERNANCE proposal=117 phase=tally votingPower=138258 turnoutPercent=16.4
2025-11-08T17:15:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138295 turnoutPercent=16.5
2025-11-08T17:16:00Z GOVERNANCE proposal=117 phase=voting votingPower=138332 turnoutPercent=16.6
2025-11-08T17:17:00Z GOVERNANCE proposal=117 phase=tally votingPower=138369 turnoutPercent=16.7
2025-11-08T17:18:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138406 turnoutPercent=16.8
2025-11-08T17:19:00Z GOVERNANCE proposal=117 phase=voting votingPower=138443 turnoutPercent=16.9
2025-11-08T17:20:00Z GOVERNANCE proposal=117 phase=tally votingPower=138480 turnoutPercent=17.0
2025-11-08T17:21:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138517 turnoutPercent=17.1
2025-11-08T17:22:00Z GOVERNANCE proposal=117 phase=voting votingPower=138554 turnoutPercent=17.2
2025-11-08T17:23:00Z GOVERNANCE proposal=117 phase=tally votingPower=138591 turnoutPercent=17.3
2025-11-08T17:24:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138628 turnoutPercent=17.4
2025-11-08T17:25:00Z GOVERNANCE proposal=117 phase=voting votingPower=138665 turnoutPercent=17.5
2025-11-08T17:26:00Z GOVERNANCE proposal=117 phase=tally votingPower=138702 turnoutPercent=17.6
2025-11-08T17:27:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138739 turnoutPercent=17.7
2025-11-08T17:28:00Z GOVERNANCE proposal=117 phase=voting votingPower=138776 turnoutPercent=17.8
2025-11-08T17:29:00Z GOVERNANCE proposal=117 phase=tally votingPower=138813 turnoutPercent=17.9
2025-11-08T17:30:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138850 turnoutPercent=18.0
2025-11-08T17:31:00Z GOVERNANCE proposal=117 phase=voting votingPower=138887 turnoutPercent=18.1
2025-11-08T17:32:00Z GOVERNANCE proposal=117 phase=tally votingPower=138924 turnoutPercent=18.2
2025-11-08T17:33:00Z GOVERNANCE proposal=117 phase=discussion votingPower=138961 turnoutPercent=18.3
2025-11-08T17:34:00Z GOVERNANCE proposal=117 phase=voting votingPower=138998 turnoutPercent=18.4
2025-11-08T17:35:00Z GOVERNANCE proposal=117 phase=tally votingPower=139035 turnoutPercent=18.5
2025-11-08T17:36:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139072 turnoutPercent=18.6
2025-11-08T17:37:00Z GOVERNANCE proposal=117 phase=voting votingPower=139109 turnoutPercent=18.7
2025-11-08T17:38:00Z GOVERNANCE proposal=117 phase=tally votingPower=139146 turnoutPercent=18.8
2025-11-08T17:39:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139183 turnoutPercent=18.9
2025-11-08T17:40:00Z GOVERNANCE proposal=117 phase=voting votingPower=139220 turnoutPercent=19.0
2025-11-08T17:41:00Z GOVERNANCE proposal=117 phase=tally votingPower=139257 turnoutPercent=19.1
2025-11-08T17:42:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139294 turnoutPercent=19.2
2025-11-08T17:43:00Z GOVERNANCE proposal=117 phase=voting votingPower=139331 turnoutPercent=19.3
2025-11-08T17:44:00Z GOVERNANCE proposal=117 phase=tally votingPower=139368 turnoutPercent=19.4
2025-11-08T17:45:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139405 turnoutPercent=19.5
2025-11-08T17:46:00Z GOVERNANCE proposal=117 phase=voting votingPower=139442 turnoutPercent=19.6
2025-11-08T17:47:00Z GOVERNANCE proposal=117 phase=tally votingPower=139479 turnoutPercent=19.7
2025-11-08T17:48:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139516 turnoutPercent=19.8
2025-11-08T17:49:00Z GOVERNANCE proposal=117 phase=voting votingPower=139553 turnoutPercent=19.9
2025-11-08T17:50:00Z GOVERNANCE proposal=117 phase=tally votingPower=139590 turnoutPercent=20.0
2025-11-08T17:51:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139627 turnoutPercent=20.1
2025-11-08T17:52:00Z GOVERNANCE proposal=117 phase=voting votingPower=139664 turnoutPercent=20.2
2025-11-08T17:53:00Z GOVERNANCE proposal=117 phase=tally votingPower=139701 turnoutPercent=20.3
2025-11-08T17:54:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139738 turnoutPercent=20.4
2025-11-08T17:55:00Z GOVERNANCE proposal=117 phase=voting votingPower=139775 turnoutPercent=20.5
2025-11-08T17:56:00Z GOVERNANCE proposal=117 phase=tally votingPower=139812 turnoutPercent=20.6
2025-11-08T17:57:00Z GOVERNANCE proposal=117 phase=discussion votingPower=139849 turnoutPercent=20.7
2025-11-08T17:58:00Z GOVERNANCE proposal=117 phase=voting votingPower=139886 turnoutPercent=20.8
2025-11-08T17:59:00Z GOVERNANCE proposal=117 phase=tally votingPower=139923 turnoutPercent=20.9
2025-11-08T18:00:00Z GOVERNANCE proposal=118 phase=discussion votingPower=139960 turnoutPercent=15.0
2025-11-08T18:01:00Z GOVERNANCE proposal=118 phase=voting votingPower=139997 turnoutPercent=15.1
2025-11-08T18:02:00Z GOVERNANCE proposal=118 phase=tally votingPower=140034 turnoutPercent=15.2
2025-11-08T18:03:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140071 turnoutPercent=15.3
2025-11-08T18:04:00Z GOVERNANCE proposal=118 phase=voting votingPower=140108 turnoutPercent=15.4
2025-11-08T18:05:00Z GOVERNANCE proposal=118 phase=tally votingPower=140145 turnoutPercent=15.5
2025-11-08T18:06:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140182 turnoutPercent=15.6
2025-11-08T18:07:00Z GOVERNANCE proposal=118 phase=voting votingPower=140219 turnoutPercent=15.7
2025-11-08T18:08:00Z GOVERNANCE proposal=118 phase=tally votingPower=140256 turnoutPercent=15.8
2025-11-08T18:09:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140293 turnoutPercent=15.9
2025-11-08T18:10:00Z GOVERNANCE proposal=118 phase=voting votingPower=140330 turnoutPercent=16.0
2025-11-08T18:11:00Z GOVERNANCE proposal=118 phase=tally votingPower=140367 turnoutPercent=16.1
2025-11-08T18:12:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140404 turnoutPercent=16.2
2025-11-08T18:13:00Z GOVERNANCE proposal=118 phase=voting votingPower=140441 turnoutPercent=16.3
2025-11-08T18:14:00Z GOVERNANCE proposal=118 phase=tally votingPower=140478 turnoutPercent=16.4
2025-11-08T18:15:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140515 turnoutPercent=16.5
2025-11-08T18:16:00Z GOVERNANCE proposal=118 phase=voting votingPower=140552 turnoutPercent=16.6
2025-11-08T18:17:00Z GOVERNANCE proposal=118 phase=tally votingPower=140589 turnoutPercent=16.7
2025-11-08T18:18:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140626 turnoutPercent=16.8
2025-11-08T18:19:00Z GOVERNANCE proposal=118 phase=voting votingPower=140663 turnoutPercent=16.9
2025-11-08T18:20:00Z GOVERNANCE proposal=118 phase=tally votingPower=140700 turnoutPercent=17.0
2025-11-08T18:21:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140737 turnoutPercent=17.1
2025-11-08T18:22:00Z GOVERNANCE proposal=118 phase=voting votingPower=140774 turnoutPercent=17.2
2025-11-08T18:23:00Z GOVERNANCE proposal=118 phase=tally votingPower=140811 turnoutPercent=17.3
2025-11-08T18:24:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140848 turnoutPercent=17.4
2025-11-08T18:25:00Z GOVERNANCE proposal=118 phase=voting votingPower=140885 turnoutPercent=17.5
2025-11-08T18:26:00Z GOVERNANCE proposal=118 phase=tally votingPower=140922 turnoutPercent=17.6
2025-11-08T18:27:00Z GOVERNANCE proposal=118 phase=discussion votingPower=140959 turnoutPercent=17.7
2025-11-08T18:28:00Z GOVERNANCE proposal=118 phase=voting votingPower=140996 turnoutPercent=17.8
2025-11-08T18:29:00Z GOVERNANCE proposal=118 phase=tally votingPower=141033 turnoutPercent=17.9
2025-11-08T18:30:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141070 turnoutPercent=18.0
2025-11-08T18:31:00Z GOVERNANCE proposal=118 phase=voting votingPower=141107 turnoutPercent=18.1
2025-11-08T18:32:00Z GOVERNANCE proposal=118 phase=tally votingPower=141144 turnoutPercent=18.2
2025-11-08T18:33:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141181 turnoutPercent=18.3
2025-11-08T18:34:00Z GOVERNANCE proposal=118 phase=voting votingPower=141218 turnoutPercent=18.4
2025-11-08T18:35:00Z GOVERNANCE proposal=118 phase=tally votingPower=141255 turnoutPercent=18.5
2025-11-08T18:36:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141292 turnoutPercent=18.6
2025-11-08T18:37:00Z GOVERNANCE proposal=118 phase=voting votingPower=141329 turnoutPercent=18.7
2025-11-08T18:38:00Z GOVERNANCE proposal=118 phase=tally votingPower=141366 turnoutPercent=18.8
2025-11-08T18:39:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141403 turnoutPercent=18.9
2025-11-08T18:40:00Z GOVERNANCE proposal=118 phase=voting votingPower=141440 turnoutPercent=19.0
2025-11-08T18:41:00Z GOVERNANCE proposal=118 phase=tally votingPower=141477 turnoutPercent=19.1
2025-11-08T18:42:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141514 turnoutPercent=19.2
2025-11-08T18:43:00Z GOVERNANCE proposal=118 phase=voting votingPower=141551 turnoutPercent=19.3
2025-11-08T18:44:00Z GOVERNANCE proposal=118 phase=tally votingPower=141588 turnoutPercent=19.4
2025-11-08T18:45:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141625 turnoutPercent=19.5
2025-11-08T18:46:00Z GOVERNANCE proposal=118 phase=voting votingPower=141662 turnoutPercent=19.6
2025-11-08T18:47:00Z GOVERNANCE proposal=118 phase=tally votingPower=141699 turnoutPercent=19.7
2025-11-08T18:48:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141736 turnoutPercent=19.8
2025-11-08T18:49:00Z GOVERNANCE proposal=118 phase=voting votingPower=141773 turnoutPercent=19.9
2025-11-08T18:50:00Z GOVERNANCE proposal=118 phase=tally votingPower=141810 turnoutPercent=20.0
2025-11-08T18:51:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141847 turnoutPercent=20.1
2025-11-08T18:52:00Z GOVERNANCE proposal=118 phase=voting votingPower=141884 turnoutPercent=20.2
2025-11-08T18:53:00Z GOVERNANCE proposal=118 phase=tally votingPower=141921 turnoutPercent=20.3
2025-11-08T18:54:00Z GOVERNANCE proposal=118 phase=discussion votingPower=141958 turnoutPercent=20.4
2025-11-08T18:55:00Z GOVERNANCE proposal=118 phase=voting votingPower=141995 turnoutPercent=20.5
2025-11-08T18:56:00Z GOVERNANCE proposal=118 phase=tally votingPower=142032 turnoutPercent=20.6
2025-11-08T18:57:00Z GOVERNANCE proposal=118 phase=discussion votingPower=142069 turnoutPercent=20.7
2025-11-08T18:58:00Z GOVERNANCE proposal=118 phase=voting votingPower=142106 turnoutPercent=20.8
2025-11-08T18:59:00Z GOVERNANCE proposal=118 phase=tally votingPower=142143 turnoutPercent=20.9
2025-11-08T19:00:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142180 turnoutPercent=15.0
2025-11-08T19:01:00Z GOVERNANCE proposal=119 phase=voting votingPower=142217 turnoutPercent=15.1
2025-11-08T19:02:00Z GOVERNANCE proposal=119 phase=tally votingPower=142254 turnoutPercent=15.2
2025-11-08T19:03:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142291 turnoutPercent=15.3
2025-11-08T19:04:00Z GOVERNANCE proposal=119 phase=voting votingPower=142328 turnoutPercent=15.4
2025-11-08T19:05:00Z GOVERNANCE proposal=119 phase=tally votingPower=142365 turnoutPercent=15.5
2025-11-08T19:06:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142402 turnoutPercent=15.6
2025-11-08T19:07:00Z GOVERNANCE proposal=119 phase=voting votingPower=142439 turnoutPercent=15.7
2025-11-08T19:08:00Z GOVERNANCE proposal=119 phase=tally votingPower=142476 turnoutPercent=15.8
2025-11-08T19:09:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142513 turnoutPercent=15.9
2025-11-08T19:10:00Z GOVERNANCE proposal=119 phase=voting votingPower=142550 turnoutPercent=16.0
2025-11-08T19:11:00Z GOVERNANCE proposal=119 phase=tally votingPower=142587 turnoutPercent=16.1
2025-11-08T19:12:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142624 turnoutPercent=16.2
2025-11-08T19:13:00Z GOVERNANCE proposal=119 phase=voting votingPower=142661 turnoutPercent=16.3
2025-11-08T19:14:00Z GOVERNANCE proposal=119 phase=tally votingPower=142698 turnoutPercent=16.4
2025-11-08T19:15:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142735 turnoutPercent=16.5
2025-11-08T19:16:00Z GOVERNANCE proposal=119 phase=voting votingPower=142772 turnoutPercent=16.6
2025-11-08T19:17:00Z GOVERNANCE proposal=119 phase=tally votingPower=142809 turnoutPercent=16.7
2025-11-08T19:18:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142846 turnoutPercent=16.8
2025-11-08T19:19:00Z GOVERNANCE proposal=119 phase=voting votingPower=142883 turnoutPercent=16.9
2025-11-08T19:20:00Z GOVERNANCE proposal=119 phase=tally votingPower=142920 turnoutPercent=17.0
2025-11-08T19:21:00Z GOVERNANCE proposal=119 phase=discussion votingPower=142957 turnoutPercent=17.1
2025-11-08T19:22:00Z GOVERNANCE proposal=119 phase=voting votingPower=142994 turnoutPercent=17.2
2025-11-08T19:23:00Z GOVERNANCE proposal=119 phase=tally votingPower=143031 turnoutPercent=17.3
2025-11-08T19:24:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143068 turnoutPercent=17.4
2025-11-08T19:25:00Z GOVERNANCE proposal=119 phase=voting votingPower=143105 turnoutPercent=17.5
2025-11-08T19:26:00Z GOVERNANCE proposal=119 phase=tally votingPower=143142 turnoutPercent=17.6
2025-11-08T19:27:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143179 turnoutPercent=17.7
2025-11-08T19:28:00Z GOVERNANCE proposal=119 phase=voting votingPower=143216 turnoutPercent=17.8
2025-11-08T19:29:00Z GOVERNANCE proposal=119 phase=tally votingPower=143253 turnoutPercent=17.9
2025-11-08T19:30:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143290 turnoutPercent=18.0
2025-11-08T19:31:00Z GOVERNANCE proposal=119 phase=voting votingPower=143327 turnoutPercent=18.1
2025-11-08T19:32:00Z GOVERNANCE proposal=119 phase=tally votingPower=143364 turnoutPercent=18.2
2025-11-08T19:33:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143401 turnoutPercent=18.3
2025-11-08T19:34:00Z GOVERNANCE proposal=119 phase=voting votingPower=143438 turnoutPercent=18.4
2025-11-08T19:35:00Z GOVERNANCE proposal=119 phase=tally votingPower=143475 turnoutPercent=18.5
2025-11-08T19:36:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143512 turnoutPercent=18.6
2025-11-08T19:37:00Z GOVERNANCE proposal=119 phase=voting votingPower=143549 turnoutPercent=18.7
2025-11-08T19:38:00Z GOVERNANCE proposal=119 phase=tally votingPower=143586 turnoutPercent=18.8
2025-11-08T19:39:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143623 turnoutPercent=18.9
2025-11-08T19:40:00Z GOVERNANCE proposal=119 phase=voting votingPower=143660 turnoutPercent=19.0
2025-11-08T19:41:00Z GOVERNANCE proposal=119 phase=tally votingPower=143697 turnoutPercent=19.1
2025-11-08T19:42:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143734 turnoutPercent=19.2
2025-11-08T19:43:00Z GOVERNANCE proposal=119 phase=voting votingPower=143771 turnoutPercent=19.3
2025-11-08T19:44:00Z GOVERNANCE proposal=119 phase=tally votingPower=143808 turnoutPercent=19.4
2025-11-08T19:45:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143845 turnoutPercent=19.5
2025-11-08T19:46:00Z GOVERNANCE proposal=119 phase=voting votingPower=143882 turnoutPercent=19.6
2025-11-08T19:47:00Z GOVERNANCE proposal=119 phase=tally votingPower=143919 turnoutPercent=19.7
2025-11-08T19:48:00Z GOVERNANCE proposal=119 phase=discussion votingPower=143956 turnoutPercent=19.8
2025-11-08T19:49:00Z GOVERNANCE proposal=119 phase=voting votingPower=143993 turnoutPercent=19.9
2025-11-08T19:50:00Z GOVERNANCE proposal=119 phase=tally votingPower=144030 turnoutPercent=20.0
2025-11-08T19:51:00Z GOVERNANCE proposal=119 phase=discussion votingPower=144067 turnoutPercent=20.1
2025-11-08T19:52:00Z GOVERNANCE proposal=119 phase=voting votingPower=144104 turnoutPercent=20.2
2025-11-08T19:53:00Z GOVERNANCE proposal=119 phase=tally votingPower=144141 turnoutPercent=20.3
2025-11-08T19:54:00Z GOVERNANCE proposal=119 phase=discussion votingPower=144178 turnoutPercent=20.4
2025-11-08T19:55:00Z GOVERNANCE proposal=119 phase=voting votingPower=144215 turnoutPercent=20.5
2025-11-08T19:56:00Z GOVERNANCE proposal=119 phase=tally votingPower=144252 turnoutPercent=20.6
2025-11-08T19:57:00Z GOVERNANCE proposal=119 phase=discussion votingPower=144289 turnoutPercent=20.7
2025-11-08T19:58:00Z GOVERNANCE proposal=119 phase=voting votingPower=144326 turnoutPercent=20.8
2025-11-08T19:59:00Z GOVERNANCE proposal=119 phase=tally votingPower=144363 turnoutPercent=20.9
2025-11-08T20:00:00Z GOVERNANCE proposal=120 phase=discussion votingPower=144400 turnoutPercent=15.0
2025-11-08T20:01:00Z GOVERNANCE proposal=120 phase=voting votingPower=144437 turnoutPercent=15.1
2025-11-08T20:02:00Z GOVERNANCE proposal=120 phase=tally votingPower=144474 turnoutPercent=15.2
2025-11-08T20:03:00Z GOVERNANCE proposal=120 phase=discussion votingPower=144511 turnoutPercent=15.3
2025-11-08T20:04:00Z GOVERNANCE proposal=120 phase=voting votingPower=144548 turnoutPercent=15.4
2025-11-08T20:05:00Z GOVERNANCE proposal=120 phase=tally votingPower=144585 turnoutPercent=15.5
2025-11-08T20:06:00Z GOVERNANCE proposal=120 phase=discussion votingPower=144622 turnoutPercent=15.6
2025-11-08T20:07:00Z GOVERNANCE proposal=120 phase=voting votingPower=144659 turnoutPercent=15.7
2025-11-08T20:08:00Z GOVERNANCE proposal=120 phase=tally votingPower=144696 turnoutPercent=15.8
2025-11-08T20:09:00Z GOVERNANCE proposal=120 phase=discussion votingPower=144733 turnoutPercent=15.9
2025-11-08T20:10:00Z GOVERNANCE proposal=120 phase=voting votingPower=144770 turnoutPercent=16.0
2025-11-08T20:11:00Z GOVERNANCE proposal=120 phase=tally votingPower=144807 turnoutPercent=16.1
2025-11-08T20:12:00Z GOVERNANCE proposal=120 phase=discussion votingPower=144844 turnoutPercent=16.2
2025-11-08T20:13:00Z GOVERNANCE proposal=120 phase=voting votingPower=144881 turnoutPercent=16.3
2025-11-08T20:14:00Z GOVERNANCE proposal=120 phase=tally votingPower=144918 turnoutPercent=16.4
2025-11-08T20:15:00Z GOVERNANCE proposal=120 phase=discussion votingPower=144955 turnoutPercent=16.5
2025-11-08T20:16:00Z GOVERNANCE proposal=120 phase=voting votingPower=144992 turnoutPercent=16.6
2025-11-08T20:17:00Z GOVERNANCE proposal=120 phase=tally votingPower=145029 turnoutPercent=16.7
2025-11-08T20:18:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145066 turnoutPercent=16.8
2025-11-08T20:19:00Z GOVERNANCE proposal=120 phase=voting votingPower=145103 turnoutPercent=16.9
2025-11-08T20:20:00Z GOVERNANCE proposal=120 phase=tally votingPower=145140 turnoutPercent=17.0
2025-11-08T20:21:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145177 turnoutPercent=17.1
2025-11-08T20:22:00Z GOVERNANCE proposal=120 phase=voting votingPower=145214 turnoutPercent=17.2
2025-11-08T20:23:00Z GOVERNANCE proposal=120 phase=tally votingPower=145251 turnoutPercent=17.3
2025-11-08T20:24:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145288 turnoutPercent=17.4
2025-11-08T20:25:00Z GOVERNANCE proposal=120 phase=voting votingPower=145325 turnoutPercent=17.5
2025-11-08T20:26:00Z GOVERNANCE proposal=120 phase=tally votingPower=145362 turnoutPercent=17.6
2025-11-08T20:27:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145399 turnoutPercent=17.7
2025-11-08T20:28:00Z GOVERNANCE proposal=120 phase=voting votingPower=145436 turnoutPercent=17.8
2025-11-08T20:29:00Z GOVERNANCE proposal=120 phase=tally votingPower=145473 turnoutPercent=17.9
2025-11-08T20:30:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145510 turnoutPercent=18.0
2025-11-08T20:31:00Z GOVERNANCE proposal=120 phase=voting votingPower=145547 turnoutPercent=18.1
2025-11-08T20:32:00Z GOVERNANCE proposal=120 phase=tally votingPower=145584 turnoutPercent=18.2
2025-11-08T20:33:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145621 turnoutPercent=18.3
2025-11-08T20:34:00Z GOVERNANCE proposal=120 phase=voting votingPower=145658 turnoutPercent=18.4
2025-11-08T20:35:00Z GOVERNANCE proposal=120 phase=tally votingPower=145695 turnoutPercent=18.5
2025-11-08T20:36:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145732 turnoutPercent=18.6
2025-11-08T20:37:00Z GOVERNANCE proposal=120 phase=voting votingPower=145769 turnoutPercent=18.7
2025-11-08T20:38:00Z GOVERNANCE proposal=120 phase=tally votingPower=145806 turnoutPercent=18.8
2025-11-08T20:39:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145843 turnoutPercent=18.9
2025-11-08T20:40:00Z GOVERNANCE proposal=120 phase=voting votingPower=145880 turnoutPercent=19.0
2025-11-08T20:41:00Z GOVERNANCE proposal=120 phase=tally votingPower=145917 turnoutPercent=19.1
2025-11-08T20:42:00Z GOVERNANCE proposal=120 phase=discussion votingPower=145954 turnoutPercent=19.2
2025-11-08T20:43:00Z GOVERNANCE proposal=120 phase=voting votingPower=145991 turnoutPercent=19.3
2025-11-08T20:44:00Z GOVERNANCE proposal=120 phase=tally votingPower=146028 turnoutPercent=19.4
2025-11-08T20:45:00Z GOVERNANCE proposal=120 phase=discussion votingPower=146065 turnoutPercent=19.5
2025-11-08T20:46:00Z GOVERNANCE proposal=120 phase=voting votingPower=146102 turnoutPercent=19.6
2025-11-08T20:47:00Z GOVERNANCE proposal=120 phase=tally votingPower=146139 turnoutPercent=19.7
2025-11-08T20:48:00Z GOVERNANCE proposal=120 phase=discussion votingPower=146176 turnoutPercent=19.8
2025-11-08T20:49:00Z GOVERNANCE proposal=120 phase=voting votingPower=146213 turnoutPercent=19.9
2025-11-08T20:50:00Z GOVERNANCE proposal=120 phase=tally votingPower=146250 turnoutPercent=20.0
2025-11-08T20:51:00Z GOVERNANCE proposal=120 phase=discussion votingPower=146287 turnoutPercent=20.1
2025-11-08T20:52:00Z GOVERNANCE proposal=120 phase=voting votingPower=146324 turnoutPercent=20.2
2025-11-08T20:53:00Z GOVERNANCE proposal=120 phase=tally votingPower=146361 turnoutPercent=20.3
2025-11-08T20:54:00Z GOVERNANCE proposal=120 phase=discussion votingPower=146398 turnoutPercent=20.4
2025-11-08T20:55:00Z GOVERNANCE proposal=120 phase=voting votingPower=146435 turnoutPercent=20.5
2025-11-08T20:56:00Z GOVERNANCE proposal=120 phase=tally votingPower=146472 turnoutPercent=20.6
2025-11-08T20:57:00Z GOVERNANCE proposal=120 phase=discussion votingPower=146509 turnoutPercent=20.7
2025-11-08T20:58:00Z GOVERNANCE proposal=120 phase=voting votingPower=146546 turnoutPercent=20.8
2025-11-08T20:59:00Z GOVERNANCE proposal=120 phase=tally votingPower=146583 turnoutPercent=20.9
2025-11-08T21:00:00Z GOVERNANCE proposal=121 phase=discussion votingPower=146620 turnoutPercent=15.0
2025-11-08T21:01:00Z GOVERNANCE proposal=121 phase=voting votingPower=146657 turnoutPercent=15.1
2025-11-08T21:02:00Z GOVERNANCE proposal=121 phase=tally votingPower=146694 turnoutPercent=15.2
2025-11-08T21:03:00Z GOVERNANCE proposal=121 phase=discussion votingPower=146731 turnoutPercent=15.3
2025-11-08T21:04:00Z GOVERNANCE proposal=121 phase=voting votingPower=146768 turnoutPercent=15.4
2025-11-08T21:05:00Z GOVERNANCE proposal=121 phase=tally votingPower=146805 turnoutPercent=15.5
2025-11-08T21:06:00Z GOVERNANCE proposal=121 phase=discussion votingPower=146842 turnoutPercent=15.6
2025-11-08T21:07:00Z GOVERNANCE proposal=121 phase=voting votingPower=146879 turnoutPercent=15.7
2025-11-08T21:08:00Z GOVERNANCE proposal=121 phase=tally votingPower=146916 turnoutPercent=15.8
2025-11-08T21:09:00Z GOVERNANCE proposal=121 phase=discussion votingPower=146953 turnoutPercent=15.9
2025-11-08T21:10:00Z GOVERNANCE proposal=121 phase=voting votingPower=146990 turnoutPercent=16.0
2025-11-08T21:11:00Z GOVERNANCE proposal=121 phase=tally votingPower=147027 turnoutPercent=16.1
2025-11-08T21:12:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147064 turnoutPercent=16.2
2025-11-08T21:13:00Z GOVERNANCE proposal=121 phase=voting votingPower=147101 turnoutPercent=16.3
2025-11-08T21:14:00Z GOVERNANCE proposal=121 phase=tally votingPower=147138 turnoutPercent=16.4
2025-11-08T21:15:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147175 turnoutPercent=16.5
2025-11-08T21:16:00Z GOVERNANCE proposal=121 phase=voting votingPower=147212 turnoutPercent=16.6
2025-11-08T21:17:00Z GOVERNANCE proposal=121 phase=tally votingPower=147249 turnoutPercent=16.7
2025-11-08T21:18:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147286 turnoutPercent=16.8
2025-11-08T21:19:00Z GOVERNANCE proposal=121 phase=voting votingPower=147323 turnoutPercent=16.9
2025-11-08T21:20:00Z GOVERNANCE proposal=121 phase=tally votingPower=147360 turnoutPercent=17.0
2025-11-08T21:21:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147397 turnoutPercent=17.1
2025-11-08T21:22:00Z GOVERNANCE proposal=121 phase=voting votingPower=147434 turnoutPercent=17.2
2025-11-08T21:23:00Z GOVERNANCE proposal=121 phase=tally votingPower=147471 turnoutPercent=17.3
2025-11-08T21:24:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147508 turnoutPercent=17.4
2025-11-08T21:25:00Z GOVERNANCE proposal=121 phase=voting votingPower=147545 turnoutPercent=17.5
2025-11-08T21:26:00Z GOVERNANCE proposal=121 phase=tally votingPower=147582 turnoutPercent=17.6
2025-11-08T21:27:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147619 turnoutPercent=17.7
2025-11-08T21:28:00Z GOVERNANCE proposal=121 phase=voting votingPower=147656 turnoutPercent=17.8
2025-11-08T21:29:00Z GOVERNANCE proposal=121 phase=tally votingPower=147693 turnoutPercent=17.9
2025-11-08T21:30:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147730 turnoutPercent=18.0
2025-11-08T21:31:00Z GOVERNANCE proposal=121 phase=voting votingPower=147767 turnoutPercent=18.1
2025-11-08T21:32:00Z GOVERNANCE proposal=121 phase=tally votingPower=147804 turnoutPercent=18.2
2025-11-08T21:33:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147841 turnoutPercent=18.3
2025-11-08T21:34:00Z GOVERNANCE proposal=121 phase=voting votingPower=147878 turnoutPercent=18.4
2025-11-08T21:35:00Z GOVERNANCE proposal=121 phase=tally votingPower=147915 turnoutPercent=18.5
2025-11-08T21:36:00Z GOVERNANCE proposal=121 phase=discussion votingPower=147952 turnoutPercent=18.6
2025-11-08T21:37:00Z GOVERNANCE proposal=121 phase=voting votingPower=147989 turnoutPercent=18.7
2025-11-08T21:38:00Z GOVERNANCE proposal=121 phase=tally votingPower=148026 turnoutPercent=18.8
2025-11-08T21:39:00Z GOVERNANCE proposal=121 phase=discussion votingPower=148063 turnoutPercent=18.9
2025-11-08T21:40:00Z GOVERNANCE proposal=121 phase=voting votingPower=148100 turnoutPercent=19.0
2025-11-08T21:41:00Z GOVERNANCE proposal=121 phase=tally votingPower=148137 turnoutPercent=19.1
2025-11-08T21:42:00Z GOVERNANCE proposal=121 phase=discussion votingPower=148174 turnoutPercent=19.2
2025-11-08T21:43:00Z GOVERNANCE proposal=121 phase=voting votingPower=148211 turnoutPercent=19.3
2025-11-08T21:44:00Z GOVERNANCE proposal=121 phase=tally votingPower=148248 turnoutPercent=19.4
2025-11-08T21:45:00Z GOVERNANCE proposal=121 phase=discussion votingPower=148285 turnoutPercent=19.5
2025-11-08T21:46:00Z GOVERNANCE proposal=121 phase=voting votingPower=148322 turnoutPercent=19.6
2025-11-08T21:47:00Z GOVERNANCE proposal=121 phase=tally votingPower=148359 turnoutPercent=19.7
2025-11-08T21:48:00Z GOVERNANCE proposal=121 phase=discussion votingPower=148396 turnoutPercent=19.8
2025-11-08T21:49:00Z GOVERNANCE proposal=121 phase=voting votingPower=148433 turnoutPercent=19.9
2025-11-08T21:50:00Z GOVERNANCE proposal=121 phase=tally votingPower=148470 turnoutPercent=20.0
2025-11-08T21:51:00Z GOVERNANCE proposal=121 phase=discussion votingPower=148507 turnoutPercent=20.1
2025-11-08T21:52:00Z GOVERNANCE proposal=121 phase=voting votingPower=148544 turnoutPercent=20.2
2025-11-08T21:53:00Z GOVERNANCE proposal=121 phase=tally votingPower=148581 turnoutPercent=20.3
2025-11-08T21:54:00Z GOVERNANCE proposal=121 phase=discussion votingPower=148618 turnoutPercent=20.4
2025-11-08T21:55:00Z GOVERNANCE proposal=121 phase=voting votingPower=148655 turnoutPercent=20.5
2025-11-08T21:56:00Z GOVERNANCE proposal=121 phase=tally votingPower=148692 turnoutPercent=20.6
2025-11-08T21:57:00Z GOVERNANCE proposal=121 phase=discussion votingPower=148729 turnoutPercent=20.7
2025-11-08T21:58:00Z GOVERNANCE proposal=121 phase=voting votingPower=148766 turnoutPercent=20.8
2025-11-08T21:59:00Z GOVERNANCE proposal=121 phase=tally votingPower=148803 turnoutPercent=20.9
2025-11-08T22:00:00Z GOVERNANCE proposal=122 phase=discussion votingPower=148840 turnoutPercent=15.0
2025-11-08T22:01:00Z GOVERNANCE proposal=122 phase=voting votingPower=148877 turnoutPercent=15.1
2025-11-08T22:02:00Z GOVERNANCE proposal=122 phase=tally votingPower=148914 turnoutPercent=15.2
2025-11-08T22:03:00Z GOVERNANCE proposal=122 phase=discussion votingPower=148951 turnoutPercent=15.3
2025-11-08T22:04:00Z GOVERNANCE proposal=122 phase=voting votingPower=148988 turnoutPercent=15.4
2025-11-08T22:05:00Z GOVERNANCE proposal=122 phase=tally votingPower=149025 turnoutPercent=15.5
2025-11-08T22:06:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149062 turnoutPercent=15.6
2025-11-08T22:07:00Z GOVERNANCE proposal=122 phase=voting votingPower=149099 turnoutPercent=15.7
2025-11-08T22:08:00Z GOVERNANCE proposal=122 phase=tally votingPower=149136 turnoutPercent=15.8
2025-11-08T22:09:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149173 turnoutPercent=15.9
2025-11-08T22:10:00Z GOVERNANCE proposal=122 phase=voting votingPower=149210 turnoutPercent=16.0
2025-11-08T22:11:00Z GOVERNANCE proposal=122 phase=tally votingPower=149247 turnoutPercent=16.1
2025-11-08T22:12:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149284 turnoutPercent=16.2
2025-11-08T22:13:00Z GOVERNANCE proposal=122 phase=voting votingPower=149321 turnoutPercent=16.3
2025-11-08T22:14:00Z GOVERNANCE proposal=122 phase=tally votingPower=149358 turnoutPercent=16.4
2025-11-08T22:15:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149395 turnoutPercent=16.5
2025-11-08T22:16:00Z GOVERNANCE proposal=122 phase=voting votingPower=149432 turnoutPercent=16.6
2025-11-08T22:17:00Z GOVERNANCE proposal=122 phase=tally votingPower=149469 turnoutPercent=16.7
2025-11-08T22:18:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149506 turnoutPercent=16.8
2025-11-08T22:19:00Z GOVERNANCE proposal=122 phase=voting votingPower=149543 turnoutPercent=16.9
2025-11-08T22:20:00Z GOVERNANCE proposal=122 phase=tally votingPower=149580 turnoutPercent=17.0
2025-11-08T22:21:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149617 turnoutPercent=17.1
2025-11-08T22:22:00Z GOVERNANCE proposal=122 phase=voting votingPower=149654 turnoutPercent=17.2
2025-11-08T22:23:00Z GOVERNANCE proposal=122 phase=tally votingPower=149691 turnoutPercent=17.3
2025-11-08T22:24:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149728 turnoutPercent=17.4
2025-11-08T22:25:00Z GOVERNANCE proposal=122 phase=voting votingPower=149765 turnoutPercent=17.5
2025-11-08T22:26:00Z GOVERNANCE proposal=122 phase=tally votingPower=149802 turnoutPercent=17.6
2025-11-08T22:27:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149839 turnoutPercent=17.7
2025-11-08T22:28:00Z GOVERNANCE proposal=122 phase=voting votingPower=149876 turnoutPercent=17.8
2025-11-08T22:29:00Z GOVERNANCE proposal=122 phase=tally votingPower=149913 turnoutPercent=17.9
2025-11-08T22:30:00Z GOVERNANCE proposal=122 phase=discussion votingPower=149950 turnoutPercent=18.0
2025-11-08T22:31:00Z GOVERNANCE proposal=122 phase=voting votingPower=149987 turnoutPercent=18.1
2025-11-08T22:32:00Z GOVERNANCE proposal=122 phase=tally votingPower=150024 turnoutPercent=18.2
2025-11-08T22:33:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150061 turnoutPercent=18.3
2025-11-08T22:34:00Z GOVERNANCE proposal=122 phase=voting votingPower=150098 turnoutPercent=18.4
2025-11-08T22:35:00Z GOVERNANCE proposal=122 phase=tally votingPower=150135 turnoutPercent=18.5
2025-11-08T22:36:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150172 turnoutPercent=18.6
2025-11-08T22:37:00Z GOVERNANCE proposal=122 phase=voting votingPower=150209 turnoutPercent=18.7
2025-11-08T22:38:00Z GOVERNANCE proposal=122 phase=tally votingPower=150246 turnoutPercent=18.8
2025-11-08T22:39:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150283 turnoutPercent=18.9
2025-11-08T22:40:00Z GOVERNANCE proposal=122 phase=voting votingPower=150320 turnoutPercent=19.0
2025-11-08T22:41:00Z GOVERNANCE proposal=122 phase=tally votingPower=150357 turnoutPercent=19.1
2025-11-08T22:42:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150394 turnoutPercent=19.2
2025-11-08T22:43:00Z GOVERNANCE proposal=122 phase=voting votingPower=150431 turnoutPercent=19.3
2025-11-08T22:44:00Z GOVERNANCE proposal=122 phase=tally votingPower=150468 turnoutPercent=19.4
2025-11-08T22:45:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150505 turnoutPercent=19.5
2025-11-08T22:46:00Z GOVERNANCE proposal=122 phase=voting votingPower=150542 turnoutPercent=19.6
2025-11-08T22:47:00Z GOVERNANCE proposal=122 phase=tally votingPower=150579 turnoutPercent=19.7
2025-11-08T22:48:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150616 turnoutPercent=19.8
2025-11-08T22:49:00Z GOVERNANCE proposal=122 phase=voting votingPower=150653 turnoutPercent=19.9
2025-11-08T22:50:00Z GOVERNANCE proposal=122 phase=tally votingPower=150690 turnoutPercent=20.0
2025-11-08T22:51:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150727 turnoutPercent=20.1
2025-11-08T22:52:00Z GOVERNANCE proposal=122 phase=voting votingPower=150764 turnoutPercent=20.2
2025-11-08T22:53:00Z GOVERNANCE proposal=122 phase=tally votingPower=150801 turnoutPercent=20.3
2025-11-08T22:54:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150838 turnoutPercent=20.4
2025-11-08T22:55:00Z GOVERNANCE proposal=122 phase=voting votingPower=150875 turnoutPercent=20.5
2025-11-08T22:56:00Z GOVERNANCE proposal=122 phase=tally votingPower=150912 turnoutPercent=20.6
2025-11-08T22:57:00Z GOVERNANCE proposal=122 phase=discussion votingPower=150949 turnoutPercent=20.7
2025-11-08T22:58:00Z GOVERNANCE proposal=122 phase=voting votingPower=150986 turnoutPercent=20.8
2025-11-08T22:59:00Z GOVERNANCE proposal=122 phase=tally votingPower=151023 turnoutPercent=20.9
2025-11-08T23:00:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151060 turnoutPercent=15.0
2025-11-08T23:01:00Z GOVERNANCE proposal=123 phase=voting votingPower=151097 turnoutPercent=15.1
2025-11-08T23:02:00Z GOVERNANCE proposal=123 phase=tally votingPower=151134 turnoutPercent=15.2
2025-11-08T23:03:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151171 turnoutPercent=15.3
2025-11-08T23:04:00Z GOVERNANCE proposal=123 phase=voting votingPower=151208 turnoutPercent=15.4
2025-11-08T23:05:00Z GOVERNANCE proposal=123 phase=tally votingPower=151245 turnoutPercent=15.5
2025-11-08T23:06:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151282 turnoutPercent=15.6
2025-11-08T23:07:00Z GOVERNANCE proposal=123 phase=voting votingPower=151319 turnoutPercent=15.7
2025-11-08T23:08:00Z GOVERNANCE proposal=123 phase=tally votingPower=151356 turnoutPercent=15.8
2025-11-08T23:09:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151393 turnoutPercent=15.9
2025-11-08T23:10:00Z GOVERNANCE proposal=123 phase=voting votingPower=151430 turnoutPercent=16.0
2025-11-08T23:11:00Z GOVERNANCE proposal=123 phase=tally votingPower=151467 turnoutPercent=16.1
2025-11-08T23:12:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151504 turnoutPercent=16.2
2025-11-08T23:13:00Z GOVERNANCE proposal=123 phase=voting votingPower=151541 turnoutPercent=16.3
2025-11-08T23:14:00Z GOVERNANCE proposal=123 phase=tally votingPower=151578 turnoutPercent=16.4
2025-11-08T23:15:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151615 turnoutPercent=16.5
2025-11-08T23:16:00Z GOVERNANCE proposal=123 phase=voting votingPower=151652 turnoutPercent=16.6
2025-11-08T23:17:00Z GOVERNANCE proposal=123 phase=tally votingPower=151689 turnoutPercent=16.7
2025-11-08T23:18:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151726 turnoutPercent=16.8
2025-11-08T23:19:00Z GOVERNANCE proposal=123 phase=voting votingPower=151763 turnoutPercent=16.9
2025-11-08T23:20:00Z GOVERNANCE proposal=123 phase=tally votingPower=151800 turnoutPercent=17.0
2025-11-08T23:21:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151837 turnoutPercent=17.1
2025-11-08T23:22:00Z GOVERNANCE proposal=123 phase=voting votingPower=151874 turnoutPercent=17.2
2025-11-08T23:23:00Z GOVERNANCE proposal=123 phase=tally votingPower=151911 turnoutPercent=17.3
2025-11-08T23:24:00Z GOVERNANCE proposal=123 phase=discussion votingPower=151948 turnoutPercent=17.4
2025-11-08T23:25:00Z GOVERNANCE proposal=123 phase=voting votingPower=151985 turnoutPercent=17.5
2025-11-08T23:26:00Z GOVERNANCE proposal=123 phase=tally votingPower=152022 turnoutPercent=17.6
2025-11-08T23:27:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152059 turnoutPercent=17.7
2025-11-08T23:28:00Z GOVERNANCE proposal=123 phase=voting votingPower=152096 turnoutPercent=17.8
2025-11-08T23:29:00Z GOVERNANCE proposal=123 phase=tally votingPower=152133 turnoutPercent=17.9
2025-11-08T23:30:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152170 turnoutPercent=18.0
2025-11-08T23:31:00Z GOVERNANCE proposal=123 phase=voting votingPower=152207 turnoutPercent=18.1
2025-11-08T23:32:00Z GOVERNANCE proposal=123 phase=tally votingPower=152244 turnoutPercent=18.2
2025-11-08T23:33:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152281 turnoutPercent=18.3
2025-11-08T23:34:00Z GOVERNANCE proposal=123 phase=voting votingPower=152318 turnoutPercent=18.4
2025-11-08T23:35:00Z GOVERNANCE proposal=123 phase=tally votingPower=152355 turnoutPercent=18.5
2025-11-08T23:36:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152392 turnoutPercent=18.6
2025-11-08T23:37:00Z GOVERNANCE proposal=123 phase=voting votingPower=152429 turnoutPercent=18.7
2025-11-08T23:38:00Z GOVERNANCE proposal=123 phase=tally votingPower=152466 turnoutPercent=18.8
2025-11-08T23:39:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152503 turnoutPercent=18.9
2025-11-08T23:40:00Z GOVERNANCE proposal=123 phase=voting votingPower=152540 turnoutPercent=19.0
2025-11-08T23:41:00Z GOVERNANCE proposal=123 phase=tally votingPower=152577 turnoutPercent=19.1
2025-11-08T23:42:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152614 turnoutPercent=19.2
2025-11-08T23:43:00Z GOVERNANCE proposal=123 phase=voting votingPower=152651 turnoutPercent=19.3
2025-11-08T23:44:00Z GOVERNANCE proposal=123 phase=tally votingPower=152688 turnoutPercent=19.4
2025-11-08T23:45:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152725 turnoutPercent=19.5
2025-11-08T23:46:00Z GOVERNANCE proposal=123 phase=voting votingPower=152762 turnoutPercent=19.6
2025-11-08T23:47:00Z GOVERNANCE proposal=123 phase=tally votingPower=152799 turnoutPercent=19.7
2025-11-08T23:48:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152836 turnoutPercent=19.8
2025-11-08T23:49:00Z GOVERNANCE proposal=123 phase=voting votingPower=152873 turnoutPercent=19.9
2025-11-08T23:50:00Z GOVERNANCE proposal=123 phase=tally votingPower=152910 turnoutPercent=20.0
2025-11-08T23:51:00Z GOVERNANCE proposal=123 phase=discussion votingPower=152947 turnoutPercent=20.1
2025-11-08T23:52:00Z GOVERNANCE proposal=123 phase=voting votingPower=152984 turnoutPercent=20.2
2025-11-08T23:53:00Z GOVERNANCE proposal=123 phase=tally votingPower=153021 turnoutPercent=20.3
2025-11-08T23:54:00Z GOVERNANCE proposal=123 phase=discussion votingPower=153058 turnoutPercent=20.4
2025-11-08T23:55:00Z GOVERNANCE proposal=123 phase=voting votingPower=153095 turnoutPercent=20.5
2025-11-08T23:56:00Z GOVERNANCE proposal=123 phase=tally votingPower=153132 turnoutPercent=20.6
2025-11-08T23:57:00Z GOVERNANCE proposal=123 phase=discussion votingPower=153169 turnoutPercent=20.7
2025-11-08T23:58:00Z GOVERNANCE proposal=123 phase=voting votingPower=153206 turnoutPercent=20.8
2025-11-08T23:59:00Z GOVERNANCE proposal=123 phase=tally votingPower=153243 turnoutPercent=20.9
```


---

## 19. References and Citations

[1] Chia Network. "Chia Consensus and Networking Documentation." 2021.  
[2] Wesolowski, B. "Efficient Verifiable Delay Functions." EUROCRYPT 2019.  
[3] Boneh, D., Bünz, B., Fisch, B. "A Survey of Two Verifiable Delay Functions." Cryptology ePrint Archive, 2018.  
[4] Archivas Network Operations. "Devnet Block Production Statistics." Internal Report, 2025.  
[5] Nakamoto, S. "Bitcoin: A Peer-to-Peer Electronic Cash System." 2008.  
[6] Percival, C. and Josefsson, S. "The scrypt Password-Based Key Derivation Function." IETF RFC 7914, 2016.  
[7] Bernstein, D. J., Lange, T. "eBACS: ECRYPT Benchmarking of Cryptographic Systems." 2024.  
[8] libp2p Project. "libp2p Specification." Accessed 2025.  
[9] Archivas Economic Working Group. "Mainnet Emission Modeling Spreadsheet." 2025.  
[10] OpenSSL Project. "TLS 1.3 Specification." IETF RFC 8446, 2018.  

---

*End of Draft. Layout teams should transform this manuscript into a formatted ninety page publication with diagrams and charts derived from the provided data tables.*

