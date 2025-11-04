# Known Issues

Active bugs and limitations in Archivas.

---

## Critical

### IBD Stuck at Low Height (v1.2.0)

**Symptom:** New nodes sync first 3 blocks then get stuck, repeatedly trying to sync from height 4 but saying "caught up to tip at height 3" even when peers are at 80,000+.

**Affected:** Clean nodes syncing from scratch

**Workaround:** Copy `data/` directory from synced node

**Root Cause:** IBD logic in `p2p/ibd.go` has bug with large gaps (>1000 blocks)

**Fix Required:**
- Detect when peer tip >> local tip (e.g., 80,000 blocks behind)
- Request blocks in larger batches (not 1-by-1)
- Proper state machine for IBD mode
- Progress tracking and resumption

**Tracked in:** https://github.com/ArchivasNetwork/archivas/issues/TBD

---

## Medium Priority

### User Transactions Not in Historical Blocks

**Symptom:** Transactions sent before node restart don't appear in block's `txs` array via API

**Affected:** `/block/<height>` and `/tx/recent` for old blocks

**Workaround:** Transactions ARE applied to state (balances update), just not visible in API

**Cause:** Blocks loaded from disk use old serialization without transaction details

**Fix:** Re-serialize old blocks with new format on load, or reindex on upgrade

---

### HTTPS API Occasionally Returns 502

**Symptom:** Nginx returns 502 Bad Gateway even when node is running

**Affected:** Public API at seed.archivas.ai

**Workaround:** Restart archivas-node service

**Cause:** HTTP server stops responding (possible deadlock or panic)

**Fix:** Add timeout/keepalive, investigate deadlock, add health checks

---

## Low Priority

### Server A Wins More Blocks Than Expected

**Symptom:** Server A (1 plot) wins ~83% of blocks vs Server C (6 plots) winning ~17%

**Expected:** 14% vs 86% based on plot count

**Cause:** Server A has localhost advantage (0ms latency vs 50-300ms for Server C via HTTPS)

**Impact:** Testnet only - demonstrates latency matters

**Fix:** Not a bug, but could add timestamp-based fairness or require minimum propagation delay

---

### `/tx/recent` Returns Empty Even With Transactions

**Symptom:** Recent transactions don't appear in `/tx/recent` endpoint

**Affected:** New blocks created after v1.2.0 deploy

**Cause:** Transaction extraction logic in `rpc/farming.go` may have type conversion issue

**Workaround:** Query specific block with `/block/<height>` to see transactions

**Fix:** Debug transaction serialization in GetRecentBlocks

---

## Resolved

### ✅ CORS Duplication (Fixed in v1.2.0)

Browser error "Access-Control-Allow-Origin: *, *" - fixed by removing Nginx CORS headers

### ✅ Farmer Field Null (Fixed in v1.2.0)

Blocks showed `"farmer": null` - fixed by updating serialization

### ✅ PoSpace Verification Bug (Fixed in v1.1.1)

Quality comparison was backwards - fixed to `quality <= difficulty`

---

## Limitations (By Design)

### No Historical Transaction Indexing

**Status:** As designed  
**Limitation:** Only recent ~100 blocks of transactions kept in memory  
**Impact:** Can't query full transaction history for an account  
**Future:** Add optional indexer module

### Fixed Fee Market

**Status:** As designed for testnet  
**Limitation:** Fee is always 0.001 RCHV (100,000 base units)  
**Impact:** No fee prioritization  
**Future:** Dynamic fee market in mainnet

### Simple VDF (Iterated SHA256)

**Status:** As designed for testnet  
**Limitation:** Not production-grade VDF  
**Impact:** Testnet only - not secure enough for mainnet  
**Future:** Wesolowski or Pietrzak VDF for mainnet

---

## Reporting Issues

Found a bug? Report it:
- **GitHub Issues:** https://github.com/ArchivasNetwork/archivas/issues
- **Include:** Version, logs, steps to reproduce
- **Label:** Use `bug` label

---

**Last Updated:** November 4, 2025

