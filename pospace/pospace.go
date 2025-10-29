package pospace

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	// PlotMagic is the magic number at the start of plot files
	PlotMagic = uint32(0x41524356) // "ARCV" in hex
	// PlotVersion is the current plot format version
	PlotVersion = uint32(1)
)

// PlotHeader contains metadata about a plot file
type PlotHeader struct {
	Magic         uint32   // Always 0x41524356 ("ARCV")
	Version       uint32   // Plot format version
	KSize         uint32   // K parameter (plot size = 2^k hashes)
	FarmerPubKey  [33]byte // Compressed secp256k1 public key
	PlotID        [32]byte // Unique plot identifier
	NumHashes     uint64   // Total number of hashes in plot
}

// PlotFile represents a Proof-of-Space plot
type PlotFile struct {
	Header PlotHeader
	Path   string
	file   *os.File
}

// Proof represents a Proof-of-Space proof
type Proof struct {
	Challenge    [32]byte // Challenge hash
	PlotID       [32]byte // Which plot
	Index        uint64   // Which hash in the plot
	Hash         [32]byte // The hash itself
	Quality      uint64   // Quality value (lower is better)
	FarmerPubKey [33]byte // Farmer's public key
}

// GeneratePlot creates a new plot file with precomputed hashes
func GeneratePlot(path string, kSize uint32, farmerPubKey []byte) error {
	if len(farmerPubKey) != 33 {
		return fmt.Errorf("farmer public key must be 33 bytes (compressed)")
	}

	// Calculate number of hashes
	numHashes := uint64(1) << kSize // 2^k

	// Generate plot ID (hash of farmer pubkey)
	plotIDHash := sha256.Sum256(farmerPubKey)

	// Create plot file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create plot file: %w", err)
	}
	defer f.Close()

	// Write header
	header := PlotHeader{
		Magic:    PlotMagic,
		Version:  PlotVersion,
		KSize:    kSize,
		PlotID:   plotIDHash,
		NumHashes: numHashes,
	}
	copy(header.FarmerPubKey[:], farmerPubKey)

	if err := binary.Write(f, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Generate and write hashes
	for i := uint64(0); i < numHashes; i++ {
		hash := computePlotHash(farmerPubKey, plotIDHash[:], i)
		if _, err := f.Write(hash[:]); err != nil {
			return fmt.Errorf("failed to write hash %d: %w", i, err)
		}

		// Progress update every 1M hashes
		if i > 0 && i%(1<<20) == 0 {
			fmt.Printf("Generated %d / %d hashes (%.1f%%)\n", i, numHashes, float64(i)*100.0/float64(numHashes))
		}
	}

	return nil
}

// computePlotHash computes a single hash for a plot entry
func computePlotHash(farmerPubKey []byte, plotID []byte, index uint64) [32]byte {
	h := sha256.New()
	h.Write(farmerPubKey)
	h.Write(plotID)
	binary.Write(h, binary.LittleEndian, index)
	return sha256.Sum256(h.Sum(nil)) // Double SHA256
}

// OpenPlot opens an existing plot file
func OpenPlot(path string) (*PlotFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plot: %w", err)
	}

	// Read header
	var header PlotHeader
	if err := binary.Read(f, binary.LittleEndian, &header); err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Validate magic
	if header.Magic != PlotMagic {
		f.Close()
		return nil, fmt.Errorf("invalid plot magic: expected %x, got %x", PlotMagic, header.Magic)
	}

	return &PlotFile{
		Header: header,
		Path:   path,
		file:   f,
	}, nil
}

// Close closes the plot file
func (p *PlotFile) Close() error {
	if p.file != nil {
		return p.file.Close()
	}
	return nil
}

// CheckChallenge checks if this plot has a winning proof for the given challenge
func (p *PlotFile) CheckChallenge(challenge [32]byte, difficultyTarget uint64) (*Proof, error) {
	// Search through the plot for qualifying hashes
	// In a real implementation, this would use a more efficient lookup structure
	// For devnet, we'll do a simple scan

	bestQuality := uint64(^uint64(0)) // max uint64
	var bestProof *Proof
	
	// Debug: Log challenge being used
	// fmt.Printf("[pospace] CheckChallenge: challenge=%x entries=%d\n", challenge[:8], p.Header.NumHashes)

	// Read through all hashes
	for i := uint64(0); i < p.Header.NumHashes; i++ {
		// Seek to hash position
		hashOffset := int64(binary.Size(p.Header)) + int64(i*32)
		if _, err := p.file.Seek(hashOffset, io.SeekStart); err != nil {
			return nil, fmt.Errorf("seek failed: %w", err)
		}

		// Read hash
		var hash [32]byte
		if _, err := io.ReadFull(p.file, hash[:]); err != nil {
			return nil, fmt.Errorf("read hash failed: %w", err)
		}

		// Compute quality
		quality := computeQuality(challenge, hash)

		// Check if this is better than our best so far
		if quality < bestQuality {
			bestQuality = quality
			bestProof = &Proof{
				Challenge:    challenge,
				PlotID:       p.Header.PlotID,
				Index:        i,
				Hash:         hash,
				Quality:      quality,
				FarmerPubKey: p.Header.FarmerPubKey,
			}
		}

		// Early exit if we found a winner
		if quality < difficultyTarget {
			return bestProof, nil
		}
	}

	// Return best proof even if it doesn't meet difficulty
	return bestProof, nil
}

// computeQuality computes the quality value for a challenge/hash pair
func computeQuality(challenge [32]byte, hash [32]byte) uint64 {
	// Mix challenge and hash
	h := sha256.New()
	h.Write(challenge[:])
	h.Write(hash[:])
	result := h.Sum(nil)

	// Take first 8 bytes as uint64 (lower is better)
	return binary.LittleEndian.Uint64(result[:8])
}

// VerifyProof verifies that a proof is valid for a given challenge
func VerifyProof(proof *Proof, challenge [32]byte, difficultyTarget uint64) bool {
	// Verify challenge matches
	if proof.Challenge != challenge {
		return false
	}

	// Recompute the hash from farmer pubkey and index
	expectedHash := computePlotHash(proof.FarmerPubKey[:], proof.PlotID[:], proof.Index)
	if expectedHash != proof.Hash {
		return false
	}

	// Recompute quality
	quality := computeQuality(challenge, proof.Hash)
	if quality != proof.Quality {
		return false
	}

	// Check difficulty
	return quality < difficultyTarget
}

