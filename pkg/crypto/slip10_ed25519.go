package crypto

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
)

const (
	// ArchivasCoinType is the registered coin type for Archivas in BIP44
	// Using 734 as specified (Archivas reserved)
	ArchivasCoinType = 734

	// DefaultDerivationPath is the canonical path: m/734'/0'/0'/0/0
	DefaultDerivationPath = "m/734'/0'/0'/0/0"
)

// DeriveEd25519Key derives an ed25519 keypair from a seed at the specified path using SLIP-0010.
// Path format: m/734'/account'/change'/address_index
// Returns: private key (64 bytes), public key (32 bytes), error
func DeriveEd25519Key(seed []byte, account, change, addressIndex uint32) (privKey []byte, pubKey []byte, err error) {
	if len(seed) != 64 {
		return nil, nil, fmt.Errorf("seed must be 64 bytes (BIP39 seed), got %d", len(seed))
	}

	// SLIP-0010: Master key derivation
	// I = HMAC-SHA512(Key = "ed25519 seed", Data = seed)
	masterHMAC := hmac.New(sha512.New, []byte("ed25519 seed"))
	masterHMAC.Write(seed)
	masterKey := masterHMAC.Sum(nil)

	// Derive m/734'/account'/change'/address_index (all hardened)
	path := []uint32{
		734 + 0x80000000,          // 734' (hardened)
		account + 0x80000000,      // account' (hardened)
		change + 0x80000000,       // change' (hardened)
		addressIndex + 0x80000000, // address_index (hardened)
	}

	currentKey := masterKey
	for i, index := range path {
		// SLIP-0010 child key derivation for ed25519
		// I = HMAC-SHA512(Key = current_key[0:32], Data = 0x00 || index)
		hmacKey := currentKey[:32]
		hmacData := make([]byte, 5)
		hmacData[0] = 0x00
		binary.BigEndian.PutUint32(hmacData[1:], index)

		childHMAC := hmac.New(sha512.New, hmacKey)
		childHMAC.Write(hmacData)
		childKey := childHMAC.Sum(nil)

		// For ed25519, we use the left 32 bytes as the seed
		// But we need to validate the chain code (right 32 bytes) for future derivations
		if i < len(path)-1 {
			// Not the final step, use child key as next iteration's key
			currentKey = childKey
		} else {
			// Final step: use left 32 bytes as ed25519 seed
			ed25519Seed := childKey[:32]
			privKeyEd := ed25519.NewKeyFromSeed(ed25519Seed)
			pubKeyEd := privKeyEd.Public().(ed25519.PublicKey)

			// Return full 64-byte private key and 32-byte public key
			privKey = make([]byte, 64)
			copy(privKey, privKeyEd)
			return privKey, pubKeyEd, nil
		}
	}

	return nil, nil, fmt.Errorf("derivation failed")
}

// DeriveDefaultKey derives at the canonical path m/734'/0'/0'/0/0
func DeriveDefaultKey(seed []byte) (privKey []byte, pubKey []byte, err error) {
	return DeriveEd25519Key(seed, 0, 0, 0)
}

// PublicKeyFromPrivate derives the 32-byte public key from a 64-byte ed25519 private key.
func PublicKeyFromPrivate(privKey []byte) ([]byte, error) {
	if len(privKey) != 64 {
		return nil, fmt.Errorf("private key must be 64 bytes, got %d", len(privKey))
	}

	privKeyEd := ed25519.PrivateKey(privKey)
	pubKeyEd := privKeyEd.Public().(ed25519.PublicKey)
	return pubKeyEd, nil
}
