#!/bin/bash
set -e

echo "════════════════════════════════════════════════════════════════"
echo "  🔧 Fixing Private Node - Removing Seed1 Corruption"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Stop the node
echo "[INFO] Stopping archivas-node-private..."
sudo systemctl stop archivas-node-private

# Clear corrupted database
echo "[INFO] Clearing corrupted database..."
sudo rm -rf /home/ubuntu/archivas/data/*

# Backup current service file
echo "[INFO] Backing up service file..."
sudo cp /etc/systemd/system/archivas-node-private.service /etc/systemd/system/archivas-node-private.service.backup

# Edit service file to remove Seed1 peer and checkpoint
echo "[INFO] Editing service file to ONLY use Seed2..."
sudo tee /etc/systemd/system/archivas-node-private.service > /dev/null << 'EOF'
[Unit]
Description=Archivas Private Node - Full Node for Local/Network Use
Documentation=https://docs.archivas.ai/farmers/private-node
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu/archivas

# Node configuration
ExecStart=/home/ubuntu/archivas/archivas-node \
  -network-id archivas-devnet-v4 \
  -db /home/ubuntu/archivas/data \
  -rpc 127.0.0.1:8080 \
  -p2p 0.0.0.0:9090 \
  -genesis /home/ubuntu/archivas/genesis/devnet.genesis.json \
  -no-peer-discovery \
  -peer seed2.archivas.ai:9090 \
  -max-peers 50

# Resource limits
LimitNOFILE=65536
MemoryMax=4G
MemoryHigh=3G
TasksMax=100

# Go runtime tuning
Environment="GOMAXPROCS=4"
Environment="GOGC=50"

Restart=always
RestartSec=10
StartLimitInterval=300
StartLimitBurst=5

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
echo "[INFO] Reloading systemd..."
sudo systemctl daemon-reload

# Start the node
echo "[INFO] Starting node (syncing from genesis via Seed2 ONLY)..."
sudo systemctl start archivas-node-private

echo ""
echo "════════════════════════════════════════════════════════════════"
echo "  ✅ Node Fixed and Restarted"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "📊 Check Status:"
echo "   sudo systemctl status archivas-node-private"
echo ""
echo "📜 View Logs:"
echo "   sudo journalctl -u archivas-node-private -f"
echo ""
echo "⏰ Wait 2-3 minutes, then verify it's syncing from Seed2:"
echo "   sudo journalctl -u archivas-node-private -n 20 | grep 'from'"
echo ""
echo "   You should see blocks from a DIFFERENT IP (not 57.129.148.132)"
echo ""

