# Seed Node Setup (seed.archivas.ai)

Production bootstrap node with HTTPS, HTTP/2, and full v1.1.0 API support.

---

## Overview

**Purpose:** Public, stable bootstrap node for Archivas network  
**Domain:** `seed.archivas.ai`  
**Server:** 57.129.148.132 (Server A)  
**Stack:** Nginx (reverse proxy) → Archivas node (RPC :8080)  
**TLS:** Let's Encrypt (automated renewal)  

---

## Prerequisites

1. **DNS Configuration**
   - Create A record: `seed.archivas.ai` → `57.129.148.132`
   - Wait for propagation (usually 5-15 minutes)
   - Verify: `dig seed.archivas.ai +short`

2. **Server Access**
   - SSH access to 57.129.148.132
   - Sudo privileges
   - Ports 80 and 443 must be accessible

3. **Firewall Configuration**
   ```bash
   # Allow public HTTP/HTTPS
   sudo ufw allow 80/tcp comment "HTTP for ACME"
   sudo ufw allow 443/tcp comment "HTTPS public RPC"
   
   # DENY direct access to node RPC (localhost only)
   sudo ufw deny 8080/tcp comment "Block external RPC access"
   
   # Allow P2P (if serving as bootnode)
   sudo ufw allow 9090/tcp comment "Archivas P2P"
   ```

4. **Running Node**
   - Archivas node must be running on `localhost:8080` (NOT `0.0.0.0`)
   - Version: v1.1.0 (Wallet API Freeze)
   - Network: archivas-devnet-v3
   
   **Critical:** Bind RPC to localhost only:
   ```bash
   ./archivas-node \
     --rpc 127.0.0.1:8080 \
     --p2p :9090 \
     --genesis genesis/devnet.genesis.json \
     --network-id archivas-devnet-v3
   ```
   
   **Never use:** `--rpc 0.0.0.0:8080` (exposes RPC to internet)

---

## Deployment Steps

### 1. Point DNS

```bash
# On your DNS provider (e.g., Cloudflare, Route53):
# Create A record:
#   Name:  seed.archivas.ai
#   Type:  A
#   Value: 57.129.148.132
#   TTL:   Auto or 3600

# Verify DNS propagation:
dig seed.archivas.ai +short
# Should return: 57.129.148.132
```

### 2. Install Nginx and Certbot

**On Server A (57.129.148.132):**

```bash
cd ~/archivas

# Run the setup script (idempotent, safe to rerun):
sudo bash scripts/setup-seed-nginx.sh

# Expected output:
# 🌱 Setting up seed.archivas.ai Nginx proxy
# ==========================================
#
# 1️⃣  Installing Nginx...
#    ✅ Nginx installed
#
# 2️⃣  Installing Certbot...
#    ✅ Certbot installed
#
# 3️⃣  Creating ACME webroot...
#    ✅ Webroot created: /var/www/seed
#
# 4️⃣  Installing Nginx site configuration...
#    ✅ Site config copied
#
# 5️⃣  Enabling site...
#    ✅ Site enabled
#
# 6️⃣  Testing Nginx configuration...
#    ✅ Nginx configuration valid
#
# 7️⃣  Configuring firewall...
#    ✅ Firewall rules added
#
# 8️⃣  Reloading Nginx...
#    ✅ Nginx reloaded
#
# ✅ Nginx setup complete!
```

### 3. Obtain TLS Certificate

**After DNS propagation:**

```bash
# Obtain certificate from Let's Encrypt:
sudo certbot --nginx \
  -d seed.archivas.ai \
  --non-interactive \
  --agree-tos \
  -m admin@archivas.ai

# Expected output:
# Saving debug log to /var/log/letsencrypt/letsencrypt.log
# Requesting a certificate for seed.archivas.ai
#
# Successfully received certificate.
# Certificate is saved at: /etc/letsencrypt/live/seed.archivas.ai/fullchain.pem
# Key is saved at:         /etc/letsencrypt/live/seed.archivas.ai/privkey.pem
# This certificate expires on YYYY-MM-DD.
```

### 4. Setup Auto-Renewal

```bash
# Install renewal hook:
sudo cp scripts/renew-cert-hook.sh /etc/letsencrypt/renewal-hooks/post/
sudo chmod +x /etc/letsencrypt/renewal-hooks/post/renew-cert-hook.sh

# Enable certbot timer (auto-renewal):
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer

# Verify timer is active:
sudo systemctl status certbot.timer

# Test renewal (dry-run):
sudo certbot renew --dry-run
```

### 5. Verify Deployment

