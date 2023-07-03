package common

import (
	"context"
	"net"

	"storj.io/drpc"
	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
)

func Serve[T any](ctx context.Context, addr string, sr func(drpc.Mux, T) error, service T) error {
	mux := drpcmux.New()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()
	if err := sr(mux, service); err != nil {
		return err
	}
	return drpcserver.New(mux).Serve(ctx, lis)
}
