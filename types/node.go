package types

import (
	"time"
)

type NodeStatus int

var (
	Healthy = NodeStatus(0)
	Unhealthy = NodeStatus(1)
)

type Node struct {
	Hostname	string
	IP 		string

	Heath		NodeStatus
	LastHeartbeat 	time.Time
	lastUtilization int
}

// NewInitialNode return a whole new node with initial information.
func NewInitialNode(ip string, hostname string) *Node{
	return &Node {
		Hostname:	hostname,
		IP: 		ip,
		Heath: 		Healthy,
		LastHeartbeat: 	time.Now(),
		lastUtilization: 0,
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