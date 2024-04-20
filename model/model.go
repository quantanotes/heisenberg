package model

type Method uint32

const (
	MethodPut    = 1
	MethodDelete = 2
)

type CommitRequest struct {
	Path string
}

type CommitResponse struct{}
