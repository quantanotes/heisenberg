package index

import (
	"context"
	"heisenberg/internal"
)

type IndexClient = internal.Client

func NewIndexClient(ctx context.Context, addr string) (*IndexClient, error) {
	c, err := internal.NewClient(ctx, addr, internal.IndexService)
	if err != nil {
		return nil, err
	}
	return c, err
}
