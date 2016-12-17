package chunk

import (
	"testing"
	"github.com/JetMuffin/whalefs/cmd"
	"github.com/JetMuffin/whalefs/master"
	"os"
)

var (
	configPath = "../conf/whale.conf"
	config, _ = cmd.NewConfig(configPath)
	m = master.NewMaster(config)
)

func init() {
	m.ListenRPC()
}

func TestNewChunkServer(t *testing.T) {
	c := NewChunkServer(config)
	defer os.RemoveAll(c.store.DataDir)
	c.Run()
}
