#!/bin/bash

RPC_URL="http://localhost:8545"

echo "=================================================="
echo "  Archivas Betanet Seed 3 - Status Dashboard"
echo "  Role: Public Gateway (Non-Farming)"
echo "=================================================="
echo ""

# Service Status
echo "🔧 Service Status:"
systemctl status archivas-betanet --no-pager | head -3
echo ""

# Network Info
echo "🌐 Network Information:"
CHAIN_ID=$(curl -s -X POST "$RPC_URL/" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
    | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
CHAIN_ID_DEC=$((16#${CHAIN_ID#0x}))

NETWORK_ID=$(curl -s -X POST "$RPC_URL/" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
    | grep -o '"result":"[^"]*"' | cut -d'"' -f4)

echo "   Chain ID: $CHAIN_ID ($CHAIN_ID_DEC)"
echo "   Network ID: $NETWORK_ID"
echo "   Network: Betanet"
echo ""

# Block Height
echo "📊 Blockchain Status:"
BLOCK_HEX=$(curl -s -X POST "$RPC_URL/" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
BLOCK_HEIGHT=$((16#${BLOCK_HEX#0x}))

echo "   Current Height: $BLOCK_HEIGHT"
echo ""

# P2P Status
echo "🔗 P2P Status:"
PEER_LOG=$(sudo journalctl -u archivas-betanet --since "1 minute ago" -n 100)
CONNECTED_PEERS=$(echo "$PEER_LOG" | grep "total peers:" | tail -1 | grep -o "total peers: [0-9]*" | cut -d' ' -f3 || echo "0")
LAST_PEER=$(echo "$PEER_LOG" | grep "connected to peer" | tail -1 | grep -o "peer [^ ]*" | cut -d' ' -f2 || echo "none")

echo "   Connected Peers: $CONNECTED_PEERS"
echo "   Last Connected: $LAST_PEER"
echo "   Max Peers: 100 (public gateway)"
echo ""

# Recent Activity
echo "📡 Recent Activity (last 20 lines):"
sudo journalctl -u archivas-betanet -n 20 --no-pager | grep -E '\[p2p\]|\[IBD\]|\[consensus\]|\[rpc\]' | tail -10
echo ""

# Port Status
echo "🔌 Port Status:"
echo -n "   P2P (30303): "
if netstat -tuln | grep -q ":30303"; then
    echo "✅ LISTENING"
else
    echo "❌ NOT LISTENING"
fi

echo -n "   RPC (8545): "
if netstat -tuln | grep -q ":8545"; then
    echo "✅ LISTENING"
else
    echo "❌ NOT LISTENING"
fi
echo ""

# Resource Usage
echo "💻 Resource Usage:"
CPU=$(ps aux | grep archivas-node | grep -v grep | awk '{print $3}')
MEM=$(ps aux | grep archivas-node | grep -v grep | awk '{print $4}')
echo "   CPU: ${CPU:-0}%"
echo "   Memory: ${MEM:-0}%"
echo ""

# Disk Usage
echo "💾 Disk Usage:"
du -sh /var/lib/archivas/betanet 2>/dev/null || echo "   No data yet"
echo ""

# Public Endpoints
echo "🌍 Public Endpoints:"
echo "   RPC: http://51.89.11.4:8545"
echo "   P2P: seed3.betanet.archivas.ai:30303"
echo "   DNS: seed3.betanet.archivas.ai → 51.89.11.4"
echo ""

echo "=================================================="
echo "  Commands:"
echo "    Restart: sudo systemctl restart archivas-betanet"
echo "    Logs:    sudo journalctl -u archivas-betanet -f"
echo "    Verify:  bash verify.sh"
echo "=================================================="

