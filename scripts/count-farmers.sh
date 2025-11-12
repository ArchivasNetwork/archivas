#!/bin/bash

# Count unique farmers from blocks over a range
# Usage: ./count-farmers.sh [num_blocks]
# Example: ./count-farmers.sh 10000  (checks last 10,000 blocks)

RPC_URL="${RPC_URL:-http://127.0.0.1:8080}"

# Get current chain height
CURRENT_HEIGHT=$(curl -s "$RPC_URL/healthz" | jq -r '.height')
echo "Current chain height: $CURRENT_HEIGHT"
echo ""

# Default: check last 10,000 blocks
NUM_BLOCKS=${1:-10000}
START_HEIGHT=$((CURRENT_HEIGHT - NUM_BLOCKS + 1))

if [ $START_HEIGHT -lt 0 ]; then
    START_HEIGHT=0
fi

echo "Checking last $NUM_BLOCKS blocks (from height $START_HEIGHT to $CURRENT_HEIGHT)"
echo ""

# Fetch blocks in batches of 1000 (max limit for /blocks/since)
TOTAL_FARMERS=""
BATCH_SIZE=1000
HEIGHT=$START_HEIGHT

while [ $HEIGHT -le $CURRENT_HEIGHT ]; do
    BATCH_END=$((HEIGHT + BATCH_SIZE - 1))
    if [ $BATCH_END -gt $CURRENT_HEIGHT ]; then
        BATCH_END=$CURRENT_HEIGHT
        BATCH_SIZE=$((CURRENT_HEIGHT - HEIGHT + 1))
    fi
    
    if [ $BATCH_SIZE -le 0 ]; then
        break
    fi
    
    echo -n "Fetching blocks $HEIGHT-$BATCH_END ($((BATCH_END - HEIGHT + 1)) blocks)... "
    
    # Use /blocks/since endpoint to get blocks from a specific height
    RESPONSE=$(curl -s "$RPC_URL/blocks/since/$HEIGHT?limit=$BATCH_SIZE" 2>/dev/null)
    FARMERS=$(echo "$RESPONSE" | jq -r '.blocks[].farmer' 2>/dev/null)
    
    if [ $? -eq 0 ] && [ -n "$FARMERS" ]; then
        FARMER_COUNT=$(echo "$FARMERS" | grep -v '^$' | wc -l)
        TOTAL_FARMERS="$TOTAL_FARMERS"$'\n'"$FARMERS"
        echo "OK ($FARMER_COUNT farmers)"
    else
        echo "Failed or no blocks"
        # If we can't get blocks, break
        break
    fi
    
    HEIGHT=$((BATCH_END + 1))
    
    # Small delay to avoid overwhelming the server
    sleep 0.2
done

echo ""
echo "=== Farmer Statistics ==="
echo ""

# Count unique farmers
UNIQUE_COUNT=$(echo "$TOTAL_FARMERS" | grep -v '^$' | sort -u | wc -l)
TOTAL_BLOCKS_CHECKED=$(echo "$TOTAL_FARMERS" | grep -v '^$' | wc -l)

echo "Unique farmers: $UNIQUE_COUNT"
echo "Total blocks checked: $TOTAL_BLOCKS_CHECKED"
echo ""

# Show all unique farmers
echo "Farmer addresses:"
echo "$TOTAL_FARMERS" | grep -v '^$' | sort -u
echo ""

# Show farmer distribution (top 20)
echo "Top 20 farmers by blocks mined:"
echo "$TOTAL_FARMERS" | grep -v '^$' | sort | uniq -c | sort -rn | head -20
echo ""

