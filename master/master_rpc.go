package master

import (
	"net"
	"strconv"
	comm "github.com/JetMuffin/whalefs/communication"
	log "github.com/Sirupsen/logrus"
	. "github.com/JetMuffin/whalefs/types"
	"net/rpc"
)

type MasterRPC struct {
	nodeManager 		*NodeManager
	blockManager 		*BlockManager
	replicationController 	*ReplicationController
}

func NewMasterRPC(nodeManager *NodeManager, blockManager *BlockManager, replicationController *ReplicationController) *MasterRPC{
	return &MasterRPC{
		nodeManager: nodeManager,
		blockManager: blockManager,
		replicationController: replicationController,
	}
}

func (c *MasterRPC) Register(message comm.RegistrationMessage, reply *NodeID) error {
	nodeID := c.nodeManager.RegisterChunkNode(message.Addr)
	log.WithField("nodeID", nodeID).Infof("Chunk node registered at %v, %v nodes totally.", message.Addr,
		len(c.nodeManager.chunks))
	*reply = nodeID

	return nil
}

func (c *MasterRPC) Heartbeat(message comm.HeartbeatMessage, reply *comm.HeartbeatResponse) error {
	log.Debugf("Heartbeat from node(%v).", message.NodeID)

	// Update node information by heartbeat
	c.nodeManager.UpdateNodeWithHeartbeat(message)

	// Tell this node to delete dead blocks
	var deadBlocks []BlockID
	var blocks []BlockID
	for _, blockID := range(message.Blocks) {
		if !c.blockManager.HasBlock(blockID) {
			deadBlocks = append(deadBlocks, blockID)
		} else {
			blocks = append(blocks, blockID)
		}
	}
	if len(deadBlocks) > 0 {
		log.Infof("Found %v inconsistent dead blocks from node %v, force to delete them.",
			len(deadBlocks), message.NodeID)
	}
	reply.DeadBlocks = deadBlocks

	// Tell this node to synchronize blocks
	syncBlocks := c.replicationController.Replicate(c.blockManager, c.nodeManager, blocks)
	if len(syncBlocks) > 0 {
		log.Infof("%v blocks will be synchorize from node %v", len(syncBlocks), message.NodeID)
	}
	reply.SyncBlocks = syncBlocks

	return nil
}

func (c *MasterRPC) SyncDone(block *BlockHeader, reply *comm.SyncDoneResponse) error {
	c.blockManager.AddBlock(block.FileID, block)
	log.Infof("Synchronize done on node %v for block %v.", block.Chunk, block.BlockID)

	return nil
}

// ListenRPC setup a RPC server on master node.
func (m *Master) ListenRPC() {
	rpc.Register(NewMasterRPC(m.nodeManager, m.blockManager, m.replicationController))
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(m.RPCPort))
	if err != nil {
		log.Fatalf("Error: listen to rpc port error: %v.", err)
	}
	log.Infof("RPC Server listen on :%v.", strconv.Itoa(m.RPCPort))
	go rpc.Accept(listener)
}