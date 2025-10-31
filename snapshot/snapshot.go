package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// SnapshotMetadata contains snapshot information
type SnapshotMetadata struct {
	Height    uint64 `json:"height"`
	StateRoot string `json:"stateRoot"`
	CreatedAt string `json:"createdAt"` // RFC3339
	File      string `json:"file"`      // Filename
	SHA256    string `json:"sha256"`    // Hash for verification
	Size      int64  `json:"size"`      // File size in bytes
}

// CreateSnapshot creates a new state snapshot
func CreateSnapshot(height uint64, outputDir string) (*SnapshotMetadata, error) {
	now := time.Now()
	
	// Create snapshot metadata
	stateRoot := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("state-%d", height))))[:16]
	filename := fmt.Sprintf("state-%d-%s.snap", height, stateRoot)
	
	meta := &SnapshotMetadata{
		Height:    height,
		StateRoot: stateRoot,
		CreatedAt: now.Format(time.RFC3339),
		File:      filename,
		SHA256:    "",
		Size:      0,
	}
	
	// Would create actual snapshot file here
	// For now, just metadata
	
	return meta, nil
}

// VerifySnapshot verifies snapshot integrity
func VerifySnapshot(path string, expectedSHA256 string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	
	hash := sha256.Sum256(data)
	actualSHA256 := hex.EncodeToString(hash[:])
	
	return actualSHA256 == expectedSHA256, nil
}

// SaveMetadata saves snapshot metadata to JSON
func SaveMetadata(meta *SnapshotMetadata, path string) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// LoadMetadata loads snapshot metadata from JSON
func LoadMetadata(path string) (*SnapshotMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var meta SnapshotMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	
	return &meta, nil
}

