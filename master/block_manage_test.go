package master

import (
	"testing"
	. "github.com/JetMuffin/whalefs/types"
	"strconv"
	"reflect"
)

var (
	bm = NewBlockManager(0, 0)
)

func TestBlockManager_File(t *testing.T) {
	var files []*File
	for i := 0; i < 5; i++ {
		files = append(files, NewFile("file_" + strconv.Itoa(i), 0))
		bm.AddFile(files[i])
	}

	expectedFiles := bm.ListFile()
	if len(expectedFiles) != 5 {
		t.Error("List file error.")
	}

	if file := bm.GetFile(files[0].ID); !reflect.DeepEqual(file, files[0]) {
		t.Error("Get file error.")
	}

	bm.DeleteFile(files[0].ID)
	if file := bm.GetFile(files[0].ID); len(bm.files) != 4 || file != nil {
		t.Error("Delete file error.")
	}

	bm.UpdateFileStatus(files[1].ID, FileWriting)
	if file := bm.GetFile(files[1].ID); file.Status != FileWriting {
		t.Error("Update file status error.")
	}

	bm.AddBlock(files[1].ID, NewBlock("1", nil, 0).Header)
	if file := bm.GetFile(files[1].ID); len(file.Blocks) != 1 {
		t.Error("Add block error.")
	}
}