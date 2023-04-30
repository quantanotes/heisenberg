package internal

import (
	"context"
	"heisenberg/internal/pb"
	"heisenberg/log"
	"net"

	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
)

func Serve(ctx context.Context, addr string, init pb.DRPCServiceServer) error {
	m := drpcmux.New()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.LogErrReturn("Serve", err)
	}
	defer lis.Close()
	err = pb.DRPCRegisterService(m, init)
	if err != nil {
		log.LogErrReturn("Serve", err)
	}
	return drpcserver.New(m).Serve(ctx, lis)
}
