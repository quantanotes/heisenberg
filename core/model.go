package core

import (
	"bytes"
	"encoding/gob"
)

type Entry struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Value  Value  `json:"value"`
}

type Value struct {
	Index  uint           `json:"index"`
	Vector []float32      `json:"vector"`
	Meta   map[string]any `json:"meta"`
}

func (v *Value) Serialise() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeserialiseValue(data []byte) (Value, error) {
	var v Value
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(&v); err != nil {
		return Value{}, err
	}
	return v, nil
}
