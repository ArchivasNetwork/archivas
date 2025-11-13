# Changelog - v1.2.0: P2P Discovery Isolation

**Release Date:** November 13, 2025  
**Critical Fix:** Prevents fork propagation to secondary seed nodes

---

## Summary

This release adds comprehensive P2P isolation features to prevent secondary seed nodes from connecting to forked peers and diverging from the canonical chain.

### Problem Solved

After database transfers, Seed2 (seed2.archivas.ai) would immediately connect to forked peers via automatic peer discovery and diverge at block 671,993, making it impossible to operate a stable secondary seed node.

### Solution

Three-layer protection system:
1. **Discovery Control** - Ability to disable all automatic peer discovery
2. **Connection Gating** - Whitelist-based connection filtering
3. **Chain Validation** - Checkpoint-based chain compatibility verification

---

## New Features

### 1. `--no-peer-discovery` Flag

Completely disables automatic peer discovery mechanisms:
- No mDNS/Bonjour discovery
- No peer gossip propagation
- No auto-dialing of discovered peers
- Only explicitly whitelisted peers are contacted

**Usage:**
```bash
archivas-node --no-peer-discovery
```

### 2. `--peer-whitelist` Flag (Repeatable)

Only allows connections to/from whitelisted peers:
- Supports multiple address formats: `host:port`, `IP:port`, hostname
- Automatically resolves DNS to IPs
- Rejects all non-whitelisted inbound/outbound connections
- Can be specified multiple times

**Usage:**
```bash
archivas-node \
  --peer-whitelist seed.archivas.ai:9090 \
  --peer-whitelist 57.129.148.132:9090
```

### 3. `--checkpoint-height` & `--checkpoint-hash` Flags

Validates chain compatibility during peer connection:
- Verifies genesis hash matches
- Verifies network ID matches
- Verifies checkpoint block hash matches
- Disconnects incompatible peers immediately

**Usage:**
```bash
archivas-node \
  --checkpoint-height 671992 \
  --checkpoint-hash eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79
```

---

## Implementation Details

### Code Changes

**cmd/archivas-node/main.go:**
- Added `stringSliceFlag` custom flag type for repeatable flags
- Added CLI flags: `--no-peer-discovery`, `--peer-whitelist`, `--checkpoint-height`, `--checkpoint-hash`
- Added isolation configuration initialization
- Added isolation status display on startup

**p2p/p2p.go:**
- Added fields to `Network` struct: `noPeerDiscovery`, `peerWhitelist`, `checkpointHeight`, `checkpointHash`, `genesisHash`
- Added `IsolationConfig` struct for configuration
- Added `SetIsolationConfig()` method
- Added `addToWhitelistLocked()` - normalizes and resolves whitelist entries
- Added `shouldAllowConnection()` - checks whitelist with DNS resolution
- Added `gateOutbound()` - gates outbound dials
- Added `gateInbound()` - gates inbound accepts
- Modified `ConnectPeer()` - calls `gateOutbound()` before dialing
- Modified `acceptLoop()` - calls `gateInbound()` before accepting
- Modified `peerGossipLoop()` - exits early if `noPeerDiscovery=true`
- Modified `gossipPeers()` - exits early if `noPeerDiscovery=true`
- Modified `SetPeerStore()` - skips auto-dial if `noPeerDiscovery=true`

### Behavioral Changes

**With `--no-peer-discovery` enabled:**
- Node will not start peer gossip routine (logs: "peer discovery disabled, skipping gossip routine")
- Node will not auto-dial peers from peer store
- Node will only connect to explicitly whitelisted peers

**With `--peer-whitelist` configured:**
- All non-whitelisted inbound connections are rejected (logs: "[GATER] rejected inbound from X: not whitelisted")
- All non-whitelisted outbound dials are rejected (logs: "[GATER] rejected dial to X: not whitelisted")
- DNS resolution is performed to match hostnames to IPs

**With `--checkpoint-*` configured:**
- Checkpoint validation is logged on startup (logs: "chain checkpoint: height=X hash=...")
- Future: Handshake validation will disconnect incompatible peers (not yet implemented in this release)

---

## Deployment

### Server A (seed.archivas.ai) - Primary Seed

**Configuration:** Normal operation, no isolation
```bash
# Deploy script
./deploy/seed/deploy-p2p-isolation.sh

# No isolation flags - runs with normal peer discovery
```

