# Farming with Seed3 (Public Gateway)

This guide shows you how to farm on Archivas Betanet using Seed3, the public gateway node. This is the easiest way to start farming without managing complex network configurations.

## What is Seed3?

Seed3 is a **public gateway node** that:

- ✅ Allows anyone to connect without whitelisting
- ✅ Provides stable, always-on connectivity
- ✅ Handles peer discovery automatically
- ✅ Perfect for community farmers and beginners
- ✅ No private node setup required

**Best for:**
- New farmers getting started
- Small-to-medium farming operations (< 1 TB)
- Users who want simple setup
- Community participation

---

## Quick Start (5 Minutes)

### Prerequisites

- Ubuntu 22.04 LTS (or similar Linux distro)
- 50+ GB free disk space
- Basic command-line knowledge

### Step 1: Install Dependencies & Build

```bash
# Install Go and dependencies
sudo apt update && sudo apt upgrade -y
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
sudo apt install -y git build-essential

# Clone and build
cd ~
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
go build -o archivas-node cmd/archivas-node/main.go
go build -o archivas-farmer cmd/archivas-farmer/main.go
go build -o archivas-wallet cmd/archivas-wallet/main.go
```

### Step 2: Start Your Node

```bash
# Create data directory
mkdir -p ~/.archivas/betanet

# Start node connected to Seed3
nohup ./archivas-node \
  --network betanet \
  --rpc 127.0.0.1:8545 \
  --p2p 0.0.0.0:30303 \
  --db ~/.archivas/betanet \
  --peer seed3.betanet.archivas.ai:30303 \
  --max-peers 50 \
  --enable-gossip > node.log 2>&1 &

# Wait 10 seconds for sync
sleep 10

# Verify node is syncing
curl -s http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq
```

### Step 3: Generate Wallet & Create Plot

```bash
# Generate wallet
./archivas-wallet new
# Save the output! You'll need the public and private keys

# Create 1GB test plot (takes ~40 seconds)
mkdir -p ~/plots
./archivas-farmer plot \
  --size 25 \
  --path ~/plots \
  --farmer-pubkey YOUR_PUBLIC_KEY_HERE
```

### Step 4: Start Farming

```bash
# Start farmer
nohup ./archivas-farmer farm \
  --node http://127.0.0.1:8545 \
  --plots ~/plots \
  --farmer-privkey YOUR_PRIVATE_KEY_HERE > farmer.log 2>&1 &

# Watch for winning blocks
tail -f farmer.log
```

**That's it!** You're now farming on Betanet! 🎉

---

## Understanding the Setup

### Why Connect to Seed3?

```
Your Node → Seed3 (Public Gateway) → Seed1 & Seed2 (Private Seeds)
```

**Benefits:**
1. **No whitelist setup** - Seed3 accepts all connections
2. **Automatic peer discovery** - Your node discovers other seeds via gossip
3. **Stable connectivity** - Seed3 is maintained 24/7
4. **Public participation** - Join the network without special permissions

### Network Architecture

```
┌─────────────────┐
│   Your Node     │ ← You control this
│   + Farmer      │
└────────┬────────┘
         │
         │ P2P Connection
         ↓
┌─────────────────┐
│     Seed3       │ ← Public Gateway (51.89.11.4)
│  Public Access  │
└────────┬────────┘
         │
    ┌────┴────┐
    ↓         ↓
┌───────┐ ┌───────┐
│ Seed1 │ │ Seed2 │ ← Private Seeds
│Farming│ │Backup │
└───────┘ └───────┘
```

---

## Detailed Setup Guide

### Part 1: Node Setup

#### Check System Requirements

```bash
# Check available disk space (need 50+ GB)
df -h ~

# Check RAM (need 4+ GB)
free -h

# Check CPU cores (need 2+)
nproc
```

#### Configure Firewall

```bash
# Allow P2P port for incoming peer connections
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp

# Enable firewall
sudo ufw enable
```

#### Start Node with Monitoring

