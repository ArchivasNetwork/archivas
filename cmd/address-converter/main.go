package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ArchivasNetwork/archivas/address"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  address-converter <arcv1_address>  - Convert Bech32 to 0x")
		fmt.Println("  address-converter <0x_address>     - Convert 0x to Bech32")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  address-converter arcv1s9m9avxdkzuv9lf6wle2r2sklcrq3ayhc8txqs")
		fmt.Println("  address-converter 0x82b6cbb5bb79b99fc734eb381511fb6f88d0e0f9")
		os.Exit(1)
	}

	input := strings.TrimSpace(os.Args[1])

	// Check if input is Bech32 (starts with arcv1)
	if strings.HasPrefix(input, "arcv1") {
		// Decode Bech32 to EVM address
		addr, err := address.DecodeARCVAddress(input, "arcv")
		if err != nil {
			fmt.Printf("Error decoding Bech32 address: %v\n", err)
			os.Exit(1)
		}

		// Convert to 0x hex
		hexAddr := "0x" + hex.EncodeToString(addr[:])
		
		fmt.Println("╔════════════════════════════════════════════════════╗")
		fmt.Println("║        Archivas Address Converter                 ║")
		fmt.Println("╠════════════════════════════════════════════════════╣")
		fmt.Printf("║ Input (Bech32):  %s\n", input)
		fmt.Printf("║ Output (0x):     %s\n", hexAddr)
		fmt.Println("╚════════════════════════════════════════════════════╝")
		fmt.Println("")
		fmt.Println("✅ Use this 0x address in MetaMask")
		fmt.Println("✅ Your MetaMask should show this address after importing private key")

	} else if strings.HasPrefix(input, "0x") || len(input) == 40 {
		// Decode hex to EVM address
		if strings.HasPrefix(input, "0x") {
			input = input[2:] // Remove 0x prefix
		}

		addrBytes, err := hex.DecodeString(input)
		if err != nil {
			fmt.Printf("Error decoding hex address: %v\n", err)
			os.Exit(1)
		}

		if len(addrBytes) != 20 {
			fmt.Printf("Error: address must be 20 bytes (40 hex chars), got %d bytes\n", len(addrBytes))
			os.Exit(1)
		}

		var addr address.EVMAddress
		copy(addr[:], addrBytes)

		// Encode to Bech32
		bech32Addr, err := address.EncodeARCVAddress(addr, "arcv")
		if err != nil {
			fmt.Printf("Error encoding Bech32 address: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("╔════════════════════════════════════════════════════╗")
		fmt.Println("║        Archivas Address Converter                 ║")
		fmt.Println("╠════════════════════════════════════════════════════╣")
		fmt.Printf("║ Input (0x):      0x%s\n", input)
		fmt.Printf("║ Output (Bech32): %s\n", bech32Addr)
		fmt.Println("╚════════════════════════════════════════════════════╝")
		fmt.Println("")
		fmt.Println("✅ Use this arcv1 address for Archivas CLI commands")
		fmt.Println("✅ This is your public receiving address")

	} else {
		fmt.Println("Error: Invalid address format")
		fmt.Println("Expected: arcv1... (Bech32) or 0x... (hex)")
		os.Exit(1)
	}
}

