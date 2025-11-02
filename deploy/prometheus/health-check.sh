#!/bin/bash
# v1.1.1: Health check script for Prometheus targets
# Run from Prometheus host to verify all targets are reachable

echo "🔍 Checking Archivas metrics endpoints..."
echo ""

# Color output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

check_endpoint() {
    name=$1
    url=$2
    
    if curl -sf --max-time 5 "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✅${NC} $name: $url"
        return 0
    else
        echo -e "${RED}❌${NC} $name: $url (connection refused or timeout)"
        return 1
    fi
}

# Check nodes
echo "📊 Nodes (port 8080):"
check_endpoint "Node 1" "http://57.129.148.132:8080/metrics"
check_endpoint "Node 2" "http://72.251.11.191:8080/metrics"
echo ""

# Check timelords
echo "⏰ Timelords (port 9101):"
check_endpoint "Timelord 1" "http://57.129.148.132:9101/metrics"
check_endpoint "Timelord 2" "http://72.251.11.191:9101/metrics"
echo ""

# Check farmers
echo "🌾 Farmers (port 9102):"
check_endpoint "Farmer A" "http://57.129.148.132:9102/metrics"
check_endpoint "Farmer C" "http://57.129.148.134:9102/metrics"
echo ""

# Check health endpoints
echo "🏥 Health endpoints:"
check_endpoint "Node A health" "http://57.129.148.132:8080/healthz"
check_endpoint "Node B health" "http://72.251.11.191:8080/healthz"
check_endpoint "Timelord A health" "http://57.129.148.132:9101/healthz"
check_endpoint "Timelord B health" "http://72.251.11.191:9101/healthz"
check_endpoint "Farmer A health" "http://57.129.148.132:9102/healthz"
check_endpoint "Farmer C health" "http://57.129.148.134:9102/healthz"
echo ""

echo "✅ Health check complete"

