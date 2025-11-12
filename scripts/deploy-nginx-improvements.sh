#!/bin/bash

# Deploy Nginx improvements for DDoS protection and better timeout handling
# Run this on Server A (seed.archivas.ai)

set -e

echo "=========================================="
echo "Deploying Nginx Improvements to Server A"
echo "=========================================="
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  This script needs sudo privileges"
    echo "Please run: sudo bash scripts/deploy-nginx-improvements.sh"
    exit 1
fi

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

echo "1. Checking current Nginx status..."
echo "-----------------------------------"
systemctl status nginx --no-pager -l | head -5
echo ""

echo "2. Backing up current Nginx config..."
echo "-------------------------------------"
BACKUP_FILE="/etc/nginx/sites-available/seed.archivas.ai.backup.$(date +%Y%m%d_%H%M%S)"
if [ -f "/etc/nginx/sites-available/seed.archivas.ai" ]; then
    cp /etc/nginx/sites-available/seed.archivas.ai "$BACKUP_FILE"
    echo "✅ Backup created: $BACKUP_FILE"
else
    echo "⚠️  Current config not found (this might be a new setup)"
fi
echo ""

echo "3. Copying updated Nginx configuration..."
echo "-----------------------------------------"
if [ -f "$PROJECT_ROOT/deploy/seed/nginx-site.conf" ]; then
    cp "$PROJECT_ROOT/deploy/seed/nginx-site.conf" /etc/nginx/sites-available/seed.archivas.ai
    echo "✅ Config file copied"
else
    echo "❌ Error: nginx-site.conf not found at $PROJECT_ROOT/deploy/seed/nginx-site.conf"
    exit 1
fi
echo ""

echo "4. Testing Nginx configuration..."
echo "--------------------------------"
if nginx -t; then
    echo "✅ Nginx configuration is valid"
else
    echo "❌ Nginx configuration test failed!"
    echo "Restoring backup..."
    if [ -f "$BACKUP_FILE" ]; then
        cp "$BACKUP_FILE" /etc/nginx/sites-available/seed.archivas.ai
        echo "✅ Backup restored"
    fi
    exit 1
fi
echo ""

echo "5. Reloading Nginx..."
echo "-------------------"
if systemctl reload nginx; then
    echo "✅ Nginx reloaded successfully"
else
    echo "❌ Failed to reload Nginx"
    exit 1
fi
echo ""

echo "6. Verifying Nginx is running..."
echo "--------------------------------"
sleep 2
if systemctl is-active --quiet nginx; then
    echo "✅ Nginx is running"
    systemctl status nginx --no-pager -l | head -5
else
    echo "❌ Nginx is not running!"
    echo "Attempting to start..."
    systemctl start nginx
fi
echo ""

echo "7. Testing endpoints..."
echo "----------------------"
echo -n "Health endpoint: "
if curl -s --max-time 5 http://localhost:8080/healthz > /dev/null 2>&1; then
    echo "✅ Responding"
else
    echo "⚠️  Not responding (node may be starting)"
fi

echo -n "Public health (via Nginx): "
if curl -s --max-time 5 https://seed.archivas.ai/healthz > /dev/null 2>&1; then
    echo "✅ Responding"
else
    echo "⚠️  Not responding"
fi
echo ""

echo "=========================================="
echo "Deployment Complete"
echo "=========================================="
echo ""
echo "Improvements deployed:"
echo "  ✅ Rate limiting on all endpoints"
echo "  ✅ Increased timeouts (180s for /blocks/range)"
echo "  ✅ DDoS protection per IP"
echo "  ✅ Specific limits for different endpoint types"
echo ""
echo "Next steps:"
echo "  1. Monitor logs: sudo tail -f /var/log/nginx/seed.archivas.ai.access.log"
echo "  2. Check for rate limiting: sudo tail -f /var/log/nginx/seed.archivas.ai.error.log"
echo "  3. Run status check: bash scripts/check-server-status.sh"
echo ""

