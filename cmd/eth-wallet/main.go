package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║   Ethereum-Compatible Wallet Generator            ║")
	fmt.Println("║   (Compatible with MetaMask & Archivas EVM)        ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println()

	// Generate a new Ethereum-style private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		os.Exit(1)
	}

	// Get private key bytes
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Derive Ethereum address (Keccak256 of uncompressed public key, last 20 bytes)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Display results
	fmt.Println("🔐 NEW WALLET GENERATED:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Printf("Private Key (hex):\n%x\n", privateKeyBytes)
	fmt.Println()
	fmt.Printf("Ethereum Address (0x):\n%s\n", address.Hex())
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⚠️  SAVE THESE SECURELY!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("✅ Import the private key into MetaMask")
	fmt.Println("✅ The address will match exactly")
	fmt.Println("✅ Use this wallet for Archivas EVM transactions")
	fmt.Println()
}

