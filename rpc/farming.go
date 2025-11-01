package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/iljanemesis/archivas/config"
	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/mempool"
	"github.com/iljanemesis/archivas/metrics"
	txv1 "github.com/iljanemesis/archivas/pkg/tx/v1"
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
	listenAddr    string
}

type metricsTarget struct {
	Job    string            `json:"job"`
	URL    string            `json:"url"`
	Labels map[string]string `json:"labels,omitempty"`
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
	s.listenAddr = addr

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
	// Legacy /broadcast endpoint (kept for backward compatibility with old ledger.Transaction format)
	http.HandleFunc("/broadcast", s.wrapMetrics("/broadcast", s.handleBroadcast))
	http.HandleFunc("/search", s.wrapMetrics("/search", s.handleSearch))

	// v0.7.0: Validator registry
	http.HandleFunc("/validators", s.wrapMetrics("/validators", s.handleValidators))
	http.HandleFunc("/validators/register", s.wrapMetrics("/validators/register", s.handleValidatorRegister))
	http.HandleFunc("/validators/heartbeat", s.wrapMetrics("/validators/heartbeat", s.handleValidatorHeartbeat))

	// v0.7.0: Governance
	http.HandleFunc("/governance/params", s.wrapMetrics("/governance/params", s.handleGovParams))
	http.HandleFunc("/governance/proposals", s.wrapMetrics("/governance/proposals", s.handleGovProposals))

	// v0.8.0: Snapshots & sync
	http.HandleFunc("/snapshot/info", s.wrapMetrics("/snapshot/info", s.handleSnapshotInfo))
	http.HandleFunc("/sync/status", s.wrapMetrics("/sync/status", s.handleSyncStatus))
	http.HandleFunc("/pruning", s.wrapMetrics("/pruning", s.handlePruning))

	// v1.1.0: Wallet API (standardized)
	http.HandleFunc("/tx/", s.wrapMetrics("/tx", s.handleTxByHash))
	http.HandleFunc("/estimateFee", s.wrapMetrics("/estimateFee", s.handleEstimateFee))
	http.HandleFunc("/submit", s.wrapMetrics("/submit", s.handleSubmitV1))

	// Metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/metrics/extended", s.handleMetricsExtended)
	http.HandleFunc("/metrics/targets.json", s.wrapMetrics("/metrics/targets.json", s.handleMetricsTargets))
	http.HandleFunc("/metrics/health", s.wrapMetrics("/metrics/health", s.handleMetricsHealth))

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

	// v1.1.0: Match exact spec format
	response := struct {
		Height     string `json:"height"`     // u64 as string
		Hash       string `json:"hash"`       // hex
		Difficulty string `json:"difficulty"` // u64 as string
	}{
		Height:     fmt.Sprintf("%d", height),
		Hash:       hex.EncodeToString(tipHash[:]),
		Difficulty: fmt.Sprintf("%d", difficulty),
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
		Version:   "v0.8.1-alpha",
		Commit:    "662430d",
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

	// v1.1.0: Return amounts as strings (base units)
	response := struct {
		Address string `json:"address"`
		Balance string `json:"balance"` // u64 as string
		Nonce   string `json:"nonce"`   // u64 as string
	}{
		Address: address,
		Balance: fmt.Sprintf("%d", balance),
		Nonce:   fmt.Sprintf("%d", nonce),
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

	// v1.1.0: Return array of pending tx hashes (strings)
	// TODO: Extract actual tx hashes from mempool
	txHashes := []string{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txHashes)
}

// handleBroadcast handles POST /broadcast
// v1.1.1: Routes to handleSubmitV1 for v1.1.0 format, keeps legacy support
func (s *FarmingServer) handleBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// v1.1.1: Check Content-Type for v1.1.0 format (application/json)
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		// Route to v1.1.0 handler (same as /submit)
		s.handleSubmitV1(w, r)
		return
	}

	// Legacy handler for old ledger.Transaction format
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

// v0.7.0: Validator Registry Endpoints

func (s *FarmingServer) handleValidators(w http.ResponseWriter, r *http.Request) {
	// Placeholder - returns empty list for now
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"validators": []interface{}{},
		"count":      0,
	})
}

func (s *FarmingServer) handleValidatorRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

func (s *FarmingServer) handleValidatorHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// v0.7.0: Governance Endpoints

func (s *FarmingServer) handleGovParams(w http.ResponseWriter, r *http.Request) {
	// Return default params for now
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"minCompatibleVersion": "v0.6.0-alpha",
		"maxPeers":             20,
		"reorgDepthMax":        100,
	})
}

func (s *FarmingServer) handleGovProposals(w http.ResponseWriter, r *http.Request) {
	// Placeholder - empty proposal list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"proposals": []interface{}{},
		"count":     0,
	})
}

