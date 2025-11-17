# Running a Betanet Node

Archivas Betanet is the active test network for the Archivas blockchain. This guide will walk you through setting up and running a Betanet node from scratch.

## What is Betanet?

Betanet is Archivas' second-generation testnet featuring:

- **Native EVM Execution** - Run Ethereum smart contracts
- **Dual Address System** - Bech32 `arcv1...` addresses + EVM `0x...` addresses
- **Proof-of-Space-and-Time Consensus** - Secure, energy-efficient block production
- **Chain ID: 1644** - Compatible with MetaMask and Ethereum tooling
- **Fast Sync** - Initial Block Download (IBD) syncs thousands of blocks in seconds
- **Auto-Reorganization** - Nodes automatically sync to the longest chain

## Prerequisites

### System Requirements

**Minimum:**
- 2 CPU cores
- 4 GB RAM
- 50 GB free disk space
- Ubuntu 22.04 LTS or later
- Stable internet connection

**Recommended:**
- 4+ CPU cores
- 8 GB RAM
- 100 GB SSD storage
- 10 Mbps+ internet

### Software Requirements

- **Go 1.21+** - For building the node binary
- **Git** - For cloning the repository
- **Build tools** - gcc, make, etc.

---

## Installation Steps

### Step 1: Install Dependencies

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version  # Should show: go version go1.21.0 linux/amd64

# Install git and build tools
sudo apt install -y git build-essential
```

### Step 2: Clone Archivas Repository

```bash
# Clone the repository
cd ~
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Verify you're on the main branch
git branch
```

### Step 3: Build the Node Binary

```bash
# Build archivas-node
go build -o archivas-node cmd/archivas-node/main.go

# Verify the binary
./archivas-node --help
```

**Expected output:**
```
Archivas Node - Blockchain node for Archivas Network

Usage:
  archivas-node [flags]

Node Flags:
  --network <name>            Network to join (betanet, devnet-legacy) [default: betanet]
  --rpc <addr>                RPC listen address (default: from network profile)
  --p2p <addr>                P2P listen address (default: from network profile)
  ...
```

### Step 4: Create Data Directory

```bash
# Create betanet data directory
mkdir -p ~/.archivas/betanet
```

---

## Running Your Node

### Option 1: Run in Foreground (For Testing)

```bash
cd ~/archivas

./archivas-node \
  --network betanet \
  --rpc 127.0.0.1:8545 \
  --p2p 0.0.0.0:30303 \
  --db ~/.archivas/betanet \
  --peer seed3.betanet.archivas.ai:30303 \
  --max-peers 50 \
  --enable-gossip
```

**What these flags mean:**
- `--network betanet` - Connect to Betanet (not Devnet)
- `--rpc 127.0.0.1:8545` - RPC listens on localhost only (secure)
- `--p2p 0.0.0.0:30303` - P2P listens on all interfaces (for peer connections)
- `--db ~/.archivas/betanet` - Database storage location
- `--peer seed3.betanet.archivas.ai:30303` - Connect to public Seed3 gateway
- `--max-peers 50` - Maximum number of peer connections
- `--enable-gossip` - Enable peer discovery via gossip protocol

### Option 2: Run in Background with nohup (Recommended)

```bash
cd ~/archivas

nohup ./archivas-node \
  --network betanet \
  --rpc 127.0.0.1:8545 \
  --p2p 0.0.0.0:30303 \
  --db ~/.archivas/betanet \
  --peer seed3.betanet.archivas.ai:30303 \
  --max-peers 50 \
  --enable-gossip > node.log 2>&1 &

# The node is now running in the background
```

**Check if it's running:**
```bash
ps aux | grep archivas-node | grep -v grep
```

**View logs:**
```bash
tail -f ~/archivas/node.log
# Press Ctrl+C to exit log view
```

### Option 3: Run in Screen Session

```bash
# Install screen
sudo apt install -y screen

# Start a screen session
screen -S archivas-node

# Inside screen, run the node
cd ~/archivas
./archivas-node \
  --network betanet \
  --rpc 127.0.0.1:8545 \
  --p2p 0.0.0.0:30303 \
  --db ~/.archivas/betanet \
  --peer seed3.betanet.archivas.ai:30303 \
  --max-peers 50 \
  --enable-gossip

