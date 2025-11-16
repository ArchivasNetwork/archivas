package snapshot

import (
	"fmt"

	"github.com/ArchivasNetwork/archivas/network"
)

// VerificationError represents a snapshot verification failure
type VerificationError struct {
	Field    string
	Expected interface{}
	Got      interface{}
	Message  string
}

func (e *VerificationError) Error() string {
	return fmt.Sprintf("verification failed: %s (expected: %v, got: %v)", e.Message, e.Expected, e.Got)
}

// VerifyManifest verifies a snapshot manifest against network configuration
// Phase 3: Ensures snapshots match the target network
func VerifyManifest(manifest *Manifest, profile *network.NetworkProfile, genesisHash string) error {
	// 1. Verify network name
	if manifest.Network != profile.Name {
		return &VerificationError{
			Field:    "network",
			Expected: profile.Name,
			Got:      manifest.Network,
			Message:  "network name mismatch",
		}
	}

	// 2. Verify chain ID
	if manifest.ChainID != profile.ChainID {
		return &VerificationError{
			Field:    "chain_id",
			Expected: profile.ChainID,
			Got:      manifest.ChainID,
			Message:  "chain ID mismatch",
		}
	}

	// 3. Verify network ID
	if manifest.NetworkID != profile.NetworkID {
		return &VerificationError{
			Field:    "network_id",
			Expected: profile.NetworkID,
			Got:      manifest.NetworkID,
			Message:  "network ID mismatch",
		}
	}

	// 4. Verify protocol version
	if manifest.ProtocolVersion != profile.ProtocolVersion {
		return &VerificationError{
			Field:    "protocol_version",
			Expected: profile.ProtocolVersion,
			Got:      manifest.ProtocolVersion,
			Message:  "protocol version mismatch",
		}
	}

	// 5. Verify genesis hash (if provided)
	if genesisHash != "" && manifest.GenesisHash != "" {
		if manifest.GenesisHash != genesisHash {
			return &VerificationError{
				Field:    "genesis_hash",
				Expected: genesisHash,
				Got:      manifest.GenesisHash,
				Message:  "genesis hash mismatch - this snapshot is from a different chain",
			}
		}
	}

	// 6. Verify required fields are present
	if manifest.Height == 0 {
		return fmt.Errorf("manifest missing required field: height")
	}
	if manifest.Hash == "" {
		return fmt.Errorf("manifest missing required field: hash (block hash)")
	}
	if manifest.SnapshotURL == "" {
		return fmt.Errorf("manifest missing required field: snapshot_url")
	}
	if manifest.ChecksumSHA256 == "" {
		return fmt.Errorf("manifest missing required field: checksum_sha256")
	}

	return nil
}

// VerifyStateRoot verifies the state root after importing a snapshot
// This ensures the imported state matches the manifest claim
func VerifyStateRoot(importedStateRoot, expectedStateRoot string) error {
	if expectedStateRoot == "" {
		// State root verification not available in manifest
		return nil
	}

	if importedStateRoot != expectedStateRoot {
		return &VerificationError{
			Field:    "state_root",
			Expected: expectedStateRoot,
			Got:      importedStateRoot,
			Message:  "state root mismatch after import - snapshot may be corrupted",
		}
	}

	return nil
}

// VerifyBlockHash verifies the block hash at the snapshot height
func VerifyBlockHash(actualHash, expectedHash string) error {
	if expectedHash == "" {
		return fmt.Errorf("expected block hash is empty")
	}

	if actualHash != expectedHash {
		return &VerificationError{
			Field:    "block_hash",
			Expected: expectedHash,
			Got:      actualHash,
			Message:  "block hash mismatch - snapshot may be from a forked chain",
		}
	}

	return nil
}

// VerificationResult holds the results of manifest verification
type VerificationResult struct {
	Valid          bool
	Errors         []error
	Warnings       []string
	ManifestInfo   *ManifestInfo
}

// ManifestInfo provides human-readable information about a manifest
type ManifestInfo struct {
	Network         string
	ChainID         string
	NetworkID       uint64
	ProtocolVersion int
	Height          uint64
	BlockHash       string
	StateRoot       string
	GenesisHash     string
	SnapshotSize    string
	CreatedAt       string
}

// NewManifestInfo creates ManifestInfo from a Manifest
func NewManifestInfo(m *Manifest) *ManifestInfo {
	return &ManifestInfo{
		Network:         m.Network,
		ChainID:         m.ChainID,
		NetworkID:       m.NetworkID,
		ProtocolVersion: m.ProtocolVersion,
		Height:          m.Height,
		BlockHash:       m.Hash,
		StateRoot:       m.StateRoot,
		GenesisHash:     m.GenesisHash,
		CreatedAt:       m.CreatedAt,
	}
}

// PerformFullVerification performs comprehensive verification of a manifest
func PerformFullVerification(
	manifest *Manifest,
	profile *network.NetworkProfile,
	genesisHash string,
) *VerificationResult {
	result := &VerificationResult{
		Valid:        true,
		Errors:       make([]error, 0),
		Warnings:     make([]string, 0),
		ManifestInfo: NewManifestInfo(manifest),
	}

	// Verify manifest
	if err := VerifyManifest(manifest, profile, genesisHash); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, err)
	}

	// Add warnings for missing optional fields
	if manifest.StateRoot == "" {
		result.Warnings = append(result.Warnings, "state_root not provided in manifest (skipping state verification)")
	}
	if manifest.GenesisHash == "" {
		result.Warnings = append(result.Warnings, "genesis_hash not provided in manifest (skipping genesis verification)")
	}

	return result
}

