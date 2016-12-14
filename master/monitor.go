package master

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

func (m *Master) checkChunkHealth() {
	for id, node := range m.nodeManager.chunks {
		if node.IsHealthy() && node.HeartbeatDuration() > m.heartbeatCheckInterval {
			log.WithField("Health", node.Heath).Errorf("node %v lost heartbeat.", id)
			m.nodeManager.LostNode(node)
		}
	}
}

func (m *Master) Monitor() {
	for {
		m.checkChunkHealth()
		time.Sleep(1 * time.Second)
	}
}