package rpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/0xsharma/compact-chain/txpool"
)

type RPCServer struct {
	Server *rpc.Server
	Addr   string
}

type RPCDomains struct {
	TxPool *txpool.TxPool
}

func NewRPCServer(addr string, domains *RPCDomains) *RPCServer {
	srv := rpc.NewServer()
	rpcServer := &RPCServer{Server: srv, Addr: addr}

	if err := rpcServer.ActivateModules(domains); err != nil {
		log.Fatalf("Couldn't activate modules. Error %s", err)
	}

	go rpcServer.Start(addr)

	return rpcServer
}

func (s *RPCServer) Start(addr string) {
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalf("Couldn't start listening on port %s. Error %s", addr, e)
	}

	log.Println("Serving RPC handler")

	err := http.Serve(l, nil)
	if err != nil {
		log.Fatalf("Error serving: %s", err)
	}

}

func (s *RPCServer) ActivateModules(domains *RPCDomains) error {
	if domains.TxPool != nil {
		if err := rpc.Register(domains.TxPool); err != nil {
			return err
		}
	}

	return nil
}
