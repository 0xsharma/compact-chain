package core

import (
	"math/big"
	"sync"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/consensus"
	"github.com/0xsharma/compact-chain/consensus/pow"
	"github.com/0xsharma/compact-chain/dbstore"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
)

type Blockchain struct {
	LastBlock *types.Block
	Consensus consensus.Consensus
	Mutex     *sync.RWMutex
	LastHash  *util.Hash
	Db        *dbstore.DB
	StateDB   *dbstore.DB
}

// defaultConsensusDifficulty is the default difficulty for the proof of work consensus.
var defaultConsensusDifficulty = 10

// NewBlockchain creates a new blockchain with the given config.
func NewBlockchain(c *config.Config) *Blockchain {
	db, err := dbstore.NewDB(c.DBDir)
	if err != nil {
		panic(err)
	}

	stateDB, err := dbstore.NewDB(c.StateDBDir)
	if err != nil {
		panic(err)
	}

	var genesis, lastBlock *types.Block

	lastBlockBytes, err := db.Get(dbstore.LastHashKey)
	if err != nil {
		genesis = CreateGenesisBlock()
		lastHash := genesis.Hash

		dbBatch := db.NewBatch()

		// Batch write to db
		dbBatch.Put([]byte(dbstore.LastHashKey), lastHash.Bytes())
		dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.HashesKey, lastHash.String())), genesis.Serialize())
		dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BlockNumberKey, genesis.Number.String())), lastHash.Bytes())

		// Commit batch to db
		err = db.WriteBatch(dbBatch)
		if err != nil {
			panic(err)
		}

		lastBlock = genesis
	} else {
		lastHash := util.ByteToHash(lastBlockBytes)
		lastBlockBytes, err = db.Get(dbstore.PrefixKey(dbstore.HashesKey, lastHash.String()))
		if err != nil {
			panic(err)
		}
		lastBlock = types.DeserializeBlock(lastBlockBytes)
	}

	var consensus consensus.Consensus

	switch c.ConsensusName {
	case "pow":
		if c.ConsensusDifficulty > 0 {
			consensus = pow.NewPOW(c.ConsensusDifficulty)
		} else {
			consensus = pow.NewPOW(defaultConsensusDifficulty)
		}
	default:
		panic("Invalid consensus algorithm")
	}

	bc := &Blockchain{LastBlock: lastBlock, Consensus: consensus, Mutex: new(sync.RWMutex), Db: db, LastHash: lastBlock.Hash, StateDB: stateDB}

	return bc
}

// AddBlock mines and adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(data []byte) {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()

	prevBlock := bc.LastBlock
	blockNumber := big.NewInt(0).Add(prevBlock.Number, big.NewInt(1))
	block := types.NewBlock(blockNumber, prevBlock.Hash, data)

	// Mine block
	minedBlock := bc.Consensus.Mine(block)

	dbBatch := bc.Db.NewBatch()

	// Batch write to db
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.HashesKey, minedBlock.Hash.String())), minedBlock.Serialize())
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BlockNumberKey, minedBlock.Number.String())), minedBlock.Hash.Bytes())
	dbBatch.Put([]byte(dbstore.LastHashKey), minedBlock.Hash.Bytes())

	// Commit batch to db
	err := bc.Db.WriteBatch(dbBatch)
	if err != nil {
		panic(err)
	}

	bc.LastBlock = minedBlock
}

// Mine the genesis block
func CreateGenesisBlock() *types.Block {
	genesis := types.NewBlock(big.NewInt(0), util.NewHash([]byte("0x0")), []byte("Genesis Block"))
	return genesis
}

// Current returns the current block in the blockchain.
func (bc *Blockchain) Current() *types.Block {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()

	return bc.LastBlock
}

// GetBlockByNumber returns the block with the given block number.
func (bc *Blockchain) GetBlockByNumber(b *big.Int) (*types.Block, error) {
	hashBytes, err := bc.Db.Get(dbstore.PrefixKey(dbstore.BlockNumberKey, b.String()))
	if err != nil {
		return nil, err
	}

	hash := util.NewHash(hashBytes)
	block, err := bc.GetBlockByHash(hash)

	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetBlockByHash returns the block with the given block hash.
func (bc *Blockchain) GetBlockByHash(h *util.Hash) (*types.Block, error) {
	blockBytes, err := bc.Db.Get(dbstore.PrefixKey(dbstore.HashesKey, h.String()))
	if err != nil {
		return nil, err
	}

	return types.DeserializeBlock(blockBytes), nil
}
