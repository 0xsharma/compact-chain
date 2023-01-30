package types

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
	nonce      *big.Int
}

func NewBlock(number *big.Int, parentHash *util.Hash, data []byte) *Block {
	block := &Block{
		number:     number,
		parentHash: parentHash,
		data:       data,
		nonce:      big.NewInt(0),
	}

	block.hash = block.DeriveHash()

	return block
}

func (dst *Block) Clone(src *Block) {
	dst.number = src.number
	dst.hash = src.hash
	dst.parentHash = src.parentHash
	dst.data = src.data
	dst.nonce = src.nonce
}

func (b *Block) DeriveHash() *util.Hash {
	blockHash := bytes.Join([][]byte{b.number.Bytes(), b.parentHash.Bytes(), b.data, b.nonce.Bytes()}, []byte{})

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

func (b *Block) SetNonce(n *big.Int) {
	b.nonce = n
}

func (b *Block) SetHash(h *util.Hash) {
	b.hash = h
}
