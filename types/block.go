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
	BlockID 	BlockID
	FileID 		FileID
	Chunk		NodeID
	Filename   	string  // Storage name of this block
	Size       	int64   // Size of block in bytes
	Replications 	int     // Total number of blocks in file
}

func NewBlock(filename string, data []byte, size int64) *Block{
	id := RandUUID()
	header := &BlockHeader{
		BlockID:  BlockID(id.Hex()),
		Filename: filename,
		Size: size,
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





