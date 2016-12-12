package master

import (
	"testing"
	"github.com/JetMuffin/whalefs/cmd"
	. "github.com/JetMuffin/whalefs/types"
	"reflect"
)

var (
	config, _ = cmd.NewConfig("../conf/whale.conf")
	master = NewMaster(config)
	node = NewInitialNode("127.0.0.1", "abcdefg")
)

func TestMaster_RegisterChunkNode(t *testing.T) {
	addr := "192.168.0.1"
	nodeID := master.RegisterChunkNode(addr)

	if master.chunks[nodeID].Addr != addr {
		t.Error("register node error.")
	}
}


func TestMaster_NodeAUDG(t *testing.T) {
	// add node
	master.AddNode(node)

	// get non-exists node
	if fake_node := master.GetNode("fake"); fake_node != nil {
		t.Error("should return nil for non-exists node id.")
	}

	// get exists node
	newNode := master.GetNode(node.ID)
	if !reflect.DeepEqual(newNode, node) {
		t.Error("error get node method.")
	}

	// update node
	newNode.Addr = "192.168.0.2"
	master.UpdateNode(newNode)
	if updatedNode := master.GetNode(newNode.ID); updatedNode.Addr != newNode.Addr {
		t.Error("update node error.")
	}

	// delete node
	master.DeleteNode(newNode.ID)
	if deletedNode := master.GetNode(newNode.ID); deletedNode != nil {
		t.Error("delete node error.")
	}
}