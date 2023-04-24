package txpool

import (
	"math/big"
	"sort"

	"github.com/0xsharma/compact-chain/core"
)

type TxPool struct {
	Transactions []*core.Transaction
}

func NewTxPool(minFee *big.Int) *TxPool {
	return &TxPool{}
}

func intToBool(n int) bool {
	return n >= 0
}

func (tp *TxPool) AddTx(tx *core.Transaction) {
	txs := append(tp.Transactions, tx)
	sort.Slice(txs, func(i, j int) bool {
		return intToBool(txs[i].Fee.Cmp(txs[j].Fee))
	})
	tp.Transactions = txs
}

func (tp *TxPool) AddTxs(txs []*core.Transaction) {
	txpoolTxs := append(tp.Transactions, txs...)
	sort.Slice(txpoolTxs, func(i, j int) bool {
		return intToBool(txpoolTxs[i].Fee.Cmp(txpoolTxs[j].Fee))
	})
	tp.Transactions = txpoolTxs
}
