package master

import (
	. "github.com/JetMuffin/whalefs/types"
)

type Master struct {
	RPCPort int
	chunks	map[NodeID]*Node
	blocks	map[BlockID]map[NodeID]bool
}

func NewMaster(port int) *Master{
	return &Master{
		RPCPort: port,
		chunks: make(map[NodeID]*Node),
		blocks: make(map[BlockID]map[NodeID]bool),
	}
}

func (master *Master) Run() {
	go master.RunRPC()
}