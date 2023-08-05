package dbstore

import (
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
)

type BlockchainDB struct {
	DB *DB
}

func NewBlockchainDB(db *DB) *BlockchainDB {
	return &BlockchainDB{DB: db}
}

func (bdb *BlockchainDB) GetLatestBlock() (*types.Block, error) {
	lastBlockHashBytes, err := bdb.DB.Get(LastHashKey)
	if err != nil {
		return nil, err
	}

	lastHash := util.ByteToHash(lastBlockHashBytes)

	lastBlockBytes, err := bdb.DB.Get(PrefixKey(HashesKey, lastHash.String()))
	if err != nil {
		return nil, err
	}

	latestBlock := types.DeserializeBlock(lastBlockBytes)

	return latestBlock, nil
}
