package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// CanonicalJSON encodes a transaction to RFC 8785 canonical JSON.
// Rules:
// - Keys sorted lexicographically
// - No whitespace
// - UTF-8 encoding
// - No duplicate keys
func CanonicalJSON(tx *Transfer) ([]byte, error) {
	// Build a map for sorting
	m := make(map[string]interface{})
	m["type"] = tx.Type
	m["from"] = tx.From
	m["to"] = tx.To
	m["amount"] = tx.Amount
	m["fee"] = tx.Fee
	m["nonce"] = tx.Nonce
	if tx.Memo != "" {
		m["memo"] = tx.Memo
	}

	// Sort keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build JSON string manually to ensure canonical format
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// Write key (JSON-escaped)
		keyJSON, err := json.Marshal(k)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal key: %w", err)
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')

		// Write value (JSON-encoded)
		valJSON, err := json.Marshal(m[k])
		if err != nil {
			return nil, fmt.Errorf("failed to marshal value: %w", err)
		}
		buf.Write(valJSON)
	}
	buf.WriteByte('}')

	return buf.Bytes(), nil
}
