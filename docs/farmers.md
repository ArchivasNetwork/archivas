# Farmer Guide - Peer Configuration

## Recommended Peer Setup

Farmers should connect to **both Seed1 and Seed2** for optimal performance and reliability.

### Why dual-peer?

- **Load distribution**: Reduces pressure on any single seed node
- **Redundancy**: If one seed goes down, the other continues
- **Faster sync**: More sources for initial blockchain download
- **Better propagation**: Improved block and challenge distribution

## Peer Configuration

### For New Farmers (Initial Sync)

```bash
archivas-farmer \
  --farmer-privkey YOUR_FARMER_PRIVKEY \
  --plot-dir ./plots \
  --p2p-peer seed.archivas.ai:30303 \
  --p2p-peer seed2.archivas.ai:30303 \
  --no-peer-discovery \
  --checkpoint-height <HEIGHT> \
  --checkpoint-hash <HASH> \
  --rpc http://localhost:8080
```

**Get checkpoint from:**
```bash
curl -s https://seed.archivas.ai:8081/chainTip | jq
```

### For Running Farmers (Standard Operation)

```bash
archivas-farmer \
  --farmer-privkey YOUR_FARMER_PRIVKEY \
  --plot-dir ./plots \
  --p2p-peer seed.archivas.ai:30303 \
  --p2p-peer seed2.archivas.ai:30303 \
  --no-peer-discovery \
  --rpc http://localhost:8080
```

### For Advanced Farmers (Multi-Peer)

If you want to connect to additional trusted peers (e.g., other validators):

```bash
archivas-farmer \
  --farmer-privkey YOUR_FARMER_PRIVKEY \
  --plot-dir ./plots \
  --p2p-peer seed.archivas.ai:30303 \
  --p2p-peer seed2.archivas.ai:30303 \
  --p2p-peer validator02.archivas.ai:30303 \
  --p2p-peer validator03.archivas.ai:30303 \
  --no-peer-discovery \
  --rpc http://localhost:8080
```

## P2P Ports

| Server | P2P Port | Purpose |
|--------|----------|---------|
| seed.archivas.ai | 30303 | Primary seed node (full history) |
| seed2.archivas.ai | 30303 | Secondary seed node (full history) |

**Note**: These are P2P ports for blockchain consensus, not HTTP/HTTPS ports.

## RPC Endpoints

For balance checks, transaction submission, and explorer queries:

| Endpoint | Protocol | Purpose |
|----------|----------|---------|
| https://seed.archivas.ai:8081 | HTTPS | Primary RPC (canonical) |
| https://seed2.archivas.ai | HTTPS | Cached RPC relay (faster reads) |

**Example queries:**

```bash
# Check chain tip (use Seed2 for cached response)
curl -s https://seed2.archivas.ai/chainTip | jq

# Check your balance (use Seed2 for cached response)
curl -s https://seed2.archivas.ai/account/YOUR_ADDRESS | jq

# Submit a transaction (both work, Seed2 proxies to Seed1)
curl -X POST https://seed2.archivas.ai/submitTx \
  -H "Content-Type: application/json" \
  -d '{"signedTx": "..."}'
```

## Firewall Configuration

If you're running a farmer, ensure these ports are accessible:

```bash
# Outbound connections needed (TCP)
30303 → seed.archivas.ai
30303 → seed2.archivas.ai
443 → seed2.archivas.ai (for HTTPS RPC)

# Inbound (optional, for other farmers to peer with you)
# If you want to act as a peer for other farmers:
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp
```

## Health Checks

### Check P2P connectivity

```bash
# From your farmer machine:
telnet seed.archivas.ai 30303
telnet seed2.archivas.ai 30303

# Both should connect (press Ctrl+C to exit)
```

### Check RPC availability

```bash
# Seed1 RPC
curl -s https://seed.archivas.ai:8081/chainTip | jq .height

# Seed2 RPC (relay)
curl -s https://seed2.archivas.ai/chainTip | jq .height

# Heights should match
```

### Check your node's peer count

```bash
# If running archivas-node locally:
curl -s http://localhost:8080/peers | jq length

# You should see 2+ peers (Seed1, Seed2, maybe others)
```

## Troubleshooting

### Cannot connect to Seed1 or Seed2

**Symptoms**: "Connection refused" or "Timeout"

**Diagnosis**:

