package core

import (
	"bytes"
	"math/big"

	"github.com/0xsharma/compact-chain/util"
)

type Transaction struct {
	From  util.Address
	To    util.Address
	Value *big.Int
	Msg   []byte
	Fee   *big.Int
}

func (tx *Transaction) Hash() *util.Hash {
	txHash := bytes.Join([][]byte{tx.From.Bytes(), tx.To.Bytes(), tx.Value.Bytes(), tx.Msg}, []byte{})

	return util.NewHash(txHash)
}
