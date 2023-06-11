package txpool

import (
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
)

type Empty struct{}

func (tp *TxPool) AddTx_RPC(args *types.Transaction, reply *types.RPCResponse) error {
	tp.AddTx(args)
	*reply = types.RPCResponse{Success: true}
	return nil
}

func (tp *TxPool) GetTxs_RPC(_ *Empty, reply *types.RPCResponse) error {
	txs := tp.GetTxs()
	responseBytes := util.EncodeToBytes(txs)

	*reply = types.RPCResponse{Success: true, Message: responseBytes}
	return nil
}
