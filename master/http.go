package master

import (
	"net/http"
	"html/template"
	log "github.com/Sirupsen/logrus"
	comm "github.com/JetMuffin/whalefs/communication"
	"strconv"
	. "github.com/JetMuffin/whalefs/types"
	"io/ioutil"
	"sort"
	"time"
	"io"
	"bytes"
)

type HTTPServer struct {
	Host 	  	string
	Port 	  	int
	blockManager 	*BlockManager
	nodeManager     *NodeManager

}

func (server *HTTPServer) Addr() string {
	return server.Host + ":" + strconv.Itoa(server.Port)
}

func (server *HTTPServer) AddrWithScheme() string {
	return "http://" + server.Addr()
}

func (server *HTTPServer) upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, err := template.ParseFiles("static/upload.html")
		if err != nil {
			log.Errorf("Unable to render templates: %v", err)
			return
		}
		data := struct {
			Files []*File
		}{
			Files: server.blockManager.ListFile(),
		}
		t.Execute(w, data)
	} else {
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Errorf("Parse form file error: %v", err)
			return
		}
		defer file.Close()

		if err != nil {
			log.Errorf("Cannot read bytes from uploaded file: %v", err)
			return
		}

		data, err := ioutil.ReadAll(file)
		fileMeta := NewFile(header.Filename, int64(len(data)))
		server.blockManager.AddFile(fileMeta)

		chunks := server.nodeManager.LeastBlocksNodes()
		// No enough chunks to store replications.
		if len(chunks) == 0 {
			log.Error("Cannot write block: no chunk server available.")
			return
		}

		node := server.nodeManager.GetNode(chunks[0])
		block := NewBlock(header.Filename, data, int64(len(data)), node.ID)
		var checksum string

		// TODO: handle the situation that this node is down.
		server.blockManager.AddBlock(fileMeta.ID, block.Header)
		client, err := comm.NewRPClient(node.Addr, 5 * time.Second)
		if err != nil {
			log.Errorf("Cannot connect to node %v: %v", node.Addr, err)
		}
		client.Connection.Call("ChunkRPC.Write", block, &checksum)
		log.WithField("checksum", checksum).Infof("Write block %v successful", block.ID)

		http.Redirect(w, r, "/upload", 301)
	}
}

func (server *HTTPServer) nodes(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/nodes.html")
	if err != nil {
		log.Errorf("Unable to render templates: %v", err)
		return
	}

	data := struct {
		Nodes []*Node
	}{
		Nodes: server.nodeManager.ListNode(),
	}
	t.Execute(w, data)
}

func (server *HTTPServer) download(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileId := FileID(r.Form.Get("id"))

	// TODO: check if the id is illegal or not
	file := server.blockManager.GetFile(fileId)

	var blocks []*BlockHeader
	for _, block := range(file.Blocks) {
		blocks = append(blocks, block)
	}
	sort.Stable(SortBlockByFunc(func(block *BlockHeader) int {
		return server.nodeManager.GetNode(block.Chunk).Connections
	}, blocks))

	// TODO: if no nodes available
	block := blocks[0]
	log.Info(block)
	node := server.nodeManager.GetNode(block.Chunk)
	log.WithField("addr", node.Addr).Infof("Try to read block from node %v", node.ID)

	client, err := comm.NewRPClient(node.Addr, 5 * time.Second)
	if err != nil {
		log.Errorf("Cannot connect to node %v: %v", node.Addr, err)
		return
	}

	var message comm.BlockMessage
	client.Connection.Call("ChunkRPC.Read", block.BlockID, &message)

	// TODO: check the checksum received.
	log.WithFields(log.Fields{
		"checksum": message.Checksum,
		"length": len(message.Data),
	}).Infof("Receive block data from node %v.", node.ID)

	b := bytes.NewBuffer(message.Data)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=" + file.Name)
	io.Copy(w, b)
}

func (server *HTTPServer) ListenAndServe()  {
	http.HandleFunc("/upload", server.upload)
	http.HandleFunc("/nodes", server.nodes)
	http.HandleFunc("/download", server.download)
	log.WithFields(log.Fields{"host": server.Host, "port": server.Port}).Info("HTTP Server start listening.")

	go http.ListenAndServe(server.Addr(), nil)
}

func NewHTTPServer(host string, port int, blockManager *BlockManager, nodeManager *NodeManager) *HTTPServer {
	return &HTTPServer{
		Host: host,
		Port: port,
		blockManager: blockManager,
		nodeManager: nodeManager,
	}
}