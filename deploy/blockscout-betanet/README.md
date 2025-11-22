# Blockscout Deployment for Archivas Betanet

This directory contains the deployment configuration for running Blockscout as the official block explorer and contract indexer for Archivas Betanet.

---

## Overview

**Blockscout** is an open-source blockchain explorer that provides:
- Block and transaction browsing
- Address balance and transaction history
- Smart contract verification and interaction
- Real-time indexing and search
- API endpoints for dApps

**Archivas Betanet Configuration:**
- Chain ID: 1644
- Currency: RCHV
- RPC: https://seed3.betanet.archivas.ai
- Explorer URL: https://explorer.betanet.archivas.ai

---

## Prerequisites

### System Requirements

- **OS:** Linux (Ubuntu 20.04+) or macOS
- **RAM:** Minimum 4GB (8GB+ recommended for production)
- **Disk:** 50GB+ free space (for database and indexed data)
- **Docker:** v20.10+
- **Docker Compose:** v2.0+

### Install Docker and Docker Compose

```bash
# Install Docker (Ubuntu)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# Log out and back in for group changes to take effect

# Verify installation
docker --version
docker compose version
```

---

## Quick Start

### Step 1: Clone Repository

```bash
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas/deploy/blockscout-betanet
```

### Step 2: Configure Environment

```bash
# Copy example environment file
cp env.example .env

# Edit configuration
nano .env
```

**Required changes:**

```env
# Set a secure database password
POSTGRES_PASSWORD=your_secure_password_here

# Generate and set secret key
# Run: openssl rand -base64 64
SECRET_KEY_BASE=your_generated_secret_key_here

# Set your domain (if using custom domain)
BLOCKSCOUT_HOST=explorer.betanet.archivas.ai
BLOCKSCOUT_PROTOCOL=https
```

### Step 3: Start Blockscout

```bash
# Start all services
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f blockscout
```

**Expected output:**

```
[info] Running Explorer.Repo Migrations
[info] == Running 1.0.0 Explorer.Repo.Migrations.CreateBlocks.change/0 forward
[info] create table blocks
[info] == Migrated 1.0.0 in 0.0s
...
[info] Access BlockScoutWeb.Endpoint at http://localhost:4000
```

### Step 4: Verify Deployment

```bash
# Check Blockscout health
curl http://localhost:4000/api/v1/health/liveness

# Expected: {"healthy":true}

# Check if indexing started
curl http://localhost:4000/api/v1/stats

# View in browser
# Navigate to: http://localhost:4000
# Or: https://explorer.betanet.archivas.ai (if using custom domain)
```

---

## Configuration Details

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RPC_HTTP_URL` | Archivas Betanet RPC endpoint | `https://seed3.betanet.archivas.ai` |
| `POSTGRES_PASSWORD` | PostgreSQL database password | `changeme` (CHANGE THIS!) |
| `SECRET_KEY_BASE` | Secret for signing sessions | Random (CHANGE THIS!) |
| `BLOCKSCOUT_HOST` | Domain for Blockscout | `explorer.betanet.archivas.ai` |
| `BLOCKSCOUT_VERSION` | Blockscout Docker image version | `latest` |

### Network Configuration

Blockscout is configured for Archivas Betanet:

- **Chain ID:** 1644 (`0x66c`)
- **Network Name:** Archivas Betanet
- **Currency Symbol:** RCHV
- **Block Transformer:** Base (standard Ethereum format)
- **EVM Version:** London

### Indexer Configuration

The following indexers are enabled:

- ✅ **Block Indexer:** Fetches and indexes blocks
- ✅ **Transaction Indexer:** Indexes transactions and receipts
- ✅ **Pending Transaction Fetcher:** Indexes mempool transactions
- ✅ **Coin Balance Indexer:** Tracks address balances
- ✅ **Token Indexer:** Indexes ERC-20/721/1155 tokens (when deployed)
- ❌ **Internal Transaction Fetcher:** Disabled (not yet supported by Archivas)

---

## Reverse Proxy Setup (Nginx)

To serve Blockscout via HTTPS with a custom domain, set up Nginx as a reverse proxy.

### Install Nginx and Certbot

```bash
sudo apt update
sudo apt install -y nginx certbot python3-certbot-nginx
```

### Create Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/blockscout
```

**Add the following:**

```nginx
upstream blockscout_backend {
    server 127.0.0.1:4000;
    keepalive 64;
}

server {
    listen 80;
    listen [::]:80;
    server_name explorer.betanet.archivas.ai;

    # Redirect HTTP to HTTPS
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name explorer.betanet.archivas.ai;

    # SSL certificates (will be configured by Certbot)
    ssl_certificate /etc/letsencrypt/live/explorer.betanet.archivas.ai/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/explorer.betanet.archivas.ai/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;

    # Proxy settings
    location / {
        proxy_pass http://blockscout_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_buffering off;
    }

    # Increase timeouts for long-running queries
    proxy_connect_timeout 300s;
    proxy_send_timeout 300s;
    proxy_read_timeout 300s;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
}
```

### Enable Configuration and Get SSL Certificate

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/blockscout /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Get SSL certificate
sudo certbot --nginx -d explorer.betanet.archivas.ai --non-interactive --agree-tos --email info@archivas.ai

# Restart Nginx
sudo systemctl restart nginx
```

