package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/iljanemesis/archivas/pkg/crypto"
	txv1 "github.com/iljanemesis/archivas/pkg/tx/v1"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "keygen":
		keygen()
	case "addr":
		addr()
	case "sign-transfer":
		signTransfer()
	case "broadcast":
		broadcast()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Archivas CLI v1.1.0

Usage:
  archivas-cli <command> [flags]

Commands:
  keygen                      Generate a new mnemonic and address
  addr <mnemonic>            Derive address from mnemonic
  sign-transfer              Sign a transfer transaction (see flags)
  broadcast <tx.json> <rpc>  Broadcast a signed transaction

Examples:
  archivas-cli keygen
  archivas-cli addr "word1 word2 ... word24"
  archivas-cli sign-transfer --from-mnemonic "..." --to arcv1... --amount 1000000000 --fee 100 --nonce 0 --out tx.json
  archivas-cli broadcast tx.json http://localhost:8080
`)
}

func keygen() {
	mnemonic, err := crypto.GenerateMnemonic()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating mnemonic: %v\n", err)
		os.Exit(1)
	}

	seed, err := crypto.SeedFromMnemonic(mnemonic, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving seed: %v\n", err)
		os.Exit(1)
	}

	privKey, pubKey, err := crypto.DeriveDefaultKey(seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving key: %v\n", err)
		os.Exit(1)
	}

	address, err := crypto.PubKeyToAddress(pubKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating address: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Mnemonic: %s\n", mnemonic)
	fmt.Printf("Address:  %s\n", address)
	fmt.Printf("PubKey:   %x\n", pubKey)
	fmt.Printf("PrivKey:  %x\n", privKey)
}

func addr() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: archivas-cli addr <mnemonic>\n")
		os.Exit(1)
	}

	mnemonic := strings.Join(os.Args[2:], " ")
	address, err := crypto.AddressFromMnemonic(mnemonic, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving address: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(address)
}

func signTransfer() {
	fromMnemonic := flag.String("from-mnemonic", "", "Sender mnemonic (24 words)")
	toAddr := flag.String("to", "", "Recipient address (arcv...)")
	amountStr := flag.String("amount", "", "Amount in base units")
	feeStr := flag.String("fee", "", "Fee in base units")
	nonceStr := flag.String("nonce", "", "Nonce (u64)")
	memo := flag.String("memo", "", "Optional memo (max 256 bytes)")
	outFile := flag.String("out", "", "Output file for signed transaction (JSON)")

	flag.CommandLine.Parse(os.Args[2:])

	if *fromMnemonic == "" || *toAddr == "" || *amountStr == "" || *feeStr == "" || *nonceStr == "" || *outFile == "" {
		fmt.Fprintf(os.Stderr, "All flags required: --from-mnemonic, --to, --amount, --fee, --nonce, --out\n")
		os.Exit(1)
	}

	// Parse amounts
	var amount, fee, nonce uint64
	if _, err := fmt.Sscanf(*amountStr, "%d", &amount); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid amount: %v\n", err)
		os.Exit(1)
	}
	if _, err := fmt.Sscanf(*feeStr, "%d", &fee); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid fee: %v\n", err)
		os.Exit(1)
	}
	if _, err := fmt.Sscanf(*nonceStr, "%d", &nonce); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid nonce: %v\n", err)
		os.Exit(1)
	}

	// Derive key from mnemonic
	seed, err := crypto.SeedFromMnemonic(*fromMnemonic, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving seed: %v\n", err)
		os.Exit(1)
	}

	privKey, pubKey, err := crypto.DeriveDefaultKey(seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving key: %v\n", err)
		os.Exit(1)
	}

	fromAddr, err := crypto.PubKeyToAddress(pubKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating address: %v\n", err)
		os.Exit(1)
	}

	// Create transaction
	tx := &txv1.Transfer{
		Type:   "transfer",
		From:   fromAddr,
		To:     *toAddr,
		Amount: amount,
		Fee:    fee,
		Nonce:  nonce,
		Memo:   *memo,
	}

	if err := tx.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid transaction: %v\n", err)
		os.Exit(1)
	}

	// Sign transaction
	stx, err := txv1.PackSignedTx(tx, privKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error signing transaction: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	file, err := os.Create(*outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(stx); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Signed transaction written to %s\n", *outFile)
	fmt.Printf("Transaction hash: %s\n", stx.Hash)
}

func broadcast() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: archivas-cli broadcast <tx.json> <rpc-url>\n")
		os.Exit(1)
	}

	txFile := os.Args[2]
	rpcURL := os.Args[3]

	// Read signed transaction
	file, err := os.Open(txFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Validate transaction format
	var stx txv1.SignedTx
	if err := json.Unmarshal(jsonData, &stx); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing transaction: %v\n", err)
		os.Exit(1)
	}

	// POST to /submit
	url := strings.TrimSuffix(rpcURL, "/") + "/submit"
	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Request failed (status %d): %s\n", resp.StatusCode, string(body))
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		fmt.Printf("Response: %s\n", string(body))
		os.Exit(1)
	}

	if ok, _ := result["ok"].(bool); ok {
		fmt.Printf("Transaction submitted successfully!\n")
		fmt.Printf("Hash: %v\n", result["hash"])
	} else {
		fmt.Fprintf(os.Stderr, "Transaction rejected: %v\n", result["error"])
		os.Exit(1)
	}
}
