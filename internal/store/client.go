package store

import (
	"context"
	"heisenberg/internal"
	"heisenberg/internal/pb"

	"storj.io/drpc/drpcconn"
)

type StoreClient struct {
	addr   string
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

func (c *StoreClient) Get(key []byte, collection []byte) *pb.Pair {
	return nil
}

func (c *StoreClient) Ping(ctx context.Context) (*pb.Pong, error) {
	return c.client.Ping(ctx, nil)
}
