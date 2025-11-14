# Seed2 Full Node - Operations Runbook

**Purpose**: Seed2 is a dual-role server providing both P2P full node services for farmers and cached RPC relay for web clients.

## Architecture

```
                    ┌─────────────────────┐
                    │   seed2.archivas.ai │
                    └──────────┬──────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
        ┌─────▼─────┐    ┌────▼────┐    ┌─────▼──────┐
        │  Node P2P │    │  Nginx  │    │  Fastify   │
        │  :30303   │    │  :443   │    │  :9090     │
        └───────────┘    └────┬────┘    └────────────┘
              │               │
        Farmers/Peers     Web Clients
```

**Roles**:
1. **Full Node (P2P)**: Participates in blockchain consensus, gossips blocks, serves peers
2. **RPC Relay**: Caches hot endpoints, proxies TX submissions to Seed1

**Ports**:
- `30303` TCP/UDP: P2P peer-to-peer networking (public)
- `8082` TCP: Node RPC (localhost only, not exposed)
- `443` TCP: HTTPS RPC relay (public, Nginx)
- `9090` TCP: Relay health/metrics service (localhost only)
- `9102` TCP: Node Prometheus metrics (restricted)

---

## Installation & Bootstrap

### Prerequisites

```bash
# Install dependencies
sudo apt-get update
sudo apt-get install -y rsync curl jq

# Create archivas user (optional, currently using root)
# sudo useradd -r -s /bin/false archivas
```

### Bootstrap from Seed1

**Option A: Database Sync (recommended for initial setup)**

```bash
# Run the bootstrap script
sudo bash /root/archivas/data/bootstrap.sh

# This will:
# 1. Stop the Seed2 node
# 2. Rsync database from Seed1
# 3. Verify data integrity
# 4. Set permissions
```

**Option B: Fresh Sync (slower, starts from genesis or checkpoint)**

```bash
# Create data directory
sudo mkdir -p /var/lib/archivas/seed2
sudo chown root:root /var/lib/archivas/seed2

# Node will sync from Seed1 peer when started
```

### Configuration

1. **Get checkpoint from Seed1**:

```bash
curl -s https://seed.archivas.ai:8081/chainTip | jq
```

2. **Create environment file**:

```bash
sudo mkdir -p /etc/archivas
sudo nano /etc/archivas/seed2-node.env
```

Paste (update values):

```bash
CHECKPOINT_HEIGHT=700000
CHECKPOINT_HASH=your_checkpoint_hash_here
SEED1_P2P=seed.archivas.ai:30303
```

3. **Install systemd unit**:

```bash
sudo cp /root/archivas/services/node-seed2/archivas-node-seed2.service \
       /etc/systemd/system/

sudo systemctl daemon-reload
sudo systemctl enable archivas-node-seed2
```

4. **Configure firewall**:

```bash
sudo bash /root/archivas/infra/firewall.seed2.sh
```

### Start the Node

```bash
# Start service
sudo systemctl start archivas-node-seed2

# Check status
sudo systemctl status archivas-node-seed2

# Monitor logs
sudo journalctl -u archivas-node-seed2 -f --no-pager
```

---

## Operations

### Start/Stop/Restart

```bash
# Start
sudo systemctl start archivas-node-seed2

# Stop (graceful shutdown)
sudo systemctl stop archivas-node-seed2

# Restart
sudo systemctl restart archivas-node-seed2

# Reload systemd config
sudo systemctl daemon-reload
```

### Health Checks

**Node RPC** (localhost):

```bash
# Chain tip
curl -s http://127.0.0.1:8082/chainTip | jq

# Expected: {"height": 700XXX, "hash": "...", "timestamp": ...}
```

**Prometheus Metrics** (if exposed):

```bash
curl -s http://127.0.0.1:9102/metrics | grep archivas
```

**P2P Connectivity**:

```bash
# Check if P2P port is listening
sudo ss -tulpn | grep 30303

# From another machine, test connectivity
telnet seed2.archivas.ai 30303
```

**Peer Count** (from logs):

```bash
sudo journalctl -u archivas-node-seed2 -n 100 --no-pager | grep -i peer
```

### Monitoring Metrics

**Key metrics to watch**:

| Metric | Target | Alert If |
|--------|--------|----------|
| Chain Height | Advancing steadily | Stalled > 60s |
| Peer Count | 20-100 | < 5 for 10m |
| CPU Usage | < 70% | > 85% for 10m |
| Memory | < 12GB | > 14GB |
| Disk I/O Wait | < 20% | > 50% sustained |
| Disk Space | > 20% free | < 10% free |

