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
)

type Explorer struct {
	nodeURL string
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
	nodeURL := flag.String("node", "http://localhost:8080", "Node RPC URL")
	port := flag.String("port", ":8082", "Explorer HTTP port")
	flag.Parse()

	explorer := &Explorer{
		nodeURL: *nodeURL,
	}

	log.Printf("🔍 Archivas Block Explorer Starting")
	log.Printf("   Node: %s", *nodeURL)
	log.Printf("   Port: %s", *port)
	log.Println()

	http.HandleFunc("/", explorer.handleHome)
	http.HandleFunc("/block/", explorer.handleBlock)
	http.HandleFunc("/address/", explorer.handleAddress)
	http.HandleFunc("/mempool", explorer.handleMempool)
	http.HandleFunc("/peers", explorer.handlePeersPage)
	http.HandleFunc("/tx/", explorer.handleTransaction)

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
        <a href="/mempool">💧 Mempool</a>
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

func (e *Explorer) getChainTip() (*ChainInfo, error) {
	resp, err := http.Get(e.nodeURL + "/chainTip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		BlockHash  []byte `json:"blockHash"`
		Height     uint64 `json:"height"`
		Difficulty uint64 `json:"difficulty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &ChainInfo{
		Height:     result.Height,
		Difficulty: result.Difficulty,
		BlockHash:  hex.EncodeToString(result.BlockHash),
	}, nil
}

func (e *Explorer) getBalance(address string) (int64, uint64, error) {
	url := fmt.Sprintf("%s/balance/%s", e.nodeURL, address)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Balance int64  `json:"balance"`
		Nonce   uint64 `json:"nonce"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	return result.Balance, result.Nonce, nil
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
	resp, err := http.Get(e.nodeURL + "/peers")
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
