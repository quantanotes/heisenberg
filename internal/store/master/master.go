package master

import (
	"context"
	"heisenberg/internal"
	"heisenberg/internal/pb"
	"heisenberg/internal/store"
	"heisenberg/log"
)

type Master struct {
	shard shard
}

func NewMaster(ctx context.Context, addr string) (*Master, error) {
	m := &Master{}
	go internal.Serve(ctx, addr, m)
	return m, nil
}

func (m *Master) ConnectNode(ctx context.Context, req string) error {
	c, err := store.NewStoreClient(ctx, req)
	if err != nil {
		return log.LogErrReturn("ConnectNode", err)
	}

	pong, err := c.Ping(ctx)
	if err != nil {
		return log.LogErrReturn("ConnectNode", err)
	}

	switch pong.Service {
	case uint32(internal.StoreService):
		id := pong.Id
		shard := *pong.Shard
		m.shard.addReplica(c, id, shard)
	default:
		return log.LogErrReturn("ConnectNode", err)
	}

	return nil
}

func (m *Master) Ping(ctx context.Context, req *pb.Empty) (*pb.Pong, error) {
	pong := &pb.Pong{
		Id:      "0",
		Master:  true,
		Service: uint32(internal.StoreService),
		Shard:   nil,
	}
	return pong, nil
}

func (m *Master) Get(ctx context.Context, req *pb.Key) (*pb.Pair, error) {
	key := req.Key
	collection := req.Collection
	shard, err := m.shard.getShard(key)
	if err != nil {
		return log.LogErrNilReturn[pb.Pair]("Get", err)
	}
	// Select random replica of given shard
	client, err := shard.choose()
	if err != nil {
		return log.LogErrNilReturn[pb.Pair]("Get", err)
	}
	res := client.Get(key, collection)
	return res, nil
}

func (m *Master) Put(ctx context.Context, req *pb.Pair) (*pb.Empty, error) {
	return nil, nil
}

func (m *Master) Delete(ctx context.Context, req *pb.Key) (*pb.Empty, error) {
	return nil, nil
}
