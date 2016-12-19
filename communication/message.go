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
	Blocks 		[]BlockID
	Utilization	int64
	Timestamp 	time.Time
}

type BlockMessage struct {
	Data []byte
	Checksum string
}

// HeartbeatResponse send from master to chunk node to do some action to
// keep consistency of cluster.
type HeartbeatResponse struct {
	DeadBlocks	[]BlockID
	SyncBlocks 	[]*SyncBlock
}

type SyncDoneResponse struct {

}

type WsInitialMessage struct {
	Type 	   string     `json:"type"`
	NodeAddr   string     `json:"node_addr"`
	NodeTime   int64      `json:"node_time"`
}

type WsMessage struct {
	Type string  `json:"type"`
}