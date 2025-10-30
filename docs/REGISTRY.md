# Archivas Node Registry

Public node discovery service for the Archivas network.

---

## Overview

The registry allows nodes to:
- Register themselves publicly
- Discover other nodes automatically
- Heartbeat to prove liveness
- Prevent network fragmentation

---

## API Reference

### POST /register

Register a new node with signed proof.

**Request:**
```json
{
  "address": "arcv1...",
  "p2pAddr": "host:port",
  "rpcAddr": "host:port",
  "networkId": "archivas-devnet-v3",
  "pubkey": "03...",
  "nonce": 12345,
  "signature": "hex..."
}
```

**Signature:** Sign `H(address|p2pAddr|rpcAddr|networkId|nonce)` with private key.

**Response:**
```json
{"status": "registered"}
```

### POST /heartbeat

Update node status (must register first).

**Request:**
```json
{
  "p2pAddr": "host:port",
  "tipHeight": 1258,
  "peerCount": 3,
  "signature": "hex..."
}
```

**Response:**
```json
{"status": "ok"}
```

### GET /peers

List active P2P addresses for bootstrapping.

**Response:**
```json
{
  "peers": ["57.129.148.132:9090", "72.251.11.191:9090"],
  "count": 2
}
```

### GET /nodes

Full node list with details.

**Response:**
```json
{
  "nodes": [
    {
      "address": "arcv1...",
      "p2pAddr": "57.129.148.132:9090",
      "rpcAddr": "57.129.148.132:8080",
      "networkId": "archivas-devnet-v3",
      "tipHeight": 1258,
      "peerCount": 2,
      "lastSeen": "2025-10-30T13:00:00Z"
    }
  ],
  "count": 1
}
```

### GET /health

Registry health check.

**Response:**
```json
{
  "ok": true,
  "activeNodes": 2,
  "totalNodes": 2
}
```

---

## Running the Registry

```bash
./archivas-registry --port :8088 --network-id archivas-devnet-v3
```

**Docker:**
```bash
docker run -p 8088:8088 archivas-registry
```

---

## Node Auto-Registration (Future)

Add to `archivas-node`:

```bash
./archivas-node \
  --registry http://registry.archivas.network:8088 \
  --auto-register
```

Node will:
1. Sign registration message with node key
2. POST /register on startup
3. POST /heartbeat every 30s

---

## Security

### Signature Verification

Registry verifies:
- Signature matches pubkey
- Pubkey derives to claimed address
- Network ID matches registry's network

### Anti-Spam

- Unique (address, p2pAddr) pairs
- TTL-based cleanup (2min heartbeat)
- Signature prevents spoofing

### Network Isolation

Nodes from different networks (`networkId`) are rejected.

---

## Client Example (Python)

```python
import hashlib
import requests
from ecdsa import SigningKey, SECP256k1

# Generate signature
privkey = SigningKey.from_string(bytes.fromhex("..."), curve=SECP256k1)
msg = f"{address}|{p2p_addr}|{rpc_addr}|{network_id}|{nonce}"
msg_hash = hashlib.sha256(msg.encode()).digest()
sig = privkey.sign_digest(msg_hash, sigencode=sigencode_der)

# Register
resp = requests.post("http://registry:8088/register", json={
    "address": address,
    "p2pAddr": p2p_addr,
    "rpcAddr": rpc_addr,
    "networkId": network_id,
    "pubkey": pubkey_hex,
    "nonce": nonce,
    "signature": sig.hex()
})
```

---

## Production Deployment

### High Availability

Run multiple registries behind load balancer:
```
Registry A (primary)
Registry B (replica)
→ Load Balancer → Public
```

### Database Backend

For persistence across restarts, use:
- BadgerDB for local storage
- PostgreSQL for multi-instance
- Redis for high-throughput

### Rate Limiting

Add nginx rate limiting:
```nginx
limit_req_zone $binary_remote_addr zone=registry:10m rate=10r/s;

location /register {
    limit_req zone=registry burst=5;
    proxy_pass http://localhost:8088;
}
```

---

## Monitoring

Registry exposes `/health` for monitoring:

**Healthcheck:**
```bash
curl http://registry:8088/health
```

**Monitor:**
- `activeNodes` should match expected count
- If activeNodes drops to 0, registry or network has issues

---

**Your network now has public node discovery!** 📋

For support: https://github.com/ArchivasNetwork/archivas/issues

