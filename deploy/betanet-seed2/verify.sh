#!/bin/bash

echo "=============================================="
echo "Archivas Betanet Seed 2 - Identity Verification"
echo "Server: 57.129.96.158"
echo "=============================================="
echo ""

PASS=0
FAIL=0

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_pass() {
    echo -e "${GREEN}✅ PASS${NC}: $1"
    ((PASS++))
}

check_fail() {
    echo -e "${RED}❌ FAIL${NC}: $1"
    ((FAIL++))
}

check_warn() {
    echo -e "${YELLOW}⚠️  WARN${NC}: $1"
}

echo "📡 Checking RPC endpoint..."
if curl -s http://localhost:8545 >/dev/null 2>&1; then
    check_pass "RPC endpoint accessible"
else
    check_fail "RPC endpoint not accessible"
    echo "Run: sudo systemctl status archivas-betanet"
    exit 1
fi

echo ""
echo "🔍 Verifying Chain Identity..."

# Check chain ID
CHAIN_ID=$(curl -s -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
    | jq -r '.result')

if [ "$CHAIN_ID" = "0x66c" ]; then
    check_pass "Chain ID = 0x66c (1644 decimal) ✓"
else
    check_fail "Chain ID = $CHAIN_ID (expected 0x66c)"
fi

# Check network version
NET_VERSION=$(curl -s -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
    | jq -r '.result')

if [ "$NET_VERSION" = "1644" ]; then
    check_pass "Network ID = 1644 ✓"
else
    check_fail "Network ID = $NET_VERSION (expected 1644)"
fi

# Check block height
BLOCK_NUM=$(curl -s -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    | jq -r '.result')

BLOCK_DEC=$((BLOCK_NUM))
if [ "$BLOCK_DEC" -ge 0 ]; then
    check_pass "Current block height = $BLOCK_DEC"
else
    check_fail "Could not fetch block height"
fi

echo ""
echo "🔗 Checking P2P connectivity..."

# Check if node service is running
if systemctl is-active --quiet archivas-betanet; then
    check_pass "Node service is running"
else
    check_fail "Node service is not running"
fi

# Check P2P port
if netstat -tuln | grep -q ":30303"; then
    check_pass "P2P port 30303 is listening"
else
    check_fail "P2P port 30303 is not listening"
fi

# Check if connected to Seed 1
if journalctl -u archivas-betanet -n 100 | grep -q "seed1.betanet.archivas.ai"; then
    check_pass "Connected to Seed 1"
else
    check_warn "No recent connection to Seed 1 in logs"
fi

echo ""
echo "🌾 Checking Farmer status..."

if systemctl is-active --quiet archivas-betanet-farmer; then
    check_pass "Farmer service is running"
    
    # Check for recent farming activity
    if journalctl -u archivas-betanet-farmer -n 50 | grep -q "plot"; then
        check_pass "Farmer is processing plots"
    else
        check_warn "No recent plot activity"
    fi
else
    check_warn "Farmer service is not running (start after creating plots)"
fi

# Check plots directory
PLOT_COUNT=$(find /mnt/plots -name "*.plot" 2>/dev/null | wc -l)
if [ "$PLOT_COUNT" -gt 0 ]; then
    check_pass "Found $PLOT_COUNT plot file(s)"
else
    check_warn "No plots found (create with: archivas-farmer plot --size 20)"
fi

echo ""
echo "🔐 Checking Firewall..."

if ufw status | grep -q "30303"; then
    check_pass "P2P port 30303 allowed in firewall"
else
    check_fail "P2P port 30303 not allowed in firewall"
fi

if ufw status | grep -q "8545"; then
    check_pass "RPC port 8545 allowed in firewall"
else
    check_warn "RPC port 8545 not allowed in firewall"
fi

echo ""
echo "=============================================="
echo "📊 Verification Summary"
echo "=============================================="
echo -e "${GREEN}Passed: $PASS${NC}"
echo -e "${RED}Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 All critical checks passed!${NC}"
    echo ""
    echo "Seed 2 is operational and syncing from Seed 1."
    exit 0
else
    echo -e "${RED}⚠️  Some checks failed. Review errors above.${NC}"
    exit 1
fi




