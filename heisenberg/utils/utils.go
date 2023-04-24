package utils

import (
	"encoding/binary"
	"encoding/json"
)

func ToJson(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func FromJson(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

func IntToBytes(val int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(val))
	return b
}
