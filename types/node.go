package types

import (
	"time"
)

type NodeID string

type NodeStatus int

var (
	Healthy = NodeStatus(0)
	Unhealthy = NodeStatus(1)
)

type Node struct {
	ID		NodeID
	Addr 		string

	Blocks		[]BlockID
	Connections	int
	Heath		NodeStatus
	LastHeartbeat 	time.Time
	LastUtilization int64
}

// NewInitialNode return a whole new node with initial information.
func NewInitialNode(addr string, id NodeID) *Node{
	return &Node {
		ID:		id,
		Addr: 		addr,
		Heath: 		Healthy,
		Connections:	0,
		LastHeartbeat: 	time.Now(),
		LastUtilization: 0,
	}
}

// IsHealthy check if a node is healthy or not.
func (node *Node) IsHealthy() bool {
	return node.Heath == Healthy
}

// HeartbeatDuration compute the duration between last heartbeat and now.
func (node *Node) HeartbeatDuration() time.Duration {
	return time.Now().Sub(node.LastHeartbeat)
}