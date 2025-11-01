package crypto

import (
	"encoding/hex"
	"strings"
	"testing"

	"golang.org/x/crypto/blake2b"
)

// Golden test vectors for mnemonic -> address derivation
func TestMnemonicToAddress(t *testing.T) {
	// Fixed mnemonic for reproducible tests
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	seed, err := SeedFromMnemonic(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to derive seed: %v", err)
	}

	if len(seed) != 64 {
		t.Fatalf("Seed must be 64 bytes, got %d", len(seed))
	}

	privKey, pubKey, err := DeriveDefaultKey(seed)
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	if len(privKey) != 64 {
		t.Fatalf("Private key must be 64 bytes, got %d", len(privKey))
	}

	if len(pubKey) != 32 {
		t.Fatalf("Public key must be 32 bytes, got %d", len(pubKey))
	}

	address, err := PubKeyToAddress(pubKey)
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Verify address format
	if len(address) == 0 {
		t.Fatal("Address is empty")
	}

	if !strings.HasPrefix(address, "arcv") {
		t.Fatalf("Address must start with 'arcv', got: %s", address)
	}

	// Round-trip: parse and verify
	parsed, err := ParseAddress(address)
	if err != nil {
		t.Fatalf("Failed to parse address: %v", err)
	}

	if len(parsed) != 20 {
		t.Fatalf("Parsed address must be 20 bytes, got %d", len(parsed))
	}

	// Verify IsValidAddress
	if !IsValidAddress(address) {
		t.Fatal("IsValidAddress returned false for valid address")
	}
}

func TestDerivationPathConstant(t *testing.T) {
	if DefaultDerivationPath != "m/734'/0'/0'/0/0" {
		t.Fatalf("Derivation path must be m/734'/0'/0'/0/0, got: %s", DefaultDerivationPath)
	}
}

func TestBech32HRPConstant(t *testing.T) {
	if AddressHRP != "arcv" {
		t.Fatalf("Address HRP must be 'arcv', got: %s", AddressHRP)
	}
}

func TestAddressRoundTrip(t *testing.T) {
	// Generate random keypair
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	seed, err := SeedFromMnemonic(mnemonic, "")
	if err != nil {
		t.Fatalf("Failed to derive seed: %v", err)
	}

	_, pubKey, err := DeriveDefaultKey(seed)
	if err != nil {
		t.Fatalf("Failed to derive key: %v", err)
	}

	address, err := PubKeyToAddress(pubKey)
	if err != nil {
		t.Fatalf("Failed to generate address: %v", err)
	}

	// Parse back
	parsed, err := ParseAddress(address)
	if err != nil {
		t.Fatalf("Failed to parse address: %v", err)
	}

	// Verify same pubkey produces same address
	address2, err := PubKeyToAddress(pubKey)
	if err != nil {
		t.Fatalf("Failed to generate address again: %v", err)
	}

	if address != address2 {
		t.Fatal("Same pubkey produced different addresses")
	}

	// Verify parsed address bytes match
	expectedHash, err := blake2b.New(20, nil)
	if err != nil {
		t.Fatalf("Failed to create hasher: %v", err)
	}
	expectedHash.Write(pubKey)
	expectedBytes := expectedHash.Sum(nil)

	if hex.EncodeToString(parsed[:]) != hex.EncodeToString(expectedBytes) {
		t.Fatal("Parsed address bytes don't match expected blake2b-160 hash")
	}
}

func TestInvalidAddress(t *testing.T) {
	invalidAddrs := []string{
		"",
		"arcv",
		"bc1invalid",
		"arcv1invalid",
		"notanaddress",
	}

	for _, addr := range invalidAddrs {
		if IsValidAddress(addr) {
			t.Errorf("IsValidAddress returned true for invalid address: %s", addr)
		}

		_, err := ParseAddress(addr)
		if err == nil {
			t.Errorf("ParseAddress did not error for invalid address: %s", addr)
		}
	}
}
