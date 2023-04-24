package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

func ToBytes(data interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)

	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FromBytes(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

func IntToBytes(val int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(val))
	return b
}
