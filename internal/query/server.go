package query

import (
	"context"
	"heisenberg/internal/index"
	"heisenberg/internal/store"
)

type QueryServer struct {
	indexClient *index.IndexClient
	storeClient *store.StoreClient
}

func NewQueryServer(ctx context.Context, indexAddr string, storeAddr string) *QueryServer {
	indexClient, err := index.NewIndexClient(ctx, indexAddr)
	if err != nil {
		return nil
	}

	storeClient, err := store.NewStoreClient(ctx, storeAddr)
	if err != nil {
		return nil
	}

	return &QueryServer{
		indexClient,
		storeClient,
	}
}
