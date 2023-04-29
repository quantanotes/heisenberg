package index

import (
	"context"
	"heisenberg/internal"
	"storj.io/drpc/drpcconn"
)

type IndexClient struct {
	client *pb.BaseClient
}

func NewIndexClient(addr string) (*IndexClient, error) {
	ic := &IndexClient{}
	ic.base, err = pb.NewBaseClient(ctx, addr, internal.IndexService, NewDRPCIndexClient)
	if err != nil {
		return nil, err
	}
	return ic, nil
}

func (ic *IndexClient) Close() {
	ic.base.Close()
}

func (ic *StoreClient) Ping(ctx context.Context) internal.Service {
	return ic.base.Ping(ctx)
}

