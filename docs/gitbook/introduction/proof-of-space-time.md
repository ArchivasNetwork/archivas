# What is Proof-of-Space-and-Time?

Proof-of-Space-and-Time (PoST) is a consensus mechanism that combines **disk space allocation** with **sequential time proofs** to secure a blockchain.

---

## The Problem with Traditional Consensus

### Proof-of-Work (Bitcoin, Ethereum pre-merge)

**How it works:** Miners race to find a hash below a target by trying billions of combinations.

**Problems:**
- ❌ **Energy waste:** Bitcoin uses ~150 TWh/year (more than Argentina)
- ❌ **Specialized hardware:** ASIC/GPU arms race
- ❌ **Centralization:** Mining pools control majority

### Proof-of-Stake (Ethereum, Cardano)

**How it works:** Validators stake capital to participate in consensus.

**Problems:**
- ❌ **Capital barriers:** Requires 32 ETH (~$100K) to run Ethereum validator
- ❌ **"Rich get richer":** Rewards proportional to stake
- ❌ **Plutocracy risk:** Wealth determines governance

---

## Proof-of-Space-and-Time Solution

PoST replaces energy and capital with **disk space** and **time**.

### Part 1: Proof-of-Space

**Concept:** Use disk space to store precomputed data (plots) that prove you allocated storage.

**How it works:**
1. **Plot Creation:** Farmers fill files with precomputed hash tables
2. **Challenge:** Network issues a random challenge each block
3. **Proof Search:** Farmers scan their plots for the best matching proof
4. **Quality Check:** The proof with lowest "quality" wins the block
5. **Reward:** Winning farmer earns block reward (20 RCHV)

**Advantages:**
- ✅ Uses commodity hardware (any HDD/SSD)
- ✅ Energy efficient (just disk I/O)
- ✅ Permissionless (anyone with storage can farm)
- ✅ Fair (more space = more chances, but diminishing returns)

**Key Insight:** You CAN'T pre-compute proofs because the challenge changes every block.

---

### Part 2: Verifiable Delay Function (VDF)

**Problem:** Without VDF, farmers could "grind" by trying many challenges in parallel.

**Solution:** VDF is a function that takes **real time** to compute and cannot be parallelized.

**How it works:**
1. **Timelord** continuously computes VDF (iterated hashing)
2. VDF output is used to generate the next challenge
3. **Sequential nature:** Cannot skip ahead or parallelize
4. **Verification:** Fast to verify, slow to compute

**Properties:**
- Takes fixed time (e.g., 10 seconds per iteration)
- Cannot be sped up with more hardware
- Provides temporal ordering
- Prevents grinding attacks

**Archivas Implementation:**
- Devnet: Iterated SHA-256 (simple, fast verification)
- Future: Wesolowski or Pietrzak VDF (more secure)

---

## How They Work Together

### Block Production Cycle

```
1. Timelord computes VDF
   ↓
2. VDF output → new challenge
   ↓
3. Farmers search plots for proofs
   ↓
4. Best proof wins block
   ↓
5. Block added to chain
   ↓
6. Repeat from step 1
```

**Security Properties:**
- **Space-hard:** Can't fake disk allocation
- **Time-hard:** Can't precompute or grind
- **Sybil-resistant:** Multiple identities don't help (same disk)
- **Fair:** Rewards proportional to space contributed

---

## Real-World Analogy

**PoW (Bitcoin):** Like a lottery where tickets are bought with electricity. The more you spend, the more tickets.

**PoS (Ethereum):** Like a lottery where tickets are bought with money. The richer you are, the more tickets.

**PoST (Archivas, Chia):** Like a lottery where tickets are based on storage space you prove you have. More disk = more tickets, but:
- Disks are reusable (not consumed like energy)
- Disk is cheap compared to capital lockup
- Anyone can participate

**VDF:** The lottery drawing happens at a fixed time. You can't rush it or skip ahead.

---

## Proven by Chia Network

**Chia Network** pioneered PoST consensus and has proven it works at scale:
- **Launched:** May 2021
- **Exabytes:** 10+ EB of storage farmed (at peak)
- **Blocks:** Millions produced without consensus failures
- **Security:** No successful 51% attacks

**Archivas uses the same fundamental model** but with:
- Simpler implementation (Go vs Python)
- Modular design (easier to extend)
- Different economic model
- Open development process

---

## Why PoST Matters

### Environmental Impact
- **Bitcoin:** ~150 TWh/year (0.5% of global electricity)
- **Archivas:** ~1 kWh/year per farmer (disk I/O only)
- **Savings:** >99.9% energy reduction

### Accessibility
- **PoW:** Requires ASIC/GPU farms (~$10K-$1M+ investment)
- **PoS:** Requires capital lockup (32 ETH = ~$100K)
- **PoST:** Requires disk space (~$50 for 1TB HDD)

### Decentralization
- **PoW:** Centralized in mining pools (top 4 control >50% Bitcoin)
- **PoS:** Centralized in exchanges (Coinbase, Binance stake for users)
- **PoST:** Geographic distribution (disk is everywhere)

---

## Technical Details

### Proof-of-Space Algorithm
1. Generate plot ID: `hash(farmer_pubkey)`
2. Fill plot with hashes: `plot[i] = hash(plot_id || i)`
3. For each challenge: compute `quality = hash(challenge || plot[best_index])`
4. If `quality < difficulty_target`, farmer wins

### VDF (Simplified)
```
seed₀ = genesis_seed
for i in 1..iterations:
    seed_i = SHA256(seed_i₋₁)
output = seed_iterations
```

**Time:** iterations × hash_time (cannot parallelize)

---

## Learn More

- **Chia Network:** https://www.chia.net/greenpaper/
- **VDF Research:** https://vdfresearch.org/
- **Archivas Consensus:** [consensus.md](../architecture/consensus.md)
- **Archivas Architecture:** [overview.md](../architecture/overview.md)

---

**Next:** Learn [Why Archivas?](why-archivas.md) and how it improves on existing PoST implementations.

