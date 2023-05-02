package internal

import (
	"context"
	"heisenberg/internal/pb"
	"net"

	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
)

func NewServer(ctx context.Context, addr string, init pb.DRPCServiceServer) (*net.Listener, *drpcserver.Server, error) {
	m := drpcmux.New()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	err = pb.DRPCRegisterService(m, init)
	if err != nil {
		return nil, nil, err
	}
	return &lis, drpcserver.New(m), nil
}
