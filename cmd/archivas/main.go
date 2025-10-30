package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/wallet"
	"golang.org/x/term"
)

const (
	DefaultKeystorePath = ".archivas/keystore.json"
	DefaultNodeURL      = "http://localhost:8080"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "wallet":
		handleWallet()
	case "tx":
		handleTx()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Archivas CLI v0.6.0")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  archivas wallet init                  Create new keystore")
	fmt.Println("  archivas wallet new <name>            Create new account")
	fmt.Println("  archivas wallet list                  List accounts")
	fmt.Println("  archivas wallet show <address>        Show account details")
	fmt.Println("  archivas wallet import                Import from mnemonic")
	fmt.Println()
	fmt.Println("  archivas tx send [flags]              Send transaction")
	fmt.Println("  archivas tx history <address>         Show tx history")
	fmt.Println()
	fmt.Println("Flags for 'tx send':")
	fmt.Println("  --from <addr>    Sender address")
	fmt.Println("  --to <addr>      Recipient address")
	fmt.Println("  --amount <units> Amount (base units)")
	fmt.Println("  --fee <units>    Fee (default: 100000)")
	fmt.Println("  --node <url>     Node URL (default: http://localhost:8080)")
}

func handleWallet() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: archivas wallet <command>")
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "init":
		cmdWalletInit()
	case "new":
		cmdWalletNew()
	case "list":
		cmdWalletList()
	case "show":
		cmdWalletShow()
	case "import":
		cmdWalletImport()
	default:
		fmt.Printf("Unknown wallet command: %s\n", subcommand)
		os.Exit(1)
	}
}

func cmdWalletInit() {
	keystorePath := getKeystorePath()

	// Check if keystore already exists
	if _, err := os.Stat(keystorePath); err == nil {
		fmt.Printf("Keystore already exists at %s\n", keystorePath)
		fmt.Println("Use 'archivas wallet import' to import from mnemonic")
		os.Exit(1)
	}

	// Generate mnemonic
	mnemonic, err := wallet.GenerateMnemonic()
	if err != nil {
		fmt.Printf("Error generating mnemonic: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔐 New Wallet Created!")
	fmt.Println()
	fmt.Println("⚠️  WRITE DOWN YOUR RECOVERY PHRASE:")
	fmt.Println()
	fmt.Println(mnemonic)
	fmt.Println()
	fmt.Println("⚠️  Store this securely! You'll need it to recover your wallet.")
	fmt.Println()

	// Get password
	fmt.Print("Enter password to encrypt keystore: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}

	// Create master seed
	seed, err := wallet.MnemonicToSeed(mnemonic, "")
	if err != nil {
		fmt.Printf("Error creating seed: %v\n", err)
		os.Exit(1)
	}

	// Create keystore
	ks := wallet.NewKeystore()
	if err := ks.Encrypt(seed, string(password)); err != nil {
		fmt.Printf("Error encrypting keystore: %v\n", err)
		os.Exit(1)
	}

	// Derive first account
	privKey, err := wallet.DeriveKey(seed, 0, 0)
	if err != nil {
		fmt.Printf("Error deriving key: %v\n", err)
		os.Exit(1)
	}

	priv := secp256k1.PrivKeyFromBytes(privKey)
	pubKey := priv.PubKey().SerializeCompressed()
	address, err := wallet.PubKeyToAddress(pubKey)
	if err != nil {
		fmt.Printf("Error deriving address: %v\n", err)
		os.Exit(1)
	}

	ks.AddAccount("primary", address, wallet.DerivationPath(0, 0))

	// Save keystore
	if err := os.MkdirAll(filepath.Dir(keystorePath), 0700); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	if err := ks.Save(keystorePath); err != nil {
		fmt.Printf("Error saving keystore: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Keystore created: %s\n", keystorePath)
	fmt.Printf("✅ First account: %s\n", address)
	fmt.Println()
	fmt.Println("Use 'archivas wallet new <name>' to create more accounts")
}

func cmdWalletNew() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: archivas wallet new <name>")
		os.Exit(1)
	}

	name := os.Args[3]
	fmt.Printf("Creating account: %s\n", name)
	fmt.Println("(Full implementation in next commit)")
}

func cmdWalletList() {
	keystorePath := getKeystorePath()

	ks, err := wallet.LoadKeystore(keystorePath)
	if err != nil {
		fmt.Printf("Error loading keystore: %v\n", err)
		fmt.Println("Run 'archivas wallet init' to create a new wallet")
		os.Exit(1)
	}

	fmt.Println("💼 Accounts:")
	fmt.Println()

	if len(ks.Accounts) == 0 {
		fmt.Println("No accounts. Use 'archivas wallet new <name>' to create one.")
		return
	}

	for i, acc := range ks.Accounts {
		fmt.Printf("%d. %s\n", i+1, acc.Name)
		fmt.Printf("   Address: %s\n", acc.Address)
		fmt.Printf("   Path: %s\n", acc.Path)
		fmt.Println()
	}
}

func cmdWalletShow() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: archivas wallet show <address>")
		os.Exit(1)
	}

	address := os.Args[3]
	nodeURL := DefaultNodeURL

	// Get balance
	resp, err := http.Get(fmt.Sprintf("%s/balance/%s", nodeURL, address))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result struct {
		Balance int64  `json:"balance"`
		Nonce   uint64 `json:"nonce"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Balance: %.8f RCHV\n", float64(result.Balance)/100000000.0)
	fmt.Printf("Nonce: %d\n", result.Nonce)
}

func cmdWalletImport() {
	fmt.Println("Enter your 24-word recovery phrase:")
	reader := bufio.NewReader(os.Stdin)
	mnemonic, _ := reader.ReadString('\n')
	mnemonic = strings.TrimSpace(mnemonic)

	// Validate
	seed, err := wallet.MnemonicToSeed(mnemonic, "")
	if err != nil {
		fmt.Printf("Invalid mnemonic: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Valid mnemonic")
	fmt.Println("(Full import implementation in next commit)")
	_ = seed
}

func handleTx() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: archivas tx <command>")
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "send":
		cmdTxSend()
	case "history":
		cmdTxHistory()
	default:
		fmt.Printf("Unknown tx command: %s\n", subcommand)
		os.Exit(1)
	}
}

func cmdTxSend() {
	sendFlags := flag.NewFlagSet("send", flag.ExitOnError)
	from := sendFlags.String("from", "", "Sender address")
	to := sendFlags.String("to", "", "Recipient address")
	amountStr := sendFlags.String("amount", "", "Amount in base units")
	feeStr := sendFlags.String("fee", "100000", "Fee in base units")
	nodeURL := sendFlags.String("node", DefaultNodeURL, "Node URL")

	sendFlags.Parse(os.Args[3:])

	if *from == "" || *to == "" || *amountStr == "" {
		fmt.Println("Error: --from, --to, and --amount required")
		os.Exit(1)
	}

	fmt.Printf("Sending from %s to %s\n", *from, *to)
	fmt.Println("(Full send implementation in next commit)")
	_, _, _ = feeStr, nodeURL, amountStr
}

func cmdTxHistory() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: archivas tx history <address>")
		os.Exit(1)
	}

	address := os.Args[3]
	fmt.Printf("Transaction history for %s\n", address)
	fmt.Println("(Will use /account/:addr/txs API)")
}

func getKeystorePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, DefaultKeystorePath)
}

