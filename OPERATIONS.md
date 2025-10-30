# Archivas Operations Guide

## Quick Reference

**Ports:**
- 8080: HTTP RPC
- 9090: P2P networking

**Paths:**
- `/opt/archivas/data` - Database
- `/opt/archivas/plots` - Plot files
- `/var/log/archivas/` - Logs
- `/opt/archivas/peers.json` - Peer persistence

**Services:**
- `archivas-node.service`
- `archivas-timelord.service`
- `archivas-farmer.service`

---

## Installation (Production)

### 1. Create User

```bash
sudo useradd -r -s /bin/false archivas
sudo mkdir -p /opt/archivas/{data,plots,logs}
sudo mkdir -p /var/log/archivas
sudo chown -R archivas:archivas /opt/archivas /var/log/archivas
```

### 2. Install Binaries

```bash
cd /tmp
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-farmer ./cmd/archivas-farmer

sudo cp archivas-* /opt/archivas/
sudo cp -r genesis /opt/archivas/
sudo chown archivas:archivas /opt/archivas/archivas-*
```

### 3. Install Systemd Units

```bash
sudo cp contrib/systemd/*.service /etc/systemd/system/
sudo systemctl daemon-reload
```

### 4. Configure Farmer (if farming)

```bash
# Create environment file
sudo nano /etc/archivas/farmer.env
```

Add:
```
FARMER_PRIVKEY=your_private_key_here
```

```bash
sudo chmod 600 /etc/archivas/farmer.env
sudo chown archivas:archivas /etc/archivas/farmer.env
```

### 5. Start Services

```bash
# Node (required)
sudo systemctl enable archivas-node
sudo systemctl start archivas-node

# Timelord (optional but recommended)
sudo systemctl enable archivas-timelord
sudo systemctl start archivas-timelord

# Farmer (if you have plots)
sudo systemctl enable archivas-farmer
sudo systemctl start archivas-farmer
```

---

## Docker Deployment

### Quick Start

```bash
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Start node + timelord
docker-compose up -d

# Check status
curl http://localhost:8080/healthz

# View logs
docker-compose logs -f node
```

### With Farming

```bash
# Set environment variable
export FARMER_PRIVKEY=your_private_key

# Uncomment farmer service in docker-compose.yml

# Start all services
docker-compose up -d

# Monitor
docker-compose logs -f
```

---

## Monitoring

### Health Check

```bash
curl http://localhost:8080/healthz
```

Response:
```json
{
  "ok": true,
  "height": 892,
  "difficulty": 1125899906842624,
  "peers": 1
}
```

### Chain Status

```bash
curl http://localhost:8080/chainTip
curl http://localhost:8080/genesisHash
curl http://localhost:8080/challenge
```

### Systemd Status

```bash
sudo systemctl status archivas-node
sudo systemctl status archivas-timelord
sudo systemctl status archivas-farmer
```

### View Logs

```bash
# Systemd journals
sudo journalctl -u archivas-node -f
sudo journalctl -u archivas-timelord -f
sudo journalctl -u archivas-farmer -f

# Log files
tail -f /var/log/archivas/node.log
tail -f /var/log/archivas/timelord.log
tail -f /var/log/archivas/farmer.log
```

### Docker Logs

```bash
docker-compose logs -f node
docker-compose logs -f timelord
docker-compose logs --tail=100 node
```

---

## Maintenance

### Update Binary

```bash
# Stop services
sudo systemctl stop archivas-node archivas-timelord archivas-farmer

# Backup
sudo cp /opt/archivas/archivas-node /opt/archivas/archivas-node.bak

# Update
cd /tmp/archivas
git pull origin main
go build -o archivas-node ./cmd/archivas-node
sudo cp archivas-node /opt/archivas/
sudo chown archivas:archivas /opt/archivas/archivas-node

# Restart
sudo systemctl start archivas-node archivas-timelord archivas-farmer
```

### Docker Update

```bash
docker-compose pull
docker-compose up -d
```

### Backup Database

