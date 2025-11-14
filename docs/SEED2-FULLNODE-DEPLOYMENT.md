# Seed2 Full Node Deployment Guide

**Version**: 1.0  
**Date**: 2025-11-14  
**Status**: Production-Ready

---

## 🎯 Overview

This guide covers the deployment of **Seed2 as a dual-role server**:

1. **Full P2P Node**: Participates in blockchain consensus, serves farmers
2. **RPC Relay**: Cached HTTPS endpoint for web clients (existing functionality)

### Goals

- ✅ Reduce load on Seed1 by distributing P2P connections
- ✅ Improve resilience (farmers can peer with Seed1 **or** Seed2)
- ✅ Faster block propagation and sync times
- ✅ Keep existing relay functionality intact

### Architecture

```
┌─────────────────────────────────────────────────┐
│            seed2.archivas.ai                    │
├─────────────────────────────────────────────────┤
│                                                 │
│  ┌──────────────┐        ┌──────────────┐      │
│  │  Full Node   │        │  Nginx Relay │      │
│  │  (P2P)       │◄───────┤  (HTTPS)     │      │
│  │  :30303      │        │  :443        │      │
│  │  :8082 (RPC) │        │              │      │
│  └──────┬───────┘        └──────────────┘      │
│         │                                       │
│         │  Peers with Seed1                     │
│         └───────────────────►  seed.archivas.ai │
│                                                 │
└─────────────────────────────────────────────────┘

Farmers connect to:
  - seed.archivas.ai:30303 (Seed1)
  - seed2.archivas.ai:30303 (Seed2) ◄── NEW

Web clients use:
  - https://seed2.archivas.ai (existing relay)
```

---

## 📦 Deliverables

All files have been created and are ready for deployment:

### 1. Systemd Unit & Config
- ✅ `services/node-seed2/archivas-node-seed2.service` - Systemd service definition
- ✅ `services/node-seed2/seed2-node.env.template` - Environment variables template

### 2. Scripts
- ✅ `infra/firewall.seed2.sh` - UFW firewall configuration (opens P2P port 30303)
- ✅ `data/bootstrap.sh` - Database bootstrap from Seed1 (rsync-based)
- ✅ `scripts/deploy_seed2_node.sh` - **Automated one-step deployment script**

### 3. Infrastructure
- ✅ `infra/nginx.seed2.conf` - Updated Nginx config with Seed2 node as fallback upstream

### 4. Documentation
- ✅ `docs/seed2-node.md` - Comprehensive operations runbook
- ✅ `docs/farmers.md` - Farmer peer configuration guide
- ✅ `docs/relay.md` - Updated relay documentation (mentions dual-peer setup)

### 5. CI/CD
- ✅ `.github/workflows/seed2-validation.yml` - Automated validation checks:
  - Systemd unit syntax
  - Nginx config syntax
  - Shellcheck for bash scripts
  - Security hardening checks
  - Port/path consistency validation

---

## 🚀 Quick Deployment (One Command)

The easiest way to deploy:

```bash
# On Seed2 server
cd /root/archivas
sudo bash scripts/deploy_seed2_node.sh
```

This script will:
1. Install dependencies (rsync, curl, jq, nginx, ufw)
2. Build archivas-node binary (if not present)
3. Create data directory `/var/lib/archivas/seed2`
4. Fetch checkpoint from Seed1
5. Create environment file `/etc/archivas/seed2-node.env`
6. Install and enable systemd unit
7. Configure firewall (open P2P port 30303)
8. Optionally bootstrap database from Seed1
9. Start the node
10. Run health checks

**Estimated time**: 
- Without bootstrap: 5-10 minutes
- With bootstrap: 30-60 minutes (depends on Seed1 database size)

---

## 📋 Manual Deployment (Step-by-Step)

If you prefer manual control or the automated script fails:

### Step 1: Prerequisites

```bash
# Install dependencies
sudo apt-get update
sudo apt-get install -y rsync curl jq golang-go git nginx certbot ufw
```

### Step 2: Build Node Binary

```bash
cd /root/archivas
go build -o archivas-node ./cmd/archivas-node
```

### Step 3: Create Data Directory

```bash
sudo mkdir -p /var/lib/archivas/seed2
sudo chown root:root /var/lib/archivas/seed2
```

### Step 4: Get Checkpoint from Seed1

```bash
CHECKPOINT_HEIGHT=$(curl -s https://seed.archivas.ai:8081/chainTip | jq -r .height)
CHECKPOINT_HASH=$(curl -s https://seed.archivas.ai:8081/chainTip | jq -r .hash)

echo "Height: $CHECKPOINT_HEIGHT"
echo "Hash: $CHECKPOINT_HASH"
```

### Step 5: Create Environment File

