#!/bin/bash
set -euo pipefail

# ============================================================================
# Archivas Private Node Setup Script (Interactive)
# ============================================================================
# 
# Complete interactive setup for private farming node
#
# Usage: sudo bash deploy/private-node/setup-private-node-interactive.sh
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
echo "  🌾 Archivas Private Node Setup (Interactive)"
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

# ============================================================================
# Step 2: Build Binaries
# ============================================================================

step "Step 2: Building Archivas binaries"

cd "$PROJECT_ROOT"

info "Building archivas-node..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-node ./cmd/archivas-node" 2>&1 | grep -v "go: downloading" || true

info "Building archivas-farmer..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-farmer ./cmd/archivas-farmer" 2>&1 | grep -v "go: downloading" || true

info "Building archivas-wallet..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-wallet ./cmd/archivas-wallet" 2>&1 | grep -v "go: downloading" || true

if [ ! -f "archivas-node" ] || [ ! -f "archivas-farmer" ] || [ ! -f "archivas-wallet" ]; then
    error "Build failed"
    exit 1
fi

info "✓ Binaries built successfully"

# ============================================================================
# Step 3: Farmer Wallet
# ============================================================================

step "Step 3: Farmer Wallet Setup"

echo "Do you have an existing farmer private key?"
echo "  1) Yes, I have a private key"
echo "  2) No, generate a new wallet"
echo ""
echo -n "Enter choice [1/2]: "
read -r WALLET_CHOICE

FARMER_PRIVKEY=""

if [ "$WALLET_CHOICE" = "1" ]; then
    echo ""
    echo -n "Enter your 64-character hex private key: "
    read -r FARMER_PRIVKEY
    
    # Validate
    if ! [[ "$FARMER_PRIVKEY" =~ ^[0-9a-fA-F]{64}$ ]]; then
        error "Invalid key format (need 64 hex chars)"
        exit 1
    fi
    
    info "✓ Private key accepted"
    
elif [ "$WALLET_CHOICE" = "2" ]; then
    echo ""
    info "Generating new wallet..."
    
    # Generate wallet
    WALLET_OUTPUT=$(su - $REAL_USER -c "cd $PROJECT_ROOT && ./archivas-wallet new 2>&1")
    
    echo ""
    echo "════════════════════════════════════════════════════════════"
    echo "$WALLET_OUTPUT"
    echo "════════════════════════════════════════════════════════════"
    echo ""
    
    warn "⚠️  SAVE THIS INFORMATION SECURELY! ⚠️"
    echo ""
    echo "Your wallet has been generated. Please save:"
    echo "  - Mnemonic phrase (for recovery)"
    echo "  - Private key (for farming)"
    echo "  - Public address (to receive RCHV)"
    echo ""
    
    # Extract private key from output
    FARMER_PRIVKEY=$(echo "$WALLET_OUTPUT" | grep -oP 'priv:\s*\K[0-9a-fA-F]{64}' || echo "")
    
    if [ -z "$FARMER_PRIVKEY" ]; then
        error "Failed to extract private key"
        echo "Please copy it from above and press Enter to continue..."
        read -r
        echo -n "Enter the private key: "
        read -r FARMER_PRIVKEY
    else
        info "✓ Private key extracted automatically"
    fi
    
    echo ""
    echo -n "Press Enter after you've saved this information..."
    read -r
else
    error "Invalid choice"
    exit 1
fi

# ============================================================================
# Step 4: Plots Setup
# ============================================================================

step "Step 4: Plots Setup"

echo "Plots are required for farming. Do you have plots already?"
echo "  1) Yes, I have existing plots"
echo "  2) No, create plots now (recommended: 100G+)"
echo "  3) Skip for now (will need plots to farm)"
echo ""
echo -n "Enter choice [1/2/3]: "
read -r PLOTS_CHOICE

PLOTS_DIR=""

if [ "$PLOTS_CHOICE" = "1" ]; then
    echo ""
    echo -n "Enter full path to your plots directory: "
    read -r PLOTS_DIR
    PLOTS_DIR="${PLOTS_DIR/#\~/$USER_HOME}"
    
    if [ ! -d "$PLOTS_DIR" ]; then
        error "Directory not found: $PLOTS_DIR"
        exit 1
    fi
    
    PLOT_COUNT=$(find "$PLOTS_DIR" -name "*.arcv" -type f 2>/dev/null | wc -l)
    PLOT_SIZE=$(du -sh "$PLOTS_DIR" 2>/dev/null | cut -f1)
    
    info "Found $PLOT_COUNT plot(s), total size: $PLOT_SIZE"
    
