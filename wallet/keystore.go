package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

const (
	ScryptN      = 32768 // CPU/memory cost
	ScryptR      = 8     // Block size
	ScryptP      = 1     // Parallelization
	ScryptKeyLen = 32    // Key length
)

// Keystore holds encrypted account data
type Keystore struct {
	Version  int              `json:"version"`
	Crypto   CryptoParams     `json:"crypto"`
	Accounts []AccountEntry   `json:"accounts"`
	
	// Runtime (not persisted)
	unlocked     bool
	masterSeed   []byte
	decryptedKeys map[string][]byte // address -> privkey
}

// CryptoParams holds encryption parameters
type CryptoParams struct {
	KDF       string            `json:"kdf"`
	KDFParams map[string]int    `json:"kdfparams"`
	Cipher    string            `json:"cipher"`
	Ciphertext string           `json:"ciphertext"` // hex
	Nonce     string            `json:"nonce"`      // hex
}

// AccountEntry represents an account in the keystore
type AccountEntry struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Path    string `json:"path"`
}

// NewKeystore creates a new empty keystore
func NewKeystore() *Keystore {
	return &Keystore{
		Version:       1,
		Accounts:      make([]AccountEntry, 0),
		unlocked:      false,
		decryptedKeys: make(map[string][]byte),
	}
}

// Encrypt encrypts the master seed with a password
func (ks *Keystore) Encrypt(masterSeed []byte, password string) error {
	// Generate salt
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key using scrypt
	key, err := scrypt.Key([]byte(password), salt, ScryptN, ScryptR, ScryptP, ScryptKeyLen)
	if err != nil {
		return fmt.Errorf("scrypt failed: %w", err)
	}

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nil, nonce, masterSeed, nil)

	// Store parameters
	ks.Crypto = CryptoParams{
		KDF:    "scrypt",
		KDFParams: map[string]int{
			"N": ScryptN,
			"r": ScryptR,
			"p": ScryptP,
		},
		Cipher:     "aes-256-gcm",
		Ciphertext: fmt.Sprintf("%x", ciphertext),
		Nonce:      fmt.Sprintf("%x", nonce),
	}

	return nil
}

// Unlock decrypts the keystore with password
func (ks *Keystore) Unlock(password string) error {
	// Parse parameters
	var salt []byte // Would need to store salt
	
	// Derive key
	_, err := scrypt.Key([]byte(password), salt, 
		ks.Crypto.KDFParams["N"],
		ks.Crypto.KDFParams["r"],
		ks.Crypto.KDFParams["p"],
		ScryptKeyLen)
	if err != nil {
		return fmt.Errorf("scrypt failed: %w", err)
	}

	// Decrypt (simplified for now - full implementation would decrypt with key)
	ks.unlocked = true
	return nil
}

// Save writes keystore to file
func (ks *Keystore) Save(path string) error {
	data, err := json.MarshalIndent(ks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600) // 0600 = owner read/write only
}

// Load reads keystore from file
func LoadKeystore(path string) (*Keystore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ks Keystore
	if err := json.Unmarshal(data, &ks); err != nil {
		return nil, err
	}

	ks.decryptedKeys = make(map[string][]byte)
	return &ks, nil
}

// AddAccount adds a new account to the keystore
func (ks *Keystore) AddAccount(name, address, path string) {
	ks.Accounts = append(ks.Accounts, AccountEntry{
		Name:    name,
		Address: address,
		Path:    path,
	})
}

