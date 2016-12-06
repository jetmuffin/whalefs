package master

import (
	"net/http"
	"html/template"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
	"io"
	"strconv"
)

type HTTPServer struct {
	Host string
	Port int
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
			log.WithField("err", err).Error("Unable to render templates.")
			return
		}
		t.Execute(w, nil)
	} else {
		file, handle, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			log.WithField("err", err).Error("Parse form file error.")
		}

		f, err := os.OpenFile(path.Join("upload", handle.Filename), os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			log.WithField("err", err).Error("Error create file.")
			return
		}
		defer f.Close()
		io.Copy(f, file)
		io.WriteString(w, "upload successful")
	}
}

func (server *HTTPServer) ListenAndServe()  {
	http.HandleFunc("/upload", server.upload)
	log.WithFields(log.Fields{"host": server.Host, "port": server.Port}).Info("HTTP Server start listening.")

	http.ListenAndServe(server.Addr(), nil)
}

func NewHTTPServer(host string, port int) *HTTPServer {
	return &HTTPServer{Host: host, Port: port}
}