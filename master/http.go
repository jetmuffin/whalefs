package master

import (
	"net/http"
	"html/template"
	log "github.com/Sirupsen/logrus"
	"io"
	"strconv"
	. "github.com/JetMuffin/whalefs/types"
	"io/ioutil"
)

type HTTPServer struct {
	Host 	  string
	Port 	  int
	blobQueue chan *Blob
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
		t.Execute(w, nil)
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

		bytes, err := ioutil.ReadAll(file)
		blob := &Blob{
			Name: header.Filename,
			Length: int64(len(bytes)),
			Content: bytes,
		}
		server.blobQueue <- blob

		io.WriteString(w, "upload successful")
	}
}

func (server *HTTPServer) ListenAndServe()  {
	http.HandleFunc("/upload", server.upload)
	log.WithFields(log.Fields{"host": server.Host, "port": server.Port}).Info("HTTP Server start listening.")

	go http.ListenAndServe(server.Addr(), nil)
}

func NewHTTPServer(host string, port int, queue chan *Blob) *HTTPServer {
	return &HTTPServer{Host: host, Port: port, blobQueue: queue}
}