package crypto

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/bech32"
	"golang.org/x/crypto/blake2b"
)

const (
	// AddressHRP is the Bech32 Human Readable Part for Archivas addresses
	AddressHRP = "arcv"
)

// PubKeyToAddress converts a 32-byte ed25519 public key to a Bech32 address.
// Process: blake2b-160(pubkey) -> Bech32 encode with HRP "arcv"
func PubKeyToAddress(pubKey []byte) (string, error) {
	if len(pubKey) != 32 {
		return "", fmt.Errorf("public key must be 32 bytes (ed25519), got %d", len(pubKey))
	}

	// Hash public key with blake2b-160 (20 bytes output)
	hash, err := blake2b.New(20, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create blake2b-160 hasher: %w", err)
	}

	hash.Write(pubKey)
	addrBytes := hash.Sum(nil)

	// Convert to 5-bit groups for bech32
	conv, err := bech32.ConvertBits(addrBytes, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("failed to convert bits: %w", err)
	}

	// Encode as bech32 with "arcv" HRP
	encoded, err := bech32.Encode(AddressHRP, conv)
	if err != nil {
		return "", fmt.Errorf("failed to encode bech32: %w", err)
	}

	return encoded, nil
}

// ParseAddress decodes a Bech32 address and returns the 20-byte hash.
func ParseAddress(addr string) ([20]byte, error) {
	var result [20]byte

	hrp, data, err := bech32.Decode(addr)
	if err != nil {
		return result, fmt.Errorf("failed to decode bech32: %w", err)
	}

	if hrp != AddressHRP {
		return result, fmt.Errorf("invalid address HRP: expected '%s', got '%s'", AddressHRP, hrp)
	}

	// Convert from 5-bit groups back to 8-bit bytes
	decoded, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return result, fmt.Errorf("failed to convert bits: %w", err)
	}

	if len(decoded) != 20 {
		return result, fmt.Errorf("invalid address length: expected 20 bytes, got %d", len(decoded))
	}

	copy(result[:], decoded)
	return result, nil
}

// IsValidAddress checks if a string is a valid Archivas Bech32 address.
func IsValidAddress(addr string) bool {
	_, err := ParseAddress(addr)
	return err == nil
}

// AddressFromMnemonic derives an address from a mnemonic at the default path.
// Convenience function: mnemonic -> seed -> key -> address
func AddressFromMnemonic(mnemonic string, passphrase string) (string, error) {
	seed, err := SeedFromMnemonic(mnemonic, passphrase)
	if err != nil {
		return "", err
	}

	_, pubKey, err := DeriveDefaultKey(seed)
	if err != nil {
		return "", err
	}

	return PubKeyToAddress(pubKey)
}
