package client


import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/JetMuffin/whalefs/cmd"
	"github.com/JetMuffin/whalefs/communication"
	log "github.com/Sirupsen/logrus"
	"time"
)

type Client struct {
	masterAddr 	string
	chunkAddr 	string
	hub 		Hub
}

func NewClient(config *cmd.Config) *Client {
	c := &Client {
		masterAddr: config.String("master_addr"),
	}
	client, err := communication.NewRPClient(c.masterAddr, 5 * time.Second)
	if err != nil {
		log.Fatalf("Cannot connect to master: %v", err)
	}
	client.Connection.Call("MasterRPC.ConnectChunk", new(interface{}), &c.chunkAddr)
	c.hub = NewHub(c.chunkAddr)
	return c
}

func(c *Client) wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket. TODO: check origin
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	conn := NewConnection(ws)
	c.hub.Register(conn)
	defer func() { c.hub.Unregister(conn) }()
	go conn.Writer()
	conn.Reader(c.hub)
}

func (c *Client) Run() {
	go c.hub.Run()
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.HandleFunc("/ws", c.wsHandler)
	if err := http.ListenAndServe(":4000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