```bash
cd ~/archivas

# Start in screen for easy monitoring
screen -S archivas-node

./archivas-node \
  --network betanet \
  --rpc 127.0.0.1:8545 \
  --p2p 0.0.0.0:30303 \
  --db ~/.archivas/betanet \
  --peer seed3.betanet.archivas.ai:30303 \
  --max-peers 50 \
  --enable-gossip

# Detach: Press Ctrl+A then D
# Reattach: screen -r archivas-node
```

#### Verify Sync Progress

```bash
# Check current height
curl -s http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result' | xargs printf "%d\n"

# Compare with Seed3
curl -s http://51.89.11.4:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result' | xargs printf "%d\n"
```

**Wait until your node matches Seed3's height!**

---

### Part 2: Wallet & Plot Creation

#### Generate Secure Wallet

```bash
./archivas-wallet new
```

**Output example:**
```
🔐 New Archivas Wallet Generated

Address:     arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs
Public Key:  024aedb4b79bf799cf484a7369b151f6fb4e1988d745e91b6a9fd9d9eb195a7359
Private Key: 1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899

⚠️  KEEP YOUR PRIVATE KEY SECRET!
```

**Important - Save These Securely:**

```bash
# Create encrypted backup
mkdir -p ~/wallet-backup
echo "Address: arcv1s9m9avx..." > ~/wallet-backup/wallet.txt
echo "Public: 024aedb4b7..." >> ~/wallet-backup/wallet.txt
echo "Private: 1cb7a7ad1c..." >> ~/wallet-backup/wallet.txt

# Encrypt with password
gpg -c ~/wallet-backup/wallet.txt
rm ~/wallet-backup/wallet.txt

# Backup to safe location (USB drive, cloud, etc.)
```

#### Choose Plot Size

For Seed3 farming, recommended sizes:

| Use Case | k Value | Size | Generation Time | Notes |
|----------|---------|------|----------------|-------|
| **Testing** | k=25 | 1 GB | 40 seconds | Quick start |
| **Home Farming** | k=28 | 8 GB | 5 minutes | Good balance |
| **Serious Farming** | k=30 | 32 GB | 20 minutes | Better odds |
| **Large Farm** | k=32 | 128 GB | 1.5 hours | Maximum efficiency |

**Recommendation:** Start with **k=28 (8 GB)** for a good balance.

#### Create Plots

```bash
# Create plots directory
mkdir -p ~/plots

# Single 8 GB plot (k=28)
./archivas-farmer plot \
  --size 28 \
  --path ~/plots \
  --farmer-pubkey YOUR_PUBLIC_KEY_HERE

# Or create multiple plots
for i in {1..4}; do
  ./archivas-farmer plot \
    --size 28 \
    --path ~/plots \
    --farmer-pubkey YOUR_PUBLIC_KEY_HERE
done

# Verify plots created
ls -lh ~/plots/
```

---

### Part 3: Start Farming

#### Start Farmer Service

```bash
cd ~/archivas

# Start in background
nohup ./archivas-farmer farm \
  --node http://127.0.0.1:8545 \
  --plots ~/plots \
  --farmer-privkey YOUR_PRIVATE_KEY_HERE > farmer.log 2>&1 &

# Save process ID
echo $! > farmer.pid
```

#### Monitor Farmer Activity

```bash
# Live logs
tail -f ~/archivas/farmer.log

# Check if farmer is running
ps aux | grep archivas-farmer | grep -v grep

# Count loaded plots
grep "Loaded plots" ~/archivas/farmer.log
```

#### Understanding Farmer Output

**Healthy farming looks like:**

```
🌾 Archivas Farmer
   Node: http://127.0.0.1:8545
   Plots: /home/ubuntu/plots
   Loaded plots:
      - plot-k28.arcv (k=28, 268435456 hashes)
      - plot-k28-2.arcv (k=28, 268435456 hashes)
      - plot-k28-3.arcv (k=28, 268435456 hashes)
      - plot-k28-4.arcv (k=28, 268435456 hashes)
🚜 Starting farming loop...

🔍 NEW HEIGHT 5200 (difficulty: 1000000)
   Challenge: 4f6822b5086d7210...
⚙️  Checking plots...
   Scanning plot plot-k28.arcv...
   Scanning plot plot-k28-2.arcv...
   Scanning plot plot-k28-3.arcv...
   Scanning plot plot-k28-4.arcv...
```