---

## Monitoring and Maintenance

### Check Service Status

```bash
# View all container status
docker compose ps

# Check Blockscout logs
docker compose logs -f blockscout

# Check database logs
docker compose logs -f db

# Check indexing progress
curl http://localhost:4000/api/v1/stats | jq
```

### Database Backup

```bash
# Backup database
docker compose exec db pg_dump -U blockscout blockscout > blockscout_backup_$(date +%Y%m%d).sql

# Restore database
cat blockscout_backup_20251121.sql | docker compose exec -T db psql -U blockscout blockscout
```

### Restart Services

```bash
# Restart Blockscout only
docker compose restart blockscout

# Restart all services
docker compose restart

# Stop all services
docker compose down

# Start with fresh database (WARNING: deletes all data)
docker compose down -v
docker compose up -d
```

### Update Blockscout

```bash
# Pull latest image
docker compose pull blockscout

# Restart with new image
docker compose up -d blockscout
```

---

## Smart Contract Verification

Blockscout supports smart contract verification via:

1. **Web UI:** Upload source code and verify manually
2. **API:** Programmatic verification via API
3. **Hardhat Plugin:** Verify during deployment

### Verify via Hardhat

In your Hardhat project:

```javascript
// hardhat.config.js
module.exports = {
  // ... other config
  etherscan: {
    apiKey: {
      betanet: "NOT_REQUIRED"
    },
    customChains: [
      {
        network: "betanet",
        chainId: 1644,
        urls: {
          apiURL: "https://explorer.betanet.archivas.ai/api",
          browserURL: "https://explorer.betanet.archivas.ai"
        }
      }
    ]
  }
};
```

**Verify a contract:**

```bash
npx hardhat verify --network betanet \
  0xCONTRACT_ADDRESS \
  "Constructor Argument 1" \
  "Constructor Argument 2"
```

---

## Troubleshooting

### Blockscout Not Starting

**Check logs:**

```bash
docker compose logs blockscout
```

**Common issues:**

1. **Database not ready:**
   - Wait for PostgreSQL to finish initializing
   - Check: `docker compose logs db`

2. **RPC connection failed:**
   - Verify `RPC_HTTP_URL` is correct
   - Test: `curl https://seed3.betanet.archivas.ai -X POST -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'`

3. **Port already in use:**
   - Change port in docker-compose.yml: `ports: - "4001:4000"`

### Indexing Not Progressing

**Check indexer status:**

```bash
curl http://localhost:4000/api/v1/stats | jq '.indexed_blocks'
```

**If stuck:**

1. **Check RPC health:**
   ```bash
   curl https://seed3.betanet.archivas.ai -X POST \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}'
   ```

2. **Restart indexer:**
   ```bash
   docker compose restart blockscout
   ```

3. **Check for errors:**
   ```bash
   docker compose logs blockscout | grep -i "error\|failed"
   ```

### Contract Verification Failing

**Check compiler version:**
- Ensure Solidity version in contract matches verification request
- Try verifying manually via Web UI first

**Check constructor arguments:**
- Must be ABI-encoded hex string
- Use Hardhat's encoding: `await ethers.utils.defaultAbiCoder.encode([...], [...])`

---

## Performance Tuning

### Database Optimization

```bash
# Edit docker-compose.yml to add PostgreSQL tuning
# Under db service, add:
command:
  - "postgres"
  - "-c"
  - "shared_buffers=256MB"
  - "-c"
  - "effective_cache_size=1GB"
  - "-c"
  - "max_connections=200"
```

### Indexer Performance

```bash
# In .env, adjust indexer batch sizes:
INDEXER_CATCHUP_BLOCKS_BATCH_SIZE=50
INDEXER_CATCHUP_BLOCKS_CONCURRENCY=20
```

---

## Useful API Endpoints

| Endpoint | Description |
|----------|-------------|
| `/api/v1/stats` | Network statistics (total blocks, transactions, etc.) |
| `/api/v1/health/liveness` | Health check (for monitoring) |
| `/api/v2/addresses/:address` | Address details (balance, txs) |
| `/api/v2/blocks` | List of recent blocks |
| `/api/v2/transactions/:hash` | Transaction details |
| `/api/v2/search` | Search for address, tx, block |

**Example:**

```bash
# Get latest stats
curl https://explorer.betanet.archivas.ai/api/v1/stats | jq

# Search for an address
curl https://explorer.betanet.archivas.ai/api/v2/search?q=0x47ea4b22029c155c835fd0a0b99f8196766f406a | jq
```

---

## Additional Resources

- [Blockscout Documentation](https://docs.blockscout.com/)
- [Blockscout GitHub](https://github.com/blockscout/blockscout)
- [Archivas Documentation](https://docs.archivas.ai)
- [Archivas Betanet RPC](https://seed3.betanet.archivas.ai)

---

## Support

For issues related to:

- **Blockscout deployment:** Open issue on [Blockscout GitHub](https://github.com/blockscout/blockscout/issues)
- **Archivas integration:** Open issue on [Archivas GitHub](https://github.com/ArchivasNetwork/archivas/issues)
- **Community help:** Join [Archivas Discord](https://discord.gg/archivas)

---

**Happy exploring! 🔍**

