# Frequently Asked Questions

---

## General

### What is Archivas?

Archivas is a Proof-of-Space-and-Time (PoST) blockchain where farmers earn RCHV by allocating disk space. It uses the same consensus model as Chia Network but with a modern Go implementation and developer-friendly API.

### Is Archivas live?

Yes! The testnet is operational at https://seed.archivas.ai with 25+ days of uptime and 64,000+ blocks.

### How is this different from Chia?

**Same consensus:** Both use PoSpace + VDF  
**Different implementation:** Go (Archivas) vs Python (Chia)  
**Different focus:** Developer API and simplicity (Archivas) vs full ecosystem (Chia)

You can think of Archivas as a streamlined, developer-focused PoST chain.

### Is this a fork of Chia?

No. Archivas is built from scratch in Go. We implement the same consensus *model* (PoSpace + VDF) but with completely independent code.

---

## Farming

### Can I farm RCHV?

Yes! Anyone can farm by:
1. Creating plots (precomputed hash tables)
2. Running the farmer software
3. Connecting to a node

See [Farmer Setup Guide](../farmers/setup-farmer.md).

### What hardware do I need?

- **Minimum:** 10 GB free disk space (for k=27 plot)
- **Recommended:** 100+ GB SSD/HDD for multiple plots
- **CPU:** Any modern CPU (plotting is CPU-bound)
- **RAM:** 2 GB minimum, 4 GB recommended

### How much can I earn?

Rewards depend on your plot size relative to total network space:

**Current network:** ~55 GB total  
**Your 10 GB:** ~18% of network = ~18% of blocks = ~3.6 RCHV/block average

With ~3,000 blocks/day × 18% = ~540 blocks/day × 20 RCHV = **10,800 RCHV/day**

*Note: Testnet only - no real value.*

### Can I farm Chia and Archivas together?

Yes! You can run both simultaneously:
- **Same disk:** Can host both plot formats
- **Different plots:** Chia plots ≠ Archivas plots (incompatible formats)
- **Same hardware:** Both use disk I/O, minimal CPU

---

## Wallets & Transactions

### How do I create a wallet?

**Option 1: TypeScript SDK**
```typescript
import { Derivation } from '@archivas/sdk';
const mnemonic = Derivation.mnemonicGenerate();
const kp = await Derivation.fromMnemonic(mnemonic);
const address = Derivation.toAddress(kp.publicKey);
```

**Option 2: CLI**
```bash
./archivas-cli keygen
```

Both generate a 24-word BIP39 mnemonic and derive an Ed25519 keypair.

### What's a mnemonic?

A 24-word phrase (e.g., "abandon abandon abandon...") that can recreate your wallet.

**Never share your mnemonic!** Anyone with it can access your RCHV.

### What do addresses look like?

Archivas uses Bech32 addresses with `arcv` prefix:
```
arcv1t3huuyd08er3yfnmk9c935rmx3wdh5j6m2uc9d
```

### How do I send RCHV?

```bash
./archivas-wallet send \
  --from-privkey <hex> \
  --to arcv1... \
  --amount 30000000000 \
  --fee 100000 \
  --node https://seed.archivas.ai
```

Amount is in base units (8 decimals), so 30000000000 = 300 RCHV.

### What are the fees?

Currently **fixed at 0.001 RCHV** per transaction.

Dynamic fee market coming in future version.

### How long do transactions take?

**~20-30 seconds** - transactions are included in the next block after submission.

---

## Technical

### What's the block time?

Target: **20 seconds**  
Actual: **~25 seconds** average  
Range: 1 second to 2 minutes

Difficulty adjusts to maintain stable block times.

### How does difficulty adjustment work?

Simplified EMA (Exponential Moving Average) retargeting:
- If blocks too fast: difficulty increases
- If blocks too slow: difficulty decreases
- Target: ~20-30 second intervals

### What's the current difficulty?

```bash
curl https://seed.archivas.ai/chainTip | jq .difficulty
# "1000000"
```

Difficulty is in the QMAX domain (0 to 1 trillion).

### Is there a max supply?