```bash
# Check DNS resolution
nslookup seed.archivas.ai
nslookup seed2.archivas.ai

# Check connectivity
ping seed.archivas.ai
ping seed2.archivas.ai

# Check port accessibility
telnet seed.archivas.ai 30303
telnet seed2.archivas.ai 30303
```

**Common fixes**:

1. **Firewall blocking outbound**: Allow TCP port 30303 outbound
2. **DNS issues**: Use IP addresses instead:
   - Seed1: `57.129.148.132:30303`
   - Seed2: `<get from hosting provider>:30303`
3. **Seed temporarily down**: Wait 5-10 minutes and retry

### Syncing very slowly

**Symptoms**: Chain height advancing slowly (< 100 blocks/minute)

**Fixes**:

1. **Add both peers** (if using only one):
   ```bash
   --p2p-peer seed.archivas.ai:30303 \
   --p2p-peer seed2.archivas.ai:30303
   ```

2. **Use checkpoint** (skip old history):
   ```bash
   --checkpoint-height <HEIGHT> \
   --checkpoint-hash <HASH>
   ```

3. **Check system resources**:
   ```bash
   htop  # CPU should have spare capacity
   df -h  # Disk should have > 20% free
   iotop  # Disk I/O should not be maxed
   ```

### Getting forked

**Symptoms**: "prev hash mismatch" in logs, stuck at specific height

**Fixes**:

1. **Stop node/farmer**:
   ```bash
   pkill archivas-farmer
   pkill archivas-node
   ```

2. **Clear corrupted data**:
   ```bash
   rm -rf ~/.archivas/data
   # Or wherever your data directory is
   ```

3. **Restart with checkpoint** (from canonical chain):
   ```bash
   # Get fresh checkpoint from Seed1
   curl -s https://seed.archivas.ai:8081/chainTip | jq
   
   # Restart with both peers
   archivas-farmer \
     --farmer-privkey YOUR_PRIVKEY \
     --plot-dir ./plots \
     --p2p-peer seed.archivas.ai:30303 \
     --p2p-peer seed2.archivas.ai:30303 \
     --no-peer-discovery \
     --checkpoint-height <HEIGHT> \
     --checkpoint-hash <HASH>
   ```

### High CPU/Memory usage

**Symptoms**: Farmer using > 80% CPU or > 4GB RAM

**Diagnosis**:

```bash
# Check process stats
ps aux | grep archivas-farmer

# Check system resources
htop
```

**Fixes**:

1. **Reduce max peers** (if running archivas-node):
   ```bash
   --max-peers 50  # Default is 100
   ```

2. **Enable aggressive GC** (if Go-based):
   ```bash
   export GOGC=50  # More aggressive garbage collection
   export GOMAXPROCS=4  # Limit CPU cores
   ```

3. **Add more RAM** (recommended: 8GB for farmer + node)

## Best Practices

### Security

1. **Never share your farmer private key**
2. **Use `--no-peer-discovery`** (only connect to trusted peers)
3. **Keep software updated** (check GitHub releases)
4. **Use checkpoint flags** to avoid syncing malicious forks

### Performance

1. **Use SSD/NVMe** for blockchain data directory
2. **Ensure good network** (> 100 Mbps, low latency)
3. **Monitor disk space** (keep > 20% free)
4. **Restart weekly** to clear memory leaks (if any)

### Reliability

1. **Use both Seed1 and Seed2** as peers
2. **Set up monitoring** (alerts if node stops syncing)
3. **Backup farmer private key** (encrypted, offline storage)
4. **Test failover** (ensure you can sync from either seed alone)

## Quick Reference

```bash
# Recommended farmer command (copy-paste ready)
archivas-farmer \
  --farmer-privkey YOUR_FARMER_PRIVKEY \
  --plot-dir ./plots \
  --p2p-peer seed.archivas.ai:30303 \
  --p2p-peer seed2.archivas.ai:30303 \
  --no-peer-discovery \
  --rpc http://localhost:8080

# Check sync status
curl -s http://localhost:8080/chainTip | jq

# Check Seed1 canonical height
curl -s https://seed.archivas.ai:8081/chainTip | jq .height

# Check Seed2 relay status
curl -s https://seed2.archivas.ai/status | jq
```

## Support

- **Documentation**: https://docs.archivas.ai
- **Discord**: https://discord.gg/archivas
- **GitHub**: https://github.com/ArchivasNetwork/archivas

---

**Last Updated**: 2025-11-14  
**Version**: 2.0 (Seed2 P2P support)

