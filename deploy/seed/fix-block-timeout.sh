#!/bin/bash
set -euo pipefail

# Fix /submitBlock timeout on Seed1 (Server A)

echo "🔧 Fixing /submitBlock timeout on Seed1..."
echo ""

# Check if already updated
if grep -q "proxy_read_timeout 120s" /etc/nginx/sites-available/archivas-rpc; then
    echo "✅ Already updated! Current timeout is 120s."
else
    echo "📝 Updating Nginx config..."
    
    # Backup
    sudo cp /etc/nginx/sites-available/archivas-rpc /etc/nginx/sites-available/archivas-rpc.backup.$(date +%Y%m%d_%H%M%S)
    
    # Update timeout in /submitBlock location
    sudo sed -i '/location = \/submitBlock/,/^  }/s/proxy_read_timeout 30s;/proxy_read_timeout 120s;\n    proxy_send_timeout 120s;/' /etc/nginx/sites-available/archivas-rpc
    
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
echo "  grep -A5 'location = /submitBlock' /etc/nginx/sites-available/archivas-rpc"

