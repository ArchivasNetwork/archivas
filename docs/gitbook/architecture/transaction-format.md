# Transaction Format

Archivas transaction structure and serialization.

## Transaction Structure

```json
{
  "type": "transfer",
  "from": "arcv1...",
  "to": "arcv1...",
  "amount": "30000000000",
  "fee": "100000",
  "nonce": "1",
  "memo": "optional"
}
```

## Fields

- `type`: "transfer" | "coinbase"
- `from`: Sender address (Bech32)
- `to`: Recipient address (Bech32)
- `amount`: Amount in base units (string)
- `fee`: Fee in base units (string)
- `nonce`: Sequence number (string)
- `memo`: Optional message (max 256 bytes)

## Serialization

**Method:** RFC 8785 Canonical JSON
- Deterministic key ordering
- No whitespace
- Consistent encoding

**Hash:** Blake2b-256(`Archivas-TxV1` || canonical_bytes)

## Signature

**Algorithm:** Ed25519
- 64-byte signature
- Signs transaction hash
- Verified with public key

---

**Next:** [Network Protocol](network-protocol.md)
