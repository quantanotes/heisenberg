package cache

import (
	"fmt"
	"heisenberg/common"
	"heisenberg/graph"
)

type bucket struct {
	name     string
	data     map[string]*common.Value
	graph    []graph.EdgeList
	iter     uint64
	freelist common.Queue[uint64]
}

func newBucket(name string) bucket {
	return bucket{
		name:     name,
		data:     make(map[string]*common.Value),
		graph:    make([]graph.EdgeList, 0),
		iter:     0,
		freelist: common.NewQueue[uint64](),
	}
}

func (b *bucket) get(key string) (*common.Value, error) {
	if v := b.data[key]; v != nil {
		return v, nil
	}
	return nil, fmt.Errorf("key %s not in bucket %s", key, b.name)
}

func (b *bucket) put(key string, vector []float32, meta common.Meta) {
	b.data[key] = &common.Value{
		Index:  b.next(),
		Vector: vector,
		Meta:   meta,
	}
}

func (b *bucket) delete(key string) error {
	if v := b.data[key]; v != nil {
		b.freelist.Push(v.Index)
		delete(b.data, key)
		return nil
	}
	return fmt.Errorf("key %s not in bucket %s", key, b.name)
}

func (b *bucket) next() uint64 {
	if index := b.freelist.Pop(); index != nil {
		return *index
	}
	b.iter++
	return b.iter
}
