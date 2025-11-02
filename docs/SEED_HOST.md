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

3. **Running Node**
   - Archivas node must be running on `localhost:8080`
   - Version: v1.1.0 (Wallet API Freeze)
   - Network: archivas-devnet-v3

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

# Chain tip:
curl https://seed.archivas.ai/chainTip

# Account balance:
curl https://seed.archivas.ai/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7

# Mempool:
curl https://seed.archivas.ai/mempool

# Submit (should return 405 for GET):
curl -i https://seed.archivas.ai/submit
# Expect: HTTP/2 405
#         allow: POST

# OPTIONS preflight:
curl -i -X OPTIONS https://seed.archivas.ai/submit
# Expect: HTTP/2 204
#         access-control-allow-methods: GET,POST,OPTIONS
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
   ↓
Archivas Node RPC (localhost:8080)
   ↓
Blockchain + P2P
```

---

## Files

- `deploy/seed/nginx-site.conf` - Nginx site configuration
- `scripts/setup-seed-nginx.sh` - Idempotent installer
- `scripts/renew-cert-hook.sh` - Certificate renewal hook
- `scripts/check-seed.sh` - Health validation

---

## Contact

For issues with seed.archivas.ai:
- GitHub Issues: https://github.com/ArchivasNetwork/archivas/issues
- Email: admin@archivas.ai

