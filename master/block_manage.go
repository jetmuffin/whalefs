package master

import (
	. "github.com/JetMuffin/whalefs/types"
)

type BlockManager struct {
	blocks			map[BlockID]map[NodeID]bool
	blockSize		int
	blockReplication	int
	blobQueue 		chan *Blob
}

func NewBlockManager(blockSize int, blockReplication int) *BlockManager {
	return &BlockManager{
		blocks: make(map[BlockID]map[NodeID]bool),
		blockSize: blockSize,
		blockReplication: blockReplication,
		blobQueue: make(chan *Blob),
	}
}