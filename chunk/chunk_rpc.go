package chunk

import (
	"net"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"net/rpc"
	. "github.com/JetMuffin/whalefs/types"
)

type ChunkRPC struct {
	blockStore *BlockStore
}

func NewChunkRPC(blockStore *BlockStore) *ChunkRPC {
	return &ChunkRPC{
		blockStore: blockStore,
	}
}

func (c *ChunkRPC) Write(block Block, checksum *string) error {
	cs, err := c.blockStore.WriteBlock(block.BlockID, block.Header.Size, block.Reader())
	err = c.blockStore.WriteChecksum(block.BlockID, cs)
	*checksum = cs
	if err != nil {
		log.Errorf("Write block %v error: %v", block.BlockID, err)
		return err
	}
	log.Infof("Successful write block %v with checksum %v", block.BlockID, cs)
	return nil
}

// ListenRPC setup a RPC server on chunk node.
func (c *ChunkServer) ListenRPC() {
	rpc.Register(NewChunkRPC(c.store))
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(c.RPCPort))
	if err != nil {
		log.Fatalf("Error: listen to rpc port error: %v.", err)
	}
	log.Infof("RPC Server listen on :%v.", strconv.Itoa(c.RPCPort))
	go rpc.Accept(listener)
}