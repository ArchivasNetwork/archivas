package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/iljanemesis/archivas/ledger"
	"github.com/iljanemesis/archivas/wallet"
)

const (
	DripAmount = 2000000000 // 20 RCHV
	RateLimit  = time.Hour  // 1 drip per hour per IP
)

type Faucet struct {
	nodeURL   string
	privKey   []byte
	address   string
	lastDrips map[string]time.Time
	mu        sync.Mutex
}

func main() {
	nodeURL := flag.String("node", "http://localhost:8080", "Node RPC URL")
	privKeyHex := flag.String("faucet-privkey", "", "Faucet private key (hex)")
	port := flag.String("port", ":8081", "HTTP server port")
	flag.Parse()

	if *privKeyHex == "" {
		log.Fatal("--faucet-privkey required")
	}

	// Parse private key
	privKeyBytes, err := hex.DecodeString(*privKeyHex)
	if err != nil || len(privKeyBytes) != 32 {
		log.Fatalf("Invalid private key: must be 32 bytes hex")
	}

	// Derive address
	priv := secp256k1.PrivKeyFromBytes(privKeyBytes)
	pubKey := priv.PubKey().SerializeCompressed()
	address, err := wallet.PubKeyToAddress(pubKey)
	if err != nil {
		log.Fatalf("Failed to derive address: %v", err)
	}

	faucet := &Faucet{
		nodeURL:   *nodeURL,
		privKey:   privKeyBytes,
		address:   address,
		lastDrips: make(map[string]time.Time),
	}

	log.Printf("💧 Archivas Faucet Starting")
	log.Printf("   Node: %s", *nodeURL)
	log.Printf("   Faucet Address: %s", address)
	log.Printf("   Drip Amount: %.8f RCHV", float64(DripAmount)/100000000.0)
	log.Printf("   Rate Limit: 1 drip per %v per IP", RateLimit)
	log.Println()

	// Check faucet balance
	balance, err := faucet.getBalance()
	if err != nil {
		log.Printf("⚠️  Warning: Could not check faucet balance: %v", err)
	} else {
		log.Printf("💰 Faucet Balance: %.8f RCHV", float64(balance)/100000000.0)
	}

	// HTTP server
	http.HandleFunc("/", faucet.handleRoot)
	http.HandleFunc("/drip", faucet.handleDrip)
	http.HandleFunc("/health", faucet.handleHealth)

	log.Printf("🌐 Starting faucet server on %s", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func (f *Faucet) handleRoot(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Archivas Faucet</title>
    <style>
        body { font-family: monospace; max-width: 600px; margin: 50px auto; padding: 20px; }
        input { width: 100%; padding: 10px; margin: 10px 0; font-family: monospace; }
        button { padding: 10px 20px; background: #4CAF50; color: white; border: none; cursor: pointer; }
        button:hover { background: #45a049; }
        .info { background: #f0f0f0; padding: 10px; margin: 10px 0; }
    </style>
</head>
<body>
    <h1>💧 Archivas Faucet</h1>
    <p>Get free test RCHV for the Archivas testnet!</p>
    
    <div class="info">
        <p><strong>Drip Amount:</strong> 20 RCHV</p>
        <p><strong>Rate Limit:</strong> 1 drip per hour per IP</p>
    </div>

    <form action="/drip" method="GET">
        <input type="text" name="address" placeholder="Enter your Archivas address (arcv1...)" required />
        <button type="submit">💧 Request RCHV</button>
    </form>

    <p style="margin-top: 30px; font-size: 12px;">
        Faucet Address: ` + f.address + `<br/>
        <a href="/health">Health Status</a>
    </p>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func (f *Faucet) handleDrip(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	// Validate address format
	if len(address) < 10 || address[:4] != "arcv" {
		http.Error(w, "Invalid Archivas address format", http.StatusBadRequest)
		return
	}

	// Check rate limit
	clientIP := r.RemoteAddr
	f.mu.Lock()
	lastDrip, exists := f.lastDrips[clientIP]
	if exists && time.Since(lastDrip) < RateLimit {
		f.mu.Unlock()
		remaining := RateLimit - time.Since(lastDrip)
		http.Error(w, fmt.Sprintf("Rate limited. Try again in %v", remaining.Round(time.Minute)), http.StatusTooManyRequests)
		return
	}
	f.mu.Unlock()

	// Get current nonce
	nonce, err := f.getNonce()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get nonce: %v", err), http.StatusInternalServerError)
		return
	}

	// Create and sign transaction
	priv := secp256k1.PrivKeyFromBytes(f.privKey)
	pubKey := priv.PubKey().SerializeCompressed()

	tx := ledger.Transaction{
		From:         f.address,
		To:           address,
		Amount:       DripAmount,
		Fee:          100000, // 0.001 RCHV fee
		Nonce:        nonce,
		SenderPubKey: pubKey,
	}

	if err := wallet.SignTransaction(&tx, f.privKey); err != nil {
		http.Error(w, fmt.Sprintf("Failed to sign: %v", err), http.StatusInternalServerError)
		return
	}

	// Submit to node
	if err := f.submitTx(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Failed to submit: %v", err), http.StatusInternalServerError)
		return
	}

	// Update rate limit
	f.mu.Lock()
	f.lastDrips[clientIP] = time.Now()
	f.mu.Unlock()

	log.Printf("💧 Dripped 20 RCHV to %s (IP: %s)", address, clientIP)

	// Success response
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Amount  string `json:"amount"`
		TxHash  string `json:"txHash"`
	}{
		Success: true,
		Message: "RCHV sent! Should arrive in next block (~20 seconds)",
		Amount:  "20.00000000 RCHV",
		TxHash:  "pending",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (f *Faucet) handleHealth(w http.ResponseWriter, r *http.Request) {
	balance, _ := f.getBalance()

	health := struct {
		OK      bool   `json:"ok"`
		Address string `json:"address"`
		Balance int64  `json:"balance"`
		Drips   int    `json:"totalDrips"`
	}{
		OK:      balance > 0,
		Address: f.address,
		Balance: balance,
		Drips:   len(f.lastDrips),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (f *Faucet) getBalance() (int64, error) {
	url := fmt.Sprintf("%s/balance/%s", f.nodeURL, f.address)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Balance int64 `json:"balance"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Balance, nil
}

func (f *Faucet) getNonce() (uint64, error) {
	url := fmt.Sprintf("%s/balance/%s", f.nodeURL, f.address)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Nonce uint64 `json:"nonce"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Nonce, nil
}

func (f *Faucet) submitTx(tx *ledger.Transaction) error {
	data, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/submitTx", f.nodeURL)
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
