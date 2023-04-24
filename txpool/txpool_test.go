package txpool

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/0xsharma/compact-chain/core"
	"github.com/0xsharma/compact-chain/util"
	"github.com/stretchr/testify/assert"
)

func NewRandomTx(t *testing.T) *core.Transaction {
	r := rand.Int63n(1000)
	return &core.Transaction{From: util.Address{}, To: util.Address{}, Value: big.NewInt(1), Msg: []byte{}, Fee: big.NewInt(r)}
}

func TestTxpoolAdd(t *testing.T) {
	txpool := NewTxPool(big.NewInt(0))

	for i := 0; i < 100; i++ {
		txpool.AddTx(NewRandomTx(t))
	}

	txs := txpool.Transactions

	for i, _ := range txs {
		if i == 0 {
			continue
		}
		assert.Equal(t, true, txs[i-1].Fee.Cmp(txs[i].Fee) >= 0)
	}
}
