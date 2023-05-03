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

func NewStoreServer(path string, id string) (*StoreServer, error) {
	s := &StoreServer{}
	store, err := loadStore(path)
	if err != nil {
		log.LogErrNilReturn[*StoreServer]("NewStoreServer", err, s.identity())
	}
	s.store = store
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

func (s *StoreServer) Ping(ctx context.Context, in *pb.Empty) (*pb.Pong, error) {
	shard := "a"
	pong := &pb.Pong{
		Id:      s.id,
		Master:  false,
		Service: uint32(internal.StoreService),
		Shard:   &shard,
	}
	return pong, nil
}

func (s *StoreServer) Connect(ctx context.Context, in *pb.Connection) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) AddShard(ctx context.Context, in *pb.Shard) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) CreateCollection(ctx context.Context, in *pb.Collection) (*pb.Empty, error) {
	collection := in.Collection
	log.Info(fmt.Sprintf("creating collection with name %s", string(collection)), s.identity())
	err := s.store.createCollection(collection)
	if err != nil {
		return log.LogErrNilReturn[pb.Empty]("CreateCollection", err, s.identity())
	}
	return nil, nil
}

func (s *StoreServer) Get(ctx context.Context, in *pb.Key) (*pb.Pair, error) {
	key := in.Key
	collection := in.Collection
	log.Info(fmt.Sprintf("getting value at key %s at collection %s", string(key), string(collection)), s.identity())
	value, err := s.store.get(key, collection)
	if err != nil {
		return log.LogErrNilReturn[pb.Pair]("Get", err, s.identity())
	}
	return &pb.Pair{Key: key, Value: value}, nil
}

func (s *StoreServer) Put(ctx context.Context, in *pb.Item) (*pb.Empty, error) {
	key := in.Key
	value := in.Value
	collection := in.Collection
	log.Info(fmt.Sprintf("putting %s at key %s at collection %s", string(value), string(key), string(collection)), s.identity())
	err := s.store.put(key, value, collection)
	if err != nil {
		return log.LogErrNilReturn[pb.Empty]("Put", err, s.identity())
	}
	return nil, nil
}

func (s *StoreServer) Delete(ctx context.Context, in *pb.Key) (*pb.Empty, error) {
	return nil, nil
}

func (s *StoreServer) identity() log.M {
	return log.M{"host": (*s.lis).Addr().String()}
}
