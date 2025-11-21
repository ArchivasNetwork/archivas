package ledger

import (
	"strings"

	"github.com/ArchivasNetwork/archivas/address"
)

// NormalizeAddress converts any address (ARCV or 0x) to a canonical hex format for internal storage.
// This ensures that ARCV and 0x addresses referring to the same account are treated identically.
//
// Returns lowercase 0x-prefixed hex address (e.g., "0xabcd...1234")
func NormalizeAddress(addr string) (string, error) {
	// Parse address (handles both ARCV and 0x)
	evmAddr, err := address.ParseAddress(addr, "arcv")
	if err != nil {
		return "", err
	}

	// Return as lowercase hex (canonical form)
	return strings.ToLower(evmAddr.Hex()), nil
}

// AddressesEqual checks if two addresses (in any format) refer to the same account.
func AddressesEqual(addr1, addr2 string) bool {
	norm1, err1 := NormalizeAddress(addr1)
	norm2, err2 := NormalizeAddress(addr2)

	if err1 != nil || err2 != nil {
		return false
	}

	return norm1 == norm2
}

