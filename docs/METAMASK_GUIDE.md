# MetaMask Integration Guide

This guide explains how to connect MetaMask to Archivas Betanet and use your Archivas wallet with MetaMask.

## 🔑 Understanding Addresses

Archivas uses a **dual address system**:

| Format | Example | Usage |
|--------|---------|-------|
| **Bech32** | `arcv1s9m9avxdkzuv9lf6wle...` | Archivas CLI, User-facing |
| **0x (EVM)** | `0x81765eb0cdb0b8c2fd3a77f2a1aa16fe...` | MetaMask, Smart Contracts |

**Important:** Both formats represent the **same address** derived from the **same private key**!

---

## 🦊 Part 1: Add Archivas Betanet to MetaMask

### Option A: Manual Configuration

1. Open MetaMask
2. Click the **network dropdown** (top left)
3. Click **"Add Network"** → **"Add a network manually"**
4. Enter network details:

```
Network Name:       Archivas Betanet
RPC URL:            http://51.89.11.4:8545
Chain ID:           1644
Currency Symbol:    RCHV
Block Explorer URL: (leave empty)
```

5. Click **"Save"**
6. Switch to Archivas Betanet

### Option B: Using Your Own Node

If you're running your own Archivas node:

```
Network Name:       Archivas Betanet (Local)
RPC URL:            http://localhost:8545
Chain ID:           1644
Currency Symbol:    RCHV
Block Explorer URL: (leave empty)
```

### Option C: Using Seed Domains

```
RPC URL: http://seed3.betanet.archivas.ai:8545
```

---

## 🔐 Part 2: Import Your Archivas Wallet

### Step 1: Get Your Private Key

When you created your Archivas wallet with `archivas-wallet new`, you received:

```
🔐 New Archivas Wallet Generated

Address:     arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs
Public Key:  024aedb4b79bf799cf484a7369b151f6fb4e1988d745e91b6a9fd9d9eb195a7359
Private Key: 1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899
             ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
             This is what you need for MetaMask!
```

### Step 2: Import into MetaMask

1. In MetaMask, click the **account icon** (top right)
2. Select **"Import Account"**
3. Choose **"Private Key"**
4. Paste your **64-character private key** (without 0x prefix)
5. Click **"Import"**

**Example:**
```
1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899
```

### Step 3: Verify Your Address

MetaMask will display your address in `0x...` format:

```
MetaMask shows:     0x81765eb0cdb0b8c2fd3a77f2a1aa16fe0608f497
Archivas shows:     arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs
```

**These are the SAME address!** ✅

---

## 🔄 Part 3: Converting Between Address Formats

### Using the Address Converter Tool

```bash
# Build the converter
cd ~/archivas
go build -o address-converter cmd/address-converter/main.go

# Convert arcv1 → 0x
./address-converter arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs

# Convert 0x → arcv1
./address-converter 0x81765eb0cdb0b8c2fd3a77f2a1aa16fe0608f497
```

### Using RPC Endpoints

```bash
# Convert arcv1 → 0x
curl -s http://51.89.11.4:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "arcv_toHexAddress",
    "params": ["arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs"],
    "id": 1
  }' | jq

# Convert 0x → arcv1
curl -s http://51.89.11.4:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "arcv_fromHexAddress",
    "params": ["0x81765eb0cdb0b8c2fd3a77f2a1aa16fe0608f497"],
    "id": 1
  }' | jq
```

---

## 💰 Part 4: Check Your Balance

### In MetaMask

1. Switch to **Archivas Betanet** network
2. Your RCHV balance should display automatically
3. If you see "0 RCHV", refresh or wait a few seconds

### Via RPC

```bash
# Check balance (replace with your address)
curl -s http://51.89.11.4:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0x81765eb0cdb0b8c2fd3a77f2a1aa16fe0608f497", "latest"],
    "id": 1
  }' | jq

# Result is in wei (smallest unit)
# 1 RCHV = 10^8 smallest units
```

---

## 📤 Part 5: Sending Transactions

### From MetaMask

1. Click **"Send"**
2. Enter recipient address (can use either format):
   - `0x...` works directly
   - `arcv1...` must be converted to `0x...` first
3. Enter amount in RCHV
4. Click **"Next"** → **"Confirm"**

### Transaction Details

```
Gas Price:   Use default (auto-calculated)
Gas Limit:   21000 for simple transfers
            Higher for contract calls
```

### Verify Transaction

After sending, MetaMask shows:
- **Transaction Hash** (0x...)
- **Status** (Pending/Confirmed)
- **Block Number**

Check on Archivas:
```bash
curl -s http://51.89.11.4:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionReceipt",
    "params": ["0x_YOUR_TX_HASH_HERE"],
    "id": 1
  }' | jq
```

---

## 🔧 Part 6: Using With Development Tools

### Hardhat Configuration

```javascript
// hardhat.config.js
module.exports = {
  networks: {
    betanet: {
      url: "http://51.89.11.4:8545",
      chainId: 1644,
      accounts: ["0x1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899"]
    }
  }
};
```

### Ethers.js

```javascript
const { ethers } = require("ethers");

// Connect to Archivas Betanet
const provider = new ethers.providers.JsonRpcProvider(
  "http://51.89.11.4:8545"
);

// Create wallet from private key
const wallet = new ethers.Wallet(
  "0x1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899",
  provider
);

// Check balance
const balance = await wallet.getBalance();
console.log("Balance:", ethers.utils.formatUnits(balance, 8), "RCHV");

// Send transaction
const tx = await wallet.sendTransaction({
  to: "0x_RECIPIENT_ADDRESS",
  value: ethers.utils.parseUnits("10", 8) // 10 RCHV
});

await tx.wait();
console.log("Transaction mined:", tx.hash);
```