**When you find a block:**

```
🎉 Found winning proof! Quality: 554427 (target: 1000000)
✅ Block submitted successfully for height 5200 (VDF t=0)
```

---

## Winning Your First Block

### How Farming Works

Every ~20 seconds, a new block height is available. Your farmer:

1. **Receives challenge** from the node
2. **Scans all plots** for a matching proof
3. **Checks quality** against network difficulty
4. **Submits proof** if quality is good enough
5. **Earns reward** if block is accepted

### Factors Affecting Success

| Factor | Impact |
|--------|--------|
| **Plot Size** | More space = more chances |
| **Network Difficulty** | Lower = easier to win |
| **Number of Plots** | More plots = more chances |
| **Connection Quality** | Faster response = better chance |

### Expected Win Rate

With 32 GB (4x k=28 plots) on a network with ~1 PB total storage:

- **Wins per day:** ~2-5 blocks
- **RCHV per block:** 20 RCHV
- **Daily earnings:** ~40-100 RCHV

*Note: Actual results vary based on network size and difficulty*

---

## Optimizing Your Farm

### 1. Add More Plots

```bash
# Create additional plots on available space
df -h ~/plots  # Check space

# Generate more plots
./archivas-farmer plot --size 28 --path ~/plots \
  --farmer-pubkey YOUR_PUBLIC_KEY

# Restart farmer to load new plots
pkill archivas-farmer
nohup ./archivas-farmer farm \
  --node http://127.0.0.1:8545 \
  --plots ~/plots \
  --farmer-privkey YOUR_PRIVATE_KEY > farmer.log 2>&1 &
```

### 2. Use Faster Storage

- **NVMe SSD:** Best performance for plot scanning
- **SATA SSD:** Good performance
- **HDD:** Slower but acceptable for completed plots

```bash
# Move plots to faster drive
sudo mkdir -p /mnt/nvme/plots
sudo chown $USER:$USER /mnt/nvme/plots
mv ~/plots/*.arcv /mnt/nvme/plots/

# Update farmer command
# --plots /mnt/nvme/plots
```

### 3. Monitor Node Health

```bash
# Check node is synced
watch -n 30 'curl -s http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\",\"params\":[],\"id\":1}" | jq'

# Check peer count
tail ~/archivas/node.log | grep "total peers"
```

### 4. Automate Restarts

Create a systemd service for reliability:

```bash
# Create service file
sudo tee /etc/systemd/system/archivas-farmer.service > /dev/null <<EOF
[Unit]
Description=Archivas Betanet Farmer
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$HOME/archivas
ExecStart=$HOME/archivas/archivas-farmer farm \
  --node http://127.0.0.1:8545 \
  --plots $HOME/plots \
  --farmer-privkey YOUR_PRIVATE_KEY_HERE
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl enable archivas-farmer
sudo systemctl start archivas-farmer

# Check status
sudo systemctl status archivas-farmer
```

---

## Troubleshooting

### Node Won't Connect to Seed3

**Error:** `Failed to connect to peer seed3.betanet.archivas.ai:30303`

**Solutions:**

1. **Check DNS resolution:**
   ```bash
   ping seed3.betanet.archivas.ai
   # Should show: 51.89.11.4
   ```

2. **Try direct IP:**
   ```bash
   --peer 51.89.11.4:30303
   ```

3. **Check firewall:**
   ```bash
   sudo ufw status
   # Ensure 30303 is allowed
   ```

4. **Use alternative seed:**
   ```bash
   --peer seed1.betanet.archivas.ai:30303
   ```

### Farmer Not Finding Blocks

**Symptom:** Farmer scans plots but never finds winning proofs

