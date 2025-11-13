package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Explorer struct {
	rpcURLs       []string
	currentRPC    int
	rpcMutex      sync.RWMutex
	rpcStatus     map[string]*RPCStatus
	requestTimeout time.Duration
	maxRetries    int
}

type RPCStatus struct {
	URL         string
	Healthy     bool
	LastCheck   time.Time
	Failures    int
	LastError   string
	mutex       sync.RWMutex
}

type ChainInfo struct {
	Height     uint64
	Difficulty uint64
	BlockHash  string
}

type BlockInfo struct {
	Height       uint64
	Hash         string
	Difficulty   uint64
	Timestamp    int64
	FarmerAddr   string
	Transactions int
}

func main() {
	nodeURL := flag.String("node", "http://localhost:8080", "Node RPC URL (can be comma-separated list)")
	port := flag.String("port", ":8082", "Explorer HTTP port")
	timeout := flag.Int("timeout", 3, "RPC request timeout in seconds")
	retries := flag.Int("retries", 2, "Max retries per request")
	flag.Parse()

	// Parse multiple RPC URLs (comma-separated)
	rpcURLs := strings.Split(*nodeURL, ",")
	for i := range rpcURLs {
		rpcURLs[i] = strings.TrimSpace(rpcURLs[i])
	}

	// Initialize RPC status tracking
	rpcStatus := make(map[string]*RPCStatus)
	for _, url := range rpcURLs {
		rpcStatus[url] = &RPCStatus{
			URL:       url,
			Healthy:   true,
			LastCheck: time.Now(),
		}
	}

	explorer := &Explorer{
		rpcURLs:        rpcURLs,
		currentRPC:     0,
		rpcStatus:      rpcStatus,
		requestTimeout: time.Duration(*timeout) * time.Second,
		maxRetries:     *retries,
	}

	// Start health check goroutine
	go explorer.healthCheckLoop()

	log.Printf("🔍 Archivas Block Explorer Starting")
	log.Printf("   RPC Endpoints: %v", rpcURLs)
	log.Printf("   Port: %s", *port)
	log.Printf("   Timeout: %ds | Retries: %d", *timeout, *retries)
	log.Println()

	http.HandleFunc("/", explorer.handleHome)
	http.HandleFunc("/block/", explorer.handleBlock)
	http.HandleFunc("/address/", explorer.handleAddress)
	http.HandleFunc("/mempool", explorer.handleMempool)
	http.HandleFunc("/peers", explorer.handlePeersPage)
	http.HandleFunc("/tx/", explorer.handleTransaction)
	http.HandleFunc("/rpc/status", explorer.handleRPCStatus)

	log.Printf("🌐 Explorer running on %s", *port)
	log.Fatal(http.ListenAndServe(*port, nil))
}

