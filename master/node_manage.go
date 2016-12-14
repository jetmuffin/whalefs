package master

import (
	. "github.com/JetMuffin/whalefs/types"
	log "github.com/Sirupsen/logrus"
	comm "github.com/JetMuffin/whalefs/communication"
	"sync"
)

type NodeManager struct {
	chunks		map[NodeID]*Node
	lock		sync.RWMutex
}

func NewNodeManager() *NodeManager{
	return &NodeManager{
		chunks: make(map[NodeID]*Node),
	}
}

func (n *NodeManager) AddNode(node *Node) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.chunks[node.ID] = node
}

func (n *NodeManager) DeleteNode(nodeID NodeID) {
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, exists := n.chunks[nodeID]; exists {
		delete(n.chunks, nodeID)
	}
}

func (n *NodeManager) UpdateNode(node *Node) {
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, exists := n.chunks[node.ID]; exists {
		n.chunks[node.ID] = node
	}
}

func (n *NodeManager) GetNode(id NodeID) *Node{
	n.lock.RLock()
	defer n.lock.RUnlock()
	if node, exists := n.chunks[id]; exists {
		return node
	} else {
		return nil
	}
}

// RegisterChunkNode generate universal unique id for a chunk node and register this node.
// to master's map.
func (n *NodeManager) RegisterChunkNode(addr string) NodeID {
	var id comm.UUID = comm.RandUUID()
	nodeID := NodeID(id.Hex())

	n.AddNode(NewInitialNode(addr, nodeID))
	return nodeID
}

// ReRegisterChunkNode re-register a lost node to node manager.
func (n *NodeManager) ReRegisterChunkNode(addr string, id NodeID) {
	n.AddNode(NewInitialNode(addr, id))
}

// UpdateNodeWithHeartbeat update node information with heartbeat message.
func (n *NodeManager) UpdateNodeWithHeartbeat(message comm.HeartbeatMessage) {
	node := n.GetNode(message.NodeID)

	// If node doesn't exist in storage, this must be a lost node.
	// Re-register it to node manager.
	if node == nil {
		n.ReRegisterChunkNode(message.Addr, message.NodeID)
	}

	node.LastHeartbeat = message.Timestamp
	node.Heath = Healthy
	n.UpdateNode(node)
}

// LostNode wipe a node from healthy node.
func (n *NodeManager) LostNode(node *Node) {
	// TODO: label node as unhealthy and enable re-registerting

	n.DeleteNode(node.ID)
	log.Infof("Node %v disconnect from master, %v nodes totally now.", node.ID, len(n.chunks))
}