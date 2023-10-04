package consensus

import (
	"math/big"

	"github.com/0xsharma/compact-chain/types"
)

// Consensus is the interface for consensus algorithms.
type Consensus interface {
	GetDifficulty() *big.Int
	SetDifficulty(d *big.Int)
	GetTarget() *big.Int
	Mine(b *types.Block, mineInterrupt chan bool) *types.Block
	Validate(b *types.Block) bool
}
