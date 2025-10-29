package wallet

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// GenerateKeypair generates a new secp256k1 keypair
func GenerateKeypair() (privKey []byte, pubKey []byte, err error) {
	// Generate private key
	priv, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Get public key
	pub := priv.PubKey()

	return priv.Serialize(), pub.SerializeCompressed(), nil
}

// PubKeyToAddress converts a public key to a bech32 address with prefix "arcv"
func PubKeyToAddress(pubKey []byte) (string, error) {
	// Hash the public key using SHA256
	hash := sha256.Sum256(pubKey)

	// Take first 20 bytes (like Ethereum does with addresses)
	addrBytes := hash[:20]

	// Convert to bech32 with "arcv" prefix (HRP - Human Readable Part)
	// Convert to 5-bit groups for bech32
	conv, err := bech32.ConvertBits(addrBytes, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("failed to convert bits: %w", err)
	}

	// Encode as bech32
	encoded, err := bech32.Encode("arcv", conv)
	if err != nil {
		return "", fmt.Errorf("failed to encode bech32: %w", err)
	}

	return encoded, nil
}

// AddressToBytes decodes a bech32 address back to bytes
func AddressToBytes(addr string) ([]byte, error) {
	hrp, data, err := bech32.Decode(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bech32: %w", err)
	}

	if hrp != "arcv" {
		return nil, fmt.Errorf("invalid address prefix: expected 'arcv', got '%s'", hrp)
	}

	// Convert from 5-bit groups back to 8-bit
	decoded, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bits: %w", err)
	}

	return decoded, nil
}

