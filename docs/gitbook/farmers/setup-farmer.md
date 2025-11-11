# Setting Up a Farmer

Step-by-step guide to start farming Archivas.

> **Windows Users:** See [Setting Up a Farmer on Windows](setup-farmer-windows.md) for Windows-specific instructions.

---

## Overview

You'll need to:
1. Build the binaries
2. Create plots
3. Start the farmer
4. Earn RCHV!

**Time:** ~30 minutes  
**Difficulty:** Intermediate

---

## Step 1: Install Dependencies

### Linux (Ubuntu/Debian)

```bash
# Update system
sudo apt-get update

# Install Go 1.21+
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### macOS

```bash
# Install Go via Homebrew
brew install go

# Verify
go version
```

---

## Step 2: Clone and Build

```bash
# Clone repository
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Build binaries
go build -o archivas-farmer ./cmd/archivas-farmer
go build -o archivas-cli ./cmd/archivas-cli

# Verify
./archivas-farmer --help
./archivas-cli --help
```

---

## Step 3: Create Your Farmer Identity

The farmer binary manages its own identity keys using **secp256k1** (32-byte private key, 33-byte compressed public key).  
Let the binary generate the keys for you and save them before plotting continues.

```bash
# Prepare a directory for plots first
mkdir -p ~/archivas-plots

# Let the farmer generate its own identity (secp256k1) and start your first plot
./archivas-farmer plot \
  --size 28 \
  --path ~/archivas-plots

# The command immediately prints something like:
#  Generated new farmer identity:
#    Address:     arcv1...
#    Public Key:  02ab... (use for --farmer-pubkey)
#    Private Key: 1f2c... (use for --farmer-privkey)
#  ⚠️  Save both keys! You'll need the private key to farm.
#
# Copy the Address, Public Key, and Private Key to a safe place.
# After the message, the same command continues generating plot-k28.arcv in ~/archivas-plots.
```

If you only needed to record the keys (for example, to plot on another machine), you can press `Ctrl+C` once you have copied them, then rerun the command later with `--farmer-pubkey <saved_key>` to start plotting.

**Important:** Do not mix these keys with the Ed25519 keys produced by `archivas-cli`.  
The farmer currently requires the secp256k1 keys printed by the step above.

---

## Step 4: Create Plots

```bash
# Create additional k28 plots (use the public key captured in Step 3)
./archivas-farmer plot \
  --size 28 \
  --path ~/archivas-plots \
  --farmer-pubkey YOUR_PUBKEY_FROM_STEP_3

# Each run drops a file such as ~/archivas-plots/plot-k28.arcv
# Repeat to add more plots (the farmer auto-increments filenames):
./archivas-farmer plot --size 28 --path ~/archivas-plots --farmer-pubkey YOUR_PUBKEY_FROM_STEP_3
```

**Notes:**
- The `--path` flag points to a directory; the farmer names the plot automatically (`plot-k28.arcv`, `plot-k28-1.arcv`, ...).
- You can run multiple plot commands in parallel if you have CPU and disk headroom.

---

## Step 5: Start Farming

```bash
# Create logs directory
mkdir -p ~/archivas/logs

# Start farmer
./archivas-farmer farm \
  --plots ~/archivas-plots \
  --node https://seed.archivas.ai \
  --farmer-privkey YOUR_PRIVKEY_FROM_STEP_3 \
  > ~/archivas/logs/farmer.log 2>&1 &

# Check it's running
ps aux | grep archivas-farmer

# Watch logs
tail -f ~/archivas/logs/farmer.log
```

**Tip:** The `--farmer-privkey` value is the 64-character hex string printed in Step 3 (32-byte secp256k1 key). If you see “must be 32 bytes”, double-check that you copied the key from the farmer output and not the Ed25519 key from `archivas-cli`.

**Expected output:**
```
👨‍🌾 Farmer Address: arcv1...
📁 Plots Directory: ~/archivas-plots
🌐 Node: https://seed.archivas.ai

✅ Loaded 3 plot(s)
   - plot-k28-1.arcv (k=28, 268435456 hashes)
   - plot-k28-2.arcv (k=28, 268435456 hashes)
   - plot-k28-3.arcv (k=28, 268435456 hashes)

