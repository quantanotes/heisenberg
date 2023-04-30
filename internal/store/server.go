package store

import (
	"context"
	"heisenberg/internal/pb"
	"heisenberg/log"
)

type StoreServer struct {
	id    string
	store store
}

func (s *StoreServer) Close() {

}

// Handshake between other compute nodes
func (s *StoreServer) ConnectNode(ctx context.Context, req string) error {
	return nil
}

func (s *StoreServer) Get(ctx context.Context, req *pb.Key) (*pb.Value, error) {
	key := req.Key
	collection := req.Collection
	res, err := s.store.get(key, collection)
	if err != nil {
		log.LogErrNilReturn[pb.Value]("Get", err)
	}
	return res, nil
}

func (s *StoreServer) Put(ctx context.Context, req *pb.Pair) error {
	return nil
}
