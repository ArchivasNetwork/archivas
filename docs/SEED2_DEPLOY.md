# Seed2 RPC Relay Deployment Guide

**Version**: 1.0  
**Target**: seed2.archivas.ai  
**Purpose**: Production-ready Nginx reverse proxy with micro-caching and failover

---

## 📋 Overview

Seed2 acts as an **RPC relay proxy** that forwards requests to Seed1 (the primary node) with:

- ✅ **Micro-caching** (1-2 seconds) to crush thundering herds
- ✅ **Rate limiting** (30 req/s per IP, bursts allowed)
- ✅ **Circuit breaker** and automatic failover
- ✅ **Connection pooling** and keepalive optimization
- ✅ **CORS headers** for web applications
- ✅ **TLS/HTTPS** with auto-renewal support

---

## 🚀 Deployment Steps

### 1. Prerequisites

```bash
# On Server D (seed2.archivas.ai)
# Ensure you have root/sudo access
sudo -i

# Update system packages
apt-get update && apt-get upgrade -y

# Install required packages
apt-get install -y nginx certbot python3-certbot-nginx curl jq
```

### 2. Clone Repository

```bash
# Clone the archivas repository
cd /root
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Or pull latest changes if already cloned
git pull origin main
```

### 3. Run Deployment Script

```bash
# Make the deployment script executable
chmod +x deploy/seed2/deploy-relay.sh

# Run the deployment
sudo ./deploy/seed2/deploy-relay.sh
```

**What the script does:**
- Installs Nginx and Certbot
- Creates cache directories
- Generates self-signed SSL certificate (temporary)
- Copies Nginx configuration
- Enables the site
- Reloads Nginx

### 4. Configure DNS

Point your DNS A record to Seed2's IP address:

```bash
# Get your server's public IP
curl -s ifconfig.me
```

**DNS Configuration:**
```
Type: A
Name: seed2.archivas.ai
Value: <SERVER_IP>
TTL: 300 (5 minutes)
```

Wait for DNS propagation (check with `dig seed2.archivas.ai`).

### 5. Obtain Real SSL Certificate

After DNS is configured:

```bash
# Get Let's Encrypt certificate
sudo certbot --nginx -d seed2.archivas.ai

# Choose options:
# - Enter email for urgent renewal notices
# - Agree to Terms of Service
# - Redirect HTTP to HTTPS: Yes
```

### 6. Enable Auto-Reload on Cert Renewal

```bash
# Copy systemd units
sudo cp deploy/seed2/nginx-reload.service /etc/systemd/system/
sudo cp deploy/seed2/nginx-reload.path /etc/systemd/system/

# Enable and start path watcher
sudo systemctl enable nginx-reload.path
sudo systemctl start nginx-reload.path

# Verify
sudo systemctl status nginx-reload.path
```

---

## ✅ Smoke Tests

### Test 1: Basic Connectivity

```bash
# Test HTTPS endpoint
curl -I https://seed2.archivas.ai/chainTip

# Expected:
# HTTP/1.1 200 OK
# X-Cache-Status: MISS (first request)
# X-Relay-Node: seed2.archivas.ai
# Access-Control-Allow-Origin: *
```

### Test 2: Caching Behavior

```bash
# First request (MISS)
echo "Request 1:"
curl -I https://seed2.archivas.ai/chainTip | grep X-Cache-Status

# Second request within 1s (HIT)
echo "Request 2:"
curl -I https://seed2.archivas.ai/chainTip | grep X-Cache-Status

# Expected output:
# Request 1: X-Cache-Status: MISS
# Request 2: X-Cache-Status: HIT
```

### Test 3: Rate Limiting

```bash
# Test rate limit (30 req/s, burst 60)
for i in {1..100}; do
  curl -s -o /dev/null -w "%{http_code}\n" https://seed2.archivas.ai/chainTip &
done | sort | uniq -c

# Expected: Some 200s, some 429s (rate limited)
```

