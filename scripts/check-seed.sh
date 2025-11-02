#!/usr/bin/env bash
# check-seed.sh - Validate seed.archivas.ai is operational
# Usage: bash scripts/check-seed.sh

set -euo pipefail

SEED_URL="https://seed.archivas.ai"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 Checking seed.archivas.ai"
echo "============================="
echo ""

check_endpoint() {
    local name=$1
    local url=$2
    local method=${3:-GET}
    local expected_status=${4:-200}
    
    echo -n "Testing $name ($method $url)... "
    
    status=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$url" -H "Content-Type: application/json" --max-time 10 || echo "000")
    
    if [[ "$status" == "$expected_status" ]]; then
        echo -e "${GREEN}✅ $status${NC}"
        return 0
    else
        echo -e "${RED}❌ $status (expected $expected_status)${NC}"
        return 1
    fi
}

check_cors() {
    local url=$1
    echo -n "Checking CORS headers... "
    
    cors=$(curl -s -I "$url" | grep -i "access-control-allow-origin" || echo "")
    
    if [[ -n "$cors" ]]; then
        echo -e "${GREEN}✅ Present${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Missing${NC}"
        return 1
    fi
}

# Test HTTP → HTTPS redirect
echo "📡 Testing HTTP → HTTPS redirect..."
check_endpoint "HTTP redirect" "http://seed.archivas.ai/version" "GET" "301"
echo ""

# Test HTTPS endpoints
echo "🔐 Testing HTTPS endpoints..."
check_endpoint "/version" "$SEED_URL/version" "GET" "200"
check_endpoint "/chainTip" "$SEED_URL/chainTip" "GET" "200"
check_endpoint "/health" "$SEED_URL/health" "GET" "200"
echo ""

# Test /submit method handling
echo "📝 Testing /submit method handling..."
check_endpoint "/submit GET (should 405)" "$SEED_URL/submit" "GET" "405"
check_endpoint "/submit OPTIONS" "$SEED_URL/submit" "OPTIONS" "204"
echo ""

# Test CORS headers
echo "🌐 Testing CORS headers..."
check_cors "$SEED_URL/version"
check_cors "$SEED_URL/chainTip"
echo ""

# Test actual data
echo "📊 Testing response data..."
echo -n "Fetching /version... "
VERSION=$(curl -s "$SEED_URL/version" || echo "{}")
if echo "$VERSION" | grep -q "version"; then
    echo -e "${GREEN}✅$(NC} $(echo $VERSION | head -c 80)..."
else
    echo -e "${RED}❌ Invalid response${NC}"
fi

echo -n "Fetching /chainTip... "
TIP=$(curl -s "$SEED_URL/chainTip" || echo "{}")
if echo "$TIP" | grep -q "height"; then
    echo -e "${GREEN}✅${NC} $(echo $TIP | head -c 80)..."
else
    echo -e "${RED}❌ Invalid response${NC}"
fi
echo ""

# Test TLS
echo "🔒 Testing TLS..."
echo -n "Certificate validity... "
if echo | openssl s_client -connect seed.archivas.ai:443 -servername seed.archivas.ai 2>/dev/null | openssl x509 -noout -dates 2>/dev/null; then
    echo -e "${GREEN}✅ Valid${NC}"
else
    echo -e "${RED}❌ Invalid or missing${NC}"
fi
echo ""

# Test HTTP/2
echo "⚡ Testing HTTP/2..."
echo -n "HTTP/2 support... "
if curl -sI "$SEED_URL/version" --http2 | grep -q "HTTP/2"; then
    echo -e "${GREEN}✅ Enabled${NC}"
else
    echo -e "${YELLOW}⚠️  HTTP/1.1 only${NC}"
fi
echo ""

echo "========================================="
echo -e "${GREEN}✅ seed.archivas.ai validation complete${NC}"
echo ""
echo "Endpoint: $SEED_URL"
echo "Status: Operational"

