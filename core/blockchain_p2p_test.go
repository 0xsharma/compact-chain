package core

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/p2p"
	"github.com/0xsharma/compact-chain/protos"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
	"github.com/stretchr/testify/assert"
)

// nolint : tparallel
func TestP2P(t *testing.T) {
	chainBlocks := []*types.Block{}

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
		P2PPort:          ":6060",
		Mine:             true,
		SignerPrivateKey: util.HexToPrivateKey("e3ddd0f483e2ef1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a6"), // Address = 0x93a63fc45341fc02ac9cce62cc5aeb5c5799403e
		BlockTime:        4,
	}

	chain := NewBlockchain(config)
	if chain.LastBlock.Number.Int64() == 0 {
		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String())
	} else {
		fmt.Println("LastNumber : ", chain.LastBlock.Number, "LastHash : ", chain.LastBlock.DeriveHash().String())
	}

	defer func() {
		chain.RPCServer.HttpServer.Shutdown(context.Background())
		chain.P2PServer.GRPCSrv.Stop()
	}()

	chainBlocks = append(chainBlocks, chain.LastBlock)

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

	// Add block 1 with empty txSet
	time.Sleep(2 * time.Second)
	chain.AddBlock([]byte(fmt.Sprintf("Block %d", chain.LastBlock.Number.Int64()+1)), []*types.Transaction{}, make(chan bool), pkey)
	fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String(), "TxCount", len(chain.LastBlock.Transactions))
	chainBlocks = append(chainBlocks, chain.LastBlock)

	// Create p2p client and connect to p2p grpc server
	conn, client := p2p.ConnectToGRPCServer("localhost" + config.P2PPort)
	defer conn.Close()

	// ASSERTIONS

	// Test LatestBlock
	r, err := client.LatestBlock(context.Background(), &protos.LatestBlockRequest{})
	if err != nil {
		t.Fatal(err)
	}

	rBlock := types.DeserializeBlock(r.EncodedBlock)
	assert.Equal(t, chain.LastBlock.Number, rBlock.Number)
	assert.Equal(t, chain.LastBlock.DeriveHash().String(), rBlock.DeriveHash().String())

	// Test Txpool Pending Transactions
	rTxpool, err := client.TxPoolPending(context.Background(), &protos.TxpoolPendingRequest{})
	if err != nil {
		t.Fatal(err)
	}

	txs := []*types.Transaction{}
	for _, tx := range rTxpool.EncodedTxs {
		txs = append(txs, types.DeserializeTransaction(tx))
	}

	assert.Equal(t, 2, len(txs))
	assert.Equal(t, tx1.Hash().String(), txs[0].Hash().String())
	assert.Equal(t, tx2.Hash().String(), txs[1].Hash().String())

	// Test GetBlocksInRange
	// Add block 2 with tx1 and tx2
	time.Sleep(2 * time.Second)
	chain.AddBlock([]byte(fmt.Sprintf("Block %d", chain.LastBlock.Number.Int64()+1)), chain.Txpool.GetTxs(), make(chan bool), pkey)
	fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String(), "TxCount", len(chain.LastBlock.Transactions))
	chainBlocks = append(chainBlocks, chain.LastBlock)

	rBlocks, err := client.BlocksInRange(context.Background(), &protos.BlocksInRangeRequest{
		StartHeight: 0,
		EndHeight:   2,
	})

	if err != nil {
		t.Fatal(err)
	}

	blocksInRange := []*types.Block{}
	for _, block := range rBlocks.EncodedBlocks {
		blocksInRange = append(blocksInRange, types.DeserializeBlock(block))
	}

	assert.Equal(t, 3, len(blocksInRange))
	assert.Equal(t, chainBlocks[2].DeriveHash(), blocksInRange[2].DeriveHash())
	assert.Equal(t, chainBlocks[1].DeriveHash(), blocksInRange[1].DeriveHash())
	assert.Equal(t, chainBlocks[0].DeriveHash(), blocksInRange[0].DeriveHash())
}
