package master

import (
	"testing"
	"github.com/JetMuffin/whalefs/cmd"
)

var (
	config, _ = cmd.NewConfig("../conf/whale.conf")
)

func TestNewMaster(t *testing.T) {
	NewMaster(config)
}
