package pow

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/0xsharma/compact-chain/executer"
	"github.com/0xsharma/compact-chain/types"
)

type POW struct {
	difficulty  *big.Int
	TxProcessor *executer.TxProcessor
}

// NewPOW creates a new proof of work consensus.
func NewPOW(difficulty int, txProcessor *executer.TxProcessor) *POW {
	return &POW{
		difficulty:  big.NewInt(int64(difficulty)),
		TxProcessor: txProcessor,
	}
}

// GetDifficulty returns the difficulty of the proof of work consensus.
func (c *POW) GetDifficulty() *big.Int {
	return c.difficulty
}

// SetDifficulty sets the difficulty of the proof of work consensus.
func (c *POW) SetDifficulty(d *big.Int) {
	c.difficulty = d
}

// GetTarget returns the target of the proof of work consensus.
func (c *POW) GetTarget() *big.Int {
	target := big.NewInt(1)
	return target.Lsh(target, uint(256-c.difficulty.Int64()))
}

// Mine mines the block with the proof of work consensus with the given difficulty.
func (c *POW) Mine(b *types.Block, mineInterrupt chan bool) *types.Block {
	nonce := big.NewInt(0)

	validTxs := []*types.Transaction{}

	for _, tx := range b.Transactions {
		if c.TxProcessor.IsValid(tx) {
			err := c.TxProcessor.ProcessTx(tx)
			if err == nil {
				validTxs = append(validTxs, tx)
			} else {
				fmt.Println("Failed to execute Tx :", "tx :", tx, "error", err)
			}
		} else {
			fmt.Println("Invalid Tx :", "tx :", tx)
		}
	}

	b.Transactions = validTxs

	for {
		select {
		case <-mineInterrupt:
			return nil

		default:
			b.SetNonce(nonce)
			hash := b.DeriveHash()

			hashBytes := hash.Bytes()
			hashBig := new(big.Int).SetBytes(hashBytes)

			if hashBig.Cmp(c.GetTarget()) < 0 {
				return b
			}

			n, _ := rand.Int(rand.Reader, big.NewInt(100))
			nonce.Add(nonce, big.NewInt(n.Int64()))
		}
	}
}

// Validate validates the block with the proof of work consensus.
func (c *POW) Validate(b *types.Block) bool {
	validTxs := []*types.Transaction{}

	for _, tx := range b.Transactions {
		if c.TxProcessor.IsValid(tx) {
			err := c.TxProcessor.ProcessTx(tx)
			if err == nil {
				validTxs = append(validTxs, tx)
			} else {
				fmt.Println("Failed to execute Tx :", "tx :", tx, "error", err)
				return false
			}
		} else {
			fmt.Println("Invalid Tx :", "tx :", tx)
			return false
		}
	}

	b.Transactions = validTxs

	hash := b.DeriveHash()
	hashBytes := hash.Bytes()
	hashBig := new(big.Int).SetBytes(hashBytes)

	if hashBig.Cmp(c.GetTarget()) > 0 {
		fmt.Println("Invalid block hash for POW :", hashBig, "target :", c.GetTarget())
		return false
	}

	return true
}
