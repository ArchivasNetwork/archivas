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
	"path/filepath"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/iljanemesis/archivas/pospace"
	"github.com/iljanemesis/archivas/wallet"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "plot":
		cmdPlot()
	case "farm":
		cmdFarm()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Archivas Farmer")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  archivas-farmer plot [flags]   Generate a plot file")
	fmt.Println("  archivas-farmer farm [flags]   Start farming")
	fmt.Println()
	fmt.Println("Plot flags:")
	fmt.Println("  --path <dir>            Directory to store plot (default: ./plots)")
	fmt.Println("  --size <k>              Plot size parameter k (2^k hashes, default: 20)")
	fmt.Println("  --farmer-pubkey <hex>   Farmer public key (compressed, 33 bytes hex)")
	fmt.Println()
	fmt.Println("Farm flags:")
	fmt.Println("  --plots <dir>           Directory containing plots (default: ./plots)")
	fmt.Println("  --node <url>            Node RPC URL (default: http://localhost:8080)")
	fmt.Println("  --farmer-privkey <hex>  Farmer PRIVATE key (32 bytes hex) ⚠️ KEEP SECRET!")
}

func cmdPlot() {
	plotFlags := flag.NewFlagSet("plot", flag.ExitOnError)
	plotPath := plotFlags.String("path", "./plots", "Plot directory")
	kSize := plotFlags.Int("size", 20, "Plot size (k parameter)")
	farmerPubKeyHex := plotFlags.String("farmer-pubkey", "", "Farmer public key (compressed, 33 bytes hex)")

	plotFlags.Parse(os.Args[2:])

	// Validate parameters
	if *kSize < 10 || *kSize > 32 {
		fmt.Println("Error: k size must be between 10 and 32")
		os.Exit(1)
	}

	// Generate or parse farmer public key
	var farmerPubKey []byte
	if *farmerPubKeyHex == "" {
		fmt.Println("No --farmer-pubkey provided, generating new keypair...")
		privKey, pubKey, err := wallet.GenerateKeypair()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating keypair: %v\n", err)
			os.Exit(1)
		}
		farmerPubKey = pubKey

		addr, _ := wallet.PubKeyToAddress(pubKey)
		fmt.Printf("Generated new farmer identity:\n")
		fmt.Printf("  Address:     %s\n", addr)
		fmt.Printf("  Public Key:  %s (use for --farmer-pubkey)\n", hex.EncodeToString(pubKey))
		fmt.Printf("  Private Key: %s (use for --farmer-privkey) ⚠️ KEEP SECRET!\n", hex.EncodeToString(privKey))
		fmt.Println()
		fmt.Println("⚠️  Save both keys! You'll need the private key to farm.")
		fmt.Println()
	} else {
		var err error
		farmerPubKey, err = hex.DecodeString(*farmerPubKeyHex)
		if err != nil || len(farmerPubKey) != 33 {
			fmt.Println("Error: --farmer-pubkey must be 33 bytes (66 hex chars) compressed public key")
			os.Exit(1)
		}
	}

	// Create plot directory
	if err := os.MkdirAll(*plotPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating plot directory: %v\n", err)
		os.Exit(1)
	}

	// Generate plot filename
	plotFile := filepath.Join(*plotPath, fmt.Sprintf("plot-k%d.arcv", *kSize))

	fmt.Printf("🌾 Generating plot with k=%d (%d hashes)\n", *kSize, uint64(1)<<*kSize)
	fmt.Printf("📁 Output: %s\n", plotFile)
	fmt.Printf("👨‍🌾 Farmer: %s\n", hex.EncodeToString(farmerPubKey))
	fmt.Println()

	start := time.Now()
	if err := pospace.GeneratePlot(plotFile, uint32(*kSize), farmerPubKey); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating plot: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(start)
	fmt.Printf("\n✅ Plot generated successfully in %v\n", duration)
	fmt.Printf("📊 Plot size: ~%.2f MB\n", float64(uint64(1)<<*kSize*32)/(1024*1024))
}

