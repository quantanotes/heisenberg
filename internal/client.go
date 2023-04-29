package internal

import (
	"context"
	"net"

	"heisenberg/internal/pb"

	"storj.io/drpc/drpcconn"
)

type Client struct {
	addr   string
	Conn   *drpcconn.Conn
	Client pb.DRPCServiceClient
}

// Base constructor for clients
func NewClient(ctx context.Context, addr string, service Service) (*Client, error) {
	conn, err := connect(addr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		addr,
		conn,
		pb.NewDRPCServiceClient(conn),
	}

	err = client.validateClient(ctx, service)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Close() {
	c.Conn.Close()
}

func connect(addr string) (*drpcconn.Conn, error) {
	tr, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, ConnectionError(addr, err)
	}
	conn := drpcconn.New(tr)
	return conn, nil
}

func (c *Client) validateClient(ctx context.Context, expected Service) error {
	pong, err := c.Client.Ping(ctx, nil)
	if err != nil {
		return ConnectionError(c.addr)
	}
	if pong.Service != uint32(expected) {
		return IncorrectServiceError(expected, Service(pong.Service))
	}
	return nil
}
