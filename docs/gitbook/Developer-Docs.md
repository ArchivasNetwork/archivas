# Developer Documentation

## Building Archivas

### Development Environment

**Requirements:**
- Go 1.21+
- Git
- Make (optional)
- 4GB RAM
- 10GB disk

**Setup:**
```bash
git clone https://github.com/ArchivasNetwork/archivas.git
cd archivas
go mod download
```

### Building from Source

**All binaries:**
```bash
go build -o bin/archivas-node ./cmd/archivas-node
go build -o bin/archivas-farmer ./cmd/archivas-farmer
go build -o bin/archivas-timelord ./cmd/archivas-timelord
go build -o bin/archivas-wallet ./cmd/archivas-wallet
```

**Individual packages:**
```bash
go build ./ledger
go build ./consensus
go build ./pospace
# etc.
```

**Run tests:**
```bash
go test ./...
```

---

## API Reference

### RPC Endpoints

#### GET /

Server status

**Response:**
```json
{
  "status": "ok",
  "message": "Archivas Devnet RPC Server"
}
```

#### GET /chainTip

Current chain tip

**Response:**
```json
{
  "blockHash": [byte array],
  "height": 78,
  "difficulty": 1125899906842624
}
```

#### GET /genesisHash

Genesis hash for network verification

**Response:**
```json
{
  "genesisHash": "11b6fedb68f1da0f..."
}
```

#### GET /balance/:address

Query account balance

**Response:**
```json
{
  "address": "arcv1q84xt5...",
  "balance": 2000000000,
  "nonce": 0
}
```

#### GET /challenge

Current mining challenge (for farmers)

**Response:**
```json
{
  "challenge": [byte array],
  "difficulty": 1125899906842624,
  "height": 79,
  "vdf": {
    "seed": "hex...",
    "iterations": 5000,
    "output": "hex..."
  }
}
```

#### POST /submitTx

Submit a signed transaction

**Request:**
```json
{
  "from": "arcv1...",
  "to": "arcv1...",
  "amount": 100000000,
  "fee": 100000,
  "nonce": 0,
  "senderPubKey": "hex...",
  "signature": "hex..."
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Transaction added to mempool"
}
```

#### POST /submitBlock

Submit a mined block (from farmer)

**Request:**
```json
{
  "proof": {
    "challenge": [byte array],
    "plotID": [byte array],
    "index": 12345,
    "hash": [byte array],
    "quality": 915713800595708,
    "farmerPubKey": [byte array]
  },
  "farmerAddr": "arcv1...",
  "farmerPubKey": "hex..."
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Block accepted"
}
```

#### POST /vdf/update

VDF update from timelord

**Request:**
```json
{
  "seed": [byte array],
  "iterations": 5000,
  "output": [byte array]
}
```

---

## Package APIs

### ledger

**World State:**
```go
func NewWorldState(genesisAlloc map[string]int64) *WorldState
func (ws *WorldState) GetBalance(addr string) int64
func (ws *WorldState) ApplyTransaction(tx Transaction) error
```

**Transactions:**
```go
type Transaction struct {
    From         string
    To           string
    Amount       int64
    Fee          int64
    Nonce        uint64
    SenderPubKey []byte
    Signature    []byte
}
```

### pospace

**Plot Management:**
```go
func GeneratePlot(path string, kSize uint32, farmerPubKey []byte) error
func OpenPlot(path string) (*PlotFile, error)
func (p *PlotFile) CheckChallenge(challenge [32]byte, difficultyTarget uint64) (*Proof, error)
func VerifyProof(proof *Proof, challenge [32]byte, difficultyTarget uint64) bool
```

### vdf

**VDF Computation:**
```go
func StepHash(input []byte) []byte
func ComputeSequential(seed []byte, iterations uint64, checkpointStep uint64) (final []byte, checkpoints [][]byte)
func VerifySequential(seed []byte, iterations uint64, claimedFinal []byte) bool
```

### wallet

**Key Management:**
```go
func GenerateKeypair() (privKey []byte, pubKey []byte, err error)
func PubKeyToAddress(pubKey []byte) (string, error)
func SignTransaction(tx *Transaction, privKey []byte) error
```

### consensus

**Difficulty:**
```go
func NewConsensus() *Consensus
func (c *Consensus) VerifyProofOfSpace(proof *Proof, challenge [32]byte) error
func (c *Consensus) UpdateDifficulty(recentBlockTimes []int64)
```

---

## Contributing

### Development Workflow

**1. Fork and clone:**
```bash
git clone https://github.com/YOUR_USERNAME/archivas.git
cd archivas
git remote add upstream https://github.com/ArchivasNetwork/archivas.git
```

**2. Create branch:**
```bash
git checkout -b feature/my-feature
```

**3. Make changes:**
```bash
# Edit code
# Run tests
go test ./...

# Build
go build ./cmd/...
```

**4. Commit:**
```bash
git add .
git commit -m "feat: add my feature"
```

