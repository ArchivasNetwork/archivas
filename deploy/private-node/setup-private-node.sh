#!/bin/bash
set -euo pipefail

# ============================================================================
# Archivas Private Node Setup Script
# ============================================================================
# 
# This script automates the setup of a private Archivas node and farmer
# for local farming. It will:
#   - Check prerequisites (Go, git)
#   - Build archivas-node and archivas-farmer
#   - Install systemd services
#   - Prompt for configuration (plots dir, farmer key)
#   - Start and enable services
#
# Usage: sudo bash deploy/private-node/setup-private-node.sh
#
# ============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
info() { echo -e "${GREEN}[INFO]${NC} $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; }
step() { echo -e "\n${BLUE}==>${NC} ${BLUE}$*${NC}\n"; }

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
   error "Please run as root (use sudo)"
   exit 1
fi

# Get the actual user (not root if using sudo)
REAL_USER="${SUDO_USER:-$(whoami)}"
USER_HOME=$(eval echo ~$REAL_USER)

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

info "Archivas Private Node Setup"
info "Project root: $PROJECT_ROOT"
info "Running as: root"
info "Target user: $REAL_USER"
info "User home: $USER_HOME"

# ============================================================================
# Step 1: Check Prerequisites
# ============================================================================

step "Step 1: Checking prerequisites..."

# Check for Go
if ! command -v go &> /dev/null; then
    error "Go is not installed!"
    echo ""
    echo "Please install Go 1.21 or later:"
    echo "  wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz"
    echo "  sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz"
    echo "  echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc"
    echo "  source ~/.bashrc"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
info "✓ Go installed: $GO_VERSION"

# Check for git
if ! command -v git &> /dev/null; then
    warn "Git not found. Installing..."
    apt-get update && apt-get install -y git
fi
info "✓ Git installed"

# Check for jq (useful for JSON parsing)
if ! command -v jq &> /dev/null; then
    warn "jq not found. Installing..."
    apt-get update && apt-get install -y jq
fi
info "✓ jq installed"

# ============================================================================
# Step 2: Build Archivas Binaries
# ============================================================================

step "Step 2: Building Archivas binaries..."

cd "$PROJECT_ROOT"

# Build as the real user
info "Building archivas-node..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-node ./cmd/archivas-node"

info "Building archivas-farmer..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-farmer ./cmd/archivas-farmer"

# Verify binaries
if [ ! -f "$PROJECT_ROOT/archivas-node" ]; then
    error "Failed to build archivas-node"
    exit 1
fi

if [ ! -f "$PROJECT_ROOT/archivas-farmer" ]; then
    error "Failed to build archivas-farmer"
    exit 1
fi

info "✓ Binaries built successfully"

# ============================================================================
# Step 3: Create Data Directory
# ============================================================================

step "Step 3: Creating data directory..."

DATA_DIR="$PROJECT_ROOT/data"
mkdir -p "$DATA_DIR"
chown -R $REAL_USER:$REAL_USER "$DATA_DIR"
info "✓ Data directory: $DATA_DIR"

# ============================================================================
# Step 4: Get Configuration from User
# ============================================================================

step "Step 4: Configuration"

# Plots directory
echo ""
echo -n "Enter the full path to your plots directory: "
read -r PLOTS_DIR

# Expand ~ if present
PLOTS_DIR="${PLOTS_DIR/#\~/$USER_HOME}"

if [ ! -d "$PLOTS_DIR" ]; then
    warn "Plots directory does not exist: $PLOTS_DIR"
    echo -n "Create it now? (y/n): "
    read -r CREATE_PLOTS
    if [[ "$CREATE_PLOTS" =~ ^[Yy]$ ]]; then
        mkdir -p "$PLOTS_DIR"
        chown -R $REAL_USER:$REAL_USER "$PLOTS_DIR"
        info "✓ Created plots directory"
    else
        error "Cannot proceed without plots directory"
        exit 1
    fi
fi

# Farmer private key
echo ""
echo "Enter your farmer private key (64 hex characters):"
echo "(Leave empty to generate a new one using archivas-wallet)"
read -r FARMER_PRIVKEY

if [ -z "$FARMER_PRIVKEY" ]; then
    info "No key provided. You can generate one using:"
    echo ""
    echo "  cd $PROJECT_ROOT"
    echo "  ./archivas-wallet new"
    echo ""
    echo "Then edit the systemd service file to add your key:"
    echo "  sudo nano /etc/systemd/system/archivas-farmer-private.service"
    echo ""
    FARMER_PRIVKEY="YOUR_PRIVATE_KEY_HERE"
else
    # Basic validation (64 hex chars)
    if ! [[ "$FARMER_PRIVKEY" =~ ^[0-9a-fA-F]{64}$ ]]; then
        warn "Private key format looks incorrect (expected 64 hex chars)"
        echo -n "Continue anyway? (y/n): "
        read -r CONTINUE
        if [[ ! "$CONTINUE" =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
fi

# Checkpoint (optional)
echo ""
echo "Fetch current checkpoint for faster sync? (recommended)"
echo -n "(y/n): "
read -r USE_CHECKPOINT

CHECKPOINT_HEIGHT=""
CHECKPOINT_HASH=""

if [[ "$USE_CHECKPOINT" =~ ^[Yy]$ ]]; then
    info "Fetching checkpoint from seed.archivas.ai..."
    
    CHECKPOINT_JSON=$(curl -sk https://seed.archivas.ai:8081/chainTip 2>/dev/null || echo "")
    
    if [ -n "$CHECKPOINT_JSON" ]; then
        CHECKPOINT_HEIGHT=$(echo "$CHECKPOINT_JSON" | jq -r '.height' 2>/dev/null || echo "")
        CHECKPOINT_HASH=$(echo "$CHECKPOINT_JSON" | jq -r '.hash' 2>/dev/null || echo "")
        
        if [ -n "$CHECKPOINT_HEIGHT" ] && [ -n "$CHECKPOINT_HASH" ]; then
            info "✓ Checkpoint: height=$CHECKPOINT_HEIGHT hash=${CHECKPOINT_HASH:0:16}..."
        else
            warn "Failed to parse checkpoint. Will sync from genesis."
            CHECKPOINT_HEIGHT=""
            CHECKPOINT_HASH=""
        fi
    else
        warn "Failed to fetch checkpoint. Will sync from genesis."
    fi
fi

# ============================================================================
# Step 5: Install Systemd Services
# ============================================================================

step "Step 5: Installing systemd services..."

# Node service
NODE_SERVICE="/etc/systemd/system/archivas-node-private.service"
info "Creating $NODE_SERVICE..."

cat > "$NODE_SERVICE" <<EOF
[Unit]
Description=Archivas Private Node - Full Node for Local Farming
Documentation=https://docs.archivas.ai/farmers/private-node
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=$REAL_USER
Group=$REAL_USER
WorkingDirectory=$PROJECT_ROOT

ExecStart=$PROJECT_ROOT/archivas-node \\
  --network-id archivas-devnet-v4 \\
  --db $DATA_DIR \\
  --rpc 127.0.0.1:8080 \\
  --p2p 0.0.0.0:9090 \\
  --genesis $PROJECT_ROOT/genesis/devnet.genesis.json \\
  --no-peer-discovery \\
  --peer-whitelist seed.archivas.ai:9090 \\
  --peer-whitelist seed2.archivas.ai:9090
EOF

# Add checkpoint if available
if [ -n "$CHECKPOINT_HEIGHT" ] && [ -n "$CHECKPOINT_HASH" ]; then
    cat >> "$NODE_SERVICE" <<EOF
 \\
  --checkpoint-height $CHECKPOINT_HEIGHT \\
  --checkpoint-hash $CHECKPOINT_HASH
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

Environment=GOMAXPROCS=2
Environment=GOGC=100

StandardOutput=journal
StandardError=journal
SyslogIdentifier=archivas-node-private

[Install]
WantedBy=multi-user.target
EOF

info "✓ Node service created"

# Farmer service
FARMER_SERVICE="/etc/systemd/system/archivas-farmer-private.service"
info "Creating $FARMER_SERVICE..."

cat > "$FARMER_SERVICE" <<EOF
[Unit]
Description=Archivas Private Farmer - Farming with Local Node
Documentation=https://docs.archivas.ai/farmers/private-node
After=archivas-node-private.service network-online.target
Wants=archivas-node-private.service network-online.target

[Service]
Type=simple
User=$REAL_USER
Group=$REAL_USER
WorkingDirectory=$PROJECT_ROOT

ExecStart=$PROJECT_ROOT/archivas-farmer farm \\
  --plots $PLOTS_DIR \\
  --node http://127.0.0.1:8080 \\
  --farmer-privkey $FARMER_PRIVKEY

Restart=always
RestartSec=10
StartLimitInterval=300
StartLimitBurst=5

LimitNOFILE=65536

StandardOutput=journal
StandardError=journal
SyslogIdentifier=archivas-farmer-private

[Install]
WantedBy=multi-user.target
EOF

info "✓ Farmer service created"

# Reload systemd
info "Reloading systemd daemon..."
systemctl daemon-reload

# ============================================================================
# Step 6: Enable and Start Services
# ============================================================================

step "Step 6: Starting services..."

# Enable services (start on boot)
info "Enabling services..."
systemctl enable archivas-node-private
systemctl enable archivas-farmer-private

# Start node first
info "Starting archivas-node-private..."
systemctl start archivas-node-private

# Wait a few seconds for node to initialize
sleep 3

# Check node status
if systemctl is-active --quiet archivas-node-private; then
    info "✓ Node started successfully"
else
    error "Node failed to start. Check logs:"
    echo "  sudo journalctl -u archivas-node-private -n 50"
    exit 1
fi

# Start farmer
info "Starting archivas-farmer-private..."
systemctl start archivas-farmer-private

# Wait a moment
sleep 2

# Check farmer status
if systemctl is-active --quiet archivas-farmer-private; then
    info "✓ Farmer started successfully"
else
    warn "Farmer may not have started. Check logs:"
    echo "  sudo journalctl -u archivas-farmer-private -n 50"
fi

# ============================================================================
# Step 7: Success & Next Steps
# ============================================================================

step "Setup Complete! 🎉"

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
info "Your private Archivas node and farmer are now running!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

echo "📊 Check Status:"
echo "  sudo systemctl status archivas-node-private"
echo "  sudo systemctl status archivas-farmer-private"
echo ""

echo "📜 View Logs:"
echo "  sudo journalctl -u archivas-node-private -f"
echo "  sudo journalctl -u archivas-farmer-private -f"
echo ""

echo "🔍 Check Node Sync:"
echo "  curl -s http://127.0.0.1:8080/chainTip | jq"
echo ""

echo "🛠️  Manage Services:"
echo "  sudo systemctl start|stop|restart archivas-node-private"
echo "  sudo systemctl start|stop|restart archivas-farmer-private"
echo ""

echo "📚 Documentation:"
echo "  https://docs.archivas.ai/farmers/private-node"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ "$FARMER_PRIVKEY" = "YOUR_PRIVATE_KEY_HERE" ]; then
    warn "IMPORTANT: You need to add your farmer private key!"
    echo ""
    echo "1. Generate a key:"
    echo "   cd $PROJECT_ROOT"
    echo "   sudo -u $REAL_USER ./archivas-wallet new"
    echo ""
    echo "2. Edit the farmer service:"
    echo "   sudo nano /etc/systemd/system/archivas-farmer-private.service"
    echo ""
    echo "3. Replace YOUR_PRIVATE_KEY_HERE with your actual key"
    echo ""
    echo "4. Restart the farmer:"
    echo "   sudo systemctl restart archivas-farmer-private"
    echo ""
fi

info "Happy farming! 🚜🌾"

