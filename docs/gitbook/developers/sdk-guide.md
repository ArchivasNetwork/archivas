# TypeScript SDK Guide

Complete guide to using `@archivas/sdk` for building applications.

---

## Installation

```bash
npm install @archivas/sdk
# or
yarn add @archivas/sdk
# or
pnpm add @archivas/sdk
```

**Requirements:**
- Node.js 18+
- TypeScript 5+ (recommended)

---

## Core Modules

### Derivation

Generate and manage wallets.

```typescript
import { Derivation } from '@archivas/sdk';

// Generate 24-word mnemonic
const mnemonic = Derivation.mnemonicGenerate();

// Derive keypair from mnemonic
const keyPair = await Derivation.fromMnemonic(mnemonic);
// keyPair.publicKey: Uint8Array (32 bytes)
// keyPair.secretKey: Uint8Array (64 bytes)

// Get address
const address = Derivation.toAddress(keyPair.publicKey);
// address: "arcv1..." (Bech32 format)
```

### Tx (Transactions)

Build, sign, and create transactions.

```typescript
import { Tx } from '@archivas/sdk';

// Build transfer
const tx = Tx.buildTransfer({
  from: 'arcv1...',
  to: 'arcv1...',
  amount: '10000000000',  // 100 RCHV
  fee: '100000',          // 0.001 RCHV
  nonce: '0',
  memo: 'Hello Archivas!' // optional
});

// Sign transaction
const { sigHex, pubHex, hashHex } = await Tx.sign(tx, secretKey);

// Or create complete signed tx
const signedTx = await Tx.createSigned(tx, secretKey);
```

### RpcClient

Query the blockchain and submit transactions.

```typescript
import { createRpcClient } from '@archivas/sdk';

const rpc = createRpcClient({
  baseUrl: 'https://seed.archivas.ai',
  timeout: 30000  // optional
});

// Get chain tip
const tip = await rpc.getChainTip();

// Get account
const account = await rpc.getAccount('arcv1...');

// Submit transaction
const result = await rpc.submit(signedTx);
```

---

## Complete Example: Send RCHV

```typescript
import { Derivation, Tx, createRpcClient } from '@archivas/sdk';

async function sendRCHV(
  mnemonic: string,
  recipientAddress: string,
  amountRCHV: number
) {
  // 1. Setup wallet
  const keyPair = await Derivation.fromMnemonic(mnemonic);
  const myAddress = Derivation.toAddress(keyPair.publicKey);
  
  // 2. Connect to RPC
  const rpc = createRpcClient({
    baseUrl: 'https://seed.archivas.ai'
  });
  
  // 3. Get current nonce
  const account = await rpc.getAccount(myAddress);
  
  // 4. Build transaction
  const amountBaseUnits = (amountRCHV * 100000000).toString();
  const tx = Tx.buildTransfer({
    from: myAddress,
    to: recipientAddress,
    amount: amountBaseUnits,
    fee: '100000',
    nonce: account.nonce
  });
  
  // 5. Sign
  const signedTx = await Tx.createSigned(tx, keyPair.secretKey);
  
  // 6. Submit
  const result = await rpc.submit(signedTx);
  
  if (result.ok) {
    console.log('✅ Transaction submitted!');
    console.log('Hash:', result.hash);
    return result.hash;
  } else {
    throw new Error(result.error);
  }
}

// Usage
await sendRCHV(
  "your 24 word mnemonic here",
  "arcv1recipient...",
  100  // 100 RCHV
);
```

---

## API Methods

### Derivation

```typescript
// Generate mnemonic
Derivation.mnemonicGenerate(): string

// Derive keypair
Derivation.fromMnemonic(mnemonic: string): Promise<KeyPair>

// Get address
Derivation.toAddress(publicKey: Uint8Array): Address

// Direct mnemonic → address
Derivation.mnemonicToAddress(mnemonic: string): Promise<Address>
```

### Tx

