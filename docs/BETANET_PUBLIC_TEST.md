# Betanet Public Test Plan

Public stress test and validation plan for 5-10 community nodes connecting via Seed3.

---

## Test Objectives

This test validates the following aspects of Archivas Betanet:

1. **P2P Connectivity:** Community nodes can discover and connect via Seed3 gateway
2. **Chain Consistency:** No forks or divergent tips across the network
3. **Address Derivation:** Unified EVM + ARCV address system works correctly
4. **Farming Compatibility:** Community farmers can submit proofs without creating forks
5. **Snapshot Bootstrap:** New nodes can bootstrap from network snapshots
6. **Network Resilience:** The network remains stable under increased peer load

---

## Test Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Archivas Betanet Topology                 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Seed1 (Private)  ←→  Seed2 (Private)  ←→  Seed3 (Public)  │
│  72.251.11.191        57.129.96.158        51.89.11.4       │
│  (Whitelist only)     (Whitelist only)     (Open)           │
│         ↑                    ↑                   ↑           │
│         └────────────────────┴───────────────────┘           │
│                            ↑                                 │
│                            │                                 │
│                ┌───────────┴───────────┐                     │
│                │                       │                     │
│         Community Node 1      Community Node 2...10         │
│         (Test Participants)   (Test Participants)           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Key Points:**
- Seed1 and Seed2 are private (whitelisted peers only)
- Seed3 is the public gateway (open to all)
- Community nodes connect ONLY via Seed3
- Seeds cross-verify each other to detect forks

---

## Prerequisites for Community Testers

### System Requirements

- **OS:** Linux (Ubuntu 22.04+ recommended) or macOS
- **RAM:** Minimum 4GB (8GB recommended)
- **Disk:** 50GB free space (for blockchain + optional plots)
- **Network:** Stable internet connection, 10 Mbps+ upload/download

### Required Software

```bash
# Install dependencies
sudo apt update
sudo apt install -y build-essential git curl jq

# Install Go 1.21+
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

---

## Community Tester Instructions

### Step 1: Clone and Build Archivas

```bash
# Clone repository
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Build node and wallet
go build -o archivas-node cmd/archivas-node/main.go
go build -o archivas-wallet cmd/archivas-wallet/main.go
go build -o archivas-farmer cmd/archivas-farmer/main.go

# Verify builds
./archivas-node --version
./archivas-wallet --version
./archivas-farmer --version
```

### Step 2: Generate Wallet (If New)

```bash
# Generate a new Ethereum-compatible wallet
./archivas-wallet new

# Expected output:
# ARCV Address: arcv1...
# EVM Address:  0x...
# Private Key:  abc123...

# IMPORTANT: Save your private key securely!
```

**Record the following for test reporting:**
- ARCV Address
- EVM Address (0x)
- Confirm they map to the same account using the converter

### Step 3: Bootstrap from Snapshot (Recommended)

**Option A: Fast sync via snapshot (recommended for testing)**

```bash
# Download and verify latest snapshot
# TODO: Add snapshot manifest URL once available
# For now, start from genesis
```

**Option B: Sync from genesis (slow, 150K+ blocks)**

```bash
# No action needed, node will IBD from Seed3
```

### Step 4: Start Node

```bash
# Create data directory
mkdir -p ~/archivas-data

# Start node connected to Seed3 (public gateway)
./archivas-node \
  --network betanet \
  --peer seed3.betanet.archivas.ai:30303 \
  --db ~/archivas-data \
  --rpc 0.0.0.0:8545 \
  --p2p 0.0.0.0:30303 \
  --enable-gossip

# Node will now sync via IBD from Seed3
```

**Expected output:**

```
[INFO] Archivas Node starting...
[INFO] Network: betanet (Chain ID: archivas-betanet-1)
[INFO] P2P listening on 0.0.0.0:30303
[INFO] RPC listening on 0.0.0.0:8545
[INFO] Connecting to peer: seed3.betanet.archivas.ai:30303
[INFO] Handshake successful with seed3.betanet.archivas.ai:30303
[INFO] Starting Initial Block Download (IBD)...
[INFO] Syncing blocks 0 - 512...
[INFO] Synced to height 512 (gap: 150000)
```

**Monitor sync progress:**

```bash
# Check current height (in another terminal)
curl -s http://localhost:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result' | xargs printf "%d\n"

