package file

import (
	"testing"
)

var (
	blockHeader = BlockHeader{
		DatanodeID: "1",
		Filename: "fake_block",
		Size: 200,
		Index: 0,
		NumberOfBlocks: 1,
	}
)

func TestGetBlockName(t *testing.T) {
	name := GetBlockName(blockHeader)
	if name != "block_0" {
		t.Error("wrong block name.")
	}
}
