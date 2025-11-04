# Consensus Mechanism

How Archivas achieves distributed agreement.

## Proof-of-Space-and-Time

**Core idea:** Blocks require BOTH disk space proof AND time proof.

### Challenge Generation
```
Challenge = SHA256(VDF_output || height)
```

### Proof Generation  
```
Quality = SHA256(Challenge || Plot_hash)
```

### Win Condition
```
IF Quality < Difficulty THEN farmer wins
```

## Difficulty Adjustment

Target: 20-30 second block times

**Algorithm:** Simplified EMA
- Too fast → difficulty increases
- Too slow → difficulty decreases
- Bounds: 1M minimum, no maximum

## Security Properties

- **Grinding-resistant:** VDF prevents precomputation
- **Sybil-resistant:** Multiple identities don't help
- **51% attack:** Requires 51% of storage
- **Nothing-at-stake:** Not applicable (PoSpace)

---

**Next:** [Block Structure](block-structure.md)
