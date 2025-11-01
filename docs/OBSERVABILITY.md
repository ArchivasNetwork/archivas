# Archivas Observability Guide

Production-grade monitoring with Prometheus and Grafana.

---

## Quick Start

### 1. Deploy Monitoring Stack

```bash
cd ops/monitoring
docker-compose up -d
```

**Access:**
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/archivas)

### 2. Enable Metrics on Nodes

Node metrics are automatically exposed at `:8080/metrics` (same port as RPC).

**Test:**
```bash
curl http://localhost:8080/metrics
curl http://57.129.148.132:8080/metrics
curl http://72.251.11.191:8080/metrics
```

### 3. Autodiscovery & Health Checks

- `GET /metrics/targets.json` – deterministic list of scrape targets. Supports `?includePeers=true&peerPort=8080` for quick peer fan-out.
- `GET /metrics/health` – JSON health summary driven by internal watchdogs. Returns `status=ok` when all gauges/counters are fresh.

```bash
# List scrape targets (self + peer inventory)
curl "http://localhost:8080/metrics/targets.json?includePeers=true"

# Watchdog health snapshot
curl http://localhost:8080/metrics/health | jq
```

---

## Metrics Reference

### Node Metrics (:8080/metrics)

| Metric | Type | Description |
|--------|------|-------------|
| `archivas_tip_height` | Gauge | Current blockchain height |
| `archivas_peer_count` | Gauge | Connected peers |
| `archivas_peers_known` | Gauge | Total known peers (connected + discovered) |
| `archivas_blocks_total` | Counter | Total blocks processed |
| `archivas_difficulty` | Gauge | Current mining difficulty |
| `archivas_submit_received_total` | Counter | Proof submissions received |
| `archivas_submit_accepted_total` | Counter | Proof submissions accepted |
| `archivas_submit_ignored_total` | Counter | Proof submissions rejected/ignored |
| `archivas_rpc_requests_total{endpoint}` | Counter | RPC requests by endpoint |
| `archivas_metrics_watchdog_triggered{metric}` | Gauge | 1 when the metric watchdog fired |
| `archivas_gossip_msgs_total` | Counter | Peer gossip messages sent |
| `archivas_gossip_addrs_received_total` | Counter | Peer addresses received |

### Timelord Metrics (:9101/metrics)

| Metric | Type | Description |
|--------|------|-------------|
| `archivas_vdf_iterations_total` | Counter | Total VDF iterations computed |
| `archivas_vdf_tick_duration_seconds` | Histogram | Time per VDF tick |
| `archivas_vdf_seed_changed_total` | Counter | VDF seed changes (new blocks) |

### Farmer Metrics (:9102/metrics)

| Metric | Type | Description |
|--------|------|-------------|
| `archivas_plots_loaded` | Gauge | Number of plot files loaded |
| `archivas_qualities_checked_total` | Counter | Total quality checks |
| `archivas_blocks_won_total` | Counter | Blocks successfully mined |
| `archivas_farming_loop_seconds` | Histogram | Time per farming iteration |

---

## Grafana Dashboards

### Pre-loaded Dashboard: "Archivas Network Overview"

**Panels:**
- Tip Height (stat)
- Connected Peers (stat)
- Difficulty (QMAX, stat)
- Blocks Mined (stat)
- Watchdogs Firing (stat)
- Tip Height Over Time (graph)
- Peer Inventory (graph)
- Difficulty Trend (graph)
- Submit Path Throughput (graph)
- RPC Requests by Endpoint (graph)
- Watchdog Status (table)

**Access:** Grafana → Dashboards → Archivas Network Overview

### Creating Custom Dashboards

