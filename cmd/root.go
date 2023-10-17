package cmd

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/core"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
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
			nodeID, _ := strconv.ParseInt(args[0], 10, 0)
			startBlockchainNode(nodeID)
		},
	}

	demoCmd = &cobra.Command{
		Use:   "demo",
		Short: "Demo the Compact-Chain node",
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
	rootCmd.AddCommand(demoCmd)
}

var (
	homePath, _ = os.UserHomeDir()
	dbPath      = homePath + "/.compact-chain/db"
	stateDbPath = homePath + "/.compact-chain/statedb"
)

func demoBlockchain() {
	config := &config.Config{
		ConsensusDifficulty: 16,
		ConsensusName:       "pow",
		DBDir:               dbPath + "demo",
		StateDBDir:          stateDbPath + "demo",
		MinFee:              big.NewInt(100),
		RPCPort:             ":1711",
		BalanceAlloc:        map[string]*big.Int{},
		P2PPort:             ":6060",
		Peers:               []string{"localhost:6061"},
		BlockTime:           2,
		SignerPrivateKey:    util.HexToPrivateKey("c3fc038a9abc0f483e2e1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a1"),
	}

	chain := core.NewBlockchain(config)
	if chain.LastBlock.Number.Int64() == 0 {
		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String())
	} else {
		fmt.Println("LastNumber : ", chain.LastBlock.Number, "LastHash : ", chain.LastBlock.DeriveHash().String())
	}

	lastNumber := chain.LastBlock.Number

	for i := lastNumber.Int64() + 1; i <= lastNumber.Int64()+10; i++ {
		time.Sleep(2 * time.Second)

		err := chain.AddBlock([]byte(fmt.Sprintf("Block %d", i)), []*types.Transaction{}, make(chan bool), config.SignerPrivateKey)
		if err != nil {
			fmt.Println("Error Adding Block", err)
		}

		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String())
	}
}

func startBlockchainNode(nodeId int64) {
	fmt.Println("Starting node", nodeId)

	config := &config.Config{
		ConsensusDifficulty: 20,
		ConsensusName:       "pow",
		DBDir:               dbPath + fmt.Sprint(nodeId),
		StateDBDir:          stateDbPath + fmt.Sprint(nodeId),
		MinFee:              big.NewInt(100),
		RPCPort:             ":1711" + fmt.Sprint(nodeId),
		BalanceAlloc:        map[string]*big.Int{},
		P2PPort:             ":6060" + fmt.Sprint(nodeId),
		Peers:               []string{"localhost:60601", "localhost:60602", "localhost:60603"},
		BlockTime:           4,
		SignerPrivateKey:    util.HexToPrivateKey("c3fc038a9abc0f483e2e1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a" + fmt.Sprint(nodeId)),
	}

	core.StartBlockchain(config)
}
