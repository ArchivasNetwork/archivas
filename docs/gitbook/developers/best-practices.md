# Developer Best Practices

Guidelines for building robust Archivas applications.

## Error Handling

Always handle RPC errors:
```typescript
try {
  const result = await rpc.submit(signedTx);
  if (!result.ok) {
    console.error('Rejected:', result.error);
  }
} catch (err) {
  if (err.status === 429) {
    // Rate limited - wait and retry
    await sleep(6000);
  }
  throw err;
}
```

## Nonce Management

Track nonces carefully:
```typescript
// Get current nonce before each tx
const account = await rpc.getAccount(address);
tx.nonce = account.nonce;

// After successful submit, nonce increments
// Next tx should use account.nonce + 1
```

## Security

1. Never expose private keys
2. Validate user input
3. Use HTTPS only
4. Rate limit your requests
5. Handle edge cases

## Testing

Test on testnet before mainnet!
- Use https://seed.archivas.ai
- Get test RCHV from faucet
- Verify all flows work

---

**Resources:** [SDK Guide](sdk-guide.md) | [API Reference](api-reference.md)
