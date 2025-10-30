# Deploy Archivas v0.2.1 - Automatic Peer Discovery

## ON BOTH SERVERS (A: 57.129.148.132, B: 72.251.11.191)

### 1. Pull and Build

```bash
cd ~/archivas
git pull origin main
go build -o archivas-node ./cmd/archivas-node
go build -o archivas-timelord ./cmd/archivas-timelord
go build -o archivas-farmer ./cmd/archivas-farmer
```

### 2. Stop Current Processes

```bash
pkill -f archivas-node
pkill -f archivas-timelord
pkill -f archivas-farmer
sleep 3
```

### 3. Start with Gossip Enabled

**SERVER A (57.129.148.132):**
```bash
nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 72.251.11.191:9090 \
  --enable-gossip \
  --gossip-interval 60s \
  --max-peers 20 \
  > logs/node-v0.2.1.log 2>&1 &

nohup ./archivas-timelord --node http://localhost:8080 > logs/timelord-v0.2.1.log 2>&1 &

nohup ./archivas-farmer farm \
  --plots ./plots-test \
  --farmer-privkey 6301448565bb68053331f80a62621a9ff0a0ac8d6863010ae5bf02237800e5c8 \
  > logs/farmer-v0.2.1.log 2>&1 &
```

**SERVER B (72.251.11.191):**
```bash
nohup ./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 57.129.148.132:9090 \
  --enable-gossip \
  --gossip-interval 60s \
  --max-peers 20 \
  > logs/node-v0.2.1.log 2>&1 &

nohup ./archivas-timelord --node http://localhost:8080 > logs/timelord-v0.2.1.log 2>&1 &
```

### 4. Verify Deployment

**Check health:**
```bash
curl http://localhost:8080/healthz
```

**Check peers:**
```bash
curl http://localhost:8080/peers
```

**Wait 65 seconds and check gossip logs:**
```bash
sleep 65
tail -50 logs/node-v0.2.1.log | grep -i gossip
```

**Expected output:**
```
[p2p] gossiped X known peers to Y connected peers
[p2p] received GOSSIP_PEERS: addrs=2 from=<IP> netID=archivas-devnet-v3
[p2p] merged N new addrs (known=M, connected=P)
```

### 5. Monitor

```bash
# Live logs
tail -f logs/node-v0.2.1.log

# Check peers file
cat data/peers.json

# Network status
curl http://localhost:8080/peers | jq
```

---

## What's New in v0.2.1

### Automatic Peer Discovery
- Nodes broadcast known peers every 60s
- Auto-connect to discovered peers (rate-limited)
- Network grows organically without manual config

### Peer Persistence
- `data/peers.json` stores known peers
- Auto-reconnect on restart
- No need to re-specify bootnodes

### Network Validation
- NetID checking prevents cross-network pollution
- Mismatched networks rejected automatically

### Rate Limiting
- Max 5 new dials per minute (configurable)
- Prevents connection storms
- Graceful degradation

---

## Testing with 3rd Node

If someone wants to join:

```bash
./archivas-node \
  --rpc :8080 \
  --p2p :9090 \
  --db ./data \
  --genesis genesis/devnet.genesis.json \
  --network-id archivas-devnet-v3 \
  --bootnodes 57.129.148.132:9090
```

Within 60 seconds:
- Node C connects to A
- A gossips "I know B at 72.251.11.191:9090"
- C auto-connects to B
- Full mesh established!

---

## Rollback (if needed)

```bash
# Disable gossip
./archivas-node ... --enable-gossip=false

# Or revert to v0.1.1
git checkout v0.1.1-devnet
go build ./cmd/archivas-node
# restart as before
```

---

**Ready to deploy!** 🚀
