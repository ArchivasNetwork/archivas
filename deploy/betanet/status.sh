#!/bin/bash

# Archivas Betanet Node Status Script
# Shows real-time status of the Betanet node

echo "═══════════════════════════════════════════════════"
echo "   Archivas Betanet Node Status"
echo "   Server: 72.251.11.191"
echo "═══════════════════════════════════════════════════"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check service status
echo "📊 Service Status:"
if systemctl is-active archivas-betanet &> /dev/null; then
    echo -e "   ${GREEN}●${NC} archivas-betanet.service - RUNNING"
    UPTIME=$(systemctl show archivas-betanet -p ActiveEnterTimestamp --value)
    echo "   ⏱️  Uptime: $(systemctl show archivas-betanet -p ActiveEnterTimestampMonotonic --value | awk '{print int($1/1000000) " seconds"}')"
else
    echo -e "   ${RED}●${NC} archivas-betanet.service - STOPPED"
fi

if systemctl is-active archivas-betanet-farmer &> /dev/null; then
    echo -e "   ${GREEN}●${NC} archivas-betanet-farmer.service - RUNNING"
else
    echo -e "   ${YELLOW}○${NC} archivas-betanet-farmer.service - STOPPED"
fi

echo ""
echo "🌐 Network Information:"
echo "   Network: betanet"
echo "   Chain ID: archivas-betanet-1"
echo "   Network ID: 102"
echo "   Protocol: v2 (EVM-enabled)"

echo ""
echo "🔌 Ports:"
if netstat -tuln | grep -q ":8545"; then
    echo -e "   ${GREEN}✓${NC} RPC: 0.0.0.0:8545 (listening)"
else
    echo -e "   ${RED}✗${NC} RPC: 0.0.0.0:8545 (not listening)"
fi

if netstat -tuln | grep -q ":30303"; then
    echo -e "   ${GREEN}✓${NC} P2P: 0.0.0.0:30303 (listening)"
else
    echo -e "   ${RED}✗${NC} P2P: 0.0.0.0:30303 (not listening)"
fi

echo ""
echo "⛓️  Blockchain Status:"
# Try to get block number
if command -v curl &> /dev/null && netstat -tuln | grep -q ":8545"; then
    BLOCK_HEX=$(curl -s -X POST http://localhost:8545/eth \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        | jq -r '.result' 2>/dev/null)
    
    if [ ! -z "$BLOCK_HEX" ] && [ "$BLOCK_HEX" != "null" ]; then
        BLOCK_NUMBER=$((16#${BLOCK_HEX#0x}))
        echo "   Current Height: $BLOCK_NUMBER"
        
        # Get chain ID
        CHAIN_ID_HEX=$(curl -s -X POST http://localhost:8545/eth \
            -H "Content-Type: application/json" \
            -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
            | jq -r '.result' 2>/dev/null)
        CHAIN_ID_DEC=$((16#${CHAIN_ID_HEX#0x}))
        echo "   Chain ID (RPC): $CHAIN_ID_DEC"
        
        # Try to get gas price
        GAS_PRICE=$(curl -s -X POST http://localhost:8545/eth \
            -H "Content-Type: application/json" \
            -d '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}' \
            | jq -r '.result' 2>/dev/null)
        
        if [ ! -z "$GAS_PRICE" ] && [ "$GAS_PRICE" != "null" ]; then
            GAS_PRICE_DEC=$((16#${GAS_PRICE#0x}))
            GAS_PRICE_GWEI=$(echo "scale=2; $GAS_PRICE_DEC / 1000000000" | bc)
            echo "   Gas Price: $GAS_PRICE_GWEI gwei"
        fi
    else
        echo "   ${YELLOW}⚠️  Node still syncing or not responding${NC}"
    fi
else
    echo "   ${YELLOW}⚠️  Unable to query blockchain status${NC}"
fi

echo ""
echo "👥 Peer Connections:"
PEER_COUNT=$(journalctl -u archivas-betanet --since "5 minutes ago" | grep -c "peer" || echo "0")
echo "   Active connections: ~$PEER_COUNT (approx from logs)"
echo "   Seed nodes:"
echo "     - seed1.betanet.archivas.ai:30303"
echo "     - seed2.betanet.archivas.ai:30303"

echo ""
echo "💾 Data Directory:"
if [ -d "/var/lib/archivas/betanet" ]; then
    SIZE=$(du -sh /var/lib/archivas/betanet 2>/dev/null | cut -f1)
    echo "   Path: /var/lib/archivas/betanet"
    echo "   Size: $SIZE"
    
    FILE_COUNT=$(find /var/lib/archivas/betanet -type f 2>/dev/null | wc -l)
    echo "   Files: $FILE_COUNT"
else
    echo -e "   ${RED}✗${NC} Data directory not found"
fi

echo ""
echo "📝 Recent Logs (last 10 lines):"
journalctl -u archivas-betanet -n 10 --no-pager 2>/dev/null || echo "   No logs available"

echo ""
echo "═══════════════════════════════════════════════════"
echo "   Quick Commands"
echo "═══════════════════════════════════════════════════"
echo ""
echo "View logs (live):           journalctl -u archivas-betanet -f"
echo "Restart node:               sudo systemctl restart archivas-betanet"
echo "Stop node:                  sudo systemctl stop archivas-betanet"
echo "Start node:                 sudo systemctl start archivas-betanet"
echo "Check service:              sudo systemctl status archivas-betanet"
echo "Run verification:           sudo bash verify.sh"
echo ""

