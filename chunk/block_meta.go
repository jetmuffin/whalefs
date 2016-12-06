package chunk

import (
	"path"
	"os"
	"hash/crc32"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/JetMuffin/whalefs/types"
)

// MetaDirectory returns block's metadata storage path.
func (store *BlockStore) MetaDirectory() string {
	return path.Join(store.DataDir, "meta")
}

// BlockCheckSumPath returns absolute path of specified block's checksum data.
func (store *BlockStore) BlockCheckSumPath(block types.BlockID) string {
	return path.Join(store.MetaDirectory(), string(block) + ".crc32")
}

// BlockCheckSum returns block's checksum data, calculate using crc32.
func (store *BlockStore) BlockCheckSum(block types.BlockID) (string, error) {
	file, err := os.Open(store.BlockStoragePath(block))
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := crc32.NewIEEE()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprint(hash.Sum32()), nil
}

// WriteChecksum writes block's checksum data to its storage path.
func (store *BlockStore) WriteChecksum(block types.BlockID, s string) error {
	return ioutil.WriteFile(store.BlockCheckSumPath(block), []byte(s), os.ModePerm)
}

// ReadChecksum reads block's checksum data from its storage path.
func (store *BlockStore) ReadChecksum(block types.BlockID) (string, error) {
	bytes, err := ioutil.ReadFile(store.BlockCheckSumPath(block))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}