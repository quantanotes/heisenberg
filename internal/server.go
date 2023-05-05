package internal

import (
	"context"
	"heisenberg/internal/pb"
	"heisenberg/log"
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

// Identify server for logging
func Identity(lis *net.Listener) log.M {
	return log.M{"host": (*lis).Addr().String()}
}
