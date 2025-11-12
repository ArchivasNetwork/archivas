#!/bin/bash

# Server A Status Check Script
# Checks node status, resources, logs, and potential DDoS indicators

echo "=========================================="
echo "Server A (seed.archivas.ai) Status Check"
echo "=========================================="
echo ""

# 1. Check if node process is running
echo "1. Node Process Status:"
echo "----------------------"
if pgrep -f archivas-node > /dev/null; then
    echo "✅ Node is running"
    ps aux | grep archivas-node | grep -v grep | awk '{print "   PID:", $2, "CPU:", $3"%", "MEM:", $4"%", "TIME:", $10}'
else
    echo "❌ Node is NOT running"
fi
echo ""

# 2. Check system resources
echo "2. System Resources:"
echo "-------------------"
echo "CPU Load:"
uptime
echo ""
echo "Memory Usage:"
free -h
echo ""
echo "Disk Usage:"
df -h / | tail -1
echo ""

# 3. Check if RPC port is listening
echo "3. Network Ports:"
echo "----------------"
if sudo ss -tlnp | grep -q ":8080"; then
    echo "✅ Port 8080 (RPC) is listening"
    sudo ss -tlnp | grep ":8080" | head -1
else
    echo "❌ Port 8080 (RPC) is NOT listening"
fi
echo ""

if sudo ss -tlnp | grep -q ":9090"; then
    echo "✅ Port 9090 (P2P) is listening"
    sudo ss -tlnp | grep ":9090" | head -1
else
    echo "❌ Port 9090 (P2P) is NOT listening"
fi
echo ""

# 4. Check local RPC endpoints
echo "4. Local RPC Endpoints:"
echo "----------------------"
echo -n "Health endpoint: "
if curl -s --max-time 5 http://localhost:8080/healthz > /dev/null 2>&1; then
    echo "✅ Responding"
    curl -s --max-time 5 http://localhost:8080/healthz | head -1
else
    echo "❌ Not responding or timeout"
fi
echo ""

echo -n "ChainTip endpoint: "
if curl -s --max-time 5 http://localhost:8080/chainTip > /dev/null 2>&1; then
    echo "✅ Responding"
    curl -s --max-time 5 http://localhost:8080/chainTip | head -1
else
    echo "❌ Not responding or timeout"
fi
echo ""

# 5. Check Nginx status
echo "5. Nginx Status:"
echo "---------------"
if systemctl is-active --quiet nginx; then
    echo "✅ Nginx is running"
    systemctl status nginx --no-pager -l | head -5
else
    echo "❌ Nginx is NOT running"
fi
echo ""

# 6. Check recent node logs for errors
echo "6. Recent Node Logs (last 50 lines):"
echo "-----------------------------------"
sudo journalctl -u archivas-node -n 50 --no-pager | tail -20
echo ""

# 7. Check for error patterns
echo "7. Error Patterns:"
echo "-----------------"
ERROR_COUNT=$(sudo journalctl -u archivas-node --since "10 minutes ago" --no-pager | grep -i "error\|panic\|fatal" | wc -l)
echo "Errors in last 10 minutes: $ERROR_COUNT"
if [ "$ERROR_COUNT" -gt 10 ]; then
    echo "⚠️  High error count - possible issue"
    sudo journalctl -u archivas-node --since "10 minutes ago" --no-pager | grep -i "error\|panic\|fatal" | tail -10
else
    echo "✅ Low error count"
fi
echo ""

# 8. Check connection counts
echo "8. Active Connections:"
echo "---------------------"
echo "TCP connections to port 8080:"
sudo ss -tn | grep ":8080" | wc -l
echo ""
echo "TCP connections to port 9090:"
sudo ss -tn | grep ":9090" | wc -l
echo ""

# 9. Check for DDoS indicators
echo "9. Potential DDoS Indicators:"
echo "----------------------------"
echo "Unique IPs connecting to port 8080 (last 100 connections):"
sudo ss -tn | grep ":8080" | awk '{print $5}' | cut -d: -f1 | sort | uniq -c | sort -rn | head -10
echo ""

# 10. Check node service status
echo "10. Systemd Service Status:"
echo "--------------------------"
sudo systemctl status archivas-node --no-pager -l | head -15
echo ""

# 11. Check Nginx access logs for patterns
echo "11. Nginx Access Log Analysis (last 100 requests):"
echo "-------------------------------------------------"
if [ -f /var/log/nginx/seed.archivas.ai.access.log ]; then
    echo "Requests per minute (last 10 minutes):"
    sudo tail -100 /var/log/nginx/seed.archivas.ai.access.log | awk '{print $4}' | cut -d: -f1-2 | uniq -c | tail -10
    echo ""
    echo "Top requesting IPs:"
    sudo tail -100 /var/log/nginx/seed.archivas.ai.access.log | awk '{print $1}' | sort | uniq -c | sort -rn | head -10
    echo ""
    echo "Status codes:"
    sudo tail -100 /var/log/nginx/seed.archivas.ai.access.log | awk '{print $9}' | sort | uniq -c | sort -rn
else
    echo "⚠️  Nginx access log not found"
fi
echo ""

# 12. Check Nginx error logs
echo "12. Nginx Error Log (last 20 lines):"
echo "-----------------------------------"
if [ -f /var/log/nginx/seed.archivas.ai.error.log ]; then
    sudo tail -20 /var/log/nginx/seed.archivas.ai.error.log
else
    echo "⚠️  Nginx error log not found"
fi
echo ""

echo "=========================================="
echo "Status Check Complete"
echo "=========================================="

