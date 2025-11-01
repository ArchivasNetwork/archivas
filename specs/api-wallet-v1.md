# Wallet API v1.1.0 Specification

**Status**: Frozen (no breaking changes allowed)

This document specifies the exact JSON response shapes and error schema for Archivas wallet API endpoints.

## Base URL

All endpoints are served at the node's RPC port (default: `:8080`).

## Common Response Format

### Success Responses
All successful responses return JSON with appropriate `Content-Type: application/json` header.

### Error Responses
Error responses have status codes:
- `400 Bad Request`: Invalid input (malformed address, invalid transaction, etc.)
- `405 Method Not Allowed`: Wrong HTTP method
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error response format:
```json
{
  "ok": false,
  "error": "<human-readable error message>"
}
```

## Endpoints

### 1. GET /account/<address>

**Description**: Get account balance and nonce.

**Path Parameters**:
- `address`: Bech32 address (e.g., `arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7`)

**Response** (200 OK):
```json
{
  "address": "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7",
  "balance": "100000000000000000",
  "nonce": "0"
}
```

**Fields**:
- `address` (string): Account address (Bech32)
- `balance` (string): Account balance in **base units** (u64 as string)
- `nonce` (string): Account nonce (u64 as string)

**Errors**:
- `400`: Invalid address format
- `404`: Account not found (may return balance=0, nonce=0)

---

### 2. GET /chainTip

**Description**: Get current chain tip (height, hash, difficulty).

**Response** (200 OK):
```json
{
  "height": "1234",
  "hash": "0xabcd...",
  "difficulty": "15000000"
}
```

**Fields**:
- `height` (string): Current chain height (u64 as string)
- `hash` (string): Tip block hash (hex, 64 characters, no `0x` prefix required but accepted)
- `difficulty` (string): Current difficulty target (u64 as string)

---

### 3. GET /mempool

**Description**: Get pending transaction hashes in mempool.

**Response** (200 OK):
```json
[
  "0xhash1...",
  "0xhash2...",
  "0xhash3..."
]
```

**Response Format**: Array of transaction hashes (strings, hex-encoded).

**Empty Mempool**: Returns `[]` (empty array).

---

### 4. GET /tx/<hash>

**Description**: Get transaction details and confirmation status.

**Path Parameters**:
- `hash`: Transaction hash (hex, 64 characters)

**Response** (200 OK, confirmed):
```json
{
  "confirmed": true,
  "height": 1234,
  "tx": { ... signed transaction JSON ... }
}
```

**Response** (200 OK, pending):
```json
{
  "confirmed": false,
  "height": null
}
```

**Fields**:
- `confirmed` (boolean): Whether transaction is confirmed in a block
- `height` (number | null): Block height if confirmed, `null` if pending
- `tx` (object, optional): Full signed transaction JSON (see wire format in `wallet-primitives.md`)

**Errors**:
- `400`: Invalid hash format (must be 32 bytes, hex-encoded)
- `404`: Transaction not found in blockchain or mempool

---

### 5. GET /estimateFee?bytes=<n>

**Description**: Estimate transaction fee based on transaction size.

**Query Parameters**:
- `bytes` (optional): Estimated transaction size in bytes (default: 256)

**Response** (200 OK):
```json
{
  "fee": "100"
}
```

**Fields**:
- `fee` (string): Estimated fee in **base units** (u64 as string)

**Algorithm**: Linear estimator (100 base units per KB, minimum 100).

**Example**:
- `bytes=256` → `fee="100"` (minimum)
- `bytes=1024` → `fee="100"` (1 KB = 100)
- `bytes=2048` → `fee="200"` (2 KB = 200)

---

### 6. POST /submit

**Description**: Submit a signed transaction to the mempool.

**Request Body**: Signed transaction JSON (see wire format in `wallet-primitives.md`)

**Request Size Limit**: 64 KB maximum

**Response** (200 OK, accepted):
```json
{
  "ok": true,
  "hash": "0xabcd..."
}
```

**Response** (200 OK, rejected):
```json
{
  "ok": false,
  "error": "Invalid signature"
}
```

**Response** (400 Bad Request):
```json
{
  "ok": false,
  "error": "Invalid request body: ..."
}
```

**Fields**:
- `ok` (boolean): Whether transaction was accepted
- `hash` (string, if `ok=true`): Transaction hash (hex)
- `error` (string, if `ok=false`): Error message

**Validation**:
1. Parse signed transaction JSON
2. Verify signature (`txv1.VerifySignedTx`)
3. Validate transaction fields (amount > 0, fee > 0, valid addresses, memo length)
4. Check nonce matches sender's current nonce
5. Check sender has sufficient balance (amount + fee)
6. Add to mempool

**Errors**:
- `400`: Invalid JSON, invalid signature, validation failure
- `413`: Request body exceeds 64 KB limit

---

## Example Workflow

### 1. Get Account Balance
```bash
curl http://localhost:8080/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
```

Response:
```json
{"address":"arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7","balance":"100000000000000000","nonce":"0"}
```

### 2. Get Chain Tip
```bash
curl http://localhost:8080/chainTip
```

Response:
```json
{"height":"1234","hash":"abcd...","difficulty":"15000000"}
```

### 3. Estimate Fee
```bash
curl http://localhost:8080/estimateFee?bytes=512
```

Response:
```json
{"fee":"100"}
```

### 4. Submit Transaction
```bash
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d @tx.json
```

Request (`tx.json`):
```json
{
  "tx": {
    "type": "transfer",
    "from": "arcv1...",
    "to": "arcv1...",
    "amount": "1000000000",
    "fee": "100",
    "nonce": "0"
  },
  "pubkey": "base64...",
  "sig": "base64...",
  "hash": "0xabcd..."
}
```

Response:
```json
{"ok":true,"hash":"0xabcd..."}
```

### 5. Check Transaction Status
```bash
curl http://localhost:8080/tx/abcd...
```

Response (confirmed):
```json
{"confirmed":true,"height":1235,"tx":{...}}
```

## Rate Limiting

Currently no rate limiting is enforced, but clients should:
- Poll endpoints (e.g., `/chainTip`, `/account`) at reasonable intervals (≥1 second)
- Batch requests when possible
- Handle `429 Too Many Requests` gracefully (if implemented)

## Versioning

This API is frozen at v1.1.0. Future changes must:
- Be additive only (new endpoints, new optional fields)
- Maintain backward compatibility
- Never change response shapes or field types
- Document breaking changes in release notes (if unavoidable)

