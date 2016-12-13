package master

import (
	"github.com/JetMuffin/whalefs/cmd"
	. "github.com/JetMuffin/whalefs/types"
	"github.com/JetMuffin/whalefs/storage"
)

type ChunkManager struct {
	chunks map[string] *Node
	store  *storage.NodeStorage
}

func NewChunkServer(config cmd.Config) *ChunkManager {
	return &ChunkManager{
		chunks: make(map[string] *Node),
		store: storage.NewNodeStore("nodes", config.String("etcd_addr")),
	}
}