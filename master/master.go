package master

import (
	. "github.com/JetMuffin/whalefs/cmd"
	"time"
)

// Master node controllers all metadata of the whole cluster.
type Master struct {
	RPCPort 		int
	httpServer		*HTTPServer
	nodeManager		*NodeManager
	blockManager		*BlockManager
	heartbeatCheckInterval 	time.Duration
}

// NewMaster returns a master.
func NewMaster(config *Config) *Master{
	master := &Master{
		RPCPort: config.Int("master_port"),
		heartbeatCheckInterval: 10 * time.Second,
		nodeManager: NewNodeManager(),
		blockManager: NewBlockManager(config.Int("block_size")),
	}
	master.httpServer = NewHTTPServer(config.String("master_ip"), config.Int("master_http_port"),
		master.blockManager.blobQueue)
	return master
}

// Run method setup all necessary goroutines.
func (master *Master) Run() {
	go master.RunRPC()
	go master.Monitor()

	// Run http server
	master.httpServer.ListenAndServe()
}