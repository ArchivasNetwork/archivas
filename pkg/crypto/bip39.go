package crypto

import (
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

// SeedFromMnemonic derives a BIP39 seed from a 24-word mnemonic and optional passphrase.
// Returns 64 bytes (512 bits) seed.
func SeedFromMnemonic(mnemonic string, passphrase string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid BIP39 mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, passphrase)
	if len(seed) != 64 {
		return nil, fmt.Errorf("invalid seed length: expected 64 bytes, got %d", len(seed))
	}

	return seed, nil
}

// GenerateMnemonic generates a new 24-word BIP39 mnemonic (256 bits entropy).
func GenerateMnemonic() (string, error) {
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

// ValidateMnemonic checks if a mnemonic string is valid according to BIP39.
func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}
