# Deploy v1.1.1 - Devnet V4 (Fresh Chain)

**Date:** 2025-11-02  
**Version:** v1.1.1  
**Network:** archivas-devnet-v4

## CRITICAL: Full Chain Reset Required

**Why:** Server A's 200k-block chain was built with pre-v0.9.0 difficulty scaling (1e15 range). The new v1.1.1 code uses QMAX=1e12 normalized quality. Old blocks **cannot be verified** with new code.

**Solution:** Both servers restart from genesis with v1.1.1 code.

---

## What Changed

### Genesis V4
- **Network ID:** `archivas-devnet-v3` → `archivas-devnet-v4`
- **Seed:** `archivas-devnet-v1-2025-10-29` → `archivas-devnet-v4-2025-11-02`
- **Added:** `protocolVersion`, `difficultyParamsID`, `initialDifficulty`
- **New Genesis Hash:** (will be computed after deployment)

### IBD (Initial Block Download)
- **Batched sync:** 512 blocks per request
- **Disk-based serving:** Blocks loaded from BadgerDB, not memory
- **Backpressure:** Max 2 concurrent IBD streams
- **Metrics:** `archivas_ibd_*` for monitoring sync progress

### Observability
- **Prometheus:** All services expose `/metrics`
- **Ports:** Node (8080), Timelord (9101), Farmer (9102)
- **Grafana:** Dashboard fixed to show block mining rate

---

## Deployment Steps

### 1. Stop All Services

**Server A (57.129.148.132):**
```bash
pkill -f archivas-node
pkill -f archivas-timelord
pkill -f archivas-farmer
```

**Server B (72.251.11.191):**
```bash
pkill -f archivas-node
pkill -f archivas-timelord
```

**Server C (57.129.148.134):**
```bash
pkill -f archivas-farmer
```

### 2. Clear Old Chain Data

**Server A:**
```bash
cd ~/archivas
rm -rf data/*
```

**Server B:**
```bash
cd ~/archivas
rm -rf data/*
```

**Server C:** (Farmer only, no DB to clear)

### 3. Update Code

**All Servers:**
```bash
cd ~/archivas
git pull origin main
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-wallet ./cmd/archivas-wallet
```

### 4. Restart Services

**Server A (57.129.148.132):**
```bash
# Node (primary)
nohup ./archivas-node \
  --rpc 0.0.0.0:8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --bootnodes "72.251.11.191:9090" \
  > logs/node.log 2>&1 &

# Timelord
nohup ./archivas-timelord \
  --node http://localhost:8080 \
  --step 500 \
  --metrics-addr 0.0.0.0:9101 \
  > logs/timelord.log 2>&1 &

# Farmer
nohup ./archivas-farmer farm \
  --plots ./plots-massive \
  --node http://localhost:8080 \
  --farmer-privkey 6301448565bb68053331f80a62621a9ff0a0ac8d6863010ae5bf02237800e5c8 \
  --metrics-addr 0.0.0.0:9102 \
  > logs/farmer.log 2>&1 &
```

**Server B (72.251.11.191):**
```bash
# Node (will sync from A via IBD)
nohup ./archivas-node \
  --rpc 0.0.0.0:8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --bootnodes "57.129.148.132:9090" \
  > logs/node.log 2>&1 &

# Timelord
nohup ./archivas-timelord \
  --node http://localhost:8080 \
  --step 500 \
  --metrics-addr 0.0.0.0:9101 \
  > logs/timelord.log 2>&1 &
```

**Server C (57.129.148.134):**
```bash
# Farmer (connects to Server A)
nohup ./archivas-farmer farm \
  --plots /home/ubuntu/archivas-plots-all \
  --node http://57.129.148.132:8080 \
  --farmer-privkey b8d7a093da6ed086d9a47d7a77dc3f75bc8399d86023fed985261af32f2c185e \
  --metrics-addr 0.0.0.0:9102 \
  > logs/farmer.log 2>&1 &
```

### 5. Verify

**Check Genesis Hash (should match on all nodes):**
```bash
# Server A:
curl -s http://localhost:8080/genesisHash

# Server B:
curl -s http://72.251.11.191:8080/genesisHash
```

**Watch Sync:**
```bash
# Server B should sync from Server A via IBD:
tail -f logs/node.log | grep -E "\[sync\]|received batch|applied"
```

**Expected output on Server B:**
```
[sync] starting IBD from height 1 via peer 57.129.148.132:9090
[p2p] received batch from=1 count=512 tip=1000
[sync] applied 512 blocks from batch, now at height 512
[sync] requesting next batch from=513
...
```

**Check Prometheus:**
```bash
curl -s http://57.129.148.132:8080/metrics | grep archivas_tip_height
curl -s http://72.251.11.191:8080/metrics | grep archivas_tip_height
```

Both should converge to the same height.

---

## Post-Deployment

### Update Prometheus Config

The network ID changed to `archivas-devnet-v4`. Update labels if needed:

```bash
# On Prometheus host:
sudo sed -i 's/archivas-devnet-v3/archivas-devnet-v4/g' /home/ubuntu/archivas/ops/monitoring/prometheus.yml
sudo docker restart archivas-prometheus
```

### Verify All Services

```bash
# From Server A:
cd ~/archivas/deploy/prometheus
bash health-check.sh
```

Expected: All 6 targets UP (2 nodes, 2 timelords, 2 farmers).

---

## Expected Results

- **Fresh chain** starting from block 0
- **Both nodes** on same chain, synced via IBD
- **7 k28 plots** farming (1.88B hashes total)
- **Blocks producing** at ~20-30s intervals
- **Grafana showing** live mining activity
- **No more verification errors**

---

## Genesis Hash

After deployment, the new genesis hash will be:
```bash
curl -s http://57.129.148.132:8080/genesisHash
```

Document this hash for future reference.

---

## Rollback

If issues arise, revert to old network:

```bash
git checkout HEAD~1 genesis/devnet.genesis.json
# Use --network-id archivas-devnet-v3
```

But the old 200k-block chain is GONE (data deleted). This is a one-way migration.

