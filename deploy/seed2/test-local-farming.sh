#!/bin/bash
set -euo pipefail

# Test Seed2 Local Node for Farming
# This script verifies that Seed2's local node has all the functionality needed for farming
# Run this BEFORE switching Nginx to use the local node

echo "=== Testing Seed2 Local Node for Farming ==="
echo ""

LOCAL_NODE="http://127.0.0.1:8082"
SEED1_NODE="https://seed.archivas.ai:8081"

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

pass() {
    echo -e "${GREEN}✓ PASS${NC}: $1"
}

fail() {
    echo -e "${RED}✗ FAIL${NC}: $1"
    return 1
}

warn() {
    echo -e "${YELLOW}⚠ WARN${NC}: $1"
}

test_passed=0
test_failed=0

# Test 1: Check if RPC is listening
echo "Test 1: RPC Server Listening"
if ss -tulnp | grep -q "127.0.0.1:8082"; then
    pass "RPC server is listening on 127.0.0.1:8082"
    ((test_passed++))
else
    fail "RPC server not listening on 127.0.0.1:8082"
    ((test_failed++))
fi
echo ""

# Test 2: /challenge endpoint (CRITICAL for farming)
echo "Test 2: /challenge Endpoint"
CHALLENGE=$(curl -s "$LOCAL_NODE/challenge" 2>/dev/null || echo "")
if [ -n "$CHALLENGE" ] && echo "$CHALLENGE" | jq -e '.challenge' >/dev/null 2>&1; then
    HEIGHT=$(echo "$CHALLENGE" | jq -r '.height')
    DIFFICULTY=$(echo "$CHALLENGE" | jq -r '.difficulty')
    pass "/challenge returns valid data (height: $HEIGHT, difficulty: $DIFFICULTY)"
    ((test_passed++))
else
    fail "/challenge endpoint not responding or invalid response"
    ((test_failed++))
fi
echo ""

# Test 3: /chainTip endpoint
echo "Test 3: /chainTip Endpoint"
CHAIN_TIP=$(curl -s "$LOCAL_NODE/chainTip" 2>/dev/null || echo "")
if [ -n "$CHAIN_TIP" ] && echo "$CHAIN_TIP" | jq -e '.height' >/dev/null 2>&1; then
    LOCAL_HEIGHT=$(echo "$CHAIN_TIP" | jq -r '.height')
    pass "/chainTip returns valid data (height: $LOCAL_HEIGHT)"
    ((test_passed++))
else
    fail "/chainTip endpoint not responding or invalid response"
    ((test_failed++))
fi
echo ""

# Test 4: /sync/status endpoint
echo "Test 4: /sync/status Endpoint"
SYNC_STATUS=$(curl -s "$LOCAL_NODE/sync/status" 2>/dev/null || echo "")
if [ -n "$SYNC_STATUS" ] && echo "$SYNC_STATUS" | jq -e '.synced' >/dev/null 2>&1; then
    SYNCED=$(echo "$SYNC_STATUS" | jq -r '.synced')
    if [ "$SYNCED" = "true" ]; then
        pass "/sync/status reports node is synced"
        ((test_passed++))
    else
        warn "/sync/status reports node is NOT synced"
        ((test_failed++))
    fi
else
    fail "/sync/status endpoint not responding or invalid response"
    ((test_failed++))
fi
echo ""

# Test 5: Compare height with Seed1
echo "Test 5: Height Comparison with Seed1"
SEED1_TIP=$(curl -sk "$SEED1_NODE/chainTip" 2>/dev/null || echo "")
if [ -n "$SEED1_TIP" ] && echo "$SEED1_TIP" | jq -e '.height' >/dev/null 2>&1; then
    SEED1_HEIGHT=$(echo "$SEED1_TIP" | jq -r '.height')
    HEIGHT_DIFF=$((SEED1_HEIGHT - LOCAL_HEIGHT))
    
    if [ $HEIGHT_DIFF -le 10 ] && [ $HEIGHT_DIFF -ge -10 ]; then
        pass "Height difference is acceptable (Seed2: $LOCAL_HEIGHT, Seed1: $SEED1_HEIGHT, diff: $HEIGHT_DIFF)"
        ((test_passed++))
    else
        warn "Height difference is large (Seed2: $LOCAL_HEIGHT, Seed1: $SEED1_HEIGHT, diff: $HEIGHT_DIFF)"
        echo "   Node may still be syncing or have an issue"
    fi
else
    warn "Could not fetch Seed1 height for comparison (may be rate limited)"
fi
echo ""

