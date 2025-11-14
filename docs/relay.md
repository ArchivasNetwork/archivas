# Archivas Seed2 - Stateless RPC Relay

**Version**: 2.0  
**Status**: Production  
**URL**: https://seed2.archivas.ai

---

## 📖 What is Seed2?

Seed2 is a **stateless RPC relay** for the Archivas network. It acts as a caching proxy layer between clients and the primary Seed1 node, providing:

✅ **Fast responses** through intelligent caching (5-10s TTL)  
✅ **High availability** with circuit breakers and automatic failover  
✅ **Rate limiting** to protect backend infrastructure  
✅ **Read optimization** for explorers, dashboards, and wallets  
✅ **TX submission proxy** with automatic retries

### What Seed2 IS:
- ✅ Read-only RPC endpoint with caching
- ✅ Transaction submission proxy
- ✅ Load balancer and circuit breaker
- ✅ Rate limiter and DDoS protection

### What Seed2 IS NOT:
- ❌ Full blockchain node (no P2P, no consensus)
- ❌ Farmer/validator peer (no block propagation)
- ❌ Data source (proxies to Seed1)
- ❌ WebSocket event streaming (use Seed1 directly)

---

## 👥 Who Should Use Seed2?

### ✅ **Explorers & Dashboards** → **YES, use Seed2 as primary**

Block explorers, network dashboards, and analytics tools should **default to Seed2** for:
- Chain tip queries (`/chainTip`)
- Recent blocks (`/blocks/recent`)
- Block lookups (`/block/{height}`)
- Account balances (`/account/{address}`)
- Transaction lookups (`/tx/{hash}`)

**Benefits:**
- 70-90% cache hit ratio = sub-second responses
- Reduced load on Seed1
- Built-in rate limiting prevents accidental DDoS

**Configuration:**
```env
# .env for explorer
RPC_PRIMARY=https://seed2.archivas.ai
RPC_FALLBACK=https://seed.archivas.ai
RPC_TIMEOUT=5000
RPC_RETRIES=2
```

---

### ✅ **Wallets & dApps** → **YES, use Seed2 for reads + TX submit**

Desktop/mobile wallets and decentralized applications should use Seed2 for:
- **Reads**: Balance checks, transaction history
- **Writes**: Transaction submission via `/submitTx`

**Benefits:**
- Fast balance lookups (10s cache)
- Automatic TX retry with backoff
- Fallback to Seed1 if Seed2 is down

**Configuration:**
```javascript
// wallet config
const RPC_CONFIG = {
  primary: 'https://seed2.archivas.ai',
  fallback: 'https://seed.archivas.ai',
  timeout: 5000,
  retries: 2
};

// For TX submission
const submitTransaction = async (signedTx) => {
  try {
    const response = await axios.post(`${RPC_CONFIG.primary}/submitTx`, signedTx);
    return response.data;
  } catch (error) {
    // Automatic retry happens at Seed2 level
    if (error.response?.status === 503) {
      // Wait and retry as indicated
      const retryAfter = error.response.data?.retryAfter || 2;
      await sleep(retryAfter * 1000);
      // Try fallback
      return axios.post(`${RPC_CONFIG.fallback}/submitTx`, signedTx);
    }
    throw error;
  }
};
```

---

### ❌ **Farmers & Validators** → **NO, connect to Seed1 P2P**

Farmers and validators **MUST NOT use Seed2** for peer-to-peer operations. They need direct P2P connections for:
- Block propagation
- Proof-of-Space challenges
- Consensus participation

**Use Seed1 for P2P:**

```bash
# Farmer configuration
./archivas-farmer \
  --farmer-privkey YOUR_PRIVKEY \
  --plot-dir ./plots \
  --p2p-peer seed.archivas.ai:9090 \
  --rpc http://localhost:8080

# Node configuration (if running your own)
./archivas-node \
  --data-dir ./data \
  --p2p 0.0.0.0:9090 \
  --peer seed.archivas.ai:9090 \
  --rpc 127.0.0.1:8080
```

