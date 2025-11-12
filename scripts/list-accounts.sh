#!/bin/bash

# List all addresses on the chain that have balances
# Usage: ./list-accounts.sh [rpc_url]

RPC_URL="${1:-http://127.0.0.1:8080}"

echo "Fetching all accounts with balances from $RPC_URL..."
echo ""

# Get accounts
RESPONSE=$(curl -s "$RPC_URL/accounts")

# Check if request succeeded
if [ $? -ne 0 ] || [ -z "$RESPONSE" ]; then
    echo "Error: Failed to fetch accounts from $RPC_URL"
    exit 1
fi

# Check if endpoint exists (might return 404)
if echo "$RESPONSE" | grep -q "404\|Not Found"; then
    echo "Error: /accounts endpoint not available. Make sure node is running latest code."
    exit 1
fi

# Extract count
COUNT=$(echo "$RESPONSE" | jq -r '.count // 0')

echo "Total accounts with balance: $COUNT"
echo ""

# Show accounts
if [ "$COUNT" -gt 0 ]; then
    echo "Accounts:"
    echo "$RESPONSE" | jq -r '.accounts[] | "\(.address) | Balance: \(.balance) | Nonce: \(.nonce)"'
    echo ""
    
    # Show summary
    TOTAL_BALANCE=$(echo "$RESPONSE" | jq '[.accounts[].balance] | add')
    echo "Total balance across all accounts: $TOTAL_BALANCE"
else
    echo "No accounts with balance found."
fi

