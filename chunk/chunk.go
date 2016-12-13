package chunk

import (
	"github.com/JetMuffin/whalefs/types"
	"time"
	log "github.com/Sirupsen/logrus"
	. "github.com/JetMuffin/whalefs/cmd"
	. "github.com/JetMuffin/whalefs/types"
	"github.com/JetMuffin/whalefs/storage"
	"os"
)

// ChunkServer is the slave node which store data blocks.
type ChunkServer struct {
	Hostname 	  string
	IP 		  string
	MasterAddr 	  string
	RPCPort	 	  int
	store		  *storage.NodeStorage

	heartbeatInterval time.Duration

	blocksToDelete	  chan types.BlockID
	deadBlocks	  []types.BlockID
}


// NewChunkServer returns a server which store data.
func NewChunkServer(config *Config) *ChunkServer {
	hostname, _ := os.Hostname()
	return &ChunkServer{
		RPCPort: config.Int("chunk_port"),
		store: storage.NewNodeStore("chunks", config.String("etcd_addr")),
		MasterAddr: config.String("master_addr"),
		IP: config.String("chunk_ip"),
		Hostname: hostname,
		heartbeatInterval: 1 * time.Second,
	}
}

// Heartbeat send chunk server's heart according to an interval.
func (c *ChunkServer) Heartbeat() {
	for {
		heartbeat(c)
		time.Sleep(c.heartbeatInterval)
	}
}

func heartbeat(chunk *ChunkServer) {
	info := &Node {
		Hostname: chunk.Hostname,
		IP: chunk.IP,
	}

	if err := chunk.store.Update(info); err != nil {
		log.Errorf("Error: cannot update node status: %v", err)
	}
}

// Run methods run up all necessary goroutines.
func (chunk *ChunkServer) Run() {
	go chunk.RunRPC()
	go chunk.Heartbeat()
}