package ledger

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

// pubKeyToAddress converts a public key to a bech32 address (same logic as wallet package)
func pubKeyToAddress(pubKey []byte) (string, error) {
	// Hash the public key using SHA256
	hash := sha256.Sum256(pubKey)

	// Take first 20 bytes
	addrBytes := hash[:20]

	// Convert to bech32 with "arcv" prefix
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

// VerifyTransactionSignature verifies that a transaction signature is valid
func VerifyTransactionSignature(tx Transaction) error {
	// Verify that SenderPubKey matches From address
	addr, err := pubKeyToAddress(tx.SenderPubKey)
	if err != nil {
		return fmt.Errorf("failed to derive address from public key: %w", err)
	}

	if addr != tx.From {
		return fmt.Errorf("public key does not match From address: expected %s, got %s", tx.From, addr)
	}

	// Parse public key
	pubKey, err := secp256k1.ParsePubKey(tx.SenderPubKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	// Parse signature
	sig, err := ecdsa.ParseDERSignature(tx.Signature)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	// Hash the transaction (same logic as wallet.HashTransaction)
	txHash := hashTransaction(tx)

	// Verify signature
	if !sig.Verify(txHash, pubKey) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}

