package types

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/0xsharma/compact-chain/util"
)

// Block is the basic unit of the blockchain.
type Block struct {
	Number     *big.Int
	Hash       *util.Hash
	ParentHash *util.Hash
	Data       []byte
	Nonce      *big.Int
}

// NewBlock creates a new block and sets the hash.
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

// Clone returns a duplicate block from the source block.
func (dst *Block) Clone(src *Block) {
	dst.Number = src.Number
	dst.Hash = src.Hash
	dst.ParentHash = src.ParentHash
	dst.Data = src.Data
	dst.Nonce = src.Nonce
}

// DeriveHash derives the hash of the block.
func (b *Block) DeriveHash() *util.Hash {
	blockHash := bytes.Join([][]byte{b.Number.Bytes(), b.ParentHash.Bytes(), b.Data, b.Nonce.Bytes()}, []byte{})

	return util.NewHash(blockHash)
}

// SetNonce sets the nonce of the block.
func (b *Block) SetNonce(n *big.Int) {
	b.Nonce = n
}

// SetHash sets the hash of the block.
func (b *Block) SetHash(h *util.Hash) {
	b.Hash = h
}

// Serialize serializes the block object into bytes.
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	if err != nil {
		panic(err)
	}

	return res.Bytes()
}

// DeserializeBlock deserializes the block bytes into a block object.
func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	if err != nil {
		panic(err)
	}

	return &block
}
