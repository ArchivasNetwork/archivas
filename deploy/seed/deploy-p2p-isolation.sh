#!/bin/bash
# Deploy P2P Isolation Features to Seed1 (Server A)
# This updates the binary with the new isolation features but doesn't enable them on Seed1

set -e

echo "=========================================="
echo "Seed1 P2P Isolation Feature Deployment"
echo "=========================================="
echo ""

# Step 1: Pull latest code
echo "1️⃣  Pulling latest code..."
cd ~/archivas
git pull

# Step 2: Build new binary
echo "2️⃣  Building archivas-node with P2P isolation features..."
go build -o archivas-node ./cmd/archivas-node

# Verify new flags are available
echo ""
echo "🔍 Verifying new flags:"
./archivas-node --help 2>&1 | grep -E "no-peer-discovery|peer-whitelist|checkpoint" || {
    echo "❌ ERROR: New flags not found!"
    exit 1
}
echo "✅ New flags detected"
echo ""

# Step 3: Stop node
echo "3️⃣  Stopping archivas-node..."
sudo systemctl stop archivas-node

# Step 4: Install new binary
echo "4️⃣  Installing new binary..."
sudo cp archivas-node /usr/local/bin/archivas-node
sudo chmod +x /usr/local/bin/archivas-node

# Step 5: Verify service file (no changes needed for Seed1)
echo "5️⃣  Verifying service configuration..."
if [ ! -f /etc/systemd/system/archivas-node.service ]; then
    echo "⚠️  Service file not found, creating default..."
    sudo tee /etc/systemd/system/archivas-node.service > /dev/null << 'EOF'
[Unit]
Description=Archivas Node (seed.archivas.ai)
Documentation=https://github.com/ArchivasNetwork/archivas
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/archivas
ExecStart=/usr/local/bin/archivas-node \
  --rpc 0.0.0.0:8080 \
  --p2p 0.0.0.0:9090 \
  --genesis /home/ubuntu/archivas/genesis/devnet.genesis.json \
  --network-id archivas-devnet-v4 \
  --db /home/ubuntu/archivas/data \
  --enable-gossip=true

Restart=always
RestartSec=5
StandardOutput=append:/var/log/archivas/node.log
StandardError=append:/var/log/archivas/node-error.log

MemoryAccounting=true
MemoryMax=4G
MemoryHigh=3G
TasksMax=2048
LimitNOFILE=65535
Environment="GOGC=50"

[Install]
WantedBy=multi-user.target
EOF
    sudo systemctl daemon-reload
fi

# Step 6: Start node
echo "6️⃣  Starting archivas-node..."
sudo systemctl start archivas-node

# Step 7: Wait and verify
echo "7️⃣  Waiting for node to start..."
sleep 10

echo ""
echo "📊 Node Status:"
sudo systemctl status archivas-node --no-pager -l | head -15
echo ""

echo "🔍 Checking node health..."
CHAIN_TIP=$(curl -s http://127.0.0.1:8080/chainTip | jq -r .height)
PEER_COUNT=$(curl -s http://127.0.0.1:8080/peers | jq '. | length')

if [ -n "$CHAIN_TIP" ] && [ "$CHAIN_TIP" != "null" ]; then
    echo "✅ Chain height: $CHAIN_TIP"
else
    echo "⚠️  Chain height not available yet"
fi

if [ -n "$PEER_COUNT" ] && [ "$PEER_COUNT" != "null" ]; then
    echo "✅ Connected peers: $PEER_COUNT"
else
    echo "⚠️  Peer count not available yet"
fi

echo ""
echo "=========================================="
echo "✅ Seed1 Deployment Complete!"
echo "=========================================="
echo ""
echo "📋 Notes:"
echo "   - Seed1 runs with normal peer discovery (no isolation)"
echo "   - New P2P isolation features available via flags"
echo "   - Seed2 can now use --no-peer-discovery to isolate"
echo ""
echo "📋 Monitoring Commands:"
echo ""
echo "  # Watch logs:"
echo "  sudo journalctl -u archivas-node -f"
echo ""
echo "  # Check height:"
echo "  curl -s http://127.0.0.1:8080/chainTip | jq"
echo ""
echo "  # Check peers:"
echo "  curl -s http://127.0.0.1:8080/peers | jq"
echo ""

