package main

import (
	"fmt"
	"os"

	"github.com/ArchivasNetwork/archivas/address"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("   Archivas Betanet - Dual Address System Demo")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	// Example 1: Create from hex
	fmt.Println("📍 Example 1: Create address from hex")
	fmt.Println("─────────────────────────────────────────────────")
	hexAddr := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"
	fmt.Printf("Input (hex):     %s\n", hexAddr)

	addr1, err := address.EVMAddressFromHex(hexAddr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Internal (EVM):  %s\n", addr1.Hex())

	arcv1, _ := address.EncodeARCVAddress(addr1, "arcv")
	fmt.Printf("External (ARCV): %s\n", arcv1)
	fmt.Println()

	// Example 2: Parse ARCV address
	fmt.Println("📍 Example 2: Parse ARCV address")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("Input (ARCV):    %s\n", arcv1)

	addr2, err := address.DecodeARCVAddress(arcv1, "arcv")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Internal (EVM):  %s\n", addr2.Hex())
	fmt.Printf("Match:           %v\n", addr1 == addr2)
	fmt.Println()

	// Example 3: Parse either format
	fmt.Println("📍 Example 3: Universal parser (accepts both formats)")
	fmt.Println("─────────────────────────────────────────────────")

	inputs := []string{
		"0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		arcv1,
	}

	for i, input := range inputs {
		parsed, err := address.ParseAddress(input, "arcv")
		if err != nil {
			fmt.Printf("  [%d] Error: %v\n", i+1, err)
			continue
		}

		inputType := "hex"
		if input[:4] == "arcv" {
			inputType = "ARCV"
		}

		fmt.Printf("  [%d] Input (%s): %s\n", i+1, inputType, input)
		fmt.Printf("      Parsed (EVM):   %s\n", parsed.Hex())
		fmt.Println()
	}

	// Example 4: Zero address
	fmt.Println("📍 Example 4: Zero address")
	fmt.Println("─────────────────────────────────────────────────")
	zero := address.ZeroAddress()
	fmt.Printf("Zero (hex):      %s\n", zero.Hex())

	zeroARCV, _ := address.EncodeARCVAddress(zero, "arcv")
	fmt.Printf("Zero (ARCV):     %s\n", zeroARCV)
	fmt.Printf("IsZero:          %v\n", zero.IsZero())
	fmt.Println()

	// Example 5: Roundtrip test
	fmt.Println("📍 Example 5: Roundtrip integrity test")
	fmt.Println("─────────────────────────────────────────────────")

	testAddrs := []string{
		"0x0000000000000000000000000000000000000000",
		"0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
		"0x1234567890abcdef1234567890abcdef12345678",
	}

	allPassed := true
	for i, testHex := range testAddrs {
		original, _ := address.EVMAddressFromHex(testHex)
		encoded, _ := address.EncodeARCVAddress(original, "arcv")
		decoded, _ := address.DecodeARCVAddress(encoded, "arcv")

		passed := original == decoded
		allPassed = allPassed && passed

		status := "✅"
		if !passed {
			status = "❌"
		}

		fmt.Printf("  [%d] %s %s\n", i+1, status, testHex)
		if !passed {
			fmt.Printf("       Original: %s\n", original.Hex())
			fmt.Printf("       After:    %s\n", decoded.Hex())
		}
	}

	fmt.Println()
	if allPassed {
		fmt.Println("✅ All roundtrip tests passed!")
	} else {
		fmt.Println("❌ Some tests failed")
	}
	fmt.Println()

	// Summary
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("   Summary")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("📝 The dual address system allows:")
	fmt.Println("   • Internal EVM execution with 0x addresses")
	fmt.Println("   • User-facing Bech32 addresses (arcv1...)")
	fmt.Println("   • Seamless conversion between formats")
	fmt.Println("   • Checksum validation via Bech32")
	fmt.Println("   • Compatibility with Ethereum tooling (0x)")
	fmt.Println()
	fmt.Println("🎯 Use cases:")
	fmt.Println("   • eth_* RPC: Use 0x addresses (Ethereum-compatible)")
	fmt.Println("   • arcv_* RPC: Use arcv1 addresses (user-friendly)")
	fmt.Println("   • CLI/wallets: Accept both formats")
	fmt.Println("   • State storage: Always use EVMAddress internally")
	fmt.Println()
}

