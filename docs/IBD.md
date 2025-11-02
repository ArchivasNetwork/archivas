# Initial Block Download (IBD)

v1.1.1: Efficient batched sync for catching up to network tip

## Overview

When a fresh node joins the network, it needs to download and verify all historical blocks from genesis to the current tip. This process is called Initial Block Download (IBD).

## How It Works

### 1. Detection

When a node receives a NEW_BLOCK announcement and detects it's more than 10 blocks behind:
```
[p2p] We're behind by 200000 blocks (our=0, peer=200000), starting IBD
```

The node automatically initiates batched sync.

### 2. Batched Requests

Instead of requesting blocks one-by-one (slow), the node requests blocks in batches:

**Client (syncing node) sends:**
```json
{
  "fromHeight": 1,
  "maxBlocks": 512
}
```

**Server (peer) responds:**
```json
{
  "fromHeight": 1,
  "count": 512,
  "blocks": [...],
  "tipHeight": 200000,
  "eof": false
}
```

### 3. Sequential Application

The node applies blocks sequentially:
- Verifies each block (PoSpace proof, difficulty, challenge)
- Updates world state (balances, nonces)
- Persists to disk
- Requests next batch

### 4. Progress Tracking

Every 100 blocks, the node logs progress:
```
[sync] progress: applied 100 blocks, now at height 100 (tip: 200000)
[sync] progress: applied 200 blocks, now at height 200 (tip: 200000)
...
```

### 5. Completion

When caught up:
```
[sync] caught up to tip at height 200000
```

## Configuration

### Batch Size

Default: 512 blocks per batch

Tunable via flag (future):
```bash
--ibd-batch-size 1024  # Larger batches = faster sync, more memory
```

### Backpressure

Max concurrent IBD streams per node: 2

If a node is already serving 2 IBD streams, it rejects new requests with an empty batch (client retries later).

## Performance

### Typical Sync Speed

- **Local network:** ~5,000 blocks/sec
- **WAN (cross-region):** ~1,000-2,000 blocks/sec
- **200k blocks:** ~2-4 minutes on fast network

### Rate Limiting

- Server adds 5ms delay between batches
- Max throughput: ~100 batches/sec = ~50,000 blocks/sec

This prevents IBD from starving RPC endpoints.

## Metrics

Monitor IBD progress via Prometheus:

```promql
# Batches requested
rate(archivas_ibd_requested_batches_total[1m])

# Batches received
rate(archivas_ibd_received_batches_total[1m])

# Blocks applied per second
rate(archivas_ibd_blocks_applied_total[1m])

# Current IBD streams
archivas_ibd_inflight
```

## Troubleshooting

### Node stuck at height 0

**Symptoms:**
```
[sync] requesting next batch from=1
[sync] requesting next batch from=1
[sync] requesting next batch from=1
```

**Causes:**
1. Peer not responding to REQUEST_BLOCKS
2. Blocks failing verification
3. Network connectivity issues

**Fix:**
```bash
# Check peer connection:
curl http://localhost:8080/peers

# Check for errors in logs:
grep -i "failed to apply\|error" logs/node.log

# Restart node to retry:
pkill archivas-node && nohup ./archivas-node ...
```

### Slow sync (<100 blocks/sec)

**Causes:**
1. Disk I/O bottleneck
2. CPU bottleneck (block verification)
3. Network latency

**Fix:**
- Increase batch size (if memory allows)
- Use SSD for block storage
- Connect to geographically closer peer

### Watchdog triggered

```
[metrics] watchdog triggered: metric=archivas_ibd_blocks_applied_total stale=5m
```

**Meaning:** No blocks have been applied in 5 minutes - sync is stalled.

**Fix:** Restart node to re-initiate IBD.

## Architecture

```
Fresh Node (height 0)              Peer Node (height 200k)
       │                                    │
       │ RequestBlocks(from=1, max=512)    │
       ├───────────────────────────────────>│
       │                                    │
       │                                    │ Load blocks 1-512 from disk
       │                                    │ (memory + BadgerDB)
       │                                    │
       │ BlocksBatch(512 blocks, tip=200k) │
       │<───────────────────────────────────┤
       │                                    │
  Apply 512 blocks                         │
  Persist to disk                          │
       │                                    │
       │ RequestBlocks(from=513, max=512)  │
       ├───────────────────────────────────>│
       │                                    │
      ...                                  ...
       │                                    │
       │ BlocksBatch(EOF=true)             │
       │<───────────────────────────────────┤
       │                                    │
   [sync] caught up!                       │
```

## Security

- **Genesis hash validation:** Verified during handshake
- **Block verification:** Every block's PoSpace proof, difficulty, and challenge are verified
- **Sequential application:** Blocks must be applied in order (height n+1 after height n)
- **No trust required:** Peer cannot forge blocks or skip verification

## Files

- `p2p/ibd.go` - Request/response handlers
- `p2p/ibd_sync.go` - Sync initiation logic
- `storage/blockchain.go` - Disk-based block serving
- `metrics/ibd_metrics.go` - IBD metrics
- `cmd/archivas-node/main.go` - OnBlocksRangeRequest implementation

