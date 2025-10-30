package wallet

import (
	"crypto/sha256"
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic creates a new 24-word BIP39 mnemonic
func GenerateMnemonic() (string, error) {
	// 256 bits entropy = 24 words
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// MnemonicToSeed converts a mnemonic to a master seed
func MnemonicToSeed(mnemonic string, password string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	// BIP39 seed derivation (with optional password for extra security)
	seed := bip39.NewSeed(mnemonic, password)
	return seed, nil
}

// DeriveKey derives a private key from seed using Archivas derivation path
// Path: m/734'/0'/account'/0/index
func DeriveKey(seed []byte, account uint32, index uint32) ([]byte, error) {
	// Simplified derivation for now (full BIP32/44 would be more complex)
	// We hash: seed || "archivas" || account || index
	
	data := make([]byte, 0, len(seed)+20)
	data = append(data, seed...)
	data = append(data, []byte("archivas")...)
	data = append(data, byte(account>>24), byte(account>>16), byte(account>>8), byte(account))
	data = append(data, byte(index>>24), byte(index>>16), byte(index>>8), byte(index))
	
	hash := sha256.Sum256(data)
	
	// Return 32-byte private key
	return hash[:], nil
}

// DerivationPath returns the string representation of the derivation path
func DerivationPath(account uint32, index uint32) string {
	return fmt.Sprintf("m/734'/0'/%d'/0/%d", account, index)
}

