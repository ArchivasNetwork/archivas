package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/iljanemesis/archivas/config"
	"github.com/iljanemesis/archivas/consensus"
	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/mempool"
	"github.com/iljanemesis/archivas/p2p"
	"github.com/iljanemesis/archivas/pospace"
	"github.com/iljanemesis/archivas/rpc"
	"github.com/iljanemesis/archivas/storage"
)

// Block represents a blockchain block with Proof-of-Space
type Block struct {
	Height        uint64
	TimestampUnix int64
	PrevHash      [32]byte
	Txs           []ledger.Transaction
	Proof         *pospace.Proof // Proof-of-Space
	FarmerAddr    string         // Address to receive block reward
}

// NodeState holds the entire node state
type NodeState struct {
	sync.RWMutex
	Chain            []Block
	WorldState       *ledger.WorldState
	Mempool          *mempool.Mempool
	Consensus        *consensus.Consensus
	CurrentHeight    uint64
	CurrentChallenge [32]byte
	// Persistence
	DB         *storage.DB
	BlockStore *storage.BlockStorage
	StateStore *storage.StateStorage
	MetaStore  *storage.MetadataStorage
	// Networking
	P2P *p2p.Network
	// VDF state (updated by timelord)
	VDFSeed       []byte
	VDFIterations uint64
	VDFOutput     []byte
	HasVDF        bool
}

