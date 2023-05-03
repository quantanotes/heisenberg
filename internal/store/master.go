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

const shardKey = "__HEISENBERG_SHARD_CONFIG"
const configKey = "__HEISENBERG_CONFIG"

type StoreMasterServer struct {
	server *drpcserver.Server
	lis    *net.Listener
	shard  *shard
	store  *store
}

func NewStoreMasterServer(path string) (*StoreMasterServer, error) {
	m := &StoreMasterServer{}
	store, err := loadStore(path)
	if err != nil {
		return log.LogErrNilReturn[StoreMasterServer]("LoadStoreMasterServer", err, m.identity())
	}
	m.store = store
	m.initConfig()
	if err != nil {
		return log.LogErrNilReturn[StoreMasterServer]("LoadStoreMasterError", err, m.identity())
	}
	m.shard = newShard()
	err = m.saveConfig()
	if err != nil {
		return log.LogErrNilReturn[StoreMasterServer]("LoadStoreMasterError", err, m.identity())
	}
	return m, nil
}

func LoadStoreMasterServer(path string) (*StoreMasterServer, error) {
	m := &StoreMasterServer{}
	store, err := loadStore(path)
	if err != nil {
		log.LogErrNilReturn[StoreMasterServer]("LoadStoreMasterServer", err, m.identity())
	}
	m.store = store
	err = m.loadConfig()
	if err != nil {
		log.LogErrNilReturn[StoreMasterServer]("LoadStoreMasterServer", err, m.identity())
	}
	return m, nil
}

func (m *StoreMasterServer) initConfig() error {
	return m.store.createCollection([]byte(configKey))
}

func (m *StoreMasterServer) loadConfig() error {
	raw, err := m.store.get([]byte(shardKey), []byte(configKey))
	if err != nil {
		return err
	}
	shard, err := shardFromBytes(raw)
	if err != nil {
		return err
	}
	m.shard = shard
	return nil
}

func (m *StoreMasterServer) saveConfig() error {
	val, err := m.shard.toByte()
	if err != nil {
		return err
	}
	return m.store.put([]byte(shardKey), val, []byte(configKey))
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
	m.saveConfig()
	m.store.close()
	(*m.lis).Close()
}

func (m *StoreMasterServer) Ping(ctx context.Context, in *pb.Empty) (*pb.Pong, error) {
	log.Info("recieved ping", m.identity())
	pong := &pb.Pong{
		Id:      "0",
		Master:  true,
		Service: uint32(internal.StoreService),
		Shard:   nil,
	}
	return pong, nil
}

func (m *StoreMasterServer) AddShard(ctx context.Context, in *pb.Shard) (*pb.Empty, error) {
	old := *m.shard // to revert if database transaction fails
	id := in.Shard
	log.Info(fmt.Sprintf("adding shard with id %s", id), m.identity())
	err := m.shard.addShard(id)
	if err != nil {
		return log.LogErrNilReturn[pb.Empty]("AddShard", err, m.identity())
	}
	err = m.saveConfig()
	if err != nil {
		m.shard = &old
		return log.LogErrNilReturn[pb.Empty]("AddShard", err, m.identity())
	}
	return nil, nil
}

func (m *StoreMasterServer) Connect(ctx context.Context, in *pb.Connection) (*pb.Empty, error) {
	addr := in.Address
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
		return log.LogErrNilReturn[pb.Empty]("Connect", internal.InvalidServiceError(), m.identity())
	}
	return nil, nil
}

func (m *StoreMasterServer) Get(ctx context.Context, in *pb.Key) (*pb.Pair, error) {
	key := in.Key
	collection := in.Collection
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

func (m *StoreMasterServer) Put(ctx context.Context, in *pb.Item) (*pb.Empty, error) {
	key := in.Key
	value := in.Value
	collection := in.Collection
	shard, err := m.shard.getShard(key)
	if err != nil {
		log.LogErrNilReturn[pb.Empty]("Put", err, m.identity())
	}
	// put to all replicas in shard, roll back all transactions if failiure on one node
	for _, replica := range shard.clients {
		err := replica.Put(ctx, key, value, collection)
		if err != nil {
			_ = 0
			// TODO : ROLL BACK TRANSACTIONS
		}
	}
	return nil, nil
}

func (m *StoreMasterServer) Delete(ctx context.Context, in *pb.Key) (*pb.Empty, error) {
	return nil, nil
}

func (m *StoreMasterServer) identity() log.M {
	host := "nil"
	if m.lis != nil {
		host = (*m.lis).Addr().String()
	}
	return log.M{"host": host}
}
