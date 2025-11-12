#!/bin/bash

# List all addresses on the chain that have balances
# Usage: ./list-accounts.sh [rpc_url]

RPC_URL="${1:-http://127.0.0.1:8080}"

echo "Fetching all accounts with balances from $RPC_URL..."
echo ""

# Get accounts
RESPONSE=$(curl -s -w "\n%{http_code}" "$RPC_URL/accounts")
HTTP_CODE=$(echo "$RESPONSE" | tail -1)
RESPONSE_BODY=$(echo "$RESPONSE" | head -n -1)

# Check if request succeeded
if [ -z "$RESPONSE_BODY" ] || [ "$HTTP_CODE" != "200" ]; then
    if [ "$HTTP_CODE" = "404" ]; then
        echo "Error: /accounts endpoint not available. Make sure node is running latest code."
    else
        echo "Error: Failed to fetch accounts from $RPC_URL (HTTP $HTTP_CODE)"
    fi
    exit 1
fi

RESPONSE="$RESPONSE_BODY"

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

