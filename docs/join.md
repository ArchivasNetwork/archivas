# Join the Archivas Public Testnet

Archivas is a **Proof-of-Space-and-Time blockchain** where anyone can contribute storage to earn RCHV rewards.

**Network:** archivas-testnet v1.0.0  
**Seed:** https://seed.archivas.ai  

---

## Quick Start

```bash
# Install Archivas
curl -sL https://get.archivas.ai/install.sh | bash

# Join the network
archivas join https://seed.archivas.ai
```

**That's it!** You're now mining RCHV! 🌾

---

## Requirements

**Hardware:**
- 50 GB free disk space (for plots)
- 2+ CPU cores
- 4 GB RAM
- Stable internet connection

**Software:**
- Linux (Ubuntu 22.04+ or Debian 12+)
- Go 1.22+ (auto-installed by installer)

---

## What Happens When You Join

1. **Downloads binaries** (node, farmer, timelord)
2. **Generates wallet** (securely stored)
3. **Creates 2x k=27 plots** (~8GB total, takes ~10 minutes)
4. **Connects to network** (syncs blockchain)
5. **Starts farming** (earns RCHV rewards!)

---

## Network Details

| Parameter | Value |
|-----------|-------|
| Network | archivas-testnet |
| Consensus | Proof-of-Space-and-Time |
| Block Time | ~25 seconds |
| Block Reward | 20 RCHV |
| Difficulty | Auto-adjusting |
| Min Plot Size | k=27 (4GB) |

---

## Monitoring

**Your Node:**
- Status: `curl http://localhost:8080/healthz`
- Balance: `curl http://localhost:8080/account/YOUR_ADDRESS`
- Blocks: `curl http://localhost:8080/recentBlocks`

**Network:**
- Explorer: https://seed.archivas.ai:8082
- Registry: https://seed.archivas.ai:8088
- Grafana: https://seed.archivas.ai:3000

---

## Earning Rewards

**How it works:**
1. Your farmer checks plots for winning proofs
2. When you find a winner (quality < difficulty)
3. You submit the proof to the network
4. You earn **20 RCHV** per block! 💰

**With 2x k=27 plots at current difficulty:**
- Expected: ~1-2 blocks per hour
- Earnings: ~40-80 RCHV per hour!

---

## Adding More Capacity

```bash
# Generate more plots
archivas plot create --size 27 --count 4

# Farmer auto-detects new plots
# More plots = more chances to win!
```

---

## Troubleshooting

**No blocks being mined:**
- Check difficulty: `curl http://localhost:8080/healthz`
- Check farmer logs: `tail -f ~/.archivas/logs/farmer.log`
- Ensure plots exist: `ls -lh ~/.archivas/plots/`

**Slow sync:**
- This is normal for initial sync
- Monitor: `watch 'curl -s http://localhost:8080/healthz | jq .height'`

**Need help:**
- GitHub: https://github.com/ArchivasNetwork/archivas/issues
- Discussions: https://github.com/ArchivasNetwork/archivas/discussions

---

**Welcome to Archivas!** 🌾

**Start mining today!** 💎

