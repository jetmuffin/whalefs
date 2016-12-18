package master

import (
	. "github.com/JetMuffin/whalefs/types"
	"sync"
)

type BlockManager struct {
	files			map[FileID]*File
	blockSize		int
	blockReplication	int
	blobQueue 		chan *Blob
	lock 			sync.RWMutex
}

func NewBlockManager(blockSize int, blockReplication int) *BlockManager {
	return &BlockManager{
		files: make(map[FileID]*File),
		blockSize: blockSize,
		blockReplication: blockReplication,
		blobQueue: make(chan *Blob),
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
		file.Blocks = append(file.Blocks, block)
		b.files[id] = file
	}
}

func (b *BlockManager) UpdateFileStatus(id FileID, status FileStatus) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if file, exists := b.files[id]; exists {
		file.Status = status
		b.files[id] = file
	}
}