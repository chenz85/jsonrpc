package server

import "errors"

var (
	Error_Server_InvalidRPCMethod    = errors.New("invalid rpc method")
	Error_Server_DuplicatedRPCMethod = errors.New("duplicated rpc method")
)
