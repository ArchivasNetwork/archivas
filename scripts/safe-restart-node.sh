#!/bin/bash
# Safe script to restart node after crash - verifies it's safe to remove lock file

set -e

DB_PATH="/home/ubuntu/archivas/data"
LOCK_FILE="$DB_PATH/LOCK"
NODE_SERVICE="archivas-node.service"

echo "=== Safe Node Restart Script ==="
echo ""

# 1. Check if node process is actually running
echo "1. Checking if node process is running..."
NODE_PID=$(pgrep -f "archivas-node" || true)
if [ ! -z "$NODE_PID" ]; then
    echo "   ⚠️  WARNING: Node process is still running (PID: $NODE_PID)"
    echo "   Checking if it's stuck..."
    
    # Check if process is responsive
    if kill -0 "$NODE_PID" 2>/dev/null; then
        echo "   ⚠️  Process exists but may be unresponsive"
        echo "   Attempting graceful shutdown first..."
        sudo systemctl stop "$NODE_SERVICE" || true
        sleep 5
        
        # Check if it's still running
        if pgrep -f "archivas-node" > /dev/null; then
            echo "   ⚠️  Process still running after stop, forcing kill..."
            sudo kill -9 "$NODE_PID" 2>/dev/null || true
            sleep 2
        fi
    fi
else
    echo "   ✓ No node process running"
fi

# 2. Check if any other process is using the database
echo ""
echo "2. Checking for other processes using database..."
DB_PROCESSES=$(sudo lsof "$DB_PATH" 2>/dev/null | grep -v "COMMAND" | awk '{print $2}' | sort -u || true)
if [ ! -z "$DB_PROCESSES" ]; then
    echo "   ⚠️  WARNING: Other processes may be using database:"
    echo "   PIDs: $DB_PROCESSES"
    for pid in $DB_PROCESSES; do
        ps -p "$pid" -o pid,cmd --no-headers 2>/dev/null || echo "   Process $pid not found"
    done
    echo ""
    echo "   ⚠️  It may not be safe to remove lock file!"
    echo "   Consider stopping these processes first."
    read -p "   Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "   Aborted."
        exit 1
    fi
else
    echo "   ✓ No other processes using database"
fi

# 3. Check if lock file exists
echo ""
echo "3. Checking for lock file..."
if [ -f "$LOCK_FILE" ]; then
    echo "   ⚠️  Lock file exists: $LOCK_FILE"
    echo "   Lock file details:"
    ls -lah "$LOCK_FILE"
    
    # Check when lock file was last modified
    LOCK_AGE=$(($(date +%s) - $(stat -c %Y "$LOCK_FILE" 2>/dev/null || echo 0)))
    echo "   Lock file age: $LOCK_AGE seconds ($(($LOCK_AGE / 60)) minutes)"
    
    if [ $LOCK_AGE -lt 10 ]; then
        echo "   ⚠️  WARNING: Lock file is very recent (< 10 seconds old)"
        echo "   Database might still be in use!"
        read -p "   Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "   Aborted."
            exit 1
        fi
    fi
else
    echo "   ✓ No lock file found (database should be safe to open)"
fi

# 4. Verify node service status
echo ""
echo "4. Checking node service status..."
if systemctl is-active --quiet "$NODE_SERVICE"; then
    echo "   ⚠️  WARNING: Node service is still active"
    echo "   Stopping service..."
    sudo systemctl stop "$NODE_SERVICE"
    sleep 3
fi

if systemctl is-failed --quiet "$NODE_SERVICE"; then
    echo "   ✓ Service is in failed state (expected after crash)"
    sudo systemctl reset-failed "$NODE_SERVICE"
fi

# 5. Create backup of database before removing lock (optional but recommended)
echo ""
echo "5. Database safety check..."
DB_SIZE=$(du -sh "$DB_PATH" 2>/dev/null | awk '{print $1}' || echo "unknown")
echo "   Database size: $DB_SIZE"
echo "   ⚠️  NOTE: Removing lock file does not modify database data"
echo "   BadgerDB will perform automatic recovery on next open"

# 6. Remove lock file if it exists
echo ""
echo "6. Removing lock file..."
if [ -f "$LOCK_FILE" ]; then
    echo "   Removing: $LOCK_FILE"
    rm -f "$LOCK_FILE"
    
    if [ ! -f "$LOCK_FILE" ]; then
        echo "   ✓ Lock file removed successfully"
    else
        echo "   ⚠️  ERROR: Failed to remove lock file"
        exit 1
    fi
else
    echo "   ✓ Lock file does not exist (nothing to remove)"
fi

# 7. Verify no lock file remains
echo ""
echo "7. Verifying lock file is gone..."
if [ ! -f "$LOCK_FILE" ]; then
    echo "   ✓ Confirmed: No lock file"
else
    echo "   ⚠️  ERROR: Lock file still exists!"
    exit 1
fi

# 8. Restart node service
echo ""
echo "8. Restarting node service..."
sudo systemctl restart "$NODE_SERVICE"

# Wait for service to start
sleep 5

# 9. Check if node started successfully
echo ""
echo "9. Checking node status..."
if systemctl is-active --quiet "$NODE_SERVICE"; then
    echo "   ✓ Node service is running"
else
    echo "   ⚠️  WARNING: Node service is not running"
    echo "   Checking logs..."
    sudo journalctl -u "$NODE_SERVICE" --since "1 minute ago" --no-pager | tail -20
    exit 1
fi

# 10. Test if node is responsive
echo ""
echo "10. Testing node responsiveness..."
sleep 2

if timeout 5 curl -s http://127.0.0.1:8080/ping > /dev/null; then
    echo "   ✓ Node is responsive (/ping works)"
else
    echo "   ⚠️  WARNING: Node is not responding to /ping"
    echo "   Check logs: sudo journalctl -u $NODE_SERVICE -f"
fi

# 11. Test /challenge endpoint
echo ""
echo "11. Testing /challenge endpoint..."
if timeout 5 curl -s http://127.0.0.1:8080/challenge > /dev/null; then
    echo "   ✓ /challenge endpoint is working"
    curl -s http://127.0.0.1:8080/challenge | jq -r '.height' | head -1
else
    echo "   ⚠️  WARNING: /challenge endpoint is not responding"
    echo "   This might indicate the node is still starting up"
    echo "   Wait a few seconds and try: curl http://127.0.0.1:8080/challenge"
fi

echo ""
echo "=== Restart Complete ==="
echo "Node should now be running. Monitor with:"
echo "  sudo journalctl -u $NODE_SERVICE -f"
echo ""
echo "Check node status:"
echo "  curl http://127.0.0.1:8080/chainTip | jq"

