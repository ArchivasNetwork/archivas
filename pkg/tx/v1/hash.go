package v1

import (
	"golang.org/x/crypto/blake2b"
)

const (
	// DomainSeparator is the domain separation string for Archivas TxV1 hashing
	DomainSeparator = "Archivas-TxV1"
)

// Hash computes the blake2b-256 hash of a canonical JSON transaction.
// Process: DomainSeparator || canonical_json_bytes -> blake2b-256
func Hash(tx *Transfer) ([32]byte, error) {
	var result [32]byte

	// Get canonical JSON
	canonJSON, err := CanonicalJSON(tx)
	if err != nil {
		return result, err
	}

	// Create blake2b-256 hasher
	hasher, err := blake2b.New256(nil)
	if err != nil {
		return result, err
	}

	// Hash: DomainSeparator || canonical_json_bytes
	hasher.Write([]byte(DomainSeparator))
	hasher.Write(canonJSON)
	hasher.Sum(result[:0])

	return result, nil
}
