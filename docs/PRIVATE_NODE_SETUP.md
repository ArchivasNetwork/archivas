# Running a Private Archivas Node

This guide explains how to run your own **private Archivas node** for farming or validation, using snapshot-based fast-sync to avoid historical corruption issues.

## Table of Contents

- [Why Run a Private Node?](#why-run-a-private-node)
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Step 1: Export a Snapshot (On Seed Server)](#step-1-export-a-snapshot-on-seed-server)
- [Step 2: Transfer the Snapshot](#step-2-transfer-the-snapshot)
- [Step 3: Import the Snapshot (On Private Node)](#step-3-import-the-snapshot-on-private-node)
- [Step 4: Start Your Private Node](#step-4-start-your-private-node)
- [Step 5: Connect Your Farmer](#step-5-connect-your-farmer)
- [Monitoring and Verification](#monitoring-and-verification)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Why Run a Private Node?

Running your own private Archivas node provides several benefits:

- **Lower Latency**: Direct local RPC access (no network hops)
- **Better Reliability**: No dependence on public RPC endpoints or their rate limits
- **Improved Farming**: Faster challenge retrieval and block submission
- **Network Resilience**: Helps decentralize the network
- **Control**: You validate your own blocks and transactions

## Overview

A **private node** is a full Archivas node that:

1. **Syncs from trusted seeds only** (`seed.archivas.ai`, `seed2.archivas.ai`)
2. **Disables peer discovery** (not publicly discoverable)
3. **Uses checkpoint validation** (starts from a recent snapshot instead of genesis)
4. **Exposes local RPC** (for your farmer/wallet/tools to connect to)

The workflow is:

```
1. Export snapshot from Seed2 at height H (e.g., current_tip - 10,000)
2. Transfer snapshot to your private server
3. Import snapshot (fast, no IBD from height 0)
4. Start private node with checkpoint and whitelist flags
5. Node syncs last ~10k blocks via IBD from seeds
6. Connect farmer to local RPC: http://127.0.0.1:8080
```

---

## Prerequisites

- **Linux server** (Ubuntu 20.04+ recommended, or similar)
- **Go 1.20+** (for building `archivas-node` and `archivas-farmer`)
- **2-4 GB RAM** (for the node)
- **50-100 GB disk space** (for database and plots)
- **Basic Linux knowledge** (SSH, systemd, firewall)
- **Access to Seed2 or a trusted snapshot** (see [Step 1](#step-1-export-a-snapshot-on-seed-server))

---

## Step 1: Export a Snapshot (On Seed Server)

**Note**: This step is typically performed by the Archivas team on **Seed2** to create a trusted snapshot. Community members can use the provided snapshot files.

If you have access to Seed2 (or another trusted node at the latest height):

```bash
# On Seed2 server
cd ~/archivas

# Get current chain tip height
CURRENT_TIP=$(curl -s http://127.0.0.1:8080/chainTip | jq -r '.height')
echo "Current tip: $CURRENT_TIP"

# Export snapshot at a safe height (e.g., tip - 10,000 for safety margin)
SNAPSHOT_HEIGHT=$((CURRENT_TIP - 10000))
echo "Exporting snapshot at height: $SNAPSHOT_HEIGHT"

# Export
./archivas-node snapshot export \
  --height $SNAPSHOT_HEIGHT \
  --out ~/archivas-snapshot-$SNAPSHOT_HEIGHT.tar.gz \
  --db ~/archivas/data \
  --network-id archivas-devnet-v4 \
  --desc "Devnet snapshot at height $SNAPSHOT_HEIGHT"

# Verify the snapshot file
ls -lh ~/archivas-snapshot-$SNAPSHOT_HEIGHT.tar.gz
```

**Output**:
```
[snapshot] Exporting snapshot at height 1200000...
[snapshot] Exporting database from /home/ubuntu/archivas/data...
[snapshot] ✓ Exported 2147483648 bytes
[snapshot] ✓ Snapshot saved to: /home/ubuntu/archivas-snapshot-1200000.tar.gz
[snapshot] Metadata: height=1200000 hash=abc123... network=archivas-devnet-v4
```

The snapshot file `archivas-snapshot-XXXXXX.tar.gz` is now ready for distribution.

---

## Step 2: Transfer the Snapshot

Transfer the snapshot file to your private server:

```bash
# From your local machine or the source server
scp user@seed2.archivas.ai:~/archivas-snapshot-1200000.tar.gz ~/

# Then to your private server
scp ~/archivas-snapshot-1200000.tar.gz user@YOUR_SERVER_IP:~/
```

Or, if you have a publicly hosted snapshot:

```bash
# On your private server
cd ~
wget https://snapshots.archivas.ai/archivas-snapshot-1200000.tar.gz
# or
curl -O https://snapshots.archivas.ai/archivas-snapshot-1200000.tar.gz
```

---

## Step 3: Import the Snapshot (On Private Node)

On your **private server**:

```bash
# Clone the Archivas repository
cd ~
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Build the node binary
go build -o archivas-node ./cmd/archivas-node

# Create data directory
mkdir -p ~/archivas/data

# Import the snapshot
./archivas-node snapshot import \
  --in ~/archivas-snapshot-1200000.tar.gz \
  --db ~/archivas/data

# The command will output the checkpoint height and hash to use in the next step
```

**Output**:
```
[snapshot] Importing snapshot from /home/ubuntu/archivas-snapshot-1200000.tar.gz...
[snapshot] Snapshot info:
  Network:     archivas-devnet-v4
  Height:      1200000
  Block Hash:  abc123def456...
  Exported At: 2025-11-14T20:00:00Z
  Type:        state-only
[snapshot] ✓ Imported 2147483648 bytes
[snapshot] ✓ Database restored to: /home/ubuntu/archivas/data
[snapshot] You can now start the node with:
  --checkpoint-height 1200000 \
  --checkpoint-hash abc123def456...
```

**Note the checkpoint height and hash** — you'll need them for Step 4.

---

## Step 4: Start Your Private Node

Now start the node with the **private node profile**:

### Option A: Manual Start (for testing)

```bash
cd ~/archivas

./archivas-node \
  --db ~/archivas/data \
  --rpc 127.0.0.1:8080 \
  --p2p 0.0.0.0:9090 \
  --genesis ~/archivas/genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --no-peer-discovery \
  --peer-whitelist seed.archivas.ai:9090 \
  --peer-whitelist seed2.archivas.ai:9090 \
  --checkpoint-height 1200000 \
  --checkpoint-hash abc123def456...
```

**Explanation of flags**:
- `--rpc 127.0.0.1:8080`: RPC binds to localhost only (secure, not publicly accessible)
- `--p2p 0.0.0.0:9090`: P2P binds to all interfaces (allows outbound connections to seeds)
- `--no-peer-discovery`: Disables automatic peer discovery (private mode)
- `--peer-whitelist`: Only connect to these trusted seeds (repeatable flag)
- `--checkpoint-height/hash`: Validates blocks against this checkpoint

You should see:

```
═══════════════════════════════════════════════════════════════
  🔐 RUNNING IN PRIVATE NODE MODE
═══════════════════════════════════════════════════════════════

  ⛔ Peer Discovery:  DISABLED
  📋 Whitelisted Peers: 2
     1. seed.archivas.ai:9090
     2. seed2.archivas.ai:9090
  📌 Checkpoint:     height=1200000 hash=abc123def456...
     (Will reject blocks that don't match checkpoint)

  ℹ️  This node will ONLY accept connections from whitelisted
     peers and will NOT be discoverable by other nodes.

═══════════════════════════════════════════════════════════════
```

The node will now:
1. Connect to the whitelisted seeds
2. Sync blocks from height 1200001 to the current tip (~10k blocks, takes 5-10 minutes)
3. Start serving RPC on `http://127.0.0.1:8080`

### Option B: Systemd Service (for production)

For a persistent, auto-starting node, create a systemd service:

```bash
sudo tee /etc/systemd/system/archivas-node-private.service > /dev/null <<EOF
[Unit]
Description=Archivas Private Node
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$USER
Group=$USER
WorkingDirectory=$HOME/archivas

ExecStart=$HOME/archivas/archivas-node \\
  --db $HOME/archivas/data \\
  --rpc 127.0.0.1:8080 \\
  --p2p 0.0.0.0:9090 \\
  --genesis $HOME/archivas/genesis/devnet.genesis.json \\
  --network-id archivas-devnet-v4 \\
  --no-peer-discovery \\
  --peer-whitelist seed.archivas.ai:9090 \\
  --peer-whitelist seed2.archivas.ai:9090 \\
  --checkpoint-height 1200000 \\
  --checkpoint-hash abc123def456... \\
  --max-peers 50

LimitNOFILE=65536
MemoryMax=4G
MemoryHigh=3G
TasksMax=100

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Replace $USER, $HOME, and checkpoint values with actual values
sudo systemctl daemon-reload
sudo systemctl enable archivas-node-private
sudo systemctl start archivas-node-private

# Check status
sudo systemctl status archivas-node-private

# View logs
sudo journalctl -u archivas-node-private -f
```

---

## Step 5: Connect Your Farmer

Once the node is synced, point your farmer to the local RPC:

```bash
cd ~/archivas

# Build the farmer binary if you haven't
go build -o archivas-farmer ./cmd/archivas-farmer

# Start farming
./archivas-farmer farm \
  --plots ~/archivas-plots \
  --node http://127.0.0.1:8080 \
  --farmer-privkey YOUR_FARMER_PRIVATE_KEY
```

Or, if using a systemd service:

```bash
sudo tee /etc/systemd/system/archivas-farmer-private.service > /dev/null <<EOF
[Unit]
Description=Archivas Private Farmer
After=network-online.target archivas-node-private.service
Wants=network-online.target archivas-node-private.service

[Service]
Type=simple
User=$USER
Group=$USER
WorkingDirectory=$HOME/archivas

ExecStart=$HOME/archivas/archivas-farmer farm \\
  --plots $HOME/archivas-plots \\
  --node http://127.0.0.1:8080 \\
  --farmer-privkey YOUR_FARMER_PRIVATE_KEY

LimitNOFILE=65536
MemoryMax=1G
MemoryHigh=512M
TasksMax=50

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable archivas-farmer-private
sudo systemctl start archivas-farmer-private
sudo systemctl status archivas-farmer-private
```

---

## Monitoring and Verification

### Check Node Sync Status

```bash
# Check chain tip
curl -s http://127.0.0.1:8080/chainTip | jq

# Compare with Seed1's tip
curl -s https://seed.archivas.ai:8081/chainTip | jq
```

Both should show the same `height`.

### Check P2P Connections

```bash
# Check logs for peer connections
sudo journalctl -u archivas-node-private -n 50 | grep "peer\|connect"
```

You should see connections to `seed.archivas.ai` and `seed2.archivas.ai` only.

### Check Farmer Status

```bash
# Check farmer logs
sudo journalctl -u archivas-farmer-private -f
```

You should see:
- "Fetching current challenge..."
- "Checking N plots..."
- "Proof found! Quality: ..." (when you win)
- "Block submitted successfully" (when accepted)

---

## Troubleshooting

### Node stuck at snapshot height

**Symptom**: Node doesn't sync beyond the snapshot height.

**Solution**:
- Check that you're connected to at least one seed: `sudo journalctl -u archivas-node-private | grep "connected"`
- Verify firewall allows outbound connections on port 9090
- Try restarting the node: `sudo systemctl restart archivas-node-private`

### "Invalid checkpoint" errors

**Symptom**: Logs show "checkpoint mismatch" or "invalid checkpoint hash".

**Solution**:
- Verify you used the correct `--checkpoint-height` and `--checkpoint-hash` from the snapshot import output
- Make sure the snapshot came from a trusted source (Seed2 or official snapshots)

### Farmer getting 504 timeouts

**Symptom**: Farmer finds proofs but block submission times out.

**Solution**:
- Verify the node is fully synced: `curl http://127.0.0.1:8080/chainTip | jq '.height'`
- Check node logs for errors: `sudo journalctl -u archivas-node-private -n 100`
- Restart the node if it's stuck

### "Database is not empty" error on import

**Symptom**: `snapshot import` fails with "database directory is not empty".

**Solution**:
- Clear the database: `rm -rf ~/archivas/data/*`
- Or use `--force` flag: `./archivas-node snapshot import --in snapshot.tar.gz --db ~/archivas/data --force`

---

## FAQ

### Q: Can I run a private node on the same machine as my farmer?

**A**: Yes! In fact, that's the recommended setup for best performance. The node and farmer both run as separate processes/services.

### Q: How much disk space do I need?

**A**: For devnet, 50 GB is sufficient. For mainnet, plan for 100+ GB and growing.

### Q: Do I need to open port 9090 on my firewall?

**A**: You need to allow **outbound** connections on port 9090 (to reach the seeds). You do NOT need to allow inbound connections (since you're running in private mode).

### Q: How often should I update my node?

**A**: Check for updates at least once per week:
```bash
cd ~/archivas
git pull
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
sudo systemctl restart archivas-node-private
sudo systemctl restart archivas-farmer-private
```

### Q: Can I use my private node for wallet transactions?

**A**: Yes! Point your wallet to `http://127.0.0.1:8080` for transaction broadcasting and balance queries.

### Q: What if I want to migrate to a new snapshot?

**A**: Stop the node, export a new snapshot from Seed2, clear your data directory, import the new snapshot, and restart with updated checkpoint values.

---

## Next Steps

- **Monitor your node**: Set up `prometheus` and `grafana` for metrics (see `docs/monitoring.md`)
- **Secure your server**: Configure UFW firewall, disable root login, use SSH keys (see `docs/security.md`)
- **Optimize farming**: Create more plots with `archivas-farmer plot create` (see `docs/gitbook/farmers/plot-creation.md`)

For more help, visit:
- **Documentation**: https://docs.archivas.ai
- **Discord**: https://discord.gg/archivas
- **GitHub**: https://github.com/ArchivasNetwork/archivas

---

**Happy Farming! 🚜**

