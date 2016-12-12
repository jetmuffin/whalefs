package communication

import (
	. "github.com/JetMuffin/whalefs/types"
	"time"
)

// RegistrationMessage is the first message send to master which includes
// chunk node's address information.
type RegistrationMessage struct {
	Addr string
}

// HeartbeatMessage is the heartbeat packet send to master, which includes
// node's metric.
// TODO add metric
type HeartbeatMessage struct {
	NodeID 		NodeID
	Addr 		string
	Timestamp 	time.Time
}

// HeartbeatResponse send from master to chunk node to do some action to
// keep consistency of cluster.
type HeartbeatResponse struct {

}