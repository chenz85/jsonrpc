package server

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

// process request data, and return response object.
// respnose object is nil if err is not nil.
func HandleRequest(data []byte) (resp_data []byte) {
	log.Printf("req data: %s\n", string(data))
	if data[0] == '[' {
		// batch request
		var resp_arr = make([]object.Response, 0, 1)
		var single_error bool
		var objs = make([]interface{}, 0, 1)
		if je := json.Unmarshal(data, &objs); je != nil {
			log.Println("parse request failed:", je)
			resp_arr = append(resp_arr, process_error(object.ErrParse))
			single_error = true
		} else if len(objs) == 0 {
			resp_arr = append(resp_arr, process_error(object.ErrInvalidRequest))
			single_error = true
		} else {
			for _, obj_val := range objs {
				obj, _ := obj_val.(map[string]interface{})
				if req, pe := object.ParseRequest(obj); pe != nil {
					resp_arr = append(resp_arr, process_error(pe))
				} else if resp := process_request(req); resp != nil {
					resp_arr = append(resp_arr, resp)
				}
			}
		}

		if resp_cnt := len(resp_arr); resp_cnt == 0 {
			// no response
			resp_data = nil
		} else if single_error {
			resp_data = resp_arr[0].JsonObject().ToJson()
		} else {
			var obj_arr = make([]object.JsonObject, resp_cnt)
			for i, resp := range resp_arr {
				obj_arr[i] = resp.JsonObject()
			}
			resp_data = object.JsonObjectArrayToJson(obj_arr)
		}
	} else {
		// single request
		var resp object.Response
		var obj = make(map[string]interface{})
		if je := json.Unmarshal(data, &obj); je != nil {
			log.Println("parse request failed:", je)
			resp = process_error(object.ErrParse)
		} else if req, pe := object.ParseRequest(obj); pe != nil {
			resp = process_error(pe)
		} else {
			resp = process_request(req)
		}

		if resp != nil {
			resp_data = resp.JsonObject().ToJson()
		}
	}
	return
}

func process_error(err object.Err) (resp object.Response) {
	resp, _ = object.NewResponse(nil, err, nil)
	return
}

func process_request(req object.Request) (resp object.Response) {
	var result interface{}
	var err object.Err
	log.Printf("req: %+v\n", req)
	if req == nil {
		err = object.ErrInvalidRequest
	} else if method, ex := get_method(req.Method()); !ex {
		err = object.ErrMethodNotFound
	} else {
		switch req.ParamType() {
		case object.RequestParamTypeNone:
			result, err = method.Invoke()
		case object.RequestParamTypeArray:
			result, err = method.InvokeA(req.ArrayParams())
		}

	}
	if _resp, re := object.NewResponse(result, err, req.Id()); re != nil {
		err = object.ErrInternalError
	} else if err != nil || !req.IsNotification() {
		resp = _resp
	}
	return
}

// register func to rpc server
// f MUST be a func
func RegisterMethod(name string, f interface{}) (err error) {
	rf := reflect.ValueOf(f)
	if !rf.IsValid() || rf.IsNil() || rf.Kind() != reflect.Func {
		err = Error_Server_InvalidRPCMethod
	} else {
		map_rpc_method(name, rf)
	}
	return
}
