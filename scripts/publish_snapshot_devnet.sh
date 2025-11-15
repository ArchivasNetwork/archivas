#!/bin/bash
# publish_snapshot_devnet.sh
# Automated snapshot export and manifest publishing for Archivas devnet
# Run this on Seed2 via cron or systemd timer

set -e

# Configuration
NODE_RPC="${NODE_RPC:-http://127.0.0.1:8080}"
NODE_BINARY="${NODE_BINARY:-/home/ubuntu/archivas/archivas-node}"
DB_PATH="${DB_PATH:-/home/ubuntu/archivas/data}"
SNAPSHOT_DIR="${SNAPSHOT_DIR:-/srv/archivas-snapshots/devnet}"
SNAPSHOT_BASE_URL="${SNAPSHOT_BASE_URL:-https://snapshots.archivas.ai/devnet}"
NETWORK_ID="${NETWORK_ID:-archivas-devnet-v4}"
SAFE_BLOCKS_BEHIND="${SAFE_BLOCKS_BEHIND:-5000}"  # Export snapshot this many blocks behind tip

echo "═══════════════════════════════════════════════════════════"
echo "  📸 Archivas Devnet Snapshot Publisher"
echo "═══════════════════════════════════════════════════════════"
echo ""

# 1. Check prerequisites
if ! command -v jq &> /dev/null; then
    echo "Error: jq is not installed. Please install it first:"
    echo "  sudo apt-get install jq"
    exit 1
fi

if [ ! -f "$NODE_BINARY" ]; then
    echo "Error: Node binary not found at $NODE_BINARY"
    exit 1
fi

if [ ! -d "$DB_PATH" ]; then
    echo "Error: Database path not found: $DB_PATH"
    exit 1
fi

# 2. Fetch current chain tip
echo "[1/7] Fetching current chain tip from $NODE_RPC..."
CURRENT_TIP=$(curl -s "$NODE_RPC/chainTip" | jq -r '.height')

if [ -z "$CURRENT_TIP" ] || [ "$CURRENT_TIP" == "null" ]; then
    echo "Error: Failed to fetch chain tip from RPC"
    exit 1
fi

echo "  Current tip: $CURRENT_TIP"

# 3. Calculate safe snapshot height
SNAPSHOT_HEIGHT=$((CURRENT_TIP - SAFE_BLOCKS_BEHIND))

if [ $SNAPSHOT_HEIGHT -lt 0 ]; then
    echo "Error: Chain height too low for safe snapshot (tip: $CURRENT_TIP, safe margin: $SAFE_BLOCKS_BEHIND)"
    exit 1
fi

echo "[2/7] Safe snapshot height: $SNAPSHOT_HEIGHT (tip - $SAFE_BLOCKS_BEHIND)"

# 4. Create snapshot directory if needed
mkdir -p "$SNAPSHOT_DIR"

# 5. Export snapshot
SNAPSHOT_FILE="$SNAPSHOT_DIR/snap-$SNAPSHOT_HEIGHT.tar.gz"
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "[3/7] Exporting snapshot at height $SNAPSHOT_HEIGHT..."
echo "  Output: $SNAPSHOT_FILE"

if [ -f "$SNAPSHOT_FILE" ]; then
    echo "  ⚠️  Snapshot already exists, skipping export"
else
    "$NODE_BINARY" snapshot export \
        --height "$SNAPSHOT_HEIGHT" \
        --out "$SNAPSHOT_FILE" \
        --db "$DB_PATH" \
        --network-id "$NETWORK_ID" \
        --desc "Devnet snapshot exported at $TIMESTAMP"
    
    echo "  ✓ Export complete"
fi

# 6. Calculate SHA256 checksum
echo "[4/7] Calculating SHA256 checksum..."
CHECKSUM=$(sha256sum "$SNAPSHOT_FILE" | awk '{print $1}')
echo "  Checksum: $CHECKSUM"

# 7. Fetch block hash at snapshot height
echo "[5/7] Fetching block hash at height $SNAPSHOT_HEIGHT..."
BLOCK_HASH=$(curl -s "$NODE_RPC/block/$SNAPSHOT_HEIGHT" | jq -r '.hash')

if [ -z "$BLOCK_HASH" ] || [ "$BLOCK_HASH" == "null" ]; then
    echo "Error: Failed to fetch block hash for height $SNAPSHOT_HEIGHT"
    exit 1
fi

echo "  Block hash: $BLOCK_HASH"

# 8. Write manifest JSON
MANIFEST_FILE="$SNAPSHOT_DIR/latest.json"
SNAPSHOT_FILENAME=$(basename "$SNAPSHOT_FILE")

echo "[6/7] Writing manifest to $MANIFEST_FILE..."

cat > "$MANIFEST_FILE" <<EOF
{
  "network": "$NETWORK_ID",
  "height": $SNAPSHOT_HEIGHT,
  "hash": "$BLOCK_HASH",
  "snapshot_url": "$SNAPSHOT_BASE_URL/$SNAPSHOT_FILENAME",
  "checksum_sha256": "$CHECKSUM",
  "exported_at": "$TIMESTAMP",
  "current_tip": $CURRENT_TIP
}
EOF

echo "  ✓ Manifest written"

# 9. Clean up old snapshots (keep last 3)
echo "[7/7] Cleaning up old snapshots (keeping last 3)..."
cd "$SNAPSHOT_DIR"
ls -t snap-*.tar.gz 2>/dev/null | tail -n +4 | xargs -r rm -v

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  ✅ Snapshot published successfully"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "  📊 Snapshot Details:"
echo "     Height:   $SNAPSHOT_HEIGHT"
echo "     Hash:     ${BLOCK_HASH:0:16}..."
echo "     File:     $SNAPSHOT_FILENAME"
echo "     Checksum: ${CHECKSUM:0:16}..."
echo "     Size:     $(du -h "$SNAPSHOT_FILE" | cut -f1)"
echo ""
echo "  🌐 Manifest URL:"
echo "     $SNAPSHOT_BASE_URL/latest.json"
echo ""
echo "  📥 Farmers can bootstrap with:"
echo "     archivas-node bootstrap --network devnet --db /var/lib/archivas"
echo ""
echo "═══════════════════════════════════════════════════════════"

