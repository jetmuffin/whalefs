package main

import (
	"github.com/JetMuffin/whalefs/cmd"
	log "github.com/Sirupsen/logrus"
	"github.com/JetMuffin/whalefs/master"
	"github.com/JetMuffin/whalefs/chunk"
)

func main() {
	var debug bool

	cli := cmd.App()
	cli.Global(func(flag cmd.Flags) {
		flag.BoolVar(&debug, "debug", false, "Show debug logs")
	})

	cli.Command("master", "Run mater node", func(flag cmd.Flags) {
		configPath := flag.String("config", "./conf/whale.conf", "path to the master config file")
		flag.Parse()

		config, err := cmd.NewConfig(*configPath)
		if err != nil {
			log.Fatalf("Unexpected error when read config file: %v", err)
		}

		m := master.NewMaster(config)
		m.Run()

		<-make(chan bool)
	})

	cli.Command("chunkserver", "Run chunk server node", func(flag cmd.Flags) {
		configPath := flag.String("config", "./conf/whale.conf", "path to the chunkserver config file")
		flag.Parse()

		config, err := cmd.NewConfig(*configPath)
		if err != nil {
			log.Fatalf("Unexpected error when read config file: %v", err)
		}

		c := chunk.NewChunkServer(config)
		c.Run()

		<-make(chan bool)
	})

	cli.Run()
}
