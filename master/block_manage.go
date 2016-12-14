package master

import (
	. "github.com/JetMuffin/whalefs/types"
)

type BlockManager struct {
	blocks			map[BlockID]map[NodeID]bool
	blockSize		int
	blobQueue 		chan *Blob
}

func NewBlockManager(blockSize int) *BlockManager {
	return &BlockManager{
		blocks: make(map[BlockID]map[NodeID]bool),
		blockSize: blockSize,
		blobQueue: make(chan *Blob),
	}
}

