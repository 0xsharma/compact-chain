package p2p

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/0xsharma/compact-chain/dbstore"
	"github.com/0xsharma/compact-chain/protos"
	"github.com/0xsharma/compact-chain/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Downloader struct {
	Peers []*Peer
	Self  string

	TxpoolCh     chan *types.Transaction
	BlockCh      chan *types.Block
	BlockchainDB *dbstore.BlockchainDB
}

type Peer struct {
	Addr        string
	ClientConn  *grpc.ClientConn
	P2PClient   protos.P2PClient
	LatestBlock *types.Block
}

func NewDownloader(self string, initPeers []string, txpoolCh chan *types.Transaction, blockCh chan *types.Block, blockchainDB *dbstore.BlockchainDB) *Downloader {
	downloader := &Downloader{
		TxpoolCh:     txpoolCh,
		BlockCh:      blockCh,
		Self:         self,
		BlockchainDB: blockchainDB,
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
		go peer.PeerBlocksLoop(d.BlockCh, *d.BlockchainDB)
		go peer.PeerTxpoolLoop(d.TxpoolCh)
	}
}

func (d *Downloader) GetPeers() []*Peer {
	return d.Peers
}

func (p *Peer) PeerBlocksLoop(blockCh chan *types.Block, blockchainDB dbstore.BlockchainDB) {
	for {
		localLatest, err := blockchainDB.GetLatestBlock()
		if err != nil {
			fmt.Println("Error Fetching Latest Block in Downloader", err)
			time.Sleep(500 * time.Millisecond)
		}

		r, err := p.P2PClient.LatestBlock(context.Background(), &protos.LatestBlockRequest{})
		if err != nil {
			time.Sleep(5000 * time.Millisecond)
			continue
		}

		rBlock := types.DeserializeBlock(r.EncodedBlock)

		// nolint : nestif
		if localLatest.Number.Int64() >= rBlock.Number.Int64() {
			if localLatest.Number.Int64() == rBlock.Number.Int64() && localLatest.DeriveHash().String() != rBlock.DeriveHash().String() {
				// send block to core.Blockchain
				blockCh <- rBlock
			} else {
				time.Sleep(500 * time.Millisecond)
				continue
			}
		} else if rBlock.Number.Int64()-localLatest.Number.Int64() > 1 {
			endHeight := uint64(rBlock.Number.Int64())
			if rBlock.Number.Int64()-localLatest.Number.Int64() > 50 {
				endHeight = uint64(localLatest.Number.Int64() + 50)
			}

			time.Sleep(1 * time.Second)
			rBlocks, err := p.P2PClient.BlocksInRange(context.Background(), &protos.BlocksInRangeRequest{
				StartHeight: uint64(localLatest.Number.Int64() + 1),
				EndHeight:   endHeight,
			})
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			for _, block := range rBlocks.EncodedBlocks {
				blockCh <- types.DeserializeBlock(block)
			}

			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			// send block to core.Blockchain
			blockCh <- rBlock
		}

		p.LatestBlock = rBlock

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
