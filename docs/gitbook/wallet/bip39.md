# BIP39 Mnemonics

24-word phrases for wallet recovery.

## Generation

```typescript
import { Derivation } from '@archivas/sdk';
const mnemonic = Derivation.mnemonicGenerate();
```

## Security

- Never share your mnemonic
- Write it down offline
- Store in safe place
- Test recovery before relying on it

See: [Wallet Primitives](primitives.md)
