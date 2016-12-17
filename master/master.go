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
	dispatcher 		*Dispatcher
	heartbeatCheckInterval 	time.Duration
}

// NewMaster returns a master.
func NewMaster(config *Config) *Master{
	master := &Master{
		RPCPort: config.Int("master_port"),
		heartbeatCheckInterval: 10 * time.Second,
		nodeManager: NewNodeManager(),
		blockManager: NewBlockManager(config.Int("block_size"), config.Int("block_replication")),
	}
	master.dispatcher = NewDispatcher(master.blockManager, master.nodeManager)
	master.httpServer = NewHTTPServer(config.String("master_ip"), config.Int("master_http_port"),
		master.blockManager.blobQueue)
	return master
}

// Run method setup all necessary goroutines.
func (master *Master) Run() {
	// Listen to RPC port
	master.ListenRPC()

	// Start monitor
	master.Monitor()

	// Dispatcher start dispatch
	master.dispatcher.Dispatch()

	// Run http server
	master.httpServer.ListenAndServe()
}