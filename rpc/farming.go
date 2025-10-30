package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/mempool"
	"github.com/iljanemesis/archivas/metrics"
	"github.com/iljanemesis/archivas/pospace"
	"github.com/iljanemesis/archivas/wallet"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NodeState interface for farming operations
type NodeState interface {
	AcceptBlock(proof *pospace.Proof, farmerAddr string, farmerPubKey []byte) error
	GetCurrentChallenge() ([32]byte, uint64, uint64)
	GetStatus() (height uint64, difficulty uint64, tipHash [32]byte)
	GetCurrentVDF() (seed []byte, iterations uint64, output []byte, hasVDF bool)
	GetGenesisHash() [32]byte
	GetPeerCount() int
	GetPeerList() (connected []string, known []string)
	GetHealthStats() interface{}
	GetRecentBlocks(count int) interface{}
	GetBlockByHeight(height uint64) (interface{}, error)
}

// FarmingServer extends Server with farming capabilities
type FarmingServer struct {
	worldState    *ledger.WorldState
	mempool       *mempool.Mempool
	nodeState     NodeState
	faucetEnabled bool
	faucetKey     []byte
	faucetAddress string
	faucetLimit   map[string]time.Time // IP -> last drip time
	faucetMutex   sync.Mutex
}

// NewFarmingServer creates a new farming-enabled RPC server
func NewFarmingServer(ws *ledger.WorldState, mp *mempool.Mempool, ns NodeState) *FarmingServer {
	return &FarmingServer{
		worldState:    ws,
		mempool:       mp,
		nodeState:     ns,
		faucetEnabled: false,
		faucetLimit:   make(map[string]time.Time),
	}
}

// EnableFaucet enables the built-in faucet with a funding private key
func (s *FarmingServer) EnableFaucet(privKeyHex string) error {
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil || len(privKeyBytes) != 32 {
		return fmt.Errorf("invalid private key: must be 32 bytes hex")
	}

	priv := secp256k1.PrivKeyFromBytes(privKeyBytes)
	pubKey := priv.PubKey().SerializeCompressed()
	address, err := wallet.PubKeyToAddress(pubKey)
	if err != nil {
		return fmt.Errorf("failed to derive address: %w", err)
	}

	s.faucetKey = privKeyBytes
	s.faucetAddress = address
	s.faucetEnabled = true

	log.Printf("[faucet] enabled with address %s", address)
	return nil
}

// Start starts the farming RPC server
func (s *FarmingServer) Start(addr string) error {
	// Original endpoints
	http.HandleFunc("/balance/", s.wrapMetrics("/balance", s.handleBalance))
	http.HandleFunc("/submitTx", s.wrapMetrics("/submitTx", s.handleSubmitTx))
	http.HandleFunc("/", s.wrapMetrics("/", s.handleRoot))

	// Farming endpoints
	http.HandleFunc("/challenge", s.wrapMetrics("/challenge", s.handleGetChallenge))
	http.HandleFunc("/submitBlock", s.wrapMetrics("/submitBlock", s.handleSubmitBlock))

	// VDF/Timelord endpoints
	http.HandleFunc("/chainTip", s.wrapMetrics("/chainTip", s.handleChainTip))
	http.HandleFunc("/vdf/update", s.wrapMetrics("/vdf/update", s.handleVDFUpdate))

	// Network endpoints
	http.HandleFunc("/genesisHash", s.wrapMetrics("/genesisHash", s.handleGenesisHash))
	http.HandleFunc("/healthz", s.wrapMetrics("/healthz", s.handleHealthz))
	http.HandleFunc("/peers", s.wrapMetrics("/peers", s.handlePeers))
	http.HandleFunc("/health", s.wrapMetrics("/health", s.handleHealthDetailed))

	// Faucet endpoint (if enabled)
	if s.faucetEnabled {
		http.HandleFunc("/faucet", s.wrapMetrics("/faucet", s.handleFaucet))
	}
	
	// Developer endpoints
	http.HandleFunc("/recentBlocks", s.wrapMetrics("/recentBlocks", s.handleRecentBlocks))
	http.HandleFunc("/block/", s.wrapMetrics("/block", s.handleBlockByHeight))
	http.HandleFunc("/version", s.wrapMetrics("/version", s.handleVersion))
	http.HandleFunc("/account/", s.wrapMetrics("/account", s.handleAccount))
	http.HandleFunc("/mempool", s.wrapMetrics("/mempool", s.handleMempoolView))
	http.HandleFunc("/broadcast", s.wrapMetrics("/broadcast", s.handleBroadcast))
	http.HandleFunc("/search", s.wrapMetrics("/search", s.handleSearch))
	
	// Metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, s.corsMiddleware(http.DefaultServeMux))
}

