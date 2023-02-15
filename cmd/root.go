package cmd

import (
	"fmt"
	"os"
	"time"

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

var (
	homePath, _ = os.UserHomeDir()
	dbPath      = homePath + "/.compact-chain/db"
)

func demoBlockchain() {
	config := &config.Config{
		ConsensusDifficulty: 16,
		ConsensusName:       "pow",
		DBDir:               dbPath,
	}

	chain := core.NewBlockchain(config)
	if chain.LastBlock.Number.Int64() == 0 {
		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.Hash.String())
	} else {
		fmt.Println("LastNumber : ", chain.LastBlock.Number, "LastHash : ", chain.LastBlock.Hash.String())
	}

	lastNumber := chain.LastBlock.Number

	for i := lastNumber.Int64() + 1; i <= lastNumber.Int64()+10; i++ {
		time.Sleep(2 * time.Second)
		chain.AddBlock([]byte(fmt.Sprintf("Block %d", i)))
		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.Hash.String())
	}
}
