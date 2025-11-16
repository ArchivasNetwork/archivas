# Betanet Phase 1: Foundation

This document describes Phase 1 of the Betanet upgrade: Network Profiles & Dual Address System.

## Overview

Phase 1 lays the groundwork for Betanet as a clean, EVM-integrated testnet. It introduces:

1. **Network Profile System** - Structured network configuration
2. **Dual Address Model** - Internal EVM addresses + External Bech32 addresses
3. **Betanet Genesis** - New genesis structure with EVM support
4. **CLI Integration** - `--network` flag for network selection

## Network Profiles

### Implementation

Network profiles are defined in `network/profile.go`:

```go
type NetworkProfile struct {
    Name             string        // e.g., "betanet", "devnet-legacy"
    ChainID          string        // e.g., "archivas-betanet-1"
    NetworkID        uint64        // Numeric ID for P2P
    ProtocolVersion  int           // Protocol version
    GenesisPath      string        // Path to genesis file
    DefaultSeeds     []string      // Seed node addresses
    DefaultRPCPort   int           // Default RPC port
    DefaultP2PPort   int           // Default P2P port
    TargetBlockTime  time.Duration // Target block time
    InitialDifficulty uint64       // Starting difficulty
    Bech32Prefix     string        // Address prefix
}
```

### Available Networks

#### Betanet (Default)

- **Name**: `betanet`
- **Chain ID**: `archivas-betanet-1`
- **Network ID**: `102`
- **Protocol Version**: `2`
- **Genesis**: `configs/genesis-betanet.json`
- **Seeds**: 
  - `seed1.betanet.archivas.ai:9090`
  - `seed2.betanet.archivas.ai:9090`
- **RPC Port**: `8545` (Ethereum-compatible)
- **Bech32 Prefix**: `arcv`

#### Devnet Legacy

- **Name**: `devnet-legacy`
- **Chain ID**: `archivas-devnet-1`
- **Network ID**: `1`
- **Protocol Version**: `1`
- **Genesis**: `genesis/devnet.genesis.json`
- **Seeds**: 
  - `seed.archivas.ai:9090`
  - `seed2.archivas.ai:30303`
- **RPC Port**: `8080`
- **Bech32 Prefix**: `arcv`

### Usage

```bash
# Run Betanet node (default)
archivas-node

# Explicitly select Betanet
archivas-node --network betanet

# Run devnet-legacy node
archivas-node --network devnet-legacy

# Override defaults
archivas-node --network betanet --rpc :9545 --p2p :10090
```

## Dual Address System

### Design

Archivas implements a **dual address model** similar to Tron:

1. **Internal (Canonical)**: 20-byte EVM addresses (`0x...`)
   - Used by EVM execution engine
   - Stored in state trie
   - Used in EVM transactions

2. **External (User-Facing)**: Bech32 encoded addresses (`arcv1...`)
   - Human-readable with checksums
   - Used in CLI, wallets, explorers
   - Cosmos-style format

### Implementation

Located in `address/` package:

#### EVMAddress Type

```go
// 20-byte canonical address
type EVMAddress [20]byte

// Create from hex
addr, _ := EVMAddressFromHex("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")

// Convert to hex
hexStr := addr.Hex() // "0x742d35cc6634c0532925a3b844bc9e7595f0beb0"

// Convert to Bech32 ARCV
arcvStr, _ := EncodeARCVAddress(addr, "arcv")
// arcvStr = "arcv1ws4rt3nxxpq9jv59wc9sghrjw4v4p0tmq..."
```

#### Address Parsing

```go
// Parse either format (0x or arcv1)
addr1, _ := ParseAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "arcv")
addr2, _ := ParseAddress("arcv1ws4rt3nxxpq9jv59wc9sghrjw4v4p0tmq...", "arcv")

// addr1 == addr2 (same canonical address)
```

#### Bech32 Encoding

Custom Bech32 implementation in `address/bech32.go`:

- No external dependencies
- BIP 173 compliant
- Checksum validation
- 5-bit to 8-bit conversion

### RPC Behavior

- **`eth_*` RPC endpoints**: Use `0x` hex addresses (Ethereum-compatible)
- **`arcv_*` RPC endpoints** (future): Accept/return `arcv1...` addresses
- All addresses are internally stored as `EVMAddress`

### Examples

```go
import "github.com/ArchivasNetwork/archivas/address"

// Create from hex
addr := address.MustParse("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0", "arcv")

// Encode to ARCV
arcvAddr, _ := address.EncodeARCVAddress(addr, "arcv")
fmt.Println(arcvAddr) // arcv1ws4rt3nxxpq9jv59wc9sghrjw4v4p0tmq...

// Decode from ARCV
decoded, _ := address.DecodeARCVAddress(arcvAddr, "arcv")
fmt.Println(decoded.Hex()) // 0x742d35cc6634c0532925a3b844bc9e7595f0beb0

// Parse either format
parsed1, _ := address.ParseAddress("0x742d35Cc...", "arcv")
parsed2, _ := address.ParseAddress("arcv1ws4rt3n...", "arcv")
// Both work!
```

