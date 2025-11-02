#!/usr/bin/env bash
# setup-seed-nginx.sh - Idempotent setup for seed.archivas.ai Nginx reverse proxy
# Usage: sudo bash scripts/setup-seed-nginx.sh

set -euo pipefail

DOMAIN="seed.archivas.ai"
WEBROOT="/var/www/seed"
NGINX_SITE="/etc/nginx/sites-available/seed.archivas.ai"
NGINX_ENABLED="/etc/nginx/sites-enabled/seed.archivas.ai"

echo "🌱 Setting up seed.archivas.ai Nginx proxy"
echo "=========================================="
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root (use sudo)"
   exit 1
fi

# Step 1: Install Nginx if not present
echo "1️⃣  Installing Nginx..."
if ! command -v nginx &> /dev/null; then
    apt-get update
    apt-get install -y nginx
    echo "   ✅ Nginx installed"
else
    echo "   ✅ Nginx already installed"
fi
echo ""

# Step 2: Install Certbot if not present
echo "2️⃣  Installing Certbot..."
if ! command -v certbot &> /dev/null; then
    apt-get install -y certbot python3-certbot-nginx
    echo "   ✅ Certbot installed"
else
    echo "   ✅ Certbot already installed"
fi
echo ""

# Step 3: Create ACME webroot
echo "3️⃣  Creating ACME webroot..."
mkdir -p "$WEBROOT/.well-known/acme-challenge"
chown -R www-data:www-data "$WEBROOT"
chmod -R 755 "$WEBROOT"
echo "   ✅ Webroot created: $WEBROOT"
echo ""

# Step 4: Copy Nginx site configuration
echo "4️⃣  Installing Nginx site configuration..."
REPO_CONF="$(dirname "$0")/../deploy/seed/nginx-site.conf"
if [[ -f "$REPO_CONF" ]]; then
    cp "$REPO_CONF" "$NGINX_SITE"
    echo "   ✅ Site config copied to $NGINX_SITE"
else
    echo "   ⚠️  Warning: $REPO_CONF not found"
    echo "   Create it manually or copy from the repository"
fi
echo ""

# Step 5: Enable site
echo "5️⃣  Enabling site..."
if [[ ! -L "$NGINX_ENABLED" ]]; then
    ln -s "$NGINX_SITE" "$NGINX_ENABLED"
    echo "   ✅ Site enabled"
else
    echo "   ✅ Site already enabled"
fi

# Remove default site if present
if [[ -L "/etc/nginx/sites-enabled/default" ]]; then
    rm -f "/etc/nginx/sites-enabled/default"
    echo "   ✅ Default site removed"
fi
echo ""

# Step 6: Test Nginx configuration
echo "6️⃣  Testing Nginx configuration..."
if nginx -t; then
    echo "   ✅ Nginx configuration valid"
else
    echo "   ❌ Nginx configuration has errors"
    exit 1
fi
echo ""

# Step 7: Configure firewall
echo "7️⃣  Configuring firewall..."
if command -v ufw &> /dev/null; then
    if ufw status | grep -q "Status: active"; then
        ufw allow 80/tcp  comment "HTTP for ACME" || true
        ufw allow 443/tcp comment "HTTPS for seed.archivas.ai" || true
        echo "   ✅ Firewall rules added"
    else
        echo "   ⏭️  UFW not active, skipping firewall rules"
    fi
else
    echo "   ⏭️  UFW not installed, skipping firewall rules"
fi
echo ""

# Step 8: Reload Nginx
echo "8️⃣  Reloading Nginx..."
systemctl reload nginx
echo "   ✅ Nginx reloaded"
echo ""

# Step 9: Instructions for TLS certificate
echo "9️⃣  Next steps:"
echo ""
echo "   A. Point DNS A record: $DOMAIN → $(hostname -I | awk '{print $1}')"
echo "   B. Wait for DNS propagation (dig $DOMAIN)"
echo "   C. Obtain TLS certificate:"
echo ""
echo "      sudo certbot --nginx -d $DOMAIN --non-interactive --agree-tos -m admin@archivas.ai"
echo ""
echo "   D. Verify:"
echo ""
echo "      bash scripts/check-seed.sh"
echo ""
echo "   E. Setup auto-renewal:"
echo ""
echo "      sudo systemctl enable certbot.timer"
echo "      sudo systemctl start certbot.timer"
echo ""
echo "✅ Nginx setup complete!"
echo ""
echo "⚠️  Note: TLS certificate must be obtained manually (step C above)"
echo "    The site will serve HTTP on port 80 (redirect to HTTPS) until then."

