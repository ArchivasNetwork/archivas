# Betanet Phase 3: P2P Identity & ETH RPC

This document describes Phase 3 of the Betanet upgrade: P2P Identity Enforcement and Ethereum-compatible RPC endpoints.

## Overview

Phase 3 adds network security and Ethereum compatibility:

1. **Snapshot Manifest Verification** - Ensures snapshots match the target network
2. **P2P Identity Enforcement** - Prevents incompatible nodes from connecting
3. **ETH RPC Endpoints** - Ethereum JSON-RPC compatibility
4. **ARCV RPC Endpoints** - Archivas-specific address conversion

## 1. Snapshot Manifest Verification

### Enhanced Manifest Structure

```go
type Manifest struct {
    // Basic info
    Network        string `json:"network"`
    Height         uint64 `json:"height"`
    Hash           string `json:"hash"`
    SnapshotURL    string `json:"snapshot_url"`
    ChecksumSHA256 string `json:"checksum_sha256"`
    
    // Phase 3: Chain identity fields
    ChainID         string `json:"chain_id"`          // "archivas-betanet-1"
    NetworkID       uint64 `json:"network_id"`        // 102
    ProtocolVersion int    `json:"protocol_version"`  // 2
    StateRoot       string `json:"state_root"`        // State root at height
    GenesisHash     string `json:"genesis_hash"`      // Genesis block hash
    
    // Metadata
    CreatedAt string `json:"created_at"`
    CreatedBy string `json:"created_by"`
}
```

### Verification Process

```go
func VerifyManifest(
    manifest *Manifest,
    profile *network.NetworkProfile,
    genesisHash string,
) error
```

**Checks**:
1. ✅ Network name matches
2. ✅ Chain ID matches
3. ✅ Network ID matches
4. ✅ Protocol version matches
5. ✅ Genesis hash matches
6. ✅ Required fields present

**Example**:
```go
// Load manifest
manifest, err := fetchManifest(manifestURL)

// Verify against Betanet profile
profile, _ := network.GetProfile("betanet")
err = VerifyManifest(manifest, profile, genesisHash)

// If error: snapshot is from wrong network/chain
if err != nil {
    log.Fatalf("Incompatible snapshot: %v", err)
}
```

### Bootstrap with Verification

```go
opts := BootstrapOptions{
    ManifestURL: "https://seed2.archivas.ai/betanet/latest.json",
    DBPath:      "./data",
    
    // Phase 3: Identity verification
    NetworkProfile: profile,
    GenesisHash:    genesisHashHex,
}

metadata, err := Bootstrap(opts)
```

**Output**:
```
[bootstrap] Fetching manifest...
[bootstrap] Manifest info:
  Network:  betanet
  Chain ID: archivas-betanet-1
  Network ID: 102
  Protocol: v2
  Height:   1250000
  Hash:     a1b2c3d4...
[bootstrap] Verifying manifest identity...
[bootstrap] ✓ Manifest verification passed
[bootstrap] Downloading snapshot...
[bootstrap] ✓ Checksum verified
[bootstrap] Importing snapshot...
[bootstrap] ✓ Bootstrap complete!
```

### Rejection Scenarios

**Wrong Network**:
```
Error: manifest verification failed: network name mismatch
Expected: betanet
Got: devnet-legacy
```

**Wrong Chain ID**:
```
Error: manifest verification failed: chain ID mismatch
Expected: archivas-betanet-1
Got: archivas-devnet-1
```

**Wrong Protocol Version**:
```
Error: manifest verification failed: protocol version mismatch
Expected: 2
Got: 1
```

**Wrong Genesis Hash**:
```
Error: manifest verification failed: genesis hash mismatch
This snapshot is from a different chain
```

## 2. P2P Identity Enforcement

### Enhanced Handshake Message

```go
type HandshakeMessage struct {
    // Phase 3: Chain identity (strict enforcement)
    GenesisHash     [32]byte // Genesis block hash
    ChainID         string   // "archivas-betanet-1"
    NetworkID       uint64   // 102
    ProtocolVersion int      // 2
    
    // Legacy compatibility
    NetworkIDLegacy    string // "102" (string)
    ProtocolVersionStr string // "v2" (string)
    DifficultyParamsID string
    
    // Informational
    NodeVersion string // "v1.0.0"
    NodeName    string // "my-node"
}
```

### Handshake Verification

```go
func VerifyHandshake(
    handshake *HandshakeMessage,
    profile *network.NetworkProfile,
    genesisHash [32]byte,
) error
```

