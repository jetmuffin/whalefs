package chunk

import (
	"path"
	"os"
	"hash/crc32"
	"fmt"
	"io"
	"io/ioutil"

	. "github.com/JetMuffin/whalefs/types"
	"encoding/json"
)

// MetaDirectory returns block's metadata storage path.
func (store *BlockStore) MetaDirectory() string {
	return path.Join(store.DataDir, "meta")
}

// BlockMetaPath returns absolute path of specified block's checksum data.
func (store *BlockStore) BlockMetaPath(blockID BlockID) string {
	return path.Join(store.MetaDirectory(), string(blockID) + ".crc32")
}

func (store *BlockStore) WriteMeta(block *BlockHeader, checksum string) error {
	block.Checksum = checksum
	jsons, err := json.Marshal(block)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.BlockMetaPath(block.BlockID), jsons, os.ModePerm)
}

// BlockCheckSum returns block's checksum data, calculate using crc32.
func (store *BlockStore) BlockCheckSum(block BlockID) (string, error) {
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

// ReadChecksum reads block's checksum data from its storage path.
func (store *BlockStore) ReadMeta(block BlockID) (*BlockHeader, error) {
	bytes, err := ioutil.ReadFile(store.BlockMetaPath(block))
	var header *BlockHeader
	err = json.Unmarshal(bytes, &header)
	if err != nil {
		return nil, err
	}
	return header, nil
}