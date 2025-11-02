# Prometheus Observability Setup

v1.1.1: Complete Prometheus metrics setup for Archivas network.

## Overview

All Archivas services now expose Prometheus metrics on dedicated ports:
- **Nodes**: Port `8080` (metrics on `/metrics` endpoint, same as RPC)
- **Timelords**: Port `9101` (dedicated metrics server)
- **Farmers**: Port `9102` (dedicated metrics server)

## Quick Setup

### 1. Deploy Updated Binaries

On each server (57.129.148.132 and 72.251.11.191):

```bash
cd ~/archivas
git pull origin main
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-farmer ./cmd/archivas-farmer
```

### 2. Configure Firewall

On each target host:

```bash
cd ~/archivas/deploy/prometheus
bash firewall-rules.sh
```

Or manually:

```bash
PROM_IP="57.129.148.132"  # Prometheus server IP
sudo ufw allow from $PROM_IP to any port 8080 proto tcp
sudo ufw allow from $PROM_IP to any port 9101 proto tcp
sudo ufw allow from $PROM_IP to any port 9102 proto tcp
sudo ufw reload
```

### 3. Restart Services

**If using systemd:**

```bash
# Update systemd units
sudo cp contrib/systemd/*.service /etc/systemd/system/
sudo systemctl daemon-reload

# Restart services
sudo systemctl restart archivas-node archivas-timelord archivas-farmer
```

**If running manually:**

```bash
# Stop old processes
pkill -f archivas-node
pkill -f archivas-timelord
pkill -f archivas-farmer

# Start node (with 0.0.0.0 binding)
nohup ./archivas-node \
  --rpc 0.0.0.0:8080 \
  --p2p :9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  > logs/node.log 2>&1 &

# Start timelord
nohup ./archivas-timelord \
  --node http://localhost:8080 \
  --step 500 \
  --metrics-addr 0.0.0.0:9101 \
  > logs/timelord.log 2>&1 &

# Start farmer
nohup ./archivas-farmer farm \
  --plots ~/archivas-plots \
  --node http://localhost:8080 \
  --farmer-privkey <YOUR_PRIVKEY> \
  --metrics-addr 0.0.0.0:9102 \
  > logs/farmer.log 2>&1 &
```

### 4. Update Prometheus Config

On Prometheus host (57.129.148.132):

```bash
# Copy config
sudo cp deploy/prometheus/prometheus.yml /etc/prometheus/prometheus.yml

# Restart Prometheus
sudo systemctl restart prometheus
```

### 5. Verify

Run health check from Prometheus host:

```bash
cd ~/archivas/deploy/prometheus
bash health-check.sh
```

Or manually:

```bash
# Check metrics endpoints
curl -s http://57.129.148.132:8080/metrics | head
curl -s http://57.129.148.132:9101/metrics | head
curl -s http://57.129.148.132:9102/metrics | head

# Check health endpoints
curl -s http://57.129.148.132:8080/healthz
curl -s http://57.129.148.132:9101/healthz
curl -s http://57.129.148.132:9102/healthz
```

## Expected Targets

After setup, Prometheus should see:

- ✅ `archivas-nodes`: 2/2 UP
  - `57.129.148.132:8080`
  - `72.251.11.191:8080`

- ✅ `archivas-timelords`: 2/2 UP
  - `57.129.148.132:9101`
  - `72.251.11.191:9101`

- ✅ `archivas-farmers`: 1/1 UP
  - `57.129.148.132:9102`

## Troubleshooting

### Connection Refused

1. **Service not running**: `ps aux | grep archivas`
2. **Wrong bind address**: Ensure `0.0.0.0:PORT` in command/flags
3. **Firewall blocking**: `sudo ufw status | grep PORT`
4. **Port already in use**: `ss -ltnp | grep :9101`

### Metrics Not Updating

1. **Watchdog check**: Metrics should update every 2s (node) or on events (timelord/farmer)
2. **Verify endpoint**: `curl http://IP:PORT/metrics`
3. **Check logs**: `journalctl -u archivas-node -f` or `tail -f logs/node.log`

## Architecture

```
┌─────────────────┐
│   Prometheus    │ (57.129.148.132:9091)
└────────┬────────┘
         │ Scrapes metrics
         │
    ┌────┴────────────────────┐
    │                         │
┌───▼─────────┐      ┌────────▼──┐
│  Node       │      │  Node      │
│  :8080      │      │  :8080     │
│  /metrics   │      │  /metrics  │
└───┬─────────┘      └────────┬──┘
    │                         │
┌───▼─────────┐      ┌────────▼──┐
│  Timelord   │      │  Timelord  │
│  :9101      │      │  :9101     │
│  /metrics   │      │  /metrics  │
└─────────────┘      └────────────┘

┌─────────────┐
│   Farmer    │
│   :9102     │
│   /metrics  │
└─────────────┘
```

