package master

import (
	. "github.com/JetMuffin/whalefs/types"
	log "github.com/Sirupsen/logrus"
	comm "github.com/JetMuffin/whalefs/communication"
)

func (m *Master) AddNode(node *Node) {
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()
	m.chunks[node.ID] = node
}

func (m *Master) DeleteNode(nodeID NodeID) {
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()
	if _, exists := m.chunks[nodeID]; exists {
		delete(m.chunks, nodeID)
	}
}

func (m *Master) UpdateNode(node *Node) {
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()
	if _, exists := m.chunks[node.ID]; exists {
		m.chunks[node.ID] = node
	}
}

func (m *Master) GetNode(id NodeID) *Node{
	m.nodeLock.RLock()
	defer m.nodeLock.RUnlock()
	if node, exists := m.chunks[id]; exists {
		return node
	} else {
		return nil
	}
}

// RegisterChunkNode generate universal unique id for a chunk node and register this node.
// to master's map.
func (m *Master) RegisterChunkNode(addr string) NodeID {
	var id comm.UUID = comm.RandUUID()
	nodeID := NodeID(id.Hex())

	m.AddNode(NewInitialNode(addr, nodeID))
	return nodeID
}

// ReRegisterChunkNode re-register a lost node to node manager.
func (m *Master) ReRegisterChunkNode(addr string, id NodeID) {
	m.AddNode(NewInitialNode(addr, id))
}

// UpdateNodeWithHeartbeat update node information with heartbeat message.
func (m *Master) UpdateNodeWithHeartbeat(message comm.HeartbeatMessage) {
	node := m.GetNode(message.NodeID)

	// If node doesn't exist in storage, this must be a lost node.
	// Re-register it to node manager.
	if node == nil {
		m.ReRegisterChunkNode(message.Addr, message.NodeID)
	}

	node.LastHeartbeat = message.Timestamp
	node.Heath = Healthy
	m.UpdateNode(node)
}

// LostNode wipe a node from healthy node.
func (m *Master) LostNode(node *Node) {
	// TODO: label node as unhealthy and enable re-registerting

	m.DeleteNode(node.ID)
	log.Infof("Node %v disconnect from master, %v nodes totally now.", node.ID, len(m.chunks))
}