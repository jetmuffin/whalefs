package master

import (
	. "github.com/JetMuffin/whalefs/types"
	"sort"
)

type NodeSort struct {
	Func     func(NodeID) int
	list  	 []NodeID
}

func (n NodeSort) Len() int {
	return len(n.list)
}

func (n NodeSort) Swap(i, j int) {
	n.list[i], n.list[j] = n.list[j], n.list[i]
}

func (n NodeSort) Less(i, j int) bool {
	return n.Func(n.list[i]) < n.Func(n.list[j])
}

func SortNodeByFunc(f func(NodeID) int, list []NodeID) sort.Interface {
	return NodeSort{f, list}
}