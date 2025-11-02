#!/usr/bin/env bash
# decommission-node.sh - Safely decommission an Archivas node
# Usage: bash scripts/decommission-node.sh

set -euo pipefail

echo "🛑 Archivas Node Decommissioning"
echo "=================================="
echo ""

# Check if running as root (needed for systemd operations)
if [[ $EUID -eq 0 ]]; then
    SUDO=""
else
    SUDO="sudo"
    echo "⚠️  Some operations require sudo access"
    echo ""
fi

# Step 1: Stop systemd services if they exist
echo "1️⃣  Stopping systemd services..."
if command -v systemctl &> /dev/null; then
    for service in archivas-node archivas-timelord archivas-farmer; do
        if systemctl is-active --quiet "$service" 2>/dev/null; then
            echo "   Stopping $service..."
            $SUDO systemctl stop "$service" || true
        fi
        
        if systemctl is-enabled --quiet "$service" 2>/dev/null; then
            echo "   Disabling $service..."
            $SUDO systemctl disable "$service" || true
        fi
    done
    echo "   ✅ Systemd services stopped and disabled"
else
    echo "   ⏭️  Systemd not found, skipping service management"
fi
echo ""

# Step 2: Kill any stray Archivas processes
echo "2️⃣  Killing stray Archivas processes..."
KILLED=0
for proc in archivas-node archivas-timelord archivas-farmer archivas-wallet archivas-cli archivas-explorer archivas-registry; do
    if pgrep -f "$proc" > /dev/null 2>&1; then
        echo "   Killing $proc processes..."
        pkill -9 -f "$proc" || true
        KILLED=1
    fi
done

if [[ $KILLED -eq 0 ]]; then
    echo "   ✅ No stray processes found"
else
    echo "   ✅ Stray processes killed"
    sleep 2
fi
echo ""

# Step 3: Remove systemd unit files
echo "3️⃣  Removing systemd unit files..."
if [[ -d /etc/systemd/system ]]; then
    UNITS_REMOVED=0
    for unit in /etc/systemd/system/archivas-*.service; do
        if [[ -f "$unit" ]]; then
            echo "   Removing $unit..."
            $SUDO rm -f "$unit"
            UNITS_REMOVED=1
        fi
    done
    
    if [[ $UNITS_REMOVED -eq 1 ]]; then
        echo "   Reloading systemd daemon..."
        $SUDO systemctl daemon-reload
        echo "   ✅ Systemd unit files removed"
    else
        echo "   ✅ No systemd unit files found"
    fi
else
    echo "   ⏭️  /etc/systemd/system not found, skipping"
fi
echo ""

# Step 4: Remove data directories
echo "4️⃣  Removing data directories..."

# User home directories
if [[ -d ~/.archivas ]]; then
    echo "   Removing ~/.archivas..."
    rm -rf ~/.archivas
fi

if [[ -d ~/archivas/data ]]; then
    echo "   Removing ~/archivas/data..."
    rm -rf ~/archivas/data
fi

# System directories (require sudo)
if [[ -d /var/lib/archivas ]]; then
    echo "   Removing /var/lib/archivas..."
    $SUDO rm -rf /var/lib/archivas || true
fi

if [[ -d /opt/archivas ]]; then
    echo "   Removing /opt/archivas..."
    $SUDO rm -rf /opt/archivas || true
fi

if [[ -d /var/log/archivas ]]; then
    echo "   Removing /var/log/archivas..."
    $SUDO rm -rf /var/log/archivas || true
fi

if [[ -d /etc/archivas ]]; then
    echo "   Removing /etc/archivas..."
    $SUDO rm -rf /etc/archivas || true
fi

echo "   ✅ Data directories removed"
echo ""

# Step 5: Verify cleanup
echo "5️⃣  Verifying cleanup..."
VERIFY_OK=1

# Check for running processes
if pgrep -f "archivas" > /dev/null 2>&1; then
    echo "   ⚠️  Warning: Archivas processes still running"
    ps aux | grep archivas | grep -v grep
    VERIFY_OK=0
else
    echo "   ✅ No Archivas processes running"
fi

# Check for remaining systemd units
if command -v systemctl &> /dev/null; then
    if systemctl list-units --all | grep -q archivas; then
        echo "   ⚠️  Warning: Archivas systemd units still present"
        systemctl list-units --all | grep archivas
        VERIFY_OK=0
    else
        echo "   ✅ No Archivas systemd units found"
    fi
fi

# Check for remaining data directories
DIRS_REMAINING=0
for dir in ~/.archivas ~/archivas/data /var/lib/archivas /opt/archivas /var/log/archivas /etc/archivas; do
    if [[ -d "$dir" ]]; then
        echo "   ⚠️  Warning: Directory still exists: $dir"
        DIRS_REMAINING=1
        VERIFY_OK=0
    fi
done

if [[ $DIRS_REMAINING -eq 0 ]]; then
    echo "   ✅ No Archivas data directories found"
fi

echo ""

# Final summary
if [[ $VERIFY_OK -eq 1 ]]; then
    echo "✅ Node decommissioned successfully!"
    echo ""
    echo "The server is now clean and ready for other uses."
else
    echo "⚠️  Decommissioning completed with warnings (see above)"
    echo ""
    echo "You may need to manually clean up remaining items."
fi

echo ""
echo "To restore this node later:"
echo "  1. Clone the repository: git clone https://github.com/ArchivasNetwork/archivas.git"
echo "  2. Follow the deployment guide in docs/"
echo ""

