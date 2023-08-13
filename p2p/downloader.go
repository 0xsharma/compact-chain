package p2p

import (
	"context"
	"log"
	"time"

	"github.com/0xsharma/compact-chain/protos"
	"github.com/0xsharma/compact-chain/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Downloader struct {
	Peers []Peer
}

type Peer struct {
	Addr        string
	ClientConn  *grpc.ClientConn
	P2PClient   protos.P2PClient
	LatestBlock *types.Block
}

func NewDownloader(initPeers []string) *Downloader {
	downloader := &Downloader{}

	for _, peer := range initPeers {
		conn, c := ConnectToGRPCServer(peer)
		downloader.Peers = append(downloader.Peers, Peer{
			Addr:       peer,
			ClientConn: conn,
			P2PClient:  c,
		})
	}

	return downloader
}

func (d *Downloader) Start() {
	for _, peer := range d.Peers {
		go peer.PeerBlocksLoop()
		go peer.PeerTxpoolLoop()
	}
}

func (d *Downloader) GetPeers() []Peer {
	return d.Peers
}

func (p *Peer) PeerBlocksLoop() {
	for {
		// TODO (0xsharma) : Add a channel to send this to core.Blockchain
		r, _ := p.P2PClient.LatestBlock(context.Background(), &protos.LatestBlockRequest{})
		rBlock := types.DeserializeBlock(r.EncodedBlock)

		p.LatestBlock = rBlock

		time.Sleep(100 * time.Millisecond)
	}
}

func (p *Peer) PeerTxpoolLoop() {
	for {
		// TODO (0xsharma) : Add a channel to send this to txpool.Txpool
		rTxpool, _ := p.P2PClient.TxPoolPending(context.Background(), &protos.TxpoolPendingRequest{})

		txs := []*types.Transaction{}

		for _, tx := range rTxpool.EncodedTxs {
			// nolint : staticcheck
			txs = append(txs, types.DeserializeTransaction(tx))
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func ConnectToGRPCServer(addr string) (*grpc.ClientConn, protos.P2PClient) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := protos.NewP2PClient(conn)

	return conn, c
}
