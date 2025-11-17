#!/bin/bash
set -e

echo "=============================================="
echo "Archivas Betanet Seed 2 Installation"
echo "Server: 57.129.96.158"
echo "=============================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

echo "📋 Step 1: Creating archivas user..."
if ! id archivas &>/dev/null; then
    useradd -r -m -s /bin/bash archivas
    echo "✅ User created"
else
    echo "✅ User already exists"
fi

echo ""
echo "📦 Step 2: Installing dependencies..."
apt update
apt install -y git golang-go build-essential net-tools ufw curl jq

echo ""
echo "🔧 Step 3: Building Archivas binaries..."
cd /home/ubuntu
if [ -d "archivas" ]; then
    cd archivas
    git fetch
    git pull
else
    git clone https://github.com/ArchivasNetwork/archivas.git
    cd archivas
fi

echo "Building archivas-node..."
go build -o archivas-node ./cmd/archivas-node

echo "Building archivas-farmer..."
go build -o archivas-farmer ./cmd/archivas-farmer

echo "Building archivas-wallet..."
go build -o archivas-wallet ./cmd/archivas-wallet

echo "Installing binaries..."
cp archivas-node /usr/local/bin/
cp archivas-farmer /usr/local/bin/
cp archivas-wallet /usr/local/bin/
chmod +x /usr/local/bin/archivas-node
chmod +x /usr/local/bin/archivas-farmer
chmod +x /usr/local/bin/archivas-wallet

echo ""
echo "📁 Step 4: Setting up directories..."
mkdir -p /opt/archivas
mkdir -p /var/lib/archivas/betanet
mkdir -p /mnt/plots
mkdir -p /etc/archivas/betanet

chown -R archivas:archivas /opt/archivas
chown -R archivas:archivas /var/lib/archivas
chown -R archivas:archivas /mnt/plots

echo ""
echo "📝 Step 5: Installing configuration..."
cp deploy/betanet-seed2/config.toml /etc/archivas/betanet/config.toml

echo ""
echo "🔥 Step 6: Configuring firewall..."
ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 30303/tcp comment 'Archivas P2P'
ufw allow 30303/udp comment 'Archivas P2P'
ufw allow 8545/tcp comment 'Archivas RPC'
ufw --force enable

echo ""
echo "🎯 Step 7: Installing systemd services..."
cp deploy/betanet-seed2/archivas-betanet.service /etc/systemd/system/
cp deploy/betanet-seed2/archivas-betanet-farmer.service /etc/systemd/system/
systemctl daemon-reload

echo ""
echo "=============================================="
echo "✅ Installation Complete!"
echo "=============================================="
echo ""
echo "📌 Next Steps:"
echo ""
echo "1. Generate a new wallet for Seed 2:"
echo "   archivas-wallet new"
echo ""
echo "2. Update farmer service with private key:"
echo "   sudo nano /etc/systemd/system/archivas-betanet-farmer.service"
echo "   (Replace YOUR_PRIVATE_KEY_HERE)"
echo ""
echo "3. Start the node:"
echo "   sudo systemctl start archivas-betanet"
echo "   sudo systemctl enable archivas-betanet"
echo ""
echo "4. Monitor sync progress:"
echo "   sudo journalctl -u archivas-betanet -f"
echo ""
echo "5. Once synced, create plots:"
echo "   cd /mnt/plots"
echo "   archivas-farmer plot --size 20"
echo ""
echo "6. Start farming:"
echo "   sudo systemctl start archivas-betanet-farmer"
echo "   sudo systemctl enable archivas-betanet-farmer"
echo ""
echo "7. Run verification:"
echo "   ./deploy/betanet-seed2/verify.sh"
echo ""




