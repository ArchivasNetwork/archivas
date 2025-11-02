#!/usr/bin/env bash
# join_devnet.sh - One-command setup to join Archivas Devnet V4
# Usage: curl https://raw.githubusercontent.com/ArchivasNetwork/archivas/main/scripts/join_devnet.sh | bash

set -euo pipefail

REPO=https://github.com/ArchivasNetwork/archivas.git
BOOT=57.129.148.132:9090
NET=archivas-devnet-v4

echo "🌾 Archivas Devnet V4 - Quick Join"
echo ""

# Clone repo
if [ ! -d "archivas" ]; then
    echo "📥 Cloning repository..."
    git clone "$REPO" archivas
else
    echo "📁 Using existing archivas directory"
    cd archivas
    git pull origin main
    cd ..
fi

cd archivas

# Build with stamped metadata
echo "🔨 Building binaries..."
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +%FT%TZ)
VERSION="v1.1.1"

go build -ldflags "-X github.com/ArchivasNetwork/archivas/internal/buildinfo.Version=$VERSION \
 -X github.com/ArchivasNetwork/archivas/internal/buildinfo.Commit=$COMMIT \
 -X github.com/ArchivasNetwork/archivas/internal/buildinfo.BuiltAt=$DATE" \
 -o archivas-node ./cmd/archivas-node

go build -ldflags "-X github.com/ArchivasNetwork/archivas/internal/buildinfo.Version=$VERSION \
 -X github.com/ArchivasNetwork/archivas/internal/buildinfo.Commit=$COMMIT \
 -X github.com/ArchivasNetwork/archivas/internal/buildinfo.BuiltAt=$DATE" \
 -o archivas-farmer ./cmd/archivas-farmer

go build -o archivas-wallet ./cmd/archivas-wallet

echo "✅ Binaries built"

# Create directories
mkdir -p data logs plots

# Generate wallet
echo ""
echo "🔐 Generating wallet..."
./archivas-wallet new > wallet.txt
cat wallet.txt

PUB=$(grep "Public Key:" wallet.txt | awk '{print $3}')
PRV=$(grep "Private Key:" wallet.txt | awk '{print $3}')

# Save keys
echo "$PRV" > .farmer-privkey
chmod 600 .farmer-privkey
echo "🔑 Private key saved to .farmer-privkey"

# Start node
echo ""
echo "🚀 Starting node..."
nohup ./archivas-node \
  --rpc 0.0.0.0:8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id $NET \
  --bootnodes "$BOOT" \
  > logs/node.log 2>&1 &

NODE_PID=$!
echo "✅ Node started (PID: $NODE_PID)"
echo "   Logs: logs/node.log"

# Wait for node to initialize
sleep 5

# Check node health
if curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
    echo "✅ Node is healthy"
else
    echo "⚠️  Node might still be starting (check logs/node.log)"
fi

# Create a small plot (k=20 for quick start, ~32MB)
echo ""
echo "📊 Creating plot (k=20, ~32MB)..."
echo "   (This takes ~2-3 seconds)"
./archivas-farmer plot --path ./plots --size 20 --farmer-pubkey $PUB

# Start farmer
echo ""
echo "🌾 Starting farmer..."
nohup ./archivas-farmer farm \
  --plots ./plots \
  --node http://localhost:8080 \
  --farmer-privkey $PRV \
  --metrics-addr 0.0.0.0:9102 \
  > logs/farmer.log 2>&1 &

FARMER_PID=$!
echo "✅ Farmer started (PID: $FARMER_PID)"
echo "   Logs: logs/farmer.log"

# Summary
echo ""
echo "═══════════════════════════════════════"
echo "✅ Archivas Devnet V4 - READY"
echo "═══════════════════════════════════════"
echo ""
echo "Node:    http://localhost:8080"
echo "Network: $NET"
echo "Wallet:  wallet.txt"
echo ""
echo "📊 Check sync status:"
echo "   curl http://localhost:8080/chainTip | jq"
echo ""
echo "📈 Monitor logs:"
echo "   tail -f logs/node.log"
echo "   tail -f logs/farmer.log"
echo ""
echo "🛑 Stop services:"
echo "   pkill archivas-node"
echo "   pkill archivas-farmer"
echo ""
echo "Happy farming! 🚜"

