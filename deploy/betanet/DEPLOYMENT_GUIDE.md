# Archivas Betanet Node Deployment Guide

**Server**: 72.251.11.191  
**Network**: Betanet (archivas-betanet-1)  
**Protocol**: v2 (EVM-enabled)  
**Date**: November 16, 2025

---

## 📋 Prerequisites

### Server Requirements

- **OS**: Ubuntu 20.04 LTS or newer
- **CPU**: 4+ cores
- **RAM**: 8+ GB
- **Storage**: 500+ GB SSD
- **Network**: 100+ Mbps, public IP
- **IP Address**: 72.251.11.191

### Software Requirements

- Go 1.19+ (installed via snap)
- Git
- curl, wget, jq
- systemd
- UFW firewall

---

## 🚀 Quick Start (Automated)

### Option 1: One-Command Installation

```bash
# Download and run installation script
curl -sSL https://raw.githubusercontent.com/ArchivasNetwork/archivas/main/deploy/betanet/install.sh | sudo bash
```

### Option 2: Manual Installation

Follow the step-by-step guide below.

---

## 📦 Step-by-Step Installation

### Step 1: System Preparation

```bash
# Update system
sudo apt-get update
sudo apt-get upgrade -y

# Install dependencies
sudo apt-get install -y \
    build-essential \
    git \
    curl \
    wget \
    jq \
    ufw \
    snapd \
    net-tools \
    bc

# Install Go
sudo snap install go --classic

# Verify
go version
```

### Step 2: Create User and Directories

```bash
# Create archivas user
sudo useradd -r -s /bin/bash -d /opt/archivas -m archivas

# Create directories
sudo mkdir -p /etc/archivas/betanet
sudo mkdir -p /var/lib/archivas/betanet
sudo mkdir -p /var/log/archivas
sudo mkdir -p /opt/archivas
sudo mkdir -p /mnt/plots

# Set ownership
sudo chown -R archivas:archivas /var/lib/archivas
sudo chown -R archivas:archivas /var/log/archivas
sudo chown -R archivas:archivas /opt/archivas
sudo chown -R archivas:archivas /mnt/plots
```

### Step 3: Clone and Build

```bash
# Clone repository
sudo -u archivas git clone https://github.com/ArchivasNetwork/archivas.git /opt/archivas/archivas

# Build node
cd /opt/archivas/archivas
sudo -u archivas go build -o archivas-node ./cmd/archivas-node/main.go

# Install binary
sudo mv archivas-node /usr/local/bin/
sudo chmod +x /usr/local/bin/archivas-node

# Verify
archivas-node version
```

### Step 4: Install Configuration

```bash
# Copy genesis file
sudo cp /opt/archivas/archivas/configs/genesis-betanet.json /etc/archivas/betanet/

# Copy config (optional, node uses CLI flags)
sudo cp /opt/archivas/archivas/deploy/betanet/config.toml /etc/archivas/betanet/

# Verify genesis
cat /etc/archivas/betanet/genesis-betanet.json | jq '.chain_id'
# Should output: "archivas-betanet-1"
```

### Step 5: Configure Firewall

```bash
# Enable firewall
sudo ufw --force enable

# Set defaults
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow required ports
sudo ufw allow 22/tcp comment 'SSH'
sudo ufw allow 30303/tcp comment 'Archivas P2P TCP'
sudo ufw allow 30303/udp comment 'Archivas P2P UDP'
sudo ufw allow 8545/tcp comment 'Archivas RPC'

# Verify
sudo ufw status numbered
```

Expected output:
```
Status: active

     To                         Action      From
     --                         ------      ----
[ 1] 22/tcp                     ALLOW IN    Anywhere
[ 2] 30303/tcp                  ALLOW IN    Anywhere
[ 3] 30303/udp                  ALLOW IN    Anywhere
[ 4] 8545/tcp                   ALLOW IN    Anywhere
```

### Step 6: Install Systemd Service

```bash
# Copy service file
sudo cp /opt/archivas/archivas/deploy/betanet/archivas-betanet.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Verify service file
sudo systemctl cat archivas-betanet
```

### Step 7: Bootstrap from Snapshot

```bash
# Bootstrap (downloads and verifies snapshot)
sudo -u archivas /usr/local/bin/archivas-node bootstrap \
    --network betanet \
    --db /var/lib/archivas/betanet
```

Expected output:
```
[bootstrap] Fetching manifest from https://seed2.betanet.archivas.ai/betanet/latest.json...
[bootstrap] Manifest info:
  Network:  betanet
  Chain ID: archivas-betanet-1
  Network ID: 102
  Protocol: v2
  Height:   1250000
[bootstrap] Verifying manifest identity...
[bootstrap] ✓ Manifest verification passed
[bootstrap] Downloading snapshot...
[bootstrap] ✓ Checksum verified
[bootstrap] Importing snapshot...
[bootstrap] ✓ Bootstrap complete!
```

