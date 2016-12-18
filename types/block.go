package types

import (
	"bytes"
	"io"
)

type BlockID string

// Block is storage unit of each file
type Block struct {
	ID  BlockID
	Header 	*BlockHeader // metadata
	Data 	[]byte	    // contents of this block
}

// BlockHeaders holds Block metadata
type BlockHeader struct {
	BlockID 	BlockID		`json:"block_id"`
	FileID 		FileID		`json:"file_id"`
	Chunk		NodeID		`json:"chunk"`
	Checksum 	string  	`json:"checksum"`
	Filename   	string   	`json:"filenam"`
	Size       	int64    	`json:"size"`
}

type SyncBlock struct {
	BlockID BlockID
	Addr  	string
	NodeID 	NodeID
}

func NewSyncBlock(blockID BlockID, addr string, nodeID NodeID) *SyncBlock{
	return &SyncBlock{
		BlockID: blockID,
		Addr: addr,
		NodeID: nodeID,
	}
}

func NewBlock(filename string, data []byte, size int64, chunk NodeID) *Block{
	id := RandUUID()
	header := &BlockHeader{
		BlockID:  BlockID(id.Hex()),
		Filename: filename,
		Size: size,
		Chunk: chunk,
	}
	block := &Block{
		ID: BlockID(id.Hex()),
		Header: header,
		Data: data,
	}
	return block
}


func (b *Block) Reader() io.Reader {
	return bytes.NewReader(b.Data)
}