**Health check commands**:

```bash
# Height check
HEIGHT=$(curl -s http://127.0.0.1:8082/chainTip | jq -r .height)
echo "Current height: $HEIGHT"

# Process status
ps aux | grep archivas-node | grep -v grep

# Resource usage
systemd-cgtop -n 1 | grep archivas-node-seed2

# Disk usage
df -h /var/lib/archivas/seed2
du -sh /var/lib/archivas/seed2
```

---

## Upgrades

### Binary Upgrade

```bash
# 1. Drain peers (reduce max-peers to discourage new connections)
# Edit /etc/systemd/system/archivas-node-seed2.service
# Change: --max-peers 100 → --max-peers 5

sudo systemctl daemon-reload
sudo systemctl restart archivas-node-seed2

# 2. Wait 2-5 minutes for peers to drop

# 3. Stop node
sudo systemctl stop archivas-node-seed2

# 4. Backup current binary
sudo cp /root/archivas/archivas-node /root/archivas/archivas-node.bak

# 5. Build/download new binary
cd /root/archivas
git pull
go build -o archivas-node ./cmd/archivas-node

# 6. Restore max-peers in service file
# Edit: --max-peers 5 → --max-peers 100
sudo systemctl daemon-reload

# 7. Start node
sudo systemctl start archivas-node-seed2

# 8. Verify
sudo journalctl -u archivas-node-seed2 -f
curl http://127.0.0.1:8082/chainTip
```

### Database Maintenance

**Weekly snapshot backup**:

```bash
# Create snapshot directory
sudo mkdir -p /backup/archivas

# Stop node
sudo systemctl stop archivas-node-seed2

# Create tarball
sudo tar -czf /backup/archivas/seed2-$(date +%Y%m%d).tar.gz \
    -C /var/lib/archivas/seed2 .

# Start node
sudo systemctl start archivas-node-seed2

# Verify backup
ls -lh /backup/archivas/
```

**Restore from backup**:

```bash
sudo systemctl stop archivas-node-seed2
sudo rm -rf /var/lib/archivas/seed2/*
sudo tar -xzf /backup/archivas/seed2-YYYYMMDD.tar.gz \
    -C /var/lib/archivas/seed2
sudo systemctl start archivas-node-seed2
```

---

## Troubleshooting

### Node won't start

**Check logs**:

```bash
sudo journalctl -u archivas-node-seed2 -n 50 --no-pager
```

**Common issues**:

1. **Port already in use**:
   ```bash
   sudo ss -tulpn | grep -E '30303|8082'
   # Kill conflicting process
   ```

2. **Database corruption**:
   ```bash
   # Check database
   ls -lh /var/lib/archivas/seed2/badger/
   # Restore from backup or re-bootstrap
   ```

3. **Permission denied**:
   ```bash
   sudo chown -R root:root /var/lib/archivas/seed2
   sudo chmod -R 755 /var/lib/archivas/seed2
   ```

### Sync stalled

**Symptoms**: Height not advancing for > 5 minutes

**Diagnosis**:

```bash
# Check peer count
sudo journalctl -u archivas-node-seed2 -n 100 | grep -i peer

# Check if Seed1 is reachable
telnet seed.archivas.ai 30303

# Check logs for errors
sudo journalctl -u archivas-node-seed2 -f | grep -i error
```

**Fix**:

```bash
# Restart with fresh peer connection
sudo systemctl restart archivas-node-seed2

# If persistent, re-bootstrap from Seed1
sudo bash /root/archivas/data/bootstrap.sh
```

### High memory usage

**Symptoms**: Memory > 14GB, OOM killer triggered

**Fix**:

```bash
# Lower GOGC (more aggressive GC)
# Edit /etc/systemd/system/archivas-node-seed2.service
# Change: Environment=GOGC=50 → Environment=GOGC=30

sudo systemctl daemon-reload
sudo systemctl restart archivas-node-seed2
```

### Peer count low

**Target**: 20-100 peers

**If < 5 peers**:

```bash
# Check firewall
sudo ufw status | grep 30303

# Check if port is publicly accessible
# From external machine:
telnet seed2.archivas.ai 30303

# Check Seed1 connectivity
telnet seed.archivas.ai 30303

# Restart node
sudo systemctl restart archivas-node-seed2
```

