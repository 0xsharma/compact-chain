package config

import "math/big"

// Config is the configuration for the Compact-Chain node.
type Config struct {
	ConsensusDifficulty int
	ConsensusName       string
	DBDir               string
	StateDBDir          string
	MinFee              *big.Int
}
