# Farming with a Private Node

This guide covers setting up a private Betanet node with farming capabilities. Private nodes are ideal for serious farmers who want full control over their infrastructure and enhanced security.

## What is a Private Node?

A private node is configured to:

- **Restrict peer connections** using a whitelist
- **Control which nodes can connect** to you
- **Prioritize farming performance** over public participation
- **Maintain a trusted peer network** with known seed nodes

Private nodes are recommended for:
- Large-scale farming operations (1+ TB)
- Farmers who want guaranteed seed connectivity
- Security-conscious operators
- Farms with multiple harvester nodes

---

## Prerequisites

Before starting, ensure you have:

- ✅ [A running Betanet node](running-betanet-node.md)
- ✅ Ubuntu 22.04 LTS or later
- ✅ At least 100 GB free disk space (for plots)
- ✅ Basic command-line knowledge

---

## Part 1: Configure Private Node

### Step 1: Stop Your Current Node (if running)

```bash
# If running in background
pkill archivas-node

# If running in screen
screen -r archivas-node
# Press Ctrl+C
```

### Step 2: Identify Trusted Peers

For Betanet, the trusted seed nodes are:

| Seed | IP Address | DNS | Type |
|------|-----------|-----|------|
| Seed1 | 72.251.11.191 | seed1.betanet.archivas.ai | Private, Farming |
| Seed2 | 57.129.96.158 | seed2.betanet.archivas.ai | Private, Backup |
| Seed3 | 51.89.11.4 | seed3.betanet.archivas.ai | Public Gateway |

### Step 3: Create Private Node Configuration

Create a config file (optional):

```bash
cd ~/archivas
nano private-node-config.toml
```

**Contents:**
```toml
[network]
name = "betanet"
chain_id = "archivas-betanet-1"

[rpc]
listen = "0.0.0.0:8545"  # Listen on all interfaces for local network access

[p2p]
listen = "0.0.0.0:30303"
max_peers = 10
enable_gossip = false  # Disable public peer discovery

# Whitelist only trusted seed nodes
whitelist = [
    "seed1.betanet.archivas.ai:30303",
    "seed2.betanet.archivas.ai:30303",
    "seed3.betanet.archivas.ai:30303"
]

[database]
path = "/var/lib/archivas/betanet"
```

### Step 4: Start Private Node

```bash
cd ~/archivas

# Option 1: Using CLI flags (recommended)
nohup ./archivas-node \
  --network betanet \
  --rpc 0.0.0.0:8545 \
  --p2p 0.0.0.0:30303 \
  --db ~/.archivas/betanet \
  --peer seed1.betanet.archivas.ai:30303 \
  --peer seed2.betanet.archivas.ai:30303 \
  --peer-whitelist seed1.betanet.archivas.ai:30303 \
  --peer-whitelist seed2.betanet.archivas.ai:30303 \
  --peer-whitelist seed3.betanet.archivas.ai:30303 \
  --max-peers 10 \
  --no-peer-discovery > private-node.log 2>&1 &
```

**What these flags do:**
- `--rpc 0.0.0.0:8545` - RPC accessible on local network (for farmer)
- `--peer seed1...` - Explicitly connect to Seed1 and Seed2
- `--peer-whitelist ...` - Only allow connections from these addresses
- `--no-peer-discovery` - Disable automatic peer discovery
- `--max-peers 10` - Limit concurrent connections

**Verify it's running:**
```bash
ps aux | grep archivas-node | grep -v grep
tail -f ~/archivas/private-node.log
```

You should see:
```
[p2p] Starting P2P listener on 0.0.0.0:30303
[p2p] Peer whitelist enabled: 3 allowed addresses
[p2p] connected to peer seed1.betanet.archivas.ai:30303
```

---

## Part 2: Generate Farming Wallet

### Step 1: Build Wallet Tool (if not already done)

```bash
cd ~/archivas
go build -o archivas-wallet cmd/archivas-wallet/main.go
```

### Step 2: Generate New Wallet

```bash
./archivas-wallet new
```

**Output:**
```
🔐 New Archivas Wallet Generated

Address:     arcv1qxy7m8vw2k9...
Public Key:  024aedb4b79bf799cf484a7369b151f6fb4e1988d745e91b6a9fd9d9eb195a7359
Private Key: 1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899

⚠️  KEEP YOUR PRIVATE KEY SECRET! Anyone with access can spend your RCHV.
```

