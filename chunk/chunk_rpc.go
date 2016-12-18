package chunk

import (
	"net"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"net/rpc"
	. "github.com/JetMuffin/whalefs/types"
	"github.com/JetMuffin/whalefs/communication"
	"bytes"
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
	cs, err := c.blockStore.WriteBlock(block.ID, block.Header.Size, block.Reader())
	err = c.blockStore.WriteChecksum(block.ID, cs)
	*checksum = cs
	if err != nil {
		log.Errorf("Write block %v error: %v", block.ID, err)
		return err
	}
	log.Infof("Successful write block %v with checksum %v", block.ID, cs)
	return nil
}

func (c *ChunkRPC) Read(blockID BlockID, reply *communication.BlockMessage) error {
	w := bytes.NewBufferString("")
	err := c.blockStore.ReadBlock(blockID, w)
	checksum, err := c.blockStore.BlockCheckSum(blockID)
	reply.Data = w.Bytes()
	reply.Checksum = checksum
	log.WithField("checksum", checksum).Infof("Read block %v", blockID)
	return err
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