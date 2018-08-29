package server

type JSONRPCServer interface {
	// register rpc method
	RegisterMethod(name string, method interface{}) (err error)
	// register rpc method mapper。method registered by RegisterMethod() before will be clear.
	// call NewRPCMethodMapper() to get a new mapper
	RegisterMapper(mapper RPCMethodMapper)
	// start rpc serve
	Serve()
}
