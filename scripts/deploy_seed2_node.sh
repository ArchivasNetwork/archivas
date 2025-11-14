#!/bin/bash
set -euo pipefail

# Seed2 Full Node Deployment Script
# Purpose: Set up Seed2 as a full P2P node + RPC relay
# Run as: sudo bash scripts/deploy_seed2_node.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== Archivas Seed2 Full Node Deployment ==="
echo "Project root: $PROJECT_ROOT"
echo ""

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

success() { echo -e "${GREEN}✓${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    error "Please run as root (use sudo)"
    exit 1
fi

# Step 1: Install dependencies
echo "Step 1: Installing dependencies..."
apt-get update -qq
apt-get install -y rsync curl jq golang-go git nginx certbot python3-certbot-nginx ufw > /dev/null 2>&1
success "Dependencies installed"

# Step 2: Check if archivas binary exists
echo ""
echo "Step 2: Checking Archivas node binary..."
if [ ! -f "$PROJECT_ROOT/archivas-node" ]; then
    warn "archivas-node binary not found. Building from source..."
    cd "$PROJECT_ROOT"
    go build -o archivas-node ./cmd/archivas-node || {
        error "Failed to build archivas-node"
        exit 1
    }
    success "Built archivas-node binary"
else
    success "archivas-node binary found"
fi

# Step 3: Create data directory
echo ""
echo "Step 3: Creating data directory..."
mkdir -p /var/lib/archivas/seed2
chown -R root:root /var/lib/archivas/seed2
chmod -R 755 /var/lib/archivas/seed2
success "Data directory created: /var/lib/archivas/seed2"

# Step 4: Get checkpoint from Seed1
echo ""
echo "Step 4: Getting checkpoint from Seed1..."
CHECKPOINT_HEIGHT=$(curl -sf --max-time 10 "https://seed.archivas.ai:8081/chainTip" | jq -r '.height' || echo "")
CHECKPOINT_HASH=$(curl -sf --max-time 10 "https://seed.archivas.ai:8081/chainTip" | jq -r '.hash' || echo "")

if [ -z "$CHECKPOINT_HEIGHT" ] || [ -z "$CHECKPOINT_HASH" ]; then
    warn "Could not fetch checkpoint from Seed1. Using placeholder values."
    CHECKPOINT_HEIGHT="0"
    CHECKPOINT_HASH="genesis"
else
    success "Checkpoint fetched: height=$CHECKPOINT_HEIGHT"
fi

# Step 5: Create environment file
echo ""
echo "Step 5: Creating environment file..."
mkdir -p /etc/archivas
cat > /etc/archivas/seed2-node.env << EOF
# Archivas Seed2 Full Node Configuration
# Auto-generated on $(date)

# Checkpoint configuration
CHECKPOINT_HEIGHT=${CHECKPOINT_HEIGHT}
CHECKPOINT_HASH=${CHECKPOINT_HASH}

# Primary peer (Seed1)
SEED1_P2P=seed.archivas.ai:30303

# Optional: Additional validator peers
# VALIDATOR_PEERS=validator02.archivas.ai:30303,validator03.archivas.ai:30303
EOF
success "Environment file created: /etc/archivas/seed2-node.env"

# Step 6: Install systemd unit
echo ""
echo "Step 6: Installing systemd unit..."
cp "$PROJECT_ROOT/services/node-seed2/archivas-node-seed2.service" \
   /etc/systemd/system/archivas-node-seed2.service
systemctl daemon-reload
systemctl enable archivas-node-seed2 > /dev/null 2>&1
success "Systemd unit installed and enabled"

# Step 7: Configure firewall
echo ""
echo "Step 7: Configuring firewall..."
if command -v ufw > /dev/null 2>&1; then
    # Check if UFW is active
    if ufw status | grep -q "Status: active"; then
        # P2P port
        ufw allow 30303/tcp comment 'Archivas P2P TCP' > /dev/null 2>&1 || true
        ufw allow 30303/udp comment 'Archivas P2P UDP' > /dev/null 2>&1 || true
        success "Firewall rules added for P2P port 30303"
    else
        warn "UFW is installed but not active. Skipping firewall configuration."
        echo "Run 'sudo bash $PROJECT_ROOT/infra/firewall.seed2.sh' to configure manually."
    fi
else
    warn "UFW not found. Firewall not configured."
fi

# Step 8: Bootstrap option
echo ""
echo "Step 8: Bootstrap database from Seed1?"
echo "Options:"
echo "  1) Bootstrap now (rsync from Seed1 - recommended, 30-60 min)"
echo "  2) Skip bootstrap (node will sync from genesis/checkpoint - slower)"
echo ""
read -p "Choose option (1 or 2): " BOOTSTRAP_CHOICE

if [ "$BOOTSTRAP_CHOICE" = "1" ]; then
    echo ""
    warn "Bootstrap requires SSH access to Seed1"
    echo "Default: root@seed.archivas.ai:/root/archivas/data"
    read -p "Press Enter to continue or Ctrl+C to abort..."
    
    bash "$PROJECT_ROOT/data/bootstrap.sh" || {
        error "Bootstrap failed. You can run it manually later:"
        echo "  sudo bash $PROJECT_ROOT/data/bootstrap.sh"
    }
else
    warn "Skipping bootstrap. Node will sync from scratch."
fi

# Step 9: Start the node
echo ""
echo "Step 9: Starting Seed2 node..."
systemctl start archivas-node-seed2

# Wait a few seconds for startup
sleep 3

if systemctl is-active --quiet archivas-node-seed2; then
    success "Seed2 node started successfully!"
else
    error "Seed2 node failed to start. Check logs:"
    echo "  sudo journalctl -u archivas-node-seed2 -n 50 --no-pager"
    exit 1
fi

# Step 10: Health checks
echo ""
echo "Step 10: Running health checks..."

# Check if P2P port is listening
if ss -tulpn | grep -q ":30303"; then
    success "P2P port 30303 is listening"
else
    warn "P2P port 30303 not listening yet (may take a moment)"
fi

# Check if RPC port is listening
if ss -tulpn | grep -q "127.0.0.1:8082"; then
    success "RPC port 8082 is listening"
    
    # Try to query chain tip
    sleep 2
    HEIGHT=$(curl -sf --max-time 5 http://127.0.0.1:8082/chainTip 2>/dev/null | jq -r '.height' || echo "")
    if [ -n "$HEIGHT" ]; then
        success "RPC responding: height=$HEIGHT"
    else
        warn "RPC not responding yet (node may still be initializing)"
    fi
else
    warn "RPC port 8082 not listening yet (may take a moment)"
fi

# Step 11: Summary
echo ""
echo "========================================"
echo "  Seed2 Full Node Deployment Complete"
echo "========================================"
echo ""
echo "✓ Node binary: $PROJECT_ROOT/archivas-node"
echo "✓ Data directory: /var/lib/archivas/seed2"
echo "✓ Systemd service: archivas-node-seed2"
echo "✓ P2P port: 30303 (TCP/UDP)"
echo "✓ RPC port: 8082 (localhost only)"
echo ""
echo "Next steps:"
echo ""
echo "1. Monitor logs:"
echo "   sudo journalctl -u archivas-node-seed2 -f"
echo ""
echo "2. Check sync status:"
echo "   curl -s http://127.0.0.1:8082/chainTip | jq"
echo ""
echo "3. Verify P2P connectivity (from another machine):"
echo "   telnet seed2.archivas.ai 30303"
echo ""
echo "4. Update Nginx relay to include Seed2 node as fallback:"
echo "   # Edit /etc/nginx/sites-available/archivas-seed2"
echo "   # Add upstream seed2_node { server 127.0.0.1:8082; }"
echo "   sudo nginx -t && sudo systemctl reload nginx"
echo ""
echo "5. Update farmer documentation:"
echo "   Farmers can now peer with: seed2.archivas.ai:30303"
echo ""
echo "Monitoring:"
echo "  - Status: sudo systemctl status archivas-node-seed2"
echo "  - Logs: sudo journalctl -u archivas-node-seed2 -f"
echo "  - Height: curl -s http://127.0.0.1:8082/chainTip | jq .height"
echo "  - Peers: sudo journalctl -u archivas-node-seed2 | grep -i peer | tail"
echo ""
echo "Troubleshooting guide: $PROJECT_ROOT/docs/seed2-node.md"
echo ""

