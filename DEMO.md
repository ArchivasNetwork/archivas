# Archivas Devnet - Demo Guide

## 🔐 Milestone 2 Complete!

Archivas now has **cryptographic ownership**. Only the holder of a private key can spend their RCHV.

## What Works

✅ World state with account balances and nonces  
✅ Transaction type for sending RCHV  
✅ Mempool for pending transactions  
✅ Block production every 20 seconds  
✅ Transaction processing and state transitions  
✅ HTTP RPC server for querying and submitting  
✅ **Wallet CLI with secp256k1 keypairs**  
✅ **Bech32 addresses (arcv prefix)**  
✅ **Transaction signing and verification**  
✅ **Cryptographic ownership enforcement**  

## Running the Node

```bash
# Build the node
go build -o archivas-node ./cmd/archivas-node

# Run it
./archivas-node
```

The node will:
- Initialize world state with genesis allocation (1B RCHV)
- Start RPC server on port 8080
- Produce blocks every 20 seconds
- Process pending transactions in each block

## Option 1: Using the Wallet CLI (Recommended)

### 1. Generate a Wallet

```bash
go run ./cmd/archivas-wallet new
```

Output:
```
🔐 New Archivas Wallet Generated

Address:     arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
Public Key:  031b2cb8df7a31f463b72468e731786df04fda65b1260e7e77fc96d11b5ec69e97
Private Key: <EXAMPLE_PRIVATE_KEY_DO_NOT_USE>

⚠️  KEEP YOUR PRIVATE KEY SECRET! Anyone with access can spend your RCHV.
```

### 2. Fund Your Wallet in Genesis

Edit `config/params_devnet.go` to include your address:

```go
var GenesisAlloc = map[string]int64{
    "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7": 1_000_000_000_00000000, // 1B RCHV
}
```

### 3. Start the Node

```bash
go build -o archivas-node ./cmd/archivas-node
./archivas-node
```

### 4. Send RCHV

```bash
go run ./cmd/archivas-wallet send \
  --from-privkey <YOUR_PRIVATE_KEY_HERE> \
  --to arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt \
  --amount 50000000000 \
  --fee 100000
```

Output:
```
📊 Current balance: 1000000000.00000000 RCHV (nonce: 0)
📝 Sending 500.00000000 RCHV from arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7 to arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt
💸 Fee: 0.00100000 RCHV
✅ success: Transaction added to mempool
⏳ Transaction will be included in the next block (~20 seconds)
```

The wallet CLI:
- ✅ Automatically derives your address from your private key
- ✅ Queries the node for your current nonce
- ✅ Signs the transaction with your private key
- ✅ Submits it to the node

### 5. Check Balances

```bash
curl --noproxy "*" http://localhost:8080/balance/arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7
curl --noproxy "*" http://localhost:8080/balance/arcv1m8z057hhhzfxep83lk77ncuk8ugc4n3p4089kt
```

## Option 2: Using the RPC API Directly

If you want to manually construct and sign transactions, you can use the RPC API directly.

### 1. Check Server Status

```bash
curl --noproxy "*" http://localhost:8080/
```

Response:
```json
{"status":"ok","message":"Archivas Devnet RPC Server"}
```

### 2. Query Account Balance

```bash
curl --noproxy "*" http://localhost:8080/balance/arcv1genesis000000000000000000
```

Response:
```json
{
  "address": "arcv1genesis000000000000000000",
  "balance": 100000000000000000,
  "nonce": 0
}
```

Balance is in base units (8 decimals). To convert to RCHV: `balance / 100000000`

### 3. Submit a Transaction

**Note:** As of Milestone 2, you must sign transactions with a private key. Use the wallet CLI instead for easier transaction submission.

Manual transaction submission requires:
- `from`: Sender bech32 address
- `to`: Recipient bech32 address
- `amount`: Amount in base units
- `fee`: Transaction fee in base units
- `nonce`: Must match sender's current nonce
- `senderPubKey`: Public key of sender (hex encoded, compressed)
- `signature`: ECDSA signature over transaction hash (DER format, hex encoded)

See `wallet/tx_sign.go` for signing implementation details.

### 4. Wait for Block Production

Transactions are included in the next block (max 20 seconds).

Watch the node output for:
```
⛏️  Produced block 1 with 2 txs (total chain length: 2)
```

### 5. Verify Balances Changed

```bash
# Check all three accounts
curl --noproxy "*" http://localhost:8080/balance/arcv1genesis000000000000000000
curl --noproxy "*" http://localhost:8080/balance/arcv1alice00000000000000000000
curl --noproxy "*" http://localhost:8080/balance/arcv1bob000000000000000000000
```

## Example Session

```bash
# Start the node
./archivas-node

# In another terminal...

# Send 100 RCHV from genesis to alice
curl --noproxy "*" -X POST http://localhost:8080/submitTx \
  -H "Content-Type: application/json" \
  -d '{"from":"arcv1genesis000000000000000000","to":"arcv1alice00000000000000000000","amount":10000000000,"fee":100000,"nonce":0}'

# Send 50 RCHV from genesis to bob
curl --noproxy "*" -X POST http://localhost:8080/submitTx \
  -H "Content-Type: application/json" \
  -d '{"from":"arcv1genesis000000000000000000","to":"arcv1bob000000000000000000000","amount":5000000000,"fee":50000,"nonce":1}'

# Wait 20 seconds for block production...

# Check balances
curl --noproxy "*" http://localhost:8080/balance/arcv1genesis000000000000000000
# Result: balance=99999984999850000, nonce=2

curl --noproxy "*" http://localhost:8080/balance/arcv1alice00000000000000000000
# Result: balance=10000000000 (100 RCHV)

curl --noproxy "*" http://localhost:8080/balance/arcv1bob000000000000000000000
# Result: balance=5000000000 (50 RCHV)
```

## 🎉 You just sent RCHV on Archivas!

Genesis account successfully transferred:
- 100 RCHV to alice
- 50 RCHV to bob
- Paid 0.00150000 RCHV in fees (burned)
- Nonce incremented from 0 → 2

## What's Next (Milestone 3)

- [ ] Proof-of-Space plot generation
- [ ] PoSpace challenge-response mechanism  
- [ ] Replace time-based blocks with PoSpace winner selection

## Future Milestones

- **Milestone 3:** Proof-of-Space implementation
- **Milestone 4:** VDF (Verifiable Delay Function)
- **Milestone 5:** Full consensus (PoSpace + Time)
- **Milestone 6:** P2P networking and multi-node testnet

