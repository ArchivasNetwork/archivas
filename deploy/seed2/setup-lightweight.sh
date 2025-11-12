#!/usr/bin/env bash
# Lightweight setup for seed2 - prevents server crashes during sync
# This version adds resource limits and swap space BEFORE starting IBD

set -euo pipefail

USER="${USER:-ubuntu}"
APP_DIR="/home/$USER/archivas"
BIN_NODE="$APP_DIR/archivas-node"
DATA_DIR="$APP_DIR/data"
LOG_DIR="/var/log/archivas"
DOMAIN="seed2.archivas.ai"
NETWORK="archivas-devnet-v4"

echo "🌾 Archivas Seed2 Setup (Lightweight - Crash-Resistant)"
echo "======================================================="
echo "Domain: $DOMAIN"
echo "Network: $NETWORK"
echo ""

# 1. Check system resources
echo "1️⃣  Checking system resources..."
TOTAL_MEM=$(free -m | awk '/^Mem:/{print $2}')
TOTAL_SWAP=$(free -m | awk '/^Swap:/{print $2}')
echo "   Total Memory: ${TOTAL_MEM}MB"
echo "   Total Swap: ${TOTAL_SWAP}MB"

if [ "$TOTAL_MEM" -lt 4096 ]; then
    echo "   ⚠️  Low memory system detected (< 4GB)"
    echo "   Will add swap space for stability"
fi
echo ""

# 2. Add swap space if needed (critical for preventing OOM)
echo "2️⃣  Configuring swap space..."
if [ "$TOTAL_SWAP" -lt 4096 ]; then
    echo "   Creating 4GB swap file..."
    sudo fallocate -l 4G /swapfile 2>/dev/null || sudo dd if=/dev/zero of=/swapfile bs=1M count=4096
    sudo chmod 600 /swapfile
    sudo mkswap /swapfile
    sudo swapon /swapfile
    # Make permanent
    if ! grep -q '/swapfile' /etc/fstab; then
        echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
    fi
    echo "   ✅ Swap space added (4GB)"
else
    echo "   ✅ Swap space already configured"
fi
echo ""

# 3. Install dependencies
echo "3️⃣  Installing dependencies..."
sudo apt-get update -y -qq
sudo apt-get install -y -qq git golang-go nginx jq ufw certbot python3-certbot-nginx
echo "   ✅ Dependencies installed"
echo ""

# 4. Clone repository
echo "4️⃣  Setting up repository..."
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

# 5. Build binaries
echo "5️⃣  Building binaries..."
CGO_ENABLED=0 go build -a -o "$BIN_NODE" ./cmd/archivas-node
echo "   ✅ Node binary built"
echo ""

# 6. Create directories
echo "6️⃣  Creating directories..."
mkdir -p "$DATA_DIR"
sudo mkdir -p "$LOG_DIR"
sudo chown $USER:$USER "$LOG_DIR"
sudo mkdir -p /var/www/seed2
sudo chown $USER:$USER /var/www/seed2
echo "   ✅ Directories created"
echo ""

# 7. Create systemd service with AGGRESSIVE resource limits
echo "7️⃣  Creating systemd service (with resource limits)..."
sudo tee /etc/systemd/system/archivas-node-seed2.service >/dev/null <<EOF
[Unit]
Description=Archivas Node (seed2 - Lightweight)
Documentation=https://github.com/ArchivasNetwork/archivas
After=network.target

[Service]
User=$USER
WorkingDirectory=$APP_DIR
ExecStart=$BIN_NODE \\
  --rpc 127.0.0.1:8080 \\
  --p2p 0.0.0.0:9090 \\
  --peer seed.archivas.ai:9090 \\
  --genesis $APP_DIR/genesis/devnet.genesis.json \\
  --network-id $NETWORK \\
  --db $DATA_DIR

Restart=always
RestartSec=5
StandardOutput=append:$LOG_DIR/archivas-node-seed2.log
StandardError=append:$LOG_DIR/archivas-node-seed2.log
LimitNOFILE=65535

# AGGRESSIVE Resource Limits (prevents server crash)
# Memory: 3GB hard limit, 2GB soft (lower for smaller servers)
MemoryMax=3G
MemoryHigh=2G
# CPU: Lower priority to keep system responsive
CPUWeight=80
Nice=10
# Task limit (prevents goroutine explosion)
TasksMax=2048
# Max processes
LimitNPROC=512

# Restart strategy
StartLimitInterval=300
StartLimitBurst=5

# Environment - MORE AGGRESSIVE GC for syncing
Environment="GOMAXPROCS=2"
Environment="GOGC=30"

# OOM behavior
OOMPolicy=continue
OOMScoreAdjust=500

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable archivas-node-seed2
echo "   ✅ Systemd service created with resource limits:"
echo "      - Memory: 3GB max, 2GB soft"
echo "      - Tasks: 2048 max"
echo "      - GOGC: 30 (very aggressive GC)"
echo ""

