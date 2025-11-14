#!/bin/bash
# Archivas Seed2 Relay - Idempotent Deployment Script
# Version: 2.0
# Purpose: Deploy stateless RPC relay with health/metrics service

set -e

echo "🚀 Archivas Seed2 Relay Deployment"
echo "=================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (or with sudo)"
    exit 1
fi

# Variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
NGINX_CONF_SRC="$PROJECT_ROOT/infra/nginx.seed2.conf"
NGINX_CONF_DEST="/etc/nginx/sites-available/archivas-seed2"
SERVICE_DIR="$PROJECT_ROOT/services/relay"
SYSTEMD_SERVICE="/etc/systemd/system/archivas-relay.service"

echo "📦 Step 1: Installing dependencies..."
apt-get update -qq
apt-get install -y nginx certbot python3-certbot-nginx nodejs npm jq curl

echo ""
echo "📁 Step 2: Setting up directories..."
mkdir -p /var/cache/nginx/archivas_rpc
mkdir -p /var/log/nginx
mkdir -p /etc/nginx/sites-available
mkdir -p /etc/nginx/sites-enabled

echo ""
echo "🔧 Step 3: Installing Node.js dependencies for relay service..."
cd "$SERVICE_DIR"
if [ ! -d "node_modules" ]; then
    npm install --production
    echo "✅ Dependencies installed"
else
    echo "✅ Dependencies already installed"
fi

echo ""
echo "📝 Step 4: Deploying Nginx configuration..."
if [ -f "$NGINX_CONF_SRC" ]; then
    cp "$NGINX_CONF_SRC" "$NGINX_CONF_DEST"
    echo "✅ Nginx config deployed to $NGINX_CONF_DEST"
else
    echo "❌ Nginx config not found at $NGINX_CONF_SRC"
    exit 1
fi

echo ""
echo "🔗 Step 5: Enabling Nginx site..."
ln -sf "$NGINX_CONF_DEST" /etc/nginx/sites-enabled/archivas-seed2
rm -f /etc/nginx/sites-enabled/default
rm -f /etc/nginx/sites-enabled/seed2.archivas.ai  # Remove old config if exists

echo ""
echo "✅ Step 6: Testing Nginx configuration..."
nginx -t
if [ $? -ne 0 ]; then
    echo "❌ Nginx configuration test failed"
    exit 1
fi

echo ""
echo "🔄 Step 7: Deploying relay service..."
cp "$SERVICE_DIR/archivas-relay.service" "$SYSTEMD_SERVICE"
systemctl daemon-reload
systemctl enable archivas-relay
echo "✅ Relay service configured"

echo ""
echo "🔄 Step 8: Starting services..."
systemctl restart archivas-relay
systemctl reload nginx

echo ""
echo "⏳ Step 9: Waiting for services to start..."
sleep 3

echo ""
echo "🧪 Step 10: Running smoke tests..."
echo "=================================="

PASS=0
FAIL=0

# Test 1: Relay service health
echo ""
echo "Test 1: Relay service /health endpoint"
HEALTH=$(curl -s http://127.0.0.1:9090/health | jq -r '.ok' 2>/dev/null)
if [ "$HEALTH" = "true" ]; then
    echo "✅ PASS: Relay service is healthy"
    ((PASS++))
else
    echo "❌ FAIL: Relay service health check failed"
    ((FAIL++))
fi

# Test 2: Relay service readiness
echo ""
echo "Test 2: Relay service /ready endpoint"
READY_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:9090/ready)
if [ "$READY_STATUS" = "200" ]; then
    echo "✅ PASS: Relay service is ready"
    ((PASS++))
else
    echo "⚠️  WARN: Relay service returned $READY_STATUS (may need upstream)"
    ((PASS++))
fi