func (e *Explorer) handleHome(w http.ResponseWriter, r *http.Request) {
	// Get chain info
	chainInfo, err := e.getChainTip()
	if err != nil {
		http.Error(w, "Failed to get chain info", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Archivas Block Explorer</title>
    <style>
        body { font-family: 'Courier New', monospace; max-width: 1000px; margin: 0 auto; padding: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; text-align: center; }
        .stats { display: flex; justify-content: space-around; margin: 30px 0; }
        .stat { background: #2a2a2a; padding: 20px; border: 1px solid #00ff00; text-align: center; }
        .stat-value { font-size: 24px; font-weight: bold; }
        .stat-label { font-size: 12px; color: #888; }
        a { color: #00ff00; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>🌾 Archivas Block Explorer</h1>
    
    <div class="stats">
        <div class="stat">
            <div class="stat-value">{{.Height}}</div>
            <div class="stat-label">BLOCK HEIGHT</div>
        </div>
        <div class="stat">
            <div class="stat-value">{{.Difficulty}}</div>
            <div class="stat-label">DIFFICULTY</div>
        </div>
        <div class="stat">
            <div class="stat-value">{{printf "%.8s..." .BlockHash}}</div>
            <div class="stat-label">TIP HASH</div>
        </div>
    </div>

    <h2>Recent Blocks</h2>
    <p style="text-align: center;">
        <a href="/block/{{.Height}}">Block {{.Height}}</a> (latest)
    </p>

    <h2>Search</h2>
    <form action="/block/" method="GET" style="text-align: center;">
        <input type="number" name="height" placeholder="Block height" style="padding: 10px; background: #2a2a2a; color: #00ff00; border: 1px solid #00ff00;" />
        <button type="submit" style="padding: 10px 20px; background: #00ff00; color: #000; border: none; cursor: pointer;">Search Block</button>
    </form>

    <form action="/address/" method="GET" style="text-align: center; margin-top: 10px;">
        <input type="text" name="addr" placeholder="Address (arcv1...)" style="padding: 10px; background: #2a2a2a; color: #00ff00; border: 1px solid #00ff00; width: 400px;" />
        <button type="submit" style="padding: 10px 20px; background: #00ff00; color: #000; border: none; cursor: pointer;">Check Balance</button>
    </form>

    <h2>Network</h2>
    <p style="text-align: center;">
        <a href="/peers">🌐 View Peers</a> | 
        <a href="/mempool">💧 Mempool</a> |
        <a href="/search">🔍 Search</a>
    </p>
    
    <h2>Try the New APIs</h2>
    <p style="text-align: center; font-size: 12px; color: #888;">
        v0.6.0 Features: HD Wallets, Account History, Smart Search
    </p>

    <p style="text-align: center; margin-top: 50px; font-size: 12px; color: #666;">
        Archivas Testnet v0.3.0 - Proof-of-Space-and-Time<br/>
        <a href="https://github.com/ArchivasNetwork/archivas">GitHub</a>
    </p>
</body>
</html>
`))

	tmpl.Execute(w, chainInfo)
}

func (e *Explorer) handleBlock(w http.ResponseWriter, r *http.Request) {
	// Parse height from URL
	heightStr := r.URL.Path[len("/block/"):]
	if r.URL.Query().Get("height") != "" {
		heightStr = r.URL.Query().Get("height")
	}

	height, err := strconv.ParseUint(heightStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid block height", http.StatusBadRequest)
		return
	}

	// For now, just show basic info (full block endpoint would need RPC addition)
	tmpl := template.Must(template.New("block").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Block {{.}} - Archivas Explorer</title>
    <style>
        body { font-family: 'Courier New', monospace; max-width: 1000px; margin: 0 auto; padding: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; }
        a { color: #00ff00; }
    </style>
</head>
<body>
    <h1>Block {{.}}</h1>
    <p><a href="/">← Back to Explorer</a></p>
    
    <p>Block details endpoint coming soon!</p>
    <p>For now, check via RPC: <code>curl http://node:8080/chainTip</code></p>
</body>
</html>
`))

	tmpl.Execute(w, height)
}

func (e *Explorer) handleAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Path[len("/address/"):]
	if r.URL.Query().Get("addr") != "" {
		address = r.URL.Query().Get("addr")
	}

	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	// Get balance
	balance, nonce, err := e.getBalance(address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get balance: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.New("address").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>{{.Address}} - Archivas Explorer</title>
    <style>
        body { font-family: 'Courier New', monospace; max-width: 1000px; margin: 0 auto; padding: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; word-break: break-all; }
        a { color: #00ff00; }
        .stat { background: #2a2a2a; padding: 15px; margin: 10px 0; border: 1px solid #00ff00; }
    </style>
</head>
<body>
    <h1>Address</h1>
    <p style="word-break: break-all;">{{.Address}}</p>
    <p><a href="/">← Back to Explorer</a></p>
    
    <div class="stat">
        <strong>Balance:</strong> {{printf "%.8f" .Balance}} RCHV
    </div>
    
    <div class="stat">
        <strong>Nonce:</strong> {{.Nonce}}
    </div>
</body>
</html>
`))

	data := struct {
		Address string
		Balance float64
		Nonce   uint64
	}{
		Address: address,
		Balance: float64(balance) / 100000000.0,
		Nonce:   nonce,
	}

	tmpl.Execute(w, data)
}

// fetchWithRetry attempts to fetch from RPC endpoints with automatic failover
func (e *Explorer) fetchWithRetry(path string) (*http.Response, string, error) {
	client := &http.Client{
		Timeout: e.requestTimeout,
	}

	var lastErr error
	triedURLs := make(map[string]bool)

	for attempt := 0; attempt <= e.maxRetries; attempt++ {
		// Get next healthy RPC endpoint
		rpcURL := e.getNextHealthyRPC()
		if rpcURL == "" {
			return nil, "", fmt.Errorf("no healthy RPC endpoints available")
		}

		// Skip if already tried this URL in this request
		if triedURLs[rpcURL] {
			continue
		}
		triedURLs[rpcURL] = true

		fullURL := rpcURL + path
		resp, err := client.Get(fullURL)

		if err == nil && resp.StatusCode < 500 {
			// Success! Mark RPC as healthy
			e.markRPCHealthy(rpcURL)
			return resp, rpcURL, nil
		}

		// Failed - mark as unhealthy and try next
		if err != nil {
			lastErr = err
			e.markRPCUnhealthy(rpcURL, err.Error())
		} else {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			e.markRPCUnhealthy(rpcURL, lastErr.Error())
		}

		log.Printf("⚠️  RPC %s failed (attempt %d/%d): %v", rpcURL, attempt+1, e.maxRetries+1, lastErr)
	}

	return nil, "", fmt.Errorf("all RPC endpoints failed, last error: %v", lastErr)
}

// getNextHealthyRPC returns the next healthy RPC URL using round-robin
func (e *Explorer) getNextHealthyRPC() string {
	e.rpcMutex.Lock()
	defer e.rpcMutex.Unlock()

	// Try to find a healthy RPC starting from current position
	for i := 0; i < len(e.rpcURLs); i++ {
		idx := (e.currentRPC + i) % len(e.rpcURLs)
		url := e.rpcURLs[idx]

		if status, ok := e.rpcStatus[url]; ok {
			status.mutex.RLock()
			healthy := status.Healthy || time.Since(status.LastCheck) > 10*time.Second
			status.mutex.RUnlock()

			if healthy {
				e.currentRPC = (idx + 1) % len(e.rpcURLs)
				return url
			}
		}
	}

	// No healthy RPCs, return first one anyway (might have recovered)
	if len(e.rpcURLs) > 0 {
		return e.rpcURLs[0]
	}

	return ""
}

// markRPCHealthy marks an RPC endpoint as healthy
func (e *Explorer) markRPCHealthy(url string) {
	if status, ok := e.rpcStatus[url]; ok {
		status.mutex.Lock()
		status.Healthy = true
		status.Failures = 0
		status.LastCheck = time.Now()
		status.LastError = ""
		status.mutex.Unlock()
	}
}

// markRPCUnhealthy marks an RPC endpoint as unhealthy
func (e *Explorer) markRPCUnhealthy(url string, errMsg string) {
	if status, ok := e.rpcStatus[url]; ok {
		status.mutex.Lock()
		status.Healthy = false
		status.Failures++
		status.LastCheck = time.Now()
		status.LastError = errMsg
		status.mutex.Unlock()
	}
}

// healthCheckLoop periodically checks RPC endpoints
func (e *Explorer) healthCheckLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, url := range e.rpcURLs {
			go e.checkRPCHealth(url)
		}
	}
}

// checkRPCHealth performs a health check on a single RPC endpoint
func (e *Explorer) checkRPCHealth(url string) {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(url + "/chainTip")
	if err != nil {
		e.markRPCUnhealthy(url, err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		e.markRPCHealthy(url)
	} else {
		e.markRPCUnhealthy(url, fmt.Sprintf("HTTP %d", resp.StatusCode))
	}
}

func (e *Explorer) getChainTip() (*ChainInfo, error) {
	resp, rpcURL, err := e.fetchWithRetry("/chainTip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_ = rpcURL // Used RPC endpoint (could display in UI)

	// v1.1.0: chainTip returns strings instead of numbers
	var result struct {
		Hash       string `json:"hash"`       // hex string
		Height     string `json:"height"`     // u64 as string
		Difficulty string `json:"difficulty"` // u64 as string
		// Legacy fields for backward compatibility
		BlockHash     []byte `json:"blockHash,omitempty"`
		HeightNum     uint64 `json:"heightNum,omitempty"`
		DifficultyNum uint64 `json:"difficultyNum,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Parse height and difficulty from strings
	var height uint64
	var difficulty uint64
	var blockHash string

	if result.Height != "" {
		// v1.1.0 format: parse from string
		var err error
		height, err = strconv.ParseUint(result.Height, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid height: %v", err)
		}
	} else if result.HeightNum != 0 {
		// Legacy format
		height = result.HeightNum
	} else {
		return nil, fmt.Errorf("missing height")
	}

	if result.Difficulty != "" {
		// v1.1.0 format: parse from string
		var err error
		difficulty, err = strconv.ParseUint(result.Difficulty, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid difficulty: %v", err)
		}
	} else if result.DifficultyNum != 0 {
		// Legacy format
		difficulty = result.DifficultyNum
	} else {
		return nil, fmt.Errorf("missing difficulty")
	}

	if result.Hash != "" {
		// v1.1.0 format: hash is already hex string
		blockHash = result.Hash
	} else if len(result.BlockHash) > 0 {
		// Legacy format: encode bytes to hex
		blockHash = hex.EncodeToString(result.BlockHash)
	} else {
		return nil, fmt.Errorf("missing block hash")
	}

	return &ChainInfo{
		Height:     height,
		Difficulty: difficulty,
		BlockHash:  blockHash,
	}, nil
}

func (e *Explorer) getBalance(address string) (int64, uint64, error) {
	path := fmt.Sprintf("/account/%s", address)
	resp, _, err := e.fetchWithRetry(path)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// v1.1.0: /account returns strings, but try /balance first for backward compatibility
	var result struct {
		Balance string `json:"balance"` // u64 as string (v1.1.0)
		Nonce   string `json:"nonce"`   // u64 as string (v1.1.0)
		// Legacy fields
		BalanceNum int64  `json:"balanceNum,omitempty"`
		NonceNum   uint64 `json:"nonceNum,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	// Parse balance and nonce from strings
	var balance int64
	var nonce uint64

	if result.Balance != "" {
		// v1.1.0 format: parse from string
		var err error
		balance, err = strconv.ParseInt(result.Balance, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid balance: %v", err)
		}
	} else if result.BalanceNum != 0 {
		// Legacy format
		balance = result.BalanceNum
	}

	if result.Nonce != "" {
		// v1.1.0 format: parse from string
		var err error
		nonce, err = strconv.ParseUint(result.Nonce, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid nonce: %v", err)
		}
	} else if result.NonceNum != 0 {
		// Legacy format
		nonce = result.NonceNum
	}

	return balance, nonce, nil
}

func (e *Explorer) handleMempool(w http.ResponseWriter, r *http.Request) {
	// TODO: Add mempool RPC endpoint to node
	tmpl := template.Must(template.New("mempool").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Mempool - Archivas Explorer</title>
    <style>
        body { font-family: 'Courier New', monospace; max-width: 1000px; margin: 0 auto; padding: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; }
        a { color: #00ff00; }
        .tx { background: #2a2a2a; padding: 10px; margin: 10px 0; border: 1px solid #00ff00; }
    </style>
</head>
<body>
    <h1>💧 Mempool</h1>
    <p><a href="/">← Back to Explorer</a></p>
    
    <p>Pending transactions will appear here.</p>
    <p><em>Note: Mempool RPC endpoint coming soon!</em></p>
</body>
</html>
`))

	tmpl.Execute(w, nil)
}

func (e *Explorer) handlePeersPage(w http.ResponseWriter, r *http.Request) {
	resp, _, err := e.fetchWithRetry("/peers")
	if err != nil {
		http.Error(w, "Failed to get peers", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var peersData struct {
		Connected []string `json:"connected"`
		Known     []string `json:"known"`
	}
	json.NewDecoder(resp.Body).Decode(&peersData)

	tmpl := template.Must(template.New("peers").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Peers - Archivas Explorer</title>
    <style>
        body { font-family: 'Courier New', monospace; max-width: 1000px; margin: 0 auto; padding: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; }
        a { color: #00ff00; }
        .peer-section { margin: 20px 0; }
        .peer { background: #2a2a2a; padding: 10px; margin: 5px 0; border: 1px solid #00ff00; }
        .connected { border-color: #00ff00; }
        .known { border-color: #888; }
    </style>
</head>
<body>
    <h1>🌐 Network Peers</h1>
    <p><a href="/">← Back to Explorer</a></p>
    
    <div class="peer-section">
        <h2>Connected Peers ({{len .Connected}})</h2>
        {{range .Connected}}
        <div class="peer connected">✅ {{.}}</div>
        {{else}}
        <p><em>No connected peers</em></p>
        {{end}}
    </div>

    <div class="peer-section">
        <h2>Known Peers ({{len .Known}})</h2>
        {{range .Known}}
        <div class="peer known">📋 {{.}}</div>
        {{else}}
        <p><em>No known peers</em></p>
        {{end}}
    </div>
</body>
</html>
`))

	tmpl.Execute(w, peersData)
}

func (e *Explorer) handleTransaction(w http.ResponseWriter, r *http.Request) {
	txHash := r.URL.Path[len("/tx/"):]

	tmpl := template.Must(template.New("tx").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Transaction {{.}} - Archivas Explorer</title>
    <style>
        body { font-family: 'Courier New', monospace; max-width: 1000px; margin: 0 auto; padding: 20px; background: #1a1a1a; color: #00ff00; }
        h1 { color: #00ff00; word-break: break-all; }
        a { color: #00ff00; }
    </style>
</head>
<body>
    <h1>Transaction</h1>
    <p style="word-break: break-all;">{{.}}</p>
    <p><a href="/">← Back to Explorer</a></p>
    
    <p>Transaction details endpoint coming soon!</p>
    <p><em>For now, transactions are included in block data.</em></p>
</body>
</html>
`))

	tmpl.Execute(w, txHash)
}

func (e *Explorer) handleRPCStatus(w http.ResponseWriter, r *http.Request) {
	e.rpcMutex.RLock()
	defer e.rpcMutex.RUnlock()

	type StatusResponse struct {
		RPC         string    `json:"rpc"`
		Healthy     bool      `json:"healthy"`
		Failures    int       `json:"failures"`
		LastCheck   time.Time `json:"last_check"`
		LastError   string    `json:"last_error,omitempty"`
	}

	var statuses []StatusResponse
	for _, url := range e.rpcURLs {
		if status, ok := e.rpcStatus[url]; ok {
			status.mutex.RLock()
			statuses = append(statuses, StatusResponse{
				RPC:       status.URL,
				Healthy:   status.Healthy,
				Failures:  status.Failures,
				LastCheck: status.LastCheck,
				LastError: status.LastError,
			})
			status.mutex.RUnlock()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rpcs":        statuses,
		"total":       len(statuses),
		"healthy":     countHealthy(statuses),
		"current_rpc": e.rpcURLs[e.currentRPC%len(e.rpcURLs)],
	})
}

func countHealthy(statuses []StatusResponse) int {
	count := 0
	for _, s := range statuses {
		if s.Healthy {
			count++
		}
	}
	return count
}

type StatusResponse struct {
	RPC         string    `json:"rpc"`
	Healthy     bool      `json:"healthy"`
	Failures    int       `json:"failures"`
	LastCheck   time.Time `json:"last_check"`
	LastError   string    `json:"last_error,omitempty"`
}
