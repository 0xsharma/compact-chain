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

func (bdb *BlockchainDB) GetBlockByHash(hash *util.Hash) (*types.Block, error) {
	blockBytes, err := bdb.DB.Get(PrefixKey(HashesKey, hash.String()))
	if err != nil {
		return nil, err
	}

	block := types.DeserializeBlock(blockBytes)

	return block, nil
}

func (bdb *BlockchainDB) GetLatestBlock() (*types.Block, error) {
	lastBlockHashBytes, err := bdb.DB.Get(LastHashKey)
	if err != nil {
		return nil, err
	}

	lastHash := util.ByteToHash(lastBlockHashBytes)

	return bdb.GetBlockByHash(lastHash)
}

func (bdb *BlockchainDB) GetBlocksInRange(start uint, end uint) ([]*types.Block, error) {
	total := end - start + 1
	blocks := make([]*types.Block, total)

	lastBlockHashBytes, err := bdb.DB.Get(LastHashKey)
	if err != nil {
		return nil, err
	}

	lastHash := util.ByteToHash(lastBlockHashBytes)

	latestBlock, err := bdb.GetBlockByHash(lastHash)
	if err != nil {
		return nil, err
	}

	blocks[total-1] = latestBlock

	for i := int(total) - 2; i >= 0; i-- {
		prevHash := blocks[i+1].ParentHash
		prevBlock, err := bdb.GetBlockByHash(prevHash)

		if err != nil {
			return nil, err
		}

		blocks[i] = prevBlock
	}

	return blocks, nil
}