🚜 Starting farming loop...
🔍 NEW HEIGHT 64500 (difficulty: 1000000)
⚙️  Checking plots...
```

---

## Step 6: Monitor Your Earnings

```bash
# Check your balance
curl https://seed.archivas.ai/account/YOUR_ADDRESS

# Watch for wins in logs
tail -f ~/archivas/logs/farmer.log | grep "Found winning"

# Expected when you win:
# 🎉 Found winning proof! Quality: 123456 (target: 1000000)
# ✅ Block submitted successfully for height 64501
```

---

## Running as a Service (Linux)

### Create systemd unit

```bash
sudo nano /etc/systemd/system/archivas-farmer.service
```

**Contents:**
```ini
[Unit]
Description=Archivas Farmer
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME/archivas
ExecStart=/home/YOUR_USERNAME/archivas/archivas-farmer farm \
  --plots /home/YOUR_USERNAME/archivas-plots \
  --node https://seed.archivas.ai \
  --farmer-privkey YOUR_PRIVKEY

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### Enable and start

```bash
sudo systemctl daemon-reload
sudo systemctl enable archivas-farmer
sudo systemctl start archivas-farmer

# Check status
sudo systemctl status archivas-farmer

# View logs
sudo journalctl -u archivas-farmer -f
```

---

## Troubleshooting

### "No plots found"

**Problem:** Farmer can't find plot files.

**Solution:**
```bash
# Check plots exist
ls -lh ~/archivas-plots/

# Verify path in farmer command
./archivas-farmer farm --plots ~/archivas-plots ...
```

### "Connection refused"

**Problem:** Can't reach node.

**Solution:**
- Verify internet connection
- Test node: `curl https://seed.archivas.ai/chainTip`
- Check firewall isn't blocking HTTPS

### "Invalid proof"

**Problem:** Proof rejected by node.

**Solution:**
- Ensure plots were created with correct farmer public key
- Check logs for specific error message
- Verify difficulty target

### Not winning blocks

**Expected behavior** if:
- Network has much more space than you
- You have small plots (k=27)
- Bad luck (probability-based)

**Check:**
- How much space do you have vs network total?
- Are plots loading correctly?
- Is farmer scanning on each challenge?

---

## Performance Optimization

### Faster Plot Scanning

```bash
# Use all CPU cores
./archivas-farmer farm \
  --plots ~/archivas-plots \
  --threads $(nproc) \
  ...
```

### Multiple Plot Directories

```bash
# Combine multiple directories
./archivas-farmer farm \
  --plots /mnt/disk1/plots,/mnt/disk2/plots,/mnt/disk3/plots \
  ...
```

### Reduce I/O

- Use SSD for plots (if possible)
- Avoid network-mounted storage
- Keep plots on local filesystem

---

## Security

### Protect Your Private Key

```bash
# Never commit private keys to git
echo "*.key" >> .gitignore

# Store in environment variable
export FARMER_PRIVKEY="your_hex_key_here"

./archivas-farmer farm \
  --plots ~/archivas-plots \
  --node https://seed.archivas.ai \
  --farmer-privkey $FARMER_PRIVKEY \
  > logs/farmer.log 2>&1 &
```

### Backup Your Mnemonic

- Write down 24 words on paper
- Store in safe place
- Never share with anyone
- Test recovery before relying on it

---

## Scaling Up

### Add More Plots

```bash
# Create additional plots
for i in {4..10}; do
  ./archivas-farmer plot \
    --size 28 \
    --path ~/archivas-plots/plot-k28-$i.arcv \
    --farmer-pubkey YOUR_PUBKEY
done
```

### Monitor Performance

```bash
# Watch farmer logs
tail -f ~/archivas/logs/farmer.log | grep -E "Found winning|Quality|NEW HEIGHT"

# Check balance growth
watch -n 30 'curl -s https://seed.archivas.ai/account/YOUR_ADDRESS | jq .balance'
```

---

## Next Steps

- [Creating Plots](creating-plots.md) - Detailed plotting guide
- [Running a Node](running-node.md) - Optional: run your own node
- [Earnings Guide](earnings.md) - Understand rewards

---

**Start farming!** You're ready to earn RCHV! 🌾