### Test 4: CORS Headers

```bash
# Test CORS preflight
curl -X OPTIONS https://seed2.archivas.ai/chainTip -v

# Expected:
# HTTP/1.1 204 No Content
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Methods: GET, POST, OPTIONS
```

### Test 5: Recent Blocks

```bash
# Test another hot endpoint
curl -s https://seed2.archivas.ai/blocks/recent?limit=5 | jq .

# Expected: JSON array of recent blocks
```

### Test 6: Account Balance

```bash
# Test account endpoint (should be cached for 2s)
curl -s https://seed2.archivas.ai/account/arcv1YOUR_ADDRESS_HERE | jq .

# Expected: {"balance":"...", "nonce":"..."}
```

### Test 7: Relay Status

```bash
# Check relay health endpoint
curl -s https://seed2.archivas.ai/relay/status | jq .

# Expected:
# {
#   "relay": "seed2.archivas.ai",
#   "backend": "seed1",
#   "cache": "enabled",
#   "status": "healthy"
# }
```

---

## 📊 Monitoring

### Check Nginx Logs

```bash
# Access logs (real-time)
sudo tail -f /var/log/nginx/seed2-access.log

# Error logs
sudo tail -f /var/log/nginx/seed2-error.log

# Filter for cache hits/misses
sudo tail -f /var/log/nginx/seed2-access.log | grep -E "HIT|MISS"

# Filter for rate limits
sudo tail -f /var/log/nginx/seed2-access.log | grep " 429 "
```

### Cache Statistics

```bash
# Check cache size
sudo du -sh /var/cache/nginx/archivas

# List cached files
sudo find /var/cache/nginx/archivas -type f | wc -l

# Clear cache (if needed)
sudo rm -rf /var/cache/nginx/archivas/*
sudo systemctl reload nginx
```

### Nginx Status

```bash
# Check Nginx status
sudo systemctl status nginx

# Test configuration
sudo nginx -t

# Reload configuration
sudo systemctl reload nginx

# Restart (if needed)
sudo systemctl restart nginx
```

### Connection Statistics

```bash
# Active connections
sudo netstat -an | grep :443 | wc -l

# Established connections to backend (Seed1)
sudo netstat -an | grep 57.129.148.132:8080 | grep ESTABLISHED

# Connection states
sudo ss -s
```

---

## 🔧 Configuration Tuning

### Adjust Cache Duration

Edit `/etc/nginx/sites-available/archivas-seed2`:

```nginx
# Hot endpoints (default: 1s)
proxy_cache_valid 200 1s;

# Account lookups (default: 2s)
proxy_cache_valid 200 2s;

# Increase if data changes less frequently
# Decrease if you need real-time updates
```

### Adjust Rate Limits

```nginx
# Per IP (default: 30 req/s, burst 60)
limit_req_zone $binary_remote_addr zone=req_ip:10m rate=30r/s;

# In location blocks:
limit_req zone=req_ip burst=60 nodelay;

# Increase if legitimate traffic is being limited
# Decrease for stricter protection
```

### Adjust Timeouts

```nginx
# Connection timeout (default: 2s)
proxy_connect_timeout 2s;

# Read timeout (default: 20s)
proxy_read_timeout 20s;

# Increase if backend is slow
# Decrease to fail fast
```

### Add More Backends

To add Seed3 or other backends:

```nginx
upstream archivas_backend {
    server 57.129.148.132:8080 max_fails=3 fail_timeout=5s;
    server SEED3_IP:8080 max_fails=3 fail_timeout=5s;
    
    keepalive 32;
}
```

---

## 🛡️ Security Hardening

### Firewall Rules

```bash
# Allow HTTPS
sudo ufw allow 443/tcp

# Allow SSH (if not already allowed)
sudo ufw allow 22/tcp

# Enable firewall
sudo ufw enable
```

### Block Abusive IPs

