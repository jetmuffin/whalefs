package storage

import (
	"time"
	. "github.com/JetMuffin/whalefs/types"
)

type NodeStorage struct {
	dir string
	store *EtcdStorage
}

func NewNodeStore(dir string, endpoints string) *NodeStorage {
	return &NodeStorage{
		dir: dir,
		store: NewEtcdStorage(endpoints),
	}
}

func (n *NodeStorage) Add(node *Node) error {
	return n.store.UpdateExpire(node.Hostname, node, time.Second * 10)
}

func (n *NodeStorage) Update(node *Node) error {
	return n.store.UpdateExpire(node.Hostname, node, time.Second * 10)
}

func (n *NodeStorage) Get(hostname string) (*Node, error) {
	value, err := n.store.Get(hostname)
	if err != nil {
		return nil, err
	}
	return value.(*Node), nil
}

func (n *NodeStorage) Watch() {

}

