# Developer Overview

Build applications on Archivas using the public API and TypeScript SDK.

---

## What You Can Build

### Wallets
- Web wallets (browser-based)
- Mobile wallets (React Native)
- Desktop wallets (Electron)
- Hardware wallet integration (future)

### Explorers
- Block browsers
- Transaction trackers
- Account analytics
- Network statistics

### Applications
- Payment processors
- Tipping systems
- Crowdfunding platforms
- NFT marketplaces (future)
- DeFi applications (future)

---

## Development Stack

### Backend: Public RPC
- **Endpoint:** https://seed.archivas.ai
- **Protocol:** HTTPS (HTTP/2)
- **Format:** JSON
- **Auth:** None required
- **CORS:** Enabled

### Frontend: TypeScript SDK
- **Package:** `@archivas/sdk`
- **Platform:** Node.js + Browser
- **Features:** Wallet, signing, RPC client
- **Types:** Full TypeScript support

### Tools
- **CLI:** `archivas-cli` (keygen, sign, broadcast)
- **Explorer:** Reference implementation
- **Monitoring:** Prometheus metrics (internal)

---

## Quick Start

### 1. Install SDK

```bash
npm install @archivas/sdk
```

### 2. Create a Wallet

```typescript
import { Derivation } from '@archivas/sdk';

const mnemonic = Derivation.mnemonicGenerate();
const keyPair = await Derivation.fromMnemonic(mnemonic);
const address = Derivation.toAddress(keyPair.publicKey);
```

### 3. Query the Chain

```typescript
import { createRpcClient } from '@archivas/sdk';

const rpc = createRpcClient({
  baseUrl: 'https://seed.archivas.ai'
});

const tip = await rpc.getChainTip();
const account = await rpc.getAccount(address);
```

### 4. Send a Transfer

```typescript
import { Tx } from '@archivas/sdk';

const tx = Tx.buildTransfer({
  from: address,
  to: 'arcv1recipient...',
  amount: '10000000000',
  fee: '100000',
  nonce: account.nonce
});

const signedTx = await Tx.createSigned(tx, keyPair.secretKey);
const result = await rpc.submit(signedTx);
```

---

## API Capabilities

### Chain Information
- Get current height, hash, difficulty
- Get genesis hash
- Query network status

### Blocks
- List recent blocks (paginated)
- Get block by height
- Get block transactions

### Accounts
- Get balance and nonce
- Query account state
- Track balance changes

### Transactions
- Submit signed transactions
- List recent transfers
- Query transaction status
- Estimate fees

---

## Architecture Patterns

### Web Wallet

```typescript
// React component
function Wallet() {
  const [balance, setBalance] = useState('0');
  
  useEffect(() => {
    const rpc = createRpcClient({
      baseUrl: 'https://seed.archivas.ai'
    });
    
    const fetchBalance = async () => {
      const account = await rpc.getAccount(myAddress);
      setBalance(account.balance);
    };
    
    fetchBalance();
    const interval = setInterval(fetchBalance, 10000);
    return () => clearInterval(interval);
  }, []);
  
  return <div>Balance: {balance} base units</div>;
}
```

### Payment Processor

```typescript
class PaymentProcessor {
  constructor(private rpc: RpcClient) {}
  
  async watchForPayment(address: string, expectedAmount: bigint): Promise<boolean> {
    const initialAccount = await this.rpc.getAccount(address);
    const initialBalance = BigInt(initialAccount.balance);
    
    // Poll every 5 seconds for 2 minutes
    for (let i = 0; i < 24; i++) {
      await new Promise(r => setTimeout(r, 5000));
      
      const account = await this.rpc.getAccount(address);
      const currentBalance = BigInt(account.balance);
      
      if (currentBalance >= initialBalance + expectedAmount) {
        return true; // Payment received!
      }
    }
    
    return false; // Timeout
  }
}
```

---

## Best Practices

### Error Handling

```typescript
try {
  const result = await rpc.submit(signedTx);
  if (!result.ok) {
    console.error('Transaction rejected:', result.error);
  }
} catch (err) {
  if (err.status === 429) {
    // Rate limited - wait and retry
    await new Promise(r => setTimeout(r, 6000));
    return retry();
  }
  throw err;
}
```

### Nonce Management

```typescript
class NonceTracker {
  private nonce: number;
  
  async getNextNonce(address: string): Promise<number> {
    const account = await rpc.getAccount(address);
    return parseInt(account.nonce);
  }
  
  async submitWithNonce(tx: TxBody, secretKey: Uint8Array) {
    const nonce = await this.getNextNonce(tx.from);
    tx.nonce = nonce.toString();
    
    const signedTx = await Tx.createSigned(tx, secretKey);
    return await rpc.submit(signedTx);
  }
}
```

### Balance Formatting

```typescript
function formatRCHV(baseUnits: string): string {
  const amount = BigInt(baseUnits);
  const rchv = Number(amount) / 100000000;
  return rchv.toFixed(8) + ' RCHV';
}

// Usage
const account = await rpc.getAccount(address);
console.log('Balance:', formatRCHV(account.balance));
// "Balance: 100.00000000 RCHV"
```

---

## Example Projects

### Minimal Web Wallet

**GitHub:** https://github.com/ArchivasNetwork/archivas-sdk  
**Demo:** Coming soon

**Features:**
- Generate wallet from mnemonic
- Display balance
- Send transfers
- Transaction history

### Block Explorer

**GitHub:** https://github.com/ArchivasNetwork/archivas-explorer  
**Live:** https://archivas-explorer-production.up.railway.app

**Features:**
- Real-time blocks
- Account lookup
- Transaction browsing
- Network stats

---

## Resources

### Documentation
- [API Reference](api-reference.md)
- [SDK Guide](sdk-guide.md)
- [Transaction Signing](transaction-signing.md)

### Code Examples
- [Building a Wallet](building-wallet.md)
- [Explorer Integration](explorer-integration.md)

### Live Services
- Public RPC: https://seed.archivas.ai
- Block Explorer: https://archivas-explorer.up.railway.app

---

## Get Help

**Questions?**
- GitHub Discussions: https://github.com/ArchivasNetwork/archivas/discussions
- GitHub Issues: https://github.com/ArchivasNetwork/archivas/issues

---

**Ready to build?** Continue to [SDK Guide](sdk-guide.md)!

