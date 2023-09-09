package json

import (
	jsoniter "github.com/json-iterator/go"
)

var j = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return j.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return j.Unmarshal(data, v)
}
