package p2p

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/0xsharma/compact-chain/dbstore"
	"github.com/0xsharma/compact-chain/protos"
	"github.com/0xsharma/compact-chain/txpool"
	"github.com/0xsharma/compact-chain/types"
	"google.golang.org/grpc"
)

const defaultP2pPort = ":6060"

type P2PServer struct {
	Port                  string
	Lis                   net.Listener
	GRPCSrv               *grpc.Server
	Peers                 []string
	P2PAddrBlockNumberMap map[string]int
	Downloader            *Downloader

	Peersmu               sync.Mutex
	P2PAddrBroadcastMapmu sync.Mutex

	BlockchainDB *dbstore.BlockchainDB
	StateDB      *dbstore.StateDB
	Txpool       *txpool.TxPool

	protos.UnimplementedP2PServer
}

type P2PMessage struct {
	From    string
	Query   string
	Message string
	Error   error
}

func NewServer(port string, initPeers []string, statedb *dbstore.StateDB, blockchainDb *dbstore.BlockchainDB, txpool *txpool.TxPool, txpoolCh chan *types.Transaction, blockCh chan *types.Block) *P2PServer {
	// sanitize p2p port
	if port == "" {
		port = defaultP2pPort
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	downloader := NewDownloader(fmt.Sprintf("localhost%s", port), initPeers, txpoolCh, blockCh, blockchainDb)
	downloader.Start()

	p2psrv := &P2PServer{
		Port:                  port,
		Lis:                   lis,
		Peers:                 initPeers,
		P2PAddrBlockNumberMap: make(map[string]int),
		GRPCSrv:               grpcSrv,
		StateDB:               statedb,
		BlockchainDB:          blockchainDb,
		Txpool:                txpool,
		Downloader:            downloader,
	}

	return p2psrv
}

func (p2psrv *P2PServer) LatestBlock(ctx context.Context, in *protos.LatestBlockRequest) (*protos.LatestBlockResponse, error) {
	latestBlock, err := p2psrv.BlockchainDB.GetLatestBlock()
	if err != nil {
		return nil, err
	}

	out := &protos.LatestBlockResponse{
		Height:       latestBlock.Number.Uint64(),
		EncodedBlock: latestBlock.Serialize(),
	}

	return out, nil
}

func (p2psrv *P2PServer) BlocksInRange(ctx context.Context, in *protos.BlocksInRangeRequest) (*protos.BlocksInRangeResponse, error) {
	total := in.EndHeight - in.StartHeight + 1
	enBlocks := make([][]byte, total)

	blocks, err := p2psrv.BlockchainDB.GetBlocksInRange(uint(in.StartHeight), uint(in.EndHeight))
	if err != nil {
		return nil, err
	}

	for i, block := range blocks {
		enBlocks[i] = block.Serialize()
	}

	out := &protos.BlocksInRangeResponse{
		EncodedBlocks: enBlocks,
	}

	return out, nil
}

func (p2psrv *P2PServer) TxPoolPending(ctx context.Context, in *protos.TxpoolPendingRequest) (*protos.TxpoolPendingResponse, error) {
	pending := p2psrv.Txpool.Transactions
	serialisedTxs := make([][]byte, len(pending))

	for i, tx := range pending {
		serialisedTxs[i] = tx.Serialize()
	}

	out := &protos.TxpoolPendingResponse{
		EncodedTxs: serialisedTxs,
	}

	return out, nil
}

func (p2psrv *P2PServer) StartServer() {
	protos.RegisterP2PServer(p2psrv.GRPCSrv, p2psrv)
	fmt.Println("Serving P2P Server on port", p2psrv.Port)

	if err := p2psrv.GRPCSrv.Serve(p2psrv.Lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
