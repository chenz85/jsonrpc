package server

import (
	"errors"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

var (
	Error_Server_InvalidRPCMethod    = errors.New("invalid rpc method")
	Error_Server_DuplicatedRPCMethod = errors.New("duplicated rpc method")
)

func Raise(err error) {
	panic(err)
}

func RaiseError(code int, message string, data interface{}) {
	panic(object.Error(code, message, data))
}

func RaiseErrorObject(err object.Err) {
	panic(err)
}
