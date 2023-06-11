package config

import (
	"math/big"
	"os"
)

var (
	homePath, _ = os.UserHomeDir()
	dbPath      = homePath + "/.compact-chain/db"
	stateDbPath = homePath + "/.compact-chain/statedb"
)

// Config is the configuration for the Compact-Chain node.
type Config struct {
	ConsensusDifficulty int
	ConsensusName       string
	DBDir               string
	StateDBDir          string
	MinFee              *big.Int
	RPCPort             string
}

func DefaultConfig() *Config {
	cfg := &Config{
		ConsensusDifficulty: 16,
		ConsensusName:       "pow",
		DBDir:               dbPath,
		StateDBDir:          stateDbPath,
		MinFee:              big.NewInt(100),
		RPCPort:             ":1711",
	}

	return cfg
}
