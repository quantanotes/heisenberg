package store

import (
	"context"
	"heisenberg/internal/pb"
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
		return nil, err
	}
	return res, nil
}

func (s *StoreServer) Put(ctx context.Context, req *pb.Pair) error {
	return nil
}
