package chunk

import (
	"github.com/JetMuffin/whalefs/types"
	"time"
)

type ChunkServer struct {
	NodeID  	  string
	Store 		  BlockStore
	Addr 		  string
	MasterAddr 	  string
	heartbeatInterval time.Duration

	blocksToDelete	  chan types.BlockID
	deadBlocks	  []types.BlockID
}


// NewChunkServer returns a server which store data.
func NewChunkServer() *ChunkServer {
	// TODO
	return nil
}