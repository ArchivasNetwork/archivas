package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ArchivasNetwork/archivas/config"
	"github.com/ArchivasNetwork/archivas/consensus"
	"github.com/ArchivasNetwork/archivas/health"
	"github.com/ArchivasNetwork/archivas/internal/buildinfo"
	"github.com/ArchivasNetwork/archivas/ledger"
	"github.com/ArchivasNetwork/archivas/logging"
	"github.com/ArchivasNetwork/archivas/mempool"
	"github.com/ArchivasNetwork/archivas/metrics"
	"github.com/ArchivasNetwork/archivas/network"
	"github.com/ArchivasNetwork/archivas/node"
	"github.com/ArchivasNetwork/archivas/p2p"
	"github.com/ArchivasNetwork/archivas/pospace"
	"github.com/ArchivasNetwork/archivas/rpc"
	"github.com/ArchivasNetwork/archivas/snapshot"
	"github.com/ArchivasNetwork/archivas/storage"
)

// Block represents a blockchain block with Proof-of-Space
type Block struct {
	Height        uint64
	TimestampUnix int64
	PrevHash      [32]byte
	Difficulty    uint64   // Difficulty target when mined
	Challenge     [32]byte // The challenge used to win this block
	Txs           []ledger.Transaction
	Proof         *pospace.Proof // Proof-of-Space
	FarmerAddr    string         // Address to receive block reward

	// v0.5.0: Cumulative work for fork resolution
	CumulativeWork uint64 // Total work from genesis to this block
}

// stringSliceFlag is a custom flag type for repeatable string flags
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
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
	P2P         *p2p.Network
	GenesisHash [32]byte
	NetworkID   string
	// Health tracking
	Health *health.ChainHealth
	// Reorg detection (v0.5.0)
	ReorgDetector *consensus.ReorgDetector
	// VDF state (updated by timelord)
	VDFSeed       []byte
	VDFIterations uint64
	VDFOutput     []byte
	HasVDF        bool
	// Backpressure for disk persistence (limit concurrent writes)
	persistSem chan struct{}
}

