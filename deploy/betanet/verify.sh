#!/bin/bash

# Archivas Betanet Identity Verification Script
# Ensures node is running on correct network with correct identity

echo "═══════════════════════════════════════════════════"
echo "   Archivas Betanet Identity Verification"
echo "   Server: 72.251.11.191"
echo "═══════════════════════════════════════════════════"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Expected values
EXPECTED_NETWORK="betanet"
EXPECTED_CHAIN_ID="archivas-betanet-1"
EXPECTED_NETWORK_ID=102
EXPECTED_PROTOCOL_VERSION=2
EXPECTED_RPC_PORT=8545
EXPECTED_P2P_PORT=30303

SUCCESS_COUNT=0
FAIL_COUNT=0

check() {
    local name=$1
    local expected=$2
    local actual=$3
    
    if [ "$expected" == "$actual" ]; then
        echo -e "${GREEN}✅ $name: $actual${NC}"
        ((SUCCESS_COUNT++))
        return 0
    else
        echo -e "${RED}❌ $name: Expected $expected, got $actual${NC}"
        ((FAIL_COUNT++))
        return 1
    fi
}

echo "🔍 Checking Genesis File..."
GENESIS_FILE="/etc/archivas/betanet/genesis-betanet.json"
if [ -f "$GENESIS_FILE" ]; then
    CHAIN_ID=$(jq -r '.chain_id' "$GENESIS_FILE")
    NETWORK_ID=$(jq -r '.network_id' "$GENESIS_FILE")
    PROTOCOL_VERSION=$(jq -r '.protocol_version' "$GENESIS_FILE")
    
    check "Chain ID" "$EXPECTED_CHAIN_ID" "$CHAIN_ID"
    check "Network ID" "$EXPECTED_NETWORK_ID" "$NETWORK_ID"
    check "Protocol Version" "$EXPECTED_PROTOCOL_VERSION" "$PROTOCOL_VERSION"
else
    echo -e "${RED}❌ Genesis file not found: $GENESIS_FILE${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "🔍 Checking Node Binary..."
if command -v archivas-node &> /dev/null; then
    echo -e "${GREEN}✅ Node binary found: $(which archivas-node)${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ Node binary not found${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "🔍 Checking Systemd Service..."
if systemctl is-enabled archivas-betanet &> /dev/null; then
    echo -e "${GREEN}✅ Service enabled${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${YELLOW}⚠️  Service not enabled${NC}"
fi

if systemctl is-active archivas-betanet &> /dev/null; then
    echo -e "${GREEN}✅ Service running${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ Service not running${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "🔍 Checking Network Ports..."
if netstat -tuln | grep -q ":$EXPECTED_RPC_PORT"; then
    echo -e "${GREEN}✅ RPC port $EXPECTED_RPC_PORT is listening${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ RPC port $EXPECTED_RPC_PORT not listening${NC}"
    ((FAIL_COUNT++))
fi

if netstat -tuln | grep -q ":$EXPECTED_P2P_PORT"; then
    echo -e "${GREEN}✅ P2P port $EXPECTED_P2P_PORT is listening${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ P2P port $EXPECTED_P2P_PORT not listening${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "🔍 Checking Firewall..."
if ufw status | grep -q "$EXPECTED_P2P_PORT"; then
    echo -e "${GREEN}✅ P2P port $EXPECTED_P2P_PORT allowed in firewall${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ P2P port $EXPECTED_P2P_PORT not allowed in firewall${NC}"
    ((FAIL_COUNT++))
fi

if ufw status | grep -q "$EXPECTED_RPC_PORT"; then
    echo -e "${GREEN}✅ RPC port $EXPECTED_RPC_PORT allowed in firewall${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ RPC port $EXPECTED_RPC_PORT not allowed in firewall${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "🔍 Checking RPC Endpoints..."
if curl -s -f http://localhost:8545/eth > /dev/null 2>&1; then
    # Try eth_chainId
    CHAIN_ID_RESPONSE=$(curl -s -X POST http://localhost:8545/eth \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' | jq -r '.result')
    
    if [ "$CHAIN_ID_RESPONSE" == "0x66c" ]; then  # 102 in hex
        echo -e "${GREEN}✅ ETH RPC responding with correct chain ID (0x66c = 1644)${NC}"
        ((SUCCESS_COUNT++))
    else
        echo -e "${RED}❌ ETH RPC responding with wrong chain ID: $CHAIN_ID_RESPONSE${NC}"
        ((FAIL_COUNT++))
    fi
else
    echo -e "${RED}❌ RPC endpoint not accessible${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "🔍 Checking Data Directory..."
if [ -d "/var/lib/archivas/betanet" ]; then
    SIZE=$(du -sh /var/lib/archivas/betanet | cut -f1)
    echo -e "${GREEN}✅ Data directory exists: $SIZE${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ Data directory not found${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "🔍 Checking Logs..."
if journalctl -u archivas-betanet --since "1 minute ago" | grep -q "Archivas"; then
    echo -e "${GREEN}✅ Logs are being generated${NC}"
    ((SUCCESS_COUNT++))
    
    # Check for network identity in logs
    if journalctl -u archivas-betanet --since "5 minutes ago" | grep -q "betanet"; then
        echo -e "${GREEN}✅ Logs confirm betanet network${NC}"
        ((SUCCESS_COUNT++))
    fi
else
    echo -e "${YELLOW}⚠️  No recent logs found${NC}"
fi

echo ""
echo "═══════════════════════════════════════════════════"
echo "   Verification Summary"
echo "═══════════════════════════════════════════════════"
echo -e "${GREEN}✅ Passed: $SUCCESS_COUNT${NC}"
echo -e "${RED}❌ Failed: $FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}🎉 All checks passed! Node is correctly configured for Betanet.${NC}"
    exit 0
else
    echo -e "${RED}⚠️  Some checks failed. Please review the output above.${NC}"
    exit 1
fi

