package cache

import (
	"heisenberg/common"
	"io"

	"github.com/hashicorp/raft"
)

type shardFSM struct {
	cache *Cache
}

type shardFSMPayload struct {
	Cmd    string
	Bucket string
	Key    string
	Vector []float32
	Meta   common.Meta
}

func newShardFSM(cache *Cache) *shardFSM {
	return &shardFSM{cache}
}

func (f *shardFSM) Apply(log *raft.Log) any {
	return nil
}

func (f *shardFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *shardFSM) Restore(snapshot io.ReadCloser) error {
	return nil
}

type globalFSM struct {
	server *Server
}
