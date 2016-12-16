package chunk

import (
	"testing"
	"github.com/JetMuffin/whalefs/cmd"
	"os"
)

var (
	config, _ = cmd.NewConfig("../conf/whale.conf")
)
func TestNewChunkServer(t *testing.T) {
	c := NewChunkServer(config)
	defer os.RemoveAll(c.store.DataDir)
	c.Run()
}
