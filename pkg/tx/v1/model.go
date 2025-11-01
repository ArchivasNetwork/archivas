package v1

import (
	"fmt"
)

// Transfer represents a v1 transfer transaction.
// All amounts are in base units (uint64).
type Transfer struct {
	Type   string `json:"type"`           // Always "transfer"
	From   string `json:"from"`           // Bech32 address (arcv...)
	To     string `json:"to"`             // Bech32 address (arcv...)
	Amount uint64 `json:"amount"`         // Base units (u64)
	Fee    uint64 `json:"fee"`            // Base units (u64)
	Nonce  uint64 `json:"nonce"`          // u64
	Memo   string `json:"memo,omitempty"` // Optional UTF-8, max 256 bytes
}

// Validate checks that the transfer transaction is valid.
func (t *Transfer) Validate() error {
	if t.Type != "transfer" {
		return fmt.Errorf("invalid type: expected 'transfer', got '%s'", t.Type)
	}

	if t.From == "" {
		return fmt.Errorf("from address is required")
	}

	if t.To == "" {
		return fmt.Errorf("to address is required")
	}

	if t.Amount == 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if t.Fee == 0 {
		return fmt.Errorf("fee must be greater than 0")
	}

	// Validate memo length (if present)
	if len(t.Memo) > 256 {
		return fmt.Errorf("memo exceeds maximum length: %d bytes (max 256)", len(t.Memo))
	}

	// Validate UTF-8 encoding
	if !isValidUTF8(t.Memo) {
		return fmt.Errorf("memo must be valid UTF-8")
	}

	return nil
}

// isValidUTF8 checks if a string is valid UTF-8.
func isValidUTF8(s string) bool {
	for _, r := range s {
		if r == 0xFFFD {
			return false // Invalid UTF-8 rune
		}
	}
	return true
}
