package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArchivasNetwork/archivas/address"
)

// ARCVHandler provides Archivas-specific RPC endpoints
// Phase 3: Address conversion and chain-specific queries
type ARCVHandler struct {
	bech32Prefix string
}

// NewARCVHandler creates a new ARCV RPC handler
func NewARCVHandler(bech32Prefix string) *ARCVHandler {
	return &ARCVHandler{
		bech32Prefix: bech32Prefix,
	}
}

// ServeHTTP handles ARCV RPC requests
func (h *ARCVHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	case "arcv_toHexAddress":
		result, err = h.toHexAddress_handler(req.Params)
	case "arcv_fromHexAddress":
		result, err = h.fromHexAddress_handler(req.Params)
	case "arcv_validateAddress":
		result, err = h.validateAddress_handler(req.Params)
	case "arcv_getAddressInfo":
		result, err = h.getAddressInfo_handler(req.Params)
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

// arcv_toHexAddress converts a Bech32 ARCV address to hex format
func (h *ARCVHandler) toHexAddress_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	arcvAddr, ok := p[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	addr, err := address.DecodeARCVAddress(arcvAddr, h.bech32Prefix)
	if err != nil {
		return "", fmt.Errorf("invalid ARCV address: %w", err)
	}

	return addr.Hex(), nil
}

// arcv_fromHexAddress converts a hex address to Bech32 ARCV format
func (h *ARCVHandler) fromHexAddress_handler(params json.RawMessage) (string, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return "", err
	}
	if len(p) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	hexAddr, ok := p[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	addr, err := address.EVMAddressFromHex(hexAddr)
	if err != nil {
		return "", fmt.Errorf("invalid hex address: %w", err)
	}

	arcvAddr, err := address.EncodeARCVAddress(addr, h.bech32Prefix)
	if err != nil {
		return "", fmt.Errorf("failed to encode ARCV address: %w", err)
	}

	return arcvAddr, nil
}

// arcv_validateAddress validates an address in either format
func (h *ARCVHandler) validateAddress_handler(params json.RawMessage) (map[string]interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing address parameter")
	}

	addrStr, ok := p[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	result := map[string]interface{}{
		"valid":  false,
		"format": "unknown",
	}

	// Try parsing as either format
	addr, err := address.ParseAddress(addrStr, h.bech32Prefix)
	if err != nil {
		result["error"] = err.Error()
		return result, nil
	}

	result["valid"] = true
	if addrStr[:2] == "0x" {
		result["format"] = "hex"
		arcvAddr, _ := address.EncodeARCVAddress(addr, h.bech32Prefix)
		result["hex"] = addr.Hex()
		result["arcv"] = arcvAddr
	} else {
		result["format"] = "bech32"
		result["hex"] = addr.Hex()
		result["arcv"] = addrStr
	}

	return result, nil
}

// arcv_getAddressInfo returns information about an address
func (h *ARCVHandler) getAddressInfo_handler(params json.RawMessage) (map[string]interface{}, error) {
	var p []interface{}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	if len(p) < 1 {
		return nil, fmt.Errorf("missing address parameter")
	}

	addrStr, ok := p[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	addr, err := address.ParseAddress(addrStr, h.bech32Prefix)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	arcvAddr, _ := address.EncodeARCVAddress(addr, h.bech32Prefix)

	result := map[string]interface{}{
		"hex":        addr.Hex(),
		"arcv":       arcvAddr,
		"isZero":     addr.IsZero(),
		"bytes":      len(addr.Bytes()),
		"network":    h.bech32Prefix,
	}

	return result, nil
}

