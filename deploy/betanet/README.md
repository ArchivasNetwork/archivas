# Archivas Betanet Deployment Package

**Target Server**: 72.251.11.191  
**Network**: Betanet (archivas-betanet-1)  
**Protocol**: v2 (EVM-enabled)

---

## 📦 Package Contents

```
betanet/
├── README.md                         # This file
├── DEPLOYMENT_GUIDE.md               # Complete deployment guide
├── config.toml                       # Node configuration (reference)
├── archivas-betanet.service          # Systemd service for node
├── archivas-betanet-farmer.service   # Systemd service for farmer
├── install.sh                        # Automated installation script
├── verify.sh                         # Identity verification script
└── status.sh                         # Node status script
```

---

## 🚀 Quick Start

### For Server 72.251.11.191

```bash
# 1. SSH into server
ssh root@72.251.11.191

# 2. Clone repository
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas/deploy/betanet

# 3. Run installation
sudo bash install.sh

# 4. Bootstrap from snapshot
sudo -u archivas archivas-node bootstrap --network betanet --db /var/lib/archivas/betanet

# 5. Start node
sudo systemctl start archivas-betanet
sudo systemctl enable archivas-betanet

# 6. Verify
sudo bash verify.sh

# 7. Check status
bash status.sh
```

---

## 📋 Files Description

### config.toml

Reference configuration file showing all available options for Betanet node. The actual node uses CLI flags defined in the systemd service.

**Key Settings**:
- Network: betanet
- Chain ID: archivas-betanet-1
- Network ID: 102
- RPC: 0.0.0.0:8545
- P2P: 0.0.0.0:30303
- Seeds: seed1.betanet.archivas.ai:30303, seed2.betanet.archivas.ai:30303

### archivas-betanet.service

Systemd service file for the Betanet node.

**Features**:
- Runs as `archivas` user
- Auto-restart on failure
- Journald logging
- Security hardening (NoNewPrivileges, PrivateTmp, ProtectSystem)
- Resource limits (65536 open files)
- Connects to Betanet seeds
- EVM enabled

**Installation**:
```bash
sudo cp archivas-betanet.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable archivas-betanet
```

### archivas-betanet-farmer.service

Systemd service for farming (optional).

**Requirements**:
- Node must be running first
- Replace `FARMER_ADDRESS_HERE` with your address
- Plots directory must exist at `/mnt/plots`

**Installation**:
```bash
# Edit farmer address first!
sudo nano archivas-betanet-farmer.service

sudo cp archivas-betanet-farmer.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable archivas-betanet-farmer
```

### install.sh

Automated installation script that:
1. ✅ Installs dependencies (Go, git, curl, etc.)
2. ✅ Creates archivas user and directories
3. ✅ Clones repository
4. ✅ Builds node and farmer binaries
5. ✅ Installs configuration files
6. ✅ Configures firewall (UFW)
7. ✅ Installs systemd services
8. ✅ Provides next steps

**Usage**:
```bash
sudo bash install.sh
```

### verify.sh

Comprehensive verification script that checks:
- ✅ Genesis file (chain ID, network ID, protocol version)
- ✅ Node binary installed
- ✅ Systemd service enabled and running
- ✅ Network ports listening (8545 RPC, 30303 P2P)
- ✅ Firewall configuration
- ✅ RPC responding with correct chain ID (0x66 = 102)
- ✅ Data directory exists
- ✅ Logs confirm betanet network

**Usage**:
```bash
sudo bash verify.sh
```

**Exit Codes**:
- `0` - All checks passed
- `1` - Some checks failed

### status.sh

Real-time status dashboard showing:
- 📊 Service status (running/stopped)
- 🌐 Network information
- 🔌 Port status
- ⛓️ Blockchain height
- 👥 Peer connections
- 💾 Data directory size
- 📝 Recent logs

**Usage**:
```bash
bash status.sh
```

---

## 🔍 Network Identity

### Betanet Configuration

```json
{
  "network": "betanet",
  "chain_id": "archivas-betanet-1",
  "network_id": 102,
  "protocol_version": 2,
  "evm_enabled": true,
  "address_format": "Bech32 (arcv1...)"
}
```

### Verification

**Genesis Hash**: Computed from `genesis-betanet.json`  
**Chain ID (hex)**: `0x66` (102 decimal)  
**RPC Port**: 8545 (Ethereum-compatible)  
**P2P Port**: 30303  

### P2P Identity Enforcement

Betanet nodes will **only** connect to other Betanet nodes:
- ✅ Genesis hash must match
- ✅ Chain ID must match
- ✅ Network ID must match (102)
- ✅ Protocol version must match (v2)

**Incompatible peers are automatically rejected!**

---

## 🌾 Farming Setup

### Prerequisites

1. Node must be fully synced
2. Plots must be created (K32 recommended)
3. Farming address (ARCV or 0x format)

### Steps

