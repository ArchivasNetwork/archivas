#!/bin/bash
set -euo pipefail

# Fix /submitBlock timeout on Seed2 (Server D)

echo "🔧 Fixing /submitBlock timeout on Seed2..."
echo ""

# Check if already updated
if grep -q "proxy_read_timeout 120s" /etc/nginx/sites-available/archivas-seed2; then
    echo "✅ Already updated! Current timeout is 120s."
else
    echo "📝 Updating Nginx config..."
    
    # Backup
    sudo cp /etc/nginx/sites-available/archivas-seed2 /etc/nginx/sites-available/archivas-seed2.backup.$(date +%Y%m%d_%H%M%S)
    
    # Update timeout in /submitBlock location
    sudo sed -i '/location = \/submitBlock/,/}/s/proxy_read_timeout 30s;/proxy_read_timeout 120s;/' /etc/nginx/sites-available/archivas-seed2
    sudo sed -i '/location = \/submitBlock/,/}/s/proxy_send_timeout 30s;/proxy_send_timeout 120s;/' /etc/nginx/sites-available/archivas-seed2
    
    echo "✅ Config updated"
fi

echo ""
echo "🧪 Testing Nginx config..."
sudo nginx -t

echo ""
echo "🔄 Reloading Nginx..."
sudo systemctl reload nginx

echo ""
echo "✅ Done! /submitBlock timeout is now 120 seconds."
echo ""
echo "Verify:"
echo "  grep -A5 'location = /submitBlock' /etc/nginx/sites-available/archivas-seed2"

