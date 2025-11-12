#!/bin/bash

# Monitor seed node stability (checks process, not RPC)
# This is better than watchdog for seed nodes under DDoS

LOG_FILE="${LOG_FILE:-/tmp/seed-node-monitor.log}"

check_node() {
    if pgrep -f archivas-node > /dev/null; then
        return 0  # Process is running
    else
        return 1  # Process is NOT running
    fi
}

log_status() {
    echo "$(date -u +"%Y-%m-%d %H:%M:%S UTC") - $1" | tee -a "$LOG_FILE"
}

# Check if node is running
if check_node; then
    log_status "✅ Node process is RUNNING"
    
    # Get process info
    PID=$(pgrep -f archivas-node)
    CPU=$(ps -p $PID -o %cpu --no-headers | tr -d ' ')
    MEM=$(ps -p $PID -o %mem --no-headers | tr -d ' ')
    
    log_status "   PID: $PID | CPU: ${CPU}% | MEM: ${MEM}%"
    
    # Check if RPC is responding (informational only)
    if curl -s --max-time 2 http://127.0.0.1:8080/ping > /dev/null 2>&1; then
        log_status "   RPC: ✅ Responding"
    else
        log_status "   RPC: ⚠️ Slow/Unresponsive (node is still running)"
    fi
    
    exit 0
else
    log_status "❌ Node process is NOT RUNNING - ACTION REQUIRED"
    log_status "   Run: sudo systemctl start archivas-node.service"
    exit 1
fi