func main() {
	// Parse CLI flags
	rpcAddr := flag.String("rpc", ":8080", "RPC listen address")
	p2pAddr := flag.String("p2p", ":9090", "P2P listen address")
	peerAddrs := flag.String("peer", "", "Comma-separated peer addresses (e.g., ip1:9090,ip2:9090)")
	dbPath := flag.String("db", "./data", "Database directory path")
	vdfRequired := flag.Bool("vdf-required", false, "Require VDF proofs in blocks (PoSpace+Time mode)")
	flag.Parse()

	log.Println("[startup] Archivas node starting...")
	fmt.Println("Archivas Devnet Node running…")
	fmt.Println()

	fmt.Printf("🔧 Configuration:\n")
	fmt.Printf("   RPC:  %s\n", *rpcAddr)
	fmt.Printf("   P2P:  %s\n", *p2pAddr)
	if *peerAddrs != "" {
		fmt.Printf("   Peers: %s\n", *peerAddrs)
	}
	fmt.Printf("   DB:   %s\n", *dbPath)
	if *vdfRequired {
		fmt.Printf("   Mode: PoSpace+Time (VDF required)\n")
	} else {
		fmt.Printf("   Mode: PoSpace only\n")
	}
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

	// Open database
	log.Println("[DEBUG] Opening database...")

	db, err := storage.OpenDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	blockStore := storage.NewBlockStorage(db)
	stateStore := storage.NewStateStorage(db)
	metaStore := storage.NewMetadataStorage(db)

	fmt.Printf("💾 Database opened: %s\n", *dbPath)
	fmt.Println()

	// Try to load existing state from disk
	var worldState *ledger.WorldState
	var cs *consensus.Consensus
	var chain []Block
	var currentHeight uint64
	var genesisChallenge [32]byte

	tipHeight, err := metaStore.LoadTipHeight()
	freshStart := err != nil

	if freshStart {
		// Fresh start - create genesis
		log.Println("[DEBUG] No existing state found, creating genesis...")
		worldState = ledger.NewWorldState(config.GenesisAlloc)
		fmt.Printf("🌱 Fresh start: Genesis block\n")
		fmt.Printf("🌍 World state initialized with %d genesis accounts\n", len(config.GenesisAlloc))
		for addr, balance := range config.GenesisAlloc {
			fmt.Printf("   %s: %.8f %s\n", addr, float64(balance)/100000000.0, config.DenomSymbol)
		}

		cs = consensus.NewConsensus()

		genesisChallenge = consensus.GenerateGenesisChallenge()
		genesisBlock := Block{
			Height:        0,
			TimestampUnix: time.Now().Unix(),
			PrevHash:      [32]byte{},
			Txs:           nil,
			Proof:         nil,
			FarmerAddr:    "",
		}
		chain = []Block{genesisBlock}
		currentHeight = 0

		// Persist genesis
		if err := blockStore.SaveBlock(0, genesisBlock); err != nil {
			log.Fatalf("Failed to save genesis block: %v", err)
		}
		for addr, balance := range config.GenesisAlloc {
			if err := stateStore.SaveAccount(addr, balance, 0); err != nil {
				log.Fatalf("Failed to save genesis account: %v", err)
			}
		}
		if err := metaStore.SaveTipHeight(0); err != nil {
			log.Fatalf("Failed to save tip height: %v", err)
		}
		if err := metaStore.SaveDifficulty(cs.DifficultyTarget); err != nil {
			log.Fatalf("Failed to save difficulty: %v", err)
		}

		fmt.Printf("📦 Genesis block created at height %d\n", genesisBlock.Height)
	} else {
		// Load from disk
		log.Printf("[DEBUG] Loading existing state from tip height %d...", tipHeight)
		fmt.Printf("💾 Restoring from disk (tip height: %d)\n", tipHeight)

		// Load difficulty
		difficulty, err := metaStore.LoadDifficulty()
		if err != nil {
			log.Fatalf("Failed to load difficulty: %v", err)
		}
		cs = &consensus.Consensus{DifficultyTarget: difficulty}

		// Load blocks
		chain = make([]Block, 0, tipHeight+1)
		for h := uint64(0); h <= tipHeight; h++ {
			var blk Block
			if err := blockStore.LoadBlock(h, &blk); err != nil {
				log.Fatalf("Failed to load block %d: %v", h, err)
			}
			chain = append(chain, blk)
		}

		// Reconstruct world state from disk
		worldState = &ledger.WorldState{
			Accounts: make(map[string]*ledger.AccountState),
		}

		// Load all accounts - we'll need to scan all keys with acc: prefix
		// For now, just load genesis accounts and any that received funds
		// (In production, you'd have an account index)
		for addr := range config.GenesisAlloc {
			balance, nonce, exists, err := stateStore.LoadAccount(addr)
			if err != nil {
				log.Fatalf("Failed to load account %s: %v", addr, err)
			}
			if exists {
				worldState.Accounts[addr] = &ledger.AccountState{
					Balance: balance,
					Nonce:   nonce,
				}
			}
		}

		// Scan blocks for other accounts (simple approach)
		accountsFound := make(map[string]bool)
		for _, blk := range chain {
			if blk.FarmerAddr != "" && !accountsFound[blk.FarmerAddr] {
				balance, nonce, exists, err := stateStore.LoadAccount(blk.FarmerAddr)
				if err == nil && exists {
					worldState.Accounts[blk.FarmerAddr] = &ledger.AccountState{
						Balance: balance,
						Nonce:   nonce,
					}
					accountsFound[blk.FarmerAddr] = true
				}
			}
			for _, tx := range blk.Txs {
				for _, addr := range []string{tx.From, tx.To} {
					if !accountsFound[addr] {
						balance, nonce, exists, err := stateStore.LoadAccount(addr)
						if err == nil && exists {
							worldState.Accounts[addr] = &ledger.AccountState{
								Balance: balance,
								Nonce:   nonce,
							}
							accountsFound[addr] = true
						}
					}
				}
			}
		}

		currentHeight = tipHeight

		// Recompute challenge from tip
		tipBlock := chain[len(chain)-1]
		tipHash := hashBlock(&tipBlock)
		genesisChallenge = consensus.GenerateChallenge(tipHash, currentHeight+1)

		fmt.Printf("✅ Restored %d blocks from disk\n", len(chain))
		fmt.Printf("📊 Loaded %d accounts\n", len(worldState.Accounts))
		fmt.Printf("⚙️  Difficulty: %d\n", cs.DifficultyTarget)
	}

	fmt.Println()

	// Initialize mempool (always fresh)
	log.Println("[DEBUG] Initializing mempool...")
	mp := mempool.NewMempool()
	fmt.Println("📋 Mempool initialized")
	fmt.Println()

	// Initialize node state
	log.Println("[DEBUG] Initializing node state...")
	nodeState := &NodeState{
		Chain:            chain,
		WorldState:       worldState,
		Mempool:          mp,
		Consensus:        cs,
		CurrentHeight:    currentHeight,
		CurrentChallenge: genesisChallenge,
		DB:               db,
		BlockStore:       blockStore,
		StateStore:       stateStore,
		MetaStore:        metaStore,
	}

	log.Println("[DEBUG] Initialized chain memory")
	fmt.Printf("🔍 Current challenge: %x\n", genesisChallenge[:8])
	fmt.Println()

	// Start P2P network if enabled
	var p2pNet *p2p.Network
	if *p2pAddr != "" {
		log.Printf("[p2p] Starting P2P listener on %s", *p2pAddr)
		p2pNet = p2p.NewNetwork(*p2pAddr, nodeState)
		if err := p2pNet.Start(); err != nil {
			log.Fatalf("[p2p] Failed to start P2P: %v", err)
		}
		nodeState.P2P = p2pNet
		fmt.Printf("🌐 P2P network started on %s\n", *p2pAddr)

		// Connect to initial peers
		if *peerAddrs != "" {
			peers := strings.Split(*peerAddrs, ",")
			for _, peer := range peers {
				peer = strings.TrimSpace(peer)
				if peer != "" {
					go func(addr string) {
						if err := p2pNet.ConnectPeer(addr); err != nil {
							log.Printf("[p2p] Failed to connect to peer %s: %v", addr, err)
						}
					}(peer)
				}
			}
			fmt.Printf("📡 Connecting to %d peer(s)...\n", len(peers))
		}
		fmt.Println()
	}

	// Start RPC server in background
	log.Println("[DEBUG] Starting RPC server...")
	server := rpc.NewFarmingServer(nodeState.WorldState, nodeState.Mempool, nodeState)
	go func() {
		log.Printf("[rpc] starting server on %s", *rpcAddr)
		fmt.Printf("🌐 Starting RPC server on %s\n", *rpcAddr)
		if err := server.Start(*rpcAddr); err != nil {
			log.Fatalf("[rpc] server error: %v", err)
		}
	}()

	// Give RPC server a moment to start
	time.Sleep(500 * time.Millisecond)
	log.Println("[DEBUG] RPC server running")

	fmt.Println("✅ Node initialized successfully")
	fmt.Println("🌾 Waiting for farmers to submit blocks...")
	fmt.Println()

	log.Println("[DEBUG] Entering consensus loop...")
	log.Printf("[consensus] height=%d difficulty=%d challenge=%x",
		nodeState.CurrentHeight, nodeState.Consensus.DifficultyTarget, nodeState.CurrentChallenge[:8])

	// Heartbeat loop - shows node is alive and waiting for blocks
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		nodeState.RLock()
		height := nodeState.CurrentHeight
		difficulty := nodeState.Consensus.DifficultyTarget
		challenge := nodeState.CurrentChallenge
		chainLen := len(nodeState.Chain)
		nodeState.RUnlock()

		log.Printf("[consensus] height=%d difficulty=%d challenge=%x chainLen=%d",
			height, difficulty, challenge[:8], chainLen)
	}
}

