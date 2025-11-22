package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/ArchivasNetwork/archivas/evm"
	"github.com/ArchivasNetwork/archivas/types"
)

// ETHHandler provides Ethereum JSON-RPC compatibility
// Phase 3: Minimal eth_* endpoint implementation
type ETHHandler struct {
	chainID      uint64
	stateDB      evm.StateDB
	getHeight    func() uint64
	getBlock     func(height uint64) (*types.Block, error)
	getBlockByHash func(hash [32]byte) (*types.Block, error)
	getReceipt   func(txHash [32]byte) (*types.Receipt, error)
	submitTx     func(tx *types.EVMTransaction) error
}

// NewETHHandler creates a new ETH RPC handler
func NewETHHandler(
	chainID uint64,
	stateDB evm.StateDB,
	getHeight func() uint64,
	getBlock func(height uint64) (*types.Block, error),
	getBlockByHash func(hash [32]byte) (*types.Block, error),
	getReceipt func(txHash [32]byte) (*types.Receipt, error),
	submitTx func(tx *types.EVMTransaction) error,
) *ETHHandler {
	return &ETHHandler{
		chainID:        chainID,
		stateDB:        stateDB,
		getHeight:      getHeight,
		getBlock:       getBlock,
		getBlockByHash: getBlockByHash,
		getReceipt:     getReceipt,
		submitTx:       submitTx,
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
	case "eth_getBlockByNumber":
		result, err = h.getBlockByNumber_handler(req.Params)
	case "eth_getBlockByHash":
		result, err = h.getBlockByHash_handler(req.Params)
	case "eth_estimateGas":
		result, err = h.estimateGas_handler(req.Params)
	case "eth_feeHistory":
		result, err = h.feeHistory_handler(req.Params)
	case "eth_syncing":
		result, err = h.syncing_handler()
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

// eth_getBlockByNumber returns block information by number
func (h *ETHHandler) getBlockByNumber_handler(params json.RawMessage) (interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing block number parameter")
	}

	blockNumStr, ok := p[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid block number parameter")
	}

	var blockNum uint64
	if blockNumStr == "latest" || blockNumStr == "pending" {
		blockNum = h.getHeight()
	} else if blockNumStr == "earliest" {
		blockNum = 0
	} else {
		// Parse hex block number
		blockNumStr = strings.TrimPrefix(blockNumStr, "0x")
		var err error
		blockNum, err = strconv.ParseUint(blockNumStr, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid block number: %v", err)
		}
	}

	// Get block from node
	block, err := h.getBlock(blockNum)
	if err != nil {
		return nil, nil // Return null for not found
	}

	// Include full transaction objects or just hashes?
	fullTx := false
	if len(p) > 1 {
		if fullTxBool, ok := p[1].(bool); ok {
			fullTx = fullTxBool
		}
	}

	// Format block response (simplified for now)
	response := map[string]interface{}{
		"number":     fmt.Sprintf("0x%x", blockNum),
		"hash":       fmt.Sprintf("0x%x", block.Hash()),
		"parentHash": fmt.Sprintf("0x%x", block.PrevHash),
		"timestamp":  fmt.Sprintf("0x%x", block.TimestampUnix),
		"difficulty": fmt.Sprintf("0x%x", block.Difficulty),
		"gasLimit":   fmt.Sprintf("0x%x", block.GasLimit),
		"gasUsed":    fmt.Sprintf("0x%x", block.GasUsed),
		"miner":      "0x0000000000000000000000000000000000000000", // PoST doesn't have miners
		"transactions": []interface{}{}, // Simplified
	}

	if !fullTx {
		// Return array of transaction hashes
		response["transactions"] = []interface{}{}
	}

	return response, nil
}

// eth_syncing returns syncing status
func (h *ETHHandler) syncing_handler() (interface{}, error) {
	// For now, always return false (not syncing)
	// In a full implementation, this would check if the node is in IBD mode
	return false, nil
}

// eth_getBlockByHash returns block information by hash
func (h *ETHHandler) getBlockByHash_handler(params json.RawMessage) (interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing block hash parameter")
	}

	blockHashStr, ok := p[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid block hash parameter")
	}

	blockHash, err := parseHash(blockHashStr)
	if err != nil {
		return nil, fmt.Errorf("invalid block hash: %v", err)
	}

	// Get block from node by hash
	block, err := h.getBlockByHash(blockHash)
	if err != nil {
		return nil, nil // Return null for not found
	}

	// Include full transaction objects or just hashes?
	fullTx := false
	if len(p) > 1 {
		if fullTxBool, ok := p[1].(bool); ok {
			fullTx = fullTxBool
		}
	}

	// Format block response (similar to getBlockByNumber)
	response := map[string]interface{}{
		"number":     fmt.Sprintf("0x%x", block.Height),
		"hash":       fmt.Sprintf("0x%x", block.Hash()),
		"parentHash": fmt.Sprintf("0x%x", block.PrevHash),
		"timestamp":  fmt.Sprintf("0x%x", block.TimestampUnix),
		"difficulty": fmt.Sprintf("0x%x", block.Difficulty),
		"gasLimit":   fmt.Sprintf("0x%x", block.GasLimit),
		"gasUsed":    fmt.Sprintf("0x%x", block.GasUsed),
		"miner":      block.FarmerAddr.Hex(), // PoST farmer address
		"stateRoot":  fmt.Sprintf("0x%x", block.StateRoot),
		"receiptsRoot": fmt.Sprintf("0x%x", block.ReceiptsRoot),
		"transactions": []interface{}{}, // Simplified
	}

	if !fullTx {
		// Return array of transaction hashes
		response["transactions"] = []interface{}{}
	}

	return response, nil
}