## Betanet Genesis

Located at `configs/genesis-betanet.json`:

### Structure

```json
{
  "chain_name": "Archivas Betanet",
  "chain_id": "archivas-betanet-1",
  "network_id": 102,
  "protocol_version": 2,
  "genesis_time": "2025-01-01T00:00:00Z",
  "evm_config": {
    "chain_id": 102,
    "homestead_block": 0,
    "eip150_block": 0,
    "eip155_block": 0,
    ...
  },
  "initial_state": {
    "state_root": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "receipts_root": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "accounts": []
  },
  "consensus_params": {
    "post": {
      "k_size": 32,
      "min_plot_size": 107374182400,
      "challenge_delay": 3,
      "signage_point_interval": 10
    }
  }
}
```

### Key Features

1. **Empty State**: Clean start, no pre-allocated balances
2. **EVM Config**: All EIPs enabled from block 0
3. **PoST Consensus**: K32 plots, standard parameters
4. **State Roots**: Empty trie hashes for clean genesis

### Loading

```go
import "github.com/ArchivasNetwork/archivas/network"

// Load genesis file
genesis, err := network.LoadGenesis("configs/genesis-betanet.json")

// Validate against profile
profile, _ := network.GetProfile("betanet")
err = network.MatchesProfile(genesis, profile)

// Compute genesis hash
genesisHash, _ := network.ComputeGenesisHash(genesis)
```

## CLI Changes

### New Flags

```
--network <name>   Network to join (betanet, devnet-legacy) [default: betanet]
```

### Flag Behavior

- **`--rpc`**: Defaults to network profile's `DefaultRPCPort`
- **`--p2p`**: Defaults to network profile's `DefaultP2PPort`
- **`--genesis`**: Defaults to network profile's `GenesisPath`
- **`--network-id`**: Defaults to network profile's `ChainID`

### Startup Sequence

1. Parse `--network` flag
2. Load network profile from registry
3. Apply profile defaults to flags
4. Validate genesis file
5. Start node with profile configuration

### Example Output

```
[network] Loading network profile: betanet
[network] Network: betanet (chain-id: archivas-betanet-1, network-id: 102, protocol: v2)
🚀 Archivas Betanet Node running…

🔧 Configuration:
   Network: betanet
   Chain ID: archivas-betanet-1
   RPC:  0.0.0.0:8545
   P2P:  0.0.0.0:9090
   DB:   ./data
   Mode: PoSpace only
```

## Testing

### Address Tests

```bash
# Run all address tests
go test ./address/... -v

# Should pass:
# - EVMAddressFromHex (with/without 0x)
# - Bech32 roundtrip encoding
# - ParseAddress (both formats)
# - Zero address handling
# - Edge cases (all zeros, all 0xFF)
# - Wrong HRP rejection
```

### Build Test

```bash
# Build node with Phase 1 changes
go build -o archivas-node ./cmd/archivas-node/main.go

# Test help
./archivas-node help

# Test network selection
./archivas-node --network betanet --help
./archivas-node --network devnet-legacy --help
```

## Next Steps: Phase 2

Phase 2 will implement:

1. **Block Header Extensions**: Add `StateRoot`, `ReceiptsRoot`, `GasUsed`, `GasLimit`
2. **EVM Engine Integration**: Wrap go-ethereum's EVM
3. **StateDB Implementation**: Merkle Patricia Trie for EVM state
4. **Transaction Types**: EVM deploy and call transactions

## Migration Path

### For Developers

- Update code to use `address.EVMAddress` instead of raw byte arrays
- Use `address.ParseAddress()` for CLI input
- Use `address.EncodeARCVAddress()` for display output

### For Node Operators

- Current devnet becomes `devnet-legacy`
- Use `--network devnet-legacy` to continue on old chain
- New deployments default to `betanet`

### For Farmers

- Betanet will require re-plotting (different network ID)
- Devnet-legacy remains available for continuity
- Address format accepts both `0x` and `arcv1` styles

## Files Created/Modified

### New Files

- `address/address.go` - Dual address implementation
- `address/address_test.go` - Comprehensive tests
- `address/bech32.go` - Bech32 encoding/decoding
- `network/profile.go` - Network profile system
- `network/genesis.go` - Genesis loader and validator
- `configs/genesis-betanet.json` - Betanet genesis file
- `docs/BETANET_PHASE1.md` - This document

### Modified Files

- `cmd/archivas-node/main.go` - Added `--network` flag, profile loading

## Summary

✅ **Completed**:
- Network profile system with Betanet & Devnet Legacy
- Dual address model (0x EVM + arcv1 Bech32)
- Betanet genesis structure
- CLI `--network` flag
- Comprehensive tests
- Documentation

🚧 **Next (Phase 2)**:
- Block header EVM extensions
- EVM engine integration
- StateDB and trie implementation

---

**Date**: 2025-11-16  
**Version**: Phase 1 Foundation  
**Status**: Complete ✅

