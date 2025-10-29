package consensus

import (
	"crypto/sha256"
	"encoding/binary"
)

// GenerateChallenge creates a challenge hash for a given block height and previous block hash
func GenerateChallenge(prevBlockHash [32]byte, height uint64) [32]byte {
	h := sha256.New()
	h.Write(prevBlockHash[:])
	binary.Write(h, binary.LittleEndian, height)
	return sha256.Sum256(h.Sum(nil)) // Double SHA256
}

// GenerateGenesisChallenge creates the challenge for height 0
func GenerateGenesisChallenge() [32]byte {
	// Genesis challenge is hash of chain name
	return sha256.Sum256([]byte("Archivas Devnet Genesis"))
}

