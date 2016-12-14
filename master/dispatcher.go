package master

import (
	. "github.com/JetMuffin/whalefs/types"
)

type Dispatcher struct {
	queue chan *Blob
	blockSize int
	manager *NodeManager
}

func (d *Dispatcher) dispatch() {
	//for {
	//	blob := <- d.queue
	//	bytesLeft := blob.Length
	//	for bytesLeft > 0 {
	//		var id comm.UUID = comm.RandUUID()
	//		blockID := BlockID(id.Hex())
	//
	//		bytesLeft = bytesLeft - d.blockSize
	//	}
	//}
}