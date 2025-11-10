# Testnet Participation Guide

## Archivas Testnet v0.1.0

**Network:** Archivas Devnet v3  
**Status:** 🟢 LIVE  
**Nodes:** 2+ (join anytime!)  
**Genesis:** Deterministic (11b6fedb68f1da...)

---

## Network Information

### Bootstrap Nodes

**Bootnode:**
```
Hostname: seed.archivas.ai
IP: 57.129.148.132
Port: 9090
RPC: https://seed.archivas.ai
```

### Genesis Verification

**Genesis Hash:** `11b6fedb68f1da0f312039cd6fae91f4dd861bea942651b0c33590013f5b8a55`  
**Network ID:** `archivas-devnet-v3`  
**Timestamp:** 1730246400 (fixed)

**Verify your node:**
```bash
curl http://localhost:8080/genesisHash
# Should return: {"genesisHash":"11b6fedb68f1da0f..."}
```

---

## Joining the Testnet

### Full Node + Farmer

**Complete setup to join and farm:**

```bash
# 1. Clone and build
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-wallet ./cmd/archivas-wallet

# 2. Generate wallet
./archivas-wallet new
# Save the output!

# 3. Create plot (using public key from wallet)
mkdir -p plots
./archivas-farmer plot \
  --size 18 \
  --path ./plots \
  --farmer-pubkey <YOUR_PUBLIC_KEY>

# 4. Start node (connects to testnet)
mkdir -p data logs
nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 57.129.148.132:9090 \
  > logs/node.log 2>&1 &

# 5. Wait for sync
sleep 10
curl http://localhost:8080/chainTip
# Check height matches network

# 6. Start farming (using private key from wallet)
nohup ./archivas-farmer farm \
  --plots ./plots \
  --farmer-privkey <YOUR_PRIVATE_KEY> \
  --node http://localhost:8080 \
  > logs/farmer.log 2>&1 &

# 7. Monitor
tail -f logs/farmer.log
```

### Sync-Only Node

**Just sync and validate (no farming):**

```bash
nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes seed.archivas.ai:9090 \
  > logs/node.log 2>&1 &

# Watch sync
tail -f logs/node.log | grep "Synced block"
```

### Timelord Node

**Help compute VDF for the network:**

```bash
# Start node first
nohup ./archivas-node ... > logs/node.log 2>&1 &

# Start timelord
nohup ./archivas-timelord \
  --node http://localhost:8080 \
  --step 500 \
  > logs/timelord.log 2>&1 &

# Monitor
tail -f logs/timelord.log
```

---

## Monitoring Your Node

### Check Sync Status

```bash
# Current height
curl http://localhost:8080/chainTip | jq .height

# Compare to network
curl https://seed.archivas.ai/chainTip | jq .height

# Should match!
```

### Check Balance

```bash
curl http://localhost:8080/balance/<YOUR_ADDRESS> | jq
```

**Response:**
```json
{
  "address": "arcv1...",
  "balance": 2000000000,  // 20.00000000 RCHV
  "nonce": 0
}
```

### Monitor Farming

```bash
# Watch farmer
tail -f logs/farmer.log

# Check recent wins
grep "Block submitted" logs/farmer.log | tail -10

# See current challenge
curl http://localhost:8080/challenge | jq
```

### Monitor Node Health

```bash
# Node log
tail -f logs/node.log

# Recent blocks
grep "Accepted block" logs/node.log | tail -10

# P2P status
grep "p2p" logs/node.log | tail -20

# Consensus heartbeat
grep "consensus" logs/node.log | tail -10
```

---

## Farming Guide

### Plot Size Selection

**Considerations:**
- Larger plots = More lottery tickets = Higher win probability
- But: Longer scan times
- Trade-off: Storage vs. scan speed

**Recommendations:**
- **Testing:** k=16 (2MB, fast scans)
- **Small farm:** k=18 (8MB)
- **Medium farm:** k=20 (32MB)
- **Large farm:** k=22+ (128MB+)

### Expected Earnings

**Factors:**
- Your plot size vs. total network plots
- Current difficulty
- Luck (lottery-based)

**Example (k=20 on small network):**
- Win rate: ~1-5 blocks/hour
- Earnings: 20-100 RCHV/hour
- Varies with network size!

### Maximizing Rewards

**Strategies:**
1. **More plots:** Create multiple k=18 or k=20 plots
2. **Uptime:** Keep node and farmer running 24/7
3. **Network:** Good internet connection for block propagation
4. **Timelord:** Run timelord to help network (no direct reward yet)

