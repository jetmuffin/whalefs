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
	blobQueue 		chan *Blob
	httpServer		*HTTPServer

	heartbeatCheckInterval 	time.Duration
	nodeLock		sync.RWMutex
}

// NewMaster returns a master.
func NewMaster(config *Config) *Master{
	blobQueue := make(chan *Blob)
	return &Master{
		RPCPort: config.Int("master_port"),
		httpServer: NewHTTPServer(config.String("master_ip"), config.Int("master_http_port"), blobQueue),
		chunks: make(map[NodeID]*Node),
		blocks: make(map[BlockID]map[NodeID]bool),
		blobQueue: blobQueue,
		heartbeatCheckInterval: 10 * time.Second,
	}
}

// Run method setup all necessary goroutines.
func (master *Master) Run() {
	go master.RunRPC()
	go master.Monitor()

	// Run http server
	master.httpServer.ListenAndServe()
}