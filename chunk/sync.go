package chunk

import (
	. "github.com/JetMuffin/whalefs/types"
	comm "github.com/JetMuffin/whalefs/communication"
	log "github.com/Sirupsen/logrus"
	"time"
	"bytes"
)

func (chunk *ChunkServer) synchronize()  {
	go func() {
		syncBlock := <- chunk.blocksToSync
		chunk.sendBlock(syncBlock)
	} ()
}

func (chunk *ChunkServer) synchronizeDone() {
	go func() {
		block := <- chunk.blockSyncDone
		chunk.alertMaster(block)
	} ()
}

func (chunk *ChunkServer) sendBlock(syncBlock *SyncBlock) {
	client, err := comm.NewRPClient(syncBlock.Addr, 5 * time.Second)

	// read block meta
	blockHeader, err := chunk.store.ReadMeta(syncBlock.BlockID)

	// read block data
	w := bytes.NewBufferString("")
	err = chunk.store.ReadBlock(syncBlock.BlockID, w)
	if err != nil {
		log.Errorf("Cannot synchronize block: %v", err)
	}

	// re-create block struct
	block := NewBlock(blockHeader.Filename, w.Bytes(), blockHeader.Size, syncBlock.NodeID)
	block.Header.FileID = blockHeader.FileID
	block.Header.Checksum = blockHeader.Checksum

	var checksum string
	client.Connection.Call("ChunkRPC.Sync", block, &checksum)
	log.WithField("checksum", checksum).Infof("Successful synchronize block %v", block.ID)
}

func (chunk *ChunkServer) alertMaster(block *BlockHeader) {
	client, err := comm.NewRPClient(chunk.MasterAddr, 5 * time.Second)
	if err != nil {
		log.Errorf("Cannot alert synchronize status: %v", err)
	}

	var message comm.SyncDoneResponse
	client.Connection.Call("MasterRPC.SyncDone", block, &message)
}