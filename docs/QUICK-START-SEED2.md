# Quick Start: Deploy Isolated Seed2

**Time Required:** ~15 minutes  
**Prerequisites:** SSH access to both servers

---

## Step 1: Deploy to Server A (Seed1) - 5 minutes

```bash
# SSH to Server A
ssh ubuntu@57.129.148.132

# Pull and deploy
cd ~/archivas
git pull
./deploy/seed/deploy-p2p-isolation.sh

# Verify
curl -s http://127.0.0.1:8080/chainTip | jq
```

**Expected:** Node restarts, continues operating normally with peer discovery enabled.

---

## Step 2: Deploy to Server D (Seed2) - 10 minutes

```bash
# SSH to Server D
ssh ubuntu@51.68.54.45

# Pull and deploy
cd ~/archivas
git pull
./deploy/seed2/deploy-isolated.sh

# Monitor deployment (auto-runs at end of script)
sudo journalctl -u archivas-node-seed2 -f
```

**Expected Logs:**
```
[p2p] peer discovery DISABLED - only whitelisted peers allowed
[p2p] peer whitelist enabled: 6 entries
[gossip] peer discovery disabled, skipping gossip routine
[p2p] connecting to peer seed.archivas.ai:9090
```

Press `Ctrl+C` to exit log viewer.

---

## Step 3: Verify Isolation - 2 minutes

```bash
# On Server D:

# Should show exactly 1 peer
curl -s http://127.0.0.1:8080/peers | jq 'length'

# Should show isolation logs
sudo journalctl -u archivas-node-seed2 -n 50 | grep -E "gossip|whitelist|DISABLED"

# Monitor sync (run for a few minutes)
watch -n 10 'echo "Seed2: $(curl -s http://127.0.0.1:8080/chainTip | jq -r .height) | Seed1: $(curl -s https://seed.archivas.ai/chainTip | jq -r .height)"'
```

**Success:**
- ✅ Exactly 1 peer connected
- ✅ Logs show "peer discovery DISABLED"
- ✅ Height increases steadily
- ✅ No "prev hash mismatch" errors

---

## Monitoring Commands (Keep These Handy)

```bash
# Watch sync progress
watch -n 5 'curl -s http://127.0.0.1:8080/chainTip | jq'

# Compare heights
watch -n 10 'echo "Seed2: $(curl -s http://127.0.0.1:8080/chainTip | jq -r .height) | Seed1: $(curl -s https://seed.archivas.ai/chainTip | jq -r .height)"'

# Check for fork
sudo journalctl -u archivas-node-seed2 -f | grep -E "prev hash|fork|reorg"

# View all logs
sudo journalctl -u archivas-node-seed2 -f
```

---

## If Something Goes Wrong

### Seed1 Issues

```bash
# Check status
ssh ubuntu@57.129.148.132
sudo systemctl status archivas-node

# View logs
sudo journalctl -u archivas-node -n 100 --no-pager

# Rollback (if needed)
sudo systemctl stop archivas-node
# restore from backup or redeploy previous version
```

### Seed2 Issues

```bash
# Simply disable
ssh ubuntu@51.68.54.45
sudo systemctl stop archivas-node-seed2
sudo systemctl disable archivas-node-seed2
```

---

## Full Documentation

For detailed information, see:
- [Full Deployment Guide](./DEPLOYMENT-P2P-ISOLATION.md)
- [Technical Design](./P2P-DISCOVERY-ISOLATION.md)
- [Changelog](../CHANGELOG-v1.2.0.md)

