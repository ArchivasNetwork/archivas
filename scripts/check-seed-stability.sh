#!/bin/bash

# Simple diagnostic script for seed node stability
# Uses basic commands: cat, grep, etc.

echo "=========================================="
echo "Seed Node Stability Check"
echo "=========================================="
echo ""

# 1. Check if node is running
echo "1. Node Process:"
echo "---------------"
if pgrep -f archivas-node > /dev/null; then
    echo "✅ Node is RUNNING"
    ps aux | grep archivas-node | grep -v grep
else
    echo "❌ Node is NOT running"
fi
echo ""

# 2. Check recent restarts
echo "2. Recent Restarts (last hour):"
echo "------------------------------"
systemctl status archivas-node.service --no-pager -l | grep -E "Active|since" | head -2
journalctl -u archivas-node.service --since "1 hour ago" --no-pager | grep -E "Stopped|Started" | tail -10
echo ""

# 3. Check watchdog logs
echo "3. Watchdog Activity (last hour):"
echo "--------------------------------"
journalctl -u archivas-watchdog.service --since "1 hour ago" --no-pager | tail -20
echo ""

# 4. Check for crashes/errors in node logs
echo "4. Node Log Errors (last 100 lines):"
echo "-----------------------------------"
if [ -f ~/archivas/logs/node-error.log ]; then
    tail -100 ~/archivas/logs/node-error.log | grep -i -E "error|panic|fatal|crash" | tail -10
else
    echo "No error log file found"
fi
echo ""

# 5. Check system resources
echo "5. System Resources:"
echo "-------------------"
echo "Memory:"
free -h | head -2
echo ""
echo "Disk:"
df -h / | tail -1
echo ""
echo "CPU Load:"
uptime
echo ""

# 6. Check RPC response
echo "6. RPC Health:"
echo "-------------"
if curl -s --max-time 5 http://127.0.0.1:8080/ping > /dev/null 2>&1; then
    echo "✅ RPC /ping responding"
    curl -s http://127.0.0.1:8080/healthz | head -1
else
    echo "❌ RPC not responding"
fi
echo ""

# 7. Check for OOM kills
echo "7. OOM Kills:"
echo "------------"
dmesg | grep -i "oom\|killed process" | tail -5
echo ""

# 8. Check database corruption
echo "8. Database Status:"
echo "------------------"
if [ -d ~/archivas/data ]; then
    echo "Database directory exists"
    ls -lh ~/archivas/data | head -5
else
    echo "No database directory found"
fi
echo ""

echo "=========================================="
echo "Check Complete"
echo "=========================================="

