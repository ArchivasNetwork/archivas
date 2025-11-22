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
	"github.com/ArchivasNetwork/archivas/address"
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
	case "rotate-farmer-key":
		cmdRotateFarmerKey()
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
	fmt.Println("  archivas-wallet rotate-farmer-key [flags]     Safely rotate farmer keys")
	fmt.Println()
	fmt.Println("Send flags:")
	fmt.Println("  --from-privkey <hex>    Private key of sender (hex encoded)")
	fmt.Println("  --to <address>          Recipient address")
	fmt.Println("  --amount <units>        Amount in base units (e.g., 10000000000 = 100 RCHV)")
	fmt.Println("  --fee <units>           Fee in base units (default: 100000)")
	fmt.Println("  --node <url>            Node RPC URL (default: http://localhost:8080)")
	fmt.Println()
	fmt.Println("Rotate farmer key flags:")
	fmt.Println("  --old-privkey <hex>     Old farmer private key (optional, for balance transfer)")
	fmt.Println("  --node <url>            Node RPC URL (default: http://localhost:8080)")
	fmt.Println("  --broadcast             Actually broadcast the transfer transaction")
}

func cmdNew() {
	// Generate keypair
	privKey, pubKey, err := wallet.GenerateKeypair()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating keypair: %v\n", err)
		os.Exit(1)
	}

	// Derive addresses using UNIFIED Ethereum-compatible derivation
	evmAddr, err := address.PrivateKeyToEVMAddress(privKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving EVM address: %v\n", err)
		os.Exit(1)
	}

	arcvAddr, err := address.EncodeARCVAddress(evmAddr, "arcv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding ARCV address: %v\n", err)
		os.Exit(1)
	}

	// Print wallet info
	fmt.Println("🔐 New Ethereum-Compatible Archivas Wallet Generated")
	fmt.Println()
	fmt.Printf("ARCV Address (Bech32): %s\n", arcvAddr)
	fmt.Printf("EVM Address (Hex):     %s\n", evmAddr.Hex())
	fmt.Printf("Public Key:            %s\n", hex.EncodeToString(pubKey))
	fmt.Printf("Private Key:           %s\n", hex.EncodeToString(privKey))
	fmt.Println()
	fmt.Println("⚠️  KEEP YOUR PRIVATE KEY SECRET! Anyone with access can spend your RCHV.")
	fmt.Println()
	fmt.Println("✅ This wallet is fully compatible with MetaMask!")
	fmt.Println("   Import the private key to MetaMask to manage your RCHV.")
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

	// Derive address from private key using UNIFIED Ethereum-compatible derivation
	evmAddr, err := address.PrivateKeyToEVMAddress(privKeyBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving EVM address: %v\n", err)
		os.Exit(1)
	}

	fromAddr, err := address.EncodeARCVAddress(evmAddr, "arcv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding ARCV address: %v\n", err)
		os.Exit(1)
	}

	// Still need public key for transaction signature
	privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
	pubKey := privKey.PubKey()
	pubKeyBytes := pubKey.SerializeCompressed()

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

