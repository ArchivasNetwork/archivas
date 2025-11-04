package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ArchivasNetwork/archivas/ledger"
	"github.com/ArchivasNetwork/archivas/wallet"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "new":
		cmdNew()
	case "send":
		cmdSend()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Archivas Wallet CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  archivas-wallet new                           Generate a new wallet")
	fmt.Println("  archivas-wallet send [flags]                  Send RCHV")
	fmt.Println()
	fmt.Println("Send flags:")
	fmt.Println("  --from-privkey <hex>    Private key of sender (hex encoded)")
	fmt.Println("  --to <address>          Recipient address")
	fmt.Println("  --amount <units>        Amount in base units (e.g., 10000000000 = 100 RCHV)")
	fmt.Println("  --fee <units>           Fee in base units (default: 100000)")
	fmt.Println("  --node <url>            Node RPC URL (default: http://localhost:8080)")
}

func cmdNew() {
	// Generate keypair
	privKey, pubKey, err := wallet.GenerateKeypair()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating keypair: %v\n", err)
		os.Exit(1)
	}

	// Derive address
	address, err := wallet.PubKeyToAddress(pubKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving address: %v\n", err)
		os.Exit(1)
	}

	// Print wallet info
	fmt.Println("🔐 New Archivas Wallet Generated")
	fmt.Println()
	fmt.Printf("Address:     %s\n", address)
	fmt.Printf("Public Key:  %s\n", hex.EncodeToString(pubKey))
	fmt.Printf("Private Key: %s\n", hex.EncodeToString(privKey))
	fmt.Println()
	fmt.Println("⚠️  KEEP YOUR PRIVATE KEY SECRET! Anyone with access can spend your RCHV.")
}

func cmdSend() {
	// Parse flags
	sendFlags := flag.NewFlagSet("send", flag.ExitOnError)
	fromPrivKeyHex := sendFlags.String("from-privkey", "", "Private key (hex)")
	to := sendFlags.String("to", "", "Recipient address")
	amountStr := sendFlags.String("amount", "", "Amount in base units")
	feeStr := sendFlags.String("fee", "100000", "Fee in base units")
	nodeURL := sendFlags.String("node", "http://localhost:8080", "Node RPC URL")

	sendFlags.Parse(os.Args[2:])

	// Validate required flags
	if *fromPrivKeyHex == "" || *to == "" || *amountStr == "" {
		fmt.Println("Error: --from-privkey, --to, and --amount are required")
		fmt.Println()
		printUsage()
		os.Exit(1)
	}

	// Parse private key
	privKeyBytes, err := hex.DecodeString(*fromPrivKeyHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid private key hex: %v\n", err)
		os.Exit(1)
	}

	// Derive public key and address from private key
	privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
	pubKey := privKey.PubKey()
	pubKeyBytes := pubKey.SerializeCompressed()

	fromAddr, err := wallet.PubKeyToAddress(pubKeyBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving address: %v\n", err)
		os.Exit(1)
	}

	// Parse amount and fee
	amount, err := strconv.ParseInt(*amountStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid amount: %v\n", err)
		os.Exit(1)
	}

	fee, err := strconv.ParseInt(*feeStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid fee: %v\n", err)
		os.Exit(1)
	}

	// Query current nonce from node
	balanceURL := fmt.Sprintf("%s/balance/%s", *nodeURL, fromAddr)
	resp, err := http.Get(balanceURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error querying balance: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error from node: %s\n", string(body))
		os.Exit(1)
	}

	var balanceResp struct {
		Address string `json:"address"`
		Balance int64  `json:"balance"`
		Nonce   uint64 `json:"nonce"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&balanceResp); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding balance response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📊 Current balance: %.8f RCHV (nonce: %d)\n", float64(balanceResp.Balance)/100000000.0, balanceResp.Nonce)

	// Build transaction
	tx := ledger.Transaction{
		From:         fromAddr,
		To:           *to,
		Amount:       amount,
		Fee:          fee,
		Nonce:        balanceResp.Nonce,
		SenderPubKey: pubKeyBytes,
	}

	// Sign transaction
	if err := wallet.SignTransaction(&tx, privKeyBytes); err != nil {
		fmt.Fprintf(os.Stderr, "Error signing transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📝 Sending %.8f RCHV from %s to %s\n", float64(amount)/100000000.0, fromAddr, *to)
	fmt.Printf("💸 Fee: %.8f RCHV\n", float64(fee)/100000000.0)

	// Submit transaction
	txJSON, err := json.Marshal(tx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding transaction: %v\n", err)
		os.Exit(1)
	}

	submitURL := fmt.Sprintf("%s/submitTx", *nodeURL)
	resp, err = http.Post(submitURL, "application/json", bytes.NewReader(txJSON))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error submitting transaction: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error from node: %s\n", string(body))
		os.Exit(1)
	}

	var submitResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &submitResp); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ %s: %s\n", submitResp.Status, submitResp.Message)
	fmt.Println("⏳ Transaction will be included in the next block (~20 seconds)")
}

