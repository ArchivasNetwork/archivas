# Betanet MetaMask and Hardhat Developer Guide

Complete guide for wallet integration and smart contract deployment on Archivas Betanet.

---

## Network Parameters

Use these parameters to connect to Archivas Betanet:

| Parameter | Value |
|-----------|-------|
| **Network Name** | Archivas Betanet |
| **RPC URL** | `https://seed3.betanet.archivas.ai` (public) or `http://localhost:8545` (local) |
| **Chain ID** | `1644` (hex: `0x66c`) |
| **Currency Symbol** | `RCHV` |
| **Block Explorer** | `https://explorer.betanet.archivas.ai` (Blockscout) |
| **Decimals** | 8 (native) / 18 (EVM-compatible) |

**Faucet:**  
Coming soon. For now, earn RCHV by running a farmer.

---

## Part 1: MetaMask Setup

### Add Archivas Betanet to MetaMask

1. **Open MetaMask** in your browser
2. **Click network dropdown** at the top
3. **Select "Add Network" or "Add a network manually"**
4. **Enter the following details:**

```
Network Name:       Archivas Betanet
New RPC URL:        https://seed3.betanet.archivas.ai
Chain ID:           1644
Currency Symbol:    RCHV
Block Explorer URL: https://explorer.betanet.archivas.ai
```

5. **Click "Save"**

You should now see "Archivas Betanet" in your network list.

---

### Import Existing Wallet

If you already have an Archivas wallet (from farming or CLI), import it to MetaMask:

1. **Get your private key** from your wallet:

```bash
# If you generated it with archivas-wallet:
archivas-wallet new
# Save the private key (64 hex characters)
```

2. **In MetaMask:**
   - Click the **account icon** (top right)
   - Select **"Import Account"**
   - Choose **"Private Key"**
   - Paste your private key
   - Click **"Import"**

3. **Verify the address matches:**

```bash
# Use the address converter to verify
cd ~/archivas
go build -o address-converter cmd/address-converter/main.go
./address-converter <YOUR_ARCV_ADDRESS>
```

Your MetaMask 0x address should match the converter output.

---

### Check Balance

Your RCHV balance should automatically display in MetaMask.

**Verify via RPC:**

```bash
curl -s https://seed3.betanet.archivas.ai -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0xYOUR_ADDRESS_HERE", "latest"],
    "id": 1
  }' | jq
```

**Note:** The balance is returned in Wei (18 decimals). To convert to RCHV:

```javascript
// In JavaScript/ethers.js
const balanceWei = "0x..."; // from RPC
const balanceRCHV = parseFloat(BigInt(balanceWei)) / 1e18;
console.log(`Balance: ${balanceRCHV} RCHV`);
```

---

### Send a Transaction

1. **In MetaMask, click "Send"**
2. **Enter recipient address** (0x or arcv1 format works via converter)
3. **Enter amount in RCHV**
4. **Review gas settings:**
   - Gas Price: 1 gwei (default)
   - Gas Limit: 21000 (simple transfer)
5. **Click "Confirm"**

**Monitor transaction:**

```bash
# Get transaction receipt
curl -s https://seed3.betanet.archivas.ai -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionReceipt",
    "params": ["0xTRANSACTION_HASH"],
    "id": 1
  }' | jq
```

---

## Part 2: Hardhat Setup

### Install Hardhat

```bash
mkdir my-archivas-project
cd my-archivas-project
npm init -y
npm install --save-dev hardhat @nomicfoundation/hardhat-toolbox
npx hardhat init
# Select "Create a TypeScript project" or "Create a JavaScript project"
```

---

### Configure Hardhat for Betanet

Create or edit `hardhat.config.js` (or `.ts`):

```javascript
require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.24",
  networks: {
    betanet: {
      url: "https://seed3.betanet.archivas.ai",
      chainId: 1644,
      accounts: [
        "YOUR_PRIVATE_KEY_HERE" // Without 0x prefix
      ],
      gasPrice: 1000000000, // 1 gwei
      gas: 8000000, // 8M gas limit per block
    },
    localhost: {
      url: "http://127.0.0.1:8545",
      chainId: 1644,
      accounts: [
        "YOUR_PRIVATE_KEY_HERE"
      ]
    }
  },
  etherscan: {
    // Blockscout verification (optional)
    apiKey: {
      betanet: "NOT_REQUIRED"
    },
    customChains: [
      {
        network: "betanet",
        chainId: 1644,
        urls: {
          apiURL: "https://explorer.betanet.archivas.ai/api",
          browserURL: "https://explorer.betanet.archivas.ai"
        }
      }
    ]
  }
};
```

**TypeScript version (`hardhat.config.ts`):**

```typescript
import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  solidity: "0.8.24",
  networks: {
    betanet: {
      url: "https://seed3.betanet.archivas.ai",
      chainId: 1644,
      accounts: [
        process.env.PRIVATE_KEY || ""
      ],
      gasPrice: 1000000000, // 1 gwei
      gas: 8000000,
    }
  }
};

export default config;
```

