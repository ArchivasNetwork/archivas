#!/bin/bash
# Nginx + SSL Setup for Archivas Betanet Seeds
# This script sets up nginx as a reverse proxy with Let's Encrypt SSL
# Run this on each seed: Seed1, Seed2, Seed3

set -e

echo "╔════════════════════════════════════════════════════╗"
echo "║   Archivas Betanet - Nginx + SSL Setup            ║"
echo "╚════════════════════════════════════════════════════╝"
echo ""

# Check if running as root or with sudo
if [[ $EUID -ne 0 ]]; then
   echo "❌ This script must be run as root or with sudo"
   exit 1
fi

# Get the domain name from user
read -p "Enter the domain for this seed (e.g., seed1.betanet.archivas.ai): " DOMAIN
read -p "Enter your email for Let's Encrypt notifications: " EMAIL

if [ -z "$DOMAIN" ] || [ -z "$EMAIL" ]; then
    echo "❌ Domain and email are required!"
    exit 1
fi

echo ""
echo "📋 Configuration:"
echo "  Domain: $DOMAIN"
echo "  Email: $EMAIL"
echo "  Backend: http://127.0.0.1:8545"
echo ""

# Step 1: Update system and install dependencies
echo "📦 Installing nginx and certbot..."
apt-get update -qq
apt-get install -y nginx certbot python3-certbot-nginx ufw

echo "✅ Dependencies installed"
echo ""

# Step 2: Configure firewall
echo "🔥 Configuring firewall..."
ufw allow 'Nginx Full'
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable
echo "✅ Firewall configured"
echo ""

# Step 3: Create nginx configuration
echo "⚙️  Creating nginx configuration..."
cat > /etc/nginx/sites-available/archivas-betanet <<EOF
# Archivas Betanet RPC Proxy
# Domain: $DOMAIN
# Backend: http://127.0.0.1:8545

upstream archivas_backend {
    server 127.0.0.1:8545;
    keepalive 64;
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name $DOMAIN;

    # Allow Let's Encrypt validation
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    # Redirect all other traffic to HTTPS
    location / {
        return 301 https://\$server_name\$request_uri;
    }
}

# HTTPS Server
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name $DOMAIN;

    # SSL certificates (will be added by certbot)
    # ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # CORS headers for MetaMask
    add_header Access-Control-Allow-Origin "*" always;
    add_header Access-Control-Allow-Methods "GET, POST, OPTIONS" always;
    add_header Access-Control-Allow-Headers "Content-Type, Authorization" always;
    add_header Access-Control-Max-Age "86400" always;

    # Handle preflight requests
    if (\$request_method = 'OPTIONS') {
        return 204;
    }

    # Logging
    access_log /var/log/nginx/archivas-betanet-access.log;
    error_log /var/log/nginx/archivas-betanet-error.log;

    # Proxy to archivas-node
    location / {
        proxy_pass http://archivas_backend;
        proxy_http_version 1.1;
        
        # Proxy headers
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Connection "";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Buffering
        proxy_buffering off;
        proxy_request_buffering off;
    }

    # Health check endpoint
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
EOF

echo "✅ Nginx configuration created"
echo ""

# Step 4: Enable the site
echo "🔗 Enabling nginx site..."
ln -sf /etc/nginx/sites-available/archivas-betanet /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
echo "✅ Site enabled"
echo ""

# Step 5: Test nginx configuration
echo "🧪 Testing nginx configuration..."
nginx -t
echo "✅ Nginx configuration valid"
echo ""

# Step 6: Restart nginx
echo "🔄 Restarting nginx..."
systemctl restart nginx
systemctl enable nginx
echo "✅ Nginx restarted"
echo ""

# Step 7: Obtain SSL certificate
echo "🔐 Obtaining SSL certificate from Let's Encrypt..."
echo "    This may take a minute..."
certbot --nginx -d $DOMAIN --non-interactive --agree-tos --email $EMAIL --redirect

if [ $? -eq 0 ]; then
    echo "✅ SSL certificate obtained and installed"
else
    echo "❌ Failed to obtain SSL certificate"
    echo "   Make sure DNS is pointing to this server!"
    exit 1
fi
echo ""

# Step 8: Set up auto-renewal
echo "⏰ Setting up automatic SSL renewal..."
systemctl enable certbot.timer
systemctl start certbot.timer
echo "✅ Auto-renewal configured"
echo ""

# Step 9: Test the setup
echo "🧪 Testing HTTPS endpoint..."
sleep 2
curl -sk https://$DOMAIN/health || echo "⚠️  Health check failed (this may be okay if archivas-node is still starting)"
echo ""

# Step 10: Final verification
echo "📊 Final Status Check:"
echo ""
systemctl status nginx --no-pager -l | head -10
echo ""

echo "╔════════════════════════════════════════════════════╗"
echo "║   ✅ Setup Complete!                              ║"
echo "╚════════════════════════════════════════════════════╝"
echo ""
echo "🌐 Your RPC endpoint is now:"
echo "   https://$DOMAIN"
echo ""
echo "🧪 Test with:"
echo "   curl -s https://$DOMAIN -X POST -H 'Content-Type: application/json' \\"
echo "     -d '{\"jsonrpc\":\"2.0\",\"method\":\"eth_chainId\",\"params\":[],\"id\":1}'"
echo ""
echo "🦊 MetaMask Settings:"
echo "   RPC URL: https://$DOMAIN"
echo "   Chain ID: 1644"
echo "   Currency: RCHV"
echo ""

