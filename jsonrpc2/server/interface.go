package server

type JSONRPCServer interface {
	Serve()
	RegisterMethod(name string, method interface{}) (err error)
}
