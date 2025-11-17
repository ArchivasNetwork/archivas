# Archivas Betanet Seed 2 Deployment

**Server:** 57.129.96.158  
**Network:** Archivas Betanet (Chain ID 1644)  
**Role:** Seed Node + Farmer  

---

## 🚀 Quick Start

### 1. Copy Deployment Files to Server

```bash
# From your local machine
scp -r deploy/betanet-seed2 ubuntu@57.129.96.158:~/
```

### 2. Run Installation

```bash
# On server 57.129.96.158
ssh ubuntu@57.129.96.158
cd ~/betanet-seed2
chmod +x install.sh verify.sh status.sh
sudo ./install.sh
```

### 3. Generate Wallet

```bash
archivas-wallet new
```

**Save the output!** You'll need:
- Private Key (for farming)
- Public Key (your farmer identity)
- ARCV Address (for receiving rewards)

### 4. Configure Farmer

```bash
sudo nano /etc/systemd/system/archivas-betanet-farmer.service
```

Replace `YOUR_PRIVATE_KEY_HERE` with your actual private key.

### 5. Start Node

```bash
sudo systemctl start archivas-betanet
sudo systemctl enable archivas-betanet
sudo journalctl -u archivas-betanet -f
```

Wait for sync to complete (check with `./status.sh`).

### 6. Create Plots

```bash
cd /mnt/plots
sudo -u archivas archivas-farmer plot --size 20
```

This creates a k=20 plot (~100GB). Repeat for more plots.

### 7. Start Farming

```bash
sudo systemctl start archivas-betanet-farmer
sudo systemctl enable archivas-betanet-farmer
sudo journalctl -u archivas-betanet-farmer -f
```

### 8. Verify Everything

```bash
./verify.sh
./status.sh
```

---

## 📡 Node Configuration

- **RPC Endpoint:** `http://57.129.96.158:8545`
- **P2P Endpoint:** `57.129.96.158:30303`
- **Bootstrap:** `seed1.betanet.archivas.ai:30303`
- **Data Dir:** `/var/lib/archivas/betanet`
- **Plots Dir:** `/mnt/plots`

---

## 🔍 Monitoring

### Check Node Status
```bash
./status.sh
```

### Check Sync Progress
```bash
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

### Check Farmer Activity
```bash
sudo journalctl -u archivas-betanet-farmer -n 50
```

### Check Network Peers
```bash
sudo journalctl -u archivas-betanet | grep "peer"
```

---

## 🛠️ Troubleshooting

### Node Won't Start
```bash
sudo systemctl status archivas-betanet
sudo journalctl -u archivas-betanet -n 100
```

### Can't Connect to Seed 1
```bash
ping seed1.betanet.archivas.ai
curl http://seed1.betanet.archivas.ai:8545
sudo ufw status
```

### Farmer Not Running
```bash
# Check if private key is set
sudo systemctl cat archivas-betanet-farmer

# Check if plots exist
ls -lh /mnt/plots/*.plot

# Check logs
sudo journalctl -u archivas-betanet-farmer -n 50
```

### Sync Taking Too Long
Check the sync status difference:
```bash
./status.sh
```

If far behind, ensure:
- Good network connection to Seed 1
- Sufficient disk I/O performance
- No firewall blocking P2P

---

## 🔐 Security

### Firewall Rules (Applied by install.sh)
```bash
ufw allow ssh
ufw allow 30303/tcp  # P2P
ufw allow 30303/udp  # P2P
ufw allow 8545/tcp   # RPC
ufw enable
```

### File Permissions
- All Archivas files owned by `archivas:archivas`
- Services run as `archivas` user (not root)
- Private keys stored in systemd service files (600 permissions)

---

## 📊 Expected Behavior

1. **Initial Sync (0-30 min)**
   - Node connects to Seed 1
   - Downloads blocks from genesis
   - Validates all block signatures
   - Builds local state

2. **Synced (30+ min)**
   - Node catches up to Seed 1
   - Validates new blocks in real-time
   - P2P gossip active

3. **Farming (once synced + plots ready)**
   - Farmer scans plots for proofs
   - Submits winning proofs to node
   - Produces blocks when eligible
   - Earns RCHV rewards

---

## 📞 Support

- **GitHub:** https://github.com/ArchivasNetwork/archivas
- **Docs:** https://docs.archivas.ai
- **Explorer:** https://betanet.archivas.ai

---

## ✅ Checklist

- [ ] Installation complete
- [ ] Wallet generated and backed up
- [ ] Farmer service configured with private key
- [ ] Node started and syncing
- [ ] Firewall configured
- [ ] Plots created
- [ ] Farmer started
- [ ] Verification passed (`./verify.sh`)
- [ ] Status dashboard working (`./status.sh`)

---

**🎉 Once all checks pass, Seed 2 is operational!**

