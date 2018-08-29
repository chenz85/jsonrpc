package main

import (
	"fmt"

	"github.com/czsilence/go/app"
	"github.com/czsilence/jsonrpc/jsonrpc2/server"
)

func main() {
	go start()
	app.HandleInterrupt()
}

func start() {
	svr := server.NewSocketServer("tcp4", ":9003")
	// svr := server.NewSocketServer("unix", "path/to/rpc.socket")
	svr.RegisterMethod("echo", func(val string) string {
		return fmt.Sprintf("you say: %s", val)
	})
	svr.Serve()
}
