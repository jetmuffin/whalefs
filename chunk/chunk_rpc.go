package chunk

import (
	"net"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"net/rpc"
)

type ChunkRPC struct {
}

func NewChunkRPC() *ChunkRPC {
	return &ChunkRPC{}
}

// ListenRPC setup a RPC server on chunk node.
func (c *ChunkServer) ListenRPC() {
	rpc.Register(NewChunkRPC())
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(c.RPCPort))
	if err != nil {
		log.Fatalf("Error: listen to rpc port error: %v.", err)
	}
	log.Infof("RPC Server listen on :%v.", strconv.Itoa(c.RPCPort))
	go rpc.Accept(listener)
}