package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: transfer-to-eth-address <from_arcv_address> <to_0x_address> <amount_rchv>")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  transfer-to-eth-address arcv1rlx93wlk26ny67zqk8eejfkl4y2az22nynqrtj 0x47Ea4b22029c155C835Fd0a0b99F8196766F406A 30000")
		os.Exit(1)
	}

	fromAddr := os.Args[1]
	toAddr := os.Args[2]
	amount := os.Args[3]

	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║        Transfer RCHV to Ethereum Address          ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("From (Archivas): %s\n", fromAddr)
	fmt.Printf("To (Ethereum):   %s\n", toAddr)
	fmt.Printf("Amount:          %s RCHV\n", amount)
	fmt.Println()
	fmt.Println("⚠️  This will create a transaction to send your farming")
	fmt.Println("   rewards to your Ethereum-compatible address.")
	fmt.Println()
	fmt.Println("Note: You'll need to sign this transaction with your")
	fmt.Println("      farming private key.")
	fmt.Println()

	// TODO: Implement transaction creation and submission
	// For now, show the curl command they need to run

	rpcURL := "https://seed3.betanet.archivas.ai"

	payload := map[string]interface{}{
		"from":   fromAddr,
		"to":     toAddr,
		"amount": amount,
	}

	payloadBytes, _ := json.MarshalIndent(payload, "", "  ")

	fmt.Println("To complete the transfer, use the Archivas RPC:")
	fmt.Println()
	fmt.Printf("curl -s %s/submitTx -X POST \\\n", rpcURL)
	fmt.Printf("  -H 'Content-Type: application/json' \\\n")
	fmt.Printf("  -d '%s'\n", string(payloadBytes))
	fmt.Println()
}

// Helper function to query balance
func queryBalance(rpcURL, address string) (string, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  []interface{}{address, "latest"},
		"id":      1,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if res, ok := result["result"].(string); ok {
		return res, nil
	}

	return "", fmt.Errorf("unexpected response format")
}

