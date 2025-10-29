package ledger

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
)

// hashTransaction creates a deterministic hash of a transaction
// Does NOT include the Signature field (since we're signing the hash)
func hashTransaction(tx Transaction) []byte {
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

