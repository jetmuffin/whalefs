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
	nodeManager *NodeManager
}

func NewMasterRPC(nodeManager *NodeManager) *MasterRPC{
	return &MasterRPC{
		nodeManager: nodeManager,
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

	// TODO send response to heartbeat
	reply = new(comm.HeartbeatResponse)
	return nil
}

// ListenRPC setup a RPC server on master node.
func (m *Master) ListenRPC() {
	rpc.Register(NewMasterRPC(m.nodeManager))
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(m.RPCPort))
	if err != nil {
		log.Fatalf("Error: listen to rpc port error: %v.", err)
	}
	log.Infof("RPC Server listen on :%v.", strconv.Itoa(m.RPCPort))
	go rpc.Accept(listener)
}