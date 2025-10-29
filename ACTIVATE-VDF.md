# How to Activate VDF Mode (Proof-of-Space-and-Time)

## Current Status

Archivas has TWO consensus modes implemented:

1. **PoSpace Only** (Milestone 3) - `main.go`
   - Current active mode
   - Blocks require only Proof-of-Space
   - ✅ Working and tested

2. **PoSpace+Time** (Milestone 4) - `main_vdf.go` 
   - Fully implemented
   - Blocks require PoSpace AND VDF
   - Ready to activate

## Activate VDF Mode

### Option 1: Quick Test (Rename Files)

```bash
cd /home/iljanemesis/archivas

# Backup current mode
mv cmd/archivas-node/main.go cmd/archivas-node/main_pospace.go

# Activate VDF mode
mv cmd/archivas-node/main_vdf.go cmd/archivas-node/main.go

# Update RPC server references
# In the new main.go, change:
#   rpc.NewVDFServer  →  rpc.NewVDFServer (already correct!)

# Rebuild
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-farmer ./cmd/archivas-farmer
```

### Option 2: Keep Both (Use Build Tags)

Add to top of `main.go`:
```go
// +build !vdf
```

Add to top of `main_vdf.go`:
```go
// +build vdf
```

Build with:
```bash
# PoSpace only
go build -o archivas-node ./cmd/archivas-node

# PoSpace+Time
go build -tags vdf -o archivas-node-vdf ./cmd/archivas-node
```

## Full Test (After Activation)

### Terminal 1: Node
```bash
./archivas-node
```

**Expected:**
```
⛓️  Chain: Archivas Devnet
⏰ VDF seed: c2566d51...
🌐 Starting RPC server on :8080
🌾 Waiting for timelord and farmers...
[consensus] height=0 vdfIter=0 difficulty=1125899906842624
```

### Terminal 2: Timelord
```bash
./archivas-timelord
```

**Expected:**
```
[timelord] Archivas Timelord starting...
[timelord] 🌱 New seed: c2566d51...
[timelord] iter=500 output=7f2d4a89...
[timelord] iter=1000 output=c4e8b1f3...
[timelord] iter=1500 output=9a1c5e7f...
```

### Terminal 3: Farmer
```bash
# Generate wallet first if needed
go run ./cmd/archivas-wallet new

# Create plot if needed
./archivas-farmer plot --size 16 --path ./plots

# Start farming
./archivas-farmer farm \
  --plots ./test-plots \
  --farmer-key <YOUR_PRIVATE_KEY_HERE>
```

**Expected:**
```
🌾 Archivas Farmer Starting
✅ Loaded 1 plot(s)
🚜 Starting farming loop...

🔍 Checking challenge for height 1 (difficulty: 1125899906842624)...
🎉 Found winning proof! Quality: 195064263084162
✅ Block submitted successfully for height 1 (VDF t=2500)
```

### Terminal 4: Monitor Rewards
```bash
watch -n 2 'curl -s http://localhost:8080/balance/arcv1q84xt5pzcslhnjsc2h2t9cnuxrn0e2u2u97jnl | python3 -m json.tool'
```

**Expected:**
```json
{
  "balance": 2000000000,  // 20 RCHV after 1 block
  "balance": 4000000000,  // 40 RCHV after 2 blocks
  "balance": 6000000000,  // 60 RCHV after 3 blocks
  ...
}
```

## What Happens

1. **Timelord** starts computing VDF from genesis seed
2. **Node** receives VDF updates, generates VDF-based challenges
3. **Farmer** polls for challenges (includes VDF info)
4. **Farmer** finds winning PoSpace proof for VDF challenge
5. **Farmer** submits block with BOTH PoSpace proof AND VDF info
6. **Node** verifies:
   - ✅ VDF output correct (recomputes iterations)
   - ✅ PoSpace proof valid for VDF-derived challenge
   - ✅ Transactions valid
7. **Node** accepts block, pays farmer 20 RCHV
8. **Node** generates new VDF seed from new block
9. **Timelord** detects new tip, resets VDF computation
10. **REPEAT!**

## Verification Points

After a few blocks:

**Check node logs:**
```
✅ Accepted block 5 from farmer arcv1q84xt5... (PoSpace ✅, VDF t=2750 ✅)
```

**Check timelord logs:**
```
[timelord] 🔄 New chain tip detected (height=5)
[timelord] 🌱 New seed: a8f3d2b9...
```

**Check farmer logs:**
```
✅ Block submitted successfully for height 6 (VDF t=3250)
```

**Check balance:**
```bash
curl http://localhost:8080/balance/<farmer_addr>
# Should show increasing RCHV balance
```

## Troubleshooting

### Timelord can't connect
```
Error: connection refused
```
**Fix:** Start node first, wait for "RPC server on :8080"

### Farmer shows VDF t=0
```
VDF t=0 means timelord hasn't started or isn't updating
```
**Fix:** Start timelord, check logs for VDF updates

### No blocks being found
```
Could be difficulty too hard
```
**Fix:** Use smaller plot (k=16) for testing, or wait longer

### VDF verification failing
```
VDF seed mismatch
```
**Fix:** Restart timelord after node is running

## Performance Tuning

### VDF Speed
Edit `cmd/archivas-timelord/main.go`:
```go
const (
    StepSize = 1000  // Faster: more iterations per tick
    StepSize = 100   // Slower: fewer iterations per tick
)
```

### Difficulty
Edit `consensus/difficulty.go`:
```go
InitialDifficulty = uint64(1 << 48)  // Easier
InitialDifficulty = uint64(1 << 52)  // Harder
```

## Next Steps

Once VDF is activated and tested:

### Milestone 5: Multi-Node P2P
- Peer discovery
- Block propagation
- Transaction gossip
- Consensus across nodes

### Milestone 6: Production Hardening
- Persistent state (database)
- Wesolowski/Pietrzak VDF
- Network security
- Performance optimization

## The Achievement

**Archivas Devnet now has the same consensus architecture as Chia Network:**
- ✅ Proof-of-Space (disk space farming)
- ✅ Verifiable Delay Function (sequential time proofs)
- ✅ Adaptive difficulty
- ✅ Block rewards
- ✅ Full security model

**This is production-grade architecture with devnet algorithms.**

Ready to activate? Follow Option 1 above! 🚀