### Web3.js

```javascript
const Web3 = require('web3');

// Connect to Archivas Betanet
const web3 = new Web3('http://51.89.11.4:8545');

// Add account
const account = web3.eth.accounts.privateKeyToAccount(
  '0x1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899'
);
web3.eth.accounts.wallet.add(account);

// Check balance
const balance = await web3.eth.getBalance(account.address);
console.log('Balance:', web3.utils.fromWei(balance, 'gwei'), 'RCHV');

// Send transaction
const tx = await web3.eth.sendTransaction({
  from: account.address,
  to: '0x_RECIPIENT_ADDRESS',
  value: web3.utils.toWei('10', 'gwei'),
  gas: 21000
});

console.log('Transaction mined:', tx.transactionHash);
```

---

## 🎮 Part 7: Smart Contracts

### Deploy a Contract

```javascript
// Using Hardhat
npx hardhat run scripts/deploy.js --network betanet

// Using Remix
// 1. Connect MetaMask to Archivas Betanet
// 2. In Remix, select "Injected Provider - MetaMask"
// 3. Deploy your contract
```

### Interact with Contracts

```javascript
const { ethers } = require("ethers");

// Contract ABI and address
const contractAddress = "0x_YOUR_CONTRACT_ADDRESS";
const abi = [/* your contract ABI */];

// Connect
const provider = new ethers.providers.JsonRpcProvider("http://51.89.11.4:8545");
const wallet = new ethers.Wallet("YOUR_PRIVATE_KEY", provider);
const contract = new ethers.Contract(contractAddress, abi, wallet);

// Call functions
const result = await contract.someFunction();
console.log("Result:", result);

// Send transactions
const tx = await contract.someWriteFunction("arg1", "arg2");
await tx.wait();
console.log("Transaction mined:", tx.hash);
```

---

## 🔍 Part 8: Troubleshooting

### Balance Not Showing

**Problem:** MetaMask shows 0 RCHV but you have balance

**Solutions:**
1. Refresh MetaMask (close and reopen)
2. Check you're on Archivas Betanet network
3. Verify address is correct:
   ```bash
   ./address-converter arcv1YOUR_ADDRESS
   ```
4. Check balance via RPC:
   ```bash
   curl http://51.89.11.4:8545 -X POST \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xYOUR_ADDRESS","latest"],"id":1}'
   ```

### Wrong Address in MetaMask

**Problem:** MetaMask shows different address than expected

**Cause:** Wrong private key imported

**Solution:**
1. Remove the account from MetaMask
2. Re-import with correct private key
3. Verify with converter tool

### Transactions Failing

**Problem:** Transactions show "Failed" status

**Common causes:**
- Insufficient balance for gas fees
- Gas limit too low
- Nonce issues (reset account in MetaMask settings)

**Solutions:**
1. Check balance covers amount + gas fees
2. Increase gas limit for complex transactions
3. Reset account: Settings → Advanced → Reset Account

### RPC Not Responding

**Problem:** MetaMask can't connect to network

**Solutions:**
1. Check RPC URL is correct
2. Verify Archivas node is running:
   ```bash
   curl http://51.89.11.4:8545 -X POST \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
   ```
3. Try alternative RPC:
   - `http://localhost:8545` (if running own node)
   - `http://seed1.betanet.archivas.ai:8545`

---

## ⚠️ Security Best Practices

### 1. Protect Your Private Key

- ✅ Never share your private key
- ✅ Store encrypted backups (use GPG)
- ✅ Use hardware wallet for large amounts
- ❌ Never paste private key in Discord/Telegram
- ❌ Never store in plain text files

### 2. Verify Addresses

Always verify recipient addresses before sending:

```bash
# Double-check address conversion
./address-converter arcv1_ADDRESS
./address-converter 0x_ADDRESS
```

### 3. Test Transactions

- Send small test amounts first
- Verify receipt before sending large amounts
- Keep transaction hashes for records

### 4. Network Verification

Always verify you're on the correct network:
- **Chain ID:** 1644
- **Network Name:** Archivas Betanet
- **RPC URL:** Matches your configuration

---

## 📋 Quick Reference

### Network Details

```
Network Name: Archivas Betanet
Chain ID:     1644 (0x66c in hex)
RPC URL:      http://51.89.11.4:8545
Currency:     RCHV
Decimals:     8
```

### Essential Commands

```bash
# Build converter
go build -o address-converter cmd/address-converter/main.go

# Convert addresses
./address-converter arcv1...
./address-converter 0x...

# Check balance
curl -s http://51.89.11.4:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xYOUR_ADDRESS","latest"],"id":1}'

# Check chain ID
curl -s http://51.89.11.4:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'
```

---

## 🆘 Support

Need help with MetaMask integration?

- **Discord:** [discord.gg/archivas](https://discord.gg/archivas)
- **Telegram:** [t.me/archivas](https://t.me/archivas)
- **GitHub Issues:** [github.com/ArchivasNetwork/archivas/issues](https://github.com/ArchivasNetwork/archivas/issues)

---

## ✅ Summary

**Key Points:**
1. ✅ Same private key works in both Archivas CLI and MetaMask
2. ✅ `arcv1...` and `0x...` are just different formats of the same address
3. ✅ Your balance appears automatically in MetaMask
4. ✅ All Ethereum tools (Hardhat, Remix, Ethers.js) work with Archivas
5. ✅ Chain ID 1644 identifies Archivas Betanet

**You're now ready to use MetaMask with Archivas!** 🎉

