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
	Peers []*Peer
	Self  string

	TxpoolCh chan *types.Transaction
	BlockCh  chan *types.Block
}

type Peer struct {
	Addr        string
	ClientConn  *grpc.ClientConn
	P2PClient   protos.P2PClient
	LatestBlock *types.Block
}

func NewDownloader(self string, initPeers []string, txpoolCh chan *types.Transaction, blockCh chan *types.Block) *Downloader {
	downloader := &Downloader{
		TxpoolCh: txpoolCh,
		BlockCh:  blockCh,
		Self:     self,
	}

	for _, peer := range initPeers {
		if peer == downloader.Self {
			continue
		}

		conn, c := ConnectToGRPCServer(peer)
		downloader.Peers = append(downloader.Peers, &Peer{
			Addr:       peer,
			ClientConn: conn,
			P2PClient:  c,
		})
	}

	return downloader
}

func (d *Downloader) Start() {
	for _, peer := range d.Peers {
		go peer.PeerBlocksLoop(d.BlockCh)
		go peer.PeerTxpoolLoop(d.TxpoolCh)
	}
}

func (d *Downloader) GetPeers() []*Peer {
	return d.Peers
}

func (p *Peer) PeerBlocksLoop(blockCh chan *types.Block) {
	for {
		r, err := p.P2PClient.LatestBlock(context.Background(), &protos.LatestBlockRequest{})
		if err != nil {
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		rBlock := types.DeserializeBlock(r.EncodedBlock)

		p.LatestBlock = rBlock

		// send block to core.Blockchain
		blockCh <- rBlock

		time.Sleep(100 * time.Millisecond)
	}
}

func (p *Peer) PeerTxpoolLoop(txpoolCh chan *types.Transaction) {
	for {
		rTxpool, err := p.P2PClient.TxPoolPending(context.Background(), &protos.TxpoolPendingRequest{})
		if err != nil {
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		for _, tx := range rTxpool.EncodedTxs {
			// send tx to txpool.Txpool
			txpoolCh <- types.DeserializeTransaction(tx)
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
