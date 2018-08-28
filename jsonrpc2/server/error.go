package server

import "errors"

var (
	Error_Server_InvalidRPCFunc    = errors.New("invalid rpc func")
	Error_Server_DuplicatedRPCFunc = errors.New("duplicated rpc func")
)
