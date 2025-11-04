# Your First Transaction

Send your first RCHV transfer on the Archivas testnet.

---

## Prerequisites

- Testnet RCHV in your wallet
- Your wallet's private key or mnemonic
- Recipient address

**Need RCHV?** See [Get Testnet RCHV](get-rchv.md)

---

## Method 1: Using archivas-wallet (CLI)

### Step 1: Check Your Balance

```bash
curl https://seed.archivas.ai/account/YOUR_ADDRESS
```

Note your current balance and nonce.

### Step 2: Send RCHV

```bash
./archivas-wallet send \
  --from-privkey YOUR_PRIVATE_KEY_HEX \
  --to RECIPIENT_ADDRESS \
  --amount 10000000000 \
  --fee 100000 \
  --node https://seed.archivas.ai
```

**Amount breakdown:**
- `10000000000` = 100 RCHV (8 decimals)
- `100000` = 0.001 RCHV (standard fee)

### Step 3: Wait for Confirmation

Transactions are included in the next block (~20-30 seconds).

### Step 4: Verify

```bash
# Check recipient received
curl https://seed.archivas.ai/account/RECIPIENT_ADDRESS

# Check your balance decreased
curl https://seed.archivas.ai/account/YOUR_ADDRESS

# Your nonce should have incremented
```

---

## Method 2: Using TypeScript SDK

### Step 1: Setup

```typescript
import { Derivation, Tx, createRpcClient } from '@archivas/sdk';

// From your mnemonic
const mnemonic = "your 24 word mnemonic here";
const keyPair = await Derivation.fromMnemonic(mnemonic);
const myAddress = Derivation.toAddress(keyPair.publicKey);

const rpc = createRpcClient({
  baseUrl: 'https://seed.archivas.ai'
});
```

### Step 2: Get Current Nonce

```typescript
const account = await rpc.getAccount(myAddress);
console.log('Balance:', account.balance);
console.log('Nonce:', account.nonce);
```

### Step 3: Build Transaction

```typescript
const tx = Tx.buildTransfer({
  from: myAddress,
  to: 'arcv1recipient...',
  amount: '10000000000',  // 100 RCHV
  fee: '100000',          // 0.001 RCHV
  nonce: account.nonce    // Current nonce
});
```

### Step 4: Sign

```typescript
const signedTx = await Tx.createSigned(tx, keyPair.secretKey);
console.log('Transaction hash:', signedTx.hash);
```

### Step 5: Submit

```typescript
const result = await rpc.submit(signedTx);
if (result.ok) {
  console.log('✅ Transaction submitted!');
  console.log('Hash:', result.hash);
} else {
  console.error('❌ Error:', result.error);
}
```

### Step 6: Wait and Verify

```typescript
// Wait 30 seconds
await new Promise(r => setTimeout(r, 30000));

// Check recipient
const recipientAccount = await rpc.getAccount('arcv1recipient...');
console.log('Recipient balance:', recipientAccount.balance);

// Check sender
const senderAccount = await rpc.getAccount(myAddress);
console.log('My balance:', senderAccount.balance);
console.log('My nonce:', senderAccount.nonce); // Should be +1
```

---

## Method 3: Using archivas-cli (Ed25519)

### Step 1: Sign Transaction

```bash
./archivas-cli sign-transfer \
  --from-mnemonic "your 24 word mnemonic" \
  --to arcv1recipient... \
  --amount 10000000000 \
  --fee 100000 \
  --nonce 0 \
  --out tx.json
```

This creates `tx.json` with the signed transaction.

### Step 2: Broadcast

```bash
./archivas-cli broadcast \
  --node https://seed.archivas.ai \
  --file tx.json
```

### Step 3: Verify

```bash
curl https://seed.archivas.ai/account/arcv1recipient...
```

---

## Common Issues

### "Invalid nonce"

**Problem:** Nonce doesn't match account's current nonce.

**Solution:**
```bash
# Get current nonce
curl https://seed.archivas.ai/account/YOUR_ADDRESS | jq .nonce

# Use that value in your transaction
```

### "Insufficient funds"

**Problem:** Balance < (amount + fee)

**Solution:**
- Reduce amount
- Or get more testnet RCHV

### "Invalid signature"

**Problem:** Wrong private key or malformed transaction.

**Solution:**
- Verify private key matches `from` address
- Ensure canonical JSON encoding (SDK handles this)

### Transaction not confirming

**Wait 60 seconds** (2-3 blocks).

If still pending:
```bash
# Check mempool
curl https://seed.archivas.ai/mempool

# If empty, transaction was included
# Check your nonce - if it incremented, tx succeeded
curl https://seed.archivas.ai/account/YOUR_ADDRESS | jq .nonce
```

---

## Transaction Explorer

View your transaction in the block explorer:

1. Go to https://archivas-explorer-production.up.railway.app
2. Click "Transactions" tab
3. Find your transfer in recent activity

---

## Amount Conversion

RCHV uses 8 decimals:

| RCHV | Base Units |
|------|------------|
| 1 RCHV | 100,000,000 |
| 10 RCHV | 1,000,000,000 |
| 100 RCHV | 10,000,000,000 |
| 0.001 RCHV | 100,000 |

**Fee:** Always `100000` (0.001 RCHV) for now.

---

## Next Steps

- [Build a Wallet](../developers/building-wallet.md)
- [API Reference](../developers/api-reference.md)
- [SDK Guide](../developers/sdk-guide.md)

---

**Congratulations!** You've sent your first Archivas transaction! 🎉

