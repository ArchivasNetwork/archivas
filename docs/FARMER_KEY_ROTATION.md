# Farmer Key Rotation Guide

This guide explains how to safely rotate your Archivas farmer private keys, especially if they have been compromised or posted publicly.

---

## Why Rotate Keys?

**Any private key that has been exposed should be considered compromised**, even on testnet. Reasons to rotate:

- ✅ Posted in chat, logs, or screenshots
- ✅ Shared in documentation or GitHub issues
- ✅ Stored on compromised systems
- ✅ Regular security hygiene (rotate every 3-6 months)
- ✅ Transitioning from test to production farming

**Even on Betanet testnet, exposed keys allow anyone to:**
- Steal your farming rewards
- Sign blocks on your behalf
- Transfer your RCHV balance

---

## Step-by-Step Rotation Process

### 1. Generate New Keypair

Use the Archivas wallet CLI to generate a new Ethereum-compatible keypair:

```bash
archivas-wallet rotate-farmer-key
```

**Expected output:**
```
╔════════════════════════════════════════════════════════════╗
║          ARCHIVAS FARMER KEY ROTATION WIZARD               ║
╚════════════════════════════════════════════════════════════╝

Step 1: Generating new farmer keypair...
✅ New farmer keypair generated!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
NEW FARMER CREDENTIALS (SAVE SECURELY):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ARCV Address: arcv1...
EVM Address:  0x...
Public Key:   02...
Private Key:  abc123...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**🔐 CRITICAL: Save the private key immediately in a secure, encrypted location.**

---

### 2. Update Farmer Service Configuration

#### For Systemd Services (Seed Nodes)

Edit your farmer service file:

```bash
sudo nano /etc/systemd/system/archivas-betanet-farmer.service
```

Update the `ExecStart` line with the **new private key**:

```ini
[Service]
ExecStart=/usr/local/bin/archivas-farmer farm \
    --plots /var/lib/archivas/plots \
    --node http://localhost:8545 \
    --farmer-privkey <NEW_PRIVATE_KEY_HERE>
```

Reload and restart the service:

```bash
sudo systemctl daemon-reload
sudo systemctl restart archivas-betanet-farmer
sudo systemctl status archivas-betanet-farmer
```

#### For Manual/Screen Farmers

Stop the current farmer process:

```bash
# If using screen:
screen -r farmer
# Press Ctrl+C to stop
exit

# If using nohup:
pkill -f archivas-farmer
```

Start the farmer with the new key:

```bash
screen -S farmer
archivas-farmer farm \
    --plots ./plots \
    --node http://localhost:8545 \
    --farmer-privkey <NEW_PRIVATE_KEY_HERE>
# Press Ctrl+A, then D to detach
```

---

### 3. Transfer Remaining Balance

If your old farmer address has RCHV balance, transfer it to the new address:

```bash
archivas-wallet rotate-farmer-key \
    --old-privkey <OLD_PRIVATE_KEY> \
    --node http://localhost:8545 \
    --broadcast
```

**This will:**
1. Generate a new keypair (as before)
2. Query the old address balance
3. Create a signed transfer transaction (balance minus fee)
4. Broadcast it to the network (if `--broadcast` flag is used)

**Example output:**
```
Step 3: Preparing balance transfer from old to new address...

Old address: arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs
Balance:     125.50000000 RCHV
Nonce:       42

Transfer prepared:
  From:   arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs
  To:     arcv1gl4ykgszns24eq6l6zstn8upjemx7sr28u0rm2
  Amount: 125.49900000 RCHV
  Fee:    0.00100000 RCHV

✅ success: Transaction submitted
⏳ Transaction will be included in the next block (~20 seconds)
```

---

### 4. Verification

#### Verify New Farmer is Running

Check that the farmer is submitting proofs with the new address:

```bash
# Check farmer logs
sudo journalctl -u archivas-betanet-farmer -f

# Or for screen/nohup:
screen -r farmer  # Press Ctrl+A, D to detach
tail -f nohup.out
```

**Look for:**
```
[FARMER] Using address: arcv1gl4ykgszns24eq6l6zstn8upjemx7sr28u0rm2 (EVM: 0x47ea4b22029c155c835fd0a0b99f8196766f406a)
[FARMER] Submitted proof for challenge height 12345
```

#### Verify Both Address Formats

Confirm your ARCV and EVM addresses match:

```bash
# Using the address converter CLI
cd ~/archivas
go build -o address-converter cmd/address-converter/main.go
./address-converter arcv1gl4ykgszns24eq6l6zstn8upjemx7sr28u0rm2
```

**Expected:**
```
╔════════════════════════════════════════════════════════════╗
║        Archivas Address Converter                         ║
╠════════════════════════════════════════════════════════════╣
║ Input (Bech32):  arcv1gl4ykgszns24eq6l6zstn8upjemx7sr28u0rm2
║ Output (0x):     0x47ea4b22029c155c835fd0a0b99f8196766f406a
╚════════════════════════════════════════════════════════════╝
```

#### Verify in MetaMask

Import your **new private key** to MetaMask:

1. Open MetaMask
2. Click account icon > "Import Account"
3. Paste the new private key
4. Verify the 0x address matches the output above
5. Check that farming rewards appear in MetaMask

#### Check Balance via RPC

```bash
# Query new ARCV address
curl -s http://localhost:8545/balance/arcv1gl4ykgszns24eq6l6zstn8upjemx7sr28u0rm2 | jq

