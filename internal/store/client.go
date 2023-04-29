package store

import (
	"context"
	"heisenberg/internal"
	"heisenberg/internal/pb"

	"storj.io/drpc/drpcconn"
)

type StoreClient struct {
	conn   *drpcconn.Conn
	client *pb.DRPCServiceClient
}

func NewStoreClient(ctx context.Context, addr string) (*StoreClient, error) {
	c, err := internal.NewClient(ctx, addr, internal.StoreService)
	if err != nil {
		return nil, err
	}
	return &StoreClient{c.Conn, &c.Client}, err
}

func (c *StoreClient) Get() *pb.Value {

}
