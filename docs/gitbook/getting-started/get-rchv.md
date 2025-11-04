# Get Testnet RCHV

How to get RCHV on the public testnet for testing.

---

## Option 1: Faucet (Coming Soon)

Public faucet will be available at:
```
https://faucet.archivas.ai
```

**Features:**
- 20 RCHV per request
- Rate limited to 1 request per hour per IP
- Instant delivery (~20 seconds)

**Status:** 🔄 Under development

---

## Option 2: Farm RCHV

Earn RCHV by farming:

1. **Setup a farmer** - See [Farmer Setup Guide](../farmers/setup-farmer.md)
2. **Create plots** - k=28 recommended (~8 GB per plot)
3. **Connect to seed.archivas.ai**
4. **Win blocks** - Earn 20 RCHV per block

**Advantages:**
- Unlimited RCHV (as long as you win blocks)
- Learn about PoSpace farming
- Support the network

---

## Option 3: Request from Community

**For now:**

1. Join [GitHub Discussions](https://github.com/ArchivasNetwork/archivas/discussions)
2. Post your testnet address
3. Community members can send you RCHV

**Example request:**
```
Hi! Testing Archivas development.
Address: arcv1...
Use case: Building a web wallet
Amount needed: 100 RCHV for testing transfers

Thanks!
```

---

## Verify You Received RCHV

```bash
# Check your balance
curl https://seed.archivas.ai/account/YOUR_ADDRESS

# Should show:
# {
#   "address": "arcv1...",
#   "balance": "10000000000",  # 100 RCHV in base units
#   "nonce": "0"
# }
```

---

## Next Steps

Once you have testnet RCHV:
- [Send Your First Transaction](first-transaction.md)
- [Build with the SDK](../developers/sdk-guide.md)
- [Explore the API](../developers/api-reference.md)

---

**Remember:** Testnet RCHV has no economic value and the chain may reset at any time!

