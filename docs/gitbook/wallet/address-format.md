# Address Format (Bech32)

Archivas addresses use Bech32 encoding.

## Format

`arcv1` + bech32_encode(blake2b_160(pubkey))

Example: `arcv1t3huuyd08er3yfnmk9c935rmx3wdh5j6m2uc9d`

## Validation

```typescript
import { isValidAddress } from '@archivas/sdk';
if (!isValidAddress(addr)) throw new Error('Invalid');
```

See: [Wallet Primitives](primitives.md)
