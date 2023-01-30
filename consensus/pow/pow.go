package pow

import (
	"math/big"

	"github.com/0xsharma/compact-chain/types"
)

type POW struct {
	difficulty *big.Int
}

func NewPOW(difficulty int) *POW {
	return &POW{
		difficulty: big.NewInt(int64(difficulty)),
	}
}

func (c *POW) GetDifficulty() *big.Int {
	return c.difficulty
}

func (c *POW) SetDifficulty(d *big.Int) {
	c.difficulty = d
}

func (c *POW) GetTarget() *big.Int {
	target := big.NewInt(1)
	return target.Lsh(target, uint(256-c.difficulty.Int64()))
}

func (c *POW) GetTargetHex() string {
	return c.GetTarget().Text(16)
}

func (c *POW) GetTargetBytes() []byte {
	return c.GetTarget().Bytes()
}

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
