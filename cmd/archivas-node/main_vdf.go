package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/iljanemesis/archivas/config"
	"github.com/iljanemesis/archivas/consensus"
	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/mempool"
	"github.com/iljanemesis/archivas/pospace"
	"github.com/iljanemesis/archivas/rpc"
	"github.com/iljanemesis/archivas/vdf"
)

// Block represents a blockchain block with Proof-of-Space and VDF
type BlockVDF struct {
	Height        uint64
	TimestampUnix int64
	PrevHash      [32]byte
	Txs           []ledger.Transaction
	Proof         *pospace.Proof // Proof-of-Space
	FarmerAddr    string         // Address to receive block reward
	// VDF fields
	VDFSeed       []byte
	VDFIterations uint64
	VDFOutput     []byte
}

// VDFState holds current VDF state
type VDFState struct {
	Seed       []byte
	Iterations uint64
	Output     []byte
	UpdatedAt  time.Time
}

// NodeStateVDF holds the entire node state with VDF
type NodeStateVDF struct {
	sync.RWMutex
	Chain            []BlockVDF
	WorldState       *ledger.WorldState
	Mempool          *mempool.Mempool
	Consensus        *consensus.Consensus
	CurrentHeight    uint64
	CurrentChallenge [32]byte
	CurrentVDF       *VDFState
}

func mainVDF() {
	log.Println("[DEBUG] Archivas node starting (VDF-enabled)...")
	fmt.Println("Archivas Devnet Node running… (Proof-of-Space-and-Time)")
	fmt.Println()

	// Display chain configuration
	log.Println("[DEBUG] Loading chain configuration...")
	fmt.Printf("⛓️  Chain: %s\n", config.ChainName)
	fmt.Printf("🆔 Chain ID: %d\n", config.ChainID)
	fmt.Printf("💰 Token: %s (decimals: %d)\n", config.DenomSymbol, config.DenomDecimals)
	fmt.Printf("⏱️  Target Block Time: %d seconds\n", config.TargetBlockTimeSeconds)
	fmt.Printf("🎁 Initial Block Reward: %d (%.8f %s)\n",
		config.InitialBlockReward,
		float64(config.InitialBlockReward)/100000000.0,
		config.DenomSymbol,
	)
	fmt.Println()

	// Initialize world state with genesis balances
	log.Println("[DEBUG] Initializing world state...")
	worldState := ledger.NewWorldState(config.LegacyGenesisAlloc)
	log.Println("[DEBUG] Loaded world state")
	fmt.Printf("🌍 World state initialized with %d genesis accounts\n", len(config.LegacyGenesisAlloc))
	for addr, balance := range config.LegacyGenesisAlloc {
		fmt.Printf("   %s: %.8f %s\n", addr, float64(balance)/100000000.0, config.DenomSymbol)
	}
	fmt.Println()

	// Initialize mempool
	log.Println("[DEBUG] Initializing mempool...")
	mp := mempool.NewMempool()
	fmt.Println("📋 Mempool initialized")

	// Initialize consensus
	log.Println("[DEBUG] Initializing consensus...")
	cs := consensus.NewConsensus()
	fmt.Printf("⚙️  Consensus initialized (difficulty: %d)\n", cs.DifficultyTarget)
	fmt.Println()

	// Create genesis block
	log.Println("[DEBUG] Creating genesis block...")
	genesisHash := [32]byte{}
	copy(genesisHash[:], []byte("genesis"))

	// Initial VDF state (genesis)
	genesisSeed := computeVDFSeed(genesisHash, 0)
	initialVDF := &VDFState{
		Seed:       genesisSeed,
		Iterations: 0,
		Output:     genesisSeed, // At iteration 0, output = seed
		UpdatedAt:  time.Now(),
	}

	genesisChallenge := computeChallengeFromVDF(initialVDF.Output, 1)
	genesisBlock := BlockVDF{
		Height:        0,
		TimestampUnix: time.Now().Unix(),
		PrevHash:      [32]byte{},
		Txs:           nil,
		Proof:         nil, // No proof needed for genesis
		FarmerAddr:    "",
		VDFSeed:       genesisSeed,
		VDFIterations: 0,
		VDFOutput:     genesisSeed,
	}

	// Initialize node state
	log.Println("[DEBUG] Initializing node state...")
	nodeState := &NodeStateVDF{
		Chain:            []BlockVDF{genesisBlock},
		WorldState:       worldState,
		Mempool:          mp,
		Consensus:        cs,
		CurrentHeight:    0,
		CurrentChallenge: genesisChallenge,
		CurrentVDF:       initialVDF,
	}

	log.Println("[DEBUG] Initialized chain memory")
	fmt.Printf("📦 Genesis block created at height %d\n", genesisBlock.Height)
	fmt.Printf("🔍 Genesis challenge: %x\n", genesisChallenge[:8])
	fmt.Printf("⏰ VDF seed: %x\n", genesisSeed[:8])
	fmt.Println()

	// Start RPC server in background
	log.Println("[DEBUG] Starting RPC server...")
	server := rpc.NewVDFServer(nodeState.WorldState, nodeState.Mempool, nodeState)
	go func() {
		log.Println("[rpc] starting server on :8080")
		fmt.Println("🌐 Starting RPC server on :8080")
		if err := server.Start(":8080"); err != nil {
			log.Fatalf("[rpc] server error: %v", err)
		}
	}()

	// Give RPC server a moment to start
	time.Sleep(500 * time.Millisecond)
	log.Println("[DEBUG] RPC server running")

	fmt.Println("✅ Node initialized successfully")
	fmt.Println("🌾 Waiting for timelord and farmers...")
	fmt.Println()

	log.Println("[DEBUG] Entering consensus loop...")
	logCurrentState(nodeState)

	// Heartbeat loop - shows node is alive and waiting for blocks
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		logCurrentState(nodeState)
	}
}

