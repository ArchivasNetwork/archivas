# IBD (Initial Block Download) Fix Specification

**Status:** Broken - Requires Complete Reimplementation  
**Priority:** Critical for Public Launch  
**Tracked:** GitHub Issue #TBD  

---

## Problem Statement

**Current Behavior:**
- New nodes connect to seed (Server A at 86,000+ blocks)
- Node gets stuck at height 3
- Logs show "caught up to tip at height 3" even though peer is 86,000 blocks ahead
- PoSpace verification fails with backwards comparison errors
- TipHeight=0 returned in some responses causing underflow

**Required Behavior:**
- New node detects it's 86,000 blocks behind
- Automatically downloads blocks in batches (512 at a time)
- Validates and persists each batch
- Shows progress: "[IBD] height 5000/86000 (5.8% complete)"
- Catches up to network tip
- Switches to normal block-by-block sync

**Why It Matters:**
- Users cannot manually copy data directories
- Public testnet requires automatic sync
- Core feature for any blockchain

---

## Root Causes Found

### 1. TipHeight=0 Bug (Fixed in v1.2.0)
**Issue:** Backpressure case returned `TipHeight: 0`  
**Impact:** Client calculated `18446744073709551613 blocks behind` (underflow)  
**Fix Applied:** Return real tip from `GetStatus()`  
**Status:** ✅ Fixed but not sufficient

### 2. Blocks Not Being Applied
**Issue:** Blocks received but validation fails  
**Error:** "quality 293581 does not meet difficulty 1875000"  
**Analysis:** 293,581 < 1,875,000 should PASS, but fails  
**Cause:** Unknown - PoSpace code shows correct `quality <= difficulty`  
**Status:** ❌ Unresolved

### 3. OnBlocksRangeRequest Returns Empty
**Issue:** Method returns 0 blocks despite 86,000 available  
**Analysis:** Code looks correct, but not executing properly  
**Status:** ❌ Needs investigation

---

## Required Implementation

### A. Server Side (Seed Node)

#### HTTP Endpoint for Block Ranges

```go
// GET /blocks/since/<height>?limit=N
func (s *Server) handleBlocksSince(w http.ResponseWriter, r *http.Request) {
    // Parse height from URL
    heightStr := r.URL.Path[len("/blocks/since/"):]
    height, _ := strconv.ParseUint(heightStr, 10, 64)
    
    // Parse limit (default 512, max 1000)
    limitStr := r.URL.Query().Get("limit")
    limit := 512
    if limitStr != "" {
        limit, _ = strconv.Atoi(limitStr)
    }
    if limit > 1000 {
        limit = 1000
    }
    
    // Get current tip
    tipHeight := ns.CurrentHeight
    
    // Get blocks from storage
    blocks := []Block{}
    for h := height; h < height+uint64(limit) && h <= tipHeight; h++ {
        block, err := ns.GetBlockByHeight(h)
        if err != nil {
            break // Stop at first missing block
        }
        blocks = append(blocks, block)
    }
    
    // Response
    response := map[string]interface{}{
        "tipHeight": tipHeight,  // MUST be real tip, never 0
        "blocks": blocks,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    
    log.Printf("[blocks-range] served %d blocks from %d (tip=%d) to %s",
        len(blocks), height, tipHeight, r.RemoteAddr)
}
```

**Critical:**
- `tipHeight` MUST be real chain tip
- NEVER return `tipHeight: 0` unless chain is actually empty
- Blocks must be sequential and ordered
- Load from disk if not in memory

#### P2P Message Handler

Fix `handleRequestBlocks` in `p2p/ibd.go`:
- Always return real tip (even in backpressure case) ✅ Done
- Actually serve blocks from `OnBlocksRangeRequest`
- Log clearly what's being sent

### B. Client Side (Syncing Node)

#### IBD Detection and Trigger

```go
// On startup, after loading local state:
func (n *Node) checkAndStartIBD() {
    localTip := n.CurrentHeight
    
    // Get remote tip from seed
    remoteTip, err := n.fetchRemoteTip()
    if err != nil {
        log.Printf("[IBD] cannot fetch remote tip: %v", err)
        return
    }
    
    gap := remoteTip - localTip
    
    if gap > 200 {
        log.Printf("[IBD] Detected large gap: local=%d remote=%d (%.1f%% behind)",
            localTip, remoteTip, float64(localTip)*100/float64(remoteTip))
        
        // Enter IBD mode
        if err := n.runIBD(remoteTip); err != nil {
            log.Printf("[IBD] failed: %v", err)
        }
    }
}
```

#### IBD Main Loop