```bash
# Run validation script:
bash scripts/check-seed.sh

# Expected output:
# 🔍 Checking seed.archivas.ai
# =============================
#
# 📡 Testing HTTP → HTTPS redirect...
# Testing HTTP redirect (GET http://seed.archivas.ai/version)... ✅ 301
#
# 🔐 Testing HTTPS endpoints...
# Testing /version (GET https://seed.archivas.ai/version)... ✅ 200
# Testing /chainTip (GET https://seed.archivas.ai/chainTip)... ✅ 200
# Testing /health (GET https://seed.archivas.ai/health)... ✅ 200
#
# 📝 Testing /submit method handling...
# Testing /submit GET (should 405) (GET https://seed.archivas.ai/submit)... ✅ 405
# Testing /submit OPTIONS (OPTIONS https://seed.archivas.ai/submit)... ✅ 204
#
# 🌐 Testing CORS headers...
# Checking CORS headers... ✅ Present
#
# ✅ seed.archivas.ai validation complete
```

---

## Testing

### Manual Tests

```bash
# Version:
curl https://seed.archivas.ai/version
# Expected: {"version":"v1.1.0",...}

# Chain tip:
curl https://seed.archivas.ai/chainTip
# Expected: {"height":"13080","hash":"...","difficulty":"1000000"}

# Account balance:
curl https://seed.archivas.ai/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
# Expected: {"address":"arcv1...","balance":"...","nonce":"0"}

# Mempool:
curl https://seed.archivas.ai/mempool
# Expected: [] (array of tx hashes)

# Submit GET (should return 405):
curl -i https://seed.archivas.ai/submit
# Expected: HTTP/2 405
#           allow: POST

# Submit OPTIONS preflight:
curl -i -X OPTIONS https://seed.archivas.ai/submit
# Expected: HTTP/2 204
#           access-control-allow-methods: POST,OPTIONS

# Verify /metrics is blocked:
curl -i https://seed.archivas.ai/metrics
# Expected: HTTP/2 404

# Verify node RPC is NOT externally accessible:
curl http://57.129.148.132:8080/chainTip
# Expected: Connection refused (port blocked by firewall)
```

### From SDK

```typescript
import { createRpcClient } from '@archivas/sdk';

const rpc = createRpcClient({ baseUrl: 'https://seed.archivas.ai' });

const tip = await rpc.getChainTip();
console.log('Height:', tip.height);
```

---

## Monitoring

### Nginx Logs

```bash
# Access log:
sudo tail -f /var/log/nginx/seed.archivas.ai.access.log

# Error log:
sudo tail -f /var/log/nginx/seed.archivas.ai.error.log
```

### Certificate Expiry

```bash
# Check certificate validity:
sudo certbot certificates

# Manual renewal (if needed):
sudo certbot renew
```

### Health Check

```bash
# Add to cron for monitoring:
echo "*/5 * * * * /home/ubuntu/archivas/scripts/check-seed.sh >> /var/log/archivas-seed-health.log 2>&1" | sudo crontab -
```

---

## Troubleshooting

### DNS Not Resolving

```bash
# Check DNS:
dig seed.archivas.ai +short

# If empty, wait for propagation or check DNS provider
```

### Certificate Issuance Fails

```bash
# Common issues:
# 1. DNS not pointing to server
# 2. Port 80 blocked
# 3. Webroot permissions

# Check:
sudo certbot certificates
sudo nginx -t
curl -I http://seed.archivas.ai/.well-known/acme-challenge/test

# Re-run:
sudo certbot --nginx -d seed.archivas.ai --non-interactive --agree-tos -m admin@archivas.ai
```

### 502 Bad Gateway

```bash
# Node RPC not running:
curl http://localhost:8080/version

# If fails, restart node:
cd ~/archivas
pkill -f archivas-node
nohup ./archivas-node --rpc 0.0.0.0:8080 ... > logs/node.log 2>&1 &
```

### CORS Issues

```bash
# Verify CORS headers:
curl -I https://seed.archivas.ai/chainTip | grep -i access-control

# Should show:
# access-control-allow-origin: *
```

---

## Rollback

### Remove Nginx Configuration

```bash
# Disable site:
sudo rm /etc/nginx/sites-enabled/seed.archivas.ai
sudo systemctl reload nginx

# Remove configuration:
sudo rm /etc/nginx/sites-available/seed.archivas.ai
```

### Revoke Certificate

```bash
# Revoke and delete:
sudo certbot revoke --cert-name seed.archivas.ai
sudo certbot delete --cert-name seed.archivas.ai
```

---

## Maintenance

### Update Nginx Configuration

