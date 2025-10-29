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

// ChallengeResponse represents the challenge info for farmers
type ChallengeResponse struct {
	Challenge  [32]byte `json:"challenge"`
	Difficulty uint64   `json:"difficulty"`
	Height     uint64   `json:"height"`
}

// BlockSubmission represents a block submission from a farmer
type BlockSubmission struct {
	Proof        *pospace.Proof `json:"proof"`
	FarmerAddr   string         `json:"farmerAddr"`
	FarmerPubKey string         `json:"farmerPubKey"` // hex encoded
}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSubmitBlock handles POST /submitBlock
func (s *FarmingServer) handleSubmitBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submission BlockSubmission
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

// ChainTipResponse represents the current chain tip
type ChainTipResponse struct {
	Height     uint64   `json:"height"`
	Hash       string   `json:"hash"` // hex-encoded
	Difficulty uint64   `json:"difficulty"`
}

// handleChainTip handles GET /chainTip (for timelord)
func (s *FarmingServer) handleChainTip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	height, difficulty, tipHash := s.nodeState.GetStatus()

	response := ChainTipResponse{
		Height:     height,
		Hash:       hex.EncodeToString(tipHash[:]),
		Difficulty: difficulty,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VDFUpdateRequest represents a VDF update from timelord
type VDFUpdateRequest struct {
	Seed       []byte `json:"seed"`
	Iterations uint64 `json:"iterations"`
	Output     []byte `json:"output"`
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

	// For now, just log the update (node in PoSpace-only mode doesn't enforce VDF)
	// In VDF mode, this would update the node's VDF state
	log.Printf("[vdf] Received VDF update: iter=%d seed=%x output=%x", 
		update.Iterations, update.Seed[:min(8, len(update.Seed))], update.Output[:min(8, len(update.Output))])

	response := SubmitTxResponse{
		Status:  "success",
		Message: "VDF update received",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