---

## Troubleshooting

### Node Issues

**Port already in use:**
```bash
# Check what's using the port
lsof -i :8080
lsof -i :9090

# Use different ports
./archivas-node --rpc :8081 --p2p :9091 ...
```

**Sync stalled:**
```bash
# Check P2P connections
grep "connected to peer" logs/node.log

# Check for errors
grep "error\|failed" logs/node.log | tail -20

# Verify bootnode connectivity
nc -zv seed.archivas.ai 9090
```

**Database corruption:**
```bash
# Nuclear option - resync from scratch
rm -rf data
nohup ./archivas-node ... &
# Will sync all blocks from peers
```

### Farmer Issues

**No blocks found:**
- **Normal!** Farming is a lottery
- Verify farmer is scanning: `tail logs/farmer.log`
- Should show: `⚙️ Checking plots... best=XXXXX`
- Quality values should change each round

**Plot not loading:**
```bash
# Verify plot file exists
ls -lh plots/*.arcv

# Check farmer log
grep "Loaded.*plot" logs/farmer.log

# Recreate if corrupted
./archivas-farmer plot --size 18 --path ./plots-new
```

**Submission errors:**
```bash
# Check node is reachable
curl http://localhost:8080/challenge

# Verify farmer connected
grep "Archivas Farmer Starting" logs/farmer.log

# Check for error messages
grep "Error submitting" logs/farmer.log
```

### Timelord Issues

**404 errors:**
```bash
# Verify node is running
curl http://localhost:8080/chainTip

# Check timelord connected
grep "Using node RPC" logs/timelord.log

# Rebuild timelord with latest code
go build -o archivas-timelord ./cmd/archivas-timelord
```

---

## Best Practices

### Security

**Do:**
- ✅ Keep private keys secure (never share!)
- ✅ Use strong passwords for servers
- ✅ Enable firewall (allow 8080, 9090)
- ✅ Monitor logs for errors
- ✅ Back up your wallet

**Don't:**
- ❌ Share private keys
- ❌ Run as root
- ❌ Expose RPC to public internet (use localhost or firewall)
- ❌ Reuse keys across networks

### Performance

**Optimize:**
- Use SSD for plots (faster seeks)
- Keep database on fast storage
- Adequate RAM (2GB+ recommended)
- Stable internet connection
- Low-latency connection to bootnodes

**Monitor:**
```bash
# CPU usage
top -p $(pgrep archivas-node)

# Disk I/O
iostat -x 5

# Network
netstat -an | grep 9090
```

---

## Advanced Topics

### Running Multiple Nodes

**On same machine:**
```bash
# Node 1
./archivas-node --rpc :8080 --p2p :9090 --db ./data1 ...

# Node 2  
./archivas-node --rpc :8081 --p2p :9091 --db ./data2 ...
```

### Custom Genesis

**For private testnet:**
```bash
# Edit genesis/custom.genesis.json
# Change timestamp, allocations, seed

# Start with custom genesis
./archivas-node --genesis genesis/custom.genesis.json --network-id my-private-net ...
```

### Timelord Competition

**Run faster timelord:**
```bash
# Increase step size (more iterations/tick)
./archivas-timelord --node http://localhost:8080 --step 1000
```

**Note:** All timelords compete; fastest one's VDF is used by network

---

## Testnet Etiquette

**Be a good citizen:**
1. Keep your node online and synced
2. Report bugs on GitHub Issues
3. Share farming results (helps calibrate difficulty)
4. Help new users in Discussions
5. Contribute code improvements

**Don't:**
1. Spam the network with invalid blocks
2. Attack or DoS other nodes
3. Hoard all the RCHV (it's testnet, share!)

---

## Getting Help

**Resources:**
- GitHub Issues: Bug reports
- GitHub Discussions: Questions and ideas
- README.md: Quick reference
- Discord: (Coming soon)

**Before asking:**
1. Check logs: `tail -100 logs/node.log`
2. Verify genesis: `curl http://localhost:8080/genesisHash`
3. Check sync: `curl http://localhost:8080/chainTip`
4. Search existing issues

**When reporting bugs:**
- Include log snippets
- Describe steps to reproduce
- Mention your setup (OS, Go version, node version)
- Check if others have same issue

---

**Ready to dive deeper?** [Developer Docs →](Developer-Docs.md)  
**Back:** [← Consensus](Consensus.md)

