# Join the Archivas Testnet

**One-command setup to join the Archivas Proof-of-Space-and-Time testnet!**

---

## Quick Start (5 Minutes)

### 1. Install Dependencies

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go git -y

# Or download from https://go.dev/dl/
```

### 2. Clone and Build

```bash
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
```

### 3. Join the Network

```bash
# Start your node (connects to bootnodes automatically!)
./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 57.129.148.132:9090,72.251.11.191:9090 \
  --enable-gossip

# Your node will:
# ✅ Connect to bootnodes
# ✅ Download all blocks (IBD)
# ✅ Discover other peers (gossip)
# ✅ Sync to tip automatically!
```

### 4. Verify You're Synced

```bash
# Check status
curl http://localhost:8080/healthz

# Should show:
# {"ok":true,"height":1878,"difficulty":...,"peers":2}
```

**That's it! You're now part of the Archivas network!** 🌾

---

## Optional: Run a Timelord (VDF Computer)

Help secure the network by computing VDFs:

```bash
./archivas-timelord --node http://localhost:8080
```

---

## Optional: Start Farming (Mining)

### Step 1: Create a Wallet

```bash
go build -o archivas-wallet ./cmd/archivas-wallet
./archivas-wallet new

# Save your private key! ⚠️
```

### Step 2: Get Test RCHV

```bash
# Use the faucet to get 20 free RCHV
curl "http://57.129.148.132:8080/faucet?address=YOUR_ADDRESS"

# Wait 20 seconds for the transaction
curl http://localhost:8080/balance/YOUR_ADDRESS
```

### Step 3: Generate a Plot

```bash
go build -o archivas-farmer ./cmd/archivas-farmer

./archivas-farmer plot \
  --path ./my-plots \
  --size 18 \
  --farmer-pubkey YOUR_PUBKEY

# This creates a plot file (~256KB for k=18)
```

### Step 4: Start Farming

```bash
./archivas-farmer farm \
  --plots ./my-plots \
  --farmer-privkey YOUR_PRIVATE_KEY \
  --node http://localhost:8080

# Your farmer will:
# ✅ Check plots every VDF update
# ✅ Submit winning proofs
# ✅ Earn 20 RCHV per block! 🚜
```

---

## Network Information

**Network Details:**
- **Network ID:** archivas-devnet-v3
- **Genesis Hash:** 11b6fedb68f1da0f312039cd6fae91f4dd861bea942651b0c33590013f5b8a55
- **Block Time:** ~20 seconds
- **Block Reward:** 20 RCHV
- **Token:** RCHV (8 decimals)

**Bootnodes:**
- 57.129.148.132:9090
- 72.251.11.191:9090

**Faucet:**
- http://57.129.148.132:8080/faucet?address=YOUR_ADDRESS
- Rate limit: 1 drip per hour (20 RCHV)

**Explorer:**
- http://57.129.148.132:8082

**Registry:**
- http://57.129.148.132:8088

**Grafana:**
- http://57.129.148.132:3000 (public read-only)

---

## Troubleshooting

### Sync is slow

Normal! Initial Block Download (IBD) takes time. Your node is downloading all 1878+ blocks.

**Monitor:**
```bash
tail -f logs/node.log | grep "Synced block"
```

### No peers connecting

**Check:**
```bash
# Firewall
sudo ufw allow 9090/tcp

# Connectivity
nc -zv 57.129.148.132 9090
```

### Chain won't advance

**Verify:**
- You're on the right genesis (check `/genesisHash`)
- Network ID matches
- Bootnodes are reachable

---

## Ports to Open

**Required:**
- 9090/tcp - P2P networking (must be public)

**Optional:**
- 8080/tcp - RPC (can be localhost-only for security)

---

## System Requirements

**Minimum:**
- CPU: 1 core
- RAM: 512MB
- Disk: 1GB
- Network: Stable internet

**Recommended:**
- CPU: 2+ cores
- RAM: 2GB
- Disk: 10GB
- Network: Low latency

---

## Getting Help

**Resources:**
- GitHub: https://github.com/ArchivasNetwork/archivas
- Issues: https://github.com/ArchivasNetwork/archivas/issues
- Discussions: https://github.com/ArchivasNetwork/archivas/discussions

**Common Questions:**
- How to farm? See above!
- How to get RCHV? Use the faucet!
- How to check balance? `curl http://localhost:8080/balance/YOUR_ADDRESS`

---

## Advanced: Run with Systemd

```bash
# Copy service file
sudo cp contrib/systemd/archivas-node.service /etc/systemd/system/

# Edit paths in the file
sudo nano /etc/systemd/system/archivas-node.service

# Enable and start
sudo systemctl enable archivas-node
sudo systemctl start archivas-node

# Check logs
sudo journalctl -u archivas-node -f
```

---

## Welcome to Archivas! 🌾

You're now part of a Proof-of-Space-and-Time blockchain network.

**What to do next:**
1. Let your node sync (watch the logs!)
2. Get free RCHV from the faucet
3. Generate a plot and start farming
4. Join the community discussions

**Happy farming!** 🚜💎

