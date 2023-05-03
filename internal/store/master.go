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
	m.shard = *newShard()
	return m, nil
}

func (m *StoreMasterServer) Run(ctx context.Context, addr string) error {
	lis, server, err := internal.NewServer(ctx, addr, m)
	if err != nil {
		log.LogErrReturn("RunStoreMasterServer", err, m.identity())
	}
	m.server = server
	m.lis = lis
	log.Info("starting master store server", m.identity())
	return m.server.Serve(ctx, *m.lis)
}

func (m *StoreMasterServer) Close() {
	log.Info("closing master store server", m.identity())
	(*m.lis).Close()
}

func (m *StoreMasterServer) Ping(ctx context.Context, msg *pb.Empty) (*pb.Pong, error) {
	log.Info("recieved ping", m.identity())
	pong := &pb.Pong{
		Id:      "0",
		Master:  true,
		Service: uint32(internal.StoreService),
		Shard:   nil,
	}
	return pong, nil
}

func (m *StoreMasterServer) AddShard(ctx context.Context, msg *pb.Shard) (*pb.Empty, error) {
	id := msg.Shard
	log.Info(fmt.Sprintf("adding shard with id %s", id), m.identity())
	err := m.shard.addShard(id)
	if err != nil {
		return log.LogErrNilReturn[pb.Empty]("AddShard", err, m.identity())
	}
	return nil, nil
}

func (m *StoreMasterServer) Connect(ctx context.Context, msg *pb.Connection) (*pb.Empty, error) {
	addr := msg.Address
	log.Info(fmt.Sprintf("connecting client to server at %s", addr), m.identity())
	c, err := NewStoreClient(ctx, addr)
	if err != nil {
		return log.LogErrNilReturn[pb.Empty]("Connect", err, m.identity())
	}

	pong, err := c.Ping(ctx)
	if err != nil {
		return log.LogErrNilReturn[pb.Empty]("Connect", err, m.identity())
	}

	switch pong.Service {
	case uint32(internal.StoreService):
		id := pong.Id
		shard := *pong.Shard
		err := m.shard.addReplica(c, id, shard)
		if err != nil {
			return log.LogErrNilReturn[pb.Empty]("Connect", err, m.identity())
		}
	default:
		return log.LogErrNilReturn[pb.Empty]("Connect", err, m.identity())
	}

	return nil, nil
}

func (m *StoreMasterServer) Get(ctx context.Context, msg *pb.Key) (*pb.Pair, error) {
	key := msg.Key
	collection := msg.Collection
	shard, err := m.shard.getShard(key)
	if err != nil {
		return log.LogErrNilReturn[pb.Pair]("Get", err, m.identity())
	}
	// Select random replica of given shard
	client, err := shard.choose()
	if err != nil {
		return log.LogErrNilReturn[pb.Pair]("Get", err, m.identity())
	}
	res := client.Get(key, collection)
	return res, nil
}

func (m *StoreMasterServer) Put(ctx context.Context, msg *pb.Pair) (*pb.Empty, error) {
	return nil, nil
}

func (m *StoreMasterServer) Delete(ctx context.Context, msg *pb.Key) (*pb.Empty, error) {
	return nil, nil
}

func (m *StoreMasterServer) identity() log.M {
	return log.M{"host": (*m.lis).Addr().String()}
}