```typescript
// Build transfer
Tx.buildTransfer(params: TransferParams): TxBody

// Get canonical bytes
Tx.canonicalBytes(tx: TxBody): Uint8Array

// Hash transaction
Tx.hash(tx: TxBody): Uint8Array

// Sign transaction
Tx.sign(tx: TxBody, secretKey: Uint8Array): Promise<{sigHex, pubHex, hashHex}>

// Create signed tx (all-in-one)
Tx.createSigned(tx: TxBody, secretKey: Uint8Array): Promise<SignedTx>

// Estimate size
Tx.estimateSize(tx: TxBody): number
```

### RpcClient

```typescript
// Chain
rpc.getChainTip(): Promise<ChainTip>

// Accounts
rpc.getAccount(address: Address): Promise<Account>
rpc.getBalance(address: Address): Promise<bigint>
rpc.getNonce(address: Address): Promise<bigint>

// Transactions
rpc.getTx(hash: string): Promise<TxDetails>
rpc.getMempool(): Promise<readonly string[]>
rpc.estimateFee(bytes: number): Promise<FeeEstimate>
rpc.submit(stx: SignedTx): Promise<SubmitResponse>
```

---

## Types

```typescript
type Address = `arcv${string}`;

interface KeyPair {
  publicKey: Uint8Array;   // 32 bytes
  secretKey: Uint8Array;   // 64 bytes
}

interface TxBody {
  type: "transfer";
  from: Address;
  to: Address;
  amount: string;   // u64 as string
  fee: string;      // u64 as string
  nonce: string;    // u64 as string
  memo?: string;    // optional, max 256 bytes
}

interface SignedTx {
  tx: TxBody;
  pubkey: string;  // hex
  sig: string;     // hex
  hash: string;    // hex
}

interface ChainTip {
  height: string;
  hash: string;
  difficulty: string;
}

interface Account {
  address: string;
  balance: string;  // base units
  nonce: string;
}
```

---

## Advanced Usage

### Memo Field

```typescript
const tx = Tx.buildTransfer({
  from: myAddress,
  to: recipientAddress,
  amount: '10000000000',
  fee: '100000',
  nonce: '0',
  memo: 'Payment for services'  // max 256 bytes UTF-8
});
```

### Custom RPC Config

```typescript
const rpc = createRpcClient({
  baseUrl: 'https://seed.archivas.ai',
  timeout: 60000,
  headers: {
    'X-Custom-Header': 'value'
  }
});
```

### Batch Queries

```typescript
const addresses = ['arcv1...', 'arcv1...', 'arcv1...'];

const balances = await Promise.all(
  addresses.map(addr => rpc.getBalance(addr))
);

console.log('Total:', balances.reduce((a, b) => a + b, 0n));
```

---

## Security

### Never Log Private Keys

```typescript
// ❌ BAD
console.log('Private key:', secretKey);

// ✅ GOOD
console.log('Address:', address);
```

### Wipe Secret Keys After Use

```typescript
// After signing
signedTx = await Tx.createSigned(tx, secretKey);

// Wipe the secret key
secretKey.fill(0);
```

### Validate Addresses

```typescript
import { isValidAddress } from '@archivas/sdk';

if (!isValidAddress(recipientAddress)) {
  throw new Error('Invalid recipient address');
}
```

---

## Testing

```typescript
import { Derivation, Tx } from '@archivas/sdk';

describe('Archivas SDK', () => {
  it('should generate valid mnemonic', () => {
    const mnemonic = Derivation.mnemonicGenerate();
    expect(mnemonic.split(' ')).toHaveLength(24);
  });
  
  it('should derive consistent address', async () => {
    const mnemonic = "abandon abandon abandon...";
    const kp1 = await Derivation.fromMnemonic(mnemonic);
    const kp2 = await Derivation.fromMnemonic(mnemonic);
    
    const addr1 = Derivation.toAddress(kp1.publicKey);
    const addr2 = Derivation.toAddress(kp2.publicKey);
    
    expect(addr1).toBe(addr2);
  });
});
```

---

## Next Steps

- [Building a Wallet](building-wallet.md)
- [Transaction Signing](transaction-signing.md)
- [API Reference](api-reference.md)

---

**Full SDK docs:** https://github.com/ArchivasNetwork/archivas-sdk

