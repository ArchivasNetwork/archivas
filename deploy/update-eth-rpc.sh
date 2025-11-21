#!/bin/bash
# Deployment script to update Archivas node with ETH RPC support
# Run this on each seed server: Seed1, Seed2, Seed3

set -e

echo "╔════════════════════════════════════════════════════╗"
echo "║   Archivas ETH RPC Update Deployment              ║"
echo "╚════════════════════════════════════════════════════╝"
echo ""

# Step 1: Pull latest code
echo "📥 Pulling latest code..."
cd ~/archivas
git pull origin main
echo "✅ Code updated"
echo ""

# Step 2: Build new binary
echo "🔨 Building new archivas-node binary..."
go build -o /tmp/archivas-node-new cmd/archivas-node/main.go
echo "✅ Binary built"
echo ""

# Step 3: Stop the node service
echo "⏸️  Stopping archivas-betanet service..."
sudo systemctl stop archivas-betanet
echo "✅ Service stopped"
echo ""

# Step 4: Replace binary
echo "🔄 Replacing archivas-node binary..."
sudo cp /tmp/archivas-node-new /usr/local/bin/archivas-node
sudo chmod +x /usr/local/bin/archivas-node
echo "✅ Binary replaced"
echo ""

# Step 5: Start the service
echo "▶️  Starting archivas-betanet service..."
sudo systemctl start archivas-betanet
echo "✅ Service started"
echo ""

# Step 6: Wait for startup
echo "⏳ Waiting for node to start..."
sleep 3
echo ""

# Step 7: Verify status
echo "📊 Checking service status..."
sudo systemctl status archivas-betanet --no-pager -l | head -20
echo ""

# Step 8: Test ETH RPC
echo "🧪 Testing ETH RPC endpoints..."
echo ""
echo "1️⃣  Testing eth_chainId:"
curl -s http://127.0.0.1:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' | jq
echo ""

echo "2️⃣  Testing eth_blockNumber:"
curl -s http://127.0.0.1:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq
echo ""

echo "3️⃣  Testing eth_getBalance (example address):"
curl -s http://127.0.0.1:8545 -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x47ea4b22029c155c835fd0a0b99f8196766f406a","latest"],"id":1}' | jq
echo ""

echo "╔════════════════════════════════════════════════════╗"
echo "║   ✅ Deployment Complete!                         ║"
echo "╚════════════════════════════════════════════════════╝"
echo ""
echo "If all tests show valid JSON responses (not errors),"
echo "then ETH RPC is working correctly!"

