package types

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/0xsharma/compact-chain/util"
	"github.com/cbergoon/merkletree"
)

// Block is the basic unit of the blockchain.
type Block struct {
	Number       *big.Int
	ParentHash   *util.Hash
	ExtraData    []byte
	Nonce        *big.Int
	Transactions []*Transaction
	TxRoot       *util.Hash
}

// NewBlock creates a new block and sets the hash.
func NewBlock(number *big.Int, parentHash *util.Hash, data []byte) *Block {
	block := &Block{
		Number:     number,
		ParentHash: parentHash,
		ExtraData:  data,
		Nonce:      big.NewInt(0),
	}

	return block
}

// Clone returns a duplicate block from the source block.
func (dst *Block) Clone(src *Block) {
	dst.Number = src.Number
	dst.ParentHash = src.ParentHash
	dst.ExtraData = src.ExtraData
	dst.Nonce = src.Nonce
}

// DeriveHash derives the hash of the block.
func (b *Block) DeriveHash() *util.Hash {
	blockHash := bytes.Join([][]byte{b.Number.Bytes(), b.ParentHash.Bytes(), b.ExtraData, b.Nonce.Bytes(), b.TxRootHash().Bytes()}, []byte{})

	return util.HashData(blockHash)
}

func (b *Block) TxRootHash() *util.Hash {
	if len(b.Transactions) == 0 {
		return util.HashData([]byte{})
	} else {
		var list []merkletree.Content
		for _, tx := range b.Transactions {
			list = append(list, tx)
		}

		//Create a new Merkle Tree from the list of Content
		t, err := merkletree.NewTree(list)
		if err != nil {
			panic(err)
		}
		mr := t.MerkleRoot()
		return util.ByteToHash(mr)
	}
}

// SetNonce sets the nonce of the block.
func (b *Block) SetNonce(n *big.Int) {
	b.Nonce = n
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