---

## Emergency Procedures

### Complete database reset

```bash
# 1. Stop node
sudo systemctl stop archivas-node-seed2

# 2. Backup current data (just in case)
sudo mv /var/lib/archivas/seed2 /var/lib/archivas/seed2.old

# 3. Re-bootstrap from Seed1
sudo bash /root/archivas/data/bootstrap.sh

# 4. Start node
sudo systemctl start archivas-node-seed2

# 5. Monitor sync
sudo journalctl -u archivas-node-seed2 -f
```

### Disable node temporarily (keep relay)

```bash
# Stop P2P node
sudo systemctl stop archivas-node-seed2

# Relay continues to serve from Seed1
# Update Nginx to remove seed2_node from read_pool if needed
```

---

## Performance Tuning

### For high-load scenarios

**Increase concurrent connections**:

```bash
# Edit /etc/systemd/system/archivas-node-seed2.service
# Change: --max-peers 100 → --max-peers 150

sudo systemctl daemon-reload
sudo systemctl restart archivas-node-seed2
```

**Optimize disk I/O**:

```bash
# Ensure data dir is on NVMe/SSD
mount | grep archivas

# Check I/O scheduler
cat /sys/block/nvme0n1/queue/scheduler
# Should be: [none] or [mq-deadline] for NVMe

# Set if needed:
echo none | sudo tee /sys/block/nvme0n1/queue/scheduler
```

**Resource limits**:

Current limits in systemd unit:
- MemoryMax: 16G
- CPUQuota: 400% (4 cores)
- TasksMax: 8192
- LimitNOFILE: 65536

Adjust as needed based on server specs.

---

## Farmer Integration

Farmers should connect to both Seed1 and Seed2 for redundancy:

```bash
archivas-node \
  --peer seed.archivas.ai:30303 \
  --peer seed2.archivas.ai:30303 \
  --no-peer-discovery \
  --checkpoint-height <HEIGHT> \
  --checkpoint-hash <HASH>
```

**Note**: Farmers should NOT use the HTTPS RPC relay for P2P - they must use the P2P port `:30303`.

---

## Alerting Rules

Recommended Prometheus alerts:

```yaml
- alert: Seed2HeightStalled
  expr: rate(archivas_chain_height[5m]) == 0
  for: 2m
  annotations:
    summary: Seed2 chain height stalled

- alert: Seed2LowPeers
  expr: archivas_peer_count < 5
  for: 10m
  annotations:
    summary: Seed2 peer count critically low

- alert: Seed2HighCPU
  expr: rate(process_cpu_seconds_total[5m]) > 3.5
  for: 10m
  annotations:
    summary: Seed2 CPU usage > 85%

- alert: Seed2LowDisk
  expr: node_filesystem_avail_bytes{mountpoint="/var/lib/archivas"} / node_filesystem_size_bytes < 0.10
  for: 5m
  annotations:
    summary: Seed2 disk < 10% free
```

---

## Logs

**View recent logs**:

```bash
sudo journalctl -u archivas-node-seed2 -n 100 --no-pager
```

**Follow logs**:

```bash
sudo journalctl -u archivas-node-seed2 -f
```

**Filter errors**:

```bash
sudo journalctl -u archivas-node-seed2 -p err -n 50
```

**Export logs**:

```bash
sudo journalctl -u archivas-node-seed2 --since "1 hour ago" > /tmp/seed2-logs.txt
```

---

## Quick Reference

| Task | Command |
|------|---------|
| Start node | `sudo systemctl start archivas-node-seed2` |
| Stop node | `sudo systemctl stop archivas-node-seed2` |
| Status | `sudo systemctl status archivas-node-seed2` |
| Logs | `sudo journalctl -u archivas-node-seed2 -f` |
| Chain height | `curl -s http://127.0.0.1:8082/chainTip \| jq .height` |
| Peer count | `sudo journalctl -u archivas-node-seed2 \| grep -i peer \| tail` |
| Disk usage | `df -h /var/lib/archivas/seed2` |
| Open ports | `sudo ss -tulpn \| grep -E '30303\|8082'` |

---

## Support

- **Documentation**: https://docs.archivas.ai
- **Logs location**: `/var/log/journal/` (via journalctl)
- **Data directory**: `/var/lib/archivas/seed2`
- **Config**: `/etc/archivas/seed2-node.env`

