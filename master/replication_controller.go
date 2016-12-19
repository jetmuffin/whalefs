package master

import (
	. "github.com/JetMuffin/whalefs/types"
)

type ReplicationController struct {
	defaultReplication 	int
}

func NewReplicationController(defaultReplication int) *ReplicationController{
	return &ReplicationController{
		defaultReplication: defaultReplication,
	}
}

func (rc *ReplicationController) Replicate(blockManager *BlockManager, nodeManager *NodeManager, blocks []BlockID) []*SyncBlock{
	var syncBlocks []*SyncBlock
	for _, id := range(blocks) {
		block := blockManager.GetBlock(id)
		file := blockManager.GetFile(block.FileID)

		// if this file is synchronizing now or it has enough replications
		if file.Status == FileSync || file.Replications >= len(nodeManager.chunks) {
			continue
		}

		nodes := nodeManager.LeastBlocksNodes()
		for _, node := range(nodes[:(len(nodeManager.chunks) - file.Replications)]) {
			syncBlocks = append(syncBlocks, NewSyncBlock(id, nodeManager.GetNode(node).Addr, NodeID(node)))
		}
		blockManager.UpdateFileStatus(file.ID, FileSync)
	}
	return syncBlocks
}


