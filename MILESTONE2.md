# Milestone 2 Complete: Cryptographic Ownership

## 🎉 RCHV is Now Owned!

Archivas now has **real cryptographic ownership**. Only the holder of a private key can spend their RCHV.

## What's New

### ✅ Wallet Package (`wallet/`)

- **Keypair Generation**: secp256k1 (same as Ethereum)
- **Address Format**: bech32 with `arcv` prefix (e.g., `arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7`)
- **Transaction Signing**: ECDSA signatures over transaction hash
- **Signature Verification**: Validates that SenderPubKey matches From address and signature is valid

### ✅ Wallet CLI (`cmd/archivas-wallet`)

Two commands:

1. **`new`** - Generate a new wallet
   ```bash
   go run ./cmd/archivas-wallet new
   ```
   
   Output:
   ```
   🔐 New Archivas Wallet Generated
   
Address:     arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
Public Key:  031b2cb8df7a31f463b72468e731786df04fda65b1260e7e77fc96d11b5ec69e97
Private Key: <EXAMPLE_PRIVATE_KEY_DO_NOT_USE>
   
   ⚠️  KEEP YOUR PRIVATE KEY SECRET!
   ```

2. **`send`** - Sign and submit a transaction
   ```bash
   go run ./cmd/archivas-wallet send \
     --from-privkey <hex> \
     --to <arcv1...> \
     --amount <base_units> \
     --fee <base_units>
   ```
   
   The wallet automatically:
   - Derives your address from your private key
   - Queries the node for your current nonce
   - Builds and signs the transaction
   - Submits it to the node

### ✅ Transaction Updates

Transactions now include:

```go
type Transaction struct {
    From         string // bech32 address
    To           string // bech32 address
    Amount       int64  // base units
    Fee          int64  // base units
    Nonce        uint64 // must match sender's current nonce
    SenderPubKey []byte // sender's public key
    Signature    []byte // secp256k1 ECDSA signature
}
```

### ✅ Signature Validation

**On Transaction Submission (RPC):**
- Verifies signature is valid
- Checks SenderPubKey matches From address
- Validates balance and nonce
- Rejects invalid transactions before they enter mempool

**On Block Production (Ledger):**
- Re-validates signature when applying transaction
- Ensures only valid, signed transactions modify state

## Complete End-to-End Test

### 1. Generate Wallets

```bash
# Wallet A (will be funded)
go run ./cmd/archivas-wallet new
# Address:     arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
# Private Key: <EXAMPLE_KEY_REPLACED_FOR_SECURITY>

# Wallet B (receiver)
go run ./cmd/archivas-wallet new
# Address:     arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt
```

### 2. Fund Wallet A in Genesis

Edit `config/params_devnet.go`:

```go
var GenesisAlloc = map[string]int64{
    "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7": 1_000_000_000_00000000, // 1B RCHV
}
```

### 3. Start the Node

```bash
go build -o archivas-node ./cmd/archivas-node
./archivas-node
```

Output:
```
🌍 World state initialized with 1 genesis accounts
   arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7: 1000000000.00000000 RCHV
🌐 Starting RPC server on :8080
```

### 4. Send RCHV from A → B

```bash
go run ./cmd/archivas-wallet send \
  --from-privkey <YOUR_PRIVATE_KEY_HERE> \
  --to arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt \
  --amount 50000000000 \
  --fee 100000
```

Output:
```
📊 Current balance: 1000000000.00000000 RCHV (nonce: 0)
📝 Sending 500.00000000 RCHV from arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7 to arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt
💸 Fee: 0.00100000 RCHV
✅ success: Transaction added to mempool
⏳ Transaction will be included in the next block (~20 seconds)
```

### 5. Verify Balances (after ~20 seconds)

```bash
# Wallet A
curl --noproxy "*" http://localhost:8080/balance/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
# {"address":"arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7","balance":99999949999900000,"nonce":1}

# Wallet B
curl --noproxy "*" http://localhost:8080/balance/arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt
# {"address":"arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt","balance":50000000000,"nonce":0}
```

**Results:**
- ✅ Wallet A: 999,999,499.999 RCHV (sent 500 + 0.001 fee)
- ✅ Wallet B: 500.00000000 RCHV (received)
- ✅ Wallet A nonce: 0 → 1

### 6. Send Another Transaction (tests nonce increment)

```bash
go run ./cmd/archivas-wallet send \
  --from-privkey <YOUR_PRIVATE_KEY_HERE> \
  --to arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt \
  --amount 25000000000 \
  --fee 150000
```

After next block:
- ✅ Wallet A: 999,999,249.9975 RCHV (nonce: 2)
- ✅ Wallet B: 750.00000000 RCHV

## Security Model

### What's Protected
✅ **Cryptographic Ownership**: Only private key holder can sign transactions  
✅ **Address Derivation**: SenderPubKey must match From address  
✅ **Replay Protection**: Nonce prevents transaction replay  
✅ **Balance Protection**: Cannot spend more than owned  
✅ **Double-Spend Prevention**: State machine enforces single spend per nonce  

### What's Still TODO
⚠️ **No P2P**: Currently single-node (no network consensus yet)  
⚠️ **No PoSpace**: Time-based block production (not Proof-of-Space)  
⚠️ **No VDF**: No Verifiable Delay Function yet  
⚠️ **In-Memory State**: State not persisted to disk  

## Files Changed/Added

### New Files
- `wallet/wallet.go` - Keypair generation, address encoding
- `wallet/tx_sign.go` - Transaction signing and verification
- `ledger/hash.go` - Transaction hashing
- `ledger/verify.go` - Signature verification for ledger
- `cmd/archivas-wallet/main.go` - Wallet CLI

### Updated Files
- `ledger/tx.go` - Added SenderPubKey field
- `ledger/apply_tx.go` - Added signature verification
- `rpc/rpc.go` - Updated submitTx to validate signatures
- `config/params_devnet.go` - Updated to use real bech32 address

## Dependencies Added
- `github.com/decred/dcrd/dcrec/secp256k1/v4` - secp256k1 crypto
- `github.com/btcsuite/btcd/btcutil/bech32` - bech32 address encoding

## Next Milestone (Milestone 3): Proof-of-Space

Now that ownership is cryptographically enforced, we can move to:

1. **Plot Generation**: Farmers create PoSpace plots
2. **Challenge-Response**: Timelord issues challenges
3. **Proof Submission**: Farmers submit PoSpace proofs
4. **Block Selection**: Replace time-based production with PoSpace winner selection

After that: **Milestone 4** (VDF) and **Milestone 5** (Full PoSpace+Time Consensus)

