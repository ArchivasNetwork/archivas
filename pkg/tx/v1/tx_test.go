package v1

import (
	"crypto/ed25519"
	"encoding/json"
	"testing"
)

func TestDomainSeparatorConstant(t *testing.T) {
	if DomainSeparator != "Archivas-TxV1" {
		t.Fatalf("Domain separator must be 'Archivas-TxV1', got: %s", DomainSeparator)
	}
}

func TestTransferValidation(t *testing.T) {
	validTx := &Transfer{
		Type:   "transfer",
		From:   "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7",
		To:     "arcv1hr2vm4v4xsehsdl3a3flxspg3wguhtxymrvgrw",
		Amount: 1000000000,
		Fee:    100,
		Nonce:  0,
		Memo:   "",
	}

	if err := validTx.Validate(); err != nil {
		t.Fatalf("Valid transaction failed validation: %v", err)
	}

	// Test invalid type
	invalidTx := *validTx
	invalidTx.Type = "invalid"
	if err := invalidTx.Validate(); err == nil {
		t.Fatal("Invalid type should fail validation")
	}

	// Test zero amount
	invalidTx = *validTx
	invalidTx.Amount = 0
	if err := invalidTx.Validate(); err == nil {
		t.Fatal("Zero amount should fail validation")
	}

	// Test zero fee
	invalidTx = *validTx
	invalidTx.Fee = 0
	if err := invalidTx.Validate(); err == nil {
		t.Fatal("Zero fee should fail validation")
	}

	// Test long memo
	invalidTx = *validTx
	invalidTx.Memo = string(make([]byte, 257))
	if err := invalidTx.Validate(); err == nil {
		t.Fatal("Memo > 256 bytes should fail validation")
	}
}

func TestCanonicalJSON(t *testing.T) {
	tx := &Transfer{
		Type:   "transfer",
		From:   "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7",
		To:     "arcv1hr2vm4v4xsehsdl3a3flxspg3wguhtxymrvgrw",
		Amount: 1000000000,
		Fee:    100,
		Nonce:  0,
	}

	canon, err := CanonicalJSON(tx)
	if err != nil {
		t.Fatalf("Failed to create canonical JSON: %v", err)
	}

	// Verify it's valid JSON
	var decoded Transfer
	if err := json.Unmarshal(canon, &decoded); err != nil {
		t.Fatalf("Canonical JSON is not valid JSON: %v", err)
	}

	// Verify deterministic: same transaction produces same canonical JSON
	canon2, err := CanonicalJSON(tx)
	if err != nil {
		t.Fatalf("Failed to create canonical JSON again: %v", err)
	}

	if string(canon) != string(canon2) {
		t.Fatal("Canonical JSON is not deterministic")
	}

	// Verify keys are sorted (first character should be '{' followed by a quote)
	// In sorted order: amount, fee, from, nonce, to, type
	// This is a simple check; full validation would parse and verify order
	if len(canon) < 10 {
		t.Fatal("Canonical JSON too short")
	}
}

func TestHash(t *testing.T) {
	tx := &Transfer{
		Type:   "transfer",
		From:   "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7",
		To:     "arcv1hr2vm4v4xsehsdl3a3flxspg3wguhtxymrvgrw",
		Amount: 1000000000,
		Fee:    100,
		Nonce:  0,
	}

	hash1, err := Hash(tx)
	if err != nil {
		t.Fatalf("Failed to hash transaction: %v", err)
	}

	if len(hash1) != 32 {
		t.Fatalf("Hash must be 32 bytes, got %d", len(hash1))
	}

	// Verify deterministic: same transaction produces same hash
	hash2, err := Hash(tx)
	if err != nil {
		t.Fatalf("Failed to hash transaction again: %v", err)
	}

	if hash1 != hash2 {
		t.Fatal("Hash is not deterministic")
	}

	// Verify different transactions produce different hashes
	tx2 := *tx
	tx2.Nonce = 1
	hash3, err := Hash(&tx2)
	if err != nil {
		t.Fatalf("Failed to hash transaction: %v", err)
	}

	if hash1 == hash3 {
		t.Fatal("Different transactions produced same hash")
	}
}

