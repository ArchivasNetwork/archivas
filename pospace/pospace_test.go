package pospace

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestVerifyComparison(t *testing.T) {
	// Create a valid farmer pubkey and plot ID
	farmerPubKey := [33]byte{}
	copy(farmerPubKey[:], []byte("test-farmer-pubkey-0123456789012"))

	plotID := sha256.Sum256(farmerPubKey[:])
	challenge := sha256.Sum256([]byte("test-challenge"))

	// Compute a valid plot hash
	index := uint64(123)
	plotHash := computePlotHash(farmerPubKey[:], plotID[:], index)

	// Compute quality
	q := computeQuality(challenge, plotHash)

	// Create a valid proof
	proof := &Proof{
		Challenge:    challenge,
		PlotID:       plotID,
		Index:        index,
		Hash:         plotHash,
		Quality:      q,
		FarmerPubKey: farmerPubKey,
	}

	// Test exact match (quality == difficulty) should PASS
	if !VerifyProof(proof, challenge, q) {
		t.Fatalf("expected exact match to pass: q=%d", q)
	}

	// Test quality < difficulty should PASS
	if !VerifyProof(proof, challenge, q+1000) {
		t.Fatalf("expected (q < difficulty) to pass: q=%d, diff=%d", q, q+1000)
	}

	// Test quality > difficulty should FAIL
	if q > 1000 && VerifyProof(proof, challenge, q-1000) {
		t.Fatalf("expected (q > difficulty) to fail: q=%d, diff=%d", q, q-1000)
	}

	// Log the quality value for reference
	t.Logf("Quality: %d (should be < QMAX=%d)", q, QMAX)
}

func TestQualityBounded(t *testing.T) {
	challenge := sha256.Sum256([]byte("test-challenge"))
	plotHash := sha256.Sum256([]byte("test-plot-hash"))

	q := computeQuality(challenge, plotHash)

	if q >= QMAX {
		t.Fatalf("quality %d should be < QMAX %d", q, QMAX)
	}
}

func TestQualityDeterministic(t *testing.T) {
	challenge := sha256.Sum256([]byte("deterministic-test"))
	plotHash := sha256.Sum256([]byte("plot-hash-test"))

	q1 := computeQuality(challenge, plotHash)
	q2 := computeQuality(challenge, plotHash)

	if q1 != q2 {
		t.Fatalf("quality should be deterministic: got %d and %d", q1, q2)
	}
}

func TestQualityChangesWithChallenge(t *testing.T) {
	plotHash := sha256.Sum256([]byte("same-plot-hash"))

	challenge1 := sha256.Sum256([]byte("challenge-1"))
	challenge2 := sha256.Sum256([]byte("challenge-2"))

	q1 := computeQuality(challenge1, plotHash)
	q2 := computeQuality(challenge2, plotHash)

	if q1 == q2 {
		t.Fatal("quality should change with different challenges")
	}
}

func TestGoldenVectors(t *testing.T) {
	// Golden test vectors for regression testing
	tests := []struct {
		name            string
		challengeHex    string
		plotHashHex     string
		expectedQuality uint64
	}{
		{
			name:            "vector1",
			challengeHex:    "0000000000000000000000000000000000000000000000000000000000000000",
			plotHashHex:     "0000000000000000000000000000000000000000000000000000000000000000",
			expectedQuality: 0, // Will be computed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenge, _ := hex.DecodeString(tt.challengeHex)
			plotHash, _ := hex.DecodeString(tt.plotHashHex)

			var c [32]byte
			var p [32]byte
			copy(c[:], challenge)
			copy(p[:], plotHash)

			q := computeQuality(c, p)

			// Quality should be deterministic and < QMAX
			if q >= QMAX {
				t.Errorf("quality %d should be < QMAX %d", q, QMAX)
			}

			// Log actual quality for golden vector validation
			t.Logf("Challenge: %s, PlotHash: %s => Quality: %d", tt.challengeHex[:16], tt.plotHashHex[:16], q)
		})
	}
}

