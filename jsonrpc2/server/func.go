package server

import (
	"reflect"
)

// rpc func data
type RPCFunc struct {
	// register name
	n string
	// func object
	rf reflect.Value
}

var (
	rpc_mapper map[string]*RPCFunc
)

func map_rpc_func(name string, reflect_func reflect.Value) (err error) {
	if _, ex := rpc_mapper[name]; ex {
		err = Error_Server_DuplicatedRPCFunc
	} else {
		rpc_mapper[name] = &RPCFunc{
			n:  name,
			rf: reflect_func,
		}
	}
	return
}
