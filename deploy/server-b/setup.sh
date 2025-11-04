#!/usr/bin/env bash
# setup-server-b.sh - Deploy Archivas node + farmer on Server B

set -euo pipefail

echo "🌾 Archivas Server B Setup"
echo "=========================="
echo ""

# Check if running on the right server
if [[ ! -d "/home/ubuntu" ]]; then
    echo "❌ This script must run on Ubuntu with /home/ubuntu"
    exit 1
fi

# Step 1: Clone repository
echo "1️⃣  Cloning Archivas repository..."
cd /home/ubuntu
if [[ -d "archivas" ]]; then
    echo "   Repository exists, pulling latest..."
    cd archivas
    git pull origin main
else
    git clone https://github.com/ArchivasNetwork/archivas.git
    cd archivas
fi
echo "   ✅ Repository ready"
echo ""

# Step 2: Build binaries
echo "2️⃣  Building binaries..."
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-farmer ./cmd/archivas-farmer
echo "   ✅ Binaries built"
echo ""

# Step 3: Create directories
echo "3️⃣  Creating directories..."
mkdir -p ~/archivas/data
mkdir -p ~/archivas/logs
mkdir -p ~/archivas/plots-b
echo "   ✅ Directories created"
echo ""

# Step 4: Generate wallet (if needed)
echo "4️⃣  Wallet setup..."
if [[ ! -f ~/archivas/wallet-b.txt ]]; then
    echo "   Generating new wallet..."
    ./archivas-cli keygen > ~/archivas/wallet-b.txt
    echo "   ✅ Wallet saved to ~/archivas/wallet-b.txt"
    echo "   📝 SAVE THIS FILE SECURELY!"
    cat ~/archivas/wallet-b.txt
else
    echo "   ✅ Wallet exists at ~/archivas/wallet-b.txt"
fi
echo ""

# Step 5: Install systemd services
echo "5️⃣  Installing systemd services..."
sudo cp deploy/server-b/archivas-node-b.service /etc/systemd/system/
sudo cp deploy/server-b/archivas-farmer-b.service /etc/systemd/system/

# Ask user to update farmer privkey
echo "   ⚠️  IMPORTANT: Edit farmer service with your private key:"
echo "   sudo nano /etc/systemd/system/archivas-farmer-b.service"
echo "   Replace FARMER_PRIVKEY_HERE with the PrivKey from wallet-b.txt"
echo ""
read -p "   Press Enter after you've updated the private key..."

sudo systemctl daemon-reload
echo "   ✅ Services installed"
echo ""

# Step 6: Start node
echo "6️⃣  Starting node..."
sudo systemctl enable archivas-node-b
sudo systemctl start archivas-node-b
echo "   ✅ Node started"
echo ""

# Wait for node to initialize
echo "7️⃣  Waiting for node to sync..."
sleep 10
echo "   Checking node status..."
curl -s http://127.0.0.1:8080/chainTip || echo "   ⚠️  Node not ready yet"
echo ""

# Step 7: Start farmer (if plots exist)
echo "8️⃣  Farmer setup..."
if [[ -n "$(ls -A ~/archivas/plots-b 2>/dev/null)" ]]; then
    echo "   Plots found, starting farmer..."
    sudo systemctl enable archivas-farmer-b
    sudo systemctl start archivas-farmer-b
    echo "   ✅ Farmer started"
else
    echo "   ⚠️  No plots found in ~/archivas/plots-b"
    echo "   Create plots with: ./archivas-farmer plot --size 28 --path ~/archivas/plots-b/plot-k28-1.arcv --farmer-pubkey YOUR_PUBKEY"
    echo "   Then: sudo systemctl start archivas-farmer-b"
fi
echo ""

# Step 8: Firewall
echo "9️⃣  Configuring firewall..."
if command -v ufw &> /dev/null; then
    sudo ufw allow 9090/tcp comment "Archivas P2P"
    sudo ufw deny 8080/tcp comment "Block external RPC"
    echo "   ✅ Firewall configured"
else
    echo "   ⏭️  UFW not installed, skipping firewall"
fi
echo ""

# Final status
echo "✅ Server B Setup Complete!"
echo ""
echo "Services:"
echo "  - archivas-node-b: sudo systemctl status archivas-node-b"
echo "  - archivas-farmer-b: sudo systemctl status archivas-farmer-b"
echo ""
echo "Verify peering:"
echo "  curl http://127.0.0.1:8080/chainTip"
echo "  curl https://seed.archivas.ai/chainTip"
echo "  (Heights should match)"
echo ""
echo "View logs:"
echo "  tail -f /var/log/archivas-node-b.log"
echo "  tail -f /var/log/archivas-farmer-b.log"

