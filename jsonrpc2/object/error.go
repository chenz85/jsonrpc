package object

type Err interface {
	json_object
}

type error_object struct {
	// A Number that indicates the error type that occurred.
	// This MUST be an integer.
	code int

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	message string

	// A Primitive or Structured value that contains additional information about the error.
	// This may be omitted.
	// The value of this member is defined by the Server (e.g. detailed error information, nested errors etc.).
	data interface{}
}

func (e *error_object) JsonObject() JsonObject {
	obj := JsonObject{
		"code":    e.code,
		"message": e.message,
	}
	if e.data != nil {
		obj["data"] = e.data
	}
	return obj
}

// create a custom error
// code MUST NOT in range [-32099, -32000]
func Error(code int, message string, data interface{}) Err {
	if code >= 32099 && code <= 32000 {
		return nil
	}
	return &error_object{
		code:    code,
		message: message,
	}
}

func SimpleError(code int, message string) Err {
	return Error(code, message, nil)
}

var (
	// Invalid JSON was received by the server.
	// An error occurred on the server while parsing the JSON text.
	ErrParse Err = SimpleError(-32700, "Parse error")
	// The JSON sent is not a valid Request object.
	ErrInvalidRequest Err = SimpleError(-32600, "Invalid Request")
	// The method does not exist / is not available.
	ErrMethodNotFound Err = SimpleError(-32601, "Method not found")
	// Invalid method parameter(s).
	ErrInvalidParams Err = SimpleError(-32602, "Invalid params")
	// Internal JSON-RPC error.
	ErrInternalError Err = SimpleError(-32603, "Internal error")
)

// error with more detail
var (
	// ErrParse
	ErrParse_MissingField_jsonrpc = Error(-32700, "Parse error", "missing field 'jsonrpc'")
	ErrParse_InvalidField_Method  = Error(-32700, "Parse error", "invalid field 'method'")

	// ErrInvalidRequest
	ErrParse_Method_ReservedFunc = Error(-32600, "Parse error", "'method' with prefix 'rpc.' is reserved")

	// ErrInvalidParams
	ErrMethod_ParamsNumNotMatch = Error(-32602, "Invalid params", "num of params is not match")
)