# Check peer count
curl -s http://localhost:8545/health | jq
```

### Step 5: Verify Chain Identity

**Confirm your node is on the correct chain:**

```bash
# Check Chain ID
curl -s http://localhost:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' | jq

# Expected: {"result":"0x66c"}  (1644 in decimal)

# Check genesis hash (should match seeds)
# TODO: Add genesis hash verification command
```

### Step 6: Verify Address Derivation

**Test that your ARCV and EVM addresses are consistent:**

```bash
# Convert ARCV to EVM
curl -s http://localhost:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "arcv_toHexAddress",
    "params": ["YOUR_ARCV_ADDRESS"],
    "id": 1
  }' | jq

# Convert EVM to ARCV
curl -s http://localhost:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "arcv_fromHexAddress",
    "params": ["YOUR_EVM_ADDRESS"],
    "id": 1
  }' | jq

# Both should point to the same account
```

### Step 7: Optional - Start Farming

**Only if you want to participate in farming during the test:**

```bash
# Create plots (k=20 = 128MB for testing)
./archivas-farmer plot \
  --path ~/archivas-plots \
  --size 20 \
  --farmer-privkey YOUR_PRIVATE_KEY

# Start farming
./archivas-farmer farm \
  --plots ~/archivas-plots \
  --node http://localhost:8545 \
  --farmer-privkey YOUR_PRIVATE_KEY

# Monitor farmer logs
# Look for: "Submitted proof for challenge height X"
```

---

## Reporting Requirements

### Report Every Hour During Test

Submit the following information via GitHub issue or Discord:

```bash
# 1. Node Height
curl -s http://localhost:8545 -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result' | xargs printf "Height: %d\n"

# 2. Peer Count
curl -s http://localhost:8545/health | jq '.peers'

# 3. Your Addresses
echo "ARCV: YOUR_ARCV_ADDRESS"
echo "EVM:  YOUR_EVM_ADDRESS"

# 4. Any errors or warnings in logs
journalctl -u archivas-node --since "1 hour ago" | grep -i "error\|warning\|reorg\|fork"
```

**Report Format (GitHub issue template):**

```markdown
## Test Report - [YOUR_NODE_NAME]

**Timestamp:** 2025-11-21 14:00 UTC

**Node Info:**
- ARCV Address: arcv1...
- EVM Address: 0x...
- Node Height: 151234
- Peer Count: 3

**Status:**
- [ ] Synced to latest tip
- [ ] No reorgs detected
- [ ] Address conversion works
- [ ] Farming active (if applicable)

**Issues/Observations:**
- None / [Describe any issues]

**Logs:**
```
[Paste any relevant log lines]
```
```

---

## Internal Monitoring Checklist

**Archivas Core Team will monitor:**

1. **Seed Health:**
   - Seed1, Seed2, Seed3 all at same tip height
   - No divergence > 2 blocks for > 60 seconds
   - Peer counts stable

2. **Genesis Verification:**
   - All nodes report same genesis hash
   - Chain ID consistent (1644)

3. **Fork Detection:**
   - Monitor for "chain reorg" logs
   - Check for prev hash mismatches
   - Verify no community node creates alternate chain

4. **Address Consistency:**
   - Verify ARCV ↔ EVM conversions match across all nodes
   - Check that farmer rewards go to correct addresses

5. **Snapshot Integrity:**
   - If snapshots are provided, verify they're on canonical chain
   - Check that bootstrapped nodes sync correctly

**Alert Triggers:**
- Seed tip divergence > 2 blocks for > 60s
- Any node reports different genesis hash
- Chain reorg > 10 blocks
- Community node creates longer chain than seeds

---

## Monitoring Script

A bash script `tools/check-network-health.sh` is provided to automate health checks:

```bash
# Run health check
./tools/check-network-health.sh

