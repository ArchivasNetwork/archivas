package address

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

// PrivateKeyToEVMAddress derives an EVM address from a private key using Ethereum standard derivation.
//
// This is the CANONICAL address derivation method for Archivas.
// All tools (farmer, wallet CLI, RPC) MUST use this function.
//
// Derivation method (Ethereum standard):
//   1. Private key (32 bytes) → ECDSA public key (secp256k1)
//   2. Serialize public key as UNCOMPRESSED (65 bytes: 0x04 prefix + 64 bytes)
//   3. Keccak256 hash of the 64-byte public key (without 0x04 prefix)
//   4. Take LAST 20 bytes of the hash
//
// This ensures compatibility with:
//   - MetaMask
//   - Ethereum tools (ethers.js, web3.js, Hardhat)
//   - Standard Ethereum wallets
//
// Example:
//   privKey := []byte{...} // 32 bytes
//   addr, err := PrivateKeyToEVMAddress(privKey)
//   if err != nil {
//       // handle error
//   }
//   fmt.Println(addr.Hex()) // 0x...
func PrivateKeyToEVMAddress(privateKeyBytes []byte) (EVMAddress, error) {
	if len(privateKeyBytes) != 32 {
		return EVMAddress{}, fmt.Errorf("private key must be 32 bytes, got %d", len(privateKeyBytes))
	}

	// Parse private key using go-ethereum's crypto package
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return EVMAddress{}, fmt.Errorf("invalid private key: %w", err)
	}

	// Derive Ethereum address using standard method
	// This uses Keccak256(uncompressed_public_key)[12:]
	ethAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Convert to our EVMAddress type
	var addr EVMAddress
	copy(addr[:], ethAddr.Bytes())

	return addr, nil
}

// PublicKeyToEVMAddress derives an EVM address from an ECDSA public key.
// Uses Ethereum standard: Keccak256(uncompressed_public_key)[12:]
//
// This is useful when you already have the public key and don't need
// to derive it from a private key.
func PublicKeyToEVMAddress(publicKey *ecdsa.PublicKey) EVMAddress {
	ethAddr := crypto.PubkeyToAddress(*publicKey)
	
	var addr EVMAddress
	copy(addr[:], ethAddr.Bytes())
	
	return addr
}

// PrivateKeyToARCVAddress derives a Bech32 ARCV address from a private key.
// This is a convenience function that combines PrivateKeyToEVMAddress and EncodeARCVAddress.
//
// Example:
//   privKey := []byte{...} // 32 bytes
//   arcvAddr, err := PrivateKeyToARCVAddress(privKey, "arcv")
//   if err != nil {
//       // handle error
//   }
//   fmt.Println(arcvAddr) // arcv1...
func PrivateKeyToARCVAddress(privateKeyBytes []byte, hrp string) (string, error) {
	if hrp == "" {
		hrp = "arcv"
	}

	// Derive EVM address
	evmAddr, err := PrivateKeyToEVMAddress(privateKeyBytes)
	if err != nil {
		return "", err
	}

	// Encode as Bech32
	return EncodeARCVAddress(evmAddr, hrp)
}

// ValidateAddressRoundTrip validates that an address can be converted between formats without loss.
// This is useful for testing address conversion integrity.
//
// Returns nil if successful, error otherwise.
func ValidateAddressRoundTrip(evmAddr EVMAddress, hrp string) error {
	if hrp == "" {
		hrp = "arcv"
	}

	// EVM → ARCV
	arcvAddr, err := EncodeARCVAddress(evmAddr, hrp)
	if err != nil {
		return fmt.Errorf("EVM→ARCV failed: %w", err)
	}

	// ARCV → EVM
	evmAddr2, err := DecodeARCVAddress(arcvAddr, hrp)
	if err != nil {
		return fmt.Errorf("ARCV→EVM failed: %w", err)
	}

	// Verify they match
	if evmAddr != evmAddr2 {
		return fmt.Errorf("round trip mismatch: %s != %s", evmAddr.Hex(), evmAddr2.Hex())
	}

	return nil
}

