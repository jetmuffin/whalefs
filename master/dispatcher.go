package master

import (
	. "github.com/JetMuffin/whalefs/types"
	comm "github.com/JetMuffin/whalefs/communication"
	log "github.com/Sirupsen/logrus"
	"time"
)

type Dispatcher struct {
	blockManager 	*BlockManager
	nodeManager 	*NodeManager
}

func NewDispatcher(blockManager *BlockManager, nodeManager *NodeManager) *Dispatcher{
	return &Dispatcher{
		blockManager: blockManager,
		nodeManager: nodeManager,
	}
}

func (d *Dispatcher) Dispatch() {
	log.Info("Dispatcher start dispatch.")
	go func() {
		for {
			blob := <- d.blockManager.blobQueue
			log.Info("Got a blob, start dispatch.")
			chunks := d.nodeManager.LeastBlocksNodes()


			if len(chunks) < d.blockManager.blockReplication {
				log.Error("Cannot write block: chunk number less than replication.")
				continue
			}
			d.blockManager.UpdateFileStatus(blob.FileID, FileWriting)

			for i := 0; i < d.blockManager.blockReplication; i++ {
				node := d.nodeManager.GetNode(chunks[i])
				block := NewBlock(blob.Name, blob.Content, blob.Length)
				block.Header.Replications = d.blockManager.blockReplication
				block.Header.Chunk = node.ID

				client, err := comm.NewRPClient(node.Addr, 5 * time.Second)
				if err != nil {
					log.Errorf("Cannot connect to node %v: %v", node.Addr, err)
				}

				var checksum string
				client.Connection.Call("ChunkRPC.Write", block, &checksum)

				d.blockManager.AddBlock(blob.FileID, block.Header)
				log.WithField("checksum", checksum).Infof("Write block %v successful", block.ID)
			}

			d.blockManager.UpdateFileStatus(blob.FileID, FileOK)
		}
	} ()
}

