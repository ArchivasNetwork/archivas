package p2p

import (
	"encoding/json"
	"testing"
)

func TestRequestBlocksMessage(t *testing.T) {
	req := RequestBlocksMessage{
		FromHeight: 100,
		MaxBlocks:  512,
	}

	// Marshal
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded RequestBlocksMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.FromHeight != req.FromHeight || decoded.MaxBlocks != req.MaxBlocks {
		t.Errorf("Round-trip mismatch: got %+v, want %+v", decoded, req)
	}
}

func TestBlocksBatchMessage(t *testing.T) {
	batch := BlocksBatchMessage{
		FromHeight: 100,
		Count:      3,
		Blocks:     []json.RawMessage{[]byte("{}"), []byte("{}"), []byte("{}")},
		TipHeight:  1000,
		EOF:        false,
	}

	// Marshal
	data, err := json.Marshal(batch)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded BlocksBatchMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.FromHeight != batch.FromHeight ||
		decoded.Count != batch.Count ||
		decoded.TipHeight != batch.TipHeight ||
		decoded.EOF != batch.EOF {
		t.Errorf("Round-trip mismatch: got %+v, want %+v", decoded, batch)
	}

	if len(decoded.Blocks) != int(batch.Count) {
		t.Errorf("Block count mismatch: got %d, want %d", len(decoded.Blocks), batch.Count)
	}
}

func TestBlocksBatchEOF(t *testing.T) {
	// Test EOF=true with empty blocks (caught up)
	batch := BlocksBatchMessage{
		FromHeight: 1000,
		Count:      0,
		Blocks:     []json.RawMessage{},
		TipHeight:  999,
		EOF:        true,
	}

	data, _ := json.Marshal(batch)
	var decoded BlocksBatchMessage
	json.Unmarshal(data, &decoded)

	if !decoded.EOF {
		t.Error("EOF should be true for empty catch-up batch")
	}

	if decoded.Count != 0 {
		t.Error("Count should be 0 for EOF batch")
	}
}

