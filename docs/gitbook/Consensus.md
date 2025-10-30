# Consensus Mechanism

## Proof-of-Space-and-Time

Archivas combines two complementary proofs to achieve secure, fair, and energy-efficient consensus:

1. **Proof-of-Space (PoSpace)** - Lottery based on disk space
2. **Verifiable Delay Function (VDF)** - Sequential time proofs

Both are required for every block.

---

## Proof-of-Space

### Overview

Farmers commit disk space by creating "plots" - large files filled with precomputed cryptographic hashes. When a block is needed, farmers search their plots for the best "proof" that satisfies the current challenge.

### Plot Generation

**Process:**
```
1. Generate farmer keypair (secp256k1)
2. Compute plot ID: plotID = SHA256(farmerPubKey)
3. For each index i in [0, 2^k):
   hash[i] = SHA256(SHA256(farmerPubKey || plotID || i))
4. Write header + hashes to .arcv file
```

**Plot Sizes:**
- k=16: 65,536 hashes (~2 MB)
- k=20: 1,048,576 hashes (~32 MB)
- k=24: 16,777,216 hashes (~512 MB)
- k=28: 268,435,456 hashes (~8 GB)

### Challenge-Response

**Each block requires:**
1. **Challenge:** `C = H(VDF_output || height)`
2. **Quality:** For each plot entry: `Q_i = H(C || hash[i])`
3. **Best Quality:** `Q_best = min(Q_0, Q_1, ..., Q_n)`
4. **Win Condition:** `Q_best < difficulty_target`

**Lower quality = Better proof**

### Verification

**Nodes verify:**
1. ✅ Challenge matches block header
2. ✅ Plot hash recomputes correctly: `H(farmerPubKey || plotID || index)`
3. ✅ Quality recomputes correctly: `H(challenge || plotHash)`
4. ✅ Quality < difficulty

**Deterministic and fast** (<1ms per proof)

---

## Verifiable Delay Functions

### Purpose

VDFs provide temporal security:
- **Grinding Resistance:** Cannot precompute alternative timelines
- **Temporal Ordering:** Blocks have provable time sequence
- **Fairness:** Challenge is unpredictable until VDF computed

### Algorithm (Devnet)

**Iterated SHA-256:**
```
y_0 = seed
y_1 = SHA256(y_0)
y_2 = SHA256(y_1)
...
y_n = SHA256(y_{n-1})
```

**Properties:**
- Inherently sequential (cannot skip iterations)
- Deterministic (same seed → same output)
- Verifiable (anyone can recompute)

**Parameters:**
- Step size: 500-1000 iterations/second
- Typical block: 1,000-5,000 iterations

### Production VDF (Future)

**Wesolowski/Pietrzak:**
- RSA or class group operations
- Succinct proofs (O(log n) verification)
- Quantum-resistant assumptions
- Faster iteration rates

---

## Adaptive Difficulty

### Goal

Maintain consistent block times (~20 seconds average) as network hashrate changes.

### Algorithm

**Inputs:**
- Recent block timestamps (10-block window)
- Current difficulty

**Computation:**
```go
avgBlockTime = (time[n] - time[n-10]) / 10
ratio = avgBlockTime / targetBlockTime  // 20 seconds

// Clamp to prevent wild swings
ratio = clamp(ratio, 0.5, 2.0)

newDifficulty = currentDifficulty * ratio

// Bounds check
newDifficulty = clamp(newDifficulty, 2^40, 2^60)
```

**Properties:**
- Responds to network changes
- Smooth adjustments (max 2x per window)
- Bounded range prevents extremes

### Difficulty Target

**Interpretation:**
```
target = MaxUint64 / difficulty
isWin = quality <= target
```

Higher difficulty = Smaller target = Harder to win

---

## Block Structure

### Block Header

```go
type Block struct {
    Height        uint64      // Block number
    TimestampUnix int64       // Unix timestamp
    PrevHash      [32]byte    // Previous block hash
    Difficulty    uint64      // Difficulty when mined
    Challenge     [32]byte    // Challenge used for PoSpace
    Txs           []Transaction
    Proof         *Proof      // PoSpace proof
    FarmerAddr    string      // Reward recipient
}
```