# 8. Create Nginx configuration (same as before)
echo "8️⃣  Creating Nginx configuration..."
sudo tee /etc/nginx/snippets/archivas-cors.conf >/dev/null <<'EOF'
add_header Access-Control-Allow-Origin "*" always;
add_header Access-Control-Allow-Methods "GET, POST, OPTIONS" always;
add_header Access-Control-Allow-Headers "Content-Type" always;
add_header Access-Control-Max-Age "86400" always;

if ($request_method = OPTIONS) {
  return 204;
}
EOF

sudo tee /etc/nginx/conf.d/archivas-ratelimit.conf >/dev/null <<'EOF'
# Rate limiting zones for seed2
limit_req_zone $binary_remote_addr zone=api_per_ip:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=challenge:10m rate=120r/m;
limit_req_zone $binary_remote_addr zone=submitblock:10m rate=60r/m;
limit_req_zone $binary_remote_addr zone=blocks:10m rate=10r/m;
limit_req_zone $binary_remote_addr zone=health:10m rate=60r/m;
limit_conn_zone $binary_remote_addr zone=concurrent:10m;
EOF

sudo tee /etc/nginx/sites-available/$DOMAIN >/dev/null <<EOF
# Nginx configuration for $DOMAIN
server {
  listen 80;
  listen [::]:80;
  server_name $DOMAIN;

  location /.well-known/acme-challenge/ {
    root /var/www/seed2;
  }

  add_header X-Content-Type-Options nosniff always;
  add_header X-Frame-Options SAMEORIGIN always;
  add_header Referrer-Policy no-referrer-when-downgrade always;

  gzip on;
  gzip_vary on;
  gzip_comp_level 6;
  gzip_types application/json text/plain application/javascript text/css;

  proxy_connect_timeout 60s;
  proxy_send_timeout 180s;
  proxy_read_timeout 180s;
  proxy_buffering on;

  client_max_body_size 1m;

  location = /metrics { return 404; }
  location /metrics/ { return 404; }

  location / {
    limit_req zone=api_per_ip burst=20 nodelay;
    limit_conn concurrent 50;
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host \$host;
    proxy_set_header X-Real-IP \$remote_addr;
    include /etc/nginx/snippets/archivas-cors.conf;
  }

  access_log /var/log/nginx/$DOMAIN.access.log;
  error_log /var/log/nginx/$DOMAIN.error.log;
}
EOF

sudo ln -sf /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/$DOMAIN
sudo nginx -t
sudo systemctl reload nginx
echo "   ✅ Nginx configured"
echo ""

# 9. Start node
echo "9️⃣  Starting node..."
sudo systemctl start archivas-node-seed2
sleep 5

if systemctl is-active --quiet archivas-node-seed2; then
  echo "   ✅ Node started successfully"
else
  echo "   ⚠️  Node failed to start. Check logs:"
  echo "   sudo journalctl -u archivas-node-seed2 -n 50"
  exit 1
fi
echo ""

# 10. Get TLS certificate
echo "🔟 Obtaining TLS certificate..."
sudo certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos -m admin@archivas.ai --redirect
echo "   ✅ TLS certificate obtained"
echo ""

# 11. Configure firewall
echo "1️⃣1️⃣  Configuring firewall..."
sudo ufw allow 80/tcp comment "HTTP"
sudo ufw allow 443/tcp comment "HTTPS"
sudo ufw allow 9090/tcp comment "P2P"
sudo ufw --force enable
echo "   ✅ Firewall configured"
echo ""

# 12. Final verification
echo "1️⃣2️⃣  Final checks..."
if curl -s http://127.0.0.1:8080/ping > /dev/null; then
  echo "   ✅ Node is responding"
else
  echo "   ⚠️  Node may still be starting up"
fi
echo ""

echo "✅ Lightweight setup complete!"
echo ""
echo "📊 Configuration:"
echo "   • Memory limit: 3GB (prevents OOM crashes)"
echo "   • Swap: 4GB (allows system to handle memory pressure)"
echo "   • GOGC: 30 (very aggressive GC during sync)"
echo "   • Task limit: 2048 (prevents goroutine explosion)"
echo ""
echo "🔍 Monitor sync progress:"
echo "   watch -n 10 'curl -s http://127.0.0.1:8080/chainTip | jq -r .height'"
echo ""
echo "📝 Check logs:"
echo "   sudo journalctl -u archivas-node-seed2 -f"
echo ""
echo "⚠️  IMPORTANT:"
echo "   • Sync will be slower but STABLE"
echo "   • Server will remain responsive"
echo "   • Estimated sync time: 4-6 hours (slower but safe)"

