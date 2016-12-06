package communication

import (
	"net/rpc"
	"io"
	"net/rpc/jsonrpc"
)

// RPCServer packaged rpc.ServerCodec, providing basic rpc methods to
// each node.
type RPCServer struct {
	codec 		  rpc.ServerCodec
	lastServiceMethod string
	lastSeq		  uint64
}

// NewRPCServer creates a jsonrpc server.
func NewRPCServer(sock io.ReadWriteCloser) *RPCServer {
	return &RPCServer{jsonrpc.NewServerCodec(sock), "", 0}
}

// ReadHeader reads method, sequence of rpc request.
func (server *RPCServer) ReadHeader() (string, error) {
	var r rpc.Request
	if err := server.codec.ReadRequestHeader(&r); err != nil {
		return "", err
	}
	server.lastServiceMethod = r.ServiceMethod
	server.lastSeq = r.Seq
	return r.ServiceMethod, nil
}

// ReadBody reads body of rpc request.
func (server *RPCServer) ReadBody(obj interface{}) error {
	return server.codec.ReadRequestBody(obj)
}

// Error response rpc client with custom errors.
func (server *RPCServer) Error(s string) error {
	var r rpc.Response
	r.ServiceMethod = server.lastServiceMethod
	r.Seq = server.lastSeq
	r.Error = s
	return server.codec.WriteResponse(&r, nil)
}

// Send response rpc client with normal response.
func (server *RPCServer) Send(obj interface{}) error {
	var r rpc.Response
	r.ServiceMethod = server.lastServiceMethod
	r.Seq = server.lastSeq
	r.Error = ""
	return server.codec.WriteResponse(&r, obj)
}