// eth_estimateGas estimates gas for a transaction
func (h *ETHHandler) estimateGas_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 1 {
		return "", fmt.Errorf("missing transaction object parameter")
	}

	callObj, ok := p[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid transaction object parameter")
	}

	// Extract transaction parameters
	// For now, implement a simplified estimation:
	// - Simple transfers: 21000 gas
	// - Contract calls/deployments: estimate based on data size

	// Check if there's contract data
	dataStr, hasData := callObj["data"].(string)
	valueStr, _ := callObj["value"].(string)

	baseGas := uint64(21000) // Base cost for a simple transfer

	if hasData && dataStr != "" && dataStr != "0x" {
		// Contract call or deployment
		data := strings.TrimPrefix(dataStr, "0x")
		dataBytes, err := hex.DecodeString(data)
		if err != nil {
			return "", fmt.Errorf("invalid data hex: %v", err)
		}

		// Estimate gas based on data size
		// Each byte of data costs:
		// - 4 gas for zero bytes
		// - 16 gas for non-zero bytes
		// Plus overhead for contract execution
		gasForData := uint64(0)
		for _, b := range dataBytes {
			if b == 0 {
				gasForData += 4
			} else {
				gasForData += 16
			}
		}

		// Check if this is a contract deployment (no "to" field)
		_, hasTo := callObj["to"]
		if !hasTo {
			// Contract deployment: base + data + overhead
			baseGas = 53000 // Deployment base cost
		} else {
			// Contract call: base + data + execution overhead
			baseGas = 21000
		}

		estimatedGas := baseGas + gasForData + 20000 // Add execution overhead
		return fmt.Sprintf("0x%x", estimatedGas), nil
	}

	// Simple value transfer
	if valueStr != "" && valueStr != "0x" && valueStr != "0x0" {
		return fmt.Sprintf("0x%x", baseGas), nil
	}

	// Default: return 21000 for simple transfer
	return fmt.Sprintf("0x%x", baseGas), nil
}

// eth_feeHistory returns fee history for recent blocks
func (h *ETHHandler) feeHistory_handler(params json.RawMessage) (map[string]interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 2 {
		return nil, fmt.Errorf("missing parameters (blockCount and newestBlock required)")
	}

	// Parse blockCount (can be number or hex string)
	var blockCount uint64
	switch v := p[0].(type) {
	case float64:
		blockCount = uint64(v)
	case string:
		blockCountStr := strings.TrimPrefix(v, "0x")
		parsed, err := strconv.ParseUint(blockCountStr, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid blockCount: %v", err)
		}
		blockCount = parsed
	default:
		return nil, fmt.Errorf("invalid blockCount type")
	}

	// Parse newestBlock (can be number, hex string, or "latest")
	var newestBlock uint64
	switch v := p[1].(type) {
	case float64:
		newestBlock = uint64(v)
	case string:
		if v == "latest" || v == "pending" {
			newestBlock = h.getHeight()
		} else if v == "earliest" {
			newestBlock = 0
		} else {
			newestBlockStr := strings.TrimPrefix(v, "0x")
			parsed, err := strconv.ParseUint(newestBlockStr, 16, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid newestBlock: %v", err)
			}
			newestBlock = parsed
		}
	default:
		return nil, fmt.Errorf("invalid newestBlock type")
	}

	// Parse rewardPercentiles (optional)
	// For now, we'll return empty rewards since we don't have historical gas price data

	// Simplified implementation:
	// Archivas uses fixed gas price (no EIP-1559 base fee yet)
	// Return constant gas price for all blocks in range

	oldestBlock := newestBlock
	if blockCount > newestBlock {
		oldestBlock = 0
	} else {
		oldestBlock = newestBlock - blockCount + 1
	}

	// Build baseFeePerGas array (one entry per block + 1 for next block)
	baseFeePerGas := make([]string, blockCount+1)
	gasUsedRatio := make([]float64, blockCount)

	// Fixed gas price: 1 gwei
	fixedGasPrice := "0x3b9aca00" // 1000000000 wei = 1 gwei

	for i := range baseFeePerGas {
		baseFeePerGas[i] = fixedGasPrice
	}

	// Calculate gas used ratio for each block
	for i := uint64(0); i < blockCount; i++ {
		blockNum := oldestBlock + i
		block, err := h.getBlock(blockNum)
		if err != nil {
			// If block not found, use 0 ratio
			gasUsedRatio[i] = 0.0
			continue
		}

		if block.GasLimit == 0 {
			gasUsedRatio[i] = 0.0
		} else {
			gasUsedRatio[i] = float64(block.GasUsed) / float64(block.GasLimit)
		}
	}

	response := map[string]interface{}{
		"oldestBlock":   fmt.Sprintf("0x%x", oldestBlock),
		"baseFeePerGas": baseFeePerGas,
		"gasUsedRatio":  gasUsedRatio,
		"reward":        [][]string{}, // Empty for now (no priority fees)
	}

	return response, nil
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

