package chunk

import (
	"time"
	log "github.com/Sirupsen/logrus"
	comm "github.com/JetMuffin/whalefs/communication"
	. "github.com/JetMuffin/whalefs/cmd"
	. "github.com/JetMuffin/whalefs/types"
)

// ChunkServer is the slave node which store data blocks.
type ChunkServer struct {
	NodeID  	  NodeID
	Addr 		  string
	MasterAddr 	  string
	RPCPort	 	  int
	store 		  *BlockStore
	rpcClient 	  *comm.RPCClient
	heartbeatInterval time.Duration

	blocksToDelete	  chan BlockID
	deadBlocks	  []BlockID
}


// NewChunkServer returns a server which store data.
func NewChunkServer(config *Config) *ChunkServer {
	chunk := &ChunkServer{
		RPCPort: config.Int("chunk_port"),
		MasterAddr: config.String("master_addr"),
		Addr: config.String("chunk_ip"),
		store: NewBlockStore(config.String("chunk_data_dir")),
		heartbeatInterval: 1 * time.Second,
	}

	client, err := comm.NewRPClient(chunk.MasterAddr, 10 * time.Second)
	if err != nil {
		log.Fatalf("Cannot connect to master: %v", err)
	}
	chunk.rpcClient = client

	blocks, err := chunk.store.ListBlocks()
	if err != nil {
		log.Errorf("Cannot list chunk blocks: %v", err)
	}
	log.Infof("Chunk data store in directory '%v', current blocks number: %v.", chunk.store.DataDir, len(blocks))

	return chunk
}

// Heartbeat send chunk server's heart according to an interval.
func (chunk *ChunkServer) Heartbeat() {
	go func() {
		for {
			// TODO: Check if the connection is closed by master or not.

			heartbeat(chunk)
			time.Sleep(chunk.heartbeatInterval)
		}
	} ()
}

func heartbeat(c *ChunkServer) {
	// if chunk server has no id, register it to master and get a node id.
	if len(c.NodeID) == 0 {
		err := c.rpcClient.Connection.Call("MasterRPC.Register", &comm.RegistrationMessage{Addr: c.Addr},
			&c.NodeID)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Registered to master(%v) and got node id %v", c.MasterAddr, c.NodeID)
		return
	}

	currentBlocks, err := c.store.ListBlocks()
	if err != nil {
		log.Errorf("Cannot list chunk blocks: %v", err)
	}

	// send heartbeat to master
	err = c.rpcClient.Connection.Call("MasterRPC.Heartbeat", comm.HeartbeatMessage{
		NodeID: 	c.NodeID,
		Addr:		c.Addr,
		Blocks: 	currentBlocks,
		Timestamp: 	time.Now(),
	}, nil)
	if err != nil {
		log.Errorf("Cannot send heart beat to master: %v", err)
	}
}

// Run methods run up all necessary goroutines.
func (chunk *ChunkServer) Run() {
	chunk.ListenRPC()
	chunk.Heartbeat()
}