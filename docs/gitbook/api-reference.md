# Public RPC API Reference

**Base URL:** `https://seed.archivas.ai`  
**Version:** v1.2.0  
**Protocol:** HTTPS (HTTP/2)  
**Format:** JSON  

---

## Authentication

No authentication required. All endpoints are publicly accessible.

**Rate Limits:**
- `/submit`: 10 requests/minute per IP (burst: 5)
- Other endpoints: No rate limit

---

## Endpoints

### Chain Information

#### GET /chainTip

Get current blockchain tip.

**Response:**
```json
{
  "height": "64000",
  "hash": "0x...",
  "difficulty": "1000000"
}
```

All numeric fields are strings (u64).

---

#### GET /genesisHash

Get genesis block hash.

**Response:**
```json
{
  "genesisHash": "de7ad6cff236a2aae89bca258445b8dc5ea390339a5af75d8492adac8a1abc84"
}
```

---

### Accounts

#### GET /account/\<address\>

Get account balance and nonce.

**Parameters:**
- `address` (path) - Bech32 address (`arcv1...`)

**Response:**
```json
{
  "address": "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7",
  "balance": "137540000000000",
  "nonce": "0"
}
```

**Fields:**
- `balance` - Balance in base units (8 decimals, so divide by 100,000,000 for RCHV)
- `nonce` - Transaction count / sequence number

---

### Blocks

#### GET /blocks/recent?limit=N

List recent blocks (v1.2.0).

**Parameters:**
- `limit` (query, optional) - Number of blocks to return (default: 20, max: 100)

**Response:**
```json
{
  "blocks": [
    {
      "height": "64000",
      "hash": "0x...",
      "timestamp": "1762210000",
      "farmer": "arcv1...",
      "miner": "arcv1...",
      "txCount": "1",
      "difficulty": "1000000"
    }
  ]
}
```

**Fields:**
- `farmer` - Primary field (block producer address)
- `miner` - Deprecated alias (remove in v1.3)
- `timestamp` - Unix timestamp

---

#### GET /block/\<height\>

Get full block details by height.

**Parameters:**
- `height` (path) - Block height (u64)

**Response:**
```json
{
  "height": "64000",
  "hash": "0x...",
  "prevHash": "0x...",
  "timestamp": "1762210000",
  "difficulty": "1000000",
  "challenge": "0x...",
  "farmer": "arcv1...",
  "farmerAddr": "arcv1...",
  "txCount": "2",
  "txs": [
    {
      "type": "coinbase",
      "from": "coinbase",
      "to": "arcv1...",
      "amount": 2000000000,
      "fee": 0,
      "nonce": 0
    },
    {
      "type": "transfer",
      "from": "arcv1...",
      "to": "arcv1...",
      "amount": 30000000000,
      "fee": 100000,
      "nonce": 1
    }
  ]
}
```

**Transaction Types:**
- `coinbase` - Block reward to farmer
- `transfer` - User-initiated transfer

---

### Transactions

#### GET /tx/\<hash\>

Get transaction details by hash.

**Parameters:**
- `hash` (path) - Transaction hash (hex)

**Response:**
```json
{
  "hash": "0x...",
  "tx": {
    "type": "transfer",
    "from": "arcv1...",
    "to": "arcv1...",
    "amount": "30000000000",
    "fee": "100000",
    "nonce": "1"
  },
  "block": {
    "height": 64000,
    "hash": "0x...",
    "timestamp": 1762210000
  }
}
```

**Note:** Currently returns placeholder data. Full implementation coming in v1.3.

---

#### GET /tx/recent?limit=N

List recent transactions (v1.2.0).

**Parameters:**
- `limit` (query, optional) - Number of transactions (default: 50, max: 200)

**Response:**
```json
{
  "txs": [
    {
      "hash": "0x...",
      "type": "transfer",
      "from": "arcv1...",
      "to": "arcv1...",
      "amount": "30000000000",
      "fee": "100000",
      "nonce": "1",
      "height": "64000",
      "timestamp": "1762210000"
    }
  ]
}
```

**Note:** Only includes user transfers (coinbase transactions filtered out).

---

#### GET /mempool

List pending transactions.

**Response:**
```json
[]
```

Returns array of pending transaction hashes (currently empty array).

---

#### POST /submit

Submit a signed transaction (v1.1.0).

**Headers:**
- `Content-Type: application/json` (required)

**Request Body:**
```json
{
  "tx": {
    "type": "transfer",
    "from": "arcv1...",
    "to": "arcv1...",
    "amount": "30000000000",
    "fee": "100000",
    "nonce": "1"
  },
  "pubkey": "0x...",
  "sig": "0x...",
  "hash": "0x..."
}
```

**Response (Success):**
```json
{
  "ok": true,
  "hash": "0x..."
}
```

**Response (Error):**
```json
{
  "ok": false,
  "error": "Invalid signature"
}
```

**Validation:**
- Signature must be valid Ed25519
- Public key must match `from` address
- Nonce must match current account nonce
- Balance must be sufficient (amount + fee)

---

### Utilities

#### GET /estimateFee?bytes=\<n\>

Estimate transaction fee.

**Parameters:**
- `bytes` (query) - Transaction size in bytes

**Response:**
```json
{
  "fee": "100000"
}
```

Currently returns fixed fee (0.001 RCHV). Dynamic fee market coming in future version.

---

## CORS

All endpoints include CORS headers for browser access:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
Access-Control-Max-Age: 86400
```

**OPTIONS requests:** Return 204 No Content with CORS headers.

---

## Error Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 204 | No Content (OPTIONS) |
| 400 | Bad Request (invalid parameters) |
| 404 | Not Found (block/tx doesn't exist) |
| 405 | Method Not Allowed |
| 415 | Unsupported Media Type (missing Content-Type) |
| 429 | Too Many Requests (rate limited) |
| 500 | Internal Server Error |
| 502 | Bad Gateway (node temporarily unavailable) |

---

## Best Practices

### For Web Apps
```typescript
import { createRpcClient } from '@archivas/sdk';

const rpc = createRpcClient({
  baseUrl: 'https://seed.archivas.ai',
  timeout: 30000
});

// Get chain info
const tip = await rpc.getChainTip();

// Get account
const account = await rpc.getAccount('arcv1...');
```

### For Direct HTTP
```bash
# Always use HTTPS
curl https://seed.archivas.ai/chainTip

# Set Content-Type for POST
curl -X POST https://seed.archivas.ai/submit \
  -H "Content-Type: application/json" \
  -d '{"tx":{...},"pubkey":"...","sig":"...","hash":"..."}'

# Handle rate limits gracefully
# Implement exponential backoff if you get 429
```

---

## Versioning

**Current:** v1.2.0

**Compatibility:**
- v1.2.0: Explorer listing endpoints added
- v1.1.0: Wallet API frozen (backward compatible)
- v1.0.0: Initial public RPC

All v1.x versions are backward compatible. Breaking changes will increment to v2.0.

---

## Support

**Issues:** https://github.com/ArchivasNetwork/archivas/issues  
**Documentation:** https://github.com/ArchivasNetwork/archivas/tree/main/docs  

For infrastructure issues with seed.archivas.ai:
- Check [docs/SEED_HOST.md](../SEED_HOST.md)
- File a GitHub issue

