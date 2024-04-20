package model

type Method uint32

const (
	MethodPut Method = iota
	MethodDelete
)

type Job int

const (
	JobDB Job = iota
	JobMaster
	JobStore
)

type PutRequest struct {
	Key, Value []byte
}

type StorePutResponse struct {
	Key          []byte
	Page, Offset uint32
}

type CommitRequest struct {
	Path string
}

type CommitResponse struct{}