### Step 8: Start the Node

```bash
# Start service
sudo systemctl start archivas-betanet

# Enable auto-start on boot
sudo systemctl enable archivas-betanet

# Check status
sudo systemctl status archivas-betanet

# View logs
sudo journalctl -u archivas-betanet -f
```

### Step 9: Verify Installation

```bash
# Run verification script
sudo bash /opt/archivas/archivas/deploy/betanet/verify.sh
```

Expected: All checks should pass ✅

### Step 10: Test RPC Endpoints

```bash
# Test ETH RPC - Get chain ID
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'

# Expected: {"jsonrpc":"2.0","result":"0x66c","id":1}
# 0x66c = 1644 (betanet network ID)

# Test ETH RPC - Get block number
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Test ARCV RPC - Convert address
curl -X POST http://localhost:8545/arcv \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"arcv_toHexAddress","params":["arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu"],"id":1}'
```

---

## 🌾 Farming Setup (Optional)

### Step 1: Prepare Plots

```bash
# Mount your plots drive
sudo mkdir -p /mnt/plots
# ... mount your drive to /mnt/plots ...

# Set permissions
sudo chown -R archivas:archivas /mnt/plots

# Verify plots
ls -lh /mnt/plots/*.plot
```

### Step 2: Get Farming Address

You can use either format:

**Bech32 format (ARCV)**:
```
arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu
```

**Hex format (0x)**:
```
0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0
```

Both formats work! The farmer will accept either.

### Step 3: Configure Farmer Service

```bash
# Edit farmer service
sudo nano /etc/systemd/system/archivas-betanet-farmer.service

# Replace FARMER_ADDRESS_HERE with your address:
#   --farmer-address arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu
# or
#   --farmer-address 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0

# Verify plots directory matches:
#   --plots /mnt/plots
```

### Step 4: Start Farmer

```bash
# Reload systemd
sudo systemctl daemon-reload

# Start farmer
sudo systemctl start archivas-betanet-farmer

# Enable auto-start
sudo systemctl enable archivas-betanet-farmer

# Check status
sudo systemctl status archivas-betanet-farmer

# View logs
sudo journalctl -u archivas-betanet-farmer -f
```

---

## 🔍 Verification & Health Checks

### Network Identity Verification

```bash
# Run comprehensive verification
sudo bash /opt/archivas/archivas/deploy/betanet/verify.sh
```

**Checks**:
- ✅ Genesis file (chain ID, network ID, protocol version)
- ✅ Node binary installed
- ✅ Service running
- ✅ Ports listening (8545 RPC, 30303 P2P)
- ✅ Firewall configured
- ✅ RPC responding with correct chain ID
- ✅ Data directory exists
- ✅ Logs confirm betanet network

### Manual Verification

```bash
# Check genesis file
jq '{chain_id, network_id, protocol_version}' /etc/archivas/betanet/genesis-betanet.json

# Expected output:
{
  "chain_id": "archivas-betanet-1",
  "network_id": 1644,
  "protocol_version": 2
}

# Check service
systemctl is-active archivas-betanet

# Check ports
netstat -tuln | grep -E '8545|30303'

# Check RPC chain ID
curl -s -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' | jq .

# Should return: {"jsonrpc":"2.0","result":"0x66c","id":1}
```

### Status Dashboard

```bash
# Run status script
bash /opt/archivas/archivas/deploy/betanet/status.sh
```

Shows:
- Service status (running/stopped)
- Network information
- Port status
- Blockchain height
- Chain ID
- Gas price
- Peer connections
- Data directory size
- Recent logs

---

## 🔧 Troubleshooting

### Node Won't Start

```bash
# Check logs
sudo journalctl -u archivas-betanet -n 100 --no-pager

# Check for port conflicts
sudo lsof -i :8545
sudo lsof -i :30303

# Verify permissions
ls -la /var/lib/archivas/betanet

# Test node manually
sudo -u archivas /usr/local/bin/archivas-node --network betanet --help
```

### Can't Connect to Peers

```bash
# Verify firewall
sudo ufw status

# Test connectivity to seeds
telnet seed1.betanet.archivas.ai 30303
telnet seed2.betanet.archivas.ai 30303

# Check P2P logs
sudo journalctl -u archivas-betanet | grep -i "p2p\|peer\|handshake"
```

### RPC Not Responding

```bash
# Check if port is listening
netstat -tuln | grep 8545

# Test locally
curl http://localhost:8545/eth

# Check logs for RPC errors
sudo journalctl -u archivas-betanet | grep -i rpc
```

### Wrong Network Identity

