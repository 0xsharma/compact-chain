package config

import (
	"crypto/ecdsa"
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
	SignerPrivateKey    *ecdsa.PrivateKey
	Mine                bool
	BalanceAlloc        map[string]*big.Int
	P2PPort             string
	Peers               []string
}

func DefaultConfig() *Config {
	cfg := &Config{
		ConsensusDifficulty: 16,
		ConsensusName:       "pow",
		DBDir:               dbPath,
		StateDBDir:          stateDbPath,
		MinFee:              big.NewInt(100),
		RPCPort:             ":1711",
		P2PPort:             ":6060",
	}

	return cfg
}
