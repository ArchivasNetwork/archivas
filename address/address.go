package address

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// EVMAddress represents a 20-byte Ethereum-style address
// This is the canonical internal representation used by the EVM and state trie
type EVMAddress [20]byte

// ZeroAddress returns an all-zero EVM address
func ZeroAddress() EVMAddress {
	return EVMAddress{}
}

// EVMAddressFromHex parses a hex string (with or without 0x prefix) into an EVMAddress
func EVMAddressFromHex(s string) (EVMAddress, error) {
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 40 {
		return EVMAddress{}, fmt.Errorf("invalid hex address length: expected 40, got %d", len(s))
	}
	
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return EVMAddress{}, fmt.Errorf("invalid hex address: %w", err)
	}
	
	var addr EVMAddress
	copy(addr[:], bytes)
	return addr, nil
}

// Hex returns the address as a 0x-prefixed hex string
func (a EVMAddress) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

// Bytes returns the address as a byte slice
func (a EVMAddress) Bytes() []byte {
	return a[:]
}

// String returns the hex representation (for logging/debugging)
func (a EVMAddress) String() string {
	return a.Hex()
}

// IsZero returns true if the address is all zeros
func (a EVMAddress) IsZero() bool {
	return a == ZeroAddress()
}

// EncodeARCVAddress converts an EVM address to a Bech32-encoded ARCV address
// Format: arcv1xxxx...
// This is the user-facing address format for Archivas
func EncodeARCVAddress(addr EVMAddress, hrp string) (string, error) {
	if hrp == "" {
		hrp = "arcv" // Default HRP
	}
	
	// Convert 20-byte address to 5-bit groups for Bech32
	converted, err := ConvertBits8To5(addr[:])
	if err != nil {
		return "", fmt.Errorf("failed to convert address bits: %w", err)
	}
	
	// Encode with Bech32
	encoded, err := bech32Encode(hrp, converted)
	if err != nil {
		return "", fmt.Errorf("failed to encode bech32: %w", err)
	}
	
	return encoded, nil
}

// DecodeARCVAddress decodes a Bech32 ARCV address back to an EVM address
// Validates HRP and checksum
func DecodeARCVAddress(s string, expectedHRP string) (EVMAddress, error) {
	if expectedHRP == "" {
		expectedHRP = "arcv" // Default HRP
	}
	
	// Decode Bech32
	hrp, data, err := bech32Decode(s)
	if err != nil {
		return EVMAddress{}, fmt.Errorf("failed to decode bech32: %w", err)
	}
	
	// Validate HRP
	if hrp != expectedHRP {
		return EVMAddress{}, fmt.Errorf("invalid HRP: expected %s, got %s", expectedHRP, hrp)
	}
	
	// Convert from 5-bit groups back to 8-bit bytes
	decoded, err := ConvertBits5To8(data)
	if err != nil {
		return EVMAddress{}, fmt.Errorf("failed to convert address bits: %w", err)
	}
	
	// Validate length
	if len(decoded) != 20 {
		return EVMAddress{}, fmt.Errorf("invalid address length: expected 20, got %d", len(decoded))
	}
	
	var addr EVMAddress
	copy(addr[:], decoded)
	return addr, nil
}

// ParseAddress parses an address string in either format:
// - 0x-prefixed hex (40 chars after 0x)
// - Bech32 ARCV address (arcv1...)
// Returns the canonical EVMAddress
func ParseAddress(s string, hrp string) (EVMAddress, error) {
	if hrp == "" {
		hrp = "arcv"
	}
	
	// Try hex format first
	if strings.HasPrefix(s, "0x") {
		return EVMAddressFromHex(s)
	}
	
	// Try Bech32 format
	if strings.HasPrefix(s, hrp+"1") {
		return DecodeARCVAddress(s, hrp)
	}
	
	return EVMAddress{}, fmt.Errorf("invalid address format: must be 0x-hex or %s1-bech32", hrp)
}

// MustParse parses an address and panics on error (for testing/constants)
func MustParse(s string, hrp string) EVMAddress {
	addr, err := ParseAddress(s, hrp)
	if err != nil {
		panic(fmt.Sprintf("failed to parse address %s: %v", s, err))
	}
	return addr
}

