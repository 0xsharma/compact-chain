package cmd

import (
	"fmt"
	"math/big"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/core"
	"github.com/spf13/cobra"
)

var (
	version = "v0.0.1"
	rootCmd = &cobra.Command{}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Compact-Chain",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the Compact-Chain node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting Compact-Chain node\n\n")
			demoBlockchain()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
}

func demoBlockchain() {
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
