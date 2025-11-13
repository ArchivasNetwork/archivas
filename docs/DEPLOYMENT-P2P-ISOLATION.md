# P2P Isolation Feature Deployment Guide

**Version:** 1.2.0  
**Date:** November 13, 2025  
**Status:** Ready for Deployment

---

## Overview

This deployment adds three new P2P isolation features to prevent fork propagation:

1. **`--no-peer-discovery`** - Disables automatic peer discovery (mDNS, DHT, gossip)
2. **`--peer-whitelist`** - Only allows connections to/from whitelisted peers
3. **`--checkpoint-*`** - Validates chain compatibility during handshake

---

## Deployment Strategy

### Phase 1: Deploy to Seed1 (Server A)
Update the binary with new features but **don't enable isolation** (Seed1 needs open discovery).

### Phase 2: Deploy to Seed2 (Server D)
Deploy with **full isolation enabled** to prevent forking.

---

## Phase 1: Deploy to Seed1 (seed.archivas.ai)

### Prerequisites
- SSH access to Server A (57.129.148.132)
- Code pushed to main branch with P2P isolation features

### Steps

```bash
# SSH to Server A
ssh ubuntu@57.129.148.132

# Run deployment script
cd ~/archivas
./deploy/seed/deploy-p2p-isolation.sh
```

### Expected Output

```
✅ Chain height: 673XXX
✅ Connected peers: X
```

### Verification

```bash
# Verify new flags are available
/usr/local/bin/archivas-node --help | grep "no-peer-discovery"

# Check node is running
curl -s http://127.0.0.1:8080/chainTip | jq

# Check logs for stability
sudo journalctl -u archivas-node -n 50 --no-pager
```

**Success Criteria:**
- ✅ Node starts successfully
- ✅ Chain height progresses
- ✅ Peers connect normally
- ✅ No error logs

---

## Phase 2: Deploy to Seed2 (seed2.archivas.ai)

### Prerequisites
- Phase 1 completed successfully
- SSH access to Server D (51.68.54.45)
- Seed1 is healthy and at latest block

### Checkpoint Configuration

**Critical:** Use the current chain tip from Seed1 as the checkpoint:

```bash
# On Server A, get current chain tip
CHECKPOINT_HEIGHT=$(curl -s http://127.0.0.1:8080/chainTip | jq -r .height)
CHECKPOINT_HASH=$(curl -s http://127.0.0.1:8080/block/$CHECKPOINT_HEIGHT | jq -r .hash)

echo "Checkpoint: height=$CHECKPOINT_HEIGHT hash=$CHECKPOINT_HASH"
```

Update `deploy/seed2/deploy-isolated.sh` with these values before deploying.

### Steps

```bash
# SSH to Server D
ssh ubuntu@51.68.54.45

# Pull latest code
cd ~/archivas
git pull

# Edit checkpoint in deployment script (if needed)
nano deploy/seed2/deploy-isolated.sh
# Update CHECKPOINT_HEIGHT and CHECKPOINT_HASH

# Run deployment script
./deploy/seed2/deploy-isolated.sh
```

### Expected Output

```
✅ Deployment Complete!

🔍 Isolation Features:
   - peer discovery DISABLED
   - peer whitelist enabled: 6 entries
   - chain checkpoint: height=671992 hash=eb9b255c...
```

### Verification

```bash
# Check logs for isolation confirmation
sudo journalctl -u archivas-node-seed2 -n 50 --no-pager | grep -E "gossip|GATER|whitelist|Isolation"

# Expected logs:
# [p2p] peer discovery DISABLED - only whitelisted peers allowed
# [p2p] peer whitelist enabled: 6 entries
# [p2p] chain checkpoint: height=671992 hash=eb9b255c
# [gossip] peer discovery disabled, skipping gossip routine
# [p2p] whitelisted: seed.archivas.ai:9090
# [p2p] whitelisted: 57.129.148.132:9090

# Check connected peers (should only be Seed1)
curl -s http://127.0.0.1:8080/peers | jq

# Expected: Only 1 peer (seed.archivas.ai or 57.129.148.132)

# Monitor sync progress
watch -n 10 'echo "Seed2: $(curl -s http://127.0.0.1:8080/chainTip | jq -r .height) | Seed1: $(curl -s https://seed.archivas.ai/chainTip | jq -r .height)"'
```

