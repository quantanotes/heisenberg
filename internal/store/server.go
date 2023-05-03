package store

import (
	"context"
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
		log.LogErrReturn("Run", err, s.identity())

	}
	s.server = server
	s.lis = lis
	log.Info("starting store server", s.identity())
	return s.server.Serve(ctx, *s.lis)
}

func (s *StoreServer) Close() {
	log.Info("closing store server", s.identity())
	(*s.lis).Close()
}

func (s *StoreServer) Ping(ctx context.Context, msg *pb.Empty) (*pb.Pong, error) {
	shard := "a"

	pong := &pb.Pong{
		Id:      s.id,
		Master:  false,
		Service: uint32(internal.StoreService),
		Shard:   &shard,
	}
	return pong, nil
}

func (s *StoreServer) Connect(ctx context.Context, msg *pb.Connection) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) AddShard(ctx context.Context, msg *pb.Shard) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) Get(ctx context.Context, msg *pb.Key) (*pb.Pair, error) {
	key := msg.Key
	collection := msg.Collection
	_, err := s.store.get(key, collection)
	if err != nil {
		log.LogErrNilReturn[pb.Pair]("Get", err, s.identity())
	}
	return nil, nil // TODO
}

func (s *StoreServer) Put(ctx context.Context, msg *pb.Pair) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) Delete(ctx context.Context, msg *pb.Key) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) identity() log.M {
	return log.M{"host": (*s.lis).Addr().String()}
}
