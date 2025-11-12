#!/usr/bin/env bash
# setup-seed2.sh - Deploy Archivas seed2 node (seed2.archivas.ai)
# Based on seed.archivas.ai configuration with simplified setup

set -euo pipefail

USER="${USER:-ubuntu}"
APP_DIR="/home/$USER/archivas"
BIN_NODE="$APP_DIR/archivas-node"
DATA_DIR="$APP_DIR/data"
LOG_DIR="/var/log/archivas"
DOMAIN="seed2.archivas.ai"
NETWORK="archivas-devnet-v4"

echo "🌾 Archivas Seed2 Setup"
echo "======================"
echo "Domain: $DOMAIN"
echo "Network: $NETWORK"
echo ""

# Install dependencies
echo "1️⃣  Installing dependencies..."
sudo apt-get update -y
sudo apt-get install -y git golang-go nginx jq ufw certbot python3-certbot-nginx
echo "   ✅ Dependencies installed"
echo ""

# Clone or update repository
echo "2️⃣  Setting up repository..."
if [ ! -d "$APP_DIR/.git" ]; then
  echo "   Cloning repository..."
  git clone https://github.com/ArchivasNetwork/archivas.git "$APP_DIR"
else
  echo "   Repository exists, updating..."
  cd "$APP_DIR"
  git fetch --all
  git reset --hard origin/main
fi
cd "$APP_DIR"
echo "   ✅ Repository ready"
echo ""

# Build binaries
echo "3️⃣  Building binaries..."
CGO_ENABLED=0 go build -a -o "$BIN_NODE" ./cmd/archivas-node
echo "   ✅ Node binary built"
echo ""

# Create directories
echo "4️⃣  Creating directories..."
mkdir -p "$DATA_DIR"
sudo mkdir -p "$LOG_DIR"
sudo chown $USER:$USER "$LOG_DIR"
sudo mkdir -p /var/www/seed2
sudo chown $USER:$USER /var/www/seed2
echo "   ✅ Directories created"
echo ""

# Create systemd service
echo "5️⃣  Creating systemd service..."
sudo tee /etc/systemd/system/archivas-node-seed2.service >/dev/null <<EOF
[Unit]
Description=Archivas Node (seed2)
Documentation=https://github.com/ArchivasNetwork/archivas
After=network.target

[Service]
User=$USER
WorkingDirectory=$APP_DIR
ExecStart=$BIN_NODE \\
  --rpc 127.0.0.1:8080 \\
  --p2p 0.0.0.0:9090 \\
  --genesis $APP_DIR/genesis/devnet.genesis.json \\
  --network-id $NETWORK \\
  --db $DATA_DIR

Restart=always
RestartSec=3
StandardOutput=append:$LOG_DIR/archivas-node-seed2.log
StandardError=append:$LOG_DIR/archivas-node-seed2.log
LimitNOFILE=65535

# Environment
Environment="GOMAXPROCS=4"

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable archivas-node-seed2
echo "   ✅ Systemd service created"
echo ""

# Create Nginx CORS snippet
echo "6️⃣  Creating Nginx CORS configuration..."
sudo tee /etc/nginx/snippets/archivas-cors.conf >/dev/null <<'EOF'
add_header Access-Control-Allow-Origin "*" always;
add_header Access-Control-Allow-Methods "GET, POST, OPTIONS" always;
add_header Access-Control-Allow-Headers "Content-Type" always;
add_header Access-Control-Max-Age "86400" always;

if ($request_method = OPTIONS) {
  return 204;
}
EOF
echo "   ✅ CORS configuration created"
echo ""