# Expected output:
# Seed1: Height 151234 (72.251.11.191)
# Seed2: Height 151234 (57.129.96.158)
# Seed3: Height 151234 (51.89.11.4)
# Max gap: 0 blocks
# Status: HEALTHY
```

**If gaps exceed 2 blocks:**

```
# Output:
# Seed1: Height 151234
# Seed2: Height 151230  <- LAGGING
# Seed3: Height 151234
# Max gap: 4 blocks
# Status: WARNING - Seed2 is behind
```

---

## Test Duration and Phases

### Phase 1: Initial Sync (2-4 hours)

**Goal:** All community nodes sync to current tip

**Success Criteria:**
- All nodes reach tip within 4 hours
- No forks detected
- Peer counts stable (2-5 peers)

### Phase 2: Steady State (24 hours)

**Goal:** Network remains stable with 5-10 active nodes

**Success Criteria:**
- All nodes stay within 2 blocks of tip
- No reorgs > 5 blocks
- Seeds remain canonical source

### Phase 3: Farming Test (12 hours)

**Goal:** Community farmers can submit proofs without causing forks

**Success Criteria:**
- Farming proofs accepted by seeds
- No alternative tips created by community farmers
- Block rewards go to correct farmer addresses

### Phase 4: Snapshot Bootstrap (2 hours)

**Goal:** New nodes can bootstrap from snapshot

**Success Criteria:**
- Snapshot manifests verify correctly
- Bootstrapped nodes sync to tip quickly
- No chain mismatches after bootstrap

---

## Known Limitations and Expected Behavior

1. **No Peer Discovery:** Community nodes must manually connect to Seed3 (no DHT yet)
2. **Fixed Gas Price:** No dynamic gas pricing (always 1 gwei)
3. **Simplified EVM:** Contract execution may have edge cases
4. **No Faucet:** Must earn RCHV via farming
5. **IBD Can Be Slow:** 150K+ blocks may take 1-2 hours to sync

**These are expected and not considered test failures.**

---

## Troubleshooting for Community Testers

### "Connection refused" to Seed3

**Solution:**
- Verify Seed3 is online: `ping seed3.betanet.archivas.ai`
- Check firewall allows outbound on port 30303
- Try alternate RPC: `https://seed3.betanet.archivas.ai` (for RPC queries)

### "Handshake failed: Chain ID mismatch"

**Solution:**
- You may be on the wrong network or wrong genesis
- Delete `~/archivas-data` and re-sync from genesis:

```bash
rm -rf ~/archivas-data
./archivas-node --network betanet --peer seed3.betanet.archivas.ai:30303
```

### "IBD stuck at height X"

**Solution:**
- Check Seed3 is responding:

```bash
curl -s https://seed3.betanet.archivas.ai -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

- Restart node to retry IBD:

```bash
pkill archivas-node
./archivas-node --network betanet --peer seed3.betanet.archivas.ai:30303
```

### "My node created a longer chain"

**CRITICAL:** This indicates a fork.

**Solution:**
1. **STOP your node immediately:**
   ```bash
   pkill archivas-node
   pkill archivas-farmer
   ```

2. **Report to Archivas team via GitHub issue or Discord**

3. **Delete local chain and resync:**
   ```bash
   rm -rf ~/archivas-data
   ./archivas-node --network betanet --peer seed3.betanet.archivas.ai:30303
   ```

4. **Do NOT continue farming until fork is resolved**

---

## Post-Test Deliverables

After the test concludes, the Core Team will publish:

1. **Test Report:**
   - Number of participants
   - Network uptime
   - Fork events (if any)
   - Performance metrics

2. **Identified Issues:**
   - Bugs discovered during test
   - GitHub issues for each bug
   - Severity and priority

3. **Next Steps:**
   - Fixes for any critical issues
   - Timeline for mainnet readiness
   - Further testing plans

---

## How to Participate

1. **Join the Discord:** [Archivas Discord Link]
2. **Announce participation in #betanet-testing channel**
3. **Follow the instructions above**
4. **Report hourly using the template**
5. **Be available for troubleshooting during test window**

**Test Window:**  
**TBD - Will be announced 48 hours in advance**

---

## Questions or Issues

- **GitHub Issues:** https://github.com/ArchivasNetwork/archivas/issues
- **Discord:** [Link to Discord]
- **Documentation:** https://docs.archivas.ai

---

**Thank you for helping test Archivas Betanet! 🚀**