```bash
# 1. Prepare plots directory
sudo mkdir -p /mnt/plots
sudo chown -R archivas:archivas /mnt/plots

# 2. Copy/mount your plots to /mnt/plots

# 3. Get your farming address
# You can use either format:
#   Bech32: arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu
#   Hex:    0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0

# 4. Edit farmer service
sudo nano /etc/systemd/system/archivas-betanet-farmer.service
# Replace: FARMER_ADDRESS_HERE

# 5. Start farmer
sudo systemctl start archivas-betanet-farmer
sudo systemctl enable archivas-betanet-farmer

# 6. Check logs
sudo journalctl -u archivas-betanet-farmer -f
```

---

## 🔥 Firewall Configuration

Required open ports:

| Port | Protocol | Purpose | Access |
|------|----------|---------|--------|
| 22 | TCP | SSH | Management |
| 30303 | TCP/UDP | P2P | Public |
| 8545 | TCP | RPC | Public* |

*Consider restricting RPC to localhost or specific IPs for security.

### UFW Commands

```bash
sudo ufw enable
sudo ufw allow 22/tcp
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp
sudo ufw allow 8545/tcp

# Or restrict RPC to localhost:
sudo ufw delete allow 8545/tcp
# Then edit service to use --rpc 127.0.0.1:8545
```

---

## 🔧 Common Commands

### Service Management

```bash
# Start/Stop/Restart
sudo systemctl start archivas-betanet
sudo systemctl stop archivas-betanet
sudo systemctl restart archivas-betanet

# Status
sudo systemctl status archivas-betanet

# Logs (live)
sudo journalctl -u archivas-betanet -f

# Logs (last 100 lines)
sudo journalctl -u archivas-betanet -n 100
```

### Health Checks

```bash
# Run verification
sudo bash verify.sh

# Check status
bash status.sh

# Test RPC
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'
```

### Maintenance

```bash
# Update node
cd /opt/archivas/archivas
sudo -u archivas git pull
sudo -u archivas go build -o /tmp/archivas-node ./cmd/archivas-node/main.go
sudo mv /tmp/archivas-node /usr/local/bin/
sudo systemctl restart archivas-betanet

# Clean logs
sudo journalctl --vacuum-time=7d

# Backup data
sudo rsync -av /var/lib/archivas/betanet /backup/location/
```

---

## 🐛 Troubleshooting

### Node Won't Start

```bash
# Check logs
sudo journalctl -u archivas-betanet -n 100

# Verify permissions
ls -la /var/lib/archivas/betanet

# Check ports
sudo lsof -i :8545
sudo lsof -i :30303
```

### Can't Connect to Peers

```bash
# Test seed connectivity
telnet seed1.betanet.archivas.ai 30303
telnet seed2.betanet.archivas.ai 30303

# Check firewall
sudo ufw status

# Check P2P logs
sudo journalctl -u archivas-betanet | grep -i p2p
```

### Wrong Network

```bash
# CRITICAL: Verify chain ID
cat /etc/archivas/betanet/genesis-betanet.json | jq '.chain_id'
# MUST output: "archivas-betanet-1"

# If wrong, you're on the wrong network!
# Stop and re-bootstrap
sudo systemctl stop archivas-betanet
sudo rm -rf /var/lib/archivas/betanet/*
sudo -u archivas archivas-node bootstrap --network betanet --db /var/lib/archivas/betanet
```

---

## ✅ Verification Checklist

Before considering deployment complete:

- [ ] Installation script completed successfully
- [ ] Node binary responds to `archivas-node version`
- [ ] Genesis file exists at `/etc/archivas/betanet/genesis-betanet.json`
- [ ] Genesis chain ID is `archivas-betanet-1`
- [ ] Genesis network ID is `102`
- [ ] Bootstrap completed successfully
- [ ] Node service is running (`systemctl is-active archivas-betanet`)
- [ ] RPC port 8545 is listening
- [ ] P2P port 30303 is listening
- [ ] Firewall configured correctly
- [ ] RPC responds with chain ID `0x66` (102)
- [ ] Verify script passes all checks
- [ ] Node is syncing blocks
- [ ] Peers are connecting
- [ ] (Optional) Farmer is running and finding challenges

---

## 📚 Additional Resources

- **Complete Deployment Guide**: `DEPLOYMENT_GUIDE.md`
- **Betanet Documentation**: `../../docs/BETANET_COMPLETE.md`
- **Phase 1 Docs**: `../../docs/BETANET_PHASE1.md`
- **Phase 2 Docs**: `../../docs/BETANET_PHASE2.md`
- **Phase 3 Docs**: `../../docs/BETANET_PHASE3.md`
- **Quick Start**: `../../BETANET_QUICKSTART.md`

---

## 🆘 Support

- **Logs**: `sudo journalctl -u archivas-betanet -f`
- **Status**: `bash status.sh`
- **Verify**: `sudo bash verify.sh`
- **Documentation**: `../../docs/`
- **GitHub**: https://github.com/ArchivasNetwork/archivas

---

**Prepared for**: Server 72.251.11.191  
**Network**: Betanet (Production Ready)  
**Date**: November 16, 2025  
**Version**: Phase 1, 2 & 3 Complete ✅

