# Deploy Archivas v0.3.0 - Observability

## Quick Deploy (Servers)

### ON BOTH SERVERS:

```bash
cd ~/archivas
git pull origin main
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-explorer ./cmd/archivas-explorer
go build -o archivas-registry ./cmd/archivas-registry

# Restart node (metrics now on :8080/metrics)
pkill -f archivas-node
sleep 3

nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes <OTHER_IP>:9090 \
  --enable-gossip \
  > logs/node-v0.3.0.log 2>&1 &

# Test metrics
sleep 5
curl http://localhost:8080/metrics | head -20
```

### SERVER A ONLY - Start Registry + Monitoring:

```bash
# Start registry
nohup ./archivas-registry \
  --port :8088 \
  --network-id archivas-devnet-v3 \
  > logs/registry.log 2>&1 &

# Start Prometheus + Grafana
cd ops/monitoring
docker-compose up -d
cd ../..

# Verify
curl http://localhost:8088/health
curl http://localhost:9090/-/healthy  # Prometheus
curl http://localhost:3000/api/health  # Grafana
```

## Access

**Metrics:**
- Server A: http://57.129.148.132:8080/metrics
- Server B: http://72.251.11.191:8080/metrics

**Monitoring:**
- Prometheus: http://57.129.148.132:9090
- Grafana: http://57.129.148.132:3000 (admin/archivas)

**Registry:**
- http://57.129.148.132:8088/nodes
- http://57.129.148.132:8088/peers

**Explorer:**
- http://57.129.148.132:8082 (already running)
- New pages: /peers, /mempool

---

## Verification

```bash
# Metrics working
curl http://localhost:8080/metrics | grep archivas_tip_height

# Grafana accessible
curl http://localhost:3000/api/health

# Registry working
curl http://localhost:8088/health
```

**Expected:**
```
archivas_tip_height 1258
{"database":"ok"}
{"ok":true,"activeNodes":0,"totalNodes":0}
```

---

## v0.3.0 Features

✅ Prometheus metrics (10+ metrics)
✅ /metrics endpoint (all binaries)
✅ Grafana dashboards
✅ Node registry
✅ Explorer upgrades
✅ Complete documentation

**Production observability complete!** 📊
