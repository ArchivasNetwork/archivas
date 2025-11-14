# Large-Scale Farming Guide

Guide for farmers with significant disk space (10+ TB).

---

## Overview

If you have **40 TB of free disk space**, you can create a substantial farming operation on Archivas. This guide covers optimal strategies for maximizing your farming potential.

> **🚨 Important for Large Farmers:** If you're farming at scale (10+ TB), we **strongly recommend** running your own [Private Node](private-node.md). Public seed APIs can experience timeouts under heavy load, and a private node gives you:
> - **Sub-second block submissions** (localhost RPC)
> - **No rate limits or 504 errors**
> - **Complete control** over your infrastructure
> - **Better uptime** during network congestion
>
> See the [Private Node Guide](private-node.md) for setup instructions.

---

## Plot Size Recommendations

### For 40 TB of Space

| Plot Size | Plot Size (GB) | Number of Plots | Total Space Used | Recommendation |
|-----------|---------------|-----------------|------------------|----------------|
| **k=28** | ~8 GB | **~5,000 plots** | ~40 TB | ✅ **Best balance** |
| **k=29** | ~16 GB | **~2,500 plots** | ~40 TB | Good for fewer files |
| **k=30** | ~32 GB | **~1,250 plots** | ~40 TB | Fewer files, slower scanning |

**Recommendation for 40 TB:** Use **k=28 plots** (~5,000 plots)
- Best balance of file count vs. scanning performance
- Faster plot creation
- Easier to manage and organize
- More "lottery tickets" = higher win probability

---

## Plotting Strategy

### Phase 1: Initial Setup (First 1-2 TB)

Start with a smaller batch to test your setup:

```bash
# Create first 100 plots (k=28, ~800 GB)
for i in {1..100}; do
  ./archivas-farmer plot \
    --size 28 \
    --path /mnt/plots/plot-k28-$i.arcv \
    --farmer-pubkey YOUR_PUBKEY &
  
  # Limit to 4 parallel plots (adjust based on CPU)
  if (( i % 4 == 0 )); then
    wait
  fi
done
```

**Why start small:**
- Verify your hardware can handle plotting
- Test disk I/O performance
- Ensure farmer loads plots correctly
- Start earning while you continue plotting

### Phase 2: Full Deployment (Remaining 38-39 TB)

Once verified, scale up:

```bash
# Create remaining plots in batches
# Adjust parallel count based on CPU cores and disk I/O

# Example: 8 parallel plots (for 16+ core CPU)
for i in {101..5000}; do
  ./archivas-farmer plot \
    --size 28 \
    --path /mnt/plots/plot-k28-$i.arcv \
    --farmer-pubkey YOUR_PUBKEY &
  
  if (( i % 8 == 0 )); then
    wait
    echo "Created $i plots so far..."
  fi
done
```

---

## Hardware Considerations

### CPU Requirements

**For 5,000 k=28 plots:**
- **Plotting:** CPU-intensive, benefits from many cores
- **Farming:** Minimal CPU, mostly disk I/O
- **Recommendation:** 8+ cores for efficient parallel plotting

**Plotting time estimates (k=28):**
- 1 core: ~30 minutes per plot
- 4 cores: ~15 minutes per plot
- 8 cores: ~10 minutes per plot
- 16 cores: ~7 minutes per plot

**Total time for 5,000 plots:**
- 8 cores, 4 parallel: ~2,500 hours (~104 days)
- 16 cores, 8 parallel: ~1,250 hours (~52 days)

### Storage Requirements

**Disk types:**
- **SSD (for plotting):** Fast plotting, then move to HDD
- **HDD (for farming):** Cost-effective for large storage
- **Hybrid approach:** Plot on SSD, farm from HDD

**Disk organization:**
```
/mnt/plots/
├── disk1/  (10 TB)
├── disk2/  (10 TB)
├── disk3/  (10 TB)
└── disk4/  (10 TB)
```

### Memory Requirements

- **Plotting:** 2-4 GB per parallel plot
- **Farming:** ~1-2 GB total (plots loaded on-demand)
- **Recommendation:** 16-32 GB RAM for comfortable operation

---

## Plotting Workflow

### Option 1: Plot on SSD, Move to HDD

```bash
# Plot on fast SSD
./archivas-farmer plot \
  --size 28 \
  --path /tmp/plot-k28-1.arcv \
  --farmer-pubkey YOUR_PUBKEY

# Move to HDD for farming
mv /tmp/plot-k28-1.arcv /mnt/plots/
```

**Benefits:**
- Faster plotting (SSD speed)
- Cost-effective farming (HDD storage)

### Option 2: Direct Plot to Final Location

```bash
# Plot directly to final location
./archivas-farmer plot \
  --size 28 \
  --path /mnt/plots/plot-k28-1.arcv \
  --farmer-pubkey YOUR_PUBKEY
```

**Benefits:**
- Simpler workflow
- No moving files
- Slower plotting (if using HDD)

---

## Organizing 5,000 Plots

### Directory Structure

```bash
# Organize by disk
mkdir -p /mnt/{disk1,disk2,disk3,disk4}/archivas-plots

# Or organize by date/batch
mkdir -p /mnt/plots/{batch-001,batch-002,batch-003,...}
```

### Multiple Plot Directories

The farmer can scan multiple directories:

