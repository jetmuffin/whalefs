package master

import (
	"testing"
	. "github.com/JetMuffin/whalefs/types"
	"reflect"
	"strconv"
)

var (
	manager = NewNodeManager()
	node = NewInitialNode("127.0.0.1", "abcdefg")
)

func TestMaster_RegisterChunkNode(t *testing.T) {
	addr := "192.168.0.1"
	nodeID := manager.RegisterChunkNode(addr)

	if manager.chunks[nodeID].Addr != addr {
		t.Error("register node error.")
	}
}


func TestMaster_NodeAUDG(t *testing.T) {
	// add node
	manager.AddNode(node)

	// get non-exists node
	if fake_node := manager.GetNode("fake"); fake_node != nil {
		t.Error("should return nil for non-exists node id.")
	}

	// get exists node
	newNode := manager.GetNode(node.ID)
	if !reflect.DeepEqual(newNode, node) {
		t.Error("error get node method.")
	}

	// update node
	newNode.Addr = "192.168.0.2"
	manager.UpdateNode(newNode)
	if updatedNode := manager.GetNode(newNode.ID); updatedNode.Addr != newNode.Addr {
		t.Error("update node error.")
	}

	// delete node
	manager.DeleteNode(newNode.ID)
	if deletedNode := manager.GetNode(newNode.ID); deletedNode != nil {
		t.Error("delete node error.")
	}
}

func TestNodeManager_LeastBlocksNodes(t *testing.T) {
	manager = NewNodeManager()
	for i := 1; i <= 5; i++ {
		node := NewInitialNode("127.0.0.1", NodeID(strconv.Itoa(i)))
		for j := 1; j <= 6 - i; j++ {
			node.Blocks = append(node.Blocks, BlockID("abc"))
		}
		manager.AddNode(node)
	}

	nodes := manager.LeastBlocksNodes()
	for i, id := range nodes {
		if len(manager.GetNode(id).Blocks) != i + 1 {
			t.Error("sort nodes by block number error.")
		}
	}
}

func TestNodeManager_LeastConnectionNodes(t *testing.T) {
	manager = NewNodeManager()
	for i := 1; i <= 5; i++ {
		node := NewInitialNode("127.0.0.1", NodeID(strconv.Itoa(i)))
		node.Connections = 5 - i
		manager.AddNode(node)
	}

	nodes := manager.LeastConnectionNodes()
	for i, id := range nodes {
		if manager.GetNode(id).Connections != i {
			t.Error("sort nodes by block number error.")
		}
	}
}