# Transaction Signing

How Archivas transactions are signed and verified.

## Signing Process

1. **Build transaction** (canonical JSON)
2. **Hash** (Blake2b-256 with domain separation)
3. **Sign hash** (Ed25519)
4. **Create SignedTx** structure
5. **Submit** to node

## Code Example

```typescript
import { Tx } from '@archivas/sdk';

const tx = Tx.buildTransfer({
  from: 'arcv1...',
  to: 'arcv1...',
  amount: '10000000000',
  fee: '100000',
  nonce: '0'
});

const signedTx = await Tx.createSigned(tx, secretKey);

// signedTx contains:
// - tx: original transaction
// - pubkey: public key (hex)
// - sig: signature (hex)
// - hash: transaction hash (hex)
```

## Verification

Node verifies:
1. Public key matches `from` address
2. Signature is valid Ed25519
3. Hash matches canonical JSON
4. Nonce matches account state
5. Balance sufficient

All verification happens server-side.

---

**Next:** [Best Practices](best-practices.md)
