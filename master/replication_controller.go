package master

import (
	. "github.com/JetMuffin/whalefs/types"
	log "github.com/Sirupsen/logrus"
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
		if file.Status == FileSync || file.Replications >= rc.defaultReplication {
			continue
		}

		nodes := nodeManager.LeastBlocksNodes()
		// if there is no enough nodes to replicate
		if len(nodes) < (rc.defaultReplication - file.Replications) {
			log.Info("Cannot add replications because no enough chunk servers for file %v, " +
				"wait for new node to register.", file.Name)
			continue
		}

		for _, node := range(nodes[:(rc.defaultReplication - file.Replications)]) {
			syncBlocks = append(syncBlocks, NewSyncBlock(id, nodeManager.GetNode(node).Addr, NodeID(node)))
		}
		blockManager.UpdateFileStatus(file.ID, FileSync)
	}
	return syncBlocks
}


