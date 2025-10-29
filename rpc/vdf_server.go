package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/mempool"
	"github.com/iljanemesis/archivas/pospace"
)

// VDFNodeState interface for VDF-enabled node operations
type VDFNodeState interface {
	AcceptBlock(proof *pospace.Proof, farmerAddr string, farmerPubKey []byte,
		vdfSeed []byte, vdfIterations uint64, vdfOutput []byte) error
	GetCurrentChallengeVDF() ([32]byte, uint64, uint64, []byte, uint64, []byte)
	GetChainTip() ([32]byte, uint64)
	UpdateVDF(seed []byte, iterations uint64, output []byte) error
}

// VDFServer extends Server with VDF capabilities
type VDFServer struct {
	worldState *ledger.WorldState
	mempool    *mempool.Mempool
	nodeState  VDFNodeState
}

// NewVDFServer creates a new VDF-enabled RPC server
func NewVDFServer(ws *ledger.WorldState, mp *mempool.Mempool, ns VDFNodeState) *VDFServer {
	return &VDFServer{
		worldState: ws,
		mempool:    mp,
		nodeState:  ns,
	}
}

// Start starts the VDF RPC server
func (s *VDFServer) Start(addr string) error {
	// Original endpoints
	http.HandleFunc("/balance/", s.handleBalance)
	http.HandleFunc("/submitTx", s.handleSubmitTx)
	http.HandleFunc("/", s.handleRoot)

	// Farming endpoints
	http.HandleFunc("/challenge", s.handleGetChallengeVDF)
	http.HandleFunc("/submitBlock", s.handleSubmitBlockVDF)

	// VDF endpoints
	http.HandleFunc("/chainTip", s.handleChainTip)
	http.HandleFunc("/vdf/update", s.handleVDFUpdate)

	return http.ListenAndServe(addr, nil)
}

// ChainTipResponse represents the current chain tip
type ChainTipResponse struct {
	BlockHash [32]byte `json:"blockHash"`
	Height    uint64   `json:"height"`
}

// VDFUpdateRequest represents a VDF update from timelord
type VDFUpdateRequest struct {
	Seed       []byte `json:"seed"`
	Iterations uint64 `json:"iterations"`
	Output     []byte `json:"output"`
}

// ChallengeVDFResponse represents the challenge info with VDF for farmers
type ChallengeVDFResponse struct {
	Challenge     [32]byte `json:"challenge"`
	Difficulty    uint64   `json:"difficulty"`
	Height        uint64   `json:"height"`
	VDFSeed       []byte   `json:"vdfSeed"`
	VDFIterations uint64   `json:"vdfIterations"`
	VDFOutput     []byte   `json:"vdfOutput"`
}

// BlockVDFSubmission represents a block submission with VDF from a farmer
type BlockVDFSubmission struct {
	Proof         *pospace.Proof `json:"proof"`
	FarmerAddr    string         `json:"farmerAddr"`
	FarmerPubKey  string         `json:"farmerPubKey"` // hex encoded
	VDFSeed       []byte         `json:"vdfSeed"`
	VDFIterations uint64         `json:"vdfIterations"`
	VDFOutput     []byte         `json:"vdfOutput"`
}

// handleChainTip handles GET /chainTip
func (s *VDFServer) handleChainTip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	blockHash, height := s.nodeState.GetChainTip()

	response := ChainTipResponse{
		BlockHash: blockHash,
		Height:    height,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleVDFUpdate handles POST /vdf/update
func (s *VDFServer) handleVDFUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var update VDFUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Update VDF state
	if err := s.nodeState.UpdateVDF(update.Seed, update.Iterations, update.Output); err != nil {
		response := SubmitTxResponse{
			Status:  "error",
			Message: fmt.Sprintf("VDF update rejected: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SubmitTxResponse{
		Status:  "success",
		Message: "VDF updated",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetChallengeVDF handles GET /challenge with VDF info
func (s *VDFServer) handleGetChallengeVDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	challenge, difficulty, height, vdfSeed, vdfIter, vdfOutput := s.nodeState.GetCurrentChallengeVDF()

	response := ChallengeVDFResponse{
		Challenge:     challenge,
		Difficulty:    difficulty,
		Height:        height,
		VDFSeed:       vdfSeed,
		VDFIterations: vdfIter,
		VDFOutput:     vdfOutput,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSubmitBlockVDF handles POST /submitBlock with VDF
func (s *VDFServer) handleSubmitBlockVDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submission BlockVDFSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Decode farmer public key
	farmerPubKey, err := hexDecode(submission.FarmerPubKey)
	if err != nil {
		http.Error(w, "Invalid farmer public key", http.StatusBadRequest)
		return
	}

	// Accept the block with VDF
	if err := s.nodeState.AcceptBlock(submission.Proof, submission.FarmerAddr, farmerPubKey,
		submission.VDFSeed, submission.VDFIterations, submission.VDFOutput); err != nil {
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

// handleBalance, handleSubmitTx, handleRoot - reuse from original server
func (s *VDFServer) handleBalance(w http.ResponseWriter, r *http.Request) {
	originalServer := &Server{
		worldState: s.worldState,
		mempool:    s.mempool,
	}
	originalServer.handleBalance(w, r)
}

func (s *VDFServer) handleSubmitTx(w http.ResponseWriter, r *http.Request) {
	originalServer := &Server{
		worldState: s.worldState,
		mempool:    s.mempool,
	}
	originalServer.handleSubmitTx(w, r)
}

func (s *VDFServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"Archivas Devnet RPC Server (PoSpace+Time)"}`)
}

// Helper function
func hexDecode(s string) ([]byte, error) {
	// Reuse from encoding/hex
	result := make([]byte, len(s)/2)
	for i := 0; i < len(result); i++ {
		hi := hexCharToByte(s[i*2])
		lo := hexCharToByte(s[i*2+1])
		if hi == 255 || lo == 255 {
			return nil, fmt.Errorf("invalid hex")
		}
		result[i] = hi<<4 | lo
	}
	return result, nil
}

func hexCharToByte(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 255
}
