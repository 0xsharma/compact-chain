package dbstore

type StateDB struct {
	DB *DB
}

func NewStateDB(db *DB) *StateDB {
	return &StateDB{DB: db}
}
