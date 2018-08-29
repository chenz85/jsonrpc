# Go implement of JSON-RPC

## What is JSON-RPC

https://www.jsonrpc.org/specification

# Server

* start server

    ```go
    import 	"github.com/czsilence/jsonrpc/jsonrpc2/server"

    // http server
	svr := server.NewHttpServer("127.0.0.1", 9002, "rpc")
	svr.RegisterMethod("echo", func(val string) string {
		return fmt.Sprintf("you say: %s", val)
	})
	svr.Serve()
    ```

* test

    ```bash
    curl -d '{"jsonrpc": "2.0", "method": "echo", "params": ["42"], "id": 1}' http://127.0.0.1:9002/rpc
    ```

    will get response if success:
    ```bash
    {"id":1,"jsonrpc":"2.0","result":"you say: 42"}
    ```

# Client

## TBD
