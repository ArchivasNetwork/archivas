#!/usr/bin/env bash
# Archivas Node Watchdog
# Monitors RPC health and auto-restarts if unresponsive
# v1.2.2

set -euo pipefail

RPC_URL="${RPC_URL:-http://127.0.0.1:8080}"
CHECK_INTERVAL="${CHECK_INTERVAL:-60}"  # Check every 60 seconds
TIMEOUT="${TIMEOUT:-5}"                 # 5 second timeout
SERVICE_NAME="${SERVICE_NAME:-archivas-node}"

LOG_PREFIX="[watchdog]"

log() {
    echo "$(date -u +"%Y-%m-%dT%H:%M:%SZ") $LOG_PREFIX $1"
}

check_rpc() {
    if curl -s --max-time "$TIMEOUT" "$RPC_URL/chainTip" > /dev/null 2>&1; then
        return 0  # Healthy
    else
        return 1  # Unhealthy
    fi
}

log "Starting watchdog for $SERVICE_NAME (checking $RPC_URL every ${CHECK_INTERVAL}s)"

consecutive_failures=0
max_failures=3

while true; do
    if check_rpc; then
        # Healthy
        if [ $consecutive_failures -gt 0 ]; then
            log "RPC recovered after $consecutive_failures failed check(s)"
        fi
        consecutive_failures=0
    else
        # Unhealthy
        consecutive_failures=$((consecutive_failures + 1))
        log "RPC check failed ($consecutive_failures/$max_failures)"
        
        if [ $consecutive_failures -ge $max_failures ]; then
            log "RPC unresponsive after $max_failures checks - restarting $SERVICE_NAME"
            
            # Restart service
            if sudo systemctl restart "$SERVICE_NAME"; then
                log "Service restarted successfully"
                consecutive_failures=0
                sleep 15  # Give it time to start
            else
                log "ERROR: Failed to restart service!"
            fi
        fi
    fi
    
    sleep "$CHECK_INTERVAL"
done

