# Archivas

> **A Proof-of-Space-and-Time L1 Blockchain. Farm RCHV with disk space.**

---

## 🌍 Public RPC

**Endpoint:** `https://seed.archivas.ai`

Archivas provides a public RPC endpoint for developers to interact with the testnet.

### Available Routes

**Account & Balance:**
- `GET /account/<address>` - Get account balance and nonce

**Chain Info:**
- `GET /chainTip` - Get current blockchain height, hash, and difficulty
- `GET /genesisHash` - Get genesis block hash

**Blocks (v1.2.0):**
- `GET /blocks/recent?limit=N` - List recent blocks (default 20, max 100)
- `GET /block/<height>` - Get block details by height

**Transactions:**
- `GET /tx/<hash>` - Get transaction details
- `GET /tx/recent?limit=N` - List recent transactions (default 50, max 200)
- `GET /mempool` - List pending transactions
- `POST /submit` - Submit a signed transaction (JSON, CORS-enabled)

**Utilities:**
- `GET /estimateFee?bytes=<n>` - Estimate transaction fee

### Examples

```bash
# Get chain tip
curl https://seed.archivas.ai/chainTip
# {"height":"13080","hash":"87e0cf03151c2debd7dab3d7143eb24f8f0281826e13ab49aaf3259c303a2810","difficulty":"1000000"}

# Get account balance
curl https://seed.archivas.ai/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
# {"address":"arcv1...","balance":"137540000000000","nonce":"0"}

# Attempt GET on /submit (returns 405 Method Not Allowed)
curl -i https://seed.archivas.ai/submit
# HTTP/2 405
# allow: POST

# List recent blocks (v1.2.0)
curl https://seed.archivas.ai/blocks/recent?limit=10
# {"blocks": [{"height":"42808","hash":"...","farmer":"arcv1...","txCount":"1",...}]}

# Get specific block
curl https://seed.archivas.ai/block/42808
# {"height":"42808","hash":"...","prevHash":"...","difficulty":"1000000",...}

# List recent transactions (v1.2.0)
curl https://seed.archivas.ai/tx/recent?limit=20
# {"txs": [{"hash":"...","from":"arcv1...","to":"arcv1...","amount":"30000000000",...}]}
```

**Note:** 
- The `/submit` endpoint only accepts POST requests with `Content-Type: application/json`.
- The `/blocks/recent` and `/tx/recent` endpoints show recent activity only (not full historical archive).

For setup and deployment details, see [docs/SEED_HOST.md](docs/SEED_HOST.md).

---

