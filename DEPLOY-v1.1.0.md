# Deployment Guide: v1.1.0 "Wallet Primitives + Public API Freeze"

## What's New

- ✅ Frozen wallet API (ed25519, BIP39, SLIP-0010, Bech32)
- ✅ New RPC endpoints: `/tx/<hash>`, `/estimateFee`, `/submit`
- ✅ Updated endpoints: `/account/<addr>`, `/chainTip`, `/mempool` (return strings for amounts)
- ✅ New CLI tool: `archivas-cli` (keygen, sign-transfer, broadcast)
- ✅ Backward compatible (existing endpoints unchanged)

## Deployment Steps

### 1. Node Server (Main RPC Server)

**On your node server** (e.g., `ns5042239` / `72.251.11.191`):

```bash
# Navigate to archivas directory
cd ~/archivas

# Pull latest code
git pull origin main  # or your branch name

# Build new node binary
go build -o archivas-node ./cmd/archivas-node

# Build new CLI tool (optional, for testing)
go build -o archivas-cli ./cmd/archivas-cli

# Stop old node process
pkill -f archivas-node

# Wait a moment for process to stop
sleep 2

# Verify it's stopped
ps aux | grep archivas-node | grep -v grep

# Restart node with your original command (preserve flags, genesis, etc.)
# Example (adjust to your actual command):
nohup ./archivas-node \
  --genesis genesis/devnet.genesis.json \
  --network-id devnet \
  --bootnodes "57.129.148.132:9090" \
  > logs/node.log 2>&1 &

# Verify it started
sleep 2
ps aux | grep archivas-node | grep -v grep
tail -20 logs/node.log
```

### 2. Verify New Endpoints

```bash
# Test updated /account endpoint (should return strings now)
curl -s http://localhost:8080/account/arcv1t3huuyd08er3yfnmk9c935rmx3wdh5j6m2uc9d | jq

# Expected format:
# {
#   "address": "arcv1...",
#   "balance": "129927999900000",  # ← Note: string, not number
#   "nonce": "1"                   # ← Note: string, not number
# }

# Test updated /chainTip endpoint
curl -s http://localhost:8080/chainTip | jq

# Expected format:
# {
#   "height": "1234",      # ← string
#   "hash": "abcd...",     # hex
#   "difficulty": "15000000"  # ← string
# }

# Test new /mempool endpoint (should return array of hashes)
curl -s http://localhost:8080/mempool | jq

# Test new /estimateFee endpoint
curl -s "http://localhost:8080/estimateFee?bytes=512" | jq

# Expected:
# {
#   "fee": "100"
# }

# Test new /tx/<hash> endpoint (will return not found for now, but shouldn't error)
curl -s http://localhost:8080/tx/abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234 | jq
```

### 3. Test CLI Tool (Optional)

```bash
# Generate a new wallet
./archivas-cli keygen

# Example output:
# Mnemonic: word1 word2 ... word24
# Address:  arcv1...
# PubKey:   abcdef...
# PrivKey:  ...

# Derive address from mnemonic
./archivas-cli addr "word1 word2 ... word24"

# Sign a transfer (example)
./archivas-cli sign-transfer \
  --from-mnemonic "word1 ... word24" \
  --to arcv1hr2vm4v4xsehsdl3a3flxspg3wguhtxymrvgrw \
  --amount 1000000000 \
  --fee 100 \
  --nonce 0 \
  --out tx.json

# View signed transaction
cat tx.json | jq

# Broadcast transaction
./archivas-cli broadcast tx.json http://localhost:8080
```

### 4. Check Node Logs

```bash
# Watch for any errors
tail -f logs/node.log | grep -i error

# Check for new endpoint registrations
grep -i "wallet\|submit\|tx" logs/node.log
```

### 5. Farmer Servers (No Changes Required)

**Farmer servers don't need updates** - they interact with the node via RPC, and we've maintained backward compatibility for existing endpoints (`/submitBlock`, `/challenge`, etc.).

However, if you want to test the new wallet primitives on farmer servers:

```bash
# On farmer server
cd ~/archivas
git pull origin main
go build -o archivas-cli ./cmd/archivas-cli

# Test keygen
./archivas-cli keygen
```

## Rollback (If Needed)

If you encounter issues:

```bash
# Stop new node
pkill -f archivas-node

# Checkout previous version
git checkout <previous-commit-hash>  # or v1.0.2 tag

# Rebuild
go build -o archivas-node ./cmd/archivas-node

# Restart
nohup ./archivas-node [your-flags] > logs/node.log 2>&1 &
```

## Verification Checklist

- [ ] Node compiles without errors
- [ ] Node starts successfully
- [ ] `/account/<addr>` returns balance/nonce as strings
- [ ] `/chainTip` returns height/hash/difficulty as strings
- [ ] `/mempool` returns array (may be empty)
- [ ] `/estimateFee?bytes=256` returns fee estimate
- [ ] `/tx/<hash>` endpoint exists (may return not found)
- [ ] `/submit` endpoint exists (test with CLI)
- [ ] Existing farming still works (`/submitBlock`, `/challenge`)
- [ ] No errors in node logs

## Troubleshooting

**Issue: "cannot find package pkg/tx/v1"**
```bash
# Make sure you pulled the latest code
git pull origin main
go mod tidy
go build -o archivas-node ./cmd/archivas-node
```

**Issue: "/submit endpoint returns error"**
- Check that the signed transaction JSON is valid
- Verify signature using the CLI: `./archivas-cli broadcast tx.json http://localhost:8080`
- Check node logs for detailed error messages

**Issue: "Node won't start after update"**
- Check database lock: `pkill -9 archivas-node` (if process is stuck)
- Verify genesis file path: `ls -la genesis/devnet.genesis.json`
- Check port conflicts: `netstat -tlnp | grep 8080`

## Next Steps

After successful deployment:

1. **Test transaction signing** with the CLI
2. **Verify backward compatibility** - existing farmers should continue working
3. **Document any issues** - report if any existing endpoints break
4. **Monitor metrics** - check `/metrics` endpoint for any anomalies

---

**Note**: This is a **non-breaking update**. All existing functionality is preserved, and new endpoints are additive only.

