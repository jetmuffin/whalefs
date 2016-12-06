package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/JetMuffin/whalefs/master"
)

func main() {
	log.SetLevel(log.DebugLevel)

	server := master.NewHTTPServer("127.0.0.1", 8888)
	go server.ListenAndServe()

	exit := make(chan bool)
	<-exit
}
