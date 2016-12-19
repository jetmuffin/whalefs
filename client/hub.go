package client

import (
"fmt"
	comm "github.com/JetMuffin/whalefs/communication"
	"encoding/json"
	"time"
)

type Hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection

	chunkAddr string
}

func NewHub(chunkAddr string) Hub {
	return Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
		connections: make(map[*Connection]bool),
		chunkAddr: chunkAddr,
	}
}

// Adds a connection to the connection map
func (h *Hub) Register(c *Connection) {
	h.register <- c
}

// Removes a connection from the connection map
func (h *Hub) Unregister(c *Connection) {
	h.unregister <- c
}

// Hub's main loop handles commands for the connection map
func (h *Hub) Run() {
	for {
		select {
		// Adds a connection
		case c := <-h.register:
			fmt.Println("Connect")
			h.connections[c] = true

			client, _ := comm.NewRPClient(h.chunkAddr, 5 * time.Second)
			var t int64
			client.Connection.Call("ChunkRPC.Time", new(interface{}), &t)
			msg, _ := json.Marshal(&comm.WsInitialMessage{
				Type: "init",
				NodeAddr: h.chunkAddr,
				NodeTime: t,
			})
			c.send <- msg
		// Removes a connection
		case c := <-h.unregister:
			fmt.Println("Disconnect")
			delete(h.connections, c)
			close(c.send)
		// Sends a mesage to each connected client
		case m := <-h.broadcast:
			fmt.Printf("Broadcasting: %s\n", m)
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}
