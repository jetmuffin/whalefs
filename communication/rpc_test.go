package communication

import (
	"time"
	"log"
	"testing"
	"net"
	"net/rpc"
)

var (
	c   *RPCClient
	err error

	dsn       = "localhost:9876"
)

func TestNewRPClient(t *testing.T) {
	go func() {
		l, e := net.Listen("tcp", ":9876")
		if e != nil {
			log.Fatal("listen error:", e)
		}
		rpc.Accept(l)
	} ()

	c, err = NewRPClient(dsn, time.Second*3)
}