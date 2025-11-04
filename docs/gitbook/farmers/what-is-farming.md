# What is Farming?

Farming is the process of earning RCHV by allocating disk space to the Archivas network.

---

## How Farming Works

### 1. Create Plots

**Plots** are large files filled with precomputed hashes.

- Size: Determined by "k-size" parameter
- k=28: ~8 GB per plot (268 million hashes)
- Reusable: Once created, plots work forever

### 2. Run Farmer Software

The farmer continuously:
- Fetches new challenges from the node
- Scans plots for matching proofs
- Submits winning proofs to earn rewards

### 3. Earn Block Rewards

When your proof wins:
- You produce the next block
- Earn 20 RCHV block reward
- Block is added to the chain

---

## Farming vs Mining

| Feature | PoW Mining | PoST Farming |
|---------|------------|--------------|
| **Resource** | Electricity | Disk space |
| **Hardware** | ASIC/GPU | HDD/SSD |
| **Energy Use** | Very high | Very low |
| **Initial Cost** | $1K-$100K | $50-$500 |
| **Reusable** | ❌ Power consumed | ✅ Disk reusable |
| **Noise/Heat** | High | None |
| **Centralization** | Mining pools | Distributed |

**Farming is like mining, but with disk space instead of energy.**

---

## Expected Earnings

Earnings depend on your space vs total network space:

**Formula:**
```
Your earnings ≈ (Your plot size / Total network size) × Total blocks × 20 RCHV
```

**Example:**
- Your plots: 50 GB (6 k28 plots)
- Network total: 55 GB
- Your share: 91%
- Blocks per day: ~3,000
- Your blocks: ~2,730 blocks/day
- **Earnings: ~54,600 RCHV/day**

*Note: Testnet only - no real value.*

---

## Hardware Requirements

### Minimum
- **Disk:** 10 GB free space
- **CPU:** Any modern CPU (1 core)
- **RAM:** 2 GB
- **Network:** 1 Mbps

### Recommended
- **Disk:** 100+ GB SSD/HDD
- **CPU:** 4+ cores (for faster plotting)
- **RAM:** 8 GB
- **Network:** 10 Mbps

### Optimal
- **Disk:** 1 TB+ HDD (many plots)
- **CPU:** 8+ cores
- **RAM:** 16 GB
- **Network:** 100 Mbps

---

## Farming Economics

### Costs
- **Electricity:** ~5W per HDD (~$0.01/day)
- **Disk:** $15/TB one-time
- **Network:** Negligible bandwidth

### Rewards
- **20 RCHV** per block won
- **Frequency:** Depends on your plot size vs network

### Profitability
**Testnet:** No economic value (testing only)  
**Mainnet:** TBD based on RCHV market price

---

## Farming Strategies

### Small Farmer (10-100 GB)
- **Strategy:** Single k28 plot, learn the system
- **Expected:** Win blocks occasionally
- **Cost:** <$50 (reuse existing disk)

### Medium Farmer (100 GB - 1 TB)
- **Strategy:** 5-10 k28 plots
- **Expected:** Win blocks regularly
- **Cost:** $50-$150

### Large Farmer (1 TB+)
- **Strategy:** 50+ k28 plots
- **Expected:** Win majority of blocks
- **Cost:** $150+

---

## Getting Started

**Next steps:**
1. [Check Hardware Requirements](hardware-requirements.md)
2. [Setup a Farmer](setup-farmer.md)
3. [Create Your First Plot](creating-plots.md)

---

## FAQ

**Q: Can I farm on my laptop?**  
A: Yes! Even a 10 GB plot can win blocks occasionally.

**Q: Do I need to keep my computer on 24/7?**  
A: Only while actively farming. Plots persist, so you can farm part-time.

**Q: Can I farm Chia and Archivas together?**  
A: Yes! But they use different plot formats (not compatible).

**Q: How often will I win blocks?**  
A: Depends on your plot size. With 50 GB on a 55 GB network, you'd win ~91% of blocks.

---

**Ready to farm?** Continue to [Hardware Requirements](hardware-requirements.md)!

