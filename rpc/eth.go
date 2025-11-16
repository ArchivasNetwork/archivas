package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/evm"
	"github.com/ArchivasNetwork/archivas/types"
)

// ETHHandler provides Ethereum JSON-RPC compatibility
// Phase 3: Minimal eth_* endpoint implementation
type ETHHandler struct {
	chainID   uint64
	stateDB   evm.StateDB
	getHeight func() uint64
	getBlock  func(height uint64) (*types.Block, error)
	getReceipt func(txHash [32]byte) (*types.Receipt, error)
	submitTx  func(tx *types.EVMTransaction) error
}

// NewETHHandler creates a new ETH RPC handler
func NewETHHandler(
	chainID uint64,
	stateDB evm.StateDB,
	getHeight func() uint64,
	getBlock func(height uint64) (*types.Block, error),
	getReceipt func(txHash [32]byte) (*types.Receipt, error),
	submitTx func(tx *types.EVMTransaction) error,
) *ETHHandler {
	return &ETHHandler{
		chainID:    chainID,
		stateDB:    stateDB,
		getHeight:  getHeight,
		getBlock:   getBlock,
		getReceipt: getReceipt,
		submitTx:   submitTx,
	}
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ServeHTTP handles ETH RPC requests
func (h *ETHHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nil, -32700, "Parse error")
		return
	}

	// Route to appropriate handler
	var result interface{}
	var err error

	switch req.Method {
	case "eth_chainId":
		result, err = h.chainID_handler()
	case "eth_blockNumber":
		result, err = h.blockNumber_handler()
	case "eth_getBalance":
		result, err = h.getBalance_handler(req.Params)
	case "eth_getCode":
		result, err = h.getCode_handler(req.Params)
	case "eth_getTransactionReceipt":
		result, err = h.getTransactionReceipt_handler(req.Params)
	case "eth_call":
		result, err = h.call_handler(req.Params)
	case "eth_sendRawTransaction":
		result, err = h.sendRawTransaction_handler(req.Params)
	case "eth_getTransactionCount":
		result, err = h.getTransactionCount_handler(req.Params)
	case "eth_gasPrice":
		result, err = h.gasPrice_handler()
	case "net_version":
		result, err = h.netVersion_handler()
	default:
		writeError(w, req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
		return
	}

	if err != nil {
		writeError(w, req.ID, -32603, err.Error())
		return
	}

	writeSuccess(w, req.ID, result)
}

// eth_chainId returns the chain ID
func (h *ETHHandler) chainID_handler() (string, error) {
	return fmt.Sprintf("0x%x", h.chainID), nil
}

// eth_blockNumber returns the latest block number
func (h *ETHHandler) blockNumber_handler() (string, error) {
	height := h.getHeight()
	return fmt.Sprintf("0x%x", height), nil
}

// eth_getBalance returns the balance of an account
func (h *ETHHandler) getBalance_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	addrStr, ok := p[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	addr, err := parseAddress(addrStr)
	if err != nil {
		return "", err
	}

	balance := h.stateDB.GetBalance(addr)
	return fmt.Sprintf("0x%x", balance), nil
}

// eth_getCode returns the code at an address
func (h *ETHHandler) getCode_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	addrStr, ok := p[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	addr, err := parseAddress(addrStr)
	if err != nil {
		return "", err
	}

	code := h.stateDB.GetCode(addr)
	if len(code) == 0 {
		return "0x", nil
	}

	return "0x" + hex.EncodeToString(code), nil
}

// eth_getTransactionReceipt returns the receipt of a transaction
func (h *ETHHandler) getTransactionReceipt_handler(params json.RawMessage) (map[string]interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing transaction hash parameter")
	}

	txHashStr, ok := p[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid transaction hash parameter")
	}

	txHash, err := parseHash(txHashStr)
	if err != nil {
		return nil, err
	}

	receipt, err := h.getReceipt(txHash)
	if err != nil {
		return nil, nil // Return null for not found
	}

	return formatReceipt(receipt), nil
}

// eth_call executes a call without creating a transaction
func (h *ETHHandler) call_handler(params json.RawMessage) (string, error) {
	// Simplified implementation - just returns 0x
	// Full implementation would execute EVM call
	return "0x", nil
}

// eth_sendRawTransaction submits a signed transaction
func (h *ETHHandler) sendRawTransaction_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 1 {
		return "", fmt.Errorf("missing transaction data parameter")
	}

	txDataStr, ok := p[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid transaction data parameter")
	}

	// Parse raw transaction
	tx, err := parseRawTransaction(txDataStr)
	if err != nil {
		return "", err
	}

	// Submit transaction
	if err := h.submitTx(tx); err != nil {
		return "", err
	}

	// Return transaction hash
	txHash := tx.Hash()
	return "0x" + hex.EncodeToString(txHash[:]), nil
}

// eth_getTransactionCount returns the nonce of an account
func (h *ETHHandler) getTransactionCount_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	addrStr, ok := p[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	addr, err := parseAddress(addrStr)
	if err != nil {
		return "", err
	}

	nonce := h.stateDB.GetNonce(addr)
	return fmt.Sprintf("0x%x", nonce), nil
}

// eth_gasPrice returns the current gas price
func (h *ETHHandler) gasPrice_handler() (string, error) {
	// Return 1 gwei as default
	return "0x3b9aca00", nil // 1000000000 wei = 1 gwei
}

// net_version returns the network ID
func (h *ETHHandler) netVersion_handler() (string, error) {
	return fmt.Sprintf("%d", h.chainID), nil
}

// Helper functions

func parseAddress(s string) (address.EVMAddress, error) {
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 40 {
		return address.ZeroAddress(), fmt.Errorf("invalid address length")
	}
	return address.EVMAddressFromHex("0x" + s)
}

func parseHash(s string) ([32]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, err
	}
	if len(bytes) != 32 {
		return [32]byte{}, fmt.Errorf("invalid hash length")
	}
	var hash [32]byte
	copy(hash[:], bytes)
	return hash, nil
}

func parseRawTransaction(s string) (*types.EVMTransaction, error) {
	// Simplified - would need full RLP decoding in production
	return nil, fmt.Errorf("transaction parsing not implemented yet")
}

func formatReceipt(receipt *types.Receipt) map[string]interface{} {
	status := "0x0"
	if receipt.Status == 1 {
		status = "0x1"
	}

	result := map[string]interface{}{
		"transactionHash":  "0x" + hex.EncodeToString(receipt.TxHash[:]),
		"blockNumber":      fmt.Sprintf("0x%x", receipt.BlockHeight),
		"transactionIndex": fmt.Sprintf("0x%x", receipt.TxIndex),
		"from":             receipt.From.Hex(),
		"gasUsed":          fmt.Sprintf("0x%x", receipt.GasUsed),
		"cumulativeGasUsed": fmt.Sprintf("0x%x", receipt.CumulativeGasUsed),
		"status":           status,
		"logs":             []interface{}{}, // Simplified
	}

	if receipt.To != nil {
		result["to"] = receipt.To.Hex()
	} else {
		result["to"] = nil
	}

	if receipt.ContractAddress != nil {
		result["contractAddress"] = receipt.ContractAddress.Hex()
	} else {
		result["contractAddress"] = nil
	}

	return result
}

func writeSuccess(w http.ResponseWriter, id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
		ID: id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // JSON-RPC errors still return 200
	json.NewEncoder(w).Encode(resp)
}

