package chunk

import (
	"testing"
	"path"
	"github.com/JetMuffin/whalefs/types"
	"os"
	"strings"
	"reflect"
	"bytes"
)

var (
	blockStore = BlockStore{DataDir: "store"}
	block = &types.BlockHeader{
		BlockID: types.BlockID("fake_block_1"),
	}

)

func TestNewBlockStore(t *testing.T) {
	b := NewBlockStore("store")
	defer os.RemoveAll(b.DataDir)

	info, err := os.Stat(b.BlocksDirectory())
	t.Log(info.IsDir())
	if err != nil {
		t.Errorf("Cannot create a new block store: %v.", err)
	}
}

func TestBlockStore_BlocksDirectory(t *testing.T) {
	expectedDir := path.Join("store", "blocks")
	if blockStore.BlocksDirectory() != expectedDir {
		t.Error("block storage directory incorrect.")
	}
}

func TestBlockStore_BlockStoragePath(t *testing.T) {
	expectedPath := path.Join("store", "blocks", string(block.BlockID))
	if blockStore.BlockStoragePath(block.BlockID) != expectedPath {
		t.Error("block storage path incorrect.")
	}
}

func TestBlockStore_MetaDirectory(t *testing.T) {
	expectedDir := path.Join("store", "meta")
	if blockStore.MetaDirectory() != expectedDir {
		t.Error("meta storage directory incorrect.")
	}
}

func TestBlockStore_BlockCheckSumPath(t *testing.T) {
	expectedPath := path.Join("store", "meta", string(block.BlockID)+".crc32")
	if blockStore.BlockMetaPath(block.BlockID) != expectedPath {
		t.Error("checksum storage path incorrect.")
	}
}

func TestBlockStore_BlockSize(t *testing.T) {
	blockId := types.BlockID("non exists id")
	size, err := blockStore.BlockSize(blockId)
	if size != -1 || err == nil {
		t.Error("size should be -1 for non exists block id.")
	}
}

func TestBlockStore_RWBlockAndMeta(t *testing.T) {
	os.MkdirAll(blockStore.BlocksDirectory(), 0700)
	os.MkdirAll(blockStore.MetaDirectory(), 0700)
	defer os.RemoveAll(blockStore.DataDir)

	// Test WriteBlock()
	r := strings.NewReader("fake content")
	checksum, err := blockStore.WriteBlock(block.BlockID, 12, r)
	if err != nil {
		t.Errorf("error when write blocks: %v.", err)
	}

	// Test BlockChecksum()
	if c, _ := blockStore.BlockCheckSum(block.BlockID); c != checksum {
		t.Error("checksum compute error.")
	}

	// Test WriteChecksum()
	if err = blockStore.WriteMeta(block, checksum); err != nil{
		t.Errorf("error when write checksum: %v", err)
	}

	// Test ReadChecksum()
	if c, _ := blockStore.ReadMeta(block.BlockID); !reflect.DeepEqual(block, c) {
		t.Error("incorrect meta read from file.")
	}

	// Test BlockSize()
	size, err := blockStore.BlockSize(block.BlockID)
	if size != 12 {
		t.Error("wrong block size.")
	}

	// Test ListBlocks()
	expectedList := []types.BlockID{block.BlockID}
	list, err := blockStore.ListBlocks()
	if !reflect.DeepEqual(expectedList, list) {
		t.Error("list blocks error.")
	}

	// Test ReadBlocks()
	w := bytes.NewBufferString("")
	err = blockStore.ReadBlock(block.BlockID, w)
	if err != nil {
		t.Errorf("error when read blocks: %v.", err)
	}
	if w.String() != "fake content" {
		t.Error("block content incorrect.")
	}
}

func TestBlockStore_Utilization(t *testing.T) {
	os.MkdirAll(blockStore.BlocksDirectory(), 0700)
	os.MkdirAll(blockStore.MetaDirectory(), 0700)
	defer os.RemoveAll(blockStore.DataDir)

	if blockStore.Utilization() != 0 {
		t.Error("block usage error.")
	}
}