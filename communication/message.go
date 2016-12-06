package communication

import "github.com/JetMuffin/whalefs/types"

type RegistrationMessage struct {
	Addr string
}

type HeartbeatMessage struct {
	NodeID 	types.NodeID
}

type HeartbeatResponse struct {

}