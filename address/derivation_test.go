package address

import (
	"encoding/hex"
	"testing"
)

// TestPrivateKeyToEVMAddress_KnownVectors tests address derivation against known Ethereum test vectors
func TestPrivateKeyToEVMAddress_KnownVectors(t *testing.T) {
	tests := []struct {
		name       string
		privateKey string
		wantAddr   string
	}{
		{
			name:       "Test vector 1",
			privateKey: "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9",
			wantAddr:   "0x39a028dfdcae40bf277ec1ec268d62665d36c073",
		},
		{
			name:       "Test vector 2",
			privateKey: "4805f457480731b65d4363b1f3be700071c7c3852c3436b869dc4bc5f29991d0",
			wantAddr:   "0x7ca9e60a6f541d0ff69b479a268a7e4390f44bb3",
		},
		{
			name:       "Test vector 3 (farming wallet)",
			privateKey: "1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899",
			wantAddr:   "0x47ea4b22029c155c835fd0a0b99f8196766f406a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privKeyBytes, err := hex.DecodeString(tt.privateKey)
			if err != nil {
				t.Fatalf("Failed to decode private key: %v", err)
			}

			gotAddr, err := PrivateKeyToEVMAddress(privKeyBytes)
			if err != nil {
				t.Fatalf("PrivateKeyToEVMAddress() error = %v", err)
			}

			if gotAddr.Hex() != tt.wantAddr && gotAddr.Hex() != toLowerCase(tt.wantAddr) {
				t.Errorf("PrivateKeyToEVMAddress() = %v, want %v", gotAddr.Hex(), tt.wantAddr)
			}
		})
	}
}

// TestPrivateKeyToARCVAddress tests ARCV address derivation
func TestPrivateKeyToARCVAddress(t *testing.T) {
	privKeyHex := "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9"
	privKeyBytes, _ := hex.DecodeString(privKeyHex)

	arcvAddr, err := PrivateKeyToARCVAddress(privKeyBytes, "arcv")
	if err != nil {
		t.Fatalf("PrivateKeyToARCVAddress() error = %v", err)
	}

	t.Logf("ARCV address: %s", arcvAddr)

	// Verify it starts with arcv1
	if len(arcvAddr) < 5 || arcvAddr[:5] != "arcv1" {
		t.Errorf("ARCV address should start with arcv1, got: %s", arcvAddr)
	}

	// Verify round trip
	evmAddr, err := DecodeARCVAddress(arcvAddr, "arcv")
	if err != nil {
		t.Fatalf("DecodeARCVAddress() error = %v", err)
	}

	expectedEvmAddr, _ := PrivateKeyToEVMAddress(privKeyBytes)
	if evmAddr != expectedEvmAddr {
		t.Errorf("Round trip failed: got %s, want %s", evmAddr.Hex(), expectedEvmAddr.Hex())
	}
}

// TestAddressRoundTrip tests EVM ↔ ARCV conversion
func TestAddressRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		evmAddr string
	}{
		{"Address 1", "0x39a028dfdcae40bf277ec1ec268d62665d36c073"},
		{"Address 2", "0x7ca9e60a6f541d0ff69b479a268a7e4390f44bb3"},
		{"Address 3", "0x47ea4b22029c155c835fd0a0b99f8196766f406a"},
		{"Zero address", "0x0000000000000000000000000000000000000000"},
		{"Max address", "0xffffffffffffffffffffffffffffffffffffffff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evmAddr, err := EVMAddressFromHex(tt.evmAddr)
			if err != nil {
				t.Fatalf("EVMAddressFromHex() error = %v", err)
			}

			// EVM → ARCV
			arcvAddr, err := EncodeARCVAddress(evmAddr, "arcv")
			if err != nil {
				t.Fatalf("EncodeARCVAddress() error = %v", err)
			}

			// ARCV → EVM
			evmAddr2, err := DecodeARCVAddress(arcvAddr, "arcv")
			if err != nil {
				t.Fatalf("DecodeARCVAddress() error = %v", err)
			}

			// Verify round trip
			if evmAddr != evmAddr2 {
				t.Errorf("Round trip failed: original=%s, after=%s", evmAddr.Hex(), evmAddr2.Hex())
			}

			t.Logf("EVM: %s → ARCV: %s → EVM: %s ✓", evmAddr.Hex(), arcvAddr, evmAddr2.Hex())
		})
	}
}

