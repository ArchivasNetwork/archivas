# Deployment Guide: v1.1.1 Hotfix - POST /submit Fix

## Critical: All Node Servers Must Be Updated

The v1.1.1 hotfix must be deployed to **ALL node servers** that accept transaction submissions. If a tester is submitting to a node that hasn't been updated, they will still get HTTP 405 errors.

## Identify All Node Servers

Check which servers are running `archivas-node`:

```bash
# On each server, check if node is running
ps aux | grep archivas-node

# Check bootnodes/IPs from your config
# Common servers:
# - ns5042239 (72.251.11.191) - ✅ Already updated
# - 57.129.148.132:9090 - ⚠️ Check if this needs update
# - Any other nodes in your network
```

## Deployment Steps for Each Node Server

**Run these commands on EVERY node server:**

```bash
# 1. Navigate to archivas directory
cd ~/archivas

# 2. Pull latest code
git pull origin main

# 3. Build new node binary
go build -o archivas-node ./cmd/archivas-node

# 4. Stop old node
pkill -f archivas-node
sleep 2

# 5. Verify stopped
ps aux | grep archivas-node | grep -v grep

# 6. Restart with your original command
# (Preserve your --genesis, --network-id, --bootnodes flags)
nohup ./archivas-node \
  --genesis genesis/devnet.genesis.json \
  --network-id devnet \
  --bootnodes "57.129.148.132:9090" \
  > logs/node.log 2>&1 &

# 7. Verify started
sleep 2
ps aux | grep archivas-node | grep -v grep
tail -20 logs/node.log
```

## Verify Hotfix is Applied

On each node, test:

```bash
# Test 1: GET /submit should return 405 with Allow: POST
curl -v http://localhost:8080/submit 2>&1 | grep -E "405|Allow"

# Expected output:
# < HTTP/1.1 405 Method Not Allowed
# < Allow: POST

# Test 2: POST /submit without Content-Type should return 415
curl -v -X POST http://localhost:8080/submit -d '{}' 2>&1 | grep -E "415|Content-Type"

# Expected output:
# < HTTP/1.1 415 Unsupported Media Type
# {"error":"Content-Type must be application/json","ok":false}
```

## Instructions for Testers

**For testers submitting transactions, they need:**

1. **Updated CLI tool** (if using `archivas-cli`):
   ```bash
   git clone https://github.com/ArchivasNetwork/archivas
   cd archivas
   git pull origin main
   go build -o archivas-cli ./cmd/archivas-cli
   ```

2. **Submit to an updated node**:
   ```bash
   ./archivas-cli broadcast tx.json http://<NODE_IP>:8080
   ```

3. **If using their own client**:
   - Must use POST method
   - Must set `Content-Type: application/json` header
   - Must not follow redirects (prevents POST→GET conversion)

## Node Servers Checklist

- [ ] **ns5042239** (72.251.11.191) - ✅ Already updated
- [ ] **57.129.148.132:9090** - ⚠️ **CHECK IF THIS NEEDS UPDATE**
- [ ] **Any other node servers** - ⚠️ **MUST BE UPDATED**

## If Tester Still Getting Errors

Ask the tester:
1. **What error are they getting?** (405, 415, 400, etc.)
2. **Which node IP are they submitting to?** (verify that node has hotfix)
3. **Are they using updated CLI?** (must have v1.1.1 code)

If they're getting 405/415, the node they're submitting to likely doesn't have the hotfix.

## Quick Fix for Single Server

If you know which server the tester is using:

```bash
# SSH to that server
ssh user@<server-ip>

# Update and restart
cd ~/archivas
git pull origin main
go build -o archivas-node ./cmd/archivas-node
pkill -f archivas-node
# Restart with original command
```

