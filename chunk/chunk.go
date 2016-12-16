package chunk

import (
	"time"
	"net"
	"net/rpc/jsonrpc"
	"net/rpc"
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

	blocks, err := chunk.store.ListBlocks()
	if err != nil {
		log.Errorf("Cannot list chunk blocks: %v", err)
	}
	log.Infof("Chunk data store in directory '%v', current blocks number: %v.", chunk.store.DataDir, len(blocks))

	return chunk
}

// Heartbeat send chunk server's heart according to an interval.
func (chunk *ChunkServer) Heartbeat() {
	for {
		heartbeat(chunk)
		time.Sleep(chunk.heartbeatInterval)
	}
}

func heartbeat(chunk *ChunkServer) {
	conn, err := net.Dial("tcp", chunk.MasterAddr)
	if err != nil {
		log.Errorf("Couldn't connect to master at %v", chunk.MasterAddr)
		chunk.NodeID = ""
		return
	}
	codec := jsonrpc.NewClientCodec(conn)
	client := rpc.NewClientWithCodec(codec)
	defer codec.Close()

	// if chunk server has no id, register it to master and get a node id.
	if len(chunk.NodeID) == 0 {
		err = client.Call("Register", &comm.RegistrationMessage{Addr: chunk.Addr}, &chunk.NodeID)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Registered to master(%v) and got node id %v", chunk.MasterAddr, chunk.NodeID)
		return
	}

	currentBlocks, err := chunk.store.ListBlocks()
	if err != nil {
		log.Errorf("Cannot list chunk blocks: %v", err)
	}

	// send heartbeat to master
	err = client.Call("Heartbeat", comm.HeartbeatMessage{
		NodeID: 	chunk.NodeID,
		Addr:		chunk.Addr,
		Blocks: 	currentBlocks,
		Timestamp: 	time.Now(),
	}, nil)
	if err != nil {
		log.Errorf("Heartbeat error: %v", err)
	}
}

// Run methods run up all necessary goroutines.
func (chunk *ChunkServer) Run() {
	go chunk.RunRPC()
	go chunk.Heartbeat()
}