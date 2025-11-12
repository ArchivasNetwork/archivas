# Seed2 Deployment Guide

This directory contains deployment scripts for `seed2.archivas.ai` (Server D).

## Quick Start

1. **Configure DNS**: Point `seed2.archivas.ai` to Server D's public IP
2. **Run setup script**:
   ```bash
   bash deploy/seed2/setup.sh
   ```
3. **Verify**:
   ```bash
   curl https://seed2.archivas.ai/healthz
   curl https://seed2.archivas.ai/chainTip | jq
   ```

## What Gets Deployed

- **Archivas Node**: Running on `127.0.0.1:8080` (RPC) and `0.0.0.0:9090` (P2P)
- **Nginx Reverse Proxy**: HTTPS on port 443, HTTP redirect on port 80
- **TLS Certificate**: Automatically obtained via Let's Encrypt
- **Rate Limiting**: Configured for API endpoints
- **CORS**: Enabled for web clients
- **Firewall**: UFW configured to allow HTTP, HTTPS, and P2P

## Configuration

### Systemd Service

- **Service**: `archivas-node-seed2.service`
- **Logs**: `/var/log/archivas/archivas-node-seed2.log`
- **Data**: `/home/ubuntu/archivas/data`

### Nginx

- **Config**: `/etc/nginx/sites-available/seed2.archivas.ai`
- **CORS**: `/etc/nginx/snippets/archivas-cors.conf`
- **Rate Limits**: `/etc/nginx/conf.d/archivas-ratelimit.conf`

### Rate Limits

- **API**: 10 requests/second per IP
- **Challenge**: 120 requests/minute (2 req/sec)
- **SubmitBlock**: 60 requests/minute (1 req/sec)
- **Blocks**: 10 requests/minute
- **Health**: 60 requests/minute
- **ChainTip**: No rate limiting (backend caching)

## Monitoring

### Check Node Status
```bash
sudo systemctl status archivas-node-seed2
```

### View Logs
```bash
# Recent logs
sudo journalctl -u archivas-node-seed2 -n 50

# Follow logs
sudo journalctl -u archivas-node-seed2 -f

# Log file
tail -f /var/log/archivas/archivas-node-seed2.log
```

### Check Nginx
```bash
# Test configuration
sudo nginx -t

# Reload configuration
sudo systemctl reload nginx

# Check access logs
sudo tail -f /var/log/nginx/seed2.archivas.ai.access.log

# Check error logs
sudo tail -f /var/log/nginx/seed2.archivas.ai.error.log
```

### Test Endpoints
```bash
# Health check
curl https://seed2.archivas.ai/healthz | jq

# Chain tip
curl https://seed2.archivas.ai/chainTip | jq

# Challenge
curl https://seed2.archivas.ai/challenge | jq
```

## Troubleshooting

### Node Not Starting

1. Check logs:
   ```bash
   sudo journalctl -u archivas-node-seed2 -n 100
   ```

2. Check if database is locked:
   ```bash
   ls -la /home/ubuntu/archivas/data/LOCK
   ```

3. Verify binary exists:
   ```bash
   ls -la /home/ubuntu/archivas/archivas-node
   ```

### Nginx Not Working

1. Test configuration:
   ```bash
   sudo nginx -t
   ```

2. Check if site is enabled:
   ```bash
   ls -la /etc/nginx/sites-enabled/ | grep seed2
   ```

3. Check error logs:
   ```bash
   sudo tail -f /var/log/nginx/seed2.archivas.ai.error.log
   ```

### TLS Certificate Issues

1. Check certificate:
   ```bash
   sudo certbot certificates
   ```

2. Renew certificate:
   ```bash
   sudo certbot renew
   ```

3. Manually obtain certificate:
   ```bash
   sudo certbot --nginx -d seed2.archivas.ai
   ```

## Updating

To update the node:

```bash
cd /home/ubuntu/archivas
git pull origin main
go build -o archivas-node ./cmd/archivas-node
sudo systemctl restart archivas-node-seed2
```

## Seed3 Setup

To set up seed3 (Server E), use the same process but change the domain:

1. Copy `deploy/seed2/` to `deploy/seed3/`
2. Update `DOMAIN="seed3.archivas.ai"` in `setup.sh`
3. Update service name to `archivas-node-seed3`
4. Run the setup script on Server E

## Multi-RPC Failover

After seed2 and seed3 are deployed, update:

- **Explorer**: Add `NEXT_PUBLIC_RPC_SECONDARY` and `NEXT_PUBLIC_RPC_TERTIARY` to `.env`
- **SDK**: Use `baseUrls` array in RPC client
- **Farmers**: Can connect to any seed node

See `../docs/MULTI-RPC-FAILOVER.md` for details.

