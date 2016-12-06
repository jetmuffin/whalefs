package master

import (
	. "github.com/JetMuffin/whalefs/types"
)

// Master node controllers all metadata of the whole cluster.
type Master struct {
	RPCPort int
	chunks	map[NodeID]*Node
	blocks	map[BlockID]map[NodeID]bool
}

// NewMaster returns a master.
func NewMaster(port int) *Master{
	return &Master{
		RPCPort: port,
		chunks: make(map[NodeID]*Node),
		blocks: make(map[BlockID]map[NodeID]bool),
	}
}

// Run method setup all necessary goroutines.
func (master *Master) Run() {
	go master.RunRPC()
}