**Why not Seed2?**
- Seed2 has **no P2P port** (no consensus)
- Seed2 **does not propagate blocks**
- Farmers connecting to Seed2 will **not receive challenges**
- Validators connecting to Seed2 will **fork**

**Farmers may use Seed2 for:**
- ✅ Checking balances/rewards (read-only queries)
- ✅ Submitting payout transactions

**Farmers must NOT use Seed2 for:**
- ❌ P2P peering (use `--p2p-peer seed.archivas.ai:9090`)
- ❌ Block propagation
- ❌ Proof challenges

---

## 🔧 Technical Details

### Cache Configuration

| Endpoint Pattern | Cache TTL | Purpose |
|-----------------|-----------|---------|
| `/chainTip` | 5 seconds | Fast tip queries |
| `/blocks/recent` | 5 seconds | Recent block lists |
| `/block/{height}` | 5 seconds | Block data (immutable) |
| `/account/{addr}` | 10 seconds | Balance checks |
| `/tx/{hash}` | 5 seconds | Transaction lookups |
| `/mempool` | **NO CACHE** | Real-time pending TXs |
| `/submitTx` | **NO CACHE** | Write operations |

### Rate Limiting

- **Per IP**: 30 requests/second
- **Burst**: 60 requests (allows temporary spikes)
- **Global**: No hard limit (relies on per-IP)
- **Response**: HTTP 429 when rate exceeded

### TX Submission Proxy

Seed2 provides intelligent TX submission with automatic retries:

**Request:**
```bash
curl -X POST https://seed2.archivas.ai/submitTx \
  -H "Content-Type: application/json" \
  -d '{"signedTx": "0x..."}'
```

**Success Response (200):**
```json
{
  "txid": "a1b2c3d4...",
  "status": "accepted"
}
```

**Retry Response (503):**
```json
{
  "error": "upstream_unavailable",
  "retryAfter": 2
}
```

**Retry Logic:**
1. First attempt (immediate)
2. Retry after 250ms if upstream error
3. Retry after 500ms if still failing
4. Final retry after 1s
5. Return 503 if all fail

### Circuit Breaker

Seed2 includes a circuit breaker to protect Seed1:

- **Threshold**: 10 upstream errors in 30 seconds
- **Action**: Trip circuit, return 503
- **Recovery**: Reset after 5 seconds
- **Monitoring**: `/status` shows circuit state

---

## 📊 Monitoring & Observability

### Health Endpoints

#### `GET /health` - Liveness Probe
```bash
curl https://seed2.archivas.ai/health
```
**Response:**
```json
{"ok": true, "timestamp": 1700000000000}
```
**Use**: Kubernetes liveness probe, uptime monitoring

---

#### `GET /ready` - Readiness Probe
```bash
curl https://seed2.archivas.ai/ready
```
**Response (healthy):**
```json
{
  "ready": true,
  "checks": {
    "upstream": {"healthy": true, "latency": 45},
    "cache": {"freeSpace": 87, "sufficient": true}
  },
  "timestamp": 1700000000000
}
```
**Response (degraded):**
```json
{
  "ready": false,
  "checks": {
    "upstream": {"healthy": false, "latency": 5000},
    "cache": {"freeSpace": 5, "sufficient": false}
  },
  "timestamp": 1700000000000
}
```
**Use**: Kubernetes readiness probe, load balancer health check

---

#### `GET /status` - Detailed Status
```bash
curl https://seed2.archivas.ai/status | jq .
```
**Response:**
```json
{
  "relay": "healthy",
  "cache": "enabled",
  "backend": "seed1",
  "upstream": {
    "url": "https://seed.archivas.ai",
    "healthy": true,
    "latency_ms": 42,
    "height": "785000"
  },
  "cache_stats": {
    "size_mb": 847.23,
    "free_space_percent": 87,
    "file_count": 1523,
    "hit_ratio_5m": 89.5,
    "hits": 8950,
    "misses": 1050,
    "expired": 500
  },
  "timestamp": 1700000000000
}
```
**Use**: Dashboards, debugging, capacity planning

