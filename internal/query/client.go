package internal

import (
	"context"
	"heisenberg/internal"
	"storj.io/drpc/drpcconn"
)

type QueryClient struct {
	client *pb.BaseClient
}

func NewQueryClient(addr string) (*QueryClient, error) {
	qc := &QueryClient{}
	qc.base, err = pb.NewBaseClient(ctx, addr, internal.QueryService, NewDRPCQueryClient)
	if err != nil {
		return nil, err
	}
	return qc, nil
}

func (qc *QueryClient) Close() {
	qc.base.Close()
}

func (qc *QueryClient) Ping(ctx context.Context) internal.Service {
	return qc.base.Ping(ctx)
}

