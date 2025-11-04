# Transaction Signing

Ed25519 signatures on canonical JSON.

## Process

1. Build transaction (canonical JSON)
2. Hash with Blake2b-256
3. Add domain separation: `Archivas-TxV1`
4. Sign hash with Ed25519
5. Return 64-byte signature

## Code

```typescript
const signedTx = await Tx.createSigned(tx, secretKey);
```

See: [Transaction Signing](../developers/transaction-signing.md)
