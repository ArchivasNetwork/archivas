#!/bin/bash

echo "=============================================="
echo "Archivas Betanet Seed 2 - Status Dashboard"
echo "Server: 57.129.96.158"
echo "=============================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}📡 Node Status${NC}"
echo "─────────────────────────────────────────────"

if systemctl is-active --quiet archivas-betanet; then
    echo -e "Service: ${GREEN}RUNNING${NC}"
else
    echo -e "Service: ${RED}STOPPED${NC}"
fi

# Get chain info
CHAIN_ID=$(curl -s -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
    | jq -r '.result' 2>/dev/null)

BLOCK_NUM=$(curl -s -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    | jq -r '.result' 2>/dev/null)

BLOCK_DEC=$((BLOCK_NUM))

echo "Chain ID: $CHAIN_ID ($(($CHAIN_ID)) decimal)"
echo "Current Block: $BLOCK_DEC"

# Get sync status from Seed 1
SEED1_BLOCK=$(curl -s -X POST http://seed1.betanet.archivas.ai:8545 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    | jq -r '.result' 2>/dev/null)

SEED1_DEC=$((SEED1_BLOCK))

if [ "$SEED1_DEC" -gt 0 ]; then
    echo "Seed 1 Block: $SEED1_DEC"
    DIFF=$((SEED1_DEC - BLOCK_DEC))
    if [ $DIFF -eq 0 ]; then
        echo -e "Sync Status: ${GREEN}FULLY SYNCED${NC}"
    elif [ $DIFF -lt 10 ]; then
        echo -e "Sync Status: ${GREEN}SYNCING${NC} (behind by $DIFF blocks)"
    else
        echo -e "Sync Status: ${YELLOW}SYNCING${NC} (behind by $DIFF blocks)"
    fi
else
    echo "Seed 1 Block: Unable to connect"
fi

echo ""
echo -e "${BLUE}🌾 Farmer Status${NC}"
echo "─────────────────────────────────────────────"

if systemctl is-active --quiet archivas-betanet-farmer; then
    echo -e "Service: ${GREEN}RUNNING${NC}"
else
    echo -e "Service: ${RED}STOPPED${NC}"
fi

PLOT_COUNT=$(find /mnt/plots -name "*.plot" 2>/dev/null | wc -l)
echo "Plots: $PLOT_COUNT"

if [ $PLOT_COUNT -gt 0 ]; then
    PLOT_SIZE=$(du -sh /mnt/plots 2>/dev/null | cut -f1)
    echo "Plot Storage: $PLOT_SIZE"
fi

echo ""
echo -e "${BLUE}🔗 Network Status${NC}"
echo "─────────────────────────────────────────────"

# Check P2P listening
if netstat -tuln | grep -q ":30303"; then
    echo -e "P2P Port 30303: ${GREEN}LISTENING${NC}"
else
    echo -e "P2P Port 30303: ${RED}NOT LISTENING${NC}"
fi

# Check RPC listening
if netstat -tuln | grep -q ":8545"; then
    echo -e "RPC Port 8545: ${GREEN}LISTENING${NC}"
else
    echo -e "RPC Port 8545: ${RED}NOT LISTENING${NC}"
fi

# Check external connectivity
if curl -s --max-time 3 http://seed1.betanet.archivas.ai:30303 >/dev/null 2>&1; then
    echo -e "Seed 1 Reachable: ${GREEN}YES${NC}"
else
    echo -e "Seed 1 Reachable: ${YELLOW}CHECK${NC}"
fi

echo ""
echo -e "${BLUE}💾 System Resources${NC}"
echo "─────────────────────────────────────────────"

# Disk usage
DISK_DATA=$(df -h /var/lib/archivas/betanet | tail -1 | awk '{print $3 "/" $2 " (" $5 ")"}')
echo "Data Dir: $DISK_DATA"

DISK_PLOTS=$(df -h /mnt/plots | tail -1 | awk '{print $3 "/" $2 " (" $5 ")"}')
echo "Plots Dir: $DISK_PLOTS"

# Memory usage
MEM=$(free -h | awk '/^Mem:/ {print $3 "/" $2}')
echo "Memory: $MEM"

echo ""
echo -e "${BLUE}📋 Recent Logs (last 5 lines)${NC}"
echo "─────────────────────────────────────────────"
journalctl -u archivas-betanet -n 5 --no-pager | tail -5

echo ""
echo "─────────────────────────────────────────────"
echo "Run 'sudo journalctl -u archivas-betanet -f' for live logs"
echo "Run 'sudo journalctl -u archivas-betanet-farmer -f' for farmer logs"
echo "Run './deploy/betanet-seed2/verify.sh' for full verification"
echo ""




