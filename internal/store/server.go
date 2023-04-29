package store

import (
	"context"
	"heisenberg/internal/pb"
)

type StoreServer struct {
	id      string
	shardId string
	master  bool
	store   store
	shard   shard
}

func (s *StoreServer) Get(ctx context.Context, req *pb.Key) (*pb.Value, error) {
	key := req.Key
	collection := req.Collection
	var res *pb.Value
	if s.master {
		shard, err := s.shard.getShard(key)
		if err != nil {
			return nil, err
		}
		client, err := shard.choose()
		if err != nil {
			return nil, err
		}
		res = client.Get(key, collection)
	} else {
		val := s.store.get(key, collection)
		res = pb.Value{
			val.Idx,
			val.Vec,
			val.Meta,
		}
	}
	return res, nil
}

func (s *StoreServer) Put() {

}

func (s *StoreServer) Close() {

}
