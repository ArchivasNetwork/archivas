package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ArchivasNetwork/archivas/ledger"
	"github.com/ArchivasNetwork/archivas/mempool"
)

// Server provides RPC interface
type Server struct {
	worldState *ledger.WorldState
	mempool    *mempool.Mempool
}

// NewServer creates a new RPC server
func NewServer(ws *ledger.WorldState, mp *mempool.Mempool) *Server {
	return &Server{
		worldState: ws,
		mempool:    mp,
	}
}

// BalanceResponse represents the response for a balance query
type BalanceResponse struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

// Note: SubmitTxRequest is no longer used - we now accept full Transaction objects with signatures

// SubmitTxResponse represents the response for a tx submission
type SubmitTxResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Start starts the HTTP RPC server
func (s *Server) Start(addr string) error {
	http.HandleFunc("/balance/", s.handleBalance)
	http.HandleFunc("/submitTx", s.handleSubmitTx)
	http.HandleFunc("/", s.handleRoot)

	return http.ListenAndServe(addr, nil)
}

// handleBalance handles GET /balance/<addr>
func (s *Server) handleBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract address from path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}
	addr := parts[2]

	balance := s.worldState.GetBalance(addr)
	nonce := s.worldState.GetNonce(addr)

	response := BalanceResponse{
		Address: addr,
		Balance: balance,
		Nonce:   nonce,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleSubmitTx handles POST /submitTx
func (s *Server) handleSubmitTx(w http.ResponseWriter, r *http.Request) {
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
			"status":  "error",
			"message": "Content-Type must be application/json",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	var tx ledger.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate signature
	if err := ledger.VerifyTransactionSignature(tx); err != nil {
		response := SubmitTxResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid signature: %v", err),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Verify sender has account and sufficient funds (simulation, doesn't modify state)
	sender := s.worldState.Accounts[tx.From]
	if sender == nil {
		response := SubmitTxResponse{
			Status:  "error",
			Message: "Sender account does not exist or has zero balance",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check nonce
	if sender.Nonce != tx.Nonce {
		response := SubmitTxResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid nonce: expected %d, got %d", sender.Nonce, tx.Nonce),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check balance
	totalCost := tx.Amount + tx.Fee
	if sender.Balance < totalCost {
		response := SubmitTxResponse{
			Status:  "error",
			Message: fmt.Sprintf("Insufficient funds: have %d, need %d", sender.Balance, totalCost),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// All checks passed - add to mempool
	s.mempool.Add(tx)

	response := SubmitTxResponse{
		Status:  "success",
		Message: "Transaction added to mempool",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Printf("📬 Transaction received: %s -> %s (%.8f RCHV)\n",
		tx.From, tx.To, float64(tx.Amount)/100000000.0)
}

// handleRoot handles GET / (status endpoint)
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","message":"Archivas Devnet RPC Server"}`)
}
