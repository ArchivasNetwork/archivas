# Getting Started with Archivas

## Quick Start Guide

This guide will help you join the Archivas testnet and start farming RCHV.

### Prerequisites

**System Requirements:**
- Ubuntu 20.04+ (or similar Linux)
- Go 1.21 or higher
- 2GB RAM minimum
- 100MB+ disk space for plots
- Public IP (optional, for running a node)

**For Farming:**
- Disk space for plots (more = better odds)
- k=16: ~2 MB (testing)
- k=20: ~32 MB (small farm)
- k=24: ~512 MB (medium farm)

### Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
```

#### 2. Build Binaries

```bash
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-wallet ./cmd/archivas-wallet
```

**Verify:**
```bash
./archivas-wallet --help
./archivas-farmer --help
./archivas-node --help
```

### Join the Testnet

#### Step 1: Generate a Wallet

```bash
./archivas-wallet new
```

**Output:**
```
🔐 New Archivas Wallet Generated

Address:     arcv1q84xt5pzcslhnjsc2h2t9cnuxrn0e2u2u97jnl
Public Key:  03457989304d0c1ecbe33bcdb2b5ae8f8f34a4d2c0f278a7ad79460c557fe98dd9
Private Key: 2fe06d47987c25e6182735dbbead5220f4d94b2c249614363d9a6c7a76a53922

⚠️  KEEP YOUR PRIVATE KEY SECRET!
```

**⚠️ Important:** Save both keys securely!
- **Public Key:** Use for creating plots
- **Private Key:** Use for farming (never share!)

#### Step 2: Create a Plot

```bash
mkdir -p plots

./archivas-farmer plot \
  --size 18 \
  --path ./plots \
  --farmer-pubkey <YOUR_PUBLIC_KEY>
```

**Options:**
- `--size 16`: 2MB plot (~65K hashes) - testing
- `--size 18`: 8MB plot (~260K hashes) - small farm
- `--size 20`: 32MB plot (~1M hashes) - medium farm

**Wait for completion:**
```
✅ Plot generated successfully in 2.3s
📊 Plot size: ~8.00 MB
```

#### Step 3: Run a Node

```bash
mkdir -p data logs

nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 57.129.148.132:9090 \
  > logs/node.log 2>&1 &
```

**Verify node is running:**
```bash
curl http://localhost:8080/chainTip
# Should show current height and block hash
```

**Watch sync:**
```bash
tail -f logs/node.log | grep "Synced block"
```

#### Step 4: Start Farming

```bash
nohup ./archivas-farmer farm \
  --plots ./plots \
  --farmer-privkey <YOUR_PRIVATE_KEY> \
  --node http://localhost:8080 \
  > logs/farmer.log 2>&1 &
```

**Watch farming:**
```bash
tail -f logs/farmer.log
```

**You'll see:**
```
⚙️  Checking plots... best=XXXXX, need=<YYYYY
🎉 Found winning proof!
✅ Block submitted successfully for height N
```

#### Step 5: Check Your Balance

```bash
curl http://localhost:8080/balance/<YOUR_ADDRESS>
```

**Response:**
```json
{
  "address": "arcv1q84xt5pzcslhnjsc2h2t9cnuxrn0e2u2u97jnl",
  "balance": 4000000000,
  "nonce": 0
}
```

Balance is in base units (8 decimals): `4000000000 = 40.00000000 RCHV`

Every block you mine = **20 RCHV!** 🌾

---

## Optional: Run a Timelord

**Note:** Only needed if you want to help compute VDF for the network.

```bash
nohup ./archivas-timelord \
  --node http://localhost:8080 \
  > logs/timelord.log 2>&1 &
```

---

## Troubleshooting

### Node won't start

**Check logs:**
```bash
tail -50 logs/node.log
```

**Common issues:**
- Port 8080 or 9090 already in use
- Genesis file not found (make sure you're in archivas/ directory)
- Database corruption (try `rm -rf data` and restart)

### Farmer not finding blocks

**This is normal!** Farming is a lottery.

**Expected time to win:**
- k=16 plot: 5-30 minutes
- k=20 plot: 1-5 minutes  
- k=24 plot: seconds to minutes

**Verify farmer is working:**
```bash
tail -f logs/farmer.log
```

Should show:
```
⚙️  Checking plots... best=XXXXX
```

If quality values are changing, farmer is working correctly!

### Node not syncing

**Check P2P connection:**
```bash
grep "p2p" logs/node.log | tail -20
```

Should show:
```
[p2p] connected to peer 57.129.148.132:9090
[p2p] peer status: height=XX
```

**Verify bootnode is reachable:**
```bash
nc -zv 57.129.148.132 9090
# Should show: Connection successful
```

---

## Network Information

### Bootstrap Nodes

**Primary:**
- IP: 57.129.148.132
- Port: 9090
- Location: US

**Secondary:**
- IP: 72.251.11.191
- Port: 9090
- Location: US

### Genesis

**File:** `genesis/devnet.genesis.json` (included in repo)  
**Hash:** `11b6fedb68f1da0f312039cd6fae91f4dd861bea942651b0c33590013f5b8a55`  
**Network ID:** `archivas-devnet-v3`  
**Timestamp:** 1730246400 (fixed)

**Verify your genesis:**
```bash
curl http://localhost:8080/genesisHash
```

Should return: `{"genesisHash":"11b6fedb68f1da0f..."}`

---

## Next Steps

**After joining:**
1. Monitor your node: `tail -f logs/node.log`
2. Watch farming: `tail -f logs/farmer.log`  
3. Check balance regularly
4. Report bugs on GitHub
5. Join community discussions

**Advanced:**
- [Run a Timelord](Testnet-Guide.md#timelord)
- [Send RCHV](Developer-Docs.md#transactions)
- [Run multiple nodes](Architecture.md#multi-node)

---

**Ready to farm?** [Testnet Guide →](Testnet-Guide.md)  
**Want to understand the tech?** [Architecture Overview →](architecture/overview.md)

