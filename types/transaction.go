package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/gob"
	"math/big"

	"github.com/0xsharma/compact-chain/util"
	"github.com/cbergoon/merkletree"
)

type Transactions []*Transaction

func (txs Transactions) Array() []*Transaction {
	return txs
}

type Transaction struct {
	From      util.Address
	To        util.Address
	Value     *big.Int
	Msg       []byte
	Fee       *big.Int
	Nonce     *big.Int
	R         *big.Int
	S         *big.Int
	PublicKey *util.CompactPublicKey
}

func (tx *Transaction) Hash() *util.Hash {
	txHash := bytes.Join([][]byte{tx.From.Bytes(), tx.To.Bytes(), tx.Value.Bytes(), tx.Msg, tx.Fee.Bytes(), tx.Nonce.Bytes()}, []byte{})

	return util.HashData(txHash)
}

func (tx *Transaction) CalculateHash() ([]byte, error) {
	return tx.Hash().Bytes(), nil
}

// Equals tests for equality of two Contents
func (tx *Transaction) Equals(other merkletree.Content) (bool, error) {
	OtherFrom := other.(*Transaction).From
	OtherTo := other.(*Transaction).To
	OtherValue := other.(*Transaction).Value.Bytes()
	OtherMsg := other.(*Transaction).Msg
	OtherFee := other.(*Transaction).Fee.Bytes()
	OtherNonce := other.(*Transaction).Nonce.Bytes()
	OtherR := other.(*Transaction).R.Bytes()
	OtherS := other.(*Transaction).S.Bytes()
	OtherPublicKey := other.(*Transaction).PublicKey

	out := tx.From == OtherFrom && tx.To == OtherTo && bytes.Equal(tx.Value.Bytes(), OtherValue) && bytes.Equal(tx.Msg, OtherMsg) && bytes.Equal(tx.Fee.Bytes(), OtherFee) && bytes.Equal(tx.Nonce.Bytes(), OtherNonce) && bytes.Equal(tx.R.Bytes(), OtherR) && bytes.Equal(tx.S.Bytes(), OtherS) && tx.PublicKey == OtherPublicKey

	return out, nil
}

func (tx *Transaction) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(tx)

	if err != nil {
		panic(err)
	}

	return res.Bytes()
}

func DeserializeTransaction(data []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&tx)
	if err != nil {
		panic(err)
	}

	return &tx
}

func (tx *Transaction) Sign(ua *util.UnlockedAccount) {
	r, s, err := ua.Sign(tx.Hash().Bytes())
	if err != nil {
		panic(err)
	}

	tx.R = r
	tx.S = s
}

func (tx *Transaction) Verify() bool {
	pubKey := tx.PublicKey.PublicKey()
	return ecdsa.Verify(pubKey, tx.Hash().Bytes(), tx.R, tx.S)
}