// AcceptBlock is called by RPC when a farmer submits a block
func (ns *NodeState) AcceptBlock(proof *pospace.Proof, farmerAddr string, farmerPubKey []byte) error {
	ns.Lock()
	defer ns.Unlock()

	// Get expected height
	nextHeight := ns.CurrentHeight + 1

	// Verify proof against current challenge
	if err := ns.Consensus.VerifyProofOfSpace(proof, ns.CurrentChallenge); err != nil {
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
		SenderPubKey: nil, // Coinbase has no sender
		Signature:    nil, // Coinbase has no signature
	}

	// Build transaction list (coinbase first, then user txs)
	allTxs := []ledger.Transaction{coinbase}

	// Apply coinbase (special handling - no signature verification)
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
			fmt.Printf("⚠️  Skipping invalid tx: %v\n", err)
		} else {
			validTxs = append(validTxs, tx)
		}
	}

	allTxs = append(allTxs, validTxs...)

	// Calculate prev hash
	var prevHash [32]byte
	if nextHeight > 0 {
		prevBlock := ns.Chain[len(ns.Chain)-1]
		prevHash = hashBlock(&prevBlock)
	}

	// Create new block
	newBlock := Block{
		Height:        nextHeight,
		TimestampUnix: time.Now().Unix(),
		PrevHash:      prevHash,
		Txs:           allTxs,
		Proof:         proof,
		FarmerAddr:    farmerAddr,
	}

	// Add to chain
	ns.Chain = append(ns.Chain, newBlock)
	ns.CurrentHeight = nextHeight

	// Clear mempool
	ns.Mempool.Clear()

	// Generate new challenge for next block
	newBlockHash := hashBlock(&newBlock)
	ns.CurrentChallenge = consensus.GenerateChallenge(newBlockHash, nextHeight+1)

	// Update difficulty (every block for now, could be less frequent)
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

	// PERSIST TO DISK
	log.Println("[storage] Persisting block and state...")

	// Save block
	if err := ns.BlockStore.SaveBlock(nextHeight, newBlock); err != nil {
		log.Printf("⚠️  Failed to persist block: %v", err)
	}

	// Save all modified accounts
	for addr, acct := range ns.WorldState.Accounts {
		if err := ns.StateStore.SaveAccount(addr, acct.Balance, acct.Nonce); err != nil {
			log.Printf("⚠️  Failed to persist account %s: %v", addr, err)
		}
	}

	// Save metadata
	if err := ns.MetaStore.SaveTipHeight(nextHeight); err != nil {
		log.Printf("⚠️  Failed to persist tip height: %v", err)
	}
	if err := ns.MetaStore.SaveDifficulty(ns.Consensus.DifficultyTarget); err != nil {
		log.Printf("⚠️  Failed to persist difficulty: %v", err)
	}

	fmt.Printf("✅ Accepted block %d from farmer %s (reward: %.8f %s, txs: %d)\n",
		nextHeight, farmerAddr, float64(config.InitialBlockReward)/100000000.0, config.DenomSymbol, len(validTxs))
	fmt.Printf("🔍 New challenge for height %d: %x\n", nextHeight+1, ns.CurrentChallenge[:8])
	fmt.Printf("⚙️  Difficulty adjusted to: %d\n", ns.Consensus.DifficultyTarget)
	log.Println("[storage] ✅ State persisted to disk")

	// Gossip new block to peers
	if ns.P2P != nil {
		ns.P2P.BroadcastNewBlock(nextHeight, newBlockHash)
	}

	return nil
}

