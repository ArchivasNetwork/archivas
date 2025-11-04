# Hardware Requirements

What you need to farm Archivas.

---

## Minimum Requirements

**Can farm with:**
- **Disk:** 10 GB free space (1 k28 plot)
- **CPU:** Any modern CPU (1 core, 2 GHz+)
- **RAM:** 2 GB
- **Network:** 1 Mbps connection
- **OS:** Linux, macOS, or Windows

**Example:** Old laptop with spare disk space.

---

## Recommended Setup

**For consistent farming:**
- **Disk:** 100 GB SSD or HDD (10+ k28 plots)
- **CPU:** 4 cores (faster plotting)
- **RAM:** 8 GB
- **Network:** 10 Mbps
- **OS:** Linux (Ubuntu 22.04+)

**Example:** Desktop PC or entry-level VPS.

---

## Optimal Configuration

**For serious farmers:**
- **Disk:** 1+ TB HDD (100+ k28 plots)
- **CPU:** 8+ cores (parallel plotting)
- **RAM:** 16 GB
- **Network:** 100 Mbps
- **OS:** Linux server

**Example:** Dedicated farming server or high-end VPS.

---

## Disk Space Details

### Plot Sizes

| k-size | Hashes | Disk Space | Win Rate (vs k28) |
|--------|--------|------------|-------------------|
| k=27 | 134M | ~4 GB | 50% |
| k=28 | 268M | ~8 GB | 100% (baseline) |
| k=29 | 537M | ~16 GB | 200% |
| k=30 | 1.07B | ~32 GB | 400% |

**Recommended:** k=28 (best size/performance ratio)

### Number of Plots

| Plots | Total Space | Network Share (current) | Blocks/Day |
|-------|-------------|-------------------------|------------|
| 1 k28 | 8 GB | ~15% | ~450 |
| 6 k28 | 48 GB | ~87% | ~2,610 |
| 12 k28 | 96 GB | ~175% | Majority |
| 50 k28 | 400 GB | Dominance | Majority |

*Based on current 55 GB network.*

---

## CPU Requirements

### Plotting
- **Task:** Generate hash tables
- **CPU-bound:** Yes (100% CPU usage)
- **Time:** ~10-30 minutes per k28 plot
- **Parallelization:** Can plot multiple in parallel

### Farming
- **Task:** Scan plots for proofs
- **CPU-bound:** Low (<5% CPU usage)
- **Time:** <1 second per challenge
- **Parallelization:** Plots scanned in parallel

**Recommendation:** 4+ cores for faster plotting, 1 core sufficient for farming.

---

## RAM Requirements

### Plotting
- **Minimum:** 2 GB per concurrent plot
- **Recommended:** 4 GB per concurrent plot

### Farming
- **Minimum:** 1 GB
- **Recommended:** 2 GB
- **Per plot:** ~10 MB overhead

**Example:** 8 GB RAM can plot 2 plots concurrently and farm 50+ plots.

---

## Storage Type

### HDD (Hard Disk Drive)
- ✅ **Cheap:** $15/TB
- ✅ **Large capacity:** 1-20 TB common
- ⚠️ **Slower:** 100-200 MB/s
- ✅ **Good for:** Farming (plots are read-sequentially)

### SSD (Solid State Drive)
- ⚠️ **Expensive:** $50-100/TB
- ⚠️ **Smaller:** 256 GB - 2 TB common
- ✅ **Fast:** 500-7000 MB/s
- ✅ **Good for:** Plotting (faster) and node database

### NVMe SSD
- ⚠️ **Most expensive:** $100-200/TB
- ✅ **Fastest:** 3000-7000 MB/s
- ✅ **Good for:** Fast plotting only

**Recommendation:** Plot on SSD, farm on HDD for best cost/performance.

---

## Network Requirements

### Bandwidth
- **Minimum:** 1 Mbps
- **Recommended:** 10 Mbps
- **Optimal:** 100 Mbps

### Data Usage
- **Block sync:** ~1 MB/day (new blocks)
- **Initial sync:** ~500 MB (60,000 blocks)
- **Total:** <10 GB/month

**Farming is network-light** - uses less bandwidth than video streaming.

---

## Operating System

### Linux (Recommended)
- ✅ **Best performance**
- ✅ **Server-friendly**
- ✅ **Systemd support**
- Tested: Ubuntu 22.04, Debian 12

### macOS
- ✅ **Works well**
- ✅ **Good for testing**
- ⚠️ **Power management** can interfere

### Windows
- ✅ **Supported**
- ⚠️ **Performance** slightly lower
- ⚠️ **Paths** use backslashes

---

## Cloud vs Home Hardware

### Cloud (VPS)
**Pros:**
- 24/7 uptime
- Fast network
- No electricity costs at home

**Cons:**
- Monthly fees ($5-50/month)
- Storage expensive
- Less decentralized

**Providers:** Hetzner, DigitalOcean, Vultr

### Home Hardware
**Pros:**
- One-time cost
- Cheaper long-term
- Supports decentralization

**Cons:**
- Electricity costs
- Uptime depends on you
- Network reliability varies

**Recommendation:** Start with cloud to learn, then home for long-term.

---

## Example Setups

### Budget Setup ($50)
- **Server:** Old laptop or spare PC
- **Disk:** 100 GB free space (reuse existing)
- **Plots:** 10 × k28
- **Cost:** $0 (reuse hardware)

### Mid-Range Setup ($200)
- **Server:** VPS (Hetzner CPX21)
- **Disk:** 500 GB
- **Plots:** 60 × k28
- **Cost:** $10/month

### Dedicated Setup ($500+)
- **Server:** Dedicated server or high-end VPS
- **Disk:** 2 TB HDD
- **Plots:** 250 × k28
- **Cost:** $30-50/month or one-time $500

---

## Performance Tips

### Faster Plotting
- Use SSD for temporary plotting directory
- Use all CPU cores: `--threads $(nproc)`
- Plot multiple k28s in parallel

### Efficient Farming
- Use HDD for plot storage (cheaper)
- Keep plots on fast filesystem (ext4, XFS)
- Avoid network-mounted storage (latency)

### Lower Costs
- Reuse existing hardware
- Buy used HDDs
- Use home server (no VPS fees)

---

## Next Steps

Ready to farm? Continue to:
1. [Setting Up a Farmer](setup-farmer.md)
2. [Creating Plots](creating-plots.md)
3. [Running a Node](running-node.md)

---

**Questions?** See [FAQ](../reference/faq.md) or ask in [GitHub Discussions](https://github.com/ArchivasNetwork/archivas/discussions).

