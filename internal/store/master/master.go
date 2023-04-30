package master

import (
	"context"
	"fmt"
	"heisenberg/internal"
	"heisenberg/internal/pb"
	"heisenberg/internal/store"
	"heisenberg/log"
)

type Master struct {
	shard shard
}

func (m *Master) ConnectNode(ctx context.Context, req string) error {
	c, err := store.NewStoreClient(ctx, req)
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
		m.shard.addReplica(c, id, shard)
	default:
		return internal.InvalidServiceError()
	}

	return nil
}

func (m *Master) Get(ctx context.Context, req *pb.Key) (*pb.Value, error) {
	key := req.Key
	collection := req.Collection
	shard, err := m.shard.getShard(key)
	if err != nil {
		return nil, err
	}
	// Select random replica of given shard
	client, err := shard.choose()
	if err != nil {
		return nil, err
	}
	res := client.Get(key, collection)
	return res, nil
}
