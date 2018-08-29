package object

import (
	"errors"
	"strings"
)

var (
	error_NoField_jsonrpc = errors.New("missing field: jsonrpc")
)

type RequestParamType int

const (
	RequestParamTypeInvalid RequestParamType = iota
	RequestParamTypeNone
	RequestParamTypeArray
)

func (t RequestParamType) String() string {
	switch t {
	case RequestParamTypeNone:
		return "none"
	case RequestParamTypeArray:
		return "array"
	default:
		return "invalid"
	}
}

type Request interface {
	json_object

	Id() interface{}
	IsNotification() bool

	Method() string
	ParamType() RequestParamType

	ArrayParams() []interface{}
}

type request_object struct {
	// A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
	jsonrpc string

	// A String containing the name of the method to be invoked. Method names that begin with the word
	// rpc followed by a period character (U+002E or ASCII 46) are reserved for rpc-internal methods
	// and extensions and MUST NOT be used for anything else.
	method string

	// A Structured value that holds the parameter values to be used during the invocation of the method.
	// This member MAY be omitted.
	params interface{}
	// type of params
	param_type RequestParamType
	// array params, ref of params
	params_arr []interface{}

	// An identifier established by the Client that MUST contain a String, Number, or NULL value if included.
	// If it is not included it is assumed to be a notification.
	// The value SHOULD normally not be Null [1] and Numbers SHOULD NOT contain fractional parts
	id interface{}
}

func (req *request_object) JsonObject() JsonObject {
	obj := JsonObject{}

	// TODO: init json obj

	return obj
}

func (e *request_object) Id() interface{} {
	return e.id
}

func (e *request_object) IsNotification() bool {
	return e.id == nil
}

func (e *request_object) Method() string {
	return e.method
}

func (e *request_object) ParamType() RequestParamType {
	return e.param_type
}

func (e *request_object) ArrayParams() []interface{} {
	return e.params_arr
}

func (e *request_object) Parse(obj map[string]interface{}) (err Err) {
	// jsonrpc
	if json_rpc, ex := obj["jsonrpc"]; !ex || json_rpc != "2.0" {
		return ErrParse_MissingField_jsonrpc
	} else {
		e.jsonrpc, _ = json_rpc.(string)
	}

	// method
	if method_val, ex := obj["method"]; !ex {
		return ErrParse_InvalidField_Method
	} else if method, ok := method_val.(string); !ok || method == "" {
		return ErrParse_InvalidField_Method
	} else if strings.HasPrefix(method, "rpc.") {
		return ErrParse_Method_ReservedFunc
	} else {
		e.method = method
	}

	// id
	if id_val, ex := obj["id"]; !ex {
		// notification
	} else {
		e.id = id_val
	}

	// params
	if params_val, ex := obj["params"]; !ex {
		e.param_type = RequestParamTypeNone
	} else if params_arr, ok := params_val.([]interface{}); ok {
		e.param_type = RequestParamTypeArray
		e.params, e.params_arr = params_arr, params_arr
	}

	return
}

func (e *request_object) Check() error {
	// TODO: data check

	if e.params == nil {
		e.param_type = RequestParamTypeNone
	}

	return nil
}

////////////////////////////////////////////////////////////////
func ParseRequest(obj map[string]interface{}) (Request, Err) {
	var req = &request_object{}
	if err := req.Parse(obj); err != nil {
		return nil, err
	} else {
		return req, nil
	}
}

// make a new request
func NewRequestA(method string, params []interface{}, id interface{}) (Request, error) {
	var req = &request_object{
		jsonrpc:    "2.0",
		method:     method,
		params:     params,
		param_type: RequestParamTypeArray,
		id:         id,
	}

	if err := req.Check(); err != nil {
		return nil, err
	} else {
		return req, nil
	}
}