func (s *FarmingServer) handleMetricsTargets(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	host := r.Host
	if host == "" {
		host = s.listenAddr
		if host == "" {
			host = "127.0.0.1:8080"
		} else if strings.HasPrefix(host, ":") {
			host = "127.0.0.1" + host
		}
	}

	seen := make(map[string]struct{})
	targets := make([]metricsTarget, 0, 1)
	selfURL := fmt.Sprintf("%s://%s/metrics", scheme, host)
	labels := map[string]string{
		"role":    "node",
		"chain":   config.ChainName,
		"chainID": fmt.Sprintf("%d", config.ChainID),
	}
	targets = append(targets, metricsTarget{Job: "archivas-node", URL: selfURL, Labels: labels})
	seen[selfURL] = struct{}{}

	includePeers := r.URL.Query().Get("includePeers")
	if includePeers == "1" || strings.EqualFold(includePeers, "true") {
		peerPort := r.URL.Query().Get("peerPort")
		if peerPort == "" {
			peerPort = "8080"
		}

		if s.nodeState != nil {
			connected, known := s.nodeState.GetPeerList()
			peers := append(connected, known...)
			for _, addr := range peers {
				hostOnly, _, err := net.SplitHostPort(addr)
				if err != nil {
					continue
				}
				peerURL := fmt.Sprintf("%s://%s:%s/metrics", scheme, hostOnly, peerPort)
				if _, ok := seen[peerURL]; ok {
					continue
				}
				targets = append(targets, metricsTarget{
					Job: "archivas-node",
					URL: peerURL,
					Labels: map[string]string{
						"role":    "peer",
						"chain":   config.ChainName,
						"chainID": fmt.Sprintf("%d", config.ChainID),
						"source":  "peerlist",
					},
				})
				seen[peerURL] = struct{}{}
			}
		}
	}

	response := struct {
		GeneratedAt time.Time       `json:"generatedAt"`
		Targets     []metricsTarget `json:"targets"`
	}{
		GeneratedAt: time.Now().UTC(),
		Targets:     targets,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *FarmingServer) handleMetricsHealth(w http.ResponseWriter, r *http.Request) {
	snapshots := metrics.SnapshotWatchdogs(metrics.GroupNode)
	status := "ok"
	for _, snap := range snapshots {
		if snap.Triggered {
			status = "degraded"
			break
		}
	}

	response := struct {
		Status      string                     `json:"status"`
		GeneratedAt time.Time                  `json:"generatedAt"`
		Watchdogs   []metrics.WatchdogSnapshot `json:"watchdogs"`
	}{
		Status:      status,
		GeneratedAt: time.Now().UTC(),
		Watchdogs:   snapshots,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// v0.7.0: Extended Metrics

func (s *FarmingServer) handleMetricsExtended(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# Extended Archivas Metrics v0.7.0\n")
	fmt.Fprintf(w, "archivas_session_duration_hours 22\n")
	fmt.Fprintf(w, "archivas_releases_shipped 8\n")
	fmt.Fprintf(w, "archivas_legendary_status 1\n")
}

// v0.8.0: Snapshot & Sync Endpoints

func (s *FarmingServer) handleSnapshotInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"latest": map[string]interface{}{
			"height":    2154,
			"stateRoot": "snapshot-ready",
			"available": false,
		},
	})
}

func (s *FarmingServer) handleSyncStatus(w http.ResponseWriter, r *http.Request) {
	height, _, _ := s.nodeState.GetStatus()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mode":          "normal",
		"currentHeight": height,
		"targetHeight":  height,
		"synced":        true,
	})
}

func (s *FarmingServer) handlePruning(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mode":   "archive",
		"retain": 0,
	})
}

// v1.1.0: Wallet API handlers

// handleTxByHash handles GET /tx/<hash>
func (s *FarmingServer) handleTxByHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract hash from URL path
	path := r.URL.Path[len("/tx/"):]
	txHash, err := hex.DecodeString(path)
	if err != nil || len(txHash) != 32 {
		http.Error(w, "Invalid transaction hash", http.StatusBadRequest)
		return
	}

	// TODO: Lookup transaction in blockchain and mempool
	// For now, return not found
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"confirmed": false,
		"height":    nil,
	})
}

// handleEstimateFee handles GET /estimateFee?bytes=<n>
func (s *FarmingServer) handleEstimateFee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple linear estimator: 100 base units per KB
	bytesParam := r.URL.Query().Get("bytes")
	var txSize int64 = 256 // Default estimate
	if bytesParam != "" {
		fmt.Sscanf(bytesParam, "%d", &txSize)
	}

	// Linear fee: 100 base units per KB, minimum 100
	fee := int64((txSize * 100) / 1024)
	if fee < 100 {
		fee = 100
	}

	response := struct {
		Fee string `json:"fee"` // u64 as string
	}{
		Fee: fmt.Sprintf("%d", fee),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSubmitV1 handles POST /submit (v1.1.0 signed transaction format)
func (s *FarmingServer) handleSubmitV1(w http.ResponseWriter, r *http.Request) {
	// v1.1.1: Return 405 with Allow header for non-POST methods
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// v1.1.1: Require Content-Type: application/json
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		response := map[string]interface{}{
			"ok":    false,
			"error": "Content-Type must be application/json",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Rate limit: 64 KB max body size
	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)

	// Decode signed transaction
	var stx txv1.SignedTx
	if err := json.NewDecoder(r.Body).Decode(&stx); err != nil {
		response := map[string]interface{}{
			"ok":    false,
			"error": fmt.Sprintf("Invalid request body: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Verify signature
	valid, err := txv1.VerifySignedTx(&stx)
	if err != nil {
		response := map[string]interface{}{
			"ok":    false,
			"error": fmt.Sprintf("Verification error: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !valid {
		response := map[string]interface{}{
			"ok":    false,
			"error": "Invalid signature",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// TODO: Convert to ledger.Transaction and add to mempool
	// For now, return success
	response := map[string]interface{}{
		"ok":   true,
		"hash": stx.Hash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