```bash
# Block specific IP
sudo iptables -A INPUT -s ABUSIVE_IP -j DROP

# Or in Nginx config:
# deny ABUSIVE_IP;

# Save iptables rules
sudo iptables-save > /etc/iptables/rules.v4
```

### Enable Fail2Ban

```bash
# Install fail2ban
sudo apt-get install -y fail2ban

# Create Nginx jail
sudo tee /etc/fail2ban/jail.d/nginx-limit.conf <<EOF
[nginx-limit]
enabled = true
filter = nginx-limit
logpath = /var/log/nginx/seed2-access.log
maxretry = 100
findtime = 60
bantime = 3600
EOF

# Create filter
sudo tee /etc/fail2ban/filter.d/nginx-limit.conf <<EOF
[Definition]
failregex = ^<HOST> .* " 429 .*$
ignoreregex =
EOF

# Restart fail2ban
sudo systemctl restart fail2ban

# Check banned IPs
sudo fail2ban-client status nginx-limit
```

---

## 🌐 Multi-RPC Explorer Setup

Update the explorer to use multiple RPC endpoints:

```bash
# Build the explorer with multi-RPC support
cd /root/archivas
go build -o archivas-explorer cmd/archivas-explorer/main.go

# Run with multiple RPCs (comma-separated)
./archivas-explorer \
  -node "https://seed.archivas.ai,https://seed2.archivas.ai,https://seed3.archivas.ai" \
  -port ":8082" \
  -timeout 3 \
  -retries 2

# Or as a systemd service:
sudo tee /etc/systemd/system/archivas-explorer.service <<EOF
[Unit]
Description=Archivas Block Explorer
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/archivas
ExecStart=/root/archivas/archivas-explorer \
  -node "https://seed.archivas.ai,https://seed2.archivas.ai" \
  -port ":8082" \
  -timeout 3 \
  -retries 2
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable archivas-explorer
sudo systemctl start archivas-explorer
```

### Test Multi-RPC Failover

```bash
# Check RPC status endpoint
curl -s http://localhost:8082/rpc/status | jq .

# Expected output:
# {
#   "rpcs": [
#     {
#       "rpc": "https://seed.archivas.ai",
#       "healthy": true,
#       "failures": 0,
#       "last_check": "2025-11-13T20:30:00Z"
#     },
#     {
#       "rpc": "https://seed2.archivas.ai",
#       "healthy": true,
#       "failures": 0,
#       "last_check": "2025-11-13T20:30:00Z"
#     }
#   ],
#   "total": 2,
#   "healthy": 2,
#   "current_rpc": "https://seed.archivas.ai"
# }
```

---

## 📈 Load Testing

### Test with Apache Bench

```bash
# Install Apache Bench
sudo apt-get install -y apache2-utils

# Test 500 requests, 50 concurrent
ab -n 500 -c 50 https://seed2.archivas.ai/chainTip

# Expected:
# - No 504 Gateway Timeouts
# - Majority should be 200 OK
# - Some might be 429 (rate limited) - this is correct!
```

### Test with wrk

```bash
# Install wrk
sudo apt-get install -y wrk

# Test 10 seconds, 10 threads, 100 connections
wrk -t10 -c100 -d10s https://seed2.archivas.ai/chainTip

# Monitor cache hit ratio during test:
sudo tail -f /var/log/nginx/seed2-access.log | grep -oP 'X-Cache-Status: \K\w+' | sort | uniq -c
```

### Stress Test (500 RPS)

```bash
# Test 500 req/s for 30 seconds
wrk -t20 -c200 -d30s --rate 500 https://seed2.archivas.ai/chainTip

# Expected:
# - Nginx should handle it smoothly
# - Cache hit ratio should be high (>90%)
# - Some requests rate limited (429)
# - No 502/504 errors
```

---

## 🚨 Troubleshooting

### Issue: 502 Bad Gateway

**Cause**: Seed1 backend is down or unreachable.