**This is normal!** On established networks, finding blocks is rare.

**Solutions:**

1. **Check network difficulty:**
   ```bash
   tail farmer.log | grep "difficulty"
   # If difficulty is high (> 1000000), blocks are harder to find
   ```

2. **Add more plots:**
   ```bash
   ./archivas-farmer plot --size 28 --path ~/plots \
     --farmer-pubkey YOUR_PUBLIC_KEY
   ```

3. **Be patient:** With 32 GB, expect 2-5 wins per day

4. **Check node is synced:**
   ```bash
   curl http://127.0.0.1:8545 -X POST \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
   ```

### Plots Not Loading

**Error:** `Failed to load plot: no such file`

**Solutions:**

```bash
# Verify plots exist
ls -lh ~/plots/

# Check permissions
chmod 644 ~/plots/*.arcv

# Verify plot integrity
file ~/plots/plot-k28.arcv
# Should show: data

# Regenerate if corrupted
rm ~/plots/plot-k28.arcv
./archivas-farmer plot --size 28 --path ~/plots \
  --farmer-pubkey YOUR_PUBLIC_KEY
```

### High CPU Usage

**Symptom:** CPU at 100% constantly

**Causes:**
- Plotting in progress (normal, temporary)
- Too many plots scanning simultaneously

**Solutions:**

1. **Wait for plotting to finish** (check with `htop`)
2. **Reduce plot count temporarily** (move some plots out)
3. **Use more CPU cores** (upgrade server)

---

## Maintenance

### Daily Checks

```bash
# Node status
curl http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Farmer status
tail -20 ~/archivas/farmer.log | grep "NEW HEIGHT"

# Check for blocks won
grep "Block submitted successfully" ~/archivas/farmer.log | wc -l
```

### Weekly Tasks

```bash
# Update node
cd ~/archivas
git pull origin main
go build -o archivas-node cmd/archivas-node/main.go
go build -o archivas-farmer cmd/archivas-farmer/main.go

# Restart services
pkill archivas-farmer
pkill archivas-node

# Start again (commands from Quick Start)
```

### Backup Strategy

```bash
# Backup wallet (critical!)
cp ~/wallet-backup/wallet.txt.gpg /backup/location/

# Backup node database (optional - can resync)
tar -czf betanet-backup-$(date +%Y%m%d).tar.gz ~/.archivas/betanet/

# No need to backup plots - can regenerate
```

---

## Upgrading to Private Node

Once you're comfortable, consider upgrading to a [private node setup](farming-with-private-node.md) for:

- ✅ Better connectivity control
- ✅ Enhanced security
- ✅ Dedicated seed connections
- ✅ Professional farming operations

---

## Support & Community

### Get Help

- **Discord:** [discord.gg/archivas](https://discord.gg/archivas) - General support
- **Telegram Farmers:** [t.me/archivas_farmers](https://t.me/archivas_farmers) - Farming-specific help
- **GitHub Issues:** [github.com/ArchivasNetwork/archivas/issues](https://github.com/ArchivasNetwork/archivas/issues) - Bug reports

### Share Your Success

Found your first block? Share it!
- Post in Discord #farming channel
- Tweet @ArchivasNetwork with #ArchivaBetanet
- Help other new farmers get started

---

## Quick Reference

### Essential Commands

```bash
# Start node
nohup ./archivas-node --network betanet --rpc 127.0.0.1:8545 \
  --p2p 0.0.0.0:30303 --db ~/.archivas/betanet \
  --peer seed3.betanet.archivas.ai:30303 --max-peers 50 \
  --enable-gossip > node.log 2>&1 &

# Start farmer
nohup ./archivas-farmer farm --node http://127.0.0.1:8545 \
  --plots ~/plots --farmer-privkey YOUR_PRIVATE_KEY > farmer.log 2>&1 &

# Check height
curl -s http://127.0.0.1:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq

# View logs
tail -f ~/archivas/farmer.log

# Stop services
pkill archivas-farmer
pkill archivas-node
```

---

**Happy Farming! 🌾**