# Query new EVM address (should show same balance)
curl -s http://localhost:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0x47ea4b22029c155c835fd0a0b99f8196766f406a", "latest"],
    "id": 1
  }' | jq
```

Both queries should return the same balance (accounting for decimal conversion: ARCV 8 decimals vs ETH 18 decimals).

---

### 5. Secure Old Key Disposal

**After confirming the new key works and balance is transferred:**

1. ✅ **Delete old private key from all config files**
2. ✅ **Remove from bash history:** `history -c && history -w`
3. ✅ **Clear old systemd env files**
4. ✅ **Overwrite old key in notes/screenshots**
5. ✅ **Update any documentation or backup scripts**

**Do NOT keep the old key "just in case" once it has been posted publicly.**

---

## Security Best Practices

### ❌ NEVER Do This

- ❌ Post private keys in public GitHub issues or Discord
- ❌ Include keys in screenshots or screen recordings
- ❌ Store keys in plain text files in cloud storage
- ❌ Email private keys without encryption
- ❌ Commit keys to Git repositories
- ❌ Share keys via unencrypted messaging

### ✅ ALWAYS Do This

- ✅ Store keys in encrypted password managers (1Password, Bitwarden)
- ✅ Use hardware wallets for mainnet farming (when available)
- ✅ Back up keys offline (encrypted USB, paper wallet)
- ✅ Use environment variables or systemd secrets for services
- ✅ Rotate keys immediately if exposure is suspected
- ✅ Test key rotation on testnet before using on mainnet
- ✅ Keep backups in multiple secure locations

---

## Troubleshooting

### "Error querying old address balance"

**Cause:** Node is not running or unreachable.

**Solution:**
```bash
# Check node is running
sudo systemctl status archivas-betanet

# Or check RPC manually
curl http://localhost:8545
```

### "Balance too low to cover transfer fee"

**Cause:** Old address has less than 0.001 RCHV.

**Solution:** No transfer needed. The balance is negligible. Just start using the new key.

### "Failed to recover public key from signature"

**Cause:** Private key format is incorrect.

**Solution:** Ensure private key is **64-character hex string** (32 bytes), not base64 or other format:

```bash
# Correct: 64 hex chars
1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899

# Wrong: includes 0x prefix
0x1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899

# Wrong: base64
HLenovxxsPzxQvfE29ARApcbyYkt65G63kcsFThD9ImQ==
```

### New Address Not Showing Rewards

**Cause:** Farmer is still using old key or plots are signed with old key.

**Solution:**
1. Verify `archivas-betanet-farmer.service` has correct `--farmer-privkey`
2. Restart farmer: `sudo systemctl restart archivas-betanet-farmer`
3. Check logs for "Using address: arcv1..." message
4. Rewards go to the farmer address **specified when plots were created**. If you want rewards to go to the new address, you may need to re-plot with the new key (for future plots).

---

## Quick Reference Commands

```bash
# Generate new key
archivas-wallet rotate-farmer-key

# Generate new key + transfer old balance
archivas-wallet rotate-farmer-key --old-privkey <OLD_KEY> --broadcast

# Update systemd farmer
sudo nano /etc/systemd/system/archivas-betanet-farmer.service
sudo systemctl daemon-reload
sudo systemctl restart archivas-betanet-farmer

# Verify farmer is running
sudo journalctl -u archivas-betanet-farmer -f

# Check new address balance
curl http://localhost:8545/balance/arcv1<NEW_ADDRESS> | jq

# Convert between address formats
./address-converter arcv1<ADDRESS>
./address-converter 0x<ADDRESS>
```

---

## Additional Resources

- [MetaMask Integration Guide](./METAMASK_GUIDE.md)
- [Running a Betanet Node](./gitbook/betanet/running-betanet-node.md)
- [Farming with Private Node](./gitbook/betanet/farming-with-private-node.md)
- [Address System Documentation](../address/README.md)

---

**Questions or issues? Open a GitHub issue or ask in the Archivas Discord.**

