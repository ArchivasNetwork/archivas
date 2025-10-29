package storage

import (
	"encoding/json"
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

// DB wraps BadgerDB for blockchain persistence
type DB struct {
	db *badger.DB
}

// OpenDB opens or creates a BadgerDB database
func OpenDB(path string) (*DB, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Disable badger's verbose logging

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &DB{db: db}, nil
}

// Close closes the database
func (db *DB) Close() error {
	return db.db.Close()
}

// Put stores a key-value pair
func (db *DB) Put(key []byte, value []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Get retrieves a value by key
func (db *DB) Get(key []byte) ([]byte, error) {
	var value []byte
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	return value, err
}

// Has checks if a key exists
func (db *DB) Has(key []byte) bool {
	err := db.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})
	return err == nil
}

// Delete removes a key
func (db *DB) Delete(key []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// PutJSON stores a JSON-encoded value
func (db *DB) PutJSON(key []byte, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return db.Put(key, data)
}

// GetJSON retrieves and decodes a JSON value
func (db *DB) GetJSON(key []byte, value interface{}) error {
	data, err := db.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
