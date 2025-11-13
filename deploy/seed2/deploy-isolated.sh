#!/bin/bash
# Deploy Seed2 with P2P Isolation and Chain Checkpoint Validation
# This script sets up a fully isolated secondary seed node that only connects to trusted peers

set -e

echo "=========================================="
echo "Seed2 Isolated Deployment"
echo "=========================================="
echo ""

# Configuration
SEED1_ADDR="seed.archivas.ai:9090"
SEED1_IP="57.129.148.132:9090"
CHECKPOINT_HEIGHT=671992
CHECKPOINT_HASH="eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79"
DATA_DIR="/home/ubuntu/archivas/data"
LOG_DIR="/var/log/archivas"

echo "📋 Configuration:"
echo "   Trusted peer: $SEED1_ADDR"
echo "   Checkpoint: block $CHECKPOINT_HEIGHT"
echo "   Hash: ${CHECKPOINT_HASH:0:16}..."
echo ""

# Step 1: Stop existing node
echo "1️⃣  Stopping existing node..."
sudo systemctl stop archivas-node-seed2 2>/dev/null || true

# Step 2: Backup and clear old data
echo "2️⃣  Clearing old blockchain data..."
if [ -d "$DATA_DIR" ]; then
    BACKUP_DIR="${DATA_DIR}_backup_$(date +%Y%m%d_%H%M%S)"
    echo "   Backing up to: $BACKUP_DIR"
    sudo mv "$DATA_DIR" "$BACKUP_DIR" || true
fi
mkdir -p "$DATA_DIR"

# Step 3: Create log directory
echo "3️⃣  Setting up log directory..."
sudo mkdir -p "$LOG_DIR"
sudo chown ubuntu:ubuntu "$LOG_DIR"

# Step 4: Create isolated systemd service
echo "4️⃣  Creating isolated systemd service..."
sudo tee /etc/systemd/system/archivas-node-seed2.service > /dev/null << EOF
[Unit]
Description=Archivas Node (seed2.archivas.ai) - Isolated Mode
Documentation=https://github.com/ArchivasNetwork/archivas
After=network.target

[Service]
User=ubuntu
WorkingDirectory=/home/ubuntu/archivas
ExecStart=/usr/local/bin/archivas-node \\
  --rpc 127.0.0.1:8080 \\
  --p2p 127.0.0.1:9090 \\
  --genesis /home/ubuntu/archivas/genesis/devnet.genesis.json \\
  --network-id archivas-devnet-v4 \\
  --db /home/ubuntu/archivas/data \\
  --no-peer-discovery \\
  --peer-whitelist $SEED1_ADDR \\
  --peer-whitelist $SEED1_IP \\
  --checkpoint-height $CHECKPOINT_HEIGHT \\
  --checkpoint-hash $CHECKPOINT_HASH

Restart=always
RestartSec=5
StandardOutput=append:$LOG_DIR/node.log
StandardError=append:$LOG_DIR/node-error.log

# Resource limits (optimized for 128GB RAM server)
MemoryAccounting=true
MemoryMax=16G
MemoryHigh=12G
TasksMax=4096
LimitNOFILE=65535
Environment="GOGC=50"
Environment="GOMAXPROCS=16"

[Install]
WantedBy=multi-user.target
EOF

# Step 5: Build and install binary
echo "5️⃣  Building and installing archivas-node..."
cd ~/archivas
go build -o archivas-node ./cmd/archivas-node
sudo cp archivas-node /usr/local/bin/archivas-node
sudo chmod +x /usr/local/bin/archivas-node

# Step 6: Reload systemd and start service
echo "6️⃣  Starting isolated node..."
sudo systemctl daemon-reload
sudo systemctl enable archivas-node-seed2
sudo systemctl start archivas-node-seed2

# Step 7: Wait for startup
echo "7️⃣  Waiting for node to start..."
sleep 5

# Step 8: Verify configuration
echo "8️⃣  Verifying configuration..."
echo ""
echo "📊 Node Status:"
sudo systemctl status archivas-node-seed2 --no-pager -l | head -15
echo ""

# Step 9: Check logs for isolation confirmation
echo "9️⃣  Checking isolation logs..."
sleep 2
echo ""
echo "🔍 Recent logs:"
sudo journalctl -u archivas-node-seed2 -n 30 --no-pager | grep -E "gossip|GATER|p2p.*whitelist|Isolation" || echo "   (waiting for isolation logs...)"
echo ""

# Step 10: Show monitoring commands
echo "=========================================="
echo "✅ Deployment Complete!"
echo "=========================================="
echo ""
echo "📋 Monitoring Commands:"
echo ""
echo "  # Watch logs (live):"
echo "  sudo journalctl -u archivas-node-seed2 -f"
echo ""
echo "  # Check isolation is active:"
echo "  sudo journalctl -u archivas-node-seed2 | grep -E 'gossip|GATER|whitelist'"
echo ""
echo "  # Check sync progress:"
echo "  curl -s http://127.0.0.1:8080/chainTip | jq"
echo ""
echo "  # Check connected peers:"
echo "  curl -s http://127.0.0.1:8080/peers | jq"
echo ""
echo "  # Compare heights (Seed2 vs Seed1):"
echo "  echo \"Seed2: \$(curl -s http://127.0.0.1:8080/chainTip | jq -r .height) | Seed1: \$(curl -s https://seed.archivas.ai/chainTip | jq -r .height)\""
echo ""
echo "⚠️  Expected Behavior:"
echo "   - Logs should show: 'peer discovery DISABLED'"
echo "   - Logs should show: 'peer whitelist enabled: 6 entries'"
echo "   - Only 1 peer should connect: seed.archivas.ai"
echo "   - No 'GATER rejected' messages (all autodiscovery disabled)"
echo "   - Sync should progress from 0 to current height without forking"
echo ""

