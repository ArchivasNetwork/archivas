# Seed2 Deployment Summary

## Quick Start

1. **Configure DNS**: Point `seed2.archivas.ai` to Server D's public IP
2. **SSH into Server D** and run:
   ```bash
   bash deploy/seed2/setup.sh
   ```
3. **Verify**:
   ```bash
   curl https://seed2.archivas.ai/healthz
   curl https://seed2.archivas.ai/chainTip | jq
   ```

## What Gets Deployed

- ✅ Archivas Node (RPC: `127.0.0.1:8080`, P2P: `0.0.0.0:9090`)
- ✅ Nginx Reverse Proxy (HTTPS on port 443)
- ✅ TLS Certificate (Let's Encrypt)
- ✅ Rate Limiting (API, Challenge, SubmitBlock, Blocks)
- ✅ CORS Headers (for web clients)
- ✅ Firewall Rules (HTTP, HTTPS, P2P)

## Network Configuration

- **RPC**: `127.0.0.1:8080` (localhost only, exposed via Nginx)
- **P2P**: `0.0.0.0:9090` (public, for peer connections)
- **Network ID**: `archivas-devnet-v4` (must match seed.archivas.ai)

## Rate Limits

- **API**: 10 requests/second per IP
- **Challenge**: 120 requests/minute (2 req/sec)
- **SubmitBlock**: 60 requests/minute (1 req/sec)
- **Blocks**: 10 requests/minute
- **Health**: 60 requests/minute
- **ChainTip**: No rate limiting (backend caching)

## Multi-RPC Failover

After seed2 is deployed, update:

1. **Explorer** (`.env`):
   ```env
   NEXT_PUBLIC_RPC_PRIMARY=https://seed.archivas.ai
   NEXT_PUBLIC_RPC_SECONDARY=https://seed2.archivas.ai
   NEXT_PUBLIC_RPC_TERTIARY=https://seed3.archivas.ai
   ```

2. **SDK** (use `baseUrls`):
   ```typescript
   const rpc = new Rpc({
     baseUrls: [
       'https://seed.archivas.ai',
       'https://seed2.archivas.ai',
       'https://seed3.archivas.ai',
     ],
   });
   ```

3. **Farmers**: Can connect to any seed node:
   ```bash
   ./archivas-farmer farm --node https://seed2.archivas.ai ...
   ```

## Files Created

- `/etc/systemd/system/archivas-node-seed2.service` - Node service
- `/etc/nginx/sites-available/seed2.archivas.ai` - Nginx config
- `/etc/nginx/snippets/archivas-cors.conf` - CORS headers
- `/etc/nginx/conf.d/archivas-ratelimit.conf` - Rate limits
- `/var/log/archivas/archivas-node-seed2.log` - Node logs
- `/var/log/nginx/seed2.archivas.ai.access.log` - Access logs
- `/var/log/nginx/seed2.archivas.ai.error.log` - Error logs

## Monitoring

```bash
# Check node status
sudo systemctl status archivas-node-seed2

# View logs
sudo journalctl -u archivas-node-seed2 -f

# Test endpoints
curl https://seed2.archivas.ai/healthz
curl https://seed2.archivas.ai/chainTip | jq
curl https://seed2.archivas.ai/challenge | jq
```

## Troubleshooting

### Node Not Starting

1. Check logs: `sudo journalctl -u archivas-node-seed2 -n 50`
2. Check database lock: `ls -la /home/ubuntu/archivas/data/LOCK`
3. Verify binary: `ls -la /home/ubuntu/archivas/archivas-node`

### Nginx Not Working

1. Test config: `sudo nginx -t`
2. Check site enabled: `ls -la /etc/nginx/sites-enabled/ | grep seed2`
3. Check error logs: `sudo tail -f /var/log/nginx/seed2.archivas.ai.error.log`

### TLS Certificate Issues

1. Check certificate: `sudo certbot certificates`
2. Renew certificate: `sudo certbot renew`
3. Manually obtain: `sudo certbot --nginx -d seed2.archivas.ai`

## Seed3 Deployment

To deploy seed3 (Server E), use the same process:

1. Copy `deploy/seed2/` to `deploy/seed3/` (already done)
2. Update `DOMAIN="seed3.archivas.ai"` in `setup.sh` (already done)
3. Run `deploy/seed3/setup.sh` on Server E
4. Configure DNS: `seed3.archivas.ai` → Server E IP

## Next Steps

1. ✅ Deploy seed2 (Server D)
2. ✅ Deploy seed3 (Server E)
3. ✅ Update Explorer with multi-RPC failover
4. ✅ Update SDK with `baseUrls` support
5. ✅ Test failover in production
6. ✅ Update documentation

## Success Criteria

- ✅ seed2.archivas.ai responds to HTTPS requests
- ✅ Node is synced and producing blocks
- ✅ Explorer fails over to seed2 when seed1 is down
- ✅ SDK automatically rotates between RPCs
- ✅ Farmers can connect to any seed node
- ✅ No single point of failure

