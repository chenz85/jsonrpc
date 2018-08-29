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
	RequestParamTypeMap
)

func (t RequestParamType) String() string {
	switch t {
	case RequestParamTypeNone:
		return "none"
	case RequestParamTypeArray:
		return "array"
	case RequestParamTypeMap:
		return "map"
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

	MapParams() map[string]interface{}
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
	// map params, ref of params
	params_map map[string]interface{}
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

func (e *request_object) MapParams() map[string]interface{} {
	return e.params_map
}

func (e *request_object) ArrayParams() []interface{} {
	return e.params_arr
}

func (e *request_object) Parse(obj map[string]interface{}) (err Err) {
	if json_rpc, ex := obj["jsonrpc"]; !ex || json_rpc != "2.0" {
		return ErrParse_MissingField_jsonrpc
	} else {
		e.jsonrpc, _ = json_rpc.(string)
	}

	if method_val, ex := obj["method"]; !ex {
		return ErrParse_InvalidField_Method
	} else if method, ok := method_val.(string); !ok || method == "" {
		return ErrParse_InvalidField_Method
	} else if strings.HasPrefix(method, "rpc.") {
		return ErrParse_Method_ReservedFunc
	} else {
		e.method = method
	}

	if id_val, ex := obj["id"]; !ex {
		// notification
	} else {
		e.id = id_val
	}

	// TODO: params

	return
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
