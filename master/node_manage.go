package master

import (
	. "github.com/JetMuffin/whalefs/types"
	comm "github.com/JetMuffin/whalefs/communication"
)

// RegisterChunkNode generate universal unique id for a chunk node and register this node
// to master's map.
func (m *Master) RegisterChunkNode(addr string) NodeID {
	var id comm.UUID = comm.RandUUID()
	nodeID := NodeID(id.Hex())

	m.chunks[nodeID] = NewInitialNode(addr)
	return nodeID
}