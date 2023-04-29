package store

import (
	"context"
)

type StoreServer struct {
	id      string
	shardId string
	master  bool
	store   store
	shard   shard
}

func (s *StoreServer) Get(ctx context.Context, req any) {
	key := req.key
	collection := req.collection
	if !s.master {
		shard, err := s.shard.getShard(key)
	}
	s.store.get(key, collection)
}

func (s *StoreServer) Put() {

}
