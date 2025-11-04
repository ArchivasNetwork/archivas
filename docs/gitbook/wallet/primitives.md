# Wallet Primitives

Archivas uses standard cryptographic primitives for wallets.

## Key Components

- **Mnemonic:** BIP39 (24 words)
- **Key Derivation:** SLIP-0010 for Ed25519
- **Signatures:** Ed25519 (64 bytes)
- **Addresses:** Bech32 (`arcv` prefix)
- **Hashing:** Blake2b-256

Full spec: [wallet-primitives.md](../../wallet-primitives.md)
