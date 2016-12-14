package master

import (
	"net"
	"strconv"
	comm "github.com/JetMuffin/whalefs/communication"
	log "github.com/Sirupsen/logrus"
)

func runMasterRPC(c net.Conn, master *Master) {
	server := comm.NewRPCServer(c)
	defer c.Close()

	method, err := server.ReadHeader()
	if err != nil {
		log.WithField("err", err).Error("unable read method from request.")
		return
	}

	switch method {
	case "Register":
		var message comm.RegistrationMessage
		if err := server.ReadBody(&message); err != nil {
			log.WithField("err", err).Error("unable to parse message body.")
			return
		}

		// TODO: re-register node when restart chunkserver.
		nodeID := master.nodeManager.RegisterChunkNode(message.Addr)
		server.Send(&nodeID)
		log.WithField("nodeID", nodeID).Infof("Chunk node registered at %v, %v nodes totally.", message.Addr, len(master.nodeManager.chunks))
	case "Heartbeat":
		var message comm.HeartbeatMessage
		if err := server.ReadBody(&message); err != nil {
			log.WithField("err", err).Error("unable to parse message body.")
			return
		}
		log.Debugf("Heartbeat from node(%v).", message.NodeID)

		// Update node information by heartbeat
		master.nodeManager.UpdateNodeWithHeartbeat(message)

		// TODO send response to heartbeat
		var resp comm.HeartbeatResponse
		server.Send(&resp)
	}
}

// RunRPC setup a RPC server on master node.
func (m *Master) RunRPC() {
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(m.RPCPort))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("RPC Server listen on :%v.", strconv.Itoa(m.RPCPort))
	for {
		peer, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go runMasterRPC(peer, m)
	}
}