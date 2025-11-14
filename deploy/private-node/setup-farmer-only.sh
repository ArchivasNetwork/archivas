#!/bin/bash
set -euo pipefail

# ============================================================================
# Archivas Farmer Setup - Farmer Only
# ============================================================================
# 
# Sets up an Archivas farmer (connects to any node)
# Can connect to local private node or public seed
#
# Usage: sudo bash deploy/private-node/setup-farmer-only.sh
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
echo "  🌾 Archivas Farmer Setup"
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

# ============================================================================
# Step 2: Build Binaries
# ============================================================================

step "Step 2: Building binaries"

cd "$PROJECT_ROOT"

info "Building archivas-farmer..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-farmer ./cmd/archivas-farmer" 2>&1 | grep -v "go: downloading" || true

info "Building archivas-wallet..."
su - $REAL_USER -c "cd $PROJECT_ROOT && go build -o archivas-wallet ./cmd/archivas-wallet" 2>&1 | grep -v "go: downloading" || true

if [ ! -f "archivas-farmer" ] || [ ! -f "archivas-wallet" ]; then
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
FARMER_PUBKEY=""

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
    
    # Derive public key
    FARMER_PUBKEY=$(su - $REAL_USER -c "cd $PROJECT_ROOT && echo '$FARMER_PRIVKEY' | ./archivas-wallet pubkey-from-privkey 2>/dev/null || echo ''")
    
    if [ -z "$FARMER_PUBKEY" ]; then
        warn "Could not derive public key automatically"
        echo -n "Enter your public key (66 hex chars, optional): "
        read -r FARMER_PUBKEY
    else
        info "✓ Derived public key: ${FARMER_PUBKEY:0:20}..."
    fi
    
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
    echo "  - Public key (for plotting)"
    echo "  - Address (to receive RCHV)"
    echo ""
    
    # Extract keys
    FARMER_PRIVKEY=$(echo "$WALLET_OUTPUT" | grep -i "private" | grep -oP '[0-9a-fA-F]{64}' | head -1 || echo "")
    FARMER_PUBKEY=$(echo "$WALLET_OUTPUT" | grep -i "public" | grep -oP '[0-9a-fA-F]{66}' | head -1 || echo "")
    
    if [ -z "$FARMER_PRIVKEY" ]; then
        error "Failed to extract private key"
        echo -n "Enter the private key: "
        read -r FARMER_PRIVKEY
    fi
    
    if [ -z "$FARMER_PUBKEY" ]; then
        warn "Failed to extract public key"
        echo -n "Enter the public key: "
        read -r FARMER_PUBKEY
    fi
    
    info "✓ Keys extracted"
    
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

echo "Plot Configuration:"
echo "  1) I have existing plots"
echo "  2) Create new plots now"
echo "  3) Skip plots (setup later)"
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
    echo "  - 10G   = Quick test (~5 minutes)"
    echo "  - 100G  = Decent farming (~30-60 minutes)"
    echo "  - 500G+ = Serious farming (~2-4 hours)"
    echo ""
    echo -n "Enter plot size (e.g., 10G, 100G, 500G): "
    read -r PLOT_SIZE
    
    echo ""
    info "Creating plot (this may take a while)..."
    info "Plot will be created in: $PLOTS_DIR"
    echo ""
    
    # Create plot with proper arguments
    PLOT_CMD="cd $PROJECT_ROOT && ./archivas-farmer plot create --plots $PLOTS_DIR --size $PLOT_SIZE"
    
    if [ -n "$FARMER_PUBKEY" ]; then
        PLOT_CMD="$PLOT_CMD --farmer-pubkey $FARMER_PUBKEY"
    fi
    
    su - $REAL_USER -c "$PLOT_CMD" || {
        error "Plot creation failed"
        warn "You can create plots later with:"
        echo "  cd $PROJECT_ROOT"
        echo "  ./archivas-farmer plot create --plots $PLOTS_DIR --size $PLOT_SIZE --farmer-pubkey $FARMER_PUBKEY"
        echo ""
        echo -n "Continue anyway? (y/n): "
        read -r CONTINUE
        if [[ ! "$CONTINUE" =~ ^[Yy]$ ]]; then
            exit 1
        fi
    }
    
    info "✓ Plot creation complete"
    
elif [ "$PLOTS_CHOICE" = "3" ]; then
    warn "Skipping plots. Farmer will not be active until plots are added."
    PLOTS_DIR="$USER_HOME/archivas-plots"
    mkdir -p "$PLOTS_DIR"
    chown -R $REAL_USER:$REAL_USER "$PLOTS_DIR"
    
    info "Empty plots directory created: $PLOTS_DIR"
