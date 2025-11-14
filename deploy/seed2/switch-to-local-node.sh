#!/bin/bash
set -euo pipefail

# Switch Seed2 Nginx to use local node instead of proxying to Seed1
# This makes farmers truly use Seed2's node, with Seed1 as backup

echo "=== Switch Seed2 Nginx to Local Node ==="
echo ""
echo "This will update Nginx to use:"
echo "  Primary: Seed2 local node (127.0.0.1:8082)"
echo "  Backup:  Seed1 Nginx (seed.archivas.ai:8081)"
echo ""

# Check if running on Seed2
if [ ! -f "/etc/nginx/sites-available/archivas-seed2" ]; then
    echo "Error: This script must be run on Seed2 server"
    echo "Nginx config not found at /etc/nginx/sites-available/archivas-seed2"
    exit 1
fi

# Backup current config
echo "1. Backing up current Nginx config..."
sudo cp /etc/nginx/sites-available/archivas-seed2 \
       /etc/nginx/sites-available/archivas-seed2.backup-$(date +%Y%m%d-%H%M%S)
echo "   ✓ Backup created"

# Update proxy_pass directives to use read_pool
echo ""
echo "2. Updating proxy_pass directives..."

sudo sed -i 's|proxy_pass http://seed1_backend;|proxy_pass http://read_pool;|g' \
    /etc/nginx/sites-available/archivas-seed2

echo "   ✓ Updated all proxy_pass directives to use read_pool"

# Test Nginx configuration
echo ""
echo "3. Testing Nginx configuration..."
if sudo nginx -t; then
    echo "   ✓ Nginx config is valid"
else
    echo "   ✗ Nginx config has errors! Restoring backup..."
    sudo cp /etc/nginx/sites-available/archivas-seed2.backup-$(date +%Y%m%d)* \
           /etc/nginx/sites-available/archivas-seed2
    exit 1
fi

# Reload Nginx
echo ""
echo "4. Reloading Nginx..."
if sudo systemctl reload nginx; then
    echo "   ✓ Nginx reloaded successfully"
else
    echo "   ✗ Nginx reload failed! Check logs:"
    echo "      sudo journalctl -u nginx -n 20"
    exit 1
fi

# Verify the change
echo ""
echo "5. Verifying configuration..."
echo ""
echo "Testing /chainTip endpoint..."
RESPONSE=$(curl -sk https://seed2.archivas.ai/chainTip 2>/dev/null || echo "")
if [ -n "$RESPONSE" ]; then
    echo "$RESPONSE" | jq -C 2>/dev/null || echo "$RESPONSE"
    echo "   ✓ Endpoint responding"
else
    echo "   ✗ No response from endpoint"
fi

echo ""
echo "Testing /challenge endpoint..."
CHALLENGE=$(curl -sk https://seed2.archivas.ai/challenge 2>/dev/null || echo "")
if [ -n "$CHALLENGE" ]; then
    echo "$CHALLENGE" | jq -C '.height' 2>/dev/null || echo "$CHALLENGE"
    echo "   ✓ Challenge endpoint responding"
else
    echo "   ✗ No response from challenge endpoint"
fi

echo ""
echo "=== ✅ Switch Complete! ==="
echo ""
echo "Farmers can now use: https://seed2.archivas.ai"
echo "  - Primary source: Seed2 local node"
echo "  - Fallback: Seed1 (automatic)"
echo ""
echo "Monitor with:"
echo "  sudo tail -f /var/log/nginx/access.log"
echo "  sudo journalctl -u archivas-node-seed2 -f"
echo ""
echo "To revert to old config:"
echo "  sudo cp /etc/nginx/sites-available/archivas-seed2.backup-* \\"
echo "         /etc/nginx/sites-available/archivas-seed2"
echo "  sudo systemctl reload nginx"

