#!/usr/bin/env bash
# Deploy final stability fix for seed node
# Addresses: goroutine accumulation, memory exhaustion, resource limits

set -euo pipefail

echo "🔧 Archivas Seed Node - Final Stability Fix"
echo "==========================================="
echo ""

# 1. Pull latest code
echo "1️⃣  Pulling latest code..."
cd /home/ubuntu/archivas
git pull
echo "   ✅ Code updated"
echo ""

# 2. Verify the fix is present
echo "2️⃣  Verifying backpressure fix..."
if grep -q "persistSem" cmd/archivas-node/main.go; then
    echo "   ✅ Backpressure mechanism found"
else
    echo "   ❌ Fix not found - code may not have pulled correctly"
    exit 1
fi
echo ""

# 3. Build node binary
echo "3️⃣  Building node binary..."
go build -o archivas-node ./cmd/archivas-node
echo "   ✅ Binary built"
echo ""

# 4. Stop node service
echo "4️⃣  Stopping node service..."
sudo systemctl stop archivas-node.service || true
sleep 3
echo "   ✅ Node stopped"
echo ""

# 5. Check for stale lock file
echo "5️⃣  Checking for stale lock file..."
if [ -f /home/ubuntu/archivas/data/LOCK ]; then
    echo "   ⚠️  Removing stale LOCK file..."
    rm -f /home/ubuntu/archivas/data/LOCK
    echo "   ✅ LOCK file removed"
else
    echo "   ✅ No stale LOCK file"
fi
echo ""

# 6. Copy binary
echo "6️⃣  Installing new binary..."
sudo cp archivas-node /usr/local/bin/archivas-node
sudo chmod +x /usr/local/bin/archivas-node
echo "   ✅ Binary installed"
echo ""

# 7. Update systemd service with resource limits
echo "7️⃣  Updating systemd service with resource limits..."
sudo cp deploy/seed/archivas-node-stable.service /etc/systemd/system/archivas-node.service
sudo systemctl daemon-reload
echo "   ✅ Service updated"
echo ""

# 8. Start node service
echo "8️⃣  Starting node service..."
sudo systemctl start archivas-node.service
sleep 5
echo "   ✅ Node started"
echo ""

# 9. Verify node is running
echo "9️⃣  Verifying node status..."
if systemctl is-active --quiet archivas-node.service; then
    echo "   ✅ Node is active"
    PID=$(pgrep -f archivas-node)
    echo "   PID: $PID"
    echo "   Uptime: $(ps -p $PID -o etime --no-headers 2>/dev/null || echo 'N/A')"
else
    echo "   ❌ Node failed to start"
    echo "   Check logs: sudo journalctl -u archivas-node.service -n 50"
    exit 1
fi
echo ""

# 10. Test RPC endpoints
echo "🔟 Testing RPC endpoints..."
echo "   /ping:"
timeout 3 curl -s http://127.0.0.1:8080/ping | jq -r '.status' || echo "   TIMEOUT"
echo "   /chainTip:"
timeout 3 curl -s http://127.0.0.1:8080/chainTip | jq -r '.height' || echo "   TIMEOUT"
echo ""

# 11. Check resource limits
echo "1️⃣1️⃣  Checking resource limits..."
PID=$(pgrep -f archivas-node)
echo "   Memory limit: $(systemctl show archivas-node.service -p MemoryMax --value)"
echo "   Task limit: $(systemctl show archivas-node.service -p TasksMax --value)"
echo "   Current memory: $(ps -p $PID -o rss --no-headers | awk '{print $1/1024 " MB"}')"
echo ""

# 12. Monitor for 30 seconds
echo "1️⃣2️⃣  Monitoring stability (30 seconds)..."
for i in {1..3}; do
    HEIGHT=$(timeout 3 curl -s http://127.0.0.1:8080/chainTip | jq -r '.height' || echo "ERROR")
    MEMORY=$(ps -p $PID -o rss --no-headers 2>/dev/null | awk '{print $1/1024 " MB"}' || echo "N/A")
    echo "   [$i/3] Height: $HEIGHT | Memory: $MEMORY"
    sleep 10
done
echo ""

echo "✅ Deployment complete!"
echo ""
echo "📊 Next steps:"
echo "   1. Monitor node: watch -n 5 'curl -s http://127.0.0.1:8080/chainTip | jq -r \".height\"'"
echo "   2. Check logs: sudo journalctl -u archivas-node.service -f"
echo "   3. Check resources: systemctl status archivas-node.service"
echo ""
echo "🛡️  Stability features enabled:"
echo "   ✓ Backpressure mechanism (max 5 concurrent disk writes)"
echo "   ✓ Memory limit: 4GB hard, 3GB soft"
echo "   ✓ Task limit: 2048 (prevents goroutine explosion)"
echo "   ✓ Aggressive GC (GOGC=50)"
echo "   ✓ OOM protection with auto-restart"

