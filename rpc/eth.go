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
	getBlock     func(height uint64) (interface{}, error) // Changed to interface{} to handle both legacy and types.Block
	getBlockByHash func(hash [32]byte) (interface{}, error) // Changed to interface{}
	getReceipt   func(txHash [32]byte) (*types.Receipt, error)
	submitTx     func(tx *types.EVMTransaction) error
	getPeerCount func() int // For net_peerCount
}

// NewETHHandler creates a new ETH RPC handler
func NewETHHandler(
	chainID uint64,
	stateDB evm.StateDB,
	getHeight func() uint64,
	getBlock func(height uint64) (interface{}, error), // Changed to interface{}
	getBlockByHash func(hash [32]byte) (interface{}, error), // Changed to interface{}
	getReceipt func(txHash [32]byte) (*types.Receipt, error),
	submitTx func(tx *types.EVMTransaction) error,
	getPeerCount func() int,
) *ETHHandler {
	return &ETHHandler{
		chainID:        chainID,
		stateDB:        stateDB,
		getHeight:      getHeight,
		getBlock:       getBlock,
		getBlockByHash: getBlockByHash,
		getReceipt:     getReceipt,
		submitTx:       submitTx,
		getPeerCount:   getPeerCount,
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
	// stateDB.GetBalance already returns balance in Wei (18 decimals)
	// via WorldStateAdapter which converts from Archivas 8-decimal format
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

	// Get block from node (returns interface{} - either *types.Block or map[string]interface{})
	blockRaw, err := h.getBlock(blockNum)
	if err != nil {
		return nil, nil // Return null for not found
	}

	// Check if block is nil
	if blockRaw == nil {
		return nil, nil
	}

	// Include full transaction objects or just hashes?
	fullTx := false
	if len(p) > 1 {
		if fullTxBool, ok := p[1].(bool); ok {
			fullTx = fullTxBool
		}
	}

	// Try to convert to legacy block format first (map)
	if legacyBlock, ok := blockRaw.(map[string]interface{}); ok {
		return convertLegacyBlockToEthBlock(legacyBlock, fullTx), nil
	}

	// Try to convert to types.Block
	if typedBlock, ok := blockRaw.(*types.Block); ok {
		return convertTypesBlockToEthBlock(typedBlock, fullTx), nil
	}

	return nil, fmt.Errorf("unexpected block type")
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

	// Get block from node by hash (returns interface{} - either *types.Block or map[string]interface{})
	blockRaw, err := h.getBlockByHash(blockHash)
	if err != nil {
		return nil, nil // Return null for not found
	}

	// Check if block is nil
	if blockRaw == nil {
		return nil, nil
	}

	// Include full transaction objects or just hashes?
	fullTx := false
	if len(p) > 1 {
		if fullTxBool, ok := p[1].(bool); ok {
			fullTx = fullTxBool
		}
	}

	// Try to convert to legacy block format first (map)
	if legacyBlock, ok := blockRaw.(map[string]interface{}); ok {
		return convertLegacyBlockToEthBlock(legacyBlock, fullTx), nil
	}

	// Try to convert to types.Block
	if typedBlock, ok := blockRaw.(*types.Block); ok {
		return convertTypesBlockToEthBlock(typedBlock, fullTx), nil
	}

	return nil, fmt.Errorf("unexpected block type")
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
		blockRaw, err := h.getBlock(blockNum)
		if err != nil || blockRaw == nil {
			// If block not found, use 0 ratio
			gasUsedRatio[i] = 0.0
			continue
		}

		// Extract gas used and gas limit depending on block type
		var gasUsed, gasLimit uint64
		
		if _, ok := blockRaw.(map[string]interface{}); ok {
			// Legacy block format - no EVM gas fields, use defaults
			gasUsed = 0
			gasLimit = 30000000 // 30M gas (standard)
		} else if typedBlock, ok := blockRaw.(*types.Block); ok {
			// types.Block format - has EVM gas fields
			gasUsed = typedBlock.GasUsed
			gasLimit = typedBlock.GasLimit
		} else {
			// Unknown format
			gasUsedRatio[i] = 0.0
			continue
		}

		if gasLimit == 0 {
			gasUsedRatio[i] = 0.0
		} else {
			gasUsedRatio[i] = float64(gasUsed) / float64(gasLimit)
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

// convertLegacyBlockToEthBlock converts Archivas legacy block format (map) to Ethereum format
func convertLegacyBlockToEthBlock(legacyBlock map[string]interface{}, fullTx bool) map[string]interface{} {
	// Extract fields from legacy block - handle both direct types and JSON float64
	var heightVal uint64
	switch v := legacyBlock["height"].(type) {
	case uint64:
		heightVal = v
	case float64:
		heightVal = uint64(v)
	case int:
		heightVal = uint64(v)
	}
	
	hashVal, _ := legacyBlock["hash"].(string)
	prevHashVal, _ := legacyBlock["prevHash"].(string)
	
	var timestampVal int64
	switch v := legacyBlock["timestamp"].(type) {
	case int64:
		timestampVal = v
	case float64:
		timestampVal = int64(v)
	case int:
		timestampVal = int64(v)
	}
	
	var difficultyVal uint64
	switch v := legacyBlock["difficulty"].(type) {
	case uint64:
		difficultyVal = v
	case float64:
		difficultyVal = uint64(v)
	case int:
		difficultyVal = uint64(v)
	}
	
	farmerAddr, _ := legacyBlock["farmerAddr"].(string)
	
	// Convert farmer address to valid EVM miner address
	// This handles ARCV Bech32 addresses, empty strings, and validates format
	minerAddr := ensureValidMinerAddress(farmerAddr)

	// For Ethereum compatibility
	ethBlock := map[string]interface{}{
		"number":          fmt.Sprintf("0x%x", heightVal),
		"hash":            ensureHexPrefix(hashVal),
		"parentHash":      ensureHexPrefix(prevHashVal),
		"nonce":           "0x0000000000000000", // PoST doesn't use nonce
		"sha3Uncles":      "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"logsBloom":       "0x" + strings.Repeat("0", 512), // Empty bloom filter
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"stateRoot":       "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"receiptsRoot":    "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"miner":           minerAddr, // Farmer acts as miner (validated EVM address)
		"difficulty":      fmt.Sprintf("0x%x", difficultyVal),
		"totalDifficulty": fmt.Sprintf("0x%x", difficultyVal), // Simplified
		"extraData":       "0x",
		"size":            "0x400", // Placeholder size
		"gasLimit":        "0x1c9c380", // 30M gas (standard)
		"gasUsed":         "0x0", // No EVM txs in legacy blocks
		"timestamp":       fmt.Sprintf("0x%x", timestampVal),
		"transactions":    []interface{}{}, // Legacy blocks don't have EVM txs
		"uncles":          []interface{}{},
	}

	return ethBlock
}

// convertTypesBlockToEthBlock converts types.Block to Ethereum format
func convertTypesBlockToEthBlock(block *types.Block, fullTx bool) map[string]interface{} {
	// Build transaction list
	transactions := make([]interface{}, 0)
	if !fullTx {
		// Just hashes
		for _, tx := range block.Txs {
			txHash := tx.Hash()
			transactions = append(transactions, "0x"+hex.EncodeToString(txHash[:]))
		}
	} else {
		// Full transaction objects (simplified for now)
		for _, tx := range block.Txs {
			txHash := tx.Hash()
			transactions = append(transactions, map[string]interface{}{
				"hash":  "0x" + hex.EncodeToString(txHash[:]),
				"nonce": fmt.Sprintf("0x%x", tx.Nonce()),
				"from":  tx.From().Hex(),
				"gas":   fmt.Sprintf("0x%x", tx.Gas()),
				// Add more fields as needed
			})
		}
	}

	// Compute block hash
	blockHash := block.Hash()
	
	// Compute transactions root (use helper or compute from block)
	txRoot := [32]byte{} // Simplified - would need proper merkle tree
	if len(block.Txs) > 0 {
		// Use first tx hash as placeholder
		txRoot = block.Txs[0].Hash()
	}

	ethBlock := map[string]interface{}{
		"number":           fmt.Sprintf("0x%x", block.Height),
		"hash":             "0x" + hex.EncodeToString(blockHash[:]),
		"parentHash":       "0x" + hex.EncodeToString(block.PrevHash[:]),
		"nonce":            "0x0000000000000000",
		"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"logsBloom":        "0x" + strings.Repeat("0", 512),
		"transactionsRoot": "0x" + hex.EncodeToString(txRoot[:]),
		"stateRoot":        "0x" + hex.EncodeToString(block.StateRoot[:]),
		"receiptsRoot":     "0x" + hex.EncodeToString(block.ReceiptsRoot[:]),
		"miner":            ensureValidMinerAddress(block.FarmerAddr.Hex()),
		"difficulty":       fmt.Sprintf("0x%x", block.Difficulty),
		"totalDifficulty":  fmt.Sprintf("0x%x", block.CumulativeWork),
		"extraData":        "0x",
		"size":             "0x400", // Placeholder - block size not tracked
		"gasLimit":         fmt.Sprintf("0x%x", block.GasLimit),
		"gasUsed":          fmt.Sprintf("0x%x", block.GasUsed),
		"timestamp":        fmt.Sprintf("0x%x", block.TimestampUnix),
		"transactions":     transactions,
		"uncles":           []interface{}{},
	}

	return ethBlock
}

// ensureHexPrefix adds 0x prefix if not present
func ensureHexPrefix(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return "0x" + s
}

// ensureValidMinerAddress converts any address format to a valid EVM miner address
// Returns zero address if conversion fails or input is invalid
func ensureValidMinerAddress(addr string) string {
	// Zero address constant
	zeroAddress := "0x0000000000000000000000000000000000000000"
	
	// Empty or zero address
	if addr == "" {
		return zeroAddress
	}
	
	// If it starts with "arcv" or "arcv1", it's a Bech32 address - convert it
	if strings.HasPrefix(addr, "arcv") {
		evmAddr, err := address.ParseAddress(addr, "arcv")
		if err != nil {
			// Failed to parse ARCV address, return zero
			return zeroAddress
		}
		addr = evmAddr.Hex()
	}
	
	// Ensure it has 0x prefix
	if !strings.HasPrefix(addr, "0x") {
		addr = "0x" + addr
	}
	
	// Validate length: must be exactly 42 chars (0x + 40 hex digits)
	if len(addr) != 42 {
		return zeroAddress
	}
	
	// Validate hex characters
	for i := 2; i < len(addr); i++ {
		c := addr[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return zeroAddress
		}
	}
	
	return strings.ToLower(addr)
}

// ==================== Blockscout-Required Methods ====================

// web3_clientVersion returns the client version
func (h *ETHHandler) web3ClientVersion_handler() (string, error) {
	return "Archivas/v1.0.0/betanet", nil
}

// netPeerCount_handler returns the number of connected peers
func (h *ETHHandler) netPeerCount_handler() (string, error) {
	peerCount := h.getPeerCount()
	return fmt.Sprintf("0x%x", peerCount), nil
}

// getLogs_handler returns logs matching the given filter
func (h *ETHHandler) getLogs_handler(params json.RawMessage) ([]interface{}, error) {
	// Parse filter params
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing filter parameter")
	}
	
	// TODO: Implement log filtering
	// For now, return empty array (no logs)
	return []interface{}{}, nil
}

// getTransactionByHash_handler returns transaction details by hash
func (h *ETHHandler) getTransactionByHash_handler(params json.RawMessage) (interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing transaction hash parameter")
	}
	
	// TODO: Implement transaction lookup by hash
	// For now, return null (transaction not found)
	return nil, nil
}

// getStorageAt_handler returns the storage value at a given address and position
func (h *ETHHandler) getStorageAt_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 2 {
		return "", fmt.Errorf("missing address or position parameter")
	}
	
	addrStr, ok := p[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}
	
	positionStr, ok := p[1].(string)
	if !ok {
		return "", fmt.Errorf("invalid position parameter")
	}
	
	// Parse address
	addr, err := parseAddress(addrStr)
	if err != nil {
		return "", err
	}
	
	// Parse storage position (32 bytes)
	position := [32]byte{}
	positionHex := strings.TrimPrefix(positionStr, "0x")
	if len(positionHex) > 64 {
		return "", fmt.Errorf("invalid position: too long")
	}
	// Pad to 64 hex chars (32 bytes)
	positionHex = strings.Repeat("0", 64-len(positionHex)) + positionHex
	posBytes, err := hex.DecodeString(positionHex)
	if err != nil {
		return "", fmt.Errorf("invalid position hex")
	}
	copy(position[:], posBytes)
	
	// Get storage from StateDB
	value := h.stateDB.GetState(addr, position)
	
	return fmt.Sprintf("0x%x", value), nil
}

// getBlockReceipts_handler returns all receipts for a block
func (h *ETHHandler) getBlockReceipts_handler(params json.RawMessage) ([]interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing block parameter")
	}
	
	// TODO: Implement receipt lookup for all transactions in block
	// For now, return empty array (no receipts)
	return []interface{}{}, nil
}

// getTransactionByBlockNumberAndIndex_handler returns transaction by block number and index
func (h *ETHHandler) getTransactionByBlockNumberAndIndex_handler(params json.RawMessage) (interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 2 {
		return nil, fmt.Errorf("missing block number or index parameter")
	}
	
	// TODO: Implement transaction lookup by block number and index
	// For now, return null (transaction not found)
	return nil, nil
}

// getTransactionByBlockHashAndIndex_handler returns transaction by block hash and index
func (h *ETHHandler) getTransactionByBlockHashAndIndex_handler(params json.RawMessage) (interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 2 {
		return nil, fmt.Errorf("missing block hash or index parameter")
	}
	
	// TODO: Implement transaction lookup by block hash and index
	// For now, return null (transaction not found)
	return nil, nil
}