**Checks**:
1. ✅ Genesis hash matches exactly
2. ✅ Chain ID matches
3. ✅ Network ID matches
4. ✅ Protocol version matches

**Any mismatch = immediate disconnection**

### Handshake Flow

```
┌─────────────┐                      ┌─────────────┐
│   Node A    │                      │   Node B    │
│  (Betanet)  │                      │  (Betanet)  │
└─────────────┘                      └─────────────┘
       │                                     │
       │─────── TCP Connection ──────────────>│
       │                                     │
       │<────── Handshake Message ───────────│
       │  {                                  │
       │    genesisHash: 0xa1b2...          │
       │    chainID: "archivas-betanet-1"   │
       │    networkID: 102                   │
       │    protocolVersion: 2               │
       │  }                                  │
       │                                     │
       │─── Verify Identity ─┐              │
       │                      │              │
       │<──── ✅ Accept ──────┘              │
       │                                     │
       │────── Handshake Message ───────────>│
       │  (send our identity)                │
       │                                     │
       │                          ┌─── Verify
       │                          │          │
       │<────────── ✅ Accept ────┘          │
       │                                     │
       │<════ Start Syncing ═══════════════>│
```

### Rejection Examples

**Scenario 1: Devnet node tries to connect to Betanet**

```
[p2p] Peer 192.168.1.100:9090 identity:
  Chain ID: archivas-devnet-1
  Network ID: 1
  Protocol: v1

[p2p] Rejecting peer: chain ID mismatch
Expected: archivas-betanet-1
Got: archivas-devnet-1

[p2p] Connection closed
```

**Scenario 2: Wrong genesis hash (forked chain)**

```
[p2p] Peer 192.168.1.101:9090 identity:
  Chain ID: archivas-betanet-1
  Network ID: 102
  Protocol: v2
  Genesis: d4e5f6a7...

[p2p] Rejecting peer: genesis hash mismatch
Expected: a1b2c3d4...
Got: d4e5f6a7...
This peer is on a different chain (fork)

[p2p] Connection closed
```

**Scenario 3: Old protocol version**

```
[p2p] Peer 192.168.1.102:9090 identity:
  Chain ID: archivas-betanet-1
  Network ID: 102
  Protocol: v1  ❌

[p2p] Rejecting peer: protocol version incompatible
Expected: v2
Got: v1
Please upgrade your node

[p2p] Connection closed
```

### Creating Handshakes

```go
// Create handshake for Betanet
profile, _ := network.GetProfile("betanet")
genesisHash := computeGenesisHash(genesis)

handshake := p2p.CreateHandshake(
    profile,
    genesisHash,
    "v1.0.0",  // node version
    "my-node", // node name
)

// Send to peer
sendHandshake(conn, handshake)
```

## 3. ETH RPC Endpoints

### Supported Endpoints

| Endpoint | Description | Status |
|----------|-------------|--------|
| `eth_chainId` | Returns chain ID | ✅ |
| `eth_blockNumber` | Latest block number | ✅ |
| `eth_getBalance` | Account balance | ✅ |
| `eth_getCode` | Contract code | ✅ |
| `eth_getTransactionReceipt` | Transaction receipt | ✅ |
| `eth_getTransactionCount` | Account nonce | ✅ |
| `eth_gasPrice` | Current gas price | ✅ |
| `eth_call` | Read-only call | ✅ |
| `eth_sendRawTransaction` | Submit signed TX | ✅ |
| `net_version` | Network ID | ✅ |

### ETHHandler Usage

```go
// Create ETH RPC handler
ethHandler := rpc.NewETHHandler(
    chainID,
    stateDB,
    getHeight,
    getBlock,
    getReceipt,
    submitTx,
)

// Mount on HTTP server
http.Handle("/eth", ethHandler)
```

### Example Requests

#### eth_chainId
```json
{
  "jsonrpc": "2.0",
  "method": "eth_chainId",
  "params": [],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": "0x66",
  "id": 1
}
```

#### eth_blockNumber
```json
{
  "jsonrpc": "2.0",
  "method": "eth_blockNumber",
  "params": [],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": "0x1312d0",
  "id": 1
}
```

#### eth_getBalance
```json
{
  "jsonrpc": "2.0",
  "method": "eth_getBalance",
  "params": [
    "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
    "latest"
  ],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": "0xde0b6b3a7640000",
  "id": 1
}
```

