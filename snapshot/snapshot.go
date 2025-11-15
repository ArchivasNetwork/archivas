package snapshot

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ArchivasNetwork/archivas/storage"
)

// Metadata holds information about a snapshot
type Metadata struct {
	Version       string    `json:"version"`
	NetworkID     string    `json:"network_id"`
	Height        uint64    `json:"height"`
	BlockHash     string    `json:"block_hash"`
	ExportedAt    time.Time `json:"exported_at"`
	ExportedBy    string    `json:"exported_by"`
	Description   string    `json:"description,omitempty"`
	DataDirs      []string  `json:"data_dirs"`      // Which DB directories are included
	SnapshotType  string    `json:"snapshot_type"`  // "full" or "state-only"
	TotalSizeBytes int64    `json:"total_size_bytes"`
}

// ExportOptions configures snapshot export
type ExportOptions struct {
	Height      uint64
	OutputPath  string
	DBPath      string // Path to database directory
	NetworkID   string
	Description string
	// If true, exports full block history; if false, only recent state
	FullHistory bool
}

// ImportOptions configures snapshot import
type ImportOptions struct {
	InputPath string
	DBPath    string
	Force     bool // Force import even if DB is non-empty
}

// Export creates a snapshot of the node state at a given height
func Export(db *storage.DB, blockStore *storage.BlockStorage, stateStore *storage.StateStorage, metaStore *storage.MetadataStorage, opts ExportOptions) error {
	fmt.Printf("[snapshot] Exporting snapshot at height %d...\n", opts.Height)

	// 1. Verify the block at the specified height exists
	if !blockStore.HasBlock(opts.Height) {
		return fmt.Errorf("no block found at height %d", opts.Height)
	}
	
	// We'll need to fetch the block hash from the node's RPC or pass it in
	// For now, use a placeholder that will be filled from the manifest
	blockHash := fmt.Sprintf("block-%d", opts.Height)

	// 2. Create output file
	outFile, err := os.Create(opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// 3. Create metadata
	metadata := Metadata{
		Version:       "1.0",
		NetworkID:     opts.NetworkID,
		Height:        opts.Height,
		BlockHash:     blockHash,
		ExportedAt:    time.Now(),
		ExportedBy:    "archivas-node",
		Description:   opts.Description,
		DataDirs:      []string{"blocks", "state", "meta"},
		SnapshotType:  "state-only",
		TotalSizeBytes: 0, // Will be calculated
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write metadata to tar
	metadataHeader := &tar.Header{
		Name:    "snapshot.json",
		Mode:    0644,
		Size:    int64(len(metadataJSON)),
		ModTime: time.Now(),
	}
	if err := tarWriter.WriteHeader(metadataHeader); err != nil {
		return fmt.Errorf("failed to write metadata header: %w", err)
	}
	if _, err := tarWriter.Write(metadataJSON); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// 4. Export the database directories
	// For now, we'll export the entire DB directory structure
	// In production, we'd want to be more selective and only export
	// what's needed to resume from the checkpoint height

	dbBasePath := opts.DBPath
	fmt.Printf("[snapshot] Exporting database from %s...\n", dbBasePath)

	// Walk the database directory and add all files to the tar
	var totalBytes int64
	err = filepath.Walk(dbBasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dbBasePath {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(dbBasePath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("failed to create tar header for %s: %w", path, err)
		}
		header.Name = filepath.Join("data", relPath)

		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header for %s: %w", path, err)
		}

		// If it's a file (not a directory), write its contents
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close()

			n, err := io.Copy(tarWriter, file)
			if err != nil {
				return fmt.Errorf("failed to copy file %s: %w", path, err)
			}
			totalBytes += n
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to export database: %w", err)
	}

	fmt.Printf("[snapshot] ✓ Exported %d bytes\n", totalBytes)
	fmt.Printf("[snapshot] ✓ Snapshot saved to: %s\n", opts.OutputPath)
	fmt.Printf("[snapshot] Metadata: height=%d hash=%s network=%s\n",
		metadata.Height, metadata.BlockHash, metadata.NetworkID)

	return nil
}

// Import restores a snapshot into the node database
func Import(opts ImportOptions) (*Metadata, error) {
	fmt.Printf("[snapshot] Importing snapshot from %s...\n", opts.InputPath)

	// 1. Check if DB directory is empty (unless --force is set)
	if !opts.Force {
		empty, err := isDirEmpty(opts.DBPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check DB directory: %w", err)
		}
		if !empty {
			return nil, fmt.Errorf("database directory %s is not empty; use --force to overwrite", opts.DBPath)
		}
	}

	// 2. Open the snapshot file
	inFile, err := os.Open(opts.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open snapshot file: %w", err)
	}
	defer inFile.Close()

	gzReader, err := gzip.NewReader(inFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	// 3. Read and parse metadata
	header, err := tarReader.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to read first tar entry: %w", err)
	}

	if header.Name != "snapshot.json" {
		return nil, fmt.Errorf("expected snapshot.json as first entry, got %s", header.Name)
	}

	metadataJSON, err := io.ReadAll(tarReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata Metadata
	if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	fmt.Printf("[snapshot] Snapshot info:\n")
	fmt.Printf("  Network:     %s\n", metadata.NetworkID)
	fmt.Printf("  Height:      %d\n", metadata.Height)
	fmt.Printf("  Block Hash:  %s\n", metadata.BlockHash)
	fmt.Printf("  Exported At: %s\n", metadata.ExportedAt.Format(time.RFC3339))
	fmt.Printf("  Type:        %s\n", metadata.SnapshotType)

	// 4. Create DB directory if it doesn't exist
	if err := os.MkdirAll(opts.DBPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create DB directory: %w", err)
	}

	// 5. Extract all files from the tar
	var totalBytes int64
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Construct target path (strip "data/" prefix)
		targetPath := filepath.Join(opts.DBPath, filepath.Clean(header.Name)[5:]) // Remove "data/" prefix

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			// Create parent directory if needed
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return nil, fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}

			// Create and write file
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return nil, fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			n, err := io.Copy(outFile, tarReader)
			outFile.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			totalBytes += n
		}
	}

	fmt.Printf("[snapshot] ✓ Imported %d bytes\n", totalBytes)
	fmt.Printf("[snapshot] ✓ Database restored to: %s\n", opts.DBPath)
	fmt.Printf("[snapshot] You can now start the node with:\n")
	fmt.Printf("  --checkpoint-height %d \\\n", metadata.Height)
	fmt.Printf("  --checkpoint-hash %s\n", metadata.BlockHash)

	return &metadata, nil
}

