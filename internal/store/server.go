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

type StoreServer struct {
	server *drpcserver.Server
	lis    *net.Listener
	id     string
	store  *store
}

func NewStoreServer(id string) (*StoreServer, error) {
	s := &StoreServer{}
	s.id = id
	return s, nil
}

func (s *StoreServer) Run(ctx context.Context, addr string) error {
	lis, server, err := internal.NewServer(ctx, addr, s)
	if err != nil {
		log.LogErrReturn("Run", err)

	}
	s.server = server
	s.lis = lis
	log.Info(fmt.Sprintf("Starting store server %s", addr), nil)
	return s.server.Serve(ctx, *s.lis)
}

func (s *StoreServer) Close() {
	log.Info(fmt.Sprintf("Closing store server %s", (*s.lis).Addr().String()), nil)
	(*s.lis).Close()
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