#### eth_getTransactionReceipt
```json
{
  "jsonrpc": "2.0",
  "method": "eth_getTransactionReceipt",
  "params": [
    "0xa1b2c3d4..."
  ],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "transactionHash": "0xa1b2c3d4...",
    "blockNumber": "0x1312d0",
    "transactionIndex": "0x0",
    "from": "0x742d35Cc...",
    "to": "0x1234567890...",
    "gasUsed": "0x5208",
    "cumulativeGasUsed": "0x5208",
    "contractAddress": null,
    "status": "0x1",
    "logs": []
  },
  "id": 1
}
```

### MetaMask Compatibility

ETH RPC endpoints allow MetaMask connection:

1. **Add Custom Network**:
   - Network Name: Archivas Betanet
   - RPC URL: `https://rpc.betanet.archivas.ai`
   - Chain ID: `102`
   - Currency Symbol: `RCHV`
   - Block Explorer: (optional)

2. **MetaMask will query**:
   - `eth_chainId` - Verify chain
   - `eth_blockNumber` - Get current height
   - `eth_getBalance` - Check balance
   - `eth_sendRawTransaction` - Submit transactions

3. **Fully compatible** with:
   - Web3.js
   - Ethers.js
   - MetaMask
   - Hardhat
   - Remix IDE

## 4. ARCV RPC Endpoints

### Archivas-Specific Endpoints

| Endpoint | Description |
|----------|-------------|
| `arcv_toHexAddress` | Convert ARCV to 0x |
| `arcv_fromHexAddress` | Convert 0x to ARCV |
| `arcv_validateAddress` | Validate either format |
| `arcv_getAddressInfo` | Get address info |

### ARCVHandler Usage

```go
// Create ARCV RPC handler
arcvHandler := rpc.NewARCVHandler("arcv")

// Mount on HTTP server
http.Handle("/arcv", arcvHandler)
```

### Example Requests

#### arcv_toHexAddress
```json
{
  "jsonrpc": "2.0",
  "method": "arcv_toHexAddress",
  "params": ["arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu"],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": "0x742d35cc6634c0532925a3b844bc9e7595f0beb0",
  "id": 1
}
```

#### arcv_fromHexAddress
```json
{
  "jsonrpc": "2.0",
  "method": "arcv_fromHexAddress",
  "params": ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": "arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu",
  "id": 1
}
```

#### arcv_validateAddress
```json
{
  "jsonrpc": "2.0",
  "method": "arcv_validateAddress",
  "params": ["arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu"],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "valid": true,
    "format": "bech32",
    "hex": "0x742d35cc6634c0532925a3b844bc9e7595f0beb0",
    "arcv": "arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu"
  },
  "id": 1
}
```

#### arcv_getAddressInfo
```json
{
  "jsonrpc": "2.0",
  "method": "arcv_getAddressInfo",
  "params": ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"],
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "hex": "0x742d35cc6634c0532925a3b844bc9e7595f0beb0",
    "arcv": "arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu",
    "isZero": false,
    "bytes": 20,
    "network": "arcv"
  },
  "id": 1
}
```

## Files Created/Modified

### New Files

**Snapshot Verification**:
- `snapshot/verify.go` - Manifest verification logic

**P2P Identity**:
- `p2p/identity.go` - Handshake verification

**ETH RPC**:
- `rpc/eth.go` - Ethereum-compatible RPC endpoints
- `rpc/arcv.go` - Archivas-specific RPC endpoints

### Modified Files

**Snapshot System**:
- `snapshot/snapshot.go` - Enhanced Manifest structure, Bootstrap verification

**P2P Protocol**:
- `p2p/protocol.go` - Enhanced HandshakeMessage structure

## Testing

### Manual Testing

```bash
# Test ETH RPC
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_chainId",
    "params": [],
    "id": 1
  }'

# Test ARCV RPC
curl -X POST http://localhost:8545/arcv \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "arcv_toHexAddress",
    "params": ["arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu"],
    "id": 1
  }'
```

## Summary

✅ **Phase 3 Complete**:
- Snapshot manifest verification (prevents wrong network snapshots)
- P2P identity enforcement (prevents incompatible peers)
- ETH RPC endpoints (Ethereum compatibility)
- ARCV RPC endpoints (address conversion)

🎯 **Security Benefits**:
- Betanet nodes can't connect to Devnet
- Forked chains are rejected
- Snapshots are verified before import
- Protocol version enforced

🚀 **User Benefits**:
- MetaMask compatibility
- Web3.js/Ethers.js support
- Address conversion utilities
- Full Ethereum ecosystem integration

---

**Date**: November 16, 2025  
**Version**: Phase 3 Complete  
**Status**: Production Ready ✅