1. Open Grafana (http://localhost:3000)
2. Click "Create" → "Dashboard"
3. Add Panel
4. Select "Prometheus" datasource
5. Enter metric query (e.g., `archivas_tip_height`)
6. Save dashboard

**Example Queries:**
```promql
# Block production rate
rate(archivas_blocks_total[5m])

# Average peers per node
avg(archivas_peer_count)

# Difficulty change over time
archivas_difficulty

# RPC requests per minute
rate(archivas_rpc_requests_total[1m])

# Metrics watchdog (any firing in last 5m)
max_over_time(archivas_metrics_watchdog_triggered[5m])
```

---

## Prometheus Configuration

### Adding New Nodes

Edit `ops/monitoring/prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'archivas-nodes'
    static_configs:
      - targets:
          - '57.129.148.132:8080'
          - '72.251.11.191:8080'
          - 'new-node-ip:8080'  # Add here
```

Reload:
```bash
docker-compose restart prometheus
```

### Scrape Interval

Default: 5 seconds (adjust in `prometheus.yml` under `global.scrape_interval`)

---

## Alerting (Future)

### Add Alertmanager

Edit `docker-compose.yml` to add:
```yaml
  alertmanager:
    image: prom/alertmanager:latest
    ports: ["9093:9093"]
```

### Example Alerts

**No Peers:**
```yaml
alert: NoPeers
expr: archivas_peer_count == 0
for: 5m
annotations:
  summary: "Node {{$labels.instance}} has no peers"
```

**Chain Stalled:**
```yaml
alert: ChainStalled
expr: increase(archivas_tip_height[5m]) == 0
for: 10m
annotations:
  summary: "Chain not advancing on {{$labels.instance}}"
```

---

## Troubleshooting

### Metrics Not Showing

**Check endpoint:**
```bash
curl http://localhost:8080/metrics | head -20
```

Should see:
```
# HELP archivas_tip_height Current blockchain height
# TYPE archivas_tip_height gauge
archivas_tip_height 1258
```

**Check Prometheus targets:**
- Open http://localhost:9090/targets
- All targets should show "UP"
- If DOWN, check firewall / network connectivity

### Grafana Shows "No Data"

1. Go to Configuration → Data Sources
2. Test Prometheus connection
3. If failed, check `docker-compose logs prometheus`
4. Verify Prometheus is scraping (check /targets)

### Historical Data Missing

Prometheus stores data in `/prometheus` volume. If you restart without volumes, data is lost.

**For persistent storage:**
```yaml
volumes:
  prometheus-data:
    driver: local
    driver_opts:
      type: none
      device: /opt/archivas/prometheus-data
      o: bind
```

---

## JSON Logging

`archivas-node`, `archivas-farmer`, and `archivas-timelord` now emit structured JSONL logs via the standard logger. Each line includes a timestamp (`ts`), component name, log level (derived from prefixes like `[DEBUG]`), and the original message.

Example:

```json
{"ts":"2025-11-01T19:33:12Z","component":"archivas-node","level":"info","msg":"[startup] Archivas node starting"}
```

- Pipe service logs through `jq` for readability: `journalctl -u archivas-node -o cat | jq`.
- Human-friendly `fmt.Printf` banner output remains for CLI UX.

---

## Production Best Practices

### 1. Secure Access

Use reverse proxy (nginx) for authentication:
```nginx
location /grafana/ {
    auth_basic "Archivas Monitoring";
    auth_basic_user_file /etc/nginx/.htpasswd;
    proxy_pass http://localhost:3000/;
}
```

### 2. Retention

Configure Prometheus retention:
```yaml
command:
  - '--storage.tsdb.retention.time=30d'
  - '--storage.tsdb.retention.size=10GB'
```

### 3. Remote Storage

For long-term storage, use remote write (e.g., Thanos, Cortex).

---

## Monitoring Your Testnet

### Key Metrics to Watch

**Health:**
- `archivas_peer_count` > 0 (not isolated)
- `archivas_tip_height` increasing (~1 block/20s)

**Performance:**
- `archivas_vdf_iterations_total` steady rate
- `archivas_farming_loop_seconds` < 2s
- `archivas_rpc_requests_total` reasonable rate

**Network:**
- `archivas_peers_known` growing (gossip working)
- All nodes at same `archivas_tip_height` (synced)

---

**Your blockchain is now production-observable!** 📊

For support: https://github.com/ArchivasNetwork/archivas/issues

