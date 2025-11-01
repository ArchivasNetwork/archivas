# Wallet Primitives (v1.1.0)

This document describes the canonical, frozen wallet interface for Archivas. **These specifications must never change** to ensure SDK and desktop app compatibility.

## Cryptographic Decisions

### Curve: Ed25519
- **Algorithm**: Ed25519 (Edwards-curve Digital Signature Algorithm)
- **Rationale**: Deterministic, fast, battle-tested
- **Key sizes**: 64-byte private key, 32-byte public key

### Mnemonic: BIP39, 24 words
- **Standard**: BIP39
- **Word count**: 24 words (256 bits entropy)
- **Language**: English wordlist

### Key Derivation: SLIP-0010 with Ed25519
- **Path format**: `m/734'/0'/0'/0/0`
  - `734'` = Archivas coin type (hardened, reserved)
  - `0'` = Account (hardened)
  - `0'` = Change (hardened, typically 0 for external)
  - `0` = Address index (hardened)
- **Derivation**: SLIP-0010 (HMAC-SHA512-based key derivation)

## Address Format

### Encoding: Bech32
- **HRP (Human Readable Part)**: `arcv`
- **Data**: blake2b-160(pubkey) = 20 bytes
- **Example**: `arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7`

### Address Generation Process
1. Derive ed25519 keypair from mnemonic at path `m/734'/0'/0'/0/0`
2. Hash the 32-byte public key with blake2b-160 → 20 bytes
3. Encode the 20-byte hash as Bech32 with HRP `arcv`

## Transaction Format (TxV1)

### Schema
```json
{
  "type": "transfer",
  "from": "<arcv...>",
  "to": "<arcv...>",
  "amount": "<u64>",
  "fee": "<u64>",
  "nonce": "<u64>",
  "memo": "<optional UTF-8, <=256 bytes>"
}
```

### Field Descriptions
- `type`: Always `"transfer"` (future: `"register_farm"` reserved)
- `from`: Sender address (Bech32, `arcv...`)
- `to`: Recipient address (Bech32, `arcv...`)
- `amount`: Transfer amount in **base units** (u64, 8 decimals for RCHV)
- `fee`: Transaction fee in **base units** (u64)
- `nonce`: Sender's nonce (u64), must match current account nonce
- `memo`: Optional UTF-8 string, maximum 256 bytes

### Serialization: RFC 8785 Canonical JSON
- **Rule**: Keys sorted lexicographically
- **Rule**: No whitespace
- **Rule**: UTF-8 encoding
- **Rule**: No duplicate keys

Example canonical JSON:
```json
{"amount":"1000000000","fee":"100","from":"arcv1...","memo":"","nonce":"0","to":"arcv1...","type":"transfer"}
```

## Transaction Hashing

### Process
1. Serialize transaction to canonical JSON (RFC 8785)
2. Prepend domain separator: `"Archivas-TxV1"`
3. Hash: `blake2b-256(domain_separator || canonical_json_bytes)`
4. Result: 32-byte transaction hash

### Domain Separation
- **Constant**: `"Archivas-TxV1"`
- **Purpose**: Prevent cross-protocol attacks

## Transaction Signing

### Process
1. Validate transaction (amount > 0, fee > 0, valid addresses, memo length)
2. Compute transaction hash (see above)
3. Sign hash with ed25519: `ed25519.Sign(private_key, hash)`
4. Result: 64-byte signature

### Wire Format (Signed Transaction)
```json
{
  "tx": { ...canonical TxV1... },
  "pubkey": "<base64, 32 bytes>",
  "sig": "<base64, 64 bytes>",
  "hash": "<hex, 32 bytes>"
}
```

### Signature Verification
1. Decode public key (base64 or hex, 32 bytes)
2. Decode signature (base64, 64 bytes)
3. Recompute transaction hash
4. Verify: `ed25519.Verify(public_key, hash, signature)`

## Units

### Base Units
- All API amounts are in **base units** (uint64)
- RCHV has 8 decimals: `1 RCHV = 100,000,000 base units`
- Client UIs format for display; APIs never return floats

### Example Conversions
- `1.5 RCHV` → `150000000` base units
- `0.00000001 RCHV` → `1` base unit

## API Endpoints

See `specs/api-wallet-v1.md` for complete endpoint specifications.

### Quick Reference
- `GET /account/<addr>` → Account balance and nonce
- `GET /chainTip` → Current chain tip (height, hash, difficulty)
- `GET /mempool` → Array of pending transaction hashes
- `GET /tx/<hash>` → Transaction details and confirmation status
- `GET /estimateFee?bytes=<n>` → Estimated fee for transaction size
- `POST /submit` → Submit signed transaction (64 KB max body)

## Constants (Immutable)

These values **must never change**:

- **Derivation Path**: `m/734'/0'/0'/0/0`
- **Domain Separator**: `"Archivas-TxV1"`
- **Bech32 HRP**: `"arcv"`
- **Mnemonic Words**: 24 (BIP39, 256 bits)
- **Memo Max Length**: 256 bytes

## Testing

See `pkg/crypto/` and `pkg/tx/v1/` for golden test vectors and invariant tests.

### Invariant Tests
- `TestDerivationPathConstant`: Path must be `m/734'/0'/0'/0/0`
- `TestDomainSeparatorConstant`: Domain must be `"Archivas-TxV1"`
- `TestBech32HRPConstant`: HRP must be `"arcv"`
- `TestRoutesExist`: All six wallet API routes must exist

## CLI Tool

The `archivas-cli` tool provides:
- `keygen`: Generate new mnemonic and address
- `addr <mnemonic>`: Derive address from mnemonic
- `sign-transfer`: Sign a transfer transaction
- `broadcast`: Submit signed transaction to node

Example:
```bash
# Generate new wallet
archivas-cli keygen

# Sign transaction
archivas-cli sign-transfer \
  --from-mnemonic "word1 ... word24" \
  --to arcv1... \
  --amount 1000000000 \
  --fee 100 \
  --nonce 0 \
  --out tx.json

# Broadcast
archivas-cli broadcast tx.json http://localhost:8080
```