```go
func (n *Node) runIBD(remoteTip uint64) error {
    const batchSize = 512
    
    for {
        localTip := n.CurrentHeight
        gap := remoteTip - localTip
        
        // Exit IBD when close to tip
        if gap <= 50 {
            log.Printf("[IBD] Complete! Height: %d (gap: %d)", localTip, gap)
            return nil
        }
        
        // Request next batch
        fromHeight := localTip + 1
        blocks, newRemoteTip, err := n.fetchBlockBatch(fromHeight, batchSize)
        if err != nil {
            log.Printf("[IBD] fetch error: %v, retrying...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        
        // Update remote tip (may have grown)
        remoteTip = newRemoteTip
        
        // Apply blocks
        for _, block := range blocks {
            if err := n.applyBlock(block); err != nil {
                return fmt.Errorf("IBD apply block %d failed: %w", block.Height, err)
            }
        }
        
        // Progress logging
        newTip := n.CurrentHeight
        pct := float64(newTip) * 100 / float64(remoteTip)
        log.Printf("[IBD] Progress: %d/%d (%.1f%% complete, %d remaining)",
            newTip, remoteTip, pct, remoteTip-newTip)
    }
}
```

#### HTTP Fetch

```go
func (n *Node) fetchBlockBatch(fromHeight uint64, limit int) (blocks []Block, remoteTip uint64, err error) {
    url := fmt.Sprintf("%s/blocks/since/%d?limit=%d", n.seedURL, fromHeight, limit)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, 0, err
    }
    defer resp.Body.Close()
    
    var result struct {
        TipHeight uint64 `json:"tipHeight"`
        Blocks []Block `json:"blocks"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, 0, err
    }
    
    // Validate tip
    if result.TipHeight == 0 {
        return nil, 0, fmt.Errorf("peer returned invalid tipHeight=0")
    }
    
    return result.Blocks, result.TipHeight, nil
}
```

### C. Block Validation During IBD

**Each block must be validated:**
1. ✅ Height is sequential (localTip + 1)
2. ✅ PrevHash matches previous block hash
3. ✅ Difficulty is reasonable (within adjustment bounds)
4. ✅ Proof-of-Space verification
5. ✅ Transaction signatures (if any)

**BUT:** Be careful with difficulty validation on old blocks. If the chain adjusted difficulty over time, historical blocks may have different difficulty than current. **Use the difficulty stored IN the block**, not current consensus difficulty.

### D. Persistence

**Write blocks immediately:**
```go
for _, block := range batch {
    // Validate
    if err := validateBlock(block); err != nil {
        return err
    }
    
    // Persist to disk IMMEDIATELY
    if err := ns.BlockStore.SaveBlock(block.Height, block); err != nil {
        return err
    }
    
    // Update in-memory state
    ns.Chain = append(ns.Chain, block)
    ns.CurrentHeight = block.Height
    
    // Also save tip to metadata store
    ns.MetaStore.SaveTipHeight(block.Height)
}
```

This way, if IBD is interrupted, we resume from the last persisted height.

---

## Testing

### Test 1: Clean Sync

```bash
# On Server B:
rm -rf data/*
./archivas-node --rpc :8080 --p2p :9090 --seed 57.129.148.132:9090

# Expected logs:
# [IBD] Detected large gap: local=0 remote=86000
# [IBD] Progress: 512/86000 (0.6% complete)
# [IBD] Progress: 5120/86000 (5.9% complete)
# [IBD] Progress: 50000/86000 (58.1% complete)
# [IBD] Complete! Height: 85950 (gap: 50)
# [sync] Following head at height 85950
```

### Test 2: Interrupted Sync

```bash
# Start IBD, kill at 50% (height 43000)
pkill archivas-node

# Restart
./archivas-node ...

# Expected:
# [IBD] Resuming from height 43000
# [IBD] Progress: 43512/86000 (50.6% complete)
# (continues from where it left off)
```

### Test 3: Moving Target

While IBD is running, new blocks are being produced on Server A. IBD must:
- Update remoteTip dynamically
- Not stop at the "old" tip
- Continue until actually caught up

---

## Current Files to Modify

1. **rpc/farming.go** - Add `/blocks/since/<height>` handler
2. **cmd/archivas-node/main.go** - Add IBD client logic
3. **p2p/ibd.go** - Fix existing IBD or replace
4. **p2p/ibd_sync.go** - Fix or remove

---

## Success Criteria

✅ Fresh node syncs from 0 → 86,000 automatically  
✅ No manual data copying required  
✅ Progress is visible in logs  
✅ Can resume if interrupted  
✅ Works for any future network (100K, 1M blocks)  
✅ Handles moving target (tip growing during sync)  

---

## Related Issues

- Server B stuck at height 3 (this spec)
- TipHeight=0 underflow (fixed)
- PoSpace verification failing on valid blocks (investigate)

---

**This spec represents the complete IBD solution for Archivas.**

