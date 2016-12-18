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
	replicationController 	*ReplicationController
	heartbeatCheckInterval 	time.Duration
}

// NewMaster returns a master.
func NewMaster(config *Config) *Master{
	master := &Master{
		RPCPort: config.Int("master_port"),
		heartbeatCheckInterval: 10 * time.Second,
		nodeManager: NewNodeManager(),
		blockManager: NewBlockManager(),

	}
	master.replicationController = NewReplicationController(config.Int("block_replication"))
	master.httpServer = NewHTTPServer(config.String("master_ip"), config.Int("master_http_port"),
		master.blockManager, master.nodeManager)
	return master
}

// Run method setup all necessary goroutines.
func (master *Master) Run() {
	// Listen to RPC port
	master.ListenRPC()

	// Start monitor
	master.Monitor()


	// Run http server
	master.httpServer.ListenAndServe()
}