package index

import (
	"context"
	"fmt"
	"heisenberg/internal"
	"heisenberg/internal/pb"
	"heisenberg/log"
	"net"

	"storj.io/drpc/drpcserver"
)

type MasterIndexServer struct {
	server *drpcserver.Server
	lis    *net.Listener
	shard  *shard
}

func NewMasterIndexServer(path string, id string) (*MasterIndexServer, error) {
	m := &MasterIndexServer{}
	m.shard = newShard()
	return m, nil
}

func (m *MasterIndexServer) Run(ctx context.Context, addr string) error {
	lis, server, err := internal.NewServer(ctx, addr, m)
	if err != nil {
		log.LogErrReturn("Run", err, m.identity())
	}
	m.server = server
	m.lis = lis
	log.Info("starting master index server", m.identity())
	return m.server.Serve(ctx, *m.lis)
}

func (m *MasterIndexServer) Close() {
	log.Info("closing master index server", m.identity())
}

func (m *MasterIndexServer) Ping(ctx context.Context, in *pb.Empty) (*pb.Pong, error) {
	log.Info("recieved ping", m.identity())
	pong := &pb.Pong{
		Id:      "0",
		Master:  true,
		Service: uint32(internal.IndexService),
		Shard:   nil,
	}
	return pong, nil
}

func (m *MasterIndexServer) Connect(ctx context.Context, in *pb.Connection) (*pb.Empty, error) {
	return nil, nil
}

func (m *MasterIndexServer) AddShard(ctx context.Context, in *pb.Shard) (*pb.Empty, error) {
	return nil, nil
}

func (m *MasterIndexServer) RemoveShard(ctx context.Context, in *pb.Shard) (*pb.Empty, error) {
	return nil, nil
}

func (m *MasterIndexServer) CreateCollection(ctx context.Context, in *pb.Collection) (*pb.Empty, error) {
	collection := in.Collection
	log.Info(fmt.Sprintf("creating collection %s", string(collection)), m.identity())
	for _, shard := range *m.shard.getReplicas() {
		for _, replica := range shard.clients {
			err := replica.CreateCollection(ctx, collection)
			if err != nil {
				_ = 0
				// TODO : ROLL BACK TRANSACTIONS
			}
		}
	}
	return nil, nil
}

func (m *MasterIndexServer) DeleteCollection(ctx context.Context, in *pb.Collection) (*pb.Empty, error) {
	collection := in.Collection
	log.Info(fmt.Sprintf("creating collection %s", string(collection)), m.identity())
	for _, shard := range *m.shard.getReplicas() {
		for _, replica := range shard.clients {
			err := replica.DeleteCollection(ctx, collection)
			if err != nil {
				_ = 0
				// TODO : ROLL BACK TRANSACTIONS
			}
		}
	}
	return nil, nil
}

func (m *MasterIndexServer) Get(ctx context.Context, in *pb.Key) (*pb.Pair, error) {
	return nil, nil
}

func (m *MasterIndexServer) Put(ctx context.Context, in *pb.Item) (*pb.Empty, error) {
	return nil, nil
}

func (m *MasterIndexServer) Delete(ctx context.Context, in *pb.Key) (*pb.Empty, error) {
	return nil, nil
}

func (m *MasterIndexServer) identity() log.M {
	return internal.Identity(m.lis)
}
