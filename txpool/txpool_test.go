package txpool

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
	"github.com/stretchr/testify/assert"
)

func NewRandomTx(t *testing.T) *types.Transaction {
	t.Helper()

	// nolint : gosec
	r := rand.Int63n(1000)

	return &types.Transaction{From: util.Address{}, To: util.Address{}, Value: big.NewInt(1), Msg: []byte{}, Fee: big.NewInt(r)}
}

func TestTxpoolAdd(t *testing.T) {
	t.Parallel()

	txpool := NewTxPool(big.NewInt(0), nil, nil)

	for i := 0; i < 100; i++ {
		txpool.AddTx(NewRandomTx(t))
	}

	txs := txpool.Transactions

	for i := range txs {
		if i == 0 {
			continue
		}
		// Check fee of txs are in desc order in txpool
		assert.Equal(t, true, txs[i-1].Fee.Cmp(txs[i].Fee) >= 0)
	}
}
