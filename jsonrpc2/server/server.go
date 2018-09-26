package server

import (
	"encoding/json"
	"log"

	"github.com/czsilence/jsonrpc/jsonrpc2/object"
)

type BaseServer struct {
	method_mapper RPCMethodMapper
}

// process request data, and return response object.
// respnose object is nil if err is not nil.
func (bs *BaseServer) HandleRequest(data []byte) (resp_data []byte) {
	log.Printf("req data: %s\n", string(data))
	if data[0] == '[' {
		// batch request
		var resp_arr = make([]object.Response, 0, 1)
		var single_error bool
		var objs = make([]interface{}, 0, 1)
		if je := json.Unmarshal(data, &objs); je != nil {
			log.Println("parse request failed:", je)
			resp_arr = append(resp_arr, bs.process_error(object.ErrParse))
			single_error = true
		} else if len(objs) == 0 {
			resp_arr = append(resp_arr, bs.process_error(object.ErrInvalidRequest))
			single_error = true
		} else {
			for _, obj_val := range objs {
				obj, _ := obj_val.(map[string]interface{})
				if req, pe := object.ParseRequest(obj); pe != nil {
					resp_arr = append(resp_arr, bs.process_error(pe))
				} else if resp := bs.process_request(req); resp != nil {
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
			resp = bs.process_error(object.ErrParse)
		} else if req, pe := object.ParseRequest(obj); pe != nil {
			resp = bs.process_error(pe)
		} else {
			resp = bs.process_request(req)
		}

		if resp != nil {
			resp_data = resp.JsonObject().ToJson()
		}
	}
	return
}

func (bs *BaseServer) process_error(err object.Err) (resp object.Response) {
	resp, _ = object.NewResponse(nil, err, nil)
	return
}

func (bs *BaseServer) process_request(req object.Request) (resp object.Response) {
	defer bs.handle_request_panic(&resp)
	var result interface{}
	var err object.Err
	log.Printf("req: %+v\n", req)
	if req == nil {
		err = object.ErrInvalidRequest
	} else if method, ex := bs.method_mapper.Get(req.Method()); !ex {
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

func (bs *BaseServer) handle_request_panic(resp *object.Response) {
	if edata := recover(); edata == nil {
		// no error
	} else if err, ok := edata.(object.Err); ok && err != nil {
		*resp = bs.process_error(err)
	} else if err, ok := edata.(error); ok && err != nil {
		*resp = bs.process_error(object.SimpleError(-1, err.Error()))
	} else if err, ok := edata.(string); ok && err != "" {
		*resp = bs.process_error(object.SimpleError(-1, err))
	} else {
		*resp = bs.process_error(object.SimpleError(-1, "unknown error"))
		log.Println("unsupported error:", edata)
	}
}

func (bs *BaseServer) RegisterMethod(name string, method interface{}) (err error) {
	if bs.method_mapper == nil {
		bs.method_mapper = NewRPCMethodMapper()
	}

	return bs.method_mapper.RegisterMethod(name, method)
}

func (bs *BaseServer) RegisterMapper(mapper RPCMethodMapper) {
	bs.method_mapper = mapper
}
