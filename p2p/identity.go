package p2p

import (
	"fmt"

	"github.com/ArchivasNetwork/archivas/network"
)

// IdentityVerificationError represents a P2P identity mismatch
type IdentityVerificationError struct {
	Field    string
	Expected interface{}
	Got      interface{}
	Message  string
}

func (e *IdentityVerificationError) Error() string {
	return fmt.Sprintf("P2P identity mismatch: %s (expected: %v, got: %v)", e.Message, e.Expected, e.Got)
}

// VerifyHandshake verifies a peer's handshake message against local configuration
// Phase 3: Enforces strict chain identity matching
func VerifyHandshake(
	handshake *HandshakeMessage,
	profile *network.NetworkProfile,
	genesisHash [32]byte,
) error {
	// 1. Verify genesis hash
	if handshake.GenesisHash != genesisHash {
		return &IdentityVerificationError{
			Field:    "genesis_hash",
			Expected: fmt.Sprintf("%x", genesisHash[:8]),
			Got:      fmt.Sprintf("%x", handshake.GenesisHash[:8]),
			Message:  "peer is on a different chain (genesis hash mismatch)",
		}
	}

	// 2. Verify chain ID
	if handshake.ChainID != profile.ChainID {
		return &IdentityVerificationError{
			Field:    "chain_id",
			Expected: profile.ChainID,
			Got:      handshake.ChainID,
			Message:  "peer chain ID does not match",
		}
	}

	// 3. Verify network ID
	if handshake.NetworkID != profile.NetworkID {
		return &IdentityVerificationError{
			Field:    "network_id",
			Expected: profile.NetworkID,
			Got:      handshake.NetworkID,
			Message:  "peer network ID does not match",
		}
	}

	// 4. Verify protocol version
	if handshake.ProtocolVersion != profile.ProtocolVersion {
		return &IdentityVerificationError{
			Field:    "protocol_version",
			Expected: profile.ProtocolVersion,
			Got:      handshake.ProtocolVersion,
			Message:  "peer protocol version is incompatible",
		}
	}

	return nil
}

// CreateHandshake creates a handshake message from network configuration
// Phase 3: Includes all identity fields
func CreateHandshake(
	profile *network.NetworkProfile,
	genesisHash [32]byte,
	nodeVersion string,
	nodeName string,
) *HandshakeMessage {
	return &HandshakeMessage{
		// Phase 3: Chain identity
		GenesisHash:     genesisHash,
		ChainID:         profile.ChainID,
		NetworkID:       profile.NetworkID,
		ProtocolVersion: profile.ProtocolVersion,
		
		// Legacy compatibility
		NetworkIDLegacy:    fmt.Sprintf("%d", profile.NetworkID),
		ProtocolVersionStr: fmt.Sprintf("v%d", profile.ProtocolVersion),
		DifficultyParamsID: "standard", // Can be extended later
		
		// Informational
		NodeVersion: nodeVersion,
		NodeName:    nodeName,
	}
}

// HandshakeVerificationResult contains the result of handshake verification
type HandshakeVerificationResult struct {
	Valid          bool
	Error          error
	PeerInfo       *PeerInfo
	IsCompatible   bool
	RejectReason   string
}

// PeerInfo contains information about a peer from handshake
type PeerInfo struct {
	ChainID         string
	NetworkID       uint64
	ProtocolVersion int
	GenesisHash     string
	NodeVersion     string
	NodeName        string
}

// ExtractPeerInfo extracts peer information from handshake
func ExtractPeerInfo(handshake *HandshakeMessage) *PeerInfo {
	return &PeerInfo{
		ChainID:         handshake.ChainID,
		NetworkID:       handshake.NetworkID,
		ProtocolVersion: handshake.ProtocolVersion,
		GenesisHash:     fmt.Sprintf("%x", handshake.GenesisHash[:8]),
		NodeVersion:     handshake.NodeVersion,
		NodeName:        handshake.NodeName,
	}
}

// PerformHandshakeVerification performs comprehensive handshake verification
func PerformHandshakeVerification(
	handshake *HandshakeMessage,
	profile *network.NetworkProfile,
	genesisHash [32]byte,
) *HandshakeVerificationResult {
	result := &HandshakeVerificationResult{
		Valid:        true,
		PeerInfo:     ExtractPeerInfo(handshake),
		IsCompatible: true,
	}

	// Verify handshake
	if err := VerifyHandshake(handshake, profile, genesisHash); err != nil {
		result.Valid = false
		result.IsCompatible = false
		result.Error = err
		result.RejectReason = err.Error()
	}

	return result
}

// ShouldAcceptPeer determines if a peer should be accepted based on handshake
func ShouldAcceptPeer(
	handshake *HandshakeMessage,
	profile *network.NetworkProfile,
	genesisHash [32]byte,
) (bool, string) {
	result := PerformHandshakeVerification(handshake, profile, genesisHash)
	
	if !result.IsCompatible {
		return false, result.RejectReason
	}
	
	return true, ""
}

// LogPeerIdentity logs peer identity information
func LogPeerIdentity(peerAddr string, handshake *HandshakeMessage) {
	fmt.Printf("[p2p] Peer %s identity:\n", peerAddr)
	fmt.Printf("  Chain ID: %s\n", handshake.ChainID)
	fmt.Printf("  Network ID: %d\n", handshake.NetworkID)
	fmt.Printf("  Protocol: v%d\n", handshake.ProtocolVersion)
	fmt.Printf("  Genesis: %x...\n", handshake.GenesisHash[:8])
	if handshake.NodeName != "" {
		fmt.Printf("  Node Name: %s\n", handshake.NodeName)
	}
	if handshake.NodeVersion != "" {
		fmt.Printf("  Node Version: %s\n", handshake.NodeVersion)
	}
}

