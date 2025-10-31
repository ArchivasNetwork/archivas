#!/usr/bin/env bash
set -euo pipefail

echo "🚀 Archivas Auto-Join Installer"
echo "================================"

# Configuration
INSTALL_DIR="/opt/archivas"
BIN_DIR="$INSTALL_DIR/bin"
LOG_DIR="$INSTALL_DIR/logs"
PLOTS_DIR="$INSTALL_DIR/plots"
ENV_FILE="$INSTALL_DIR/archivas.env"
PLOTS_YAML="$INSTALL_DIR/plots.yaml"

# CPU pinning defaults
TL_CORES="${TL_CORES:-0-3}"
FARMER_CORES="${FARMER_CORES:-4-15}"

# System packages
PKGS="curl wget ca-certificates coreutils jq util-linux taskset"

echo "📦 Installing dependencies..."
sudo mkdir -p "$BIN_DIR" "$LOG_DIR" "$PLOTS_DIR" "$INSTALL_DIR/autojoin"
sudo apt-get update -qq
sudo apt-get install -y $PKGS

echo "📥 Downloading binaries..."
# For now, use local builds (in production, download from releases)
if [ -f "./archivas-timelord" ]; then
  sudo cp archivas-timelord "$BIN_DIR/"
  sudo chmod +x "$BIN_DIR/archivas-timelord"
  echo "  ✅ archivas-timelord"
fi

if [ -f "./archivas-farmer" ]; then
  sudo cp archivas-farmer "$BIN_DIR/"
  sudo chmod +x "$BIN_DIR/archivas-farmer"
  echo "  ✅ archivas-farmer"
fi

if [ -f "./archivas" ]; then
  sudo cp archivas "$BIN_DIR/"
  sudo chmod +x "$BIN_DIR/archivas"
  echo "  ✅ archivas (wallet CLI)"
fi

echo "⚙️  Creating configuration files..."
# Create env file
if [ ! -f "$ENV_FILE" ]; then
  cat > /tmp/archivas.env.tmp <<'EOF'
# Archivas Auto-Join Configuration
ARCHIVAS_WALLET_ADDR=arcv1REPLACE_WITH_YOUR_ADDRESS
ARCHIVAS_FARMER_PRIVKEY=REPLACE_WITH_YOUR_PRIVATE_KEY
ARCHIVAS_NODE_RPC=http://57.129.148.132:8080
ARCHIVAS_RPC_TOKEN=

# Timelord Performance
TL_ITER_PER_TICK=3000
TL_TICK_MS=40

# Farmer Performance  
FARMER_POLL_MS=200
FARMER_CHUNK_BYTES=8388608
FARMER_HEALTH_INTERVAL=5s
EOF
  sudo mv /tmp/archivas.env.tmp "$ENV_FILE"
  sudo chmod 600 "$ENV_FILE"
  echo "  ✅ Created $ENV_FILE (edit before starting!)"
else
  echo "  ℹ️  $ENV_FILE exists (not overwriting)"
fi

# Create plots.yaml
if [ ! -f "$PLOTS_YAML" ]; then
  cat > /tmp/plots.yaml.tmp <<EOF
plots:
  - $PLOTS_DIR/k28-1
  - $PLOTS_DIR/k28-2
EOF
  sudo mv /tmp/plots.yaml.tmp "$PLOTS_YAML"
  echo "  ✅ Created $PLOTS_YAML"
else
  echo "  ℹ️  $PLOTS_YAML exists (not overwriting)"
fi

echo "🔧 Installing systemd units..."
# Timelord service
cat > /tmp/archivas-timelord.service.tmp <<EOF
[Unit]
Description=Archivas Timelord
After=network-online.target
Wants=network-online.target

[Service]
User=root
EnvironmentFile=$ENV_FILE
WorkingDirectory=$INSTALL_DIR
CPUAffinity=$TL_CORES
Environment=GOGC=100
ExecStart=$BIN_DIR/archivas-timelord --node \${ARCHIVAS_NODE_RPC}
Restart=always
RestartSec=2
StandardOutput=append:$LOG_DIR/timelord.log
StandardError=append:$LOG_DIR/timelord.err
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF
sudo mv /tmp/archivas-timelord.service.tmp /etc/systemd/system/archivas-timelord.service

# Farmer template service
cat > /tmp/archivas-farmer@.service.tmp <<'EOFSERVICE'
[Unit]
Description=Archivas Farmer %i
After=network-online.target
Wants=network-online.target

[Service]
User=root
EnvironmentFile=/opt/archivas/archivas.env
WorkingDirectory=/opt/archivas
CPUAffinity=4-15
Environment=GOGC=100
ExecStart=/opt/archivas/bin/archivas-farmer farm \
  --plots /opt/archivas/plots/k28-1,/opt/archivas/plots/k28-2 \
  --farmer-privkey ${ARCHIVAS_FARMER_PRIVKEY} \
  --node ${ARCHIVAS_NODE_RPC}
Restart=always
RestartSec=2
StandardOutput=append:/opt/archivas/logs/farmer-%i.log
StandardError=append:/opt/archivas/logs/farmer-%i.err
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOFSERVICE
sudo mv /tmp/archivas-farmer@.service.tmp /etc/systemd/system/archivas-farmer@.service

sudo systemctl daemon-reload
echo "  ✅ Systemd units installed"

echo ""
echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "1. Edit configuration:"
echo "   sudo nano $ENV_FILE"
echo ""
echo "2. (Optional) Generate plots:"
echo "   sudo bash $INSTALL_DIR/autojoin/create-plots.sh 2 28"
echo ""
echo "3. Start services:"
echo "   sudo systemctl enable --now archivas-timelord"
echo "   sudo systemctl enable --now archivas-farmer@1"
echo ""
echo "4. Monitor:"
echo "   tail -f $LOG_DIR/timelord.log"
echo "   tail -f $LOG_DIR/farmer-1.log"
echo ""
echo "🌾 Archivas Auto-Join Ready!"

