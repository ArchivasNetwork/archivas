package types

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash256 computes SHA256 hash of data
func Hash256(data []byte) [32]byte {
	return sha256.Sum256(data)
}

// HashToString converts a 32-byte hash to hex string
func HashToString(hash [32]byte) string {
	return "0x" + hex.EncodeToString(hash[:])
}

// StringToHash converts a hex string to 32-byte hash
func StringToHash(s string) ([32]byte, error) {
	var hash [32]byte
	
	// Remove 0x prefix if present
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}
	
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return hash, err
	}
	
	copy(hash[:], bytes)
	return hash, nil
}

// EmptyHash returns an all-zero hash
func EmptyHash() [32]byte {
	return [32]byte{}
}

// IsEmptyHash checks if a hash is all zeros
func IsEmptyHash(hash [32]byte) bool {
	return hash == EmptyHash()
}

