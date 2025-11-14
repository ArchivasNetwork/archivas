#!/bin/bash
set -euo pipefail

# ============================================================================
# Archivas Private Node Setup - Node Only
# ============================================================================
# 
# Sets up a private Archivas node (no farming)
# Use this if you want to run your own node for RPC access
#
# Usage: sudo bash deploy/private-node/setup-node-only.sh
#
# ============================================================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Logging
info() { echo -e "${GREEN}[INFO]${NC} $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; }
step() { echo -e "\n${CYAN}═══${NC} ${BLUE}$*${NC}\n"; }

# Check root
if [ "$EUID" -ne 0 ]; then 
   error "Please run as root (use sudo)"
   exit 1
fi

# Get real user
REAL_USER="${SUDO_USER:-$(whoami)}"
USER_HOME=$(eval echo ~$REAL_USER)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo ""
echo "════════════════════════════════════════════════════════════"
echo "  🌐 Archivas Private Node Setup"
echo "════════════════════════════════════════════════════════════"
echo ""
info "User: $REAL_USER"
info "Home: $USER_HOME"
info "Project: $PROJECT_ROOT"
echo ""

# ============================================================================
# Step 1: Prerequisites
# ============================================================================

step "Step 1: Checking prerequisites"

# Go
if ! command -v go &> /dev/null; then
    error "Go not installed!"
    echo "Install Go 1.21+:"
    echo "  wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz"
    echo "  sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz"
    echo "  export PATH=\$PATH:/usr/local/go/bin"
    exit 1
fi
GO_VERSION=$(go version | awk '{print $3}')
info "✓ Go: $GO_VERSION"

# Git
if ! command -v git &> /dev/null; then
    info "Installing git..."
    apt-get update && apt-get install -y git
fi
info "✓ Git installed"

# jq
if ! command -v jq &> /dev/null; then
    info "Installing jq..."
    apt-get update && apt-get install -y jq
fi
info "✓ jq installed"

# curl
if ! command -v curl &> /dev/null; then
    info "Installing curl..."
    apt-get update && apt-get install -y curl
fi
info "✓ curl installed"

# ============================================================================
# Step 2: Build Node Binary
# ============================================================================

step "Step 2: Building archivas-node"

cd "$PROJECT_ROOT"

info "Building archivas-node..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-node ./cmd/archivas-node" 2>&1 | grep -v "go: downloading" || true

if [ ! -f "archivas-node" ]; then
    error "Build failed"
    exit 1
fi

info "✓ Binary built successfully"

# ============================================================================
# Step 3: Configuration
# ============================================================================

step "Step 3: Node configuration"

# RPC binding
echo "RPC Configuration:"
echo "  1) Bind to localhost only (127.0.0.1) - Most secure"
echo "  2) Bind to all interfaces (0.0.0.0) - Accessible from network"
echo ""
echo -n "Enter choice [1/2] (default: 1): "
read -r RPC_CHOICE

RPC_BIND="127.0.0.1:8080"
if [ "$RPC_CHOICE" = "2" ]; then
    RPC_BIND="0.0.0.0:8080"
    warn "RPC will be accessible from the network!"
    warn "Make sure to configure firewall appropriately."
fi

info "RPC will bind to: $RPC_BIND"

# P2P configuration
echo ""
echo "P2P Configuration:"
echo "Recommended: 0.0.0.0:9090 (allows connections from seeds)"
echo ""
echo -n "P2P bind address [default: 0.0.0.0:9090]: "
read -r P2P_BIND

if [ -z "$P2P_BIND" ]; then
    P2P_BIND="0.0.0.0:9090"
fi

info "P2P will bind to: $P2P_BIND"

# Checkpoint
echo ""
echo "Fetch checkpoint for fast sync? (recommended)"
echo -n "(y/n): "
read -r USE_CHECKPOINT

CHECKPOINT_HEIGHT=""
CHECKPOINT_HASH=""

if [[ "$USE_CHECKPOINT" =~ ^[Yy]$ ]]; then
    info "Fetching checkpoint from seed.archivas.ai..."
    
    CHECKPOINT_JSON=$(curl -sk --max-time 10 https://seed.archivas.ai/chainTip 2>/dev/null || echo "")
    
    if [ -n "$CHECKPOINT_JSON" ]; then
        CHECKPOINT_HEIGHT=$(echo "$CHECKPOINT_JSON" | jq -r '.height' 2>/dev/null || echo "")
        CHECKPOINT_HASH=$(echo "$CHECKPOINT_JSON" | jq -r '.hash' 2>/dev/null || echo "")
        
        if [ -n "$CHECKPOINT_HEIGHT" ] && [ "$CHECKPOINT_HEIGHT" != "null" ]; then
            info "✓ Checkpoint: height=$CHECKPOINT_HEIGHT"
        else
            warn "Failed to fetch checkpoint. Will sync from genesis."
            CHECKPOINT_HEIGHT=""
            CHECKPOINT_HASH=""
        fi
    else
        warn "Seed unreachable. Will sync from genesis."
    fi
fi

# ============================================================================
# Step 4: Install Systemd Service
# ============================================================================

step "Step 4: Installing systemd service"

DATA_DIR="$PROJECT_ROOT/data"
mkdir -p "$DATA_DIR"
chown -R $REAL_USER:$REAL_USER "$DATA_DIR"

NODE_SERVICE="/etc/systemd/system/archivas-node-private.service"
info "Creating $NODE_SERVICE..."

cat > "$NODE_SERVICE" <<EOF
[Unit]
Description=Archivas Private Node - Full Node for Local/Network Use
Documentation=https://docs.archivas.ai/farmers/private-node
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$REAL_USER
Group=$REAL_USER
WorkingDirectory=$PROJECT_ROOT

ExecStart=$PROJECT_ROOT/archivas-node \\
  -network-id archivas-devnet-v4 \\
  -db $DATA_DIR \\
  -rpc $RPC_BIND \\
  -p2p $P2P_BIND \\
  -genesis $PROJECT_ROOT/genesis/devnet.genesis.json \\
  -no-peer-discovery \\
  -peer seed.archivas.ai:9090 \\
  -peer seed2.archivas.ai:9090
EOF

# Add checkpoint if available
if [ -n "$CHECKPOINT_HEIGHT" ] && [ -n "$CHECKPOINT_HASH" ]; then
    cat >> "$NODE_SERVICE" <<EOF
 \\
  -checkpoint-height $CHECKPOINT_HEIGHT \\
  -checkpoint-hash $CHECKPOINT_HASH
EOF
else
    echo "" >> "$NODE_SERVICE"
fi

cat >> "$NODE_SERVICE" <<EOF

Restart=always
RestartSec=10
StartLimitInterval=300
StartLimitBurst=5

LimitNOFILE=65536
MemoryMax=4G
MemoryHigh=3G
TasksMax=100

Environment=GOMAXPROCS=2
Environment=GOGC=100

StandardOutput=journal
StandardError=journal
SyslogIdentifier=archivas-node-private

[Install]
WantedBy=multi-user.target
EOF

info "✓ Service file created"

# Reload systemd
systemctl daemon-reload

# ============================================================================
# Step 5: Start Service
# ============================================================================

step "Step 5: Starting node"

info "Enabling service..."
systemctl enable archivas-node-private

info "Starting archivas-node-private..."
systemctl start archivas-node-private

sleep 3

if systemctl is-active --quiet archivas-node-private; then
    info "✓ Node started successfully"
else
    error "Node failed to start"
    echo ""
    echo "Logs:"
    journalctl -u archivas-node-private -n 30 --no-pager
    exit 1
fi

# ============================================================================
# Success
# ============================================================================

echo ""
echo "════════════════════════════════════════════════════════════"
echo "  ✅ Private Node Setup Complete!"
echo "════════════════════════════════════════════════════════════"
echo ""

echo "📊 Node Status:"
echo "   sudo systemctl status archivas-node-private"
echo ""

echo "📜 View Logs:"
echo "   sudo journalctl -u archivas-node-private -f"
echo ""

echo "🔍 Check Sync Progress:"
echo "   curl -s http://${RPC_BIND}/chainTip | jq"
echo ""

echo "🌐 RPC Endpoint:"
echo "   http://${RPC_BIND}"
echo ""

if [ "$RPC_BIND" = "0.0.0.0:8080" ]; then
    LOCAL_IP=$(hostname -I | awk '{print $1}')
    echo "   Accessible from network at: http://${LOCAL_IP}:8080"
    echo ""
    warn "⚠️  Remember to configure firewall if needed!"
fi

echo "📚 Next Steps:"
echo "   - Wait for node to sync (check with curl command above)"
echo "   - To setup a farmer, run:"
echo "     sudo bash deploy/private-node/setup-farmer-only.sh"
echo ""

info "Node is running! 🚀"
echo ""

