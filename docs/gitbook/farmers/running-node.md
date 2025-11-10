# Running a Node

Run your own Archivas node for farming or development.

---

## Why Run a Node?

### Benefits

**For Farmers:**
- ✅ Faster proof submission (localhost latency)
- ✅ No dependency on public RPC
- ✅ Higher block win rate
- ✅ Support network decentralization

**For Developers:**
- ✅ Full control over RPC
- ✅ Access to all endpoints (including internal)
- ✅ Historical data access
- ✅ Custom modifications

**For Network:**
- ✅ More geographic distribution
- ✅ Stronger consensus
- ✅ Better censorship resistance

---

## Requirements

### Minimum
- **Disk:** 5 GB (blockchain data)
- **RAM:** 4 GB
- **CPU:** 2 cores
- **Network:** 5 Mbps

### Recommended
- **Disk:** 20 GB (for growth)
- **RAM:** 8 GB
- **CPU:** 4 cores
- **Network:** 20 Mbps

---

## Installation

### Step 1: Build Node

```bash
# Clone repository
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Build node
go build -o archivas-node ./cmd/archivas-node

# Verify
./archivas-node --help
```

### Step 2: Get Genesis File

```bash
# Already included in repo
cat genesis/devnet.genesis.json
```

The genesis file defines the initial network state. The node will compute and display the genesis hash on startup.

---

## Running the Node

### Basic Command

**Important:** To connect to the Archivas network, you **must** specify at least one bootnode using `--bootnodes`.

```bash
./archivas-node \
  --rpc 127.0.0.1:8080 \
  --p2p :9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --bootnodes seed.archivas.ai:9090
```

**Required flags:**
- `--genesis`: Path to genesis file (required on first start)
- `--network-id`: Network identifier (default: `archivas-devnet-v3`; **must specify `archivas-devnet-v4`** for current devnet)
- `--bootnodes`: **Required** - Bootnode address to connect to the network
  - Use: `seed.archivas.ai:9090` (or `57.129.148.132:9090`)
  - Example: `--bootnodes seed.archivas.ai:9090`

**Optional flags:**
- `--rpc`: RPC listen address (default: `:8080` binds to all interfaces; use `127.0.0.1:8080` for localhost only)
- `--p2p`: P2P listen address (default: `:9090` binds to all interfaces)
- `--db`: Database directory path (default: `./data`)
- `--peer`: Additional peer addresses (e.g., `ip1:9090,ip2:9090`)
- `--vdf-required`: Require VDF proofs in blocks (PoSpace+Time mode)
- `--enable-gossip`: Enable automatic peer discovery (default: `true`)
- `--max-peers`: Maximum number of peer connections (default: `20`)

**Note:** Without `--bootnodes`, your node will start but won't connect to the network. Once connected to bootnodes, the node will automatically discover additional peers via gossip and sync blocks via IBD (Initial Block Download).

### As Background Service

```bash
# Create logs directory
mkdir -p logs

# Start in background
./archivas-node \
  --rpc 127.0.0.1:8080 \
  --p2p :9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --bootnodes seed.archivas.ai:9090 \
  > logs/node.log 2>&1 &

# Check it's running
ps aux | grep archivas-node

# Monitor logs
tail -f logs/node.log
```

---

## Initial Sync

When you first start, the node will sync from peers:

```
💾 Database opened: ./data
🌱 Fresh start from genesis file
   Genesis Hash: <hash>
   Network ID: archivas-devnet-v4

🔄 Syncing from peers...
📥 Downloaded blocks 0-1000 (1.2 MB)
📥 Downloaded blocks 1000-2000 (1.1 MB)
...
✅ Sync complete! Height: <current_height>
```

**Initial sync time:** 5-30 minutes (depends on network speed)

---

## Systemd Service (Linux)

### Create service file

```bash
sudo nano /etc/systemd/system/archivas-node.service
```

**Contents:**
```ini
[Unit]
Description=Archivas Node
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME/archivas

ExecStart=/home/YOUR_USERNAME/archivas/archivas-node \
  --rpc 127.0.0.1:8080 \
  --p2p :9090 \
  --genesis /home/YOUR_USERNAME/archivas/genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --bootnodes seed.archivas.ai:9090

Restart=always
RestartSec=10
StandardOutput=append:/var/log/archivas/node.log
StandardError=append:/var/log/archivas/node-error.log

[Install]
WantedBy=multi-user.target
```

### Enable and start