func cmdFarm() {
	farmFlags := flag.NewFlagSet("farm", flag.ExitOnError)
	plotsDir := farmFlags.String("plots", "./plots", "Plots directory")
	nodeURL := farmFlags.String("node", "http://localhost:8080", "Node RPC URL")
	farmerPrivKeyHex := farmFlags.String("farmer-privkey", "", "Farmer PRIVATE key (32 bytes hex) ⚠️ KEEP SECRET!")

	farmFlags.Parse(os.Args[2:])

	if *farmerPrivKeyHex == "" {
		fmt.Println("Error: --farmer-privkey is required for farming")
		fmt.Println("Generate a wallet first: ./archivas-wallet new")
		os.Exit(1)
	}

	// Parse farmer private key
	privKeyBytes, err := hex.DecodeString(*farmerPrivKeyHex)
	if err != nil || len(privKeyBytes) != 32 {
		fmt.Println("Error: --farmer-privkey must be 32 bytes (64 hex chars)")
		fmt.Println("Get your private key from: ./archivas-wallet new")
		os.Exit(1)
	}

	// Derive farmer address
	farmerAddr, farmerPubKey, err := deriveFarmerAddress(privKeyBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deriving farmer address: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🌾 Archivas Farmer Starting")
	fmt.Printf("👨‍🌾 Farmer Address: %s\n", farmerAddr)
	fmt.Printf("📁 Plots Directory: %s\n", *plotsDir)
	fmt.Printf("🌐 Node: %s\n", *nodeURL)
	fmt.Println()

	// Load plots
	plots, err := loadPlots(*plotsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading plots: %v\n", err)
		os.Exit(1)
	}

	if len(plots) == 0 {
		fmt.Println("⚠️  No plots found! Generate plots first with: archivas-farmer plot")
		os.Exit(1)
	}

	fmt.Printf("✅ Loaded %d plot(s)\n", len(plots))
	for _, p := range plots {
		fmt.Printf("   - %s (k=%d, %d hashes)\n", filepath.Base(p.Path), p.Header.KSize, p.Header.NumHashes)
	}
	fmt.Println()
	fmt.Println("🚜 Starting farming loop...")
	fmt.Println()

	// Farming loop
	ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds
	defer ticker.Stop()

	var lastHeight uint64

	for range ticker.C {
		// Get current challenge from node
		challengeInfo, err := getChallenge(*nodeURL)
		if err != nil {
			fmt.Printf("⚠️  Error getting challenge: %v\n", err)
			continue
		}

		// Log when height changes
		if challengeInfo.Height != lastHeight {
			fmt.Printf("\n🔍 NEW HEIGHT %d (difficulty: %d)\n", challengeInfo.Height, challengeInfo.Difficulty)
			lastHeight = challengeInfo.Height
		}

		// Always check plots (VDF changes constantly = new challenge!)
		fmt.Printf("⚙️  Checking plots...")
		if challengeInfo.VDF != nil {
			fmt.Printf("   VDF: iter=%d\n", challengeInfo.VDF.Iterations)
		}

		// Check all plots
		var bestProof *pospace.Proof
		for _, plot := range plots {
			fmt.Printf("   Scanning plot %s...\n", filepath.Base(plot.Path))
			proof, err := plot.CheckChallenge(challengeInfo.Challenge, challengeInfo.Difficulty)
			if err != nil {
				fmt.Printf("⚠️  Error checking plot %s: %v\n", filepath.Base(plot.Path), err)
				continue
			}

			if proof != nil && (bestProof == nil || proof.Quality < bestProof.Quality) {
				bestProof = proof
			}
		}

		if bestProof != nil && bestProof.Quality < challengeInfo.Difficulty {
			fmt.Printf("🎉 Found winning proof! Quality: %d (target: %d)\n", bestProof.Quality, challengeInfo.Difficulty)
			
			// Submit block with VDF info
			if err := submitBlock(*nodeURL, bestProof, farmerAddr, farmerPubKey, privKeyBytes, challengeInfo); err != nil {
				fmt.Printf("❌ Error submitting block: %v\n", err)
			} else {
				vdfIter := uint64(0)
				if challengeInfo.VDF != nil {
					vdfIter = challengeInfo.VDF.Iterations
				}
				fmt.Printf("✅ Block submitted successfully for height %d (VDF t=%d)\n", challengeInfo.Height, vdfIter)
			}
		} else {
			bestQ := uint64(0)
			if bestProof != nil {
				bestQ = bestProof.Quality
			}
			fmt.Printf(" best=%d, need=<%d\n", bestQ, challengeInfo.Difficulty)
		}
	}
}

type ChallengeInfo struct {
	Challenge  [32]byte `json:"challenge"`
	Difficulty uint64   `json:"difficulty"`
	Height     uint64   `json:"height"`
	VDF        *struct {
		Seed       string `json:"seed"`       // hex-encoded
		Iterations uint64 `json:"iterations"`
		Output     string `json:"output"`     // hex-encoded
	} `json:"vdf,omitempty"`
}

func getChallenge(nodeURL string) (*ChallengeInfo, error) {
	resp, err := http.Get(nodeURL + "/challenge")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var info ChallengeInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode challenge: %w", err)
	}

	// Debug: verify challenge decoded
	if len(info.Challenge) == 0 {
		return nil, fmt.Errorf("challenge is empty after decode")
	}

	return &info, nil
}

func submitBlock(nodeURL string, proof *pospace.Proof, farmerAddr string, farmerPubKey []byte, privKey []byte, challenge *ChallengeInfo) error {
	// Create block submission (VDF fields only if available)
	submission := map[string]interface{}{
		"proof":        proof,
		"farmerAddr":   farmerAddr,
		"farmerPubKey": hex.EncodeToString(farmerPubKey),
	}
	
	// Add VDF info if present (for PoSpace+Time mode)
	if challenge.VDF != nil {
		vdfSeed, _ := hex.DecodeString(challenge.VDF.Seed)
		vdfOutput, _ := hex.DecodeString(challenge.VDF.Output)
		submission["vdfSeed"] = vdfSeed
		submission["vdfIterations"] = challenge.VDF.Iterations
		submission["vdfOutput"] = vdfOutput
	}

	data, err := json.Marshal(submission)
	if err != nil {
		return err
	}

	resp, err := http.Post(nodeURL+"/submitBlock", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func loadPlots(dir string) ([]*pospace.PlotFile, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var plots []*pospace.PlotFile
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != ".arcv" {
			continue
		}

		plotPath := filepath.Join(dir, f.Name())
		plot, err := pospace.OpenPlot(plotPath)
		if err != nil {
			fmt.Printf("⚠️  Skipping %s: %v\n", f.Name(), err)
			continue
		}

		plots = append(plots, plot)
	}

	return plots, nil
}

func deriveFarmerAddress(privKey []byte) (string, []byte, error) {
	// Import secp256k1 to derive public key from private key
	// We need to add this import at the top
	privKeyObj := secp256k1.PrivKeyFromBytes(privKey)
	pubKey := privKeyObj.PubKey()
	pubKeyBytes := pubKey.SerializeCompressed()

	addr, err := wallet.PubKeyToAddress(pubKeyBytes)
	if err != nil {
		return "", nil, err
	}

	return addr, pubKeyBytes, nil
}
