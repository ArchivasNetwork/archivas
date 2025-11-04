# Building a Wallet

Build a complete wallet application for Archivas.

## Wallet Architecture

**Components needed:**
1. Key management (BIP39 + Ed25519)
2. Address derivation (Bech32)
3. Transaction signing (Ed25519 + Blake2b)
4. RPC client (query balances, submit txs)
5. UI (web, mobile, or desktop)

## Using the SDK

```typescript
import { Derivation, Tx, createRpcClient } from '@archivas/sdk';

class ArchivasWallet {
  private keyPair: KeyPair;
  private address: Address;
  private rpc: RpcClient;
  
  async initialize(mnemonic: string) {
    this.keyPair = await Derivation.fromMnemonic(mnemonic);
    this.address = Derivation.toAddress(this.keyPair.publicKey);
    this.rpc = createRpcClient({
      baseUrl: 'https://seed.archivas.ai'
    });
  }
  
  async getBalance(): Promise<bigint> {
    return await this.rpc.getBalance(this.address);
  }
  
  async send(to: Address, amount: string): Promise<string> {
    const account = await this.rpc.getAccount(this.address);
    
    const tx = Tx.buildTransfer({
      from: this.address,
      to,
      amount,
      fee: '100000',
      nonce: account.nonce
    });
    
    const signedTx = await Tx.createSigned(tx, this.keyPair.secretKey);
    const result = await this.rpc.submit(signedTx);
    
    if (!result.ok) throw new Error(result.error);
    return result.hash!;
  }
}
```

## Security Best Practices

1. **Never log private keys**
2. **Encrypt mnemonics at rest**
3. **Use HTTPS only**
4. **Validate addresses before sending**
5. **Show confirmation before signing**

---

**Full guide:** [SDK Guide](sdk-guide.md)