**Store private key in `.env`:**

```bash
# .env
PRIVATE_KEY=your_64_char_private_key_here
```

```bash
# Install dotenv
npm install --save-dev dotenv

# Load in hardhat.config.js
require('dotenv').config();
```

---

### Example Smart Contract

Create `contracts/Greeter.sol`:

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract Greeter {
    string private greeting;
    address public owner;

    event GreetingChanged(string oldGreeting, string newGreeting);

    constructor(string memory _greeting) {
        greeting = _greeting;
        owner = msg.sender;
    }

    function greet() public view returns (string memory) {
        return greeting;
    }

    function setGreeting(string memory _greeting) public {
        require(msg.sender == owner, "Only owner can change greeting");
        emit GreetingChanged(greeting, _greeting);
        greeting = _greeting;
    }

    function getOwner() public view returns (address) {
        return owner;
    }
}
```

**Compile the contract:**

```bash
npx hardhat compile
```

---

### Deployment Script

Create `scripts/deploy-greeter.js`:

```javascript
const hre = require("hardhat");

async function main() {
  const [deployer] = await hre.ethers.getSigners();

  console.log("Deploying Greeter with account:", deployer.address);
  console.log("Account balance:", (await hre.ethers.provider.getBalance(deployer.address)).toString());

  const Greeter = await hre.ethers.getContractFactory("Greeter");
  const greeter = await Greeter.deploy("Hello, Archivas Betanet!");

  await greeter.waitForDeployment();

  const greeterAddress = await greeter.getAddress();
  console.log("Greeter deployed to:", greeterAddress);

  // Verify deployment
  const greeting = await greeter.greet();
  console.log("Initial greeting:", greeting);

  console.log("\nVerify contract on Blockscout:");
  console.log(`https://explorer.betanet.archivas.ai/address/${greeterAddress}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

**TypeScript version (`scripts/deploy-greeter.ts`):**

```typescript
import { ethers } from "hardhat";

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("Deploying Greeter with account:", deployer.address);
  console.log("Account balance:", (await ethers.provider.getBalance(deployer.address)).toString());

  const Greeter = await ethers.getContractFactory("Greeter");
  const greeter = await Greeter.deploy("Hello, Archivas Betanet!");

  await greeter.waitForDeployment();

  const greeterAddress = await greeter.getAddress();
  console.log("Greeter deployed to:", greeterAddress);

  // Verify deployment
  const greeting = await greeter.greet();
  console.log("Initial greeting:", greeting);

  console.log("\nVerify contract on Blockscout:");
  console.log(`https://explorer.betanet.archivas.ai/address/${greeterAddress}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

---

### Deploy to Betanet

```bash
npx hardhat run scripts/deploy-greeter.js --network betanet
```

**Expected output:**

```
Deploying Greeter with account: 0x47ea4b22029c155c835fd0a0b99f8196766f406a
Account balance: 125500000000000000000
Greeter deployed to: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Initial greeting: Hello, Archivas Betanet!

Verify contract on Blockscout:
https://explorer.betanet.archivas.ai/address/0x5FbDB2315678afecb367f032d93F642f64180aa3
```

---

### Interact with Deployed Contract

Create `scripts/interact-greeter.js`:

```javascript
const hre = require("hardhat");

async function main() {
  const contractAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"; // From deployment

  const Greeter = await hre.ethers.getContractFactory("Greeter");
  const greeter = Greeter.attach(contractAddress);

  // Read greeting
  const greeting = await greeter.greet();
  console.log("Current greeting:", greeting);

  // Change greeting
  console.log("\nChanging greeting...");
  const tx = await greeter.setGreeting("Archivas is awesome!");
  console.log("Transaction hash:", tx.hash);

  // Wait for confirmation
  await tx.wait();
  console.log("Transaction confirmed!");

  // Read new greeting
  const newGreeting = await greeter.greet();
  console.log("New greeting:", newGreeting);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

**Run interaction:**

```bash
npx hardhat run scripts/interact-greeter.js --network betanet
```

---

### Verify Transaction Receipt

```bash
# Using curl
curl -s https://seed3.betanet.archivas.ai -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionReceipt",
    "params": ["0xTRANSACTION_HASH_FROM_ABOVE"],
    "id": 1
  }' | jq

# Using Hardhat
npx hardhat run scripts/get-receipt.js --network betanet
```

**Expected receipt:**

```json
{
  "transactionHash": "0x...",
  "blockNumber": "0x1234",
  "from": "0x47ea4b22029c155c835fd0a0b99f8196766f406a",
  "to": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "gasUsed": "0x7530",
  "status": "0x1",
  "logs": [...]
}
```

---

### View Contract on Blockscout

After deployment, your contract will be indexed by Blockscout:

1. **Go to Blockscout:**  
   `https://explorer.betanet.archivas.ai/address/0xYOUR_CONTRACT_ADDRESS`

2. **View contract details:**
   - Transaction history
   - Internal transactions
   - Events/logs
   - Contract bytecode

3. **Verify and publish source code:**

```bash
npx hardhat verify --network betanet 0xYOUR_CONTRACT_ADDRESS "Hello, Archivas Betanet!"
```

**Note:** Blockscout verification requires the contract source code and constructor arguments.

---

## Part 3: ethers.js Integration

### Install ethers.js

```bash
npm install ethers
```

### Connect to Betanet

```javascript
const { ethers } = require("ethers");

// Connect to Archivas Betanet
const provider = new ethers.JsonRpcProvider("https://seed3.betanet.archivas.ai");

// Or connect to local node
// const provider = new ethers.JsonRpcProvider("http://localhost:8545");

// Create wallet from private key
const privateKey = "YOUR_PRIVATE_KEY_HERE";
const wallet = new ethers.Wallet(privateKey, provider);

console.log("Wallet address:", wallet.address);

// Get balance
async function getBalance() {
  const balance = await provider.getBalance(wallet.address);
  console.log("Balance:", ethers.formatEther(balance), "RCHV");
}

getBalance();
```

### Send Transaction

```javascript
async function sendTransaction() {
  const tx = {
    to: "0xRECIPIENT_ADDRESS",
    value: ethers.parseEther("1.0"), // 1 RCHV
  };

  const txResponse = await wallet.sendTransaction(tx);
  console.log("Transaction hash:", txResponse.hash);

  // Wait for confirmation
  const receipt = await txResponse.wait();
  console.log("Transaction confirmed in block:", receipt.blockNumber);
}

sendTransaction();
```

### Call Contract

```javascript
async function callContract() {
  const contractAddress = "0xYOUR_CONTRACT_ADDRESS";
  const contractABI = [
    "function greet() public view returns (string memory)",
    "function setGreeting(string memory _greeting) public"
  ];

  const contract = new ethers.Contract(contractAddress, contractABI, wallet);

  // Read data
  const greeting = await contract.greet();
  console.log("Greeting:", greeting);

  // Write data
  const tx = await contract.setGreeting("New greeting!");
  await tx.wait();
  console.log("Greeting updated!");
}

callContract();
```

---

## Troubleshooting

### "Insufficient funds" Error

**Problem:** Your account doesn't have enough RCHV to cover gas costs.

**Solution:**
- Earn RCHV by running a farmer
- Check balance: `curl http://localhost:8545/balance/YOUR_ARCV_ADDRESS`
- Wait for farming rewards (1 RCHV per block found)

### "Nonce too low" Error

**Problem:** Transaction nonce is outdated.

**Solution:**
```javascript
// Get current nonce
const nonce = await provider.getTransactionCount(wallet.address);
console.log("Current nonce:", nonce);

// Send with explicit nonce
const tx = await wallet.sendTransaction({
  to: "0x...",
  value: ethers.parseEther("1.0"),
  nonce: nonce
});
```

### "Chain ID mismatch" Error

**Problem:** Hardhat/MetaMask is configured with wrong chain ID.

**Solution:**
- Verify `chainId: 1644` in hardhat.config.js
- Verify MetaMask network settings
- Check current chain ID:

```bash
curl -s https://seed3.betanet.archivas.ai -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' | jq
```

Expected: `{"result":"0x66c"}` (1644 in hex)

### Contract Not Showing in Blockscout

**Problem:** Newly deployed contract not indexed yet.

**Solution:**
- Wait 1-2 minutes for indexing
- Refresh Blockscout page
- Check transaction was successful: `eth_getTransactionReceipt`

---

## Address Compatibility Notes

**Archivas uses a dual address system:**

- **ARCV (Bech32):** `arcv1gl4ykgszns24eq6l6zstn8upjemx7sr28u0rm2`  
  Used by Archivas CLI and farmer

- **EVM (Hex):** `0x47ea4b22029c155c835fd0a0b99f8196766f406a`  
  Used by MetaMask, Hardhat, ethers.js

**Both represent the same account!**

Convert between formats:

```bash
./address-converter arcv1gl4ykgszns24eq6l6zstn8upjemx7sr28u0rm2
./address-converter 0x47ea4b22029c155c835fd0a0b99f8196766f406a
```

**When to use which:**
- **MetaMask/Hardhat:** Always use 0x hex addresses
- **Archivas CLI (`archivas-wallet`):** Can use either format
- **Farmer config:** Use private key (same for both)
- **RPC queries:** Both formats work (internally normalized)

---

## Additional Resources

- [Archivas Betanet Explorer](https://explorer.betanet.archivas.ai)
- [Running a Betanet Node](./gitbook/betanet/running-betanet-node.md)
- [Farmer Key Rotation](./FARMER_KEY_ROTATION.md)
- [MetaMask Integration Guide](./METAMASK_GUIDE.md)
- [Hardhat Documentation](https://hardhat.org/docs)
- [ethers.js Documentation](https://docs.ethers.org/v6/)

---

**Questions or issues? Open a GitHub issue or ask in the Archivas Discord.**

