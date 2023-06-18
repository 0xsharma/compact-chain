package types

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/0xsharma/compact-chain/util"
)

type Transactions []*Transaction

func (txs Transactions) Array() []*Transaction {
	return txs
}

type Transaction struct {
	From  util.Address
	To    util.Address
	Value *big.Int
	Msg   []byte
	Fee   *big.Int
	Nonce *big.Int
	R     *big.Int
	S     *big.Int
}

func (tx *Transaction) Hash() *util.Hash {
	txHash := bytes.Join([][]byte{tx.From.Bytes(), tx.To.Bytes(), tx.Value.Bytes(), tx.Msg, tx.Fee.Bytes(), tx.Nonce.Bytes()}, []byte{})

	return util.NewHash(txHash)
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
