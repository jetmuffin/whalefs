package types

import "time"

type NodeID string

type Node struct {
	Addr 		string
	lastHeartbeat 	time.Time
	lastUtilization int
}

func NewInitialNode(addr string) *Node{
	return &Node {
		Addr: addr,
		lastHeartbeat: time.Now(),
		lastUtilization: 0,
	}
}