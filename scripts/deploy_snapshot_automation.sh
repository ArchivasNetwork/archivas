#!/bin/bash
# Archivas Snapshot Publishing - Production Deployment
# Run this script on Seed2 to enable automated snapshot publishing

set -e

echo "=============================================="
echo "  Archivas Snapshot Publishing Deployment"
echo "=============================================="
echo ""

# Verify we're on the right server
if [ ! -f "/etc/nginx/sites-available/archivas-seed2" ]; then
    echo "❌ Error: This script must be run on Seed2"
    exit 1
fi

# Create log directory
echo "📁 Creating log directory..."
sudo mkdir -p /var/log/archivas
sudo chown ubuntu:ubuntu /var/log/archivas

# Verify the publish script exists
if [ ! -f "/home/ubuntu/archivas/scripts/publish_snapshot_devnet.sh" ]; then
    echo "❌ Error: publish_snapshot_devnet.sh not found"
    echo "   Please ensure the repo is up to date: cd ~/archivas && git pull"
    exit 1
fi

# Make the script executable
chmod +x /home/ubuntu/archivas/scripts/publish_snapshot_devnet.sh

# Test run the script
echo ""
echo "🧪 Running test snapshot publish..."
if /home/ubuntu/archivas/scripts/publish_snapshot_devnet.sh; then
    echo "✅ Test snapshot published successfully!"
else
    echo "❌ Test snapshot failed"
    exit 1
fi

# Set up cron job (runs every 6 hours at minute 0)
echo ""
echo "⏰ Setting up cron job..."
CRON_LINE="0 */6 * * * /home/ubuntu/archivas/scripts/publish_snapshot_devnet.sh >> /var/log/archivas/snapshot-publish.log 2>&1"

# Remove existing cron job if present
crontab -l 2>/dev/null | grep -v "publish_snapshot_devnet.sh" | crontab - 2>/dev/null || true

# Add new cron job
(crontab -l 2>/dev/null; echo "$CRON_LINE") | crontab -

echo "✅ Cron job configured:"
crontab -l | grep publish_snapshot_devnet.sh

# Verify snapshot is accessible
echo ""
echo "🌐 Verifying snapshot accessibility..."
if curl -s -f https://seed2.archivas.ai/devnet/latest.json > /dev/null; then
    echo "✅ Snapshot manifest is accessible at https://seed2.archivas.ai/devnet/latest.json"
    
    # Show the latest snapshot info
    echo ""
    echo "📦 Current snapshot:"
    curl -s https://seed2.archivas.ai/devnet/latest.json | jq .
else
    echo "❌ Snapshot manifest not accessible"
    exit 1
fi

echo ""
echo "=============================================="
echo "  ✅ DEPLOYMENT COMPLETE!"
echo "=============================================="
echo ""
echo "📅 Schedule: Snapshots will be published every 6 hours"
echo "   (00:00, 06:00, 12:00, 18:00 UTC)"
echo ""
echo "📝 Logs: /var/log/archivas/snapshot-publish.log"
echo ""
echo "🔍 Monitor: tail -f /var/log/archivas/snapshot-publish.log"
echo ""
echo "🧪 Test bootstrap: archivas-node bootstrap --network devnet"
echo ""

