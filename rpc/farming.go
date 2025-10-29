package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/mempool"
	"github.com/iljanemesis/archivas/pospace"
)

// NodeState interface for farming operations
type NodeState interface {
	AcceptBlock(proof *pospace.Proof, farmerAddr string, farmerPubKey []byte) error
	GetCurrentChallenge() ([32]byte, uint64, uint64)
	GetStatus() (height uint64, difficulty uint64, tipHash [32]byte)
	GetCurrentVDF() (seed []byte, iterations uint64, output []byte, hasVDF bool)
}

// FarmingServer extends Server with farming capabilities
type FarmingServer struct {
	worldState *ledger.WorldState
	mempool    *mempool.Mempool
	nodeState  NodeState
}

// NewFarmingServer creates a new farming-enabled RPC server
func NewFarmingServer(ws *ledger.WorldState, mp *mempool.Mempool, ns NodeState) *FarmingServer {
	return &FarmingServer{
		worldState: ws,
		mempool:    mp,
		nodeState:  ns,
	}
}

// Start starts the farming RPC server
func (s *FarmingServer) Start(addr string) error {
	// Original endpoints
	http.HandleFunc("/balance/", s.handleBalance)
	http.HandleFunc("/submitTx", s.handleSubmitTx)
	http.HandleFunc("/", s.handleRoot)

	// Farming endpoints
	http.HandleFunc("/challenge", s.handleGetChallenge)
	http.HandleFunc("/submitBlock", s.handleSubmitBlock)
	
	// VDF/Timelord endpoints
	http.HandleFunc("/chainTip", s.handleChainTip)
	http.HandleFunc("/vdf/update", s.handleVDFUpdate)

	return http.ListenAndServe(addr, nil)
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