**Critical fields for sync:**
- `Difficulty` - Difficulty target when block was mined
- `Challenge` - Challenge used to find PoSpace proof

**These enable deterministic verification independent of node's current VDF state!**

### PoSpace Proof

```go
type Proof struct {
    Challenge    [32]byte  // Challenge (matches block header)
    PlotID       [32]byte  // Which plot
    Index        uint64    // Which hash in plot
    Hash         [32]byte  // The plot hash
    Quality      uint64    // Computed quality
    FarmerPubKey [33]byte  // Farmer's public key
}
```

---

## Block Validation

### Producer Validation (Real-time)

When node receives block from farmer:

1. ✅ **Height:** `block.Height == currentHeight + 1`
2. ✅ **PrevHash:** `block.PrevHash == hash(tipBlock)`
3. ✅ **Difficulty:** `block.Difficulty == currentDifficulty`
4. ✅ **Challenge:** `block.Challenge == currentChallenge`
5. ✅ **PoSpace:** `VerifyProof(block.Proof, block.Challenge, block.Difficulty)`
6. ✅ **Transactions:** All signatures valid, balances sufficient
7. ✅ **Apply:** Update state, persist block, gossip to peers

### Sync Validation (Historical)

When syncing from peer:

1. ✅ **Height:** Sequential (n+1)
2. ✅ **PrevHash:** Chain linkage
3. ✅ **PoSpace:** Verify using `block.Challenge` and `block.Difficulty` (NOT current state!)
4. ✅ **Difficulty:** (Optional) Recompute from history and verify match
5. ✅ **Apply:** Update state, persist block

**Key insight:** Blocks are self-contained with their own validation parameters!

---

## Chain Selection

**Current (Devnet):**
- Longest chain rule
- Chain with greatest height wins

**Future (Production):**
- Cumulative difficulty: `Σ(block.Difficulty)`
- Heaviest chain wins
- Finality thresholds

---

## Economic Model

### Block Rewards

**Current:**
- 20 RCHV per block
- No transaction fees (burned)

**Future:**
- Halving schedule (TBD)
- Fee market
- Block reward + fees to farmer

### Token Supply

**Devnet:**
- Genesis: 1B RCHV (test allocation)
- Block rewards: 20 RCHV/block
- Unlimited (testnet)

**Mainnet:**
- TBD (likely Bitcoin-like halving)
- Terminal supply: ~42M RCHV (example)
- Foundation allocation: 10% (example)

---

## Consensus Parameters

### Devnet v3

```
Chain ID: 1616
Network ID: archivas-devnet-v3
Block Time: ~20 seconds
Block Reward: 20 RCHV
Initial Difficulty: 2^50 (1,125,899,906,842,624)
Difficulty Range: 2^40 to 2^60
Adjustment Window: 10 blocks
Max Adjustment: 2x per window
```

### VDF Parameters

```
Algorithm: Iterated SHA-256
Step Size: 500 iterations/second
Typical Iterations: 1,000-5,000 per block
Verification: O(n) recomputation
```

---

## Security Analysis

### Attack Scenarios

**Grinding Attack:**
- **Vector:** Try many alternative blocks to find favorable chain
- **Defense:** VDF takes real time, cannot be parallelized
- **Cost:** Must redo all VDF work for alternative timeline

**51% Attack:**
- **Vector:** Control majority of disk space
- **Defense:** Distributed plot ownership
- **Cost:** Proportional to network size

**Nothing-at-Stake:**
- **Vector:** Bet on multiple forks simultaneously  
- **Defense:** Disk space is committed resource
- **Cost:** Cannot use same plots for multiple chains

**Sybil Attack:**
- **Vector:** Create many fake identities
- **Defense:** Genesis hash validation, peer scoring
- **Cost:** Still need real disk space to win blocks

### Assumptions

**Trust:**
- Cryptographic primitives (SHA-256, secp256k1)
- Majority of disk space is honest
- VDF computation is sequential

**Network:**
- Eventual connectivity between honest nodes
- Bounded message delays
- Some peers maintain full history

---

**Next:** [Testnet Guide →](Testnet-Guide.md)  
**Back:** [← Architecture](Architecture.md)

