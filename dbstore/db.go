package dbstore

import (
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	LastHashKey    = "lh"
	HashesKey      = "hs"
	BlockNumberKey = "bn"
)

func PrefixKey(prefix string, str string) string {
	return prefix + str
}

type DB struct {
	dbPath  string
	LevelDb *leveldb.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return &DB{dbPath: dbPath, LevelDb: db}, nil
}

func (db *DB) Get(key string) ([]byte, error) {
	value, err := db.LevelDb.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (db *DB) Put(key string, value []byte) error {
	err := db.LevelDb.Put([]byte(key), value, nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Delete(key string) error {
	err := db.LevelDb.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Has(key string) (bool, error) {
	has, err := db.LevelDb.Has([]byte(key), nil)
	if err != nil {
		return false, err
	}
	return has, nil
}

func (db *DB) NewBatch() *leveldb.Batch {
	return new(leveldb.Batch)
}

func (db *DB) WriteBatch(batch *leveldb.Batch) error {
	err := db.LevelDb.Write(batch, nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() error {
	err := db.LevelDb.Close()
	if err != nil {
		return err
	}
	return nil
}
