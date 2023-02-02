package dbstore

import (
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	homePath, _ = os.UserHomeDir()
	dbPath      = homePath + "/.compact-chain/db"
)

type DB struct {
	dbPath  string
	levelDb *leveldb.DB
}

func NewDB() (*DB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return &DB{dbPath: dbPath, levelDb: db}, nil
}

func (db *DB) Get(key string) ([]byte, error) {
	value, err := db.levelDb.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (db *DB) Put(key string, value []byte) error {
	err := db.levelDb.Put([]byte(key), value, nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Delete(key string) error {
	err := db.levelDb.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Has(key string) (bool, error) {
	has, err := db.levelDb.Has([]byte(key), nil)
	if err != nil {
		return false, err
	}
	return has, nil
}

func (db *DB) Batch() *leveldb.Batch {
	return new(leveldb.Batch)
}

func (db *DB) Close() error {
	err := db.levelDb.Close()
	if err != nil {
		return err
	}
	return nil
}
