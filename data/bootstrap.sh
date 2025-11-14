#!/bin/bash
set -euo pipefail

# Bootstrap Seed2 node from Seed1 database
# This script syncs the blockchain data from Seed1 to Seed2 for fast initial sync
# Run as: sudo bash data/bootstrap.sh

SEED1_HOST="${SEED1_HOST:-seed.archivas.ai}"
SEED1_USER="${SEED1_USER:-ubuntu}"
SEED2_DATA_DIR="${SEED2_DATA_DIR:-/var/lib/archivas/seed2}"
SEED1_DATA_DIR="${SEED1_DATA_DIR:-/home/ubuntu/archivas/data}"

echo "=== Archivas Seed2 Bootstrap ==="
echo "This will sync blockchain data from Seed1 to Seed2"
echo ""
echo "Source: ${SEED1_USER}@${SEED1_HOST}:${SEED1_DATA_DIR}"
echo "Target: ${SEED2_DATA_DIR}"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

# Stop Seed2 node if running
echo "Stopping Seed2 node service..."
systemctl stop archivas-node-seed2 || true

# Create data directory
echo "Creating data directory..."
mkdir -p "${SEED2_DATA_DIR}"
chown -R root:root "${SEED2_DATA_DIR}"

# Check available disk space
echo "Checking disk space..."
AVAILABLE_GB=$(df -BG "${SEED2_DATA_DIR}" | awk 'NR==2 {print $4}' | sed 's/G//')
echo "Available disk space: ${AVAILABLE_GB}GB"
if [ "${AVAILABLE_GB}" -lt 50 ]; then
    echo "WARNING: Less than 50GB available. Blockchain data may be large."
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Get Seed1 database hash (for verification)
echo "Getting Seed1 database info..."
SEED1_HEIGHT=$(curl -sf "https://${SEED1_HOST}:8081/chainTip" | jq -r '.height' || echo "unknown")
echo "Seed1 current height: ${SEED1_HEIGHT}"

# Rsync database from Seed1
echo "Syncing database from Seed1 (this may take 30-60 minutes)..."
echo "Command: rsync -avz --progress ${SEED1_USER}@${SEED1_HOST}:${SEED1_DATA_DIR}/ ${SEED2_DATA_DIR}/"
echo ""

rsync -avz --progress \
    --exclude='*.log' \
    --exclude='*.lock' \
    "${SEED1_USER}@${SEED1_HOST}:${SEED1_DATA_DIR}/" \
    "${SEED2_DATA_DIR}/"

# Verify data directory
echo ""
echo "Verifying data directory..."
du -sh "${SEED2_DATA_DIR}"
ls -lh "${SEED2_DATA_DIR}/"

# Check database integrity
if [ -d "${SEED2_DATA_DIR}/badger" ]; then
    echo "BadgerDB directory found: ${SEED2_DATA_DIR}/badger"
    DB_SIZE=$(du -sh "${SEED2_DATA_DIR}/badger" | awk '{print $1}')
    echo "Database size: ${DB_SIZE}"
else
    echo "WARNING: BadgerDB directory not found!"
fi

# Set permissions
echo "Setting permissions..."
chown -R root:root "${SEED2_DATA_DIR}"
chmod -R 755 "${SEED2_DATA_DIR}"

echo ""
echo "=== Bootstrap complete ==="
echo "Database synced from Seed1 (height: ${SEED1_HEIGHT})"
echo ""
echo "Next steps:"
echo "1. Update /etc/archivas/seed2-node.env with checkpoint height and hash"
echo "2. Start node: sudo systemctl start archivas-node-seed2"
echo "3. Monitor logs: sudo journalctl -u archivas-node-seed2 -f"
echo "4. Verify sync: curl http://127.0.0.1:8082/chainTip"

