package core

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/consensus"
	"github.com/0xsharma/compact-chain/consensus/pow"
	"github.com/0xsharma/compact-chain/dbstore"
	"github.com/0xsharma/compact-chain/executer"
	"github.com/0xsharma/compact-chain/p2p"
	"github.com/0xsharma/compact-chain/rpc"
	"github.com/0xsharma/compact-chain/txpool"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
)

type Blockchain struct {
	LastBlock    *types.Block
	Consensus    consensus.Consensus
	Mutex        *sync.RWMutex
	LastHash     *util.Hash
	BlockchainDb *dbstore.BlockchainDB
	StateDB      *dbstore.StateDB
	Txpool       *txpool.TxPool
	RPCServer    *rpc.RPCServer
	TxProcessor  *executer.TxProcessor
	Signer       *util.Address
	P2PServer    *p2p.P2PServer

	TxpoolCh     chan *types.Transaction
	BlockCh      chan *types.Block
	TxpoolChSize int
	BlockChSize  int

	MineInterrupt     chan bool
	MineInterruptSize int
}

// defaultConsensusDifficulty is the default difficulty for the proof of work consensus.
var defaultConsensusDifficulty = 10

// defaultTxpoolChSize is the default size of the txpool channel.
var defaultTxpoolChSize = 1000

// defaultBlockChSize is the default size of the block channel.
var defaultBlockChSize = 100

// defaultMineInterruptSize is the default size of the mine interrupt channel.
var defaultMineInterruptSize = 100

// NewBlockchain creates a new blockchain with the given config.
func NewBlockchain(c *config.Config) *Blockchain {
	dbInstance, err := dbstore.NewDBInstance(c.DBDir)
	if err != nil {
		panic(err)
	}

	blockchainDB := dbstore.NewBlockchainDB(dbInstance)

	stateDBInstance, err := dbstore.NewDBInstance(c.StateDBDir)
	if err != nil {
		panic(err)
	}

	stateDB := dbstore.NewStateDB(stateDBInstance)

	var genesis, lastBlock *types.Block

	lastBlockHashBytes, err := blockchainDB.DB.Get(dbstore.LastHashKey)
	if err != nil {
		genesis = CreateGenesisBlock(c.BalanceAlloc, stateDB.DB)
		lastHash := genesis.DeriveHash()

		dbBatch := blockchainDB.DB.NewBatch()

		// Batch write to db
		dbBatch.Put([]byte(dbstore.LastHashKey), lastHash.Bytes())
		dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.HashesKey, lastHash.String())), genesis.Serialize())
		dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BlockNumberKey, genesis.Number.String())), lastHash.Bytes())

		// Commit batch to db
		err = blockchainDB.DB.WriteBatch(dbBatch)
		if err != nil {
			panic(err)
		}

		lastBlock = genesis
	} else {
		lastHash := util.ByteToHash(lastBlockHashBytes)
		lastBlockBytes, err := blockchainDB.DB.Get(dbstore.PrefixKey(dbstore.HashesKey, lastHash.String()))
		if err != nil {
			panic(err)
		}
		lastBlock = types.DeserializeBlock(lastBlockBytes)
	}

	var txProcessor *executer.TxProcessor

	if c.Mine && c.SignerPrivateKey != nil {
		p := c.SignerPrivateKey.PublicKey
		txProcessor = executer.NewTxProcessor(stateDB.DB, c.MinFee, util.PublicKeyToAddress(&p))
	}

	var consensus consensus.Consensus

	switch c.ConsensusName {
	case "pow":
		if c.ConsensusDifficulty > 0 {
			consensus = pow.NewPOW(c.ConsensusDifficulty, txProcessor)
		} else {
			consensus = pow.NewPOW(defaultConsensusDifficulty, txProcessor)
		}
	default:
		panic("Invalid consensus algorithm")
	}

	txpoolChSize := defaultTxpoolChSize
	blockChSize := defaultBlockChSize
	mineInterruptSize := defaultMineInterruptSize

	txpoolCh := make(chan *types.Transaction, txpoolChSize)
	blockCh := make(chan *types.Block, blockChSize)
	mineInterrupt := make(chan bool, mineInterruptSize)

	bc_txpool := txpool.NewTxPool(c.MinFee, stateDB.DB, txpoolCh)

	rpcDomains := &rpc.RPCDomains{
		TxPool: bc_txpool,
	}
	rpcServer := rpc.NewRPCServer(c.RPCPort, rpcDomains)

	p2pServer := p2p.NewServer(c.P2PPort, c.Peers, stateDB, blockchainDB, bc_txpool, txpoolCh, blockCh)
	go p2pServer.StartServer()

	bc := &Blockchain{LastBlock: lastBlock,
		Consensus:     consensus,
		Mutex:         new(sync.RWMutex),
		BlockchainDb:  blockchainDB,
		LastHash:      lastBlock.DeriveHash(),
		StateDB:       stateDB,
		Txpool:        bc_txpool,
		TxProcessor:   txProcessor,
		RPCServer:     rpcServer,
		P2PServer:     p2pServer,
		TxpoolCh:      txpoolCh,
		BlockCh:       blockCh,
		MineInterrupt: mineInterrupt,
	}

	return bc
}

