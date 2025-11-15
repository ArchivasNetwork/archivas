# Archivas Snapshot System - Deployment Guide

## Overview

The Archivas snapshot system allows new nodes to bootstrap from a recent chain state instead of syncing from genesis. This reduces sync time from hours/days to minutes.

## System Components

1. **Snapshot Export** (`snapshot` package)
   - Exports blockchain state at a safe height
   - Creates compressed tarball with metadata
   - Generates SHA256 checksum

2. **Snapshot Publishing** (`scripts/publish_snapshot_devnet.sh`)
   - Automated script that runs on Seed2
   - Publishes snapshots every 6 hours
   - Serves via HTTPS at `https://seed2.archivas.ai/devnet/`

3. **Bootstrap Command** (`archivas-node bootstrap`)
   - Downloads snapshot from Seed2
   - Verifies checksum
   - Extracts database
   - Node ready to sync from snapshot height

## Architecture

```
┌─────────────────┐
│     Seed2       │
│  (Full Node)    │
└────────┬────────┘
         │
         │ Every 6 hours:
         │ 1. Export snapshot
         │ 2. Compress (~500MB)
         │ 3. Generate checksum
         │ 4. Update latest.json
         │ 5. Clean old snapshots
         │
         ▼
┌─────────────────┐
│  Nginx Server   │
│ seed2.archivas  │
└────────┬────────┘
         │
         │ HTTPS
         │
         ▼
┌─────────────────┐
│   New Nodes     │
│ ./archivas-node │
│    bootstrap    │
└─────────────────┘
```

## Production Deployment

### Prerequisites

✅ Seed2 is running and synced
✅ Nginx is configured to serve `/devnet/` location
✅ SSL certificate is valid for `seed2.archivas.ai`
✅ Repository is up to date on Seed2

### Deployment Steps

#### 1. SSH into Seed2

```bash
ssh ubuntu@seed2.archivas.ai
```

#### 2. Update the Repository

```bash
cd ~/archivas
git pull origin main
```

#### 3. Run the Deployment Script

```bash
cd ~/archivas
./scripts/deploy_snapshot_automation.sh
```

This script will:
- ✅ Create log directory
- ✅ Test snapshot generation
- ✅ Configure cron job (every 6 hours)
- ✅ Verify HTTPS accessibility

#### 4. Verify Deployment

Check that the snapshot is accessible:

```bash
curl -s https://seed2.archivas.ai/devnet/latest.json | jq
```

Expected output:
```json
{
  "network": "devnet",
  "height": 1247470,
  "block_hash": "d3182add55edf1d4b1c06999523b270759fb567a4fa1d4954bf4e591fb35120b",
  "snapshot_url": "https://seed2.archivas.ai/devnet/snap-1247470.tar.gz",
  "checksum": "sha256:abc123...",
  "size_bytes": 504238592,
  "created_at": "2025-11-15T12:00:00Z"
}
```

#### 5. Monitor Logs

```bash
# Watch the cron job logs
tail -f /var/log/archivas/snapshot-publish.log

# Check cron job is active
crontab -l | grep publish_snapshot
```

### Schedule

Snapshots are published automatically every 6 hours:
- 00:00 UTC
- 06:00 UTC
- 12:00 UTC
- 18:00 UTC

## Testing the Bootstrap

### On a New Server

```bash
# Clone the repo
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Build the node
go build -o archivas-node ./cmd/archivas-node

# Bootstrap from snapshot
./archivas-node bootstrap --network devnet

# Start the node
./archivas-node \
  -db ~/archivas/data \
  -genesis ~/archivas/genesis/devnet.genesis.json \
  -p2p 0.0.0.0:9090 \
  -rpc 127.0.0.1:8080 \
  -peer seed.archivas.ai:9090 \
  -peer seed2.archivas.ai:9090
```

### Verification

1. **Check the node started at snapshot height:**
   ```bash
   curl -s http://127.0.0.1:8080/chainTip | jq
   ```

2. **Verify chain consensus:**
   ```bash
   # Get block hash at snapshot height from new node
   NEW_NODE_HASH=$(curl -s http://127.0.0.1:8080/block/<height> | jq -r '.hash')
   
   # Get same block from Seed2
   SEED2_HASH=$(curl -s https://seed2.archivas.ai/block/<height> | jq -r '.hash')
   
   # They should match
   [ "$NEW_NODE_HASH" = "$SEED2_HASH" ] && echo "✅ Chain consensus verified"
   ```

3. **Check sync progress:**
   ```bash
   curl -s http://127.0.0.1:8080/sync/status | jq
   ```

## Maintenance

### Manual Snapshot Generation

If you need to generate a snapshot manually:

```bash
ssh ubuntu@seed2.archivas.ai
cd ~/archivas
./scripts/publish_snapshot_devnet.sh
```

### Cleaning Old Snapshots

The publish script automatically keeps only the 3 most recent snapshots. To manually clean:

