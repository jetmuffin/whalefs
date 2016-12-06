package master

import (
	"testing"
	"net/http"
	"time"
)

var (
	server = NewHTTPServer("127.0.0.1", 9999)
)

func TestUpload(t *testing.T) {
	go server.ListenAndServe()

	time.Sleep(1000)
	res, err := http.Get(server.AddrWithScheme() + "/upload")
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("error when get upload page: %v.", err)
	}
}