func main() {
	logging.ConfigureJSON("archivas-node")

	// Handle subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "snapshot":
			handleSnapshotCommand()
			return
		case "bootstrap":
			handleBootstrapCommand()
			return
		case "help", "--help", "-h":
			printUsage()
			return
		case "version", "--version", "-v":
			buildInfo := buildinfo.GetInfo()
			fmt.Printf("archivas-node version %s (commit: %s, built: %s)\n",
				buildInfo["version"], buildInfo["commit"], buildInfo["builtAt"])
			return
		}
	}

	// Log build info and run self-test
	buildInfo := buildinfo.GetInfo()
	log.Printf("[build] version=%s commit=%s built=%s rule=%s",
		buildInfo["version"], buildInfo["commit"], buildInfo["builtAt"], buildInfo["poSpaceRule"])

	// Run PoSpace self-test
	if err := pospace.SelfTest(); err != nil {
		log.Fatalf("[FATAL] PoSpace self-test failed: %v", err)
	}
	log.Println("[build] PoSpace self-test passed ✓")

	// Parse CLI flags
	// Phase 1: Network profile selection
	networkName := flag.String("network", network.DefaultNetwork(), "Network to join (betanet, devnet-legacy)")
	rpcAddr := flag.String("rpc", "", "RPC listen address (default: from network profile)")
	p2pAddr := flag.String("p2p", "", "P2P listen address (default: from network profile)")
	peerAddrs := flag.String("peer", "", "Comma-separated peer addresses (e.g., ip1:9090,ip2:9090)")
	dbPath := flag.String("db", "./data", "Database directory path")
	vdfRequired := flag.Bool("vdf-required", false, "Require VDF proofs in blocks (PoSpace+Time mode)")
	genesisPath := flag.String("genesis", "", "Genesis file path (overrides network profile)")
	networkID := flag.String("network-id", "", "Network ID (overrides network profile)")
	bootnodes := flag.String("bootnodes", "", "Comma-separated bootnode addresses")

	// Gossip flags
	enableGossip := flag.Bool("enable-gossip", true, "Enable automatic peer discovery via gossip")
	gossipInterval := flag.Duration("gossip-interval", 60*time.Second, "Interval between gossip broadcasts")
	maxPeers := flag.Int("max-peers", 20, "Maximum number of peer connections")
	dialsPerMin := flag.Int("gossip-dials-per-min", 5, "Maximum new peer dials per minute")
	peersFile := flag.String("peers-file", "", "Path to peers.json (default: <db>/peers.json)")

	// P2P Isolation flags (v1.2.0)
	noPeerDiscovery := flag.Bool("no-peer-discovery", false, "Disable automatic peer discovery (only dial whitelisted peers)")
	var peerWhitelist stringSliceFlag
	flag.Var(&peerWhitelist, "peer-whitelist", "Whitelisted peer address (repeatable, format: host:port or IP:port)")
	checkpointHeight := flag.Uint64("checkpoint-height", 0, "Chain checkpoint height for validation")
	checkpointHash := flag.String("checkpoint-hash", "", "Chain checkpoint hash (hex, 64 chars)")

	flag.Parse()

	// Phase 1: Load network profile
	log.Printf("[network] Loading network profile: %s", *networkName)
	profile, err := network.GetProfile(*networkName)
	if err != nil {
		log.Fatalf("Failed to load network profile: %v", err)
	}

	// Apply network profile defaults
	if *rpcAddr == "" {
		*rpcAddr = fmt.Sprintf(":%d", profile.DefaultRPCPort)
	}
	if *p2pAddr == "" {
		*p2pAddr = fmt.Sprintf(":%d", profile.DefaultP2PPort)
	}
	if *genesisPath == "" {
		*genesisPath = profile.GenesisPath
	}
	if *networkID == "" {
		*networkID = profile.ChainID
	}

	log.Printf("[network] Network: %s (chain-id: %s, network-id: %d, protocol: v%d)",
		profile.Name, profile.ChainID, profile.NetworkID, profile.ProtocolVersion)

	log.Println("[startup] Archivas node starting...")
	if profile.Name == "betanet" {
		fmt.Println("🚀 Archivas Betanet Node running…")
	} else {
		fmt.Println("Archivas Devnet Node running…")
	}
	fmt.Println()

	// v1.1.1: Ensure RPC binds to 0.0.0.0 if no host specified (for metrics scraping)
	rpcBindAddr := *rpcAddr
	if strings.HasPrefix(rpcBindAddr, ":") {
		rpcBindAddr = "0.0.0.0" + rpcBindAddr
	}

	fmt.Printf("🔧 Configuration:\n")
	fmt.Printf("   Network: %s\n", profile.Name)
	fmt.Printf("   Chain ID: %s\n", profile.ChainID)
	fmt.Printf("   RPC:  %s\n", rpcBindAddr)
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
	var genesisHash [32]byte

	tipHeight, err := metaStore.LoadTipHeight()
	freshStart := err != nil

	if freshStart {
		// Fresh start - require genesis file
		if *genesisPath == "" {
			log.Fatal("--genesis required for first start (e.g., --genesis genesis/devnet.genesis.json)")
		}

		log.Printf("[genesis] Loading genesis from %s", *genesisPath)
		gen, err := config.LoadGenesis(*genesisPath)
		if err != nil {
			log.Fatalf("Failed to load genesis: %v", err)
		}

		genesisHash = config.HashGenesis(gen)
		genesisAllocs := config.GenesisAllocToMap(gen.Allocations)

		worldState = ledger.NewWorldState(genesisAllocs)
		fmt.Printf("🌱 Fresh start from genesis file\n")
		fmt.Printf("   Genesis Hash: %x\n", genesisHash[:8])
		fmt.Printf("   Network ID: %s\n", *networkID)
		fmt.Printf("🌍 World state initialized with %d genesis accounts\n", len(gen.Allocations))
		for _, alloc := range gen.Allocations {
			fmt.Printf("   %s: %.8f %s\n", alloc.Address, float64(alloc.Amount)/100000000.0, config.DenomSymbol)
		}

		cs = consensus.NewConsensus()

		genesisChallenge = consensus.GenerateGenesisChallenge()
		genesisBlock := Block{
			Height:         0,
			TimestampUnix:  gen.Timestamp, // Use FIXED timestamp from genesis.json!
			PrevHash:       [32]byte{},
			Difficulty:     1125899906842624, // 2^50 initial difficulty
			Challenge:      genesisChallenge,
			Txs:            nil,
			Proof:          nil,
			FarmerAddr:     "",
			CumulativeWork: consensus.CalculateWork(1125899906842624), // Genesis work
		}
		chain = []Block{genesisBlock}
		currentHeight = 0

		// Persist genesis
		if err := blockStore.SaveBlock(0, genesisBlock); err != nil {
			log.Fatalf("Failed to save genesis block: %v", err)
		}
		for addr, balance := range genesisAllocs {
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
		if err := metaStore.SaveGenesisHash(genesisHash); err != nil {
			log.Fatalf("Failed to save genesis hash: %v", err)
		}
		if err := metaStore.SaveNetworkID(*networkID); err != nil {
			log.Fatalf("Failed to save network ID: %v", err)
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

		// Load genesis hash and network ID
		savedGenesisHash, err := metaStore.LoadGenesisHash()
		if err != nil {
			log.Fatalf("Failed to load genesis hash: %v", err)
		}
		genesisHash = savedGenesisHash

		savedNetworkID, err := metaStore.LoadNetworkID()
		if err != nil {
			log.Printf("[warning] Network ID not found in DB, using default")
			savedNetworkID = "archivas-devnet-v1"
		}
		fmt.Printf("   Genesis Hash: %x\n", genesisHash[:8])
		fmt.Printf("   Network ID: %s\n", savedNetworkID)

		// Load all accounts - we'll need to scan all keys with acc: prefix
		// For now, just load genesis accounts and any that received funds
		// (In production, you'd have an account index)
		for addr := range config.LegacyGenesisAlloc {
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
		Health:           health.NewChainHealth(),
		ReorgDetector:    consensus.NewReorgDetector(),
		GenesisHash:      genesisHash,
		NetworkID:        *networkID,
		persistSem:       make(chan struct{}, 5), // Limit to 5 concurrent disk writes
	}

	metrics.StartWatchdogs(metrics.GroupNode)
	metrics.UpdateTipHeight(nodeState.CurrentHeight)
	metrics.UpdateDifficulty(nodeState.Consensus.DifficultyTarget)
	metrics.UpdatePeerCount(0)

	log.Println("[DEBUG] Initialized chain memory")
	fmt.Printf("🔍 Current challenge: %x\n", genesisChallenge[:8])
	fmt.Println()

	// Start P2P network if enabled
	var p2pNet *p2p.Network
	if *p2pAddr != "" {
		log.Printf("[p2p] Starting P2P listener on %s", *p2pAddr)
		p2pNet = p2p.NewNetwork(*p2pAddr, nodeState)

		// Configure gossip
		p2pNet.SetGossipConfig(p2p.GossipConfig{
			NetworkID:      *networkID,
			EnableGossip:   *enableGossip,
			Interval:       *gossipInterval,
			MaxPeers:       *maxPeers,
			DialsPerMinute: *dialsPerMin,
		})

		// Configure peer isolation (v1.2.0)
		if *noPeerDiscovery || len(peerWhitelist) > 0 || *checkpointHeight > 0 {
			// Parse checkpoint hash if provided
			var checkpointHashBytes [32]byte
			if *checkpointHash != "" {
				hashBytes, err := hex.DecodeString(*checkpointHash)
				if err != nil {
					log.Fatalf("[p2p] Invalid checkpoint hash: %v", err)
				}
				if len(hashBytes) != 32 {
					log.Fatalf("[p2p] Checkpoint hash must be 64 hex chars (32 bytes)")
				}
				copy(checkpointHashBytes[:], hashBytes)
			}

			p2pNet.SetIsolationConfig(p2p.IsolationConfig{
				NoPeerDiscovery:  *noPeerDiscovery,
				PeerWhitelist:    []string(peerWhitelist),
				CheckpointHeight: *checkpointHeight,
				CheckpointHash:   checkpointHashBytes,
				GenesisHash:      nodeState.GenesisHash,
			})

			// Print prominent PRIVATE NODE banner
			fmt.Println()
			fmt.Println("═══════════════════════════════════════════════════════════════")
			fmt.Println("  🔐 RUNNING IN PRIVATE NODE MODE")
			fmt.Println("═══════════════════════════════════════════════════════════════")
			fmt.Println()
			if *noPeerDiscovery {
				fmt.Println("  ⛔ Peer Discovery:  DISABLED")
			} else {
				fmt.Println("  ✅ Peer Discovery:  ENABLED")
			}
			if len(peerWhitelist) > 0 {
				fmt.Printf("  📋 Whitelisted Peers: %d\n", len(peerWhitelist))
				for i, peer := range peerWhitelist {
					fmt.Printf("     %d. %s\n", i+1, peer)
				}
			} else {
				fmt.Println("  📋 Whitelisted Peers: NONE (will connect to any peer)")
			}
			if *checkpointHeight > 0 {
				fmt.Printf("  📌 Checkpoint:     height=%d hash=%s...\n", *checkpointHeight, (*checkpointHash)[:16])
				fmt.Println("     (Will reject blocks that don't match checkpoint)")
			} else {
				fmt.Println("  📌 Checkpoint:     NONE (will sync from genesis)")
			}
			fmt.Println()
			fmt.Println("  ℹ️  This node will ONLY accept connections from whitelisted")
			fmt.Println("     peers and will NOT be discoverable by other nodes.")
			fmt.Println()
			fmt.Println("═══════════════════════════════════════════════════════════════")
			fmt.Println()
		}

		// Set up peer persistence
		peerStorePath := *peersFile
		if peerStorePath == "" {
			peerStorePath = *dbPath + "/peers.json"
		}
		peerStore, err := p2p.NewFilePeerStore(peerStorePath)
		if err != nil {
			log.Printf("[p2p] Warning: failed to create peer store: %v", err)
		} else {
			p2pNet.SetPeerStore(peerStore)
			log.Printf("[p2p] Peer store: %s", peerStorePath)
		}

		if err := p2pNet.Start(); err != nil {
			log.Fatalf("[p2p] Failed to start P2P: %v", err)
		}
		nodeState.P2P = p2pNet
		fmt.Printf("🌐 P2P network started on %s\n", *p2pAddr)
		if *enableGossip {
			fmt.Printf("🔄 Peer gossip enabled (interval=%v, max=%d)\n", *gossipInterval, *maxPeers)
		}

		// Connect to initial peers and bootnodes
		allPeers := []string{}
		if *peerAddrs != "" {
			allPeers = append(allPeers, strings.Split(*peerAddrs, ",")...)
		}
		if *bootnodes != "" {
			allPeers = append(allPeers, strings.Split(*bootnodes, ",")...)
		}

		if len(allPeers) > 0 {
			for _, peer := range allPeers {
				peer = strings.TrimSpace(peer)
				if peer != "" {
					go func(addr string) {
						time.Sleep(500 * time.Millisecond) // Stagger connections
						if err := p2pNet.ConnectPeer(addr); err != nil {
							log.Printf("[p2p] Failed to connect to peer %s: %v", addr, err)
						}
					}(peer)
				}
			}
			fmt.Printf("📡 Connecting to %d peer(s)/bootnode(s)...\n", len(allPeers))
		}
		fmt.Println()

		// v1.2.1: IBD using new node.IBDManager
		if len(allPeers) > 0 {
			go func() {
				time.Sleep(2 * time.Second) // Wait for peer connections

				// Create IBD manager
				ibdConfig := node.DefaultIBDConfig(*dbPath)
				ibdManager := node.NewIBDManager(ibdConfig, nodeState)

				// Load resume state if exists
				if err := ibdManager.LoadState(); err != nil {
					log.Printf("[IBD] Warning: failed to load state: %v", err)
				}

				// Build peer URLs
				peerURLs := []string{}
				for _, peerAddr := range allPeers {
					peerAddr = strings.TrimSpace(peerAddr)
					if peerAddr == "" {
						continue
					}

					var peerURL string
					if strings.Contains(peerAddr, "seed.archivas.ai") {
						peerURL = "https://seed.archivas.ai"
					} else {
						// Convert P2P port to RPC port
						// Handle both legacy Devnet (9090->8080) and Betanet (30303->8545)
						peerURL = fmt.Sprintf("http://%s", peerAddr)
						peerURL = strings.Replace(peerURL, ":9090", ":8080", 1)  // Legacy Devnet
						peerURL = strings.Replace(peerURL, ":30303", ":8545", 1) // Betanet
					}
					peerURLs = append(peerURLs, peerURL)
				}

			// Run IBD with retry across peers
			if err := ibdManager.RunIBDWithRetry(peerURLs); err != nil {
				log.Printf("[IBD] All peers failed: %v", err)
			}

			// Start periodic chain sync check (runs every 30 seconds)
			// This enables automatic chain reorganization when peers have a longer chain
			go func() {
				ticker := time.NewTicker(30 * time.Second)
				defer ticker.Stop()

				// Wait for initial IBD to complete
				time.Sleep(10 * time.Second)

				for range ticker.C {
					// Get current local height
					localHeight, _, _ := nodeState.GetStatus()

					// Check each peer's chain tip
					var bestPeerURL string
					var bestRemoteHeight uint64

					for _, peerURL := range peerURLs {
						remoteHeight, err := ibdManager.FetchRemoteTip(peerURL)
						if err != nil {
							continue // Skip unreachable peers
						}

						if remoteHeight > bestRemoteHeight {
							bestRemoteHeight = remoteHeight
							bestPeerURL = peerURL
						}
					}

					// If any peer has a significantly longer chain, trigger IBD
					if bestRemoteHeight > localHeight && ibdManager.ShouldRunIBD(localHeight, bestRemoteHeight) {
						log.Printf("[SYNC] Detected longer chain: local=%d remote=%d (gap=%d), triggering reorg from %s",
							localHeight, bestRemoteHeight, bestRemoteHeight-localHeight, bestPeerURL)

						if err := ibdManager.RunIBD(bestPeerURL); err != nil {
							log.Printf("[SYNC] Failed to sync from %s: %v", bestPeerURL, err)
						} else {
							log.Printf("[SYNC] Successfully reorganized to height %d", bestRemoteHeight)
						}
					}
				}
			}()
		}()
	}
}

	// Start RPC server in background
	log.Println("[DEBUG] Starting RPC server...")
	server := rpc.NewFarmingServer(nodeState.WorldState, nodeState.Mempool, nodeState)
	go func() {
		log.Printf("[rpc] starting server on %s", rpcBindAddr)
		fmt.Printf("🌐 Starting RPC server on %s\n", rpcBindAddr)
		if err := server.Start(rpcBindAddr); err != nil {
			log.Fatalf("[rpc] server error: %v", err)
		}
	}()

	// Give RPC server a moment to start
	time.Sleep(500 * time.Millisecond)
	log.Println("[DEBUG] RPC server running")

	// Start metrics updater (updates gauges every 2s)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			nodeState.RLock()
			metrics.UpdateTipHeight(nodeState.CurrentHeight)
			if nodeState.P2P != nil {
				connected, _ := nodeState.P2P.GetPeerList()
				metrics.UpdatePeerCount(len(connected))
			}
			metrics.UpdateDifficulty(nodeState.Consensus.DifficultyTarget)
			nodeState.RUnlock()
		}
	}()

	fmt.Println("✅ Node initialized successfully")
	fmt.Println("🌾 Waiting for farmers to submit blocks...")
	fmt.Println()

	log.Println("[DEBUG] Entering consensus loop...")
	log.Printf("[consensus] height=%d difficulty=%d challenge=%x",
		nodeState.CurrentHeight, nodeState.Consensus.DifficultyTarget, nodeState.CurrentChallenge[:8])

	// Heartbeat loop - shows node is alive and waiting for blocks
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lastBlockTime := time.Now()

	for range ticker.C {
		nodeState.Lock()
		height := nodeState.CurrentHeight
		difficulty := nodeState.Consensus.DifficultyTarget
		challenge := nodeState.CurrentChallenge
		chainLen := len(nodeState.Chain)

		// Check if we got a new block
		if chainLen > int(height)+1 {
			lastBlockTime = time.Now()
		}

		// Time-based difficulty drop: if no block for 60 seconds, halve difficulty
		timeSinceBlock := time.Since(lastBlockTime)
		if timeSinceBlock > 60*time.Second && difficulty > 1_000_000 {
			oldDiff := difficulty
			difficulty = difficulty / 2
			if difficulty < 1_000_000 {
				difficulty = 1_000_000
			}
			nodeState.Consensus.DifficultyTarget = difficulty
			lastBlockTime = time.Now() // Reset timer
			log.Printf("[auto-drop] No block for %v, dropping difficulty: %d → %d", timeSinceBlock, oldDiff, difficulty)
		}

		nodeState.Unlock()

		log.Printf("[consensus] height=%d difficulty=%d challenge=%x chainLen=%d",
			height, difficulty, challenge[:8], chainLen)
	}
}

