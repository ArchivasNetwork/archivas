# Run a Private Node for Farming

## What is a Private Node?

A **private node** is an Archivas full node that you run on your own hardware specifically for farming. Unlike public seed nodes that serve the entire network, a private node:

- ✅ **Syncs from trusted seeds** (seed.archivas.ai and seed2.archivas.ai)
- ✅ **Validates blocks locally** on your machine
- ✅ **Exposes RPC only to your farmer** (localhost)
- ✅ **Does not accept public P2P connections**
- ✅ **Provides lower latency** for block submissions

## Why Run a Private Node?

### Benefits for Farmers

1. **🚀 Lower Latency**
   - Your farmer connects to `http://127.0.0.1:8080` instead of a remote seed
   - Block submission happens in microseconds, not seconds
   - No network congestion or public API rate limits

2. **🛡️ RPC Resilience**
   - No downtime when public seeds are under heavy load
   - Your node stays synced 24/7 regardless of public API status
   - Complete control over your infrastructure

3. **🔍 Better Monitoring**
   - Direct access to node logs and metrics
   - See exactly what your node is doing
   - Debug issues without relying on public seed operators

4. **✅ Enhanced Validation**
   - Verify blocks locally before submitting proofs
   - Trust your own node's view of the chain
   - Detect forks or network issues immediately

### When to Use a Private Node

**✅ Recommended for:**
- Large-scale farmers (10+ TB of plots)
- Farmers with multiple machines
- Users who want maximum reliability
- Farmers experiencing frequent 504 timeouts

**⚠️ Optional for:**
- Small-scale farmers (< 1 TB)
- Casual farmers on residential internet
- Users comfortable with public seed dependency

## Prerequisites

Before setting up a private node, ensure you have:

- **Linux server or VPS** (Ubuntu 20.04+ or Debian 11+ recommended)
- **Minimum specs:**
  - 2+ CPU cores
  - 4 GB RAM (8 GB recommended for fast sync)
  - 50 GB free disk space (grows over time)
  - Stable internet connection (10 Mbps+ recommended)
