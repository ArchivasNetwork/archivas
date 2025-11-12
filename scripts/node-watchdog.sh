#!/usr/bin/env bash
# Archivas Node Watchdog
# Monitors RPC health and auto-restarts if unresponsive
# v1.2.2

set -euo pipefail

RPC_URL="${RPC_URL:-http://127.0.0.1:8080}"
CHECK_INTERVAL="${CHECK_INTERVAL:-60}"  # Check every 60 seconds
TIMEOUT="${TIMEOUT:-30}"                # 30 second timeout (increased for heavy load)
SERVICE_NAME="${SERVICE_NAME:-archivas-node}"
MAX_FAILURES="${MAX_FAILURES:-5}"       # Allow 5 consecutive failures before restart (was 3)

LOG_PREFIX="[watchdog]"

log() {
    echo "$(date -u +"%Y-%m-%dT%H:%M:%SZ") $LOG_PREFIX $1"
}

check_rpc() {
    # Use /healthz endpoint which is lighter than /chainTip
    # Check both endpoint and process to avoid false positives
    if pgrep -f "archivas-node" > /dev/null 2>&1; then
        # Process is running, check if RPC responds (with increased timeout for heavy load)
        if curl -s --max-time "$TIMEOUT" "$RPC_URL/healthz" > /dev/null 2>&1; then
            return 0  # Healthy
        else
            return 1  # RPC unresponsive but process running
        fi
    else
        return 1  # Process not running
    fi
}

log "Starting watchdog for $SERVICE_NAME (checking $RPC_URL/healthz every ${CHECK_INTERVAL}s, timeout ${TIMEOUT}s, max failures ${MAX_FAILURES})"

consecutive_failures=0
max_failures="${MAX_FAILURES}"

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

