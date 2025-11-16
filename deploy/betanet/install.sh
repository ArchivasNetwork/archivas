#!/bin/bash
set -e

# Archivas Betanet Node Installation Script
# Server: 72.251.11.191
# Network: betanet

echo "═══════════════════════════════════════════════════"
echo "   Archivas Betanet Node Installation"
echo "   Server: 72.251.11.191"
echo "═══════════════════════════════════════════════════"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}❌ Please run as root (sudo)${NC}"
    exit 1
fi

echo "📦 Step 1: Installing dependencies..."
apt-get update
apt-get install -y \
    build-essential \
    git \
    curl \
    wget \
    jq \
    ufw \
    snapd

# Install Go
if ! command -v go &> /dev/null; then
    echo "📦 Installing Go..."
    snap install go --classic
else
    echo "✅ Go already installed: $(go version)"
fi

echo ""
echo "👤 Step 2: Creating archivas user..."
if ! id "archivas" &>/dev/null; then
    useradd -r -s /bin/bash -d /opt/archivas -m archivas
    echo "✅ User 'archivas' created"
else
    echo "✅ User 'archivas' already exists"
fi

echo ""
echo "📁 Step 3: Creating directories..."
mkdir -p /etc/archivas/betanet
mkdir -p /var/lib/archivas/betanet
mkdir -p /var/log/archivas
mkdir -p /opt/archivas
mkdir -p /mnt/plots

# Set ownership
chown -R archivas:archivas /var/lib/archivas
chown -R archivas:archivas /var/log/archivas
chown -R archivas:archivas /opt/archivas
chown -R archivas:archivas /mnt/plots

echo "✅ Directories created"

echo ""
echo "📥 Step 4: Cloning Archivas repository..."
if [ ! -d "/opt/archivas/archivas" ]; then
    su - archivas -c "cd /opt/archivas && git clone https://github.com/ArchivasNetwork/archivas.git"
    echo "✅ Repository cloned"
else
    echo "✅ Repository already exists, pulling latest..."
    su - archivas -c "cd /opt/archivas/archivas && git pull"
fi

echo ""
echo "🔨 Step 5: Building Archivas node..."
su - archivas -c "cd /opt/archivas/archivas && go build -o /tmp/archivas-node ./cmd/archivas-node/main.go"
mv /tmp/archivas-node /usr/local/bin/archivas-node
chmod +x /usr/local/bin/archivas-node
echo "✅ Node built and installed to /usr/local/bin/archivas-node"

echo ""
echo "🔨 Step 6: Building Archivas farmer..."
if [ -d "/opt/archivas/archivas/cmd/archivas-farmer" ]; then
    su - archivas -c "cd /opt/archivas/archivas && go build -o /tmp/archivas-farmer ./cmd/archivas-farmer/main.go"
    mv /tmp/archivas-farmer /usr/local/bin/archivas-farmer
    chmod +x /usr/local/bin/archivas-farmer
    echo "✅ Farmer built and installed to /usr/local/bin/archivas-farmer"
else
    echo "⚠️  Farmer source not found, skipping"
fi

echo ""
echo "📋 Step 7: Installing configuration files..."
# Copy genesis file
cp /opt/archivas/archivas/configs/genesis-betanet.json /etc/archivas/betanet/
echo "✅ Genesis file installed"

# Copy config.toml
if [ -f "/opt/archivas/archivas/deploy/betanet/config.toml" ]; then
    cp /opt/archivas/archivas/deploy/betanet/config.toml /etc/archivas/betanet/
    echo "✅ Config file installed"
fi

echo ""
echo "🔥 Step 8: Configuring firewall..."
ufw --force enable
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp comment 'SSH'
ufw allow 30303/tcp comment 'Archivas P2P TCP'
ufw allow 30303/udp comment 'Archivas P2P UDP'
ufw allow 8545/tcp comment 'Archivas RPC'
ufw status numbered
echo "✅ Firewall configured"

echo ""
echo "⚙️  Step 9: Installing systemd services..."
cp /opt/archivas/archivas/deploy/betanet/archivas-betanet.service /etc/systemd/system/
systemctl daemon-reload
echo "✅ Node service installed"

if [ -f "/opt/archivas/archivas/deploy/betanet/archivas-betanet-farmer.service" ]; then
    cp /opt/archivas/archivas/deploy/betanet/archivas-betanet-farmer.service /etc/systemd/system/
    echo "✅ Farmer service installed"
fi

echo ""
echo "═══════════════════════════════════════════════════"
echo "   ✅ Installation Complete!"
echo "═══════════════════════════════════════════════════"
echo ""
echo "📝 Next Steps:"
echo ""
echo "1. Bootstrap from snapshot:"
echo "   sudo -u archivas /usr/local/bin/archivas-node bootstrap --network betanet --db /var/lib/archivas/betanet"
echo ""
echo "2. Start the node:"
echo "   sudo systemctl start archivas-betanet"
echo "   sudo systemctl enable archivas-betanet"
echo ""
echo "3. Check logs:"
echo "   sudo journalctl -u archivas-betanet -f"
echo ""
echo "4. Check status:"
echo "   sudo systemctl status archivas-betanet"
echo ""
echo "5. (Optional) Configure farmer:"
echo "   - Edit /etc/systemd/system/archivas-betanet-farmer.service"
echo "   - Replace FARMER_ADDRESS_HERE with your address"
echo "   - sudo systemctl start archivas-betanet-farmer"
echo ""
echo "6. Verify node:"
echo "   curl http://localhost:8545/eth -X POST -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"eth_chainId\",\"params\":[],\"id\":1}'"
echo ""

