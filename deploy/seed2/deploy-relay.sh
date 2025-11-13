#!/bin/bash
# Deploy Seed2 Relay RPC - seed2.archivas.ai
# This script sets up Nginx as a reverse proxy to Seed1

set -e

echo "🚀 Deploying Seed2 Relay RPC..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "❌ Please run as root (or with sudo)"
    exit 1
fi

# Variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NGINX_CONF_SRC="$SCRIPT_DIR/nginx-seed2-relay.conf"
NGINX_CONF_DEST="/etc/nginx/sites-available/archivas-seed2"
NGINX_CONF_LINK="/etc/nginx/sites-enabled/archivas-seed2"
CACHE_DIR="/var/cache/nginx/archivas"
SSL_DIR="/etc/nginx/ssl"
LOG_DIR="/var/log/nginx"

echo "📦 Installing Nginx..."
apt-get update -qq
apt-get install -y nginx certbot python3-certbot-nginx

echo "📁 Creating directories..."
mkdir -p "$CACHE_DIR"
mkdir -p "$SSL_DIR"
mkdir -p "$LOG_DIR"

echo "🔐 Generating self-signed SSL certificate (temporary)..."
if [ ! -f "$SSL_DIR/seed2.archivas.ai.crt" ]; then
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout "$SSL_DIR/seed2.archivas.ai.key" \
        -out "$SSL_DIR/seed2.archivas.ai.crt" \
        -subj "/CN=seed2.archivas.ai/O=Archivas/C=US"
    chmod 600 "$SSL_DIR/seed2.archivas.ai.key"
    chmod 644 "$SSL_DIR/seed2.archivas.ai.crt"
    echo "✅ Self-signed certificate created"
    echo "⚠️  Run 'sudo certbot --nginx -d seed2.archivas.ai' to get a real certificate"
else
    echo "✅ SSL certificate already exists"
fi

echo "📝 Copying Nginx configuration..."
cp "$NGINX_CONF_SRC" "$NGINX_CONF_DEST"

echo "🔗 Enabling site..."
ln -sf "$NGINX_CONF_DEST" "$NGINX_CONF_LINK"

echo "🧹 Removing default site..."
rm -f /etc/nginx/sites-enabled/default

echo "✅ Testing Nginx configuration..."
nginx -t

echo "🔄 Reloading Nginx..."
systemctl reload nginx
systemctl enable nginx

echo ""
echo "✅ Seed2 Relay RPC deployed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Point DNS: seed2.archivas.ai → $(curl -s ifconfig.me 2>/dev/null || echo 'THIS_SERVER_IP')"
echo "2. Get real SSL cert: sudo certbot --nginx -d seed2.archivas.ai"
echo "3. Test the relay:"
echo "   curl -I https://seed2.archivas.ai/chainTip"
echo "   curl -I https://seed2.archivas.ai/relay/status"
echo ""
echo "🔍 Check logs:"
echo "   sudo tail -f /var/log/nginx/seed2-access.log"
echo "   sudo tail -f /var/log/nginx/seed2-error.log"
echo ""
echo "🎯 Cache stats:"
echo "   sudo du -sh /var/cache/nginx/archivas"
echo ""

