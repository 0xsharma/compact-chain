package main

import (
	"math/big"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/core"
)

func main() {
	config := &config.Config{
		ConsensusDifficulty: 16,
		ConsensusName:       "pow",
	}

	chain := core.NewBlockchain(config)

	chain.AddBlock([]byte("Block 1"))
	chain.AddBlock([]byte("Block 2"))
	chain.AddBlock([]byte("Block 3"))

	currentNumber := int(chain.Current().Number().Int64())

	for i := 0; i <= currentNumber; i++ {
		block := chain.GetBlockByNumber(big.NewInt(int64(i)))
		println("BlockNumber : ", block.Number().String())
		println("BlockHash : ", block.Hash().String())
		println("ParentHash : ", block.ParentHash().String())
		println("BlockData : ", string(block.Data()))
		println()
	}
}
