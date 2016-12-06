package master

import (
	"testing"
)

var (
	master = NewMaster()
)

func TestMaster_RegisterChunkNode(t *testing.T) {
	addr := "192.168.0.1"
	nodeID := master.RegisterChunkNode(addr)

	if master.chunks[nodeID].Addr != addr {
		t.Error("register node error.")
	}
}
