package store

import (
	"context"
	"fmt"
	"heisenberg/internal"
	"heisenberg/internal/pb"
)

type StoreServer struct {
	id      string
	shardId string
	master  bool
	store   store
	shard   *shard
	replica *replica
}

func (s *StoreServer) Close() {

}

func (s *StoreServer) ConnectNode(ctx context.Context, req string) error {
	c, err := NewStoreClient(ctx, req)
	if err != nil {
		return fmt.Errorf("ConnectNode error %v", err)
	}

	pong, err := c.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ConnectNode error %v", err)
	}

	switch pong.Service {
	case uint32(internal.StoreService):
		// If master add as shard, if shard add as replica
		if s.master {

		} else {

		}
		return nil
	default:
		return internal.InvalidServiceError()
	}
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

		// Select random replica of given shard
		client, err := shard.choose()
		if err != nil {
			return nil, err
		}

		res = client.Get(key, collection)
	} else {
		var err error
		res, err = s.store.get(key, collection)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (s *StoreServer) Put() {

}