func logCurrentState(ns *NodeStateVDF) {
	ns.RLock()
	defer ns.RUnlock()

	vdfAge := time.Since(ns.CurrentVDF.UpdatedAt)
	log.Printf("[consensus] height=%d difficulty=%d challenge=%x vdfIter=%d vdfAge=%v chainLen=%d",
		ns.CurrentHeight, ns.Consensus.DifficultyTarget, ns.CurrentChallenge[:8],
		ns.CurrentVDF.Iterations, vdfAge.Round(time.Millisecond), len(ns.Chain))
}

// AcceptBlock is called by RPC when a farmer submits a block
func (ns *NodeStateVDF) AcceptBlock(proof *pospace.Proof, farmerAddr string, farmerPubKey []byte,
	vdfSeed []byte, vdfIterations uint64, vdfOutput []byte) error {
	ns.Lock()
	defer ns.Unlock()

	// Get expected height
	nextHeight := ns.CurrentHeight + 1

	// Verify VDF fields match current VDF state
	if !bytesEqual(vdfSeed, ns.CurrentVDF.Seed) {
		return fmt.Errorf("VDF seed mismatch: expected %x, got %x", ns.CurrentVDF.Seed[:8], vdfSeed[:8])
	}

	// Verify VDF output
	if !vdf.VerifySequential(vdfSeed, vdfIterations, vdfOutput) {
		return fmt.Errorf("VDF verification failed")
	}

	// Verify proof against VDF-derived challenge
	vdfChallenge := computeChallengeFromVDF(vdfOutput, nextHeight)
	if err := ns.Consensus.VerifyProofOfSpace(proof, vdfChallenge); err != nil {
		return fmt.Errorf("invalid proof: %w", err)
	}

	// Get pending transactions
	pending := ns.Mempool.Pending()

	// Create coinbase transaction (block reward to farmer)
	coinbase := ledger.Transaction{
		From:         "coinbase",
		To:           farmerAddr,
		Amount:       config.InitialBlockReward,
		Fee:          0,
		Nonce:        0,
		SenderPubKey: nil,
		Signature:    nil,
	}

	// Build transaction list
	allTxs := []ledger.Transaction{coinbase}

	// Apply coinbase
	receiver, ok := ns.WorldState.Accounts[farmerAddr]
	if !ok {
		receiver = &ledger.AccountState{Balance: 0, Nonce: 0}
		ns.WorldState.Accounts[farmerAddr] = receiver
	}
	receiver.Balance += config.InitialBlockReward

	// Apply user transactions
	validTxs := []ledger.Transaction{}
	for _, tx := range pending {
		err := ns.WorldState.ApplyTransaction(tx)
		if err != nil {
			log.Printf("⚠️  Skipping invalid tx: %v", err)
		} else {
			validTxs = append(validTxs, tx)
		}
	}

	allTxs = append(allTxs, validTxs...)

	// Calculate prev hash
	var prevHash [32]byte
	if nextHeight > 0 {
		prevBlock := ns.Chain[len(ns.Chain)-1]
		prevHash = hashBlockVDF(&prevBlock)
	}

	// Create new block
	newBlock := BlockVDF{
		Height:        nextHeight,
		TimestampUnix: time.Now().Unix(),
		PrevHash:      prevHash,
		Txs:           allTxs,
		Proof:         proof,
		FarmerAddr:    farmerAddr,
		VDFSeed:       vdfSeed,
		VDFIterations: vdfIterations,
		VDFOutput:     vdfOutput,
	}

	// Add to chain
	ns.Chain = append(ns.Chain, newBlock)
	ns.CurrentHeight = nextHeight

	// Clear mempool
	ns.Mempool.Clear()

	// Generate new VDF seed for next block
	newBlockHash := hashBlockVDF(&newBlock)
	newVDFSeed := computeVDFSeed(newBlockHash, nextHeight)

	// Reset VDF state (timelord will restart from this seed)
	ns.CurrentVDF = &VDFState{
		Seed:       newVDFSeed,
		Iterations: 0,
		Output:     newVDFSeed,
		UpdatedAt:  time.Now(),
	}

	// Generate challenge for next height
	ns.CurrentChallenge = computeChallengeFromVDF(newVDFSeed, nextHeight+1)

	// Update difficulty
	if len(ns.Chain) >= 2 {
		recentTimes := make([]int64, 0, consensus.DifficultyAdjustmentWindow)
		startIdx := len(ns.Chain) - consensus.DifficultyAdjustmentWindow
		if startIdx < 0 {
			startIdx = 0
		}
		for i := startIdx; i < len(ns.Chain); i++ {
			recentTimes = append(recentTimes, ns.Chain[i].TimestampUnix)
		}
		ns.Consensus.UpdateDifficulty(recentTimes)
	}

	log.Printf("✅ Accepted block %d from farmer %s (PoSpace ✅, VDF t=%d ✅, reward: %.8f %s, txs: %d)",
		nextHeight, farmerAddr, vdfIterations, float64(config.InitialBlockReward)/100000000.0, config.DenomSymbol, len(validTxs))
	log.Printf("🔍 New VDF seed for height %d: %x", nextHeight+1, newVDFSeed[:8])
	log.Printf("⚙️  Difficulty adjusted to: %d", ns.Consensus.DifficultyTarget)

	return nil
}