# Test 6: /submit endpoint exists (without actually submitting)
echo "Test 6: /submit Endpoint Exists"
# Test with invalid data to see if endpoint processes requests
SUBMIT_TEST=$(curl -s -X POST "$LOCAL_NODE/submit" \
    -H "Content-Type: application/json" \
    -d '{"invalid":"test"}' \
    -w "\nHTTP_CODE:%{http_code}" 2>/dev/null || echo "")

if echo "$SUBMIT_TEST" | grep -q "HTTP_CODE"; then
    HTTP_CODE=$(echo "$SUBMIT_TEST" | grep "HTTP_CODE" | cut -d: -f2)
    # Any response code (400, 415, 500) means endpoint exists
    if [ "$HTTP_CODE" = "000" ]; then
        fail "/submit endpoint not responding"
        ((test_failed++))
    else
        pass "/submit endpoint is accessible (HTTP $HTTP_CODE)"
        ((test_passed++))
    fi
else
    warn "/submit endpoint test inconclusive"
fi
echo ""

# Test 7: Response time (should be fast for local queries)
echo "Test 7: Response Time"
START_TIME=$(date +%s%N)
curl -s "$LOCAL_NODE/chainTip" >/dev/null 2>&1
END_TIME=$(date +%s%N)
RESPONSE_MS=$(( (END_TIME - START_TIME) / 1000000 ))

if [ $RESPONSE_MS -lt 100 ]; then
    pass "Response time is excellent (${RESPONSE_MS}ms)"
    ((test_passed++))
elif [ $RESPONSE_MS -lt 500 ]; then
    pass "Response time is good (${RESPONSE_MS}ms)"
    ((test_passed++))
else
    warn "Response time is slow (${RESPONSE_MS}ms)"
fi
echo ""

# Test 8: Check for recent errors in logs
echo "Test 8: Recent Node Errors"
if command -v journalctl >/dev/null 2>&1; then
    ERROR_COUNT=$(sudo journalctl -u archivas-node-seed2 --since "10 minutes ago" 2>/dev/null | grep -ic "error\|panic\|fatal" || echo "0")
    if [ "$ERROR_COUNT" -lt 5 ]; then
        pass "Few or no errors in recent logs ($ERROR_COUNT in last 10 min)"
        ((test_passed++))
    else
        warn "Multiple errors detected in logs ($ERROR_COUNT in last 10 min)"
        echo "   Check logs: sudo journalctl -u archivas-node-seed2 -n 50"
    fi
else
    warn "Cannot check logs (journalctl not available)"
fi
echo ""

# Test 9: Memory usage
echo "Test 9: Node Memory Usage"
if command -v ps >/dev/null 2>&1; then
    MEMORY_MB=$(ps aux | grep archivas-node | grep -v grep | awk '{print $6/1024}' | head -1 || echo "0")
    if (( $(echo "$MEMORY_MB > 0" | bc -l) )); then
        if (( $(echo "$MEMORY_MB < 2000" | bc -l) )); then
            pass "Memory usage is healthy (${MEMORY_MB}MB)"
            ((test_passed++))
        else
            warn "Memory usage is high (${MEMORY_MB}MB)"
        fi
    fi
else
    warn "Cannot check memory usage"
fi
echo ""

# Test 10: P2P connectivity
echo "Test 10: P2P Port Status"
if ss -tulnp | grep -q ":30303"; then
    pass "P2P port 30303 is listening"
    ((test_passed++))
else
    warn "P2P port 30303 is NOT listening"
    echo "   Farmers may not be able to peer with this node"
fi
echo ""

# Summary
echo "========================================"
echo "           TEST SUMMARY"
echo "========================================"
echo -e "Tests Passed: ${GREEN}$test_passed${NC}"
echo -e "Tests Failed: ${RED}$test_failed${NC}"
echo ""

if [ $test_failed -eq 0 ]; then
    echo -e "${GREEN}✅ ALL TESTS PASSED!${NC}"
    echo ""
    echo "Seed2's local node is ready for farming!"
    echo ""
    echo "Next steps:"
    echo "  1. (Optional) Test with a real farmer against 127.0.0.1:8082"
    echo "  2. Run: sudo bash deploy/seed2/switch-to-local-node.sh"
    echo "  3. Monitor: sudo tail -f /var/log/nginx/access.log"
    echo ""
    exit 0
else
    echo -e "${RED}❌ SOME TESTS FAILED${NC}"
    echo ""
    echo "Do NOT switch Nginx to use local node yet!"
    echo "Investigate and fix the issues above first."
    echo ""
    exit 1
fi