// corsMiddleware adds CORS headers for external developers
func (s *FarmingServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// wrapMetrics wraps an HTTP handler to increment request metrics
func (s *FarmingServer) wrapMetrics(endpoint string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.RPCRequests.WithLabelValues(endpoint).Inc()
		handler(w, r)
	}
}

// Types moved to rpc/types.go to avoid duplication

// handleGetChallenge handles GET /challenge
func (s *FarmingServer) handleGetChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	challenge, difficulty, height := s.nodeState.GetCurrentChallenge()

	response := ChallengeResponse{
		Challenge:  challenge,
		Difficulty: difficulty,
		Height:     height,
	}

	// Include VDF info if available (for PoSpace+Time farming)
	if seed, iterations, output, hasVDF := s.nodeState.GetCurrentVDF(); hasVDF {
		response.VDF = &VDFInfo{
			Seed:       hex.EncodeToString(seed),
			Iterations: iterations,
			Output:     hex.EncodeToString(output),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSubmitBlock handles POST /submitBlock
func (s *FarmingServer) handleSubmitBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode directly to pospace.Proof to avoid interface{} issues
	var submission struct {
		Proof        *pospace.Proof `json:"proof"`
		FarmerAddr   string         `json:"farmerAddr"`
		FarmerPubKey string         `json:"farmerPubKey"`
	}

	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Decode farmer public key
	farmerPubKey, err := hex.DecodeString(submission.FarmerPubKey)
	if err != nil {
		http.Error(w, "Invalid farmer public key", http.StatusBadRequest)
		return
	}

	// Accept the block
	if err := s.nodeState.AcceptBlock(submission.Proof, submission.FarmerAddr, farmerPubKey); err != nil {
		response := SubmitTxResponse{
			Status:  "error",
			Message: fmt.Sprintf("Block rejected: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SubmitTxResponse{
		Status:  "success",
		Message: "Block accepted",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleBalance handles GET /balance/<addr> (reused from original server)
func (s *FarmingServer) handleBalance(w http.ResponseWriter, r *http.Request) {
	// Same implementation as original Server
	originalServer := &Server{
		worldState: s.worldState,
		mempool:    s.mempool,
	}
	originalServer.handleBalance(w, r)
}

// handleSubmitTx handles POST /submitTx (reused from original server)
func (s *FarmingServer) handleSubmitTx(w http.ResponseWriter, r *http.Request) {
	// Same implementation as original Server
	originalServer := &Server{
		worldState: s.worldState,
		mempool:    s.mempool,
	}
	originalServer.handleSubmitTx(w, r)
}

// handleRoot handles GET / (status endpoint)
func (s *FarmingServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"Archivas Devnet RPC Server (Farming Enabled)"}`)
}

// handleChainTip handles GET /chainTip (for timelord)
func (s *FarmingServer) handleChainTip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	height, difficulty, tipHash := s.nodeState.GetStatus()

	response := ChainTipResponse{
		BlockHash:  tipHash,
		Height:     height,
		Difficulty: difficulty,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleVDFUpdate handles POST /vdf/update (for timelord)
func (s *FarmingServer) handleVDFUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var update VDFUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Store VDF update in node state (for /challenge to include)
	// In full VDF mode, this would also verify and enforce VDF proofs
	if updater, ok := s.nodeState.(interface {
		UpdateVDFState(seed []byte, iterations uint64, output []byte)
	}); ok {
		updater.UpdateVDFState(update.Seed, update.Iterations, update.Output)
	}

	seedPreview := update.Seed
	if len(seedPreview) > 8 {
		seedPreview = seedPreview[:8]
	}
	outputPreview := update.Output
	if len(outputPreview) > 8 {
		outputPreview = outputPreview[:8]
	}

	log.Printf("[vdf] Received VDF update: iter=%d seed=%x output=%x",
		update.Iterations, seedPreview, outputPreview)

	response := SubmitTxResponse{
		Status:  "success",
		Message: "VDF update received",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGenesisHash handles GET /genesisHash
func (s *FarmingServer) handleGenesisHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	genesisHash := s.nodeState.GetGenesisHash()

	response := struct {
		GenesisHash string `json:"genesisHash"`
	}{
		GenesisHash: hex.EncodeToString(genesisHash[:]),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleHealthz handles GET /healthz
func (s *FarmingServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	height, difficulty, _ := s.nodeState.GetStatus()
	peerCount := s.nodeState.GetPeerCount()

	response := struct {
		OK         bool   `json:"ok"`
		Height     uint64 `json:"height"`
		Difficulty uint64 `json:"difficulty"`
		Peers      int    `json:"peers"`
	}{
		OK:         true,
		Height:     height,
		Difficulty: difficulty,
		Peers:      peerCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handlePeers handles GET /peers
func (s *FarmingServer) handlePeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	connected, known := s.nodeState.GetPeerList()

	response := struct {
		Connected []string `json:"connected"`
		Known     []string `json:"known"`
	}{
		Connected: connected,
		Known:     known,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleHealthDetailed handles GET /health (detailed chain health)
func (s *FarmingServer) handleHealthDetailed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	height, difficulty, _ := s.nodeState.GetStatus()
	peerCount := s.nodeState.GetPeerCount()
	healthStats := s.nodeState.GetHealthStats()

	response := struct {
		OK          bool        `json:"ok"`
		Height      uint64      `json:"height"`
		Difficulty  uint64      `json:"difficulty"`
		Peers       int         `json:"peers"`
		HealthStats interface{} `json:"healthStats"`
	}{
		OK:          true,
		Height:      height,
		Difficulty:  difficulty,
		Peers:       peerCount,
		HealthStats: healthStats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleFaucet handles GET /faucet?address=<addr>
func (s *FarmingServer) handleFaucet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !s.faucetEnabled {
		http.Error(w, "Faucet not enabled", http.StatusServiceUnavailable)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	// Validate address format
	if len(address) < 10 || address[:4] != "arcv" {
		http.Error(w, "Invalid Archivas address format", http.StatusBadRequest)
		return
	}

	// Rate limit by IP (1 drip per hour)
	clientIP := r.RemoteAddr
	s.faucetMutex.Lock()
	lastDrip, exists := s.faucetLimit[clientIP]
	if exists && time.Since(lastDrip) < time.Hour {
		s.faucetMutex.Unlock()
		remaining := time.Hour - time.Since(lastDrip)
		http.Error(w, fmt.Sprintf("Rate limited. Try again in %v", remaining.Round(time.Minute)), http.StatusTooManyRequests)
		return
	}
	s.faucetMutex.Unlock()

	// Get current nonce
	balance := s.worldState.GetBalance(s.faucetAddress)
	nonce := s.worldState.GetNonce(s.faucetAddress)

	// Check faucet has funds
	dripAmount := int64(2000000000) // 20 RCHV
	fee := int64(100000)            // 0.001 RCHV
	if balance < dripAmount+fee {
		http.Error(w, fmt.Sprintf("Faucet empty (balance: %.8f RCHV)", float64(balance)/100000000.0), http.StatusServiceUnavailable)
		return
	}

	// Create and sign transaction
	priv := secp256k1.PrivKeyFromBytes(s.faucetKey)
	pubKey := priv.PubKey().SerializeCompressed()

	tx := ledger.Transaction{
		From:         s.faucetAddress,
		To:           address,
		Amount:       dripAmount,
		Fee:          fee,
		Nonce:        nonce,
		SenderPubKey: pubKey,
	}

	if err := wallet.SignTransaction(&tx, s.faucetKey); err != nil {
		http.Error(w, fmt.Sprintf("Failed to sign: %v", err), http.StatusInternalServerError)
		return
	}

	// Add to mempool
	s.mempool.Add(tx)

	// Update rate limit
	s.faucetMutex.Lock()
	s.faucetLimit[clientIP] = time.Now()
	s.faucetMutex.Unlock()

	log.Printf("[faucet] 💧 Dripped 20 RCHV to %s (IP: %s)", address, clientIP)

	// Success response
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Amount  string `json:"amount"`
		TxHash  string `json:"txHash,omitempty"`
	}{
		Success: true,
		Message: "20 RCHV sent! Should arrive in next block (~20 seconds)",
		Amount:  "20.00000000 RCHV",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleRecentBlocks handles GET /recentBlocks?count=10
func (s *FarmingServer) handleRecentBlocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := 10
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		fmt.Sscanf(countStr, "%d", &count)
	}
	
	if count > 100 {
		count = 100 // Max 100
	}

	blocks := s.nodeState.GetRecentBlocks(count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"blocks": blocks,
		"count":  count,
	})
}

// handleBlockByHeight handles GET /block/<height>
func (s *FarmingServer) handleBlockByHeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse height from URL
	heightStr := r.URL.Path[len("/block/"):]
	var height uint64
	if _, err := fmt.Sscanf(heightStr, "%d", &height); err != nil {
		http.Error(w, "Invalid height", http.StatusBadRequest)
		return
	}

	block, err := s.nodeState.GetBlockByHeight(height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

// handleVersion handles GET /version
func (s *FarmingServer) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		Network   string `json:"network"`
		Consensus string `json:"consensus"`
	}{
		Version:   "v0.5.0-alpha",
		Commit:    "dev", // Would be replaced with git commit hash in production
		Network:   "archivas-devnet-v3",
		Consensus: "Proof-of-Space-and-Time",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleAccount handles GET /account/<addr> and /account/<addr>/txs
func (s *FarmingServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse address from URL
	path := r.URL.Path[len("/account/"):]
	parts := strings.Split(path, "/")
	address := parts[0]

	// Check if requesting tx history
	if len(parts) > 1 && parts[1] == "txs" {
		s.handleAccountTxs(w, r, address)
		return
	}

	// Get account info
	balance := s.worldState.GetBalance(address)
	nonce := s.worldState.GetNonce(address)

	response := struct {
		Address string `json:"address"`
		Balance int64  `json:"balance"`
		Nonce   uint64 `json:"nonce"`
	}{
		Address: address,
		Balance: balance,
		Nonce:   nonce,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleAccountTxs handles GET /account/<addr>/txs
func (s *FarmingServer) handleAccountTxs(w http.ResponseWriter, r *http.Request, address string) {
	// TODO: Implement tx history indexing
	// For now, return empty list
	response := struct {
		Address string        `json:"address"`
		Txs     []interface{} `json:"txs"`
		Count   int           `json:"count"`
	}{
		Address: address,
		Txs:     []interface{}{},
		Count:   0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMempoolView handles GET /mempool
func (s *FarmingServer) handleMempoolView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	txs := s.mempool.GetAll()

	response := struct {
		Count int                    `json:"count"`
		Txs   []ledger.Transaction   `json:"txs"`
	}{
		Count: len(txs),
		Txs:   txs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleBroadcast handles POST /broadcast
func (s *FarmingServer) handleBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx ledger.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid transaction", http.StatusBadRequest)
		return
	}

	// Add to mempool
	s.mempool.Add(tx)

	response := struct {
		Status string `json:"status"`
		TxHash string `json:"txHash,omitempty"`
	}{
		Status: "accepted",
		TxHash: "pending", // Would compute actual tx hash
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSearch handles GET /search?q=<query>
func (s *FarmingServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	// Detect type
	var resultType string
	var result interface{}

	if strings.HasPrefix(query, "arcv") {
		// Address
		resultType = "address"
		result = map[string]string{"address": query, "redirect": "/account/" + query}
	} else if len(query) == 64 {
		// Likely a hash (block or tx)
		resultType = "hash"
		result = map[string]string{"hash": query, "type": "block_or_tx"}
	} else {
		// Try as height
		resultType = "height"
		result = map[string]string{"height": query, "redirect": "/block/" + query}
	}

	response := struct {
		Query  string      `json:"query"`
		Type   string      `json:"type"`
		Result interface{} `json:"result"`
	}{
		Query:  query,
		Type:   resultType,
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
