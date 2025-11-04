# Creating Plots

Detailed guide to plot generation for Archivas farming.

---

## What is a Plot?

A **plot** is a file containing precomputed hash tables used for Proof-of-Space farming.

**Properties:**
- Created once, used forever
- Cannot be faked or compressed
- Size determined by k-parameter
- Tied to your farmer public key

---

## Plot Generation

### Basic Command

```bash
./archivas-farmer plot \
  --size 28 \
  --path ~/archivas-plots/plot-k28.arcv \
  --farmer-pubkey YOUR_PUBLIC_KEY_HEX
```

### Parameters

- `--size` (k-size): Plot size parameter (27-30)
- `--path`: Output file location
- `--farmer-pubkey`: Your public key (from `archivas-cli keygen`)

### Progress

```
🌾 Creating plot k=28
📏 Size: 268435456 hashes (~8 GB)
🔑 Farmer: 6301448565bb68053331f80a62621a9ff0a0ac8d6863010ae5bf02237800e5c8

⏳ Plotting...
[████████████████████████████████] 100% (15m 32s)

✅ Plot created: plot-k28.arcv
💾 Size: 8.59 GB
🎯 Ready to farm!
```

---

## K-Size Selection

### k=27 (Small)
- **Hashes:** 134,217,728
- **Size:** ~4 GB
- **Time:** 5-15 minutes
- **Use:** Testing, small farmers

### k=28 (Recommended)
- **Hashes:** 268,435,456
- **Size:** ~8 GB
- **Time:** 10-30 minutes
- **Use:** Production farming

### k=29 (Large)
- **Hashes:** 536,870,912
- **Size:** ~16 GB
- **Time:** 20-60 minutes
- **Use:** Dedicated farmers

### k=30 (Very Large)
- **Hashes:** 1,073,741,824
- **Size:** ~32 GB
- **Time:** 40-120 minutes
- **Use:** Large-scale operations

**Recommendation:** Start with k=28.

---

## Plotting Performance

### CPU Impact

Plotting is **CPU-intensive**:
- Uses 100% of allocated cores
- More cores = faster plotting
- Can plot multiple files in parallel

**Example timings** (k=28):
- 1 core: ~30 minutes
- 4 cores: ~15 minutes
- 8 cores: ~10 minutes

### Parallel Plotting

```bash
# Plot 3 files simultaneously (if you have 12+ cores)
./archivas-farmer plot --size 28 --path plot-1.arcv --farmer-pubkey KEY &
./archivas-farmer plot --size 28 --path plot-2.arcv --farmer-pubkey KEY &
./archivas-farmer plot --size 28 --path plot-3.arcv --farmer-pubkey KEY &

# Wait for all to complete
wait
```

### Temporary Storage

Use fast storage (SSD) for plotting, then move to slow storage (HDD) for farming:

```bash
# Plot on SSD
./archivas-farmer plot \
  --size 28 \
  --path /tmp/plot-k28.arcv \
  --farmer-pubkey KEY

# Move to HDD for farming
mv /tmp/plot-k28.arcv ~/archivas-plots/
```

---

## Plot Management

### Organizing Plots

```bash
# By size
mkdir -p ~/plots/{k27,k28,k29}
./archivas-farmer plot --size 28 --path ~/plots/k28/plot-1.arcv ...

# By date
mkdir -p ~/plots/2025-11
./archivas-farmer plot --size 28 --path ~/plots/2025-11/plot-1.arcv ...

# By server
mkdir -p ~/plots/server-a
./archivas-farmer plot --size 28 --path ~/plots/server-a/plot-1.arcv ...
```

### Verifying Plots

```bash
# Check plot file size
ls -lh ~/archivas-plots/plot-k28.arcv
# Should be ~8-9 GB for k=28

# Test plot loads correctly
./archivas-farmer farm \
  --plots ~/archivas-plots \
  --node https://seed.archivas.ai \
  --farmer-privkey KEY

# Should show: "✅ Loaded N plot(s)"
```

### Deleting Bad Plots

```bash
# If plot is corrupted or wrong key
rm ~/archivas-plots/bad-plot.arcv

# Recreate
./archivas-farmer plot --size 28 --path ~/archivas-plots/new-plot.arcv ...
```

---

## Batch Plotting

### Create 10 plots

```bash
#!/bin/bash
# create-plots.sh

PUBKEY="your_public_key_here"
PLOT_DIR="$HOME/archivas-plots"

mkdir -p "$PLOT_DIR"

for i in {1..10}; do
  echo "Creating plot $i/10..."
  ./archivas-farmer plot \
    --size 28 \
    --path "$PLOT_DIR/plot-k28-$i.arcv" \
    --farmer-pubkey "$PUBKEY"
done

echo "✅ All plots created!"
ls -lh "$PLOT_DIR"
```

Run:
```bash
chmod +x create-plots.sh
./create-plots.sh
```

---

## Plot Portability

### Moving Plots

Plots are **portable** - you can:
- Move between directories
- Move between servers
- Rename files
- Copy to backup locations

**Plot data is self-contained.**

### Sharing Plots

**DO NOT share plots** between different farmer addresses!

Each plot is tied to ONE farmer public key. If you use the same plot with different keys, it won't work.

---

## Storage Estimates

### How many plots?

| Total Disk | k27 | k28 | k29 |
|------------|-----|-----|-----|
| 50 GB | 12 | 6 | 3 |
| 100 GB | 25 | 12 | 6 |
| 500 GB | 125 | 62 | 31 |
| 1 TB | 250 | 125 | 62 |
| 2 TB | 500 | 250 | 125 |

**Recommendation:** Fill 80% of available space (leave room for OS/logs).

---

## Plot Quality

All correctly-generated plots have the same **probability** of finding proofs.

**There is NO such thing as a "better" or "worse" plot** - it's purely random which plots win.

More plots = more lottery tickets = higher win probability.

---

## Advanced Options

### Custom Plot ID

```bash
# Plots are identified by farmer pubkey by default
# For advanced users, you can specify custom plot IDs

./archivas-farmer plot \
  --size 28 \
  --path plot-custom.arcv \
  --farmer-pubkey KEY \
  --plot-id CUSTOM_HEX
```

**Warning:** Only use if you know what you're doing!

---

## Troubleshooting

### "Out of disk space"

**Problem:** Not enough space for plot.

**Solution:**
- Free up disk space
- Use smaller k-size (k=27)
- Use different disk

### "Plot creation failed"

**Problem:** Error during plotting.

**Solution:**
- Check disk isn't full
- Verify farmer-pubkey is valid hex
- Ensure write permissions

### "Farmer can't load plots"

**Problem:** Plots created but farmer doesn't see them.

**Solution:**
```bash
# Check files exist
ls -lh ~/archivas-plots/

# Verify path in farm command
./archivas-farmer farm --plots ~/archivas-plots ...
```

---

## Next Steps

- [Running a Node](running-node.md) - Optional: run your own node
- [Earnings Guide](earnings.md) - Understand rewards
- [Troubleshooting](troubleshooting.md) - Common issues

---

**Plots created?** Start [farming](setup-farmer.md#step-5-start-farming)! 🌾

