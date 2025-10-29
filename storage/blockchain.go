package storage

import (
	"encoding/binary"
	"fmt"
)

// Key prefixes for different data types
var (
	PrefixBlock      = []byte("blk:") // blk:<height> → block data
	PrefixAccount    = []byte("acc:") // acc:<address> → account state
	KeyTipHeight     = []byte("meta:tip_height")
	KeyDifficulty    = []byte("meta:difficulty")
	KeyVDFSeed       = []byte("meta:vdf_seed")
	KeyVDFIterations = []byte("meta:vdf_iterations")
	KeyVDFOutput     = []byte("meta:vdf_output")
)

// BlockStorage handles block persistence
type BlockStorage struct {
	db *DB
}

// NewBlockStorage creates a new block storage
func NewBlockStorage(db *DB) *BlockStorage {
	return &BlockStorage{db: db}
}

// SaveBlock persists a block
func (bs *BlockStorage) SaveBlock(height uint64, blockData interface{}) error {
	key := makeBlockKey(height)
	return bs.db.PutJSON(key, blockData)
}

// LoadBlock retrieves a block
func (bs *BlockStorage) LoadBlock(height uint64, blockData interface{}) error {
	key := makeBlockKey(height)
	return bs.db.GetJSON(key, blockData)
}

// HasBlock checks if a block exists
func (bs *BlockStorage) HasBlock(height uint64) bool {
	key := makeBlockKey(height)
	return bs.db.Has(key)
}

// AccountState represents stored account state
type AccountState struct {
	Balance int64  `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

// StateStorage handles world state persistence
type StateStorage struct {
	db *DB
}

// NewStateStorage creates a new state storage
func NewStateStorage(db *DB) *StateStorage {
	return &StateStorage{db: db}
}

// SaveAccount persists an account state
func (ss *StateStorage) SaveAccount(address string, balance int64, nonce uint64) error {
	key := makeAccountKey(address)
	state := AccountState{
		Balance: balance,
		Nonce:   nonce,
	}
	return ss.db.PutJSON(key, state)
}

// LoadAccount retrieves an account state
func (ss *StateStorage) LoadAccount(address string) (balance int64, nonce uint64, exists bool, err error) {
	key := makeAccountKey(address)

	var state AccountState
	err = ss.db.GetJSON(key, &state)
	if err != nil {
		// Key not found is not an error - account just doesn't exist
		return 0, 0, false, nil
	}

	return state.Balance, state.Nonce, true, nil
}

// MetadataStorage handles chain metadata
type MetadataStorage struct {
	db *DB
}

// NewMetadataStorage creates a new metadata storage
func NewMetadataStorage(db *DB) *MetadataStorage {
	return &MetadataStorage{db: db}
}

// SaveTipHeight saves the current tip height
func (ms *MetadataStorage) SaveTipHeight(height uint64) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, height)
	return ms.db.Put(KeyTipHeight, data)
}

// LoadTipHeight loads the tip height
func (ms *MetadataStorage) LoadTipHeight() (uint64, error) {
	data, err := ms.db.Get(KeyTipHeight)
	if err != nil {
		return 0, err
	}
	if len(data) != 8 {
		return 0, fmt.Errorf("invalid tip height data")
	}
	return binary.BigEndian.Uint64(data), nil
}

// SaveDifficulty saves the current difficulty
func (ms *MetadataStorage) SaveDifficulty(difficulty uint64) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, difficulty)
	return ms.db.Put(KeyDifficulty, data)
}

// LoadDifficulty loads the difficulty
func (ms *MetadataStorage) LoadDifficulty() (uint64, error) {
	data, err := ms.db.Get(KeyDifficulty)
	if err != nil {
		return 0, err
	}
	if len(data) != 8 {
		return 0, fmt.Errorf("invalid difficulty data")
	}
	return binary.BigEndian.Uint64(data), nil
}

// SaveVDFState saves the VDF state
func (ms *MetadataStorage) SaveVDFState(seed []byte, iterations uint64, output []byte) error {
	if err := ms.db.Put(KeyVDFSeed, seed); err != nil {
		return err
	}

	iterData := make([]byte, 8)
	binary.BigEndian.PutUint64(iterData, iterations)
	if err := ms.db.Put(KeyVDFIterations, iterData); err != nil {
		return err
	}

	return ms.db.Put(KeyVDFOutput, output)
}

// LoadVDFState loads the VDF state
func (ms *MetadataStorage) LoadVDFState() (seed []byte, iterations uint64, output []byte, err error) {
	seed, err = ms.db.Get(KeyVDFSeed)
	if err != nil {
		return nil, 0, nil, err
	}

	iterData, err := ms.db.Get(KeyVDFIterations)
	if err != nil {
		return nil, 0, nil, err
	}
	iterations = binary.BigEndian.Uint64(iterData)

	output, err = ms.db.Get(KeyVDFOutput)
	if err != nil {
		return nil, 0, nil, err
	}

	return seed, iterations, output, nil
}

// Helper functions to create keys
func makeBlockKey(height uint64) []byte {
	key := make([]byte, len(PrefixBlock)+8)
	copy(key, PrefixBlock)
	binary.BigEndian.PutUint64(key[len(PrefixBlock):], height)
	return key
}

func makeAccountKey(address string) []byte {
	key := make([]byte, len(PrefixAccount)+len(address))
	copy(key, PrefixAccount)
	copy(key[len(PrefixAccount):], []byte(address))
	return key
}
