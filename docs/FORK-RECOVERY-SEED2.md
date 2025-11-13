# Seed2 Fork Recovery - November 13, 2025

## Problem Summary

Seed2 (seed2.archivas.ai) forked from the canonical chain at block **665,772** after initial IBD completed successfully.

### Root Cause

The node was configured with `--p2p 0.0.0.0:9090`, which allowed **inbound P2P connections** from any peer. Two forked farmer nodes connected:
- `72.251.11.191:9090`
- `57.129.148.132:9090`

These peers fed Seed2 invalid blocks, causing it to diverge from the canonical chain maintained by Seed1 (seed.archivas.ai).

### Fork Detection

Binary search revealed:
- **Last good block**: 665,771 ✅
- **First bad block**: 665,772 ❌
- **Divergence**: Hashes for 665,772 didn't match between Seed1 and Seed2

## Solution

### 1. P2P Isolation

Modified the systemd service to **disable inbound P2P connections**:

```bash
sudo sed -i 's/--p2p 0\.0\.0\.0:9090/--p2p 127.0.0.1:9090/' /etc/systemd/system/archivas-node-seed2.service
```

**Before:**
```
--p2p 0.0.0.0:9090  # Accepts connections from anyone
```

**After:**
```
--p2p 127.0.0.1:9090  # Only listens locally (no inbound connections)
--peer seed.archivas.ai:9090  # Only trusts canonical seed
```

This ensures Seed2 **only** syncs from Seed1 and cannot accept blocks from potentially forked peers.

### 2. Database Reset

Since BadgerDB doesn't support partial block deletion, we nuked the entire database:

```bash
sudo systemctl stop archivas-node-seed2
rm -rf ~/archivas/data/*
sudo systemctl daemon-reload
sudo systemctl restart archivas-node-seed2
```

### 3. Verification

After restart, confirmed proper isolation:

```bash
curl -s http://127.0.0.1:8080/peers | jq
```

**Expected output:**
```json
{
  "connected": [
    "seed.archivas.ai:9090"
  ],
  "known": []
}
```

✅ **Result**: Only connected to seed.archivas.ai, no peer discovery enabled.

## Monitoring

### Check Sync Progress

```bash
# Watch height increase
watch -n 10 'curl -s https://seed2.archivas.ai/chainTip | jq'

# Check for any unexpected peers
curl -s https://seed2.archivas.ai/peers | jq

# Monitor logs
sudo journalctl -u archivas-node-seed2 -f
```

### Verify No Fork

To verify Seed2 is on the correct chain, compare block hashes at any height:

```bash
HEIGHT=100000
echo "Seed1 hash: $(curl -s "https://seed.archivas.ai/block/$HEIGHT" | jq -r '.hash')"
echo "Seed2 hash: $(curl -s "https://seed2.archivas.ai/block/$HEIGHT" | jq -r '.hash')"
```

Hashes should match! ✅

## Expected Behavior

- **Sync time**: 4-6 hours (with aggressive GC settings)
- **Sync rate**: 200-300 blocks/minute
- **Peers**: Always exactly 1 (seed.archivas.ai:9090)
- **Fork risk**: Zero (no inbound connections, only trusted peer)

## Future Prevention

For any new seed nodes (e.g., Seed3):

1. **Always use P2P isolation**: `--p2p 127.0.0.1:9090`
2. **Always specify trusted peer**: `--peer seed.archivas.ai:9090`
3. **Never use**: `--p2p 0.0.0.0:9090` (unless you're the canonical seed)
4. **Monitor peers regularly**: Ensure no unexpected connections

## Lessons Learned

1. **Inbound P2P connections are risky** for non-canonical seeds during early network stages when forks are common.
2. **P2P isolation** is essential for secondary seed nodes to prevent accepting blocks from forked peers.
3. **Peer whitelisting** (or blacklisting) would be a valuable future feature for the node software.
4. **Fork detection** should be automated - consider adding chain tip hash comparison to health checks.

## Related Files

- **Systemd service**: `/etc/systemd/system/archivas-node-seed2.service`
- **Setup script**: `/home/iljanemesis/archivas/deploy/seed2/setup-production.sh`
- **Nginx config**: `/etc/nginx/sites-available/seed2.archivas.ai`

## Status

- **Date fixed**: November 13, 2025 09:55 UTC
- **Current status**: Syncing from genesis (height ~2,560)
- **Expected completion**: ~6 hours
- **Fork risk**: Eliminated ✅