---

#### `GET /metrics` - Prometheus Metrics
```bash
curl https://seed2.archivas.ai/metrics
```
**Response (sample):**
```
# HELP seed2_cache_hits_total Total number of cache hits
# TYPE seed2_cache_hits_total counter
seed2_cache_hits_total{endpoint="/chainTip"} 12543

# HELP seed2_cache_miss_total Total number of cache misses
# TYPE seed2_cache_miss_total counter
seed2_cache_miss_total{endpoint="/chainTip"} 1234

# HELP seed2_upstream_latency_ms Upstream request latency
# TYPE seed2_upstream_latency_ms histogram
seed2_upstream_latency_ms_bucket{le="50"} 8234
seed2_upstream_latency_ms_bucket{le="100"} 9123
seed2_upstream_latency_ms_bucket{le="500"} 9876
seed2_upstream_latency_ms_sum 456789
seed2_upstream_latency_ms_count 10000

# HELP seed2_tx_submit_success_total Successful TX submissions
# TYPE seed2_tx_submit_success_total counter
seed2_tx_submit_success_total 234

# HELP seed2_tx_submit_failed_total Failed TX submissions
# TYPE seed2_tx_submit_failed_total counter
seed2_tx_submit_failed_total 12
```
**Use**: Prometheus scraping, Grafana dashboards

---

## 🚀 Deployment

### For System Administrators

```bash
# Deploy Seed2 relay (idempotent)
cd /root/archivas
sudo ./scripts/deploy_seed2.sh

# Check status
sudo systemctl status nginx
sudo systemctl status archivas-relay

# View logs
sudo journalctl -u archivas-relay -f
sudo tail -f /var/log/nginx/seed2-access.log

# Check metrics
curl http://localhost:9090/status | jq .
```

### For Application Developers

```javascript
// Multi-RPC client with automatic failover
class ArchivasRPC {
  constructor() {
    this.endpoints = [
      'https://seed2.archivas.ai',  // Primary (cached)
      'https://seed.archivas.ai'     // Fallback (direct)
    ];
    this.currentIndex = 0;
  }

  async request(method, params, retries = 2) {
    for (let attempt = 0; attempt <= retries; attempt++) {
      const endpoint = this.endpoints[this.currentIndex];
      try {
        const response = await axios.post(`${endpoint}/${method}`, params, {
          timeout: 5000
        });
        return response.data;
      } catch (error) {
        console.warn(`RPC ${endpoint} failed:`, error.message);
        
        // Try next endpoint
        this.currentIndex = (this.currentIndex + 1) % this.endpoints.length;
        
        if (attempt === retries) {
          throw new Error('All RPC endpoints failed');
        }
        
        // Exponential backoff
        await this.sleep(Math.pow(2, attempt) * 250);
      }
    }
  }

  sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  // Helper methods
  async getChainTip() {
    return this.request('chainTip', {});
  }

  async getAccount(address) {
    return this.request(`account/${address}`, {});
  }

  async submitTx(signedTx) {
    return this.request('submitTx', signedTx);
  }
}

// Usage
const rpc = new ArchivasRPC();
const tip = await rpc.getChainTip();
console.log(`Current height: ${tip.height}`);
```

---

## 🔐 Security Considerations

### Rate Limiting
- **Per-IP limits** prevent abuse from single sources
- **Burst allowance** handles legitimate traffic spikes
- **429 responses** inform clients to back off

### Attack Mitigation
- **Basic WAF** blocks common attack patterns (`.env`, `wp-admin`, etc.)
- **Connection limits** prevent resource exhaustion
- **Timeouts** prevent slow-loris attacks

