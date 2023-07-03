package common

import (
	"context"
	"net"

	"storj.io/drpc/drpcconn"
)

func NewConn(ctx context.Context, addr string) (*drpcconn.Conn, error) {
	rawconn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn := drpcconn.New(rawconn)
	return conn, nil
}
