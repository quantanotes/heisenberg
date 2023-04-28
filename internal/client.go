package internal

import (
	Errors "heisenberg/errors"
	"net"

	"storj.io/drpc/drpcconn"
)

func NewConnection(addr string) (*drpcconn.Conn, error) {
	tr, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, Errors.ConnectionError(addr, err)
	}

	conn := drpcconn.New(tr)
	return conn, nil
}
