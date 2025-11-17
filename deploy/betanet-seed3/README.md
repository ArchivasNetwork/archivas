# Archivas Betanet - Seed 3 (Public Gateway)

**Server:** 51.89.11.4  
**DNS:** seed3.betanet.archivas.ai  
**Role:** Public Gateway Node (Non-Farming)  
**Network:** Betanet (chain-id: archivas-betanet-1, network-id: 1644)

---

## 🌍 Overview

Seed 3 is the **public gateway** for Archivas Betanet. It serves external nodes, farmers, and developers who want to:
- Sync the Betanet blockchain
- Connect to the P2P network
- Access public RPC endpoints
- Build applications on Archivas

**This node does NOT farm blocks** - it purely relays and serves the network.

---

## 🔧 Architecture

### Three-Tier Betanet Design:

1. **Seed 1 (72.251.11.191)** - Core canonical node, private, farming
2. **Seed 2 (57.129.96.158)** - Backup & snapshot server, semi-private
3. **Seed 3 (51.89.11.4)** - **Public gateway, open P2P/RPC** ← YOU ARE HERE

---

## 📋 Deployment Files

This package contains:

- `config.toml` - Node configuration reference
- `archivas-betanet.service` - Systemd service (node only, no farming)
- `install.sh` - Automated installation script
- `verify.sh` - Identity verification script
- `status.sh` - Node status dashboard
- `README.md` - This file

---

## 🚀 Quick Deploy

```bash
# 1. Clone repository
cd ~
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas/deploy/betanet-seed3

# 2. Run installation
sudo bash install.sh

# 3. Verify identity
bash verify.sh

# 4. Monitor status
bash status.sh
```

---

## 🔍 Key Configuration

- **P2P:** 0.0.0.0:30303 (PUBLIC, no whitelist)
- **RPC:** 0.0.0.0:8545 (PUBLIC)
- **Network:** betanet (1644)
- **Genesis:** configs/genesis-betanet.json
- **Peers:** seed1.betanet.archivas.ai, seed2.betanet.archivas.ai
- **Farming:** DISABLED (relay only)

---

## 🌐 For External Users

If you're not part of the Archivas team and want to join Betanet:

**Connect to this seed:**
```bash
archivas-node \
  --network betanet \
  --rpc 0.0.0.0:8545 \
  --p2p 0.0.0.0:30303 \
  --db ./betanet-data \
  --peer seed3.betanet.archivas.ai:30303
```

See `docs/BETANET_QUICKSTART.md` for full public node setup.

---

## 🔐 Security

### Firewall Rules:
- Port 30303 (P2P): OPEN to public
- Port 8545 (RPC): OPEN to public
- Port 22 (SSH): Restricted to admin IPs

### DDoS Protection:
- `--max-peers 100` (limit connections)
- Consider Cloudflare or rate limiting for RPC

---

## 📊 Monitoring

```bash
# Check service
sudo systemctl status archivas-betanet

# View logs
sudo journalctl -u archivas-betanet -f

# Check connected peers
curl -s http://localhost:8545/ | jq

# Check block height
curl -X POST http://localhost:8545/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq
```

---

## 🆘 Support

- **GitHub:** https://github.com/ArchivasNetwork/archivas
- **Discord:** [Your Discord link]
- **Docs:** https://docs.archivas.ai

---

**Last Updated:** 2025-11-17  
**Maintainer:** Archivas Network Team

