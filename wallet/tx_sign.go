package wallet

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/ArchivasNetwork/archivas/ledger"
)

// HashTransaction creates a deterministic hash of a transaction
// Does NOT include the Signature field (since we're signing the hash)
func HashTransaction(tx ledger.Transaction) []byte {
	var buf bytes.Buffer

	// Write all fields except Signature
	buf.WriteString(tx.From)
	buf.WriteString(tx.To)
	binary.Write(&buf, binary.BigEndian, tx.Amount)
	binary.Write(&buf, binary.BigEndian, tx.Fee)
	binary.Write(&buf, binary.BigEndian, tx.Nonce)
	buf.Write(tx.SenderPubKey)

	// Hash the serialized data
	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}

// SignTransaction signs a transaction with a private key
func SignTransaction(tx *ledger.Transaction, privKey []byte) error {
	// Parse private key
	priv := secp256k1.PrivKeyFromBytes(privKey)

	// Hash the transaction
	txHash := HashTransaction(*tx)

	// Sign the hash
	sig := ecdsa.Sign(priv, txHash)

	// Serialize signature to DER format
	tx.Signature = sig.Serialize()

	return nil
}

// VerifyTransactionSignature verifies that a transaction was signed by the owner of the From address
func VerifyTransactionSignature(tx ledger.Transaction) (bool, error) {
	// Verify that SenderPubKey matches From address
	addr, err := PubKeyToAddress(tx.SenderPubKey)
	if err != nil {
		return false, fmt.Errorf("failed to derive address from public key: %w", err)
	}

	if addr != tx.From {
		return false, fmt.Errorf("public key does not match From address: expected %s, got %s", tx.From, addr)
	}

	// Parse public key
	pubKey, err := secp256k1.ParsePubKey(tx.SenderPubKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Parse signature
	sig, err := ecdsa.ParseDERSignature(tx.Signature)
	if err != nil {
		return false, fmt.Errorf("failed to parse signature: %w", err)
	}

	// Hash the transaction
	txHash := HashTransaction(tx)

	// Verify signature
	valid := sig.Verify(txHash, pubKey)
	return valid, nil
}