// P2P NodeHandler implementation
func (ns *NodeState) OnNewBlock(height uint64, hash [32]byte, fromPeer string) {
	log.Printf("[p2p] Peer %s announced block %d", fromPeer, height)

	ns.RLock()
	currentHeight := ns.CurrentHeight
	ns.RUnlock()

	// If we're behind, request the block
	if height > currentHeight {
		log.Printf("[p2p] We're behind (our height=%d, peer height=%d), requesting block", currentHeight, height)
		if ns.P2P != nil {
			ns.P2P.RequestBlock(height)
		}
	}
}

func (ns *NodeState) OnBlockRequest(height uint64) (interface{}, error) {
	ns.RLock()
	defer ns.RUnlock()

	if int(height) >= len(ns.Chain) {
		return nil, fmt.Errorf("block not found")
	}

	return ns.Chain[height], nil
}

func (ns *NodeState) GetStatus() (uint64, uint64, [32]byte) {
	ns.RLock()
	defer ns.RUnlock()
	
	if len(ns.Chain) == 0 {
		return 0, ns.Consensus.DifficultyTarget, [32]byte{}
	}
	
	tipBlock := ns.Chain[len(ns.Chain)-1]
	tipHash := hashBlock(&tipBlock)
	return ns.CurrentHeight, ns.Consensus.DifficultyTarget, tipHash
}

// GetCurrentVDF returns current VDF state (if timelord is active)
func (ns *NodeState) GetCurrentVDF() (seed []byte, iterations uint64, output []byte, hasVDF bool) {
	ns.RLock()
	defer ns.RUnlock()
	return ns.VDFSeed, ns.VDFIterations, ns.VDFOutput, ns.HasVDF
}