elif [ "$PLOTS_CHOICE" = "2" ]; then
    echo ""
    echo -n "Enter path for new plots directory (will be created): "
    read -r PLOTS_DIR
    PLOTS_DIR="${PLOTS_DIR/#\~/$USER_HOME}"
    
    mkdir -p "$PLOTS_DIR"
    chown -R $REAL_USER:$REAL_USER "$PLOTS_DIR"
    
    echo ""
    echo "Plot size options:"
    echo "  - Minimum: 10G (testing)"
    echo "  - Recommended: 100G+"
    echo "  - Each 100G = ~1 k32 plot"
    echo ""
    echo -n "Enter plot size (e.g., 10G, 100G, 500G): "
    read -r PLOT_SIZE
    
    echo ""
    info "Creating plot (this may take a while)..."
    echo ""
    
    su - $REAL_USER -c "cd $PROJECT_ROOT && ./archivas-farmer plot create \
        --plots $PLOTS_DIR \
        --size $PLOT_SIZE \
        --farmer-pubkey \$(echo '$FARMER_PRIVKEY' | ./archivas-wallet pubkey-from-privkey 2>/dev/null || echo '')" || {
        
        # Fallback if pubkey extraction fails
        warn "Using private key directly for plotting..."
        su - $REAL_USER -c "cd $PROJECT_ROOT && ./archivas-farmer plot create \
            --plots $PLOTS_DIR \
            --size $PLOT_SIZE" || {
            error "Plot creation failed"
            exit 1
        }
    }
    
    info "✓ Plot created successfully"
    
elif [ "$PLOTS_CHOICE" = "3" ]; then
    warn "Skipping plots. Farmer will not be active until plots are added."
    PLOTS_DIR="$USER_HOME/archivas-plots"
    mkdir -p "$PLOTS_DIR"
    chown -R $REAL_USER:$REAL_USER "$PLOTS_DIR"
else
    error "Invalid choice"
    exit 1
fi

# ============================================================================
# Step 5: Node Configuration
# ============================================================================

step "Step 5: Node configuration"

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
# Step 6: Install Systemd Services
# ============================================================================

step "Step 6: Installing systemd services"

DATA_DIR="$PROJECT_ROOT/data"
mkdir -p "$DATA_DIR"
chown -R $REAL_USER:$REAL_USER "$DATA_DIR"

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
  -network-id archivas-devnet-v4 \\
  -db $DATA_DIR \\
  -rpc 127.0.0.1:8080 \\
  -p2p 0.0.0.0:9090 \\
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
MemoryMax=1G
MemoryHigh=512M
TasksMax=50

StandardOutput=journal
StandardError=journal
SyslogIdentifier=archivas-farmer-private

[Install]
WantedBy=multi-user.target
EOF

info "✓ Farmer service created"

# Reload systemd
systemctl daemon-reload

# ============================================================================
# Step 7: Start Services
# ============================================================================

step "Step 7: Starting services"

# Enable services
info "Enabling services..."
systemctl enable archivas-node-private
systemctl enable archivas-farmer-private

# Start node
info "Starting archivas-node-private..."
systemctl start archivas-node-private

sleep 3

if systemctl is-active --quiet archivas-node-private; then
    info "✓ Node started"
else
    error "Node failed to start"
    echo ""
    echo "Logs:"
    journalctl -u archivas-node-private -n 30 --no-pager
    exit 1
fi

# Start farmer
info "Starting archivas-farmer-private..."
systemctl start archivas-farmer-private

sleep 2

if systemctl is-active --quiet archivas-farmer-private; then
    info "✓ Farmer started"
else
    warn "Farmer may have issues. Checking logs..."
    journalctl -u archivas-farmer-private -n 10 --no-pager
fi

# ============================================================================
# Success
# ============================================================================

echo ""
echo "════════════════════════════════════════════════════════════"
echo "  ✅ Setup Complete!"
echo "════════════════════════════════════════════════════════════"
echo ""

echo "📊 Check Status:"
echo "   sudo systemctl status archivas-node-private"
echo "   sudo systemctl status archivas-farmer-private"
echo ""

echo "📜 View Logs:"
echo "   sudo journalctl -u archivas-node-private -f"
echo "   sudo journalctl -u archivas-farmer-private -f"
echo ""

echo "🔍 Check Node Sync:"
echo "   curl -s http://127.0.0.1:8080/chainTip | jq"
echo ""

echo "🌾 Your farmer is now running locally - no more 504 timeouts!"
echo ""

info "Happy farming! 🚜"
echo ""

