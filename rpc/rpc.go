package rpc

import (
	"log"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/0xsharma/compact-chain/txpool"
)

type RPCServer struct {
	Server     *rpc.Server
	Addr       string
	HttpServer *http.Server
}

type RPCDomains struct {
	TxPool *txpool.TxPool
}

var registerOnce sync.Once

func NewRPCServer(addr string, domains *RPCDomains) *RPCServer {
	srv := rpc.NewServer()
	rpcServer := &RPCServer{Server: srv, Addr: addr}

	if err := rpcServer.ActivateModules(domains); err != nil {
		log.Fatalf("Couldn't activate modules. Error %s", err)
	}

	rpcServer.Start(addr)

	return rpcServer
}

func (s *RPCServer) Start(addr string) {
	registerOnce.Do(func() {
		rpc.HandleHTTP()
	})

	// nolint : gosec
	srv := &http.Server{Addr: addr}

	log.Println("Serving RPC handler")

	go func() {
		// nolint : gosec
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("Error serving: %s", err)
		}
	}()

	s.HttpServer = srv
}

func (s *RPCServer) ActivateModules(domains *RPCDomains) error {
	if domains.TxPool != nil {
		// nolint : errcheck
		rpc.Register(domains.TxPool)
	}

	return nil
}
