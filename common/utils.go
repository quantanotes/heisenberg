package common

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

func IntToBytes(i int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func ToBytes(val interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func FromBytes[T interface{}](data []byte, val *T) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(val)
	if err != nil {
		return err
	}
	return nil
}