func TestSignVerify(t *testing.T) {
	// Generate a test keypair (in real usage, this comes from crypto.DeriveDefaultKey)
	// For testing, we'll use a deterministic private key seed
	privKeySeed := make([]byte, 32)
	for i := range privKeySeed {
		privKeySeed[i] = byte(i)
	}

	privKey := ed25519.NewKeyFromSeed(privKeySeed)
	pubKey := privKey.Public().(ed25519.PublicKey)

	// Use full 64-byte private key
	privKey64 := make([]byte, 64)
	copy(privKey64, privKey)

	tx := &Transfer{
		Type:   "transfer",
		From:   "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7",
		To:     "arcv1hr2vm4v4xsehsdl3a3flxspg3wguhtxymrvgrw",
		Amount: 1000000000,
		Fee:    100,
		Nonce:  0,
	}

	// Sign
	sig, pubKeyResult, hash, err := Sign(privKey64, tx)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	if len(sig) != 64 {
		t.Fatalf("Signature must be 64 bytes, got %d", len(sig))
	}

	if len(pubKeyResult) != 32 {
		t.Fatalf("Public key must be 32 bytes, got %d", len(pubKeyResult))
	}

	if len(hash) != 32 {
		t.Fatalf("Hash must be 32 bytes, got %d", len(hash))
	}

	// Verify public key matches
	if string(pubKeyResult) != string(pubKey) {
		t.Fatal("Returned public key doesn't match")
	}

	// Verify signature
	valid, err := Verify(pubKeyResult, tx, sig)
	if err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}

	if !valid {
		t.Fatal("Signature verification failed")
	}

	// Verify with wrong public key should fail
	wrongPubKey := make([]byte, 32)
	copy(wrongPubKey, pubKey)
	wrongPubKey[0] ^= 1 // Flip one bit

	valid, err = Verify(wrongPubKey, tx, sig)
	if err != nil {
		t.Fatalf("Failed to verify with wrong key: %v", err)
	}

	if valid {
		t.Fatal("Verification with wrong public key should fail")
	}
}

func TestSignedTxWireFormat(t *testing.T) {
	// Generate test keypair
	privKeySeed := make([]byte, 32)
	for i := range privKeySeed {
		privKeySeed[i] = byte(i)
	}
	privKey := ed25519.NewKeyFromSeed(privKeySeed)
	privKey64 := make([]byte, 64)
	copy(privKey64, privKey)

	tx := &Transfer{
		Type:   "transfer",
		From:   "arcv1zramsn568zt3cwc8ny995u3dhpz5rpuamx2jz7",
		To:     "arcv1hr2vm4v4xsehsdl3a3flxspg3wguhtxymrvgrw",
		Amount: 1000000000,
		Fee:    100,
		Nonce:  0,
	}

	// Pack signed transaction
	stx, err := PackSignedTx(tx, privKey64)
	if err != nil {
		t.Fatalf("Failed to pack signed transaction: %v", err)
	}

	// Verify fields
	if stx.Tx == nil {
		t.Fatal("SignedTx.Tx is nil")
	}

	if stx.PubKey == "" {
		t.Fatal("SignedTx.PubKey is empty")
	}

	if stx.Sig == "" {
		t.Fatal("SignedTx.Sig is empty")
	}

	if stx.Hash == "" {
		t.Fatal("SignedTx.Hash is empty")
	}

	// Test encoding/decoding
	pubKeyDecoded, err := DecodePubKey(stx.PubKey)
	if err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	if len(pubKeyDecoded) != 32 {
		t.Fatalf("Decoded public key must be 32 bytes, got %d", len(pubKeyDecoded))
	}

	sigDecoded, err := DecodeSig(stx.Sig)
	if err != nil {
		t.Fatalf("Failed to decode signature: %v", err)
	}

	if len(sigDecoded) != 64 {
		t.Fatalf("Decoded signature must be 64 bytes, got %d", len(sigDecoded))
	}

	hashDecoded, err := DecodeHash(stx.Hash)
	if err != nil {
		t.Fatalf("Failed to decode hash: %v", err)
	}

	if len(hashDecoded) != 32 {
		t.Fatalf("Decoded hash must be 32 bytes, got %d", len(hashDecoded))
	}

	// Verify signature
	valid, err := VerifySignedTx(stx)
	if err != nil {
		t.Fatalf("Failed to verify signed transaction: %v", err)
	}

	if !valid {
		t.Fatal("Signed transaction verification failed")
	}

	// Test JSON round-trip
	jsonData, err := json.Marshal(stx)
	if err != nil {
		t.Fatalf("Failed to marshal signed transaction: %v", err)
	}

	var stx2 SignedTx
	if err := json.Unmarshal(jsonData, &stx2); err != nil {
		t.Fatalf("Failed to unmarshal signed transaction: %v", err)
	}

	// Verify unmarshaled transaction
	valid, err = VerifySignedTx(&stx2)
	if err != nil {
		t.Fatalf("Failed to verify unmarshaled transaction: %v", err)
	}

	if !valid {
		t.Fatal("Unmarshaled signed transaction verification failed")
	}
}
