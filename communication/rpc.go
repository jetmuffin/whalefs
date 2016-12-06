package communication

import (
	"net/rpc"
	"io"
	"net/rpc/jsonrpc"
)

type RPCServer struct {
	codec 		  rpc.ServerCodec
	lastServiceMethod string
	lastSeq		  uint64
}

func NewRPCServer(sock io.ReadWriteCloser) *RPCServer {
	return &RPCServer{jsonrpc.NewServerCodec(sock), "", 0}
}

func (server *RPCServer) ReadHeader() (string, error) {
	var r rpc.Request
	if err := server.codec.ReadRequestHeader(&r); err != nil {
		return "", err
	}
	server.lastServiceMethod = r.ServiceMethod
	server.lastSeq = r.Seq
	return r.ServiceMethod, nil
}

func (server *RPCServer) ReadBody(obj interface{}) error {
	return server.codec.ReadRequestBody(obj)
}

func (server *RPCServer) Error(s string) error {
	var r rpc.Response
	r.ServiceMethod = server.lastServiceMethod
	r.Seq = server.lastSeq
	r.Error = s
	return server.codec.WriteResponse(&r, nil)
}

func (server *RPCServer) Send(obj interface{}) error {
	var r rpc.Response
	r.ServiceMethod = server.lastServiceMethod
	r.Seq = server.lastSeq
	r.Error = ""
	return server.codec.WriteResponse(&r, obj)
}