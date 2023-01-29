package core

import (
	"bytes"
	"math/big"

	"github.com/0xsharma/compact-chain/util"
)

type Block struct {
	number     *big.Int
	hash       *util.Hash
	parentHash *util.Hash
	data       []byte
}

func NewBlock(number *big.Int, parentHash *util.Hash, data []byte) *Block {
	block := &Block{
		number:     number,
		parentHash: parentHash,
		data:       data,
	}

	block.hash = block.DeriveHash()

	return block
}

func (b *Block) DeriveHash() *util.Hash {
	blockHash := bytes.Join([][]byte{b.number.Bytes(), b.parentHash.Bytes(), b.data}, []byte{})

	return util.NewHash(blockHash)
}

func (b *Block) Number() *big.Int {
	return b.number
}

func (b *Block) Hash() *util.Hash {
	return b.hash
}

func (b *Block) ParentHash() *util.Hash {
	return b.parentHash
}

func (b *Block) Data() []byte {
	return b.data
}