[![Build Status](https://github.com/ArchivasNetwork/archivas/workflows/Build%20and%20Test/badge.svg)](https://github.com/ArchivasNetwork/archivas/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)

**Archivas** is a Layer 1 blockchain secured by disk space and sequential time, not energy-intensive computation or capital lockup. Anyone with storage can participate in consensus and earn RCHV.

⚠️ **EXPERIMENTAL TESTNET** - This is research software. Not audited. Not for production use. Use at your own risk.

---

## Quick Start

### Build

```bash
# Clone repository
git clone https://github.com/ArchivasNetwork/archivas
cd archivas

# Download dependencies
go mod download

# Build all binaries
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-wallet ./cmd/archivas-wallet
```

### Farm Your First RCHV

```bash
# 1. Generate a wallet
./archivas-wallet new

# 2. Create a plot
./archivas-farmer plot --size 18 --path ./plots

# 3. Start the node
./archivas-node

# 4. Start farming
./archivas-farmer farm \
  --plots ./plots \
  --farmer-key <your_private_key_from_step_1>

# 5. Check your balance
curl http://localhost:8080/balance/<your_address_from_step_1>
```

Every block you farm = **20 RCHV**! 🌾

---

## What is Archivas?

Archivas uses **Proof-of-Space-and-Time** consensus:

- 🌾 **Proof-of-Space** - Farmers allocate disk space, create "plots"
- ⏰ **Verifiable Delay Functions** - Timelords compute sequential time proofs

Blocks require BOTH a winning disk space proof AND a valid time proof.

**Why this matters:**
- ✅ **Permissionless** - Anyone with disk can farm (no ASICs, no capital requirements)
- ✅ **Energy Efficient** - No wasteful computation (disk only)
- ✅ **Fair** - Disk space determines odds, not wealth
- ✅ **Secure** - VDF prevents grinding, PoSpace prevents centralization

**Same consensus class as Chia Network.**

---

## Features

- 🌾 **Proof-of-Space Farming** - Tested with real blocks
- ⏰ **VDF/Timelord** - Temporal security (ready to activate)
- 🔐 **Cryptographic Wallets** - secp256k1, bech32 addresses
- ✍️ **Transaction Signing** - ECDSA signatures
- 💾 **Persistent Storage** - BadgerDB, crash recovery
- 📊 **Adaptive Difficulty** - Maintains ~20s block times
- 🌐 **P2P Networking** - Multi-node capable (ready to activate)

---

## Status

🟢 **Devnet Operational**

| Feature | Status | Tested |
|---------|--------|--------|
| Wallet generation | ✅ Working | ✅ Yes |
| Transaction signing | ✅ Working | ✅ Yes |
| Plot generation | ✅ Working | ✅ Yes |
| PoSpace farming | ✅ Working | ✅ Yes (120 RCHV earned!) |
| Block rewards | ✅ Working | ✅ Yes |
| Persistent storage | ✅ Working | ✅ Yes (restart verified!) |
| VDF/Timelord | ✅ Implemented | ⏸️ Ready to activate |
| P2P networking | ✅ Implemented | ⏸️ Ready to activate |

---

## Chain Parameters

| Parameter | Value |
|-----------|-------|
| Chain ID | 1616 |
| Native Token | RCHV |
| Decimals | 8 |
| Block Time | ~20 seconds |
| Block Reward | 20.00000000 RCHV |
| Address Prefix | arcv (bech32) |
| Consensus | Proof-of-Space (+Time ready) |

---

## How It Works

### Proof-of-Space

1. **Create Plots** - Precompute large hash tables on disk
2. **Get Challenge** - Network broadcasts challenge for each new block
3. **Search Plots** - Find the best proof (lowest quality hash)
4. **Win Block** - If your proof beats difficulty, you produce the block
5. **Earn Reward** - Receive 20 RCHV for the block

More disk space = more lottery tickets = higher win probability.

### Verifiable Delay Functions (VDF Mode)

1. **Timelord** computes sequential function (can't be parallelized)
2. **VDF Output** derives the PoSpace challenge
3. **Cannot Skip** - Must compute all iterations sequentially
4. **Prevents Grinding** - Can't precompute alternative timelines

Blocks require BOTH PoSpace AND VDF = **Chia-class security**.

---

## Documentation

- **[START-HERE.md](START-HERE.md)** - Navigation guide
- **[STATUS.md](STATUS.md)** - Current technical status  
- **[JOURNEY.md](JOURNEY.md)** - Complete development story
- **[MILESTONE3.md](MILESTONE3.md)** - Farming guide (tested!)
- **[MILESTONE5-PERSISTENCE.md](MILESTONE5-PERSISTENCE.md)** - Storage guide
- **[MILESTONE6-P2P.md](MILESTONE6-P2P.md)** - Networking guide
- **[ACTIVATE-VDF.md](ACTIVATE-VDF.md)** - VDF activation instructions
- **[docs/WHITEPAPER-OUTLINE.md](docs/WHITEPAPER-OUTLINE.md)** - Technical whitepaper

---

## Architecture

```
┌─────────────────┐
│ archivas-node   │ ← Validates PoSpace+VDF, manages chain
└────────┬────────┘
         │
    ┌────┴─────────────────┐
    │                      │
┌───▼──────────┐  ┌───────▼────────┐
│  timelord    │  │     farmer     │
│  (VDF)       │  │   (PoSpace)    │
└──────────────┘  └────────────────┘
```

---

## Test Results

**Farming Test (60 seconds):**
- Blocks found: 6
- RCHV earned: 120.00000000 RCHV
- Difficulty: Adapted 5 times
- Plot size: k=16 (2MB)

**Persistence Test:**
- Blocks before restart: 7
- Node killed & restarted
- Blocks after restart: 7 ✅
- Balances: 100% preserved ✅
- Recovery time: <100ms

**All core features verified end-to-end.** ✅

---

## Roadmap

- [x] **Phase 1: Devnet** - Core blockchain ✅ COMPLETE
- [ ] **Phase 2: Testnet** - Multi-node P2P network ⏸️ READY
- [ ] **Phase 3: Public Testnet** - Community participation 🚧 Q1 2026
- [ ] **Phase 4: Mainnet** - Security audit, public launch 📋 Q2-Q3 2026

---

## Security Disclaimer

⚠️ **EXPERIMENTAL SOFTWARE - USE AT YOUR OWN RISK**

- This is research/testnet software
- NOT security audited
- NOT for production use
- NOT financial advice
- May contain bugs
- Private keys are YOUR responsibility
- RCHV has NO monetary value (testnet only)

**Do not:**
- Use on mainnet (doesn't exist yet)
- Store real value
- Treat as financial instrument
- Use without understanding risks

**Do:**
- Test, experiment, learn
- Report bugs
- Contribute improvements
- Have fun farming! 🌾

---

## Contributing

Archivas is open source! Contributions welcome.

**Focus areas:**
- P2P networking activation
- Performance optimization
- Testing and security
- Documentation
- Block explorer

See [MILESTONE6-P2P.md](MILESTONE6-P2P.md) for P2P integration guide.

---

## Community

- **GitHub:** [github.com/ArchivasNetwork/archivas](https://github.com/ArchivasNetwork/archivas)
- **Discussions:** [GitHub Discussions](https://github.com/ArchivasNetwork/archivas/discussions)
- **Issues:** [Report bugs](https://github.com/ArchivasNetwork/archivas/issues)

---

## License

MIT License - see [LICENSE](LICENSE) file

---

## Acknowledgments

Inspired by:
- **Chia Network** - PoSpace+Time consensus model
- **Filecoin** - Storage-based consensus
- **Bitcoin** - Cryptographic security

Built with Go, BadgerDB, secp256k1, and bech32.

---

<p align="center">
  <strong>Archivas: Farming the future of decentralized storage</strong> 🌾
</p>

## Decommissioning Nodes

### When to Decommission

Decommission an Archivas node when:
- **Non-syncing node** - Stuck at an old height and unable to catch up (e.g., Server B at height 3 with 5800+ block gap)
- **Redundant node** - Network has sufficient nodes and this one is no longer needed
- **Hardware reallocation** - Server needed for other purposes
- **Network upgrade** - Node cannot upgrade to new protocol version

### Decommission Process

**⚠️ Warning:** This will permanently delete all blockchain data, plots, and configuration. Plots can be recreated, but sync history will be lost.

**Command:**
```bash
cd ~/archivas
bash scripts/decommission-node.sh
```

**What it does:**
1. Stops all Archivas systemd services (node, timelord, farmer)
2. Disables services to prevent auto-start
3. Kills any stray processes
4. Removes data directories (`~/.archivas`, `/var/lib/archivas`, `/opt/archivas`, etc.)
5. Removes systemd unit files
6. Runs `systemctl daemon-reload`
7. Verifies cleanup completed successfully

**Safe to rerun** - The script is idempotent and can be run multiple times.

### After Decommissioning

The server will be clean and ready for:
- Redeployment as a fresh node (sync from genesis)
- Other applications
- Deallocation/shutdown

To redeploy later, clone the repository and follow the setup guide in `docs/`.

---