func StartBlockchain(config *config.Config) {
	chain := NewBlockchain(config)
	if chain.LastBlock.Number.Int64() == 0 {
		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String())
	} else {
		fmt.Println("LastNumber : ", chain.LastBlock.Number, "LastHash : ", chain.LastBlock.DeriveHash().String())
	}

	blockTime := config.BlockTime

	go chain.ImportBlockLoop()

	// Manual sleep to let it connect to peers
	time.Sleep(4 * time.Second)

	for {
		start := time.Now()
		lastBlockNumber := chain.LastBlock.Number

		shouldSleep := true

		err := chain.AddBlock([]byte(fmt.Sprintf("Block %d", lastBlockNumber.Int64()+1)), chain.Txpool.GetTxs(), chain.MineInterrupt, config.SignerPrivateKey)
		if err != nil {
			shouldSleep = false
		}

		var delay float64

		// Adjust delay to maintain block time
		elapsed := time.Since(start)
		if elapsed.Seconds() < float64(blockTime) {
			delay = float64(blockTime) - elapsed.Seconds()

			if shouldSleep && delay > 0 {
				time.Sleep(time.Duration(delay) * time.Second)
			}
		}
	}
}

func (bc *Blockchain) ImportBlockLoop() {
	// nolint : gosimple
	for {
		select {
		case block := <-bc.BlockCh:
			err := bc.AddExternalBlock(block)
			if err == nil {
				bc.MineInterrupt <- true
			}
		}
	}
}

// AddBlock mines and adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(data []byte, txs []*types.Transaction, mineInterrupt chan bool, signerPrivateKey *ecdsa.PrivateKey) error {
	start := time.Now()

	prevBlock := bc.LastBlock
	blockNumber := big.NewInt(0).Add(prevBlock.Number, big.NewInt(1))
	block := types.NewBlock(blockNumber, prevBlock.DeriveHash(), data)

	block.Transactions = txs

	// Mine block
	minedBlock := bc.Consensus.Mine(block, mineInterrupt)
	if minedBlock == nil {
		return errors.New("Mining interrupted")
	}

	ua := util.NewUnlockedAccount(signerPrivateKey)
	minedBlock.Sign(ua)

	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()

	dbBatch := bc.BlockchainDb.DB.NewBatch()

	// Batch write to db
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.HashesKey, minedBlock.DeriveHash().String())), minedBlock.Serialize())
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BlockNumberKey, minedBlock.Number.String())), minedBlock.DeriveHash().Bytes())
	dbBatch.Put([]byte(dbstore.LastHashKey), minedBlock.DeriveHash().Bytes())

	// Commit batch to db
	err := bc.BlockchainDb.DB.WriteBatch(dbBatch)
	if err != nil {
		panic(err)
	}

	bc.LastBlock = minedBlock
	elapsed := time.Since(start)

	for _, tx := range minedBlock.Transactions {
		// nolint : errcheck
		bc.Txpool.RemoveTx(tx)
	}

	fmt.Println("Mined block", block.Number, block.DeriveHash().String(), "Elapsed", prettySeconds(elapsed.Seconds()), "data", string(block.ExtraData), "TxCount", len(block.Transactions))

	return nil
}

func prettySeconds(f float64) string {
	return fmt.Sprintf("%.2fsec", f)
}

