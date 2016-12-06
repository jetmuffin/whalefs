package chunk

import (
	"io/ioutil"
	"os"
	"path"
	"errors"

	"github.com/JetMuffin/whalefs/file"
)

var (
	ErrBlockNotFound = errors.New("Block not found.")
)

func (c *ChunkServer) WriteBlock(block file.Block) error {
	list, err := ioutil.ReadDir(c.RootDir)
	storePath := path.Join(c.RootDir, block.Header.Filename, file.GetBlockName(block.Header))
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

func (c *ChunkServer) ReadBlock(header file.BlockHeader) (*file.Block, error){
	storePath := path.Join(c.RootDir, header.Filename, file.GetBlockName(header))

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
			return &file.Block {Header: header, Data: bytes}, nil
		}
	}
	return nil, ErrBlockNotFound
}