**Important:** 
- ✅ Save the **private key** securely (you'll need it for farming)
- ✅ Save the **public key** (you'll need it for plotting)
- ✅ Save the **ARCV address** (this is where rewards go)
- ⚠️ **NEVER** share your private key with anyone!

### Step 3: Backup Your Wallet

```bash
# Create secure backup directory
mkdir -p ~/archivas-backup
chmod 700 ~/archivas-backup

# Save wallet details to encrypted file
echo "Address: arcv1qxy7m8vw2k9..." > ~/archivas-backup/wallet.txt
echo "Public Key: 024aedb4b79bf..." >> ~/archivas-backup/wallet.txt
echo "Private Key: 1cb7a7ad1c75b0..." >> ~/archivas-backup/wallet.txt

# Encrypt the file
gpg -c ~/archivas-backup/wallet.txt
rm ~/archivas-backup/wallet.txt  # Remove unencrypted version

# Backup to USB drive or cloud storage
```

---

## Part 3: Create Plots

Plots are the disk space you dedicate to farming. Larger plots = more chances to win blocks.

### Step 1: Choose Plot Size

| k Value | Plot Size | Generation Time | Recommended For |
|---------|-----------|----------------|-----------------|
| k=20 | ~32 MB | 5 seconds | Testing only |
| k=25 | ~1 GB | 40 seconds | Testing, small farms |
| k=28 | ~8 GB | 5 minutes | Home farming |
| k=30 | ~32 GB | 20 minutes | Serious farming |
| k=32 | ~128 GB | 1.5 hours | Large farms |

**Recommendation for private nodes:** Start with **k=30 (32 GB)** for serious farming.

### Step 2: Create Plot Directory

```bash
# Create plots directory on a fast drive
mkdir -p /mnt/plots  # Or ~/plots if no separate drive

# Check available space
df -h /mnt/plots
```

### Step 3: Build Farmer Binary

```bash
cd ~/archivas
go build -o archivas-farmer cmd/archivas-farmer/main.go
```

### Step 4: Generate Plot

```bash
# For 32 GB plot (k=30)
./archivas-farmer plot \
  --size 30 \
  --path /mnt/plots \
  --farmer-pubkey YOUR_PUBLIC_KEY_HERE
```

**Replace `YOUR_PUBLIC_KEY_HERE`** with the public key from Step 2.

**Expected output:**
```
🌾 Generating plot with k=30 (1073741824 hashes)
📁 Output: /mnt/plots/plot-k30.arcv
👨‍🌾 Farmer: 024aedb4b79bf...

Generated 33554432 / 1073741824 hashes (3.1%)
Generated 67108864 / 1073741824 hashes (6.2%)
...
Generated 1073741824 / 1073741824 hashes (100.0%)

✅ Plot generated successfully in 1234.56s
📊 Plot size: ~32768.00 MB
```

### Step 5: Create Multiple Plots (Optional)

For increased farming power, create multiple plots:

```bash
# Create 3x 32 GB plots (96 GB total)
for i in {1..3}; do
  ./archivas-farmer plot \
    --size 30 \
    --path /mnt/plots \
    --farmer-pubkey YOUR_PUBLIC_KEY_HERE
done
```

**Verify plots:**
```bash
ls -lh /mnt/plots/
```

Expected:
```
plot-k30.arcv  (32 GB)
plot-k30-2.arcv  (32 GB)
plot-k30-3.arcv  (32 GB)
```

---

## Part 4: Start Farming

### Step 1: Start the Farmer

```bash
cd ~/archivas

nohup ./archivas-farmer farm \
  --node http://127.0.0.1:8545 \
  --plots /mnt/plots \
  --farmer-privkey YOUR_PRIVATE_KEY_HERE > farmer.log 2>&1 &
```

**Replace `YOUR_PRIVATE_KEY_HERE`** with your private key from Part 2, Step 2.

### Step 2: Monitor Farmer Activity

```bash
# View live logs
tail -f ~/archivas/farmer.log
```

**Expected output:**
```
🌾 Archivas Farmer
   Node: http://127.0.0.1:8545
   Plots: /mnt/plots
   Loaded plots:
      - plot-k30.arcv (k=30, 1073741824 hashes)
      - plot-k30-2.arcv (k=30, 1073741824 hashes)
      - plot-k30-3.arcv (k=30, 1073741824 hashes)
🚜 Starting farming loop...

🔍 NEW HEIGHT 5200 (difficulty: 1000000)
   Challenge: 4f6822b5086d7210...
⚙️  Checking plots...
   Scanning plot plot-k30.arcv...
   Scanning plot plot-k30-2.arcv...
   Scanning plot plot-k30-3.arcv...
```

### Step 3: Wait for Winning Proof

When you find a winning proof:

```
🎉 Found winning proof! Quality: 826432 (target: 1000000)
✅ Block submitted successfully for height 5201 (VDF t=0)
```

**Congratulations!** You just mined a block! 🎊

---

## Monitoring Your Farm

### Check Farmer Status

```bash
# Is farmer running?
ps aux | grep archivas-farmer | grep -v grep

# View recent activity
tail -50 ~/archivas/farmer.log | grep "Found winning proof"
```

### Check Node Status

```bash
# Current block height
curl -s http://127.0.0.1:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq

# Check balance (requires block explorer or wallet)
# Coming soon...
```

### Monitor Logs

```bash
# Node logs
tail -f ~/archivas/private-node.log

# Farmer logs
tail -f ~/archivas/farmer.log

# Check for errors
grep -i error ~/archivas/farmer.log
```

---

## Performance Optimization

### 1. Use Fast Storage

- **NVMe SSD:** Best for plotting and farming
- **SSD:** Good for farming, acceptable for plotting
- **HDD:** Only for completed plots storage (slow for plotting)

### 2. Optimize Plot Storage

```bash
# Move plots to fastest drive
sudo mkdir -p /mnt/nvme/plots
sudo chown $USER:$USER /mnt/nvme/plots
mv /mnt/plots/*.arcv /mnt/nvme/plots/

# Update farmer command
# --plots /mnt/nvme/plots
```

### 3. Increase System Limits

```bash
# Allow more open files (for many plots)
echo "* soft nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf

# Reboot or re-login for changes to take effect
```

### 4. Monitor Resource Usage

```bash
# CPU and memory
htop

# Disk I/O
iotop

# Network
iftop
```

---

## Troubleshooting

### Farmer Cannot Connect to Node

**Error:** `Failed to connect to http://127.0.0.1:8545`

**Solutions:**
1. Verify node is running:
   ```bash
   ps aux | grep archivas-node
   curl http://127.0.0.1:8545
   ```

2. Check if RPC is listening:
   ```bash
   sudo netstat -tulpn | grep 8545
   ```

3. Ensure node started with `--rpc 0.0.0.0:8545` (not `127.0.0.1`)

### No Winning Proofs Found

**Symptom:** Farmer scans plots but never finds winning proofs

**Causes:**
1. **High network difficulty** - Normal on established networks
2. **Small plots** - Increase plot size (k=30 or higher)
3. **Wrong challenge** - Node might be out of sync

**Solutions:**
1. Check node is synced:
   ```bash
   curl -s http://127.0.0.1:8545 -X POST \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
   ```

2. Add more plots:
   ```bash
   ./archivas-farmer plot --size 30 --path /mnt/plots \
     --farmer-pubkey YOUR_PUBLIC_KEY
   ```

3. Check farmer logs for errors:
   ```bash
   grep -i error ~/archivas/farmer.log
   ```

### Plot Generation Fails

**Error:** `Failed to create plot file: permission denied`

**Solution:**
```bash
# Fix permissions
sudo chown -R $USER:$USER /mnt/plots
chmod 755 /mnt/plots
```

### High CPU Usage

**Symptom:** Farmer uses 100% CPU constantly

**Causes:**
- Plotting in progress (normal)
- Too many plots scanning simultaneously

**Solutions:**
1. Wait for plotting to finish
2. Reduce concurrent scans (not yet configurable)
3. Use lower k value for testing

---

## Security Best Practices

### 1. Protect Your Private Key

```bash
# Store in encrypted file
echo "YOUR_PRIVATE_KEY" | gpg -c > ~/private-key.gpg

# Use in farmer with decryption
PRIVKEY=$(gpg -d ~/private-key.gpg 2>/dev/null)
./archivas-farmer farm --farmer-privkey $PRIVKEY ...
```

### 2. Firewall Configuration

```bash
# Block all incoming except SSH and P2P
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp
sudo ufw enable

# Do NOT open RPC port 8545 to internet
```

### 3. Regular Backups

```bash
# Backup wallet
cp ~/archivas-backup/wallet.txt.gpg /backup/location/

# Backup plots (optional - can regenerate)
rsync -avz /mnt/plots/ /backup/plots/
```

### 4. Monitor for Unauthorized Access

```bash
# Check who's logged in
w

# Check recent logins
last -20

# Check failed login attempts
sudo grep "Failed password" /var/log/auth.log
```

---

## Upgrading Your Node

When new versions are released:

```bash
# Stop services
pkill archivas-farmer
pkill archivas-node

# Update code
cd ~/archivas
git pull origin main

# Rebuild binaries
go build -o archivas-node cmd/archivas-node/main.go
go build -o archivas-farmer cmd/archivas-farmer/main.go

# Restart (database and plots preserved)
nohup ./archivas-node --network betanet ... > private-node.log 2>&1 &
nohup ./archivas-farmer farm ... > farmer.log 2>&1 &
```

---

## Next Steps

- **Scale up:** Add more plots to increase winning chances
- **Monitor rewards:** Track your farming income
- **Join community:** Share experiences with other farmers

---

## Support

Need help with private farming?

- **Discord:** [discord.gg/archivas](https://discord.gg/archivas)
- **Telegram Farmers Group:** [t.me/archivas_farmers](https://t.me/archivas_farmers)
- **GitHub Issues:** [github.com/ArchivasNetwork/archivas/issues](https://github.com/ArchivasNetwork/archivas/issues)

