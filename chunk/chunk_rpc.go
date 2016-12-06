package chunk

import (
	"net"
	comm "github.com/JetMuffin/whalefs/communication"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

func runChunkRPC(c net.Conn, chunk *ChunkServer) {
	server := comm.NewRPCServer(c)
	defer c.Close()

	method, err := server.ReadHeader()
	if err != nil {
		log.WithField("err", err).Error("unable read method from request.")
		return
	}

	switch method {
	}
}

func (chunk *ChunkServer) RunRPC() {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(chunk.RPCPort))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("RPC Server listen on :%v.", strconv.Itoa(chunk.RPCPort))
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go runChunkRPC(conn, chunk)
	}
}