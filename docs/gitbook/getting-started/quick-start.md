# Quick Start

Get started with Archivas in 5 minutes!

---

## Choose Your Path

### 🔧 I want to build apps
→ Use the [Public RPC API](#option-1-use-public-api)  
→ Install the [TypeScript SDK](#option-2-use-typescript-sdk)

### 🌾 I want to farm RCHV
→ [Setup a Farming Node](../farmers/setup-farmer.md)

### 👁️ I want to explore the chain
→ [Block Explorer](https://archivas-explorer-production.up.railway.app)

---

## Option 1: Use Public API

**No setup required!** Query the chain directly.

```bash
# Get current height
curl https://seed.archivas.ai/chainTip

# Get an account balance
curl https://seed.archivas.ai/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7

# List recent blocks
curl https://seed.archivas.ai/blocks/recent?limit=10

# Get specific block
curl https://seed.archivas.ai/block/64000
```

**API Documentation:** [API Reference](../developers/api-reference.md)

---

## Option 2: Use TypeScript SDK

### Install

```bash
npm install @archivas/sdk
```

### Create a Wallet

```typescript
import { Derivation } from '@archivas/sdk';

// Generate mnemonic
const mnemonic = Derivation.mnemonicGenerate();
console.log('Save this:', mnemonic);

// Derive keypair
const keyPair = await Derivation.fromMnemonic(mnemonic);

// Get address
const address = Derivation.toAddress(keyPair.publicKey);
console.log('Your address:', address);
```

### Query the Chain

```typescript
import { createRpcClient } from '@archivas/sdk';

const rpc = createRpcClient({
  baseUrl: 'https://seed.archivas.ai'
});

// Get chain info
const tip = await rpc.getChainTip();
console.log('Height:', tip.height);

// Get balance
const account = await rpc.getAccount(address);
console.log('Balance:', account.balance, 'base units');
```

**SDK Guide:** [TypeScript SDK](../developers/sdk-guide.md)

---

## Option 3: Explore with the Block Explorer

**Live Explorer:** https://archivas-explorer-production.up.railway.app

**Features:**
- Real-time block updates
- Account balances
- Transaction history
- Farmer leaderboard
- Network statistics

---

## Next Steps

### For Developers
1. Read [Developer Overview](../developers/overview.md)
2. Review [API Reference](../developers/api-reference.md)
3. Try [Building a Wallet](../developers/building-wallet.md)

### For Farmers
1. Check [Hardware Requirements](../farmers/hardware-requirements.md)
2. Follow [Farmer Setup Guide](../farmers/setup-farmer.md)
3. Learn about [Creating Plots](../farmers/creating-plots.md)

### For Everyone
1. Join GitHub discussions
2. Explore the testnet
3. Read the [FAQ](../reference/faq.md)

---

## Get Testnet RCHV

**Coming soon:** Public faucet

**For now:** Ask in GitHub Discussions for testnet RCHV to experiment with transfers.

---

## Resources

- **Public RPC:** https://seed.archivas.ai
- **Block Explorer:** https://archivas-explorer-production.up.railway.app
- **GitHub:** https://github.com/ArchivasNetwork/archivas
- **TypeScript SDK:** https://github.com/ArchivasNetwork/archivas-sdk
- **Grafana:** http://57.129.148.132:3001

---

**Ready to build?** Continue to [Developer Overview](../developers/overview.md)!