```bash
sudo mkdir -p /etc/archivas
sudo tee /etc/archivas/seed2-node.env > /dev/null << EOF
CHECKPOINT_HEIGHT=$CHECKPOINT_HEIGHT
CHECKPOINT_HASH=$CHECKPOINT_HASH
SEED1_P2P=seed.archivas.ai:30303
EOF
```

### Step 6: Install Systemd Unit

```bash
sudo cp services/node-seed2/archivas-node-seed2.service \
       /etc/systemd/system/archivas-node-seed2.service

sudo systemctl daemon-reload
sudo systemctl enable archivas-node-seed2
```

### Step 7: Configure Firewall

```bash
# Allow P2P port
sudo ufw allow 30303/tcp comment 'Archivas P2P'
sudo ufw allow 30303/udp comment 'Archivas P2P'

# Or run the full firewall script
sudo bash infra/firewall.seed2.sh
```

### Step 8: Bootstrap (Optional but Recommended)

**Option A: Bootstrap from Seed1** (30-60 minutes, faster initial sync):

```bash
sudo bash data/bootstrap.sh
```

This will rsync the database from Seed1.

**Option B: Sync from Scratch** (slower, but doesn't require SSH to Seed1):

Skip this step. The node will sync from genesis/checkpoint when started.

### Step 9: Start the Node

```bash
sudo systemctl start archivas-node-seed2
```

### Step 10: Verify

```bash
# Check service status
sudo systemctl status archivas-node-seed2

# View logs
sudo journalctl -u archivas-node-seed2 -f

# Check RPC
curl -s http://127.0.0.1:8082/chainTip | jq

# Check P2P port
sudo ss -tulpn | grep 30303
```

---

## 🔍 Post-Deployment Verification

### 1. Node is Running

```bash
sudo systemctl status archivas-node-seed2

# Expected: "active (running)"
```

### 2. P2P Port is Accessible

From another machine:

```bash
telnet seed2.archivas.ai 30303

# Expected: Connection established (press Ctrl+C to exit)
```

### 3. RPC is Responding

```bash
curl -s http://127.0.0.1:8082/chainTip | jq

# Expected: {"height": <number>, "hash": "...", ...}
```

### 4. Chain Height is Advancing

```bash
# Wait 30 seconds between checks
curl -s http://127.0.0.1:8082/chainTip | jq .height
# Wait 30 seconds
curl -s http://127.0.0.1:8082/chainTip | jq .height

# Expected: Height should increase
```

### 5. Peer Count is Healthy

```bash
sudo journalctl -u archivas-node-seed2 -n 100 | grep -i peer

# Expected: Should see peer connections (at least Seed1)
```

### 6. Update Nginx (Optional Failover)

To enable Seed2 node as a fallback for read requests:

```bash
# The updated nginx.seed2.conf already includes:
# - upstream seed2_node { server 127.0.0.1:8082; }
# - upstream read_pool { ... with seed2_node as backup }

# Copy updated config
sudo cp /root/archivas/infra/nginx.seed2.conf \
       /etc/nginx/sites-available/archivas-seed2

# Test and reload
sudo nginx -t
sudo systemctl reload nginx
```

---

## 👨‍🌾 Farmer Configuration

Once Seed2 is deployed, update farmer commands to include both peers:

### Before (single peer):

```bash
archivas-farmer \
  --farmer-privkey YOUR_PRIVKEY \
  --plot-dir ./plots \
  --p2p-peer seed.archivas.ai:30303 \
  --no-peer-discovery
```

### After (dual-peer, recommended):

```bash
archivas-farmer \
  --farmer-privkey YOUR_PRIVKEY \
  --plot-dir ./plots \
  --p2p-peer seed.archivas.ai:30303 \
  --p2p-peer seed2.archivas.ai:30303 \
  --no-peer-discovery
```

**Benefits**:
- Load distributed across both seeds
- If Seed1 is overloaded, Seed2 handles the farmer
- Faster sync during IBD
- Better block propagation

### Communicate to Farmers

Update documentation, Discord, and announcements with:

```
🚀 New P2P Seed Node Available!

Farmers can now connect to seed2.archivas.ai:30303 in addition to seed.archivas.ai:30303

Recommended peer configuration:
  --p2p-peer seed.archivas.ai:30303
  --p2p-peer seed2.archivas.ai:30303

This improves resilience and reduces load on the primary seed.
```

---

## 📊 Monitoring

### Key Metrics

| Metric | Target | Alert If |
|--------|--------|----------|
| Chain Height | Advancing | Stalled > 60s |
| Peer Count | 20-100 | < 5 for 10m |
| CPU Usage | < 70% | > 85% for 10m |
| Memory | < 12GB | > 14GB |
| Disk I/O | < 20% wait | > 50% sustained |
| Disk Free | > 20% | < 10% |

### Health Check Commands

```bash
# Chain height
curl -s http://127.0.0.1:8082/chainTip | jq .height

# Peer count (from logs)
sudo journalctl -u archivas-node-seed2 -n 100 | grep -i peer | tail

# Resource usage
systemd-cgtop -n 1 | grep archivas-node-seed2

# Disk usage
df -h /var/lib/archivas/seed2
du -sh /var/lib/archivas/seed2

# Service status
sudo systemctl status archivas-node-seed2
```

### Logs

```bash
# Follow logs
sudo journalctl -u archivas-node-seed2 -f

# Last 100 lines
sudo journalctl -u archivas-node-seed2 -n 100 --no-pager

# Filter errors
sudo journalctl -u archivas-node-seed2 -p err
```

---

## 🆘 Troubleshooting

See `docs/seed2-node.md` for comprehensive troubleshooting guide.

### Quick Fixes

**Node won't start**:
```bash
# Check logs for errors
sudo journalctl -u archivas-node-seed2 -n 50 --no-pager

# Check if ports are in use
sudo ss -tulpn | grep -E '30303|8082'

# Verify binary exists
ls -lh /root/archivas/archivas-node
```

**Sync stalled**:
```bash
# Restart node
sudo systemctl restart archivas-node-seed2

# Check if Seed1 is reachable
telnet seed.archivas.ai 30303
```

**Low peer count**:
```bash
# Check firewall
sudo ufw status | grep 30303

# Restart node
sudo systemctl restart archivas-node-seed2
```

---

## 🔄 Upgrades

### Binary Upgrade

```bash
# Build new binary
cd /root/archivas
git pull
go build -o archivas-node ./cmd/archivas-node

# Restart node
sudo systemctl restart archivas-node-seed2

# Verify
curl -s http://127.0.0.1:8082/chainTip | jq
```

### Database Maintenance

```bash
# Weekly backup (before node upgrade)
sudo systemctl stop archivas-node-seed2
sudo tar -czf /backup/seed2-$(date +%Y%m%d).tar.gz \
    -C /var/lib/archivas/seed2 .
sudo systemctl start archivas-node-seed2
```

---

## ✅ Acceptance Criteria

- [ ] P2P port 30303 is publicly accessible (TCP/UDP)
- [ ] Node syncs to tip and maintains sync
- [ ] Peer count > 5 (at minimum, Seed1 + a few farmers)
- [ ] Farmers can successfully peer with Seed2
- [ ] RPC on localhost:8082 responds correctly
- [ ] Nginx relay continues to serve HTTPS RPC as before
- [ ] Chain height advances steadily
- [ ] No memory leaks or crashes for 24h+
- [ ] Documentation updated (farmers.md, relay.md, seed2-node.md)
- [ ] CI validation passes

---

## 📚 Additional Resources

- **Operations Runbook**: `docs/seed2-node.md`
- **Farmer Guide**: `docs/farmers.md`
- **Relay Documentation**: `docs/relay.md`
- **Deployment Script**: `scripts/deploy_seed2_node.sh`
- **Bootstrap Script**: `data/bootstrap.sh`
- **Firewall Script**: `infra/firewall.seed2.sh`

---

## 🎬 Next Steps After Deployment

1. **Monitor for 24-48 hours**:
   - Watch logs for errors
   - Verify peer connections stable
   - Check resource usage (CPU, memory, disk)

2. **Announce to community**:
   - Update Discord/Telegram/X
   - Provide farmer configuration example
   - Share benefits (resilience, load distribution)

3. **Update public documentation**:
   - Website docs
   - GitHub README
   - FAQ

4. **Set up alerting** (if not already done):
   - Prometheus metrics export on :9102
   - Grafana dashboard
   - PagerDuty/alerting for height stalls, peer drops

5. **Gradual rollout**:
   - Test with a few farmers first
   - Monitor Seed1 load reduction
   - Gradually inform all farmers

---

## 🔐 Security Notes

- P2P port 30303 is **publicly accessible** (required for farmers to connect)
- RPC port 8082 is **localhost only** (not exposed)
- HTTPS relay on 443 continues to use existing SSL cert
- Systemd unit includes security hardening (NoNewPrivileges, PrivateTmp, resource limits)
- No private keys or sensitive data in config files

---

## 🏁 Summary

**Before**: Seed1 handled all P2P + RPC traffic → overloaded  
**After**: Seed1 + Seed2 share P2P load → resilient and scalable

**Capacity Impact**:
- Seed1 P2P load: ~50% reduction (split with Seed2)
- Seed2: New P2P node + existing relay functionality
- Farmers: 2x resilience (can peer with either seed)

**Zero Downtime**: Seed1 continues operating normally; Seed2 is an additive improvement.

---

**Deployment Status**: ✅ Ready for Production  
**Version**: 1.0  
**Last Updated**: 2025-11-14

