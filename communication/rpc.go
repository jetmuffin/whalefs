package communication

import (
	"net/rpc"
	"time"
	"net"
)

type RPCClient struct {
	Connection *rpc.Client
}

func NewRPClient(dsn string, timeout time.Duration) (*RPCClient, error) {
	connection, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		return nil, err
	}
	return &RPCClient{Connection: rpc.NewClient(connection)}, nil
}