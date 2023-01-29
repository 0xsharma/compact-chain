package main

import (
	"fmt"
	"math/big"

	"github.com/0xsharma/compact-chain/core"
	"github.com/0xsharma/compact-chain/util"
)

func main() {
	block := &core.Block{
		Number:     big.NewInt(1),
		Hash:       &util.Hash{},
		ParentHash: &util.Hash{},
		Data:       []byte("Hello World"),
		Miner:      &util.Address{},
	}

	fmt.Printf("%+v", block)
}
