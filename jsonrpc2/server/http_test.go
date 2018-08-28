package server_test

import (
	"fmt"
	"testing"

	"github.com/czsilence/jsonrpc/jsonrpc2/server"
)

func TestHttpServer(t *testing.T) {
	server.HandleFunc("echo", func(val string) string {
		return fmt.Sprintf("you say: %s", val)
	})
	server.StartHttpServer("", 9002, "rpc")
}
