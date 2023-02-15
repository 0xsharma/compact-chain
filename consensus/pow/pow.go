package pow

import (
	"math/big"

	"github.com/0xsharma/compact-chain/types"
)

type POW struct {
	difficulty *big.Int
}

// NewPOW creates a new proof of work consensus.
func NewPOW(difficulty int) *POW {
	return &POW{
		difficulty: big.NewInt(int64(difficulty)),
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
func (c *POW) Mine(b *types.Block) *types.Block {
	nonce := big.NewInt(0)

	for {
		b.SetNonce(nonce)
		hash := b.DeriveHash()

		hashBytes := hash.Bytes()
		hashBig := new(big.Int).SetBytes(hashBytes)

		if hashBig.Cmp(c.GetTarget()) < 0 {
			b.SetHash(hash)
			break
		}

		nonce.Add(nonce, big.NewInt(1))
	}

	return b
}
