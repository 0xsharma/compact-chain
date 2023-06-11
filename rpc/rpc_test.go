package rpc

import (
	"log"
	"math/big"
	"net/rpc"
	"testing"
	"time"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/txpool"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
	"github.com/stretchr/testify/assert"
)

var empty struct{}

func TestTxpoolRPC(t *testing.T) {
	t.Parallel()
	rpcPort := ":1711"

	txpool := txpool.NewTxPool(config.DefaultConfig().MinFee)
	NewRPCServer(rpcPort, &RPCDomains{TxPool: txpool})
	time.Sleep(2 * time.Second)

	// Send add Transacation request 1
	tx1 := newTransaction(t, []byte{0x01}, []byte{0x02}, "hello", 100, 100)
	res, err := SendRpcRequest(t, "TxPool.AddTx_RPC", tx1, rpcPort)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := res.(types.RPCResponse); !ok {
		t.Fatal("expected RPCResponse", "got", res)
	}

	// Send add Transacation request 2
	tx2 := newTransaction(t, []byte{0x02}, []byte{0x03}, "hello1", 101, 101)
	res, err = SendRpcRequest(t, "TxPool.AddTx_RPC", tx2, rpcPort)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := res.(types.RPCResponse); !ok {
		t.Fatal("expected RPCResponse", "got", res)
	}

	// Send get Transactions request
	res, err = SendRpcRequest(t, "TxPool.GetTxs_RPC", empty, rpcPort)
	if err != nil {
		t.Fatal(err)
	}

	var txs *types.Transactions
	if resTxs, ok := res.(types.RPCResponse); !ok {
		t.Fatal("expected RPCResponse", "got", res)
	} else {
		txs, err = util.DecodeFromBytes[types.Transactions](resTxs.Message)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Tx2 Fee > Tx1 Fee ; its index will be higher
	// Assert Transactions That we sent to pool and received from pool are same
	assert.Equal(t, tx2, txs.Array()[0])
	assert.Equal(t, tx1, txs.Array()[1])

}

func SendRpcRequest(t *testing.T, method string, params interface{}, addr string) (interface{}, error) {
	t.Helper()
	hostname := "localhost"

	client, err := rpc.DialHTTP("tcp", hostname+addr)
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	var reply types.RPCResponse

	err = client.Call(method, params, &reply)
	if err != nil {
		log.Fatal("error: ", err)
	}

	return reply, nil
}

func newTransaction(t *testing.T, from, to []byte, msg string, fee, value int64) *types.Transaction {
	t.Helper()
	return &types.Transaction{
		From:  *util.NewAddress(from),
		To:    *util.NewAddress(to),
		Msg:   []byte(msg),
		Fee:   big.NewInt(fee),
		Value: big.NewInt(value),
	}
}
