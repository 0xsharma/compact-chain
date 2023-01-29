package core

import (
	"math/big"

	"github.com/0xsharma/compact-chain/util"
)

type Block struct {
	Number     *big.Int
	Hash       *util.Hash
	ParentHash *util.Hash
	Data       []byte
	Miner      *util.Address
}
