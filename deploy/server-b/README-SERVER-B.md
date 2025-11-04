# Server B Setup Guide

Server B is a secondary Archivas node + farmer that peers with Server A (seed.archivas.ai).

---

## Overview

**Purpose:** Secondary validator and farmer node  
**Network:** archivas-devnet-v4  
**Peers:** Connects to Server A (57.129.148.132:9090)  
**Services:** archivas-node-b, archivas-farmer-b  

---

## Setup Steps

### 1. Clone and Build

```bash
# On Server B:
cd /home/ubuntu
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Pull latest code
git pull origin main

# Build binaries
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer

# Verify
./archivas-node --help
./archivas-farmer --help
```

### 2. Create Directories

```bash
mkdir -p ~/archivas/data
mkdir -p ~/archivas/logs
mkdir -p ~/archivas/plots-b
```

### 3. Install Systemd Services

Copy service files to systemd:

```bash
# Node service
sudo cp deploy/server-b/archivas-node-b.service /etc/systemd/system/

# Farmer service
sudo cp deploy/server-b/archivas-farmer-b.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload
```

### 4. Start Services

```bash
# Enable services
sudo systemctl enable archivas-node-b
sudo systemctl enable archivas-farmer-b

# Start node first
sudo systemctl start archivas-node-b

# Wait for sync (check logs)
sudo journalctl -u archivas-node-b -f

# Once synced, start farmer
sudo systemctl start archivas-farmer-b
```

---

## Service Management

### Check Status

```bash
sudo systemctl status archivas-node-b archivas-farmer-b
```

### View Logs

```bash
# Node logs
tail -f /var/log/archivas-node-b.log

# Farmer logs
tail -f /var/log/archivas-farmer-b.log

# Or via journalctl
sudo journalctl -u archivas-node-b -f
sudo journalctl -u archivas-farmer-b -f
```

### Restart Services

```bash
sudo systemctl restart archivas-node-b
sudo systemctl restart archivas-farmer-b
```

### Stop Services

```bash
sudo systemctl stop archivas-farmer-b
sudo systemctl stop archivas-node-b
```

---

## Peering with Server A

Server B automatically connects to Server A via the `--bootnodes` flag:

```
--bootnodes 57.129.148.132:9090
```

**Verify peering:**
```bash
# On Server B, check local RPC
curl http://127.0.0.1:8080/chainTip

# Should match Server A
curl https://seed.archivas.ai/chainTip

# Both should show same height (or within 1-2 blocks)
```

---

## Farming

### Create Plots

```bash
# Generate wallet
./archivas-cli keygen
# Save the privkey

# Create k28 plots
for i in {1..3}; do
  ./archivas-farmer plot \
    --size 28 \
    --path ~/archivas/plots-b/plot-k28-$i.arcv \
    --farmer-pubkey YOUR_PUBKEY
done
```

### Update Farmer Service

Edit the service to use your private key:

```bash
sudo nano /etc/systemd/system/archivas-farmer-b.service

# Update the --farmer-privkey line
# Then:
sudo systemctl daemon-reload
sudo systemctl restart archivas-farmer-b
```

---

## Monitoring

### Check if Winning Blocks

```bash
tail -f /var/log/archivas-farmer-b.log | grep "Found winning"
```

### Check Balance

```bash
curl https://seed.archivas.ai/account/YOUR_FARMER_ADDRESS
```

### Network Distribution

With 2 farmers, you should see blocks alternating between:
- Server A farmer
- Server B farmer

Check recent blocks:
```bash
curl https://seed.archivas.ai/blocks/recent?limit=10 | jq '.blocks[] | {height, farmer: .farmer[:20]}'
```

---

## Firewall

```bash
# Allow P2P from Server A
sudo ufw allow from 57.129.148.132 to any port 9090 proto tcp

# Or allow P2P from anywhere
sudo ufw allow 9090/tcp

# Keep RPC local-only (don't allow 8080 from external)
sudo ufw deny 8080/tcp
```

---

## Troubleshooting

### "Cannot sync blocks"

Check peer connection:
```bash
grep "connected to peer" /var/log/archivas-node-b.log
```

Should show connection to 57.129.148.132:9090.

### "Genesis hash mismatch"

Ensure Server B uses the same genesis file:
```bash
curl http://127.0.0.1:8080/genesisHash
curl https://seed.archivas.ai/genesisHash
# Should match: de7ad6cff236a2aa...
```

### Farmer not finding plots

```bash
ls -la ~/archivas/plots-b/
# Ensure plot files exist

# Check farmer logs for load errors
tail -50 /var/log/archivas-farmer-b.log | grep -i plot
```

---

## Decommissioning

If you need to stop and remove Server B:

```bash
sudo systemctl stop archivas-farmer-b archivas-node-b
sudo systemctl disable archivas-farmer-b archivas-node-b
sudo rm /etc/systemd/system/archivas-{node,farmer}-b.service
sudo systemctl daemon-reload
```

---

**Server B is now a full Archivas validator and farmer!** 🌾

