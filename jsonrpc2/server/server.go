package server

import (
	"log"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

// process request data, and return response object.
// respnose object is nil if err is not nil.
func ProcessRequest(data string) (resp object.Response, err object.Err) {
	log.Printf("data: %+v\n", string(data))
	err = object.SimpleError(100, "test error")
	return
}
