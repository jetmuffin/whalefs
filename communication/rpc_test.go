package communication

import (
	"testing"
	"net"
	"sync"
	"fmt"
)

var (
	server 	*RPCServer
	once 	sync.Once
	stop 	bool
)

func setupTestRPCServer() {
	listener, _ := net.Listen("tcp", ":8888")
	for !stop {
		peer, _ := listener.Accept()
		server = NewRPCServer(peer)
		fmt.Println(stop)
	}
}

func TestRPCServer(t *testing.T) {
	//go once.Do(setupTestRPCServer)
	//client, err := rpc.Dial("tcp", "127.0.0.1:8888")
	//if err != nil {
	//	t.Fatal("dialing", err)
	//}
	//defer client.Close()
	//
	//var resp string
	//client.Call("test_method", "args", &resp)
	//method, err := server.ReadHeader()
	//if err != nil {
	//	t.Errorf("error when read request head: %v", err)
	//}
	//if method != "test_method" {
	//	t.Error("read method error")
	//}
}
