# Archivas Architecture

## System Overview

Archivas consists of four main components that work together to create a complete blockchain network:

```
┌─────────────────────────────────────────────┐
│           Archivas Network                  │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────────┐         ┌──────────────┐ │
│  │    Node      │◄───────►│   Farmer     │ │
│  │  (Validator) │         │  (PoSpace)   │ │
│  └──────┬───────┘         └──────────────┘ │
│         │                                   │
│         ▼                                   │
│  ┌──────────────┐         ┌──────────────┐ │
│  │  Timelord    │         │    Wallet    │ │
│  │    (VDF)     │         │  (Key Mgmt)  │ │
│  └──────────────┘         └──────────────┘ │
│                                             │
└─────────────────────────────────────────────┘
```

### 1. Node (archivas-node)

**Role:** Validator and chain coordinator

**Responsibilities:**
- Maintain blockchain state
- Validate Proof-of-Space proofs
- Validate VDF proofs
- Process transactions
- Manage mempool
- Persist blocks to database
- Gossip blocks to peers (P2P)
- Serve RPC API

**Key Components:**
- **Ledger:** Account balances, transaction processing
- **Consensus:** Difficulty adjustment, challenge generation
- **Storage:** BadgerDB persistence
- **RPC:** HTTP API for farmers/wallets
- **P2P:** TCP networking for block propagation

### 2. Timelord (archivas-timelord)

**Role:** VDF computer (provides temporal security)

**Responsibilities:**
- Continuously compute VDF from chain tip
- Derive VDF seed: `SHA256(blockHash || height)`
- Iterate VDF: `y_n = SHA256(y_{n-1})`
- Publish VDF updates to node via `/vdf/update`
- Reset on new blocks (new canonical tip)

**Properties:**
- Sequential computation (cannot be parallelized)
- Provides temporal ordering
- Prevents grinding attacks
- Enforces real elapsed time

### 3. Farmer (archivas-farmer)

**Role:** Block producer (finds Proof-of-Space)

**Responsibilities:**
- Create plots (precomputed hash tables)
- Poll node for current challenge
- Search plots for best quality proof
- Submit winning blocks to node
- Earn RCHV rewards

**Farming Process:**
1. Get challenge from `/challenge` endpoint
2. Scan all plots for best quality: `Q = H(challenge || plot_hash)`
3. If quality < difficulty → WIN!
4. Submit block with proof via `/submitBlock`
5. Receive 20 RCHV reward

### 4. Wallet (archivas-wallet)

**Role:** Key management and transactions

**Responsibilities:**
- Generate secp256k1 keypairs
- Derive bech32 addresses (arcv prefix)
- Sign transactions
- Query balances
- Send RCHV

**Commands:**
- `archivas-wallet new` - Generate new wallet
- `archivas-wallet send` - Send RCHV (coming soon)

---

## Data Flow

### Block Production Flow

```
1. Timelord computes VDF
   ↓
2. Node receives VDF update
   ↓
3. Node generates challenge = H(VDF_output || height)
   ↓
4. Farmer polls /challenge
   ↓
5. Farmer scans plots for best proof
   ↓
6. Farmer finds winning proof (quality < difficulty)
   ↓
7. Farmer submits block to node
   ↓
8. Node validates: PoSpace + VDF + transactions
   ↓
9. Node accepts block, pays farmer 20 RCHV
   ↓
10. Node gossips block to peers (P2P)
   ↓
11. Peers receive and validate block
   ↓
12. Peers sync block to their chain
   ↓
13. Timelord detects new tip, resets VDF
   ↓
[CYCLE REPEATS]
```

### Multi-Node Sync Flow

```
Node A (at height 50)         Node B (at height 0)
        │                              │
        │◄────── P2P Connect ──────────┤
        │                              │
        │─────── STATUS (h=50) ───────►│
        │                              │
        │◄────── GET_BLOCK(1) ──────────┤
        │                              │
        │─────── BLOCK_DATA(1) ────────►│
        │                              │
        │                         Verify & Apply
        │                              │
        │◄────── GET_BLOCK(2) ──────────┤
        │                              │
        │─────── BLOCK_DATA(2) ────────►│
        │                              │
        │                         [Continues 3→50]
        │                              │
        │                         ✅ Synced!
```

---

## Package Structure

### Core Packages