### SSL/TLS
- **Let's Encrypt** certificates (auto-renewed)
- **TLS 1.2+** only
- **Strong ciphers** (no RC4, 3DES, etc.)

---

## 📈 Performance

### Expected Metrics

| Metric | Target | Acceptable |
|--------|--------|------------|
| Cache Hit Ratio | 80-90% | >70% |
| P99 Latency | <100ms | <500ms |
| Throughput | 100+ RPS | 50+ RPS |
| Uptime | 99.9% | 99.5% |
| Error Rate | <0.1% | <1% |

### Capacity Planning

**Current Configuration:**
- Cache size: 1GB
- Connection limit: 50 per IP
- Rate limit: 30 req/s per IP
- File descriptors: 65536

**Estimated Capacity:**
- ~100 RPS sustained (with caching)
- ~1000 concurrent connections
- ~1M requests/day

**Scaling Up:**
- Increase cache size: Edit `max_size` in Nginx config
- Increase rate limits: Edit `limit_req_zone` rate
- Add more relays: Deploy Seed3, Seed4, etc.

---

## 🛠️ Troubleshooting

### Issue: High cache miss ratio (<50%)

**Diagnosis:**
```bash
curl https://seed2.archivas.ai/status | jq '.cache_stats.hit_ratio_5m'
```

**Cause**: Cache TTL too low or cache eviction

**Fix**:
```bash
# Increase cache size
sudo sed -i 's/max_size=1g/max_size=2g/' /etc/nginx/sites-available/archivas-seed2
sudo systemctl reload nginx
```

---

### Issue: Slow responses (>1s)

**Diagnosis:**
```bash
curl -w "\nTime: %{time_total}s\n" https://seed2.archivas.ai/chainTip
```

**Cause**: Upstream Seed1 slow or cache miss

**Fix**:
```bash
# Check upstream health
curl https://seed2.archivas.ai/status | jq '.upstream'

# Check Seed1 directly
curl -w "\nTime: %{time_total}s\n" https://seed.archivas.ai/chainTip
```

---

### Issue: 503 Service Unavailable

**Diagnosis:**
```bash
curl https://seed2.archivas.ai/ready
```

**Cause**: Upstream down or cache disk full

**Fix**:
```bash
# Check upstream
curl https://seed.archivas.ai/chainTip

# Check cache disk space
df -h /var/cache/nginx/archivas_rpc

# Clear cache if needed
sudo rm -rf /var/cache/nginx/archivas_rpc/*
sudo systemctl reload nginx
```

---

### Issue: 429 Too Many Requests

**Diagnosis**:
```bash
# Check your request rate
curl -I https://seed2.archivas.ai/chainTip
# Look for: X-RateLimit-* headers
```

**Cause**: Exceeding 30 req/s per IP

**Fix**:
- Implement client-side rate limiting
- Add delays between requests
- Cache responses client-side
- Request rate limit increase (contact team)

---

## 📚 Additional Resources

- **API Documentation**: https://docs.archivas.ai/api
- **Grafana Dashboard**: `/infra/grafana/seed2.json`
- **Source Code**: https://github.com/ArchivasNetwork/archivas
- **Support**: https://github.com/ArchivasNetwork/archivas/issues

---

## 🎯 Quick Reference

### For Explorers
```
Primary: https://seed2.archivas.ai
Fallback: https://seed.archivas.ai
```

### For Wallets
```
Reads: https://seed2.archivas.ai
Writes: https://seed2.archivas.ai/submitTx
```

### For Farmers
```
P2P: seed.archivas.ai:9090 (Seed1 ONLY)
RPC: https://seed2.archivas.ai (balance checks only)
```

### For Validators
```
P2P Peers: seed.archivas.ai:9090, OTHER_PEERS
RPC: http://localhost:8080 (own node)
```

---

**Version**: 2.0  
**Last Updated**: 2025-11-13  
**Maintained By**: Archivas Network Team

