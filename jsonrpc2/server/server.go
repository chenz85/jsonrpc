package server

import (
	"log"
	"reflect"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

// process request data, and return response object.
// respnose object is nil if err is not nil.
func ProcessRequest(data string) (resp object.Response, err object.Err) {
	log.Printf("data: %+v\n", string(data))
	err = object.SimpleError(100, "test error")
	return
}

// register func to rpc server
// f MUST be a func
func HandleFunc(name string, f interface{}) (err error) {
	rf := reflect.ValueOf(f)
	if !rf.IsValid() || rf.IsNil() || rf.Kind() != reflect.Func {
		err = Error_Server_InvalidRPCFunc
	} else {
		map_rpc_func(name, rf)
	}
	return
}
