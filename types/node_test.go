package types

import "testing"

var (
	node = NewInitialNode("127.0.0.1", "abcdefg")
)

func TestNode_IsHealthy(t *testing.T) {
	if !node.IsHealthy() {
		t.Error("node is healthy function error.")
	}
}

func TestNode_HeartbeatDuration(t *testing.T) {
	if node.HeartbeatDuration() < 0 {
		t.Error("heart beat duration compute error.")
	}
}
