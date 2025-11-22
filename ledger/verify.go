package ledger

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ArchivasNetwork/archivas/address"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	decredEcdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
)

// pubKeyToARCVAddress converts a compressed or uncompressed public key to an ARCV bech32 address
// using the CANONICAL Ethereum-compatible derivation (Keccak256-based).
//
// This replaces the old SHA256-based derivation to ensure consistency across all components.
func pubKeyToARCVAddress(pubKeyBytes []byte) (string, error) {
	// Parse the public key (handles both compressed and uncompressed formats)
	pubKey, err := secp256k1.ParsePubKey(pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %w", err)
	}

	// Convert decred public key to Go's standard ecdsa.PublicKey
	ecdsaPubKey := ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     pubKey.X(),
		Y:     pubKey.Y(),
	}

	// Derive EVM address using canonical Ethereum method (Keccak256)
	evmAddr := address.PublicKeyToEVMAddress(&ecdsaPubKey)

	// Encode as ARCV Bech32
	arcvAddr, err := address.EncodeARCVAddress(evmAddr, "arcv")
	if err != nil {
		return "", fmt.Errorf("failed to encode ARCV address: %w", err)
	}

	return arcvAddr, nil
}

// VerifyTransactionSignature verifies that a transaction signature is valid
// using the CANONICAL Ethereum-compatible address derivation.
func VerifyTransactionSignature(tx Transaction) error {
	// Verify that SenderPubKey matches From address using canonical derivation
	addr, err := pubKeyToARCVAddress(tx.SenderPubKey)
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
	sig, err := decredEcdsa.ParseDERSignature(tx.Signature)
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