```bash
# Stop node first!
sudo systemctl stop archivas-node

# Backup
sudo tar -czf archivas-backup-$(date +%Y%m%d).tar.gz /opt/archivas/data

# Restart
sudo systemctl start archivas-node
```

---

## Firewall Configuration

### UFW (Ubuntu)

```bash
# Allow RPC (localhost only recommended)
sudo ufw allow from 127.0.0.1 to any port 8080

# Allow P2P (public)
sudo ufw allow 9090/tcp

# Enable
sudo ufw enable
```

### Firewalld (CentOS/RHEL)

```bash
sudo firewall-cmd --permanent --add-port=9090/tcp
sudo firewall-cmd --reload
```

---

## Troubleshooting

### Node won't start

```bash
# Check logs
sudo journalctl -u archivas-node -n 50

# Common issues:
# - Port 8080/9090 in use
# - Database corruption
# - Missing genesis file
```

### Peer connection issues

```bash
# Check P2P port
sudo netstat -tlnp | grep 9090

# Test connectivity
nc -zv 57.129.148.132 9090

# Check peers
curl http://localhost:8080/healthz | jq .peers
```

### Sync stalled

```bash
# Check peer status
grep "peer.*status" /var/log/archivas/node.log | tail -5

# Restart to reconnect
sudo systemctl restart archivas-node
```

### High CPU usage

```bash
# Check what's consuming
top -u archivas

# Farmer scanning is CPU-intensive (normal)
# Timelord VDF computation (normal)
```

---

## Performance Tuning

### Increase File Limits

```bash
# Edit /etc/security/limits.conf
archivas soft nofile 65535
archivas hard nofile 65535
```

### Optimize Database

```bash
# BadgerDB compaction (automatic)
# Logs will show periodic compactions
grep "compaction" /var/log/archivas/node.log
```

---

## Security Best Practices

### Minimize Attack Surface

```bash
# Only expose P2P publicly
# Keep RPC on localhost or behind firewall
sudo ufw deny 8080/tcp  # Block public RPC
sudo ufw allow 9090/tcp  # Allow P2P
```

### Protect Private Keys

```bash
# Farmer private key
sudo chmod 600 /etc/archivas/farmer.env
sudo chown root:root /etc/archivas/farmer.env

# Never commit private keys to version control!
```

### Regular Updates

```bash
# Subscribe to releases
# Update within 24h of security releases
# Test on non-production node first
```

---

## Metrics & Monitoring

### Health Endpoint

**HTTP GET /healthz**

Returns:
- `ok`: boolean
- `height`: current block height
- `difficulty`: current difficulty
- `peers`: connected peer count

**Use for:**
- Load balancer health checks
- Uptime monitoring (UptimeRobot, etc.)
- Alerting (if peers = 0, alert!)

### Prometheus Metrics (Coming Soon)

**HTTP GET /metrics**

Exposes:
- `archivas_tip_height` - Current block height
- `archivas_peer_count` - Connected peers
- `archivas_blocks_total` - Total blocks mined
- `archivas_vdf_iterations_total` - VDF computation progress

---

## Common Tasks

### Check Sync Status

```bash
curl http://localhost:8080/chainTip
curl http://57.129.148.132:8080/chainTip
# Compare heights
```

### Restart Services Safely

```bash
# Graceful stop
sudo systemctl stop archivas-farmer
sudo systemctl stop archivas-timelord  
sudo systemctl stop archivas-node

# Wait for shutdown
sleep 5

# Start in order
sudo systemctl start archivas-node
sleep 2
sudo systemctl start archivas-timelord
sudo systemctl start archivas-farmer
```

### Reset to Genesis

```bash
sudo systemctl stop archivas-node
sudo rm -rf /opt/archivas/data
sudo systemctl start archivas-node
# Will sync from network
```

---

## Support

**Issues:** https://github.com/ArchivasNetwork/archivas/issues  
**Discussions:** https://github.com/ArchivasNetwork/archivas/discussions  
**Docs:** https://github.com/ArchivasNetwork/archivas/tree/main/docs

---

**Archivas Operations Guide** - Keep your node running! 🌾

