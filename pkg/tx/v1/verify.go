package v1

import (
	"crypto/ed25519"
	"fmt"
)

// Verify verifies an ed25519 signature on a transaction.
// pubKey: 32-byte ed25519 public key
// sig: 64-byte ed25519 signature
func Verify(pubKey []byte, tx *Transfer, sig []byte) (bool, error) {
	if len(pubKey) != 32 {
		return false, fmt.Errorf("public key must be 32 bytes (ed25519), got %d", len(pubKey))
	}

	if len(sig) != 64 {
		return false, fmt.Errorf("signature must be 64 bytes (ed25519), got %d", len(sig))
	}

	// Validate transaction
	if err := tx.Validate(); err != nil {
		return false, fmt.Errorf("invalid transaction: %w", err)
	}

	// Compute transaction hash
	txHash, err := Hash(tx)
	if err != nil {
		return false, fmt.Errorf("failed to hash transaction: %w", err)
	}

	// Verify signature
	pubKeyEd := ed25519.PublicKey(pubKey)
	return ed25519.Verify(pubKeyEd, txHash[:], sig), nil
}

// VerifySignedTx verifies a SignedTx structure.
func VerifySignedTx(stx *SignedTx) (bool, error) {
	// Decode public key
	pubKey, err := DecodePubKey(stx.PubKey)
	if err != nil {
		return false, fmt.Errorf("invalid public key: %w", err)
	}

	// Decode signature
	sig, err := DecodeSig(stx.Sig)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %w", err)
	}

	// Decode and verify hash matches
	expectedHash, err := Hash(stx.Tx)
	if err != nil {
		return false, fmt.Errorf("failed to hash transaction: %w", err)
	}

	decodedHash, err := DecodeHash(stx.Hash)
	if err != nil {
		return false, fmt.Errorf("invalid hash: %w", err)
	}

	if decodedHash != expectedHash {
		return false, fmt.Errorf("hash mismatch: transaction hash does not match")
	}

	// Verify signature
	return Verify(pubKey, stx.Tx, sig)
}
