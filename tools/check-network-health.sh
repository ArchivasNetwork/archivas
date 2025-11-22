#!/bin/bash

# Archivas Betanet Network Health Check Script
# Queries all three seed nodes and compares their chain tips
# Fails if any seed diverges by more than 2 blocks

set -euo pipefail

# Seed node RPC endpoints
SEED1="https://seed1.betanet.archivas.ai"
SEED2="https://seed2.betanet.archivas.ai"
SEED3="https://seed3.betanet.archivas.ai"

# Alert threshold (blocks)
MAX_GAP=2

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Archivas Betanet Network Health Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Function to query block number from a seed
get_block_number() {
    local url="$1"
    local response
    response=$(curl -s -m 5 "$url" -X POST \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' 2>/dev/null)
    
    if [ $? -ne 0 ] || [ -z "$response" ]; then
        echo "ERROR"
        return 1
    fi
    
    local hex_height
    hex_height=$(echo "$response" | jq -r '.result' 2>/dev/null)
    
    if [ "$hex_height" == "null" ] || [ -z "$hex_height" ]; then
        echo "ERROR"
        return 1
    fi
    
    # Convert hex to decimal
    printf "%d" "$hex_height" 2>/dev/null || echo "ERROR"
}

# Query all seeds
echo "Querying seed nodes..."
echo ""

SEED1_HEIGHT=$(get_block_number "$SEED1")
SEED2_HEIGHT=$(get_block_number "$SEED2")
SEED3_HEIGHT=$(get_block_number "$SEED3")

# Display results
echo "┌─────────────────────────────────────────────────┐"
echo "│ Seed Node Heights                               │"
echo "├─────────────────────────────────────────────────┤"

if [ "$SEED1_HEIGHT" == "ERROR" ]; then
    echo -e "│ Seed1: ${RED}UNREACHABLE${NC}                              │"
else
    printf "│ %-50s│\n" "Seed1: $SEED1_HEIGHT (72.251.11.191)"
fi

if [ "$SEED2_HEIGHT" == "ERROR" ]; then
    echo -e "│ Seed2: ${RED}UNREACHABLE${NC}                              │"
else
    printf "│ %-50s│\n" "Seed2: $SEED2_HEIGHT (57.129.96.158)"
fi

if [ "$SEED3_HEIGHT" == "ERROR" ]; then
    echo -e "│ Seed3: ${RED}UNREACHABLE${NC}                              │"
else
    printf "│ %-50s│\n" "Seed3: $SEED3_HEIGHT (51.89.11.4)"
fi

echo "└─────────────────────────────────────────────────┘"
echo ""

# Calculate min and max heights (excluding errors)
heights=()
[ "$SEED1_HEIGHT" != "ERROR" ] && heights+=("$SEED1_HEIGHT")
[ "$SEED2_HEIGHT" != "ERROR" ] && heights+=("$SEED2_HEIGHT")
[ "$SEED3_HEIGHT" != "ERROR" ] && heights+=("$SEED3_HEIGHT")

if [ ${#heights[@]} -eq 0 ]; then
    echo -e "${RED}CRITICAL: All seeds are unreachable!${NC}"
    exit 2
fi

# Sort heights to find min and max
IFS=$'\n' sorted=($(sort -n <<<"${heights[*]}"))
unset IFS

min_height="${sorted[0]}"
max_height="${sorted[-1]}"
gap=$((max_height - min_height))

echo "Gap Analysis:"
echo "  Min height: $min_height"
echo "  Max height: $max_height"
echo "  Gap:        $gap blocks"
echo ""

# Determine health status
if [ ${#heights[@]} -lt 3 ]; then
    echo -e "${YELLOW}WARNING: One or more seeds are unreachable${NC}"
    exit 1
elif [ $gap -gt $MAX_GAP ]; then
    echo -e "${RED}CRITICAL: Seed divergence detected (gap > ${MAX_GAP} blocks)${NC}"
    echo ""
    echo "Recommended actions:"
    echo "  1. Check seed node logs for errors"
    echo "  2. Verify no forks have occurred"
    echo "  3. Investigate which seed is out of sync"
    echo ""
    
    # Identify lagging seed(s)
    if [ "$SEED1_HEIGHT" != "ERROR" ] && [ "$SEED1_HEIGHT" -lt $max_height ]; then
        echo -e "  ${YELLOW}→ Seed1 is lagging by $((max_height - SEED1_HEIGHT)) blocks${NC}"
    fi
    if [ "$SEED2_HEIGHT" != "ERROR" ] && [ "$SEED2_HEIGHT" -lt $max_height ]; then
        echo -e "  ${YELLOW}→ Seed2 is lagging by $((max_height - SEED2_HEIGHT)) blocks${NC}"
    fi
    if [ "$SEED3_HEIGHT" != "ERROR" ] && [ "$SEED3_HEIGHT" -lt $max_height ]; then
        echo -e "  ${YELLOW}→ Seed3 is lagging by $((max_height - SEED3_HEIGHT)) blocks${NC}"
    fi
    
    exit 1
else
    echo -e "${GREEN}✓ HEALTHY: All seeds are in sync (gap ≤ ${MAX_GAP} blocks)${NC}"
    exit 0
fi

