package store

import (
	"context"
	"heisenberg/internal"

	"storj.io/drpc/drpcconn"
)

type StoreClient struct {
	base  *pb.BaseClient
}

func NewStoreClient(ctx context.Context, addr string) (*StoreClient, error) {
	sc := &StoreClient{}
	sc.base, err = pb.NewBaseClient(ctx, addr, internal.StoreService, NewDRPCStoreClient)
	return sc, nil
}

func (sc *StoreClient) Close() {
	sc.base.Close()
}

func (sc *StoreClient) Ping(ctx context.Context) internal.Service {
	return sc.base.Ping(ctx)
}