```bash
ssh ubuntu@seed2.archivas.ai
ls -lt /srv/archivas-snapshots/devnet/snap-*.tar.gz | tail -n +4 | awk '{print $9}' | xargs rm -f
```

### Updating the Publish Script

1. Make changes to `scripts/publish_snapshot_devnet.sh`
2. Commit and push to GitHub
3. Pull on Seed2:
   ```bash
   ssh ubuntu@seed2.archivas.ai
   cd ~/archivas
   git pull
   ```
4. The cron job will use the updated script on next run

## Troubleshooting

### Snapshot Not Accessible

**Problem:** `curl https://seed2.archivas.ai/devnet/latest.json` returns 404

**Solutions:**
1. Check if snapshot exists:
   ```bash
   ls -lh /srv/archivas-snapshots/devnet/
   ```

2. Verify Nginx config:
   ```bash
   sudo nginx -t
   sudo systemctl reload nginx
   ```

3. Check Nginx logs:
   ```bash
   sudo tail -f /var/log/nginx/error.log
   ```

### Bootstrap Fails with Checksum Error

**Problem:** Bootstrap fails with "checksum mismatch"

**Solutions:**
1. Snapshot was corrupted during download - retry
2. Snapshot was updated during download - retry
3. Check disk space on both servers

### Cron Job Not Running

**Problem:** No new snapshots being created

**Solutions:**
1. Check cron status:
   ```bash
   systemctl status cron
   ```

2. Verify cron job:
   ```bash
   crontab -l | grep publish_snapshot
   ```

3. Check logs:
   ```bash
   tail -100 /var/log/archivas/snapshot-publish.log
   ```

4. Test manual run:
   ```bash
   /home/ubuntu/archivas/scripts/publish_snapshot_devnet.sh
   ```

## Monitoring

### Key Metrics to Monitor

1. **Snapshot File Size**
   - Should be ~500MB (will grow with chain size)
   - Check: `ls -lh /srv/archivas-snapshots/devnet/snap-*.tar.gz`

2. **Snapshot Age**
   - Should be < 6 hours old
   - Check: `ls -lt /srv/archivas-snapshots/devnet/latest.json`

3. **HTTP Response Time**
   - Should be < 1 second for manifest
   - Check: `time curl -s https://seed2.archivas.ai/devnet/latest.json`

4. **Download Speed**
   - Should complete in < 5 minutes on good connection
   - Check: `time curl -o /tmp/test.tar.gz https://seed2.archivas.ai/devnet/snap-*.tar.gz`

### Alerts to Set Up

Consider setting up alerts for:
- Snapshot age > 8 hours (cron job might have failed)
- Snapshot file size > 1GB or < 100MB (potential issue)
- HTTPS endpoint returning 404 or 500
- Disk space < 10GB on `/srv/archivas-snapshots`

## Security Considerations

### HTTPS Only
- ✅ Snapshots are served over HTTPS
- ✅ Checksums are verified on download
- ✅ No authentication required (public snapshots)

### Checksum Verification
- Every snapshot includes SHA256 checksum
- Bootstrap command verifies before extraction
- Prevents tampering and corruption

### Server Access
- Only Seed2 needs write access to `/srv/archivas-snapshots/`
- Nginx serves read-only
- Cron job runs as `ubuntu` user

## Performance

### Disk Space Requirements

**On Seed2:**
- Current database: ~3GB
- 3 snapshots @ 500MB each: ~1.5GB
- **Total: ~5GB minimum**

**On New Nodes:**
- Snapshot download: ~500MB
- Extracted database: ~3GB
- **Total: ~4GB minimum**

### Network Bandwidth

**Per Bootstrap:**
- Download: ~500MB
- Time: 1-5 minutes (depending on connection)

**Expected Load:**
- Assume 10 bootstraps/day: ~5GB/day
- Monthly: ~150GB/month
- Well within typical server limits

## Rollout Plan

### Phase 1: Soft Launch (Current)
✅ Deploy to production
✅ Test with internal servers
✅ Update documentation
⏳ Monitor for 24-48 hours

### Phase 2: Public Announcement
- Announce in Discord/Telegram
- Update all farmer guides
- Create video tutorial
- Monitor adoption

### Phase 3: Optimization
- Tune snapshot frequency based on usage
- Consider compression algorithm improvements
- Add mainnet support when ready

## Success Criteria

✅ Snapshots generated every 6 hours without errors
✅ Snapshots accessible via HTTPS
✅ Bootstrap completes in < 10 minutes
✅ New nodes start at snapshot height
✅ Chain consensus verified across all nodes
✅ No degradation in Seed2 performance

## Support

For issues or questions:
- GitHub Issues: https://github.com/ArchivasNetwork/archivas/issues
- Discord: [Archivas Discord]
- Email: support@archivas.ai

---

**Last Updated:** 2025-11-15
**Version:** 1.0
**Maintainer:** Archivas Core Team