**config/** - Chain parameters and genesis
- `params_devnet.go` - Chain constants
- `genesis.go` - Genesis loading and hashing

**ledger/** - State and transactions
- `state.go` - Account balances and nonces
- `tx.go` - Transaction structure
- `apply_tx.go` - State transitions
- `verify.go` - Signature verification

**pospace/** - Proof-of-Space
- `pospace.go` - Plot generation and verification
- `CheckChallenge()` - Search plots for proof
- `VerifyProof()` - Validate PoSpace proof

**vdf/** - Verifiable Delay Functions
- `vdf.go` - Iterated SHA-256 implementation
- `ComputeSequential()` - Advance VDF
- `VerifySequential()` - Validate VDF output

**consensus/** - Difficulty and challenges
- `difficulty.go` - Adaptive difficulty algorithm
- `challenge.go` - Challenge generation
- `consensus.go` - PoSpace verification wrapper

**storage/** - Persistence layer
- `storage.go` - BadgerDB wrapper
- `blockchain.go` - Block and state persistence

**rpc/** - HTTP API
- `rpc.go` - Core endpoints (balance, submitTx)
- `farming.go` - Farming endpoints (challenge, submitBlock)
- `vdf_server.go` - VDF endpoints (chainTip, vdf/update)

**p2p/** - Networking
- `p2p.go` - TCP P2P implementation
- `protocol.go` - Message types
- `sync.go` - Block synchronization logic

**wallet/** - Cryptography
- `wallet.go` - Key generation, address derivation
- `tx_sign.go` - Transaction signing

### Binary Commands

**cmd/archivas-node/** - Full validator
- Loads genesis or restores from disk
- Validates PoSpace + VDF proofs
- Processes transactions
- Gossips blocks to peers
- Serves RPC API

**cmd/archivas-farmer/** - Block producer
- `plot` subcommand: Generate plots
- `farm` subcommand: Search for proofs and submit blocks

**cmd/archivas-timelord/** - VDF computer
- Polls node for chain tip
- Computes VDF continuously
- Posts updates to node

**cmd/archivas-wallet/** - Key management
- `new` subcommand: Generate wallet
- `send` subcommand: Send RCHV (coming soon)

---

## Database Schema

### BadgerDB Keys

**Blocks:**
- `blk:<height>` → Block JSON

**Accounts:**
- `acc:<address>` → {balance, nonce}

**Metadata:**
- `meta:tip_height` → Current chain tip
- `meta:difficulty` → Current difficulty
- `meta:genesis_hash` → Genesis hash
- `meta:network_id` → Network identifier
- `meta:vdf_seed` → VDF seed
- `meta:vdf_iterations` → VDF iteration count
- `meta:vdf_output` → VDF output

---

## Network Protocol

### RPC Endpoints

**Wallet/Query:**
- `GET /` - Server status
- `GET /balance/<address>` - Query balance
- `POST /submitTx` - Submit transaction

**Farming:**
- `GET /challenge` - Current challenge + VDF info
- `POST /submitBlock` - Submit winning block

**VDF/Timelord:**
- `GET /chainTip` - Current tip (height, hash, difficulty)
- `POST /vdf/update` - VDF update from timelord

**Network:**
- `GET /genesisHash` - Genesis hash for verification

### P2P Messages

**Handshake:**
- `STATUS` - Exchange height, difficulty, tipHash

**Block Propagation:**
- `NEW_BLOCK` - Announce new block (height, hash)
- `GET_BLOCK` - Request block by height
- `BLOCK_DATA` - Send full block data

**Keepalive:**
- `PING` / `PONG` - Connection health

**Protocol:** Newline-delimited JSON over TCP

---

## Security Model

### Cryptographic Security
- **secp256k1:** Industry-standard elliptic curve
- **ECDSA:** Digital signatures for transactions
- **SHA-256:** Hashing for challenges and VDF
- **Bech32:** Error-detecting address format

### Consensus Security
- **Proof-of-Space:** Cannot fake disk space (must precompute plots)
- **VDF:** Cannot skip iterations (sequential time required)
- **Challenge:** Derived from VDF (unpredictable until computed)
- **Difficulty:** Adaptive (maintains ~20s block times)

### Attack Resistance
- **Grinding:** Prevented by VDF (cannot recompute alternative timelines)
- **Precomputation:** VDF seed changes per block
- **Nothing-at-Stake:** Disk space is committed resource
- **Sybil:** P2P uses genesis hash validation

---

**Next:** [Consensus Details →](architecture/consensus.md)  
**Back:** [← Introduction](introduction/proof-of-space-time.md)

