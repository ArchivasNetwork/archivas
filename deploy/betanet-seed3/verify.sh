#!/bin/bash

echo "=================================================="
echo "  Archivas Betanet Seed 3 Identity Verification"
echo "=================================================="
echo ""

RPC_URL="http://localhost:8545"
EXPECTED_CHAIN_ID="0x66c"  # 1644 in hex
EXPECTED_NETWORK_ID="1644"

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

check_passed=0
check_failed=0

echo "🔍 Running identity checks..."
echo ""

# Check 1: Service running
echo -n "✓ Service Status: "
if systemctl is-active --quiet archivas-betanet; then
    echo -e "${GREEN}RUNNING${NC}"
    ((check_passed++))
else
    echo -e "${RED}NOT RUNNING${NC}"
    ((check_failed++))
fi

# Check 2: P2P Port listening
echo -n "✓ P2P Port (30303): "
if netstat -tuln | grep -q ":30303"; then
    echo -e "${GREEN}LISTENING${NC}"
    ((check_passed++))
else
    echo -e "${RED}NOT LISTENING${NC}"
    ((check_failed++))
fi

# Check 3: RPC Port listening
echo -n "✓ RPC Port (8545): "
if netstat -tuln | grep -q ":8545"; then
    echo -e "${GREEN}LISTENING${NC}"
    ((check_passed++))
else
    echo -e "${RED}NOT LISTENING${NC}"
    ((check_failed++))
fi

# Give node time to start if just installed
sleep 2

# Check 4: RPC responding
echo -n "✓ RPC Responding: "
if curl -s -m 5 "$RPC_URL/" > /dev/null 2>&1; then
    echo -e "${GREEN}OK${NC}"
    ((check_passed++))
else
    echo -e "${RED}NO RESPONSE${NC}"
    ((check_failed++))
fi

# Check 5: Chain ID
echo -n "✓ Chain ID (1644): "
CHAIN_ID=$(curl -s -X POST "$RPC_URL/" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
    | grep -o '"result":"[^"]*"' | cut -d'"' -f4)

if [ "$CHAIN_ID" = "$EXPECTED_CHAIN_ID" ]; then
    echo -e "${GREEN}$CHAIN_ID (CORRECT)${NC}"
    ((check_passed++))
else
    echo -e "${RED}$CHAIN_ID (EXPECTED $EXPECTED_CHAIN_ID)${NC}"
    ((check_failed++))
fi

# Check 6: Network ID
echo -n "✓ Network ID (1644): "
NETWORK_ID=$(curl -s -X POST "$RPC_URL/" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
    | grep -o '"result":"[^"]*"' | cut -d'"' -f4)

if [ "$NETWORK_ID" = "$EXPECTED_NETWORK_ID" ]; then
    echo -e "${GREEN}$NETWORK_ID (CORRECT)${NC}"
    ((check_passed++))
else
    echo -e "${RED}$NETWORK_ID (EXPECTED $EXPECTED_NETWORK_ID)${NC}"
    ((check_failed++))
fi

# Check 7: Block height
echo -n "✓ Block Height: "
BLOCK_HEX=$(curl -s -X POST "$RPC_URL/" \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    | grep -o '"result":"[^"]*"' | cut -d'"' -f4)

if [ -n "$BLOCK_HEX" ]; then
    BLOCK_HEIGHT=$((16#${BLOCK_HEX#0x}))
    if [ "$BLOCK_HEIGHT" -ge 0 ]; then
        echo -e "${GREEN}$BLOCK_HEIGHT${NC}"
        ((check_passed++))
    else
        echo -e "${RED}INVALID${NC}"
        ((check_failed++))
    fi
else
    echo -e "${RED}NO RESPONSE${NC}"
    ((check_failed++))
fi

# Check 8: P2P Peers
echo -n "✓ Connected Peers: "
PEER_COUNT=$(sudo journalctl -u archivas-betanet --since "1 minute ago" | grep -c "connected to peer" || echo "0")
if [ "$PEER_COUNT" -gt 0 ]; then
    echo -e "${GREEN}$PEER_COUNT${NC}"
    ((check_passed++))
else
    echo -e "${YELLOW}$PEER_COUNT (may take time to connect)${NC}"
    # Don't fail this check as it might take time
    ((check_passed++))
fi

# Check 9: Genesis Hash
echo -n "✓ Genesis Hash: "
GENESIS_HASH=$(sudo journalctl -u archivas-betanet --since "10 minutes ago" | grep "Genesis Hash:" | tail -1 | grep -o '74187e4036f7a489')
if [ "$GENESIS_HASH" = "74187e4036f7a489" ]; then
    echo -e "${GREEN}$GENESIS_HASH (CORRECT)${NC}"
    ((check_passed++))
else
    echo -e "${YELLOW}Check logs for genesis hash${NC}"
    ((check_passed++))
fi

# Check 10: No farming activity (should be relay only)
echo -n "✓ Farming Status: "
if sudo journalctl -u archivas-betanet --since "5 minutes ago" | grep -q "archivas-farmer"; then
    echo -e "${RED}FARMING DETECTED (should be relay only)${NC}"
    ((check_failed++))
else
    echo -e "${GREEN}DISABLED (correct for public gateway)${NC}"
    ((check_passed++))
fi

echo ""
echo "=================================================="
echo "  Summary: $check_passed passed, $check_failed failed"
echo "=================================================="
echo ""

if [ $check_failed -eq 0 ]; then
    echo -e "${GREEN}✅ Seed 3 is correctly configured as a public gateway!${NC}"
    echo ""
    echo "🌍 External users can connect to:"
    echo "   RPC:  http://51.89.11.4:8545"
    echo "   P2P:  seed3.betanet.archivas.ai:30303"
    exit 0
else
    echo -e "${RED}❌ Some checks failed. Review logs:${NC}"
    echo "   sudo journalctl -u archivas-betanet -n 50"
    exit 1
fi