**Fix**:
```bash
# Check if Seed1 is responding
curl -I http://57.129.148.132:8080/chainTip

# Check Nginx error logs
sudo tail -50 /var/log/nginx/seed2-error.log

# Verify upstream configuration
sudo nginx -t

# Check network connectivity
ping 57.129.148.132
traceroute 57.129.148.132
```

### Issue: Cache not working (always MISS)

**Cause**: Cache directory permissions or configuration issue.

**Fix**:
```bash
# Check cache directory permissions
sudo ls -la /var/cache/nginx/

# Ensure nginx can write
sudo chown -R www-data:www-data /var/cache/nginx/archivas

# Clear cache and restart
sudo rm -rf /var/cache/nginx/archivas/*
sudo systemctl restart nginx
```

### Issue: High 429 rate limits

**Cause**: Legitimate traffic exceeding rate limits.

**Fix**:
```bash
# Increase rate limits in config
sudo nano /etc/nginx/sites-available/archivas-seed2

# Change:
# limit_req_zone $binary_remote_addr zone=req_ip:10m rate=30r/s;
# To:
# limit_req_zone $binary_remote_addr zone=req_ip:10m rate=50r/s;

# And increase burst:
# limit_req zone=req_ip burst=100 nodelay;

# Reload
sudo systemctl reload nginx
```

### Issue: SSL certificate renewal failed

**Cause**: DNS not pointing to server, or firewall blocking port 80.

**Fix**:
```bash
# Check DNS
dig seed2.archivas.ai

# Check firewall
sudo ufw status

# Manually renew
sudo certbot renew --dry-run

# Check certbot logs
sudo journalctl -u certbot -n 50
```

---

## 📚 Architecture Diagram

```
┌─────────────────┐
│   Explorer      │
│   Frontend      │
└────────┬────────┘
         │ HTTPS
         ↓
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Seed1 (57...2) │────→│  Seed2 (Relay)  │←────│  Seed3 (Relay)  │
│  Primary Node   │     │  Nginx Proxy    │     │  Nginx Proxy    │
│  Port: 8080     │     │  Port: 443      │     │  Port: 443      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         ↑                       ↑                       ↑
         │                       │                       │
         └───────────────────────┴───────────────────────┘
                    Round-Robin + Failover
```

---

## 📋 Acceptance Checklist

- [ ] DNS points to Seed2 IP
- [ ] SSL certificate installed and valid
- [ ] `/chainTip` returns 200 OK
- [ ] Second request shows `X-Cache-Status: HIT`
- [ ] Rate limiting works (429 after ~60 req/s)
- [ ] CORS headers present
- [ ] `/relay/status` returns healthy
- [ ] 500 RPS load test passes without 504s
- [ ] Explorer failover works (tested by stopping Seed1 temporarily)
- [ ] Logs show cache hit ratio >70%

---

## 🎯 Production Checklist

- [ ] Seed1 hardened and stable
- [ ] Seed2 relay deployed and tested
- [ ] Seed3 relay deployed (optional)
- [ ] Explorer uses multi-RPC
- [ ] DNS records configured
- [ ] SSL certificates auto-renewing
- [ ] Monitoring in place (logs, alerts)
- [ ] Backup plan documented
- [ ] Firewall rules configured
- [ ] Rate limits tuned for production
- [ ] Load tested at 500+ RPS

---

## 🔗 Related Documentation

- [P2P Discovery Isolation](./P2P-DISCOVERY-ISOLATION.md)
- [Deployment Guide](./DEPLOYMENT-P2P-ISOLATION.md)
- [Multi-RPC Failover](./MULTI-RPC-FAILOVER.md)
- [Nginx Configuration](../deploy/seed2/nginx-seed2-relay.conf)

---

## 📞 Support

- **GitHub**: https://github.com/ArchivasNetwork/archivas
- **Issues**: https://github.com/ArchivasNetwork/archivas/issues
- **Docs**: https://docs.archivas.ai

---

**Last Updated**: 2025-11-13  
**Maintainer**: Archivas Network Team

