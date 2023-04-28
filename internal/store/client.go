package store

import (
	"context"
	"heisenberg/internal"

	"storj.io/drpc/drpcconn"
)

type StoreClient struct {
	conn   *drpcconn.Conn
	client DRPCStoreClient
}

func NewStoreClient(ctx context.Context, addr string) (*StoreClient, error) {
	sc := &StoreClient{}
	var err error
	sc, err = internal.NewClient(ctx, addr, internal.StoreService)

}

func (sc *StoreClient) Close() {
	sc.conn.Close()
}

func (sc *StoreClient) Ping(ctx context.Context) internal.Service {
	pong, err := sc.client.Ping(ctx, nil)
	if err != nil {
		return internal.NoneService
	}
	return internal.Service{Code: pong.Service.Code, Name: pong.Service.Name}
}