# Detach from screen: Press Ctrl+A, then D
# Reattach later: screen -r archivas-node
```

---

## Understanding the Sync Process

### Initial Block Download (IBD)

When you first start your node, it will sync the blockchain:

```
[IBD] Starting sync: local=0 remote=5200 (0.0% behind, 5200 blocks to download)
[IBD] Downloading: height 512/5200 (9.8% complete, 4688 blocks remaining)
[IBD] Downloading: height 1024/5200 (19.7% complete, 4176 blocks remaining)
...
[IBD] Complete! Synced to height 5200 (gap: 0 blocks)
[IBD] Successfully synced via http://seed3.betanet.archivas.ai:8545
```

**Sync speed:** Betanet can sync 5000+ blocks in under 1 second on fast connections!

### Auto-Sync & Chain Reorganization

Your node automatically checks for longer chains every 30 seconds:

```
[SYNC] Detected longer chain: local=5200 remote=5215 (gap=15)
[IBD] Starting sync: local=5200 remote=5215
[IBD] Complete! Synced to height 5215
[SYNC] Successfully reorganized to height 5215
```

This ensures your node always stays on the canonical chain.

### P2P Peer Discovery

Your node discovers other peers via gossip:

```
[p2p] connected to peer seed3.betanet.archivas.ai:30303
[p2p] received GOSSIP_PEERS: addrs=4 from=seed3.betanet.archivas.ai:30303
[p2p] merged 4 new addrs (known=4, connected=1)
[p2p] connecting to peer seed1.betanet.archivas.ai:30303
[p2p] connecting to peer seed2.betanet.archivas.ai:30303
```

---

## Verifying Your Node

### Check Block Height

```bash
curl -s http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq
```

**Expected output:**
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": "0x1450"  // 5200 in hex
}
```

### Check Chain ID

```bash
curl -s http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' | jq
```

**Expected output:**
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": "0x66c"  // 1644 in hex
}
```

### Check Network ID

```bash
curl -s http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' | jq
```

**Expected output:**
```json
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": "1644"
}
```

### Check Peer Count

```bash
# View logs to see peer connections
tail ~/archivas/node.log | grep "total peers"
```

---

## Firewall Configuration

If you want other nodes to connect to you, open the P2P port:

```bash
# Allow P2P connections
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp

# RPC should stay localhost-only (don't open 8545)
```

---

## Stopping Your Node

### If Running in Foreground:
Press `Ctrl+C` - The node will shut down gracefully.

### If Running in Background:
```bash
# Graceful shutdown
pkill archivas-node

# Verify it stopped
ps aux | grep archivas-node | grep -v grep
```

### If Running in Screen:
```bash
# Reattach to screen
screen -r archivas-node

# Press Ctrl+C to stop
# Or press Ctrl+A then D to detach without stopping
```

---

## Troubleshooting

### Port Already in Use

**Error:** `bind: address already in use`

**Solution:** Another process is using port 30303 or 8545. Find and stop it:
```bash
# Check what's using the port
sudo netstat -tulpn | grep 30303
sudo netstat -tulpn | grep 8545

# Stop the conflicting process
kill <PID>
```

### Cannot Connect to Seed3

**Error:** `Failed to connect to peer seed3.betanet.archivas.ai:30303`

**Solutions:**
1. Check your internet connection
2. Try using Seed1 instead:
   ```bash
   --peer seed1.betanet.archivas.ai:30303
   ```
3. Check if DNS is resolving:
   ```bash
   ping seed3.betanet.archivas.ai
   ```

### Sync is Stuck

**Symptom:** Node shows same height for a long time

**Solutions:**
1. Check if seeds are reachable:
   ```bash
   curl http://seed3.betanet.archivas.ai:8545 -X POST \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
   ```

2. Restart your node - it will resume from last synced height

3. Check logs for errors:
   ```bash
   tail -50 ~/archivas/node.log
   ```

### Database Corruption

**Symptom:** Node crashes with database errors

**Solution:** Wipe database and resync:
```bash
# Stop the node
pkill archivas-node

# Remove database
rm -rf ~/.archivas/betanet/*

# Restart node - will sync from scratch
```

---

## Next Steps

Now that your node is running:

1. **[Start Farming →](farming-with-seed3.md)** - Earn RCHV by farming
2. **[Run a Private Node →](farming-with-private-node.md)** - Set up a private farming node
3. **[Explore the API →](../api/overview.md)** - Build applications on Betanet

---

## Betanet Network Information

| Property | Value |
|----------|-------|
| **Network Name** | Archivas Betanet |
| **Chain ID** | 1644 (0x66c) |
| **Network ID** | 1644 |
| **Protocol Version** | 2 |
| **Genesis Hash** | `74187e4036f7a489` |
| **Target Block Time** | 20 seconds |
| **Currency Symbol** | RCHV |
| **Block Explorer** | Coming soon |

### Public Endpoints

- **Seed3 (Public Gateway):** `51.89.11.4:30303` (P2P), `51.89.11.4:8545` (RPC)
- **DNS:** `seed3.betanet.archivas.ai`

---

## Support

Need help? Join our community:

- **Discord:** [discord.gg/archivas](https://discord.gg/archivas)
- **Telegram:** [t.me/archivas](https://t.me/archivas)
- **GitHub Issues:** [github.com/ArchivasNetwork/archivas/issues](https://github.com/ArchivasNetwork/archivas/issues)

