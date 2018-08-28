package server_test

import (
	"testing"

	"github.com/czsilence/jsonrpc/jsonrpc2/server"
)

func TestHttpServer(t *testing.T) {
	server.StartHttpServer("", 9002, "rpc")
}
