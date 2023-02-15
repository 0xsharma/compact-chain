package types

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/0xsharma/compact-chain/util"
)

type Block struct {
	Number     *big.Int
	Hash       *util.Hash
	ParentHash *util.Hash
	Data       []byte
	Nonce      *big.Int
}

func NewBlock(number *big.Int, parentHash *util.Hash, data []byte) *Block {
	block := &Block{
		Number:     number,
		ParentHash: parentHash,
		Data:       data,
		Nonce:      big.NewInt(0),
	}

	block.Hash = block.DeriveHash()

	return block
}

func (dst *Block) Clone(src *Block) {
	dst.Number = src.Number
	dst.Hash = src.Hash
	dst.ParentHash = src.ParentHash
	dst.Data = src.Data
	dst.Nonce = src.Nonce
}

func (b *Block) DeriveHash() *util.Hash {
	blockHash := bytes.Join([][]byte{b.Number.Bytes(), b.ParentHash.Bytes(), b.Data, b.Nonce.Bytes()}, []byte{})

	return util.NewHash(blockHash)
}

// func (b *Block) Number() *big.Int {
// 	return b.number
// }

// func (b *Block) Hash() *util.Hash {
// 	return b.hash
// }

// func (b *Block) ParentHash() *util.Hash {
// 	return b.parentHash
// }

// func (b *Block) Data() []byte {
// 	return b.data
// }

func (b *Block) SetNonce(n *big.Int) {
	b.Nonce = n
}

func (b *Block) SetHash(h *util.Hash) {
	b.Hash = h
}

// func (b *Block) Nonce() *big.Int {
// 	return b.nonce
// }

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	if err != nil {
		panic(err)
	}

	return res.Bytes()
}

func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	if err != nil {
		panic(err)
	}

	return &block
}