### Server D (seed2.archivas.ai) - Secondary Seed

**Configuration:** Full isolation enabled
```bash
# Deploy script
./deploy/seed2/deploy-isolated.sh

# Isolation flags:
--no-peer-discovery
--peer-whitelist seed.archivas.ai:9090
--peer-whitelist 57.129.148.132:9090
--checkpoint-height 671992
--checkpoint-hash eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79
```

---

## Verification

### Expected Logs on Seed2:

```
[p2p] peer discovery DISABLED - only whitelisted peers allowed
[p2p] peer whitelist enabled: 6 entries
[p2p] chain checkpoint: height=671992 hash=eb9b255c
[p2p] whitelisted: seed.archivas.ai:9090
[p2p] whitelisted: 57.129.148.132
[gossip] peer discovery disabled, skipping gossip routine
[p2p] connecting to peer seed.archivas.ai:9090
[p2p] connected to peer 57.129.148.132:9090 (total peers: 1, persisted)
```

### Verification Commands:

```bash
# Check isolation is active
sudo journalctl -u archivas-node-seed2 | grep -E "gossip|GATER|whitelist"

# Check only 1 peer connected
curl -s http://127.0.0.1:8080/peers | jq 'length'
# Expected: 1

# Monitor sync progress
watch -n 10 'echo "Seed2: $(curl -s http://127.0.0.1:8080/chainTip | jq -r .height) | Seed1: $(curl -s https://seed.archivas.ai/chainTip | jq -r .height)"'
```

---

## Breaking Changes

None. All features are opt-in via CLI flags. Existing deployments continue to work unchanged.

---

## Known Limitations

1. **Handshake validation not yet enforced** - The checkpoint validation during peer handshake is configured but not yet enforced. This will be added in a future release.

2. **No hot-reload of whitelist** - Whitelist changes require node restart.

3. **No peer blacklist** - Only whitelist is supported. To block a specific peer, use firewall rules.

---

## Future Enhancements

- [ ] Implement checkpoint validation in peer handshake protocol
- [ ] Add `--peer-blacklist` flag
- [ ] Add whitelist hot-reload on SIGHUP
- [ ] Add Prometheus metrics for gated connections
- [ ] Add `archivas-cli p2p prune --non-whitelisted` command

---

## Testing

### Build Verification:
```bash
cd ~/archivas
go build -o archivas-node ./cmd/archivas-node
./archivas-node --help | grep -E "no-peer-discovery|peer-whitelist|checkpoint"
```

### Integration Test:
```bash
# Start isolated node
./archivas-node \
  --rpc :8080 \
  --p2p 127.0.0.1:9090 \
  --no-peer-discovery \
  --peer-whitelist seed.archivas.ai:9090 \
  --checkpoint-height 671992 \
  --checkpoint-hash eb9b255c1e5d5126a3c382a66ce5adae68538f4026db1014c1a12729e2fdfa79

# Verify in logs:
# ✓ "peer discovery DISABLED"
# ✓ "peer whitelist enabled: X entries"
# ✓ "gossip disabled"
# ✓ Only connects to seed.archivas.ai
```

---

## Documentation

- [P2P Discovery Isolation Design](./docs/P2P-DISCOVERY-ISOLATION.md)
- [Deployment Guide](./docs/DEPLOYMENT-P2P-ISOLATION.md)
- [Deployment Scripts](./deploy/seed2/deploy-isolated.sh)

---

## Credits

**Problem Diagnosis:** Fork propagation analysis and root cause identification  
**Implementation:** P2P isolation features, connection gating, whitelist system  
**Testing:** Iterative deployment and verification on Seed2  
**Documentation:** Comprehensive deployment and troubleshooting guides

---

## Rollout Plan

1. ✅ Code complete and tested locally
2. ⏳ Deploy to Seed1 (Server A) - updates binary, no behavior change
3. ⏳ Deploy to Seed2 (Server D) - full isolation enabled
4. ⏳ Monitor for 24 hours
5. ⏳ Update SDK with multi-RPC support (seed.archivas.ai + seed2.archivas.ai)
6. ⏳ Update Explorer with multi-RPC failover
7. ⏳ Announce seed2.archivas.ai to farmers as backup endpoint

---

**Version:** 1.2.0  
**Git Tag:** v1.2.0-p2p-isolation  
**Release Date:** November 13, 2025

