# Archivas Auto-Join Guide

Quick setup for Timelord + Farmer on any server.

---

## Quick Start

### 1. Run Installer

```bash
cd ~/archivas
sudo bash deploy/autojoin/archivas-autojoin.sh
```

### 2. Configure

```bash
sudo nano /opt/archivas/archivas.env
```

Set:
- `ARCHIVAS_WALLET_ADDR` - Your arcv1... address
- `ARCHIVAS_FARMER_PRIVKEY` - Your private key
- `ARCHIVAS_NODE_RPC` - Node URL (default: http://57.129.148.132:8080)

### 3. Generate Plots (Optional)

```bash
# Create 2x k=28 plots (~8GB each)
sudo bash /opt/archivas/autojoin/create-plots.sh 2 28

# This takes ~8-10 minutes total
```

### 4. Start Services

```bash
# Start Timelord (pinned to CPU 0-3)
sudo systemctl enable --now archivas-timelord

# Start Farmer (pinned to CPU 4-15)
sudo systemctl enable --now archivas-farmer@1

# Check status
sudo systemctl status archivas-timelord
sudo systemctl status archivas-farmer@1
```

### 5. Monitor

```bash
# Timelord logs
tail -f /opt/archivas/logs/timelord.log

# Farmer logs
tail -f /opt/archivas/logs/farmer-1.log | grep -E "🎉|best=|NEW HEIGHT"
```

---

## Adding More Farmers

```bash
# Start additional farmer instances
sudo systemctl enable --now archivas-farmer@2
sudo systemctl enable --now archivas-farmer@3

# Each runs independently, same wallet
```

---

## CPU Pinning

Edit systemd units to change CPU affinity:

```bash
sudo systemctl edit archivas-timelord
# Add: CPUAffinity=0-7

sudo systemctl daemon-reload
sudo systemctl restart archivas-timelord
```

---

## Troubleshooting

**No blocks being mined:**
- Check difficulty: `curl http://NODE:8080/healthz`
- If difficulty > 10M, wait for auto-drop
- Check farmer logs for "best=" outputs

**Timelord not computing:**
- Check if node is reachable
- Verify `ARCHIVAS_NODE_RPC` in archivas.env

**Plots not loading:**
- Check plots.yaml paths match actual directories
- Verify plot files exist: `ls -lh /opt/archivas/plots/*/`

---

## Uninstall

```bash
sudo bash /opt/archivas/autojoin/uninstall-archivas.sh

# With plots deletion:
sudo bash /opt/archivas/autojoin/uninstall-archivas.sh --purge
```

---

**Auto-Join makes it easy to add mining capacity!** 🚜

