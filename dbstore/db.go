package dbstore

import (
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	LastHashKey    = "lh" // Last hash key ( lastHash -> hash)
	HashesKey      = "hs" // Hashes key (hash->block)
	BlockNumberKey = "bn" // Block number key (blockNumber -> hash)
	BalanceKey     = "bl" // Balance key (address -> balance)
	NonceKey       = "nc" // Nonce key (address -> nonce)
)

// PrefixKey prefixes a string with another string.
func PrefixKey(prefix string, str string) string {
	return prefix + str
}

// DB is a wrapper around leveldb.
type DB struct {
	dbPath  string
	LevelDb *leveldb.DB
}

// NewDB creates a new DB instance.
func NewDBInstance(dbPath string) (*DB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	return &DB{dbPath: dbPath, LevelDb: db}, nil
}

// Get returns the value for the given key.
func (db *DB) Get(key string) ([]byte, error) {
	value, err := db.LevelDb.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Put sets the value for the given key.
func (db *DB) Put(key string, value []byte) error {
	err := db.LevelDb.Put([]byte(key), value, nil)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the value for the given key.
func (db *DB) Delete(key string) error {
	err := db.LevelDb.Delete([]byte(key), nil)
	if err != nil {
		return err
	}

	return nil
}

// Has returns true if the given key exists.
func (db *DB) Has(key string) (bool, error) {
	has, err := db.LevelDb.Has([]byte(key), nil)
	if err != nil {
		return false, err
	}

	return has, nil
}

// NewBatch creates a new batch.
func (db *DB) NewBatch() *leveldb.Batch {
	return new(leveldb.Batch)
}

// WriteBatch writes a batch to the db.
func (db *DB) WriteBatch(batch *leveldb.Batch) error {
	err := db.LevelDb.Write(batch, nil)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the db.
func (db *DB) Close() error {
	err := db.LevelDb.Close()
	if err != nil {
		return err
	}

	return nil
}
