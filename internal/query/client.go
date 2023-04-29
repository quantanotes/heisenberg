package query

import (
	"context"
	"heisenberg/internal"
)

type QueryClient = internal.Client

func NewQueryClient(ctx context.Context, addr string) (*QueryClient, error) {
	c, err := internal.NewClient(ctx, addr, internal.QueryService)
	if err != nil {
		return nil, err
	}
	return c, err
}
