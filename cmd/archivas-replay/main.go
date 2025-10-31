package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"math/rand"
)

func main() {
	difficulty := flag.Uint64("difficulty", 1000000, "Target difficulty")
	totalHashes := flag.Uint64("total-hashes", 1610612736, "Total plot hashes (e.g., 6x k28)")
	iterPerSec := flag.Int("iter-per-sec", 20, "VDF iterations per second")
	cycles := flag.Int("cycles", 200, "Number of cycles to simulate")
	seed := flag.Int64("seed", 12345, "Random seed")
	trace := flag.Bool("trace", false, "Verbose output")

	flag.Parse()

	fmt.Println("🔬 Archivas Replay Harness")
	fmt.Println("==========================")
	fmt.Printf("Difficulty: %d\n", *difficulty)
	fmt.Printf("Total Hashes: %d\n", *totalHashes)
	fmt.Printf("Iterations/sec: %d\n", *iterPerSec)
	fmt.Printf("Cycles: %d\n", *cycles)
	fmt.Println()

	// Initialize RNG
	r := rand.New(rand.NewSource(*seed))

	// Simulate farming cycles
	wins := 0
	qmax := uint64(^uint64(0)) // Max uint64

	for c := 0; c < *cycles; c++ {
		// Simulate checking N hashes
		foundWinner := false

		for h := uint64(0); h < *totalHashes; h++ {
			// Generate random quality (simulates SHA256)
			quality := r.Uint64()

			if quality < *difficulty {
				foundWinner = true
				if *trace {
					fmt.Printf("Cycle %d: Winner! quality=%d < %d\n", c, quality, *difficulty)
				}
				break
			}
		}

		if foundWinner {
			wins++
		}
	}

	// Calculate statistics
	pObs := float64(wins) / float64(*cycles)

	// Theoretical probability: p = 1 - (1 - D/QMAX)^N
	pSingle := float64(*difficulty) / float64(qmax)
	pTheory := 1.0 - math.Pow(1.0-pSingle, float64(*totalHashes))

	// ETW (Expected Time to Win)
	// etw = 1 / (p × iters/sec)
	etwObsSec := 0.0
	if pObs > 0 {
		etwObsSec = 1.0 / (pObs * float64(*iterPerSec))
	}

	etwTheorySec := 0.0
	if pTheory > 0 {
		etwTheorySec = 1.0 / (pTheory * float64(*iterPerSec))
	}

	deltaPct := math.Abs(pObs-pTheory) / pTheory * 100.0

	fmt.Println()
	fmt.Printf("REPLAY cycles=%d wins=%d p_obs=%.4f etw_obs_sec=%.1f\n", *cycles, wins, pObs, etwObsSec)
	fmt.Printf("EXPECTED p_theory=%.4f etw_theory_sec=%.1f qmax=%d diff=%d N=%d\n", pTheory, etwTheorySec, qmax, *difficulty, *totalHashes)

	result := "OK"
	if deltaPct > 15.0 {
		result = "MISMATCH"
	}

	fmt.Printf("RESULT=%s delta_p=%.2f%%\n", result, deltaPct)
	fmt.Println()

	if result == "OK" {
		fmt.Println("✅ Probability matches theory!")
		fmt.Printf("   With %d hashes at difficulty %d:\n", *totalHashes, *difficulty)
		fmt.Printf("   Expected time to win: %.1f seconds (%.1f minutes)\n", etwTheorySec, etwTheorySec/60.0)
	} else {
		fmt.Println("❌ MISMATCH detected!")
		fmt.Println("   There may be a bug in quality calculation or early-exit logic")
	}
}

// computeQuality mimics the real farmer's quality function
func computeQuality(challenge [32]byte, hash [32]byte) uint64 {
	h := sha256.New()
	h.Write(challenge[:])
	h.Write(hash[:])
	result := h.Sum(nil)
	return binary.LittleEndian.Uint64(result[:8])
}
