package chunk

import (
	"io/ioutil"
	"os"
	"path"
	"io"
	"hash/crc32"
	"fmt"

	"github.com/JetMuffin/whalefs/types"
	"path/filepath"
	log "github.com/Sirupsen/logrus"
)

// BlockStore is a block storage manage object.
type BlockStore struct {
	DataDir string
}

// NewBlockStore returns a new store with given data directory.
func NewBlockStore(directory string) *BlockStore {
	blockStore := &BlockStore{
		DataDir: directory,
	}

	// if directory does not exist, create it.
	_, err := os.Stat(blockStore.DataDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(blockStore.BlocksDirectory(), 0700)
		err = os.MkdirAll(blockStore.MetaDirectory(), 0700)
		if err != nil {
			log.Fatal("Cannot create store data directory, please check the chunk_data_dir configuration.")
		}
	}
	return blockStore
}

// BlocksDirectory returns blocks storage directory path.
func (store *BlockStore) BlocksDirectory() string {
	return path.Join(store.DataDir, "blocks")
}

// BlockStoragePath returns absolute storage path of specified block.
func (store *BlockStore) BlockStoragePath(block types.BlockID) string {
	return path.Join(store.BlocksDirectory(), string(block))
}

// BlockSize return bytes size of block.
func (store *BlockStore) BlockSize(block types.BlockID) (int64, error){
	fileInfo, err := os.Stat(store.BlockStoragePath(block))
	if err != nil {
		return -1, err
	}
	return fileInfo.Size(), nil
}

// ListBlocks returns all blocks' id in current block store data directory.
func (store *BlockStore) ListBlocks() ([]types.BlockID, error) {
	files, err := ioutil.ReadDir(store.BlocksDirectory())
	if err != nil {
		return nil, err
	}

	var blocks []types.BlockID
	for _, f := range files {
		blocks = append(blocks, types.BlockID(f.Name()))
	}
	return blocks, nil
}

// WriteBlock write block data from given io reader to storage file and return its checksum.
func (store *BlockStore) WriteBlock(block types.BlockID, size int64, r io.Reader) (string, error) {
	file, err := os.Create(store.BlockStoragePath(block))
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := crc32.NewIEEE()
	_, err = io.CopyN(file, io.TeeReader(r, hash), size)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(hash.Sum32()), nil
}

// ReadBlock read block data from storage file to io writer.
func (store *BlockStore) ReadBlock(block types.BlockID, w io.Writer) error {
	file, err := os.Open(store.BlockStoragePath(block))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}
	return nil
}

// Utilization returns total disk usage of store data directory.
func (store *BlockStore) Utilization() int64 {
	var size int64
	err := filepath.Walk(store.DataDir, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return -1
	} else {
		return size
	}
}