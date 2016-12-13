package master

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

func (m *Master) Monitor() {
	for {
		for id, node := range m.chunks {
			if node.IsHealthy() && node.HeartbeatDuration() > m.heartbeatCheckInterval {
				log.WithField("Health", node.Heath).Errorf("node %v lost heartbeat.", id)
			}
		}
		time.Sleep(1 * time.Second)
	}
}