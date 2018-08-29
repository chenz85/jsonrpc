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
	// type of func object
	rft reflect.Type
	// param type list
	rfpt []reflect.Type
}

func (m *RPCMethod) Invoke() (result interface{}, err object.Err) {
	if m.rft.NumIn() != 0 {
		err = object.ErrMethod_ParamsNumNotMatch
	} else {
		var result_vals = m.rf.Call(nil)
		result = m.return_values(result_vals)
	}
	return
}
func (m *RPCMethod) InvokeA(params []interface{}) (result interface{}, err object.Err) {
	if m.rft.NumIn() != len(params) {
		err = object.ErrMethod_ParamsNumNotMatch
	} else {
		var param_vals = make([]reflect.Value, len(params))
		for i, p := range params {
			if reflect.TypeOf(p) != m.rfpt[i] {
				err = object.ErrMethod_WrongParamsType
				return
			}
			param_vals[i] = reflect.ValueOf(p)
		}
		var result_vals = m.rf.Call(param_vals)
		result = m.return_values(result_vals)
	}
	return
}

func (m *RPCMethod) return_values(vals []reflect.Value) (result interface{}) {
	if result_num := m.rf.Type().NumOut(); result_num == 0 {
		result = nil
	} else if result_num == 1 {
		result = vals[0].Interface()
	} else {
		var results = make([]interface{}, result_num)
		for i, rv := range vals {
			results[i] = rv.Interface()
		}
		result = results
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
		var method = &RPCMethod{
			n:   name,
			rf:  reflect_func,
			rft: reflect_func.Type(),
		}

		method.rfpt = make([]reflect.Type, method.rft.NumIn())
		for i := 0; i < method.rft.NumIn(); i++ {
			method.rfpt[i] = method.rft.In(i)
		}

		rpc_mapper[name] = method
	}
	return
}

func get_method(name string) (method *RPCMethod, exist bool) {
	method, exist = rpc_mapper[name]
	return
}
