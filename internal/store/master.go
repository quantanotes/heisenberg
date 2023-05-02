package store

import (
	"context"
	"fmt"
	"heisenberg/internal"
	"heisenberg/internal/pb"
	"heisenberg/log"
	"net"

	"storj.io/drpc/drpcserver"
)

type StoreMasterServer struct {
	server *drpcserver.Server
	lis    *net.Listener
	shard  shard
}

func NewStoreMasterServer() (*StoreMasterServer, error) {
	m := &StoreMasterServer{}
	return m, nil
}

func (m *StoreMasterServer) Run(ctx context.Context, addr string) error {
	lis, server, err := internal.NewServer(ctx, addr, m)
	if err != nil {
		log.LogErrReturn("RunStoreMasterServer", err)

	}
	m.server = server
	m.lis = lis
	log.Info(fmt.Sprintf("Starting master store server %s", addr), nil)
	return m.server.Serve(ctx, *m.lis)
}

func (m *StoreMasterServer) Close() {
	log.Info(fmt.Sprintf("Closing master store server %s", (*m.lis).Addr().String()), nil)
	(*m.lis).Close()
}

func (s *StoreMasterServer) Ping(ctx context.Context, req *pb.Empty) (*pb.Pong, error) {
	pong := &pb.Pong{
		Id:      "0",
		Master:  true,
		Service: uint32(internal.StoreService),
		Shard:   nil,
	}
	return pong, nil
}

func (s *StoreMasterServer) ConnectNode(ctx context.Context, req string) error {
	c, err := NewStoreClient(ctx, req)
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
		s.shard.addReplica(c, id, shard)
	default:
		return log.LogErrReturn("ConnectNode", err)
	}

	return nil
}

func (s *StoreMasterServer) Get(ctx context.Context, req *pb.Key) (*pb.Pair, error) {
	key := req.Key
	collection := req.Collection
	shard, err := s.shard.getShard(key)
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

func (s *StoreMasterServer) Put(ctx context.Context, req *pb.Pair) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreMasterServer) Delete(ctx context.Context, req *pb.Key) (*pb.Empty, error) {
	return nil, nil
}
