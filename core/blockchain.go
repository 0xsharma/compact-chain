package core

import (
	"math/big"
	"sync"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/consensus"
	"github.com/0xsharma/compact-chain/consensus/pow"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
)

type Blockchain struct {
	blocks    []*types.Block
	consensus consensus.Consensus
	mutex     *sync.RWMutex
}

var defaultAccordDifficulty = 10

func NewBlockchain(c *config.Config) *Blockchain {
	genesis := CreateGenesisBlock()
	var consensus consensus.Consensus

	switch c.Accord {
	case "pow":
		if c.AccordDifficulty > 0 {
			consensus = pow.NewPOW(c.AccordDifficulty)
		} else {
			consensus = pow.NewPOW(defaultAccordDifficulty)
		}
	default:
		panic("Invalid consensus algorithm")

	}
	return &Blockchain{[]*types.Block{genesis}, consensus, new(sync.RWMutex)}
}

func (bc *Blockchain) AddBlock(data []byte) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	prevBlock := bc.blocks[len(bc.blocks)-1]
	blockNumber := big.NewInt(0).Add(prevBlock.Number(), big.NewInt(1))
	block := types.NewBlock(blockNumber, prevBlock.Hash(), data)
	minedBlock := bc.consensus.Mine(block)
	bc.blocks = append(bc.blocks, minedBlock)
}

func CreateGenesisBlock() *types.Block {
	return types.NewBlock(big.NewInt(0), util.NewHashFromHex("0x0"), []byte("Genesis Block"))
}

func (bc *Blockchain) Current() *types.Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) GetBlockByNumber(b *big.Int) *types.Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.blocks[b.Int64()]
}
