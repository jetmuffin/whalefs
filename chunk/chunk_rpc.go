package chunk

import (
	"net"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"net/rpc"
	. "github.com/JetMuffin/whalefs/types"
	"github.com/JetMuffin/whalefs/communication"
	"bytes"
	"time"
)

type ChunkRPC struct {
	blockStore *BlockStore
	blockSyncDone chan *BlockHeader
}

func NewChunkRPC(blockStore *BlockStore, blockSyncDone chan *BlockHeader) *ChunkRPC {
	return &ChunkRPC{
		blockStore: blockStore,
		blockSyncDone: blockSyncDone,
	}
}

func (c *ChunkRPC) Write(block Block, checksum *string) error {
	cs, err := c.blockStore.WriteBlock(block.ID, block.Header.Size, block.Reader())
	err = c.blockStore.WriteMeta(block.Header, cs)
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

func (c *ChunkRPC) Sync(block Block, checksum *string) error {
	cs, err := c.blockStore.WriteBlock(block.ID, block.Header.Size, block.Reader())
	err = c.blockStore.WriteMeta(block.Header, cs)
	if err != nil {
		log.Errorf("Write block %v error: %v", block.ID, err)
		return err
	}
	log.Infof("Successful synchronize block %v with checksum %v", block.ID, cs)

	c.blockSyncDone <- block.Header
	return nil
}

func (c *ChunkRPC) Time(args interface{}, timestamp *int64) error {
	*timestamp = time.Now().Unix()
	return nil
}

// ListenRPC setup a RPC server on chunk node.
func (c *ChunkServer) ListenRPC() {
	rpc.Register(NewChunkRPC(c.store, c.blockSyncDone))
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(c.RPCPort))
	if err != nil {
		log.Fatalf("Error: listen to rpc port error: %v.", err)
	}
	log.Infof("RPC Server listen on :%v.", strconv.Itoa(c.RPCPort))
	go rpc.Accept(listener)
}