```bash
# After modifying deploy/seed/nginx-site.conf:
cd ~/archivas
git pull origin main

# Re-run setup (idempotent):
sudo bash scripts/setup-seed-nginx.sh

# Test and reload:
sudo nginx -t
sudo systemctl reload nginx
```

### Rotate Certificates

Automatic via certbot timer. Certificates auto-renew 30 days before expiry.

---

## Security

- ✅ HTTPS enforced (HTTP redirects to HTTPS)
- ✅ TLS 1.2+ only
- ✅ HTTP/2 enabled
- ✅ Security headers (X-Frame-Options, X-Content-Type-Options, HSTS)
- ✅ CORS configured (permissive for public API)
- ✅ Method validation (/submit POST only)
- ✅ Rate limiting (can add `limit_req` if needed)

---

## Architecture

```
Internet
   ↓
DNS (seed.archivas.ai → 57.129.148.132)
   ↓
Nginx (:80 redirect, :443 TLS termination)
   │
   ├─ Rate limiting (10 req/min on /submit)
   ├─ CORS headers
   ├─ Block /metrics (internal only)
   │
   ↓
Archivas Node RPC (127.0.0.1:8080 - localhost ONLY)
   ↓
Blockchain + P2P (:9090)
```

**Security Layers:**
1. Nginx blocks `/metrics` from public access
2. Node RPC bound to `127.0.0.1` only (not accessible externally)
3. Rate limiting on `/submit` (10 req/min per IP, burst 5)
4. Firewall blocks port 8080 from external access
5. TLS 1.2+ with modern ciphers only

---

## Advanced Configuration

### Rate Limiting (Optional Enhancement)

The default Nginx config includes rate limiting for `/submit`. To adjust:

```nginx
# In deploy/seed/nginx-site.conf at top level:
limit_req_zone $binary_remote_addr zone=submit:10m rate=10r/m;

# In location = /submit:
limit_req zone=submit burst=5 nodelay;
```

Adjust `rate=10r/m` (requests per minute) and `burst=5` as needed.

### Node Configuration for Seed Host

Optional: Configure node to advertise as bootnode:

```bash
./archivas-node \
  --rpc 127.0.0.1:8080 \
  --p2p :9090 \
  --advertised-addr seed.archivas.ai:9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3
```

This allows other nodes to discover and connect via DNS.

---

## Files

- `deploy/seed/nginx-site.conf` - Nginx site configuration
- `scripts/setup-seed-nginx.sh` - Idempotent installer
- `scripts/renew-cert-hook.sh` - Certificate renewal hook
- `scripts/check-seed.sh` - Health validation

---

---

## Metrics After Hardening

After v1.1.1-infra, `/metrics` is **intentionally blocked** from public access for security.

### Public Access (Blocked)

```bash
curl https://seed.archivas.ai/metrics
# HTTP/2 404 (expected - this is correct)
```

### Localhost Scraping (Correct Method)

Prometheus must scrape `127.0.0.1` instead of public IPs:

**File:** `/etc/prometheus/prometheus.yml`

```yaml
scrape_configs:
  - job_name: 'archivas-nodes'
    static_configs:
      - targets: ['127.0.0.1:8080']  # ✅ Correct
        # NOT: ['57.129.148.132:8080']  # ❌ Wrong (blocked)
```

Restart Prometheus:
```bash
sudo systemctl restart prometheus
```

Verify:
```bash
# Check targets:
curl http://localhost:9090/targets

# Test scrape:
curl http://127.0.0.1:8080/metrics | head -20
```

### Multi-Server Monitoring

For Server C (or other nodes), use Prometheus **remote_write**:

**On Server C:**
```yaml
# /etc/prometheus/prometheus.yml
scrape_configs:
  - job_name: 'archivas-nodes'
    static_configs:
      - targets: ['127.0.0.1:8080']

remote_write:
  - url: http://57.129.148.132:9090/api/v1/write
```

See [docs/METRICS_LOCALHOST.md](METRICS_LOCALHOST.md) for complete guide.

### Troubleshooting

**"Connection refused" errors:**
1. Check service is running: `ps aux | grep archivas`
2. Check binding: `sudo lsof -i :8080`
3. Update Prometheus targets to use `127.0.0.1`

**Grafana shows "No data":**
1. Check Prometheus targets: `http://localhost:9090/targets`
2. Increase time range to "Last 1 hour"
3. Test query: `curl 'http://localhost:9090/api/v1/query?query=archivas_tip_height'`

---

## Contact

For issues with seed.archivas.ai:
- GitHub Issues: https://github.com/ArchivasNetwork/archivas/issues
- Email: admin@archivas.ai