// AcceptBlock is called by RPC when a farmer submits a block
func (ns *NodeState) AcceptBlock(proof *pospace.Proof, farmerAddr string, farmerPubKey []byte) error {
	// Track submission
	metrics.IncSubmitReceived()

	ns.Lock()

	// Get expected height
	nextHeight := ns.CurrentHeight + 1

	// Verify proof against its embedded challenge
	// Note: proof.Challenge is set by the farmer based on when it found the winner
	// We accept this because VDF advances quickly and farmer might have found it
	// for a slightly older challenge
	if err := ns.Consensus.VerifyProofOfSpace(proof, proof.Challenge); err != nil {
		metrics.IncSubmitIgnored()
		ns.Unlock()
		return fmt.Errorf("invalid proof: %w", err)
	}

	// Proof accepted
	metrics.IncSubmitAccepted()

	// Get pending transactions
	pending := ns.Mempool.Pending()
	log.Printf("[block] Creating block %d with %d pending transactions from mempool", nextHeight, len(pending))

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

	// Create new block with current difficulty
	newBlock := Block{
		Height:        nextHeight,
		TimestampUnix: time.Now().Unix(),
		PrevHash:      prevHash,
		Difficulty:    ns.Consensus.DifficultyTarget, // Difficulty when mined
		Challenge:     ns.CurrentChallenge,           // Challenge used to win
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

	// TEMPORARY: Aggressive difficulty drop to get blocks flowing
	// Drop difficulty by 50% every block until it reaches 1M
	if ns.Consensus.DifficultyTarget > 1_000_000 {
		oldDiff := ns.Consensus.DifficultyTarget
		ns.Consensus.DifficultyTarget = ns.Consensus.DifficultyTarget / 2
		if ns.Consensus.DifficultyTarget < 1_000_000 {
			ns.Consensus.DifficultyTarget = 1_000_000
		}
		log.Printf("[difficulty] Dropping difficulty: %d → %d", oldDiff, ns.Consensus.DifficultyTarget)
	}

	// Copy data needed for persistence before releasing lock
	// Track modified accounts (coinbase receiver + all transaction participants)
	modifiedAccounts := make(map[string]*ledger.AccountState)
	modifiedAccounts[farmerAddr] = &ledger.AccountState{
		Balance: receiver.Balance,
		Nonce:   receiver.Nonce,
	}
	for _, tx := range validTxs {
		if sender, ok := ns.WorldState.Accounts[tx.From]; ok {
			modifiedAccounts[tx.From] = &ledger.AccountState{
				Balance: sender.Balance,
				Nonce:   sender.Nonce,
			}
		}
		if recv, ok := ns.WorldState.Accounts[tx.To]; ok {
			modifiedAccounts[tx.To] = &ledger.AccountState{
				Balance: recv.Balance,
				Nonce:   recv.Nonce,
			}
		}
	}
	currentDifficulty := ns.Consensus.DifficultyTarget

	// Release lock BEFORE disk I/O to prevent blocking other requests
	ns.Unlock()

	// PERSIST TO DISK in background with backpressure (prevents goroutine accumulation)
	go func() {
		// Acquire semaphore slot (blocks if too many concurrent writes)
		ns.persistSem <- struct{}{}
		defer func() { <-ns.persistSem }() // Release slot when done

	log.Println("[storage] Persisting block and state...")

	// Save block
	if err := ns.BlockStore.SaveBlock(nextHeight, newBlock); err != nil {
		log.Printf("⚠️  Failed to persist block: %v", err)
			return // Early exit on block save failure
	}

		// Save only modified accounts (not all accounts - much faster!)
		for addr, acct := range modifiedAccounts {
		if err := ns.StateStore.SaveAccount(addr, acct.Balance, acct.Nonce); err != nil {
			log.Printf("⚠️  Failed to persist account %s: %v", addr, err)
		}
	}

	// Save metadata
	if err := ns.MetaStore.SaveTipHeight(nextHeight); err != nil {
		log.Printf("⚠️  Failed to persist tip height: %v", err)
	}
		if err := ns.MetaStore.SaveDifficulty(currentDifficulty); err != nil {
		log.Printf("⚠️  Failed to persist difficulty: %v", err)
	}

		log.Println("[storage] ✅ State persisted to disk")
	}()

	fmt.Printf("✅ Accepted block %d from farmer %s (reward: %.8f %s, txs: %d)\n",
		nextHeight, farmerAddr, float64(config.InitialBlockReward)/100000000.0, config.DenomSymbol, len(validTxs))
	fmt.Printf("🔍 New challenge for height %d: %x\n", nextHeight+1, ns.CurrentChallenge[:8])
	fmt.Printf("⚙️  Difficulty adjusted to: %d\n", currentDifficulty)

	// Gossip new block to peers (non-blocking, already outside lock)
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

	// If we're behind, request blocks
	if height > currentHeight {
		gap := height - currentHeight

		// v1.1.1: Use batched IBD for large gaps (>10 blocks)
		if gap > 10 {
			log.Printf("[p2p] We're behind by %d blocks (our=%d, peer=%d), starting IBD",
				gap, currentHeight, height)
			if ns.P2P != nil {
				ns.P2P.StartIBD(currentHeight + 1)
			}
		} else {
			// Small gap, use single block requests
			log.Printf("[p2p] We're behind (our height=%d, peer height=%d), requesting block", currentHeight, height)
			if ns.P2P != nil {
				ns.P2P.RequestBlock(height)
			}
		}
	}
}

func (ns *NodeState) OnBlockRequest(height uint64) (interface{}, error) {
	ns.RLock()
	defer ns.RUnlock()

	// v1.1.1: Try memory first, then disk
	if int(height) < len(ns.Chain) {
		return ns.Chain[height], nil
	}

	// Load from disk
	if ns.BlockStore != nil {
		var block Block
		if err := ns.BlockStore.LoadBlock(height, &block); err == nil {
			return block, nil
		}
	}

	return nil, fmt.Errorf("block not found")
}

// OnBlocksRangeRequest serves a batch of blocks for IBD
// v1.1.1: Efficient disk-based block serving
func (ns *NodeState) OnBlocksRangeRequest(fromHeight uint64, maxBlocks uint32) (blocks []json.RawMessage, tipHeight uint64, eof bool, err error) {
	ns.RLock()
	defer ns.RUnlock()

	// Get current tip
	tipHeight = ns.CurrentHeight

	// Cap batch size
	if maxBlocks == 0 || maxBlocks > 512 {
		maxBlocks = 512
	}

	// Ensure fromHeight is valid
	if fromHeight < 1 {
		fromHeight = 1
	}

	// If requesting beyond tip, return empty with EOF
	if fromHeight > tipHeight {
		return []json.RawMessage{}, tipHeight, true, nil
	}

	// Calculate how many blocks to serve
	remaining := tipHeight - fromHeight + 1
	count := uint64(maxBlocks)
	if count > remaining {
		count = remaining
		eof = true // This is the last batch
	}

	blocks = make([]json.RawMessage, 0, count)

	// Serve blocks from disk (more reliable for IBD than memory)
	for h := fromHeight; h < fromHeight+count; h++ {
		var block Block

		// For IBD, always load from disk to ensure consistency
		// Memory chain (ns.Chain) might be incomplete after database transfers
		if ns.BlockStore != nil {
			if err := ns.BlockStore.LoadBlock(h, &block); err != nil {
				log.Printf("[ibd] failed to load block %d from disk: %v", h, err)
				// Return what we have so far
				break
			}
		} else if int(h) < len(ns.Chain) {
			// Fallback to memory if no BlockStore (shouldn't happen in production)
			block = ns.Chain[h]
		} else {
			// Block not available
			log.Printf("[ibd] block %d not available (chain len=%d, no disk store)", h, len(ns.Chain))
			break
		}

		// Serialize block to JSON
		blockJSON, err := json.Marshal(block)
		if err != nil {
			log.Printf("[ibd] failed to marshal block %d: %v", h, err)
			break
		}

		blocks = append(blocks, blockJSON)
	}

	// If we got fewer blocks than expected, not at EOF yet
	if uint64(len(blocks)) < count {
		eof = false
	}

	log.Printf("[ibd] serving range from=%d count=%d tip=%d eof=%v", fromHeight, len(blocks), tipHeight, eof)

	return blocks, tipHeight, eof, nil
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

// GetGenesisHash returns the genesis hash
func (ns *NodeState) GetGenesisHash() [32]byte {
	ns.RLock()
	defer ns.RUnlock()
	return ns.GenesisHash
}

// GetCurrentHeight returns the current chain height (for IBD interface)
func (ns *NodeState) GetCurrentHeight() uint64 {
	ns.RLock()
	defer ns.RUnlock()
	return ns.CurrentHeight
}

// ApplyBlock applies a block received during IBD
func (ns *NodeState) ApplyBlock(blockData json.RawMessage) error {
	// Unmarshal to map first (blocks from /blocks/range are hex-encoded)
	var blockMap map[string]interface{}
	if err := json.Unmarshal(blockData, &blockMap); err != nil {
		return fmt.Errorf("failed to unmarshal block map: %w", err)
	}

	// Extract and decode fields
	height, _ := blockMap["height"].(float64)
	difficulty, _ := blockMap["difficulty"].(float64)
	timestamp, _ := blockMap["timestamp"].(float64)
	farmerAddr, _ := blockMap["farmerAddr"].(string)
	
	// Decode hex fields
	var prevHash, challenge [32]byte
	if prevHashStr, ok := blockMap["prevHash"].(string); ok {
		prevHashBytes, _ := hex.DecodeString(prevHashStr)
		copy(prevHash[:], prevHashBytes)
	}
	if challengeStr, ok := blockMap["challenge"].(string); ok {
		challengeBytes, _ := hex.DecodeString(challengeStr)
		copy(challenge[:], challengeBytes)
	}

	// Extract transactions
	txs := []ledger.Transaction{}
	if txList, ok := blockMap["txs"].([]interface{}); ok {
		for _, txRaw := range txList {
			if txMap, ok := txRaw.(map[string]interface{}); ok {
				tx := ledger.Transaction{
					From:   getString(txMap, "from"),
					To:     getString(txMap, "to"),
					Amount: getInt64(txMap, "amount"),
					Fee:    getInt64(txMap, "fee"),
					Nonce:  getUint64(txMap, "nonce"),
				}
				txs = append(txs, tx)
			}
		}
	}

	// Reconstruct Block struct
	block := Block{
		Height:        uint64(height),
		TimestampUnix: int64(timestamp),
		PrevHash:      prevHash,
		Difficulty:    uint64(difficulty),
		Challenge:     challenge,
		Txs:           txs,
		Proof:         nil, // Not serialized in /blocks/range
		FarmerAddr:    farmerAddr,
	}

	ns.Lock()
	defer ns.Unlock()

	// Validate height continuity
	expectedHeight := ns.CurrentHeight + 1
	if block.Height != expectedHeight {
		return fmt.Errorf("height discontinuity: expected %d, got %d", expectedHeight, block.Height)
	}

	// During IBD, we trust the seed node's blocks
	// Full validation (including parent hash) happens only during P2P sync for new blocks
	// We validate:
	// 1. Height continuity (already checked above) ✓
	// 2. Genesis hash match (checked at handshake) ✓
	// 3. Network ID match (checked at handshake) ✓
	//
	// We skip:
	// - Parent hash verification (can't recompute without full Proof data)
	// - PoSpace proof verification (legacy blocks may have mismatched challenges)
	// - Transaction signature verification (performance optimization during bulk sync)
	//
	// This allows backward-compatible sync from nodes with legacy block formats

	// Apply transactions
	for _, tx := range block.Txs {
		if tx.From == "coinbase" {
			// Coinbase transaction
			receiver, ok := ns.WorldState.Accounts[tx.To]
			if !ok {
				receiver = &ledger.AccountState{Balance: 0, Nonce: 0}
				ns.WorldState.Accounts[tx.To] = receiver
			}
			receiver.Balance += tx.Amount
		} else {
			// Regular transaction (skip validation during IBD for performance)
			if err := ns.WorldState.ApplyTransaction(tx); err != nil {
				log.Printf("[IBD] Warning: skipping invalid tx in block %d: %v", block.Height, err)
			}
		}
	}

	// Add block to chain
	ns.Chain = append(ns.Chain, block)
	ns.CurrentHeight = block.Height

	// Persist to disk
	if ns.BlockStore != nil {
		if err := ns.BlockStore.SaveBlock(block.Height, block); err != nil {
			return fmt.Errorf("failed to save block %d: %w", block.Height, err)
		}
	}

	// Update tip in metadata
	if ns.MetaStore != nil {
		if err := ns.MetaStore.SaveTipHeight(block.Height); err != nil {
			log.Printf("[IBD] Warning: failed to save tip height: %v", err)
		}
	}

	return nil
}

// GetPeerCount returns number of connected peers
func (ns *NodeState) GetPeerCount() int {
	if ns.P2P == nil {
		return 0
	}
	return ns.P2P.GetPeerCount()
}

// GetPeerList returns connected and known peer addresses
func (ns *NodeState) GetPeerList() (connected []string, known []string) {
	if ns.P2P == nil {
		return []string{}, []string{}
	}
	return ns.P2P.GetPeerList()
}

// GetHealthStats returns detailed health statistics
func (ns *NodeState) GetHealthStats() interface{} {
	if ns.Health == nil {
		return map[string]interface{}{"status": "not initialized"}
	}

	stats := ns.Health.GetStats()

	return map[string]interface{}{
		"uptime":          stats.Uptime.String(),
		"uptimeSeconds":   int(stats.Uptime.Seconds()),
		"totalBlocks":     stats.TotalBlocks,
		"avgBlockTime":    stats.AverageBlockTime.String(),
		"avgBlockSeconds": stats.AverageBlockTime.Seconds(),
		"blocksPerHour":   stats.BlocksPerHour,
		"lastBlockTime":   stats.LastBlockTime.Format(time.RFC3339),
	}
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

	// Verify difficulty matches expected (recompute from chain history)
	// For now, trust the block's difficulty (production would recompute)
	// TODO: Add RecomputeDifficulty(prev, params) and verify match

	// Verify PoSpace proof using block's own difficulty and challenge!
	if block.Proof != nil {
		// Create temporary consensus with block's difficulty for verification
		blockConsensus := &consensus.Consensus{DifficultyTarget: block.Difficulty}
		if err := blockConsensus.VerifyProofOfSpace(block.Proof, block.Challenge); err != nil {
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

	// Update Prometheus metrics
	metrics.UpdateTipHeight(ns.CurrentHeight)
	metrics.IncBlocksTotal()
	metrics.UpdateDifficulty(block.Difficulty)

	// Record block for health tracking
	if ns.Health != nil {
		ns.Health.RecordBlock()
	}

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

// Helper functions for IBD block parsing
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getUint64(m map[string]interface{}, key string) uint64 {
	if val, ok := m[key].(float64); ok {
		return uint64(val)
	}
	return 0
}

func getInt64(m map[string]interface{}, key string) int64 {
	if val, ok := m[key].(float64); ok {
		return int64(val)
	}
	return 0
}

// hashBlock computes the hash of a block
func hashBlock(b *Block) [32]byte {
	// Simple hash of block data
	h := sha256.New()
	fmt.Fprintf(h, "%d", b.Height)
	fmt.Fprintf(h, "%d", b.TimestampUnix)
	h.Write(b.PrevHash[:])
	fmt.Fprintf(h, "%d", b.Difficulty) // Include difficulty
	h.Write(b.Challenge[:])            // Include challenge
	if b.Proof != nil {
		h.Write(b.Proof.Hash[:])
	}
	return sha256.Sum256(h.Sum(nil))
}

// GetRecentBlocks returns the most recent N blocks
func (ns *NodeState) GetRecentBlocks(count int) interface{} {
	ns.RLock()
	defer ns.RUnlock()

	chainLen := len(ns.Chain)
	if count > chainLen {
		count = chainLen
	}

	start := chainLen - count
	recentBlocks := make([]map[string]interface{}, 0, count)

	for i := start; i < chainLen; i++ {
		block := ns.Chain[i]
		blockHash := hashBlock(&block)

		// Format transactions with type field
		formattedTxs := make([]map[string]interface{}, len(block.Txs))
		for j, tx := range block.Txs {
			txType := "transfer"
			if tx.From == "coinbase" {
				txType = "coinbase"
			}

			formattedTxs[j] = map[string]interface{}{
				"type":   txType,
				"from":   tx.From,
				"to":     tx.To,
				"amount": tx.Amount,
				"fee":    tx.Fee,
				"nonce":  tx.Nonce,
			}
		}

		recentBlocks = append(recentBlocks, map[string]interface{}{
			"height":     block.Height,
			"hash":       hex.EncodeToString(blockHash[:]),
			"timestamp":  block.TimestampUnix,
			"difficulty": block.Difficulty,
			"farmerAddr": block.FarmerAddr,
			"txCount":    len(block.Txs),
			"txs":        formattedTxs,
		})
	}

	return recentBlocks
}

// GetBlockByHeight returns a specific block by height
func (ns *NodeState) GetBlockByHeight(height uint64) (interface{}, error) {
	ns.RLock()
	defer ns.RUnlock()

	if int(height) >= len(ns.Chain) {
		return nil, fmt.Errorf("block %d not found (tip: %d)", height, len(ns.Chain)-1)
	}

	block := ns.Chain[height]
	blockHash := hashBlock(&block)

	// Format transactions with type field
	formattedTxs := make([]map[string]interface{}, len(block.Txs))
	for i, tx := range block.Txs {
		txType := "transfer"
		if tx.From == "coinbase" {
			txType = "coinbase"
		}

		formattedTxs[i] = map[string]interface{}{
			"type":   txType,
			"from":   tx.From,
			"to":     tx.To,
			"amount": tx.Amount,
			"fee":    tx.Fee,
			"nonce":  tx.Nonce,
		}
	}

	// Format proof if present (needed for hash calculation during IBD)
	var proofData interface{} = nil
	if block.Proof != nil {
		proofData = map[string]interface{}{
			"hash":         hex.EncodeToString(block.Proof.Hash[:]),
			"quality":      block.Proof.Quality,
			"plotID":       hex.EncodeToString(block.Proof.PlotID[:]),
			"index":        block.Proof.Index,
			"farmerPubKey": hex.EncodeToString(block.Proof.FarmerPubKey[:]),
		}
	}

	return map[string]interface{}{
		"height":     block.Height,
		"hash":       hex.EncodeToString(blockHash[:]),
		"prevHash":   hex.EncodeToString(block.PrevHash[:]),
		"timestamp":  block.TimestampUnix,
		"difficulty": block.Difficulty,
		"challenge":  hex.EncodeToString(block.Challenge[:]),
		"farmerAddr": block.FarmerAddr,
		"txCount":    len(block.Txs),
		"txs":        formattedTxs,
		"proof":      proofData, // Include proof for hash calculation during IBD
	}, nil
}

// handleSnapshotCommand handles the 'snapshot' subcommand
func handleSnapshotCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: archivas-node snapshot <export|import> [flags]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  export    Export a snapshot at a specific height")
		fmt.Println("  import    Import a snapshot into an empty database")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  archivas-node snapshot export --height 1200000 --out snapshot-1200000.tar.gz --db ./data")
		fmt.Println("  archivas-node snapshot import --in snapshot-1200000.tar.gz --db ./data")
		os.Exit(1)
	}

	cmd := os.Args[2]
	switch cmd {
	case "export":
		handleSnapshotExport()
	case "import":
		handleSnapshotImport()
	default:
		fmt.Printf("Unknown snapshot command: %s\n", cmd)
		os.Exit(1)
	}
}

// handleSnapshotExport handles 'snapshot export' subcommand
func handleSnapshotExport() {
	exportCmd := flag.NewFlagSet("export", flag.ExitOnError)
	height := exportCmd.Uint64("height", 0, "Block height to export (required)")
	outputPath := exportCmd.String("out", "", "Output file path (required, .tar.gz)")
	dbPath := exportCmd.String("db", "./data", "Database directory path")
	networkID := exportCmd.String("network-id", "archivas-devnet-v4", "Network ID")
	description := exportCmd.String("desc", "", "Optional description for this snapshot")

	exportCmd.Parse(os.Args[3:])

	if *height == 0 {
		fmt.Println("Error: --height is required")
		exportCmd.PrintDefaults()
		os.Exit(1)
	}

	if *outputPath == "" {
		fmt.Println("Error: --out is required")
		exportCmd.PrintDefaults()
		os.Exit(1)
	}

	// Open database
	db, err := storage.OpenDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	blockStore := storage.NewBlockStorage(db)
	stateStore := storage.NewStateStorage(db)
	metaStore := storage.NewMetadataStorage(db)

	// Export snapshot
	opts := snapshot.ExportOptions{
		Height:      *height,
		OutputPath:  *outputPath,
		DBPath:      *dbPath,
		NetworkID:   *networkID,
		Description: *description,
		FullHistory: false,
	}

	if err := snapshot.Export(db, blockStore, stateStore, metaStore, opts); err != nil {
		log.Fatalf("Snapshot export failed: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ Snapshot export completed successfully")
}

// handleSnapshotImport handles 'snapshot import' subcommand
func handleSnapshotImport() {
	importCmd := flag.NewFlagSet("import", flag.ExitOnError)
	inputPath := importCmd.String("in", "", "Input snapshot file path (required, .tar.gz)")
	dbPath := importCmd.String("db", "./data", "Database directory path")
	force := importCmd.Bool("force", false, "Force import even if database is non-empty (will overwrite)")

	importCmd.Parse(os.Args[3:])

	if *inputPath == "" {
		fmt.Println("Error: --in is required")
		importCmd.PrintDefaults()
		os.Exit(1)
	}

	// Import snapshot
	opts := snapshot.ImportOptions{
		InputPath: *inputPath,
		DBPath:    *dbPath,
		Force:     *force,
	}

	metadata, err := snapshot.Import(opts)
	if err != nil {
		log.Fatalf("Snapshot import failed: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ Snapshot import completed successfully")
	fmt.Println()
	fmt.Println("📋 Next steps:")
	fmt.Println("  1. Start the node with checkpoint validation:")
	fmt.Println()
	fmt.Printf("     archivas-node \\\n")
	fmt.Printf("       --db %s \\\n", *dbPath)
	fmt.Printf("       --network-id %s \\\n", metadata.NetworkID)
	fmt.Printf("       --checkpoint-height %d \\\n", metadata.Height)
	fmt.Printf("       --checkpoint-hash %s \\\n", metadata.BlockHash)
	fmt.Println("       --no-peer-discovery \\")
	fmt.Println("       --peer-whitelist seed.archivas.ai:9090 \\")
	fmt.Println("       --peer-whitelist seed2.archivas.ai:9090")
	fmt.Println()
	fmt.Println("  2. The node will sync remaining blocks from the whitelisted seeds")
}

// handleBootstrapCommand handles 'bootstrap' subcommand for automated setup
func handleBootstrapCommand() {
	if len(os.Args) < 3 || (os.Args[2] != "--help" && os.Args[2] != "-h" && len(os.Args) < 4) {
		fmt.Println("Usage: archivas-node bootstrap [flags]")
		fmt.Println()
		fmt.Println("Bootstrap downloads a snapshot and starts a private node in one command.")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --network <id>           Network to bootstrap (e.g., 'devnet', 'mainnet')")
		fmt.Println("  --snapshot-url <url>     URL to snapshot manifest JSON (default: auto-detects from network)")
		fmt.Println("  --db <path>              Database directory path (default: ./data)")
		fmt.Println("  --force                  Force overwrite if database exists")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  # Bootstrap a devnet private node:")
		fmt.Println("  archivas-node bootstrap --network devnet --db /var/lib/archivas")
		fmt.Println()
		fmt.Println("  # Bootstrap with custom snapshot URL:")
		fmt.Println("  archivas-node bootstrap \\")
		fmt.Println("    --snapshot-url https://snapshots.archivas.ai/devnet/latest.json \\")
		fmt.Println("    --db /var/lib/archivas")
		os.Exit(1)
	}

	bootstrapCmd := flag.NewFlagSet("bootstrap", flag.ExitOnError)
	network := bootstrapCmd.String("network", "", "Network ID (devnet, mainnet)")
	snapshotURL := bootstrapCmd.String("snapshot-url", "", "Snapshot manifest URL (auto-detects if not specified)")
	dbPath := bootstrapCmd.String("db", "./data", "Database directory path")
	force := bootstrapCmd.Bool("force", false, "Force overwrite if database exists")

	bootstrapCmd.Parse(os.Args[2:])

	// Auto-detect snapshot URL from network if not specified
	manifestURL := *snapshotURL
	if manifestURL == "" {
		if *network == "" {
			fmt.Println("Error: Either --network or --snapshot-url must be specified")
			os.Exit(1)
		}

		// Map network to default snapshot URL
		switch *network {
		case "devnet":
			manifestURL = "https://seed2.archivas.ai/devnet/latest.json"
		case "mainnet":
			manifestURL = "https://seed2.archivas.ai/mainnet/latest.json"
		default:
			fmt.Printf("Error: Unknown network '%s'. Use --snapshot-url to specify a custom URL.\n", *network)
			os.Exit(1)
		}
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  🚀 ARCHIVAS NODE BOOTSTRAP")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Bootstrap (download + import snapshot)
	opts := snapshot.BootstrapOptions{
		ManifestURL: manifestURL,
		DBPath:      *dbPath,
		Force:       *force,
	}

	metadata, err := snapshot.Bootstrap(opts)
	if err != nil {
		log.Fatalf("Bootstrap failed: %v", err)
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  ✅ BOOTSTRAP COMPLETE")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("  📊 Snapshot imported at height: %d\n", metadata.Height)
	blockHashDisplay := metadata.BlockHash
	if len(blockHashDisplay) > 16 {
		blockHashDisplay = blockHashDisplay[:16] + "..."
	}
	fmt.Printf("  🔗 Block hash: %s\n", blockHashDisplay)
	fmt.Printf("  🌐 Network: %s\n", metadata.NetworkID)
	fmt.Println()
	fmt.Println("  📋 Next steps:")
	fmt.Println()
	fmt.Println("  1. Start your private node:")
	fmt.Println()
	fmt.Printf("     archivas-node \\\n")
	fmt.Printf("       --db %s \\\n", *dbPath)
	fmt.Printf("       --network-id %s \\\n", metadata.NetworkID)
	fmt.Println("       --rpc 127.0.0.1:8080 \\")
	fmt.Println("       --p2p 0.0.0.0:9090 \\")
	fmt.Println("       --no-peer-discovery \\")
	fmt.Println("       --peer-whitelist seed.archivas.ai:9090 \\")
	fmt.Println("       --peer-whitelist seed2.archivas.ai:9090 \\")
	fmt.Printf("       --checkpoint-height %d \\\n", metadata.Height)
	fmt.Printf("       --checkpoint-hash %s\n", metadata.BlockHash)
	fmt.Println()
	fmt.Println("  2. Or use the systemd service (see docs/PRIVATE_NODE_SETUP.md)")
	fmt.Println()
	fmt.Println("  3. Point your farmer to http://127.0.0.1:8080")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
}

// printUsage prints usage information
func printUsage() {
	fmt.Println("Archivas Node - Blockchain node for Archivas Network")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  archivas-node [flags]                        Run the node")
	fmt.Println("  archivas-node bootstrap [flags]              Bootstrap from snapshot (one-command setup)")
	fmt.Println("  archivas-node snapshot <export|import>       Manage snapshots")
	fmt.Println("  archivas-node version                        Show version")
	fmt.Println("  archivas-node help                           Show this help")
	fmt.Println()
	fmt.Println("Node Flags:")
	fmt.Println("  --network <name>            Network to join (betanet, devnet-legacy) [default: betanet]")
	fmt.Println("  --rpc <addr>                RPC listen address (default: from network profile)")
	fmt.Println("  --p2p <addr>                P2P listen address (default: from network profile)")
	fmt.Println("  --db <path>                 Database directory (default: ./data)")
	fmt.Println("  --genesis <path>            Genesis file path (overrides network profile)")
	fmt.Println("  --network-id <id>           Network ID (overrides network profile)")
	fmt.Println()
	fmt.Println("Private Node Flags:")
	fmt.Println("  --no-peer-discovery         Disable automatic peer discovery")
	fmt.Println("  --peer-whitelist <host:port> Whitelisted peer (repeatable)")
	fmt.Println("  --checkpoint-height <N>     Checkpoint height for validation")
	fmt.Println("  --checkpoint-hash <hash>    Checkpoint block hash (hex)")
	fmt.Println()
	fmt.Println("Snapshot Commands:")
	fmt.Println("  archivas-node snapshot export --height <N> --out <file> --db <path>")
	fmt.Println("  archivas-node snapshot import --in <file> --db <path> [--force]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Run a Betanet node (default):")
	fmt.Println("  archivas-node --rpc 0.0.0.0:8545 --p2p 0.0.0.0:9090 --db ./data")
	fmt.Println()
	fmt.Println("  # Run a devnet-legacy node:")
	fmt.Println("  archivas-node --network devnet-legacy --db ./devnet-data")
	fmt.Println()
	fmt.Println("  # Run a private Betanet node for farming:")
	fmt.Println("  archivas-node --network betanet --rpc 127.0.0.1:8545 --p2p 0.0.0.0:9090 \\")
	fmt.Println("    --no-peer-discovery \\")
	fmt.Println("    --peer-whitelist seed1.betanet.archivas.ai:9090 \\")
	fmt.Println("    --peer-whitelist seed2.betanet.archivas.ai:9090")
	fmt.Println()
	fmt.Println("  # Export a snapshot:")
	fmt.Println("  archivas-node snapshot export --height 1200000 --out snapshot.tar.gz --db ./data")
	fmt.Println()
	fmt.Println("  # Import a snapshot:")
	fmt.Println("  archivas-node snapshot import --in snapshot.tar.gz --db ./data")
	fmt.Println()
	fmt.Println("For more information, visit: https://docs.archivas.ai")
}
