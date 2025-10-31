#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="/opt/archivas"
BIN="$INSTALL_DIR/bin/archivas-farmer"
OUT="$INSTALL_DIR/plots"
COUNT="${1:-2}"
SIZE="${2:-28}"

echo "🌾 Archivas Plot Generator"
echo "=========================="
echo "Creating $COUNT plots of size k=$SIZE"
echo ""

for i in $(seq 1 "$COUNT"); do
  DIR="$OUT/k${SIZE}-${i}"
  mkdir -p "$DIR"
  
  echo "📁 Generating plot $i/$COUNT (k=$SIZE) in $DIR..."
  "$BIN" plot --size "$SIZE" --path "$DIR"
  echo "  ✅ Plot $i complete"
done

echo ""
echo "✅ All plots generated!"
echo ""
echo "Update $INSTALL_DIR/plots.yaml to include new directories"
echo "Then reload farmer:"
echo "  sudo kill -HUP \$(pgrep -f 'archivas-farmer farm')"