- **Go 1.21 or later** installed ([install guide](https://golang.org/doc/install))
- **Basic Linux terminal skills** (ssh, systemd, logs)

## Quick Start

### Option 1: Automated Setup (Recommended)

We provide an automated setup script that handles everything:

```bash
# Clone the repository
cd ~
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Run the automated setup
sudo bash deploy/private-node/setup-private-node.sh
```

The script will:
1. Check for Go installation
2. Build `archivas-node` and `archivas-farmer`
3. Install systemd services
4. Prompt for your plots directory and farmer private key
5. Start and enable both services
6. Show you how to check status and logs

**Skip to [Verify Your Setup](#verify-your-setup) after running the script.**

### Option 2: Manual Setup

If you prefer manual control, follow these steps:

#### Step 1: Clone and Build Archivas

```bash
# Clone the repository
cd ~
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Build the node and farmer
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer

# Verify binaries
./archivas-node --version
./archivas-farmer --version
```

#### Step 2: Start the Private Node

Create a data directory for your node:

```bash
mkdir -p ~/archivas/data
```

Start the node with private node configuration:

```bash
./archivas-node \
  --network-id archivas-devnet-v4 \
  --db ~/archivas/data \
  --rpc 127.0.0.1:8080 \
  --p2p 0.0.0.0:9090 \
  --genesis ~/archivas/genesis/devnet.genesis.json \
  --no-peer-discovery \
  --peer-whitelist seed.archivas.ai:9090 \
  --peer-whitelist seed2.archivas.ai:9090
```

**Explanation of flags:**
- `--network-id archivas-devnet-v4`: Connect to the devnet
- `--db ~/archivas/data`: Store blockchain data locally
- `--rpc 127.0.0.1:8080`: Bind RPC to localhost only (secure)
- `--p2p 0.0.0.0:9090`: Listen for P2P on all interfaces
- `--no-peer-discovery`: Disable automatic peer discovery
- `--peer-whitelist`: Only connect to trusted seeds

**🚀 Recommended: Use Snapshot Bootstrap for Instant Sync**

Instead of syncing from genesis (which can take hours or days), use the snapshot bootstrap feature to start at a recent height in minutes:

```bash
# Download and import a recent snapshot from Seed2
./archivas-node bootstrap --network devnet

# This will:
# 1. Download ~500MB snapshot from seed2.archivas.ai
# 2. Verify the checksum
# 3. Extract the database to ~/archivas/data
# 4. Start your node at a recent height (e.g., 1.2M blocks)
```

**After bootstrap completes**, start your node normally:

```bash
./archivas-node \
  --network-id archivas-devnet-v4 \
  --db ~/archivas/data \
  --rpc 127.0.0.1:8080 \
  --p2p 0.0.0.0:9090 \
  --genesis ~/archivas/genesis/devnet.genesis.json \
  --no-peer-discovery \
  --peer-whitelist seed.archivas.ai:9090 \
  --peer-whitelist seed2.archivas.ai:9090
```

Your node will resume syncing from the snapshot height instead of starting from block 0!

**Time saved:** ⏱️ From **hours/days** to **minutes**

#### Step 3: Verify Node is Syncing

In another terminal, check the node's status:

```bash
# Check current height
curl -s http://127.0.0.1:8080/chainTip | jq

# Should show:
# {
#   "height": "1200050",
#   "hash": "def456...",
#   "difficulty": "1000000"
# }
```

Watch the node logs:

```bash
# If running in terminal (Ctrl+C to stop)
# Node will print sync progress:
# ✅ Synced to height 1200100
# ⚙️ Validating blocks...
```

Wait until the node catches up to the current network height.

#### Step 4: Point Your Farmer to the Private Node

Once the node is synced, start your farmer:

```bash
cd ~/archivas

./archivas-farmer farm \
  --plots /path/to/your/plots \
  --node http://127.0.0.1:8080 \
  --farmer-privkey YOUR_PRIVATE_KEY_HERE
```

**Important:** Replace `/path/to/your/plots` with your actual plots directory and `YOUR_PRIVATE_KEY_HERE` with your farmer private key.

You should see:

```
🌾 Archivas Farmer Starting
👨‍🌾 Farmer Address: arcv1...
📁 Plots Directory: /path/to/your/plots
🌐 Node: http://127.0.0.1:8080

✅ Loaded N plot(s)
🚜 Starting farming loop...
```

#### Step 5: Make Services Persistent with Systemd

To keep your node and farmer running after logout and reboots, set up systemd services.

**Create the node service:**

```bash
sudo nano /etc/systemd/system/archivas-node-private.service
```

Paste this content (adjust paths as needed):

```systemd
[Unit]
Description=Archivas Private Node
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=YOUR_USERNAME
Group=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME/archivas

ExecStart=/home/YOUR_USERNAME/archivas/archivas-node \
  --network-id archivas-devnet-v4 \
  --db /home/YOUR_USERNAME/archivas/data \
  --rpc 127.0.0.1:8080 \
  --p2p 0.0.0.0:9090 \
  --genesis /home/YOUR_USERNAME/archivas/genesis/devnet.genesis.json \
  --no-peer-discovery \
  --peer-whitelist seed.archivas.ai:9090 \
  --peer-whitelist seed2.archivas.ai:9090

Restart=always
RestartSec=10
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

Replace `YOUR_USERNAME` with your actual username (e.g., `ubuntu`).

**Create the farmer service:**

```bash
sudo nano /etc/systemd/system/archivas-farmer-private.service
```

Paste this content (adjust paths and keys):

```systemd
[Unit]
Description=Archivas Private Farmer
After=archivas-node-private.service
Wants=archivas-node-private.service

[Service]
Type=simple
User=YOUR_USERNAME
Group=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME/archivas

ExecStart=/home/YOUR_USERNAME/archivas/archivas-farmer farm \
  --plots /path/to/your/plots \
  --node http://127.0.0.1:8080 \
  --farmer-privkey YOUR_PRIVATE_KEY_HERE

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Enable and start the services:**

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable (start on boot)
sudo systemctl enable archivas-node-private
sudo systemctl enable archivas-farmer-private

# Start now
sudo systemctl start archivas-node-private
sudo systemctl start archivas-farmer-private
```

## Verify Your Setup

Check that both services are running:

```bash
# Check node status
sudo systemctl status archivas-node-private

# Check farmer status
sudo systemctl status archivas-farmer-private

# Verify RPC is working
curl -s http://127.0.0.1:8080/chainTip | jq
```

Watch the logs:

```bash
# Node logs
sudo journalctl -u archivas-node-private -f

# Farmer logs
sudo journalctl -u archivas-farmer-private -f
```

You should see:
- **Node logs**: Block validation, P2P connections to seeds
- **Farmer logs**: Plot scanning, winning proofs, successful block submissions

## Firewall Configuration

For security and optimal performance, configure your firewall:

### Recommended Firewall Rules

```bash
# Allow SSH (if remote)
sudo ufw allow 22/tcp

# Allow outbound connections (for syncing from seeds)
sudo ufw default allow outgoing

# OPTIONAL: If you want to accept P2P connections from other nodes
# (Not needed for pure private node)
# sudo ufw allow 9090/tcp
# sudo ufw allow 9090/udp

# Block inbound RPC (keep it localhost-only)
# No rule needed - RPC bound to 127.0.0.1

# Enable firewall
sudo ufw enable
```

### RPC Binding Best Practices

**For a pure private node** (recommended):
```bash
--rpc 127.0.0.1:8080  # Localhost only, most secure
```

**If you need remote access** (advanced):
```bash
--rpc 0.0.0.0:8080  # Listen on all interfaces
# WARNING: Make sure to add authentication or firewall rules!
```

**For production:**
- Use `127.0.0.1:8080` and access via SSH tunnel
- Or use a reverse proxy (Nginx) with authentication
- Never expose RPC to the public internet without protection

## Monitoring Your Private Node

### Health Checks

Create a simple monitoring script:

```bash
#!/bin/bash
# ~/check-node-health.sh

# Check if node is responding
if curl -sf http://127.0.0.1:8080/chainTip >/dev/null; then
  echo "✅ Node RPC is responding"
else
  echo "❌ Node RPC is down!"
  exit 1
fi

# Check if node is syncing
HEIGHT=$(curl -s http://127.0.0.1:8080/chainTip | jq -r '.height')
echo "Current height: $HEIGHT"

# Check if farmer is running
if systemctl is-active --quiet archivas-farmer-private; then
  echo "✅ Farmer is running"
else
  echo "❌ Farmer is not running!"
  exit 1
fi
```

Run it periodically:

```bash
chmod +x ~/check-node-health.sh
watch -n 10 ~/check-node-health.sh
```

### Key Metrics to Monitor

1. **Sync Status**
   ```bash
   curl -s http://127.0.0.1:8080/chainTip | jq '.height'
   ```
   Compare with public seed height to ensure you're in sync.

2. **P2P Connections**
   ```bash
   sudo journalctl -u archivas-node-private -n 100 | grep "peer"
   ```
   Should show connections to seed.archivas.ai and seed2.archivas.ai.

3. **Farmer Block Submissions**
   ```bash
   sudo journalctl -u archivas-farmer-private -n 100 | grep "Block submitted"
   ```
   Count successful submissions.

4. **Resource Usage**
   ```bash
   # CPU and memory
   top -bn1 | grep archivas

   # Disk usage
   du -sh ~/archivas/data
   ```

## Troubleshooting

### Node Won't Sync

**Problem:** Node stuck at low height or not receiving blocks.

**Solutions:**
1. Check P2P connectivity:
   ```bash
   sudo journalctl -u archivas-node-private -n 50 | grep "peer\|connection"
   ```
   You should see connections to seeds.

2. Verify peer whitelist flags:
   ```bash
   ps aux | grep archivas-node
   ```
   Ensure `--peer-whitelist` flags are present.

3. Try adding a checkpoint:
   ```bash
   # Get current height from seed
   curl -s https://seed.archivas.ai:8081/chainTip

   # Update your systemd service with checkpoint flags
   sudo systemctl edit archivas-node-private
   ```

### Farmer Gets 504 Timeouts

**Problem:** Farmer shows `HTTP 504` errors when submitting blocks.

**Solutions:**
1. Check if node is synced:
   ```bash
   curl -s http://127.0.0.1:8080/chainTip
   ```
   Node must be at current network height.

2. Check node logs for errors:
   ```bash
   sudo journalctl -u archivas-node-private -n 100 | grep "error\|ERROR"
   ```

3. Restart the node if stuck:
   ```bash
   sudo systemctl restart archivas-node-private
   ```

### High CPU or Memory Usage

**Problem:** Node using too many resources.

**Solutions:**
1. Reduce memory limits in systemd:
   ```systemd
   [Service]
   MemoryMax=2G
   MemoryHigh=1.5G
   ```

2. Lower Go garbage collector aggressiveness:
   ```systemd
   Environment=GOGC=50
   ```

3. Limit concurrent connections:
   ```bash
   --max-peers 20
   ```

### Node Restarts Frequently

**Problem:** Node crashes and systemd keeps restarting it.

**Check logs:**
```bash
sudo journalctl -u archivas-node-private -n 200 | tail -50
```

**Common causes:**
- Out of memory (add more RAM or swap)
- Database corruption (backup data, delete, resync)
- Peer issues (check `--peer-whitelist` flags)

## Upgrading Your Private Node

To update to the latest version:

```bash
# Stop services
sudo systemctl stop archivas-farmer-private
sudo systemctl stop archivas-node-private

# Update code
cd ~/archivas
git pull

# Rebuild
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer

# Restart services
sudo systemctl start archivas-node-private
sudo systemctl start archivas-farmer-private

# Verify
sudo systemctl status archivas-node-private
sudo systemctl status archivas-farmer-private
```

## Advanced Configuration

### Multiple Farmers on One Node

You can run multiple farmers pointing to the same private node:

```bash
# Farmer 1
./archivas-farmer farm \
  --plots /mnt/disk1/plots \
  --node http://127.0.0.1:8080 \
  --farmer-privkey KEY1

# Farmer 2 (different terminal)
./archivas-farmer farm \
  --plots /mnt/disk2/plots \
  --node http://127.0.0.1:8080 \
  --farmer-privkey KEY2
```

Each farmer will independently scan its plots and submit proofs.

### Running Node on Different Port

If port 8080 is in use:

```bash
--rpc 127.0.0.1:8081

# Update farmer
--node http://127.0.0.1:8081
```

### Using Environment Variables for Secrets

Instead of hardcoding your private key in the systemd file:

```bash
# Create env file
echo "FARMER_PRIVKEY=your_key_here" | sudo tee /etc/archivas/farmer.env

# Update systemd service
[Service]
EnvironmentFile=/etc/archivas/farmer.env
ExecStart=/home/USERNAME/archivas/archivas-farmer farm \
  --plots /path/to/plots \
  --node http://127.0.0.1:8080 \
  --farmer-privkey ${FARMER_PRIVKEY}
```

## FAQ

**Q: Do I still need to connect to seed nodes?**  
A: Yes! Your private node syncs from seed.archivas.ai and seed2.archivas.ai. It doesn't mine or create new blocks on its own - it validates blocks from the network.

**Q: Can I use my private node as a public seed?**  
A: Not with this configuration. Private nodes use `--no-peer-discovery` and only connect to trusted seeds. To run a public seed, see [Seed Node Setup](../operations/seed-node.md).

**Q: How much disk space will I need?**  
A: Currently ~50 GB. The blockchain grows over time, so plan for 100+ GB long-term.

**Q: Can I run this on Raspberry Pi?**  
A: Technically yes, but not recommended. Initial sync is CPU-intensive. Use a Pi 4 with 8GB RAM minimum.

**Q: What if my node falls behind?**  
A: It will automatically catch up by syncing from the seeds. If you're offline for days, sync may take a few hours.

**Q: Do I need to open port 9090?**  
A: No. For a private node, outbound connections are enough. You only need inbound 9090 if you want to help relay blocks to other nodes.

## Support

If you need help:

1. **Check logs first:**
   ```bash
   sudo journalctl -u archivas-node-private -n 200
   sudo journalctl -u archivas-farmer-private -n 200
   ```

2. **Visit our Discord:** [discord.gg/archivas](https://discord.gg/archivas)

3. **GitHub Issues:** [github.com/ArchivasNetwork/archivas/issues](https://github.com/ArchivasNetwork/archivas/issues)

4. **Documentation:** [docs.archivas.ai](https://docs.archivas.ai)

---

## Next Steps

- [Large-Scale Farming Guide](large-scale-farming.md) - Optimize for 10+ TB
- [Troubleshooting](troubleshooting.md) - Common farmer issues
- [Seed Node Setup](../operations/seed-node.md) - Run a public seed (advanced)
- [Monitoring Guide](../operations/monitoring.md) - Set up Prometheus + Grafana