else
    error "Invalid choice"
    exit 1
fi

# ============================================================================
# Step 5: Node Configuration
# ============================================================================

step "Step 5: Node connection"

echo "Which node should the farmer connect to?"
echo "  1) Local private node (http://127.0.0.1:8080)"
echo "  2) Public seed (https://seed.archivas.ai)"
echo "  3) Custom URL"
echo ""
echo -n "Enter choice [1/2/3]: "
read -r NODE_CHOICE

NODE_URL=""

if [ "$NODE_CHOICE" = "1" ]; then
    NODE_URL="http://127.0.0.1:8080"
    
    # Check if local node is running
    if ! systemctl is-active --quiet archivas-node-private 2>/dev/null; then
        warn "Local private node service not found or not running"
        echo ""
        echo "To setup a local node first, run:"
        echo "  sudo bash deploy/private-node/setup-node-only.sh"
        echo ""
        echo -n "Continue anyway? (y/n): "
        read -r CONTINUE
        if [[ ! "$CONTINUE" =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        info "✓ Local node detected"
    fi
    
elif [ "$NODE_CHOICE" = "2" ]; then
    NODE_URL="https://seed.archivas.ai"
    warn "Using public seed - you may experience higher latency"
    
elif [ "$NODE_CHOICE" = "3" ]; then
    echo ""
    echo -n "Enter node URL (e.g., http://192.168.1.100:8080): "
    read -r NODE_URL
    
    if [ -z "$NODE_URL" ]; then
        error "URL cannot be empty"
        exit 1
    fi
else
    error "Invalid choice"
    exit 1
fi

info "Farmer will connect to: $NODE_URL"

# ============================================================================
# Step 6: Install Systemd Service
# ============================================================================

step "Step 6: Installing systemd service"

FARMER_SERVICE="/etc/systemd/system/archivas-farmer-private.service"
info "Creating $FARMER_SERVICE..."

cat > "$FARMER_SERVICE" <<EOF
[Unit]
Description=Archivas Farmer - Farming Proof of Space
Documentation=https://docs.archivas.ai/farmers/
After=network-online.target
Wants=network-online.target
EOF

# If using local node, add dependency
if [ "$NODE_URL" = "http://127.0.0.1:8080" ]; then
    cat >> "$FARMER_SERVICE" <<EOF
After=archivas-node-private.service
Wants=archivas-node-private.service
EOF
fi

cat >> "$FARMER_SERVICE" <<EOF

[Service]
Type=simple
User=$REAL_USER
Group=$REAL_USER
WorkingDirectory=$PROJECT_ROOT

ExecStart=$PROJECT_ROOT/archivas-farmer farm \\
  --plots $PLOTS_DIR \\
  --node $NODE_URL \\
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

info "✓ Service file created"

# Reload systemd
systemctl daemon-reload

# ============================================================================
# Step 7: Start Service
# ============================================================================

step "Step 7: Starting farmer"

info "Enabling service..."
systemctl enable archivas-farmer-private

info "Starting archivas-farmer-private..."
systemctl start archivas-farmer-private

sleep 2

if systemctl is-active --quiet archivas-farmer-private; then
    info "✓ Farmer started successfully"
else
    warn "Farmer may have issues. Checking logs..."
    journalctl -u archivas-farmer-private -n 20 --no-pager
    echo ""
    warn "Check if plots exist and node is accessible"
fi

# ============================================================================
# Success
# ============================================================================

echo ""
echo "════════════════════════════════════════════════════════════"
echo "  ✅ Farmer Setup Complete!"
echo "════════════════════════════════════════════════════════════"
echo ""

echo "📊 Farmer Status:"
echo "   sudo systemctl status archivas-farmer-private"
echo ""

echo "📜 View Logs:"
echo "   sudo journalctl -u archivas-farmer-private -f"
echo ""

echo "🔍 Check Plots:"
echo "   ls -lh $PLOTS_DIR/"
echo ""

echo "🌐 Node Connection:"
echo "   $NODE_URL"
echo ""

echo "📚 Useful Commands:"
echo "   Stop farmer:   sudo systemctl stop archivas-farmer-private"
echo "   Restart farmer: sudo systemctl restart archivas-farmer-private"
echo "   View metrics:  curl -s http://localhost:9102/metrics"
echo ""

info "Happy farming! 🚜🌾"
echo ""

