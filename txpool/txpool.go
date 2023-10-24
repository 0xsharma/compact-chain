package txpool

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/0xsharma/compact-chain/dbstore"
	"github.com/0xsharma/compact-chain/types"
	"github.com/golang/groupcache/lru"
)

var (
	ErrInvalidTransaction = errors.New("invalid transaction")
)

type TxPool struct {
	MinFee       *big.Int
	State        *dbstore.DB
	Transactions []*types.Transaction

	TxPoolCh chan *types.Transaction

	LatestIncludedTxs *lru.Cache
}

func NewTxPool(minFee *big.Int, db *dbstore.DB, txpoolCh chan *types.Transaction) *TxPool {
	if db == nil {
		fmt.Println("DB is nil, running in mock mode for tests")
	}

	txpool := &TxPool{
		MinFee:            minFee,
		State:             db,
		TxPoolCh:          txpoolCh,
		LatestIncludedTxs: lru.New(1000),
	}

	go txpool.loop()

	return txpool
}

func intToBool(n int) bool {
	return n >= 0
}

func (txp *TxPool) loop() {
	// nolint : gosimple
	for {
		select {
		case tx := <-txp.TxPoolCh:
			txp.AddTx(tx)
		}
	}
}

func (txp *TxPool) IsValid(tx *types.Transaction) bool {
	if txp.State == nil {
		return true
	}

	if tx.Fee.Cmp(txp.MinFee) < 0 {
		return false
	}

	from := tx.From

	signOk := tx.Verify()
	if !signOk {
		return false
	}

	balance, err := txp.State.Get(dbstore.PrefixKey(dbstore.BalanceKey, from.String()))
	if err != nil {
		return false
	}

	balanceBig := new(big.Int).SetBytes(balance)

	// Add Fee to Value
	totalValue := big.NewInt(0).Add(tx.Value, tx.Fee)

	// nolint : gosimple
	if balanceBig.Cmp(totalValue) < 0 {
		return false
	}

	// TODO : Write nonce logic in txpool and enable this check

	// var nonceBig *big.Int

	// nonce, err := txp.State.Get(dbstore.PrefixKey(dbstore.NonceKey, from.String()))
	// if err != nil {
	// 	nonceBig = big.NewInt(-1)
	// } else {
	// 	nonceBig = new(big.Int).SetBytes(nonce)
	// }

	// if big.NewInt(0).Sub(tx.Nonce, nonceBig).Cmp(big.NewInt(1)) != 0 {
	// 	return false
	// }

	return true
}

func (tp *TxPool) AddTx(tx *types.Transaction) {
	if !tp.IsValid(tx) {
		return
	}

	_, ok := tp.LatestIncludedTxs.Get(tx.Hash().String())
	if ok {
		return
	}

	for _, tx2 := range tp.Transactions {
		if tx2.Hash().String() == tx.Hash().String() {
			return
		}
	}

	txs := append(tp.Transactions, tx)
	sort.Slice(txs, func(i, j int) bool {
		return intToBool(txs[i].Fee.Cmp(txs[j].Fee))
	})

	tp.Transactions = txs
}

func (tp *TxPool) AddTxs(txs []*types.Transaction) {
	validTxs := make([]*types.Transaction, 0, len(txs))

	for _, tx := range txs {
		_, ok := tp.LatestIncludedTxs.Get(tx.Hash().String())
		if ok {
			continue
		}

		for _, tx2 := range tp.Transactions {
			if tx2.Hash().String() == tx.Hash().String() {
				return
			}
		}

		if tp.IsValid(tx) {
			validTxs = append(validTxs, tx)
		}
	}

	txpoolTxs := append(tp.Transactions, validTxs...)
	sort.Slice(txpoolTxs, func(i, j int) bool {
		return intToBool(txpoolTxs[i].Fee.Cmp(txpoolTxs[j].Fee))
	})

	tp.Transactions = txpoolTxs
}

// remove transaction from txpool
func (tp *TxPool) RemoveTx(tx *types.Transaction) error {
	for i, tx2 := range tp.Transactions {
		if tx2.Hash().String() == tx.Hash().String() {
			tp.Transactions = append(tp.Transactions[:i], tp.Transactions[i+1:]...)
			return nil
		}
	}

	return errors.New("transaction not found")
}

func (tp *TxPool) GetTxs() []*types.Transaction {
	txs := tp.Transactions
	for _, tx := range txs {
		tp.LatestIncludedTxs.Add(tx.Hash().String(), []byte{})
	}

	return txs
}
