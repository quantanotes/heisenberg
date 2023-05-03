package store

import (
	"context"
	"heisenberg/internal"
	"heisenberg/internal/pb"

	"storj.io/drpc/drpcconn"
)

type StoreClient struct {
	Addr   string
	conn   *drpcconn.Conn
	client pb.DRPCServiceClient
}

func NewStoreClient(ctx context.Context, addr string) (*StoreClient, error) {
	c, err := internal.NewClient(ctx, addr, internal.StoreService)
	if err != nil {
		return nil, err
	}
	return &StoreClient{addr, c.Conn, c.Client}, err
}

func (c *StoreClient) Close() {
	c.conn.Close()
}

func (c *StoreClient) Connect(ctx context.Context, addr string) {
	c.client.Connect(ctx, &pb.Connection{Address: addr})
}

func (c *StoreClient) AddShard(ctx context.Context, id string) {
	c.client.AddShard(ctx, &pb.Shard{Shard: id})
}

func (c *StoreClient) CreateCollection(ctx context.Context, colleciton []byte) error {
	_, err := c.client.CreateCollection(ctx, &pb.Collection{Collection: colleciton})
	return err
}

func (c *StoreClient) Get(ctx context.Context, key []byte, collection []byte) (*pb.Pair, error) {
	pair, err := c.client.Get(ctx, &pb.Key{Key: key, Collection: collection})
	return pair, err
}

func (c *StoreClient) Put(ctx context.Context, key []byte, value []byte, collection []byte) error {
	_, err := c.client.Put(ctx, &pb.Item{Key: key, Value: value, Collection: collection})
	return err
}

func (c *StoreClient) Ping(ctx context.Context) (*pb.Pong, error) {
	return c.client.Ping(ctx, nil)
}