```bash
# Create log directory
sudo mkdir -p /var/log/archivas
sudo chown YOUR_USERNAME:YOUR_USERNAME /var/log/archivas

# Enable service
sudo systemctl daemon-reload
sudo systemctl enable archivas-node
sudo systemctl start archivas-node

# Check status
sudo systemctl status archivas-node

# View logs
sudo journalctl -u archivas-node -f
```

---

## Connecting Farmer to Local Node

Once your node is running:

```bash
# Point farmer to local node
./archivas-farmer farm \
  --plots ~/archivas-plots \
  --node http://localhost:8080 \
  --farmer-privkey YOUR_PRIVKEY \
  > logs/farmer.log 2>&1 &
```

**Advantages:**
- Faster proof submission (<1ms vs 50-300ms)
- No rate limits
- More control

---

## Monitoring Your Node

### Check Status

```bash
# Via RPC
curl http://localhost:8080/chainTip

# Via logs
tail -f logs/node.log | grep -E "Accepted block|NEW_BLOCK|height"
```

### Metrics

```bash
# Prometheus metrics (local node only)
curl http://localhost:8080/metrics

# Key metrics:
# - archivas_tip_height
# - archivas_peer_count
# - archivas_difficulty
# - archivas_blocks_total
```

**Note:** The `/metrics` endpoint is available on local nodes. The public seed node (`seed.archivas.ai`) blocks `/metrics` for security reasons.

---

## Troubleshooting

### "Cannot acquire directory lock"

**Problem:** Another node process is already running.

**Solution:**
```bash
# Find and kill old process
pkill -f archivas-node

# Wait
sleep 2

# Restart
./archivas-node ...
```

### "No peers connected"

**Problem:** Node is isolated.

**Solution:**
- **Check you specified `--bootnodes`** - This is required to connect to the network
- Verify bootnode address is correct: `seed.archivas.ai:9090` (or `57.129.148.132:9090`)
- Check firewall allows port 9090 (outbound for connections, inbound for accepting peers)
- Verify internet connection
- Test bootnode connectivity: `telnet seed.archivas.ai 9090` or `nc -zv seed.archivas.ai 9090`
- Once connected to the bootnode, gossip will automatically discover additional peers

### "Sync stuck"

**Problem:** Not downloading blocks.

**Solution:**
```bash
# Check logs for errors
tail -100 logs/node.log | grep -i error

# Verify genesis hash matches
curl http://localhost:8080/genesisHash
# Compare with the hash shown in node startup logs

# Restart node
pkill -f archivas-node
./archivas-node ...
```

---

## Node Upgrades

### Update to New Version

```bash
# Pull latest code
cd archivas
git pull origin main

# Rebuild
go build -o archivas-node ./cmd/archivas-node

# Restart (data persists)
pkill -f archivas-node
./archivas-node \
  --rpc 127.0.0.1:8080 \
  --p2p :9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --bootnodes seed.archivas.ai:9090 \
  > logs/node.log 2>&1 &
```

**Node will resume from last saved height** - no re-sync needed!

---

## Security

### RPC Binding

**Important:** The default `--rpc :8080` binds to all interfaces (`0.0.0.0:8080`), exposing RPC to the network.

**Public node (not recommended for most users):**
```bash
--rpc 0.0.0.0:8080  # ❌ Exposes to internet
# or
--rpc :8080  # Same as above (default)
```

**Private node (recommended):**
```bash
--rpc 127.0.0.1:8080  # ✅ Localhost only
```

**Why localhost?** Most users only need RPC for local farmers or development. Exposing RPC publicly requires rate limiting, authentication, and security hardening (see production deployment guides).

### Firewall

```bash
# Allow P2P
sudo ufw allow 9090/tcp

# Block RPC from external
sudo ufw deny 8080/tcp

# Enable firewall
sudo ufw enable
```

---

## Data Management

### Database Location

Default: `./data`

Custom:
```bash
--db /var/lib/archivas/data
```

### Database Size

- **Genesis:** ~1 MB
- **After 10,000 blocks:** ~500 MB
- **After 64,000 blocks:** ~2 GB
- **Growth:** ~30 MB per 1,000 blocks

### Backup

```bash
# Stop node
systemctl stop archivas-node

# Backup database
tar -czf archivas-backup-$(date +%Y%m%d).tar.gz data/

# Restart
systemctl start archivas-node
```

---

## Next Steps

- [Earnings Guide](earnings.md) - Understand rewards
- [Troubleshooting](troubleshooting.md) - Common issues
- [Operations Guide](../operations/deployment.md) - Production deployment

---

**Node running?** Connect your [farmer](setup-farmer.md#step-5-start-farming)!