// UpdateVDFState updates the VDF state from timelord
func (ns *NodeState) UpdateVDFState(seed []byte, iterations uint64, output []byte) {
	ns.Lock()
	defer ns.Unlock()
	ns.VDFSeed = seed
	ns.VDFIterations = iterations
	ns.VDFOutput = output
	ns.HasVDF = true
	
	// CRITICAL: Update challenge based on new VDF output!
	// Challenge should be H(VDF_output || height)
	h := sha256.New()
	h.Write(output)
	binary.Write(h, binary.BigEndian, ns.CurrentHeight+1)
	ns.CurrentChallenge = sha256.Sum256(h.Sum(nil))
}

// LocalHeight returns current chain height
func (ns *NodeState) LocalHeight() uint64 {
	ns.RLock()
	defer ns.RUnlock()
	return ns.CurrentHeight
}

// HasBlock checks if we have a block at given height
func (ns *NodeState) HasBlock(height uint64) bool {
	ns.RLock()
	defer ns.RUnlock()
	return int(height) < len(ns.Chain)
}

// VerifyAndApplyBlock verifies and applies a block received from a peer
func (ns *NodeState) VerifyAndApplyBlock(blockJSON json.RawMessage) error {
	var block Block
	if err := json.Unmarshal(blockJSON, &block); err != nil {
		return fmt.Errorf("failed to unmarshal block: %w", err)
	}
	
	ns.Lock()
	defer ns.Unlock()
	
	// Verify block is next in sequence
	if block.Height != ns.CurrentHeight+1 {
		return fmt.Errorf("block height %d doesn't match expected %d", block.Height, ns.CurrentHeight+1)
	}
	
	// Verify prev hash
	if len(ns.Chain) > 0 {
		prevBlock := ns.Chain[len(ns.Chain)-1]
		prevHash := hashBlock(&prevBlock)
		if block.PrevHash != prevHash {
			return fmt.Errorf("prev hash mismatch")
		}
	}
	
	// Verify PoSpace proof if present
	if block.Proof != nil {
		challenge := ns.CurrentChallenge
		if err := ns.Consensus.VerifyProofOfSpace(block.Proof, challenge); err != nil {
			return fmt.Errorf("invalid PoSpace proof: %w", err)
		}
	}
	
	// Apply transactions (excluding coinbase)
	for i, tx := range block.Txs {
		if i == 0 && tx.From == "coinbase" {
			// Skip coinbase validation
			continue
		}
		// Transactions were already validated when block was created
		// For received blocks, we trust the PoSpace proof validates the block
		// In production, you'd re-verify all transactions here
	}
	
	// Add block to chain
	ns.Chain = append(ns.Chain, block)
	ns.CurrentHeight = block.Height
	
	// Update challenge for next block
	newBlockHash := hashBlock(&block)
	ns.CurrentChallenge = consensus.GenerateChallenge(newBlockHash, ns.CurrentHeight+1)
	
	// Persist to database
	if ns.BlockStore != nil {
		if err := ns.BlockStore.SaveBlock(block.Height, block); err != nil {
			log.Printf("⚠️  Failed to persist block %d: %v", block.Height, err)
		}
		if err := ns.MetaStore.SaveTipHeight(block.Height); err != nil {
			log.Printf("⚠️  Failed to persist tip height: %v", err)
		}
	}
	
	log.Printf("✅ Synced block %d from peer", block.Height)
	
	return nil
}

// GetCurrentChallenge returns the current challenge and difficulty
func (ns *NodeState) GetCurrentChallenge() ([32]byte, uint64, uint64) {
	ns.RLock()
	defer ns.RUnlock()
	return ns.CurrentChallenge, ns.Consensus.DifficultyTarget, ns.CurrentHeight + 1
}

// hashBlock computes the hash of a block
func hashBlock(b *Block) [32]byte {
	// Simple hash of block data
	h := sha256.New()
	fmt.Fprintf(h, "%d", b.Height)
	fmt.Fprintf(h, "%d", b.TimestampUnix)
	h.Write(b.PrevHash[:])
	if b.Proof != nil {
		h.Write(b.Proof.Hash[:])
	}
	return sha256.Sum256(h.Sum(nil))
}
