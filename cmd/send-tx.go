package cmd

import (
	"fmt"
	"log"
	"math/big"
	"net/rpc"

	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
)

type sendTxConfig struct {
	PrivateKey string
	To         string
	Value      int64
	RPCAddr    string
	Nonce      int64
}

func SendTx(sendTxCfg *sendTxConfig) {
	ua := util.NewUnlockedAccount(util.HexToPrivateKey(sendTxCfg.PrivateKey))
	from := ua.Address()

	tx := &types.Transaction{
		From:  *from,
		To:    *util.StringToAddress(sendTxCfg.To),
		Value: big.NewInt(sendTxCfg.Value),
		Msg:   []byte("hello"),
		Fee:   big.NewInt(1000),
		Nonce: big.NewInt(sendTxCfg.Nonce),
	}
	tx.Sign(ua)

	fmt.Printf("%+v\n", tx)
	res, err := SendRpcRequest("TxPool.AddTx_RPC", tx, sendTxCfg.RPCAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

func SendRpcRequest(method string, params interface{}, rpcAddr string) (interface{}, error) {
	client, err := rpc.DialHTTP("tcp", rpcAddr)
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
