# Public Testnet Access

Access the Archivas public testnet without running your own node.

---

## Public Services

### RPC Endpoint
**URL:** `https://seed.archivas.ai`  
**Status:** 🟢 Live  
**Protocol:** HTTPS (HTTP/2)  
**Authentication:** None required  

### Block Explorer
**URL:** https://archivas-explorer-production.up.railway.app  
**Features:**
- Real-time block updates
- Account lookup
- Transaction browsing
- Network statistics

### Metrics Dashboard
**URL:** http://57.129.148.132:3001  
**Login:** admin / admin  
**Panels:**
- Chain height
- Peer count
- Difficulty
- Block production rate

---

## Quick Test

Test the testnet in your terminal:

```bash
# Get current height
curl https://seed.archivas.ai/chainTip

# Get account balance  
curl https://seed.archivas.ai/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7

# List recent blocks
curl https://seed.archivas.ai/blocks/recent?limit=5

# View mempool
curl https://seed.archivas.ai/mempool
```

---

## Network Information

**Network ID:** `archivas-devnet-v4`  
**Genesis Hash:** `de7ad6cff236a2aae89bca258445b8dc5ea390339a5af75d8492adac8a1abc84`  
**Chain ID:** 1616  
**Block Reward:** 20 RCHV  
**Target Block Time:** 20 seconds  

---

## Rate Limits

| Endpoint | Limit |
|----------|-------|
| `/submit` | 10 requests/minute per IP |
| Other endpoints | Unlimited |

---

## Next Steps

- [Get Testnet RCHV](get-rchv.md)
- [Send Your First Transaction](first-transaction.md)
- [Developer Quickstart](../developer-quickstart.md)