func (bc *Blockchain) RemoveLastBlock() {
	if bc.LastBlock.Number.Int64() == 0 {
		// Cannot remove genesis block
		return
	}

	lastBlockParentHash := bc.LastBlock.ParentHash

	dbBatch := bc.BlockchainDb.DB.NewBatch()

	// Batch write to db
	dbBatch.Delete([]byte(dbstore.PrefixKey(dbstore.HashesKey, bc.LastBlock.DeriveHash().String())))
	dbBatch.Delete([]byte(dbstore.PrefixKey(dbstore.BlockNumberKey, bc.LastBlock.Number.String())))
	dbBatch.Put([]byte(dbstore.LastHashKey), lastBlockParentHash.Bytes())

	// Commit batch to db
	err := bc.BlockchainDb.DB.WriteBatch(dbBatch)
	if err != nil {
		panic(err)
	}

	for _, tx := range bc.LastBlock.Transactions {
		err := bc.TxProcessor.RollbackTx(tx)
		if err != nil {
			fmt.Println("Failed to rollback Tx :", "tx :", tx, "error", err)
		}
	}

	newLastBlock, err := bc.BlockchainDb.GetBlockByHash(lastBlockParentHash)
	if err != nil {
		panic(err)
	}

	bc.LastBlock = newLastBlock
}

// AddBlock mines and adds a new block to the blockchain.
func (bc *Blockchain) AddExternalBlock(block *types.Block) error {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()

	currentLatestBlock := bc.LastBlock
	externalBlock := block

	if currentLatestBlock.Number.Int64() > externalBlock.Number.Int64() {
		fmt.Println("Invalid block number", block.Number, block.DeriveHash().String(), "LastBlock", bc.LastBlock.Number, bc.LastBlock.DeriveHash().String())
		return fmt.Errorf("Invalid block number")
	}

	if currentLatestBlock.Number.Int64() == externalBlock.Number.Int64() {
		if currentLatestBlock.Nonce.Int64() > externalBlock.Nonce.Int64() {
			return fmt.Errorf("Better Block already exists")
		} else if currentLatestBlock.Nonce.Int64() < externalBlock.Nonce.Int64() {
			fmt.Println("REORG : Better remote Block found", block.Number, block.DeriveHash().String())
			bc.RemoveLastBlock()
		}
	}

	if block.ParentHash.String() != bc.LastBlock.DeriveHash().String() {
		fmt.Println("Invalid parent hash", block.Number, block.ParentHash, "LastBlock", bc.LastBlock.Number, bc.LastBlock.DeriveHash().String())
		return fmt.Errorf("Invalid parent hash")
	}

	// Validate block
	if valid := bc.Consensus.Validate(block); !valid {
		fmt.Println("Invalid block", block.Number, block.DeriveHash().String())
		return fmt.Errorf("Invalid block")
	}

	dbBatch := bc.BlockchainDb.DB.NewBatch()

	// Batch write to db
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.HashesKey, block.DeriveHash().String())), block.Serialize())
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BlockNumberKey, block.Number.String())), block.DeriveHash().Bytes())
	dbBatch.Put([]byte(dbstore.LastHashKey), block.DeriveHash().Bytes())

	// Commit batch to db
	err := bc.BlockchainDb.DB.WriteBatch(dbBatch)
	if err != nil {
		panic(err)
	}

	for _, tx := range block.Transactions {
		// nolint : errcheck
		bc.Txpool.RemoveTx(tx)
	}

	bc.LastBlock = block
	fmt.Println("Imported block", block.Number, block.DeriveHash().String(), "TxCount", len(block.Transactions))

	return nil
}

// Mine the genesis block and do initial balance allocation.
func CreateGenesisBlock(balanceAlloc map[string]*big.Int, db *dbstore.DB) *types.Block {
	genesis := types.NewBlock(big.NewInt(0), util.HashData([]byte("0x0")), []byte("Genesis Block"))

	dbBatch := db.NewBatch()

	for address, balance := range balanceAlloc {
		fmt.Println("Allocating", balance, "to", address)
		dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BalanceKey, address)), balance.Bytes())
	}

	// Commit batch to db
	err := db.WriteBatch(dbBatch)
	if err != nil {
		panic(err)
	}

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
	hashBytes, err := bc.BlockchainDb.DB.Get(dbstore.PrefixKey(dbstore.BlockNumberKey, b.String()))
	if err != nil {
		return nil, err
	}

	hash := util.HashData(hashBytes)
	block, err := bc.GetBlockByHash(hash)

	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetBlockByHash returns the block with the given block hash.
func (bc *Blockchain) GetBlockByHash(h *util.Hash) (*types.Block, error) {
	blockBytes, err := bc.BlockchainDb.DB.Get(dbstore.PrefixKey(dbstore.HashesKey, h.String()))
	if err != nil {
		return nil, err
	}

	return types.DeserializeBlock(blockBytes), nil
}
