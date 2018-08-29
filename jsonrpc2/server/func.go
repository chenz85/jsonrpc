package server

import (
	"reflect"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

// rpc method data
type RPCMethod struct {
	// register name
	n string
	// func object
	rf reflect.Value
}

func (m *RPCMethod) Invoke() (result interface{}, err object.Err) {
	if m.rf.Type().NumIn() != 0 {
		err = object.ErrMethod_ParamsNumNotMatch
	}
	return
}
func (m *RPCMethod) InvokeA(params []interface{}) (result interface{}, err object.Err) {
	if m.rf.Type().NumIn() != len(params) {
		err = object.ErrMethod_ParamsNumNotMatch
	}
	return
}

var (
	rpc_mapper map[string]*RPCMethod = make(map[string]*RPCMethod)
)

func map_rpc_method(name string, reflect_func reflect.Value) (err error) {
	if _, ex := rpc_mapper[name]; ex {
		err = Error_Server_DuplicatedRPCMethod
	} else {
		rpc_mapper[name] = &RPCMethod{
			n:  name,
			rf: reflect_func,
		}
	}
	return
}

func get_method(name string) (method *RPCMethod, exist bool) {
	method, exist = rpc_mapper[name]
	return
}