func cmdRotateFarmerKey() {
	// Parse flags
	rotateFlags := flag.NewFlagSet("rotate-farmer-key", flag.ExitOnError)
	oldPrivKeyHex := rotateFlags.String("old-privkey", "", "Old farmer private key (for balance transfer)")
	nodeURL := rotateFlags.String("node", "http://localhost:8080", "Node RPC URL")
	broadcast := rotateFlags.Bool("broadcast", false, "Actually broadcast the transfer transaction")

	rotateFlags.Parse(os.Args[2:])

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║          ARCHIVAS FARMER KEY ROTATION WIZARD               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Step 1: Generate new keypair
	fmt.Println("Step 1: Generating new farmer keypair...")
	newPrivKey, newPubKey, err := wallet.GenerateKeypair()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating keypair: %v\n", err)
		os.Exit(1)
	}

	newEVMAddr, err := address.PrivateKeyToEVMAddress(newPrivKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving EVM address: %v\n", err)
		os.Exit(1)
	}

	newARCVAddr, err := address.EncodeARCVAddress(newEVMAddr, "arcv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding ARCV address: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ New farmer keypair generated!")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("NEW FARMER CREDENTIALS (SAVE SECURELY):")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("ARCV Address: %s\n", newARCVAddr)
	fmt.Printf("EVM Address:  %s\n", newEVMAddr.Hex())
	fmt.Printf("Public Key:   %s\n", hex.EncodeToString(newPubKey))
	fmt.Printf("Private Key:  %s\n", hex.EncodeToString(newPrivKey))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Step 2: Output systemd service snippet
	fmt.Println("Step 2: Update your farmer service configuration")
	fmt.Println()
	fmt.Println("For systemd service, update the ExecStart line:")
	fmt.Println()
	fmt.Printf("ExecStart=/usr/local/bin/archivas-farmer farm \\\n")
	fmt.Printf("    --plots /var/lib/archivas/plots \\\n")
	fmt.Printf("    --node http://localhost:8545 \\\n")
	fmt.Printf("    --farmer-privkey %s\n", hex.EncodeToString(newPrivKey))
	fmt.Println()
	fmt.Println("Then restart the farmer:")
	fmt.Println("  sudo systemctl restart archivas-betanet-farmer")
	fmt.Println()

	// Step 3: Prepare balance transfer if old key provided
	if *oldPrivKeyHex != "" {
		fmt.Println("Step 3: Preparing balance transfer from old to new address...")
		fmt.Println()

		oldPrivKeyBytes, err := hex.DecodeString(*oldPrivKeyHex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid old private key hex: %v\n", err)
			os.Exit(1)
		}

		oldEVMAddr, err := address.PrivateKeyToEVMAddress(oldPrivKeyBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deriving old EVM address: %v\n", err)
			os.Exit(1)
		}

		oldARCVAddr, err := address.EncodeARCVAddress(oldEVMAddr, "arcv")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding old ARCV address: %v\n", err)
			os.Exit(1)
		}

		// Query old address balance
		balanceURL := fmt.Sprintf("%s/balance/%s", *nodeURL, oldARCVAddr)
		resp, err := http.Get(balanceURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error querying old address balance: %v\n", err)
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

		fmt.Printf("Old address: %s (EVM: %s)\n", oldARCVAddr, oldEVMAddr.Hex())
		fmt.Printf("Balance:     %.8f RCHV\n", float64(balanceResp.Balance)/100000000.0)
		fmt.Printf("Nonce:       %d\n", balanceResp.Nonce)
		fmt.Println()

		if balanceResp.Balance == 0 {
			fmt.Println("✅ Old address has zero balance. No transfer needed.")
			return
		}

		// Calculate transfer amount (balance minus fee)
		fee := int64(100000) // 0.001 RCHV
		if balanceResp.Balance <= fee {
			fmt.Println("⚠️  Balance too low to cover transfer fee. No transfer possible.")
			return
		}

		transferAmount := balanceResp.Balance - fee

		// Build transaction
		oldPrivKey := secp256k1.PrivKeyFromBytes(oldPrivKeyBytes)
		oldPubKey := oldPrivKey.PubKey()
		oldPubKeyBytes := oldPubKey.SerializeCompressed()

		tx := ledger.Transaction{
			From:         oldARCVAddr,
			To:           newARCVAddr,
			Amount:       transferAmount,
			Fee:          fee,
			Nonce:        balanceResp.Nonce,
			SenderPubKey: oldPubKeyBytes,
		}

		// Sign transaction
		if err := wallet.SignTransaction(&tx, oldPrivKeyBytes); err != nil {
			fmt.Fprintf(os.Stderr, "Error signing transaction: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Transfer prepared:\n")
		fmt.Printf("  From:   %s\n", oldARCVAddr)
		fmt.Printf("  To:     %s\n", newARCVAddr)
		fmt.Printf("  Amount: %.8f RCHV\n", float64(transferAmount)/100000000.0)
		fmt.Printf("  Fee:    %.8f RCHV\n", float64(fee)/100000000.0)
		fmt.Println()

		if *broadcast {
			fmt.Println("Broadcasting transfer transaction...")

			txJSON, err := json.Marshal(tx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding transaction: %v\n", err)
				os.Exit(1)
			}

			submitURL := fmt.Sprintf("%s/submitTx", *nodeURL)
			resp, err := http.Post(submitURL, "application/json", bytes.NewReader(txJSON))
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
		} else {
			fmt.Println("⚠️  Transfer NOT broadcast (use --broadcast flag to send)")
			fmt.Println()
			fmt.Println("To broadcast this transfer manually, run:")
			fmt.Printf("  archivas-wallet send --from-privkey %s --to %s --amount %d --fee %d\n",
				*oldPrivKeyHex, newARCVAddr, transferAmount, fee)
		}
	} else {
		fmt.Println("Step 3: No old private key provided, skipping balance transfer.")
		fmt.Println()
		fmt.Println("If you want to transfer RCHV from your old farmer address, run:")
		fmt.Println("  archivas-wallet rotate-farmer-key --old-privkey <OLD_KEY> --broadcast")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⚠️  SECURITY REMINDERS:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("1. Delete or securely destroy the old private key")
	fmt.Println("2. Never post private keys in logs, screenshots, or chat")
	fmt.Println("3. Store new private key in encrypted backup")
	fmt.Println("4. Verify new address can be imported to MetaMask")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