# Create rate limiting configuration
echo "7️⃣  Creating rate limiting configuration..."
sudo tee /etc/nginx/conf.d/archivas-ratelimit.conf >/dev/null <<'EOF'
# Rate limiting zones for seed2
limit_req_zone $binary_remote_addr zone=api_per_ip:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=challenge:10m rate=120r/m;  # 2 req/sec
limit_req_zone $binary_remote_addr zone=submitblock:10m rate=60r/m;  # 1 req/sec
limit_req_zone $binary_remote_addr zone=blocks:10m rate=10r/m;      # Heavy endpoint
limit_req_zone $binary_remote_addr zone=health:10m rate=60r/m;      # Health checks

# Connection limiting
limit_conn_zone $binary_remote_addr zone=concurrent:10m;
EOF
echo "   ✅ Rate limiting configuration created"
echo ""

# Create Nginx site configuration
echo "8️⃣  Creating Nginx site configuration..."
sudo tee /etc/nginx/sites-available/$DOMAIN >/dev/null <<EOF
# Nginx configuration for $DOMAIN
# Reverse proxy to Archivas node RPC (localhost:8080)

# HTTP server - redirect to HTTPS + ACME challenge
server {
  listen 80;
  listen [::]:80;
  server_name $DOMAIN;

  # ACME challenge for Let's Encrypt
  location /.well-known/acme-challenge/ {
    root /var/www/seed2;
  }

  # Redirect all other traffic to HTTPS
  location / {
    return 301 https://\$host\$request_uri;
  }
}