// isDirEmpty checks if a directory is empty or doesn't exist
func isDirEmpty(path string) (bool, error) {
	// If directory doesn't exist, consider it "empty"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	return len(entries) == 0, nil
}

// Manifest represents a snapshot manifest from a URL
type Manifest struct {
	Network       string `json:"network"`
	Height        uint64 `json:"height"`
	Hash          string `json:"hash"`
	SnapshotURL   string `json:"snapshot_url"`
	ChecksumSHA256 string `json:"checksum_sha256"`
}

// BootstrapOptions configures automated snapshot bootstrap
type BootstrapOptions struct {
	ManifestURL string
	DBPath      string
	Force       bool
}

// Bootstrap downloads a snapshot from a manifest URL, verifies it, and imports it
func Bootstrap(opts BootstrapOptions) (*Metadata, error) {
	fmt.Printf("[bootstrap] Fetching manifest from %s...\n", opts.ManifestURL)

	// 1. Download manifest
	manifest, err := fetchManifest(opts.ManifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}

	fmt.Printf("[bootstrap] Manifest info:\n")
	fmt.Printf("  Network:  %s\n", manifest.Network)
	fmt.Printf("  Height:   %d\n", manifest.Height)
	fmt.Printf("  Hash:     %s\n", manifest.Hash[:16]+"...")
	fmt.Printf("  Snapshot: %s\n", manifest.SnapshotURL)

	// 2. Download snapshot to temp file
	tempFile, err := os.CreateTemp("", "archivas-snapshot-*.tar.gz")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	fmt.Printf("[bootstrap] Downloading snapshot...\n")
	checksum, err := downloadFile(manifest.SnapshotURL, tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to download snapshot: %w", err)
	}

	// 3. Verify checksum
	if manifest.ChecksumSHA256 != "" {
		if checksum != manifest.ChecksumSHA256 {
			return nil, fmt.Errorf("checksum mismatch: expected %s, got %s", manifest.ChecksumSHA256, checksum)
		}
		fmt.Printf("[bootstrap] ✓ Checksum verified: %s\n", checksum[:16]+"...")
	} else {
		fmt.Printf("[bootstrap] ⚠️  No checksum in manifest, skipping verification\n")
	}

	// 4. Import snapshot
	fmt.Printf("[bootstrap] Importing snapshot...\n")
	tempFile.Close() // Close before import reads it

	importOpts := ImportOptions{
		InputPath: tempFile.Name(),
		DBPath:    opts.DBPath,
		Force:     opts.Force,
	}

	metadata, err := Import(importOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to import snapshot: %w", err)
	}

	fmt.Printf("[bootstrap] ✓ Bootstrap complete!\n")
	return metadata, nil
}

// fetchManifest downloads and parses a manifest JSON
func fetchManifest(url string) (*Manifest, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// downloadFile downloads a file from a URL and returns its SHA256 checksum
func downloadFile(url string, dest *os.File) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Calculate SHA256 while downloading
	hasher := sha256.New()
	writer := io.MultiWriter(dest, hasher)

	written, err := io.Copy(writer, resp.Body)
	if err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("[bootstrap] Downloaded %d bytes\n", written)

	return checksum, nil
}
