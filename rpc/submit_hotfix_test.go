package rpc

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	txv1 "github.com/ArchivasNetwork/archivas/pkg/tx/v1"
)

// TestSubmitMethodHandling tests that /submit returns 405 for non-POST methods
func TestSubmitMethodHandling(t *testing.T) {
	// Create minimal server
	server := &FarmingServer{}

	// Test GET /submit (should return 405 with Allow: POST)
	req := httptest.NewRequest("GET", "/submit", nil)
	w := httptest.NewRecorder()
	server.handleSubmitV1(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}

	// Check Allow header
	allow := w.Header().Get("Allow")
	if allow != "POST" {
		t.Errorf("Expected Allow: POST, got %s", allow)
	}

	// Test PUT /submit
	req = httptest.NewRequest("PUT", "/submit", nil)
	w = httptest.NewRecorder()
	server.handleSubmitV1(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}

	allow = w.Header().Get("Allow")
	if allow != "POST" {
		t.Errorf("Expected Allow: POST, got %s", allow)
	}
}

// TestSubmitContentTypeValidation tests that /submit requires Content-Type: application/json
func TestSubmitContentTypeValidation(t *testing.T) {
	server := &FarmingServer{}

	// Test POST /submit without Content-Type
	req := httptest.NewRequest("POST", "/submit", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()
	server.handleSubmitV1(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected 415, got %d", w.Code)
	}

	// Check response format
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if ok, _ := response["ok"].(bool); ok {
		t.Error("Expected ok: false")
	}

	if errorMsg, _ := response["error"].(string); errorMsg != "Content-Type must be application/json" {
		t.Errorf("Expected error message about Content-Type, got: %s", errorMsg)
	}

	// Test POST /submit with wrong Content-Type
	req = httptest.NewRequest("POST", "/submit", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "text/plain")
	w = httptest.NewRecorder()
	server.handleSubmitV1(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected 415, got %d", w.Code)
	}
}

// TestSubmitValidJSON tests that POST /submit with valid JSON doesn't return 405 or 415
func TestSubmitValidJSON(t *testing.T) {
	server := &FarmingServer{}

	// Create a minimal valid signed transaction
	tx := &txv1.Transfer{
		Type:   "transfer",
		From:   "arcv1test",
		To:     "arcv1test2",
		Amount: 1000000000,
		Fee:    100,
		Nonce:  0,
	}

	// Sign transaction (using a dummy key for testing)
	privKeySeed := make([]byte, 32)
	for i := range privKeySeed {
		privKeySeed[i] = byte(i)
	}
	privKey := ed25519.NewKeyFromSeed(privKeySeed)
	privKey64 := make([]byte, 64)
	copy(privKey64, privKey)

	stx, err := txv1.PackSignedTx(tx, privKey64)
	if err != nil {
		t.Fatalf("Failed to pack signed transaction: %v", err)
	}

	jsonData, err := json.Marshal(stx)
	if err != nil {
		t.Fatalf("Failed to marshal signed transaction: %v", err)
	}

	// Test POST /submit with valid JSON and Content-Type
	req := httptest.NewRequest("POST", "/submit", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleSubmitV1(w, req)

	// Should not return 405 or 415
	if w.Code == http.StatusMethodNotAllowed {
		t.Error("POST /submit returned 405 Method Not Allowed")
	}

	if w.Code == http.StatusUnsupportedMediaType {
		t.Error("POST /submit returned 415 Unsupported Media Type")
	}

	// Should return some response (200, 400, etc., but not 405/415)
	if w.Code < 200 || w.Code >= 500 {
		t.Errorf("Unexpected status code: %d", w.Code)
	}
}

// TestBroadcastCompatibility tests that POST /broadcast works with v1.1.0 format
func TestBroadcastCompatibility(t *testing.T) {
	server := &FarmingServer{}

	// Create a minimal valid signed transaction
	tx := &txv1.Transfer{
		Type:   "transfer",
		From:   "arcv1test",
		To:     "arcv1test2",
		Amount: 1000000000,
		Fee:    100,
		Nonce:  0,
	}

	privKeySeed := make([]byte, 32)
	for i := range privKeySeed {
		privKeySeed[i] = byte(i)
	}
	privKey := ed25519.NewKeyFromSeed(privKeySeed)
	privKey64 := make([]byte, 64)
	copy(privKey64, privKey)

	stx, err := txv1.PackSignedTx(tx, privKey64)
	if err != nil {
		t.Fatalf("Failed to pack signed transaction: %v", err)
	}

	jsonData, err := json.Marshal(stx)
	if err != nil {
		t.Fatalf("Failed to marshal signed transaction: %v", err)
	}

	// Test POST /broadcast with v1.1.0 format (should route to handleSubmitV1)
	req := httptest.NewRequest("POST", "/broadcast", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleBroadcast(w, req)

	// Should not return 405 or 415
	if w.Code == http.StatusMethodNotAllowed {
		t.Error("POST /broadcast returned 405 Method Not Allowed")
	}

	if w.Code == http.StatusUnsupportedMediaType {
		t.Error("POST /broadcast returned 415 Unsupported Media Type")
	}

	// Should return some response (200, 400, etc.)
	if w.Code < 200 || w.Code >= 500 {
		t.Errorf("Unexpected status code: %d", w.Code)
	}
}