# HTTPS server - reverse proxy to node RPC
server {
  listen 443 ssl http2;
  listen [::]:443 ssl http2;
  server_name $DOMAIN;

  # TLS certificates (managed by certbot)
  # These will be added by certbot
  # ssl_certificate     /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
  # ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;
  ssl_protocols TLSv1.2 TLSv1.3;
  ssl_prefer_server_ciphers on;
  ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';

  # Security headers
  add_header X-Content-Type-Options nosniff always;
  add_header X-Frame-Options SAMEORIGIN always;
  add_header Referrer-Policy no-referrer-when-downgrade always;
  add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

  # Compression
  gzip on;
  gzip_vary on;
  gzip_comp_level 6;
  gzip_types application/json text/plain application/javascript text/css;
  gzip_min_length 1000;

  # Timeouts
  proxy_connect_timeout 60s;
  proxy_send_timeout    180s;
  proxy_read_timeout    180s;
  proxy_buffering on;
  proxy_buffer_size 4k;
  proxy_buffers 8 4k;

  # Disable large uploads
  client_max_body_size 1m;

  # Block internal metrics from public access
  location = /metrics {
    return 404;
  }

  location /metrics/ {
    return 404;
  }

  # Health check endpoints (lenient rate limiting)
  location = /healthz {
    limit_req zone=health burst=10 nodelay;
    proxy_pass http://127.0.0.1:8080/healthz;
    proxy_set_header Host              \$host;
    proxy_set_header X-Real-IP         \$remote_addr;
    proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  location = /health {
    limit_req zone=health burst=10 nodelay;
    proxy_pass http://127.0.0.1:8080/health;
    proxy_set_header Host              \$host;
    proxy_set_header X-Real-IP         \$remote_addr;
    proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  # ChainTip endpoint (no rate limiting - backend caching handles load)
  location = /chainTip {
    proxy_pass http://127.0.0.1:8080/chainTip;
    proxy_set_header Host              \$host;
    proxy_set_header X-Real-IP         \$remote_addr;
    proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;
    proxy_connect_timeout 10s;
    proxy_send_timeout 15s;
    proxy_read_timeout 15s;
    proxy_next_upstream off;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  # Challenge endpoint (farmers check frequently)
  location = /challenge {
    limit_req zone=challenge burst=20 nodelay;
    limit_conn concurrent 20;
    proxy_pass http://127.0.0.1:8080/challenge;
    proxy_set_header Host              \$host;
    proxy_set_header X-Real-IP         \$remote_addr;
    proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;
    proxy_connect_timeout 10s;
    proxy_send_timeout 15s;
    proxy_read_timeout 15s;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  # SubmitBlock endpoint (farmers submit blocks)
  location = /submitBlock {
    limit_req zone=submitblock burst=10 nodelay;
    limit_conn concurrent 10;
    proxy_pass http://127.0.0.1:8080/submitBlock;
    proxy_set_header Host              \$host;
    proxy_set_header X-Real-IP         \$remote_addr;
    proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;
    proxy_connect_timeout 10s;
    proxy_send_timeout 30s;
    proxy_read_timeout 30s;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  # Blocks range endpoint (heavy, needs rate limiting)
  location = /blocks/range {
    limit_req zone=blocks burst=2 nodelay;
    limit_conn concurrent 5;
    proxy_pass http://127.0.0.1:8080/blocks/range;
    proxy_set_header Host              \$host;
    proxy_set_header X-Real-IP         \$remote_addr;
    proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;
    proxy_connect_timeout 60s;
    proxy_send_timeout 180s;
    proxy_read_timeout 180s;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  # General API endpoints (catch-all)
  location / {
    limit_req zone=api_per_ip burst=20 nodelay;
    limit_conn concurrent 50;
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host              \$host;
    proxy_set_header X-Real-IP         \$remote_addr;
    proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto \$scheme;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  # Access logs
  access_log /var/log/nginx/$DOMAIN.access.log;
  error_log /var/log/nginx/$DOMAIN.error.log;
}
EOF

# Enable site
sudo ln -sf /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/$DOMAIN
echo "   ✅ Nginx configuration created"
echo ""

# Test Nginx configuration
echo "9️⃣  Testing Nginx configuration..."
sudo nginx -t
echo "   ✅ Nginx configuration valid"
echo ""

# Start node (before getting TLS cert)
echo "🔟 Starting node..."
sudo systemctl restart archivas-node-seed2
sleep 5

# Check if node is running
if systemctl is-active --quiet archivas-node-seed2; then
  echo "   ✅ Node started successfully"
else
  echo "   ⚠️  Node failed to start. Check logs:"
  echo "   sudo journalctl -u archivas-node-seed2 -n 50"
  exit 1
fi
echo ""

# Reload Nginx
echo "1️⃣1️⃣  Reloading Nginx..."
sudo systemctl reload nginx
echo "   ✅ Nginx reloaded"
echo ""

# Get TLS certificate
echo "1️⃣2️⃣  Obtaining TLS certificate..."
sudo certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos -m admin@archivas.ai --redirect
echo "   ✅ TLS certificate obtained"
echo ""

# Configure firewall
echo "1️⃣3️⃣  Configuring firewall..."
sudo ufw allow 80/tcp comment "HTTP"
sudo ufw allow 443/tcp comment "HTTPS"
sudo ufw allow 9090/tcp comment "P2P"
sudo ufw --force enable
echo "   ✅ Firewall configured"
echo ""

# Final verification
echo "1️⃣4️⃣  Verifying setup..."
sleep 3

echo "   Checking node health..."
if curl -s http://127.0.0.1:8080/healthz > /dev/null; then
  echo "   ✅ Node is healthy"
else
  echo "   ⚠️  Node health check failed"
fi

echo "   Checking HTTPS endpoint..."
if curl -s "https://$DOMAIN/healthz" > /dev/null; then
  echo "   ✅ HTTPS endpoint is working"
else
  echo "   ⚠️  HTTPS endpoint check failed (DNS might not be configured yet)"
fi

echo ""
echo "✅ Seed2 deployment complete!"
echo ""
echo "Next steps:"
echo "1. Configure DNS: $DOMAIN → $(curl -s ifconfig.me || echo 'YOUR_SERVER_IP')"
echo "2. Verify: curl https://$DOMAIN/chainTip | jq"
echo "3. Check logs: sudo journalctl -u archivas-node-seed2 -f"
echo "4. Update Explorer and SDK to use $DOMAIN as secondary RPC"
echo ""
echo "Node status:"
systemctl status archivas-node-seed2 --no-pager | head -10

