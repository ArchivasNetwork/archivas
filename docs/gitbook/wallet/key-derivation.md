# Key Derivation (SLIP-0010)

Ed25519 hierarchical deterministic key derivation.

## Path

Archivas uses: `m/734'/0'/0'/0/0`

- 734: Archivas coin type
- Hardened derivation for security

## Implementation

```typescript
const kp = await Derivation.fromMnemonic(mnemonic);
```

See: [Wallet Primitives](primitives.md)
