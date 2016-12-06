package chunk

import (
	"testing"
	"os"
	"path"
	"io/ioutil"
	"reflect"
	"github.com/JetMuffin/whalefs/file"
)

var (
	data = []byte{1}
	chunkServer = NewChunkServer("test")
	blockHeader = file.BlockHeader{
		DatanodeID: "1",
		Filename: "fake_block",
		Size: 200,
		Index: 0,
		NumberOfBlocks: 1,
	}
	blockNotFoundHeader = file.BlockHeader {
		Filename: "not_found_block",
	}
	block = file.Block{
		Header: blockHeader,
		Data: data,
	}
)

func TestWriteBlock(t *testing.T) {
	defer func() {
		os.RemoveAll(chunkServer.RootDir)
	} ()

	err := chunkServer.WriteBlock(block)
	if err == nil {
		t.Error("should throw chunkserver not exists error.")
	}

	err = os.Mkdir(chunkServer.RootDir, 0700)
	if err != nil {
		t.Error(err.Error())
	}

	err = chunkServer.WriteBlock(block)
	if err != nil {
		t.Errorf("write block with error: %v.", err)
	}

	_, err = os.Stat(path.Join(chunkServer.RootDir, blockHeader.Filename))
	if err != nil && !os.IsExist(err) {
		t.Error("file directory has not been created.")
	}

	storePath := path.Join(chunkServer.RootDir, blockHeader.Filename, file.GetBlockName(blockHeader))
	_, err = os.Stat(storePath)
	if err != nil && !os.IsExist(err) {
		t.Error("block has not been written.")
	}
}

func TestReadBlock(t *testing.T)  {
	defer func() {
		os.RemoveAll(chunkServer.RootDir)
	} ()

	os.MkdirAll(path.Join(chunkServer.RootDir, blockHeader.Filename), 0700)
	ioutil.WriteFile(path.Join(chunkServer.RootDir, blockHeader.Filename, file.GetBlockName(blockHeader)), data, os.ModePerm)

	blockGot, err := chunkServer.ReadBlock(blockHeader)
	if err != nil {
		t.Errorf("read block error: %v.", err)
	} else if (!reflect.DeepEqual(blockGot.Data, data)){
		t.Error("block data should be same with which were written.")
	}

	_, err = chunkServer.ReadBlock(blockNotFoundHeader)
	if err == nil {
		t.Error("should throw block not found error.")
	}
}