// TestValidateAddressRoundTrip tests the validation helper
func TestValidateAddressRoundTrip(t *testing.T) {
	evmAddr, _ := EVMAddressFromHex("0x39a028dfdcae40bf277ec1ec268d62665d36c073")

	err := ValidateAddressRoundTrip(evmAddr, "arcv")
	if err != nil {
		t.Errorf("ValidateAddressRoundTrip() error = %v", err)
	}
}

// TestConsistentDerivation ensures the same private key always produces the same address
func TestConsistentDerivation(t *testing.T) {
	privKeyHex := "1cb7a7ad1c75b0dcf142f7c4dbd01102971bc9892dae91badf472c35843f4899"
	privKeyBytes, _ := hex.DecodeString(privKeyHex)

	// Derive address 100 times
	var addrs []EVMAddress
	for i := 0; i < 100; i++ {
		addr, err := PrivateKeyToEVMAddress(privKeyBytes)
		if err != nil {
			t.Fatalf("Derivation %d failed: %v", i, err)
		}
		addrs = append(addrs, addr)
	}

	// All should match
	firstAddr := addrs[0]
	for i, addr := range addrs {
		if addr != firstAddr {
			t.Errorf("Derivation %d produced different address: %s != %s", i, addr.Hex(), firstAddr.Hex())
		}
	}

	t.Logf("✓ 100 derivations all produced: %s", firstAddr.Hex())
}

// TestParseAddress_BothFormats tests parsing of both hex and ARCV addresses
func TestParseAddress_BothFormats(t *testing.T) {
	privKeyHex := "fe4e0b573e892b9abc7692782aa14bc7560e6e7637948403827a6242dcf0f2b9"
	privKeyBytes, _ := hex.DecodeString(privKeyHex)

	// Derive canonical address
	canonicalAddr, _ := PrivateKeyToEVMAddress(privKeyBytes)
	arcvAddr, _ := EncodeARCVAddress(canonicalAddr, "arcv")

	tests := []struct {
		name      string
		input     string
		wantAddr  EVMAddress
		wantError bool
	}{
		{
			name:     "Parse hex address",
			input:    canonicalAddr.Hex(),
			wantAddr: canonicalAddr,
		},
		{
			name:     "Parse ARCV address",
			input:    arcvAddr,
			wantAddr: canonicalAddr,
		},
		{
			name:      "Invalid format",
			input:     "invalid_address",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddr, err := ParseAddress(tt.input, "arcv")
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseAddress() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseAddress() error = %v", err)
			}

			if gotAddr != tt.wantAddr {
				t.Errorf("ParseAddress() = %s, want %s", gotAddr.Hex(), tt.wantAddr.Hex())
			}
		})
	}
}

// TestInvalidPrivateKeys tests error handling for invalid inputs
func TestInvalidPrivateKeys(t *testing.T) {
	tests := []struct {
		name       string
		privateKey []byte
	}{
		{"Empty", []byte{}},
		{"Too short", []byte{1, 2, 3}},
		{"Too long", make([]byte, 33)},
		{"All zeros", make([]byte, 32)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := PrivateKeyToEVMAddress(tt.privateKey)
			if err == nil {
				t.Errorf("PrivateKeyToEVMAddress() expected error for %s, got nil", tt.name)
			}
		})
	}
}

// Helper function to convert to lowercase for case-insensitive comparison
func toLowerCase(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'F' {
			c = c + ('a' - 'A')
		}
		result[i] = c
	}
	return string(result)
}

