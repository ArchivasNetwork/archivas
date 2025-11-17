#!/bin/bash
set -e

echo "=================================================="
echo "  Archivas Betanet Seed 3 (Public Gateway)"
echo "  Server: 51.89.11.4"
echo "  Role: Public P2P/RPC Node (Non-Farming)"
echo "=================================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

echo "📦 Step 1: Installing dependencies..."
apt update
apt install -y git golang-go build-essential net-tools ufw

echo "👤 Step 2: Creating archivas user..."
if ! id -u archivas &>/dev/null; then
    useradd -r -m -s /bin/bash archivas
    echo "✅ User 'archivas' created"
else
    echo "✅ User 'archivas' already exists"
fi

echo "🔨 Step 3: Building binaries..."
cd /home/ubuntu
if [ ! -d "archivas" ]; then
    sudo -u ubuntu git clone https://github.com/ArchivasNetwork/archivas.git
fi
cd archivas
sudo -u ubuntu git pull
sudo -u ubuntu go build -o archivas-node ./cmd/archivas-node

echo "📁 Step 4: Installing binaries..."
cp archivas-node /usr/local/bin/
chmod +x /usr/local/bin/archivas-node
echo "✅ Binary installed: /usr/local/bin/archivas-node"

echo "📂 Step 5: Setting up directories..."
mkdir -p /opt/archivas
mkdir -p /var/lib/archivas/betanet
cp -r /home/ubuntu/archivas/configs /opt/archivas/
chown -R archivas:archivas /opt/archivas
chown -R archivas:archivas /var/lib/archivas
echo "✅ Directories created and permissions set"

echo "🔥 Step 6: Configuring firewall..."
ufw --force enable
ufw allow 22/tcp comment 'SSH'
ufw allow 30303/tcp comment 'Archivas P2P'
ufw allow 30303/udp comment 'Archivas P2P'
ufw allow 8545/tcp comment 'Archivas RPC'
ufw reload
echo "✅ Firewall configured (P2P: 30303, RPC: 8545)"

echo "⚙️  Step 7: Installing systemd service..."
cat > /etc/systemd/system/archivas-betanet.service << 'EOF'
[Unit]
Description=Archivas Betanet Seed 3 (Public Gateway)
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=archivas
Group=archivas
WorkingDirectory=/opt/archivas
ExecStart=/usr/local/bin/archivas-node \
  --network betanet \
  --rpc 0.0.0.0:8545 \
  --p2p 0.0.0.0:30303 \
  --db /var/lib/archivas/betanet \
  --genesis /opt/archivas/configs/genesis-betanet.json \
  --peer seed1.betanet.archivas.ai:30303 \
  --peer seed2.betanet.archivas.ai:30303 \
  --max-peers 100 \
  --enable-gossip

# Performance
LimitNOFILE=65536
LimitNPROC=4096

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=true

# Restart policy
Restart=always
RestartSec=10s
StartLimitInterval=600s
StartLimitBurst=5

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=archivas-betanet-seed3

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable archivas-betanet
echo "✅ Service installed and enabled"

echo "🚀 Step 8: Starting Archivas Betanet Seed 3..."
systemctl start archivas-betanet
sleep 3

echo ""
echo "=================================================="
echo "  ✅ Installation Complete!"
echo "=================================================="
echo ""
echo "📊 Check status:"
echo "   sudo systemctl status archivas-betanet"
echo ""
echo "📜 View logs:"
echo "   sudo journalctl -u archivas-betanet -f"
echo ""
echo "🔍 Verify identity:"
echo "   bash verify.sh"
echo ""
echo "🌐 Node Info:"
echo "   - Network: Betanet (1644)"
echo "   - RPC: http://51.89.11.4:8545"
echo "   - P2P: 51.89.11.4:30303"
echo "   - Role: Public Gateway (Non-Farming)"
echo ""

