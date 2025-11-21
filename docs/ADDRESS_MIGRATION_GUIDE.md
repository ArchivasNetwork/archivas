# Address Derivation Migration Guide

## Summary

Archivas has migrated from a custom address derivation system to the **Ethereum standard** for full MetaMask compatibility.

---

## What Changed?

### Before (Old System)
- **Derivation:** Private key → secp256k1 public key (compressed) → **SHA256** → first 20 bytes
- **Result:** Addresses were NOT compatible with MetaMask
- **Issue:** Same private key produced DIFFERENT addresses in Archivas vs. MetaMask

### After (New System)
- **Derivation:** Private key → secp256k1 public key (uncompressed) → **Keccak256** → last 20 bytes
- **Result:** Addresses are fully compatible with MetaMask
- **Benefit:** Same private key produces IDENTICAL addresses in all tools

---

## Unified Derivation

All Archivas tools now use **ONE** canonical function:

```go
address.PrivateKeyToEVMAddress(privateKeyBytes)
```

This function:
1. Takes a 32-byte private key
2. Derives the ECDSA public key (secp256k1)
3. Computes Keccak256 of the uncompressed public key
4. Returns the last 20 bytes as the EVM address

---

## Impact on Users

### ✅ New Wallets (Generated After Update)
- **Fully compatible** with MetaMask
- Import private key → MetaMask shows correct address
- Farm rewards appear in MetaMask automatically
- Use with any Ethereum tool (ethers.js, Hardhat, etc.)

### ⚠️ Old Wallets (Generated Before Update)
- **NOT compatible** with MetaMask by default
- Old private key derives to DIFFERENT address in MetaMask
- Farm rewards went to old Archivas-derived address
- **Solution:** Generate new Ethereum-compatible wallet

---

## How to Migrate

### Option 1: Generate New Wallet (Recommended)

```bash
cd ~/archivas
go build -o archivas-wallet cmd/archivas-wallet/*.go
./archivas-wallet new
```

**Output:**
```
ARCV Address (Bech32): arcv1...
EVM Address (Hex):     0x...
Private Key:           abc123...
```

✅ This wallet is fully MetaMask-compatible!

### Option 2: Use eth-wallet Tool

```bash
cd ~/archivas
go build -o eth-wallet cmd/eth-wallet/main.go
./eth-wallet
```

**Output:**
```
Private Key: def456...
Ethereum Address: 0x...
```

Import this private key to MetaMask → works perfectly!

### Option 3: Transfer Old Balance to New Wallet

If you have RCHV in an old wallet:

1. Generate new Ethereum-compatible wallet
2. Note the new ARCV address
3. Use old wallet to send RCHV to new address:

```bash
./archivas-wallet send \
  --from-privkey OLD_PRIVATE_KEY \
  --to NEW_ARCV_ADDRESS \
  --amount 3048000000000 \
  --node https://seed3.betanet.archivas.ai
```

---

## For Farmers

### Update Your Farmer

If you're currently farming with an old wallet:

1. **Generate new Ethereum-compatible wallet** (see above)
2. **Update farmer private key:**

```bash
# Stop farmer
sudo systemctl stop archivas-betanet-farmer

# Edit farmer service
sudo nano /etc/systemd/system/archivas-betanet-farmer.service

# Replace --farmer-privkey with NEW private key

# Restart farmer
sudo systemctl daemon-reload
sudo systemctl start archivas-betanet-farmer
```

3. **Verify new address:**

```bash
sudo journalctl -u archivas-betanet-farmer -f
```

Should show: `👨‍🌾 Farmer Address: arcv1...` (new Ethereum-compatible address)

✅ **Future rewards will be MetaMask-compatible!**

---

## Backward Compatibility

The WorldState has been updated to support **both old and new addresses**:

- **Old ARCV-keyed accounts** (from before update) remain accessible
- **New EVM-keyed accounts** (from after update) work perfectly
- **Querying works with both formats:**
  - Query with `arcv1...` → finds account
  - Query with `0x...` → finds same account
  - Both formats reference the same balance

### Example:

```bash
# Old ARCV address (legacy derivation)
curl https://seed3.betanet.archivas.ai/balance/arcv1rlx93wlk26ny67zqk8eejfkl4y2az22nynqrtj

# New EVM address (Ethereum derivation) - for a DIFFERENT account
curl https://seed3.betanet.archivas.ai/balance/0x39a028dfdcae40bf277ec1ec268d62665d36c073

# Both return correct balances!
```

---

## Testing Address Derivation

Verify that your tools are using the unified system:

```bash
cd ~/archivas

# Run address tests
go test ./address/... -v

# Check wallet generates Ethereum addresses
./archivas-wallet new

# Verify farmer uses same derivation
# (Check farmer logs for derived address)
```

**Expected:** Wallet and farmer produce **identical** addresses for the same private key.

---

## Technical Details

### Derivation Comparison

| Method | Hash Function | Public Key Format | Bytes Used | MetaMask Compatible |
|--------|---------------|-------------------|------------|-------------------|
| **Old (SHA256)** | SHA-256 | Compressed (33 bytes) | First 20 | ❌ No |
| **New (Keccak256)** | Keccak-256 | Uncompressed (64 bytes) | Last 20 | ✅ Yes |

### Code Locations

- **Canonical derivation:** `address/derivation.go`
- **Tests:** `address/derivation_test.go`
- **Wallet CLI:** `cmd/archivas-wallet/main.go`
- **Farmer:** `cmd/archivas-farmer/main.go`
- **WorldState helpers:** `ledger/address_helpers.go`
- **Backward compatibility:** `ledger/state.go`

---

## Troubleshooting

### Issue: MetaMask shows 0 balance, but Archivas shows balance

**Cause:** Your RCHV is on an old Archivas-derived address, not the Ethereum-compatible one.

**Solution:**
1. Generate new Ethereum-compatible wallet
2. Transfer balance from old wallet to new
3. Import new wallet to MetaMask

### Issue: Farmer rewards not showing in MetaMask

**Cause:** Farmer is using old private key, farming to old address.

**Solution:**
1. Generate new Ethereum-compatible wallet
2. Update farmer service with new private key
3. Restart farmer
4. Future rewards will be MetaMask-compatible

### Issue: Same private key, different addresses in wallet vs. MetaMask

**Cause:** Wallet was generated before the update (old derivation).

**Solution:** Generate new wallet with updated `archivas-wallet` tool.

---

## Summary

✅ **New System Benefits:**
- Full MetaMask compatibility
- Same address across all tools
- Standard Ethereum derivation
- Works with any Ethereum tooling

⚠️ **Migration Required:**
- Old wallets need new generation
- Farmers should update private keys
- Transfers from old → new address if needed

📚 **Resources:**
- MetaMask Guide: `docs/METAMASK_QUICK_START.md`
- MetaMask One-Pager: `docs/METAMASK_ONE_PAGER.md`
- Full Documentation: `docs/METAMASK_GUIDE.md`

---

**Questions?** Open an issue on GitHub or check the documentation.

