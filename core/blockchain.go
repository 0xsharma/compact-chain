package core

import (
	"math/big"
	"sync"

	"github.com/0xsharma/compact-chain/util"
)

type Blockchain struct {
	blocks []*Block
	mutex  *sync.RWMutex
}

func NewBlockchain() *Blockchain {
	genesis := CreateGenesisBlock()

	return &Blockchain{[]*Block{genesis}, new(sync.RWMutex)}
}

func (bc *Blockchain) AddBlock(data []byte) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	prevBlock := bc.blocks[len(bc.blocks)-1]
	blockNumber := big.NewInt(0).Add(prevBlock.Number(), big.NewInt(1))
	block := NewBlock(blockNumber, prevBlock.Hash(), data)
	bc.blocks = append(bc.blocks, block)
}

func CreateGenesisBlock() *Block {
	return NewBlock(big.NewInt(0), util.NewHashFromHex("0x0"), []byte("Genesis Block"))
}

func (bc *Blockchain) Current() *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) GetBlockByNumber(b *big.Int) *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return bc.blocks[b.Int64()]
}
