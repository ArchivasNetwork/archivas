package v1

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Sign signs a transaction with an ed25519 private key.
// privKey: 64-byte ed25519 private key
// Returns: signature (64 bytes), public key (32 bytes), transaction hash (32 bytes)
func Sign(privKey []byte, tx *Transfer) (sig []byte, pubKey []byte, hash [32]byte, err error) {
	if len(privKey) != 64 {
		return nil, nil, hash, fmt.Errorf("private key must be 64 bytes (ed25519), got %d", len(privKey))
	}

	// Validate transaction
	if err := tx.Validate(); err != nil {
		return nil, nil, hash, fmt.Errorf("invalid transaction: %w", err)
	}

	// Compute transaction hash
	txHash, err := Hash(tx)
	if err != nil {
		return nil, nil, hash, fmt.Errorf("failed to hash transaction: %w", err)
	}
	hash = txHash

	// Sign the hash with ed25519
	privKeyEd := ed25519.PrivateKey(privKey)
	sig = ed25519.Sign(privKeyEd, hash[:])

	// Get public key
	pubKeyEd := privKeyEd.Public().(ed25519.PublicKey)
	pubKey = make([]byte, 32)
	copy(pubKey, pubKeyEd)

	return sig, pubKey, hash, nil
}

// SignedTx represents a signed transaction in wire format.
type SignedTx struct {
	Tx     *Transfer `json:"tx"`
	PubKey string    `json:"pubkey"` // base64 or hex encoded 32-byte public key
	Sig    string    `json:"sig"`    // base64 encoded 64-byte signature
	Hash   string    `json:"hash"`   // hex encoded 32-byte transaction hash
}

// EncodePubKey encodes a public key for wire format (using base64 for compactness).
func EncodePubKey(pubKey []byte) string {
	return base64.StdEncoding.EncodeToString(pubKey)
}

// DecodePubKey decodes a public key from wire format.
func DecodePubKey(encoded string) ([]byte, error) {
	// Try base64 first
	if decoded, err := base64.StdEncoding.DecodeString(encoded); err == nil {
		if len(decoded) == 32 {
			return decoded, nil
		}
	}

	// Try hex
	if decoded, err := hex.DecodeString(encoded); err == nil {
		if len(decoded) == 32 {
			return decoded, nil
		}
	}

	return nil, fmt.Errorf("invalid public key encoding: must be 32 bytes (base64 or hex)")
}

// EncodeSig encodes a signature for wire format (base64).
func EncodeSig(sig []byte) string {
	return base64.StdEncoding.EncodeToString(sig)
}

// DecodeSig decodes a signature from wire format.
func DecodeSig(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding: %w", err)
	}

	if len(decoded) != 64 {
		return nil, fmt.Errorf("signature must be 64 bytes, got %d", len(decoded))
	}

	return decoded, nil
}

// EncodeHash encodes a transaction hash for wire format (hex).
func EncodeHash(hash [32]byte) string {
	return hex.EncodeToString(hash[:])
}

// DecodeHash decodes a transaction hash from wire format.
func DecodeHash(encoded string) ([32]byte, error) {
	var result [32]byte
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		return result, fmt.Errorf("invalid hash encoding: %w", err)
	}

	if len(decoded) != 32 {
		return result, fmt.Errorf("hash must be 32 bytes, got %d", len(decoded))
	}

	copy(result[:], decoded)
	return result, nil
}