Not yet defined. Current emission is 20 RCHV per block with no halving schedule implemented.

Mainnet will have a defined economic model.

### What cryptography does Archivas use?

- **Signatures:** Ed25519 (64-byte signatures)
- **Hashing:** Blake2b-256  
- **Key Derivation:** BIP39 + SLIP-0010
- **Addresses:** Bech32 encoding
- **Transaction Hashing:** RFC 8785 canonical JSON + Blake2b

---

## Network & API

### What's the RPC endpoint?

```
https://seed.archivas.ai
```

Open to everyone, no authentication required.

### Is there rate limiting?

Yes:
- `/submit`: 10 requests/minute per IP (burst: 5)
- Other endpoints: No limit

### Does the API support CORS?

Yes! All endpoints include CORS headers for browser apps:
```
Access-Control-Allow-Origin: *
```

### Is there a block explorer?

Yes: https://archivas-explorer-production.up.railway.app

Shows real-time blocks, transactions, and account balances.

### Can I run my own node?

Yes! Build from source:
```bash
git clone https://github.com/ArchivasNetwork/archivas
cd archivas
go build -o archivas-node ./cmd/archivas-node
./archivas-node --rpc 127.0.0.1:8080 --p2p :9090 \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4
```

See [Running a Node](../farmers/running-node.md).

---

## Development

### Is Archivas open source?

Yes! MIT License.

**Repositories:**
- Core: https://github.com/ArchivasNetwork/archivas
- SDK: https://github.com/ArchivasNetwork/archivas-sdk  
- Explorer: https://github.com/ArchivasNetwork/archivas-explorer

### What language is it written in?

**Core blockchain:** Go 1.24+  
**TypeScript SDK:** TypeScript 5+  
**Block Explorer:** Next.js 16 + TypeScript

### Can I contribute?

Absolutely! 

1. Fork the repository
2. Make your changes
3. Submit a pull request

See [How to Contribute](../contributing/how-to-contribute.md).

### Is there a bug bounty?

Not yet. Coming after security audit.

For now, report bugs via GitHub Issues.

---

## Testnet vs Mainnet

### Is this production-ready?

**Testnet:** Yes - stable and operational  
**Mainnet:** Not yet - requires security audit

The testnet is suitable for:
- Development and testing
- Learning about PoST
- Building applications

**Do not** use for storing real value.

### When is mainnet?

**Planned:** Q2-Q3 2026

**Requirements before mainnet:**
- Security audit
- Economic model finalized
- State pruning implemented
- Snapshot sync working
- 10+ geographically distributed farmers
- Advanced VDF (Wesolowski)

### Will testnet RCHV have value?

**No.** Testnet RCHV has no economic value and may be reset at any time.

Mainnet will launch a new chain with a clean genesis.

### Can I keep my testnet wallet?

Your **wallet (mnemonic)** will work on mainnet - same Ed25519 derivation.

Your **balance** will NOT transfer - mainnet starts fresh.

---

## Troubleshooting

### "Connection refused" error

The node might be down. Check:
```bash
curl https://seed.archivas.ai/chainTip
```

If this fails, the seed node is temporarily unavailable. Try again in a few minutes.

### "Invalid signature" error

Your transaction signature is incorrect. Ensure:
- Correct private key for the `from` address
- Correct nonce (matches current account nonce)
- Canonical JSON encoding

### Transaction not confirming

Check:
1. Mempool: `curl https://seed.archivas.ai/mempool`
2. Account nonce: Did it increment?
3. Wait 60 seconds (2-3 blocks)

### Farmer not finding proofs

- Check plot quality (k=28 recommended)
- Verify difficulty isn't too high
- Ensure farmer is connected to node
- Check logs for errors

---

## Get Help

**GitHub Discussions:** https://github.com/ArchivasNetwork/archivas/discussions  
**GitHub Issues:** https://github.com/ArchivasNetwork/archivas/issues  
**Documentation:** This GitBook!

---

**Still have questions?** Ask in [GitHub Discussions](https://github.com/ArchivasNetwork/archivas/discussions)!

