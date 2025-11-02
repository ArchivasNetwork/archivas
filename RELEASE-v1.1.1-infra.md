# Release Notes: v1.1.1-infra

**Release Date:** 2025-11-02  
**Tag:** v1.1.1-infra  
**Type:** Infrastructure-only (no protocol changes)

---

## Summary

v1.1.1-infra adds production-grade infrastructure for the Archivas public testnet, deploying a hardened seed node at `https://seed.archivas.ai` with TLS, rate limiting, and comprehensive security measures.

**This release contains NO changes to consensus, database schema, or RPC response formats.** It is fully compatible with v1.1.0 wallet API.

---

## 🌍 Public Seed Node

**Live Endpoint:** `https://seed.archivas.ai`

The Archivas public RPC is now available for developers to interact with the testnet.

### Available Routes

- `GET /account/<address>` - Get account balance and nonce
- `GET /chainTip` - Get current blockchain height, hash, and difficulty
- `GET /mempool` - List pending transactions
- `GET /tx/<hash>` - Get transaction details
- `GET /estimateFee?bytes=<n>` - Estimate transaction fee
- `POST /submit` - Submit a signed transaction (JSON, CORS-enabled)

### Examples

```bash
# Get chain tip
curl https://seed.archivas.ai/chainTip
# {"height":"14870","hash":"df4d75f99f0f9a55c836cc4c517e5da04fb9d3cb18b12e3340f2c4b73e6322ab","difficulty":"1000000"}

# Get account balance
curl https://seed.archivas.ai/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
# {"address":"arcv1...","balance":"...","nonce":"0"}

# Test method enforcement (returns 405)
curl -i https://seed.archivas.ai/submit
# HTTP/2 405
# allow: POST
```

---

## ✨ New Features

### Nginx Reverse Proxy

- **TLS Termination:** Let's Encrypt certificate (valid until 2026-01-31)
- **HTTP/2:** Enabled for performance
- **HTTPS Only:** HTTP redirects to HTTPS
- **Auto-Renewal:** Certificate renews automatically 30 days before expiry

### Security Hardening

- **Localhost Binding:** Node RPC bound to `127.0.0.1:8080` only (not externally accessible)
- **Metrics Protection:** `/metrics` endpoint blocked from public access (returns 404)
- **Firewall Rules:** Port 8080 blocked from external access
- **Rate Limiting:** 10 requests/min per IP on `/submit` endpoint (burst: 5)
- **Security Headers:**
  - `Strict-Transport-Security: max-age=31536000; includeSubDomains`
  - `X-Frame-Options: SAMEORIGIN`
  - `X-Content-Type-Options: nosniff`
  - `Referrer-Policy: no-referrer-when-downgrade`

### CORS Configuration

- **GET Endpoints:** Permissive CORS (`Access-Control-Allow-Origin: *`)
- **POST /submit:** CORS-enabled for JSON submissions
- **OPTIONS Support:** Preflight requests handled correctly

### Deployment Scripts

- **`scripts/setup-seed-nginx.sh`** - Idempotent installer for Nginx and certbot
- **`scripts/renew-cert-hook.sh`** - Post-renewal hook (reloads Nginx)
- **`scripts/check-seed.sh`** - Comprehensive health validation

### Documentation

- **`docs/SEED_HOST.md`** - Complete seed node setup guide
  - DNS configuration
  - Firewall rules
  - Node binding requirements
  - Rate limiting configuration
  - Security best practices
  - Troubleshooting
- **README.md** - Updated with public RPC endpoint section

---

## 🔒 Security

### Architecture

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

### Security Layers

1. **Nginx blocks `/metrics`** from public access
2. **Node RPC bound to `127.0.0.1`** only (not accessible externally)
3. **Rate limiting on `/submit`** (10 req/min per IP, burst 5)
4. **Firewall blocks port 8080** from external access
5. **TLS 1.2+ with modern ciphers** only

---

## ✅ Verification

All health checks passing as of deployment:

- ✅ HTTP → HTTPS redirect (301)
- ✅ TLS certificate valid
- ✅ `/version` → 200
- ✅ `/chainTip` → 200
- ✅ `/health` → 200
- ✅ `/submit` GET → 405 (POST only)
- ✅ `/submit` OPTIONS → 204 (CORS preflight)
- ✅ `/metrics` → 404 (blocked)
- ✅ CORS headers present
- ✅ Rate limiting active

**Chain Height at Deploy:** 14,870+

---

## 📦 Files Changed

**Infrastructure:**
- `deploy/seed/nginx-site.conf` - Nginx reverse proxy configuration
- `scripts/setup-seed-nginx.sh` - Deployment automation
- `scripts/renew-cert-hook.sh` - Certificate renewal hook
- `scripts/check-seed.sh` - Health validation script

**Documentation:**
- `docs/SEED_HOST.md` - Complete setup guide
- `README.md` - Public RPC section
- `CHANGELOG.md` - Release notes

**No changes to:**
- ❌ Consensus logic
- ❌ Database schema
- ❌ RPC response formats
- ❌ Transaction validation
- ❌ Block verification
- ❌ P2P protocol

---

## 🔧 Compatibility

**Fully compatible with v1.1.0:**
- All RPC endpoints unchanged
- Response formats identical
- Wallet API frozen
- Transaction format unchanged

**Upgrade Path:**
- No database migration required
- No genesis change
- No node restart required (infrastructure only)
- Existing clients continue to work

---

## 📊 Deployment Details

**Server:** 57.129.148.132 (Server A)  
**DNS:** seed.archivas.ai  
**Certificate:** Let's Encrypt (expires 2026-01-31)  
**Nginx Version:** 1.24.0  
**Node Version:** v1.1.0  
**Network:** archivas-devnet-v3  

---

## 🚀 Getting Started

### For Developers

```bash
# Test the endpoint
curl https://seed.archivas.ai/chainTip

# Using TypeScript SDK (coming soon)
npm install @archivas/sdk

import { createRpcClient } from '@archivas/sdk';
const rpc = createRpcClient({ baseUrl: 'https://seed.archivas.ai' });
const tip = await rpc.getChainTip();
```

### For Node Operators

See [docs/SEED_HOST.md](docs/SEED_HOST.md) for complete deployment instructions.

---

## 📝 Notes

- **Production Ready:** This infrastructure is production-grade and suitable for public use
- **No Breaking Changes:** Fully backward compatible with v1.1.0
- **Security Audited:** All configurations follow industry best practices
- **Monitoring:** Prometheus metrics available internally (blocked from public)

---

## 🔗 Links

- **Live Endpoint:** https://seed.archivas.ai
- **GitHub Release:** https://github.com/ArchivasNetwork/archivas/releases/tag/v1.1.1-infra
- **Documentation:** https://github.com/ArchivasNetwork/archivas/blob/main/docs/SEED_HOST.md
- **TypeScript SDK:** https://github.com/ArchivasNetwork/archivas-sdk

---

## ⚠️ Important

This is an **infrastructure-only** release. No changes to protocol, consensus, or API.

**If you are running a node:**
- No action required
- Your node continues to work as before
- Optional: Deploy your own seed node using the provided scripts

**If you are developing on Archivas:**
- Update your RPC endpoint to `https://seed.archivas.ai`
- All API endpoints remain the same
- Rate limits: 10 req/min on `/submit`

---

**Release verified and deployed successfully!** 🎉