**Success Criteria:**
- ✅ Logs show "peer discovery DISABLED"
- ✅ Only 1 peer connected (seed.archivas.ai)
- ✅ No "GATER rejected" messages
- ✅ Sync progresses without forking
- ✅ Chain heights match within ~10 blocks after full sync

---

## Monitoring & Troubleshooting

### Monitoring Commands

```bash
# On Seed2, watch sync progress
watch -n 5 'curl -s http://127.0.0.1:8080/chainTip | jq'

# Compare Seed1 vs Seed2 heights
watch -n 10 'echo "Seed2: $(curl -s http://127.0.0.1:8080/chainTip | jq -r .height) | Seed1: $(curl -s https://seed.archivas.ai/chainTip | jq -r .height)"'

# Check for fork indicators
sudo journalctl -u archivas-node-seed2 -f | grep -E "prev hash mismatch|fork|reorg"

# Monitor peer connections
watch -n 30 'curl -s http://127.0.0.1:8080/peers | jq "length"'
```

### Troubleshooting

#### Issue: Seed2 connects to multiple peers

**Symptom:**
```bash
curl -s http://127.0.0.1:8080/peers | jq 'length'
# Returns: 2 or more
```

**Cause:** Whitelist includes IPs that resolve to forked nodes.

**Fix:**
```bash
# Edit service file to only include Seed1
sudo nano /etc/systemd/system/archivas-node-seed2.service

# Ensure only these whitelists:
#   --peer-whitelist seed.archivas.ai:9090 \
#   --peer-whitelist 57.129.148.132:9090 \

sudo systemctl daemon-reload
sudo systemctl restart archivas-node-seed2
```

#### Issue: Sync stalls at a specific height

**Symptom:**
```bash
# Height doesn't progress for >5 minutes
```

**Cause:** Possible fork or peer connectivity issue.

**Fix:**
```bash
# Check logs for fork indicators
sudo journalctl -u archivas-node-seed2 -n 100 | grep -E "prev hash|fork"

# Verify peer connection
curl -s http://127.0.0.1:8080/peers | jq

# If forked, perform manual verification:
SEED2_HEIGHT=$(curl -s http://127.0.0.1:8080/chainTip | jq -r .height)
SEED2_HASH=$(curl -s http://127.0.0.1:8080/block/$SEED2_HEIGHT | jq -r .hash)
SEED1_HASH=$(curl -s https://seed.archivas.ai/block/$SEED2_HEIGHT | jq -r .hash)

if [ "$SEED2_HASH" != "$SEED1_HASH" ]; then
    echo "❌ FORK DETECTED at height $SEED2_HEIGHT"
    echo "   Seed2: $SEED2_HASH"
    echo "   Seed1: $SEED1_HASH"
    # Requires full resync
    sudo systemctl stop archivas-node-seed2
    rm -rf ~/archivas/data/*
    sudo systemctl start archivas-node-seed2
fi
```

#### Issue: "GATER rejected" messages in logs

**Symptom:**
```bash
sudo journalctl -u archivas-node-seed2 | grep "GATER rejected"
# Shows rejected connection attempts
```

**Expected:** This is **normal** if other nodes try to connect to Seed2.

**If concerning:** Verify P2P listen address is `127.0.0.1:9090` (not `0.0.0.0:9090`).

#### Issue: No peers connecting

**Symptom:**
```bash
curl -s http://127.0.0.1:8080/peers | jq 'length'
# Returns: 0
```

**Cause:** Whitelist might not include correct address format.