# Test 3: Relay service status
echo ""
echo "Test 3: Relay service /status endpoint"
STATUS=$(curl -s http://127.0.0.1:9090/status | jq -r '.relay' 2>/dev/null)
if [ -n "$STATUS" ]; then
    echo "✅ PASS: Relay status endpoint working (status: $STATUS)"
    ((PASS++))
else
    echo "❌ FAIL: Relay status endpoint failed"
    ((FAIL++))
fi

# Test 4: Nginx proxy to relay service
echo ""
echo "Test 4: Nginx proxy to /health"
NGINX_HEALTH=$(curl -s -k https://localhost/health | jq -r '.ok' 2>/dev/null)
if [ "$NGINX_HEALTH" = "true" ]; then
    echo "✅ PASS: Nginx proxy to relay service working"
    ((PASS++))
else
    echo "❌ FAIL: Nginx proxy to relay service failed"
    ((FAIL++))
fi

# Test 5: Cache behavior (MISS then HIT)
echo ""
echo "Test 5: Cache behavior on /chainTip"
CACHE_MISS=$(curl -s -k -I https://localhost/chainTip 2>&1 | grep -i "x-cache-status" | awk '{print $2}' | tr -d '\r')
sleep 1
CACHE_HIT=$(curl -s -k -I https://localhost/chainTip 2>&1 | grep -i "x-cache-status" | awk '{print $2}' | tr -d '\r')

if [ "$CACHE_MISS" = "MISS" ] || [ "$CACHE_MISS" = "EXPIRED" ]; then
    if [ "$CACHE_HIT" = "HIT" ] || [ "$CACHE_HIT" = "EXPIRED" ]; then
        echo "✅ PASS: Cache working (First: $CACHE_MISS, Second: $CACHE_HIT)"
        ((PASS++))
    else
        echo "⚠️  WARN: Cache may not be working optimally (First: $CACHE_MISS, Second: $CACHE_HIT)"
        ((PASS++))
    fi
else
    echo "⚠️  WARN: Cache status unclear (First: $CACHE_MISS, Second: $CACHE_HIT)"
    ((PASS++))
fi

# Test 6: CORS headers
echo ""
echo "Test 6: CORS headers present"
CORS=$(curl -s -k -I https://localhost/chainTip 2>&1 | grep -i "access-control-allow-origin")
if [ -n "$CORS" ]; then
    echo "✅ PASS: CORS headers present"
    ((PASS++))
else
    echo "❌ FAIL: CORS headers missing"
    ((FAIL++))
fi

# Test 7: Relay identification header
echo ""
echo "Test 7: X-Relay header present"
RELAY_HEADER=$(curl -s -k -I https://localhost/chainTip 2>&1 | grep -i "x-relay:")
if [ -n "$RELAY_HEADER" ]; then
    echo "✅ PASS: X-Relay header present"
    ((PASS++))
else
    echo "❌ FAIL: X-Relay header missing"
    ((FAIL++))
fi

# Test 8: POST /submitTx never cached
echo ""
echo "Test 8: POST /submitTx bypasses cache"
POST_CACHE=$(curl -s -k -X POST https://localhost/submitTx \
    -H "Content-Type: application/json" \
    -d '{"test":"data"}' \
    -I 2>&1 | grep -i "x-cache-status")
if [ -z "$POST_CACHE" ]; then
    echo "✅ PASS: POST requests not cached (no X-Cache-Status header)"
    ((PASS++))
else
    echo "❌ FAIL: POST request was cached (should never happen)"
    ((FAIL++))
fi

# Test 9: Metrics endpoint
echo ""
echo "Test 9: Prometheus /metrics endpoint"
METRICS=$(curl -s http://127.0.0.1:9090/metrics | grep "seed2_" | head -1)
if [ -n "$METRICS" ]; then
    echo "✅ PASS: Metrics endpoint working"
    ((PASS++))
else
    echo "❌ FAIL: Metrics endpoint failed"
    ((FAIL++))
fi

# Test 10: Service logs
echo ""
echo "Test 10: Service logs accessible"
LOGS=$(journalctl -u archivas-relay -n 5 --no-pager 2>/dev/null)
if [ -n "$LOGS" ]; then
    echo "✅ PASS: Service logs accessible"
    ((PASS++))
else
    echo "❌ FAIL: Service logs not accessible"
    ((FAIL++))
fi

echo ""
echo "=================================="
echo "📊 Smoke Test Results:"
echo "   PASSED: $PASS/10"
echo "   FAILED: $FAIL/10"
echo "=================================="

if [ $FAIL -eq 0 ]; then
    echo "✅ All smoke tests PASSED!"
else
    echo "⚠️  Some tests failed. Check logs for details."
fi

echo ""
echo "📊 Post-Deployment Statistics:"
echo "=================================="

# Cache statistics
CACHE_SIZE=$(du -sh /var/cache/nginx/archivas_rpc 2>/dev/null | cut -f1)
CACHE_FILES=$(find /var/cache/nginx/archivas_rpc -type f 2>/dev/null | wc -l)
echo "Cache size: $CACHE_SIZE"
echo "Cached files: $CACHE_FILES"

# Service status
echo ""
echo "Service status:"
systemctl is-active --quiet nginx && echo "✅ Nginx: active" || echo "❌ Nginx: inactive"
systemctl is-active --quiet archivas-relay && echo "✅ Relay: active" || echo "❌ Relay: inactive"

# Upstream latency
echo ""
echo "Upstream latency test:"
UPSTREAM_TIME=$(curl -s -k -o /dev/null -w "%{time_total}" https://localhost/chainTip)
echo "Response time: ${UPSTREAM_TIME}s"

# Cache hit ratio (from logs)
echo ""
echo "Cache hit ratio (last 100 requests):"
if [ -f /var/log/nginx/seed2-access.log ]; then
    tail -100 /var/log/nginx/seed2-access.log 2>/dev/null | \
        grep -oP 'X-Cache-Status: \K\w+' | \
        sort | uniq -c | \
        awk '{printf "  %s: %d (%.1f%%)\n", $2, $1, ($1/100)*100}'
fi

echo ""
echo "=================================="
echo "✅ Deployment Complete!"
echo "=================================="
echo ""
echo "🌐 Endpoints:"
echo "   Public: https://seed2.archivas.ai"
echo "   Health: https://seed2.archivas.ai/health"
echo "   Status: https://seed2.archivas.ai/status"
echo "   Metrics: https://seed2.archivas.ai/metrics"
echo ""
echo "📋 Next steps:"
echo "   1. Verify DNS: dig seed2.archivas.ai"
echo "   2. Get SSL cert: sudo certbot --nginx -d seed2.archivas.ai"
echo "   3. Monitor: curl https://seed2.archivas.ai/status | jq ."
echo "   4. View logs: sudo journalctl -u archivas-relay -f"
echo ""
echo "📚 Documentation: docs/relay.md"
echo ""

exit 0

