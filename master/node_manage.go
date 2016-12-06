package master

import (
	. "github.com/JetMuffin/whalefs/types"
	comm "github.com/JetMuffin/whalefs/communication"
)

func (m *Master) RegisterChunkNode(addr string) NodeID {
	var id comm.UUID = comm.RandUUID()
	nodeID := NodeID(id.Hex())

	m.chunks[nodeID] = NewInitialNode(addr)
	return nodeID
}