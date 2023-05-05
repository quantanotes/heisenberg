package store

import (
	"context"
	"heisenberg/internal"
)

type StoreClient = internal.Client

func NewStoreClient(ctx context.Context, addr string) (*StoreClient, error) {
	c, err := internal.NewClient(ctx, addr, internal.StoreService)
	if err != nil {
		return nil, err
	}
	return (*StoreClient)(c), nil
}
