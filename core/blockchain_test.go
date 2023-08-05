package core

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/dbstore"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
	"github.com/stretchr/testify/assert"
)

func TestBlockchainStateBalance(t *testing.T) {
	t.Parallel()

	to1 := util.BytesToAddress([]byte{0x01})
	to2 := util.BytesToAddress([]byte{0x02})

	config := &config.Config{
		ConsensusDifficulty: 16,
		ConsensusName:       "pow",
		DBDir:               t.TempDir(),
		StateDBDir:          t.TempDir(),
		MinFee:              big.NewInt(100),
		RPCPort:             ":1711",
		BalanceAlloc: map[string]*big.Int{
			"0xa52c981eee8687b5e4afd69aa5006548c24d7685": big.NewInt(1000000000000000000), // Allocating funds to 0xa52c981eee8687b5e4afd69aa5006548c24d7685
		},
		Mine:             true,
		SignerPrivateKey: util.HexToPrivateKey("e3ddd0f483e2ef1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a6"), // Address = 0x93a63fc45341fc02ac9cce62cc5aeb5c5799403e
	}

	chain := NewBlockchain(config)
	if chain.LastBlock.Number.Int64() == 0 {
		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String())
	} else {
		fmt.Println("LastNumber : ", chain.LastBlock.Number, "LastHash : ", chain.LastBlock.DeriveHash().String())
	}

	pkey := util.HexToPrivateKey("c3fc038a9abc0f483e2e1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a6") // Address = 0xa52c981eee8687b5e4afd69aa5006548c24d7685
	ua := util.NewUnlockedAccount(pkey)

	// tx1
	tx1 := newTransaction(t, ua.Address().Bytes(), to1.Bytes(), "hello", 200, 1000, 0)
	tx1.Sign(ua)

	// tx2
	tx2 := newTransaction(t, ua.Address().Bytes(), to2.Bytes(), "hello", 100, 2000, 1)
	tx2.Sign(ua)

	// Add tx1 and tx2 to txpool
	chain.Txpool.AddTxs([]*types.Transaction{tx1, tx2})

	// Add block 1
	time.Sleep(2 * time.Second)
	chain.AddBlock([]byte(fmt.Sprintf("Block %d", chain.LastBlock.Number.Int64()+1)), []*types.Transaction{})

	fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String(), "TxCount", len(chain.LastBlock.Transactions))

	//Add block 2
	time.Sleep(2 * time.Second)
	chain.AddBlock([]byte(fmt.Sprintf("Block %d", chain.LastBlock.Number.Int64()+1)), chain.Txpool.GetTxs())

	fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String(), "TxCount", len(chain.LastBlock.Transactions))

	// Assertions
	balanceSender, err := chain.StateDB.DB.Get(dbstore.PrefixKey(dbstore.BalanceKey, ua.Address().String()))
	if err != nil {
		t.Fatal(err)
	}

	balanceSenderBig := new(big.Int).SetBytes(balanceSender)
	assert.Equal(t, big.NewInt(999999999999997000), balanceSenderBig)

	balanceTo1, err := chain.StateDB.DB.Get(dbstore.PrefixKey(dbstore.BalanceKey, to1.String()))
	if err != nil {
		t.Fatal(err)
	}

	balanceTo1Big := new(big.Int).SetBytes(balanceTo1)
	assert.Equal(t, big.NewInt(1000), balanceTo1Big)

	balanceTo2, err := chain.StateDB.DB.Get(dbstore.PrefixKey(dbstore.BalanceKey, to2.String()))
	if err != nil {
		t.Fatal(err)
	}

	balanceTo2Big := new(big.Int).SetBytes(balanceTo2)
	assert.Equal(t, big.NewInt(2000), balanceTo2Big)

	addressMiner := util.NewUnlockedAccount(config.SignerPrivateKey)

	balanceMiner, err := chain.StateDB.DB.Get(dbstore.PrefixKey(dbstore.BalanceKey, addressMiner.Address().String()))
	if err != nil {
		t.Fatal(err)
	}

	balanceMinerBig := new(big.Int).SetBytes(balanceMiner)
	assert.Equal(t, big.NewInt(300), balanceMinerBig)
}

func newTransaction(t *testing.T, from, to []byte, msg string, fee, value int64, nonce int64) *types.Transaction {
	t.Helper()

	return &types.Transaction{
		From:  *util.BytesToAddress(from),
		To:    *util.BytesToAddress(to),
		Msg:   []byte(msg),
		Fee:   big.NewInt(fee),
		Value: big.NewInt(value),
		Nonce: big.NewInt(nonce),
	}
}
