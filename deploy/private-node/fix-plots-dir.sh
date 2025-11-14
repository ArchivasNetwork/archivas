#!/bin/bash

set -e

echo "════════════════════════════════════════════════════════════"
echo "🔧 Archivas Private Farmer - Plots Directory Fix"
echo "════════════════════════════════════════════════════════════"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (use sudo)"
    exit 1
fi

SERVICE_FILE="/etc/systemd/system/archivas-farmer-private.service"

if [ ! -f "$SERVICE_FILE" ]; then
    echo "❌ Service file not found: $SERVICE_FILE"
    exit 1
fi

echo "Current plots directory configuration:"
grep "plots" "$SERVICE_FILE" | grep -v "^#" || true
echo ""

echo "Available plots directories:"
echo ""
for dir in /home/*/archivas-plots*; do
    if [ -d "$dir" ]; then
        count=$(find "$dir" -name "*.arcv" -type f 2>/dev/null | wc -l)
        size=$(du -sh "$dir" 2>/dev/null | cut -f1)
        echo "  📁 $dir"
        echo "     Plots: $count files, Total size: $size"
        echo ""
    fi
done

read -p "Enter the full path to your plots directory: " PLOTS_DIR

if [ ! -d "$PLOTS_DIR" ]; then
    echo "⚠️  Warning: Directory doesn't exist: $PLOTS_DIR"
    read -p "Continue anyway? (y/n): " confirm
    if [ "$confirm" != "y" ]; then
        echo "Cancelled."
        exit 0
    fi
fi

echo ""
echo "Updating service file..."

# Backup first
cp "$SERVICE_FILE" "$SERVICE_FILE.backup.$(date +%s)"

# Update the plots directory
sed -i "s|--plots [^ ]*|--plots $PLOTS_DIR|g" "$SERVICE_FILE"

echo "✓ Service file updated"
echo ""

echo "New configuration:"
grep "plots" "$SERVICE_FILE" | grep -v "^#"
echo ""

read -p "Restart the farmer service now? (y/n): " restart
if [ "$restart" = "y" ]; then
    echo "Reloading systemd..."
    systemctl daemon-reload
    
    echo "Restarting archivas-farmer-private..."
    systemctl restart archivas-farmer-private
    
    echo ""
    echo "Service status:"
    systemctl status archivas-farmer-private --no-pager -l || true
    
    echo ""
    echo "Recent logs:"
    journalctl -u archivas-farmer-private -n 15 --no-pager
fi

echo ""
echo "════════════════════════════════════════════════════════════"
echo "✅ Done!"
echo "════════════════════════════════════════════════════════════"

