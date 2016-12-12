package master

import (
	. "github.com/JetMuffin/whalefs/types"
	. "github.com/JetMuffin/whalefs/cmd"
	"time"
	"sync"
)

// Master node controllers all metadata of the whole cluster.
type Master struct {
	RPCPort 		int
	chunks			map[NodeID]*Node
	blocks			map[BlockID]map[NodeID]bool

	heartbeatCheckInterval 	time.Duration
	nodeLock		sync.RWMutex
}

// NewMaster returns a master.
func NewMaster(config *Config) *Master{
	return &Master{
		RPCPort: config.Int("master_port"),
		chunks: make(map[NodeID]*Node),
		blocks: make(map[BlockID]map[NodeID]bool),
		heartbeatCheckInterval: 10 * time.Second,
	}
}

// Run method setup all necessary goroutines.
func (master *Master) Run() {
	go master.RunRPC()
	go master.Monitor()
}