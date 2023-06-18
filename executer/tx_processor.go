package executer

import (
	"errors"
	"math/big"

	"github.com/0xsharma/compact-chain/dbstore"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
)

var (
	ErrInvalidTransaction = errors.New("invalid transaction")
)

type TxProcessor struct {
	MinFee *big.Int
	State  *dbstore.DB
	Signer *util.Address
}

func NewTxProcessor(state *dbstore.DB, minFee *big.Int, signer *util.Address) *TxProcessor {
	return &TxProcessor{
		MinFee: minFee,
		State:  state,
		Signer: signer,
	}
}

func (txp *TxProcessor) IsValid(tx *types.Transaction) bool {
	from := tx.From
	balance, err := txp.State.Get(dbstore.PrefixKey(dbstore.BalanceKey, from.String()))

	if err != nil {
		return false
	}

	balanceBig := new(big.Int).SetBytes(balance)

	// Add Fee to Value
	totalValue := big.NewInt(0).Add(tx.Value, tx.Fee)

	if balanceBig.Cmp(totalValue) < 0 {
		return false
	}

	var nonceBig *big.Int

	nonce, err := txp.State.Get(dbstore.PrefixKey(dbstore.NonceKey, from.String()))
	if err != nil {
		nonceBig = big.NewInt(0)
	} else {
		nonceBig = new(big.Int).SetBytes(nonce)
	}

	if big.NewInt(0).Sub(tx.Nonce, nonceBig).Cmp(big.NewInt(1)) != 0 {
		return false
	}

	return true
}

// ProcessTx processes a transaction and returns the transaction fee.
func (txp *TxProcessor) ProcessTx(tx *types.Transaction) error {
	if !txp.IsValid(tx) {
		return ErrInvalidTransaction
	}

	from := tx.From
	to := tx.To
	value := tx.Value

	dbBatch := txp.State.NewBatch()

	// Get sender balance.
	senderBalance, err := txp.State.Get(dbstore.PrefixKey(dbstore.BalanceKey, from.String()))
	if err != nil {
		return err
	}

	sendBalanceBig := new(big.Int).SetBytes(senderBalance)

	// Get receiver balance.
	receiverBalance, err := txp.State.Get(dbstore.PrefixKey(dbstore.BalanceKey, to.String()))
	if err != nil {
		return err
	}

	receiverBalanceBig := new(big.Int).SetBytes(receiverBalance)

	// Update sender balance.
	sendBalanceBig.Sub(sendBalanceBig, value)
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BalanceKey, from.String())), sendBalanceBig.Bytes())

	// Update receiver balance.
	receiverBalanceBig.Add(receiverBalanceBig, value)
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BalanceKey, to.String())), receiverBalanceBig.Bytes())

	// Update sender nonce.
	nonce, err := txp.State.Get(dbstore.PrefixKey(dbstore.NonceKey, from.String()))
	if err != nil {
		return err
	}

	nonceBig := new(big.Int).SetBytes(nonce)
	nonceBig.Add(nonceBig, big.NewInt(1))
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.NonceKey, from.String())), nonceBig.Bytes())

	// Update Miner Fee.
	dbBatch.Put([]byte(dbstore.PrefixKey(dbstore.BalanceKey, txp.Signer.String())), tx.Fee.Bytes())

	// Commit batch to db
	err = txp.State.WriteBatch(dbBatch)
	if err != nil {
		panic(err)
	}

	return nil
}
