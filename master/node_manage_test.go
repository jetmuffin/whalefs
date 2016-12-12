package master

import (
	"testing"
	"github.com/JetMuffin/whalefs/cmd"
)

var (
	config, _ = cmd.NewConfig("../conf/whale.conf")
	master = NewMaster(config)
)

func TestMaster_RegisterChunkNode(t *testing.T) {
	addr := "192.168.0.1"
	nodeID := master.RegisterChunkNode(addr)

	if master.chunks[nodeID].Addr != addr {
		t.Error("register node error.")
	}
}
