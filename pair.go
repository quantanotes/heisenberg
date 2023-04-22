package main

import (
	"bytes"
	"encoding/gob"
)

type value struct {
	Vec  []float32 `json:"vec"`
	Meta any       `json:"meta"`
}

type pair struct {
	K string `json:"k"`
	V value  `json:"v"`
}

func (v *value) toBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func valueFromBytes(data []byte) (*value, error) {
	v := value{}
	dec := gob.NewDecoder(bytes.NewReader(data))

	err := dec.Decode(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