```bash
./archivas-farmer farm \
  --plots /mnt/disk1/archivas-plots,/mnt/disk2/archivas-plots,/mnt/disk3/archivas-plots,/mnt/disk4/archivas-plots \
  --node https://seed.archivas.ai \
  --farmer-privkey YOUR_PRIVKEY
```

---

## Performance Optimization

### Parallel Plotting

**Find optimal parallel count:**

```bash
# Test with different parallel counts
# Monitor CPU, disk I/O, and memory

# Start conservative (4 parallel)
# Increase if CPU/disk can handle more
# Decrease if system becomes unresponsive
```

**Guidelines:**
- **CPU-bound:** Parallel count = CPU cores
- **I/O-bound:** Parallel count = 2-4x disk count
- **Memory-bound:** Parallel count = RAM / (2-4 GB per plot)

### Plot Scanning Performance

**For 5,000 plots:**
- Scanning all plots takes time
- Farmer scans plots on each new challenge
- More plots = more scanning time

**Optimization tips:**
- Use fast storage (SSD) for plots if possible
- Distribute plots across multiple disks
- The farmer scans plots sequentially - faster disks = faster scanning

```bash
./archivas-farmer farm \
  --plots /mnt/plots \
  --node https://seed.archivas.ai \
  --farmer-privkey YOUR_PRIVKEY
```

---

## Expected Earnings

### Probability Calculation

**With 5,000 k=28 plots:**
- Each plot: 268,435,456 hashes
- Total space: ~40 TB
- Block reward: 20 RCHV per block
- Block time: ~20 seconds

**Win probability per block:**
- Depends on network total space
- More space = higher probability
- Probability = Your Space / Total Network Space

**Example:**
- If network has 400 TB total
- Your share: 40 TB / 400 TB = 10%
- Expected: ~1 block per 10 blocks = 2 blocks per minute
- Earnings: ~40 RCHV per minute = ~2,400 RCHV per hour

**Note:** This is theoretical. Actual earnings depend on:
- Network total space
- Difficulty adjustments
- Luck factor

---

## Monitoring Large Operations

### Plot Creation Progress

```bash
# Count created plots
find /mnt/plots -name "plot-k28-*.arcv" | wc -l

# Check total size
du -sh /mnt/plots

# Monitor plotting processes
ps aux | grep archivas-farmer | grep plot
```

### Farming Status

```bash
# Check farmer is running
ps aux | grep "archivas-farmer farm"

# Monitor logs for wins
tail -f ~/archivas-logs/farmer.log | grep -E "Found winning|Block submitted"

# Check loaded plots
tail -f ~/archivas-logs/farmer.log | grep "Loaded.*plot"
```

### Balance Tracking

```bash
# Check balance periodically
watch -n 60 'curl -s https://seed.archivas.ai/account/YOUR_ADDRESS | jq .balance'
```

---

## Scaling Up Gradually

### Recommended Approach

1. **Week 1:** Create 100-200 plots, start farming
2. **Week 2-4:** Continue plotting while farming
3. **Month 2-3:** Scale to full 5,000 plots
4. **Ongoing:** Monitor and optimize

**Benefits:**
- Start earning immediately
- Test system stability
- Adjust strategy based on results
- Spread plotting load over time

---

## Troubleshooting Large Operations

### "Too many open files"

**Problem:** System limit on open file descriptors.

**Solution:**
```bash
# Increase limit
ulimit -n 100000

# Or add to /etc/security/limits.conf
* soft nofile 100000
* hard nofile 100000
```

### Slow Plot Scanning

**Problem:** Scanning 5,000 plots takes too long.

**Solutions:**
- Use faster storage (SSD)
- Distribute across multiple disks
- Consider k=29 or k=30 for fewer files (fewer files = faster scanning)
- The farmer scans plots sequentially, so faster I/O helps

### Disk I/O Bottleneck

**Problem:** Disk can't keep up with plotting/farming.

**Solutions:**
- Reduce parallel plotting count
- Use separate disks for plotting vs. farming
- Use faster disks (SSD for active plots)
- Stagger plotting operations

---

## Cost-Benefit Analysis

### For 40 TB Operation

**Costs:**
- Storage: 40 TB HDD ~$800-1,200
- Electricity: ~50-100W for farming
- Time: Plotting takes weeks/months

**Benefits:**
- Potential for significant RCHV earnings
- One-time plotting, permanent farming
- Passive income after setup

**ROI depends on:**
- RCHV value
- Network total space
- Your share percentage
- Block rewards

---

## Best Practices

1. **Start small:** Test with 100-200 plots first
2. **Monitor closely:** Watch for errors, performance issues
3. **Backup keys:** Store farmer private key securely
4. **Organize plots:** Use clear directory structure
5. **Document setup:** Keep notes on your configuration
6. **Scale gradually:** Don't rush to fill all 40 TB at once
7. **Optimize continuously:** Adjust based on performance

---

## Next Steps

- [Setting Up a Farmer](setup-farmer.md) - Basic setup guide
- [Creating Plots](creating-plots.md) - Detailed plotting guide
- [Earning & Rewards](earnings.md) - Understand rewards
- [Troubleshooting](troubleshooting.md) - Common issues

---

**Ready to scale up?** Start with a small batch, verify everything works, then gradually fill your 40 TB! 🌾

