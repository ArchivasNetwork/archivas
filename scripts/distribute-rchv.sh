#!/usr/bin/env bash
# distribute-rchv.sh - Send RCHV to multiple addresses using archivas-cli
# Usage: ./distribute-rchv.sh <mnemonic-24-words> [amount-per-recipient] [rpc-url]

set -euo pipefail

if [ $# -lt 1 ]; then
    echo "❌ Usage: $0 '<mnemonic-24-words>' [amount-per-recipient] [rpc-url]"
    echo ""
    echo "Example:"
    echo "  $0 'word1 word2 ... word24' 20000 http://127.0.0.1:8080"
    echo ""
    echo "Defaults:"
    echo "  amount-per-recipient: 20000 RCHV"
    echo "  rpc-url: http://127.0.0.1:8080"
    exit 1
fi

MNEMONIC="$1"
AMOUNT_PER_RECIPIENT="${2:-20000}"
RPC_URL="${3:-http://127.0.0.1:8080}"
FEE=100000  # 0.001 RCHV standard fee
CLI_BIN="${CLI_BIN:-archivas-cli}"

# Check if archivas-cli exists
if ! command -v "$CLI_BIN" &> /dev/null; then
    echo "❌ Error: archivas-cli not found in PATH"
    echo "   Build it with: cd /home/ubuntu/archivas && go build -o archivas-cli ./cmd/archivas-cli"
    echo "   Or set CLI_BIN env var: export CLI_BIN=/path/to/archivas-cli"
    exit 1
fi

# Recipients list (29 addresses, 1 duplicate removed)
RECIPIENTS=(
    "arcv1hxjc5snn3sv47wskh79ezt2dkdnwh03lcm36tx"
    "arcv10zfdrk6mwtmddda5vue25ctv2z7taw4l0rups5"
    "arcv1dxy7lpgyk99zfrtaul6kyssvqzfr4dack0nw8x"
    "arcv1zxxmdzjp5y9ae7m6frdu7uqh8zs95szve5q039"
    "arcv1lvklfm467p5snquavj2avgusu73h8q74hyhf4m"
    "arcv103eve4ztlp4u9yr9z04709zct8nvdn97npqp00"
    "arcv1yd0zyhje38f42yf2ts7jnws95vqa82szakd2ym"
    "arcv1slzdg88rfz6hxj7n3fjq8xtglkkjxfk94pqlau"
    "arcv1cdhunmpkungm3u274zwz70pl24fhzpfc4m7ud6"
    "arcv1zut2kfps5ks5xgxz5p6gw6wg9eu8qjgwqlv6g5"
    "arcv1z7k56ma4589xnlvk9uea034klq6zp095472x53"
    "arcv1kff5jmqqet7j0zfjyacyqjparnpvlss3q2wx7c"
    "arcv1hrplyjfvkejumf7gqpeanez7urhtp2p0a8fr79"
    "arcv196k5kwrd5zgtcgfxjrcr8cfxxx2f67vvnrdcqy"
    "arcv1pawnty0vekd5uuu8k7908lczncgdyvw7nqp8xw"
    "arcv1mfed00qyh373f622frdfk7zd9ge4h7jagqv2yy"
    "arcv155dafkxvu5mzxqf0ulul9ndyhccw5fvy33y9g0"
    "arcv1wlkfwewd08dcjrg3tt5hzn3lpn7h00y40c22f3"
    "arcv1svg999qmmhfzc9q0qdxy0gv53ul7j7ug0wwc7d"
    "arcv1epy2ykvxp9kns9nxh69q4lgrt4x0205a3s6kwu"
    "arcv1n7tgcr09uqt9q47m6synthk4l4h65mvp4hgfzl"
    "arcv1w93nj99yem6fyq4rn7h0hzr3lpd2p8tmz0xwrr"
    "arcv1cxmutvxlmr08uz0ezn6h0hwskr4kh4yus239cz"
    "arcv1mzl5w0cz64pcm4a29wxgulc6prak9w7txmyu55"
    "arcv1cvy99ah68gw9r0tpk7l662mwpz32ar46zjgjh3"
    "arcv1pm5jazzsvgujrqwa68xufdrrhxxre43dc0x99w"
    "arcv1ynet278eevqwprzusjrcmjz9m66gaf5v7lkegp"
    "arcv1knr2mr05me3kkkzy5z476cq632ezju0dx6snnz"
)

TOTAL_RECIPIENTS=${#RECIPIENTS[@]}
AMOUNT_BASE_UNITS=$((AMOUNT_PER_RECIPIENT * 100000000))  # Convert RCHV to base units
TOTAL_AMOUNT=$((TOTAL_RECIPIENTS * AMOUNT_PER_RECIPIENT))

# Derive sender address from mnemonic
echo "🔑 Deriving sender address from mnemonic..."
FROM_ADDRESS=$("$CLI_BIN" addr $MNEMONIC 2>/dev/null)
if [ -z "$FROM_ADDRESS" ]; then
    echo "❌ Error: Failed to derive address from mnemonic"
    exit 1
fi
echo "   Sender address: $FROM_ADDRESS"
echo ""

echo "╔════════════════════════════════════════════════════════════╗"
echo "║         RCHV Distribution Script                           ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "📊 Distribution Summary:"
echo "   From Address:         $FROM_ADDRESS"
echo "   Recipients:           $TOTAL_RECIPIENTS addresses"
echo "   Amount per recipient: $AMOUNT_PER_RECIPIENT RCHV"
echo "   Fee per tx:           0.001 RCHV"
echo "   Total amount:         $TOTAL_AMOUNT RCHV"
echo "   Total fees:           $((TOTAL_RECIPIENTS * FEE / 100000000)) RCHV"
echo "   Grand total:          $((TOTAL_AMOUNT + (TOTAL_RECIPIENTS * FEE / 100000000))) RCHV"
echo "   RPC URL:              $RPC_URL"
echo ""

# Check sender balance and nonce
echo "🔍 Checking sender balance and nonce..."
ACCOUNT_DATA=$(curl -s "$RPC_URL/accounts" | jq -r ".accounts[] | select(.address == \"$FROM_ADDRESS\")")

if [ -z "$ACCOUNT_DATA" ]; then
    echo "❌ Error: Could not find account data for address $FROM_ADDRESS"
    exit 1
fi

SENDER_BALANCE=$(echo "$ACCOUNT_DATA" | jq -r '.balance')
SENDER_NONCE=$(echo "$ACCOUNT_DATA" | jq -r '.nonce')
SENDER_BALANCE_RCHV=$((SENDER_BALANCE / 100000000))

echo "   Balance: $SENDER_BALANCE_RCHV RCHV"
echo "   Nonce:   $SENDER_NONCE"
echo ""

REQUIRED_AMOUNT=$((TOTAL_AMOUNT + (TOTAL_RECIPIENTS * FEE / 100000000)))
if [ "$SENDER_BALANCE_RCHV" -lt "$REQUIRED_AMOUNT" ]; then
    echo "❌ Error: Insufficient balance!"
    echo "   Required: $REQUIRED_AMOUNT RCHV (includes fees)"
    echo "   Available: $SENDER_BALANCE_RCHV RCHV"
    echo "   Shortfall: $((REQUIRED_AMOUNT - SENDER_BALANCE_RCHV)) RCHV"
    exit 1
fi

echo "✅ Sufficient balance available"
echo ""
echo "⚠️  WARNING: This will send $TOTAL_AMOUNT RCHV to $TOTAL_RECIPIENTS addresses"
echo "             (Plus ~$((TOTAL_RECIPIENTS / 1000)) RCHV in transaction fees)"
echo ""
read -p "Do you want to proceed? (yes/no): " -r CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "❌ Distribution cancelled"
    exit 0
fi

echo ""
echo "🚀 Starting distribution..."
echo ""

SUCCESS_COUNT=0
FAIL_COUNT=0
CURRENT_NONCE=$SENDER_NONCE
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

for i in "${!RECIPIENTS[@]}"; do
    RECIPIENT="${RECIPIENTS[$i]}"
    RECIPIENT_NUM=$((i + 1))
    TX_FILE="$TMP_DIR/tx_${RECIPIENT_NUM}.json"
    
    echo -n "[$RECIPIENT_NUM/$TOTAL_RECIPIENTS] $RECIPIENT... "
    
    # Sign transaction
    if ! "$CLI_BIN" sign-transfer \
        --from-mnemonic "$MNEMONIC" \
        --to "$RECIPIENT" \
        --amount "$AMOUNT_BASE_UNITS" \
        --fee "$FEE" \
        --nonce "$CURRENT_NONCE" \
        --out "$TX_FILE" > /dev/null 2>&1; then
        echo "❌ (sign failed)"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        continue
    fi
    
    # Broadcast transaction
    if "$CLI_BIN" broadcast "$TX_FILE" "$RPC_URL" > /dev/null 2>&1; then
        echo "✅"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        CURRENT_NONCE=$((CURRENT_NONCE + 1))
    else
        echo "❌ (broadcast failed)"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    # Small delay to avoid overwhelming the RPC
    sleep 0.5
done

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║         Distribution Complete                              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo "📊 Results:"
echo "   ✅ Successful: $SUCCESS_COUNT / $TOTAL_RECIPIENTS"
echo "   ❌ Failed: $FAIL_COUNT / $TOTAL_RECIPIENTS"
echo ""

if [ $FAIL_COUNT -gt 0 ]; then
    echo "⚠️  Some transactions failed."
    echo "   You may need to manually retry failed transactions."
    echo "   Current nonce after successful txs: $CURRENT_NONCE"
    exit 1
fi

echo "✅ All distributions completed successfully!"
echo ""
echo "🔍 Verify distributions:"
echo "   curl -s $RPC_URL/accounts | jq -r '.accounts[] | select(.address | IN(\"arcv1hxjc5snn3sv47wskh79ezt2dkdnwh03lcm36tx\", \"arcv10zfdrk6mwtmddda5vue25ctv2z7taw4l0rups5\")) | \"\(.address): \(.balance / 100000000) RCHV\"'"