**5. Push and PR:**
```bash
git push origin feature/my-feature
# Create PR on GitHub
```

### Code Style

**Follow Go conventions:**
- `go fmt` all code
- Run `go vet`
- Add comments for exported functions
- Write tests for new features

**Commit messages:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation
- `refactor:` - Code refactoring
- `test:` - Tests

### Testing

**Unit tests:**
```bash
go test ./pospace
go test ./vdf
go test ./consensus
```

**Integration tests:**
```bash
# Start local testnet
./scripts/test-local-network.sh
```

**Manual testing:**
```bash
# Build and run
go build -o archivas-node ./cmd/archivas-node
./archivas-node --genesis genesis/devnet.genesis.json ...
```

---

## Building Features

### Adding RPC Endpoints

**Example: Add /status endpoint**

```go
// In rpc/farming.go

func (s *FarmingServer) Start(addr string) error {
    // ... existing endpoints
    http.HandleFunc("/status", s.handleStatus)
    return http.ListenAndServe(addr, nil)
}

func (s *FarmingServer) handleStatus(w http.ResponseWriter, r *http.Request) {
    status := struct {
        Height     uint64 `json:"height"`
        Difficulty uint64 `json:"difficulty"`
        Peers      int    `json:"peers"`
    }{
        Height:     s.nodeState.LocalHeight(),
        Difficulty: s.nodeState.GetDifficulty(),
        Peers:      s.nodeState.GetPeerCount(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}
```

### Adding P2P Messages

**Example: Add PING/PONG**

```go
// In p2p/protocol.go

const (
    MsgTypePing MessageType = 10
    MsgTypePong MessageType = 11
)

type PingMessage struct {
    Timestamp int64 `json:"timestamp"`
}

// In p2p/p2p.go

func (n *Network) handlePing(peer *Peer, payload json.RawMessage) {
    var ping PingMessage
    json.Unmarshal(payload, &ping)
    
    pong := PingMessage{Timestamp: time.Now().Unix()}
    n.SendMessage(peer, MsgTypePong, pong)
}
```

### Adding Consensus Rules

**Example: New difficulty algorithm**

```go
// In consensus/difficulty.go

func CalculateDifficultyV2(recentBlocks []*Block) uint64 {
    // Your improved algorithm
    // Must be deterministic!
    return newDifficulty
}
```

---

## Architecture Decisions

### Why Go?

- Systems language (performance)
- Strong stdlib (crypto, networking)
- Good concurrency (goroutines)
- Fast compilation
- Static binaries

### Why BadgerDB?

- Embedded (no separate server)
- Fast (LSM tree)
- ACID transactions
- Go-native
- Battle-tested

### Why TCP P2P?

- Simple and reliable
- Easy to debug
- Works everywhere
- Can upgrade to libp2p later

### Why SHA-256 VDF?

**Devnet:**
- Simple to implement
- Easy to understand
- Good for testing
- No dependencies

**Mainnet will use:**
- Wesolowski or Pietrzak VDF
- RSA/class group operations
- Succinct proofs
- Quantum-resistant

---

## Performance Considerations

### Bottlenecks

**Current:**
- Plot scanning: O(n) through all hashes
- VDF: Single-threaded (sequential by design)
- P2P: No multiplexing
- Storage: No pruning

**Optimizations:**
- Plot indexing (k1/k2 tables like Chia)
- VDF hardware acceleration
- P2P connection pooling
- State snapshots

### Scalability

**Block size:** ~1-10 KB currently  
**TPS:** ~10-100 (depends on block time)  
**State size:** ~1 MB per 1000 blocks  
**Network:** Supports 10-100 nodes easily

**Future:**
- Sharding
- Rollups
- State channels
- Light clients

---

## Security Considerations

### Threat Model

**Assumptions:**
- SHA-256 is collision-resistant
- secp256k1 is secure
- Majority of disk space is honest
- VDF is sequential

**Risks:**
- Bugs in consensus logic
- P2P vulnerabilities
- Database corruption
- Key management

### Audit Status

**Current:** Not audited (alpha software!)

**Before mainnet:**
- External security audit
- Formal verification of consensus
- Penetration testing
- Economic analysis

---

## Roadmap for Developers

### Short-term

**Testnet improvements:**
- [ ] Peer persistence (reconnect on restart)
- [ ] Block explorer API
- [ ] Prometheus metrics
- [ ] Light client protocol
- [ ] Better error messages

### Medium-term

**Production features:**
- [ ] Wesolowski VDF
- [ ] Plot compression
- [ ] State snapshots
- [ ] WASM smart contracts
- [ ] Cross-chain bridges

### Long-term

**Ecosystem:**
- [ ] DEX on Archivas
- [ ] NFT standard
- [ ] Oracles
- [ ] L2 rollups
- [ ] Mobile wallets

---

**Next:** [Roadmap →](Roadmap.md)  
**Back:** [← Testnet Guide](Testnet-Guide.md)

