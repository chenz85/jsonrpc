package object

import "encoding/json"

type JsonObject map[string]interface{}

func (obj JsonObject) ToJson() []byte {
	if data, err := json.Marshal(obj); err != nil {
		return nil
	} else {
		return data
	}
}

type json_object interface {
	JsonObject() JsonObject
}

func JsonObjectArrayToJson(arr []JsonObject) []byte {
	if data, err := json.Marshal(arr); err != nil {
		return nil
	} else {
		return data
	}
}
