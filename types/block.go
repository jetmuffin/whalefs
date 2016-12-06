package types

import (
	"strconv"
)

type BlockID string

// Block is storage unit of each file
type Block struct {
	Header 	BlockHeader // metadata
	Data 	[]byte	    // contents of this block
}

// BlockHeaders holds Block metadata
type BlockHeader struct {
	DatanodeID 	string	// ID of datanode which store this block
	Filename   	string  // Storage name of this block
	Size       	int64   // Size of block in bytes
	Index   	int     // Index in file
	NumberOfBlocks 	int     // Total number of blocks in file
}

func GetBlockName(header BlockHeader) string {
	return "block_" + strconv.Itoa(header.Index)
}