// UpdateVDF is called by timelord to update VDF state
func (ns *NodeStateVDF) UpdateVDF(seed []byte, iterations uint64, output []byte) error {
	ns.Lock()
	defer ns.Unlock()

	// Verify seed matches current expected seed
	if !bytesEqual(seed, ns.CurrentVDF.Seed) {
		return fmt.Errorf("VDF seed mismatch")
	}

	// Verify the VDF computation
	if !vdf.VerifySequential(seed, iterations, output) {
		return fmt.Errorf("VDF verification failed")
	}

	// Update VDF state
	ns.CurrentVDF.Iterations = iterations
	ns.CurrentVDF.Output = output
	ns.CurrentVDF.UpdatedAt = time.Now()

	// Update challenge based on new VDF output
	ns.CurrentChallenge = computeChallengeFromVDF(output, ns.CurrentHeight+1)

	return nil
}

// GetChainTip returns current chain tip
func (ns *NodeStateVDF) GetChainTip() ([32]byte, uint64, uint64) {
	ns.RLock()
	defer ns.RUnlock()
	
	if len(ns.Chain) == 0 {
		return [32]byte{}, 0, ns.Consensus.DifficultyTarget
	}
	
	tip := ns.Chain[len(ns.Chain)-1]
	return hashBlockVDF(&tip), tip.Height, ns.Consensus.DifficultyTarget
}

// GetCurrentChallenge returns the current challenge and difficulty for farmers
func (ns *NodeStateVDF) GetCurrentChallengeVDF() ([32]byte, uint64, uint64, []byte, uint64, []byte) {
	ns.RLock()
	defer ns.RUnlock()
	return ns.CurrentChallenge, ns.Consensus.DifficultyTarget, ns.CurrentHeight + 1,
		ns.CurrentVDF.Seed, ns.CurrentVDF.Iterations, ns.CurrentVDF.Output
}

// Helper functions
func computeVDFSeed(blockHash [32]byte, height uint64) []byte {
	h := sha256.New()
	h.Write(blockHash[:])
	binary.Write(h, binary.BigEndian, height)
	return h.Sum(nil)
}

func computeChallengeFromVDF(vdfOutput []byte, height uint64) [32]byte {
	h := sha256.New()
	h.Write(vdfOutput)
	binary.Write(h, binary.BigEndian, height)
	return sha256.Sum256(h.Sum(nil))
}

func hashBlockVDF(b *BlockVDF) [32]byte {
	h := sha256.New()
	fmt.Fprintf(h, "%d", b.Height)
	fmt.Fprintf(h, "%d", b.TimestampUnix)
	h.Write(b.PrevHash[:])
	if b.Proof != nil {
		h.Write(b.Proof.Hash[:])
	}
	h.Write(b.VDFOutput)
	return sha256.Sum256(h.Sum(nil))
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
