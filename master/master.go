package master

import (
	. "github.com/JetMuffin/whalefs/types"
)

type Master struct {
	chunks	map[NodeID]*Node
	blocks map[BlockID]map[NodeID]bool
}

func NewMaster() *Master{
	return &Master{
		chunks: make(map[NodeID]*Node),
		blocks: make(map[BlockID]map[NodeID]bool),
	}
}