```bash
# This is critical! Verify:
cat /etc/archivas/betanet/genesis-betanet.json | jq '.chain_id'
# MUST be: "archivas-betanet-1"

# If wrong, you're on the wrong network!
# Stop node and re-bootstrap
sudo systemctl stop archivas-betanet
sudo rm -rf /var/lib/archivas/betanet/*
sudo -u archivas archivas-node bootstrap --network betanet --db /var/lib/archivas/betanet
```

### Database Corruption

```bash
# Stop node
sudo systemctl stop archivas-betanet

# Backup current data
sudo mv /var/lib/archivas/betanet /var/lib/archivas/betanet.backup

# Re-bootstrap
sudo -u archivas archivas-node bootstrap --network betanet --db /var/lib/archivas/betanet

# Start node
sudo systemctl start archivas-betanet
```

---

## 📊 Monitoring

### Real-Time Logs

```bash
# Follow logs
sudo journalctl -u archivas-betanet -f

# Filter for errors
sudo journalctl -u archivas-betanet | grep -i error

# Filter for P2P
sudo journalctl -u archivas-betanet | grep -i p2p

# Filter for sync progress
sudo journalctl -u archivas-betanet | grep -i "height\|block"
```

### Performance Metrics

```bash
# Check CPU/Memory usage
ps aux | grep archivas-node

# Disk usage
du -sh /var/lib/archivas/betanet

# Network connections
netstat -an | grep -E ':8545|:30303' | wc -l
```

### Log Rotation

Logs are automatically managed by journald. To configure:

```bash
# Edit journald config
sudo nano /etc/systemd/journald.conf

# Set limits:
SystemMaxUse=1G
SystemKeepFree=2G
MaxRetentionSec=1week

# Restart journald
sudo systemctl restart systemd-journald
```

---

## 🔐 Security Best Practices

### 1. Firewall Configuration

```bash
# Minimal required ports only
sudo ufw status numbered

# Should show ONLY:
# - 22 (SSH)
# - 30303 (P2P)
# - 8545 (RPC - consider restricting)
```

### 2. RPC Access Control

For production, consider restricting RPC access:

```bash
# Option 1: Localhost only (most secure)
# Edit systemd service:
sudo nano /etc/systemd/system/archivas-betanet.service
# Change: --rpc 0.0.0.0:8545
# To: --rpc 127.0.0.1:8545

# Option 2: Specific IPs only
sudo ufw delete allow 8545/tcp
sudo ufw allow from 1.2.3.4 to any port 8545 proto tcp
```

### 3. User Isolation

```bash
# Verify services run as archivas user
ps aux | grep archivas-node

# Should show: archivas  (not root!)
```

### 4. Regular Updates

```bash
# Update node software
cd /opt/archivas/archivas
sudo -u archivas git pull
sudo -u archivas go build -o /tmp/archivas-node ./cmd/archivas-node/main.go
sudo mv /tmp/archivas-node /usr/local/bin/archivas-node
sudo systemctl restart archivas-betanet
```

---

## 📚 Quick Reference

### Service Management

```bash
# Start
sudo systemctl start archivas-betanet

# Stop
sudo systemctl stop archivas-betanet

# Restart
sudo systemctl restart archivas-betanet

# Status
sudo systemctl status archivas-betanet

# Enable auto-start
sudo systemctl enable archivas-betanet

# Disable auto-start
sudo systemctl disable archivas-betanet

# View logs
sudo journalctl -u archivas-betanet -f
```

### Node Commands

```bash
# Version
archivas-node version

# Help
archivas-node help

# Bootstrap
archivas-node bootstrap --network betanet --db /var/lib/archivas/betanet

# Manual start (for testing)
archivas-node --network betanet --rpc 0.0.0.0:8545 --p2p 0.0.0.0:30303
```

### RPC Queries

```bash
# Chain ID
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'

# Block number
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Get balance
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xADDRESS","latest"],"id":1}'
```

---

## ✅ Post-Deployment Checklist

- [ ] Node binary installed and working
- [ ] Genesis file in place (betanet)
- [ ] Firewall configured and enabled
- [ ] Systemd service installed
- [ ] Bootstrap completed successfully
- [ ] Node started and running
- [ ] RPC responding on port 8545
- [ ] P2P listening on port 30303
- [ ] Chain ID verified (1644 / 0x66c)
- [ ] Peers connecting
- [ ] Blocks syncing
- [ ] Logs clean (no errors)
- [ ] Verification script passes
- [ ] (Optional) Farmer configured and running

---

## 🆘 Support

- **Documentation**: `/opt/archivas/archivas/docs/`
- **Logs**: `sudo journalctl -u archivas-betanet -f`
- **Status**: `bash status.sh`
- **Verify**: `bash verify.sh`
- **GitHub**: https://github.com/ArchivasNetwork/archivas
- **Discord**: (your discord link)

---

**Last Updated**: November 16, 2025  
**Version**: Betanet Phase 1, 2 & 3  
**Status**: Production Ready ✅

