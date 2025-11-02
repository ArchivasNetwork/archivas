# Rollback to v1.1.0 (Stable Release)

**Date:** 2025-11-02  
**Reason:** v1.1.1 introduced deployment complexity (IBD, build cache issues, Server B sync failures)  
**Action:** Revert to last stable release (v1.1.0 - Wallet API Freeze)

---

## What v1.1.0 Includes

✅ **Wallet Primitives:**
- Ed25519 keypairs (BIP39, SLIP-0010)
- Bech32 addresses (`arcv` prefix)
- Transaction v1 schema (canonical JSON, domain separation)
- Full signing and verification

✅ **RPC Endpoints:**
- `GET /account/<addr>` - Balance and nonce (as strings)
- `GET /chainTip` - Current tip (height, hash, difficulty as strings)
- `GET /mempool` - Pending transactions
- `GET /tx/<hash>` - Transaction details
- `GET /estimateFee?bytes=<n>` - Fee estimation
- `POST /submit` - Submit signed transaction
- `POST /submitTx` - Legacy endpoint (backward compatible)

✅ **CLI Tools:**
- `archivas-cli` - keygen, addr, sign-transfer, broadcast
- `archivas-wallet` - Legacy wallet tool

✅ **Tested and Stable:**
- Transaction submission working
- Explorer compatible
- Metrics exposed (`archivas_blocks_total`, `archivas_tip_height`, etc.)
- Multi-server deployment tested

---

## Deployment Commands

### All Servers

```bash
cd ~/archivas
git fetch --tags
git checkout v1.1.0
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-wallet ./cmd/archivas-wallet
```

### Server A (57.129.148.132): Node + Timelord + Farmer

```bash
pkill -f archivas-node
pkill -f archivas-timelord
pkill -f archivas-farmer

# Start node (keep existing data):
nohup ./archivas-node \
  --rpc 0.0.0.0:8080 \
  --p2p :9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes "72.251.11.191:9090" \
  > logs/node.log 2>&1 &

# Start timelord:
nohup ./archivas-timelord \
  --node http://localhost:8080 \
  --step 500 \
  > logs/timelord.log 2>&1 &

# Start farmer:
nohup ./archivas-farmer farm \
  --plots ./plots-massive \
  --node http://localhost:8080 \
  --farmer-privkey <YOUR_KEY> \
  > logs/farmer.log 2>&1 &
```

### Server B (72.251.11.191): Node + Timelord

```bash
pkill -f archivas-node
pkill -f archivas-timelord

# Start node:
nohup ./archivas-node \
  --rpc 0.0.0.0:8080 \
  --p2p :9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes "57.129.148.132:9090" \
  > logs/node.log 2>&1 &

# Start timelord:
nohup ./archivas-timelord \
  --node http://localhost:8080 \
  --step 500 \
  > logs/timelord.log 2>&1 &
```

### Server C (57.129.148.134): Farmer

```bash
pkill -f archivas-farmer

nohup ./archivas-farmer farm \
  --plots /home/ubuntu/archivas-plots-all \
  --node http://57.129.148.132:8080 \
  --farmer-privkey <YOUR_KEY> \
  > logs/farmer.log 2>&1 &
```

---

## Verification

```bash
# Check version:
curl -s http://57.129.148.132:8080/version | jq

# Check chain tip:
curl -s http://57.129.148.132:8080/chainTip | jq

# Test wallet transfer:
./archivas-wallet send \
  --from-privkey <KEY> \
  --to arcv1... \
  --amount 1000000000 \
  --fee 100000 \
  --node http://57.129.148.132:8080

# Check metrics:
curl -s http://57.129.148.132:8080/metrics | grep archivas_blocks_total
curl -s http://57.129.148.132:8080/metrics | grep archivas_tip_height
```

---

## Network Details

- **Network ID:** `archivas-devnet-v3` (not v4)
- **Genesis Hash:** (from existing chain)
- **No chain wipe** - keeps existing blocks and state
- **Farmers keep existing plots**

---

## What's NOT in v1.1.0

❌ IBD (Initial Block Download) - batched sync  
❌ Build stamping and self-tests  
❌ Advanced observability (watchdogs, etc.)  
❌ Devnet V4 genesis  

These will be re-introduced in a future stable release after proper testing.

---

## Rollback Complete When

- All 3 servers running v1.1.0
- Metrics updating normally
- Farmers mining blocks
- Wallet transactions work
- No sync errors in logs

