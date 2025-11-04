# Developer Quickstart

Build applications on Archivas in 5 minutes.

---

## Using the Public RPC

**Endpoint:** `https://seed.archivas.ai`

No setup required - the public RPC is ready to use!

### Quick Test

```bash
# Get chain tip
curl https://seed.archivas.ai/chainTip
# {"height":"64000","hash":"...","difficulty":"1000000"}

# Get account balance
curl https://seed.archivas.ai/account/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
# {"address":"arcv1...","balance":"...","nonce":"0"}

# List recent blocks
curl https://seed.archivas.ai/blocks/recent?limit=10
# {"blocks":[...]}
```

---

## TypeScript SDK

### Installation

```bash
npm install @archivas/sdk
# or
yarn add @archivas/sdk
```

### Generate a Wallet

```typescript
import { Derivation } from '@archivas/sdk';

// Generate 24-word mnemonic
const mnemonic = Derivation.mnemonicGenerate();
console.log('Mnemonic:', mnemonic);

// Derive keypair
const keyPair = await Derivation.fromMnemonic(mnemonic);

// Get address
const address = Derivation.toAddress(keyPair.publicKey);
console.log('Address:', address);  // arcv1...
```

### Query the Chain

```typescript
import { createRpcClient } from '@archivas/sdk';

const rpc = createRpcClient({
  baseUrl: 'https://seed.archivas.ai'
});

// Get chain tip
const tip = await rpc.getChainTip();
console.log('Height:', tip.height);

// Get balance
const balance = await rpc.getBalance('arcv1...');
console.log('Balance:', balance.toString(), 'base units');
```

### Send a Transfer

```typescript
import { Derivation, Tx, createRpcClient } from '@archivas/sdk';

// Setup
const mnemonic = "your 24 word mnemonic here";
const keyPair = await Derivation.fromMnemonic(mnemonic);
const address = Derivation.toAddress(keyPair.publicKey);

const rpc = createRpcClient({ 
  baseUrl: 'https://seed.archivas.ai'
});

// Get current nonce
const account = await rpc.getAccount(address);

// Build transaction
const tx = Tx.buildTransfer({
  from: address,
  to: 'arcv1recipient...',
  amount: '30000000000',  // 300 RCHV (8 decimals)
  fee: '100000',          // 0.001 RCHV
  nonce: account.nonce
});

// Sign
const signedTx = await Tx.createSigned(tx, keyPair.secretKey);

// Submit
const result = await rpc.submit(signedTx);
console.log('Transaction hash:', result.hash);
```

---

## API Endpoints

All endpoints return JSON with numeric fields as strings.

### Chain Info
- `GET /chainTip` - Current height, hash, difficulty
- `GET /genesisHash` - Genesis block hash

### Blocks
- `GET /blocks/recent?limit=N` - Recent blocks (max 100)
- `GET /block/<height>` - Block by height

### Transactions  
- `GET /tx/recent?limit=N` - Recent transactions (max 200)
- `GET /tx/<hash>` - Transaction by hash
- `GET /mempool` - Pending transactions
- `POST /submit` - Submit signed transaction

### Accounts
- `GET /account/<address>` - Balance and nonce

### Utilities
- `GET /estimateFee?bytes=<n>` - Fee estimation

---

## Example: Simple Web Wallet

```typescript
import { Derivation, Tx, createRpcClient } from '@archivas/sdk';
import { useState, useEffect } from 'react';

function Wallet() {
  const [balance, setBalance] = useState('0');
  const [address, setAddress] = useState('');
  
  useEffect(() => {
    const init = async () => {
      // Load or generate wallet
      const mnemonic = localStorage.getItem('mnemonic') || 
                       Derivation.mnemonicGenerate();
      localStorage.setItem('mnemonic', mnemonic);
      
      const kp = await Derivation.fromMnemonic(mnemonic);
      const addr = Derivation.toAddress(kp.publicKey);
      setAddress(addr);
      
      // Get balance
      const rpc = createRpcClient({ 
        baseUrl: 'https://seed.archivas.ai' 
      });
      const bal = await rpc.getBalance(addr);
      setBalance((Number(bal) / 100000000).toFixed(8));
    };
    
    init();
  }, []);
  
  return (
    <div>
      <h2>Archivas Wallet</h2>
      <p>Address: {address}</p>
      <p>Balance: {balance} RCHV</p>
    </div>
  );
}
```

---

## Rate Limits

- **`/submit`:** 10 requests/minute per IP (burst: 5)
- **Other endpoints:** No limit

Handle rate limits gracefully:
```typescript
async function submitWithRetry(rpc, signedTx, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await rpc.submit(signedTx);
    } catch (err) {
      if (err.status === 429 && i < maxRetries - 1) {
        await new Promise(r => setTimeout(r, 6000)); // Wait 6s
        continue;
      }
      throw err;
    }
  }
}
```

---

## Next Steps

- **Build a wallet:** See [Building a Wallet](building-wallet.md)
- **Integrate payments:** See [Transaction Signing](transaction-signing.md)
- **Explore the chain:** https://archivas-explorer-production.up.railway.app
- **Read API docs:** [api-reference.md](api-reference.md)

---

## Resources

- **SDK Repository:** https://github.com/ArchivasNetwork/archivas-sdk
- **Core Repository:** https://github.com/ArchivasNetwork/archivas
- **Explorer Repository:** https://github.com/ArchivasNetwork/archivas-explorer
- **API Docs:** [api-reference.md](api-reference.md)
