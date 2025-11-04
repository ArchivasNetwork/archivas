# Architecture Overview

Archivas system architecture and components.

## System Diagram

```
                     ┌─────────────┐
                     │   Farmer    │
                     │  (Plots)    │
                     └──────┬──────┘
                            │ Proofs
                            ▼
┌──────────┐        ┌─────────────┐        ┌──────────────┐
│ Timelord │──VDF──▶│    Node     │◀──P2P──│  Other Nodes │
│  (VDF)   │        │ (Consensus) │        │              │
└──────────┘        └──────┬──────┘        └──────────────┘
                           │
                           │ RPC/HTTPS
                           ▼
                    ┌─────────────┐
                    │   Clients   │
                    │ (SDK, Apps) │
                    └─────────────┘
```

## Core Components

### Node
- Block validation
- State management (accounts, balances)
- P2P networking
- RPC server
- Storage engine (BadgerDB)

### Farmer
- Plot scanning
- Proof generation
- Block submission
- Reward collection

### Timelord (Optional)
- VDF computation
- Challenge generation
- Temporal ordering

### SDK/Clients
- Wallet management
- Transaction signing
- RPC queries

## Data Flow

1. Timelord computes VDF → Challenge
2. Farmers scan plots → Find proofs
3. Best proof → Submit to Node
4. Node validates → Add block
5. Broadcast to Peers → Sync network
6. Farmer receives → 20 RCHV reward

---

**Details:** [Consensus](consensus.md) | [Block Structure](block-structure.md)
