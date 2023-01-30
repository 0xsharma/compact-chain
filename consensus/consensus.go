package consensus

import (
	"math/big"

	"github.com/0xsharma/compact-chain/types"
)

type Consensus interface {
	GetDifficulty() *big.Int
	SetDifficulty(d *big.Int)
	GetTarget() *big.Int
	GetTargetHex() string
	GetTargetBytes() []byte
	Mine(b *types.Block) *types.Block
}
