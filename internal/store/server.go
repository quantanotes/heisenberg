package store

import (
	"context"
	"fmt"
	"heisenberg/internal"
	"heisenberg/internal/pb"
	"heisenberg/log"
)

type StoreServer struct {
	id      string
	shardId string
	master  bool
	store   store
	shard   *shard
}

func (s *StoreServer) Close() {

}

// Handshake between other compute nodes
func (s *StoreServer) ConnectNode(ctx context.Context, req string) error {
	c, err := NewStoreClient(ctx, req)
	if err != nil {
		err := fmt.Errorf("@ConnectNode, %v", err)
		log.Error(err.Error(), nil)
		return err
	}

	pong, err := c.Ping(ctx)
	if err != nil {
		err := fmt.Errorf("@ConnectNode, %v", err)
		log.Error(err.Error(), nil)
		return err
	}

	switch pong.Service {
	case uint32(internal.StoreService):
		id := pong.Id
		shard := *pong.Shard
		if s.master {
			s.shard.addReplica(c, id, shard)
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
