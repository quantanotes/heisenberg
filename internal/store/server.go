package store

import (
	"context"
	"heisenberg/internal"
	"heisenberg/internal/pb"
	"heisenberg/log"
)

type StoreServer struct {
	id    string
	store store
}

func RunStoreServer(ctx context.Context, addr string) {
	m := &StoreServer{}
	go internal.Serve(ctx, addr, m)
}

func (s *StoreServer) Ping(ctx context.Context, req *pb.Empty) (*pb.Pong, error) {
	pong := &pb.Pong{
		Id:      s.id,
		Master:  false,
		Service: uint32(internal.StoreService),
		Shard:   nil,
	}
	return pong, nil
}

func (s *StoreServer) ConnectNode(ctx context.Context, req string) error {
	return nil
}

func (s *StoreServer) Get(ctx context.Context, req *pb.Key) (*pb.Pair, error) {
	key := req.Key
	collection := req.Collection
	_, err := s.store.get(key, collection)
	if err != nil {
		log.LogErrNilReturn[pb.Pair]("Get", err)
	}
	return nil, nil // TODO
}

func (s *StoreServer) Put(ctx context.Context, req *pb.Pair) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) Delete(ctx context.Context, req *pb.Key) (*pb.Empty, error) {
	return nil, nil
}