**Fix:**
```bash
# Check DNS resolution
dig +short seed.archivas.ai
# Should return: 57.129.148.132

# Check logs for whitelist entries
sudo journalctl -u archivas-node-seed2 | grep "whitelisted"

# Manually dial peer
curl -X POST http://127.0.0.1:8080/admin/dial -d '{"peer": "seed.archivas.ai:9090"}'
```

---

## Rollback Procedure

### If Seed1 has issues after deployment:

```bash
# SSH to Server A
ssh ubuntu@57.129.148.132

# Restore previous binary (if backed up)
sudo systemctl stop archivas-node
sudo cp ~/archivas/archivas-node.backup /usr/local/bin/archivas-node
sudo systemctl start archivas-node
```

### If Seed2 has issues:

```bash
# SSH to Server D
ssh ubuntu@51.68.54.45

# Simply stop and disable
sudo systemctl stop archivas-node-seed2
sudo systemctl disable archivas-node-seed2
```

---

## Post-Deployment Tasks

### Update Documentation

1. Update farmer guides to include seed2.archivas.ai as backup RPC
2. Update SDK to support multiple RPC endpoints
3. Update Explorer to support failover

### Announce to Community

```
🎉 Archivas Network Update

We've deployed a secondary seed node (seed2.archivas.ai) with enhanced fork protection:

✅ Peer isolation prevents fork propagation
✅ Chain checkpoint validation
✅ Backup RPC endpoint for farmers

Farmers can now configure failover:
- Primary: https://seed.archivas.ai
- Backup: https://seed2.archivas.ai

No action required - your farmer will continue working as normal.
```

---

## Success Metrics

### After 24 Hours

- [ ] Seed1 uptime: 100%
- [ ] Seed2 uptime: 100%
- [ ] Seed2 synced within 10 blocks of Seed1
- [ ] No fork events on Seed2
- [ ] No critical errors in logs
- [ ] Peer count on Seed2: 1 (seed.archivas.ai only)

### After 1 Week

- [ ] Both nodes stable with no manual interventions
- [ ] Seed2 serves RPC traffic successfully
- [ ] No fork propagation observed
- [ ] Community reports no issues

---

## Reference

### New CLI Flags

```bash
--no-peer-discovery
  Disable automatic peer discovery (mDNS, DHT, gossip)
  
--peer-whitelist <addr>
  Whitelisted peer address (repeatable)
  Formats: host:port, IP:port
  Example: --peer-whitelist seed.archivas.ai:9090
  
--checkpoint-height <N>
  Chain checkpoint height for validation
  Example: --checkpoint-height 671992
  
--checkpoint-hash <hex>
  Chain checkpoint hash (64 hex chars, 32 bytes)
  Example: --checkpoint-hash eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79
```

### Checkpoint Data

**Current Checkpoint (as of deployment):**
- Height: 671,992
- Hash: `eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79`

### Server Details

**Seed1 (Server A):**
- IP: 57.129.148.132
- Domain: seed.archivas.ai
- Role: Primary seed, normal peer discovery
- Config: No isolation flags

**Seed2 (Server D):**
- IP: 51.68.54.45
- Domain: seed2.archivas.ai
- Role: Secondary seed, isolated
- Config: `--no-peer-discovery` + whitelist + checkpoint

---

## Related Documents

- [P2P Discovery Isolation Design](./P2P-DISCOVERY-ISOLATION.md)
- [Fork Recovery Documentation](./FORK-RECOVERY-SEED2.md) (archived)
- [Deployment Scripts](../deploy/seed2/deploy-isolated.sh)

---

**Deployment Sign-Off:**

- [ ] Code reviewed and tested locally
- [ ] Checkpoint verified against Seed1
- [ ] Deployment scripts tested
- [ ] Rollback procedure documented
- [ ] Monitoring dashboard configured
- [ ] On-call engineer notified

**Deployed By:** _____________  
**Date:** _____________  
**Sign-Off:** _____________

