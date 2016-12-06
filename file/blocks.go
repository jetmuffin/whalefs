package chunk

import (
	"io/ioutil"
	"os"
	"path"
	"errors"
	"strconv"
)

var (
	ErrBlockNotFound = errors.New("Block not found.")
)

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

func getBlockName(header BlockHeader) string {
	return "block_" + strconv.Itoa(header.Index)
}

func (c *ChunkServer) WriteBlock(block Block) error {
	list, err := ioutil.ReadDir(c.RootDir)
	storePath := path.Join(c.RootDir, block.Header.Filename, getBlockName(block.Header))
	if err != nil {
		return err
	}

	for _, dir := range list {
		if dir.Name() == block.Header.Filename {
			ioutil.WriteFile(storePath, block.Data, os.ModeAppend)
			return nil
		}
	}

	os.Mkdir(path.Join(c.RootDir, block.Header.Filename), 0700)
	err = ioutil.WriteFile(storePath, block.Data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChunkServer) ReadBlock(header BlockHeader) (*Block, error){
	storePath := path.Join(c.RootDir, header.Filename, getBlockName(header))

	list, err := ioutil.ReadDir(c.RootDir)
	if err != nil {
		return nil, err
	}

	for _, dir := range list {
		if dir.Name() == header.Filename {
			bytes, err := ioutil.ReadFile(storePath)
			if err != nil {
				return nil, err
			}
			return &Block {Header: header, Data: bytes}, nil
		}
	}
	return nil, ErrBlockNotFound
}




