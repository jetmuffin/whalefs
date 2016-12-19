package master

import (
	. "github.com/JetMuffin/whalefs/types"
	"sync"
)

type BlockManager struct {
	files			map[FileID]*File
	blocks 			map[BlockID]*BlockHeader
	lock 			sync.RWMutex
}

func NewBlockManager() *BlockManager {
	return &BlockManager{
		files: make(map[FileID]*File),
		blocks: make(map[BlockID]*BlockHeader),
	}
}

func (b *BlockManager) AddFile(file *File) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.files[file.ID] = file
}

func (b *BlockManager) DeleteFile(id FileID) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, exists := b.files[id]; exists {
		delete(b.files, id)
	}
}

func (b *BlockManager) GetFile(id FileID) *File{
	b.lock.RLock()
	defer b.lock.RUnlock()
	if file, exists := b.files[id]; exists {
		return file
	} else {
		return nil
	}
}

func (b *BlockManager) ListFile() []*File {
	var files []*File
	for _, file := range(b.files) {
		files = append(files, file)
	}
	return files
}

func (b *BlockManager) AddBlock(id FileID, block *BlockHeader) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if file, exists := b.files[id]; exists {
		// TODO: ensure thread safe for blocks.
		file.Blocks[block.BlockID] = block
		file.Replications = len(file.Blocks)
		block.FileID = id
		b.blocks[block.BlockID] = block
	}
}

func (b *BlockManager) DeleteBlock(id FileID, blockID BlockID) {
	// TODO: ensure thread safe for deletion.
	b.lock.Lock()
	defer b.lock.Unlock()
	if file, exists := b.files[id]; exists {
		if _, bexists := file.Blocks[blockID]; bexists {
			delete(file.Blocks, blockID)
			file.Replications = len(file.Blocks)
		}
	}
	if _, exists := b.blocks[blockID]; exists {
		delete(b.blocks, blockID)
	}
}

func (b *BlockManager) GetBlock(blockID BlockID) *BlockHeader {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if block, exists := b.blocks[blockID]; exists {
		return block
	}
	return nil
}

func (b *BlockManager) HasBlock(blockID BlockID) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if _, exists := b.blocks[blockID]; exists {
		return true
	}
	return false
}

func (b *BlockManager) UpdateFileStatus(id FileID, status FileStatus) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if file, exists := b.files[id]; exists {
		file.Status = status
		b.files[id] = file
	}
}