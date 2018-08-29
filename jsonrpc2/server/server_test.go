package server_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/czsilence/jsonrpc/jsonrpc2/server"
)

func TestHandleRequest(t *testing.T) {
	// examples from: https://www.jsonrpc.org/specification
	var values = []struct {
		request  string
		response string
	}{
		{
			`{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			`{"jsonrpc": "2.0", "result": 19, "id": 1}`,
		},
		{
			`{"jsonrpc": "2.0", "method": "subtract", "params": [23, 42], "id": 2}`,
			`{"jsonrpc": "2.0", "result": -19, "id": 2}`,
		},
		// {
		// 	`{"jsonrpc": "2.0", "method": "subtract", "params": {"subtrahend": 23, "minuend": 42}, "id": 3}`,
		// 	`{"jsonrpc": "2.0", "result": 19, "id": 3}`,
		// },
		// {
		// 	`{"jsonrpc": "2.0", "method": "subtract", "params": {"minuend": 42, "subtrahend": 23}, "id": 4}`,
		// 	`{"jsonrpc": "2.0", "result": 19, "id": 4}`,
		// },
		{
			`{"jsonrpc": "2.0", "method": "update", "params": [1,2,3,4,5]}`,
			``,
		},
		{
			`{"jsonrpc": "2.0", "method": "foobar"}`,
			``,
		},
		{
			`{"jsonrpc": "2.0", "method": "foobar1", "id": "1"}`,
			`{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
		},
		{
			`{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]`,
			`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
		},
		{
			`{"jsonrpc": "2.0", "method": 1, "params": "bar"}`,
			`{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		},
		{
			`[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
				{"jsonrpc": "2.0", "method"
			]`,
			`{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
		},
		{
			`[]`,
			`{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		},
		{
			`[1]`,
			`[{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}]`,
		},
		{
			`[1,2,3]`,
			`[
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}
			]`,
		},
		{
			`[
				{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]},
				{"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": "2"},
				{"foo": "boo"},
				{"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": "5"},
				{"jsonrpc": "2.0", "method": "get_data", "id": "9"} 
			]`,
			`[
				{"jsonrpc": "2.0", "result": 7, "id": "1"},
				{"jsonrpc": "2.0", "result": 19, "id": "2"},
				{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null},
				{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "5"},
				{"jsonrpc": "2.0", "result": ["hello", 5], "id": "9"}
			]`,
		},
		{
			`[
				{"jsonrpc": "2.0", "method": "notify_sum", "params": [1,2,4]},
				{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]}
			]`,
			``,
		},
	}

	server.RegisterMethod("subtract", func(a, b float64) float64 {
		t.Log("subtract", a, b)
		return a - b
	})

	server.RegisterMethod("update", func(a, b, c, d, e float64) {
		t.Log("update", a, b, c, d, e)
	})

	server.RegisterMethod("foobar", func() {
		t.Log("foobar")
	})

	server.RegisterMethod("sum", func(a, b, c float64) float64 {
		t.Log("sum", a, b, c)
		return a + b + c
	})

	server.RegisterMethod("notify_hello", func(v float64) {
		t.Log("notify_hello", v)
	})

	server.RegisterMethod("get_data", func() interface{} {
		t.Log("get_data")
		return []interface{}{"hello", 5}
	})

	server.RegisterMethod("notify_sum", func(a, b, c float64) {
		t.Log("notify_sum", a, b, c)
	})

	for i, value := range values {
		resp := server.HandleRequest([]byte(value.request))
		// t.Logf("request (%d): %s", i, value.request)
		if len(resp) == 0 && value.response == "" {
			continue
		} else if len(resp) == 0 && value.response != "" || len(resp) != 0 && value.response == "" {
			t.Errorf("response not match (#%d) (type): Got: %s, Expect: %s", i, string(resp), value.response)
		} else if resp[0] != value.response[0] {
			t.Errorf("response not match (#%d) (type): Got: %s, Expect: %s", i, string(resp), value.response)
		} else {
			var a, b interface{}
			var batch bool = resp[0] == '['
			if batch {
				a = make([]interface{}, 0, 1)
				b = make([]interface{}, 0, 1)
			} else {
				a = make(map[string]interface{})
				b = make(map[string]interface{})
			}
			if ea, eb := json.Unmarshal(resp, &a), json.Unmarshal([]byte(value.response), &b); ea != nil || eb != nil {
				t.Error(ea)
				t.Error(eb)
				t.Errorf("response not match (#%d) (parse): Got: %s, Expect: %s", i, string(resp), value.response)
			} else {
				ignore_data(a)
				if !obj_cmpr(a, b) {
					t.Errorf("response not match (#%d) (cmpr): Got: %s, Expect: %s", i, string(resp), value.response)
				}
			}
		}

	}
}

func ignore_data(obj interface{}) {
	switch obj := obj.(type) {
	case []interface{}:
		for _, o := range obj {
			ignore_data(o)
		}
	case map[string]interface{}:
		if fe_val, ex := obj["error"]; !ex {
		} else if fe, ok := fe_val.(map[string]interface{}); ok {
			delete(fe, "data")
		}
	}
}

func obj_cmpr(a, b interface{}) bool {
	ra, rb := reflect.ValueOf(a), reflect.ValueOf(b)
	if ra.IsValid() != rb.IsValid() {
		fmt.Println("diff1")
		return false
	} else if ra.Kind() != rb.Kind() {
		fmt.Println("diff2")
		return false
	}

	switch ra.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		return ra.Int() == rb.Int()
	case reflect.Uint, reflect.Uint32, reflect.Uint64, reflect.Uint16, reflect.Uint8:
		return ra.Uint() == rb.Uint()
	case reflect.Float32, reflect.Float64:
		return ra.Float() == rb.Float()
	case reflect.Bool:
		return ra.Bool() == rb.Bool()
	case reflect.String:
		return ra.String() == rb.String()
	case reflect.Array, reflect.Slice:
		if ra.Len() != rb.Len() {
			fmt.Println("diff3")
			return false
		} else {
			for i := 0; i < ra.Len(); i++ {
				if !obj_cmpr(ra.Index(i).Interface(), rb.Index(i).Interface()) {
					return false
				}
			}
		}
	case reflect.Map:
		keys := ra.MapKeys()
		if len(keys) != len(rb.MapKeys()) {
			fmt.Println("diff4")
			return false
		}

		for _, k := range keys {
			fa, fb := ra.MapIndex(k), rb.MapIndex(k)
			if fa.IsValid() != fb.IsValid() {
				fmt.Println("diff5")
				return false
			} else if !obj_cmpr(fa.Interface(), fb.Interface()) {
				return false
			}
		}
	}
	return true
}
