# Betanet Quick Start Guide

**Version**: Phase 1, 2 & 3 Complete  
**Status**: Production Ready ✅

---

## 🚀 Quick Commands

```bash
# Run Betanet node (default)
archivas-node

# Run with explicit network
archivas-node --network betanet

# Run Devnet Legacy
archivas-node --network devnet-legacy

# Show help
archivas-node help

# Show version
archivas-node version
```

---

## 📦 Build & Install

```bash
# Clone repository
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas

# Build
go build -o archivas-node ./cmd/archivas-node/main.go

# Install (optional)
sudo mv archivas-node /usr/local/bin/

# Verify
archivas-node version
```

---

## 🌐 Networks

### Betanet (Default)

```bash
archivas-node --network betanet
```

- **Chain ID**: `archivas-betanet-1`
- **Network ID**: `102`
- **Protocol**: `v2`
- **RPC Port**: `8545`
- **P2P Port**: `9090`
- **Features**: Full EVM support

### Devnet Legacy

```bash
archivas-node --network devnet-legacy
```

- **Chain ID**: `archivas-devnet-1`
- **Network ID**: `1`
- **Protocol**: `v1`
- **RPC Port**: `8080`
- **P2P Port**: `9090`
- **Features**: Legacy (no EVM)

---

## 🔧 Configuration Examples

### Public Seed Node

```bash
archivas-node \
  --network betanet \
  --rpc 0.0.0.0:8545 \
  --p2p 0.0.0.0:9090 \
  --db ./data
```

### Private Node (Farming)

```bash
archivas-node \
  --network betanet \
  --rpc 127.0.0.1:8545 \
  --p2p 0.0.0.0:9090 \
  --no-peer-discovery \
  --peer-whitelist seed1.betanet.archivas.ai:9090 \
  --peer-whitelist seed2.betanet.archivas.ai:9090
```

### Bootstrap from Snapshot

```bash
# Bootstrap automatically downloads and imports snapshot
archivas-node bootstrap --network betanet
```

---

## 🔌 MetaMask Setup

### Add Betanet to MetaMask

1. Open MetaMask
2. Click network dropdown
3. Select "Add Network"
4. Enter details:
   - **Network Name**: Archivas Betanet
   - **RPC URL**: `https://rpc.betanet.archivas.ai`
   - **Chain ID**: `102`
   - **Currency Symbol**: `RCHV`
   - **Block Explorer**: (optional)
5. Click "Save"

---

## 📍 Address Formats

### Internal (EVM)

```
0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0
```

- 20-byte hex address
- Used by EVM internally
- Compatible with Ethereum tools

### External (Bech32)

```
arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu
```

- User-facing address
- Checksum validation
- Cosmos-style

### Both Work!

```bash
# CLI accepts both formats
archivas-node --farmer-addr 0x742d35Cc...
archivas-node --farmer-addr arcv1wskntnr...
```

---

## 🧪 Test the Installation

### Run Address Demo

```bash
go run examples/address/address_demo.go
```

**Expected Output**:
```
═══════════════════════════════════════════════════
   Archivas Betanet - Dual Address System Demo
═══════════════════════════════════════════════════

Internal (EVM):  0x742d35cc6634c0532925a3b844bc9e7595f0beb0
External (ARCV): arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu
✅ All roundtrip tests passed!
```

### Run EVM Demo

```bash
go run examples/evm/evm_demo.go
```

**Expected Output**:
```
═══════════════════════════════════════════════════
   Archivas Betanet - EVM Engine Demo
═══════════════════════════════════════════════════

Transaction 1: Simple Transfer
📊 Status: ✅ Success
📊 Gas used: 21000

Transaction 2: Contract Deployment
📊 Status: ✅ Success
📊 Contract address: 0xa0fcc4706ad995b41eb0edc2937a2d97caf95083
✅ All transactions executed successfully!
```

### Run Tests

```bash
# Address tests
go test ./address/... -v

# EVM tests
go test ./evm/... -v

# All tests
go test ./... -v
```

**Expected**: All tests pass ✅

---

## 🔐 RPC Endpoints

### ETH RPC (Port 8545)

```bash
# eth_chainId
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'

# eth_blockNumber
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# eth_getBalance
curl -X POST http://localhost:8545/eth \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"eth_getBalance",
    "params":["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0","latest"],
    "id":1
  }'
```

### ARCV RPC (Address Conversion)

```bash
# Convert ARCV to 0x
curl -X POST http://localhost:8545/arcv \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"arcv_toHexAddress",
    "params":["arcv1wskntnrxxnq9x2f95wuyf0y7wk2lp04s47qnwu"],
    "id":1
  }'

# Convert 0x to ARCV
curl -X POST http://localhost:8545/arcv \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "method":"arcv_fromHexAddress",
    "params":["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"],
    "id":1
  }'
```

---

## 📚 Documentation

- **Phase 1**: [`docs/BETANET_PHASE1.md`](docs/BETANET_PHASE1.md)
- **Phase 2**: [`docs/BETANET_PHASE2.md`](docs/BETANET_PHASE2.md)
- **Phase 3**: [`docs/BETANET_PHASE3.md`](docs/BETANET_PHASE3.md)
- **Complete**: [`docs/BETANET_COMPLETE.md`](docs/BETANET_COMPLETE.md)
- **Progress**: [`docs/BETANET_PROGRESS.md`](docs/BETANET_PROGRESS.md)

---

## 🐛 Troubleshooting

### Node won't start

```bash
# Check if port is already in use
lsof -i :8545
lsof -i :9090

# Try different ports
archivas-node --rpc :9545 --p2p :10090
```

### Can't connect to peers

```bash
# Check firewall
sudo ufw status
sudo ufw allow 9090

# Use explicit peers
archivas-node --peer seed1.betanet.archivas.ai:9090
```

### Wrong network

```bash
# Explicitly set network
archivas-node --network betanet

# Check current network
archivas-node --network betanet | grep "Network:"
```

### Database issues

```bash
# Reset database
rm -rf ./data
archivas-node bootstrap --network betanet
```

---

## ✅ Checklist

Before going to production:

- [ ] Build successful: `go build ./cmd/archivas-node/main.go`
- [ ] Tests passing: `go test ./address/... ./evm/... -v`
- [ ] Demos working: `go run examples/*/main.go`
- [ ] Node starts: `./archivas-node --network betanet`
- [ ] RPC accessible: `curl http://localhost:8545/eth -d '...'`
- [ ] MetaMask connects
- [ ] Peers connecting (check logs)
- [ ] Blocks syncing (check height)

---

## 🆘 Support

- **Documentation**: `docs/`
- **Issues**: GitHub Issues
- **Discord**: (link)
- **Twitter**: (link)

---

**Last Updated**: November 16, 2025  
**Version**: Betanet Phase 1, 2 & 3 Complete  
**Status**: ✅ Production Ready

