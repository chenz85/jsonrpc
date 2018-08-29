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
	server.RegisterMethod("echo", func(val string) string {
		return fmt.Sprintf("you say: %s", val)
	})
	server.StartHttpServer("", 9002, "rpc")
}
