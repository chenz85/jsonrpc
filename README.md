# Go implement of JSON-RPC

## What is JSON-RPC

https://www.jsonrpc.org/specification

# Server

* start http server

    ```go
    import 	"github.com/czsilence/jsonrpc/jsonrpc2/server"

	svr := server.NewHttpServer("127.0.0.1", 9002, "rpc")
	svr.RegisterMethod("echo", func(val string) string {
		return fmt.Sprintf("you say: %s", val)
	})
	svr.Serve()
    ```

* start socket server

    ```go
    import 	"github.com/czsilence/jsonrpc/jsonrpc2/server"

    // listen on a tcp host
    svr := server.NewSocketServer("tcp4", ":9003")
    // or linsten on a unix socket. The path to rpc.socket must be exist
	// svr := server.NewSocketServer("unix", "path/to/rpc.socket")
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
